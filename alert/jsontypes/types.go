package jsontypes

type AlertList struct {
	Lists []*AlertSymbolGroup `json:"lists"`
}

type AlertSymbolGroup struct {
	CounterId string       `json:"counter_id"`
	Code      string       `json:"code"`
	Market    string       `json:"market"`
	Name      string       `json:"name"`
	Price     string       `json:"price"`
	Chg       string       `json:"chg"`
	PChg      string       `json:"p_chg"`
	Product   string       `json:"product"`
	Items     []*AlertItem `json:"indicators"`
}

type AlertItem struct {
	Id          string                 `json:"id"`
	IndicatorId string                 `json:"indicator_id"`
	Enabled     bool                   `json:"enabled"`
	Frequency   int32                  `json:"frequency"`
	Scope       int32                  `json:"scope"`
	Text        string                 `json:"text"`
	State       []int32                `json:"state"`
	ValueMap    map[string]interface{} `json:"value_map"`
}
