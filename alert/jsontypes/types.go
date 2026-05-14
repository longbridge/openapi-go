package jsontypes

import "strings"

type AlertList struct {
	Lists []*AlertSymbolGroup `json:"lists"`
}

type AlertSymbolGroup struct {
	CounterId string      `json:"counter_id"`
	Code      string      `json:"code"`
	Market    string      `json:"market"`
	Name      string      `json:"name"`
	Price     string      `json:"price"`
	Chg       string      `json:"chg"`
	PChg      string      `json:"p_chg"`
	Product   string      `json:"product"`
	Items     []*AlertItem `json:"indicators"`
}

type AlertItem struct {
	Id          string         `json:"id"`
	IndicatorId string         `json:"indicator_id"`
	Enabled     bool           `json:"enabled"`
	Frequency   int32          `json:"frequency"`
	Scope       int32          `json:"scope"`
	Text        string         `json:"text"`
	State       []int32        `json:"state"`
	ValueMap    map[string]interface{} `json:"value_map"`
}

// CounterIDToSymbol converts a counter_id like "ST/HK/700" to symbol "700.HK".
func CounterIDToSymbol(counterID string) string {
	parts := strings.SplitN(counterID, "/", 3)
	if len(parts) == 3 {
		return parts[2] + "." + parts[1]
	}
	return counterID
}

// SymbolToCounterID converts a symbol like "700.HK" to counter_id "ST/HK/700".
func SymbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])
	return "ST/" + market + "/" + code
}
