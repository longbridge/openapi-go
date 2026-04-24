package trade

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go"
)

type (
	OrderType            string
	OrderSide            string
	OutsideRTH           string // Outside regular trading hours
	OrderStatus          string
	Market               string
	OrderTag             string
	TriggerStatus        string
	TriggerPriceType     string
	BalanceType          int32
	OfDirection          int32
	TimeType             string
	CommissionFreeStatus string
	DeductionStatus      string
	ChargeCategoryCode   string
	Currency             string
)

const (
	// Balance type
	BalanceTypeUnknown BalanceType = 0 // unknown type
	BalanceTypeCash    BalanceType = 1 // cash
	BalanceTypeStock   BalanceType = 2 // stock
	BalanceTypeFund    BalanceType = 3 // fund

	// Outflow direction
	OfDirectionUnkown OfDirection = 0
	OfDirectionOut    OfDirection = 1 // outflow
	OfDirectionIn     OfDirection = 2 // inflow

	// Time force in type
	TimeTypeDay TimeType = "Day" // Day Order
	TimeTypeGTC TimeType = "GTC" // Good Til Canceled Order
	TimeTypeGTD TimeType = "GTD" // Good Til Date Order

	// Order type
	OrderTypeLO      OrderType = "LO"      // Limit Order
	OrderTypeELO     OrderType = "ELO"     // Enhanced Limit Order
	OrderTypeMO      OrderType = "MO"      // Market Order
	OrderTypeAO      OrderType = "AO"      // At-auction Order
	OrderTypeALO     OrderType = "ALO"     // At-auction Limit Order
	OrderTypeODD     OrderType = "ODD"     // Odd Lots Order
	OrderTypeLIT     OrderType = "LIT"     // Limit If Touched
	OrderTypeMIT     OrderType = "MIT"     // Market If Touched
	OrderTypeTSLPAMT OrderType = "TSLPAMT" // Trailing Limit If Touched (Trailing Amount)
	OrderTypeTSLPPCT OrderType = "TSLPPCT" // Trailing Limit If Touched (Trailing Percent)
	OrderTypeTSMAMT  OrderType = "TSMAMT"  // Trailing Market If Touched (Trailing Amount)
	OrderTypeTSMPCT  OrderType = "TSMPCT"  // Trailing Market If Touched (Trailing Percent)
	OrderTypeSLO     OrderType = "SLO"     // SLO order

	// Order side
	OrderSideBuy  OrderSide = "Buy"
	OrderSideSell OrderSide = "Sell"

	// Order status
	OrderNotReported          OrderStatus = "NotReported"
	OrderReplacedNotReported  OrderStatus = "ReplacedNotReported"
	OrderProtectedNotReported OrderStatus = "ProtectedNotReported"
	OrderVarietiesNotReported OrderStatus = "VarietiesNotReported"
	OrderFilledStatus         OrderStatus = "FilledStatus"
	OrderWaitToNew            OrderStatus = "WaitToNew"
	OrderNewStatus            OrderStatus = "NewStatus"
	OrderWaitToReplace        OrderStatus = "WaitToReplace"
	OrderPendingReplaceStatus OrderStatus = "PendingReplaceStatus"
	OrderReplacedStatus       OrderStatus = "ReplacedStatus"
	OrderPartialFilledStatus  OrderStatus = "PartialFilledStatus"
	OrderWaitToCancel         OrderStatus = "WaitToCancel"
	OrderPendingCancelStatus  OrderStatus = "PendingCancelStatus"
	OrderRejectedStatus       OrderStatus = "RejectedStatus"
	OrderCanceledStatus       OrderStatus = "CanceledStatus"
	OrderExpiredStatus        OrderStatus = "ExpiredStatus"
	OrderPartialWithdrawal    OrderStatus = "PartialWithdrawal"

	// Outside RTH
	OutsideRTHOnly      OutsideRTH = "RTH_ONLY"          // Regular trading hour only
	OutsideRTHAny       OutsideRTH = "ANY_TIME"          // Any time
	OutsideRTHOvernight OutsideRTH = "OVERNIGHT"         // Overnight
	OutsideRTHUnknown   OutsideRTH = "UnknownOutsideRth" // Default is UnknownOutsideRth when the order is not a US stock

	// Order tag
	OrderTagNormal       OrderTag = "Normal"       // Normal Order
	OrderTagLongTerm     OrderTag = "GTC"          // Long term Order
	OrderTagGrey         OrderTag = "Grey"         // Grey Order
	OrderTagMarginCall   OrderTag = "MarginCall"   // Force Selling
	OrderTagOffline      OrderTag = "Offline"      // OTC
	OrderTagCreditor     OrderTag = "Creditor"     // Option Exercise Long
	OrderTagDebtor       OrderTag = "Debtor"       // Option Exercise Short
	OrderTagNonExercise  OrderTag = "NonExercise"  // Wavier Of Option Exercise
	OrderTagAllocatedSub OrderTag = "AllocatedSub" // Trade Allocation

	// Trigger status
	TriggerStatusDeactive TriggerStatus = "DEACTIVE"
	TriggerStatusActive   TriggerStatus = "ACTIVE"
	TriggerStatusReleased TriggerStatus = "RELEASED"

	// Trigger price type
	TriggerPriceTypeLimit  TriggerPriceType = "LIT" // Limit If Touched
	TriggerPriceTypeMarket TriggerPriceType = "MIT" // Market If Touched

	// Commission-free Status
	CommissionFreeStatusNone      CommissionFreeStatus = "None"
	CommissionFreeStatusCaculated CommissionFreeStatus = "Calculated" // Commission-free amount to be calculated
	CommissionFreeStatusPending   CommissionFreeStatus = "Pending"    // Pending commission-free
	CommissionFreeStatusReady     CommissionFreeStatus = "Ready"      // Commission-free applied

	// Deduction status/Cashback Status
	DeductionStatusNone    DeductionStatus = "NONE"
	DeductionStatusNoData  DeductionStatus = "NO_DATA"
	DeductionStatusPending DeductionStatus = "PENDING"
	DeductionStatusDone    DeductionStatus = "DONE"

	// Charge category code
	ChargeCategoryCodeUnknown    ChargeCategoryCode = "UNKNOWN"
	ChargeCategoryCodeBrokerFees ChargeCategoryCode = "BROKER_FEES"
	ChargeCategoryCodeThirdFees  ChargeCategoryCode = "THIRD_FEES"

	// Currency
	CurrencyHKD     Currency = "HKD"
	CurrencyUSD     Currency = "USD"
	CurrencyCNH     Currency = "CNH"
	CurrencyDefault Currency = ""
)

