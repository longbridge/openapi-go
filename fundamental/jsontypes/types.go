// Package jsontypes contains raw JSON wire-format structs for the fundamental
// data API. Fields use the exact JSON field names from the API and are not
// exported with Go-idiomatic names; conversion happens in the parent package.
package jsontypes

import "encoding/json"

// ── financial_report ─────────────────────────────────────────────

// FinancialReports is the raw response for GET /v1/quote/financial-reports.
type FinancialReports struct {
	List json.RawMessage `json:"list"`
}

// ── dividend ─────────────────────────────────────────────────────

// DividendList is the raw response for GET /v1/quote/dividends and
// GET /v1/quote/dividends/details.
type DividendList struct {
	List []DividendItem `json:"list"`
}

// DividendItem is a single dividend / distribution event.
type DividendItem struct {
	CounterID   string `json:"counter_id"`
	ID          string `json:"id"`
	Desc        string `json:"desc"`
	RecordDate  string `json:"record_date"`
	ExDate      string `json:"ex_date"`
	PaymentDate string `json:"payment_date"`
}

// ── institution_rating ────────────────────────────────────────────

// InstitutionRatingLatest is the raw response for
// GET /v1/quote/institution-rating-latest.
type InstitutionRatingLatest struct {
	Evaluate       RatingEvaluate `json:"evaluate"`
	Target         RatingTarget   `json:"target"`
	IndustryID     int64          `json:"industry_id"`
	IndustryName   string         `json:"industry_name"`
	IndustryRank   int32          `json:"industry_rank"`
	IndustryTotal  int32          `json:"industry_total"`
	IndustryMean   int32          `json:"industry_mean"`
	IndustryMedian int32          `json:"industry_median"`
}

