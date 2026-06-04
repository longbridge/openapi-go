package fundamental

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// ── enums ─────────────────────────────────────────────────────────

// FinancialReportKind identifies which financial statement to fetch.
type FinancialReportKind int

const (
	// FinancialReportKindIncomeStatement fetches the income statement.
	FinancialReportKindIncomeStatement FinancialReportKind = iota
	// FinancialReportKindBalanceSheet fetches the balance sheet.
	FinancialReportKindBalanceSheet
	// FinancialReportKindCashFlow fetches the cash flow statement.
	FinancialReportKindCashFlow
	// FinancialReportKindAll fetches all statements.
	FinancialReportKindAll
)

// FinancialReportPeriod identifies the reporting period.
type FinancialReportPeriod int

const (
	// FinancialReportPeriodAnnual is the annual report ("af").
	FinancialReportPeriodAnnual FinancialReportPeriod = iota
	// FinancialReportPeriodSemiAnnual is the semi-annual report ("saf").
	FinancialReportPeriodSemiAnnual
	// FinancialReportPeriodQ1 is the Q1 report.
	FinancialReportPeriodQ1
	// FinancialReportPeriodQ2 is the Q2 report.
	FinancialReportPeriodQ2
	// FinancialReportPeriodQ3 is the Q3 report.
	FinancialReportPeriodQ3
	// FinancialReportPeriodQuarterlyFull is the full quarterly report ("qf").
	FinancialReportPeriodQuarterlyFull
	// FinancialReportPeriodThreeQ is the three-quarter report ("3q").
	FinancialReportPeriodThreeQ
)

// InstitutionRecommend encodes an analyst's consensus recommendation.
type InstitutionRecommend int

const (
	// InstitutionRecommendUnknown is an unknown recommendation.
	InstitutionRecommendUnknown InstitutionRecommend = iota
	// InstitutionRecommendStrongBuy is a strong-buy recommendation.
	InstitutionRecommendStrongBuy
	// InstitutionRecommendBuy is a buy recommendation.
	InstitutionRecommendBuy
	// InstitutionRecommendHold is a hold recommendation.
	InstitutionRecommendHold
	// InstitutionRecommendSell is a sell recommendation.
	InstitutionRecommendSell
	// InstitutionRecommendStrongSell is a strong-sell recommendation.
	InstitutionRecommendStrongSell
	// InstitutionRecommendUnderperform is an underperform recommendation.
	InstitutionRecommendUnderperform
	// InstitutionRecommendNoOpinion indicates no analyst opinion.
	InstitutionRecommendNoOpinion
)

// institutionRecommendFromString converts the API string to the Go enum.
func institutionRecommendFromString(s string) InstitutionRecommend {
	switch s {
	case "strong_buy":
		return InstitutionRecommendStrongBuy
	case "buy":
		return InstitutionRecommendBuy
	case "hold":
		return InstitutionRecommendHold
	case "sell":
		return InstitutionRecommendSell
	case "strong_sell":
		return InstitutionRecommendStrongSell
	case "underperform":
		return InstitutionRecommendUnderperform
	case "no_opinion":
		return InstitutionRecommendNoOpinion
	default:
		return InstitutionRecommendUnknown
	}
}

// ── financial_report ─────────────────────────────────────────────

// FinancialReports is the response for FundamentalContext.FinancialReport.
// The List field contains deeply-nested indicator/account/value data keyed
// by report kind ("IS", "BS", "CF"). The exact structure varies and is
// preserved as raw JSON.
type FinancialReports struct {
	// Raw nested financial data keyed by report kind.
	List json.RawMessage
}

// ── dividend ─────────────────────────────────────────────────────

// DividendList is the response for FundamentalContext.Dividend and
// FundamentalContext.DividendDetail.
type DividendList struct {
	// List of dividend events.
	List []DividendItem
}

// DividendItem is a single dividend / distribution event.
type DividendItem struct {
	// Security symbol, e.g. "700.HK".
	Symbol string
	// Internal record ID.
	ID string
	// Human-readable description, e.g. "每股派息 5.3 HKD".
	Desc string
	// Record / book-close date, e.g. "2026.05.18".
	RecordDate string
	// Ex-dividend date, e.g. "2026.05.15".
	ExDate string
	// Payment date, e.g. "2026.06.01".
	PaymentDate string
}

// ── institution_rating ────────────────────────────────────────────

// InstitutionRating is the combined analyst-rating response for
// FundamentalContext.InstitutionRating.
type InstitutionRating struct {
	// Latest snapshot.
	Latest InstitutionRatingLatest
	// Consensus summary.
	Summary InstitutionRatingSummary
}

// InstitutionRatingLatest is the latest analyst-rating snapshot.
type InstitutionRatingLatest struct {
	// Rating distribution counts and date range.
	Evaluate RatingEvaluate
	// Target price range.
	Target RatingTarget
	// Industry classification ID.
	IndustryID int64
	// Industry name.
	IndustryName string
	// Rank of this security within the industry (1 = highest).
	IndustryRank int32
	// Total number of securities in the industry.
	IndustryTotal int32
	// Mean analyst count in the industry.
	IndustryMean int32
	// Median analyst count in the industry.
	IndustryMedian int32
}

