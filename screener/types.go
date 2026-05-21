// Package screener provides a client for the Longbridge Screener OpenAPI.
// It covers stock screener strategies, indicator search, and pre-defined
// recommendation strategies.
package screener

import "encoding/json"

// RecommendStrategiesResponse holds the raw data for recommended screener
// strategies from GET /v1/quote/screener/strategies/recommend.
type RecommendStrategiesResponse struct {
	Data json.RawMessage `json:"data"`
}

// UserStrategiesResponse holds the raw data for the current user's screener
// strategies from GET /v1/quote/screener/strategies/mine.
type UserStrategiesResponse struct {
	Data json.RawMessage `json:"data"`
}

// StrategyResponse holds the raw data for a single screener strategy from
// GET /v1/quote/screener/strategy.
type StrategyResponse struct {
	Data json.RawMessage `json:"data"`
}

// ScreenerSearchResponse holds the raw search results from
// POST /v1/quote/screener/search.
type ScreenerSearchResponse struct {
	Data json.RawMessage `json:"data"`
}

// ScreenerIndicatorsResponse holds the raw list of screener indicators from
// GET /v1/quote/screener/indicators.
type ScreenerIndicatorsResponse struct {
	Data json.RawMessage `json:"data"`
}
