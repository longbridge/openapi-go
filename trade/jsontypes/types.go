package jsontypes

// Execution is execution details
type Execution struct {
	OrderId     string `json:"order_id"`
	TradeId     string `json:"trade_id"`
	Symbol      string `json:"symbol"`
	TradeDoneAt int64  `json:"trade_done_at,string"`
	Quantity    string `json:"quantity"`
	Price       string `json:"price"`
}

// Executions has a Execution list
type Executions struct {
	Trades []*Execution `json:"trades"`
}

// AllExecutionsResponse is the response for get all executions request
type AllExecutionsResponse struct {
	HasMore bool         `json:"has_more"`
	Trades  []*Execution `json:"trades"`
}

type SubmitOrderResponse struct {
	OrderId string `json:"order_id"`
}

// Orders has a Order details
type Orders struct {
	HasMore bool     `json:"has_more"`
	Orders  []*Order `json:"orders"`
}

// Order is order details
type Order struct {
	OrderId          string `json:"order_id"`
	Status           string `json:"status"`
	StockName        string `json:"stock_name"`
	Quantity         string `json:"quantity"`
	ExecutedQuantity string `json:"executed_quantity"`
	Price            string `json:"price"`
	ExecutedPrice    string `json:"executed_price"`
	SubmittedAt      string `json:"submmited_at"`
	Side             string `json:"side"`
	Symbol           string `json:"symbol"`
	OrderType        string `json:"order_type"`
	LastDone         string `json:"last_done"`
	TriggerPrice     string `json:"trigger_price"`
	Msg              string `json:"msg"`
	Tag              string `json:"tag"`
	TimeInForce      string `json:"time_in_force"`
	ExpireDate       string `json:"expire_date"`
	UpdatedAt        string `json:"updated_at"`
	TriggerAt        string `json:"trigger_at"`
	TrailingAmount   string `json:"trailing_amount"`
	TrailingPercent  string `json:"trailing_percent"`
	LimitOffset      string `json:"limit_offset"`
	TriggerStatus    string `json:"trigger_status"`
	Currency         string `json:"currency"`
	OutsideRth       string `json:"outside_rth"`
	Remark           string `json:"remark"`
}

// AccountBalances has a AccountBalance list
type AccountBalances struct {
	List []*AccountBalance `json:"list"`
}

// AccountBalance is user account balance
type AccountBalance struct {
	TotalCash              string      `json:"total_cash"`
	MaxFinanceAmount       string      `json:"max_finance_amount"`
	RemainingFinanceAmount string      `json:"remaining_finance_amount"`
	RiskLevel              string      `json:"risk_level"`
	MarginCall             string      `json:"margin_call"`
	NetAssets              string      `json:"net_assets"`
	InitMargin             string      `json:"init_margin"`
	MaintenanceMargin      string      `json:"maintenance_margin"`
	Currency               string      `json:"currency"`
	BuyPower               string      `json:"buy_power"`
	CashInfos              []*CashInfo `json:"cash_infos"`
}

// FundPositions has a FundPosition list
type FundPositions struct {
	List []*FundPositionChannel `json:"list"`
}

// FundPositionChannel is a account channel's fund position details
type FundPositionChannel struct {
	AccountChannel string          `json:"account_channel"`
	Positions      []*FundPosition `json:"fund_info"`
}

// FundPosition is fund position details
type FundPosition struct {
	Symbol               string `json:"symbol"`
	CurrentNetAssetValue string `json:"current_net_asset_value"`
	NetAssetValueDay     int64  `json:"net_asset_value_day,string"` // timestamp
	SymbolName           string `json:"symbol_name"`
	Currency             string `json:"currency"`
	CostNetAssetValue    string `json:"cost_net_asset_value"`
	HoldingUnits         string `json:"holding_units"`
}

// StockPositions has a StockPosition list
type StockPositions struct {
	List []*StockPositionChannel `json:"list"`
}

// StockPositionChannel is a account channel's stock positions details
type StockPositionChannel struct {
	AccountChannel string           `json:"account_channel"`
	Positions      []*StockPosition `json:"stock_info"`
}

