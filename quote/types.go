package quote

import (
	"time"

	quotev1 "github.com/longportapp/openapi-protobufs/gen/go/quote"
	"github.com/shopspring/decimal"

	"github.com/longportapp/openapi-go"
	"github.com/longportapp/openapi-go/internal/util"
	"github.com/longportapp/openapi-go/quote/jsontypes"
)

type (
	TradeStatus             int32
	TradeSession            int32
	TradeSessionType        int32
	EventType               int8
	SubType                 uint8
	Period                  int32
	AdjustType              int32
	CalcIndex               int32
	WarrantStatus           int32
	WarrantSortBy           int32
	WarrantSortOrder        int32
	WarrantType             int32
	WarrantExpiryDateType   int32
	WarrantInOutBoundsType  int32
	WarrantLanguage         int32
	SecurityListCategory    string
	WatchlistUpdateMode     string
	CandlestickTradeSession int32
)

const (
	// SubType
	SubTypeUnknown SubType = SubType(quotev1.SubType_UNKNOWN_TYPE)
	SubTypeQuote   SubType = SubType(quotev1.SubType_QUOTE)
	SubTypeDepth   SubType = SubType(quotev1.SubType_DEPTH)
	SubTypeBrokers SubType = SubType(quotev1.SubType_BROKERS)
	SubTypeTrade   SubType = SubType(quotev1.SubType_TRADE)

	// SubEvent
	EventQuote EventType = iota
	EventBroker
	EventTrade
	EventDepth

	// Period
	PeriodOneMinute     = Period(quotev1.Period_ONE_MINUTE)
	PeriodFiveMinute    = Period(quotev1.Period_FIVE_MINUTE)
	PeriodFifteenMinute = Period(quotev1.Period_FIFTEEN_MINUTE)
	PeriodThirtyMinute  = Period(quotev1.Period_THIRTY_MINUTE)
	PeriodSixtyMinute   = Period(quotev1.Period_SIXTY_MINUTE)
	PeriodDay           = Period(quotev1.Period_DAY)
	PeriodWeek          = Period(quotev1.Period_WEEK)
	PeriodMonth         = Period(quotev1.Period_MONTH)
	PeriodYear          = Period(quotev1.Period_YEAR)

	// AdjustType
	AdjustTypeNo      = AdjustType(quotev1.AdjustType_NO_ADJUST)
	AdjustTypeForward = AdjustType(quotev1.AdjustType_FORWARD_ADJUST)

	// TradeSession
	TradeSessionNormal    = TradeSession(quotev1.TradeSession_NORMAL_TRADE)
	TradeSessionPreTrade  = TradeSession(quotev1.TradeSession_PRE_TRADE)
	TradeSessionPostTrade = TradeSession(quotev1.TradeSession_POST_TRADE)
	TradeSessionOvernight = TradeSession(quotev1.TradeSession_OVERNIGHT_TRADE)

	// CandlestickTradeSessionNormal includes normal trade session
	CandlestickTradeSessionNormal = CandlestickTradeSession(0)
	// CandlestickTradeSessionAll includes all trade sessions
	CandlestickTradeSessionAll = CandlestickTradeSession(100)

	// CalcIndex
	CalcIndexUnknown               CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_UNKNOWN)
	CalcIndexLastDone              CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_LAST_DONE)
	CalcIndexChangeVal             CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_CHANGE_VAL)
	CalcIndexChangeRate            CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_CHANGE_RATE)
	CalcIndexVolume                CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_VOLUME)
	CalcIndexTurnover              CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_TURNOVER)
	CalcIndexYtdChangeRate         CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_YTD_CHANGE_RATE)
	CalcIndexTurnoverRate          CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_TURNOVER_RATE)
	CalcIndexTotalMarketValue      CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_TOTAL_MARKET_VALUE)
	CalcIndexCapitalFlow           CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_CAPITAL_FLOW)
	CalcIndexAmplitude             CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_AMPLITUDE)
	CalcIndexVolumeRatio           CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_VOLUME_RATIO)
	CalcIndexPeTTMRatio            CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_PE_TTM_RATIO)
	CalcIndexPbRatio               CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_PB_RATIO)
	CalcIndexDividendRatioTTM      CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_DIVIDEND_RATIO_TTM)
	CalcIndexFiveDayChangeRate     CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_FIVE_DAY_CHANGE_RATE)
	CalcIndexTenDayChangeRate      CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_TEN_DAY_CHANGE_RATE)
	CalcIndexHalfYearChangeRate    CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_HALF_YEAR_CHANGE_RATE)
	CalcIndexFiveMinutesChangeRate CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_FIVE_MINUTES_CHANGE_RATE)
	CalcIndexExpiryDate            CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_EXPIRY_DATE)
	CalcIndexStrikePrice           CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_STRIKE_PRICE)
	CalcIndexUpperStrikePrice      CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_UPPER_STRIKE_PRICE)
	CalcIndexLowerStrikePrice      CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_LOWER_STRIKE_PRICE)
	CalcIndexOutstandingQTY        CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_OUTSTANDING_QTY)
	CalcIndexOutstandingRatio      CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_OUTSTANDING_RATIO)
	CalcIndexPremium               CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_PREMIUM)
	CalcIndexItmOtm                CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_ITM_OTM)
	CalcIndexImpliedVolatility     CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_IMPLIED_VOLATILITY)
	CalcIndexWarrantDelta          CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_WARRANT_DELTA)
	CalcIndexCallPrice             CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_CALL_PRICE)
	CalcIndexToCallPrice           CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_TO_CALL_PRICE)
	CalcIndexEffectiveLeverage     CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_EFFECTIVE_LEVERAGE)
	CalcIndexLeverageRatio         CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_LEVERAGE_RATIO)
	CalcIndexConversionRatio       CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_CONVERSION_RATIO)
	CalcIndexBalancePoint          CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_BALANCE_POINT)
	CalcIndexOpenInterest          CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_OPEN_INTEREST)
	CalcIndexDELTA                 CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_DELTA)
	CalcIndexGAMMA                 CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_GAMMA)
	CalcIndexTHETA                 CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_THETA)
	CalcIndexVEGA                  CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_VEGA)
	CalcIndexRHO                   CalcIndex = CalcIndex(quotev1.CalcIndex_CALCINDEX_RHO)
)

