package jsontypes

type CalendarEventsResponse struct {
	Date string               `json:"date"`
	List []*CalendarDateGroup `json:"list"`
}

type CalendarDateGroup struct {
	Date  string               `json:"date"`
	Count int32                `json:"count"`
	Infos []*CalendarEventInfo `json:"infos"`
}

type CalendarEventInfo struct {
	CounterId           string             `json:"counter_id"`
	Market              string             `json:"market"`
	Content             string             `json:"content"`
	CounterName         string             `json:"counter_name"`
	DateType            string             `json:"date_type"`
	Date                string             `json:"date"`
	ChartUid            string             `json:"chart_uid"`
	DataKv              []*CalendarDataKv  `json:"data_kv"`
	EventType           string             `json:"type"`
	Datetime            string             `json:"datetime"`
	Icon                string             `json:"icon"`
	Star                int32              `json:"star"`
	Id                  string             `json:"id"`
	FinancialMarketTime string             `json:"financial_market_time"`
	Currency            string             `json:"currency"`
	ActivityType        string             `json:"activity_type"`
}

type CalendarDataKv struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ValueType string `json:"type"`
	ValueRaw  string `json:"value_raw"`
}
