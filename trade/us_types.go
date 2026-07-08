package trade

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/longbridge/openapi-go/internal/counter"
)

// unixTimestamp is a time.Time that unmarshals from JSON unix-second values
// expressed either as a JSON number (1783455324) or as a decimal string
// ("1783455324"). RFC3339 strings are accepted as a fallback. An empty string,
// "0", or null maps to the zero time.Time.
type unixTimestamp time.Time

func (t *unixTimestamp) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Strip JSON string quotes if present.
	if len(s) >= 2 && s[0] == '"' {
		s = s[1 : len(s)-1]
	}
	if s == "" || s == "null" || s == "0" {
		*t = unixTimestamp(time.Time{})
		return nil
	}
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		*t = unixTimestamp(time.Unix(n, 0).UTC())
		return nil
	}
	if parsed, err := time.Parse(time.RFC3339, s); err == nil {
		*t = unixTimestamp(parsed.UTC())
		return nil
	}
	*t = unixTimestamp(time.Time{})
	return nil
}

func (t unixTimestamp) Time() time.Time { return time.Time(t) }

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
	SubmittedAt         unixTimestamp `json:"submitted_at"`
	UpdatedAt           unixTimestamp `json:"updated_at"`
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

// USOrderHistory is one state-transition entry within USOrderDetail.OrderHistories.
type USOrderHistory struct {
	Price            string `json:"price"`
	Qty              string `json:"qty"`
	Status           string `json:"status"`
	Msg              string `json:"msg"`
	Time             string `json:"time"` // unix seconds as string
	IsManually       bool   `json:"is_manually"`
	ExecType         int32  `json:"exec_type"`
	OppPartyID       string `json:"opp_party_id"`
	TrdMatchID       string `json:"trd_match_id"`
	Operator         string `json:"operator"`
	OpEntrustWay     string `json:"op_entrust_way"`
	CxlRejResponseTo int32  `json:"cxl_rej_response_to"`
	WithdrawalReason string `json:"withdrawal_reason"`
	OppName          string `json:"opp_name"`
	ExecID           string `json:"exec_id"`
}

// USButtonControl holds the action-button state for an order.
type USButtonControl struct {
	Withdraw      int32    `json:"withdraw"`
	Replace       int32    `json:"replace"`
	Exceptionable []string `json:"exceptionable"`
}

// USChargeItem is one fee category within USChargeDetail.
type USChargeItem struct {
	Code int32    `json:"code"`
	Name string   `json:"name"`
	Fees []string `json:"fees"`
}

// USChargeDetail holds the fee breakdown for an order.
type USChargeDetail struct {
	Currency    string        `json:"currency"`
	TotalAmount string        `json:"total_amount"`
	Items       []USChargeItem `json:"items"`
}

// USAttachedOrder is one bracket/conditional order attached to the main order.
type USAttachedOrder struct {
	AttachedTypeDisplay int32  `json:"attached_type_display"`
	ExecutedQty         string `json:"executed_qty"`
	Quantity            string `json:"quantity"`
	Status              string `json:"status"`
	TriggerPrice        string `json:"trigger_price"`
	OrderID             string `json:"order_id"`
	GTD                 string `json:"gtd"`
	TimeInForce         int32  `json:"time_in_force"`
	Tag                 int32  `json:"tag"`
	ActivateOrderType   string `json:"activate_order_type"`
	ActivateRTH         int32  `json:"activate_rth"`
	SubmitPrice         string `json:"submit_price"`
	// Symbol is the user-facing trading symbol (e.g. "NKE.US"), converted from CounterID.
	Symbol    string `json:"-"`
	CounterID string `json:"counter_id"`
	Withdrawn bool   `json:"withdrawn"`
}

