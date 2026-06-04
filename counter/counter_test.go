package counter

import (
	"os"
	"path/filepath"
	"testing"
)

// resetCache clears the in-memory cache state and redirects the cache file to
// the given directory, so each test starts from a clean, isolated cache.
func resetCache(dir string) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cacheDirOverride = dir
	cacheSet = nil
	cacheLoaded = false
}

func TestSymbolToCounterID(t *testing.T) {
	cases := []struct {
		symbol string
		want   string
	}{
		{"TSLA.US", "ST/US/TSLA"},
		{"700.HK", "ST/HK/700"},
		{"00700.HK", "ST/HK/700"},
		{"09988.HK", "ST/HK/9988"},
		{"000001.SZ", "ST/SZ/000001"},
		{"SPY.US", "ETF/US/SPY"},
		{"QQQ.US", "ETF/US/QQQ"},
		{"DRAM.US", "ETF/US/DRAM"},
		{"SPY.us", "ETF/US/SPY"},
		{"NODOT", "NODOT"},
		{".DJI.US", "IX/US/.DJI"},
		{".VIX.US", "IX/US/.VIX"},
		{".IXIC.US", "IX/US/.IXIC"},
		{".SPX.US", "IX/US/.SPX"},
		{"HSI.HK", "IX/HK/HSI"},
		{"10005.HK", "WT/HK/10005"},
	}
	for _, c := range cases {
		if got := SymbolToCounterID(c.symbol); got != c.want {
			t.Errorf("SymbolToCounterID(%q) = %q, want %q", c.symbol, got, c.want)
		}
	}
}

func TestIsETF(t *testing.T) {
	for _, s := range []string{"QQQ.US", "SPY.US", "DRAM.US"} {
		if !IsETF(s) {
			t.Errorf("IsETF(%q) = false, want true", s)
		}
	}
	for _, s := range []string{"TSLA.US", "HSI.HK", "700.HK"} {
		if IsETF(s) {
			t.Errorf("IsETF(%q) = true, want false", s)
		}
	}
}

func TestIndexSymbolToCounterID(t *testing.T) {
	if got := IndexSymbolToCounterID("HSI.HK"); got != "IX/HK/HSI" {
		t.Errorf("IndexSymbolToCounterID(HSI.HK) = %q, want IX/HK/HSI", got)
	}
}

func TestCounterIDToSymbol(t *testing.T) {
	cases := []struct {
		counterID string
		want      string
	}{
		{"IX/US/.DJI", ".DJI.US"},
		{"IX/HK/HSI", "HSI.HK"},
		{"ST/US/TSLA", "TSLA.US"},
		{"NODOT", "NODOT"},
	}
	for _, c := range cases {
		if got := CounterIDToSymbol(c.counterID); got != c.want {
			t.Errorf("CounterIDToSymbol(%q) = %q, want %q", c.counterID, got, c.want)
		}
	}
}

func TestRoundtrip(t *testing.T) {
	cid := SymbolToCounterID("TSLA.US")
	if got := CounterIDToSymbol(cid); got != "TSLA.US" {
		t.Errorf("roundtrip = %q, want TSLA.US", got)
	}
}

func TestLookupKnownSpecial(t *testing.T) {
	resetCache(t.TempDir())
	if got, ok := LookupCounterID("QQQ.US"); !ok || got != "ETF/US/QQQ" {
		t.Errorf("LookupCounterID(QQQ.US) = %q,%v want ETF/US/QQQ,true", got, ok)
	}
	if got, ok := LookupCounterID("HSI.HK"); !ok || got != "IX/HK/HSI" {
		t.Errorf("LookupCounterID(HSI.HK) = %q,%v want IX/HK/HSI,true", got, ok)
	}
	if got, ok := LookupCounterID(".DJI.US"); !ok || got != "IX/US/.DJI" {
		t.Errorf("LookupCounterID(.DJI.US) = %q,%v want IX/US/.DJI,true", got, ok)
	}
	if _, ok := LookupCounterID("TSLA.US"); ok {
		t.Errorf("LookupCounterID(TSLA.US) ok = true, want false")
	}
	if _, ok := LookupCounterID("NODOT"); ok {
		t.Errorf("LookupCounterID(NODOT) ok = true, want false")
	}
}

func TestCacheCounterIDsRoundtrip(t *testing.T) {
	dir := t.TempDir()
	resetCache(dir)

	// Unknown symbol falls back to ST/ before caching.
	if _, ok := LookupCounterID("FAKE9.US"); ok {
		t.Fatalf("LookupCounterID(FAKE9.US) ok = true, want false")
	}
	if got := SymbolToCounterID("FAKE9.US"); got != "ST/US/FAKE9" {
		t.Fatalf("SymbolToCounterID(FAKE9.US) = %q, want ST/US/FAKE9", got)
	}

	// After caching remote-resolved entries, lookups return them — including
	// backend-confirmed ST/ entries.
	CacheCounterIDs([]string{"ETF/US/FAKE9", "ST/US/FAKE8"})
	if got, ok := LookupCounterID("FAKE9.US"); !ok || got != "ETF/US/FAKE9" {
		t.Errorf("LookupCounterID(FAKE9.US) = %q,%v want ETF/US/FAKE9,true", got, ok)
	}
	if got := SymbolToCounterID("FAKE9.US"); got != "ETF/US/FAKE9" {
		t.Errorf("SymbolToCounterID(FAKE9.US) = %q, want ETF/US/FAKE9", got)
	}
	if got, ok := LookupCounterID("FAKE8.US"); !ok || got != "ST/US/FAKE8" {
		t.Errorf("LookupCounterID(FAKE8.US) = %q,%v want ST/US/FAKE8,true", got, ok)
	}

	// Persisted to disk as one counter_id per line, lexicographically sorted.
	saved, err := os.ReadFile(filepath.Join(dir, "counter-ids.csv"))
	if err != nil {
		t.Fatalf("read cache file: %v", err)
	}
	if string(saved) != "ETF/US/FAKE9\nST/US/FAKE8\n" {
		t.Errorf("cache file = %q, want %q", string(saved), "ETF/US/FAKE9\nST/US/FAKE8\n")
	}
}

func TestCacheCounterIDsNoNewEntriesSkipsWrite(t *testing.T) {
	dir := t.TempDir()
	resetCache(dir)
	path := filepath.Join(dir, "counter-ids.csv")

	CacheCounterIDs([]string{"ETF/US/FAKE1"})
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	// Re-caching the same entry should not rewrite the file.
	CacheCounterIDs([]string{"ETF/US/FAKE1"})
	info2, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if !info.ModTime().Equal(info2.ModTime()) {
		t.Errorf("cache file rewritten when no new entries were added")
	}
}
