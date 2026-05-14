package dca

import "github.com/shopspring/decimal"

// DCAFrequency is the investment frequency for a DCA plan.
type DCAFrequency string

const (
	DCAFrequencyDaily       DCAFrequency = "Daily"
	DCAFrequencyWeekly      DCAFrequency = "Weekly"
	DCAFrequencyFortnightly DCAFrequency = "Fortnightly"
	DCAFrequencyMonthly     DCAFrequency = "Monthly"
)

// DCAStatus is the status of a DCA plan.
type DCAStatus string

const (
	DCAStatusActive    DCAStatus = "Active"
	DCAStatusSuspended DCAStatus = "Suspended"
	DCAStatusFinished  DCAStatus = "Finished"
)

type DcaList struct {
	Plans []*DcaPlan
}

type DcaPlan struct {
	PlanId             string
	Status             string
	Symbol             string
	MemberId           string
	Aaid               string
	AccountChannel     string
	DisplayAccount     string
	Market             string
	PerInvestAmount    *decimal.Decimal
	InvestFrequency    string
	InvestDayOfWeek    string
	InvestDayOfMonth   string
	AllowMarginFinance bool
	AlterHours         string
	CreatedAt          string
	UpdatedAt          string
	NextTrdDate        string
	StockName          string
	CumAmount          *decimal.Decimal
	IssueNumber        int64
	AverageCost        *decimal.Decimal
	CumProfit          *decimal.Decimal
}

type DcaStats struct {
	ActiveCount    string
	FinishedCount  string
	SuspendedCount string
	NearestPlans   []*DcaPlan
	RestDays       string
	TotalAmount    *decimal.Decimal
	TotalProfit    *decimal.Decimal
}

type DcaSupportList struct {
	Infos []*DcaSupportInfo
}

type DcaSupportInfo struct {
	Symbol               string
	SupportRegularSaving bool
}

type DcaHistoryResponse struct {
	Records []*DcaHistoryRecord
	HasMore bool
}

type DcaHistoryRecord struct {
	CreatedAt      string
	OrderId        string
	Status         string
	Action         string
	OrderType      string
	ExecutedQty    *decimal.Decimal
	ExecutedPrice  *decimal.Decimal
	ExecutedAmount *decimal.Decimal
	RejectedReason string
	Symbol         string
}

type DcaCreateResult struct {
	PlanId string
}

type DcaCalcDateResult struct {
	TradeDate string
}
