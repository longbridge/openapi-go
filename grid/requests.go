package grid

import (
	"github.com/shopspring/decimal"

	openapi "github.com/longbridge/openapi-go"
)

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
