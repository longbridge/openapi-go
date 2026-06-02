package jsontypes

import "encoding/json"

// CalendarEventsResponse is the raw JSON response for /v1/quote/finance_calendar.
type CalendarEventsResponse struct {
	// Start date of the query window
	Date string `json:"date"`
	// Per-day event groups
	List []CalendarDateGroup `json:"list"`
	// Next page cursor; empty string means no more pages
	NextDate string `json:"next_date"`
}

// CalendarDateGroup holds all events for a single calendar date.
type CalendarDateGroup struct {
	// Date string, e.g. "2025-05-02"
	Date string `json:"date"`
	// Total event count for this date
	Count int32 `json:"count"`
	// Event details
	Infos []CalendarEventInfo `json:"infos"`
}

// CalendarEventInfo represents one financial calendar event.
type CalendarEventInfo struct {
	// Security symbol (mapped from counter_id by the API)
	Symbol string `json:"counter_id"`
	// Market, e.g. "HK"
	Market string `json:"market"`
	// Event content description
	Content string `json:"content"`
	// Security name
	CounterName string `json:"counter_name"`
	// Date type label, e.g. "盘前"
	DateType string `json:"date_type"`
	// Event date string, e.g. "2025.05.02"
	Date string `json:"date"`
	// Chart UID (may be empty)
	ChartUID string `json:"chart_uid"`
	// Structured data key-value pairs
	DataKV []CalendarDataKv `json:"data_kv"`
	// Event type code, e.g. "financial"
	EventType string `json:"type"`
	// Event datetime (unix timestamp string)
	Datetime string `json:"datetime"`
	// Icon URL
	Icon string `json:"icon"`
	// Importance star rating (0–3)
	Star int32 `json:"star"`
	// Associated live stream (usually null)
	Live *json.RawMessage `json:"live"`
	// Internal event ID
	ID string `json:"id"`
	// Financial market session time string
	FinancialMarketTime string `json:"financial_market_time"`
	// Currency
	Currency string `json:"currency"`
	// Extended data (structure varies by event type)
	Ext *json.RawMessage `json:"ext"`
	// Activity type code
	ActivityType string `json:"activity_type"`
}

// CalendarDataKv is one key-value data pair within a calendar event.
type CalendarDataKv struct {
	// Key (may be empty)
	Key string `json:"key"`
	// Formatted display value
	Value string `json:"value"`
	// Value type code, e.g. "estimate_eps"
	ValueType string `json:"type"`
	// Raw numeric value as string (may be empty or non-numeric)
	ValueRaw string `json:"value_raw"`
}