// Warrant status
const (
	WarrantSuspend    WarrantStatus = iota + 2 // can't trade
	WarrantPapareList                          // wait to be listed
	WarrantNormal                              // Tradable
)

const (
	WarrantLastDone WarrantSortBy = iota
	WarrantChangeRate
	WarrantChangeVal
	WarrantVolume
	WarrantTurnover
	WarrantExpiryDate
	WarrantStrikePrice
	WarrantUpperStrikePrice
	WarrantLowerStrikePrice
	WarrantOutstandingQty
	WarrantOutstandingRatio
	WarrantPremiun
	WarrantItmOtm
	WarrantImpliedVolatility
	WarrantDelta
	WarrantCallPrice
	WarrantToCallPrice
	WarrantEffectiveLeverage
	WarrantLeverageRatio
	WarrantConversionRatio
	WarrantBalancePoint
	WarrantSortStatus
)

const (
	WarrantAsc WarrantSortOrder = iota
	WarrantDesc
)

const (
	WarrantCall WarrantType = iota
	WarrantPut
	WarrantBull
	WarrantBear
	WarrantInline
)

const (
	WarrantLT3    WarrantExpiryDateType = iota + 1 // Less than three months
	WarrantBT3_6                                   // between three and six months
	WarrantBT6_12                                  // between six and twelve months
	WarrantGT12                                    // greate than twelve months
)

const (
	WarrantInBounds WarrantInOutBoundsType = iota + 1
	WarrantOutBounds
)

