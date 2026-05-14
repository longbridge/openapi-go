package market

import (
	"context"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/market/jsontypes"
)

// MarketContext is a client for market data (broker holdings, A/H premium, etc.).
type MarketContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a MarketContext from a Config.
func NewFromCfg(cfg *config.Config) (*MarketContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &MarketContext{httpClient: httpClient}, nil
}

// NewFromEnv creates a MarketContext from environment variables.
func NewFromEnv() (*MarketContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

func parseDecimalOpt(s string) *decimal.Decimal {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return &d
}

// MarketStatus returns the current trading status for all markets.
func (c *MarketContext) MarketStatus(ctx context.Context) (*MarketStatusResponse, error) {
	var resp jsontypes.MarketStatusResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/market-status", nil, &resp); err != nil {
		return nil, err
	}
	result := &MarketStatusResponse{}
	for _, item := range resp.MarketTime {
		result.MarketTime = append(result.MarketTime, &MarketTimeItem{
			Market:           item.Market,
			TradeStatus:      item.TradeStatus,
			Timestamp:        item.Timestamp,
			DelayTradeStatus: item.DelayTradeStatus,
			DelayTimestamp:   item.DelayTimestamp,
			SubStatus:        item.SubStatus,
			DelaySubStatus:   item.DelaySubStatus,
		})
	}
	return result, nil
}

// BrokerHolding returns the top broker holdings for a security.
func (c *MarketContext) BrokerHolding(ctx context.Context, symbol string, period BrokerHoldingPeriod) (*BrokerHoldingTop, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	values.Add("type", string(period))
	var resp jsontypes.BrokerHoldingTop
	if err := c.httpClient.Get(ctx, "/v1/quote/broker-holding", values, &resp); err != nil {
		return nil, err
	}
	result := &BrokerHoldingTop{UpdatedAt: resp.UpdatedAt}
	for _, e := range resp.Buy {
		result.Buy = append(result.Buy, &BrokerHoldingEntry{
			Name: e.Name, PartiNumber: e.PartiNumber,
			Chg: parseDecimalOpt(e.Chg), Strong: e.Strong,
		})
	}
	for _, e := range resp.Sell {
		result.Sell = append(result.Sell, &BrokerHoldingEntry{
			Name: e.Name, PartiNumber: e.PartiNumber,
			Chg: parseDecimalOpt(e.Chg), Strong: e.Strong,
		})
	}
	return result, nil
}

// BrokerHoldingDetail returns full broker holding details for a security.
func (c *MarketContext) BrokerHoldingDetail(ctx context.Context, symbol string) (*BrokerHoldingDetail, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.BrokerHoldingDetail
	if err := c.httpClient.Get(ctx, "/v1/quote/broker-holding/detail", values, &resp); err != nil {
		return nil, err
	}
	result := &BrokerHoldingDetail{UpdatedAt: resp.UpdatedAt}
	for _, item := range resp.List {
		di := &BrokerHoldingDetailItem{
			Name: item.Name, PartiNumber: item.PartiNumber, Strong: item.Strong,
		}
		if item.Ratio != nil {
			di.Ratio = convertHoldingChanges(item.Ratio)
		}
		if item.Shares != nil {
			di.Shares = convertHoldingChanges(item.Shares)
		}
		result.List = append(result.List, di)
	}
	return result, nil
}

func convertHoldingChanges(c *jsontypes.BrokerHoldingChanges) *BrokerHoldingChanges {
	return &BrokerHoldingChanges{
		Value: parseDecimalOpt(c.Value),
		Chg1:  parseDecimalOpt(c.Chg1),
		Chg5:  parseDecimalOpt(c.Chg5),
		Chg20: parseDecimalOpt(c.Chg20),
		Chg60: parseDecimalOpt(c.Chg60),
	}
}

// BrokerHoldingDaily returns daily holding history for a specific broker.
func (c *MarketContext) BrokerHoldingDaily(ctx context.Context, symbol, brokerId string) (*BrokerHoldingDailyHistory, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	values.Add("parti_number", brokerId)
	var resp jsontypes.BrokerHoldingDailyHistory
	if err := c.httpClient.Get(ctx, "/v1/quote/broker-holding/daily", values, &resp); err != nil {
		return nil, err
	}
	result := &BrokerHoldingDailyHistory{}
	for _, item := range resp.List {
		result.List = append(result.List, &BrokerHoldingDailyItem{
			Date:    item.Date,
			Holding: parseDecimalOpt(item.Holding),
			Ratio:   parseDecimalOpt(item.Ratio),
			Chg:     parseDecimalOpt(item.Chg),
		})
	}
	return result, nil
}

// AhPremium returns A/H premium K-line data for a dual-listed security.
func (c *MarketContext) AhPremium(ctx context.Context, symbol string, period AhPremiumPeriod, count int32) (*AhPremiumKlines, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	values.Add("line_type", string(period))
	values.Add("line_num", strconv.FormatInt(int64(count), 10))
	var resp jsontypes.AhPremiumKlines
	if err := c.httpClient.Get(ctx, "/v1/quote/ahpremium/klines", values, &resp); err != nil {
		return nil, err
	}
	result := &AhPremiumKlines{}
	for _, k := range resp.Klines {
		result.Klines = append(result.Klines, &AhPremiumKline{
			Aprice:        parseDecimalOpt(k.Aprice),
			Apreclose:     parseDecimalOpt(k.Apreclose),
			Hprice:        parseDecimalOpt(k.Hprice),
			Hpreclose:     parseDecimalOpt(k.Hpreclose),
			CurrencyRate:  parseDecimalOpt(k.CurrencyRate),
			AhpremiumRate: parseDecimalOpt(k.AhpremiumRate),
			PriceSpread:   parseDecimalOpt(k.PriceSpread),
			Timestamp:     k.Timestamp,
		})
	}
	return result, nil
}