// StockPosition is user stock position details
type StockPosition struct {
	Symbol            string `json:"symbol"`
	SymbolName        string `json:"symbol_name"`
	Quantity          string `json:"quantity"`
	AvailableQuantity string `json:"available_quantity"`
	Currency          string `json:"currency"`
	CostPrice         string `json:"cost_price"`
	Market            string `json:"market"`
}

// CashFlows has a CashFlow list
type CashFlows struct {
	List []*CashFlow `json:"list"`
}

// CashFlow is cash flow details
type CashFlow struct {
	TransactionFlowName string `json:"transaction_flow_name"`
	Direction           int32  `json:"direction"`
	BusinessType        int32  `json:"business_type"`
	Balance             string `json:"balance"`
	Currency            string `json:"currency"`
	BusinessTime        string `json:"business_time"`
	Symbol              string `json:"symbol"`
	Description         string `json:"description"`
}

// CashInfo
type CashInfo struct {
	WithdrawCash  string `json:"withdraw_cash"`
	AvailableCash string `json:"available_cash"`
	FrozenCash    string `json:"frozen_cash"`
	SettlingCash  string `json:"settling_cash"`
	Currency      string `json:"currency"`
}

// PushEvent is quote context callback event
type PushEvent struct {
	Event string            `json:"event"`
	Data  *PushOrderChanged `json:"data"`
}

// PushOrderChanged is order change event details
type PushOrderChanged struct {
	AccountNo        string `json:"account_no"`
	Currency         string `json:"currency"`
	ExecutedPrice    string `json:"executed_price"`
	ExecutedQuantity string `json:"executed_quantity"`
	LastPrice        string `json:"last_price"`
	LastShare        string `json:"last_share"`
	LimitOffset      string `json:"limit_offset"`
	Msg              string `json:"msg"`
	OrderId          string `json:"order_id"`
	OrderType        string `json:"order_type"`
	Side             string `json:"side"`
	Status           string `json:"status"`
	StockName        string `json:"stock_name"`
	SubmittedAt      string `json:"submitted_at"`
	Price            string `json:"submitted_price"`
	Quantity         string `json:"submitted_quantity"`
	Symbol           string `json:"symbol"`
	Tag              string `json:"tag"`
	TrailingAmount   string `json:"trailing_amount"`
	TrailingPercent  string `json:"trailing_percent"`
	TriggerAt        string `json:"trigger_at"`
	TriggerPrice     string `json:"trigger_price"`
	TriggerStatus    string `json:"trigger_status"`
	UpdatedAt        string `json:"updated_at"`
	Remark           string `json:"remark"`
}

type ReplaceOrder struct {
	OrderId         string `json:"order_id"`
	Quantity        uint64 `json:"quantity,string"`
	Price           string `json:"price"`
	TriggerPrice    string `json:"trigger_price,omitempty"`
	LimitOffset     string `json:"limit_offset,omitempty"`
	TrailingAmount  string `json:"trailing_ammount,omitempty"`
	TrailingPercent string `json:"trailing_percent,omitempty"`
	Remark          string `json:"remark"`
}

type SubmitOrder struct {
	Symbol            string `json:"symbol"`
	OrderType         string `json:"order_type"`
	Side              string `json:"side"`
	SubmittedQuantity uint64 `json:"submitted_quantity,string"`
	SubmittedPrice    string `json:"submitted_price,omitempty"`
	TriggerPrice      string `json:"trigger_price,omitempty"`
	LimitOffset       string `json:"limit_offset,omitempty"`
	TrailingAmount    string `json:"trailing_amount,omitempty"`
	TrailingPercent   string `json:"trailing_percent,omitempty"`
	ExpireDate        string `json:"expire_date,omitempty"`
	OutsideRTH        string `json:"outside_rth,omitempty"`
	Remark            string `json:"remark,omitempty"`
	TimeInForce       string `json:"time_in_force"`
}

