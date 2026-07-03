package quote

import (
	"context"
	"net/url"

	"github.com/longbridge/openapi-go/internal/counter"
)

// CryptoOverview holds the market overview for a single cryptocurrency.
type CryptoOverview struct {
	CounterID          string `json:"counter_id"`
	Name               string `json:"name"`
	Ticker             string `json:"ticker"`
	BaseAsset          string `json:"base_asset"`
	Currency           string `json:"currency"`
	AllTimeHigh        string `json:"all_time_high"`
	AllTimeHighDate    string `json:"all_time_high_date"`
	AllTimeLow         string `json:"all_time_low"`
	AllTimeLowDate     string `json:"all_time_low_date"`
	IpoDate            string `json:"ipo_date"`
	IssuePrice         string `json:"issue_price"`
	Shares             string `json:"shares"`
	OfficialWebAddress string `json:"official_web_address"`
	Logo               string `json:"logo"`
	WikiURL            string `json:"wiki_url"`
	Profile            string `json:"profile"`
}

// CryptoOverview returns market overview data for a cryptocurrency.
//
// symbol must be in PAIR.EXCHANGE format, e.g. "BTCUSD.BKKT" → VA/BKKT/BTCUSD (US DC).
// Uses counter.SymbolToID for conversion, consistent with all other symbol-based methods.
//
// Path: GET /v1/gemini/us/crypto-overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *QuoteContext) CryptoOverview(ctx context.Context, symbol string) (*CryptoOverview, error) {
	if err := c.opts.httpClient.CheckRegion("/v1/gemini/us/crypto-overview", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counter.SymbolToID(symbol))
	var resp CryptoOverview
	if err := c.opts.httpClient.Get(ctx, "/v1/gemini/us/crypto-overview", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
