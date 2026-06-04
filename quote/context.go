package quote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go"
	"github.com/longbridge/openapi-go/config"
	counterpkg "github.com/longbridge/openapi-go/counter"
	"github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/counter"
	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/longbridge"
	"github.com/longbridge/openapi-go/quote/jsontypes"
)

// QuoteContext is a client for interacting with Longbridge Quote OpenAPI
// Longbridge Quote OpenAPI document is https://open.longbridge.com/en/docs/quote/overview
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	// handle err
//	// do subscribe
//	err = qctx.Subscribe(context.Background(), []string{"AAPL.US"}, []quote.SubType{quote.SubTypeQuote, quote.SubTypeTrade, quote.SubTypeDepth, quote.SubTypeBrokers}, true)
//	// handle err
//
//	qctx.OnQuote(func (quote *quote.PushQuote) {
//	  // quote callback
//
//	})
//	qctx.OnTrade(func (trade *quote.PushTrade) {
//	  // trade callback
//	})
//	qctx.OnTrade(func (trade *quote.PushDepth) {
//	  // depth callback
//	})
//	qctx.OnBrokers(func (trade *quote.PushBrokers) {
//	  // brokers callback
//	})

type QuoteContext struct {
	opts *Options
	core *core
}

// Profile obtain the user quote profile
func (c *QuoteContext) Profile() *UserProfile {
	return c.core.Profile()
}

// OnQuote set callback function which will be called when server push quote events.
func (c *QuoteContext) OnQuote(f func(*PushQuote)) {
	c.core.SetQuoteHandler(f)
}

// OnTrade set callback function which will be called when server push trade events.
func (c *QuoteContext) OnTrade(f func(*PushTrade)) {
	c.core.SetTradeHandler(f)
}

// OnDepth set callback function which will be called when server push depth events.
func (c *QuoteContext) OnDepth(f func(*PushDepth)) {
	c.core.SetDepthHandler(f)
}

// OnBrokers set callback function which will be called when server push brokers events.
func (c *QuoteContext) OnBrokers(f func(*PushBrokers)) {
	c.core.SetBrokersHandler(f)
}

// Subscribe quote
// Reference: https://open.longbridge.com/en/docs/quote/subscribe/subscribe
func (c *QuoteContext) Subscribe(ctx context.Context, symbols []string, subTypes []SubType, isFirstPush bool) (err error) {
	return c.core.Subscribe(ctx, symbols, subTypes, isFirstPush)
}

// Unsubscribe quote
// Reference: https://open.longbridge.com/en/docs/quote/subscribe/unsubscribe
func (c *QuoteContext) Unsubscribe(ctx context.Context, unSubAll bool, symbols []string, subTypes []SubType) (err error) {
	return c.core.Unsubscribe(ctx, unSubAll, symbols, subTypes)
}

// Subscriptions obtain the subscription information.
// Reference: https://open.longbridge.com/en/docs/quote/subscribe/subscription
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	subs, err := qctx.Subscriptions(context.Background())
func (c *QuoteContext) Subscriptions(ctx context.Context) (subscriptions map[string][]SubType, err error) {
	return c.core.Subscriptions(ctx)
}

// StaticInfo obtain the basic information of securities.
// Reference: https://open.longbridge.com/en/docs/quote/pull/static
//
// Example:
//
//	 qctx, err := quote.NewFromEnv()
//	 infos, err := qctx.StaticInfo(context.Background(), []string{"AAPL.US"})

func (c *QuoteContext) StaticInfo(ctx context.Context, symbols []string) (staticInfos []*StaticInfo, err error) {
	return c.core.StaticInfo(ctx, symbols)
}

// Quote obtain the real-time quotes of securities, and supports all types of securities.
// Reference: https://open.longbridge.com/en/docs/quote/pull/quote
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	quotes, err := qctx.Quote(context.Background(), []string{"AAPL.US"})
func (c *QuoteContext) Quote(ctx context.Context, symbols []string) (quotes []*SecurityQuote, err error) {
	return c.core.Quote(ctx, symbols)
}

