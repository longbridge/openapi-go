package grid

import (
	"github.com/shopspring/decimal"
)

// GridOrders is the response for the grid orders list request.
type GridOrders struct {
	GridOrder []*GridOrder `json:"grid_order"`
	HasMore   bool         `json:"has_more"`
}

// GridTriggerHistory is the response for the grid trigger history request.
type GridTriggerHistory struct {
	TriggerOrders []*TriggerOrder `json:"trigger_orders"`
	HasMore       bool            `json:"has_more"`
}

// GridOrder is a grid trading order (element of the list / by-ids responses).
type GridOrder struct {
	OrderId              string           `json:"order_id"`
	Symbol               string           `json:"symbol"`
	StockName            string           `json:"stock_name"`
	Market               string           `json:"market"`
	Status               string           `json:"status"`
	GridStatus           string           `json:"grid_status"`
	SubmittedBasePrice   *decimal.Decimal `json:"submitted_base_price"`
	CurrentBasePrice     *decimal.Decimal `json:"current_base_price"`
	PreTriggerBasePrice  *decimal.Decimal `json:"pre_trigger_base_price"`
	PostTriggerBasePrice *decimal.Decimal `json:"post_trigger_base_price"`
	UpperLimitPrice      *decimal.Decimal `json:"upper_limit_price"`
	LowerLimitPrice      *decimal.Decimal `json:"lower_limit_price"`
	TriggerPriceType     TriggerPriceType `json:"trigger_price_type"` // Spread / Percent
	TriggerSpreadUp      *decimal.Decimal `json:"trigger_spread_up"`
	TriggerSpreadDown    *decimal.Decimal `json:"trigger_spread_down"`
	TriggerPercentUp     *decimal.Decimal `json:"trigger_percent_up"`
	TriggerPercentDown   *decimal.Decimal `json:"trigger_percent_down"`
	PullbackPercent      *decimal.Decimal `json:"pullback_percent"`
	PullbackSpread       *decimal.Decimal `json:"pullback_spread"`
	ReboundPercent       *decimal.Decimal `json:"rebound_percent"`
	ReboundSpread        *decimal.Decimal `json:"rebound_spread"`
	TriggerSellOrderType string           `json:"trigger_sell_order_type"`
	TriggerBuyOrderType  string           `json:"trigger_buy_order_type"`
	TriggerSellDepth     int32            `json:"trigger_sell_depth"`
	TriggerBuyDepth      int32            `json:"trigger_buy_depth"`
	TriggerQuantity      *decimal.Decimal `json:"trigger_quantity"`
	TriggerSellQuantity  *decimal.Decimal `json:"trigger_sell_quantity"`
	TriggerBuyQuantity   *decimal.Decimal `json:"trigger_buy_quantity"`
	UpperLimitQuantity   *decimal.Decimal `json:"upper_limit_quantity"`
	LowerLimitQuantity   *decimal.Decimal `json:"lower_limit_quantity"`
	UpperLimitEvent      GridLimitEvent   `json:"upper_limit_event"`
	LowerLimitEvent      GridLimitEvent   `json:"lower_limit_event"`
	MultipleTrigger      bool             `json:"multiple_trigger"`
	TriggerTimes         int32            `json:"trigger_times"`
	TotalBuyQuantity     *decimal.Decimal `json:"total_buy_quantity"`
	TotalSellQuantity    *decimal.Decimal `json:"total_sell_quantity"`
	TotalProfitBalance   *decimal.Decimal `json:"total_profit_balance"`
	SettlementCurrency   string           `json:"settlement_currency"`
	TimeInForce          GridTimeInForce  `json:"time_in_force"` // Day / GTC / GTD
	Gtd                  string           `json:"gtd"`
	CreatedAt            string           `json:"created_at"`
	Rth                  int32            `json:"rth"`
	SupportShortsell     bool             `json:"support_shortsell"`
	GridOrderTypeUp      string           `json:"grid_order_type_up"`   // GMO / GLO / GTG
	GridOrderTypeDown    string           `json:"grid_order_type_down"` // GMO / GLO / GTG
}

// GridOrderSubOrder is a triggered sub-order carried in the grid order detail.
type GridOrderSubOrder struct {
	Id          string           `json:"id"`
	Price       *decimal.Decimal `json:"price"`
	OrderType   string           `json:"order_type"`
	Quantity    *decimal.Decimal `json:"quantity"`
	ExecutedQty *decimal.Decimal `json:"executed_qty"`
	Action      int32            `json:"action"`
	Status      string           `json:"status"`
	SubmittedAt string           `json:"submitted_at"`
	Rth         int32            `json:"rth"`
}

// GridOrderHistory is a grid order lifecycle-history entry.
type GridOrderHistory struct {
	HistoryId     string `json:"history_id"`
	CreatedAt     string `json:"created_at"`
	Status        string `json:"status"`
	SuspendReason string `json:"suspend_reason"`
	Reason        string `json:"reason"`
}