// RatingEvaluate holds analyst rating distribution counts.
type RatingEvaluate struct {
	// Number of "Buy" ratings.
	Buy int32
	// Number of "Strong Buy" / "Outperform" ratings.
	Over int32
	// Number of "Hold" / "Neutral" ratings.
	Hold int32
	// Number of "Underperform" ratings.
	Under int32
	// Number of "Sell" ratings.
	Sell int32
	// Number of "No Opinion" ratings.
	NoOpinion int32
	// Total analyst count.
	Total int32
	// Window start (unix timestamp string; "0" means unset).
	StartDate string
	// Window end (unix timestamp string; "0" means unset).
	EndDate string
}

// RatingTarget holds the analyst target price range.
type RatingTarget struct {
	// Highest price target.
	HighestPrice *decimal.Decimal
	// Lowest price target.
	LowestPrice *decimal.Decimal
	// Previous close price.
	PrevClose *decimal.Decimal
	// Window start (unix timestamp string).
	StartDate string
	// Window end (unix timestamp string).
	EndDate string
}

// InstitutionRatingSummary is the consensus summary.
type InstitutionRatingSummary struct {
	// Currency symbol, e.g. "HK$".
	CcySymbol string
	// Change vs previous period.
	Change *decimal.Decimal
	// Simplified rating distribution.
	Evaluate RatingSummaryEvaluate
	// Overall recommendation.
	Recommend InstitutionRecommend
	// Consensus target price.
	Target *decimal.Decimal
	// Last updated display string.
	UpdatedAt string
}

// RatingSummaryEvaluate is the simplified rating distribution for the
// consensus summary.
type RatingSummaryEvaluate struct {
	// Number of "Buy" ratings.
	Buy int32
	// Date of the latest update.
	Date string
	// Number of "Hold" ratings.
	Hold int32
	// Number of "Sell" ratings.
	Sell int32
	// Number of "Strong Buy" ratings.
	StrongBuy int32
	// Number of "Underperform" ratings.
	Under int32
}

// ── institution_rating_detail ─────────────────────────────────────

// InstitutionRatingDetail is the response for
// FundamentalContext.InstitutionRatingDetail.
type InstitutionRatingDetail struct {
	// Currency symbol, e.g. "HK$".
	CcySymbol string
	// Historical rating distribution time-series.
	Evaluate InstitutionRatingDetailEvaluate
	// Historical target price time-series.
	Target InstitutionRatingDetailTarget
}

// InstitutionRatingDetailEvaluate holds the historical rating distribution
// time-series.
type InstitutionRatingDetailEvaluate struct {
	// Weekly snapshots ordered from oldest to newest.
	List []InstitutionRatingDetailEvaluateItem
}

// InstitutionRatingDetailEvaluateItem is one weekly rating distribution
// snapshot.
type InstitutionRatingDetailEvaluateItem struct {
	// Number of "Buy" ratings.
	Buy int32
	// Date in "2021/05/14" format.
	Date string
	// Number of "Hold" ratings.
	Hold int32
	// Number of "Sell" ratings.
	Sell int32
	// Number of "Strong Buy" / "Outperform" ratings.
	StrongBuy int32
	// Number of "No Opinion" ratings.
	NoOpinion int32
	// Number of "Underperform" ratings.
	Under int32
}

// InstitutionRatingDetailTarget holds the historical target price time-series.
type InstitutionRatingDetailTarget struct {
	// Prediction accuracy ratio (may be nil).
	DataPercent *decimal.Decimal
	// Overall prediction accuracy percentage.
	PredictionAccuracy *decimal.Decimal
	// Last updated display string.
	UpdatedAt string
	// Weekly target price snapshots.
	List []InstitutionRatingDetailTargetItem
}

// InstitutionRatingDetailTargetItem is one weekly target price snapshot.
type InstitutionRatingDetailTargetItem struct {
	// Average target price.
	AvgTarget *decimal.Decimal
	// Date in "2021/05/16" format.
	Date string
	// Highest target price.
	MaxTarget *decimal.Decimal
	// Lowest target price.
	MinTarget *decimal.Decimal
	// Whether the stock price reached the target.
	Meet bool
	// Actual stock price at this date.
	Price *decimal.Decimal
	// Unix timestamp string.
	Timestamp string
}

// ── forecast_eps ──────────────────────────────────────────────────

// ForecastEps is the response for FundamentalContext.ForecastEps.
type ForecastEps struct {
	// EPS forecast snapshots ordered by ForecastStartDate ascending.
	Items []ForecastEpsItem
}

// ForecastEpsItem is one EPS forecast snapshot.
type ForecastEpsItem struct {
	// Median EPS estimate.
	ForecastEpsMedian *decimal.Decimal
	// Mean EPS estimate.
	ForecastEpsMean *decimal.Decimal
	// Lowest EPS estimate.
	ForecastEpsLowest *decimal.Decimal
	// Highest EPS estimate.
	ForecastEpsHighest *decimal.Decimal
	// Total number of forecasting institutions.
	InstitutionTotal int32
	// Number of institutions that raised their estimate.
	InstitutionUp int32
	// Number of institutions that lowered their estimate.
	InstitutionDown int32
	// Forecast window start.
	ForecastStartDate time.Time
	// Forecast window end.
	ForecastEndDate time.Time
}

