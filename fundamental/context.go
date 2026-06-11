// Package fundamental provides a client for the Longbridge Fundamental
// OpenAPI. It covers financial reports, analyst ratings, dividends, EPS
// forecasts, consensus estimates, valuation metrics, company overview,
// executives, shareholders, fund holders, corporate actions, investor
// relations, operating reports, buyback data, and stock ratings.
package fundamental

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/config"
	counterpkg "github.com/longbridge/openapi-go/counter"
	"github.com/longbridge/openapi-go/fundamental/jsontypes"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/counter"
)

// FundamentalContext is a client for the Longbridge Fundamental OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	fctx, err := fundamental.NewFromCfg(conf)
//	reports, err := fctx.FinancialReport(context.Background(), "700.HK",
//	    fundamental.FinancialReportKindAll, nil)
type FundamentalContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a FundamentalContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*FundamentalContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &FundamentalContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a FundamentalContext configured from environment
// variables.
func NewFromEnv() (*FundamentalContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// ─── helpers ───────────────────────────────────────────────────────────────

// symbolToCounterID converts a symbol like "TSLA.US" to a counter_id like
// "ST/US/TSLA". All symbols are treated as equities (ST prefix).
func symbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])
	return fmt.Sprintf("ST/%s/%s", market, code)
}

func counterIDToSymbol(counterID string) string { return counter.IDToSymbol(counterID) }

// decimalFromString parses a decimal string; returns nil for empty strings or
// unparseable values.
func decimalFromString(s string) *decimal.Decimal {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return &d
}

// decimalFromStringZero parses a decimal string; returns zero for empty strings.
func decimalFromStringZero(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}

// ─── FinancialReport ───────────────────────────────────────────────────────

// FinancialReport fetches financial reports for a security.
//
// Path: GET /v1/quote/financial-reports
func (c *FundamentalContext) FinancialReport(
	ctx context.Context,
	symbol string,
	kind FinancialReportKind,
	period *FinancialReportPeriod,
) (*FinancialReports, error) {
	kindStr := financialReportKindStr(kind)
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("kind", kindStr)
	if period != nil {
		q.Set("report", financialReportPeriodStr(*period))
	}
	var resp jsontypes.FinancialReports
	if err := c.httpClient.Get(ctx, "/v1/quote/financial-reports", q, &resp); err != nil {
		return nil, err
	}
	return &FinancialReports{List: json.RawMessage(resp.List)}, nil
}

func financialReportKindStr(k FinancialReportKind) string {
	switch k {
	case FinancialReportKindIncomeStatement:
		return "IS"
	case FinancialReportKindBalanceSheet:
		return "BS"
	case FinancialReportKindCashFlow:
		return "CF"
	default:
		return "ALL"
	}
}

func financialReportPeriodStr(p FinancialReportPeriod) string {
	switch p {
	case FinancialReportPeriodAnnual:
		return "af"
	case FinancialReportPeriodSemiAnnual:
		return "saf"
	case FinancialReportPeriodQ1:
		return "q1"
	case FinancialReportPeriodQ2:
		return "q2"
	case FinancialReportPeriodQ3:
		return "q3"
	case FinancialReportPeriodQuarterlyFull:
		return "qf"
	case FinancialReportPeriodThreeQ:
		return "3q"
	default:
		return "af"
	}
}

// ─── InstitutionRating ────────────────────────────────────────────────────

// InstitutionRating fetches analyst ratings for a security by combining the
// latest snapshot and the consensus summary.
//
// Paths: GET /v1/quote/institution-rating-latest
//
//	GET /v1/quote/institution-ratings
func (c *FundamentalContext) InstitutionRating(
	ctx context.Context,
	symbol string,
) (*InstitutionRating, error) {
	cid := symbolToCounterID(symbol)
	q := url.Values{}
	q.Set("counter_id", cid)

	type result struct {
		latest  jsontypes.InstitutionRatingLatest
		summary jsontypes.InstitutionRatingSummary
		latErr  error
		sumErr  error
	}

	ch := make(chan result, 1)
	go func() {
		var r result
		r.latErr = c.httpClient.Get(ctx, "/v1/quote/institution-rating-latest", q, &r.latest)
		r.sumErr = c.httpClient.Get(ctx, "/v1/quote/institution-ratings", q, &r.summary)
		ch <- r
	}()
	r := <-ch

	if r.latErr != nil {
		return nil, r.latErr
	}
	if r.sumErr != nil {
		return nil, r.sumErr
	}

	out := &InstitutionRating{
		Latest:  convertInstitutionRatingLatest(&r.latest),
		Summary: convertInstitutionRatingSummary(&r.summary),
	}
	return out, nil
}

// InstitutionRatingDetail fetches historical analyst rating details for a
// security.
//
// Path: GET /v1/quote/institution-ratings/detail
func (c *FundamentalContext) InstitutionRatingDetail(
	ctx context.Context,
	symbol string,
) (*InstitutionRatingDetail, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.InstitutionRatingDetail
	if err := c.httpClient.Get(ctx, "/v1/quote/institution-ratings/detail", q, &resp); err != nil {
		return nil, err
	}
	return convertInstitutionRatingDetail(&resp), nil
}

// ─── Dividend ─────────────────────────────────────────────────────────────

// Dividend fetches dividend history for a security.
//
// Path: GET /v1/quote/dividends
func (c *FundamentalContext) Dividend(
	ctx context.Context,
	symbol string,
) (*DividendList, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.DividendList
	if err := c.httpClient.Get(ctx, "/v1/quote/dividends", q, &resp); err != nil {
		return nil, err
	}
	return convertDividendList(&resp), nil
}

// DividendDetail fetches detailed dividend information for a security.
//
// Path: GET /v1/quote/dividends/details
func (c *FundamentalContext) DividendDetail(
	ctx context.Context,
	symbol string,
) (*DividendList, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.DividendList
	if err := c.httpClient.Get(ctx, "/v1/quote/dividends/details", q, &resp); err != nil {
		return nil, err
	}
	return convertDividendList(&resp), nil
}

// ─── ForecastEps ──────────────────────────────────────────────────────────

// ForecastEps fetches EPS forecasts for a security.
//
// Path: GET /v1/quote/forecast-eps
func (c *FundamentalContext) ForecastEps(
	ctx context.Context,
	symbol string,
) (*ForecastEps, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.ForecastEps
	if err := c.httpClient.Get(ctx, "/v1/quote/forecast-eps", q, &resp); err != nil {
		return nil, err
	}
	return convertForecastEps(&resp), nil
}

// ─── Consensus ────────────────────────────────────────────────────────────

// Consensus fetches financial consensus estimates for a security.
//
// Path: GET /v1/quote/financial-consensus-detail
func (c *FundamentalContext) Consensus(
	ctx context.Context,
	symbol string,
) (*FinancialConsensus, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.FinancialConsensus
	if err := c.httpClient.Get(ctx, "/v1/quote/financial-consensus-detail", q, &resp); err != nil {
		return nil, err
	}
	return convertFinancialConsensus(&resp), nil
}

// ─── Valuation ────────────────────────────────────────────────────────────

