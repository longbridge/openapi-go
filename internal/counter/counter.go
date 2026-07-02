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

// IDToSymbol converts a slash-format counter_id like "ST/US/TSLA" to a symbol
// like "TSLA.US". Returns the input unchanged if it is not in the expected format.
func IDToSymbol(counterID string) string {
	parts := strings.SplitN(counterID, "/", 3)
	if len(parts) == 3 {
		return parts[2] + "." + parts[1]
	}
	return counterID
}

// SymbolToID converts a user-facing symbol like "TSLA.US" to its counter_id.
// ETF, index, and warrant symbols are resolved to their correct prefix
// (e.g. "SPY.US" → "ETF/US/SPY"); ordinary stocks default to "ST/{MARKET}/{CODE}".
func SymbolToID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])

	// Check special prefixes first (ETF, index, warrant).
	if id, ok := getSpecial()[code]; ok {
		return id
	}
	return "ST/" + market + "/" + code
}