// RatingEvaluate holds analyst rating distribution counts.
type RatingEvaluate struct {
	Buy       int32  `json:"buy"`
	Over      int32  `json:"over"`
	Hold      int32  `json:"hold"`
	Under     int32  `json:"under"`
	Sell      int32  `json:"sell"`
	NoOpinion int32  `json:"no_opinion"`
	Total     int32  `json:"total"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// RatingTarget holds analyst target price range.
type RatingTarget struct {
	HighestPrice string `json:"highest_price"`
	LowestPrice  string `json:"lowest_price"`
	PrevClose    string `json:"prev_close"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

// InstitutionRatingSummary is the raw response for
// GET /v1/quote/institution-ratings.
type InstitutionRatingSummary struct {
	CcySymbol string                 `json:"ccy_symbol"`
	Change    string                 `json:"change"`
	Evaluate  RatingSummaryEvaluate  `json:"evaluate"`
	Recommend string                 `json:"recommend"`
	Target    string                 `json:"target"`
	UpdatedAt string                 `json:"updated_at"`
}

// RatingSummaryEvaluate is the simplified rating distribution for the
// consensus summary.
type RatingSummaryEvaluate struct {
	Buy       int32  `json:"buy"`
	Date      string `json:"date"`
	Hold      int32  `json:"hold"`
	Sell      int32  `json:"sell"`
	StrongBuy int32  `json:"strong_buy"`
	Under     int32  `json:"under"`
}

// ── institution_rating_detail ─────────────────────────────────────

// InstitutionRatingDetail is the raw response for
// GET /v1/quote/institution-ratings/detail.
type InstitutionRatingDetail struct {
	CcySymbol string                          `json:"ccy_symbol"`
	Evaluate  InstitutionRatingDetailEvaluate `json:"evaluate"`
	Target    InstitutionRatingDetailTarget   `json:"target"`
}

// InstitutionRatingDetailEvaluate holds a historical rating distribution
// time-series.
type InstitutionRatingDetailEvaluate struct {
	List []InstitutionRatingDetailEvaluateItem `json:"list"`
}

// InstitutionRatingDetailEvaluateItem is one weekly rating distribution
// snapshot.
type InstitutionRatingDetailEvaluateItem struct {
	Buy       int32  `json:"buy"`
	Date      string `json:"date"`
	Hold      int32  `json:"hold"`
	Sell      int32  `json:"sell"`
	StrongBuy int32  `json:"strong_buy"`
	NoOpinion int32  `json:"no_opinion"`
	Under     int32  `json:"under"`
}

// InstitutionRatingDetailTarget holds the historical target price time-series.
type InstitutionRatingDetailTarget struct {
	DataPercent        *string                              `json:"data_percent"`
	PredictionAccuracy string                               `json:"prediction_accuracy"`
	UpdatedAt          string                               `json:"updated_at"`
	List               []InstitutionRatingDetailTargetItem `json:"list"`
}

// InstitutionRatingDetailTargetItem is one weekly target price snapshot.
type InstitutionRatingDetailTargetItem struct {
	AvgTarget string `json:"avg_target"`
	Date      string `json:"date"`
	MaxTarget string `json:"max_target"`
	MinTarget string `json:"min_target"`
	Meet      bool   `json:"meet"`
	Price     string `json:"price"`
	Timestamp string `json:"timestamp"`
}

// ── forecast_eps ──────────────────────────────────────────────────

// ForecastEps is the raw response for GET /v1/quote/forecast-eps.
type ForecastEps struct {
	Items []ForecastEpsItem `json:"items"`
}

// ForecastEpsItem is one EPS forecast snapshot.
type ForecastEpsItem struct {
	ForecastEpsMedian  string `json:"forecast_eps_median"`
	ForecastEpsMean    string `json:"forecast_eps_mean"`
	ForecastEpsLowest  string `json:"forecast_eps_lowest"`
	ForecastEpsHighest string `json:"forecast_eps_highest"`
	InstitutionTotal   int32  `json:"institution_total"`
	InstitutionUp      int32  `json:"institution_up"`
	InstitutionDown    int32  `json:"institution_down"`
	ForecastStartDate  json.Number `json:"forecast_start_date"` // API returns string timestamp
	ForecastEndDate    json.Number `json:"forecast_end_date"` // API returns string timestamp
}

// ── consensus ─────────────────────────────────────────────────────

// FinancialConsensus is the raw response for
// GET /v1/quote/financial-consensus-detail.
type FinancialConsensus struct {
	List          []ConsensusReport `json:"list"`
	CurrentIndex  int32             `json:"current_index"`
	Currency      string            `json:"currency"`
	OptPeriods    []string          `json:"opt_periods"`
	CurrentPeriod string            `json:"current_period"`
}

// ConsensusReport is the consensus data for one fiscal period.
type ConsensusReport struct {
	FiscalYear   int32             `json:"fiscal_year"`
	FiscalPeriod string            `json:"fiscal_period"`
	PeriodText   string            `json:"period_text"`
	Details      []ConsensusDetail `json:"details"`
}

// ConsensusDetail is the consensus estimate for one financial metric.
type ConsensusDetail struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	Desc       string `json:"description"`
	Actual     string `json:"actual"`
	Estimate   string `json:"estimate"`
	CompValue  string `json:"comp_value"`
	CompDesc   string `json:"comp_desc"`
	Comp       string `json:"comp"`
	IsReleased bool   `json:"is_released"`
}

// ── valuation ─────────────────────────────────────────────────────

// ValuationData is the raw response for GET /v1/quote/valuation.
type ValuationData struct {
	Metrics ValuationMetricsData `json:"metrics"`
}

// ValuationMetricsData holds PE/PB/PS/dividend yield metric containers.
type ValuationMetricsData struct {
	PE     *ValuationMetricData `json:"pe"`
	PB     *ValuationMetricData `json:"pb"`
	PS     *ValuationMetricData `json:"ps"`
	DvdYld *ValuationMetricData `json:"dvd_yld"`
}

// ValuationMetricData holds the historical time-series for one valuation metric.
type ValuationMetricData struct {
	Desc   string          `json:"desc"`
	High   string          `json:"high"`
	Low    string          `json:"low"`
	Median string          `json:"median"`
	List   []ValuationPoint `json:"list"`
}

// ValuationPoint is one valuation data point.
type ValuationPoint struct {
	Timestamp json.Number `json:"timestamp"` // API returns either int64 or string
	Value     string      `json:"value"`
}

// ── valuation_history ─────────────────────────────────────────────

// ValuationHistoryResponse is the raw response for
// GET /v1/quote/valuation/detail.
type ValuationHistoryResponse struct {
	History ValuationHistoryData `json:"history"`
}

// ValuationHistoryData holds the historical valuation metrics container.
type ValuationHistoryData struct {
	Metrics ValuationHistoryMetrics `json:"metrics"`
}

