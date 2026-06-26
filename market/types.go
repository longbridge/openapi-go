package market

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// RankListItem is one item in the popularity rank list.
type RankListItem struct {
	// Symbol — converted from counter_id (e.g. "MU.US")
	Symbol string
	// Code — ticker code (e.g. "MU")
	Code string
	Name         string
	LastDone     string
	// Chg — price change as decimal ratio (e.g. 0.0252 = +2.52%)
	Chg          string
	// Change — absolute price change
	Change       string
	Inflow       string
	MarketCap    string
	Industry     string
	PrePostPrice string
	PrePostChg   string
	Amplitude    string
	FiveDayChg   string
	TurnoverRate string
	VolumeRate   string
	PbTtm        string
}

// TopMoversStock holds stock info for a top-movers event.
type TopMoversStock struct {
	// Symbol — converted from counter_id
	Symbol   string
	Code     string
	Name     string
	FullName string
	Change   string
	LastDone string
	Market   string
	Labels   []string
	Logo     string
}

// TopMoversEvent is one top-movers event.
type TopMoversEvent struct {
	// Timestamp — RFC 3339
	Timestamp   string
	AlertReason string
	AlertType   int64
	Stock       TopMoversStock
	// Post — associated news article (raw JSON, complex structure; nil when no news)
	Post json.RawMessage
}

// MarketStatusResponse holds the current trading status for all markets.
type MarketStatusResponse struct {
	// Per-market trading status items
	MarketTime []MarketTimeItem
}

// MarketTimeItem is the trading status for one market.
type MarketTimeItem struct {
	// Market code, e.g. "HK", "US", "CN"
	Market string
	// Market trade status. See TradeStatus for the code table.
	TradeStatus TradeStatus
	// Current market time
	Timestamp time.Time
	// Delayed-quote market trade status. See TradeStatus for the code table.
	DelayTradeStatus TradeStatus
	// Delayed-quote market time
	DelayTimestamp time.Time
	// Sub-status code
	SubStatus int32
	// Delayed-quote sub-status code
	DelaySubStatus int32
}

// BrokerHoldingPeriod is the lookback period for broker holding queries.
type BrokerHoldingPeriod int

const (
	// BrokerHoldingPeriodRct1 is the 1-day lookback period.
	BrokerHoldingPeriodRct1 BrokerHoldingPeriod = iota
	// BrokerHoldingPeriodRct5 is the 5-day lookback period.
	BrokerHoldingPeriodRct5
	// BrokerHoldingPeriodRct20 is the 20-day lookback period.
	BrokerHoldingPeriodRct20
	// BrokerHoldingPeriodRct60 is the 60-day lookback period.
	BrokerHoldingPeriodRct60
)

// BrokerHoldingTop holds the top broker buy/sell leaders for a security.
type BrokerHoldingTop struct {
	// Top brokers by net buying
	Buy []BrokerHoldingEntry
	// Top brokers by net selling
	Sell []BrokerHoldingEntry
	// Last updated (may be zero if not provided)
	UpdatedAt string
}

// BrokerHoldingEntry is one broker entry in a top-holding list.
type BrokerHoldingEntry struct {
	// Broker name
	Name string
	// Participant number / broker code
	PartiNumber string
	// Net change in shares held (nil if not available)
	Chg *decimal.Decimal
	// Whether this is a "strengthening" broker
	Strong bool
}

// BrokerHoldingDetail is the full broker holding detail for a security.
type BrokerHoldingDetail struct {
	// Full list of broker holdings
	List []BrokerHoldingDetailItem
	// Last updated (may be empty if not provided)
	UpdatedAt string
}

// BrokerHoldingDetailItem is one broker's full holding detail.
type BrokerHoldingDetailItem struct {
	// Broker name
	Name string
	// Participant number / broker code
	PartiNumber string
	// Holding ratio changes over various periods
	Ratio BrokerHoldingChanges
	// Share count changes over various periods
	Shares BrokerHoldingChanges
	// Whether this is a "strengthening" broker
	Strong bool
}

