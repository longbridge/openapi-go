package portfolio

import (
	"context"
	"net/url"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/portfolio/jsontypes"
)

// PortfolioContext is a client for portfolio analytics (exchange rates, P&L analysis).
type PortfolioContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a PortfolioContext from a Config.
func NewFromCfg(cfg *config.Config) (*PortfolioContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &PortfolioContext{httpClient: httpClient}, nil
}

// NewFromEnv creates a PortfolioContext from environment variables.
func NewFromEnv() (*PortfolioContext, error) {
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

// ExchangeRate returns exchange rates for supported currencies.
func (c *PortfolioContext) ExchangeRate(ctx context.Context) (*ExchangeRates, error) {
	var resp jsontypes.ExchangeRates
	if err := c.httpClient.Get(ctx, "/v1/asset/exchange_rates", nil, &resp); err != nil {
		return nil, err
	}
	result := &ExchangeRates{}
	for _, e := range resp.Exchanges {
		result.Exchanges = append(result.Exchanges, &ExchangeRate{
			AverageRate:   e.AverageRate,
			BaseCurrency:  e.BaseCurrency,
			BidRate:       e.BidRate,
			OfferRate:     e.OfferRate,
			OtherCurrency: e.OtherCurrency,
		})
	}
	return result, nil
}

// ProfitAnalysis returns the portfolio P&L analysis. start and end are optional
// date strings in YYYY-MM-DD format; pass empty string to omit.
func (c *PortfolioContext) ProfitAnalysis(ctx context.Context, start, end string) (*ProfitAnalysis, error) {
	var (
		summaryResp jsontypes.ProfitAnalysisSummary
		sublistResp jsontypes.ProfitAnalysisSublist
		summaryErr  error
		sublistErr  error
		wg          sync.WaitGroup
	)

	summaryValues := url.Values{}
	if start != "" {
		summaryValues.Add("start_date", start)
	}
	if end != "" {
		summaryValues.Add("end_date", end)
	}

	sublistValues := url.Values{}
	sublistValues.Add("profit_or_loss", "all")
	if start != "" {
		sublistValues.Add("start_date", start)
	}
	if end != "" {
		sublistValues.Add("end_date", end)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		summaryErr = c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis-summary", summaryValues, &summaryResp)
	}()
	go func() {
		defer wg.Done()
		sublistErr = c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis-sublist", sublistValues, &sublistResp)
	}()
	wg.Wait()

	if summaryErr != nil {
		return nil, summaryErr
	}
	if sublistErr != nil {
		return nil, sublistErr
	}

	result := &ProfitAnalysis{
		Summary: convertSummary(&summaryResp),
		Sublist: convertSublist(&sublistResp),
	}
	return result, nil
}

func convertSummary(resp *jsontypes.ProfitAnalysisSummary) *ProfitAnalysisSummary {
	s := &ProfitAnalysisSummary{
		Currency:          resp.Currency,
		CurrentTotalAsset: parseDecimalOpt(resp.CurrentTotalAsset),
		StartDate:         resp.StartDate,
		EndDate:           resp.EndDate,
		StartTime:         resp.StartTime,
		EndTime:           resp.EndTime,
		EndingAssetValue:  parseDecimalOpt(resp.EndingAssetValue),
		InitialAssetValue: parseDecimalOpt(resp.InitialAssetValue),
		InvestAmount:      parseDecimalOpt(resp.InvestAmount),
		IsTraded:          resp.IsTraded,
		SumProfit:         parseDecimalOpt(resp.SumProfit),
		SumProfitRate:     parseDecimalOpt(resp.SumProfitRate),
	}
	if resp.Profits != nil {
		s.Profits = &ProfitSummaryBreakdown{
			Stock:                       parseDecimalOpt(resp.Profits.Stock),
			Fund:                        parseDecimalOpt(resp.Profits.Fund),
			Crypto:                      parseDecimalOpt(resp.Profits.Crypto),
			Mmf:                         parseDecimalOpt(resp.Profits.Mmf),
			Other:                       parseDecimalOpt(resp.Profits.Other),
			CumulativeTransactionAmount: parseDecimalOpt(resp.Profits.CumulativeTransactionAmount),
			TradeOrderNum:               resp.Profits.TradeOrderNum,
			TradeStockNum:               resp.Profits.TradeStockNum,
			Ipo:                         parseDecimalOpt(resp.Profits.Ipo),
			IpoHit:                      resp.Profits.IpoHit,
			IpoSubscription:             resp.Profits.IpoSubscription,
		}
		for _, si := range resp.Profits.SummaryInfo {
			s.Profits.SummaryInfo = append(s.Profits.SummaryInfo, &ProfitSummaryInfo{
				AssetType:     si.AssetType,
				ProfitMax:     si.ProfitMax,
				ProfitMaxName: si.ProfitMaxName,
				LossMax:       si.LossMax,
				LossMaxName:   si.LossMaxName,
			})
		}
	}
	return s
}

func convertSublist(resp *jsontypes.ProfitAnalysisSublist) *ProfitAnalysisSublist {
	s := &ProfitAnalysisSublist{
		Start:       resp.Start,
		End:         resp.End,
		StartDate:   resp.StartDate,
		EndDate:     resp.EndDate,
		UpdatedAt:   resp.UpdatedAt,
		UpdatedDate: resp.UpdatedDate,
	}
	for _, item := range resp.Items {
		s.Items = append(s.Items, &ProfitAnalysisItem{
			Name:              item.Name,
			Market:            item.Market,
			IsHolding:         item.IsHolding,
			Profit:            parseDecimalOpt(item.Profit),
			ProfitRate:        parseDecimalOpt(item.ProfitRate),
			ClearanceTimes:    item.ClearanceTimes,
			ItemType:          item.ItemType,
			Currency:          item.Currency,
			Symbol:            util.CounterIDToSymbol(item.CounterId),
			HoldingPeriod:     item.HoldingPeriod,
			SecurityCode:      item.SecurityCode,
			Isin:              item.Isin,
			UnderlyingProfit:  parseDecimalOpt(item.UnderlyingProfit),
			DerivativesProfit: parseDecimalOpt(item.DerivativesProfit),
			OrderProfit:       parseDecimalOpt(item.OrderProfit),
		})
	}
	return s
}

// ProfitAnalysisByMarket returns paginated P&L analysis filtered by market.
func (c *PortfolioContext) ProfitAnalysisByMarket(ctx context.Context, market, start, end, currency string, page, size int32) (*ProfitAnalysisByMarket, error) {
	values := url.Values{}
	values.Add("page", strconv.FormatInt(int64(page), 10))
	values.Add("size", strconv.FormatInt(int64(size), 10))
	if market != "" {
		values.Add("market", market)
	}
	if start != "" {
		values.Add("start_date", start)
	}
	if end != "" {
		values.Add("end_date", end)
	}
	if currency != "" {
		values.Add("currency", currency)
	}
	var resp jsontypes.ProfitAnalysisByMarket
	if err := c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis/by-market", values, &resp); err != nil {
		return nil, err
	}
	result := &ProfitAnalysisByMarket{
		Profit:  parseDecimalOpt(resp.Profit),
		HasMore: resp.HasMore,
	}
	for _, item := range resp.StockItems {
		result.StockItems = append(result.StockItems, &ProfitAnalysisByMarketItem{
			Code: item.Code, Name: item.Name,
			Market: item.Market, Profit: parseDecimalOpt(item.Profit),
		})
	}
	return result, nil
}

