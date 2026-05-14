package fundamental

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// FinancialReportKind identifies the type of financial statement.
type FinancialReportKind string

const (
	FinancialReportKindIncomeStatement FinancialReportKind = "IS"
	FinancialReportKindBalanceSheet    FinancialReportKind = "BS"
	FinancialReportKindCashFlow        FinancialReportKind = "CF"
	FinancialReportKindAll             FinancialReportKind = "ALL"
)

// FinancialReportPeriod identifies the reporting period.
type FinancialReportPeriod string

const (
	FinancialReportPeriodAnnual        FinancialReportPeriod = "af"
	FinancialReportPeriodSemiAnnual    FinancialReportPeriod = "saf"
	FinancialReportPeriodQ1            FinancialReportPeriod = "q1"
	FinancialReportPeriodQ2            FinancialReportPeriod = "q2"
	FinancialReportPeriodQ3            FinancialReportPeriod = "q3"
	FinancialReportPeriodQuarterlyFull FinancialReportPeriod = "qf"
	FinancialReportPeriodThreeQ        FinancialReportPeriod = "3q"
)

type FinancialReports struct {
	List json.RawMessage
}

type DividendList struct {
	List []*DividendItem
}

type DividendItem struct {
	Symbol      string
	Id          string
	Desc        string
	RecordDate  string
	ExDate      string
	PaymentDate string
}

type InstitutionRating struct {
	Latest  *InstitutionRatingLatest
	Summary *InstitutionRatingSummary
}

type InstitutionRatingLatest struct {
	Evaluate       RatingEvaluate
	Target         RatingTarget
	IndustryId     int64
	IndustryName   string
	IndustryRank   int32
	IndustryTotal  int32
	IndustryMean   int32
	IndustryMedian int32
}

type RatingEvaluate struct {
	Buy       int32
	Over      int32
	Hold      int32
	Under     int32
	Sell      int32
	NoOpinion int32
	Total     int32
	StartDate string
	EndDate   string
}

type RatingTarget struct {
	HighestPrice *decimal.Decimal
	LowestPrice  *decimal.Decimal
	PrevClose    *decimal.Decimal
	StartDate    string
	EndDate      string
}

type InstitutionRatingSummary struct {
	CcySymbol string
	Change    *decimal.Decimal
	Evaluate  RatingSummaryEvaluate
	Recommend string
	Target    *decimal.Decimal
	UpdatedAt string
}

type RatingSummaryEvaluate struct {
	Buy      int32
	Date     string
	Hold     int32
	Sell     int32
	StrongBuy int32
	Under    int32
}

type InstitutionRatingDetail struct {
	CcySymbol string
	Evaluate  InstitutionRatingDetailEvaluate
	Target    InstitutionRatingDetailTarget
}

type InstitutionRatingDetailEvaluate struct {
	List []*InstitutionRatingDetailEvaluateItem
}

type InstitutionRatingDetailEvaluateItem struct {
	Buy       int32
	Date      string
	Hold      int32
	Sell      int32
	StrongBuy int32
	NoOpinion int32
	Under     int32
}

type InstitutionRatingDetailTarget struct {
	DataPercent        *decimal.Decimal
	PredictionAccuracy *decimal.Decimal
	UpdatedAt          string
	List               []*InstitutionRatingDetailTargetItem
}

type InstitutionRatingDetailTargetItem struct {
	AvgTarget *decimal.Decimal
	Date      string
	MaxTarget *decimal.Decimal
	MinTarget *decimal.Decimal
	Meet      bool
	Price     *decimal.Decimal
	Timestamp string
}

type ForecastEps struct {
	Items []*ForecastEpsItem
}

type ForecastEpsItem struct {
	ForecastEpsMedian  *decimal.Decimal
	ForecastEpsMean    *decimal.Decimal
	ForecastEpsLowest  *decimal.Decimal
	ForecastEpsHighest *decimal.Decimal
	InstitutionTotal   int32
	InstitutionUp      int32
	InstitutionDown    int32
	ForecastStartDate  int64
	ForecastEndDate    int64
}

