package jsontypes

type MarketStatusResponse struct {
	MarketTime []*MarketTimeItem `json:"market_time"`
}

type MarketTimeItem struct {
	Market            string `json:"market"`
	TradeStatus       int32  `json:"trade_status"`
	Timestamp         string `json:"timestamp"`
	DelayTradeStatus  int32  `json:"delay_trade_status"`
	DelayTimestamp    string `json:"delay_timestamp"`
	SubStatus         int32  `json:"sub_status"`
	DelaySubStatus    int32  `json:"delay_sub_status"`
}

type BrokerHoldingTop struct {
	Buy       []*BrokerHoldingEntry `json:"buy"`
	Sell      []*BrokerHoldingEntry `json:"sell"`
	UpdatedAt string                `json:"updated_at"`
}

type BrokerHoldingEntry struct {
	Name        string `json:"name"`
	PartiNumber string `json:"parti_number"`
	Chg         string `json:"chg"`
	Strong      bool   `json:"strong"`
}

type BrokerHoldingDetail struct {
	List      []*BrokerHoldingDetailItem `json:"list"`
	UpdatedAt string                     `json:"updated_at"`
}

type BrokerHoldingDetailItem struct {
	Name        string               `json:"name"`
	PartiNumber string               `json:"parti_number"`
	Ratio       *BrokerHoldingChanges `json:"ratio"`
	Shares      *BrokerHoldingChanges `json:"shares"`
	Strong      bool                 `json:"strong"`
}

type BrokerHoldingChanges struct {
	Value string `json:"value"`
	Chg1  string `json:"chg_1"`
	Chg5  string `json:"chg_5"`
	Chg20 string `json:"chg_20"`
	Chg60 string `json:"chg_60"`
}

type BrokerHoldingDailyHistory struct {
	List []*BrokerHoldingDailyItem `json:"list"`
}

type BrokerHoldingDailyItem struct {
	Date    string `json:"date"`
	Holding string `json:"holding"`
	Ratio   string `json:"ratio"`
	Chg     string `json:"chg"`
}

type AhPremiumKlines struct {
	Klines []*AhPremiumKline `json:"klines"`
}

type AhPremiumKline struct {
	Aprice       string `json:"aprice"`
	Apreclose    string `json:"apreclose"`
	Hprice       string `json:"hprice"`
	Hpreclose    string `json:"hpreclose"`
	CurrencyRate string `json:"currency_rate"`
	AhpremiumRate string `json:"ahpremium_rate"`
	PriceSpread  string `json:"price_spread"`
	Timestamp    string `json:"timestamp"`
}

type TradeStatsResponse struct {
	Statistics *TradeStatistics   `json:"statistics"`
	Trades     []*TradePriceLevel `json:"trades"`
}

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

type TradePriceLevel struct {
	BuyAmount     string `json:"buy_amount"`
	NeutralAmount string `json:"neutral_amount"`
	Price         string `json:"price"`
	SellAmount    string `json:"sell_amount"`
}

type AnomalyResponse struct {
	AllOff  bool           `json:"all_off"`
	Changes []*AnomalyItem `json:"changes"`
}

type AnomalyItem struct {
	CounterId    string   `json:"counter_id"`
	Name         string   `json:"name"`
	AlertName    string   `json:"alert_name"`
	AlertTime    int64    `json:"alert_time"`
	ChangeValues []string `json:"change_values"`
	Emotion      int32    `json:"emotion"`
}

type IndexConstituents struct {
	FallNum int32              `json:"fall_num"`
	FlatNum int32              `json:"flat_num"`
	RiseNum int32              `json:"rise_num"`
	Stocks  []*ConstituentStock `json:"stocks"`
}

type ConstituentStock struct {
	CounterId         string   `json:"counter_id"`
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
