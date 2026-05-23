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

// defaultReturns are the column keys always included in a screener search
// request body.
var defaultReturns = []string{
	"filter_prevclose",
	"filter_prevchg",
	"filter_marketcap",
	"filter_salesgrowthyoy",
	"filter_pettm",
	"filter_pbmrq",
	"filter_industry",
}

// ScreenerContext is a client for the Longbridge Screener OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	sctx, err := screener.NewFromCfg(conf)
//	recs, err := sctx.ScreenerRecommendStrategies(context.Background(), "US")
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

// stripFilterPrefix removes the "filter_" prefix from s, returning the
// remainder, or s unchanged if the prefix is absent.
func stripFilterPrefix(s string) string {
	return strings.TrimPrefix(s, "filter_")
}

// ensureFilterPrefix returns s unchanged if it already starts with "filter_",
// otherwise it prepends that prefix.
func ensureFilterPrefix(s string) string {
	if strings.HasPrefix(s, "filter_") {
		return s
	}
	return "filter_" + s
}

// ─── ScreenerRecommendStrategies ──────────────────────────────────────────────

// ScreenerRecommendStrategies fetches the list of recommended screener strategies.
//
// Path: GET /v1/quote/ai/screener/strategies/recommend
func (c *ScreenerContext) ScreenerRecommendStrategies(ctx context.Context, market string) (*RecommendStrategiesResponse, error) {
	q := url.Values{}
	q.Set("market", market)
	var resp json.RawMessage
	if err := c.httpClient.Get(ctx, "/v1/quote/ai/screener/strategies/recommend", q, &resp); err != nil {
		return nil, err
	}
	return &RecommendStrategiesResponse{Data: resp}, nil
}

// ─── ScreenerUserStrategies ───────────────────────────────────────────────────

// ScreenerUserStrategies fetches the current user's saved screener strategies.
//
// Path: GET /v1/quote/ai/screener/strategies/mine
func (c *ScreenerContext) ScreenerUserStrategies(ctx context.Context, market string) (*UserStrategiesResponse, error) {
	q := url.Values{}
	q.Set("market", market)
	var resp json.RawMessage
	if err := c.httpClient.Get(ctx, "/v1/quote/ai/screener/strategies/mine", q, &resp); err != nil {
		return nil, err
	}
	return &UserStrategiesResponse{Data: resp}, nil
}

// ─── ScreenerStrategy ─────────────────────────────────────────────────────────