// BrokerHoldingChanges holds the value and period-over-period changes for a holding metric.
type BrokerHoldingChanges struct {
	// Current value
	Value *decimal.Decimal
	// 1-day change
	Chg1 *decimal.Decimal
	// 5-day change
	Chg5 *decimal.Decimal
	// 20-day change
	Chg20 *decimal.Decimal
	// 60-day change
	Chg60 *decimal.Decimal
}

// BrokerHoldingDailyHistory is the daily holding history for a specific broker.
type BrokerHoldingDailyHistory struct {
	// Daily broker holding records
	List []BrokerHoldingDailyItem
}

// BrokerHoldingDailyItem is one day's broker holding record.
type BrokerHoldingDailyItem struct {
	// Date in "2026.05.05" format
	Date string
	// Total shares held (nil if not available)
	Holding *decimal.Decimal
	// Holding ratio as a decimal (nil if not available)
	Ratio *decimal.Decimal
	// Change vs previous day (nil if not available)
	Chg *decimal.Decimal
}

// AhPremiumPeriod is the K-line period for A/H premium queries.
type AhPremiumPeriod int

const (
	// AhPremiumPeriodMin1 is the 1-minute period.
	AhPremiumPeriodMin1 AhPremiumPeriod = iota
	// AhPremiumPeriodMin5 is the 5-minute period.
	AhPremiumPeriodMin5
	// AhPremiumPeriodMin15 is the 15-minute period.
	AhPremiumPeriodMin15
	// AhPremiumPeriodMin30 is the 30-minute period.
	AhPremiumPeriodMin30
	// AhPremiumPeriodMin60 is the 60-minute period.
	AhPremiumPeriodMin60
	// AhPremiumPeriodDay is the daily period.
	AhPremiumPeriodDay
	// AhPremiumPeriodWeek is the weekly period.
	AhPremiumPeriodWeek
	// AhPremiumPeriodMonth is the monthly period.
	AhPremiumPeriodMonth
	// AhPremiumPeriodYear is the yearly period.
	AhPremiumPeriodYear
)

// lineType converts the AhPremiumPeriod to the API's line_type parameter value.
func (p AhPremiumPeriod) lineType() string {
	switch p {
	case AhPremiumPeriodMin1:
		return "1"
	case AhPremiumPeriodMin5:
		return "5"
	case AhPremiumPeriodMin15:
		return "15"
	case AhPremiumPeriodMin30:
		return "30"
	case AhPremiumPeriodMin60:
		return "60"
	case AhPremiumPeriodDay:
		return "1000"
	case AhPremiumPeriodWeek:
		return "2000"
	case AhPremiumPeriodMonth:
		return "3000"
	case AhPremiumPeriodYear:
		return "4000"
	default:
		return "1000"
	}
}

// AhPremiumKlines holds A/H premium K-line data for a dual-listed security.
type AhPremiumKlines struct {
	// K-line data points
	Klines []AhPremiumKline
}

// AhPremiumIntraday holds A/H premium intraday data for a dual-listed security.
type AhPremiumIntraday struct {
	// Intraday data points
	Klines []AhPremiumKline
}

// AhPremiumKline is one A/H premium data point.
type AhPremiumKline struct {
	// A-share price
	Aprice decimal.Decimal
	// A-share previous close
	Apreclose decimal.Decimal
	// H-share price
	Hprice decimal.Decimal
	// H-share previous close
	Hpreclose decimal.Decimal
	// CNY/HKD exchange rate
	CurrencyRate decimal.Decimal
	// A/H premium rate (negative = H-share at premium)
	AhpremiumRate decimal.Decimal
	// Price spread
	PriceSpread decimal.Decimal
	// Data point timestamp
	Timestamp time.Time
}

// TradeStatsResponse holds buy/sell/neutral trade statistics for a security.
type TradeStatsResponse struct {
	// Summary statistics
	Statistics TradeStatistics
	// Per-price-level breakdown
	Trades []TradePriceLevel
}

