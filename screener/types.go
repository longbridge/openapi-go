// Package screener provides a client for the Longbridge Screener OpenAPI.
// It covers stock screener strategies, indicator search, and pre-defined
// recommendation strategies.
package screener

import "encoding/json"

// RecommendStrategiesResponse holds the raw data for recommended screener
// strategies from GET /v1/quote/ai/screener/strategies/recommend.
type RecommendStrategiesResponse struct {
	Data json.RawMessage `json:"data"`
}

// UserStrategiesResponse holds the raw data for the current user's screener
// strategies from GET /v1/quote/ai/screener/strategies/mine.
type UserStrategiesResponse struct {
	Data json.RawMessage `json:"data"`
}

// StrategyResponse holds the raw data for a single screener strategy from
// GET /v1/quote/ai/screener/strategy/{id}.
type StrategyResponse struct {
	Data json.RawMessage `json:"data"`
}

// ScreenerSearchResponse holds the raw search results from
// POST /v1/quote/ai/screener/search.
type ScreenerSearchResponse struct {
	Data json.RawMessage `json:"data"`
}

// ScreenerIndicatorsResponse holds the raw list of screener indicators from
// GET /v1/quote/ai/screener/indicators.
type ScreenerIndicatorsResponse struct {
	Data json.RawMessage `json:"data"`
}

// ── ScreenerCondition ─────────────────────────────────────────────

// ScreenerCondition is a filter condition for ScreenerSearch Mode B.
type ScreenerCondition struct {
	// Key is the indicator key without "filter_" prefix, e.g. "pettm", "roe", "macd_day".
	Key string `json:"key"`
	// Min is the lower bound (empty string = no lower bound).
	Min string `json:"min"`
	// Max is the upper bound (empty string = no upper bound).
	Max string `json:"max"`
	// TechValues holds technical indicator params (nil map for fundamental indicators).
	// Example: {"category": "goldenfork", "period": "day"}
	TechValues map[string]string `json:"tech_values,omitempty"`
}
