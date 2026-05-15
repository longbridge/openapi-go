package dca

import (
	"time"

	"github.com/shopspring/decimal"
)

// DCAFrequency specifies how often a DCA plan executes.
type DCAFrequency int

const (
	// DCAFrequencyDaily executes the plan every trading day.
	DCAFrequencyDaily DCAFrequency = iota
	// DCAFrequencyWeekly executes the plan once per week.
	DCAFrequencyWeekly
	// DCAFrequencyFortnightly executes the plan once every two weeks.
	DCAFrequencyFortnightly
	// DCAFrequencyMonthly executes the plan once per month (default).
	DCAFrequencyMonthly
)

// String returns the API wire string for a DCAFrequency.
func (f DCAFrequency) String() string {
	switch f {
	case DCAFrequencyDaily:
		return "Daily"
	case DCAFrequencyWeekly:
		return "Weekly"
	case DCAFrequencyFortnightly:
		return "Fortnightly"
	default:
		return "Monthly"
	}
}

// DCAStatus represents the lifecycle state of a DCA plan.
type DCAStatus int

const (
	// DCAStatusActive is the default running state.
	DCAStatusActive DCAStatus = iota
	// DCAStatusSuspended means the plan is paused.
	DCAStatusSuspended
	// DCAStatusFinished means the plan has been permanently stopped.
	DCAStatusFinished
)

// String returns the API wire string for a DCAStatus.
func (s DCAStatus) String() string {
	switch s {
	case DCAStatusSuspended:
		return "Suspended"
	case DCAStatusFinished:
		return "Finished"
	default:
		return "Active"
	}
}

// DcaPlan is the idiomatic Go representation of a DCA investment plan.
type DcaPlan struct {
	// PlanID is the unique identifier of the plan.
	PlanID string
	// Status is the current lifecycle state of the plan.
	Status DCAStatus
	// Symbol is the security symbol in Longbridge format (e.g. "700.HK").
	Symbol string
	// MemberID is the user's member identifier.
	MemberID string
	// Aaid is the account asset ID.
	Aaid string
	// AccountChannel identifies the brokerage channel.
	AccountChannel string
	// DisplayAccount is the masked account display string.
	DisplayAccount string
	// Market is the market identifier.
	Market string
	// PerInvestAmount is the amount invested per execution.
	PerInvestAmount decimal.Decimal
	// Frequency is the recurrence interval.
	Frequency DCAFrequency
	// DayOfWeek is the scheduled day of week for weekly plans (e.g. "Monday").
	DayOfWeek string
	// DayOfMonth is the scheduled day of month for monthly plans (e.g. "15").
	DayOfMonth string
	// AllowMarginFinance indicates whether margin financing is enabled.
	AllowMarginFinance bool
	// AlterHours is the advance reminder hours before execution.
	AlterHours string
	// CreatedAt is the creation timestamp of the plan.
	CreatedAt string
	// UpdatedAt is the last update timestamp.
	UpdatedAt string
	// NextTrdDate is the next projected trade date.
	NextTrdDate string
	// StockName is the display name of the security.
	StockName string
	// CumAmount is the total cumulative invested amount (nil if not available).
	CumAmount *decimal.Decimal
	// IssueNumber is the number of executions completed.
	IssueNumber int64
	// AverageCost is the average cost per share (nil if not available).
	AverageCost *decimal.Decimal
	// CumProfit is the cumulative profit/loss (nil if not available).
	CumProfit *decimal.Decimal
}

// DcaList holds a list of DCA plans.
type DcaList struct {
	Plans []*DcaPlan
}

// DcaStats holds aggregated statistics across DCA plans.
type DcaStats struct {
	// ActiveCount is the number of active plans.
	ActiveCount string
	// FinishedCount is the number of finished plans.
	FinishedCount string
	// SuspendedCount is the number of suspended plans.
	SuspendedCount string
	// NearestPlans is the list of plans executing soonest.
	NearestPlans []*DcaPlan
	// RestDays is the number of rest days remaining.
	RestDays string
	// TotalAmount is the total invested across all plans (nil if not available).
	TotalAmount *decimal.Decimal
	// TotalProfit is the total profit/loss across all plans (nil if not available).
	TotalProfit *decimal.Decimal
}

// DcaSupportInfo holds DCA eligibility for a single security.
type DcaSupportInfo struct {
	// Symbol is the security symbol.
	Symbol string
	// SupportRegularSaving indicates whether DCA is supported.
	SupportRegularSaving bool
}

// DcaHistoryRecord is a single DCA execution event.
type DcaHistoryRecord struct {
	// CreatedAt is the execution timestamp string.
	CreatedAt string
	// OrderID is the order identifier.
	OrderID string
	// Status is the execution status string.
	Status string
	// Action is the action taken.
	Action string
	// OrderType is the order type used.
	OrderType string
	// ExecutedQty is the number of shares executed (nil if not filled).
	ExecutedQty *decimal.Decimal
	// ExecutedPrice is the execution price (nil if not filled).
	ExecutedPrice *decimal.Decimal
	// ExecutedAmount is the total execution amount (nil if not filled).
	ExecutedAmount *decimal.Decimal
	// RejectedReason explains why the order was rejected, if applicable.
	RejectedReason string
	// Symbol is the security symbol.
	Symbol string
}

// DcaHistoryResponse is the paginated response for execution history.
type DcaHistoryResponse struct {
	Records []*DcaHistoryRecord
	HasMore bool
}

// DcaCreateResult is returned when creating or updating a plan.
type DcaCreateResult struct {
	PlanID string
}

// DcaCalcDateResult holds the next projected trade date.
type DcaCalcDateResult struct {
	TradeDate time.Time
}