// ── consensus ─────────────────────────────────────────────────────

// FinancialConsensus is the response for FundamentalContext.Consensus.
type FinancialConsensus struct {
	// Per-period consensus reports.
	List []ConsensusReport
	// Index into List of the most recently released period.
	CurrentIndex int32
	// Reporting currency, e.g. "HKD".
	Currency string
	// Available period types, e.g. ["qf", "saf", "af"].
	OptPeriods []string
	// Currently returned period type.
	CurrentPeriod string
}

// ConsensusReport is the consensus report for one fiscal period.
type ConsensusReport struct {
	// Fiscal year, e.g. 2025.
	FiscalYear int32
	// Fiscal period code, e.g. "Q4".
	FiscalPeriod string
	// Human-readable period label, e.g. "Q4 FY2025".
	PeriodText string
	// Per-metric consensus details.
	Details []ConsensusDetail
}

// ConsensusDetail is the consensus estimate for one financial metric.
type ConsensusDetail struct {
	// Metric key, e.g. "revenue", "eps".
	Key string
	// Display name.
	Name string
	// Metric description.
	Description string
	// Actual reported value (nil if not yet released).
	Actual *decimal.Decimal
	// Consensus estimate value.
	Estimate *decimal.Decimal
	// Actual minus estimate.
	CompValue *decimal.Decimal
	// Beat/miss description, e.g. "超出预期".
	CompDesc string
	// Comparison result code for colour coding.
	Comp string
	// Whether the actual results have been published.
	IsReleased bool
}

// ── valuation ─────────────────────────────────────────────────────

// ValuationData is the response for FundamentalContext.Valuation.
type ValuationData struct {
	// Valuation metrics (PE / PB / PS / dividend yield).
	Metrics ValuationMetricsData
}

// ValuationMetricsData holds all valuation metrics containers.
type ValuationMetricsData struct {
	// Price-to-Earnings ratio history.
	PE *ValuationMetricData
	// Price-to-Book ratio history.
	PB *ValuationMetricData
	// Price-to-Sales ratio history.
	PS *ValuationMetricData
	// Dividend yield history.
	DvdYld *ValuationMetricData
}

// ValuationMetricData holds the historical time-series for one valuation
// metric.
type ValuationMetricData struct {
	// Human-readable description with current value and percentile.
	Desc string
	// Historical high value.
	High *decimal.Decimal
	// Historical low value.
	Low *decimal.Decimal
	// Historical median value.
	Median *decimal.Decimal
	// Historical data points.
	List []ValuationPoint
}

// ValuationPoint is one valuation data point.
type ValuationPoint struct {
	// Date of the data point.
	Timestamp time.Time
	// Metric value.
	Value *decimal.Decimal
}

// ── valuation_history ─────────────────────────────────────────────

// ValuationHistoryResponse is the response for
// FundamentalContext.ValuationHistory.
type ValuationHistoryResponse struct {
	// Historical valuation data.
	History ValuationHistoryData
}

// ValuationHistoryData holds the historical valuation metrics container.
type ValuationHistoryData struct {
	// Historical metrics (PE / PB / PS).
	Metrics ValuationHistoryMetrics
}

// ValuationHistoryMetrics holds PE/PB/PS historical data.
type ValuationHistoryMetrics struct {
	// Price-to-Earnings history.
	PE *ValuationHistoryMetric
	// Price-to-Book history.
	PB *ValuationHistoryMetric
	// Price-to-Sales history.
	PS *ValuationHistoryMetric
}

// ValuationHistoryMetric holds the historical data for one valuation metric
// including statistical bounds.
type ValuationHistoryMetric struct {
	// Human-readable description.
	Desc string
	// Historical high over the period.
	High *decimal.Decimal
	// Historical low over the period.
	Low *decimal.Decimal
	// Historical median over the period.
	Median *decimal.Decimal
	// Historical data points.
	List []ValuationPoint
}

// ── industry_valuation ────────────────────────────────────────────

// IndustryValuationList is the response for
// FundamentalContext.IndustryValuation.
type IndustryValuationList struct {
	// List of peer securities with their valuation data.
	List []IndustryValuationItem
}

// IndustryValuationItem holds valuation data for one peer security.
type IndustryValuationItem struct {
	// Security symbol, e.g. "700.HK".
	Symbol string
	// Company name.
	Name string
	// Reporting currency.
	Currency string
	// Total assets.
	Assets *decimal.Decimal
	// Book value per share.
	Bps *decimal.Decimal
	// Earnings per share.
	Eps *decimal.Decimal
	// Dividends per share.
	Dps *decimal.Decimal
	// Dividend yield.
	DivYld *decimal.Decimal
	// Dividend payout ratio.
	DivPayoutRatio *decimal.Decimal
	// 5-year average dividends per share.
	FiveYAvgDps *decimal.Decimal
	// Current PE ratio.
	PE *decimal.Decimal
	// Historical PE/PB/PS snapshots.
	History []IndustryValuationHistory
}

