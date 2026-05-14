// Package market provides a client for the Longbridge Market OpenAPI.
// It covers market status, broker holdings, A/H premium, trade statistics,
// market anomalies, and index constituents.
package market

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/market/jsontypes"
	httplib "github.com/longbridge/openapi-go/http"
)

// MarketContext is a client for the Longbridge Market OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	mctx, err := market.NewFromCfg(conf)
//	status, err := mctx.MarketStatus(context.Background())
type MarketContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a MarketContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*MarketContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &MarketContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a MarketContext configured from environment variables.
func NewFromEnv() (*MarketContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// MarketStatus returns the current trading status for all markets.
//
// Path: GET /v1/quote/market-status
func (m *MarketContext) MarketStatus(ctx context.Context) (*MarketStatusResponse, error) {
	var resp jsontypes.MarketStatusResponse
	if err := m.httpClient.Get(ctx, "/v1/quote/market-status", url.Values{}, &resp); err != nil {
		return nil, err
	}
	out := &MarketStatusResponse{
		MarketTime: make([]MarketTimeItem, 0, len(resp.MarketTime)),
	}
	for _, item := range resp.MarketTime {
		out.MarketTime = append(out.MarketTime, MarketTimeItem{
			Market:           item.Market,
			TradeStatus:      item.TradeStatus,
			Timestamp:        parseTimestampString(item.Timestamp),
			DelayTradeStatus: item.DelayTradeStatus,
			DelayTimestamp:   parseTimestampString(item.DelayTimestamp),
			SubStatus:        item.SubStatus,
			DelaySubStatus:   item.DelaySubStatus,
		})
	}
	return out, nil
}

// BrokerHolding returns the top broker holdings (buy/sell leaders) for a security.
//
// Path: GET /v1/quote/broker-holding
func (m *MarketContext) BrokerHolding(ctx context.Context, symbol string, period BrokerHoldingPeriod) (*BrokerHoldingTop, error) {
	var resp jsontypes.BrokerHoldingTop
	params := url.Values{}
	params.Set("counter_id", symbolToCounterID(symbol))
	params.Set("type", period.toAPIString())
	if err := m.httpClient.Get(ctx, "/v1/quote/broker-holding", params, &resp); err != nil {
		return nil, err
	}
	return convertBrokerHoldingTop(&resp), nil
}

// BrokerHoldingDetail returns the full broker holding details for a security.
//
// Path: GET /v1/quote/broker-holding/detail
func (m *MarketContext) BrokerHoldingDetail(ctx context.Context, symbol string) (*BrokerHoldingDetail, error) {
	var resp jsontypes.BrokerHoldingDetail
	params := url.Values{}
	params.Set("counter_id", symbolToCounterID(symbol))
	if err := m.httpClient.Get(ctx, "/v1/quote/broker-holding/detail", params, &resp); err != nil {
		return nil, err
	}
	return convertBrokerHoldingDetail(&resp), nil
}

// BrokerHoldingDaily returns the daily holding history for a specific broker.
//
// Path: GET /v1/quote/broker-holding/daily
func (m *MarketContext) BrokerHoldingDaily(ctx context.Context, symbol string, brokerID string) (*BrokerHoldingDailyHistory, error) {
	var resp jsontypes.BrokerHoldingDailyHistory
	params := url.Values{}
	params.Set("counter_id", symbolToCounterID(symbol))
	params.Set("parti_number", brokerID)
	if err := m.httpClient.Get(ctx, "/v1/quote/broker-holding/daily", params, &resp); err != nil {
		return nil, err
	}
	out := &BrokerHoldingDailyHistory{
		List: make([]BrokerHoldingDailyItem, 0, len(resp.List)),
	}
	for _, item := range resp.List {
		out.List = append(out.List, BrokerHoldingDailyItem{
			Date:    item.Date,
			Holding: parseOptionalDecimal(item.Holding),
			Ratio:   parseOptionalDecimal(item.Ratio),
			Chg:     parseOptionalDecimal(item.Chg),
		})
	}
	return out, nil
}

// AhPremium returns the A/H premium K-line data for a dual-listed security.
//
// Path: GET /v1/quote/ahpremium/klines
func (m *MarketContext) AhPremium(ctx context.Context, symbol string, period AhPremiumPeriod, count uint32) (*AhPremiumKlines, error) {
	var resp jsontypes.AhPremiumKlines
	params := url.Values{}
	params.Set("counter_id", symbolToCounterID(symbol))
	params.Set("line_type", period.lineType())
	params.Set("line_num", fmt.Sprintf("%d", count))
	if err := m.httpClient.Get(ctx, "/v1/quote/ahpremium/klines", params, &resp); err != nil {
		return nil, err
	}
	out := &AhPremiumKlines{
		Klines: convertAhPremiumKlines(resp.Klines),
	}
	return out, nil
}

// AhPremiumIntraday returns the A/H premium intraday data for a dual-listed security.
//
// Path: GET /v1/quote/ahpremium/timeshares
func (m *MarketContext) AhPremiumIntraday(ctx context.Context, symbol string) (*AhPremiumIntraday, error) {
	var resp jsontypes.AhPremiumIntraday
	params := url.Values{}
	params.Set("counter_id", symbolToCounterID(symbol))
	params.Set("days", "1")
	if err := m.httpClient.Get(ctx, "/v1/quote/ahpremium/timeshares", params, &resp); err != nil {
		return nil, err
	}
	out := &AhPremiumIntraday{
		Klines: convertAhPremiumKlines(resp.Klines),
	}
	return out, nil
}

// TradeStats returns buy/sell/neutral trade statistics for a security.
//
// Path: GET /v1/quote/trades-statistics
func (m *MarketContext) TradeStats(ctx context.Context, symbol string) (*TradeStatsResponse, error) {
	var resp jsontypes.TradeStatsResponse
	params := url.Values{}
	params.Set("counter_id", symbolToCounterID(symbol))
	if err := m.httpClient.Get(ctx, "/v1/quote/trades-statistics", params, &resp); err != nil {
		return nil, err
	}
	out := &TradeStatsResponse{
		Statistics: TradeStatistics{
			Avgprice:    parseDecimal(resp.Statistics.Avgprice),
			Buy:         parseDecimal(resp.Statistics.Buy),
			Neutral:     parseDecimal(resp.Statistics.Neutral),
			Preclose:    parseDecimal(resp.Statistics.Preclose),
			Sell:        parseDecimal(resp.Statistics.Sell),
			Timestamp:   resp.Statistics.Timestamp,
			TotalAmount: parseDecimal(resp.Statistics.TotalAmount),
			TradeDate:   resp.Statistics.TradeDate,
			TradesCount: resp.Statistics.TradesCount,
		},
		Trades: make([]TradePriceLevel, 0, len(resp.Trades)),
	}
	for _, t := range resp.Trades {
		out.Trades = append(out.Trades, TradePriceLevel{
			BuyAmount:     parseDecimal(t.BuyAmount),
			NeutralAmount: parseDecimal(t.NeutralAmount),
			Price:         parseDecimal(t.Price),
			SellAmount:    parseDecimal(t.SellAmount),
		})
	}
	return out, nil
}

// Anomaly returns market anomaly alerts (unusual price/volume events) for a market.
//
// Path: GET /v1/quote/changes
func (m *MarketContext) Anomaly(ctx context.Context, market string) (*AnomalyResponse, error) {
	var resp jsontypes.AnomalyResponse
	params := url.Values{}
	params.Set("market", strings.ToUpper(market))
	params.Set("category", "0")
	if err := m.httpClient.Get(ctx, "/v1/quote/changes", params, &resp); err != nil {
		return nil, err
	}
	out := &AnomalyResponse{
		AllOff:  resp.AllOff,
		Changes: make([]AnomalyItem, 0, len(resp.Changes)),
	}
	for _, item := range resp.Changes {
		out.Changes = append(out.Changes, AnomalyItem{
			Symbol:       counterIDToSymbol(item.CounterID),
			Name:         item.Name,
			AlertName:    item.AlertName,
			AlertTime:    item.AlertTime,
			ChangeValues: item.ChangeValues,
			Emotion:      item.Emotion,
		})
	}
	return out, nil
}

// Constituent returns the constituent stocks for an index.
//
// symbol should be an index symbol such as "HSI.HK".
//
// Path: GET /v1/quote/index-constituents
func (m *MarketContext) Constituent(ctx context.Context, symbol string) (*IndexConstituents, error) {
	var resp jsontypes.IndexConstituents
	params := url.Values{}
	params.Set("counter_id", indexSymbolToCounterID(symbol))
	if err := m.httpClient.Get(ctx, "/v1/quote/index-constituents", params, &resp); err != nil {
		return nil, err
	}
	out := &IndexConstituents{
		FallNum: resp.FallNum,
		FlatNum: resp.FlatNum,
		RiseNum: resp.RiseNum,
		Stocks:  make([]ConstituentStock, 0, len(resp.Stocks)),
	}
	for _, s := range resp.Stocks {
		out.Stocks = append(out.Stocks, ConstituentStock{
			Symbol:            counterIDToSymbol(s.CounterID),
			Name:              s.Name,
			LastDone:          parseOptionalDecimal(s.LastDone),
			PrevClose:         parseOptionalDecimal(s.PrevClose),
			Inflow:            parseOptionalDecimal(s.Inflow),
			Balance:           parseOptionalDecimal(s.Balance),
			Amount:            parseOptionalDecimal(s.Amount),
			TotalShares:       parseOptionalDecimal(s.TotalShares),
			Tags:              s.Tags,
			Intro:             s.Intro,
			Market:            s.Market,
			CirculatingShares: parseOptionalDecimal(s.CirculatingShares),
			Delay:             s.Delay,
			Chg:               parseOptionalDecimal(s.Chg),
			TradeStatus:       s.TradeStatus,
		})
	}
	return out, nil
}

// --- helpers ---

// toAPIString converts a BrokerHoldingPeriod to the API's type parameter value.
func (p BrokerHoldingPeriod) toAPIString() string {
	switch p {
	case BrokerHoldingPeriodRct5:
		return "rct_5"
	case BrokerHoldingPeriodRct20:
		return "rct_20"
	case BrokerHoldingPeriodRct60:
		return "rct_60"
	default:
		return "rct_1"
	}
}

// symbolToCounterID converts a symbol like "700.HK" to a counter ID like "700_HK".
// This mirrors the Rust symbol_to_counter_id utility.
func symbolToCounterID(symbol string) string {
	return strings.Replace(symbol, ".", "_", 1)
}

// indexSymbolToCounterID converts an index symbol like "HSI.HK" to a counter ID like "IX_HSI_HK".
// This mirrors the Rust index_symbol_to_counter_id utility.
func indexSymbolToCounterID(symbol string) string {
	parts := strings.SplitN(symbol, ".", 2)
	if len(parts) == 2 {
		return fmt.Sprintf("IX_%s_%s", parts[0], parts[1])
	}
	return symbol
}

// counterIDToSymbol converts a counter ID like "700_HK" back to a symbol "700.HK".
func counterIDToSymbol(counterID string) string {
	// Find the last underscore that separates market suffix
	idx := strings.LastIndex(counterID, "_")
	if idx > 0 {
		return counterID[:idx] + "." + counterID[idx+1:]
	}
	return counterID
}

// parseOptionalDecimal parses a decimal string, returning nil if the string is empty or zero-ish.
func parseOptionalDecimal(s string) *decimal.Decimal {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return &d
}

// parseDecimal parses a decimal string, returning zero on error.
func parseDecimal(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// parseTimestampString parses a unix-second timestamp string to time.Time.
func parseTimestampString(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	var sec int64
	for _, c := range s {
		if c < '0' || c > '9' {
			return time.Time{}
		}
		sec = sec*10 + int64(c-'0')
	}
	return time.Unix(sec, 0).UTC()
}

// convertBrokerHoldingEntry converts a raw jsontypes entry to a public entry.
func convertBrokerHoldingEntry(item jsontypes.BrokerHoldingEntry) BrokerHoldingEntry {
	return BrokerHoldingEntry{
		Name:        item.Name,
		PartiNumber: item.PartiNumber,
		Chg:         parseOptionalDecimal(item.Chg),
		Strong:      item.Strong,
	}
}

// convertBrokerHoldingChanges converts raw jsontypes changes to public changes.
func convertBrokerHoldingChanges(c jsontypes.BrokerHoldingChanges) BrokerHoldingChanges {
	return BrokerHoldingChanges{
		Value: parseOptionalDecimal(c.Value),
		Chg1:  parseOptionalDecimal(c.Chg1),
		Chg5:  parseOptionalDecimal(c.Chg5),
		Chg20: parseOptionalDecimal(c.Chg20),
		Chg60: parseOptionalDecimal(c.Chg60),
	}
}

// convertBrokerHoldingTop converts a raw jsontypes response to the public type.
func convertBrokerHoldingTop(resp *jsontypes.BrokerHoldingTop) *BrokerHoldingTop {
	out := &BrokerHoldingTop{
		UpdatedAt: resp.UpdatedAt,
		Buy:       make([]BrokerHoldingEntry, 0, len(resp.Buy)),
		Sell:      make([]BrokerHoldingEntry, 0, len(resp.Sell)),
	}
	for _, item := range resp.Buy {
		out.Buy = append(out.Buy, convertBrokerHoldingEntry(item))
	}
	for _, item := range resp.Sell {
		out.Sell = append(out.Sell, convertBrokerHoldingEntry(item))
	}
	return out
}

// convertBrokerHoldingDetail converts a raw jsontypes response to the public type.
func convertBrokerHoldingDetail(resp *jsontypes.BrokerHoldingDetail) *BrokerHoldingDetail {
	out := &BrokerHoldingDetail{
		UpdatedAt: resp.UpdatedAt,
		List:      make([]BrokerHoldingDetailItem, 0, len(resp.List)),
	}
	for _, item := range resp.List {
		out.List = append(out.List, BrokerHoldingDetailItem{
			Name:        item.Name,
			PartiNumber: item.PartiNumber,
			Ratio:       convertBrokerHoldingChanges(item.Ratio),
			Shares:      convertBrokerHoldingChanges(item.Shares),
			Strong:      item.Strong,
		})
	}
	return out
}

// convertAhPremiumKlines converts raw jsontypes kline slices to public kline slices.
func convertAhPremiumKlines(raw []jsontypes.AhPremiumKline) []AhPremiumKline {
	out := make([]AhPremiumKline, 0, len(raw))
	for _, k := range raw {
		out = append(out, AhPremiumKline{
			Aprice:        parseDecimal(k.Aprice),
			Apreclose:     parseDecimal(k.Apreclose),
			Hprice:        parseDecimal(k.Hprice),
			Hpreclose:     parseDecimal(k.Hpreclose),
			CurrencyRate:  parseDecimal(k.CurrencyRate),
			AhpremiumRate: parseDecimal(k.AhpremiumRate),
			PriceSpread:   parseDecimal(k.PriceSpread),
			Timestamp:     time.Unix(k.Timestamp, 0).UTC(),
		})
	}
	return out
}