// OptionQuote obtain the real-time quotes of US stock options, including the option-specific data.
// Reference: https://open.longbridge.com/en/docs/quote/pull/option-quote
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	optionQuotes, err := qctx.OptionQuote(context.Background(), []string{"AAPL240531P192500.US"})
func (c *QuoteContext) OptionQuote(ctx context.Context, symbols []string) (optionQuotes []*OptionQuote, err error) {
	return c.core.OptionQuote(ctx, symbols)
}

// WarrantQuote obtain the real-time quotes of HK warrants, including the warrant-specific data.
// Reference: https://open.longbridge.com/en/docs/quote/pull/warrant-quote
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	warrantQuotes, err := qctx.WarrantQuote(context.Background(), []string{"9988.HK"})
func (c *QuoteContext) WarrantQuote(ctx context.Context, symbols []string) (warrantQuotes []*WarrantQuote, err error) {
	return c.core.WarrantQuote(ctx, symbols)
}

// Depth obtain the depth data of security.
// Reference: https://open.longbridge.com/en/docs/quote/pull/depth
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	depth, err := qctx.Depth(context.Background(), []string{"AAPL.HK"})
func (c *QuoteContext) Depth(ctx context.Context, symbol string) (securityDepth *SecurityDepth, err error) {
	return c.core.Depth(ctx, symbol)
}

// Brokers obtain the real-time broker queue data of security.
// Reference: https://open.longbridge.com/en/docs/quote/pull/brokers
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	brokers, err := qctx.Brokers(context.Background(), "AAPL.HK")
func (c *QuoteContext) Brokers(ctx context.Context, symbol string) (securityBrokers *SecurityBrokers, err error) {
	return c.core.Brokers(ctx, symbol)
}

// Participants obtain participant IDs data (which can be synchronized once a day).
// Reference: https://open.longbridge.com/en/docs/quote/pull/broker-ids
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	participants, err := qctx.Participants(context.Background())
func (c *QuoteContext) Participants(ctx context.Context) (infos []*ParticipantInfo, err error) {
	return c.core.Participants(ctx)
}

// Trades obtain the trades data of security.
// Reference: https://open.longbridge.com/en/docs/quote/pull/trade
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	trades, err := qctx.Trades(context.Background())
func (c *QuoteContext) Trades(ctx context.Context, symbol string, count int32) (trades []*Trade, err error) {
	return c.core.Trades(ctx, symbol, count)
}

// Intraday obtain the intraday data of security.
// Reference: https://open.longbridge.com/en/docs/quote/pull/intraday
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	trades, err := qctx.Intraday(context.Background(), "AAPL.US")
func (c *QuoteContext) Intraday(ctx context.Context, symbol string) (lines []*IntradayLine, err error) {
	return c.core.Intraday(ctx, symbol)
}

// Candlesticks obtain the candlestick data of security.
// Reference: https://open.longbridge.com/en/docs/quote/pull/candlestick
// periond support values:
//   - quote.PeriodOneMinute - 1m
//   - quote.PeriodFiveMinute - 5m
//   - quote.PeriodFifteenMinute - 15m
//   - quote.PeriodThirtyMinute - 40m
//   - quote.PeriodDay - 1 day
//   - quote.PeriodWeek - 1 week
//   - quote.PeriodMonth - 1 month
//   - quote.PeriodYear - 1 year
//
// adjustType support values:
//   - quote.AdjustTypeNo
//   - quote.AdjustTypeForward
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	klines, err := qctx.Candlesticks(context.Background(), "AAPL.US", quote.PeriodDay, 10, quote.AdjustTypeNo)
func (c *QuoteContext) Candlesticks(ctx context.Context, symbol string, period Period, count int32, adjustType AdjustType) (sticks []*Candlestick, err error) {
	return c.core.Candlesticks(ctx, symbol, period, count, adjustType)
}