// Valuation fetches valuation metrics (PE/PB/PS/dividend yield) for a security.
//
// Path: GET /v1/quote/valuation
func (c *FundamentalContext) Valuation(
	ctx context.Context,
	symbol string,
) (*ValuationData, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("indicator", "pe")
	q.Set("range", "1")
	var resp jsontypes.ValuationData
	if err := c.httpClient.Get(ctx, "/v1/quote/valuation", q, &resp); err != nil {
		return nil, err
	}
	return convertValuationData(&resp), nil
}

// ValuationHistory fetches historical valuation data for a security.
//
// Path: GET /v1/quote/valuation/detail
func (c *FundamentalContext) ValuationHistory(
	ctx context.Context,
	symbol string,
) (*ValuationHistoryResponse, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.ValuationHistoryResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/valuation/detail", q, &resp); err != nil {
		return nil, err
	}
	return convertValuationHistoryResponse(&resp), nil
}

// ─── IndustryValuation ───────────────────────────────────────────────────

// IndustryValuation fetches valuation comparison against industry peers.
//
// Path: GET /v1/quote/industry-valuation-comparison
func (c *FundamentalContext) IndustryValuation(
	ctx context.Context,
	symbol string,
) (*IndustryValuationList, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.IndustryValuationList
	if err := c.httpClient.Get(ctx, "/v1/quote/industry-valuation-comparison", q, &resp); err != nil {
		return nil, err
	}
	return convertIndustryValuationList(&resp), nil
}

// IndustryValuationDist fetches valuation distribution within the industry.
//
// Path: GET /v1/quote/industry-valuation-distribution
func (c *FundamentalContext) IndustryValuationDist(
	ctx context.Context,
	symbol string,
) (*IndustryValuationDist, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.IndustryValuationDist
	if err := c.httpClient.Get(ctx, "/v1/quote/industry-valuation-distribution", q, &resp); err != nil {
		return nil, err
	}
	return convertIndustryValuationDist(&resp), nil
}

// ─── Company ──────────────────────────────────────────────────────────────

// Company fetches company overview information.
//
// Path: GET /v1/quote/comp-overview
func (c *FundamentalContext) Company(
	ctx context.Context,
	symbol string,
) (*CompanyOverview, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.CompanyOverview
	if err := c.httpClient.Get(ctx, "/v1/quote/comp-overview", q, &resp); err != nil {
		return nil, err
	}
	return convertCompanyOverview(&resp), nil
}

// ─── Executive ────────────────────────────────────────────────────────────

// Executive fetches executive and board member information.
//
// Path: GET /v1/quote/company-professionals
func (c *FundamentalContext) Executive(
	ctx context.Context,
	symbol string,
) (*ExecutiveList, error) {
	q := url.Values{}
	q.Set("counter_ids", symbolToCounterID(symbol))
	var resp jsontypes.ExecutiveList
	if err := c.httpClient.Get(ctx, "/v1/quote/company-professionals", q, &resp); err != nil {
		return nil, err
	}
	return convertExecutiveList(&resp), nil
}

// ─── Shareholder ──────────────────────────────────────────────────────────

// Shareholder fetches major shareholders for a security.
//
// Path: GET /v1/quote/shareholders
func (c *FundamentalContext) Shareholder(
	ctx context.Context,
	symbol string,
) (*ShareholderList, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.ShareholderList
	if err := c.httpClient.Get(ctx, "/v1/quote/shareholders", q, &resp); err != nil {
		return nil, err
	}
	return convertShareholderList(&resp), nil
}

// ─── FundHolder ───────────────────────────────────────────────────────────

// FundHolder fetches funds and ETFs that hold a security.
//
// Path: GET /v1/quote/fund-holders
func (c *FundamentalContext) FundHolder(
	ctx context.Context,
	symbol string,
) (*FundHolders, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.FundHolders
	if err := c.httpClient.Get(ctx, "/v1/quote/fund-holders", q, &resp); err != nil {
		return nil, err
	}
	return convertFundHolders(&resp), nil
}

// ─── CorpAction ───────────────────────────────────────────────────────────

// CorpAction fetches corporate actions (dividends, splits, buybacks, etc.).
//
// Path: GET /v1/quote/company-act
func (c *FundamentalContext) CorpAction(
	ctx context.Context,
	symbol string,
) (*CorpActions, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("req_type", "1")
	q.Set("version", "3")
	var resp jsontypes.CorpActions
	if err := c.httpClient.Get(ctx, "/v1/quote/company-act", q, &resp); err != nil {
		return nil, err
	}
	return convertCorpActions(&resp), nil
}

// ─── InvestRelation ───────────────────────────────────────────────────────

// InvestRelation fetches investor relations / investment holdings.
//
// Path: GET /v1/quote/invest-relations
func (c *FundamentalContext) InvestRelation(
	ctx context.Context,
	symbol string,
) (*InvestRelations, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("count", "0")
	var resp jsontypes.InvestRelations
	if err := c.httpClient.Get(ctx, "/v1/quote/invest-relations", q, &resp); err != nil {
		return nil, err
	}
	return convertInvestRelations(&resp), nil
}

// ─── Operating ────────────────────────────────────────────────────────────

// Operating fetches operating metrics and financial report summaries.
//
// Path: GET /v1/quote/operatings
func (c *FundamentalContext) Operating(
	ctx context.Context,
	symbol string,
) (*OperatingList, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.OperatingList
	if err := c.httpClient.Get(ctx, "/v1/quote/operatings", q, &resp); err != nil {
		return nil, err
	}
	return convertOperatingList(&resp), nil
}

// ─── Buyback ──────────────────────────────────────────────────────────────

// Buyback fetches buyback data for a security.
//
// Path: GET /v1/quote/buy-backs
func (c *FundamentalContext) Buyback(
	ctx context.Context,
	symbol string,
) (*BuybackData, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.BuybackData
	if err := c.httpClient.Get(ctx, "/v1/quote/buy-backs", q, &resp); err != nil {
		return nil, err
	}
	return convertBuybackData(&resp), nil
}

// ─── Ratings ──────────────────────────────────────────────────────────────

// Ratings fetches stock ratings for a security.
//
// Path: GET /v1/quote/ratings
func (c *FundamentalContext) Ratings(
	ctx context.Context,
	symbol string,
) (*StockRatings, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.StockRatings
	if err := c.httpClient.Get(ctx, "/v1/quote/ratings", q, &resp); err != nil {
		return nil, err
	}
	return convertStockRatings(&resp), nil
}

// ─── internal converters ──────────────────────────────────────────────────

func convertDividendList(j *jsontypes.DividendList) *DividendList {
	out := &DividendList{
		List: make([]DividendItem, 0, len(j.List)),
	}
	for _, item := range j.List {
		out.List = append(out.List, DividendItem{
			Symbol:      counterIDToSymbol(item.CounterID),
			ID:          item.ID,
			Desc:        item.Desc,
			RecordDate:  item.RecordDate,
			ExDate:      item.ExDate,
			PaymentDate: item.PaymentDate,
		})
	}
	return out
}

func convertRatingEvaluate(j jsontypes.RatingEvaluate) RatingEvaluate {
	return RatingEvaluate{
		Buy:       j.Buy,
		Over:      j.Over,
		Hold:      j.Hold,
		Under:     j.Under,
		Sell:      j.Sell,
		NoOpinion: j.NoOpinion,
		Total:     j.Total,
		StartDate: j.StartDate,
		EndDate:   j.EndDate,
	}
}

