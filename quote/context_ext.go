package quote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// symbolToCounterID converts a symbol like "700.HK" to counter-id format "ST/HK/700".
func symbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])
	return "ST/" + market + "/" + code
}

// FinancialReport returns financial report data for a symbol.
//
// Reference: GET /v1/quote/financial-reports
func (c *QuoteContext) FinancialReport(ctx context.Context, symbol string, kind, reportType *string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	if kind != nil {
		values.Add("kind", *kind)
	}
	if reportType != nil {
		values.Add("report", *reportType)
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/financial-reports", values, &resp)
	return resp, err
}

// InstitutionRatings returns institution ratings for a symbol.
//
// Reference: GET /v1/quote/institution-ratings
func (c *QuoteContext) InstitutionRatings(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/institution-ratings", values, &resp)
	return resp, err
}

// InstitutionRatingLatest returns the latest institution rating for a symbol.
//
// Reference: GET /v1/quote/institution-rating-latest
func (c *QuoteContext) InstitutionRatingLatest(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/institution-rating-latest", values, &resp)
	return resp, err
}

// InstitutionRatingDetail returns institution rating detail for a symbol.
//
// Reference: GET /v1/quote/institution-ratings/detail
func (c *QuoteContext) InstitutionRatingDetail(ctx context.Context, symbol string, page, pageSize *int) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	if page != nil {
		values.Add("page", fmt.Sprintf("%d", *page))
	}
	if pageSize != nil {
		values.Add("page_size", fmt.Sprintf("%d", *pageSize))
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/institution-ratings/detail", values, &resp)
	return resp, err
}

// Dividends returns dividend data for a symbol.
//
// Reference: GET /v1/quote/dividends
func (c *QuoteContext) Dividends(ctx context.Context, symbol string, startDate, endDate *string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	if startDate != nil {
		values.Add("start_date", *startDate)
	}
	if endDate != nil {
		values.Add("end_date", *endDate)
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/dividends", values, &resp)
	return resp, err
}

// DividendDetail returns detail for a specific dividend of a symbol.
//
// Reference: GET /v1/quote/dividends/details
func (c *QuoteContext) DividendDetail(ctx context.Context, symbol, dividendID string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	values.Add("dividend_id", dividendID)
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/dividends/details", values, &resp)
	return resp, err
}

// ForecastEPS returns EPS forecast data for a symbol.
//
// Reference: GET /v1/quote/forecast-eps
func (c *QuoteContext) ForecastEPS(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/forecast-eps", values, &resp)
	return resp, err
}

// FinancialConsensus returns financial consensus detail for a symbol.
//
// Reference: GET /v1/quote/financial-consensus-detail
func (c *QuoteContext) FinancialConsensus(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/financial-consensus-detail", values, &resp)
	return resp, err
}

// Valuation returns valuation data for a symbol.
//
// Reference: GET /v1/quote/valuation
func (c *QuoteContext) Valuation(ctx context.Context, symbol string, period *string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	if period != nil {
		values.Add("period", *period)
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/valuation", values, &resp)
	return resp, err
}

// ValuationHistory returns valuation history for a symbol.
//
// Reference: GET /v1/quote/valuation/detail
func (c *QuoteContext) ValuationHistory(ctx context.Context, symbol string, period *string, count *int) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	if period != nil {
		values.Add("period", *period)
	}
	if count != nil {
		values.Add("count", fmt.Sprintf("%d", *count))
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/valuation/detail", values, &resp)
	return resp, err
}

// IndustryValuation returns industry valuation comparison for a symbol.
//
// Reference: GET /v1/quote/industry-valuation-comparison
func (c *QuoteContext) IndustryValuation(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/industry-valuation-comparison", values, &resp)
	return resp, err
}

// IndustryValuationDistribution returns industry valuation distribution for a symbol.
//
// Reference: GET /v1/quote/industry-valuation-distribution
func (c *QuoteContext) IndustryValuationDistribution(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/industry-valuation-distribution", values, &resp)
	return resp, err
}

// CompanyOverview returns company overview for a symbol.
//
// Reference: GET /v1/quote/comp-overview
func (c *QuoteContext) CompanyOverview(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/comp-overview", values, &resp)
	return resp, err
}

// CompanyExecutives returns company executives for a symbol.
//
// Reference: GET /v1/quote/company-professionals
func (c *QuoteContext) CompanyExecutives(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/company-professionals", values, &resp)
	return resp, err
}

// Shareholders returns shareholders data for a symbol.
//
// Reference: GET /v1/quote/shareholders
func (c *QuoteContext) Shareholders(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/shareholders", values, &resp)
	return resp, err
}

// FundHolders returns fund holders for a symbol.
//
// Reference: GET /v1/quote/fund-holders
func (c *QuoteContext) FundHolders(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/fund-holders", values, &resp)
	return resp, err
}

// CorporateActions returns corporate actions for a symbol.
//
// Reference: GET /v1/quote/company-act
func (c *QuoteContext) CorporateActions(ctx context.Context, symbol string, actionType *string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	if actionType != nil {
		values.Add("action_type", *actionType)
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/company-act", values, &resp)
	return resp, err
}

// InvestorRelations returns investor relations data for a symbol.
//
// Reference: GET /v1/quote/invest-relations
func (c *QuoteContext) InvestorRelations(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/invest-relations", values, &resp)
	return resp, err
}

// OperatingData returns operating data for a symbol.
//
// Reference: GET /v1/quote/operatings
func (c *QuoteContext) OperatingData(ctx context.Context, symbol string, period *string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("counter_id", symbolToCounterID(symbol))
	if period != nil {
		values.Add("period", *period)
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/operatings", values, &resp)
	return resp, err
}

// MarketStatus returns market status for a market.
//
// Reference: GET /v1/quote/market-status
func (c *QuoteContext) MarketStatus(ctx context.Context, market string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("market", market)
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/market-status", values, &resp)
	return resp, err
}

// BrokerHolding returns broker holding data for a symbol.
//
// Reference: GET /v1/quote/broker-holding
func (c *QuoteContext) BrokerHolding(ctx context.Context, symbol string, period *string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	if period != nil {
		values.Add("period", *period)
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/broker-holding", values, &resp)
	return resp, err
}

// BrokerHoldingDetail returns broker holding detail for a symbol.
//
// Reference: GET /v1/quote/broker-holding/detail
func (c *QuoteContext) BrokerHoldingDetail(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/broker-holding/detail", values, &resp)
	return resp, err
}

// BrokerHoldingDaily returns daily broker holding for a symbol and broker.
//
// Reference: GET /v1/quote/broker-holding/daily
func (c *QuoteContext) BrokerHoldingDaily(ctx context.Context, symbol, brokerID string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	values.Add("broker_id", brokerID)
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/broker-holding/daily", values, &resp)
	return resp, err
}

// AHPremiumKlines returns AH premium klines for a symbol.
//
// Reference: GET /v1/quote/ahpremium/klines
func (c *QuoteContext) AHPremiumKlines(ctx context.Context, symbol string, period *string, count *int) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	if period != nil {
		values.Add("period", *period)
	}
	if count != nil {
		values.Add("count", fmt.Sprintf("%d", *count))
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/ahpremium/klines", values, &resp)
	return resp, err
}

// AHPremiumTimeshares returns AH premium timeshares for a symbol.
//
// Reference: GET /v1/quote/ahpremium/timeshares
func (c *QuoteContext) AHPremiumTimeshares(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/ahpremium/timeshares", values, &resp)
	return resp, err
}

// TradeStatistics returns trade statistics for a symbol.
//
// Reference: GET /v1/quote/trades-statistics
func (c *QuoteContext) TradeStatistics(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/trades-statistics", values, &resp)
	return resp, err
}

// MarketAnomaly returns market anomaly data for a market.
//
// Reference: GET /v1/quote/changes
func (c *QuoteContext) MarketAnomaly(ctx context.Context, market string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("market", market)
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/changes", values, &resp)
	return resp, err
}

// IndexConstituents returns constituents for an index symbol.
//
// Reference: GET /v1/quote/index-constituents
func (c *QuoteContext) IndexConstituents(ctx context.Context, symbol string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/index-constituents", values, &resp)
	return resp, err
}

// FinanceCalendar returns finance calendar events for a market.
//
// Reference: GET /v1/quote/finance_calendar
func (c *QuoteContext) FinanceCalendar(ctx context.Context, market string, startDate, endDate *string) (json.RawMessage, error) {
	values := url.Values{}
	values.Add("market", market)
	if startDate != nil {
		values.Add("start_date", *startDate)
	}
	if endDate != nil {
		values.Add("end_date", *endDate)
	}
	var resp json.RawMessage
	err := c.opts.httpClient.Get(ctx, "/v1/quote/finance_calendar", values, &resp)
	return resp, err
}
