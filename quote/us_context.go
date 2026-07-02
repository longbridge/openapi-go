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

// cryptoSymbolToCounterID converts a crypto symbol in PAIR.EXCHANGE format
// to the internal VA/{EXCHANGE}/{PAIR} counter_id.
//
// Example: "BTCUSD.HAS" → "VA/HAS/BTCUSD"
func cryptoSymbolToCounterID(symbol string) string {
	if idx := strings.LastIndex(symbol, "."); idx > 0 {
		pair := symbol[:idx]
		exchange := strings.ToUpper(symbol[idx+1:])
		return "VA/" + exchange + "/" + pair
	}
	// No exchange suffix — pass through as-is for forward compatibility.
	return symbol
}

// CryptoOverview returns market overview data for a cryptocurrency.
//
// symbol must be in PAIR.EXCHANGE format, e.g. "BTCUSD.HAS" → VA/HAS/BTCUSD.
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