// Execution is execution details
type Execution struct {
	OrderId     string
	TradeId     string
	Symbol      string
	TradeDoneAt time.Time
	Quantity    string
	Price       *decimal.Decimal
}

// Executions has a Execution list
type Executions struct {
	Trades []*Execution
}

// Orders has a Order details
type Orders struct {
	HasMore bool
	Orders  []*Order
}

// Order is order details
type Order struct {
	OrderId          string
	Status           OrderStatus
	StockName        string
	Quantity         string
	ExecutedQuantity string
	Price            *decimal.Decimal
	ExecutedPrice    *decimal.Decimal
	SubmittedAt      string
	Side             OrderSide
	Symbol           string
	OrderType        OrderType
	LastDone         *decimal.Decimal
	TriggerPrice     *decimal.Decimal
	Msg              string
	Tag              OrderTag
	TimeInForce      TimeType
	ExpireDate       string
	UpdatedAt        string
	TriggerAt        string
	TrailingAmount   *decimal.Decimal
	TrailingPercent  *decimal.Decimal
	LimitOffset      *decimal.Decimal
	TriggerStatus    TriggerStatus
	Currency         string
	OutsideRth       OutsideRTH
	Remark           string
}

type OrderChargeItem struct {
	Code ChargeCategoryCode
	Name string
	Fees []OrderChargeFee
}

type OrderChargeDetail struct {
	TotalAmount decimal.Decimal
	Currency    string
	Items       []OrderChargeItem
}

type OrderChargeFee struct {
	Code ChargeCategoryCode
	Name string
	Fees []OrderChargeFee
}

type OrderHistoryDetail struct {
	// Executed price for executed orders, submitted price for expired,
	// canceled, rejected orders, etc.
	Price decimal.Decimal
	// Executed quantity for executed orders, remaining quantity for expired,
	// canceled, rejected orders, etc.
	Quantity int64
	Status   OrderStatus
	Msg      string // Execution or error message
	Time     string // Occurrence time
}

type OrderDetail struct {
	OrderId                  string
	Status                   OrderStatus
	StockName                string
	Quantity                 int64 // Submitted quantity
	ExecutedQuantity         int64
	Price                    *decimal.Decimal // Submitted price
	ExecutedPrice            *decimal.Decimal
	SubmittedAt              string    // Submitted time
	Side                     OrderSide /// Order side
	Symbol                   string
	OrderType                OrderType
	LastDone                 *decimal.Decimal
	TriggerPrice             *decimal.Decimal
	Msg                      string // Rejected Message or remark
	Tag                      OrderTag
	TimeInForce              TimeType
	ExpireDate               string
	UpdatedAt                string
	TriggerAt                string // Conditional order trigger time
	TrailingAmount           *decimal.Decimal
	TrailingPercent          *decimal.Decimal
	LimitOffset              *decimal.Decimal
	TriggerStatus            TriggerStatus
	Currency                 string
	OutsideRth               OutsideRTH /// Enable or disable outside regular trading hours
	Remark                   string
	FreeStatus               CommissionFreeStatus
	FreeAmount               *decimal.Decimal
	FreeCurrency             string
	DeductionsStatus         DeductionStatus
	DeductionsAmount         *decimal.Decimal
	DeductionsCurrency       string
	PlatformDeductedStatus   DeductionStatus
	PlatformDeductedAmount   *decimal.Decimal
	PlatformDeductedCurrency string
	History                  OrderHistoryDetail
	ChargeDetail             OrderChargeDetail
}

