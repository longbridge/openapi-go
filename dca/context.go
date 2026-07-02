// Package dca provides a client for the Longbridge DCA (dollar-cost averaging) OpenAPI.
// It supports creating, updating, pausing, resuming, and stopping DCA plans,
// as well as querying execution history, statistics, and support eligibility.
package dca

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/dca/jsontypes"
	httplib "github.com/longbridge/openapi-go/http"
)

// DCAContext is a client for the Longbridge DCA OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	dctx, err := dca.NewFromCfg(conf)
//	list, err := dctx.List(ctx, nil, nil)
type DCAContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a DCAContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*DCAContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &DCAContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a DCAContext configured from environment variables.
func NewFromEnv() (*DCAContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// symbolToCounterID converts a Longbridge symbol (e.g. "700.HK") to the
// counter_id format expected by the DCA API (e.g. "ST/HK/700").
// Symbols without a dot separator are returned unchanged.
func symbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])
	return fmt.Sprintf("ST/%s/%s", market, code)
}

// counterIDToSymbol converts a counter_id (e.g. "ST/HK/700" or "ETF/US/SPY")
// back to a Longbridge symbol (e.g. "700.HK" or "SPY.US").
func counterIDToSymbol(counterID string) string {
	parts := strings.SplitN(counterID, "/", 3)
	if len(parts) == 3 {
		return fmt.Sprintf("%s.%s", parts[2], parts[1])
	}
	return counterID
}

// decimalFromString parses a decimal string; returns nil for empty strings.
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

// statusFromString converts an API status string to DCAStatus.
func statusFromString(s string) DCAStatus {
	switch s {
	case "Suspended":
		return DCAStatusSuspended
	case "Finished":
		return DCAStatusFinished
	default:
		return DCAStatusActive
	}
}

// frequencyFromString converts an API frequency string to DCAFrequency.
func frequencyFromString(s string) DCAFrequency {
	switch s {
	case "Daily":
		return DCAFrequencyDaily
	case "Weekly":
		return DCAFrequencyWeekly
	case "Fortnightly":
		return DCAFrequencyFortnightly
	default:
		return DCAFrequencyMonthly
	}
}

// convertPlan converts a jsontypes.DcaPlan to the idiomatic DcaPlan.
func convertPlan(j *jsontypes.DcaPlan) *DcaPlan {
	d, _ := decimal.NewFromString(j.PerInvestAmount)
	return &DcaPlan{
		PlanID:             j.PlanID,
		Status:             statusFromString(j.Status),
		Symbol:             counterIDToSymbol(j.CounterID),
		MemberID:           j.MemberID,
		Aaid:               j.Aaid,
		AccountChannel:     j.AccountChannel,
		DisplayAccount:     j.DisplayAccount,
		Market:             j.Market,
		PerInvestAmount:    d,
		Frequency:          frequencyFromString(j.InvestFrequency),
		DayOfWeek:          j.InvestDayOfWeek,
		DayOfMonth:         j.InvestDayOfMonth,
		AllowMarginFinance: j.AllowMarginFinance,
		AlterHours:         j.AlterHours.String(),
		CreatedAt:          j.CreatedAt,
		UpdatedAt:          j.UpdatedAt,
		NextTrdDate:        j.NextTrdDate,
		StockName:          j.StockName,
		CumAmount:          decimalFromString(j.CumAmount),
		IssueNumber:        j.IssueNumber,
		AverageCost:        decimalFromString(j.AverageCost),
		CumProfit:          decimalFromString(j.CumProfit),
	}
}

// checkAP returns an error when the client's credentials are for the US data center.
// All DCA APIs are served exclusively by the AP data center.
func (d *DCAContext) checkAP(path string) error {
	return d.httpClient.CheckRegion(path, "AP")
}

// List returns the caller's DCA plans, optionally filtered by status and/or symbol.
//
// Path: GET /v1/dailycoins/query
func (d *DCAContext) List(ctx context.Context, status *DCAStatus, symbol *string) (*DcaList, error) {
	if err := d.checkAP("/v1/dailycoins/query"); err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("page", "1")
	params.Set("limit", "100")
	if status != nil {
		params.Set("status", status.String())
	}
	if symbol != nil {
		params.Set("counter_id", symbolToCounterID(*symbol))
	}

	var resp jsontypes.DcaList
	if err := d.httpClient.Get(ctx, "/v1/dailycoins/query", params, &resp); err != nil {
		return nil, err
	}
	plans := make([]*DcaPlan, 0, len(resp.Plans))
	for _, p := range resp.Plans {
		plans = append(plans, convertPlan(p))
	}
	return &DcaList{Plans: plans}, nil
}

// CreateOptions holds optional parameters for creating a DCA plan.
type CreateOptions struct {
	// DayOfWeek is required for Weekly/Fortnightly plans (e.g. "Monday").
	DayOfWeek string
	// DayOfMonth is required for Monthly plans (e.g. 15).
	DayOfMonth *uint32
	// AllowMargin enables margin financing for this plan.
	AllowMargin bool
}

