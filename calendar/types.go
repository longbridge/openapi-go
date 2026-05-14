// Package calendar provides a client for the Longbridge Financial Calendar OpenAPI.
// It covers earnings reports, dividends, stock splits, IPOs, macro data releases,
// market closures, shareholder meetings, and stock mergers.
package calendar

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// CalendarCategory identifies the type of financial calendar event to query.
type CalendarCategory int

const (
	// CalendarCategoryReport represents earnings report events.
	CalendarCategoryReport CalendarCategory = iota
	// CalendarCategoryDividend represents dividend distribution events.
	CalendarCategoryDividend
	// CalendarCategorySplit represents stock split events.
	CalendarCategorySplit
	// CalendarCategoryIpo represents initial public offering events.
	CalendarCategoryIpo
	// CalendarCategoryMacroData represents macro-economic data releases.
	CalendarCategoryMacroData
	// CalendarCategoryClosed represents market closure days.
	CalendarCategoryClosed
	// CalendarCategoryMeeting represents shareholder or analyst meeting events.
	CalendarCategoryMeeting
	// CalendarCategoryMerge represents stock consolidation or merger events.
	CalendarCategoryMerge
)

// calendarCategoryString maps CalendarCategory values to their API string representation.
var calendarCategoryString = map[CalendarCategory]string{
	CalendarCategoryReport:    "report",
	CalendarCategoryDividend:  "dividend",
	CalendarCategorySplit:     "split",
	CalendarCategoryIpo:       "ipo",
	CalendarCategoryMacroData: "macrodata",
	CalendarCategoryClosed:    "closed",
	CalendarCategoryMeeting:   "meeting",
	CalendarCategoryMerge:     "merge",
}

// String returns the API wire string for a CalendarCategory.
func (c CalendarCategory) String() string {
	if s, ok := calendarCategoryString[c]; ok {
		return s
	}
	return "report"
}

// CalendarEventsResponse is the top-level response from the finance_calendar endpoint.
type CalendarEventsResponse struct {
	// Start date of the query window, e.g. "2025-05-01"
	Date string
	// Per-day event groups
	List []CalendarDateGroup
}

// CalendarDateGroup holds all events for a single calendar date.
type CalendarDateGroup struct {
	// Date string, e.g. "2025-05-02"
	Date string
	// Total event count for this date
	Count int32
	// Individual event records
	Infos []CalendarEventInfo
}

// CalendarEventInfo represents one financial calendar event.
type CalendarEventInfo struct {
	// Security symbol (converted from counter_id)
	Symbol string
	// Market identifier, e.g. "HK"
	Market string
	// Human-readable event content description
	Content string
	// Security name
	CounterName string
	// Date type label, e.g. "盘前" (pre-market)
	DateType string
	// Event date string, e.g. "2025.05.02"
	Date string
	// Chart UID (may be empty)
	ChartUID string
	// Structured key-value data attached to the event
	DataKV []CalendarDataKv
	// Event type code, e.g. "financial"
	EventType string
	// Event datetime as unix timestamp string
	Datetime string
	// Icon URL
	Icon string
	// Importance star rating, 0–3
	Star int32
	// Raw live stream JSON (usually nil)
	Live *json.RawMessage
	// Internal event ID
	ID string
	// Financial market session time string
	FinancialMarketTime string
	// Currency
	Currency string
	// Extended data (structure varies by event type; nil when absent)
	Ext *json.RawMessage
	// Activity type code
	ActivityType string
}

// CalendarDataKv is one key-value data pair within a calendar event.
type CalendarDataKv struct {
	// Key label (may be empty)
	Key string
	// Formatted display value string
	Value string
	// Value type code, e.g. "estimate_eps"
	ValueType string
	// Raw numeric value; nil when the field is absent or non-numeric
	ValueRaw *decimal.Decimal
}