// GridOrderDetail is the detail of a grid trading order.
type GridOrderDetail struct {
	OrderId             string               `json:"order_id"`
	Symbol              string               `json:"symbol"`
	StockName           string               `json:"stock_name"`
	Status              string               `json:"status"`
	GridStatus          string               `json:"grid_status"`
	SuspendReason       string               `json:"suspend_reason"`
	SleepingReason      string               `json:"sleeping_reason"`
	SubmittedBasePrice  *decimal.Decimal     `json:"submitted_base_price"`
	CurrentBasePrice    *decimal.Decimal     `json:"current_base_price"`
	UpperLimitPrice     *decimal.Decimal     `json:"upper_limit_price"`
	LowerLimitPrice     *decimal.Decimal     `json:"lower_limit_price"`
	TriggerPriceType    TriggerPriceType     `json:"trigger_price_type"`
	TriggerSpreadUp     *decimal.Decimal     `json:"trigger_spread_up"`
	TriggerSpreadDown   *decimal.Decimal     `json:"trigger_spread_down"`
	TriggerPercentUp    *decimal.Decimal     `json:"trigger_percent_up"`
	TriggerPercentDown  *decimal.Decimal     `json:"trigger_percent_down"`
	PullbackPercent     *decimal.Decimal     `json:"pullback_percent"`
	PullbackSpread      *decimal.Decimal     `json:"pullback_spread"`
	ReboundPercent      *decimal.Decimal     `json:"rebound_percent"`
	ReboundSpread       *decimal.Decimal     `json:"rebound_spread"`
	MultipleTrigger     bool                 `json:"multiple_trigger"`
	TimeInForce         GridTimeInForce      `json:"time_in_force"`
	TriggerQuantity     *decimal.Decimal     `json:"trigger_quantity"`
	TriggerSellQuantity *decimal.Decimal     `json:"trigger_sell_quantity"`
	TriggerBuyQuantity  *decimal.Decimal     `json:"trigger_buy_quantity"`
	UpperLimitQuantity  *decimal.Decimal     `json:"upper_limit_quantity"`
	LowerLimitQuantity  *decimal.Decimal     `json:"lower_limit_quantity"`
	UpperLimitEvent     GridLimitEvent       `json:"upper_limit_event"`
	LowerLimitEvent     GridLimitEvent       `json:"lower_limit_event"`
	TriggerSellDepth    int32                `json:"trigger_sell_depth"`
	TriggerBuyDepth     int32                `json:"trigger_buy_depth"`
	CreatedAt           string               `json:"created_at"`
	UpdatedAt           string               `json:"updated_at"`
	SettlementCurrency  string               `json:"settlement_currency"`
	ExpireTime          string               `json:"expire_time"`
	Gtd                 string               `json:"gtd"`
	GridSubOrders       []*GridOrderSubOrder `json:"grid_sub_orders"`
	SubHasMore          bool                 `json:"sub_has_more"`
	GridOrderHistory    []*GridOrderHistory  `json:"grid_order_history"`
	HistoryHasMore      bool                 `json:"history_has_more"`
	SupportShortsell    bool                 `json:"support_shortsell"`
	Rth                 int32                `json:"rth"`
	GridOrderTypeUp     string               `json:"grid_order_type_up"`
	GridOrderTypeDown   string               `json:"grid_order_type_down"`
}

// TriggerOrder is a grid trigger-history entry (one triggered order).
type TriggerOrder struct {
	Id            string           `json:"id"`
	Status        string           `json:"status"`
	Name          string           `json:"name"`
	Symbol        string           `json:"symbol"`
	Price         *decimal.Decimal `json:"price"`
	Quantity      *decimal.Decimal `json:"quantity"`
	ExecutedPrice *decimal.Decimal `json:"executed_price"`
	ExecutedQty   *decimal.Decimal `json:"executed_qty"`
	SubmittedAt   string           `json:"submitted_at"`
	Action        int32            `json:"action"`
	OrderType     string           `json:"order_type"`
	TriggerPrice  *decimal.Decimal `json:"trigger_price"`
	Msg           string           `json:"msg"`
	Currency      string           `json:"currency"`
	LastDone      *decimal.Decimal `json:"last_done"`
	UpdatedAt     string           `json:"updated_at"`
	TimeInForce   GridTimeInForce  `json:"time_in_force"`
	Gtd           string           `json:"gtd"`
	TriggerAt     string           `json:"trigger_at"`
	TriggerStatus int32            `json:"trigger_status"`
}

// GridBidSize is a price-step (bid-size) rule entry from the symbol-info response.
type GridBidSize struct {
	StrProceed *decimal.Decimal `json:"str_proceed"`
	EndProceed *decimal.Decimal `json:"end_proceed"`
	BidSize    *decimal.Decimal `json:"bid_size"`
}

// GridChannelInfo is the channel / authorization info nested in symbol-info.
type GridChannelInfo struct {
	StrategyGranted    bool     `json:"strategy_granted"`
	SupportRth         bool     `json:"support_rth"`
	Currency           string   `json:"currency"`
	SettlementCurrency []string `json:"settlement_currency"`
}

// GridSymbolInfo is the /v1/orders/info response — the security (symbol) info
// used to build a grid order.
type GridSymbolInfo struct {
	Name        string           `json:"name"`
	LastDone    *decimal.Decimal `json:"last_done"`
	LotSize     *decimal.Decimal `json:"lot_size"`
	BuyLotSize  *decimal.Decimal `json:"buy_lot_size"`
	SellLotSize *decimal.Decimal `json:"sell_lot_size"`
	BidSizes    []*GridBidSize   `json:"bid_sizes"`
	ChannelInfo GridChannelInfo  `json:"channel_info"`
}
