package jsontypes

import "encoding/json"

type FinancialReports struct {
	List json.RawMessage `json:"list"`
}

type DividendList struct {
	List []*DividendItem `json:"list"`
}

type DividendItem struct {
	CounterId   string `json:"counter_id"`
	Id          string `json:"id"`
	Desc        string `json:"desc"`
	RecordDate  string `json:"record_date"`
	ExDate      string `json:"ex_date"`
	PaymentDate string `json:"payment_date"`
}

type InstitutionRatingLatest struct {
	Evaluate       RatingEvaluate `json:"evaluate"`
	Target         RatingTarget   `json:"target"`
	IndustryId     int64          `json:"industry_id"`
	IndustryName   string         `json:"industry_name"`
	IndustryRank   int32          `json:"industry_rank"`
	IndustryTotal  int32          `json:"industry_total"`
	IndustryMean   int32          `json:"industry_mean"`
	IndustryMedian int32          `json:"industry_median"`
}

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

type RatingTarget struct {
	HighestPrice string `json:"highest_price"`
	LowestPrice  string `json:"lowest_price"`
	PrevClose    string `json:"prev_close"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

type InstitutionRatingSummary struct {
	CcySymbol string                `json:"ccy_symbol"`
	Change    string                `json:"change"`
	Evaluate  RatingSummaryEvaluate `json:"evaluate"`
	Recommend string                `json:"recommend"`
	Target    string                `json:"target"`
	UpdatedAt string                `json:"updated_at"`
}

type RatingSummaryEvaluate struct {
	Buy      int32  `json:"buy"`
	Date     string `json:"date"`
	Hold     int32  `json:"hold"`
	Sell     int32  `json:"sell"`
	StrongBuy int32 `json:"strong_buy"`
	Under    int32  `json:"under"`
}

type InstitutionRatingDetail struct {
	CcySymbol string                          `json:"ccy_symbol"`
	Evaluate  InstitutionRatingDetailEvaluate `json:"evaluate"`
	Target    InstitutionRatingDetailTarget   `json:"target"`
}

type InstitutionRatingDetailEvaluate struct {
	List []*InstitutionRatingDetailEvaluateItem `json:"list"`
}

type InstitutionRatingDetailEvaluateItem struct {
	Buy       int32  `json:"buy"`
	Date      string `json:"date"`
	Hold      int32  `json:"hold"`
	Sell      int32  `json:"sell"`
	StrongBuy int32  `json:"strong_buy"`
	NoOpinion int32  `json:"no_opinion"`
	Under     int32  `json:"under"`
}

type InstitutionRatingDetailTarget struct {
	DataPercent        string                               `json:"data_percent"`
	PredictionAccuracy string                               `json:"prediction_accuracy"`
	UpdatedAt          string                               `json:"updated_at"`
	List               []*InstitutionRatingDetailTargetItem `json:"list"`
}

type InstitutionRatingDetailTargetItem struct {
	AvgTarget string `json:"avg_target"`
	Date      string `json:"date"`
	MaxTarget string `json:"max_target"`
	MinTarget string `json:"min_target"`
	Meet      bool   `json:"meet"`
	Price     string `json:"price"`
	Timestamp string `json:"timestamp"`
}

type ForecastEps struct {
	Items []*ForecastEpsItem `json:"items"`
}

type ForecastEpsItem struct {
	ForecastEpsMedian  string `json:"forecast_eps_median"`
	ForecastEpsMean    string `json:"forecast_eps_mean"`
	ForecastEpsLowest  string `json:"forecast_eps_lowest"`
	ForecastEpsHighest string `json:"forecast_eps_highest"`
	InstitutionTotal   int32  `json:"institution_total"`
	InstitutionUp      int32  `json:"institution_up"`
	InstitutionDown    int32  `json:"institution_down"`
	ForecastStartDate  int64  `json:"forecast_start_date"`
	ForecastEndDate    int64  `json:"forecast_end_date"`
}

type FinancialConsensus struct {
	List          []*ConsensusReport `json:"list"`
	CurrentIndex  int32              `json:"current_index"`
	Currency      string             `json:"currency"`
	OptPeriods    []string           `json:"opt_periods"`
	CurrentPeriod string             `json:"current_period"`
}

type ConsensusReport struct {
	FiscalYear   int32              `json:"fiscal_year"`
	FiscalPeriod string             `json:"fiscal_period"`
	PeriodText   string             `json:"period_text"`
	Details      []*ConsensusDetail `json:"details"`
}

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

type ValuationData struct {
	Metrics ValuationMetricsData `json:"metrics"`
}

type ValuationMetricsData struct {
	Pe     *ValuationMetricData `json:"pe"`
	Pb     *ValuationMetricData `json:"pb"`
	Ps     *ValuationMetricData `json:"ps"`
	DvdYld *ValuationMetricData `json:"dvd_yld"`
}

type ValuationMetricData struct {
	Desc   string           `json:"desc"`
	High   string           `json:"high"`
	Low    string           `json:"low"`
	Median string           `json:"median"`
	List   []*ValuationPoint `json:"list"`
}

type ValuationPoint struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}

type ValuationHistoryResponse struct {
	History ValuationHistoryData `json:"history"`
}

type ValuationHistoryData struct {
	Metrics ValuationHistoryMetrics `json:"metrics"`
}

type ValuationHistoryMetrics struct {
	Pe *ValuationHistoryMetric `json:"pe"`
	Pb *ValuationHistoryMetric `json:"pb"`
	Ps *ValuationHistoryMetric `json:"ps"`
}

type ValuationHistoryMetric struct {
	Desc   string           `json:"desc"`
	High   string           `json:"high"`
	Low    string           `json:"low"`
	Median string           `json:"median"`
	List   []*ValuationPoint `json:"list"`
}

type IndustryValuationList struct {
	List []*IndustryValuationItem `json:"list"`
}

type IndustryValuationItem struct {
	CounterId      string                      `json:"counter_id"`
	Name           string                      `json:"name"`
	Currency       string                      `json:"currency"`
	Assets         string                      `json:"assets"`
	Bps            string                      `json:"bps"`
	Eps            string                      `json:"eps"`
	Dps            string                      `json:"dps"`
	DivYld         string                      `json:"div_yld"`
	DivPayoutRatio string                      `json:"div_payout_ratio"`
	FiveYAvgDps    string                      `json:"five_y_avg_dps"`
	Pe             string                      `json:"pe"`
	History        []*IndustryValuationHistory `json:"history"`
}

type IndustryValuationHistory struct {
	Date string `json:"date"`
	Pe   string `json:"pe"`
	Pb   string `json:"pb"`
	Ps   string `json:"ps"`
}

type IndustryValuationDist struct {
	Pe *ValuationDist `json:"pe"`
	Pb *ValuationDist `json:"pb"`
	Ps *ValuationDist `json:"ps"`
}

type ValuationDist struct {
	Low       string `json:"low"`
	High      string `json:"high"`
	Median    string `json:"median"`
	Value     string `json:"value"`
	Ranking   string `json:"ranking"`
	RankIndex string `json:"rank_index"`
	RankTotal string `json:"rank_total"`
}

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

type ExecutiveList struct {
	ProfessionalList []*ExecutiveGroup `json:"professional_list"`
}

type ExecutiveGroup struct {
	CounterId     string          `json:"counter_id"`
	ForwardUrl    string          `json:"forward_url"`
	Total         int32           `json:"total"`
	Professionals []*Professional `json:"professionals"`
}

type Professional struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	NameZhcn  string `json:"name_zhcn"`
	NameEn    string `json:"name_en"`
	Title     string `json:"title"`
	Biography string `json:"biography"`
	Photo     string `json:"photo"`
	WikiUrl   string `json:"wiki_url"`
}

type ShareholderList struct {
	Shareholders []*Shareholder `json:"shareholder_list"`
	ForwardUrl   string         `json:"forward_url"`
	Total        int32          `json:"total"`
}

type Shareholder struct {
	ShareholderId   string              `json:"shareholder_id"`
	ShareholderName string              `json:"shareholder_name"`
	InstitutionType string              `json:"institution_type"`
	PercentOfShares string              `json:"percent_of_shares"`
	SharesChanged   string              `json:"shares_changed"`
	ReportDate      string              `json:"report_date"`
	Stocks          []*ShareholderStock `json:"stocks"`
}

type ShareholderStock struct {
	CounterId string `json:"counter_id"`
	Code      string `json:"code"`
	Market    string `json:"market"`
	Chg       string `json:"chg"`
}

type FundHolders struct {
	Lists []*FundHolder `json:"lists"`
}

type FundHolder struct {
	Code          string `json:"code"`
	CounterId     string `json:"counter_id"`
	Currency      string `json:"currency"`
	Name          string `json:"name"`
	PositionRatio string `json:"position_ratio"`
	ReportDate    string `json:"report_date"`
}

type CorpActions struct {
	Items []*CorpActionItem `json:"items"`
}

type CorpActionItem struct {
	Id           string          `json:"id"`
	Date         string          `json:"date"`
	DateStr      string          `json:"date_str"`
	DateType     string          `json:"date_type"`
	DateZone     string          `json:"date_zone"`
	ActType      string          `json:"act_type"`
	ActDesc      string          `json:"act_desc"`
	Action       string          `json:"action"`
	Recent       bool            `json:"recent"`
	IsDelay      bool            `json:"is_delay"`
	DelayContent string          `json:"delay_content"`
	Live         *CorpActionLive `json:"live"`
	Security     json.RawMessage `json:"security"`
}

type CorpActionLive struct {
	Id        string          `json:"id"`
	Status    json.RawMessage `json:"status"`
	StartedAt string          `json:"started_at"`
	Name      string          `json:"name"`
	Icon      string          `json:"icon"`
}

type InvestRelations struct {
	ForwardUrl       string            `json:"forward_url"`
	InvestSecurities []*InvestSecurity `json:"invest_securities"`
}

type InvestSecurity struct {
	CompanyId       string `json:"company_id"`
	CompanyName     string `json:"company_name"`
	CompanyNameEn   string `json:"company_name_en"`
	CompanyNameZhcn string `json:"company_name_zhcn"`
	CounterId       string `json:"counter_id"`
	Currency        string `json:"currency"`
	PercentOfShares string `json:"percent_of_shares"`
	SharesRank      string `json:"shares_rank"`
	SharesValue     string `json:"shares_value"`
}

type OperatingList struct {
	List []*OperatingItem `json:"list"`
}

type OperatingItem struct {
	Id        string             `json:"id"`
	Report    string             `json:"report"`
	Title     string             `json:"title"`
	Txt       string             `json:"txt"`
	Latest    bool               `json:"latest"`
	Keywords  json.RawMessage    `json:"keywords"`
	WebUrl    string             `json:"web_url"`
	Financial OperatingFinancial `json:"financial"`
}

type OperatingFinancial struct {
	Code       string               `json:"code"`
	CounterId  string               `json:"counter_id"`
	Currency   string               `json:"currency"`
	Name       string               `json:"name"`
	Region     string               `json:"region"`
	Report     string               `json:"report"`
	ReportTxt  string               `json:"report_txt"`
	Indicators []*OperatingIndicator `json:"indicators"`
}

type OperatingIndicator struct {
	FieldName      string `json:"field_name"`
	IndicatorName  string `json:"indicator_name"`
	IndicatorValue string `json:"indicator_value"`
	Yoy            string `json:"yoy"`
}

type BuybackData struct {
	RecentBuybacks *RecentBuybacks      `json:"recent_buybacks"`
	BuybackHistory []*BuybackHistoryItem `json:"buyback_history"`
	BuybackRatios  []*BuybackRatios     `json:"buyback_ratios"`
}

type RecentBuybacks struct {
	Currency           string `json:"currency"`
	NetBuybackTtm      string `json:"net_buyback_ttm"`
	NetBuybackYieldTtm string `json:"net_buyback_yield_ttm"`
}

type BuybackHistoryItem struct {
	FiscalYear           string `json:"fiscal_year"`
	FiscalYearRange      string `json:"fiscal_year_range"`
	NetBuyback           string `json:"net_buyback"`
	NetBuybackYield      string `json:"net_buyback_yield"`
	NetBuybackGrowthRate string `json:"net_buyback_growth_rate"`
	Currency             string `json:"currency"`
}

type BuybackRatios struct {
	NetBuybackPayoutRatio     string `json:"net_buyback_payout_ratio"`
	NetBuybackToCashflowRatio string `json:"net_buyback_to_cashflow_ratio"`
}

type StockRatings struct {
	StyleTxtName        string            `json:"style_txt_name"`
	ScaleTxtName        string            `json:"scale_txt_name"`
	ReportPeriodTxt     string            `json:"report_period_txt"`
	MultiScore          json.RawMessage   `json:"multi_score"`
	MultiLetter         string            `json:"multi_letter"`
	MultiScoreChange    int32             `json:"multi_score_change"`
	IndustryName        string            `json:"industry_name"`
	IndustryRank        json.RawMessage   `json:"industry_rank"`
	IndustryTotal       json.RawMessage   `json:"industry_total"`
	IndustryMeanScore   json.RawMessage   `json:"industry_mean_score"`
	IndustryMedianScore json.RawMessage   `json:"industry_median_score"`
	Ratings             []*RatingCategory `json:"ratings"`
}

type RatingCategory struct {
	Kind          int32                      `json:"type"`
	SubIndicators []*RatingSubIndicatorGroup `json:"sub_indicators"`
}

type RatingSubIndicatorGroup struct {
	Indicator     RatingIndicator       `json:"indicator"`
	SubIndicators []*RatingLeafIndicator `json:"sub_indicators"`
}

type RatingIndicator struct {
	Name   string          `json:"name"`
	Score  json.RawMessage `json:"score"`
	Letter string          `json:"letter"`
}

type RatingLeafIndicator struct {
	Name      string          `json:"name"`
	Value     string          `json:"value"`
	ValueType string          `json:"value_type"`
	Score     json.RawMessage `json:"score"`
	Letter    string          `json:"letter"`
}
