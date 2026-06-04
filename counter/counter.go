// Package counter provides symbol ↔ counter_id conversion utilities.
//
// A counter_id is the internal instrument identifier used by the Longbridge
// backend, e.g. "ST/US/TSLA", "ETF/US/SPY", "IX/HK/HSI", "WT/HK/10005". These
// helpers convert between user-facing symbols (e.g. "TSLA.US", "700.HK",
// ".DJI.US") and counter IDs, using an embedded ETF + index + warrant
// directory to pick the right prefix.
//
// The embedded directory may lag behind newly listed instruments. Entries
// resolved remotely (see quote.QuoteContext.ResolveCounterIds) are persisted
// to a local cache file and consulted on subsequent lookups.
package counter

import (
	"bufio"
	"embed"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

//go:embed US-ETF.csv US-IX.csv US-WT.csv
var directoryFS embed.FS

// specialCounterIDs returns the embedded ETF + index + warrant directory as a
// set, loaded lazily once.
var specialCounterIDs = sync.OnceValue(func() map[string]struct{} {
	set := make(map[string]struct{})
	for _, name := range []string{"US-ETF.csv", "US-IX.csv", "US-WT.csv"} {
		data, err := directoryFS.ReadFile(name)
		if err != nil {
			continue
		}
		sc := bufio.NewScanner(strings.NewReader(string(data)))
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line != "" {
				set[line] = struct{}{}
			}
		}
	}
	return set
})

// ── remote-resolved counter_id cache ──────────────────────────────

var (
	cacheMu     sync.RWMutex
	cacheSet    map[string]struct{}
	cacheLoaded bool

	// cacheDirOverride lets tests redirect the cache file away from the real
	// user cache directory. When empty, the path is derived from
	// LONGBRIDGE_CACHE_DIR / HOME as documented in cacheFilePath.
	cacheDirOverride string
)

// cacheFilePath returns the cache file path:
// "$LONGBRIDGE_CACHE_DIR/counter-ids.csv", defaulting to
// "~/.longbridge/cache/counter-ids.csv" (one counter_id per line, same format
// as the embedded directory files). Returns "" when no home directory can be
// determined.
func cacheFilePath() string {
	if cacheDirOverride != "" {
		return filepath.Join(cacheDirOverride, "counter-ids.csv")
	}
	if dir := os.Getenv("LONGBRIDGE_CACHE_DIR"); dir != "" {
		return filepath.Join(dir, "counter-ids.csv")
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".longbridge", "cache", "counter-ids.csv")
}

// cachedCounterIDs returns the in-memory remote-resolved cache set, loading it
// from disk on first access. The caller must hold cacheMu (read or write).
func loadCacheLocked() {
	if cacheLoaded {
		return
	}
	cacheLoaded = true
	cacheSet = make(map[string]struct{})
	path := cacheFilePath()
	if path == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			cacheSet[line] = struct{}{}
		}
	}
}

// CacheCounterIDs merges remotely resolved counter IDs into the local cache (in
// memory and on disk), so subsequent SymbolToCounterID / LookupCounterID calls
// resolve them without another network round trip. The cache file is rewritten
// with one counter_id per line in lexicographic order. Writing is skipped when
// no new entry is added.
func CacheCounterIDs(counterIDs []string) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	loadCacheLocked()

	before := len(cacheSet)
	for _, id := range counterIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			cacheSet[id] = struct{}{}
		}
	}
	if len(cacheSet) == before {
		return
	}

	path := cacheFilePath()
	if path == "" {
		return
	}
	if parent := filepath.Dir(path); parent != "" {
		_ = os.MkdirAll(parent, 0o755)
	}
	lines := make([]string, 0, len(cacheSet))
	for id := range cacheSet {
		lines = append(lines, id)
	}
	sort.Strings(lines)
	_ = os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

// cacheContains reports whether the given counter_id is present in the
// remote-resolved cache.
func cacheContains(counterID string) bool {
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	loadCacheLocked()
	_, ok := cacheSet[counterID]
	return ok
}

