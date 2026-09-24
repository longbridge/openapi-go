package signal

import (
	"net/url"
	"strconv"
	"time"
)

// SignalsOptions filters SignalContext.Signals. Every field is a filter; leaving
// one at its zero value removes that filter.
type SignalsOptions struct {
	// SymbolName filters by security symbol, e.g. AAPL.US or 700.HK.
	SymbolName string
	// StrategyID filters by strategy id, e.g. buffett-value. Preferred over the
	// deprecated StrategyName; takes precedence when both are provided.
	StrategyID string
	// StrategyName filters by strategy name. Deprecated in favour of StrategyID.
	StrategyName string
	// CatalystName filters by the name of the factor that triggered the signal,
	// e.g. EARNINGS_RELEASED or macd_12_26_9 — the triggering fact's factor name,
	// not the display label in Signal.KeyCatalyst.
	CatalystName string
	// CatalystType filters by the catalyst type, e.g. News, Fundamental, Technical.
	CatalystType string
	// StartTime only returns signals created at or after this time (zero = no bound).
	StartTime time.Time
	// EndTime only returns signals created at or before this time (zero = no bound).
	EndTime time.Time
	// Limit is the maximum number of results to return (nil = server default of 20).
	Limit *int32
	// Offset is the number of results to skip for pagination (nil = 0).
	Offset *int32
}

// Values renders the options into a query string.
func (o *SignalsOptions) Values() url.Values {
	v := url.Values{}
	if o == nil {
		return v
	}
	if o.SymbolName != "" {
		v.Set("symbol_name", o.SymbolName)
	}
	if o.StrategyID != "" {
		v.Set("strategy_id", o.StrategyID)
	}
	if o.StrategyName != "" {
		v.Set("strategy_name", o.StrategyName)
	}
	if o.CatalystName != "" {
		v.Set("catalyst_name", o.CatalystName)
	}
	if o.CatalystType != "" {
		v.Set("catalyst_type", o.CatalystType)
	}
	if !o.StartTime.IsZero() {
		v.Set("start_time", o.StartTime.UTC().Format(time.RFC3339))
	}
	if !o.EndTime.IsZero() {
		v.Set("end_time", o.EndTime.UTC().Format(time.RFC3339))
	}
	if o.Limit != nil {
		v.Set("limit", strconv.Itoa(int(*o.Limit)))
	}
	if o.Offset != nil {
		v.Set("offset", strconv.Itoa(int(*o.Offset)))
	}
	return v
}

// SecurityFactsOptions filters SignalContext.SecurityFacts.
type SecurityFactsOptions struct {
	// Symbol is the security to query, e.g. AAPL.US or 700.HK. Required.
	Symbol string
	// BeginTime is the start of the query window (zero = earliest available).
	BeginTime time.Time
	// EndTime is the end of the query window (zero = latest available).
	EndTime time.Time
	// Limit is the maximum number of facts to return; when more facts fall inside
	// the window, only the latest Limit are returned (nil = server default of 100).
	Limit *int32
}

// Values renders the options into a query string.
func (o *SecurityFactsOptions) Values() url.Values {
	v := url.Values{}
	if o == nil {
		return v
	}
	v.Set("symbol", o.Symbol)
	if !o.BeginTime.IsZero() {
		v.Set("begin_time", o.BeginTime.UTC().Format(time.RFC3339))
	}
	if !o.EndTime.IsZero() {
		v.Set("end_time", o.EndTime.UTC().Format(time.RFC3339))
	}
	if o.Limit != nil {
		v.Set("limit", strconv.Itoa(int(*o.Limit)))
	}
	return v
}