// ProfitAnalysisDetail returns P&L detail for a specific security.
func (c *PortfolioContext) ProfitAnalysisDetail(ctx context.Context, symbol, start, end string) (*ProfitAnalysisDetail, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	if start != "" {
		values.Add("start_date", start)
	}
	if end != "" {
		values.Add("end_date", end)
	}
	var resp jsontypes.ProfitAnalysisDetail
	if err := c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis/detail", values, &resp); err != nil {
		return nil, err
	}
	result := &ProfitAnalysisDetail{
		Profit:      parseDecimalOpt(resp.Profit),
		Name:        resp.Name,
		UpdatedAt:   resp.UpdatedAt,
		UpdatedDate: resp.UpdatedDate,
		Currency:    resp.Currency,
		DefaultTag:  resp.DefaultTag,
		Start:       resp.Start,
		End:         resp.End,
		StartDate:   resp.StartDate,
		EndDate:     resp.EndDate,
	}
	if resp.UnderlyingDetails != nil {
		result.UnderlyingDetails = convertProfitDetails(resp.UnderlyingDetails)
	}
	if resp.DerivativePnlDetails != nil {
		result.DerivativePnlDetails = convertProfitDetails(resp.DerivativePnlDetails)
	}
	return result, nil
}

func convertProfitDetails(d *jsontypes.ProfitDetails) *ProfitDetails {
	pd := &ProfitDetails{
		HoldingValue:             parseDecimalOpt(d.HoldingValue),
		Profit:                   parseDecimalOpt(d.Profit),
		CumulativeCreditedAmount: parseDecimalOpt(d.CumulativeCreditedAmount),
		CumulativeDebitedAmount:  parseDecimalOpt(d.CumulativeDebitedAmount),
		CumulativeFeeAmount:      parseDecimalOpt(d.CumulativeFeeAmount),
		ShortHoldingValue:        parseDecimalOpt(d.ShortHoldingValue),
		LongHoldingValue:         parseDecimalOpt(d.LongHoldingValue),
		HoldingValueAtBeginning:  parseDecimalOpt(d.HoldingValueAtBeginning),
		HoldingValueAtEnding:     parseDecimalOpt(d.HoldingValueAtEnding),
	}
	for _, e := range d.CreditedDetails {
		pd.CreditedDetails = append(pd.CreditedDetails, &ProfitDetailEntry{Describe: e.Describe, Amount: parseDecimalOpt(e.Amount)})
	}
	for _, e := range d.DebitedDetails {
		pd.DebitedDetails = append(pd.DebitedDetails, &ProfitDetailEntry{Describe: e.Describe, Amount: parseDecimalOpt(e.Amount)})
	}
	for _, e := range d.FeeDetails {
		pd.FeeDetails = append(pd.FeeDetails, &ProfitDetailEntry{Describe: e.Describe, Amount: parseDecimalOpt(e.Amount)})
	}
	return pd
}