// AhPremiumIntraday returns A/H premium intraday data for a dual-listed security.
func (c *MarketContext) AhPremiumIntraday(ctx context.Context, symbol string) (*AhPremiumKlines, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	values.Add("days", "1")
	var resp jsontypes.AhPremiumKlines
	if err := c.httpClient.Get(ctx, "/v1/quote/ahpremium/timeshares", values, &resp); err != nil {
		return nil, err
	}
	result := &AhPremiumKlines{}
	for _, k := range resp.Klines {
		result.Klines = append(result.Klines, &AhPremiumKline{
			Aprice:        parseDecimalOpt(k.Aprice),
			Apreclose:     parseDecimalOpt(k.Apreclose),
			Hprice:        parseDecimalOpt(k.Hprice),
			Hpreclose:     parseDecimalOpt(k.Hpreclose),
			CurrencyRate:  parseDecimalOpt(k.CurrencyRate),
			AhpremiumRate: parseDecimalOpt(k.AhpremiumRate),
			PriceSpread:   parseDecimalOpt(k.PriceSpread),
			Timestamp:     k.Timestamp,
		})
	}
	return result, nil
}

// TradeStats returns buy/sell/neutral trade statistics for a security.
func (c *MarketContext) TradeStats(ctx context.Context, symbol string) (*TradeStatsResponse, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	var resp jsontypes.TradeStatsResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/trades-statistics", values, &resp); err != nil {
		return nil, err
	}
	result := &TradeStatsResponse{}
	if resp.Statistics != nil {
		result.Statistics = &TradeStatistics{
			Avgprice:    parseDecimalOpt(resp.Statistics.Avgprice),
			Buy:         parseDecimalOpt(resp.Statistics.Buy),
			Neutral:     parseDecimalOpt(resp.Statistics.Neutral),
			Preclose:    parseDecimalOpt(resp.Statistics.Preclose),
			Sell:        parseDecimalOpt(resp.Statistics.Sell),
			Timestamp:   resp.Statistics.Timestamp,
			TotalAmount: parseDecimalOpt(resp.Statistics.TotalAmount),
			TradeDate:   resp.Statistics.TradeDate,
			TradesCount: resp.Statistics.TradesCount,
		}
	}
	for _, t := range resp.Trades {
		result.Trades = append(result.Trades, &TradePriceLevel{
			BuyAmount:     parseDecimalOpt(t.BuyAmount),
			NeutralAmount: parseDecimalOpt(t.NeutralAmount),
			Price:         parseDecimalOpt(t.Price),
			SellAmount:    parseDecimalOpt(t.SellAmount),
		})
	}
	return result, nil
}

// Anomaly returns market anomaly alerts (unusual price/volume events) for a market.
func (c *MarketContext) Anomaly(ctx context.Context, market string) (*AnomalyResponse, error) {
	values := url.Values{}
	values.Add("market", market)
	values.Add("category", "0")
	var resp jsontypes.AnomalyResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/changes", values, &resp); err != nil {
		return nil, err
	}
	result := &AnomalyResponse{AllOff: resp.AllOff}
	for _, item := range resp.Changes {
		result.Changes = append(result.Changes, &AnomalyItem{
			Symbol:       util.CounterIDToSymbol(item.CounterId),
			Name:         item.Name,
			AlertName:    item.AlertName,
			AlertTime:    item.AlertTime,
			ChangeValues: item.ChangeValues,
			Emotion:      item.Emotion,
		})
	}
	return result, nil
}

// Constituent returns the constituent stocks for an index symbol (e.g. "HSI.HK").
func (c *MarketContext) Constituent(ctx context.Context, symbol string) (*IndexConstituents, error) {
	values := url.Values{}
	values.Add("counter_id", util.IndexSymbolToCounterID(symbol))
	var resp jsontypes.IndexConstituents
	if err := c.httpClient.Get(ctx, "/v1/quote/index-constituents", values, &resp); err != nil {
		return nil, err
	}
	result := &IndexConstituents{
		FallNum: resp.FallNum,
		FlatNum: resp.FlatNum,
		RiseNum: resp.RiseNum,
	}
	for _, s := range resp.Stocks {
		result.Stocks = append(result.Stocks, &ConstituentStock{
			Symbol:            util.CounterIDToSymbol(s.CounterId),
			Name:              s.Name,
			LastDone:          parseDecimalOpt(s.LastDone),
			PrevClose:         parseDecimalOpt(s.PrevClose),
			Inflow:            parseDecimalOpt(s.Inflow),
			Balance:           parseDecimalOpt(s.Balance),
			Amount:            parseDecimalOpt(s.Amount),
			TotalShares:       parseDecimalOpt(s.TotalShares),
			Tags:              s.Tags,
			Intro:             s.Intro,
			Market:            s.Market,
			CirculatingShares: parseDecimalOpt(s.CirculatingShares),
			Delay:             s.Delay,
			Chg:               parseDecimalOpt(s.Chg),
			TradeStatus:       s.TradeStatus,
		})
	}
	return result, nil
}
