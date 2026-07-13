package quote

import (
	"context"
	"net/url"

	"github.com/longbridge/openapi-go/internal/counter"
)

// CryptoOverview holds the market overview for a single cryptocurrency.
type CryptoOverview struct {
	// Symbol is the user-facing trading-pair symbol (e.g. "BTCUSD.BKKT"),
	// converted from the API's counter_id field (e.g. "VA/BKKT/BTCUSD").
	Symbol             string `json:"symbol"`
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

// cryptoRawOverview is the raw API shape before symbol conversion.
type cryptoRawOverview struct {
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
// symbol must be in PAIR.EXCHANGE format, e.g. "BTCUSD.BKKT" (US DC uses BKKT).
// The returned Symbol field is converted back from the API's counter_id.
//
// Path: GET /v1/us/gemini/crypto-overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *QuoteContext) CryptoOverview(ctx context.Context, symbol string) (*CryptoOverview, error) {
	if err := c.opts.httpClient.CheckRegion("/v1/us/gemini/crypto-overview", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counter.SymbolToID(symbol))
	var raw cryptoRawOverview
	if err := c.opts.httpClient.Get(ctx, "/v1/us/gemini/crypto-overview", q, &raw); err != nil {
		return nil, err
	}
	return &CryptoOverview{
		Symbol:             counter.IDToSymbol(raw.CounterID),
		Name:               raw.Name,
		Ticker:             raw.Ticker,
		BaseAsset:          raw.BaseAsset,
		Currency:           raw.Currency,
		AllTimeHigh:        raw.AllTimeHigh,
		AllTimeHighDate:    raw.AllTimeHighDate,
		AllTimeLow:         raw.AllTimeLow,
		AllTimeLowDate:     raw.AllTimeLowDate,
		IpoDate:            raw.IpoDate,
		IssuePrice:         raw.IssuePrice,
		Shares:             raw.Shares,
		OfficialWebAddress: raw.OfficialWebAddress,
		Logo:               raw.Logo,
		WikiURL:            raw.WikiURL,
		Profile:            raw.Profile,
	}, nil
}