// ValuationHistoryMetrics holds PE/PB/PS historical data.
type ValuationHistoryMetrics struct {
	PE *ValuationHistoryMetric `json:"pe"`
	PB *ValuationHistoryMetric `json:"pb"`
	PS *ValuationHistoryMetric `json:"ps"`
}

// ValuationHistoryMetric holds the historical data for one valuation metric
// including statistical bounds.
type ValuationHistoryMetric struct {
	Desc   string          `json:"desc"`
	High   string          `json:"high"`
	Low    string          `json:"low"`
	Median string          `json:"median"`
	List   []ValuationPoint `json:"list"`
}

// ── industry_valuation ────────────────────────────────────────────

// IndustryValuationList is the raw response for
// GET /v1/quote/industry-valuation-comparison.
type IndustryValuationList struct {
	List []IndustryValuationItem `json:"list"`
}

// IndustryValuationItem holds valuation data for one peer security.
type IndustryValuationItem struct {
	CounterID      string                    `json:"counter_id"`
	Name           string                    `json:"name"`
	Currency       string                    `json:"currency"`
	Assets         string                    `json:"assets"`
	Bps            string                    `json:"bps"`
	Eps            string                    `json:"eps"`
	Dps            string                    `json:"dps"`
	DivYld         string                    `json:"div_yld"`
	DivPayoutRatio string                    `json:"div_payout_ratio"`
	FiveYAvgDps    string                    `json:"five_y_avg_dps"`
	PE             string                    `json:"pe"`
	History        []IndustryValuationHistory `json:"history"`
}

// IndustryValuationHistory is a historical valuation snapshot for a peer.
type IndustryValuationHistory struct {
	Date string `json:"date"`
	PE   string `json:"pe"`
	PB   string `json:"pb"`
	PS   string `json:"ps"`
}

// ── industry_valuation_dist ───────────────────────────────────────

// IndustryValuationDist is the raw response for
// GET /v1/quote/industry-valuation-distribution.
type IndustryValuationDist struct {
	PE *ValuationDist `json:"pe"`
	PB *ValuationDist `json:"pb"`
	PS *ValuationDist `json:"ps"`
}

// ValuationDist holds distribution statistics for one valuation metric within
// an industry.
type ValuationDist struct {
	Low       string `json:"low"`
	High      string `json:"high"`
	Median    string `json:"median"`
	Value     string `json:"value"`
	Ranking   string `json:"ranking"`
	RankIndex string `json:"rank_index"`
	RankTotal string `json:"rank_total"`
}

// ── company ───────────────────────────────────────────────────────

// CompanyOverview is the raw response for GET /v1/quote/comp-overview.
type CompanyOverview struct {
	Name           string `json:"name"`
	CompanyName    string `json:"company_name"`
	Founded        string `json:"founded"`
	ListingDate    string `json:"listing_date"`
	Market         string `json:"market"`
	Region         string `json:"region"`
	Address        string `json:"address"`
	OfficeAddress  string `json:"office_address"`
	Website        string `json:"website"`
	IssuePrice     string `json:"issue_price"`
	SharesOffered  string `json:"shares_offered"`
	Chairman       string `json:"chairman"`
	Secretary      string `json:"secretary"`
	AuditInst      string `json:"audit_inst"`
	Category       string `json:"category"`
	YearEnd        string `json:"year_end"`
	Employees      string `json:"employees"`
	Phone          string `json:"Phone"`
	Fax            string `json:"fax"`
	Email          string `json:"email"`
	LegalRepr      string `json:"legal_repr"`
	Manager        string `json:"manager"`
	BusLicense     string `json:"bus_license"`
	AccountingFirm string `json:"accounting_firm"`
	SecuritiesRep  string `json:"securities_rep"`
	LegalCounsel   string `json:"legal_counsel"`
	ZipCode        string `json:"zip_code"`
	Ticker         string `json:"ticker"`
	Icon           string `json:"icon"`
	Profile        string `json:"profile"`
	AdsRatio       string `json:"ads_ratio"`
	Sector         int32  `json:"sector"`
}

// ── executive ─────────────────────────────────────────────────────

// ExecutiveList is the raw response for GET /v1/quote/company-professionals.
type ExecutiveList struct {
	ProfessionalList []ExecutiveGroup `json:"professional_list"`
}

