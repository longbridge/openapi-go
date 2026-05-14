package util

import "strings"

// CounterIDToSymbol converts a counter_id (e.g. "ST/HK/700") to a symbol (e.g. "700.HK").
func CounterIDToSymbol(counterID string) string {
	parts := strings.SplitN(counterID, "/", 3)
	if len(parts) == 3 {
		return parts[2] + "." + parts[1]
	}
	return counterID
}

// SymbolToCounterID converts a symbol (e.g. "700.HK") to a counter_id (e.g. "ST/HK/700").
func SymbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])
	return "ST/" + market + "/" + code
}

// IndexSymbolToCounterID converts an index symbol (e.g. "HSI.HK") to counter_id ("IX/HK/HSI").
func IndexSymbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])
	return "IX/" + market + "/" + code
}
