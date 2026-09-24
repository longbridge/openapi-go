package signal

import (
	"encoding/json"
	"time"
)

// Outlook is the direction a strategy expects the security to take.
type Outlook string

const (
	OutlookUnknown       Outlook = "Unknown"
	OutlookStrongBullish Outlook = "Strong bullish" // outlook_enum 1
	OutlookBullish       Outlook = "Bullish"        // outlook_enum 2
	OutlookNeutral       Outlook = "Neutral"        // outlook_enum 3
	OutlookBearish       Outlook = "Bearish"        // outlook_enum 4
	OutlookStrongBearish Outlook = "Strong bearish" // outlook_enum 5
)

// SignalStatus is where a signal is in its lifecycle.
type SignalStatus int32

const (
	// SignalStatusUnknown is a status this SDK does not know yet.
	SignalStatusUnknown SignalStatus = -1
	// SignalStatusPending is generated, not yet published.
	SignalStatusPending SignalStatus = 0
	// SignalStatusActive is published and current — the only status the API serves today.
	SignalStatusActive SignalStatus = 1
	// SignalStatusDeleted is deleted.
	SignalStatusDeleted SignalStatus = 2
	// SignalStatusAiFailed means the strategy analysis failed to generate.
	SignalStatusAiFailed SignalStatus = 3
	// SignalStatusFilteredByManual means a human reviewer filtered it out.
	SignalStatusFilteredByManual SignalStatus = 4
	// SignalStatusAiSubmitFailed means the strategy analysis failed to submit.
	SignalStatusAiSubmitFailed SignalStatus = 5
)

// FactType is the kind of a fact, and of the source that produced it.
type FactType string

const (
	FactTypeUnknown     FactType = "Unknown"
	FactTypeNews        FactType = "News"
	FactTypeFundamental FactType = "Fundamental"
	FactTypeTechnical   FactType = "Technical"
)

// FactDirection is the side a fact or factor points to.
type FactDirection string

const (
	FactDirectionUnknown FactDirection = "" // wire empty string
	FactDirectionLong    FactDirection = "long"
	FactDirectionShort   FactDirection = "short"
	FactDirectionNeutral FactDirection = "neutral"
)

// Signal is a strategy's take on a security, triggered by a catalyst.
type Signal struct {
	// ID is the signal ID, e.g. sign_992_1a00c9425c3_48ab. Pass it to
	// SignalContext.Signal for the full record.
	ID string
	// Symbol is the security symbol, e.g. 992.HK.
	Symbol string
	// CompanyName is the company name.
	CompanyName string
	// Market is the market the security trades in, e.g. HK.
	Market string
	// Title is the signal headline.
	Title string
	// Summary is a natural-language summary of the signal, in Markdown.
	Summary string
	// StrategyID is the strategy ID that produced the signal.
	StrategyID string
	// StrategyName is the strategy name that produced the signal.
	StrategyName string
	// RecommendBy is who recommended the signal; empty for strategy-generated signals.
	RecommendBy string
	// Expression is the strategy expression, e.g. 992.HK:GROWTH:long.
	Expression string
	// KeyFactID is the ID of the fact that triggered the signal.
	KeyFactID string
	// KeyCatalyst is the display label of the catalyst that triggered the signal,
	// e.g. "Q1 Revenue Surge". Filtering with SignalsOptions.CatalystName matches
	// the underlying factor name instead.
	KeyCatalyst string
	// AnalysisPrice is the price the analysis was based on.
	AnalysisPrice float64
	// ConservativePrice is the conservative-scenario target price.
	ConservativePrice float64
	// BenchmarkPrice is the benchmark-scenario target price.
	BenchmarkPrice float64
	// OptimisticPrice is the optimistic-scenario target price.
	OptimisticPrice float64
	// Outlook is the outlook the strategy takes on the security.
	Outlook Outlook
	// OutlookDesc is the localized rendering of Outlook in the caller's language.
	OutlookDesc string
	// Status is where the signal is in its lifecycle.
	Status SignalStatus
	// JsonData is the full analysis behind the signal, as a verbatim JSON document
	// (strategy-specific shape).
	JsonData string
	// CreatedAt is the creation time.
	CreatedAt time.Time
	// UpdatedAt is the last update time.
	UpdatedAt time.Time
}

