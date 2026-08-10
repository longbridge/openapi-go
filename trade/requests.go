package trade

import (
	"time"

	"github.com/longbridge/openapi-go"
	"github.com/shopspring/decimal"
)

type GetHistoryExecutions struct {
	Symbol  string    // optional
	StartAt time.Time // optional
	EndAt   time.Time // optional
}

type GetTodayExecutions struct {
	Symbol  string // optional
	OrderId string // optional
}

type GetAllExecutions struct {
	Symbol  string    // optional
	OrderId string    // optional
	StartAt time.Time // optional
	EndAt   time.Time // optional
	Page    int64     // optional
}

type GetHistoryOrders struct {
	Symbol  string         // optional
	Status  []OrderStatus  // optional
	Side    OrderSide      // optional
	Market  openapi.Market // optional
	StartAt int64          // optional
	EndAt   int64          // optional
}

type GetTodayOrders struct {
	Symbol string         // optional
	Status []OrderStatus  // optional
	Side   OrderSide      // optional
	Market openapi.Market // optional
}

type GetFundPositions struct {
	Symbols []string // optional
}

type GetStockPositions struct {
	Symbols []string // optional
}

type GetCashFlow struct {
	StartAt      int64 // start timestamp , required
	EndAt        int64 // end timestamp, required
	BusinessType BalanceType
	Symbol       string
	Page         int64
	Size         int64
}

type ReplaceOrder struct {
	OrderId         string          // required
	Quantity        uint64          // required
	Price           decimal.Decimal // LO / ELO / ALO / ODD / LIT Order Required
	TriggerPrice    decimal.Decimal // LIT / MIT Order Required
	LimitOffset     decimal.Decimal // TSLPAMT / TSLPPCT Order Required
	TrailingAmount  decimal.Decimal // TSLPAMT / TSMAMT Order Required
	TrailingPercent decimal.Decimal // TSLPPCT / TSMAPCT Order Required
	Remark          string
}

type SubmitOrder struct {
	Symbol            string          // required
	OrderType         OrderType       // required
	Side              OrderSide       // required
	SubmittedQuantity uint64          // required
	SubmittedPrice    decimal.Decimal // LO / ELO / ALO / ODD / LIT Order Required
	TriggerPrice      decimal.Decimal // LIT / MIT Order Required
	LimitOffset       decimal.Decimal // TSLPAMT / TSLPPCT Order Required
	TrailingAmount    decimal.Decimal // TSLPAMT / TSMAMT Order Required
	TrailingPercent   decimal.Decimal // TSLPPCT / TSMAPCT Order Required
	ExpireDate        *time.Time      // required when time_in_force is GTD
	OutsideRTH        OutsideRTH
	Remark            string
	TimeInForce       TimeType // required
}

type GetEstimateMaxPurchaseQuantity struct {
	Symbol    string
	OrderType OrderType
	Price     decimal.Decimal
	Currency  string
	OrderId   string
	Side      OrderSide
}

type GetAccountBalance struct {
	Currency Currency // optional
}

// GridTradeRule is the grid trading rule for submit / replace requests.
// Prices and quantities use decimal.Decimal; unset (zero-value) decimals are
// omitted from the request. Enum-like fields are raw integers whose code
// tables are documented inline.
type GridTradeRule struct {
	SubmittedBasePrice decimal.Decimal // Base price the grid is anchored to
	UpperLimitPrice    decimal.Decimal // Upper price bound
	LowerLimitPrice    decimal.Decimal // Lower price bound
	TriggerPriceType   int32           // Trigger price type (1 = spread, 2 = percent)
	TriggerSpreadUp    decimal.Decimal // Upward trigger spread (absolute)
	TriggerSpreadDown  decimal.Decimal // Downward trigger spread (absolute)
	TriggerPercentUp   decimal.Decimal // Upward trigger percent
	TriggerPercentDown decimal.Decimal // Downward trigger percent
	MultipleTrigger    bool            // Whether a single grid level may trigger multiple times
	TimeInForce        int32           // Time in force (0 = Day, 1 = GTC, 6 = GTD)
	UpperLimitQuantity decimal.Decimal // Quantity handled when the upper bound is reached
	LowerLimitQuantity decimal.Decimal // Quantity handled when the lower bound is reached
	ExpireTime         int64           // Expiry time (unix seconds), used with GTD
	UpperLimitEvent    int32           // Action when the upper bound is reached (1 / 2)
	LowerLimitEvent    int32           // Action when the lower bound is reached (1 / 2)
	TriggerSellDepth   int32           // Sell-side order-book depth (-5..5, 0 = use GridOrderTypeUp)
	TriggerBuyDepth    int32           // Buy-side order-book depth (-5..5, 0 = use GridOrderTypeDown)
	TriggerQuantity    decimal.Decimal // Quantity per trigger
	SupportShortsell   bool            // Whether short selling is allowed
	Rth                int32           // Regular trading hours flag (0 / 1 / 2)
	GridOrderTypeUp    string          // Sell-side order type when depth is 0 (GMO / GLO / GTG)
	GridOrderTypeDown  string          // Buy-side order type when depth is 0 (GMO / GLO / GTG)
}

// SubmitGridOrder is the request for submitting a grid trading order.
type SubmitGridOrder struct {
	Symbol             string        // required, e.g. 700.HK
	SettlementCurrency string        // required, e.g. HKD
	GridTradingRule    GridTradeRule // required
}

// ReplaceGridOrder is the request for replacing a grid trading order.
type ReplaceGridOrder struct {
	OrderId         string        // required
	GridTradingRule GridTradeRule // required
}

// GetGridOrders is the request for the grid orders list.
type GetGridOrders struct {
	Page      int32          // optional
	Limit     int32          // optional
	Market    openapi.Market // optional
	Status    string         // optional, comma-joined (e.g. "Performing,Suspended")
	Symbol    string         // optional, e.g. 700.HK
	SortBy    string         // optional
	SortOrder string         // optional
}

// GetGridOrderDetail is the request for a grid order detail.
type GetGridOrderDetail struct {
	OrderId   string // required
	HistoryId string // optional, history cursor for paging
	Limit     int32  // optional, history page size
}

// GetGridTriggerHistory is the request for grid trigger history.
// Note: the required parameter is named grid_order_id (not order_id).
type GetGridTriggerHistory struct {
	GridOrderId string // required
	Page        int32  // optional
	Limit       int32  // optional
}
