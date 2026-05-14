// Package portfolio provides a client for the Longbridge Portfolio Analytics API.
// It covers exchange rates and profit/loss analysis.
package portfolio

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/portfolio/jsontypes"
	httplib "github.com/longbridge/openapi-go/http"
)

// PortfolioContext is a client for the Longbridge Portfolio Analytics API.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	pctx, err := portfolio.NewFromCfg(conf)
//	rates, err := pctx.ExchangeRate(context.Background())
type PortfolioContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a PortfolioContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*PortfolioContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &PortfolioContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a PortfolioContext configured from environment variables.
func NewFromEnv() (*PortfolioContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// ExchangeRate returns the current exchange rates for supported currencies.
//
// Path: GET /v1/asset/exchange_rates
func (c *PortfolioContext) ExchangeRate(ctx context.Context) (*ExchangeRates, error) {
	var resp jsontypes.ExchangeRates
	if err := c.httpClient.Get(ctx, "/v1/asset/exchange_rates", url.Values{}, &resp); err != nil {
		return nil, err
	}
	out := &ExchangeRates{
		Exchanges: make([]ExchangeRate, 0, len(resp.Exchanges)),
	}
	for _, r := range resp.Exchanges {
		out.Exchanges = append(out.Exchanges, ExchangeRate{
			AverageRate:   r.AverageRate,
			BaseCurrency:  r.BaseCurrency,
			BidRate:       r.BidRate,
			OfferRate:     r.OfferRate,
			OtherCurrency: r.OtherCurrency,
		})
	}
	return out, nil
}

// ProfitAnalysisOptions are the optional parameters for ProfitAnalysis.
type ProfitAnalysisOptions struct {
	// Start date in YYYY-MM-DD format (optional).
	Start string
	// End date in YYYY-MM-DD format (optional).
	End string
}

// ProfitAnalysis returns the portfolio P&L analysis (summary + per-security breakdown).
// It concurrently calls two endpoints and merges the results.
//
// Paths:
//   - GET /v1/portfolio/profit-analysis-summary
//   - GET /v1/portfolio/profit-analysis-sublist
func (c *PortfolioContext) ProfitAnalysis(ctx context.Context, opts *ProfitAnalysisOptions) (*ProfitAnalysis, error) {
	if opts == nil {
		opts = &ProfitAnalysisOptions{}
	}
	startTS := dateToUnixOpt(opts.Start)
	endTS := dateToUnixEndOpt(opts.End)

	summaryQ := url.Values{}
	if startTS != nil {
		summaryQ.Set("start", fmt.Sprintf("%d", *startTS))
	}
	if endTS != nil {
		summaryQ.Set("end", fmt.Sprintf("%d", *endTS))
	}

	sublistQ := url.Values{}
	sublistQ.Set("profit_or_loss", "all")
	if startTS != nil {
		sublistQ.Set("start", fmt.Sprintf("%d", *startTS))
	}
	if endTS != nil {
		sublistQ.Set("end", fmt.Sprintf("%d", *endTS))
	}

	type summaryResult struct {
		resp jsontypes.ProfitAnalysisSummary
		err  error
	}
	type sublistResult struct {
		resp jsontypes.ProfitAnalysisSublist
		err  error
	}

	sumCh := make(chan summaryResult, 1)
	subCh := make(chan sublistResult, 1)

	go func() {
		var r jsontypes.ProfitAnalysisSummary
		err := c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis-summary", summaryQ, &r)
		sumCh <- summaryResult{r, err}
	}()
	go func() {
		var r jsontypes.ProfitAnalysisSublist
		err := c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis-sublist", sublistQ, &r)
		subCh <- sublistResult{r, err}
	}()

	sumRes := <-sumCh
	subRes := <-subCh

	if sumRes.err != nil {
		return nil, sumRes.err
	}
	if subRes.err != nil {
		return nil, subRes.err
	}

	summary, err := convertProfitAnalysisSummary(&sumRes.resp)
	if err != nil {
		return nil, err
	}
	sublist := convertProfitAnalysisSublist(&subRes.resp)

	return &ProfitAnalysis{
		Summary: *summary,
		Sublist: *sublist,
	}, nil
}

// ProfitAnalysisByMarketOptions are the parameters for ProfitAnalysisByMarket.
type ProfitAnalysisByMarketOptions struct {
	// Market filter, e.g. "HK", "US" (optional).
	Market string
	// Start date in YYYY-MM-DD format (optional).
	Start string
	// End date in YYYY-MM-DD format (optional).
	End string
	// Currency filter (optional).
	Currency string
	// Page number (1-based).
	Page uint32
	// Page size.
	Size uint32
}

// ProfitAnalysisByMarket returns paginated P&L analysis filtered by market.
//
// Path: GET /v1/portfolio/profit-analysis/by-market
func (c *PortfolioContext) ProfitAnalysisByMarket(ctx context.Context, opts *ProfitAnalysisByMarketOptions) (*ProfitAnalysisByMarket, error) {
	if opts == nil {
		opts = &ProfitAnalysisByMarketOptions{}
	}
	q := url.Values{}
	q.Set("page", fmt.Sprintf("%d", opts.Page))
	q.Set("size", fmt.Sprintf("%d", opts.Size))
	if opts.Market != "" {
		q.Set("market", opts.Market)
	}
	if opts.Currency != "" {
		q.Set("currency", opts.Currency)
	}
	if ts := dateToUnixOpt(opts.Start); ts != nil {
		q.Set("start", fmt.Sprintf("%d", *ts))
	}
	if ts := dateToUnixEndOpt(opts.End); ts != nil {
		q.Set("end", fmt.Sprintf("%d", *ts))
	}

	var resp jsontypes.ProfitAnalysisByMarket
	if err := c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis/by-market", q, &resp); err != nil {
		return nil, err
	}
	return convertProfitAnalysisByMarket(&resp)
}