// Create creates a new DCA plan and returns the result containing the new plan ID.
//
// Path: POST /v1/dailycoins/create
func (d *DCAContext) Create(ctx context.Context, symbol string, amount string, frequency DCAFrequency, opts *CreateOptions) (*DcaCreateResult, error) {
	if err := d.checkAP("/v1/dailycoins/create"); err != nil {
		return nil, err
	}
	if opts == nil {
		opts = &CreateOptions{}
	}
	allowMargin := 0
	if opts.AllowMargin {
		allowMargin = 1
	}
	body := map[string]interface{}{
		"counter_id":           symbolToCounterID(symbol),
		"per_invest_amount":    amount,
		"invest_frequency":     frequency.String(),
		"allow_margin_finance": allowMargin,
	}
	if opts.DayOfWeek != "" {
		body["invest_day_of_week"] = opts.DayOfWeek
	}
	if opts.DayOfMonth != nil {
		body["invest_day_of_month"] = fmt.Sprintf("%d", *opts.DayOfMonth)
	}

	var resp jsontypes.DcaCreateResult
	if err := d.httpClient.Post(ctx, "/v1/dailycoins/create", body, &resp); err != nil {
		return nil, err
	}
	return &DcaCreateResult{PlanID: resp.PlanID}, nil
}

// UpdateOptions holds the fields that can be updated on an existing DCA plan.
// Nil/zero fields are not sent.
type UpdateOptions struct {
	Amount      *string
	Frequency   *DCAFrequency
	DayOfWeek   *string
	DayOfMonth  *uint32
	AllowMargin *bool
}

// Update modifies an existing DCA plan.
//
// Path: POST /v1/dailycoins/update
func (d *DCAContext) Update(ctx context.Context, planID string, opts *UpdateOptions) (*DcaCreateResult, error) {
	if err := d.checkAP("/v1/dailycoins/update"); err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"plan_id": planID,
	}
	if opts != nil {
		if opts.Amount != nil {
			body["per_invest_amount"] = *opts.Amount
		}
		if opts.Frequency != nil {
			body["invest_frequency"] = opts.Frequency.String()
		}
		if opts.DayOfWeek != nil {
			body["invest_day_of_week"] = *opts.DayOfWeek
		}
		if opts.DayOfMonth != nil {
			body["invest_day_of_month"] = fmt.Sprintf("%d", *opts.DayOfMonth)
		}
		if opts.AllowMargin != nil {
			m := 0
			if *opts.AllowMargin {
				m = 1
			}
			body["allow_margin_finance"] = m
		}
	}

	var resp jsontypes.DcaCreateResult
	if err := d.httpClient.Post(ctx, "/v1/dailycoins/update", body, &resp); err != nil {
		return nil, err
	}
	return &DcaCreateResult{PlanID: resp.PlanID}, nil
}

// Pause suspends an active DCA plan.
//
// Path: POST /v1/dailycoins/toggle
func (d *DCAContext) Pause(ctx context.Context, planID string) error {
	if err := d.checkAP("/v1/dailycoins/toggle"); err != nil {
		return err
	}
	body := map[string]interface{}{
		"plan_id": planID,
		"status":  "Suspended",
	}
	return d.httpClient.Post(ctx, "/v1/dailycoins/toggle", body, nil)
}

// Resume activates a suspended DCA plan.
//
// Path: POST /v1/dailycoins/toggle
func (d *DCAContext) Resume(ctx context.Context, planID string) error {
	if err := d.checkAP("/v1/dailycoins/toggle"); err != nil {
		return err
	}
	body := map[string]interface{}{
		"plan_id": planID,
		"status":  "Active",
	}
	return d.httpClient.Post(ctx, "/v1/dailycoins/toggle", body, nil)
}

// Stop permanently finishes a DCA plan.
//
// Path: POST /v1/dailycoins/toggle
func (d *DCAContext) Stop(ctx context.Context, planID string) error {
	if err := d.checkAP("/v1/dailycoins/toggle"); err != nil {
		return err
	}
	body := map[string]interface{}{
		"plan_id": planID,
		"status":  "Finished",
	}
	return d.httpClient.Post(ctx, "/v1/dailycoins/toggle", body, nil)
}

// History returns the execution history for a DCA plan.
//
// Path: GET /v1/dailycoins/query-records
func (d *DCAContext) History(ctx context.Context, planID string, page int, limit int) (*DcaHistoryResponse, error) {
	if err := d.checkAP("/v1/dailycoins/query-records"); err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Set("plan_id", planID)
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("limit", fmt.Sprintf("%d", limit))

	var resp jsontypes.DcaHistoryResponse
	if err := d.httpClient.Get(ctx, "/v1/dailycoins/query-records", params, &resp); err != nil {
		return nil, err
	}
	records := make([]*DcaHistoryRecord, 0, len(resp.Records))
	for _, r := range resp.Records {
		records = append(records, &DcaHistoryRecord{
			CreatedAt:      r.CreatedAt,
			OrderID:        r.OrderID,
			Status:         r.Status,
			Action:         r.Action,
			OrderType:      r.OrderType,
			ExecutedQty:    decimalFromString(r.ExecutedQty),
			ExecutedPrice:  decimalFromString(r.ExecutedPrice),
			ExecutedAmount: decimalFromString(r.ExecutedAmount),
			RejectedReason: r.RejectedReason,
			Symbol:         counterIDToSymbol(r.CounterID),
		})
	}
	return &DcaHistoryResponse{
		Records: records,
		HasMore: resp.HasMore,
	}, nil
}

