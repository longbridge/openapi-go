// Package jsontypes contains the raw JSON wire types for the DCA API.
// These types match the exact JSON field names returned by the Longbridge API.
// Use the parent dca package for idiomatic Go types.
package jsontypes

import "encoding/json"

// DcaList is the API response for listing DCA plans.
type DcaList struct {
	Plans []*DcaPlan `json:"plans"`
}

// DcaPlan is the raw JSON representation of a single DCA plan.
type DcaPlan struct {
	PlanID             string `json:"plan_id"`
	Status             string `json:"status"`
	CounterID          string `json:"counter_id"`
	MemberID           string `json:"member_id"`
	Aaid               string `json:"aaid"`
	AccountChannel     string `json:"account_channel"`
	DisplayAccount     string `json:"display_account"`
	Market             string `json:"market"`
	PerInvestAmount    string `json:"per_invest_amount"`
	InvestFrequency    string `json:"invest_frequency"`
	InvestDayOfWeek    string `json:"invest_day_of_week"`
	InvestDayOfMonth   string `json:"invest_day_of_month"`
	AllowMarginFinance bool   `json:"allow_margin_finance"`
	AlterHours         json.Number `json:"alter_hours"` // API returns int or string
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	NextTrdDate        string `json:"next_trd_date"`
	StockName          string `json:"stock_name"`
	CumAmount          string `json:"cum_amount"`
	IssueNumber        int64  `json:"issue_number"`
	AverageCost        string `json:"average_cost"`
	CumProfit          string `json:"cum_profit"`
}

// DcaStats is the raw JSON response for DCA statistics.
type DcaStats struct {
	ActiveCount    string     `json:"active_count"`
	FinishedCount  string     `json:"finished_count"`
	SuspendedCount string     `json:"suspended_count"`
	NearestPlans   []*DcaPlan `json:"nearest_plans"`
	RestDays       string     `json:"rest_days"`
	TotalAmount    string     `json:"total_amount"`
	TotalProfit    string     `json:"total_profit"`
}

// DcaSupportList is the raw JSON response for batch check-support.
type DcaSupportList struct {
	Infos []*DcaSupportInfo `json:"infos"`
}

// DcaSupportInfo is the raw JSON representation of DCA support for a security.
type DcaSupportInfo struct {
	CounterID            string `json:"counter_id"`
	SupportRegularSaving bool   `json:"support_regular_saving"`
}

// DcaHistoryResponse is the raw JSON response for DCA execution history.
type DcaHistoryResponse struct {
	Records []*DcaHistoryRecord `json:"records"`
	HasMore bool                `json:"has_more"`
}

// DcaHistoryRecord is the raw JSON representation of one DCA execution record.
type DcaHistoryRecord struct {
	CreatedAt      string `json:"created_at"`
	OrderID        string `json:"order_id"`
	Status         string `json:"status"`
	Action         string `json:"action"`
	OrderType      string `json:"order_type"`
	ExecutedQty    string `json:"executed_qty"`
	ExecutedPrice  string `json:"executed_price"`
	ExecutedAmount string `json:"executed_amount"`
	RejectedReason string `json:"rejected_reason"`
	CounterID      string `json:"counter_id"`
}

// DcaCreateResult is the raw JSON response for create/update DCA plan.
type DcaCreateResult struct {
	PlanID string `json:"plan_id"`
}

// DcaCalcDateResult is the raw JSON response for calc-trd-date.
type DcaCalcDateResult struct {
	TradeDate string `json:"trade_date"`
}