// IndustryValuationHistory is a historical valuation snapshot for an industry
// peer.
type IndustryValuationHistory struct {
	// Unix timestamp string.
	Date string
	// Price-to-Earnings ratio.
	PE *decimal.Decimal
	// Price-to-Book ratio.
	PB *decimal.Decimal
	// Price-to-Sales ratio.
	PS *decimal.Decimal
}

// ── industry_valuation_dist ───────────────────────────────────────

// IndustryValuationDist is the response for
// FundamentalContext.IndustryValuationDist.
type IndustryValuationDist struct {
	// PE ratio distribution within the industry.
	PE *ValuationDist
	// PB ratio distribution within the industry.
	PB *ValuationDist
	// PS ratio distribution within the industry.
	PS *ValuationDist
}

// ValuationDist holds distribution statistics for one valuation metric within
// an industry.
type ValuationDist struct {
	// Minimum value in the industry.
	Low *decimal.Decimal
	// Maximum value in the industry.
	High *decimal.Decimal
	// Median value in the industry.
	Median *decimal.Decimal
	// Current value of the queried security.
	Value *decimal.Decimal
	// Percentile ranking (0–1 range).
	Ranking *decimal.Decimal
	// Ordinal rank index (1-based, as string).
	RankIndex string
	// Total number of securities in the industry (as string).
	RankTotal string
}

// ── company ───────────────────────────────────────────────────────

// CompanyOverview is the response for FundamentalContext.Company.
type CompanyOverview struct {
	// Short name, e.g. "腾讯控股".
	Name string
	// Full legal name.
	CompanyName string
	// Founding date.
	Founded string
	// Listing date.
	ListingDate string
	// Primary listing market display name.
	Market string
	// Market region code, e.g. "HK".
	Region string
	// Registered address.
	Address string
	// Principal office address.
	OfficeAddress string
	// Company website.
	Website string
	// IPO issue price.
	IssuePrice *decimal.Decimal
	// Number of shares offered at IPO.
	SharesOffered string
	// Chairman name.
	Chairman string
	// Company secretary name.
	Secretary string
	// Auditing institution.
	AuditInst string
	// Company classification category.
	Category string
	// Fiscal year end, e.g. "12 月 31 日".
	YearEnd string
	// Number of employees.
	Employees string
	// Phone number.
	Phone string
	// Fax number.
	Fax string
	// Investor relations email.
	Email string
	// Legal representative.
	LegalRepr string
	// CEO / Managing Director.
	Manager string
	// Business licence number.
	BusLicense string
	// Accounting firm.
	AccountingFirm string
	// Securities representative.
	SecuritiesRep string
	// Legal counsel.
	LegalCounsel string
	// Postal code.
	ZipCode string
	// Exchange ticker code, e.g. "00700".
	Ticker string
	// URL to the company's logo icon.
	Icon string
	// Business profile / description.
	Profile string
	// ADS ratio (may be empty).
	AdsRatio string
	// Industry sector code.
	Sector int32
}

// ── executive ─────────────────────────────────────────────────────

// ExecutiveList is the response for FundamentalContext.Executive.
type ExecutiveList struct {
	// Groups of executives per security (usually one group).
	ProfessionalList []ExecutiveGroup
}

// ExecutiveGroup holds executives for one security.
type ExecutiveGroup struct {
	// Security symbol.
	Symbol string
	// Link to the company wiki page.
	ForwardURL string
	// Total number of executives.
	Total int32
	// Individual executive entries.
	Professionals []Professional
}

// Professional is one executive / board member.
type Professional struct {
	// Internal wiki person ID.
	ID string
	// Full name.
	Name string
	// Full name in Simplified Chinese.
	NameZhCN string
	// Full name in English.
	NameEn string
	// Job title, e.g. "Co-Founder, Chairman & CEO".
	Title string
	// Biography text.
	Biography string
	// URL to the person's photo.
	Photo string
	// URL to the wiki profile page.
	WikiURL string
}

// ── shareholder ───────────────────────────────────────────────────

// ShareholderList is the response for FundamentalContext.Shareholder.
type ShareholderList struct {
	// List of major shareholders.
	ShareholderList []Shareholder
	// Link to the full shareholder page.
	ForwardURL string
	// Total number of shareholders returned.
	Total int32
}

// Shareholder is one major shareholder.
type Shareholder struct {
	// Internal shareholder ID.
	ShareholderID string
	// Shareholder name.
	ShareholderName string
	// Institution type (may be empty).
	InstitutionType string
	// Percentage of shares held.
	PercentOfShares *decimal.Decimal
	// Change in shares held (positive = bought, negative = sold).
	SharesChanged *decimal.Decimal
	// Date of the most recent filing, e.g. "2026-05-04".
	ReportDate string
	// Other securities held by this shareholder (cross-holdings).
	Stocks []ShareholderStock
}