type FinancialConsensus struct {
	List          []*ConsensusReport
	CurrentIndex  int32
	Currency      string
	OptPeriods    []string
	CurrentPeriod string
}

type ConsensusReport struct {
	FiscalYear   int32
	FiscalPeriod string
	PeriodText   string
	Details      []*ConsensusDetail
}

type ConsensusDetail struct {
	Key        string
	Name       string
	Desc       string
	Actual     *decimal.Decimal
	Estimate   *decimal.Decimal
	CompValue  *decimal.Decimal
	CompDesc   string
	Comp       string
	IsReleased bool
}

type ValuationData struct {
	Metrics ValuationMetricsData
}

type ValuationMetricsData struct {
	Pe     *ValuationMetricData
	Pb     *ValuationMetricData
	Ps     *ValuationMetricData
	DvdYld *ValuationMetricData
}

type ValuationMetricData struct {
	Desc   string
	High   *decimal.Decimal
	Low    *decimal.Decimal
	Median *decimal.Decimal
	List   []*ValuationPoint
}

type ValuationPoint struct {
	Timestamp int64
	Value     *decimal.Decimal
}

type ValuationHistoryResponse struct {
	History ValuationHistoryData
}

type ValuationHistoryData struct {
	Metrics ValuationHistoryMetrics
}

type ValuationHistoryMetrics struct {
	Pe *ValuationHistoryMetric
	Pb *ValuationHistoryMetric
	Ps *ValuationHistoryMetric
}

type ValuationHistoryMetric struct {
	Desc   string
	High   *decimal.Decimal
	Low    *decimal.Decimal
	Median *decimal.Decimal
	List   []*ValuationPoint
}

type IndustryValuationList struct {
	List []*IndustryValuationItem
}

type IndustryValuationItem struct {
	Symbol         string
	Name           string
	Currency       string
	Assets         *decimal.Decimal
	Bps            *decimal.Decimal
	Eps            *decimal.Decimal
	Dps            *decimal.Decimal
	DivYld         *decimal.Decimal
	DivPayoutRatio *decimal.Decimal
	FiveYAvgDps    *decimal.Decimal
	Pe             *decimal.Decimal
	History        []*IndustryValuationHistory
}

type IndustryValuationHistory struct {
	Date string
	Pe   *decimal.Decimal
	Pb   *decimal.Decimal
	Ps   *decimal.Decimal
}

type IndustryValuationDist struct {
	Pe *ValuationDist
	Pb *ValuationDist
	Ps *ValuationDist
}

type ValuationDist struct {
	Low       *decimal.Decimal
	High      *decimal.Decimal
	Median    *decimal.Decimal
	Value     *decimal.Decimal
	Ranking   *decimal.Decimal
	RankIndex string
	RankTotal string
}

type CompanyOverview struct {
	Name           string
	CompanyName    string
	Founded        string
	ListingDate    string
	Market         string
	Region         string
	Address        string
	OfficeAddress  string
	Website        string
	IssuePrice     *decimal.Decimal
	SharesOffered  string
	Chairman       string
	Secretary      string
	AuditInst      string
	Category       string
	YearEnd        string
	Employees      string
	Phone          string
	Fax            string
	Email          string
	LegalRepr      string
	Manager        string
	BusLicense     string
	AccountingFirm string
	SecuritiesRep  string
	LegalCounsel   string
	ZipCode        string
	Ticker         string
	Icon           string
	Profile        string
	AdsRatio       string
	Sector         int32
}

type ExecutiveList struct {
	ProfessionalList []*ExecutiveGroup
}

type ExecutiveGroup struct {
	Symbol        string
	ForwardUrl    string
	Total         int32
	Professionals []*Professional
}

type Professional struct {
	Id        string
	Name      string
	NameZhcn  string
	NameEn    string
	Title     string
	Biography string
	Photo     string
	WikiUrl   string
}

