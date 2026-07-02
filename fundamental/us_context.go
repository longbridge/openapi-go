package fundamental

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

// ErrUSOnly is kept as a sentinel for callers using errors.Is; the actual
// returned error is *http.RegionRestrictedError with the same meaning.
var ErrUSOnly = errors.New("longbridge: this API is only available for US accounts")

// CompanyOverview returns the US company summary for the given counter_id.
//
// counterID format: "ST/US/AAPL"
//
// Path: GET /v1/stock-info/company-overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) CompanyOverview(ctx context.Context, counterID string) (*USCompanyOverview, error) {
	if err := c.httpClient.CheckRegion("/v1/stock-info/company-overview", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	var resp USCompanyOverview
	if err := c.httpClient.Get(ctx, "/v1/stock-info/company-overview", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ValuationOverview returns the US valuation snapshot (PE/PB/PS) for the given counter_id.
//
// Path: GET /v1/stock-info/valuation-overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) ValuationOverview(ctx context.Context, counterID string) (*ValuationOverview, error) {
	if err := c.httpClient.CheckRegion("/v1/stock-info/valuation-overview", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	var resp ValuationOverview
	if err := c.httpClient.Get(ctx, "/v1/stock-info/valuation-overview", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FinancialOverview returns the US financial overview (revenue, net income, EPS, cash flow)
// for the given counter_id and report period.
//
// report: "annual" or "quarterly"
//
// Path: GET /v1/stock-info/finn-overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) FinancialOverview(ctx context.Context, counterID, report string) (FinancialOverview, error) {
	if err := c.httpClient.CheckRegion("/v1/stock-info/finn-overview", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	q.Set("report", report)
	var resp FinancialOverview
	if err := c.httpClient.Get(ctx, "/v1/stock-info/finn-overview", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// FinancialStatementV3 returns the US financial statement detail (IS/BS/CF) for a given
// counter_id, statement kind, and report period.
//
// kind: "IS" (income statement), "BS" (balance sheet), "CF" (cash flow)
// report: "annual" or "quarterly"
//
// Path: GET /v1/us/quote/financials/statements
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) FinancialStatementV3(ctx context.Context, counterID, kind, report string) (*FinancialStatement, error) {
	if err := c.httpClient.CheckRegion("/v1/us/quote/financials/statements", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	q.Set("kind", kind)
	q.Set("report", report)
	var resp FinancialStatement
	if err := c.httpClient.Get(ctx, "/v1/us/quote/financials/statements", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// KeyFinancialMetrics returns key financial ratios (ROE, gross/net margin, debt ratio) for
// the given counter_id and report period.
//
// report: "annual" or "quarterly"
//
// Path: GET /v1/stock-info/fin-keyfactor
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) KeyFinancialMetrics(ctx context.Context, counterID, report string) (KeyFinancialMetrics, error) {
	if err := c.httpClient.CheckRegion("/v1/stock-info/fin-keyfactor", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	q.Set("report", report)
	var resp KeyFinancialMetrics
	if err := c.httpClient.Get(ctx, "/v1/stock-info/fin-keyfactor", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AnalystConsensus returns analyst consensus estimates (EPS and revenue forecasts) for
// the given counter_id and report period.
//
// report: "annual" or "quarterly"
//
// Path: GET /v1/stock-info/fin-consensus
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) AnalystConsensus(ctx context.Context, counterID, report string) (AnalystConsensus, error) {
	if err := c.httpClient.CheckRegion("/v1/stock-info/fin-consensus", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	q.Set("report", report)
	var resp AnalystConsensus
	if err := c.httpClient.Get(ctx, "/v1/stock-info/fin-consensus", q, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ETFDividendInfo returns dividend history for a US ETF.
//
// counterID format: "ST/US/SPY"
//
// Path: GET /v1/stock-info/etf-dividend-info
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) ETFDividendInfo(ctx context.Context, counterID string) (*ETFDividendInfo, error) {
	if err := c.httpClient.CheckRegion("/v1/stock-info/etf-dividend-info", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	var resp ETFDividendInfo
	if err := c.httpClient.Get(ctx, "/v1/stock-info/etf-dividend-info", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CompanyDividends returns historical dividend payments for a US stock.
//
// counterID format: "ST/US/AAPL"
//
// Path: GET /v1/stock-info/company-dividends
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) CompanyDividends(ctx context.Context, counterID string) (*USCompanyDividends, error) {
	if err := c.httpClient.CheckRegion("/v1/stock-info/company-dividends", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	var resp USCompanyDividends
	if err := c.httpClient.Get(ctx, "/v1/stock-info/company-dividends", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ETFFiles returns the document list (prospectus, annual report, etc.) for a US ETF.
//
// counterID format: "ST/US/SPY"
// size: number of files to return; pass nil for all (defaults to 0 = all on the server).
//
// Path: GET /v1/stock-info/etf-files
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *FundamentalContext) ETFFiles(ctx context.Context, counterID string, size *int32) (*ETFFilesResponse, error) {
	if err := c.httpClient.CheckRegion("/v1/stock-info/etf-files", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("counter_id", counterID)
	if size != nil {
		q.Set("size", fmt.Sprintf("%d", *size))
	}
	var resp ETFFilesResponse
	if err := c.httpClient.Get(ctx, "/v1/stock-info/etf-files", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