// ShareholderStock is a security in an institutional shareholder's
// cross-holdings.
type ShareholderStock struct {
	// Security symbol of the cross-held stock.
	Symbol string
	// Ticker code, e.g. "BLK".
	Code string
	// Market, e.g. "US".
	Market string
	// Day change percentage, e.g. "-0.32%".
	Chg string
}

// ── fund_holder ───────────────────────────────────────────────────

// FundHolders is the response for FundamentalContext.FundHolder.
type FundHolders struct {
	// Funds and ETFs that hold the queried security.
	Lists []FundHolder
}

// FundHolder is a fund or ETF that holds the queried security.
type FundHolder struct {
	// Fund/ETF ticker code, e.g. "513050".
	Code string
	// Fund/ETF symbol.
	Symbol string
	// Reporting currency, e.g. "CNY".
	Currency string
	// Fund/ETF full name.
	Name string
	// Position ratio as a percentage decimal.
	PositionRatio decimal.Decimal
	// Report date, e.g. "2025.12.31".
	ReportDate string
}

// ── corp_action ───────────────────────────────────────────────────

// CorpActions is the response for FundamentalContext.CorpAction.
type CorpActions struct {
	// Corporate action events.
	Items []CorpActionItem
}

// CorpActionItem is one corporate action event.
type CorpActionItem struct {
	// Internal event ID.
	ID string
	// Date in YYYYMMDD format, e.g. "20260601".
	Date string
	// Short display date, e.g. "06.01".
	DateStr string
	// Date type label, e.g. "派息日", "除权日".
	DateType string
	// Time zone description, e.g. "北京时间".
	DateZone string
	// Event category, e.g. "分配方案".
	ActType string
	// Human-readable event description.
	ActDesc string
	// Machine-readable action code, e.g. "DividendExDate".
	Action string
	// Whether this is a recent event.
	Recent bool
	// Whether publication was delayed.
	IsDelay bool
	// Delay announcement content (if IsDelay is true).
	DelayContent string
	// Associated live stream (if any).
	Live *CorpActionLive
	// Associated security info (rarely populated; preserved as raw JSON).
	Security *json.RawMessage
}

// CorpActionLive is the live stream associated with a corporate action.
type CorpActionLive struct {
	// Live stream ID.
	ID string
	// Status (raw JSON; API may return int or string).
	// 1=preview, 2=live, 3=ended, 4=replay, 5=processing.
	Status json.RawMessage
	// Start time.
	StartedAt string
	// Stream title.
	Name string
	// Icon URL.
	Icon string
}

// ── invest_relation ───────────────────────────────────────────────

// InvestRelations is the response for FundamentalContext.InvestRelation.
type InvestRelations struct {
	// Link to the full investor-relations page.
	ForwardURL string
	// Securities in which the queried company holds a stake.
	InvestSecurities []InvestSecurity
}

// InvestSecurity is a security in which the queried company has an investment
// stake.
type InvestSecurity struct {
	// Internal company ID (string form; may be "0").
	CompanyID string
	// Company name (locale-aware).
	CompanyName string
	// Company name in English.
	CompanyNameEn string
	// Company name in Simplified Chinese.
	CompanyNameZhCN string
	// Security symbol of the invested company.
	Symbol string
	// Reporting currency.
	Currency string
	// Percentage of shares held.
	PercentOfShares *decimal.Decimal
	// Shareholder rank, e.g. "1" = largest shareholder.
	SharesRank string
	// Market value of the holding.
	SharesValue *decimal.Decimal
}

// ── operating ─────────────────────────────────────────────────────

// OperatingList is the response for FundamentalContext.Operating.
type OperatingList struct {
	// List of operating summary reports.
	List []OperatingItem
}

// OperatingItem is one operating summary report (annual / quarterly).
type OperatingItem struct {
	// Internal report ID.
	ID string
	// Report period code, e.g. "af" (annual), "qf" (quarterly).
	Report string
	// Report title, e.g. "2025 财年年报".
	Title string
	// Management discussion text.
	Txt string
	// Whether this is the most recent report.
	Latest bool
	// Keyword tags (structure undocumented; usually empty).
	Keywords []json.RawMessage
	// URL to the full community report page.
	WebURL string
	// Key financial metrics extracted from the report.
	Financial OperatingFinancial
}

// OperatingFinancial holds key financial metrics extracted from an operating
// report.
type OperatingFinancial struct {
	// Ticker code (may be empty).
	Code string
	// Symbol in CODE.MARKET format (may be empty).
	Symbol string
	// Reporting currency.
	Currency string
	// Company name.
	Name string
	// Market region.
	Region string
	// Report period code.
	Report string
	// Report period display text.
	ReportTxt string
	// Financial indicators.
	Indicators []OperatingIndicator
}

// OperatingIndicator is one financial indicator in an operating report.
type OperatingIndicator struct {
	// Field name key, e.g. "operating_revenue".
	FieldName string
	// Display name, e.g. "营业收入".
	IndicatorName string
	// Formatted value, e.g. "8217 亿".
	IndicatorValue string
	// Year-over-year change.
	Yoy *decimal.Decimal
}

