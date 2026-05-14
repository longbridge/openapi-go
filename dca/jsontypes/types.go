package jsontypes

type DcaList struct {
	Plans []*DcaPlan `json:"plans"`
}

type DcaPlan struct {
	PlanId             string `json:"plan_id"`
	Status             string `json:"status"`
	CounterId          string `json:"counter_id"`
	MemberId           string `json:"member_id"`
	Aaid               string `json:"aaid"`
	AccountChannel     string `json:"account_channel"`
	DisplayAccount     string `json:"display_account"`
	Market             string `json:"market"`
	PerInvestAmount    string `json:"per_invest_amount"`
	InvestFrequency    string `json:"invest_frequency"`
	InvestDayOfWeek    string `json:"invest_day_of_week"`
	InvestDayOfMonth   string `json:"invest_day_of_month"`
	AllowMarginFinance bool   `json:"allow_margin_finance"`
	AlterHours         string `json:"alter_hours"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	NextTrdDate        string `json:"next_trd_date"`
	StockName          string `json:"stock_name"`
	CumAmount          string `json:"cum_amount"`
	IssueNumber        int64  `json:"issue_number"`
	AverageCost        string `json:"average_cost"`
	CumProfit          string `json:"cum_profit"`
}

type DcaStats struct {
	ActiveCount    string     `json:"active_count"`
	FinishedCount  string     `json:"finished_count"`
	SuspendedCount string     `json:"suspended_count"`
	NearestPlans   []*DcaPlan `json:"nearest_plans"`
	RestDays       string     `json:"rest_days"`
	TotalAmount    string     `json:"total_amount"`
	TotalProfit    string     `json:"total_profit"`
}

type DcaSupportList struct {
	Infos []*DcaSupportInfo `json:"infos"`
}

type DcaSupportInfo struct {
	CounterId            string `json:"counter_id"`
	SupportRegularSaving bool   `json:"support_regular_saving"`
}

type DcaHistoryResponse struct {
	Records []*DcaHistoryRecord `json:"records"`
	HasMore bool                `json:"has_more"`
}

type DcaHistoryRecord struct {
	CreatedAt      string `json:"created_at"`
	OrderId        string `json:"order_id"`
	Status         string `json:"status"`
	Action         string `json:"action"`
	OrderType      string `json:"order_type"`
	ExecutedQty    string `json:"executed_qty"`
	ExecutedPrice  string `json:"executed_price"`
	ExecutedAmount string `json:"executed_amount"`
	RejectedReason string `json:"rejected_reason"`
	CounterId      string `json:"counter_id"`
}

type DcaCreateResult struct {
	PlanId string `json:"plan_id"`
}

type DcaCalcDateResult struct {
	TradeDate string `json:"trade_date"`
}
