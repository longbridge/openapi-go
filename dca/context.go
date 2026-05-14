package dca

import (
	"context"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/dca/jsontypes"
)

// DCAContext is a client for dollar-cost averaging (recurring investment) plan management.
type DCAContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a DCAContext from a Config.
func NewFromCfg(cfg *config.Config) (*DCAContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &DCAContext{httpClient: httpClient}, nil
}

// NewFromEnv creates a DCAContext from environment variables.
func NewFromEnv() (*DCAContext, error) {
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

func convertPlan(p *jsontypes.DcaPlan) *DcaPlan {
	return &DcaPlan{
		PlanId:             p.PlanId,
		Status:             p.Status,
		Symbol:             util.CounterIDToSymbol(p.CounterId),
		MemberId:           p.MemberId,
		Aaid:               p.Aaid,
		AccountChannel:     p.AccountChannel,
		DisplayAccount:     p.DisplayAccount,
		Market:             p.Market,
		PerInvestAmount:    parseDecimalOpt(p.PerInvestAmount),
		InvestFrequency:    p.InvestFrequency,
		InvestDayOfWeek:    p.InvestDayOfWeek,
		InvestDayOfMonth:   p.InvestDayOfMonth,
		AllowMarginFinance: p.AllowMarginFinance,
		AlterHours:         p.AlterHours,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
		NextTrdDate:        p.NextTrdDate,
		StockName:          p.StockName,
		CumAmount:          parseDecimalOpt(p.CumAmount),
		IssueNumber:        p.IssueNumber,
		AverageCost:        parseDecimalOpt(p.AverageCost),
		CumProfit:          parseDecimalOpt(p.CumProfit),
	}
}

// List returns all DCA plans, optionally filtered by status or symbol.
// Pass empty string/nil to omit a filter.
func (c *DCAContext) List(ctx context.Context, status DCAStatus, symbol string) (*DcaList, error) {
	values := url.Values{}
	if status != "" {
		values.Add("status", string(status))
	}
	if symbol != "" {
		values.Add("counter_id", util.SymbolToCounterID(symbol))
	}
	var resp jsontypes.DcaList
	if err := c.httpClient.Get(ctx, "/v1/dailycoins/query", values, &resp); err != nil {
		return nil, err
	}
	result := &DcaList{}
	for _, p := range resp.Plans {
		result.Plans = append(result.Plans, convertPlan(p))
	}
	return result, nil
}

// Create creates a new DCA plan.
func (c *DCAContext) Create(ctx context.Context, symbol, amount string, frequency DCAFrequency, dayOfWeek string, dayOfMonth int32, allowMargin bool) (*DcaCreateResult, error) {
	body := map[string]interface{}{
		"counter_id":           util.SymbolToCounterID(symbol),
		"per_invest_amount":    amount,
		"invest_frequency":     string(frequency),
		"allow_margin_finance": allowMargin,
	}
	if dayOfWeek != "" {
		body["invest_day_of_week"] = dayOfWeek
	}
	if dayOfMonth > 0 {
		body["invest_day_of_month"] = strconv.FormatInt(int64(dayOfMonth), 10)
	}
	var resp jsontypes.DcaCreateResult
	if err := c.httpClient.Post(ctx, "/v1/dailycoins/create", body, &resp); err != nil {
		return nil, err
	}
	return &DcaCreateResult{PlanId: resp.PlanId}, nil
}

// Update modifies an existing DCA plan. Pass empty strings/nil to leave fields unchanged.
func (c *DCAContext) Update(ctx context.Context, planId string, amount string, frequency DCAFrequency, dayOfWeek string, dayOfMonth int32, allowMargin *bool) (*DcaCreateResult, error) {
	body := map[string]interface{}{"plan_id": planId}
	if amount != "" {
		body["per_invest_amount"] = amount
	}
	if frequency != "" {
		body["invest_frequency"] = string(frequency)
	}
	if dayOfWeek != "" {
		body["invest_day_of_week"] = dayOfWeek
	}
	if dayOfMonth > 0 {
		body["invest_day_of_month"] = strconv.FormatInt(int64(dayOfMonth), 10)
	}
	if allowMargin != nil {
		body["allow_margin_finance"] = *allowMargin
	}
	var resp jsontypes.DcaCreateResult
	if err := c.httpClient.Post(ctx, "/v1/dailycoins/update", body, &resp); err != nil {
		return nil, err
	}
	return &DcaCreateResult{PlanId: resp.PlanId}, nil
}

// Pause suspends an active DCA plan.
func (c *DCAContext) Pause(ctx context.Context, planId string) error {
	return c.httpClient.Post(ctx, "/v1/dailycoins/toggle", map[string]interface{}{
		"plan_id": planId,
		"status":  "Suspended",
	}, nil)
}

// Resume reactivates a suspended DCA plan.
func (c *DCAContext) Resume(ctx context.Context, planId string) error {
	return c.httpClient.Post(ctx, "/v1/dailycoins/toggle", map[string]interface{}{
		"plan_id": planId,
		"status":  "Active",
	}, nil)
}

// Stop terminates a DCA plan permanently.
func (c *DCAContext) Stop(ctx context.Context, planId string) error {
	return c.httpClient.Post(ctx, "/v1/dailycoins/toggle", map[string]interface{}{
		"plan_id": planId,
		"status":  "Finished",
	}, nil)
}

// History returns execution history for a DCA plan.
func (c *DCAContext) History(ctx context.Context, planId string, page, limit int32) (*DcaHistoryResponse, error) {
	values := url.Values{}
	values.Add("plan_id", planId)
	values.Add("page", strconv.FormatInt(int64(page), 10))
	values.Add("limit", strconv.FormatInt(int64(limit), 10))
	var resp jsontypes.DcaHistoryResponse
	if err := c.httpClient.Get(ctx, "/v1/dailycoins/query-records", values, &resp); err != nil {
		return nil, err
	}
	result := &DcaHistoryResponse{HasMore: resp.HasMore}
	for _, r := range resp.Records {
		result.Records = append(result.Records, &DcaHistoryRecord{
			CreatedAt:      r.CreatedAt,
			OrderId:        r.OrderId,
			Status:         r.Status,
			Action:         r.Action,
			OrderType:      r.OrderType,
			ExecutedQty:    parseDecimalOpt(r.ExecutedQty),
			ExecutedPrice:  parseDecimalOpt(r.ExecutedPrice),
			ExecutedAmount: parseDecimalOpt(r.ExecutedAmount),
			RejectedReason: r.RejectedReason,
			Symbol:         util.CounterIDToSymbol(r.CounterId),
		})
	}
	return result, nil
}

// Stats returns DCA portfolio statistics, optionally filtered by symbol.
func (c *DCAContext) Stats(ctx context.Context, symbol string) (*DcaStats, error) {
	values := url.Values{}
	if symbol != "" {
		values.Add("counter_id", util.SymbolToCounterID(symbol))
	}
	var resp jsontypes.DcaStats
	if err := c.httpClient.Get(ctx, "/v1/dailycoins/statistic", values, &resp); err != nil {
		return nil, err
	}
	result := &DcaStats{
		ActiveCount:    resp.ActiveCount,
		FinishedCount:  resp.FinishedCount,
		SuspendedCount: resp.SuspendedCount,
		RestDays:       resp.RestDays,
		TotalAmount:    parseDecimalOpt(resp.TotalAmount),
		TotalProfit:    parseDecimalOpt(resp.TotalProfit),
	}
	for _, p := range resp.NearestPlans {
		result.NearestPlans = append(result.NearestPlans, convertPlan(p))
	}
	return result, nil
}

// CheckSupport checks whether DCA is supported for the given symbols.
func (c *DCAContext) CheckSupport(ctx context.Context, symbols []string) (*DcaSupportList, error) {
	counterIds := make([]string, len(symbols))
	for i, s := range symbols {
		counterIds[i] = util.SymbolToCounterID(s)
	}
	var resp jsontypes.DcaSupportList
	if err := c.httpClient.Post(ctx, "/v1/dailycoins/batch-check-support", map[string]interface{}{
		"counter_ids": counterIds,
	}, &resp); err != nil {
		return nil, err
	}
	result := &DcaSupportList{}
	for _, info := range resp.Infos {
		result.Infos = append(result.Infos, &DcaSupportInfo{
			Symbol:               util.CounterIDToSymbol(info.CounterId),
			SupportRegularSaving: info.SupportRegularSaving,
		})
	}
	return result, nil
}

// CalcDate calculates the next projected investment date for a DCA plan configuration.
func (c *DCAContext) CalcDate(ctx context.Context, symbol string, frequency DCAFrequency, dayOfWeek string, dayOfMonth int32) (*DcaCalcDateResult, error) {
	body := map[string]interface{}{
		"counter_id":       util.SymbolToCounterID(symbol),
		"invest_frequency": string(frequency),
	}
	if dayOfWeek != "" {
		body["invest_day_of_week"] = dayOfWeek
	}
	if dayOfMonth > 0 {
		body["invest_day_of_month"] = strconv.FormatInt(int64(dayOfMonth), 10)
	}
	var resp jsontypes.DcaCalcDateResult
	if err := c.httpClient.Post(ctx, "/v1/dailycoins/calc-trd-date", body, &resp); err != nil {
		return nil, err
	}
	return &DcaCalcDateResult{TradeDate: resp.TradeDate}, nil
}

// SetReminder sets the reminder notification hours before DCA execution.
func (c *DCAContext) SetReminder(ctx context.Context, hours string) error {
	return c.httpClient.Post(ctx, "/v1/dailycoins/update-alter-hours", map[string]interface{}{
		"alter_hours": hours,
	}, nil)
}