// ScreenerStrategy fetches a single screener strategy by ID.
//
// Path: GET /v1/quote/ai/screener/strategy/{id}
//
// The "filter_" prefix is stripped from every filters[].key before
// returning so callers see clean keys like "pettm" instead of
// "filter_pettm".
func (c *ScreenerContext) ScreenerStrategy(ctx context.Context, id int64) (*StrategyResponse, error) {
	path := fmt.Sprintf("/v1/quote/ai/screener/strategy/%d", id)
	var raw map[string]interface{}
	if err := c.httpClient.Get(ctx, path, url.Values{}, &raw); err != nil {
		return nil, err
	}
	// Strip filter_ prefix from filter.filters[].key
	if filterObj, ok := raw["filter"].(map[string]interface{}); ok {
		if filters, ok := filterObj["filters"].([]interface{}); ok {
			for _, f := range filters {
				if fmap, ok := f.(map[string]interface{}); ok {
					if k, ok := fmap["key"].(string); ok {
						fmap["key"] = stripFilterPrefix(k)
					}
				}
			}
		}
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	return &StrategyResponse{Data: json.RawMessage(b)}, nil
}

// ─── ScreenerSearch ───────────────────────────────────────────────────────────

// ScreenerSearch executes a screener search.
//
// Path: POST /v1/quote/ai/screener/search
//
// market is the market code, e.g. "US" or "HK".
// strategyID is optional; pass nil to use custom conditions instead.
// conditions is a list of "KEY:MIN:MAX" strings used in Mode B (strategyID == nil).
// show is an optional list of extra return columns to include.
// page is 0-indexed; size is the page size.
//
// Mode A (strategyID given): the strategy is fetched from
// GET /v1/quote/ai/screener/strategy/{id}, its filter.filters[] are
// forwarded to the search endpoint, and market is taken from the strategy
// response.
//
// Mode B (strategyID nil): conditions drive the filters and the supplied
// market is used directly.
//
// The "filter_" prefix is stripped from every items[].indicators[].key in
// the response before it is returned.
func (c *ScreenerContext) ScreenerSearch(
	ctx context.Context,
	market string,
	strategyID *int64,
	conditions []string,
	show []string,
	page, size uint32,
) (*ScreenerSearchResponse, error) {
	var effectiveMarket string
	var filters []map[string]interface{}

	if strategyID != nil {
		// Mode A: fetch strategy from AI endpoint
		path := fmt.Sprintf("/v1/quote/ai/screener/strategy/%d", *strategyID)
		var strategy map[string]interface{}
		if err := c.httpClient.Get(ctx, path, url.Values{}, &strategy); err != nil {
			return nil, err
		}
		mkt := ""
		if v, ok := strategy["market"].(string); ok {
			mkt = strings.ToUpper(v)
		}
		if mkt == "" || mkt == "-" {
			mkt = "US"
		}
		effectiveMarket = mkt

		if filterObj, ok := strategy["filter"].(map[string]interface{}); ok {
			if rawFilters, ok := filterObj["filters"].([]interface{}); ok {
				for _, f := range rawFilters {
					fmap, ok := f.(map[string]interface{})
					if !ok {
						continue
					}
					key, _ := fmap["key"].(string)
					if key == "" {
						continue
					}
					min, _ := fmap["min"].(string)
					max, _ := fmap["max"].(string)
					techValues := fmap["tech_values"]
					if techValues == nil {
						techValues = map[string]interface{}{}
					}
					filters = append(filters, map[string]interface{}{
						"key":         key,
						"min":         min,
						"max":         max,
						"tech_values": techValues,
					})
				}
			}
		}
	} else {
		// Mode B: custom conditions ("KEY:MIN:MAX")
		effectiveMarket = market
		for _, cond := range conditions {
			parts := strings.SplitN(cond, ":", 3)
			if len(parts) == 0 || parts[0] == "" {
				continue
			}
			key := parts[0]
			min := ""
			if len(parts) > 1 {
				min = parts[1]
			}
			max := ""
			if len(parts) > 2 {
				max = parts[2]
			}
			filters = append(filters, map[string]interface{}{
				"key":         key,
				"min":         min,
				"max":         max,
				"tech_values": map[string]interface{}{},
			})
		}
	}

	// Build returns list: always include defaultReturns, then filter keys, then show
	returnsSet := make(map[string]struct{})
	returnsList := make([]string, 0, len(defaultReturns)+len(filters)+len(show))
	addReturn := func(k string) {
		k = ensureFilterPrefix(k)
		if _, exists := returnsSet[k]; !exists {
			returnsSet[k] = struct{}{}
			returnsList = append(returnsList, k)
		}
	}
	for _, r := range defaultReturns {
		addReturn(r)
	}
	for _, f := range filters {
		if k, ok := f["key"].(string); ok && k != "" {
			addReturn(k)
		}
	}
	for _, s := range show {
		addReturn(s)
	}

	body := map[string]interface{}{
		"market":  effectiveMarket,
		"filters": filters,
		"returns": returnsList,
		"page":    page,
		"size":    size,
	}
	if filters == nil {
		body["filters"] = []interface{}{}
	}

	var raw map[string]interface{}
	if err := c.httpClient.Post(ctx, "/v1/quote/ai/screener/search", body, &raw); err != nil {
		return nil, err
	}

	// Strip filter_ prefix from items[].indicators[].key
	if items, ok := raw["items"].([]interface{}); ok {
		for _, item := range items {
			imap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if indicators, ok := imap["indicators"].([]interface{}); ok {
				for _, ind := range indicators {
					indmap, ok := ind.(map[string]interface{})
					if !ok {
						continue
					}
					if k, ok := indmap["key"].(string); ok {
						indmap["key"] = stripFilterPrefix(k)
					}
				}
			}
		}
	}

	b, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	return &ScreenerSearchResponse{Data: json.RawMessage(b)}, nil
}

// ─── ScreenerIndicators ───────────────────────────────────────────────────────

// ScreenerIndicators fetches the list of available screener indicators.
//
// Path: GET /v1/quote/ai/screener/indicators
//
// Post-processing applied before returning:
//   - "filter_" prefix is stripped from every groups[].indicators[].key
//   - tech_values is built from tech_indicators:
//     {tech_key: [{value, label}]}
func (c *ScreenerContext) ScreenerIndicators(ctx context.Context) (*ScreenerIndicatorsResponse, error) {
	var raw map[string]interface{}
	if err := c.httpClient.Get(ctx, "/v1/quote/ai/screener/indicators", url.Values{}, &raw); err != nil {
		return nil, err
	}

	if groups, ok := raw["groups"].([]interface{}); ok {
		for _, group := range groups {
			gmap, ok := group.(map[string]interface{})
			if !ok {
				continue
			}
			indicators, ok := gmap["indicators"].([]interface{})
			if !ok {
				continue
			}
			for _, ind := range indicators {
				indmap, ok := ind.(map[string]interface{})
				if !ok {
					continue
				}
				// Strip filter_ prefix
				if k, ok := indmap["key"].(string); ok {
					indmap["key"] = stripFilterPrefix(k)
				}
				// Build tech_values from tech_indicators
				techInds, ok := indmap["tech_indicators"].([]interface{})
				if !ok || len(techInds) == 0 {
					continue
				}
				tv := make(map[string]interface{})
				for _, ti := range techInds {
					timap, ok := ti.(map[string]interface{})
					if !ok {
						continue
					}
					techKey, _ := timap["tech_key"].(string)
					if techKey == "" {
						continue
					}
					var opts []map[string]interface{}
					if items, ok := timap["tech_items"].([]interface{}); ok {
						for _, item := range items {
							imap, ok := item.(map[string]interface{})
							if !ok {
								continue
							}
							val, _ := imap["item_value"].(string)
							label, _ := imap["item_name"].(string)
							opts = append(opts, map[string]interface{}{
								"value": val,
								"label": label,
							})
						}
					}
					tv[techKey] = opts
				}
				if len(tv) > 0 {
					indmap["tech_values"] = tv
				}
			}
		}
	}

	b, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	return &ScreenerIndicatorsResponse{Data: json.RawMessage(b)}, nil
}
