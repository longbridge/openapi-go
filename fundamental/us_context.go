package fundamental

import (
	"context"
	"fmt"
	"net/url"
)

// CompanyOverview returns the US company summary for the given symbol (e.g. "AAPL.US").
//
// Path: GET /v1/us/stock-info/company-overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) CompanyOverview(ctx context.Context, symbol string) (*USCompanyOverview, error) {
	if err := c.httpClient.CheckRegion("/v1/us/stock-info/company-overview", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp USCompanyOverview
	if err := c.httpClient.Get(ctx, "/v1/us/stock-info/company-overview", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ValuationOverview returns the US valuation snapshot (PE/PB/PS) for the given symbol (e.g. "AAPL.US").
//
// Path: GET /v1/us/stock-info/valuation-overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) ValuationOverview(ctx context.Context, symbol string) (*ValuationOverview, error) {
	if err := c.httpClient.CheckRegion("/v1/us/stock-info/valuation-overview", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp ValuationOverview
	if err := c.httpClient.Get(ctx, "/v1/us/stock-info/valuation-overview", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FinancialOverview returns the US financial overview (revenue, net income, EPS, cash flow)
// for the given symbol (e.g. "AAPL.US") and report period.
//
// report: "annual" or "quarterly"
//
// Path: GET /v1/us/stock-info/finn-overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) FinancialOverview(ctx context.Context, symbol, report string) (*FinancialOverview, error) {
	if err := c.httpClient.CheckRegion("/v1/us/stock-info/finn-overview", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("report", report)
	var resp FinancialOverview
	if err := c.httpClient.Get(ctx, "/v1/us/stock-info/finn-overview", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FinancialStatement returns the US financial statement detail for the given
// symbol (e.g. "AAPL.US"), statement kind, and report period.
// Only applicable to stocks — ETFs do not have financial statements.
//
// kind: "IS" (income statement), "BS" (balance sheet), "CF" (cash flow)
// report: "q1" (Q1), "qf" (quarterly), "saf" (semi-annual), "3q" (Q3), "af" (annual)
//
// Path: GET /v1/us/quote/financials/statements
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) FinancialStatement(ctx context.Context, symbol, kind, report string) (*FinancialStatement, error) {
	if err := c.httpClient.CheckRegion("/v1/us/quote/financials/statements", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("kind", kind)
	q.Set("report", report)
	var resp FinancialStatement
	if err := c.httpClient.Get(ctx, "/v1/us/quote/financials/statements", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// KeyFinancialMetrics returns key financial ratios (ROE, gross/net margin, debt ratio) for
// the given symbol (e.g. "AAPL.US") and report period.
//
// report: "q1" (Q1), "qf" (quarterly), "saf" (semi-annual), "3q" (Q3), "af" (annual)
//
// Path: GET /v1/us/stock-info/fin-keyfactor
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) KeyFinancialMetrics(ctx context.Context, symbol, report string) (*KeyFinancialMetrics, error) {
	if err := c.httpClient.CheckRegion("/v1/us/stock-info/fin-keyfactor", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("report", report)
	var resp KeyFinancialMetrics
	if err := c.httpClient.Get(ctx, "/v1/us/stock-info/fin-keyfactor", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AnalystConsensus returns analyst consensus estimates (EPS and revenue forecasts) for
// the given symbol (e.g. "AAPL.US") and report period.
//
// report: "q1" (Q1), "qf" (quarterly), "saf" (semi-annual), "3q" (Q3), "af" (annual)
//
// Path: GET /v1/us/stock-info/fin-consensus
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) AnalystConsensus(ctx context.Context, symbol, report string) (*AnalystConsensus, error) {
	if err := c.httpClient.CheckRegion("/v1/us/stock-info/fin-consensus", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	q.Set("report", report)
	var resp AnalystConsensus
	if err := c.httpClient.Get(ctx, "/v1/us/stock-info/fin-consensus", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ETFDividendInfo returns dividend history for a US ETF (e.g. "SPY.US").
//
// Path: GET /v1/us/stock-info/etf-dividend-info
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) ETFDividendInfo(ctx context.Context, symbol string) (*ETFDividendInfo, error) {
	if err := c.httpClient.CheckRegion("/v1/us/stock-info/etf-dividend-info", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp ETFDividendInfo
	if err := c.httpClient.Get(ctx, "/v1/us/stock-info/etf-dividend-info", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CompanyDividends returns historical dividend payments for a US stock (e.g. "AAPL.US").
//
// Path: GET /v1/us/stock-info/company-dividends
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) CompanyDividends(ctx context.Context, symbol string) (*USCompanyDividends, error) {
	if err := c.httpClient.CheckRegion("/v1/us/stock-info/company-dividends", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	var resp USCompanyDividends
	if err := c.httpClient.Get(ctx, "/v1/us/stock-info/company-dividends", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ETFFiles returns the document list (prospectus, annual report, etc.) for a US ETF (e.g. "SPY.US").
//
// size: number of files to return; pass nil for all.
//
// Path: GET /v1/us/stock-info/etf-files
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) ETFFiles(ctx context.Context, symbol string, size *int32) (*ETFFilesResponse, error) {
	if err := c.httpClient.CheckRegion("/v1/us/stock-info/etf-files", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", symbolToCounterID(symbol))
	if size != nil {
		q.Set("size", fmt.Sprintf("%d", *size))
	}
	var resp ETFFilesResponse
	if err := c.httpClient.Get(ctx, "/v1/us/stock-info/etf-files", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