func convertRatingTarget(j jsontypes.RatingTarget) RatingTarget {
	return RatingTarget{
		HighestPrice: decimalFromString(j.HighestPrice),
		LowestPrice:  decimalFromString(j.LowestPrice),
		PrevClose:    decimalFromString(j.PrevClose),
		StartDate:    j.StartDate,
		EndDate:      j.EndDate,
	}
}

func convertInstitutionRatingLatest(j *jsontypes.InstitutionRatingLatest) InstitutionRatingLatest {
	return InstitutionRatingLatest{
		Evaluate:       convertRatingEvaluate(j.Evaluate),
		Target:         convertRatingTarget(j.Target),
		IndustryID:     j.IndustryID,
		IndustryName:   j.IndustryName,
		IndustryRank:   j.IndustryRank,
		IndustryTotal:  j.IndustryTotal,
		IndustryMean:   j.IndustryMean,
		IndustryMedian: j.IndustryMedian,
	}
}

func convertInstitutionRatingSummary(j *jsontypes.InstitutionRatingSummary) InstitutionRatingSummary {
	return InstitutionRatingSummary{
		CcySymbol: j.CcySymbol,
		Change:    decimalFromString(j.Change),
		Evaluate: RatingSummaryEvaluate{
			Buy:       j.Evaluate.Buy,
			Date:      j.Evaluate.Date,
			Hold:      j.Evaluate.Hold,
			Sell:      j.Evaluate.Sell,
			StrongBuy: j.Evaluate.StrongBuy,
			Under:     j.Evaluate.Under,
		},
		Recommend: institutionRecommendFromString(j.Recommend),
		Target:    decimalFromString(j.Target),
		UpdatedAt: j.UpdatedAt,
	}
}

func convertInstitutionRatingDetail(j *jsontypes.InstitutionRatingDetail) *InstitutionRatingDetail {
	evalItems := make([]InstitutionRatingDetailEvaluateItem, 0, len(j.Evaluate.List))
	for _, item := range j.Evaluate.List {
		evalItems = append(evalItems, InstitutionRatingDetailEvaluateItem{
			Buy:       item.Buy,
			Date:      item.Date,
			Hold:      item.Hold,
			Sell:      item.Sell,
			StrongBuy: item.StrongBuy,
			NoOpinion: item.NoOpinion,
			Under:     item.Under,
		})
	}

	targetItems := make([]InstitutionRatingDetailTargetItem, 0, len(j.Target.List))
	for _, item := range j.Target.List {
		targetItems = append(targetItems, InstitutionRatingDetailTargetItem{
			AvgTarget: decimalFromString(item.AvgTarget),
			Date:      item.Date,
			MaxTarget: decimalFromString(item.MaxTarget),
			MinTarget: decimalFromString(item.MinTarget),
			Meet:      item.Meet,
			Price:     decimalFromString(item.Price),
			Timestamp: item.Timestamp,
		})
	}

	var dataPercent *decimal.Decimal
	if j.Target.DataPercent != nil {
		dataPercent = decimalFromString(*j.Target.DataPercent)
	}

	return &InstitutionRatingDetail{
		CcySymbol: j.CcySymbol,
		Evaluate: InstitutionRatingDetailEvaluate{
			List: evalItems,
		},
		Target: InstitutionRatingDetailTarget{
			DataPercent:        dataPercent,
			PredictionAccuracy: decimalFromString(j.Target.PredictionAccuracy),
			UpdatedAt:          j.Target.UpdatedAt,
			List:               targetItems,
		},
	}
}

func convertForecastEps(j *jsontypes.ForecastEps) *ForecastEps {
	items := make([]ForecastEpsItem, 0, len(j.Items))
	for _, item := range j.Items {
		items = append(items, ForecastEpsItem{
			ForecastEpsMedian:  decimalFromString(item.ForecastEpsMedian),
			ForecastEpsMean:    decimalFromString(item.ForecastEpsMean),
			ForecastEpsLowest:  decimalFromString(item.ForecastEpsLowest),
			ForecastEpsHighest: decimalFromString(item.ForecastEpsHighest),
			InstitutionTotal:   item.InstitutionTotal,
			InstitutionUp:      item.InstitutionUp,
			InstitutionDown:    item.InstitutionDown,
			ForecastStartDate:  time.Unix(parseTimestampNumber(item.ForecastStartDate), 0).UTC(),
			ForecastEndDate:    time.Unix(parseTimestampNumber(item.ForecastEndDate), 0).UTC(),
		})
	}
	return &ForecastEps{Items: items}
}

func convertConsensusDetail(j jsontypes.ConsensusDetail) ConsensusDetail {
	return ConsensusDetail{
		Key:         j.Key,
		Name:        j.Name,
		Description: j.Desc,
		Actual:      decimalFromString(j.Actual),
		Estimate:    decimalFromString(j.Estimate),
		CompValue:   decimalFromString(j.CompValue),
		CompDesc:    j.CompDesc,
		Comp:        j.Comp,
		IsReleased:  j.IsReleased,
	}
}

func convertFinancialConsensus(j *jsontypes.FinancialConsensus) *FinancialConsensus {
	reports := make([]ConsensusReport, 0, len(j.List))
	for _, r := range j.List {
		details := make([]ConsensusDetail, 0, len(r.Details))
		for _, d := range r.Details {
			details = append(details, convertConsensusDetail(d))
		}
		reports = append(reports, ConsensusReport{
			FiscalYear:   r.FiscalYear,
			FiscalPeriod: r.FiscalPeriod,
			PeriodText:   r.PeriodText,
			Details:      details,
		})
	}
	return &FinancialConsensus{
		List:          reports,
		CurrentIndex:  j.CurrentIndex,
		Currency:      j.Currency,
		OptPeriods:    j.OptPeriods,
		CurrentPeriod: j.CurrentPeriod,
	}
}

func convertValuationPoints(jps []jsontypes.ValuationPoint) []ValuationPoint {
	out := make([]ValuationPoint, 0, len(jps))
	for _, p := range jps {
		out = append(out, ValuationPoint{
			Timestamp: time.Unix(parseTimestampNumber(p.Timestamp), 0).UTC(),
			Value:     decimalFromString(p.Value),
		})
	}
	return out
}

func convertValuationMetricData(j *jsontypes.ValuationMetricData) *ValuationMetricData {
	if j == nil {
		return nil
	}
	return &ValuationMetricData{
		Desc:   j.Desc,
		High:   decimalFromString(j.High),
		Low:    decimalFromString(j.Low),
		Median: decimalFromString(j.Median),
		List:   convertValuationPoints(j.List),
	}
}

func convertValuationData(j *jsontypes.ValuationData) *ValuationData {
	return &ValuationData{
		Metrics: ValuationMetricsData{
			PE:     convertValuationMetricData(j.Metrics.PE),
			PB:     convertValuationMetricData(j.Metrics.PB),
			PS:     convertValuationMetricData(j.Metrics.PS),
			DvdYld: convertValuationMetricData(j.Metrics.DvdYld),
		},
	}
}

