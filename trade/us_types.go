package trade

import "errors"

// ErrUSOnly is returned when a US-only API is called with a non-US token.
var ErrUSOnly = errors.New("longbridge: this API is only available for US accounts (us_ token required)")

// QueryUSOrdersRequest is the request body for QueryUSOrders.
type QueryUSOrdersRequest struct {
	AccountChannel string   `json:"account_channel"`
	Action         int32    `json:"action"`
	StartAt        float64  `json:"start_at"`
	EndAt          float64  `json:"end_at"`
	CounterIDs     []string `json:"counter_ids"`
	SecurityTypes  []string `json:"security_types"`
	QueryType      int32    `json:"query_type"`
	Page           int32    `json:"page"`
	Limit          int32    `json:"limit"`
	QueryVersion   float64  `json:"query_version"`
}

// QueryUSOrdersResponse holds the paged list of US orders.
type QueryUSOrdersResponse struct {
	Orders []map[string]interface{} `json:"orders"`
}

// AttachedOrder is a take-profit or stop-loss sub-order attached to a US order.
type AttachedOrder struct {
	OrderID      string `json:"order_id"`
	Type         string `json:"type"`
	Side         string `json:"side"`
	Price        string `json:"price"`
	TrailAmount  string `json:"trail_amount"`
	TrailPercent string `json:"trail_percent"`
	Status       string `json:"status"`
}

// USOrderDetailResponse is the response for USOrderDetail.
// The raw order detail fields are passed through as-is; attached_orders is
// populated only when isAttached=true.
type USOrderDetailResponse struct {
	Raw            map[string]interface{} `json:"-"`
	AttachedOrders []AttachedOrder        `json:"attached_orders,omitempty"`
}

// USStockPosition is a stock holding in the US account.
type USStockPosition struct {
	Symbol             string `json:"symbol"`
	Name               string `json:"name"`
	Quantity           string `json:"quantity"`
	AvailableQuantity  string `json:"available_quantity"`
	Currency           string `json:"currency"`
	CostPrice          string `json:"cost_price"`
	MarketValue        string `json:"market_value"`
	UnrealizedPL       string `json:"unrealized_pl"`
	UnrealizedPLRatio  string `json:"unrealized_pl_ratio"`
	LastDone           string `json:"last_done"`
	PrevClose          string `json:"prev_close"`
	ChangeRate         string `json:"change_rate"`
	NightLastDone      string `json:"night_last_done"`
	PretradeClose      string `json:"pretrade_close"`
	TradeStatus        string `json:"trade_status"`
	IndividualQuantity string `json:"individual_quantity"`
}

// USOptionPosition is an option holding in the US account.
type USOptionPosition struct {
	Symbol             string `json:"symbol"`
	StrikePrice        string `json:"strike_price"`
	DueDate            string `json:"due_date"`
	ContractMultiplier int32  `json:"contract_multiplier"`
	Type               string `json:"type"`
	Quantity           string `json:"quantity"`
	MarketValue        string `json:"market_value"`
	UnrealizedPL       string `json:"unrealized_pl"`
}

// USCryptoPosition is a cryptocurrency holding in the US account.
type USCryptoPosition struct {
	Symbol       string `json:"symbol"`
	Quantity     string `json:"quantity"`
	MarketValue  string `json:"market_value"`
	UnrealizedPL string `json:"unrealized_pl"`
	CostPrice    string `json:"cost_price"`
}

// USBuyPower holds purchasing power breakdown for a US account.
type USBuyPower struct {
	CashBuyPower      string `json:"cash_buy_power"`
	OvernightBuyPower string `json:"overnight_buy_power"`
	DayTradeBuyPower  string `json:"day_trade_buy_power"`
	OptionBuyPower    string `json:"option_buy_power"`
	CryptoBuyPower    string `json:"crypto_buy_power"`
}

// USAssetOverview is the full US account asset snapshot.
type USAssetOverview struct {
	AccountType     string            `json:"account_type"`
	NetAssets       string            `json:"net_assets"`
	TotalCash       string            `json:"total_cash"`
	UnrealizedPL    string            `json:"unrealized_pl"`
	Positions       []USStockPosition `json:"positions"`
	OptionPositions []USOptionPosition `json:"option_positions"`
	MultiLegs       []map[string]interface{} `json:"multi_legs"`
	CryptoPositions []USCryptoPosition `json:"crypto_positions"`
	BuyPower        USBuyPower        `json:"buy_power"`
}

// RealizedPLItem is a single realized P&L entry by symbol.
type RealizedPLItem struct {
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	RealizedPL   string `json:"realized_pl"`
	QuantitySold string `json:"quantity_sold"`
	AvgCost      string `json:"avg_cost"`
	AvgSellPrice string `json:"avg_sell_price"`
}

// USRealizedPL is the response for USRealizedPL.
type USRealizedPL struct {
	TotalRealizedPL string           `json:"total_realized_pl"`
	Currency        string           `json:"currency"`
	Items           []RealizedPLItem `json:"items"`
}