const (
	WarrantZH_CN WarrantLanguage = iota
	WarrantEN
	WarrantHK_CN
)

const (
	Overnight SecurityListCategory = "Overnight"
)

const (
	// AddWatchlist will add securities to watchlist group
	AddWatchlist WatchlistUpdateMode = "add"
	// RemoveWatchlist will remove securities from watchlist group
	RemoveWatchlist WatchlistUpdateMode = "remove"
	// ReplaceWatchlist will replace securities in watchlist group
	// For example:
	//   If the original group contains APPL.US, BABA.US, TSLA.US,
	//   and it is updated with ["BABA.US", "AAPL.US", "MSFT.US"],
	//   it will become BABA.US, AAPL.US, MSFT.US,
	//   removing TSLA.US and adding MSFT.US, while adjusting the order of BABA.US and AAPL.US.
	ReplaceWatchlist WatchlistUpdateMode = "replace"
)

// PushQuote is quote info push from server
type PushQuote struct {
	Symbol       string
	Sequence     int64
	LastDone     *decimal.Decimal
	Open         *decimal.Decimal
	High         *decimal.Decimal
	Low          *decimal.Decimal
	Timestamp    int64
	Volume       int64
	Turnover     *decimal.Decimal
	TradeStatus  TradeStatus
	TradeSession TradeSessionType
}

// PushDepth is depth info push from server
type PushDepth struct {
	Symbol   string
	Sequence int64
	Ask      []*Depth
	Bid      []*Depth
}

// PushBrokers is brokers info push from server
type PushBrokers struct {
	Symbol     string
	Sequence   int64
	AskBrokers []*Brokers
	BidBrokers []*Brokers
}

// PushTrade is trade info push from server
type PushTrade struct {
	Symbol   string
	Sequence int64
	Trade    []*Trade
}

// Depth store depth details
type Depth struct {
	Position int32
	Price    *decimal.Decimal
	Volume   int64
	OrderNum int64
}

type Brokers struct {
	Position  int32
	BrokerIds []int32
}

// Trade store trade details
type Trade struct {
	Price     string
	Volume    int64
	Timestamp int64
	// TradeType
	// HK
	//
	// - `*` - Overseas trade
	// - `D` - Odd-lot trade
	// - `M` - Non-direct off-exchange trade
	// - `P` - Late trade (Off-exchange previous day)
	// - `U` - Auction trade
	// - `X` - Direct off-exchange trade
	// - `Y` - Automatch internalized
	// - `<empty string>` -  Automatch normal
	//
	// US
	//
	// - `<empty string>` - Regular sale
	// - `A` - Acquisition
	// - `B` - Bunched trade
	// - `D` - Distribution
	// - `F` - Intermarket sweep
	// - `G` - Bunched sold trades
	// - `H` - Price variation trade
	// - `I` - Odd lot trade
	// - `K` - Rule 155 trde(NYSE MKT)
	// - `M` - Market center close price
	// - `P` - Prior reference price
	// - `Q` - Market center open price
	// - `S` - Split trade
	// - `V` - Contingent trade
	// - `W` - Average price trade
	// - `X` - Cross trade
	// - `1` - Stopped stock(Regular trade)
	TradeType    string
	Direction    int32
	TradeSession TradeSession
}

// StaticInfo store static details
type StaticInfo struct {
	Symbol            string
	NameCn            string
	NameEn            string
	NameHk            string
	Exchange          string
	Currency          string
	LotSize           int32
	TotalShares       int64
	CirculatingShares int64
	HkShares          int64
	Eps               *decimal.Decimal
	EpsTtm            *decimal.Decimal
	Bps               *decimal.Decimal
	DividendYield     string
	StockDerivatives  []int32
}

// Issuer to save issuer id
type Issuer struct {
	ID     int32
	NameCn string
	NameEn string
	NameHk string
}

