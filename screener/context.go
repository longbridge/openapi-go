// Package screener provides a client for the Longbridge Screener OpenAPI.
// It covers stock screener strategies, indicator search, and pre-defined
// recommendation strategies.
package screener

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// ScreenerContext is a client for the Longbridge Screener OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	sctx, err := screener.NewFromCfg(conf)
//	recs, err := sctx.ScreenerRecommendStrategies(context.Background())
type ScreenerContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a ScreenerContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*ScreenerContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &ScreenerContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a ScreenerContext configured from environment variables.
func NewFromEnv() (*ScreenerContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// symbolToCounterID converts a symbol like "TSLA.US" to a counter_id like
// "ST/US/TSLA". All symbols are treated as equities (ST prefix).
func symbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])
	return fmt.Sprintf("ST/%s/%s", market, code)
}

// ─── ScreenerRecommendStrategies ──────────────────────────────────────────────

// ScreenerRecommendStrategies fetches the list of recommended screener strategies.
//
// Path: GET /v1/quote/screener/strategies/recommend
func (c *ScreenerContext) ScreenerRecommendStrategies(ctx context.Context) (*RecommendStrategiesResponse, error) {
	var resp json.RawMessage
	if err := c.httpClient.Get(ctx, "/v1/quote/screener/strategies/recommend", url.Values{}, &resp); err != nil {
		return nil, err
	}
	return &RecommendStrategiesResponse{Data: resp}, nil
}

// ─── ScreenerUserStrategies ───────────────────────────────────────────────────

// ScreenerUserStrategies fetches the current user's saved screener strategies.
//
// Path: GET /v1/quote/screener/strategies/mine
func (c *ScreenerContext) ScreenerUserStrategies(ctx context.Context) (*UserStrategiesResponse, error) {
	var resp json.RawMessage
	if err := c.httpClient.Get(ctx, "/v1/quote/screener/strategies/mine", url.Values{}, &resp); err != nil {
		return nil, err
	}
	return &UserStrategiesResponse{Data: resp}, nil
}

// ─── ScreenerStrategy ─────────────────────────────────────────────────────────

// ScreenerStrategy fetches a single screener strategy by ID.
//
// Path: GET /v1/quote/screener/strategy
func (c *ScreenerContext) ScreenerStrategy(ctx context.Context, id int64) (*StrategyResponse, error) {
	q := url.Values{}
	q.Set("id", fmt.Sprintf("%d", id))
	var resp json.RawMessage
	if err := c.httpClient.Get(ctx, "/v1/quote/screener/strategy", q, &resp); err != nil {
		return nil, err
	}
	return &StrategyResponse{Data: resp}, nil
}

// ─── ScreenerSearch ───────────────────────────────────────────────────────────

// ScreenerSearch executes a screener search.
//
// Path: POST /v1/quote/screener/search
//
// market is the market code, e.g. "US" or "HK".
// strategyID is optional; pass nil to search without a strategy filter.
// page and size control pagination (1-indexed).
func (c *ScreenerContext) ScreenerSearch(ctx context.Context, market string, strategyID *int64, page, size uint32) (*ScreenerSearchResponse, error) {
	body := map[string]interface{}{
		"market": market,
		"page":   page,
		"size":   size,
	}
	if strategyID != nil {
		body["strategy_id"] = *strategyID
	}
	var resp json.RawMessage
	if err := c.httpClient.Post(ctx, "/v1/quote/screener/search", body, &resp); err != nil {
		return nil, err
	}
	return &ScreenerSearchResponse{Data: resp}, nil
}

// ─── ScreenerIndicators ───────────────────────────────────────────────────────

// ScreenerIndicators fetches the list of available screener indicators.
//
// Path: GET /v1/quote/screener/indicators
func (c *ScreenerContext) ScreenerIndicators(ctx context.Context) (*ScreenerIndicatorsResponse, error) {
	var resp json.RawMessage
	if err := c.httpClient.Get(ctx, "/v1/quote/screener/indicators", url.Values{}, &resp); err != nil {
		return nil, err
	}
	return &ScreenerIndicatorsResponse{Data: resp}, nil
}