func convertValuationHistoryMetric(j *jsontypes.ValuationHistoryMetric) *ValuationHistoryMetric {
	if j == nil {
		return nil
	}
	return &ValuationHistoryMetric{
		Desc:   j.Desc,
		High:   decimalFromString(j.High),
		Low:    decimalFromString(j.Low),
		Median: decimalFromString(j.Median),
		List:   convertValuationPoints(j.List),
	}
}

func convertValuationHistoryResponse(j *jsontypes.ValuationHistoryResponse) *ValuationHistoryResponse {
	return &ValuationHistoryResponse{
		History: ValuationHistoryData{
			Metrics: ValuationHistoryMetrics{
				PE: convertValuationHistoryMetric(j.History.Metrics.PE),
				PB: convertValuationHistoryMetric(j.History.Metrics.PB),
				PS: convertValuationHistoryMetric(j.History.Metrics.PS),
			},
		},
	}
}

func convertIndustryValuationHistory(j jsontypes.IndustryValuationHistory) IndustryValuationHistory {
	return IndustryValuationHistory{
		Date: j.Date,
		PE:   decimalFromString(j.PE),
		PB:   decimalFromString(j.PB),
		PS:   decimalFromString(j.PS),
	}
}

func convertIndustryValuationList(j *jsontypes.IndustryValuationList) *IndustryValuationList {
	out := &IndustryValuationList{
		List: make([]IndustryValuationItem, 0, len(j.List)),
	}
	for _, item := range j.List {
		history := make([]IndustryValuationHistory, 0, len(item.History))
		for _, h := range item.History {
			history = append(history, convertIndustryValuationHistory(h))
		}
		out.List = append(out.List, IndustryValuationItem{
			Symbol:         counterIDToSymbol(item.CounterID),
			Name:           item.Name,
			Currency:       item.Currency,
			Assets:         decimalFromString(item.Assets),
			Bps:            decimalFromString(item.Bps),
			Eps:            decimalFromString(item.Eps),
			Dps:            decimalFromString(item.Dps),
			DivYld:         decimalFromString(item.DivYld),
			DivPayoutRatio: decimalFromString(item.DivPayoutRatio),
			FiveYAvgDps:    decimalFromString(item.FiveYAvgDps),
			PE:             decimalFromString(item.PE),
			History:        history,
		})
	}
	return out
}

func convertValuationDist(j *jsontypes.ValuationDist) *ValuationDist {
	if j == nil {
		return nil
	}
	return &ValuationDist{
		Low:       decimalFromString(j.Low),
		High:      decimalFromString(j.High),
		Median:    decimalFromString(j.Median),
		Value:     decimalFromString(j.Value),
		Ranking:   decimalFromString(j.Ranking),
		RankIndex: j.RankIndex,
		RankTotal: j.RankTotal,
	}
}

func convertIndustryValuationDist(j *jsontypes.IndustryValuationDist) *IndustryValuationDist {
	return &IndustryValuationDist{
		PE: convertValuationDist(j.PE),
		PB: convertValuationDist(j.PB),
		PS: convertValuationDist(j.PS),
	}
}

func convertCompanyOverview(j *jsontypes.CompanyOverview) *CompanyOverview {
	return &CompanyOverview{
		Name:           j.Name,
		CompanyName:    j.CompanyName,
		Founded:        j.Founded,
		ListingDate:    j.ListingDate,
		Market:         j.Market,
		Region:         j.Region,
		Address:        j.Address,
		OfficeAddress:  j.OfficeAddress,
		Website:        j.Website,
		IssuePrice:     decimalFromString(j.IssuePrice),
		SharesOffered:  j.SharesOffered,
		Chairman:       j.Chairman,
		Secretary:      j.Secretary,
		AuditInst:      j.AuditInst,
		Category:       j.Category,
		YearEnd:        j.YearEnd,
		Employees:      j.Employees,
		Phone:          j.Phone,
		Fax:            j.Fax,
		Email:          j.Email,
		LegalRepr:      j.LegalRepr,
		Manager:        j.Manager,
		BusLicense:     j.BusLicense,
		AccountingFirm: j.AccountingFirm,
		SecuritiesRep:  j.SecuritiesRep,
		LegalCounsel:   j.LegalCounsel,
		ZipCode:        j.ZipCode,
		Ticker:         j.Ticker,
		Icon:           j.Icon,
		Profile:        j.Profile,
		AdsRatio:       j.AdsRatio,
		Sector:         j.Sector,
	}
}

func convertExecutiveList(j *jsontypes.ExecutiveList) *ExecutiveList {
	groups := make([]ExecutiveGroup, 0, len(j.ProfessionalList))
	for _, g := range j.ProfessionalList {
		profs := make([]Professional, 0, len(g.Professionals))
		for _, p := range g.Professionals {
			profs = append(profs, Professional{
				ID:        p.ID,
				Name:      p.Name,
				NameZhCN:  p.NameZhCN,
				NameEn:    p.NameEn,
				Title:     p.Title,
				Biography: p.Biography,
				Photo:     p.Photo,
				WikiURL:   p.WikiURL,
			})
		}
		groups = append(groups, ExecutiveGroup{
			Symbol:        counterIDToSymbol(g.CounterID),
			ForwardURL:    g.ForwardURL,
			Total:         g.Total,
			Professionals: profs,
		})
	}
	return &ExecutiveList{ProfessionalList: groups}
}

func convertShareholderList(j *jsontypes.ShareholderList) *ShareholderList {
	shareholders := make([]Shareholder, 0, len(j.ShareholderList))
	for _, s := range j.ShareholderList {
		stocks := make([]ShareholderStock, 0, len(s.Stocks))
		for _, st := range s.Stocks {
			stocks = append(stocks, ShareholderStock{
				Symbol: counterIDToSymbol(st.CounterID),
				Code:   st.Code,
				Market: st.Market,
				Chg:    st.Chg,
			})
		}
		shareholders = append(shareholders, Shareholder{
			ShareholderID:   s.ShareholderID,
			ShareholderName: s.ShareholderName,
			InstitutionType: s.InstitutionType,
			PercentOfShares: decimalFromString(s.PercentOfShares),
			SharesChanged:   decimalFromString(s.SharesChanged),
			ReportDate:      s.ReportDate,
			Stocks:          stocks,
		})
	}
	return &ShareholderList{
		ShareholderList: shareholders,
		ForwardURL:      j.ForwardURL,
		Total:           j.Total,
	}
}

func convertFundHolders(j *jsontypes.FundHolders) *FundHolders {
	lists := make([]FundHolder, 0, len(j.Lists))
	for _, h := range j.Lists {
		lists = append(lists, FundHolder{
			Code:          h.Code,
			Symbol:        counterIDToSymbol(h.CounterID),
			Currency:      h.Currency,
			Name:          h.Name,
			PositionRatio: decimalFromStringZero(h.PositionRatio),
			ReportDate:    h.ReportDate,
		})
	}
	return &FundHolders{Lists: lists}
}

func convertCorpActionLive(j *jsontypes.CorpActionLive) *CorpActionLive {
	if j == nil {
		return nil
	}
	return &CorpActionLive{
		ID:        j.ID,
		Status:    json.RawMessage(j.Status),
		StartedAt: j.StartedAt,
		Name:      j.Name,
		Icon:      j.Icon,
	}
}