// ProfitAnalysisDetailOptions are the parameters for ProfitAnalysisDetail.
type ProfitAnalysisDetailOptions struct {
	// Symbol, e.g. "TSLA.US".
	Symbol string
	// Start date in YYYY-MM-DD format (optional).
	Start string
	// End date in YYYY-MM-DD format (optional).
	End string
}

// ProfitAnalysisDetail returns P&L detail for a specific security.
//
// Path: GET /v1/portfolio/profit-analysis/detail
func (c *PortfolioContext) ProfitAnalysisDetail(ctx context.Context, opts *ProfitAnalysisDetailOptions) (*ProfitAnalysisDetail, error) {
	if opts == nil {
		return nil, errors.New("opts must not be nil")
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(opts.Symbol))
	if ts := dateToUnixOpt(opts.Start); ts != nil {
		q.Set("start", fmt.Sprintf("%d", *ts))
	}
	if ts := dateToUnixEndOpt(opts.End); ts != nil {
		q.Set("end", fmt.Sprintf("%d", *ts))
	}

	var resp jsontypes.ProfitAnalysisDetail
	if err := c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis/detail", q, &resp); err != nil {
		return nil, err
	}
	return convertProfitAnalysisDetail(&resp)
}

// ProfitAnalysisFlowsOptions are the parameters for ProfitAnalysisFlows.
type ProfitAnalysisFlowsOptions struct {
	// Symbol, e.g. "TSLA.US".
	Symbol string
	// Page number (1-based).
	Page uint32
	// Page size.
	Size uint32
	// Whether to include derivative flows.
	Derivative bool
	// Start date in YYYY-MM-DD format (optional).
	Start string
	// End date in YYYY-MM-DD format (optional).
	End string
}

// ProfitAnalysisFlows returns paginated P&L flow records for a security.
//
// Path: GET /v1/portfolio/profit-analysis/flows
func (c *PortfolioContext) ProfitAnalysisFlows(ctx context.Context, opts *ProfitAnalysisFlowsOptions) (*ProfitAnalysisFlows, error) {
	if opts == nil {
		return nil, errors.New("opts must not be nil")
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(opts.Symbol))
	q.Set("page", fmt.Sprintf("%d", opts.Page))
	q.Set("size", fmt.Sprintf("%d", opts.Size))
	if opts.Derivative {
		q.Set("derivative", "true")
	} else {
		q.Set("derivative", "false")
	}
	if opts.Start != "" {
		q.Set("start", opts.Start)
	}
	if opts.End != "" {
		q.Set("end", opts.End)
	}

	var resp jsontypes.ProfitAnalysisFlows
	if err := c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis/flows", q, &resp); err != nil {
		return nil, err
	}
	return convertProfitAnalysisFlows(&resp)
}

// ── internal converters ───────────────────────────────────────────