// ProfitAnalysisFlows returns paginated P&L flow records for a security.
func (c *PortfolioContext) ProfitAnalysisFlows(ctx context.Context, symbol string, page, size int32, derivative bool, start, end string) (*ProfitAnalysisFlows, error) {
	values := url.Values{}
	values.Add("counter_id", util.SymbolToCounterID(symbol))
	values.Add("page", strconv.FormatInt(int64(page), 10))
	values.Add("size", strconv.FormatInt(int64(size), 10))
	if derivative {
		values.Add("derivative", "true")
	} else {
		values.Add("derivative", "false")
	}
	if start != "" {
		values.Add("start_date", start)
	}
	if end != "" {
		values.Add("end_date", end)
	}
	var resp jsontypes.ProfitAnalysisFlows
	if err := c.httpClient.Get(ctx, "/v1/portfolio/profit-analysis/flows", values, &resp); err != nil {
		return nil, err
	}
	result := &ProfitAnalysisFlows{HasMore: resp.HasMore}
	for _, f := range resp.FlowsList {
		result.FlowsList = append(result.FlowsList, &FlowItem{
			ExecutedDate:      f.ExecutedDate,
			ExecutedTimestamp: f.ExecutedTimestamp,
			Code:              f.Code,
			Direction:         f.Direction,
			ExecutedQuantity:  parseDecimalOpt(f.ExecutedQuantity),
			ExecutedPrice:     parseDecimalOpt(f.ExecutedPrice),
			ExecutedCost:      parseDecimalOpt(f.ExecutedCost),
			Describe:          f.Describe,
		})
	}
	return result, nil
}