func convertCorpActions(j *jsontypes.CorpActions) *CorpActions {
	items := make([]CorpActionItem, 0, len(j.Items))
	for _, item := range j.Items {
		var sec *json.RawMessage
		if item.Security != nil {
			raw := json.RawMessage(*item.Security)
			sec = &raw
		}
		items = append(items, CorpActionItem{
			ID:           item.ID,
			Date:         item.Date,
			DateStr:      item.DateStr,
			DateType:     item.DateType,
			DateZone:     item.DateZone,
			ActType:      item.ActType,
			ActDesc:      item.ActDesc,
			Action:       item.Action,
			Recent:       item.Recent,
			IsDelay:      item.IsDelay,
			DelayContent: item.DelayContent,
			Live:         convertCorpActionLive(item.Live),
			Security:     sec,
		})
	}
	return &CorpActions{Items: items}
}

func convertInvestRelations(j *jsontypes.InvestRelations) *InvestRelations {
	secs := make([]InvestSecurity, 0, len(j.InvestSecurities))
	for _, s := range j.InvestSecurities {
		secs = append(secs, InvestSecurity{
			CompanyID:       s.CompanyID,
			CompanyName:     s.CompanyName,
			CompanyNameEn:   s.CompanyNameEn,
			CompanyNameZhCN: s.CompanyNameZhCN,
			Symbol:          counterIDToSymbol(s.CounterID),
			Currency:        s.Currency,
			PercentOfShares: decimalFromString(s.PercentOfShares),
			SharesRank:      s.SharesRank,
			SharesValue:     decimalFromString(s.SharesValue),
		})
	}
	return &InvestRelations{
		ForwardURL:       j.ForwardURL,
		InvestSecurities: secs,
	}
}

func convertOperatingList(j *jsontypes.OperatingList) *OperatingList {
	items := make([]OperatingItem, 0, len(j.List))
	for _, item := range j.List {
		indicators := make([]OperatingIndicator, 0, len(item.Financial.Indicators))
		for _, ind := range item.Financial.Indicators {
			indicators = append(indicators, OperatingIndicator{
				FieldName:      ind.FieldName,
				IndicatorName:  ind.IndicatorName,
				IndicatorValue: ind.IndicatorValue,
				Yoy:            decimalFromString(ind.Yoy),
			})
		}
		items = append(items, OperatingItem{
			ID:      item.ID,
			Report:  item.Report,
			Title:   item.Title,
			Txt:     item.Txt,
			Latest:  item.Latest,
			Keywords: item.Keywords,
			WebURL:  item.WebURL,
			Financial: OperatingFinancial{
				Code:     item.Financial.Code,
				Symbol:   counterIDToSymbol(item.Financial.CounterID),
				Currency: item.Financial.Currency,
				Name:       item.Financial.Name,
				Region:     item.Financial.Region,
				Report:     item.Financial.Report,
				ReportTxt:  item.Financial.ReportTxt,
				Indicators: indicators,
			},
		})
	}
	return &OperatingList{List: items}
}

func convertBuybackData(j *jsontypes.BuybackData) *BuybackData {
	out := &BuybackData{}

	if j.RecentBuybacks != nil {
		out.RecentBuybacks = &RecentBuybacks{
			Currency:           j.RecentBuybacks.Currency,
			NetBuybackTTM:      decimalFromString(j.RecentBuybacks.NetBuybackTTM),
			NetBuybackYieldTTM: decimalFromString(j.RecentBuybacks.NetBuybackYieldTTM),
		}
	}

	out.BuybackHistory = make([]BuybackHistoryItem, 0, len(j.BuybackHistory))
	for _, h := range j.BuybackHistory {
		out.BuybackHistory = append(out.BuybackHistory, BuybackHistoryItem{
			FiscalYear:           h.FiscalYear,
			FiscalYearRange:      h.FiscalYearRange,
			NetBuyback:           decimalFromString(h.NetBuyback),
			NetBuybackYield:      decimalFromString(h.NetBuybackYield),
			NetBuybackGrowthRate: decimalFromString(h.NetBuybackGrowthRate),
			Currency:             h.Currency,
		})
	}

	out.BuybackRatios = make([]BuybackRatios, 0, len(j.BuybackRatios))
	for _, r := range j.BuybackRatios {
		out.BuybackRatios = append(out.BuybackRatios, BuybackRatios{
			NetBuybackPayoutRatio:     decimalFromString(r.NetBuybackPayoutRatio),
			NetBuybackToCashflowRatio: decimalFromString(r.NetBuybackToCashflowRatio),
		})
	}

	return out
}

func convertRatingCategory(j jsontypes.RatingCategory) RatingCategory {
	subGroups := make([]RatingSubIndicatorGroup, 0, len(j.SubIndicators))
	for _, g := range j.SubIndicators {
		leaves := make([]RatingLeafIndicator, 0, len(g.SubIndicators))
		for _, l := range g.SubIndicators {
			leaves = append(leaves, RatingLeafIndicator{
				Name:      l.Name,
				Value:     l.Value,
				ValueType: l.ValueType,
				Score:     json.RawMessage(l.Score),
				Letter:    l.Letter,
			})
		}
		subGroups = append(subGroups, RatingSubIndicatorGroup{
			Indicator: RatingIndicator{
				Name:   g.Indicator.Name,
				Score:  json.RawMessage(g.Indicator.Score),
				Letter: g.Indicator.Letter,
			},
			SubIndicators: leaves,
		})
	}
	return RatingCategory{
		Kind:          j.Kind,
		SubIndicators: subGroups,
	}
}

func convertStockRatings(j *jsontypes.StockRatings) *StockRatings {
	ratings := make([]RatingCategory, 0, len(j.Ratings))
	for _, r := range j.Ratings {
		ratings = append(ratings, convertRatingCategory(r))
	}
	return &StockRatings{
		StyleTxtName:        j.StyleTxtName,
		ScaleTxtName:        j.ScaleTxtName,
		ReportPeriodTxt:     j.ReportPeriodTxt,
		MultiScore:          json.RawMessage(j.MultiScore),
		MultiLetter:         j.MultiLetter,
		MultiScoreChange:    j.MultiScoreChange,
		IndustryName:        j.IndustryName,
		IndustryRank:        json.RawMessage(j.IndustryRank),
		IndustryTotal:       json.RawMessage(j.IndustryTotal),
		IndustryMeanScore:   json.RawMessage(j.IndustryMeanScore),
		IndustryMedianScore: json.RawMessage(j.IndustryMedianScore),
		Ratings:             ratings,
	}
}

// ─── ShareholderTop ───────────────────────────────────────────────────────────

// ShareholderTop fetches the top shareholders list for a security.
//
// Path: GET /v1/quote/shareholders/top
func (c *FundamentalContext) ShareholderTop(
	ctx context.Context,
	symbol string,
) (*ShareholderTopResponse, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	if err := c.httpClient.Get(ctx, "/v1/quote/shareholders/top", q, &resp); err != nil {
		return nil, err
	}
	return &ShareholderTopResponse{Data: resp}, nil
}

// ─── ShareholderDetail ────────────────────────────────────────────────────────