func convertProfitAnalysisSummary(j *jsontypes.ProfitAnalysisSummary) (*ProfitAnalysisSummary, error) {
	s := &ProfitAnalysisSummary{
		Currency:          j.Currency,
		CurrentTotalAsset: decimalFromStr(j.CurrentTotalAsset),
		StartDate:         j.StartDate,
		EndDate:           j.EndDate,
		StartTime:         j.StartTime,
		EndTime:           j.EndTime,
		EndingAssetValue:  decimalFromStr(j.EndingAssetValue),
		InitialAssetValue: decimalFromStr(j.InitialAssetValue),
		InvestAmount:      decimalFromStr(j.InvestAmount),
		IsTraded:          j.IsTraded,
		SumProfit:         decimalFromStr(j.SumProfit),
		SumProfitRate:     decimalFromStr(j.SumProfitRate),
	}
	s.Profits = convertProfitSummaryBreakdown(&j.Profits)
	return s, nil
}

func convertProfitSummaryBreakdown(j *jsontypes.ProfitSummaryBreakdown) ProfitSummaryBreakdown {
	b := ProfitSummaryBreakdown{
		Stock:                       decimalFromStr(j.Stock),
		Fund:                        decimalFromStr(j.Fund),
		Crypto:                      decimalFromStr(j.Crypto),
		Mmf:                         decimalFromStr(j.Mmf),
		Other:                       decimalFromStr(j.Other),
		CumulativeTransactionAmount: decimalFromStr(j.CumulativeTransactionAmount),
		TradeOrderNum:               j.TradeOrderNum,
		TradeStockNum:               j.TradeStockNum,
		Ipo:                         decimalFromStr(j.Ipo),
		IpoHit:                      j.IpoHit,
		IpoSubscription:             j.IpoSubscription,
	}
	for _, info := range j.SummaryInfo {
		b.SummaryInfo = append(b.SummaryInfo, ProfitSummaryInfo{
			AssetType:     assetTypeFromString(info.AssetType),
			ProfitMax:     info.ProfitMax,
			ProfitMaxName: info.ProfitMaxName,
			LossMax:       info.LossMax,
			LossMaxName:   info.LossMaxName,
		})
	}
	return b
}

func convertProfitAnalysisSublist(j *jsontypes.ProfitAnalysisSublist) *ProfitAnalysisSublist {
	s := &ProfitAnalysisSublist{
		Start:       j.Start,
		End:         j.End,
		StartDate:   j.StartDate,
		EndDate:     j.EndDate,
		UpdatedAt:   j.UpdatedAt,
		UpdatedDate: j.UpdatedDate,
		Items:       make([]ProfitAnalysisItem, 0, len(j.Items)),
	}
	for _, item := range j.Items {
		s.Items = append(s.Items, ProfitAnalysisItem{
			Name:              item.Name,
			Market:            item.Market,
			IsHolding:         item.IsHolding,
			Profit:            decimalFromStr(item.Profit),
			ProfitRate:        decimalFromStr(item.ProfitRate),
			ClearanceTimes:    item.ClearanceTimes,
			ItemType:          assetTypeFromString(item.ItemType),
			Currency:          item.Currency,
			Symbol:            counterIDToSymbol(item.Symbol),
			HoldingPeriod:     item.HoldingPeriod,
			SecurityCode:      item.SecurityCode,
			Isin:              item.Isin,
			UnderlyingProfit:  decimalFromStr(item.UnderlyingProfit),
			DerivativesProfit: decimalFromStr(item.DerivativesProfit),
			OrderProfit:       decimalFromStr(item.OrderProfit),
		})
	}
	return s
}

func convertProfitAnalysisDetail(j *jsontypes.ProfitAnalysisDetail) (*ProfitAnalysisDetail, error) {
	return &ProfitAnalysisDetail{
		Profit:               decimalFromStr(j.Profit),
		UnderlyingDetails:    convertProfitDetails(&j.UnderlyingDetails),
		DerivativePnlDetails: convertProfitDetails(&j.DerivativePnlDetails),
		Name:                 j.Name,
		UpdatedAt:            j.UpdatedAt,
		UpdatedDate:          j.UpdatedDate,
		Currency:             j.Currency,
		DefaultTag:           j.DefaultTag,
		Start:                j.Start,
		End:                  j.End,
		StartDate:            j.StartDate,
		EndDate:              j.EndDate,
	}, nil
}