// HistoryCandlesticksByDate obtains the history candlestick data of security within a date range.
// Reference: https://open.longbridge.com/en/docs/quote/pull/history-candlestick
// periond support values:
//   - quote.PeriodOneMinute - 1m
//   - quote.PeriodFiveMinute - 5m
//   - quote.PeriodFifteenMinute - 15m
//   - quote.PeriodThirtyMinute - 40m
//   - quote.PeriodDay - 1 day
//   - quote.PeriodWeek - 1 week
//   - quote.PeriodMonth - 1 month
//   - quote.PeriodYear - 1 year
//
// adjustType support values:
//   - quote.AdjustTypeNo
//   - quote.AdjustTypeForward
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	dateTime := time.Date(2022, 5, 10, 11, 10, 0, 0, time.UTC)
//	klines, err := qctx.HistoryCandlesticksByOffset(context.Background(), "AAPL.US", quote.PeriodDay, quote.AdjustTypeNo, true, &dateTime, 100)
func (c *QuoteContext) HistoryCandlesticksByOffset(ctx context.Context, symbol string, period Period, adjustType AdjustType, isForward bool, dateTime *time.Time, count int32, opts ...CandlestickRequestOption) (sticks []*Candlestick, err error) {
	return c.core.HistoryCandlesticksByOffset(ctx, symbol, period, adjustType, isForward, dateTime, count, opts...)
}

// HistoryCandlesticksByOffset obtains the history candlestick data of security after or before an offset time.
// Reference: https://open.longbridge.com/en/docs/quote/pull/history-candlestick
// periond support values:
//   - quote.PeriodOneMinute - 1m
//   - quote.PeriodFiveMinute - 5m
//   - quote.PeriodFifteenMinute - 15m
//   - quote.PeriodThirtyMinute - 40m
//   - quote.PeriodDay - 1 day
//   - quote.PeriodWeek - 1 week
//   - quote.PeriodMonth - 1 month
//   - quote.PeriodYear - 1 year
//
// adjustType support values:
//   - quote.AdjustTypeNo
//   - quote.AdjustTypeForward
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	startDate := time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)
//	endDate := time.Date(2022, 6, 10, 0, 0, 0, 0, time.UTC)
//	klines, err := qctx.HistoryCandlesticksByDate(context.Background(), "AAPL.US", quote.PeriodDay, quote.AdjustTypeNo, &startDate, &endDate)
func (c *QuoteContext) HistoryCandlesticksByDate(ctx context.Context, symbol string, period Period, adjustType AdjustType, startDate *time.Time, endDate *time.Time, opts ...CandlestickRequestOption) (sticks []*Candlestick, err error) {
	return c.core.HistoryCandlesticksByDate(ctx, symbol, period, adjustType, startDate, endDate, opts...)
}

// OptionChainExpiryDateList obtain the the list of expiration dates of option chain
// Reference: https://open.longbridge.com/en/docs/quote/pull/optionchain-date
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	list, err := qctx.OptionChainExpiryDateList(context.Background(), "AAPL.US")
func (c *QuoteContext) OptionChainExpiryDateList(ctx context.Context, symbol string) (times []time.Time, err error) {
	return c.core.OptionChainExpiryDateList(ctx, symbol)
}

// OptionChainInfoByDate obtain a list of option securities by the option chain expiry date.
// Reference: https://open.longbridge.com/en/docs/quote/pull/optionchain-date-strike
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	date := time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)
//	list, err := qctx.OptionChainInfoByDate(context.Background(), "AAPL.US", &date)
func (c *QuoteContext) OptionChainInfoByDate(ctx context.Context, symbol string, expiryDate *time.Time) (strikePriceInfos []*StrikePriceInfo, err error) {
	return c.core.OptionChainInfoByDate(ctx, symbol, expiryDate)
}

// WarrantIssuers obtain the warrant issuer IDs data (which can be synchronized once a day).
// Reference: https://open.longbridge.com/en/docs/quote/pull/issuer
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	infos, err := qctx.WarrantIssuers(context.Background())
func (c *QuoteContext) WarrantIssuers(ctx context.Context) (infos []*IssuerInfo, err error) {
	return c.core.WarrantIssuers(ctx)
}

