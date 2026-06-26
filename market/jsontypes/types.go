// Package jsontypes holds raw JSON-deserialization structs for the market package.
// These types use exact API field names with json tags and are not part of the
// public surface; callers should use the types in the parent market package.
package jsontypes

// MarketStatusResponse is the raw JSON response for GET /v1/quote/market-status.
type MarketStatusResponse struct {
	MarketTime []MarketTimeItem `json:"market_time"`
}

// MarketTimeItem is the raw JSON representation of one market's trading status.
type MarketTimeItem struct {
	Market           string `json:"market"`
	// Raw market trade status code. See market.TradeStatus for the code table.
	TradeStatus      int32  `json:"trade_status"`
	Timestamp        string `json:"timestamp"`
	// Raw delayed market trade status code.
	DelayTradeStatus int32  `json:"delay_trade_status"`
	DelayTimestamp   string `json:"delay_timestamp"`
	SubStatus        int32  `json:"sub_status"`
	DelaySubStatus   int32  `json:"delay_sub_status"`
}

// BrokerHoldingTop is the raw JSON response for GET /v1/quote/broker-holding.
type BrokerHoldingTop struct {
	Buy       []BrokerHoldingEntry `json:"buy"`
	Sell      []BrokerHoldingEntry `json:"sell"`
	UpdatedAt string               `json:"updated_at"`
}

// BrokerHoldingEntry is one broker entry in the top-holding list.
type BrokerHoldingEntry struct {
	Name        string `json:"name"`
	PartiNumber string `json:"parti_number"`
	Chg         string `json:"chg"`
	Strong      bool   `json:"strong"`
}

// BrokerHoldingDetail is the raw JSON response for GET /v1/quote/broker-holding/detail.
type BrokerHoldingDetail struct {
	List      []BrokerHoldingDetailItem `json:"list"`
	UpdatedAt string                    `json:"updated_at"`
}

// BrokerHoldingDetailItem is one broker's full holding detail.
type BrokerHoldingDetailItem struct {
	Name        string               `json:"name"`
	PartiNumber string               `json:"parti_number"`
	Ratio       BrokerHoldingChanges `json:"ratio"`
	Shares      BrokerHoldingChanges `json:"shares"`
	Strong      bool                 `json:"strong"`
}

// BrokerHoldingChanges holds the value and period-over-period changes for a holding metric.
type BrokerHoldingChanges struct {
	Value string `json:"value"`
	Chg1  string `json:"chg_1"`
	Chg5  string `json:"chg_5"`
	Chg20 string `json:"chg_20"`
	Chg60 string `json:"chg_60"`
}

// BrokerHoldingDailyHistory is the raw JSON response for GET /v1/quote/broker-holding/daily.
type BrokerHoldingDailyHistory struct {
	List []BrokerHoldingDailyItem `json:"list"`
}

// BrokerHoldingDailyItem is one day's broker holding record.
type BrokerHoldingDailyItem struct {
	Date    string `json:"date"`
	Holding string `json:"holding"`
	Ratio   string `json:"ratio"`
	Chg     string `json:"chg"`
}

// AhPremiumKlines is the raw JSON response for GET /v1/quote/ahpremium/klines.
type AhPremiumKlines struct {
	Klines []AhPremiumKline `json:"klines"`
}

// AhPremiumIntraday is the raw JSON response for GET /v1/quote/ahpremium/timeshares.
type AhPremiumIntraday struct {
	Klines []AhPremiumKline `json:"klines"`
}

// AhPremiumKline is one A/H premium data point.
type AhPremiumKline struct {
	Aprice       string `json:"aprice"`
	Apreclose    string `json:"apreclose"`
	Hprice       string `json:"hprice"`
	Hpreclose    string `json:"hpreclose"`
	CurrencyRate string `json:"currency_rate"`
	AhpremiumRate string `json:"ahpremium_rate"`
	PriceSpread  string `json:"price_spread"`
	Timestamp    int64  `json:"timestamp"`
}

// TradeStatsResponse is the raw JSON response for GET /v1/quote/trades-statistics.
type TradeStatsResponse struct {
	Statistics TradeStatistics  `json:"statistics"`
	Trades     []TradePriceLevel `json:"trades"`
}

// TradeStatistics holds summary buy/sell/neutral statistics.
type TradeStatistics struct {
	Avgprice    string   `json:"avgprice"`
	Buy         string   `json:"buy"`
	Neutral     string   `json:"neutral"`
	Preclose    string   `json:"preclose"`
	Sell        string   `json:"sell"`
	Timestamp   string   `json:"timestamp"`
	TotalAmount string   `json:"total_amount"`
	TradeDate   []string `json:"trade_date"`
	TradesCount string   `json:"trades_count"`
}

// TradePriceLevel is trade volume at one price level.
type TradePriceLevel struct {
	BuyAmount     string `json:"buy_amount"`
	NeutralAmount string `json:"neutral_amount"`
	Price         string `json:"price"`
	SellAmount    string `json:"sell_amount"`
}

// AnomalyResponse is the raw JSON response for GET /v1/quote/changes.
type AnomalyResponse struct {
	AllOff  bool          `json:"all_off"`
	Changes []AnomalyItem `json:"changes"`
}

// AnomalyItem is one market anomaly event.
type AnomalyItem struct {
	CounterID    string   `json:"counter_id"`
	Name         string   `json:"name"`
	AlertName    string   `json:"alert_name"`
	AlertTime    int64    `json:"alert_time"`
	ChangeValues []string `json:"change_values"`
	Emotion      int32    `json:"emotion"`
}

// IndexConstituents is the raw JSON response for GET /v1/quote/index-constituents.
type IndexConstituents struct {
	FallNum int32              `json:"fall_num"`
	FlatNum int32              `json:"flat_num"`
	RiseNum int32              `json:"rise_num"`
	Stocks  []ConstituentStock `json:"stocks"`
}

// ConstituentStock is one constituent stock of an index.
type ConstituentStock struct {
	CounterID         string   `json:"counter_id"`
	Name              string   `json:"name"`
	LastDone          string   `json:"last_done"`
	PrevClose         string   `json:"prev_close"`
	Inflow            string   `json:"inflow"`
	Balance           string   `json:"balance"`
	Amount            string   `json:"amount"`
	TotalShares       string   `json:"total_shares"`
	Tags              []string `json:"tags"`
	Intro             string   `json:"intro"`
	Market            string   `json:"market"`
	CirculatingShares string   `json:"circulating_shares"`
	Delay             bool     `json:"delay"`
	Chg               string   `json:"chg"`
	TradeStatus       int32    `json:"trade_status"`
}
