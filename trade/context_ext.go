package trade

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ProfitAnalysisSummary returns profit analysis summary.
//
// Reference: GET /v1/portfolio/profit-analysis-summary
func (c *TradeContext) ProfitAnalysisSummary(ctx context.Context, currency, startDate, endDate *string) (json.RawMessage, error) {
	values := url.Values{}
	if currency != nil {
		values.Add("currency", *currency)
	}
	if startDate != nil {
		values.Add("start_date", *startDate)
	}
	if endDate != nil {
		values.Add("end_date", *endDate)
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/portfolio/profit-analysis-summary", values, &resp)
	return resp, err
}

// ProfitAnalysisSublist returns profit analysis sub-list.
//
// Reference: GET /v1/portfolio/profit-analysis-sublist
func (c *TradeContext) ProfitAnalysisSublist(ctx context.Context, currency, startDate, endDate *string, page, pageSize *int) (json.RawMessage, error) {
	values := url.Values{}
	if currency != nil {
		values.Add("currency", *currency)
	}
	if startDate != nil {
		values.Add("start_date", *startDate)
	}
	if endDate != nil {
		values.Add("end_date", *endDate)
	}
	if page != nil {
		values.Add("page", fmt.Sprintf("%d", *page))
	}
	if pageSize != nil {
		values.Add("page_size", fmt.Sprintf("%d", *pageSize))
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/portfolio/profit-analysis-sublist", values, &resp)
	return resp, err
}

// ProfitAnalysisDetail returns profit analysis detail for a symbol.
//
// Reference: GET /v1/portfolio/profit-analysis/detail
func (c *TradeContext) ProfitAnalysisDetail(ctx context.Context, symbol string, currency, startDate, endDate *string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	if currency != nil {
		values.Add("currency", *currency)
	}
	if startDate != nil {
		values.Add("start_date", *startDate)
	}
	if endDate != nil {
		values.Add("end_date", *endDate)
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/portfolio/profit-analysis/detail", values, &resp)
	return resp, err
}
