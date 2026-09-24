// Package fund provides a client for the Longbridge fund (mutual fund)
// OpenAPI. It covers the fund catalog / market data (hot funds, fund list,
// filters, detail, analysis, trend, returns, performance, net value,
// holdings), the user's fund positions, and fund orders / trading.
//
// The fund identifier is exposed as symbol (e.g. a fund code); the gateway
// converts it to the backend counter_id.
package fund

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/fund/jsontypes"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
)

// FundContext is a client for the Longbridge fund (mutual fund) OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	fctx, err := fund.NewFromCfg(conf)
//	funds, err := fctx.HotFunds(ctx)
type FundContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a FundContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*FundContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &FundContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a FundContext configured from environment variables.
func NewFromEnv() (*FundContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// ----- fund catalog / market data -----

// HotFunds returns the hot-selling fund list.
func (c *FundContext) HotFunds(ctx context.Context) (funds []*HotFund, err error) {
	resp := &jsontypes.HotFundsResponse{}
	if err = c.httpClient.Get(ctx, "/v1/fund/hot-funds", url.Values{}, resp); err != nil {
		return
	}
	err = util.Copy(&funds, resp.List)
	return
}

// Funds returns the fund list.
func (c *FundContext) Funds(ctx context.Context, params *GetFunds) (funds []*FundBrief, err error) {
	body := jsontypes.GetFundsBody{}
	if params != nil {
		body.Filter = params.Filter
		body.QuickIds = params.QuickIds
		body.TimeInterval = params.TimeInterval
	}
	resp := &jsontypes.FundsResponse{}
	if err = c.httpClient.Post(ctx, "/v1/fund/funds", body, resp); err != nil {
		return
	}
	err = util.Copy(&funds, resp.Funds)
	return
}

// Filters returns the fund list filter options.
func (c *FundContext) Filters(ctx context.Context) (filters *FundFilters, err error) {
	resp := &jsontypes.FundFilters{}
	if err = c.httpClient.Get(ctx, "/v1/fund/filters", url.Values{}, resp); err != nil {
		return
	}
	filters = &FundFilters{}
	err = util.Copy(filters, resp)
	return
}

// Detail returns the fund detail.
func (c *FundContext) Detail(ctx context.Context, symbol string) (detail *FundDetail, err error) {
	resp := &jsontypes.FundDetail{}
	path := fmt.Sprintf("/v1/fund/funds/%s", symbol)
	if err = c.httpClient.Get(ctx, path, url.Values{}, resp); err != nil {
		return
	}
	detail = &FundDetail{}
	err = util.Copy(detail, resp)
	return
}

// Analysis returns the fund analysis (level 1).
func (c *FundContext) Analysis(ctx context.Context, symbol string, params *GetFundAnalysis) (analysis *FundAnalysis, err error) {
	resp := &jsontypes.FundAnalysis{}
	path := fmt.Sprintf("/v1/fund/funds/%s/analysis", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	analysis = &FundAnalysis{}
	err = util.Copy(analysis, resp)
	return
}

// AnalysisDetail returns the fund analysis detail (level 2).
func (c *FundContext) AnalysisDetail(ctx context.Context, symbol string, params *GetFundAnalysis) (detail *FundAnalysisDetail, err error) {
	resp := &jsontypes.FundAnalysisDetail{}
	path := fmt.Sprintf("/v1/fund/funds/%s/analysis/detail", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	detail = &FundAnalysisDetail{}
	err = util.Copy(detail, resp)
	return
}

// Trend returns the fund trend chart.
func (c *FundContext) Trend(ctx context.Context, symbol string, params *GetFundAnalysis) (trend *FundTrend, err error) {
	resp := &jsontypes.FundTrend{}
	path := fmt.Sprintf("/v1/fund/funds/%s/trend", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	trend = &FundTrend{}
	err = util.Copy(trend, resp)
	return
}

// AnnualReturns returns the fund annual returns.
func (c *FundContext) AnnualReturns(ctx context.Context, symbol string, params *FundPage) (returns []*FundAnnualReturn, err error) {
	resp := &jsontypes.AnnualReturnsResponse{}
	path := fmt.Sprintf("/v1/fund/funds/%s/returns/annual", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	err = util.Copy(&returns, resp.List)
	return
}

// QuarterlyReturns returns the fund quarterly returns.
func (c *FundContext) QuarterlyReturns(ctx context.Context, symbol string, params *FundPage) (returns []*FundQuarterlyReturn, err error) {
	resp := &jsontypes.QuarterlyReturnsResponse{}
	path := fmt.Sprintf("/v1/fund/funds/%s/returns/quarterly", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	err = util.Copy(&returns, resp.List)
	return
}

// Performance returns the fund performance figures.
func (c *FundContext) Performance(ctx context.Context, symbol string) (performance []*FundPerformance, err error) {
	resp := &jsontypes.PerformanceResponse{}
	path := fmt.Sprintf("/v1/fund/funds/%s/performance", symbol)
	if err = c.httpClient.Get(ctx, path, url.Values{}, resp); err != nil {
		return
	}
	err = util.Copy(&performance, resp.Value)
	return
}

// PerformanceComparison returns the fund performance comparison.
func (c *FundContext) PerformanceComparison(ctx context.Context, symbol string, params *GetFundAnalysis) (comparison *FundPerformanceComparison, err error) {
	resp := &jsontypes.FundPerformanceComparison{}
	path := fmt.Sprintf("/v1/fund/funds/%s/performance/comparison", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	comparison = &FundPerformanceComparison{}
	err = util.Copy(comparison, resp)
	return
}

// Nav returns the fund latest net value.
func (c *FundContext) Nav(ctx context.Context, symbol string) (nav []*FundNavValue, err error) {
	resp := &jsontypes.NavResponse{}
	path := fmt.Sprintf("/v1/fund/funds/%s/nav", symbol)
	if err = c.httpClient.Get(ctx, path, url.Values{}, resp); err != nil {
		return
	}
	err = util.Copy(&nav, resp.Value)
	return
}

// NavHistory returns the fund historical net value (paged).
func (c *FundContext) NavHistory(ctx context.Context, symbol string, params *FundPage) (nav []*FundNavValue, err error) {
	resp := &jsontypes.NavHistoryResponse{}
	path := fmt.Sprintf("/v1/fund/funds/%s/nav-history", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	err = util.Copy(&nav, resp.HistoryValue)
	return
}

// NavRange returns the fund historical net value by relative time range.
func (c *FundContext) NavRange(ctx context.Context, symbol string, params *FundNavRange) (nav []*FundNavValue, err error) {
	resp := &jsontypes.NavHistoryResponse{}
	path := fmt.Sprintf("/v1/fund/funds/%s/nav-range", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	err = util.Copy(&nav, resp.HistoryValue)
	return
}

// Holdings returns a fund's top-10 holdings.
func (c *FundContext) Holdings(ctx context.Context, symbol string, params *GetFundHoldings) (holdings *FundHoldings, err error) {
	resp := &jsontypes.FundHoldings{}
	path := fmt.Sprintf("/v1/fund/funds/%s/holdings", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	holdings = &FundHoldings{}
	err = util.Copy(holdings, resp)
	return
}

// StockHoldings returns the stocks held by a fund (reverse lookup).
func (c *FundContext) StockHoldings(ctx context.Context, symbol string, params *GetFundStockHoldings) (holdings []*FundStockHolding, err error) {
	resp := &jsontypes.StockHoldingsResponse{}
	path := fmt.Sprintf("/v1/fund/funds/%s/stock-holdings", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	err = util.Copy(&holdings, resp.Lists)
	return
}

// ----- user fund positions -----

// Positions returns the user's fund positions overview.
func (c *FundContext) Positions(ctx context.Context, params *GetFundPositions) (positions *FundPositions, err error) {
	resp := &jsontypes.FundPositions{}
	if err = c.httpClient.Get(ctx, "/v1/asset/funds", params.Values(), resp); err != nil {
		return
	}
	positions = &FundPositions{}
	err = util.Copy(positions, resp)
	return
}

// Position returns the user's single fund position detail.
func (c *FundContext) Position(ctx context.Context, symbol string, params *GetFundPosition) (detail *FundPositionDetail, err error) {
	resp := &jsontypes.FundPositionDetail{}
	path := fmt.Sprintf("/v1/asset/funds/%s", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	detail = &FundPositionDetail{}
	err = util.Copy(detail, resp)
	return
}

// PositionPerformance returns the performance figures of a held fund.
func (c *FundContext) PositionPerformance(ctx context.Context, symbol string) (performance []*FundPositionPerformance, err error) {
	resp := &jsontypes.PositionPerformanceResponse{}
	path := fmt.Sprintf("/v1/asset/funds/%s/performance", symbol)
	if err = c.httpClient.Get(ctx, path, url.Values{}, resp); err != nil {
		return
	}
	err = util.Copy(&performance, resp.Value)
	return
}

// PositionProfits returns the cumulative-profit series of a held fund.
func (c *FundContext) PositionProfits(ctx context.Context, symbol string, params *GetFundPositionProfits) (profits *FundPositionProfits, err error) {
	resp := &jsontypes.FundPositionProfits{}
	path := fmt.Sprintf("/v1/asset/funds/%s/profits", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	profits = &FundPositionProfits{}
	err = util.Copy(profits, resp)
	return
}

// PositionNav returns the net-value history of a held fund.
func (c *FundContext) PositionNav(ctx context.Context, symbol string, params *FundNavRange) (nav []*FundPositionNav, err error) {
	resp := &jsontypes.PositionNavResponse{}
	path := fmt.Sprintf("/v1/asset/funds/%s/nav-history", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	err = util.Copy(&nav, resp.HistoryValue)
	return
}

// PositionDividends returns the dividend records of a held fund.
func (c *FundContext) PositionDividends(ctx context.Context, symbol string, params *GetFundPositionDividends) (dividends *FundDividends, err error) {
	resp := &jsontypes.FundDividends{}
	path := fmt.Sprintf("/v1/asset/funds/%s/dividends", symbol)
	if err = c.httpClient.Get(ctx, path, params.Values(), resp); err != nil {
		return
	}
	dividends = &FundDividends{}
	err = util.Copy(dividends, resp)
	return
}

// ----- fund orders & trading -----

// Orders returns the user's fund orders (also serves as the trade / execution
// record).
func (c *FundContext) Orders(ctx context.Context, params *GetFundOrders) (orders []*FundOrder, err error) {
	resp := &jsontypes.OrdersResponse{}
	if err = c.httpClient.Get(ctx, "/v1/fund/orders", params.Values(), resp); err != nil {
		return
	}
	err = util.Copy(&orders, resp.Orders)
	return
}

// Order returns a fund order detail.
func (c *FundContext) Order(ctx context.Context, orderId int64) (detail *FundOrderDetail, err error) {
	resp := &jsontypes.FundOrderDetail{}
	path := fmt.Sprintf("/v1/fund/orders/%d", orderId)
	if err = c.httpClient.Get(ctx, path, url.Values{}, resp); err != nil {
		return
	}
	detail = &FundOrderDetail{}
	err = util.Copy(detail, resp)
	return
}

// Transactions returns the user's fund transactions (cash-flow records).
func (c *FundContext) Transactions(ctx context.Context, params *GetFundTransactions) (transactions []*FundTransaction, err error) {
	resp := &jsontypes.TransactionsResponse{}
	if err = c.httpClient.Get(ctx, "/v1/fund/transactions", params.Values(), resp); err != nil {
		return
	}
	err = util.Copy(&transactions, resp.List)
	return
}

// ValidateOrder validates a fund order before submitting.
func (c *FundContext) ValidateOrder(ctx context.Context, params *ValidateFundOrder) (validation *FundOrderValidation, err error) {
	body := jsontypes.ValidateFundOrderBody{
		Symbol:         params.Symbol,
		Action:         params.Action,
		Currency:       params.Currency,
		Amount:         params.Amount,
		Units:          params.Units,
		DividendOption: params.DividendOption,
		FundSource:     params.FundSource,
		AccountChannel: params.AccountChannel,
	}
	resp := &jsontypes.FundOrderValidation{}
	if err = c.httpClient.Post(ctx, "/v1/fund/orders/validate", body, resp); err != nil {
		return
	}
	validation = &FundOrderValidation{}
	err = util.Copy(validation, resp)
	return
}

// SubmitOrder submits a fund order (buy / sell).
func (c *FundContext) SubmitOrder(ctx context.Context, params *SubmitFundOrder) (order *FundOrderSubmitResponse, err error) {
	body := jsontypes.SubmitFundOrderBody{
		Symbol:         params.Symbol,
		Action:         params.Action,
		Currency:       params.Currency,
		Amount:         params.Amount,
		Units:          params.Units,
		DividendOption: params.DividendOption,
		Fee:            params.Fee,
		IsSellAll:      params.IsSellAll,
		Remark:         params.Remark,
		TradeMethod:    params.TradeMethod,
	}
	resp := &jsontypes.FundOrderSubmitResponse{}
	if err = c.httpClient.Post(ctx, "/v1/fund/orders", body, resp); err != nil {
		return
	}
	order = &FundOrderSubmitResponse{}
	err = util.Copy(order, resp)
	return
}

// CancelOrder cancels (withdraws) a fund order.
func (c *FundContext) CancelOrder(ctx context.Context, orderId int64) (err error) {
	body := jsontypes.CancelFundOrderBody{UtId: orderId}
	path := fmt.Sprintf("/v1/fund/orders/%d/cancel", orderId)
	return c.httpClient.Post(ctx, path, body, nil)
}
