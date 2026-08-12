// Package jsontypes contains the raw JSON wire types for the grid trading API.
// These types match the exact JSON field names returned by the Longbridge API.
// Use the parent grid package for idiomatic Go types.
package jsontypes

// GridTradeRule is the grid trading rule request body used by submit / replace.
// Prices and quantities are serialized as strings; enum-like fields are raw
// integers. Optional string fields are omitted when empty.
type GridTradeRule struct {
	SubmittedBasePrice string `json:"submitted_base_price,omitempty"`
	UpperLimitPrice    string `json:"upper_limit_price,omitempty"`
	LowerLimitPrice    string `json:"lower_limit_price,omitempty"`
	TriggerPriceType   int32  `json:"trigger_price_type"`
	TriggerSpreadUp    string `json:"trigger_spread_up,omitempty"`
	TriggerSpreadDown  string `json:"trigger_spread_down,omitempty"`
	TriggerPercentUp   string `json:"trigger_percent_up,omitempty"`
	TriggerPercentDown string `json:"trigger_percent_down,omitempty"`
	MultipleTrigger    bool   `json:"multiple_trigger"`
	TimeInForce        int32  `json:"time_in_force"`
	UpperLimitQuantity string `json:"upper_limit_quantity,omitempty"`
	LowerLimitQuantity string `json:"lower_limit_quantity,omitempty"`
	ExpireTime         int64  `json:"expire_time"`
	UpperLimitEvent    int32  `json:"upper_limit_event"`
	LowerLimitEvent    int32  `json:"lower_limit_event"`
	TriggerSellDepth   int32  `json:"trigger_sell_depth"`
	TriggerBuyDepth    int32  `json:"trigger_buy_depth"`
	TriggerQuantity    string `json:"trigger_quantity,omitempty"`
	SupportShortsell   bool   `json:"support_shortsell"`
	Rth                int32  `json:"rth"`
	GridOrderTypeUp    string `json:"grid_order_type_up,omitempty"`
	GridOrderTypeDown  string `json:"grid_order_type_down,omitempty"`
}

// SubmitGridOrder is the request body for submitting a grid trading order.
type SubmitGridOrder struct {
	Symbol             string        `json:"symbol"`
	SettlementCurrency string        `json:"settlement_currency"`
	GridTradingRule    GridTradeRule `json:"grid_trading_rule"`
}

// ReplaceGridOrder is the request body for replacing a grid trading order.
type ReplaceGridOrder struct {
	OrderId         string        `json:"order_id"`
	GridTradingRule GridTradeRule `json:"grid_trading_rule"`
}

// GridOrderIdBody is the shared request body for cancel / suspend / restart.
type GridOrderIdBody struct {
	OrderId string `json:"order_id"`
}

// GridOrderIdsBody is the request body for querying grid orders by IDs.
type GridOrderIdsBody struct {
	OrderIds []string `json:"order_ids"`
}

// SubmitStrategyQuestionnaire is the request body for the strategy
// risk-disclosure questionnaire record.
type SubmitStrategyQuestionnaire struct {
	Type  string            `json:"type"`
	Items map[string]string `json:"items"`
}

// SubmitGridOrderResponse is the response for submit grid trading order.
type SubmitGridOrderResponse struct {
	OrderId string `json:"order_id"`
}

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
	OrderId              string `json:"order_id"`
	Symbol               string `json:"symbol"`
	StockName            string `json:"stock_name"`
	Market               string `json:"market"`
	Status               string `json:"status"`
	GridStatus           string `json:"grid_status"`
	SubmittedBasePrice   string `json:"submitted_base_price"`
	CurrentBasePrice     string `json:"current_base_price"`
	PreTriggerBasePrice  string `json:"pre_trigger_base_price"`
	PostTriggerBasePrice string `json:"post_trigger_base_price"`
	UpperLimitPrice      string `json:"upper_limit_price"`
	LowerLimitPrice      string `json:"lower_limit_price"`
	TriggerPriceType     int32  `json:"trigger_price_type"`
	TriggerSpreadUp      string `json:"trigger_spread_up"`
	TriggerSpreadDown    string `json:"trigger_spread_down"`
	TriggerPercentUp     string `json:"trigger_percent_up"`
	TriggerPercentDown   string `json:"trigger_percent_down"`
	PullbackPercent      string `json:"pullback_percent"`
	PullbackSpread       string `json:"pullback_spread"`
	ReboundPercent       string `json:"rebound_percent"`
	ReboundSpread        string `json:"rebound_spread"`
	TriggerSellOrderType string `json:"trigger_sell_order_type"`
	TriggerBuyOrderType  string `json:"trigger_buy_order_type"`
	TriggerSellDepth     int32  `json:"trigger_sell_depth"`
	TriggerBuyDepth      int32  `json:"trigger_buy_depth"`
	TriggerQuantity      string `json:"trigger_quantity"`
	TriggerSellQuantity  string `json:"trigger_sell_quantity"`
	TriggerBuyQuantity   string `json:"trigger_buy_quantity"`
	UpperLimitQuantity   string `json:"upper_limit_quantity"`
	LowerLimitQuantity   string `json:"lower_limit_quantity"`
	UpperLimitEvent      int32  `json:"upper_limit_event"`
	LowerLimitEvent      int32  `json:"lower_limit_event"`
	MultipleTrigger      bool   `json:"multiple_trigger"`
	TriggerTimes         int32  `json:"trigger_times"`
	TotalBuyQuantity     string `json:"total_buy_quantity"`
	TotalSellQuantity    string `json:"total_sell_quantity"`
	TotalProfitBalance   string `json:"total_profit_balance"`
	SettlementCurrency   string `json:"settlement_currency"`
	TimeInForce          int32  `json:"time_in_force"`
	Gtd                  string `json:"gtd"`
	CreatedAt            string `json:"created_at"`
	Rth                  int32  `json:"rth"`
	SupportShortsell     bool   `json:"support_shortsell"`
	GridOrderTypeUp      string `json:"grid_order_type_up"`
	GridOrderTypeDown    string `json:"grid_order_type_down"`
}

