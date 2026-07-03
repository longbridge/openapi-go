package trade

import "time"

// GetUSHistoryOrders is the request for QueryUSOrders, modelled after
// GetHistoryOrders for HK/CN orders.
//
// QueryType: 0=all (includes Rejected), 1=pending, 2=history (filled only).
// Default (QueryType=0) matches what the app shows as "past orders".
type GetUSHistoryOrders struct {
	Symbol    string    // optional — user-facing symbol e.g. "AAPL.US" or "DOGEUSD.BKKT"
	Side      OrderSide // optional — Buy / Sell; zero value = all directions
	StartAt   int64     // optional — unix seconds; zero = last 90 days
	EndAt     int64     // optional — unix seconds; zero = now
	QueryType int32     // 0=all, 1=pending, 2=filled history; default 0
	Page      int32     // 1-based page number; default 1
	Limit     int32     // page size; default 20
}

// QueryUSOrdersRequest is an alias for GetUSHistoryOrders.
type QueryUSOrdersRequest = GetUSHistoryOrders

// GetUSRealizedPL is the request for USRealizedPL.
type GetUSRealizedPL struct {
	Currency string // required — e.g. "USD"
	// Category filters by asset type: "ALL", "STOCK", "OPTION", "CRYPTO".
	// Zero value / empty = all categories.
	Category string
}