// ExecutiveGroup holds executives for one security.
type ExecutiveGroup struct {
	CounterID     string         `json:"counter_id"`
	ForwardURL    string         `json:"forward_url"`
	Total         int32          `json:"total"`
	Professionals []Professional `json:"professionals"`
}

// Professional is one executive / board member.
type Professional struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	NameZhCN   string `json:"name_zhcn"`
	NameEn     string `json:"name_en"`
	Title      string `json:"title"`
	Biography  string `json:"biography"`
	Photo      string `json:"photo"`
	WikiURL    string `json:"wiki_url"`
}

// ── shareholder ───────────────────────────────────────────────────

// ShareholderList is the raw response for GET /v1/quote/shareholders.
type ShareholderList struct {
	ShareholderList []Shareholder `json:"shareholder_list"`
	ForwardURL      string        `json:"forward_url"`
	Total           int32         `json:"total"`
}

// Shareholder is one major shareholder.
type Shareholder struct {
	ShareholderID   string            `json:"shareholder_id"`
	ShareholderName string            `json:"shareholder_name"`
	InstitutionType string            `json:"institution_type"`
	PercentOfShares string            `json:"percent_of_shares"`
	SharesChanged   string            `json:"shares_changed"`
	ReportDate      string            `json:"report_date"`
	Stocks          []ShareholderStock `json:"stocks"`
}

// ShareholderStock is a security in an institutional shareholder's
// cross-holdings.
type ShareholderStock struct {
	CounterID string `json:"counter_id"`
	Code      string `json:"code"`
	Market    string `json:"market"`
	Chg       string `json:"chg"`
}

// ── fund_holder ───────────────────────────────────────────────────

// FundHolders is the raw response for GET /v1/quote/fund-holders.
type FundHolders struct {
	Lists []FundHolder `json:"lists"`
}

// FundHolder is a fund or ETF that holds the queried security.
type FundHolder struct {
	Code          string `json:"code"`
	CounterID     string `json:"counter_id"`
	Currency      string `json:"currency"`
	Name          string `json:"name"`
	PositionRatio string `json:"position_ratio"`
	ReportDate    string `json:"report_date"`
}

// ── corp_action ───────────────────────────────────────────────────

// CorpActions is the raw response for GET /v1/quote/company-act.
type CorpActions struct {
	Items []CorpActionItem `json:"items"`
}

// CorpActionItem is one corporate action event.
type CorpActionItem struct {
	ID           string           `json:"id"`
	Date         string           `json:"date"`
	DateStr      string           `json:"date_str"`
	DateType     string           `json:"date_type"`
	DateZone     string           `json:"date_zone"`
	ActType      string           `json:"act_type"`
	ActDesc      string           `json:"act_desc"`
	Action       string           `json:"action"`
	Recent       bool             `json:"recent"`
	IsDelay      bool             `json:"is_delay"`
	DelayContent string           `json:"delay_content"`
	Live         *CorpActionLive  `json:"live"`
	Security     *json.RawMessage `json:"security"`
}

// CorpActionLive is the live stream associated with a corporate action.
type CorpActionLive struct {
	ID        string          `json:"id"`
	Status    json.RawMessage `json:"status"`
	StartedAt string          `json:"started_at"`
	Name      string          `json:"name"`
	Icon      string          `json:"icon"`
}

// ── invest_relation ───────────────────────────────────────────────

// InvestRelations is the raw response for GET /v1/quote/invest-relations.
type InvestRelations struct {
	ForwardURL       string          `json:"forward_url"`
	InvestSecurities []InvestSecurity `json:"invest_securities"`
}

// InvestSecurity is a security in which the queried company has an investment
// stake.
type InvestSecurity struct {
	CompanyID      string `json:"company_id"`
	CompanyName    string `json:"company_name"`
	CompanyNameEn  string `json:"company_name_en"`
	CompanyNameZhCN string `json:"company_name_zhcn"`
	CounterID      string `json:"counter_id"`
	Currency       string `json:"currency"`
	PercentOfShares string `json:"percent_of_shares"`
	SharesRank     string `json:"shares_rank"`
	SharesValue    string `json:"shares_value"`
}

// ── operating ─────────────────────────────────────────────────────

