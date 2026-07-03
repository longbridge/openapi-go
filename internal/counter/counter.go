package counter

import (
	_ "embed"
	"strings"
	"sync"
)

//go:embed US-ETF.csv
var usEtfCSV string

//go:embed US-IX.csv
var usIxCSV string

//go:embed US-WT.csv
var usWtCSV string

// specialCounterIDs is the set of counter_ids (e.g. "ETF/US/SPY") that need
// a non-ST prefix. Populated once on first use.
var (
	specialOnce sync.Once
	specialIDs  map[string]string // "SPY" -> "ETF/US/SPY"
)

func loadSpecial() map[string]string {
	m := make(map[string]string, 8000)
	for _, csv := range []string{usEtfCSV, usIxCSV, usWtCSV} {
		for _, line := range strings.Split(csv, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// counter_id format: PREFIX/MARKET/CODE (e.g. ETF/US/SPY)
			parts := strings.SplitN(line, "/", 3)
			if len(parts) == 3 {
				m[parts[2]] = line
			}
		}
	}
	return m
}

func getSpecial() map[string]string {
	specialOnce.Do(func() {
		specialIDs = loadSpecial()
	})
	return specialIDs
}

// IDToSymbol converts a slash-format counter_id to a user-facing symbol.
//
// Conversions:
//   - "ST/US/TSLA"    → "TSLA.US"
//   - "ETF/US/SPY"    → "SPY.US"
//   - "VA/HAS/BTCUSD" → "BTCUSD.HAS"  (crypto: PAIR.EXCHANGE)
//
// Returns the input unchanged if it is not in the expected format.
func IDToSymbol(counterID string) string {
	parts := strings.SplitN(counterID, "/", 3)
	if len(parts) != 3 {
		return counterID
	}
	prefix, market, code := parts[0], parts[1], parts[2]
	// Crypto (VA prefix): return PAIR.EXCHANGE format
	if strings.EqualFold(prefix, "VA") {
		return code + "." + strings.ToUpper(market)
	}
	return code + "." + market
}

// cryptoExchanges is the set of known crypto exchange identifiers used as
// symbol suffixes (e.g. "BTCUSD.BKKT" → "VA/BKKT/BTCUSD").
//
// Region mapping:
//   - US DC: BKKT
//   - HK DC: HAS, OSL
var cryptoExchanges = map[string]bool{
	"BKKT": true, // US
	"HAS":  true, // HK
	"OSL":  true, // HK
}

// SymbolToID converts a user-facing symbol to its internal counter_id.
//
// Conversions:
//   - "TSLA.US"     → "ST/US/TSLA"    (stock)
//   - "SPY.US"      → "ETF/US/SPY"    (ETF, from embedded list)
//   - "BTCUSD.HAS"  → "VA/HAS/BTCUSD" (crypto: suffix is a known exchange)
//
// Stocks default to "ST/{MARKET}/{CODE}"; ETF/index/warrant are resolved via
// embedded CSV lists.
func SymbolToID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])

	// Known crypto exchange suffix → VA/{EXCHANGE}/{PAIR}
	if cryptoExchanges[market] {
		return "VA/" + market + "/" + code
	}

	// Check special prefixes (ETF, index, warrant).
	if id, ok := getSpecial()[code]; ok {
		return id
	}
	return "ST/" + market + "/" + code
}