// rsplitDot splits a symbol on its last "." into (code, market). The second
// return value reports whether a "." was found.
func rsplitDot(symbol string) (code, market string, ok bool) {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return "", "", false
	}
	return symbol[:idx], symbol[idx+1:], true
}

// isAllDigits reports whether s consists solely of ASCII digits.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// LookupCounterID looks up a symbol in the local directory only (embedded
// special set, the remote-resolved cache, and leading-dot index notation). It
// returns (counterID, true) on a match, or ("", false) when the symbol is
// unknown locally — i.e. SymbolToCounterID would fall back to the default "ST/"
// prefix, which may be wrong for newly listed ETFs / indexes / warrants.
func LookupCounterID(symbol string) (string, bool) {
	code, market, ok := rsplitDot(symbol)
	if !ok {
		return "", false
	}
	market = strings.ToUpper(market)
	if strings.HasPrefix(code, ".") {
		return "IX/" + market + "/" + code, true
	}
	if market == "HK" && isAllDigits(code) {
		code = strings.TrimLeft(code, "0")
	}
	special := specialCounterIDs()
	for _, prefix := range []string{"ETF", "IX", "WT"} {
		candidate := prefix + "/" + market + "/" + code
		if _, found := special[candidate]; found {
			return candidate, true
		}
	}
	for _, prefix := range []string{"ETF", "IX", "WT", "ST"} {
		candidate := prefix + "/" + market + "/" + code
		if cacheContains(candidate) {
			return candidate, true
		}
	}
	return "", false
}

// SymbolToCounterID converts a user-supplied symbol (e.g. "TSLA.US", "700.HK",
// ".DJI.US", "HSI.HK") to a counter_id (e.g. "ST/US/TSLA", "ST/HK/700",
// "IX/US/.DJI", "IX/HK/HSI").
//
// Leading-dot symbols (e.g. ".DJI.US") are US market indexes and always map to
// "IX/". All other symbols are checked against the embedded ETF + index +
// warrant set and the remote-resolved cache; a matching entry is returned
// as-is. Unmatched symbols default to "ST/". An input without a "." is returned
// unchanged.
func SymbolToCounterID(symbol string) string {
	code, market, ok := rsplitDot(symbol)
	if !ok {
		return symbol
	}
	if counterID, found := LookupCounterID(symbol); found {
		return counterID
	}
	market = strings.ToUpper(market)
	// Strip leading zeros from numeric HK codes (e.g. "00700" → "700"). Other
	// markets keep their codes verbatim (A-share codes such as "000001.SZ"
	// have significant leading zeros).
	if market == "HK" && isAllDigits(code) {
		code = strings.TrimLeft(code, "0")
	}
	return "ST/" + market + "/" + code
}

// IndexSymbolToCounterID converts an index symbol (e.g. "HSI.HK") to a
// counter_id (e.g. "IX/HK/HSI"), always using the "IX/" prefix. An input
// without a "." is returned unchanged.
func IndexSymbolToCounterID(symbol string) string {
	code, market, ok := rsplitDot(symbol)
	if !ok {
		return symbol
	}
	return "IX/" + strings.ToUpper(market) + "/" + code
}

// CounterIDToSymbol converts a counter_id (e.g. "ST/US/TSLA", "ETF/US/SPY",
// "IX/US/.DJI", "ST/HK/700") back to a display symbol (e.g. "TSLA.US",
// "SPY.US", ".DJI.US", "700.HK").
//
// US index counter IDs ("IX/US/...") preserve the leading dot in the code part
// (e.g. "IX/US/.DJI" → ".DJI.US"). An input that is not in the three-segment
// counter_id format is returned unchanged.
func CounterIDToSymbol(counterID string) string {
	parts := strings.SplitN(counterID, "/", 3)
	if len(parts) == 3 {
		return parts[2] + "." + parts[1]
	}
	return counterID
}

// IsETF reports whether a user-supplied symbol resolves to an ETF (e.g.
// "QQQ.US", "SPY.US").
//
// Determined by checking the embedded special counter_id set: a symbol is an
// ETF when SymbolToCounterID maps it to an "ETF/..." counter_id.
func IsETF(symbol string) bool {
	return strings.HasPrefix(SymbolToCounterID(symbol), "ETF/")
}
