// Compare every quote pull API over the WebSocket vs HTTP transports.
//
// Reads credentials from LONGBRIDGE_* env (or .env). Point LONGBRIDGE_HTTP_URL /
// LONGBRIDGE_QUOTE_URL at an environment that serves the /quote/* routes. For
// each call prints the arguments, both results (JSON), latency and whether
// WS == HTTP.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/longbridge/openapi-go"
	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/quote"
)

var errCode = regexp.MustCompile(`\b(\d{3,6})\b`)

func dump(v interface{}, err error) string {
	if err != nil {
		code := "?"
		if m := errCode.FindStringSubmatch(err.Error()); m != nil {
			code = m[1]
		}
		return fmt.Sprintf("Err(code=%s, message=%q)", code, err.Error())
	}
	b, e := json.MarshalIndent(v, "", "  ")
	if e != nil {
		return "marshal error: " + e.Error()
	}
	return "Ok(\n" + string(b) + "\n)"
}

func compare(label, args string, call func(ctx *quote.QuoteContext) (interface{}, error), ws, http *quote.QuoteContext) {
	t := time.Now()
	a := dump(call(ws))
	wsMs := time.Since(t).Milliseconds()
	t = time.Now()
	b := dump(call(http))
	httpMs := time.Since(t).Milliseconds()
	fmt.Printf("=== %s\nREQ %s\nWS_MS %d\nHTTP_MS %d\nEQUAL %v\nWS<<\n%s\n>>WS\nHTTP<<\n%s\n>>HTTP\n", label, args, wsMs, httpMs, a == b, a, b)
}