// ShareholderDetail fetches the holding detail for a specific shareholder.
//
// Path: GET /v1/quote/shareholders/holding
func (c *FundamentalContext) ShareholderDetail(
	ctx context.Context,
	symbol string,
	objectID int64,
) (*ShareholderDetailResponse, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("object_id", strconv.FormatInt(objectID, 10))
	var resp json.RawMessage
	if err := c.httpClient.Get(ctx, "/v1/quote/shareholders/holding", q, &resp); err != nil {
		return nil, err
	}
	return &ShareholderDetailResponse{Data: resp}, nil
}

// ─── ValuationComparison ──────────────────────────────────────────────────────

// ValuationComparison fetches valuation comparison data for a symbol against
// a set of peer symbols.
//
// Path: GET /v1/quote/compare/valuation
//
// comparisonSymbols is a list of peer symbols (e.g. ["MSFT.US", "GOOG.US"])
// that are converted to counter_ids and serialized as a JSON array string in
// the comparison_counter_ids query parameter.
func (c *FundamentalContext) ValuationComparison(
	ctx context.Context,
	symbol string,
	currency string,
	comparisonSymbols []string,
) (*ValuationComparisonResponse, error) {
	counterIDs := make([]string, 0, len(comparisonSymbols))
	for _, s := range comparisonSymbols {
		counterIDs = append(counterIDs, symbolToCounterID(s))
	}
	counterIDsJSON, err := json.Marshal(counterIDs)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("currency", currency)
	q.Set("comparison_counter_ids", string(counterIDsJSON))
	var raw struct {
		List []struct {
			CounterID   string `json:"counter_id"`
			Name        string `json:"name"`
			Currency    string `json:"currency"`
			MarketValue string `json:"market_value"`
			PriceClose  string `json:"price_close"`
			Pe          string `json:"pe"`
			Pb          string `json:"pb"`
			Ps          string `json:"ps"`
			Roe         string `json:"roe"`
			Eps         string `json:"eps"`
			Bps         string `json:"bps"`
			Dps         string `json:"dps"`
			DivYld      string `json:"div_yld"`
			Assets      string `json:"assets"`
			History     []struct {
				Date string `json:"date"`
				Pe   string `json:"pe"`
				Pb   string `json:"pb"`
				Ps   string `json:"ps"`
			} `json:"history"`
		} `json:"list"`
	}
	if err := c.httpClient.Get(ctx, "/v1/quote/compare/valuation", q, &raw); err != nil {
		return nil, err
	}
	items := make([]*ValuationComparisonItem, 0, len(raw.List))
	for _, it := range raw.List {
		history := make([]*ValuationHistoryPoint, 0, len(it.History))
		for _, h := range it.History {
			history = append(history, &ValuationHistoryPoint{
				Date: unixSecsToRFC3339(h.Date),
				Pe:   h.Pe,
				Pb:   h.Pb,
				Ps:   h.Ps,
			})
		}
		items = append(items, &ValuationComparisonItem{
			Symbol:      counterIDToSymbol(it.CounterID),
			Name:        it.Name,
			Currency:    it.Currency,
			MarketValue: it.MarketValue,
			PriceClose:  it.PriceClose,
			Pe:          it.Pe,
			Pb:          it.Pb,
			Ps:          it.Ps,
			Roe:         it.Roe,
			Eps:         it.Eps,
			Bps:         it.Bps,
			Dps:         it.Dps,
			DivYld:      it.DivYld,
			Assets:      it.Assets,
			History:     history,
		})
	}
	return &ValuationComparisonResponse{List: items}, nil
}

// ─── EtfAssetAllocation ───────────────────────────────────────────────────

// EtfAssetAllocation fetches the ETF asset allocation (holdings / regional /
// asset class / industry) for an ETF symbol.
//
// The symbol is converted to its counter_id using the directory-aware
// counter package (so e.g. "QQQ.US" maps to "ETF/US/QQQ").
//
// Path: GET /v1/quote/etf-asset-allocation
func (c *FundamentalContext) EtfAssetAllocation(
	ctx context.Context,
	symbol string,
) (*AssetAllocationResponse, error) {
	q := url.Values{}
	q.Set("counter_id", counterpkg.SymbolToCounterID(symbol))
	var resp jsontypes.AssetAllocationResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/etf-asset-allocation", q, &resp); err != nil {
		return nil, err
	}
	return convertAssetAllocationResponse(&resp), nil
}

func convertHoldingDetail(j *jsontypes.HoldingDetail) *HoldingDetail {
	if j == nil {
		return nil
	}
	return &HoldingDetail{
		IndustryID:      j.IndustryID,
		IndustryName:    j.IndustryName,
		Index:           j.Index,
		IndexName:       j.IndexName,
		HoldingType:     j.HoldingType,
		HoldingTypeName: j.HoldingTypeName,
	}
}

func convertAssetAllocationResponse(j *jsontypes.AssetAllocationResponse) *AssetAllocationResponse {
	groups := make([]*AssetAllocationGroup, 0, len(j.Info))
	for _, g := range j.Info {
		items := make([]*AssetAllocationItem, 0, len(g.Lists))
		for _, item := range g.Lists {
			var symbol string
			if item.CounterID != "" {
				symbol = counterIDToSymbol(item.CounterID)
			}
			items = append(items, &AssetAllocationItem{
				Name:          item.Name,
				Code:          item.Code,
				PositionRatio: item.PositionRatio,
				Symbol:        symbol,
				NameLocales:   item.NameLocales,
				HoldingDetail: convertHoldingDetail(item.HoldingDetail),
			})
		}
		groups = append(groups, &AssetAllocationGroup{
			ReportDate: g.ReportDate,
			AssetType:  ElementType(g.AssetType),
			Lists:      items,
		})
	}
	return &AssetAllocationResponse{Info: groups}
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

// parseTimestampNumber converts a json.Number (int or quoted string) to int64 Unix seconds.
func parseTimestampNumber(n json.Number) int64 {
	s := n.String()
	// strip surrounding quotes if the API returned a JSON string
	s = strings.Trim(s, `"`)
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

// ─── BusinessSegments ─────────────────────────────────────────────────────

// BusinessSegments fetches the latest business segment breakdown for a security.
//
// Path: GET /v1/quote/fundamentals/business-segments
func (c *FundamentalContext) BusinessSegments(
	ctx context.Context,
	symbol string,
) (*BusinessSegments, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.BusinessSegments
	if err := c.httpClient.Get(ctx, "/v1/quote/fundamentals/business-segments", q, &resp); err != nil {
		return nil, err
	}
	return convertBusinessSegments(&resp), nil
}

// BusinessSegmentsHistory fetches historical business segment breakdowns for a
// security.
//
// Path: GET /v1/quote/fundamentals/business-segments/history
func (c *FundamentalContext) BusinessSegmentsHistory(
	ctx context.Context,
	symbol string,
	report string,
	cate string,
) (*BusinessSegmentsHistory, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	if report != "" {
		q.Set("report", report)
	}
	if cate != "" {
		q.Set("cate", cate)
	}
	var resp jsontypes.BusinessSegmentsHistory
	if err := c.httpClient.Get(ctx, "/v1/quote/fundamentals/business-segments/history", q, &resp); err != nil {
		return nil, err
	}
	return convertBusinessSegmentsHistory(&resp), nil
}

// ─── InstitutionRatingViews ───────────────────────────────────────────────

// InstitutionRatingViews fetches historical institutional rating views for a
// security.
//
// Path: GET /v1/quote/ratings/institutional
func (c *FundamentalContext) InstitutionRatingViews(
	ctx context.Context,
	symbol string,
) (*InstitutionRatingViews, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp jsontypes.InstitutionRatingViews
	if err := c.httpClient.Get(ctx, "/v1/quote/ratings/institutional", q, &resp); err != nil {
		return nil, err
	}
	return convertInstitutionRatingViews(&resp), nil
}

// ─── IndustryRank ─────────────────────────────────────────────────────────

// IndustryRank fetches the industry rank for a market.
//
// Path: GET /v1/quote/industry/rank
//
// indicator is an IndustryRankIndicator constant ("0"–"7").
// sortType is an IndustryRankSortType constant ("0" ascending, "1" descending).
func (c *FundamentalContext) IndustryRank(
	ctx context.Context,
	market string,
	indicator IndustryRankIndicator,
	sortType IndustryRankSortType,
	limit int,
) (*IndustryRankResponse, error) {
	q := url.Values{}
	q.Set("market", market)
	q.Set("indicator", string(indicator))
	q.Set("sort_type", string(sortType))
	q.Set("limit", strconv.Itoa(limit))
	var resp jsontypes.IndustryRankResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/industry/rank", q, &resp); err != nil {
		return nil, err
	}
	return convertIndustryRankResponse(&resp), nil
}

// ─── IndustryPeers ────────────────────────────────────────────────────────

// IndustryPeers fetches the industry peer chain for a security or industry.
//
// Path: GET /v1/quote/industries/peers
//
// counterID may be a regular symbol like "AAPL.US" (auto-converted) or an
// industry counter ID like "BK/US/123" (passed through as-is when it contains
// a "/").
func (c *FundamentalContext) IndustryPeers(
	ctx context.Context,
	counterID string,
	market string,
	industryID string,
) (*IndustryPeersResponse, error) {
	// pass industry counter IDs (BK/xx/xx, IN/xx/xx, etc.) through as-is
	cid := counterID
	if !strings.Contains(counterID, "/") {
		cid = symbolToCounterID(counterID)
	}
	q := url.Values{}
	q.Set("type", "1")
	q.Set("market", market)
	q.Set("industry_id", industryID)
	q.Set("counter_id", cid)
	var resp jsontypes.IndustryPeersResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/industries/peers", q, &resp); err != nil {
		return nil, err
	}
	return convertIndustryPeersResponse(&resp), nil
}