// Stats returns aggregated DCA statistics, optionally filtered to a single symbol.
//
// Path: GET /v1/dailycoins/statistic
func (d *DCAContext) Stats(ctx context.Context, symbol *string) (*DcaStats, error) {
	if err := d.checkAP("/v1/dailycoins/statistic"); err != nil {
		return nil, err
	}
	params := url.Values{}
	if symbol != nil {
		params.Set("counter_id", symbolToCounterID(*symbol))
	}

	var resp jsontypes.DcaStats
	if err := d.httpClient.Get(ctx, "/v1/dailycoins/statistic", params, &resp); err != nil {
		return nil, err
	}
	nearestPlans := make([]*DcaPlan, 0, len(resp.NearestPlans))
	for _, p := range resp.NearestPlans {
		nearestPlans = append(nearestPlans, convertPlan(p))
	}
	return &DcaStats{
		ActiveCount:    resp.ActiveCount,
		FinishedCount:  resp.FinishedCount,
		SuspendedCount: resp.SuspendedCount,
		NearestPlans:   nearestPlans,
		RestDays:       resp.RestDays,
		TotalAmount:    decimalFromString(resp.TotalAmount),
		TotalProfit:    decimalFromString(resp.TotalProfit),
	}, nil
}

// CheckSupport checks which of the provided symbols support DCA plans.
//
// Path: POST /v1/dailycoins/batch-check-support
func (d *DCAContext) CheckSupport(ctx context.Context, symbols []string) ([]*DcaSupportInfo, error) {
	if err := d.checkAP("/v1/dailycoins/batch-check-support"); err != nil {
		return nil, err
	}
	counterIDs := make([]string, len(symbols))
	for i, s := range symbols {
		counterIDs[i] = symbolToCounterID(s)
	}
	body := map[string]interface{}{
		"counter_ids": counterIDs,
	}

	var resp jsontypes.DcaSupportList
	if err := d.httpClient.Post(ctx, "/v1/dailycoins/batch-check-support", body, &resp); err != nil {
		return nil, err
	}
	infos := make([]*DcaSupportInfo, 0, len(resp.Infos))
	for _, info := range resp.Infos {
		infos = append(infos, &DcaSupportInfo{
			Symbol:               counterIDToSymbol(info.CounterID),
			SupportRegularSaving: info.SupportRegularSaving,
		})
	}
	return infos, nil
}

// CalcDateOptions holds optional schedule parameters for calc-date.
type CalcDateOptions struct {
	DayOfWeek  string
	DayOfMonth *uint32
}

// CalcDate calculates the next projected trade date for the given schedule parameters.
//
// Path: POST /v1/dailycoins/calc-trd-date
func (d *DCAContext) CalcDate(ctx context.Context, symbol string, frequency DCAFrequency, opts *CalcDateOptions) (*DcaCalcDateResult, error) {
	if err := d.checkAP("/v1/dailycoins/calc-trd-date"); err != nil {
		return nil, err
	}
	body := map[string]interface{}{
		"counter_id":       symbolToCounterID(symbol),
		"invest_frequency": frequency.String(),
	}
	if opts != nil {
		if opts.DayOfWeek != "" {
			body["invest_day_of_week"] = opts.DayOfWeek
		}
		if opts.DayOfMonth != nil {
			body["invest_day_of_month"] = fmt.Sprintf("%d", *opts.DayOfMonth)
		}
	}

	var resp jsontypes.DcaCalcDateResult
	if err := d.httpClient.Post(ctx, "/v1/dailycoins/calc-trd-date", body, &resp); err != nil {
		return nil, err
	}
	t, err := parseTradeDate(resp.TradeDate)
	if err != nil {
		return nil, fmt.Errorf("dca: parse trade_date %q: %w", resp.TradeDate, err)
	}
	return &DcaCalcDateResult{TradeDate: t}, nil
}

// SetReminder updates the advance reminder hours for DCA execution notifications.
// hours must be one of "1", "6", or "12".
//
// Path: POST /v1/dailycoins/update-alter-hours
func (d *DCAContext) SetReminder(ctx context.Context, hours string) error {
	if err := d.checkAP("/v1/dailycoins/update-alter-hours"); err != nil {
		return err
	}
	body := map[string]interface{}{
		"alter_hours": hours,
	}
	return d.httpClient.Post(ctx, "/v1/dailycoins/update-alter-hours", body, nil)
}

// parseTradeDate parses a trade date string in "YYYY-MM-DD" format.
func parseTradeDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	// API may return "YYYY-MM-DD" or a Unix timestamp string.
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse trade_date %q: %w", s, err)
	}
	return time.Unix(ts, 0).UTC(), nil
}