// AccountBalances has a AccountBalance list
type AccountBalances struct {
	List []*AccountBalance
}

// AccountBalance is user account balance
type AccountBalance struct {
	TotalCash              *decimal.Decimal
	MaxFinanceAmount       *decimal.Decimal
	RemainingFinanceAmount *decimal.Decimal
	RiskLevel              string
	MarginCall             *decimal.Decimal
	NetAssets              *decimal.Decimal // net asset
	InitMargin             *decimal.Decimal // initial margin
	MaintenanceMargin      *decimal.Decimal // maintenance margin
	Currency               string
	BuyPower               *decimal.Decimal
	CashInfos              []*CashInfo
}

// FundPositions has a FundPosition list
type FundPositions struct {
	List []*FundPositionChannel
}

// FundPositionChannel is a account channel's fund position details
type FundPositionChannel struct {
	AccountChannel string
	Positions      []*FundPosition
}

// FundPosition is fund position details
type FundPosition struct {
	Symbol               string
	CurrentNetAssetValue *decimal.Decimal
	NetAssetValueDay     int64
	SymbolName           string
	Currency             string
	CostNetAssetValue    *decimal.Decimal
	HoldingUnits         *decimal.Decimal
}

// StockPositions has a StockPosition list
type StockPositions struct {
	List []*StockPositionChannel
}

// StockPositionChannel is a account channel's stock positions details
type StockPositionChannel struct {
	AccountChannel string
	Positions      []*StockPosition
}

// StockPosition is user stock position details
type StockPosition struct {
	Symbol            string
	SymbolName        string
	Quantity          string
	AvailableQuantity string
	Currency          string
	CostPrice         *decimal.Decimal
	Market            openapi.Market
}

// CashFlows has a CashFlow list
type CashFlows struct {
	List []*CashFlow
}

// CashFlow is cash flow details
type CashFlow struct {
	TransactionFlowName string
	Direction           OfDirection
	BusinessType        BalanceType
	Balance             *decimal.Decimal
	Currency            string
	BusinessTime        string
	Symbol              string
	Description         string
}

// CashInfo
type CashInfo struct {
	WithdrawCash  *decimal.Decimal
	AvailableCash *decimal.Decimal
	FrozenCash    *decimal.Decimal
	SettlingCash  *decimal.Decimal
	Currency      string
}

// PushEvent is quote context callback event
type PushEvent struct {
	Event string
	Data  *PushOrderChanged
}

// PushOrderChanged is order change event details
type PushOrderChanged struct {
	AccountNo        string
	Currency         string
	ExecutedPrice    *decimal.Decimal
	ExecutedQuantity *decimal.Decimal
	LastPrice        *decimal.Decimal
	LastShare        *decimal.Decimal
	LimitOffset      string
	Msg              string
	OrderId          string
	OrderType        OrderType
	Side             OrderSide
	Status           OrderStatus
	StockName        string
	SubmittedAt      string
	Price            *decimal.Decimal
	Quantity         *decimal.Decimal
	Symbol           string
	Tag              OrderTag
	TrailingAmount   *decimal.Decimal
	TrailingPercent  string
	TriggerAt        string
	TriggerPrice     *decimal.Decimal
	TriggerStatus    TriggerStatus
	UpdatedAt        string
	Remark           string
}

// SubResponse is subscribe function response
type SubResponse struct {
	Success []string
	Fail    []*SubResponseFail
	Current []string
}

// SubResponseFail contains subscribe failed reason
type SubResponseFail struct {
	Topic  string
	Reason string
}

type UnsubResponse struct {
	Current []string
}

// MarginRatio contains some ratio
type MarginRatio struct {
	ImFactor *decimal.Decimal // Initial margin ratio
	MmFactor *decimal.Decimal // Maintain the initial margin ratio
	FmFactor *decimal.Decimal // Forced close-out margin ratio
}

// EstimateMaxPurchaseQuantity is response for estimate maximum purchase quantity
type EstimateMaxPurchaseQuantityResponse struct {
	CashMaxQty   int64 // Cash available quantity
	MarginMaxQty int64 // Margin available quantity
}