// OperatingList is the raw response for GET /v1/quote/operatings.
type OperatingList struct {
	List []OperatingItem `json:"list"`
}

// OperatingItem is one operating summary report (annual / quarterly).
type OperatingItem struct {
	ID        string              `json:"id"`
	Report    string              `json:"report"`
	Title     string              `json:"title"`
	Txt       string              `json:"txt"`
	Latest    bool                `json:"latest"`
	Keywords  []json.RawMessage   `json:"keywords"`
	WebURL    string              `json:"web_url"`
	Financial OperatingFinancial  `json:"financial"`
}

// OperatingFinancial holds key financial metrics extracted from an operating
// report.
type OperatingFinancial struct {
	Code       string               `json:"code"`
	CounterID  string               `json:"counter_id"`
	Currency   string               `json:"currency"`
	Name       string               `json:"name"`
	Region     string               `json:"region"`
	Report     string               `json:"report"`
	ReportTxt  string               `json:"report_txt"`
	Indicators []OperatingIndicator `json:"indicators"`
}

// OperatingIndicator is one financial indicator in an operating report.
type OperatingIndicator struct {
	FieldName      string `json:"field_name"`
	IndicatorName  string `json:"indicator_name"`
	IndicatorValue string `json:"indicator_value"`
	Yoy            string `json:"yoy"`
}

// ── buyback ───────────────────────────────────────────────────────

// BuybackData is the raw response for GET /v1/quote/buy-backs.
type BuybackData struct {
	RecentBuybacks  *RecentBuybacks     `json:"recent_buybacks"`
	BuybackHistory  []BuybackHistoryItem `json:"buyback_history"`
	BuybackRatios   []BuybackRatios     `json:"buyback_ratios"`
}

// RecentBuybacks is the TTM (trailing twelve months) buyback summary.
type RecentBuybacks struct {
	Currency             string `json:"currency"`
	NetBuybackTTM        string `json:"net_buyback_ttm"`
	NetBuybackYieldTTM   string `json:"net_buyback_yield_ttm"`
}

// BuybackHistoryItem is one historical annual buyback data point.
type BuybackHistoryItem struct {
	FiscalYear            string `json:"fiscal_year"`
	FiscalYearRange       string `json:"fiscal_year_range"`
	NetBuyback            string `json:"net_buyback"`
	NetBuybackYield       string `json:"net_buyback_yield"`
	NetBuybackGrowthRate  string `json:"net_buyback_growth_rate"`
	Currency              string `json:"currency"`
}

// BuybackRatios holds buyback payout and cash-flow ratios.
type BuybackRatios struct {
	NetBuybackPayoutRatio       string `json:"net_buyback_payout_ratio"`
	NetBuybackToCashflowRatio   string `json:"net_buyback_to_cashflow_ratio"`
}

// ── ratings ───────────────────────────────────────────────────────

// StockRatings is the raw response for GET /v1/quote/ratings.
type StockRatings struct {
	StyleTxtName       string              `json:"style_txt_name"`
	ScaleTxtName       string              `json:"scale_txt_name"`
	ReportPeriodTxt    string              `json:"report_period_txt"`
	MultiScore         json.RawMessage     `json:"multi_score"`
	MultiLetter        string              `json:"multi_letter"`
	MultiScoreChange   int32               `json:"multi_score_change"`
	IndustryName       string              `json:"industry_name"`
	IndustryRank       json.RawMessage     `json:"industry_rank"`
	IndustryTotal      json.RawMessage     `json:"industry_total"`
	IndustryMeanScore  json.RawMessage     `json:"industry_mean_score"`
	IndustryMedianScore json.RawMessage    `json:"industry_median_score"`
	Ratings            []RatingCategory    `json:"ratings"`
}

// RatingCategory is one rating category (e.g. growth, profitability).
type RatingCategory struct {
	Kind          int32                    `json:"type"`
	SubIndicators []RatingSubIndicatorGroup `json:"sub_indicators"`
}

// RatingSubIndicatorGroup is a group of sub-indicators under one category.
type RatingSubIndicatorGroup struct {
	Indicator     RatingIndicator      `json:"indicator"`
	SubIndicators []RatingLeafIndicator `json:"sub_indicators"`
}