type MarginRatio struct {
	ImFactor string `json:"im_factor,omitempty"`
	MmFactor string `json:"mm_factor,omitempty"`
	FmFactor string `json:"fm_factor,omitempty"`
}

type OrderChargeItem struct {
	Code string           `json:"code"`
	Name string           `json:"name"`
	Fees []OrderChargeFee `json:"fees"`
}

type OrderChargeDetail struct {
	TotalAmount string            `json:"total_amount"`
	Currency    string            `json:"currency"`
	Items       []OrderChargeItem `json:"items"`
}

type OrderChargeFee struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type OrderHistoryDetail struct {
	// Executed price for executed orders, submitted price for expired,
	// canceled, rejected orders, etc.
	Price string `json:"price"`
	// Executed quantity for executed orders, remaining quantity for expired,
	// canceled, rejected orders, etc.
	Quantity string `json:"quantity"`
	Status   string `json:"status"`
	Msg      string `json:"msg"`  // Execution or error message
	Time     string `json:"time"` // Occurrence time
}

// OrderDetail is for get order detail response
type OrderDetail struct {
	OrderId                  string               `json:"order_id"`
	Status                   string               `json:"status"`
	StockName                string               `json:"stock_name"`
	Quantity                 string               `json:"quantity"` // Submitted quantity
	ExecutedQuantity         string               `json:"executed_quantity"`
	Price                    string               `json:"price"` // Submitted price
	ExecutedPrice            string               `json:"executed_price"`
	SubmittedAt              string               `json:"submitted_at"`
	Side                     string               `json:"side"` // Order side
	Symbol                   string               `json:"symbol"`
	OrderType                string               `json:"order_type"`
	LastDone                 string               `json:"last_done"`
	TriggerPrice             string               `json:"trigger_price"`
	Msg                      string               `json:"msg"` // Rejected Message or remark
	Tag                      string               `json:"tag"`
	TimeInForce              string               `json:"time_in_force"`
	ExpireDate               string               `json:"expire_date"`
	UpdatedAt                string               `json:"updated_at"`
	TriggerAt                string               `json:"trigger_at"` // Conditional order trigger time
	TrailingAmount           string               `json:"trailing_amount"`
	TrailingPercent          string               `json:"trailing_precent"`
	LimitOffset              string               `json:"limit_offset"`
	TriggerStatus            string               `json:"trigger_status"`
	Currency                 string               `json:"currency"`
	OutsideRth               string               `json:"outside_rth"` // Enable or disable outside regular trading hours
	Remark                   string               `json:"remark"`
	FreeStatus               string               `json:"free_status"`
	FreeAmount               string               `json:"free_amount"`
	FreeCurrency             string               `json:"free_currency"`
	DeductionsStatus         string               `json:"deductions_status"`
	DeductionsAmount         string               `json:"deductions_amount"`
	DeductionsCurrency       string               `json:"deductions_currency"`
	PlatformDeductedStatus   string               `json:"platform_deducted_status"`
	PlatformDeductedAmount   string               `json:"platform_deducted_amount"`
	PlatformDeductedCurrency string               `json:"platform_deducted_currency"`
	History                  []OrderHistoryDetail `json:"history"`
	ChargeDetail             OrderChargeDetail    `json:"charge_detail"`
}

// EstimateMaxPurchaseQuantity is response for estimate maximum purchase quantity
type EstimateMaxPurchaseQuantityResponse struct {
	CashMaxQty   string `json:"cash_max_qty"`   // Cash available quantity
	MarginMaxQty string `json:"margin_max_qty"` // Margin available quantity
}

// ==================== Grid trading ====================

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

// GridOrderInfo is the /v1/orders/info response used by the grid order window.
type GridOrderInfo struct {
	Name         string          `json:"name"`
	LastDone     string          `json:"last_done"`
	LotSize      string          `json:"lot_size"`
	BuyLotSize   string          `json:"buy_lot_size"`
	SellLotSize  string          `json:"sell_lot_size"`
	BidSizes     []*GridBidSize  `json:"bid_sizes"`
	ChannelInfos GridChannelInfo `json:"channel_infos"`
}