// ─── FinancialReportSnapshot ──────────────────────────────────────────────

// FinancialReportSnapshot fetches a financial report snapshot (earnings
// snapshot) for a security.
//
// Path: GET /v1/quote/financials/earnings-snapshot
func (c *FundamentalContext) FinancialReportSnapshot(
	ctx context.Context,
	symbol string,
	report string,
	fiscalYear int,
	fiscalPeriod string,
) (*FinancialReportSnapshot, error) {
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	if report != "" {
		q.Set("report", report)
	}
	if fiscalYear != 0 {
		q.Set("fiscal_year", strconv.Itoa(fiscalYear))
	}
	if fiscalPeriod != "" {
		q.Set("fiscal_period", fiscalPeriod)
	}
	var resp jsontypes.FinancialReportSnapshot
	if err := c.httpClient.Get(ctx, "/v1/quote/financials/earnings-snapshot", q, &resp); err != nil {
		return nil, err
	}
	return convertFinancialReportSnapshot(&resp), nil
}

// ─── new converters ───────────────────────────────────────────────────────

func convertBusinessSegments(j *jsontypes.BusinessSegments) *BusinessSegments {
	items := make([]BusinessSegmentItem, 0, len(j.Business))
	for _, b := range j.Business {
		items = append(items, BusinessSegmentItem{Name: b.Name, Percent: b.Percent})
	}
	return &BusinessSegments{
		Date:     j.Date,
		Total:    j.Total,
		Currency: j.Currency,
		Business: items,
	}
}

func convertBusinessSegmentHistoryItems(js []jsontypes.BusinessSegmentHistoryItem) []BusinessSegmentHistoryItem {
	out := make([]BusinessSegmentHistoryItem, 0, len(js))
	for _, b := range js {
		out = append(out, BusinessSegmentHistoryItem{Name: b.Name, Percent: b.Percent, Value: b.Value})
	}
	return out
}

func convertBusinessSegmentsHistory(j *jsontypes.BusinessSegmentsHistory) *BusinessSegmentsHistory {
	items := make([]BusinessSegmentsHistoricalItem, 0, len(j.Historical))
	for _, h := range j.Historical {
		items = append(items, BusinessSegmentsHistoricalItem{
			Date:      h.Date,
			Total:     h.Total,
			Currency:  h.Currency,
			Business:  convertBusinessSegmentHistoryItems(h.Business),
			Regionals: convertBusinessSegmentHistoryItems(h.Regionals),
		})
	}
	return &BusinessSegmentsHistory{Historical: items}
}

func convertInstitutionRatingViews(j *jsontypes.InstitutionRatingViews) *InstitutionRatingViews {
	items := make([]InstitutionRatingViewItem, 0, len(j.Elist))
	for _, e := range j.Elist {
		items = append(items, InstitutionRatingViewItem{
			Date:  time.Unix(parseTimestampNumber(e.Date), 0).UTC(),
			Buy:   e.Buy,
			Over:  e.Over,
			Hold:  e.Hold,
			Under: e.Under,
			Sell:  e.Sell,
			Total: e.Total,
		})
	}
	return &InstitutionRatingViews{Elist: items}
}

func convertIndustryRankResponse(j *jsontypes.IndustryRankResponse) *IndustryRankResponse {
	groups := make([]IndustryRankGroup, 0, len(j.Items))
	for _, g := range j.Items {
		items := make([]IndustryRankItem, 0, len(g.Lists))
		for _, it := range g.Lists {
			items = append(items, IndustryRankItem{
				Name:          it.Name,
				CounterID:     it.CounterID,
				Chg:           it.Chg,
				LeadingName:   it.LeadingName,
				LeadingTicker: it.LeadingTicker,
				LeadingChg:    it.LeadingChg,
				ValueName:     it.ValueName,
				ValueData:     it.ValueData,
			})
		}
		groups = append(groups, IndustryRankGroup{Lists: items})
	}
	return &IndustryRankResponse{Items: groups}
}

func convertIndustryPeerNode(j *jsontypes.IndustryPeerNode) *IndustryPeerNode {
	if j == nil {
		return nil
	}
	next := make([]IndustryPeerNode, 0, len(j.Next))
	for i := range j.Next {
		if converted := convertIndustryPeerNode(&j.Next[i]); converted != nil {
			next = append(next, *converted)
		}
	}
	return &IndustryPeerNode{
		Name:      j.Name,
		CounterID: j.CounterID,
		StockNum:  j.StockNum,
		Chg:       j.Chg,
		YtdChg:    j.YtdChg,
		Next:      next,
	}
}

func convertIndustryPeersResponse(j *jsontypes.IndustryPeersResponse) *IndustryPeersResponse {
	return &IndustryPeersResponse{
		Top: IndustryPeersTop{
			Name:   j.Top.Name,
			Market: j.Top.Market,
		},
		Chain: convertIndustryPeerNode(j.Chain),
	}
}