// RatingIndicator is a rating indicator node.
type RatingIndicator struct {
	Name   string          `json:"name"`
	Score  json.RawMessage `json:"score"`
	Letter string          `json:"letter"`
}

// RatingLeafIndicator is a leaf rating indicator with a raw value.
type RatingLeafIndicator struct {
	Name      string          `json:"name"`
	Value     string          `json:"value"`
	ValueType string          `json:"value_type"`
	Score     json.RawMessage `json:"score"`
	Letter    string          `json:"letter"`
}

// ── business_segments ─────────────────────────────────────────────

// BusinessSegments is the raw response for
// GET /v1/quote/fundamentals/business-segments.
type BusinessSegments struct {
	Date     string                `json:"date"`
	Total    string                `json:"total"`
	Currency string                `json:"currency"`
	Business []BusinessSegmentItem `json:"business"`
}

// BusinessSegmentItem is one business segment entry (latest snapshot).
type BusinessSegmentItem struct {
	Name    string `json:"name"`
	Percent string `json:"percent"`
}

// BusinessSegmentsHistory is the raw response for
// GET /v1/quote/fundamentals/business-segments/history.
type BusinessSegmentsHistory struct {
	Historical []BusinessSegmentsHistoricalItem `json:"historical"`
}

// BusinessSegmentsHistoricalItem is one historical business segments snapshot.
type BusinessSegmentsHistoricalItem struct {
	Date      string                       `json:"date"`
	Total     string                       `json:"total"`
	Currency  string                       `json:"currency"`
	Business  []BusinessSegmentHistoryItem `json:"business"`
	Regionals []BusinessSegmentHistoryItem `json:"regionals"`
}

// BusinessSegmentHistoryItem is one business/regional segment entry in a
// historical snapshot.
type BusinessSegmentHistoryItem struct {
	Name    string `json:"name"`
	Percent string `json:"percent"`
	Value   string `json:"value"`
}

// ── institution_rating_views ──────────────────────────────────────

// InstitutionRatingViews is the raw response for
// GET /v1/quote/ratings/institutional.
type InstitutionRatingViews struct {
	Elist []InstitutionRatingViewItem `json:"elist"`
}

// InstitutionRatingViewItem is one historical rating distribution snapshot.
type InstitutionRatingViewItem struct {
	Date  json.Number `json:"date"` // int64 or quoted string timestamp
	Buy   string      `json:"buy"`
	Over  string      `json:"over"`
	Hold  string      `json:"hold"`
	Under string      `json:"under"`
	Sell  string      `json:"sell"`
	Total string      `json:"total"`
}

// ── industry_rank ─────────────────────────────────────────────────

// IndustryRankResponse is the raw response for GET /v1/quote/industry/rank.
type IndustryRankResponse struct {
	Items []IndustryRankGroup `json:"items"`
}

// IndustryRankGroup is a group of ranked industry items.
type IndustryRankGroup struct {
	Lists []IndustryRankItem `json:"lists"`
}

// IndustryRankItem is one ranked industry item.
type IndustryRankItem struct {
	Name          string `json:"name"`
	CounterID     string `json:"counter_id"`
	Chg           string `json:"chg"`
	LeadingName   string `json:"leading_name"`
	LeadingTicker string `json:"leading_ticker"`
	LeadingChg    string `json:"leading_chg"`
	ValueName     string `json:"value_name"`
	ValueData     string `json:"value_data"`
}

// ── industry_peers ────────────────────────────────────────────────

// IndustryPeersResponse is the raw response for
// GET /v1/quote/industries/peers.
type IndustryPeersResponse struct {
	Top   IndustryPeersTop  `json:"top"`
	Chain *IndustryPeerNode `json:"chain"`
}

// IndustryPeersTop holds the top-level industry info.
type IndustryPeersTop struct {
	Name   string `json:"name"`
	Market string `json:"market"`
}

// IndustryPeerNode is a node in the recursive industry peer chain.
type IndustryPeerNode struct {
	Name      string             `json:"name"`
	CounterID string             `json:"counter_id"`
	StockNum  int32              `json:"stock_num"`
	Chg       string             `json:"chg"`
	YtdChg    string             `json:"ytd_chg"`
	Next      []IndustryPeerNode `json:"next"`
}

// ── financial_report_snapshot ─────────────────────────────────────

