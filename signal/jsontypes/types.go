// Package jsontypes holds raw JSON-deserialization structs for the signal package.
// These types use the exact API field names with json tags and are not part of the
// public surface; callers should use the types in the parent signal package.
package jsontypes

import (
	"strconv"
	"strings"
	"time"
)

// MilliTimestamp is a unix-millisecond timestamp that the API may send as either a
// JSON number or a JSON string. It tolerates both.
type MilliTimestamp int64

// UnmarshalJSON accepts a bare number (123), a quoted number ("123"), or null.
func (m *MilliTimestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(strings.TrimSpace(string(b)), `"`)
	if s == "" || s == "null" {
		*m = 0
		return nil
	}
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		*m = MilliTimestamp(v)
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*m = MilliTimestamp(int64(f))
	return nil
}

// Time converts the millisecond timestamp to a time.Time (zero value when unset).
func (m MilliTimestamp) Time() time.Time {
	if m == 0 {
		return time.Time{}
	}
	return time.UnixMilli(int64(m))
}

// SignalsResponse is the raw response for GET /v1/signals.
type SignalsResponse struct {
	Signals []Signal `json:"signals"`
	Total   int32    `json:"total"`
}

// SignalResponse is the raw response for GET /v1/signals/{signal_id}.
type SignalResponse struct {
	Signal Signal `json:"signal"`
}

// Signal is the raw wire type for one strategy signal.
type Signal struct {
	ID                string         `json:"id"`
	Symbol            string         `json:"symbol"`
	CompanyName       string         `json:"company_name"`
	Market            string         `json:"market"`
	Title             string         `json:"title"`
	Summary           string         `json:"summary"`
	StrategyID        string         `json:"strategy_id"`
	StrategyName      string         `json:"strategy_name"`
	RecommendBy       string         `json:"recommend_by"`
	Expression        string         `json:"expression"`
	KeyFactID         string         `json:"key_fact_id"`
	KeyCatalyst       string         `json:"key_catalyst"`
	AnalysisPrice     float64        `json:"analysis_price"`
	ConservativePrice float64        `json:"conservative_price"`
	BenchmarkPrice    float64        `json:"benchmark_price"`
	OptimisticPrice   float64        `json:"optimistic_price"`
	Outlook           string         `json:"outlook"`
	OutlookDesc       string         `json:"outlook_desc"`
	Status            int32          `json:"status"`
	JsonData          string         `json:"json_data"`
	CreatedAtMs       MilliTimestamp `json:"created_at"` // unix milliseconds (string or number on the wire)
	UpdatedAtMs       MilliTimestamp `json:"updated_at"` // unix milliseconds (string or number on the wire)
}

// SecurityFactsResponse is the raw response for GET /v1/facts/security_facts.
type SecurityFactsResponse struct {
	Facts []SecurityFact `json:"facts"`
}

// SecurityFact is the raw wire type for a fact (catalyst) event.
type SecurityFact struct {
	FactID      string           `json:"fact_id"`
	FactType    string           `json:"fact_type"`
	Direction   string           `json:"direction"`
	OccurTime   string           `json:"occur_time"` // RFC3339
	SymbolsInfo []FactSymbol     `json:"symbols_info"`
	Factors     []FactFactor     `json:"factors"`
	DataSource  []FactDataSource `json:"data_source"`
	NlInfo      FactNlInfo       `json:"nl_info"`
}

// FactSymbol is a security a fact is about.
type FactSymbol struct {
	Symbol       string `json:"symbol"`
	SecurityName string `json:"security_name"`
}

// FactFactor is one factor that contributed to a fact.
type FactFactor struct {
	Name               string           `json:"name"`
	FactorGroups       []string         `json:"factor_groups"`
	LongShortDirection string           `json:"long_short_direction"`
	TriggerCondition   string           `json:"trigger_condition"`
	AnomalyDetection   AnomalyDetection `json:"anomaly_detection"`
}

// AnomalyDetection is the outcome of the anomaly test behind a factor.
type AnomalyDetection struct {
	AnomalyResult     string            `json:"anomaly_result"`
	SignificanceLevel string            `json:"significance_level"`
	TestMethod        string            `json:"test_method"`
	Thresholds        AnomalyThresholds `json:"thresholds"`
}

// AnomalyThresholds are the thresholds an anomaly test was scored against.
type AnomalyThresholds struct {
	Low    string `json:"low"`
	Medium string `json:"medium"`
	High   string `json:"high"`
}

// FactDataSource is where a fact came from.
type FactDataSource struct {
	SourceName string `json:"source_name"`
	SourceType string `json:"type"`
	URL        string `json:"url"`
	Icon       string `json:"icon"`
}

// FactNlInfo is the natural-language rendering of a fact.
type FactNlInfo struct {
	Title      string `json:"title"`
	SubTitle   string `json:"sub_title"`
	Summary    string `json:"summary"`
	InvestAnal string `json:"invest_anal"`
	EliExplain string `json:"eli_explain"`
}