// SignalsResponse is a page of signals returned by SignalContext.Signals.
type SignalsResponse struct {
	// Signals are the signals on this page.
	Signals []*Signal
	// Total is the total number of signals matching the filters, for paging with
	// SignalsOptions.Offset.
	Total int32
}

// SecurityFact is a fact (catalyst) event: something that happened to a security,
// with the factors, sources and prose behind it. Facts are what strategies react
// to — a signal names the fact that triggered it in Signal.KeyFactID.
type SecurityFact struct {
	// FactID is the fact ID, e.g. technical_rsi_14_short_1783674041337603409.
	FactID string
	// FactType is what kind of fact this is.
	FactType FactType
	// Direction is the side the fact points to.
	Direction FactDirection
	// OccurTime is when the fact occurred.
	OccurTime time.Time
	// SymbolsInfo are the securities the fact is about.
	SymbolsInfo []FactSymbol
	// Factors are the factors that contributed to the fact.
	Factors []FactFactor
	// DataSource is where the fact came from.
	DataSource []FactDataSource
	// NlInfo is the natural-language rendering of the fact.
	NlInfo FactNlInfo
}

// FactSymbol is a security a fact is about.
type FactSymbol struct {
	Symbol       string
	SecurityName string
}

// FactFactor is one factor that contributed to a fact.
type FactFactor struct {
	// Name is the factor name, e.g. rsi_14.
	Name string
	// FactorGroups are the groups the factor belongs to, e.g. MOMENTUM.
	FactorGroups []string
	// LongShortDirection is the side the factor points to.
	LongShortDirection FactDirection
	// TriggerCondition is the condition that fired the factor.
	TriggerCondition string
	// AnomalyDetection is the anomaly test behind the factor.
	AnomalyDetection AnomalyDetection
}

// AnomalyDetection is the outcome of the anomaly test behind a factor. Fields are
// empty for factors that did not run one.
type AnomalyDetection struct {
	AnomalyResult     string
	SignificanceLevel string
	TestMethod        string
	Thresholds        AnomalyThresholds
}

// AnomalyThresholds are the thresholds an anomaly test was scored against.
type AnomalyThresholds struct {
	Low    string
	Medium string
	High   string
}

// FactDataSource is where a fact came from.
type FactDataSource struct {
	SourceName string
	SourceType FactType
	URL        string
	Icon       string
}

// FactNlInfo is the natural-language rendering of a fact, in the caller's language.
type FactNlInfo struct {
	// Title is the headline.
	Title string
	// SubTitle is the sub-headline.
	SubTitle string
	// Summary is what happened, as a JSON array of {tag, value} entries carried in
	// a string. Use SummaryTags to read it.
	Summary string
	// InvestAnal is what it may mean for an investor, in the same JSON-in-a-string
	// shape as Summary. Use InvestAnalTags to read it.
	InvestAnal string
	// EliExplain is a plain-language walk-through of the fact, in the same
	// JSON-in-a-string shape as Summary. Use EliExplainTags to read it.
	EliExplain string
}

// NlTag is one {tag, value} entry from a natural-language field.
type NlTag struct {
	// Tag is what the entry is about, e.g. RSI.
	Tag string `json:"tag"`
	// Value is the prose.
	Value string `json:"value"`
}

// SummaryTags parses Summary into its {tag, value} entries. Returns an empty list
// when the field is empty or not the expected JSON — the raw string stays available.
func (n FactNlInfo) SummaryTags() []NlTag {
	return parseNlTags(n.Summary)
}

// InvestAnalTags parses InvestAnal into its {tag, value} entries. Returns an empty
// list when the field is empty or not the expected JSON.
func (n FactNlInfo) InvestAnalTags() []NlTag {
	return parseNlTags(n.InvestAnal)
}

// EliExplainTags parses EliExplain into its {tag, value} entries. Returns an empty
// list when the field is empty or not the expected JSON.
func (n FactNlInfo) EliExplainTags() []NlTag {
	return parseNlTags(n.EliExplain)
}

func parseNlTags(s string) []NlTag {
	if s == "" {
		return nil
	}
	var tags []NlTag
	if err := json.Unmarshal([]byte(s), &tags); err != nil {
		return nil
	}
	return tags
}