// ── buyback ───────────────────────────────────────────────────────

// BuybackData is the response for FundamentalContext.Buyback.
type BuybackData struct {
	// Most recent buyback summary (TTM).
	RecentBuybacks *RecentBuybacks
	// Historical annual buyback data.
	BuybackHistory []BuybackHistoryItem
	// Buyback payout and cash-flow ratios.
	BuybackRatios []BuybackRatios
}

// RecentBuybacks is the TTM (trailing twelve months) buyback summary.
type RecentBuybacks struct {
	// Reporting currency.
	Currency string
	// Net buyback amount TTM.
	NetBuybackTTM *decimal.Decimal
	// Net buyback yield TTM.
	NetBuybackYieldTTM *decimal.Decimal
}

// BuybackHistoryItem is one historical annual buyback data point.
type BuybackHistoryItem struct {
	// Fiscal year label, e.g. "FY2024".
	FiscalYear string
	// Fiscal year date range string.
	FiscalYearRange string
	// Net buyback amount.
	NetBuyback *decimal.Decimal
	// Net buyback yield.
	NetBuybackYield *decimal.Decimal
	// Year-over-year net buyback growth rate.
	NetBuybackGrowthRate *decimal.Decimal
	// Reporting currency.
	Currency string
}

// BuybackRatios holds buyback payout and cash-flow ratios.
type BuybackRatios struct {
	// Net buyback payout ratio.
	NetBuybackPayoutRatio *decimal.Decimal
	// Net buyback to free cash-flow ratio.
	NetBuybackToCashflowRatio *decimal.Decimal
}

// ── ratings ───────────────────────────────────────────────────────

// StockRatings is the response for FundamentalContext.Ratings.
type StockRatings struct {
	// Style display name.
	StyleTxtName string
	// Scale display name.
	ScaleTxtName string
	// Report period display text.
	ReportPeriodTxt string
	// Composite score (raw JSON; may be int, float, or null).
	MultiScore json.RawMessage
	// Composite score letter grade.
	MultiLetter string
	// Score change vs previous period.
	MultiScoreChange int32
	// Industry name.
	IndustryName string
	// Industry rank (raw JSON; may be int or null).
	IndustryRank json.RawMessage
	// Total securities in the industry (raw JSON; may be int or null).
	IndustryTotal json.RawMessage
	// Industry mean score (raw JSON; may be float or null).
	IndustryMeanScore json.RawMessage
	// Industry median score (raw JSON; may be float or null).
	IndustryMedianScore json.RawMessage
	// Detailed rating categories.
	Ratings []RatingCategory
}

// RatingCategory is one rating category (e.g. growth, profitability).
type RatingCategory struct {
	// Category type code.
	Kind int32
	// Sub-indicator groups within this category.
	SubIndicators []RatingSubIndicatorGroup
}

// RatingSubIndicatorGroup is a group of sub-indicators under one category
// indicator.
type RatingSubIndicatorGroup struct {
	// Parent indicator for this group.
	Indicator RatingIndicator
	// Leaf sub-indicators.
	SubIndicators []RatingLeafIndicator
}

// RatingIndicator is a rating indicator node (may be a parent or a leaf).
type RatingIndicator struct {
	// Indicator display name.
	Name string
	// Score (raw JSON; may be int, float, or null).
	Score json.RawMessage
	// Letter grade.
	Letter string
}

// RatingLeafIndicator is a leaf rating indicator with a raw value.
type RatingLeafIndicator struct {
	// Indicator display name.
	Name string
	// Formatted value string.
	Value string
	// Value type hint, e.g. "percent".
	ValueType string
	// Score (raw JSON; may be int, float, or null).
	Score json.RawMessage
	// Letter grade.
	Letter string
}

// ── industry_rank enums ───────────────────────────────────────────

// IndustryRankIndicator identifies the metric used for industry ranking.
type IndustryRankIndicator string

const (
	// IndustryRankIndicator0 is indicator 0.
	IndustryRankIndicator0 IndustryRankIndicator = "0"
	// IndustryRankIndicator1 is indicator 1.
	IndustryRankIndicator1 IndustryRankIndicator = "1"
	// IndustryRankIndicator2 is indicator 2.
	IndustryRankIndicator2 IndustryRankIndicator = "2"
	// IndustryRankIndicator3 is indicator 3.
	IndustryRankIndicator3 IndustryRankIndicator = "3"
	// IndustryRankIndicator4 is indicator 4.
	IndustryRankIndicator4 IndustryRankIndicator = "4"
	// IndustryRankIndicator5 is indicator 5.
	IndustryRankIndicator5 IndustryRankIndicator = "5"
	// IndustryRankIndicator6 is indicator 6.
	IndustryRankIndicator6 IndustryRankIndicator = "6"
	// IndustryRankIndicator7 is indicator 7.
	IndustryRankIndicator7 IndustryRankIndicator = "7"
)

// IndustryRankSortType specifies the sort direction for industry ranking.
type IndustryRankSortType string

