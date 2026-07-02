package quote

import (
	"context"
	"net/url"
	"strings"
)

// CryptoOverview holds the market overview for a single cryptocurrency.
type CryptoOverview struct {
	Name               string      `json:"name"`
	Ticker             string      `json:"ticker"`
	Currency           string      `json:"currency"`
	AllTimeHigh        string      `json:"all_time_high"`
	AllTimeHighDate    string      `json:"all_time_high_date"`
	AllTimeLow         string      `json:"all_time_low"`
	AllTimeLowDate     string      `json:"all_time_low_date"`
	IpoDate            string      `json:"ipo_date"`
	IssuePrice         string      `json:"issue_price"`
	Shares             string      `json:"shares"`
	OfficialWebAddress string      `json:"official_web_address"`
	Profile            interface{} `json:"profile"`
}

// cryptoSymbolToCounterID converts a user-facing crypto trading pair to the
// internal VA/HAS/... counter_id format.
//
// Accepts: "BTCUSD" or "BTC/USD" → "VA/HAS/BTCUSD"
// Rule: remove "/" separator and prepend "VA/HAS/".
func cryptoSymbolToCounterID(symbol string) string {
	return "VA/HAS/" + strings.ReplaceAll(symbol, "/", "")
}

// CryptoOverview returns market overview data for a cryptocurrency.
//
// symbol accepts "BTCUSD" or "BTC/USD" format; it is converted to the
// internal VA/HAS/BTCUSDT counter_id automatically.
//
// Path: GET /v1/gemini/us/crypto-overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *QuoteContext) CryptoOverview(ctx context.Context, symbol string) (*CryptoOverview, error) {
	if err := c.opts.httpClient.CheckRegion("/v1/gemini/us/crypto-overview", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", cryptoSymbolToCounterID(symbol))
	var resp CryptoOverview
	if err := c.opts.httpClient.Get(ctx, "/v1/gemini/us/crypto-overview", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