// OptionQuote to option quote details
type OptionQuote struct {
	Symbol       string
	LastDone     *decimal.Decimal
	PrevClose    *decimal.Decimal
	Open         *decimal.Decimal
	High         *decimal.Decimal
	Low          *decimal.Decimal
	Timestamp    int64
	Volume       int64
	Turnover     *decimal.Decimal
	TradeStatus  TradeStatus
	OptionExtend *OptionExtend
}

// OptionExtend is option extended properties
type OptionExtend struct {
	ImpliedVolatility    string
	OpenInterest         int64
	ExpiryDate           string // YYMMDD
	StrikePrice          *decimal.Decimal
	ContractMultiplier   string
	ContractType         string
	ContractSize         string
	Direction            string
	HistoricalVolatility string
	UnderlyingSymbol     string
}

// StrikePriceInfo is strike price details
type StrikePriceInfo struct {
	Price      *decimal.Decimal
	CallSymbol string
	PutSymbol  string
	Standard   bool
}

// WarrantExtend is warrant extended properties
type WarrantExtend struct {
	ImpliedVolatility string
	ExpiryDate        string
	LastTradeDate     string
	OutstandingRatio  string
	OutstandingQty    int64
	ConversionRatio   string
	Category          string
	StrikePrice       *decimal.Decimal
	UpperStrikePrice  *decimal.Decimal
	LowerStrikePrice  *decimal.Decimal
	CallPrice         *decimal.Decimal
	UnderlyingSymbol  string
}

// WarrantQuote is warrant quote details
type WarrantQuote struct {
	Symbol        string
	LastDone      *decimal.Decimal
	PrevClose     *decimal.Decimal
	Open          *decimal.Decimal
	High          *decimal.Decimal
	Low           *decimal.Decimal
	Timestamp     int64
	Volume        int64
	Turnover      *decimal.Decimal
	TradeStatus   TradeStatus
	WarrantExtend *WarrantExtend
}

type WarrantFilter struct {
	SortBy     WarrantSortBy
	SortOrder  WarrantSortOrder
	SortOffset int32
	SortCount  int32
	Type       []WarrantType
	Issuer     []int32
	ExpiryDate []WarrantExpiryDateType
	PriceType  []WarrantInOutBoundsType
	Status     []WarrantStatus
	Language   WarrantLanguage
}

// WarrantInfo is info of warrant asset
type WarrantInfo struct {
	Symbol            string
	Name              string
	LastDone          *decimal.Decimal
	ChangeRate        *decimal.Decimal
	ChangeVal         *decimal.Decimal
	Volume            int64
	Turnover          *decimal.Decimal
	ExpiryDate        string
	StrikePrice       *decimal.Decimal
	UpperStrikePrice  *decimal.Decimal
	LowerStrikePrice  *decimal.Decimal
	OutstandingQty    *decimal.Decimal
	OutstandingRatio  *decimal.Decimal
	Premium           *decimal.Decimal
	ItmOtm            *decimal.Decimal
	ImpliedVolatility *decimal.Decimal
	Delta             *decimal.Decimal
	CallPrice         *decimal.Decimal
	ToCallPrice       *decimal.Decimal
	EffectiveLeverage *decimal.Decimal
	LeverageRatio     *decimal.Decimal
	ConversionRatio   *decimal.Decimal
	BalancePoint      *decimal.Decimal
	Status            WarrantStatus
}

// TradeDate
type TradeDate struct {
	Date          string
	TradeDateType int32 // 0 full day, 1 morning only, 2 afternoon only(not happened before)
}

// Candlestick is candlestick details
type Candlestick struct {
	Close     *decimal.Decimal
	Open      *decimal.Decimal
	Low       *decimal.Decimal
	High      *decimal.Decimal
	Volume    int64
	Turnover  *decimal.Decimal
	Timestamp int64
}

