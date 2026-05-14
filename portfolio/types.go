package portfolio

import "github.com/shopspring/decimal"

type ExchangeRates struct {
	Exchanges []*ExchangeRate
}

type ExchangeRate struct {
	AverageRate   float64
	BaseCurrency  string
	BidRate       float64
	OfferRate     float64
	OtherCurrency string
}

type ProfitAnalysis struct {
	Summary *ProfitAnalysisSummary
	Sublist *ProfitAnalysisSublist
}

type ProfitAnalysisSummary struct {
	Currency          string
	CurrentTotalAsset *decimal.Decimal
	StartDate         string
	EndDate           string
	StartTime         string
	EndTime           string
	EndingAssetValue  *decimal.Decimal
	InitialAssetValue *decimal.Decimal
	InvestAmount      *decimal.Decimal
	IsTraded          bool
	SumProfit         *decimal.Decimal
	SumProfitRate     *decimal.Decimal
	Profits           *ProfitSummaryBreakdown
}

type ProfitSummaryBreakdown struct {
	Stock                       *decimal.Decimal
	Fund                        *decimal.Decimal
	Crypto                      *decimal.Decimal
	Mmf                         *decimal.Decimal
	Other                       *decimal.Decimal
	CumulativeTransactionAmount *decimal.Decimal
	TradeOrderNum               string
	TradeStockNum               string
	Ipo                         *decimal.Decimal
	IpoHit                      int32
	IpoSubscription             int32
	SummaryInfo                 []*ProfitSummaryInfo
}

type ProfitSummaryInfo struct {
	AssetType     string
	ProfitMax     string
	ProfitMaxName string
	LossMax       string
	LossMaxName   string
}

type ProfitAnalysisSublist struct {
	Start       string
	End         string
	StartDate   string
	EndDate     string
	UpdatedAt   string
	UpdatedDate string
	Items       []*ProfitAnalysisItem
}

type ProfitAnalysisItem struct {
	Name              string
	Market            string
	IsHolding         bool
	Profit            *decimal.Decimal
	ProfitRate        *decimal.Decimal
	ClearanceTimes    int64
	ItemType          string
	Currency          string
	Symbol            string
	HoldingPeriod     string
	SecurityCode      string
	Isin              string
	UnderlyingProfit  *decimal.Decimal
	DerivativesProfit *decimal.Decimal
	OrderProfit       *decimal.Decimal
}

type ProfitAnalysisDetail struct {
	Profit               *decimal.Decimal
	UnderlyingDetails    *ProfitDetails
	DerivativePnlDetails *ProfitDetails
	Name                 string
	UpdatedAt            string
	UpdatedDate          string
	Currency             string
	DefaultTag           int32
	Start                string
	End                  string
	StartDate            string
	EndDate              string
}

type ProfitDetails struct {
	HoldingValue                *decimal.Decimal
	Profit                      *decimal.Decimal
	CumulativeCreditedAmount    *decimal.Decimal
	CreditedDetails             []*ProfitDetailEntry
	CumulativeDebitedAmount     *decimal.Decimal
	DebitedDetails              []*ProfitDetailEntry
	CumulativeFeeAmount         *decimal.Decimal
	FeeDetails                  []*ProfitDetailEntry
	ShortHoldingValue           *decimal.Decimal
	LongHoldingValue            *decimal.Decimal
	HoldingValueAtBeginning     *decimal.Decimal
	HoldingValueAtEnding        *decimal.Decimal
}

type ProfitDetailEntry struct {
	Describe string
	Amount   *decimal.Decimal
}

type ProfitAnalysisByMarket struct {
	Profit     *decimal.Decimal
	HasMore    bool
	StockItems []*ProfitAnalysisByMarketItem
}

type ProfitAnalysisByMarketItem struct {
	Code   string
	Name   string
	Market string
	Profit *decimal.Decimal
}

type ProfitAnalysisFlows struct {
	FlowsList []*FlowItem
	HasMore   bool
}

type FlowItem struct {
	ExecutedDate     string
	ExecutedTimestamp string
	Code             string
	Direction        string
	ExecutedQuantity *decimal.Decimal
	ExecutedPrice    *decimal.Decimal
	ExecutedCost     *decimal.Decimal
	Describe         string
}