// TradeStatistics holds summary buy/sell/neutral trade statistics.
type TradeStatistics struct {
	// Volume-weighted average price
	Avgprice decimal.Decimal
	// Total buy volume (shares)
	Buy decimal.Decimal
	// Total neutral / unknown-direction volume
	Neutral decimal.Decimal
	// Previous close price
	Preclose decimal.Decimal
	// Total sell volume (shares)
	Sell decimal.Decimal
	// Data timestamp (unix timestamp string, raw from API)
	Timestamp string
	// Total trading volume (shares)
	TotalAmount decimal.Decimal
	// Unix timestamps for the last N trading days (raw strings)
	TradeDate []string
	// Total number of trades (raw string from API)
	TradesCount string
}

// TradePriceLevel holds trade volume at one price level.
type TradePriceLevel struct {
	// Buy volume at this price
	BuyAmount decimal.Decimal
	// Neutral (unknown direction) volume at this price
	NeutralAmount decimal.Decimal
	// Price level
	Price decimal.Decimal
	// Sell volume at this price
	SellAmount decimal.Decimal
}

// AnomalyResponse holds market anomaly alerts for a market.
type AnomalyResponse struct {
	// Whether anomaly alerts are globally disabled
	AllOff bool
	// List of market anomaly events
	Changes []AnomalyItem
}

// AnomalyItem is one market anomaly event (e.g. large block trade, margin buying surge).
type AnomalyItem struct {
	// Security symbol (e.g. "700.HK")
	Symbol string
	// Security name
	Name string
	// Anomaly type name, e.g. "大宗交易", "融资买入"
	AlertName string
	// Time of the anomaly (unix timestamp in milliseconds)
	AlertTime int64
	// Change values associated with the anomaly event
	ChangeValues []string
	// Sentiment direction: 1 = positive/up, 2 = negative/down
	Emotion int32
}

// IndexConstituents holds the constituent stocks for an index.
type IndexConstituents struct {
	// Number of constituent stocks that fell today
	FallNum int32
	// Number of constituent stocks unchanged today
	FlatNum int32
	// Number of constituent stocks that rose today
	RiseNum int32
	// Constituent stock details
	Stocks []ConstituentStock
}

// ConstituentStock is one constituent stock of an index.
type ConstituentStock struct {
	// Security symbol (e.g. "700.HK")
	Symbol string
	// Security name
	Name string
	// Latest price (nil if not available)
	LastDone *decimal.Decimal
	// Previous close (nil if not available)
	PrevClose *decimal.Decimal
	// Net capital inflow today (nil if not available)
	Inflow *decimal.Decimal
	// Turnover amount (nil if not available)
	Balance *decimal.Decimal
	// Trading volume in shares (nil if not available)
	Amount *decimal.Decimal
	// Total shares outstanding (nil if not available)
	TotalShares *decimal.Decimal
	// Tags, e.g. ["领涨龙头"]
	Tags []string
	// Brief description
	Intro string
	// Market, e.g. "HK"
	Market string
	// Circulating shares (nil if not available)
	CirculatingShares *decimal.Decimal
	// Whether this is a delayed quote
	Delay bool
	// Day change percentage (nil if not available)
	Chg *decimal.Decimal
	// Raw trade status code
	TradeStatus int32
}

// TopMoversResponse is the response for MarketContext.TopMovers.
type TopMoversResponse struct {
	Events []*TopMoversEvent
	// NextParams — pagination cursor; pass to next call to get next page
	NextParams json.RawMessage
}

// RankCategoriesResponse holds the raw data for rank categories from
// GET /v1/quote/market/rank/categories.
type RankCategoriesResponse struct {
	Data json.RawMessage `json:"data"`
}

// RankListResponse is the response for MarketContext.RankList.
type RankListResponse struct {
	Bmp   bool
	Lists []*RankListItem
}