// Quote is quote details
type Quote struct {
	Symbol       string
	LastDone     *decimal.Decimal
	Open         *decimal.Decimal
	High         *decimal.Decimal
	Low          *decimal.Decimal
	Timestamp    int64
	Volume       int64
	Turnover     *decimal.Decimal
	TradeStatus  TradeStatus
	TradeSession TradeSessionType
}

// SecurityQuote is quote details with pre market and post market
type SecurityQuote struct {
	Symbol          string
	LastDone        *decimal.Decimal
	PrevClose       *decimal.Decimal
	Open            *decimal.Decimal
	High            *decimal.Decimal
	Low             *decimal.Decimal
	Timestamp       int64
	Volume          int64
	Turnover        *decimal.Decimal
	TradeStatus     TradeStatus
	PreMarketQuote  *PrePostQuote
	PostMarketQuote *PrePostQuote
	OverNightQuote  *PrePostQuote
}

// PrePostQuote is pre or post quote details
type PrePostQuote struct {
	LastDone  *decimal.Decimal
	Timestamp int64
	Volume    int64
	Turnover  *decimal.Decimal
	High      *decimal.Decimal
	Low       *decimal.Decimal
	PrevClose *decimal.Decimal
}

// SecurityDepth
type SecurityDepth struct {
	Symbol string
	Ask    []*Depth
	Bid    []*Depth
}

// SecurityBrokers is security brokers details
type SecurityBrokers struct {
	Symbol     string
	AskBrokers []*Brokers
	BidBrokers []*Brokers
}

// ParticipantInfo has all participant brokers
type ParticipantInfo struct {
	BrokerIds         []int32
	ParticipantNameCn string
	ParticipantNameEn string
	ParticipantNameHk string
}

// IntradayLine is k line
type IntradayLine struct {
	Price     *decimal.Decimal
	Timestamp int64
	Volume    int64
	Turnover  *decimal.Decimal
	AvgPrice  *decimal.Decimal
}

// IssuerInfo is issuer infomation
type IssuerInfo struct {
	Id     int32
	NameCn string
	NameEn string
	NameHk string
}

// MarketTradingSession is market's session details
type MarketTradingSession struct {
	Market       openapi.Market
	TradeSession []*TradePeriod
}

// TradePeriod
type TradePeriod struct {
	BegTime      int32
	EndTime      int32
	TradeSession TradeSession
}

// MarketTradingDay contains market open trade days
type MarketTradingDay struct {
	TradeDay     []time.Time
	HalfTradeDay []time.Time
}

// WatchedSecurity the security is watched
type WatchedSecurity struct {
	Symbol    string
	Market    string
	Name      string
	Price     *decimal.Decimal
	WatchedAt int64 // timestamp
}

// WatchedGroup a group of the security is watched
type WatchedGroup struct {
	Id        int64  // group id
	Name      string // group name
	Securites []*WatchedSecurity
}

// CapitalDistribution information
type CapitalDistribution struct {
	Symbol     string
	Timestamp  int64   // data update timestamp
	CapitalIn  Capital // inflow capital data
	CapitalOut Capital // outflow capital data
}

// Capital infomartion
type Capital struct {
	Large  *decimal.Decimal
	Medium *decimal.Decimal
	Small  *decimal.Decimal
}

// CapitalFlowLine
type CapitalFlowLine struct {
	Inflow    *decimal.Decimal
	Timestamp int64
}

