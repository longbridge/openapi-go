package quote

import (
	"context"
	"errors"
	"net/url"
)

// ErrUSOnly is returned when a US-only API is called with a non-US token.
var ErrUSOnly = errors.New("longbridge: this API is only available for US accounts (us_ token required)")

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

// CryptoOverview returns market overview data for a cryptocurrency.
//
// counterID is the crypto counter_id, e.g. "CY/US/BTC".
//
// Path: GET /v1/gemini/crypto-overview
// US token required; returns ErrUSOnly for non-US credentials.
func (c *QuoteContext) CryptoOverview(ctx context.Context, counterID string) (*CryptoOverview, error) {
	if !c.opts.httpClient.IsUS() {
		return nil, ErrUSOnly
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	var resp CryptoOverview
	if err := c.opts.httpClient.Get(ctx, "/v1/gemini/crypto-overview", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