// USOrderDetail is the full typed order object within USOrderDetailResponse.
// submitted_at and done_at are raw unix-second strings (not converted to time.Time,
// since this struct is also reused for CurrentAttachedOrder which may vary).
type USOrderDetail struct {
	ID                         string           `json:"id"`
	AAID                       string           `json:"aaid"`
	AccountChannel             string           `json:"account_channel"`
	Action                     int32            `json:"action"`
	// Symbol is the user-facing trading symbol (e.g. "NKE.US"), converted from CounterID.
	Symbol                     string           `json:"-"`
	// UnderlyingSymbol is the user-facing underlying symbol (options only), converted from UnderlyingCounterID.
	UnderlyingSymbol           string           `json:"-"`
	CounterID                  string           `json:"counter_id"`
	UnderlyingCounterID        string           `json:"underlying_counter_id"`
	SecurityType               string           `json:"security_type"`
	Name                       string           `json:"name"`
	Currency                   string           `json:"currency"`
	TradeCurrency              string           `json:"trade_currency"`
	OrderType                  string           `json:"order_type"`
	Status                     string           `json:"status"`
	Price                      string           `json:"price"`
	Quantity                   string           `json:"quantity"`
	ExecutedQty                string           `json:"executed_qty"`
	ExecutedPrice              string           `json:"executed_price"`
	ExecutedAmount             string           `json:"executed_amount"`
	OperateDirection           string           `json:"operate_direction"`
	TimeInForce                int32            `json:"time_in_force"`
	GTD                        string           `json:"gtd"`
	Tag                        int32            `json:"tag"`
	Msg                        string           `json:"msg"`
	ForceOnlyRTH               int32            `json:"force_only_rth"`
	SubmittedAt                string           `json:"submitted_at"` // unix seconds string
	DoneAt                     string           `json:"done_at"`      // unix seconds string
	TriggerPrice               string           `json:"trigger_price"`
	TriggerAt                  string           `json:"trigger_at"`
	TriggerStatus              int32            `json:"trigger_status"`
	TriggerExchange            string           `json:"trigger_exchange"`
	TriggerLastDone            string           `json:"trigger_last_done"`
	TriggerCount               int32            `json:"trigger_count"`
	TailingAmount              string           `json:"tailing_amount"`
	TailingPercent             string           `json:"tailing_percent"`
	LimitOffset                string           `json:"limit_offset"`
	LimitDepthLevel            int32            `json:"limit_depth_level"`
	MarketPrice                string           `json:"market_price"`
	SubmittedAmount            string           `json:"submitted_amount"`
	EstimatedFee               string           `json:"estimated_fee"`
	FreeStatus                 int32            `json:"free_status"`
	FreeAmount                 string           `json:"free_amount"`
	FreeCurrency               string           `json:"free_currency"`
	DeductionsStatus           int32            `json:"deductions_status"`
	DeductionsAmount           string           `json:"deductions_amount"`
	DeductionsCurrency         string           `json:"deductions_currency"`
	PlatformDeductionsStatus   int32            `json:"platform_deductions_status"`
	PlatformDeductionsAmount   string           `json:"platform_deductions_amount"`
	PlatformDeductionsCurrency string           `json:"platform_deductions_currency"`
	DisplayAccount             string           `json:"display_account"`
	SettlementAccount          string           `json:"settlement_account"`
	SettlementChannel          string           `json:"settlement_channel"`
	CustomerName               string           `json:"customer_name"`
	RealName                   string           `json:"real_name"`
	EnName                     string           `json:"en_name"`
	JointRealName              string           `json:"joint_real_name"`
	JointEnName                string           `json:"joint_en_name"`
	OrgID                      string           `json:"org_id"`
	BCAN                       string           `json:"bcan"`
	OpEntrustWay               int32            `json:"op_entrust_way"`
	OpEntrustWayName           string           `json:"op_entrust_way_name"`
	Remark                     string           `json:"remark"`
	Notice                     string           `json:"notice"`
	ShortSellType              int32            `json:"short_sell_type"`
	PloyType                   string           `json:"ploy_type"` // API returns string e.g. "0"
	PloyID                     string           `json:"ploy_id"`
	PloyStatus                 string           `json:"ploy_status"`
	Trend                      int32            `json:"trend"`
	WithdrawalReason           string           `json:"withdrawal_reason"`
	ActivateOrderType          string           `json:"activate_order_type"`
	ActivateRTH                int32            `json:"activate_rth"`
	SubmitPrice                string           `json:"submit_price"`
	ContractDirection          string           `json:"contract_direction"`
	StrikePrice                string           `json:"strike_price"`
	ContractSize               string           `json:"contract_size"`
	MonitorPrice               string           `json:"monitor_price"`
	ButtonControl              USButtonControl   `json:"button_control"`
	ChargeDetail               *USChargeDetail   `json:"charge_detail"`
	AttachedOrders             []USAttachedOrder `json:"attached_orders"`
	OrderHistories             []USOrderHistory  `json:"order_histories"`
}

// UnmarshalJSON converts counter_id / underlying_counter_id fields to
// user-facing Symbol / UnderlyingSymbol after standard JSON deserialization.
func (o *USOrderDetail) UnmarshalJSON(b []byte) error {
	type raw USOrderDetail
	if err := json.Unmarshal(b, (*raw)(o)); err != nil {
		return err
	}
	o.Symbol = counter.IDToSymbol(o.CounterID)
	o.UnderlyingSymbol = counter.IDToSymbol(o.UnderlyingCounterID)
	for i := range o.AttachedOrders {
		o.AttachedOrders[i].Symbol = counter.IDToSymbol(o.AttachedOrders[i].CounterID)
	}
	return nil
}