// SecurityCalcIndex the infomation of calculate indexes's security
type SecurityCalcIndex struct {
	Symbol                string
	LastDone              *decimal.Decimal
	ChangeVal             *decimal.Decimal
	ChangeRate            *decimal.Decimal
	Volume                int64
	Turnover              *decimal.Decimal
	YtdChangeRate         *decimal.Decimal
	TurnoverRate          *decimal.Decimal
	TotalMarketValue      *decimal.Decimal
	CapitalFlow           *decimal.Decimal
	Amplitude             *decimal.Decimal
	VolumeRatio           *decimal.Decimal
	PeTtmRatio            *decimal.Decimal
	PbRatio               *decimal.Decimal
	DividendRatioTtm      *decimal.Decimal
	FiveDayChangeRate     *decimal.Decimal
	TenDayChangeRate      *decimal.Decimal
	HalfYearChangeRate    *decimal.Decimal
	FiveMinutesChangeRate *decimal.Decimal
	ExpiryDate            string
	StrikePrice           *decimal.Decimal
	UpperStrikePrice      *decimal.Decimal
	LowerStrikePrice      *decimal.Decimal
	OutstandingQty        *decimal.Decimal
	OutstandingRatio      *decimal.Decimal
	Premium               *decimal.Decimal
	ItmOtm                *decimal.Decimal
	ImpliedVolatility     *decimal.Decimal
	WarrantDelta          *decimal.Decimal
	CallPrice             *decimal.Decimal
	ToCallPrice           *decimal.Decimal
	EffectiveLeverage     *decimal.Decimal
	LeverageRatio         *decimal.Decimal
	ConversionRatio       *decimal.Decimal
	BalancePoint          *decimal.Decimal
	OpenInterest          int64
	Delta                 *decimal.Decimal
	Gamma                 *decimal.Decimal
	Theta                 *decimal.Decimal
	Vega                  *decimal.Decimal
	Rho                   *decimal.Decimal
}

// Security is base info contains symbol and name
type Security = jsontypes.Security

// doRatio process some ratio fields
func doRatio(calcIndex *SecurityCalcIndex) {
	calcIndex.ChangeRate = util.Percent(calcIndex.ChangeRate)
	calcIndex.YtdChangeRate = util.Percent(calcIndex.YtdChangeRate)
	calcIndex.Turnover = util.Percent(calcIndex.TurnoverRate)
	calcIndex.Amplitude = util.Percent(calcIndex.Amplitude)
	calcIndex.FiveDayChangeRate = util.Percent(calcIndex.FiveDayChangeRate)
	calcIndex.TenDayChangeRate = util.Percent(calcIndex.TenDayChangeRate)
	calcIndex.HalfYearChangeRate = util.Percent(calcIndex.HalfYearChangeRate)
	calcIndex.FiveMinutesChangeRate = util.Percent(calcIndex.FiveMinutesChangeRate)
	calcIndex.OutstandingRatio = util.Percent(calcIndex.OutstandingRatio)
	calcIndex.Premium = util.Percent(calcIndex.Premium)
	calcIndex.ItmOtm = util.Percent(calcIndex.ItmOtm)
	calcIndex.ImpliedVolatility = util.Percent(calcIndex.ImpliedVolatility)
	calcIndex.ToCallPrice = util.Percent(calcIndex.ToCallPrice)
}

type UserProfile struct {
	MemberId                int64
	QuoteLevel              string
	SubscribeLimit          int32
	HistoryCandlestickLimit int32
	RateLimit               []*RateLimit
	QuoteLevelDetail        *UserQuoteLevelDetail
}

type RateLimit struct {
	Cmd   uint32
	Limit int32
	Burst int32
}

type UserQuoteLevelDetail struct {
	ByPackageKey map[string]*PackageDetail
	ByMarketCode map[string]*MarketPackageDetail
}

type PackageDetail struct {
	PackageKey string
	Limit      int32
	Burst      int32
}

type MarketPackageDetail struct {
	MarketCode string
	Limit      int32
	Burst      int32
}

// CandlestickRequestOption is the option for the candlestick request
type CandlestickRequestOption func(*quotev1.SecurityHistoryCandlestickRequest)

// CandlestickRequestTradeSession sets the trade session of the candlestick request
func CandlestickRequestTradeSession(session CandlestickTradeSession) CandlestickRequestOption {
	return func(req *quotev1.SecurityHistoryCandlestickRequest) {
		req.TradeSession = int32(session)
	}
}