const (
	// IndustryRankSortTypeAscending sorts ascending.
	IndustryRankSortTypeAscending IndustryRankSortType = "0"
	// IndustryRankSortTypeDescending sorts descending.
	IndustryRankSortTypeDescending IndustryRankSortType = "1"
)

// ── business_segments ─────────────────────────────────────────────

// BusinessSegments is the response for FundamentalContext.BusinessSegments.
type BusinessSegments struct {
	// Report date.
	Date string
	// Total revenue.
	Total string
	// Reporting currency.
	Currency string
	// Business segment breakdown.
	Business []BusinessSegmentItem
}

// BusinessSegmentItem is one business segment entry (latest snapshot).
type BusinessSegmentItem struct {
	// Segment name.
	Name string
	// Percentage of total revenue.
	Percent string
}

// BusinessSegmentsHistory is the response for
// FundamentalContext.BusinessSegmentsHistory.
type BusinessSegmentsHistory struct {
	// Historical snapshots.
	Historical []BusinessSegmentsHistoricalItem
}

// BusinessSegmentsHistoricalItem is one historical business segments snapshot.
type BusinessSegmentsHistoricalItem struct {
	// Report date.
	Date string
	// Total revenue.
	Total string
	// Reporting currency.
	Currency string
	// Business segment breakdown.
	Business []BusinessSegmentHistoryItem
	// Regional breakdown.
	Regionals []BusinessSegmentHistoryItem
}

// BusinessSegmentHistoryItem is one business/regional segment entry in a
// historical snapshot.
type BusinessSegmentHistoryItem struct {
	// Segment name.
	Name string
	// Percentage of total.
	Percent string
	// Absolute value.
	Value string
}

// ── institution_rating_views ──────────────────────────────────────

// InstitutionRatingViews is the response for
// FundamentalContext.InstitutionRatingViews.
type InstitutionRatingViews struct {
	// Historical rating distribution snapshots.
	Elist []InstitutionRatingViewItem
}

// InstitutionRatingViewItem is one historical rating distribution snapshot.
type InstitutionRatingViewItem struct {
	// Date of the snapshot.
	Date time.Time
	// Number of "Buy" ratings.
	Buy string
	// Number of "Outperform" ratings.
	Over string
	// Number of "Hold" ratings.
	Hold string
	// Number of "Underperform" ratings.
	Under string
	// Number of "Sell" ratings.
	Sell string
	// Total analyst count.
	Total string
}

// ── industry_rank ─────────────────────────────────────────────────

// IndustryRankResponse is the response for FundamentalContext.IndustryRank.
type IndustryRankResponse struct {
	// Grouped rank items.
	Items []IndustryRankGroup
}

// IndustryRankGroup is a group of ranked industry items.
type IndustryRankGroup struct {
	// Items in this group.
	Lists []IndustryRankItem
}

// IndustryRankItem is one ranked industry item.
type IndustryRankItem struct {
	// Industry / sector name.
	Name string
	// Counter ID of the industry.
	CounterID string
	// Change percentage.
	Chg string
	// Name of the leading stock.
	LeadingName string
	// Ticker of the leading stock.
	LeadingTicker string
	// Change percentage of the leading stock.
	LeadingChg string
	// Value label name.
	ValueName string
	// Value data.
	ValueData string
}

// ── industry_peers ────────────────────────────────────────────────

// IndustryPeersResponse is the response for FundamentalContext.IndustryPeers.
type IndustryPeersResponse struct {
	// Top-level industry node info.
	Top IndustryPeersTop
	// Root peer chain node (nil if no data).
	Chain *IndustryPeerNode
}

// IndustryPeersTop holds the top-level industry info.
type IndustryPeersTop struct {
	// Industry name.
	Name string
	// Market code.
	Market string
}

// IndustryPeerNode is a node in the recursive industry peer chain.
type IndustryPeerNode struct {
	// Node name.
	Name string
	// Counter ID.
	CounterID string
	// Number of stocks in this node.
	StockNum int32
	// Change percentage.
	Chg string
	// Year-to-date change.
	YtdChg string
	// Child nodes (recursive).
	Next []IndustryPeerNode
}

// ── financial_report_snapshot ─────────────────────────────────────

// FinancialReportSnapshot is the response for
// FundamentalContext.FinancialReportSnapshot.
type FinancialReportSnapshot struct {
	// Company name.
	Name string
	// Ticker code.
	Ticker string
	// Fiscal period start date.
	FpStart string
	// Fiscal period end date.
	FpEnd string
	// Reporting currency.
	Currency string
	// Report description.
	ReportDesc string
	// Forecast revenue.
	FoRevenue *SnapshotForecastMetric
	// Forecast EBIT.
	FoEbit *SnapshotForecastMetric
	// Forecast EPS.
	FoEps *SnapshotForecastMetric
	// Reported revenue.
	FrRevenue *SnapshotReportedMetric
	// Reported net profit.
	FrProfit *SnapshotReportedMetric
	// Reported operating cash flow.
	FrOperateCash *SnapshotReportedMetric
	// Reported investing cash flow.
	FrInvestCash *SnapshotReportedMetric
	// Reported financing cash flow.
	FrFinanceCash *SnapshotReportedMetric
	// Reported total assets.
	FrTotalAssets *SnapshotReportedMetric
	// Reported total liabilities.
	FrTotalLiability *SnapshotReportedMetric
	// ROE TTM.
	FrRoeTtm string
	// Profit margin.
	FrProfitMargin string
	// Profit margin TTM.
	FrProfitMarginTtm string
	// Asset turnover TTM.
	FrAssetTurnTtm string
	// Leverage TTM.
	FrLeverageTtm string
	// Debt-to-assets ratio.
	FrDebtAssetsRatio string
}