// USOrderDetailResponse is the response for USOrderDetail.
// CurrentAttachedOrder is the active bracket/conditional sub-order, or nil.
// CurrentMillisecond is the server timestamp at response time.
type USOrderDetailResponse struct {
	Order                *USOrderDetail `json:"order"`
	CurrentAttachedOrder *USOrderDetail `json:"current_attached_order"`
	CurrentMillisecond   string         `json:"current_millisecond"`
}

// ── USAssetOverview ────────────────────────────────────────────────────────

// USStockEntry is one stock/equity position in USAssetOverview.
type USStockEntry struct {
	// Symbol is the ticker code returned by the API (e.g. "AAPL"). See FullSymbol for the qualified form.
	Symbol                     string `json:"symbol"`
	// FullSymbol is the user-facing qualified symbol (e.g. "AAPL.US"), converted from CounterID.
	FullSymbol                 string `json:"-"`
	AssetType                  string `json:"asset_type"`
	Quantity                   string `json:"quantity"`
	Currency                   string `json:"currency"`
	AverageCost                string `json:"average_cost"`
	Market                     string `json:"market"`
	CounterID                  string `json:"counter_id"`
	TradeStatus                string `json:"trade_status"`
	PrevClose                  string `json:"prev_close"`
	LastDone                   string `json:"last_done"`
	MarketPrice                string `json:"market_price"`
	PretradeClose              string `json:"pretrade_close"`
	StockInvestOfToday         string `json:"stock_invest_of_today"`
	TodayPL                    string `json:"today_pl"`
	PretradeStockInvestOfToday string `json:"pretrade_stock_invest_of_today"`
	PretradeTodayPL            string `json:"pretrade_today_pl"`
	NightLastDone              string `json:"night_last_done"`
	NightPrevClose             string `json:"night_prev_close"`
	PositionSide               string `json:"position_side"`
	OpenPositionTime           string `json:"open_position_time"`
	Name                       string `json:"name"`
	IndustryCounterID          string `json:"industry_counter_id"`
	IndustryName               string `json:"industry_name"`
}

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
	AccountType       string             `json:"account_type"`
	AssetTimestamp    string             `json:"asset_timestamp"`
	CashBuyPower      string             `json:"cash_buy_power"`
	OvernightBuyPower string             `json:"overnight_buy_power"`
	Currency          string             `json:"currency"`
	CashList          []USCashEntry      `json:"cash_list"`
	StockList         []USStockEntry     `json:"stock_list"`
	OptionList        []interface{}      `json:"option_list"`
	CryptoList        []usRawCryptoEntry `json:"crypto_list"`
	MultiLeg          interface{}        `json:"multi_leg"`
}

// USAssetOverview is the US account asset snapshot.
type USAssetOverview struct {
	AccountType       string          `json:"account_type"`
	AssetTimestamp    time.Time       `json:"-"`
	CashBuyPower      string          `json:"cash_buy_power"`
	OvernightBuyPower string          `json:"overnight_buy_power"`
	Currency          string          `json:"currency"`
	CashList          []USCashEntry   `json:"cash_list"`
	StockList         []USStockEntry  `json:"stock_list"`
	OptionList        []interface{}   `json:"option_list"`
	CryptoList        []USCryptoEntry `json:"crypto_list"`
	MultiLeg          interface{}     `json:"multi_leg"`
}

// ── USRealizedPL ───────────────────────────────────────────────────────────

// USRealizedPLMetric is one time-period metric within a USRealizedPLEntry.
type USRealizedPLMetric struct {
	Amount string `json:"amount"`
	Period int32  `json:"period"`
	Rate   string `json:"rate"`
}

// USRealizedPLEntry is one asset-category entry in USRealizedPLResponse.
// Category values: 0=all, 1=stock, 2=option, 3=crypto.
type USRealizedPLEntry struct {
	Category int32                `json:"category"`
	Currency string               `json:"currency"`
	Metrics  []USRealizedPLMetric `json:"metrics"`
}

// USRealizedPLResponse is the response for GET /v1/us/assets/pl/realized.
type USRealizedPLResponse struct {
	RealizedPLList []USRealizedPLEntry `json:"realized_pl_list"`
}