func convertSnapshotForecastMetric(j *jsontypes.SnapshotForecastMetric) *SnapshotForecastMetric {
	if j == nil {
		return nil
	}
	return &SnapshotForecastMetric{
		Value:    j.Value,
		Yoy:      j.Yoy,
		CmpDesc:  j.CmpDesc,
		EstValue: j.EstValue,
	}
}

func convertSnapshotReportedMetric(j *jsontypes.SnapshotReportedMetric) *SnapshotReportedMetric {
	if j == nil {
		return nil
	}
	return &SnapshotReportedMetric{Value: j.Value, Yoy: j.Yoy}
}

func convertFinancialReportSnapshot(j *jsontypes.FinancialReportSnapshot) *FinancialReportSnapshot {
	return &FinancialReportSnapshot{
		Name:              j.Name,
		Ticker:            j.Ticker,
		FpStart:           j.FpStart,
		FpEnd:             j.FpEnd,
		Currency:          j.Currency,
		ReportDesc:        j.ReportDesc,
		FoRevenue:         convertSnapshotForecastMetric(j.FoRevenue),
		FoEbit:            convertSnapshotForecastMetric(j.FoEbit),
		FoEps:             convertSnapshotForecastMetric(j.FoEps),
		FrRevenue:         convertSnapshotReportedMetric(j.FrRevenue),
		FrProfit:          convertSnapshotReportedMetric(j.FrProfit),
		FrOperateCash:     convertSnapshotReportedMetric(j.FrOperateCash),
		FrInvestCash:      convertSnapshotReportedMetric(j.FrInvestCash),
		FrFinanceCash:     convertSnapshotReportedMetric(j.FrFinanceCash),
		FrTotalAssets:     convertSnapshotReportedMetric(j.FrTotalAssets),
		FrTotalLiability:  convertSnapshotReportedMetric(j.FrTotalLiability),
		FrRoeTtm:          j.FrRoeTtm,
		FrProfitMargin:    j.FrProfitMargin,
		FrProfitMarginTtm: j.FrProfitMarginTtm,
		FrAssetTurnTtm:    j.FrAssetTurnTtm,
		FrLeverageTtm:     j.FrLeverageTtm,
		FrDebtAssetsRatio: j.FrDebtAssetsRatio,
	}
}

// ─── Macroeconomic ────────────────────────────────────────────────────

// MacroeconomicIndicators fetches the list of available macroeconomic indicators.
//
// Pass country to filter by country code (e.g. MacroeconomicCountryUS).
// Pass nil for all countries.
//
// Path: GET /v1/quote/macrodata
func (c *FundamentalContext) MacroeconomicIndicators(
	ctx context.Context,
	country *MacroeconomicCountry,
	offset *int32,
	limit *int32,
) (*MacroeconomicIndicatorListResponse, error) {
	q := url.Values{}
	if country != nil {
		q.Set("country", macroeconomicCountryToAPIValue(*country))
	}
	if offset != nil {
		q.Set("offset", fmt.Sprintf("%d", *offset))
	}
	if limit != nil {
		q.Set("limit", fmt.Sprintf("%d", *limit))
	}
	var resp jsontypes.MacroeconomicIndicatorListResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/macrodata", q, &resp); err != nil {
		return nil, err
	}
	out := make([]MacroeconomicIndicator, 0, len(resp.Data))
	for _, item := range resp.Data {
		out = append(out, convertMacroeconomicIndicator(&item))
	}
	return &MacroeconomicIndicatorListResponse{Data: out, Count: resp.Count}, nil
}

// Macroeconomic fetches historical data for a specific macroeconomic indicator.
//
// startDate and endDate are date strings in "YYYY-MM-DD" format.
// startDate is sent as YYYY-MM-DDT00:00:00Z; endDate is sent as YYYY-MM-DDT23:59:59Z.
//
// Path: GET /v1/quote/macrodata/{indicator_code}
func (c *FundamentalContext) Macroeconomic(
	ctx context.Context,
	indicatorCode string,
	startDate *string,
	endDate *string,
	offset *int32,
	limit *int32,
) (*MacroeconomicResponse, error) {
	q := url.Values{}
	if startDate != nil {
		q.Set("start_time", *startDate+"T00:00:00Z")
	}
	if endDate != nil {
		q.Set("end_time", *endDate+"T23:59:59Z")
	}
	if offset != nil {
		q.Set("offset", fmt.Sprintf("%d", *offset))
	}
	if limit != nil {
		q.Set("limit", fmt.Sprintf("%d", *limit))
	}
	var resp jsontypes.MacroeconomicResponse
	path := "/v1/quote/macrodata/" + indicatorCode
	if err := c.httpClient.Get(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	data := make([]Macroeconomic, 0, len(resp.Data))
	for _, d := range resp.Data {
		data = append(data, convertMacroeconomic(&d))
	}
	return &MacroeconomicResponse{
		Info:  convertMacroeconomicIndicator(&resp.Info),
		Data:  data,
		Count: resp.Count,
	}, nil
}

func convertMultiLanguageText(j jsontypes.MultiLanguageText) MultiLanguageText {
	return MultiLanguageText{
		English:            j.English,
		SimplifiedChinese:  j.SimplifiedChinese,
		TraditionalChinese: j.TraditionalChinese,
	}
}

func parseOptionalRFC3339(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	t = t.UTC()
	return &t
}

func convertMacroeconomicIndicator(j *jsontypes.MacroeconomicIndicator) MacroeconomicIndicator {
	return MacroeconomicIndicator{
		IndicatorCode:    j.IndicatorCode,
		SourceOrg:        j.SourceOrg,
		Country:          j.Country,
		Name:             convertMultiLanguageText(j.Name),
		AdjustmentFactor: j.AdjustmentFactor,
		Periodicity:      j.Periodicity,
		Category:         j.Category,
		Describe:         convertMultiLanguageText(j.Describe),
		Importance:       j.Importance,
		StartDate:        parseOptionalRFC3339(j.StartDate),
	}
}

func convertMacroeconomic(j *jsontypes.Macroeconomic) Macroeconomic {
	return Macroeconomic{
		Period:        j.Period,
		ReleaseAt:     parseOptionalRFC3339(j.ReleaseAt),
		ActualValue:   j.ActualValue,
		PreviousValue: j.PreviousValue,
		ForecastValue: j.ForecastValue,
		RevisedValue:  j.RevisedValue,
		NextReleaseAt: parseOptionalRFC3339(j.NextReleaseAt),
		Unit:          convertMultiLanguageText(j.Unit),
		UnitPrefix:    convertMultiLanguageText(j.UnitPrefix),
	}
}

func macroeconomicCountryToAPIValue(c MacroeconomicCountry) string {
	switch c {
	case MacroeconomicCountryHK:
		return "Hong Kong SAR China"
	case MacroeconomicCountryCN:
		return "China (Mainland)"
	case MacroeconomicCountryUS:
		return "United States"
	case MacroeconomicCountryEU:
		return "Euro Zone"
	case MacroeconomicCountryJP:
		return "Japan"
	case MacroeconomicCountrySG:
		return "Singapore"
	default:
		return string(c)
	}
}