// SnapshotForecastMetric is a forecast metric in the financial report snapshot.
type SnapshotForecastMetric struct {
	// Actual value.
	Value string
	// Year-over-year change.
	Yoy string
	// Beat/miss description.
	CmpDesc string
	// Consensus estimate value.
	EstValue string
}

// SnapshotReportedMetric is a reported metric in the financial report snapshot.
type SnapshotReportedMetric struct {
	// Actual value.
	Value string
	// Year-over-year change.
	Yoy string
}

// ── ShareholderTopResponse ────────────────────────────────────────

// ShareholderTopResponse holds the raw data for the top shareholders list.
// The Data field contains the JSON payload from GET /v1/quote/shareholders/top.
type ShareholderTopResponse struct {
	Data json.RawMessage
}

// ── ShareholderDetailResponse ─────────────────────────────────────

// ShareholderDetailResponse holds the raw data for a single shareholder's
// holding details from GET /v1/quote/shareholders/holding.
type ShareholderDetailResponse struct {
	Data json.RawMessage
}

// ── ValuationComparisonResponse ───────────────────────────────────

// ValuationHistoryPoint is one historical valuation data point.
type ValuationHistoryPoint struct {
	// Date — RFC 3339 (converted from Unix timestamp)
	Date string
	Pe   string
	Pb   string
	Ps   string
}

// ValuationComparisonItem is one security in the valuation comparison.
type ValuationComparisonItem struct {
	// Symbol — converted from counter_id (e.g. "AAPL.US")
	Symbol      string
	Name        string
	Currency    string
	MarketValue string
	PriceClose  string
	Pe          string
	Pb          string
	Ps          string
	Roe         string
	Eps         string
	Bps         string
	Dps         string
	DivYld      string
	Assets      string
	History     []*ValuationHistoryPoint
}

// ValuationComparisonResponse is the response for FundamentalContext.ValuationComparison.
type ValuationComparisonResponse struct {
	List []*ValuationComparisonItem
}

// ── EtfAssetAllocation ────────────────────────────────────────────

// ElementType identifies the kind of an ETF asset allocation group.
type ElementType int32

const (
	// ElementTypeUnknown is an unknown / unrecognized element type.
	ElementTypeUnknown ElementType = 0
	// ElementTypeHoldings groups the ETF's individual holdings.
	ElementTypeHoldings ElementType = 1
	// ElementTypeRegional groups holdings by region.
	ElementTypeRegional ElementType = 2
	// ElementTypeAssetClass groups holdings by asset class.
	ElementTypeAssetClass ElementType = 3
	// ElementTypeIndustry groups holdings by industry.
	ElementTypeIndustry ElementType = 4
)

// HoldingDetail is the holding detail of an ETF asset allocation element
// (holdings only).
type HoldingDetail struct {
	// IndustryID is the industry ID.
	IndustryID string
	// IndustryName is the industry name.
	IndustryName string
	// Index is the index counter ID (e.g. "BK/US/CP99000").
	Index string
	// IndexName is the index name.
	IndexName string
	// HoldingType is the holding type (e.g. "E" for stock).
	HoldingType string
	// HoldingTypeName is the holding type name.
	HoldingTypeName string
}

// AssetAllocationItem is one element of an ETF asset allocation group.
type AssetAllocationItem struct {
	// Name is the element name.
	Name string
	// Code is the security code (holdings only, e.g. "NVDA").
	Code string
	// PositionRatio is the position ratio (e.g. "0.0861114").
	PositionRatio string
	// Symbol is the security symbol (holdings only, e.g. "NVDA.US"), converted
	// from the API's counter_id. Empty for non-holdings groups.
	Symbol string
	// NameLocales maps a locale to the localized name (e.g. "zh-CN" → "英伟达").
	NameLocales map[string]string
	// HoldingDetail is the holding detail (holdings only); nil otherwise.
	HoldingDetail *HoldingDetail
}

// AssetAllocationGroup is one ETF asset allocation group (grouped by element
// type).
type AssetAllocationGroup struct {
	// ReportDate is the report date (e.g. "20260601").
	ReportDate string
	// AssetType is the element type of this group.
	AssetType ElementType
	// Lists are the elements of this group.
	Lists []*AssetAllocationItem
}

// AssetAllocationResponse is the response for FundamentalContext.EtfAssetAllocation.
type AssetAllocationResponse struct {
	// Info are the asset allocation groups.
	Info []*AssetAllocationGroup
}