// FinancialReportSnapshot is the raw response for
// GET /v1/quote/financials/earnings-snapshot.
type FinancialReportSnapshot struct {
	Name              string                  `json:"name"`
	Ticker            string                  `json:"ticker"`
	FpStart           string                  `json:"fp_start"`
	FpEnd             string                  `json:"fp_end"`
	Currency          string                  `json:"currency"`
	ReportDesc        string                  `json:"report_desc"`
	FoRevenue         *SnapshotForecastMetric `json:"fo_revenue"`
	FoEbit            *SnapshotForecastMetric `json:"fo_ebit"`
	FoEps             *SnapshotForecastMetric `json:"fo_eps"`
	FrRevenue         *SnapshotReportedMetric `json:"fr_revenue"`
	FrProfit          *SnapshotReportedMetric `json:"fr_profit"`
	FrOperateCash     *SnapshotReportedMetric `json:"fr_operate_cash"`
	FrInvestCash      *SnapshotReportedMetric `json:"fr_invest_cash"`
	FrFinanceCash     *SnapshotReportedMetric `json:"fr_finance_cash"`
	FrTotalAssets     *SnapshotReportedMetric `json:"fr_total_assets"`
	FrTotalLiability  *SnapshotReportedMetric `json:"fr_total_liability"`
	FrRoeTtm          string                  `json:"fr_roe_ttm"`
	FrProfitMargin    string                  `json:"fr_profit_margin"`
	FrProfitMarginTtm string                  `json:"fr_profit_margin_ttm"`
	FrAssetTurnTtm    string                  `json:"fr_asset_turn_ttm"`
	FrLeverageTtm     string                  `json:"fr_leverage_ttm"`
	FrDebtAssetsRatio string                  `json:"fr_debt_assets_ratio"`
}

// SnapshotForecastMetric is a forecast metric in the financial report snapshot.
type SnapshotForecastMetric struct {
	Value    string `json:"value"`
	Yoy      string `json:"yoy"`
	CmpDesc  string `json:"cmp_desc"`
	EstValue string `json:"est_value"`
}

// SnapshotReportedMetric is a reported metric in the financial report snapshot.
type SnapshotReportedMetric struct {
	Value string `json:"value"`
	Yoy   string `json:"yoy"`
}

// ── shareholder_top ──────────────────────────────────────────────

// ShareholderTopResponse is the raw response for GET /v1/quote/shareholders/top.
type ShareholderTopResponse struct {
	Data json.RawMessage `json:"data"`
}

// ── shareholder_detail ───────────────────────────────────────────

// ShareholderDetailResponse is the raw response for GET /v1/quote/shareholders/holding.
type ShareholderDetailResponse struct {
	Data json.RawMessage `json:"data"`
}

// ── valuation_comparison ─────────────────────────────────────────

// ValuationComparisonResponse is the raw response for GET /v1/quote/compare/valuation.
type ValuationComparisonResponse struct {
	Data json.RawMessage `json:"data"`
}

// ── etf_asset_allocation ─────────────────────────────────────────

// AssetAllocationResponse is the raw response for
// GET /v1/quote/etf-asset-allocation.
type AssetAllocationResponse struct {
	Info []AssetAllocationGroup `json:"info"`
}

// AssetAllocationGroup is one ETF asset allocation group (grouped by element
// type).
type AssetAllocationGroup struct {
	ReportDate string                `json:"report_date"`
	AssetType  int32                 `json:"asset_type"`
	Lists      []AssetAllocationItem `json:"lists"`
}

// AssetAllocationItem is one element of an ETF asset allocation group.
type AssetAllocationItem struct {
	Name          string            `json:"name"`
	Code          string            `json:"code"`
	PositionRatio string            `json:"position_ratio"`
	CounterID     string            `json:"counter_id"`
	NameLocales   map[string]string `json:"name_locales_map"`
	HoldingDetail *HoldingDetail    `json:"holding_detail"`
}

// HoldingDetail is the holding detail of an ETF asset allocation element
// (holdings only).
type HoldingDetail struct {
	IndustryID      string `json:"industry_id"`
	IndustryName    string `json:"industry_name"`
	Index           string `json:"index"`
	IndexName       string `json:"index_name"`
	HoldingType     string `json:"holding_type"`
	HoldingTypeName string `json:"holding_type_name"`
}