func main() {
	wsCfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	wsCfg.QuoteTransport = "ws"
	ws, err := quote.NewFromCfg(wsCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	httpCfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	httpCfg.QuoteTransport = "http"
	http, err := quote.NewFromCfg(httpCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer http.Close()

	ctx := context.Background()
	all := []quote.CandlestickRequestOption{quote.CandlestickRequestTradeSession(quote.CandlestickTradeSessionAll)}
	date := func(y int, m time.Month, d int) *time.Time {
		t := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		return &t
	}

	// Resolve an option and a warrant symbol to use below.
	expiries, err := ws.OptionChainExpiryDateList(ctx, "AAPL.US")
	if err != nil || len(expiries) == 0 {
		log.Fatalf("option expiries: %v", err)
	}
	d := expiries[0]
	optionSymbol := fmt.Sprintf("AAPL%02d%02d%02dC340000.US", d.Year()%100, int(d.Month()), d.Day())
	warrants, err := ws.WarrantList(ctx, "700.HK", quote.WarrantFilter{SortBy: quote.WarrantLastDone, SortOrder: quote.WarrantDesc}, quote.WarrantEN)
	if err != nil || len(warrants) == 0 {
		log.Fatalf("warrants: %v", err)
	}
	warrantSymbol := warrants[0].Symbol

	type c = *quote.QuoteContext
	compare("static_info", `["700.HK","AAPL.US"]`, func(q c) (interface{}, error) { return q.StaticInfo(ctx, []string{"700.HK", "AAPL.US"}) }, ws, http)
	compare("quote", `["700.HK","AAPL.US"]`, func(q c) (interface{}, error) { return q.Quote(ctx, []string{"700.HK", "AAPL.US"}) }, ws, http)
	compare("option_quote", optionSymbol, func(q c) (interface{}, error) { return q.OptionQuote(ctx, []string{optionSymbol}) }, ws, http)
	compare("warrant_quote", warrantSymbol, func(q c) (interface{}, error) { return q.WarrantQuote(ctx, []string{warrantSymbol}) }, ws, http)
	for _, s := range []string{"700.HK", "AAPL.US"} {
		s := s
		compare("depth "+s, s, func(q c) (interface{}, error) { return q.Depth(ctx, s) }, ws, http)
	}
	compare("brokers 700.HK", "700.HK", func(q c) (interface{}, error) { return q.Brokers(ctx, "700.HK") }, ws, http)
	compare("participants", "-", func(q c) (interface{}, error) { return q.Participants(ctx) }, ws, http)
	compare("trades 700.HK", "700.HK, count=20", func(q c) (interface{}, error) { return q.Trades(ctx, "700.HK", 20) }, ws, http)
	compare("trades AAPL.US", "AAPL.US, count=50", func(q c) (interface{}, error) { return q.Trades(ctx, "AAPL.US", 50) }, ws, http)
	compare("intraday AAPL.US", "AAPL.US", func(q c) (interface{}, error) { return q.Intraday(ctx, "AAPL.US") }, ws, http)
	compare("intraday .VIX.US (index)", ".VIX.US", func(q c) (interface{}, error) { return q.Intraday(ctx, ".VIX.US") }, ws, http)
	compare("intraday 700.HK", "700.HK", func(q c) (interface{}, error) { return q.Intraday(ctx, "700.HK") }, ws, http)
	compare("candlesticks AAPL.US", "AAPL.US, Day, 20, NoAdjust", func(q c) (interface{}, error) {
		return q.Candlesticks(ctx, "AAPL.US", quote.PeriodDay, 20, quote.AdjustTypeNo)
	}, ws, http)
	compare("candlesticks 700.HK 1m", "700.HK, OneMinute, 30, ForwardAdjust", func(q c) (interface{}, error) {
		return q.Candlesticks(ctx, "700.HK", quote.PeriodOneMinute, 30, quote.AdjustTypeForward)
	}, ws, http)
	compare("candlesticks AAPL.US 1m", "AAPL.US, OneMinute, 1000, NoAdjust", func(q c) (interface{}, error) {
		return q.Candlesticks(ctx, "AAPL.US", quote.PeriodOneMinute, 1000, quote.AdjustTypeNo)
	}, ws, http)
	compare("history_candlesticks_by_offset 700.HK", "700.HK, Day, NoAdjust, backward, 2026-06-01, 10", func(q c) (interface{}, error) {
		return q.HistoryCandlesticksByOffset(ctx, "700.HK", quote.PeriodDay, quote.AdjustTypeNo, false, date(2026, 6, 1), 10)
	}, ws, http)
	for _, fwd := range []bool{false, true} {
		fwd := fwd
		compare(fmt.Sprintf("history_candlesticks_by_offset AAPL.US 1m All forward=%v", fwd),
			fmt.Sprintf("AAPL.US, OneMinute, NoAdjust, forward=%v, 2026-10-08 12:00, 600, All", fwd), func(q c) (interface{}, error) {
				t := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
				return q.HistoryCandlesticksByOffset(ctx, "AAPL.US", quote.PeriodOneMinute, quote.AdjustTypeNo, fwd, &t, 600, all...)
			}, ws, http)
	}
	compare("history_candlesticks_by_date AAPL.US", "AAPL.US, Week, ForwardAdjust, 2026-01-01..2026-03-01, All", func(q c) (interface{}, error) {
		return q.HistoryCandlesticksByDate(ctx, "AAPL.US", quote.PeriodWeek, quote.AdjustTypeForward, date(2026, 1, 1), date(2026, 3, 1), all...)
	}, ws, http)
	compare("history_candlesticks_by_date AAPL.US 1m All", "AAPL.US, OneMinute, NoAdjust, 2026-10-08..2026-10-09, All", func(q c) (interface{}, error) {
		return q.HistoryCandlesticksByDate(ctx, "AAPL.US", quote.PeriodOneMinute, quote.AdjustTypeNo, date(2026, 10, 8), date(2026, 10, 9), all...)
	}, ws, http)
	compare("history_candlesticks_by_date AAPL.US 1m All (uncapped)", "AAPL.US, OneMinute, NoAdjust, 2026-10-09..2026-10-09, All", func(q c) (interface{}, error) {
		return q.HistoryCandlesticksByDate(ctx, "AAPL.US", quote.PeriodOneMinute, quote.AdjustTypeNo, date(2026, 10, 9), date(2026, 10, 9), all...)
	}, ws, http)
	compare("option_chain_expiry_date_list AAPL.US", "AAPL.US", func(q c) (interface{}, error) { return q.OptionChainExpiryDateList(ctx, "AAPL.US") }, ws, http)
	compare("warrant_issuers", "-", func(q c) (interface{}, error) { return q.WarrantIssuers(ctx) }, ws, http)
	compare("warrant_list 700.HK", "700.HK, ChangeRate, Asc, no filters", func(q c) (interface{}, error) {
		return q.WarrantList(ctx, "700.HK", quote.WarrantFilter{SortBy: quote.WarrantChangeRate, SortOrder: quote.WarrantAsc}, quote.WarrantEN)
	}, ws, http)
	compare("trading_session", "-", func(q c) (interface{}, error) { return q.TradingSession(ctx) }, ws, http)
	compare("trading_days HK", "HK, 2026-10-01..2026-10-31", func(q c) (interface{}, error) {
		return q.TradingDays(ctx, openapi.MarketHK, date(2026, 10, 1), date(2026, 10, 31))
	}, ws, http)
	compare("capital_flow 700.HK", "700.HK", func(q c) (interface{}, error) { return q.CapitalFlow(ctx, "700.HK") }, ws, http)
	compare("capital_distribution 700.HK", "700.HK", func(q c) (interface{}, error) { return q.CapitalDistribution(ctx, "700.HK") }, ws, http)
	compare("calc_index", "[700.HK, AAPL.US, option, warrant], all indexes", func(q c) (interface{}, error) {
		return q.CalcIndex(ctx, []string{"700.HK", "AAPL.US", optionSymbol, warrantSymbol}, []quote.CalcIndex{
			quote.CalcIndexLastDone, quote.CalcIndexChangeVal, quote.CalcIndexChangeRate, quote.CalcIndexVolume, quote.CalcIndexTurnover,
			quote.CalcIndexYtdChangeRate, quote.CalcIndexTurnoverRate, quote.CalcIndexTotalMarketValue, quote.CalcIndexCapitalFlow,
			quote.CalcIndexAmplitude, quote.CalcIndexVolumeRatio, quote.CalcIndexPeTTMRatio, quote.CalcIndexPbRatio, quote.CalcIndexDividendRatioTTM,
			quote.CalcIndexFiveDayChangeRate, quote.CalcIndexTenDayChangeRate, quote.CalcIndexHalfYearChangeRate, quote.CalcIndexFiveMinutesChangeRate,
			quote.CalcIndexExpiryDate, quote.CalcIndexStrikePrice, quote.CalcIndexUpperStrikePrice, quote.CalcIndexLowerStrikePrice,
			quote.CalcIndexOutstandingQTY, quote.CalcIndexOutstandingRatio, quote.CalcIndexPremium, quote.CalcIndexItmOtm, quote.CalcIndexImpliedVolatility,
			quote.CalcIndexWarrantDelta, quote.CalcIndexCallPrice, quote.CalcIndexToCallPrice, quote.CalcIndexEffectiveLeverage, quote.CalcIndexLeverageRatio,
			quote.CalcIndexConversionRatio, quote.CalcIndexBalancePoint, quote.CalcIndexOpenInterest,
			quote.CalcIndexDELTA, quote.CalcIndexGAMMA, quote.CalcIndexTHETA, quote.CalcIndexVEGA, quote.CalcIndexRHO,
		})
	}, ws, http)

	// Error paths: the business error code must match across transports.
	compare("ERR depth NOTEXIST.XX", "NOTEXIST.XX", func(q c) (interface{}, error) { return q.Depth(ctx, "NOTEXIST.XX") }, ws, http)
	compare("ERR candlesticks count=5000", "700.HK, Day, 5000", func(q c) (interface{}, error) {
		return q.Candlesticks(ctx, "700.HK", quote.PeriodDay, 5000, quote.AdjustTypeNo)
	}, ws, http)
	_ = os.Stdout.Sync()
}
