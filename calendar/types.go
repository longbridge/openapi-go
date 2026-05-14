package calendar

// CalendarCategory is the type of financial calendar event.
type CalendarCategory string

const (
	CalendarCategoryReport   CalendarCategory = "report"
	CalendarCategoryDividend CalendarCategory = "dividend"
	CalendarCategorySplit    CalendarCategory = "split"
	CalendarCategoryIpo      CalendarCategory = "ipo"
	CalendarCategoryMacroData CalendarCategory = "macrodata"
	CalendarCategoryClosed   CalendarCategory = "closed"
	CalendarCategoryMeeting  CalendarCategory = "meeting"
	CalendarCategoryMerge    CalendarCategory = "merge"
)

// CalendarEventsResponse is the response for FinanceCalendar.
type CalendarEventsResponse struct {
	Date string
	List []*CalendarDateGroup
}

// CalendarDateGroup groups events by date.
type CalendarDateGroup struct {
	Date  string
	Count int32
	Infos []*CalendarEventInfo
}

// CalendarEventInfo is a single financial calendar event.
type CalendarEventInfo struct {
	Symbol              string
	Market              string
	Content             string
	CounterName         string
	DateType            string
	Date                string
	ChartUid            string
	DataKv              []*CalendarDataKv
	EventType           string
	Datetime            string
	Icon                string
	Star                int32
	Id                  string
	FinancialMarketTime string
	Currency            string
	ActivityType        string
}

// CalendarDataKv is a key-value data pair in a calendar event.
type CalendarDataKv struct {
	Key       string
	Value     string
	ValueType string
	ValueRaw  string
}