type ShareholderList struct {
	Shareholders []*Shareholder
	ForwardUrl   string
	Total        int32
}

type Shareholder struct {
	ShareholderId   string
	ShareholderName string
	InstitutionType string
	PercentOfShares *decimal.Decimal
	SharesChanged   *decimal.Decimal
	ReportDate      string
	Stocks          []*ShareholderStock
}

type ShareholderStock struct {
	Symbol string
	Code   string
	Market string
	Chg    string
}

type FundHolders struct {
	Lists []*FundHolder
}

type FundHolder struct {
	Code          string
	Symbol        string
	Currency      string
	Name          string
	PositionRatio *decimal.Decimal
	ReportDate    string
}

type CorpActions struct {
	Items []*CorpActionItem
}

type CorpActionItem struct {
	Id           string
	Date         string
	DateStr      string
	DateType     string
	DateZone     string
	ActType      string
	ActDesc      string
	Action       string
	Recent       bool
	IsDelay      bool
	DelayContent string
	Live         *CorpActionLive
	Security     json.RawMessage
}

type CorpActionLive struct {
	Id        string
	Status    json.RawMessage
	StartedAt string
	Name      string
	Icon      string
}

type InvestRelations struct {
	ForwardUrl       string
	InvestSecurities []*InvestSecurity
}

type InvestSecurity struct {
	CompanyId       string
	CompanyName     string
	CompanyNameEn   string
	CompanyNameZhcn string
	Symbol          string
	Currency        string
	PercentOfShares *decimal.Decimal
	SharesRank      string
	SharesValue     *decimal.Decimal
}

type OperatingList struct {
	List []*OperatingItem
}

type OperatingItem struct {
	Id        string
	Report    string
	Title     string
	Txt       string
	Latest    bool
	Keywords  json.RawMessage
	WebUrl    string
	Financial OperatingFinancial
}

type OperatingFinancial struct {
	Code       string
	CounterId  string
	Currency   string
	Name       string
	Region     string
	Report     string
	ReportTxt  string
	Indicators []*OperatingIndicator
}

type OperatingIndicator struct {
	FieldName      string
	IndicatorName  string
	IndicatorValue string
	Yoy            *decimal.Decimal
}

type BuybackData struct {
	RecentBuybacks *RecentBuybacks
	BuybackHistory []*BuybackHistoryItem
	BuybackRatios  []*BuybackRatios
}

type RecentBuybacks struct {
	Currency           string
	NetBuybackTtm      *decimal.Decimal
	NetBuybackYieldTtm *decimal.Decimal
}

type BuybackHistoryItem struct {
	FiscalYear           string
	FiscalYearRange      string
	NetBuyback           *decimal.Decimal
	NetBuybackYield      *decimal.Decimal
	NetBuybackGrowthRate *decimal.Decimal
	Currency             string
}

type BuybackRatios struct {
	NetBuybackPayoutRatio     *decimal.Decimal
	NetBuybackToCashflowRatio *decimal.Decimal
}

type StockRatings struct {
	StyleTxtName        string
	ScaleTxtName        string
	ReportPeriodTxt     string
	MultiScore          json.RawMessage
	MultiLetter         string
	MultiScoreChange    int32
	IndustryName        string
	IndustryRank        json.RawMessage
	IndustryTotal       json.RawMessage
	IndustryMeanScore   json.RawMessage
	IndustryMedianScore json.RawMessage
	Ratings             []*RatingCategory
}

type RatingCategory struct {
	Kind          int32
	SubIndicators []*RatingSubIndicatorGroup
}

type RatingSubIndicatorGroup struct {
	Indicator     RatingIndicator
	SubIndicators []*RatingLeafIndicator
}

type RatingIndicator struct {
	Name   string
	Score  json.RawMessage
	Letter string
}

type RatingLeafIndicator struct {
	Name      string
	Value     string
	ValueType string
	Score     json.RawMessage
	Letter    string
}