// WarrantIssuers obtain the quotes of HK warrants, and supports sorting and filtering.
// Reference: https://open.longbridge.com/en/docs/quote/pull/warrant-filter
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	list, err := qctx.WarrantList(context.Background(), "9988.HK", quote.WarrantFilter{
//	  SortBy: quote.WarrantLastDone,
//	  SortOrder quote.WarrantAsc,
//	  SortOffset: 0,
//	  SortCount: 20,
//	  Type: []quote.WarrantType{quote.WarrantCall, quote.WarrantPut},
//	}, quote. WarrantZH_CN)
func (c *QuoteContext) WarrantList(ctx context.Context, symbol string, config WarrantFilter, lang WarrantLanguage) (infos []*WarrantInfo, err error) {
	return c.core.WarrantList(ctx, symbol, config, lang)
}

// TradingSession obtain the daily trading hours of each market.
// Reference: https://open.longbridge.com/en/docs/quote/pull/trade-session
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	sess, err := qctx.TradingSession(context.Background())
func (c *QuoteContext) TradingSession(ctx context.Context) (sessions []*MarketTradingSession, err error) {
	return c.core.TradingSession(ctx)
}

// TradingDays obtain the trading days of the market.
// Reference: https://open.longbridge.com/en/docs/quote/pull/trade-day
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	begin := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
//	end := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
//	days, err := qctx.TradingDays(context.Background(), openapi.MarketUS, &begin, &end)
func (c *QuoteContext) TradingDays(ctx context.Context, market openapi.Market, begin *time.Time, end *time.Time) (days *MarketTradingDay, err error) {
	return c.core.TradingDays(ctx, market, begin, end)
}

// CapitalDistribution is used to obtain the daily capital distribution of security.
// Reference: https://open.longbridge.com/en/docs/quote/pull/capital-distribution
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	dist, err := qctx.CapitalDistribution(context.Background(), "700.HK")
func (c *QuoteContext) CapitalDistribution(ctx context.Context, symbol string) (capitalDib CapitalDistribution, err error) {
	return c.core.CapitalDistribution(ctx, symbol)
}

// CapitalFlow is used to obtain the daily capital flow intraday of security.
// Reference: https://open.longbridge.com/en/docs/quote/pull/capital-flow-intraday
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	flowlines, err := qctx.CapitalFlow(context.Background(), "700.HK")
func (c *QuoteContext) CapitalFlow(ctx context.Context, symbol string) (capitalFlowLines []CapitalFlowLine, err error) {
	return c.core.CapitalFlow(ctx, symbol)
}

// CalcIndex is used to obtain the calculate indexes of securities.
// Reference: https://open.longbridge.com/en/docs/quote/pull/calc-index
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	flowlines, err := qctx.CapitalFlow(context.Background(), "700.HK")
func (c *QuoteContext) CalcIndex(ctx context.Context, symbols []string, indexes []CalcIndex) (calcIndexes []*SecurityCalcIndex, err error) {
	return c.core.CalcIndex(ctx, symbols, indexes)
}

// RealtimeQuote to get quote infomations on local store
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	flowlines, err := qctx.RealtimeQuote(context.Background(), []string{"700.HK", "9988.HK"})
func (c *QuoteContext) RealtimeQuote(ctx context.Context, symbols []string) ([]*Quote, error) {
	return c.core.RealtimeQuote(ctx, symbols)
}

// RealtimeDepth to get depth infomations on local store
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	depth, err := qctx.RealtimeDepth(context.Background(), "700.HK")
func (c *QuoteContext) RealtimeDepth(ctx context.Context, symbol string) (*SecurityDepth, error) {
	return c.core.RealtimeDepth(ctx, symbol)
}

// RealtimeTrades to get trade infomations on local store
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	trades, err := qctx.RealtimeTrades(context.Background(), "700.HK")
func (c *QuoteContext) RealtimeTrades(ctx context.Context, symbol string) ([]*Trade, error) {
	return c.core.RealtimeTrades(ctx, symbol)
}

// RealtimeBrokers to get broker infomations on local store
//
// Example:
//
//	qctx, err := quote.NewFromEnv()
//	brokers, err := qctx.RealtimeBrokers(context.Background(), "700.HK")
func (c *QuoteContext) RealtimeBrokers(ctx context.Context, symbol string) (*SecurityBrokers, error) {
	return c.core.RealtimeBrokers(ctx, symbol)
}

// CreateWatchlistGroup use to create watchlist group. Doc: https://open.longbridge.com/en/docs/quote/individual/watchlist_create_group
//
// Example:
//
//	qctx, err := quote.NewFromCfg(conf)
//	// handle error
//	err = qctx.CreateWatchlistGroup(context.Background(), "test", []string{"AAPL.US"})
//	qctx, err := quote.NewFromCfg(conf)
//	// ignore handle error
//	err = qctx.CreateWatchlistGroup(context.Background(), "test", []string{"AAPL.US"})
func (c *QuoteContext) CreateWatchlistGroup(ctx context.Context, name string, symbols []string) (gid int64, err error) {
	var resp struct {
		ID int64 `json:"id,string"`
	}
	err = c.opts.httpClient.Post(ctx, "/v1/watchlist/groups", map[string]interface{}{
		"name":       name,
		"securities": symbols,
	}, &resp)

	gid = resp.ID
	return
}

// DeleteWatchlistGroup use to delete watchlist group.
// If `purge` set to true, the securities in the group will be unfollowed even in the *other* groups.
// If `purge` set to false, the securities in the group will remain in the *other* groups.
// Doc: https://open.longbridge.com/en/docs/quote/individual/watchlist_delete_group
//
// Example:
//
//  qctx, err := quote.NewFromCfg(conf)
//  // ignore handle error
//  err = qctx.DeleteWatchlistGroup(context.Background(), 1, false)

func (c *QuoteContext) DeleteWatchlistGroup(ctx context.Context, id int64, purge bool) (err error) {
	var resp struct{}
	err = c.opts.httpClient.Delete(ctx, "/v1/watchlist/groups", nil, &resp, http.WithBody(map[string]interface{}{
		"id":    id,
		"purge": purge,
	}))
	return
}

// UpdateWatchlistGroup use to update watchlist group.
// Doc: https://open.longbridge.com/en/docs/quote/individual/watchlist_update_group
//
// Example:
//
//  qctx, err := quote.NewFromCfg(conf)
//  // ignore handle error
//  err = qctx.UpdateWatchlistGroup(context.Background(), 1, "test", []string{"AAPL.US"}, quote.AddWatchlist)

func (c *QuoteContext) UpdateWatchlistGroup(ctx context.Context, id int64, name string, symbols []string, mode WatchlistUpdateMode) (err error) {
	var resp struct{}
	err = c.opts.httpClient.Put(ctx, "/v1/watchlist/groups", map[string]interface{}{
		"id":         id,
		"name":       name,
		"securities": symbols,
		"mode":       mode,
	}, &resp)
	return
}

// WatchedGroups to get all watchlist groups.
// Reference: https://open.longbridge.com/en/docs/quote/individual/watchlist_groups
//
// Example:
//
//  qctx, err := quote.NewFromCfg(conf)
//  // ignore handle error
//  err = qctx.WatchedGroups(context.Background())

func (c *QuoteContext) WatchedGroups(ctx context.Context) (groupList []*WatchedGroup, err error) {
	var resp jsontypes.WatchedGroupList
	err = c.opts.httpClient.Get(ctx, "/v1/watchlist/groups", nil, &resp)
	if err != nil {
		return
	}

	err = util.Copy(&groupList, resp.Groups)
	return
}

// Filings returns the filings list for a symbol.
func (c *QuoteContext) Filings(ctx context.Context, symbol string) (items []*FilingItem, err error) {
	var resp jsontypes.FilingList
	values := url.Values{}
	values.Add("symbol", symbol)
	err = c.opts.httpClient.Get(ctx, "/v1/quote/filings", values, &resp)
	if err != nil {
		return
	}
	items = make([]*FilingItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &FilingItem{
			Id:          item.Id,
			Title:       item.Title,
			Description: item.Description,
			FileName:    item.FileName,
			FileUrls:    item.FileUrls,
			PublishAt:   time.Unix(item.PublishAt, 0).UTC(),
		})
	}
	return
}

// ShortPositions returns short interest / short position data for a US or HK security.
//
// Market is auto-detected from the symbol suffix:
//   - ".HK" → GET /v1/quote/short-positions/hk (HKEX daily data)
//   - otherwise → GET /v1/quote/short-positions/us (FINRA bi-monthly data)
//
// count controls the number of records returned (1–100, default 20).
func (c *QuoteContext) ShortPositions(ctx context.Context, symbol string, count uint32) (*ShortPositionsResponse, error) {
	isHK := strings.HasSuffix(strings.ToUpper(symbol), ".HK")
	path := "/v1/quote/short-positions/us"
	if isHK {
		path = "/v1/quote/short-positions/hk"
	}
	values := url.Values{}
	values.Set("counter_id", quoteSymbolToCounterID(symbol))
	values.Set("last_timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	values.Set("count", fmt.Sprintf("%d", count))
	// Response: {"counter_id": "ST/US/AAPL", "data": [{...}]}
	var outer struct {
		CounterID string                       `json:"counter_id"`
		Data      []map[string]json.RawMessage `json:"data"`
	}
	if err := c.opts.httpClient.Get(ctx, path, values, &outer); err != nil {
		return nil, err
	}
	items := make([]*ShortPositionsItem, 0, len(outer.Data))
	for _, r := range outer.Data {
		items = append(items, &ShortPositionsItem{
			Timestamp:           unixSecsToRFC3339(rawStr(r, "timestamp")),
			Rate:                rawStr(r, "rate"),
			Close:               rawStr(r, "close"),
			CurrentSharesShort:  rawStr(r, "current_shares_short"),
			AvgDailyShareVolume: rawStr(r, "avg_daily_share_volume"),
			DaysToCover:         rawStr(r, "days_to_cover"),
			Amount:              rawStr(r, "amount"),
			Balance:             rawStr(r, "balance"),
			Cost:                rawStr(r, "cost"),
		})
	}
	return &ShortPositionsResponse{Data: items}, nil
}

// OptionVolume returns aggregated call/put volume stats for a security.
// Path: GET /v1/quote/option-volume-stats
func (c *QuoteContext) OptionVolume(ctx context.Context, symbol string) (*OptionVolumeStats, error) {
	var resp jsontypes.OptionVolumeStats
	values := url.Values{}
	values.Add("symbol", symbol)
	if err := c.opts.httpClient.Get(ctx, "/v1/quote/option-volume-stats", values, &resp); err != nil {
		return nil, err
	}
	return &OptionVolumeStats{
		CallVolume: resp.CallVolume,
		PutVolume:  resp.PutVolume,
	}, nil
}

// OptionVolumeDaily returns daily option volume data for a security within a time range.
// Path: GET /v1/quote/option-volume-stats/daily
func (c *QuoteContext) OptionVolumeDaily(ctx context.Context, symbol string, start time.Time, end time.Time) ([]*DailyOptionVolume, error) {
	var resp jsontypes.OptionVolumeDaily
	values := url.Values{}
	values.Add("symbol", symbol)
	values.Add("start", start.Format("2006-01-02"))
	values.Add("end", end.Format("2006-01-02"))
	if err := c.opts.httpClient.Get(ctx, "/v1/quote/option-volume-stats/daily", values, &resp); err != nil {
		return nil, err
	}
	result := make([]*DailyOptionVolume, 0, len(resp.Stats))
	for _, s := range resp.Stats {
		result = append(result, &DailyOptionVolume{
			Symbol:                   counter.IDToSymbol(s.Symbol),
			Timestamp:                s.Timestamp,
			TotalVolume:              s.TotalVolume,
			TotalPutVolume:           s.TotalPutVolume,
			TotalCallVolume:          s.TotalCallVolume,
			PutCallVolumeRatio:       s.PutCallVolumeRatio,
			TotalOpenInterest:        s.TotalOpenInterest,
			TotalPutOpenInterest:     s.TotalPutOpenInterest,
			TotalCallOpenInterest:    s.TotalCallOpenInterest,
			PutCallOpenInterestRatio: s.PutCallOpenInterestRatio,
		})
	}
	return result, nil
}

// UpdatePinned pins or unpins securities in the watchlist.
// Path: POST /v1/watchlist/pinned
func (c *QuoteContext) UpdatePinned(ctx context.Context, mode PinnedMode, symbols []string) error {
	var resp struct{}
	return c.opts.httpClient.Post(ctx, "/v1/watchlist/pinned", map[string]interface{}{
		"mode":       mode.String(),
		"securities": symbols,
	}, &resp)
}

// SecurityList used to list securities. Doc: https://open.longbridge.com/en/docs/quote/security/security_list
func (c *QuoteContext) SecurityList(ctx context.Context, market openapi.Market, category SecurityListCategory) (list []*Security, err error) {
	var resp jsontypes.SecurityList
	values := url.Values{}
	values.Add("market", string(market))
	values.Add("category", string(category))

	err = c.opts.httpClient.Get(ctx, "/v1/quote/get_security_list", values, &resp)
	if err != nil {
		return
	}
	err = util.Copy(&list, resp.List)
	return
}

// Close
func (c *QuoteContext) Close() error {
	return c.core.Close()
}

// Deprecated: NewFromEnv return QuoteContext, use NewFromCfg plz
func NewFormEnv() (*QuoteContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, err
	}
	return NewFromCfg(cfg)
}

// NewFromCfg return QuoteContext with config.Config
func NewFromCfg(cfg *config.Config) (*QuoteContext, error) {
	httpClient, err := http.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	lbOpts := longbridge.NewOptions(
		longbridge.WithAuthTimeout(cfg.AuthTimeout),
		longbridge.WithTimeout(cfg.Timeout),
		longbridge.WithReadBufferSize(cfg.ReadBufferSize),
		longbridge.WithReadQueueSize(cfg.ReadQueueSize),
		longbridge.WithWriteQueueSize(cfg.WriteQueueSize),
		longbridge.WithMinGzipSize(cfg.MinGzipSize),
	)
	return New(
		WithQuoteURL(cfg.QuoteUrl),
		WithHttpClient(httpClient),
		WithLogLevel(cfg.LogLevel),
		WithLogger(cfg.Logger()),
		WithLbOptions(lbOpts),
		WithEnableOvernight(cfg.EnableOvernight),
		WithLanguage(cfg.Language),
	)
}

// New return QuoteContext with option.
// A connection will be created with quote server.
func New(opt ...Option) (*QuoteContext, error) {
	opts := newOptions(opt...)
	core, err := newCore(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create core")
	}
	tc := &QuoteContext{
		opts: opts,
		core: core,
	}
	return tc, nil
}


// ShortTrades returns short trade records for a HK or US security.
//
// The endpoint is automatically chosen based on the symbol suffix:
//   - ".HK" → GET /v1/quote/short-trades/hk
//   - ".US" → GET /v1/quote/short-trades/us
func (c *QuoteContext) ShortTrades(ctx context.Context, symbol string, count uint32) (*ShortTradesResponse, error) {
	values := url.Values{}
	values.Set("counter_id", quoteSymbolToCounterID(symbol))
	values.Set("last_timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	values.Set("page_size", fmt.Sprintf("%d", count))

	path := "/v1/quote/short-trades/hk"
	if strings.HasSuffix(strings.ToUpper(symbol), ".US") {
		path = "/v1/quote/short-trades/us"
	}
	// Response: {"counter_id": "ST/HK/700", "data": [{...}]}
	var outer struct {
		CounterID string                       `json:"counter_id"`
		Data      []map[string]json.RawMessage `json:"data"`
	}
	if err := c.opts.httpClient.Get(ctx, path, values, &outer); err != nil {
		return nil, err
	}
	items := make([]*ShortTradesItem, 0, len(outer.Data))
	for _, r := range outer.Data {
		items = append(items, &ShortTradesItem{
			Timestamp:   unixSecsToRFC3339(rawStr(r, "timestamp")),
			Rate:        rawStr(r, "rate"),
			Close:       rawStr(r, "close"),
			NusAmount:   rawStr(r, "nus_amount"),
			NyAmount:    rawStr(r, "ny_amount"),
			TotalAmount: rawStr(r, "total_amount"),
			Amount:      rawStr(r, "amount"),
			Balance:     rawStr(r, "balance"),
		})
	}
	return &ShortTradesResponse{Data: items}, nil
}

// SymbolToCounterIds batch-converts symbols to counter IDs via the remote API.
//
// It returns a map of symbol → counter_id (e.g. "DRAM.US" → "ETF/US/DRAM").
// Symbols the backend does not recognize are omitted from the result.
//
// Path: POST /v1/quote/symbol-to-counter-ids
func (c *QuoteContext) SymbolToCounterIds(ctx context.Context, symbols []string) (map[string]string, error) {
	body := map[string]interface{}{
		"ticker_regions": symbols,
	}
	var resp struct {
		List map[string]string `json:"list"`
	}
	if err := c.opts.httpClient.Post(ctx, "/v1/quote/symbol-to-counter-ids", body, &resp); err != nil {
		return nil, err
	}
	if resp.List == nil {
		resp.List = map[string]string{}
	}
	return resp.List, nil
}

// ResolveCounterIds resolves counter IDs for symbols, local-first with remote
// fallback.
//
// Symbols found in the embedded ETF / index / warrant directory (or in the
// local cache of previous remote resolutions) are resolved without network
// access. The remaining symbols are resolved in one batch via
// SymbolToCounterIds and the results are persisted to the local cache for
// subsequent lookups. Symbols the backend does not recognize fall back to the
// default "ST/" conversion (and are not cached).
func (c *QuoteContext) ResolveCounterIds(ctx context.Context, symbols []string) (map[string]string, error) {
	result := make(map[string]string, len(symbols))
	var unknown []string
	for _, symbol := range symbols {
		if counterID, ok := counterpkg.LookupCounterID(symbol); ok {
			result[symbol] = counterID
		} else {
			unknown = append(unknown, symbol)
		}
	}
	if len(unknown) > 0 {
		resolved, err := c.SymbolToCounterIds(ctx, unknown)
		if err != nil {
			return nil, err
		}
		cacheIDs := make([]string, 0, len(resolved))
		for _, id := range resolved {
			cacheIDs = append(cacheIDs, id)
		}
		counterpkg.CacheCounterIDs(cacheIDs)
		for _, symbol := range unknown {
			if counterID, ok := resolved[symbol]; ok {
				result[symbol] = counterID
			} else {
				result[symbol] = counterpkg.SymbolToCounterID(symbol)
			}
		}
	}
	return result, nil
}

// quoteSymbolToCounterID converts "AAPL.US" → "ST/US/AAPL" for endpoints
// that require the internal counter_id format.
func quoteSymbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	return fmt.Sprintf("ST/%s/%s", strings.ToUpper(symbol[idx+1:]), symbol[:idx])
}

// rawStr extracts a string value from a map of raw JSON values.
// If the value is a JSON string it is unquoted; otherwise the raw bytes are
// returned with surrounding quotes stripped.
func rawStr(m map[string]json.RawMessage, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	var s string
	if err := json.Unmarshal(v, &s); err == nil {
		return s
	}
	return strings.Trim(string(v), `"`)
}

// unixSecsToRFC3339 converts a Unix-seconds string to an RFC 3339 timestamp.
// If the string cannot be parsed it is returned unchanged.
func unixSecsToRFC3339(s string) string {
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return s
	}
	return time.Unix(ts, 0).UTC().Format(time.RFC3339)
}