func convertProfitDetails(j *jsontypes.ProfitDetails) ProfitDetails {
	return ProfitDetails{
		HoldingValue:             decimalFromStr(j.HoldingValue),
		Profit:                   decimalFromStr(j.Profit),
		CumulativeCreditedAmount: decimalFromStr(j.CumulativeCreditedAmount),
		CreditedDetails:          convertProfitDetailEntries(j.CreditedDetails),
		CumulativeDebitedAmount:  decimalFromStr(j.CumulativeDebitedAmount),
		DebitedDetails:           convertProfitDetailEntries(j.DebitedDetails),
		CumulativeFeeAmount:      decimalFromStr(j.CumulativeFeeAmount),
		FeeDetails:               convertProfitDetailEntries(j.FeeDetails),
		ShortHoldingValue:        decimalFromStr(j.ShortHoldingValue),
		LongHoldingValue:         decimalFromStr(j.LongHoldingValue),
		HoldingValueAtBeginning:  decimalFromStr(j.HoldingValueAtBeginning),
		HoldingValueAtEnding:     decimalFromStr(j.HoldingValueAtEnding),
	}
}

func convertProfitDetailEntries(jj []jsontypes.ProfitDetailEntry) []ProfitDetailEntry {
	out := make([]ProfitDetailEntry, 0, len(jj))
	for _, e := range jj {
		out = append(out, ProfitDetailEntry{
			Describe: e.Describe,
			Amount:   decimalFromStr(e.Amount),
		})
	}
	return out
}

func convertProfitAnalysisByMarket(j *jsontypes.ProfitAnalysisByMarket) (*ProfitAnalysisByMarket, error) {
	out := &ProfitAnalysisByMarket{
		Profit:     decimalFromStr(j.Profit),
		HasMore:    j.HasMore,
		StockItems: make([]ProfitAnalysisByMarketItem, 0, len(j.StockItems)),
	}
	for _, item := range j.StockItems {
		out.StockItems = append(out.StockItems, ProfitAnalysisByMarketItem{
			Code:   item.Code,
			Name:   item.Name,
			Market: item.Market,
			Profit: decimalFromStr(item.Profit),
		})
	}
	return out, nil
}

func convertProfitAnalysisFlows(j *jsontypes.ProfitAnalysisFlows) (*ProfitAnalysisFlows, error) {
	out := &ProfitAnalysisFlows{
		HasMore:   j.HasMore,
		FlowsList: make([]FlowItem, 0, len(j.FlowsList)),
	}
	for _, f := range j.FlowsList {
		out.FlowsList = append(out.FlowsList, FlowItem{
			ExecutedDate:      f.ExecutedDate,
			ExecutedTimestamp: f.ExecutedTimestamp,
			Code:              f.Code,
			Direction:         flowDirectionFromString(f.Direction),
			ExecutedQuantity:  decimalFromStr(f.ExecutedQuantity),
			ExecutedPrice:     decimalFromStr(f.ExecutedPrice),
			ExecutedCost:      decimalFromStr(f.ExecutedCost),
			Describe:          f.Describe,
		})
	}
	return out, nil
}

// ── helpers ───────────────────────────────────────────────────────

// decimalFromStr parses a decimal string; returns nil for empty or invalid values.
func decimalFromStr(s string) *decimal.Decimal {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return &d
}

// dateToUnixOpt converts an optional "YYYY-MM-DD" string to a unix timestamp
// at midnight UTC. Returns nil if the input is empty or invalid.
func dateToUnixOpt(date string) *int64 {
	if date == "" {
		return nil
	}
	parts := strings.Split(date, "-")
	if len(parts) != 3 {
		return nil
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil
	}
	ts := t.UTC().Unix()
	return &ts
}

// dateToUnixEndOpt converts an optional "YYYY-MM-DD" string to an end-of-day
// unix timestamp (23:59:59 UTC). Returns nil if the input is empty or invalid.
func dateToUnixEndOpt(date string) *int64 {
	ts := dateToUnixOpt(date)
	if ts == nil {
		return nil
	}
	end := *ts + 86399
	return &end
}

// symbolToCounterID converts a symbol like "TSLA.US" to a counter_id like "ST/US/TSLA".
func symbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])
	return fmt.Sprintf("ST/%s/%s", market, code)
}

// counterIDToSymbol converts a counter_id like "ST/US/TSLA" to a symbol like "TSLA.US".
func counterIDToSymbol(counterID string) string {
	parts := strings.SplitN(counterID, "/", 3)
	if len(parts) == 3 {
		return fmt.Sprintf("%s.%s", parts[2], parts[1])
	}
	return counterID
}
