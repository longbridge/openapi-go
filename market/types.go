package market

import "github.com/shopspring/decimal"

// BrokerHoldingPeriod is the lookback period for broker holding data.
type BrokerHoldingPeriod string

const (
	BrokerHoldingPeriodRct1  BrokerHoldingPeriod = "rct_1"
	BrokerHoldingPeriodRct5  BrokerHoldingPeriod = "rct_5"
	BrokerHoldingPeriodRct20 BrokerHoldingPeriod = "rct_20"
	BrokerHoldingPeriodRct60 BrokerHoldingPeriod = "rct_60"
)

// AhPremiumPeriod is the K-line period for A/H premium data.
type AhPremiumPeriod string

const (
	AhPremiumPeriodMin1  AhPremiumPeriod = "1"
	AhPremiumPeriodMin5  AhPremiumPeriod = "5"
	AhPremiumPeriodMin15 AhPremiumPeriod = "15"
	AhPremiumPeriodMin30 AhPremiumPeriod = "30"
	AhPremiumPeriodMin60 AhPremiumPeriod = "60"
	AhPremiumPeriodDay   AhPremiumPeriod = "1000"
	AhPremiumPeriodWeek  AhPremiumPeriod = "2000"
	AhPremiumPeriodMonth AhPremiumPeriod = "3000"
	AhPremiumPeriodYear  AhPremiumPeriod = "4000"
)

type MarketStatusResponse struct {
	MarketTime []*MarketTimeItem
}

type MarketTimeItem struct {
	Market           string
	TradeStatus      int32
	Timestamp        string
	DelayTradeStatus int32
	DelayTimestamp   string
	SubStatus        int32
	DelaySubStatus   int32
}

type BrokerHoldingTop struct {
	Buy       []*BrokerHoldingEntry
	Sell      []*BrokerHoldingEntry
	UpdatedAt string
}

type BrokerHoldingEntry struct {
	Name        string
	PartiNumber string
	Chg         *decimal.Decimal
	Strong      bool
}

type BrokerHoldingDetail struct {
	List      []*BrokerHoldingDetailItem
	UpdatedAt string
}

type BrokerHoldingDetailItem struct {
	Name        string
	PartiNumber string
	Ratio       *BrokerHoldingChanges
	Shares      *BrokerHoldingChanges
	Strong      bool
}

type BrokerHoldingChanges struct {
	Value *decimal.Decimal
	Chg1  *decimal.Decimal
	Chg5  *decimal.Decimal
	Chg20 *decimal.Decimal
	Chg60 *decimal.Decimal
}

type BrokerHoldingDailyHistory struct {
	List []*BrokerHoldingDailyItem
}

type BrokerHoldingDailyItem struct {
	Date    string
	Holding *decimal.Decimal
	Ratio   *decimal.Decimal
	Chg     *decimal.Decimal
}

type AhPremiumKlines struct {
	Klines []*AhPremiumKline
}

type AhPremiumKline struct {
	Aprice        *decimal.Decimal
	Apreclose     *decimal.Decimal
	Hprice        *decimal.Decimal
	Hpreclose     *decimal.Decimal
	CurrencyRate  *decimal.Decimal
	AhpremiumRate *decimal.Decimal
	PriceSpread   *decimal.Decimal
	Timestamp     string
}

type TradeStatsResponse struct {
	Statistics *TradeStatistics
	Trades     []*TradePriceLevel
}

type TradeStatistics struct {
	Avgprice    *decimal.Decimal
	Buy         *decimal.Decimal
	Neutral     *decimal.Decimal
	Preclose    *decimal.Decimal
	Sell        *decimal.Decimal
	Timestamp   string
	TotalAmount *decimal.Decimal
	TradeDate   []string
	TradesCount string
}

type TradePriceLevel struct {
	BuyAmount     *decimal.Decimal
	NeutralAmount *decimal.Decimal
	Price         *decimal.Decimal
	SellAmount    *decimal.Decimal
}

type AnomalyResponse struct {
	AllOff  bool
	Changes []*AnomalyItem
}

type AnomalyItem struct {
	Symbol       string
	Name         string
	AlertName    string
	AlertTime    int64
	ChangeValues []string
	Emotion      int32
}

type IndexConstituents struct {
	FallNum int32
	FlatNum int32
	RiseNum int32
	Stocks  []*ConstituentStock
}

type ConstituentStock struct {
	Symbol            string
	Name              string
	LastDone          *decimal.Decimal
	PrevClose         *decimal.Decimal
	Inflow            *decimal.Decimal
	Balance           *decimal.Decimal
	Amount            *decimal.Decimal
	TotalShares       *decimal.Decimal
	Tags              []string
	Intro             string
	Market            string
	CirculatingShares *decimal.Decimal
	Delay             bool
	Chg               *decimal.Decimal
	TradeStatus       int32
}
