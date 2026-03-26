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