// USOrder is one order entry in QueryUSOrdersResponse.
// counter_id fields are converted to user-facing symbol format.
type USOrder struct {
	// OrderID is the unique order identifier (field "id" in raw response).
	// Use this with USOrderDetail.
	OrderID            string    `json:"id"`
	AAID               string    `json:"aaid"`
	AccountChannel     string    `json:"account_channel"`
	// Action: 1=buy, 2=sell
	Action             int32     `json:"action"`
	// Symbol is converted from counter_id (e.g. "VA/BKKT/DOGEUSD" → "DOGEUSD.BKKT")
	Symbol             string    `json:"symbol"`
	// UnderlyingSymbol is converted from underlying_counter_id (options only)
	UnderlyingSymbol   string    `json:"underlying_symbol"`
	Code               string    `json:"code"`
	Name               string    `json:"name"`
	SecurityType       string    `json:"security_type"`
	Currency           string    `json:"currency"`
	TradeCurrency      string    `json:"trade_currency"`
	OrderType          string    `json:"order_type"`
	Status             string    `json:"status"`
	Price              string    `json:"price"`
	Quantity           string    `json:"quantity"`
	ExecutedQty        string    `json:"executed_qty"`
	ExecutedPrice      string    `json:"executed_price"`
	ExecutedAmount     string    `json:"executed_amount"`
	OperateDirection   string    `json:"operate_direction"`
	TimeInForce        int32     `json:"time_in_force"`
	GTD                string    `json:"gtd"`
	SubmittedAt        time.Time `json:"submitted_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Msg                string    `json:"msg"`
	Report             string    `json:"report"`
	// Options-specific fields
	ContractDirection  string    `json:"contract_direction"`
	ContractDueDate    string    `json:"contract_due_date"`
	StrikePrice        string    `json:"strike_price"`
	// Trailing stop fields
	TailingAmount      string    `json:"tailing_amount"`
	TailingPercent     string    `json:"tailing_percent"`
	// Trigger / conditional order fields
	TriggerPrice       string    `json:"trigger_price"`
	TriggerStatus      int32     `json:"trigger_status"`
	TriggerAt          string    `json:"trigger_at"`
	TriggerExchange    string    `json:"trigger_exchange"`
	TriggerLastDone    string    `json:"trigger_last_done"`
	TriggerCount       int32     `json:"trigger_count"`
	// Other
	LotSize            string    `json:"lot_size"`
	LimitOffset        string    `json:"limit_offset"`
	LimitDepthLevel    int32     `json:"limit_depth_level"`
	MarketPrice        string    `json:"market_price"`
	LastDone           string    `json:"last_done"`
	MonitorPrice             string    `json:"monitor_price"`
	SubmittedAmount          string    `json:"submitted_amount"`
	PlatformDeductionsStatus int32     `json:"platform_deductions_status"`
	PloyID                   string    `json:"ploy_id"`
	PloyType                 string    `json:"ploy_type"`
	TickerSize               string    `json:"ticker_size"`
	CurrentMillisecond       string    `json:"current_millisecond"`
	OrgID                    string    `json:"org_id"`
	Tag                      int32     `json:"tag"`
	ForceOnlyRTH       int32     `json:"force_only_rth"`
	DeductionsStatus   int32     `json:"deductions_status"`
	FreeStatus         int32     `json:"free_status"`
	Trend              int32     `json:"trend"`
}

// usRawOrder is the raw deserialization shape before symbol conversion.
type usRawOrder struct {
	ID                  string `json:"id"`
	AAID                string `json:"aaid"`
	AccountChannel      string `json:"account_channel"`
	Action              int32  `json:"action"`
	CounterID           string `json:"counter_id"`
	UnderlyingCounterID string `json:"underlying_counter_id"`
	Code                string `json:"code"`
	Name                string `json:"name"`
	SecurityType        string `json:"security_type"`
	Currency            string `json:"currency"`
	TradeCurrency       string `json:"trade_currency"`
	OrderType           string `json:"order_type"`
	Status              string `json:"status"`
	Price               string `json:"price"`
	Quantity            string `json:"quantity"`
	ExecutedQty         string `json:"executed_qty"`
	ExecutedPrice       string `json:"executed_price"`
	ExecutedAmount      string `json:"executed_amount"`
	OperateDirection    string `json:"operate_direction"`
	TimeInForce         int32  `json:"time_in_force"`
	GTD                 string `json:"gtd"`
	SubmittedAt         string `json:"submitted_at"`
	UpdatedAt           string `json:"updated_at"`
	Msg                 string `json:"msg"`
	Report              string `json:"report"`
	ContractDirection   string `json:"contract_direction"`
	ContractDueDate     string `json:"contract_due_date"`
	StrikePrice         string `json:"strike_price"`
	TailingAmount       string `json:"tailing_amount"`
	TailingPercent      string `json:"tailing_percent"`
	TriggerPrice        string `json:"trigger_price"`
	TriggerStatus       int32  `json:"trigger_status"`
	TriggerAt           string `json:"trigger_at"`
	TriggerExchange     string `json:"trigger_exchange"`
	TriggerLastDone     string `json:"trigger_last_done"`
	TriggerCount        int32  `json:"trigger_count"`
	LotSize             string `json:"lot_size"`
	LimitOffset         string `json:"limit_offset"`
	LimitDepthLevel     int32  `json:"limit_depth_level"`
	MarketPrice         string `json:"market_price"`
	LastDone            string `json:"last_done"`
	MonitorPrice              string `json:"monitor_price"`
	SubmittedAmount           string `json:"submitted_amount"`
	PlatformDeductionsStatus  int32  `json:"platform_deductions_status"`
	PloyID                    string `json:"ploy_id"`
	PloyType                  string `json:"ploy_type"`
	TickerSize                string `json:"ticker_size"`
	CurrentMillisecond        string `json:"current_millisecond"`
	OrgID                     string `json:"org_id"`
	Tag                       int32  `json:"tag"`
	ForceOnlyRTH        int32  `json:"force_only_rth"`
	DeductionsStatus    int32  `json:"deductions_status"`
	FreeStatus          int32  `json:"free_status"`
	Trend               int32  `json:"trend"`
}

// QueryUSOrdersResponse holds the paged list of US orders.
type QueryUSOrdersResponse struct {
	Orders     []USOrder `json:"orders"`
	TotalCount int32     `json:"total_count"`
}

// usRawQueryUSOrdersResponse is the raw deserialization shape.
type usRawQueryUSOrdersResponse struct {
	Orders     []usRawOrder `json:"orders"`
	TotalCount int32        `json:"total_count"`
}

// USOrderHistory is one entry in the order history list.
type USOrderHistory struct {
	ExecType  int    `json:"exec_type"`
	Status    string `json:"status"`
	Price     string `json:"price"`
	Qty       string `json:"qty"`
	Time      string `json:"time"`
	Msg       string `json:"msg"`
}

// USOrderDetailResponse is the response for USOrderDetail.
// Order contains the full order object (counter_id is NOT converted to symbol here;
// use QueryUSOrders for converted Symbol fields).
// OrderHistories contains state transitions.
type USOrderDetailResponse struct {
	Order                map[string]interface{} `json:"order"`
	OrderHistories       []USOrderHistory       `json:"order_histories"`
	CurrentAttachedOrder map[string]interface{} `json:"current_attached_order"`
}

// ── USAssetOverview ────────────────────────────────────────────────────────
// Note: only cash and crypto positions are exposed; stock/option/multi-leg
// positions from the full API response are not currently captured.

// USCashEntry is one currency cash entry in USAssetOverview.
type USCashEntry struct {
	Currency      string `json:"currency"`
	FrozenBuyCash string `json:"frozen_buy_cash"`
	Outstanding   string `json:"outstanding"`
	SettledCash   string `json:"settled_cash"`
	TotalAmount   string `json:"total_amount"`
	TotalCash     string `json:"total_cash"`
}

// USCryptoEntry is one cryptocurrency holding in USAssetOverview.
type USCryptoEntry struct {
	AssetType   string `json:"asset_type"`
	AverageCost string `json:"average_cost"`
	// Symbol is the user-facing trading-pair symbol (e.g. "BTCUSD.BKKT"),
	// converted from the API's counter_id field (e.g. "VA/BKKT/BTCUSD").
	Symbol      string `json:"symbol"`
	Currency    string `json:"currency"`
	IndustryName string `json:"industry_name"`
}

// usRawCryptoEntry is the raw API shape before symbol conversion.
type usRawCryptoEntry struct {
	AssetType   string `json:"asset_type"`
	AverageCost string `json:"average_cost"`
	CounterID   string `json:"counter_id"`
	Currency    string `json:"currency"`
	IndustryName string `json:"industry_name"`
}

// usRawAssetOverview is the raw API shape before field conversion.
type usRawAssetOverview struct {
	AccountType    string             `json:"account_type"`
	AssetTimestamp string             `json:"asset_timestamp"`
	CashBuyPower   string             `json:"cash_buy_power"`
	CashList       []USCashEntry      `json:"cash_list"`
	CryptoList     []usRawCryptoEntry `json:"crypto_list"`
}

// USAssetOverview is the US account asset snapshot.
type USAssetOverview struct {
	AccountType    string          `json:"account_type"`
	AssetTimestamp time.Time       `json:"asset_timestamp"`
	CashBuyPower   string          `json:"cash_buy_power"`
	CashList       []USCashEntry   `json:"cash_list"`
	CryptoList     []USCryptoEntry `json:"crypto_list"`
}

// ── USRealizedPL ───────────────────────────────────────────────────────────

// USRealizedPLMetric is one time-period metric within a USRealizedPLEntry.
type USRealizedPLMetric struct {
	Amount string `json:"amount"`
	Period int32  `json:"period"`
	Rate   string `json:"rate"`
}

// USRealizedPLEntry is one asset-category entry in USRealizedPL.
// Category values: 0=all, 1=stock, 2=option, 3=crypto.
type USRealizedPLEntry struct {
	Category int32                `json:"category"`
	Currency string               `json:"currency"`
	Metrics  []USRealizedPLMetric `json:"metrics"`
}

// USRealizedPL is the response for GET /v1/us/assets/pl/realized.
type USRealizedPL struct {
	RealizedPLList []USRealizedPLEntry `json:"realized_pl_list"`
}