// GridOrderSubOrder is a triggered sub-order carried in the grid order detail.
type GridOrderSubOrder struct {
	Id          string `json:"id"`
	Price       string `json:"price"`
	OrderType   string `json:"order_type"`
	Quantity    string `json:"quantity"`
	ExecutedQty string `json:"executed_qty"`
	Action      int32  `json:"action"`
	Status      string `json:"status"`
	SubmittedAt string `json:"submitted_at"`
	Rth         int32  `json:"rth"`
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
	SubmittedBasePrice  string               `json:"submitted_base_price"`
	CurrentBasePrice    string               `json:"current_base_price"`
	UpperLimitPrice     string               `json:"upper_limit_price"`
	LowerLimitPrice     string               `json:"lower_limit_price"`
	TriggerPriceType    int32                `json:"trigger_price_type"`
	TriggerSpreadUp     string               `json:"trigger_spread_up"`
	TriggerSpreadDown   string               `json:"trigger_spread_down"`
	TriggerPercentUp    string               `json:"trigger_percent_up"`
	TriggerPercentDown  string               `json:"trigger_percent_down"`
	PullbackPercent     string               `json:"pullback_percent"`
	PullbackSpread      string               `json:"pullback_spread"`
	ReboundPercent      string               `json:"rebound_percent"`
	ReboundSpread       string               `json:"rebound_spread"`
	MultipleTrigger     bool                 `json:"multiple_trigger"`
	TimeInForce         int32                `json:"time_in_force"`
	TriggerQuantity     string               `json:"trigger_quantity"`
	TriggerSellQuantity string               `json:"trigger_sell_quantity"`
	TriggerBuyQuantity  string               `json:"trigger_buy_quantity"`
	UpperLimitQuantity  string               `json:"upper_limit_quantity"`
	LowerLimitQuantity  string               `json:"lower_limit_quantity"`
	UpperLimitEvent     int32                `json:"upper_limit_event"`
	LowerLimitEvent     int32                `json:"lower_limit_event"`
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
	Id            string `json:"id"`
	Status        string `json:"status"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	ExecutedPrice string `json:"executed_price"`
	ExecutedQty   string `json:"executed_qty"`
	SubmittedAt   string `json:"submitted_at"`
	Action        int32  `json:"action"`
	OrderType     string `json:"order_type"`
	TriggerPrice  string `json:"trigger_price"`
	Msg           string `json:"msg"`
	Currency      string `json:"currency"`
	LastDone      string `json:"last_done"`
	UpdatedAt     string `json:"updated_at"`
	TimeInForce   int32  `json:"time_in_force"`
	Gtd           string `json:"gtd"`
	TriggerAt     string `json:"trigger_at"`
	TriggerStatus int32  `json:"trigger_status"`
}

// GridBidSize is a price-step (bid-size) rule entry from the order-info response.
type GridBidSize struct {
	StrProceed string `json:"str_proceed"`
	EndProceed string `json:"end_proceed"`
	BidSize    string `json:"bid_size"`
}

// GridChannelInfo is the channel / authorization info nested in order-info.
type GridChannelInfo struct {
	StrategyGranted    bool     `json:"strategy_granted"`
	SupportRth         bool     `json:"support_rth"`
	Currency           string   `json:"currency"`
	SettlementCurrency []string `json:"settlement_currency"`
}

// GridSymbolInfo is the /v1/orders/info response — the security (symbol) info
// used to build a grid order.
type GridSymbolInfo struct {
	Name        string          `json:"name"`
	LastDone    string          `json:"last_done"`
	LotSize     string          `json:"lot_size"`
	BuyLotSize  string          `json:"buy_lot_size"`
	SellLotSize string          `json:"sell_lot_size"`
	BidSizes    []*GridBidSize  `json:"bid_sizes"`
	ChannelInfo GridChannelInfo `json:"channel_infos"`
}
