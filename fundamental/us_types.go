package fundamental

// RankTag is an industry-rank label returned by CompanyOverview.
type RankTag struct {
	Key           string `json:"key"`
	Location      int32  `json:"location"`
	Title         string `json:"title"`
	Text          string `json:"text"`
	RankType      int32  `json:"rank_type"`
	HighlightText string `json:"highlight_text"`
}

// USSharelistItem is one entry in USCompanyOverview.ShareList.
type USSharelistItem struct {
	Chg  string `json:"chg"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// USCompanyOverview is the US company summary snapshot.
type USCompanyOverview struct {
	Intro       string            `json:"intro"`
	MarketCap   string            `json:"market_cap"`
	CcySymbol   string            `json:"ccy_symbol"`
	TopRankTags []RankTag         `json:"top_rank_tags"`
	ShareList   []USSharelistItem `json:"sharelist"`
	DetailURL   string            `json:"detail_url"`
}

// ValuationMetric is one valuation indicator entry within ValuationOverview.Metrics.
// Keys include "pe", "pb", "ps", etc.
type ValuationMetric struct {
	Circle         string `json:"circle"`
	Part           string `json:"part"`
	Metric         string `json:"metric"`
	Desc           string `json:"desc"`
	IndustryMedian string `json:"industry_median"`
}

// ValuationOverview is the US valuation snapshot.
type ValuationOverview struct {
	Metrics    map[string]ValuationMetric `json:"metrics"`
	Indicator  string                     `json:"indicator"`
	Range      int32                      `json:"range"`
	Date       string                     `json:"date"`
	CcySymbol  string                     `json:"ccy_symbol"`
	AIChatData AIChatData                 `json:"aichat_data"`
	AISummary  string                     `json:"ai_summary"`
}

// USReportPeriod identifies a reporting period (quarter/half/annual).
type USReportPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	ReportTxt string `json:"report_txt"`
}

// FinancialISItem is one income-statement entry in FinancialOverview.
type FinancialISItem struct {
	Revenue   string         `json:"revenue"`
	NetIncome string         `json:"net_income"`
	NetMargin string         `json:"net_margin"`
	Report    USReportPeriod `json:"report"`
}

// FinancialBSItem is one balance-sheet entry in FinancialOverview.
type FinancialBSItem struct {
	DebtAssetsRatio  string         `json:"debt_assets_ratio"`
	TotalAssets      string         `json:"total_assets"`
	TotalLiabilities string         `json:"total_liabilities"`
	Report           USReportPeriod `json:"report"`
}

// FinancialCFItem is one cash-flow entry in FinancialOverview.
type FinancialCFItem struct {
	Operating string         `json:"operating"`
	Investing string         `json:"investing"`
	Financing string         `json:"financing"`
	Report    USReportPeriod `json:"report"`
}

// FinancialOverview is the US financial overview containing income statement,
// balance sheet, and cash flow summaries by reporting period.
type FinancialOverview struct {
	CcySymbol  string            `json:"ccy_symbol"`
	ReportType string            `json:"report_type"`
	ISList     []FinancialISItem `json:"is_list"`
	BSList     []FinancialBSItem `json:"bs_list"`
	CFList     []FinancialCFItem `json:"cf_list"`
}

// FinancialStatementField is one line item within a FinancialStatementPeriod.
type FinancialStatementField struct {
	DisplayOrder int32  `json:"display_order"`
	Field        string `json:"field"`
	ID           string `json:"id"`
	Level        int64  `json:"level"`
	Name         string `json:"name"`
	Value        string `json:"value"`
	ValueType    string `json:"value_type"`
	YoY          string `json:"yoy"`
}

// FinancialStatementPeriod is one reporting period in FinancialStatement.
type FinancialStatementPeriod struct {
	FfPeriod  string                    `json:"ff_period"`
	FfYear    int32                     `json:"ff_year"`
	Fields    []FinancialStatementField `json:"fields"`
	FpEnd     string                    `json:"fp_end"`
	ReportTxt string                    `json:"report_txt"`
	RptDate   string                    `json:"rpt_date"`
}

// FinancialStatement is the US financial statement (IS/BS/CF).
// kind controls which statement is returned: "IS" (income), "BS" (balance sheet), "CF" (cash flow).
// EmptyFields lists field IDs that the API could not populate.
type FinancialStatement struct {
	Currency    string                     `json:"currency"`
	Report      string                     `json:"report"`
	List        []FinancialStatementPeriod `json:"list"`
	EmptyFields []string                   `json:"empty_fields"`
}

// KeyMetricItem is one period entry in KeyFinancialMetrics.
type KeyMetricItem struct {
	FfPeriod  string      `json:"ff_period"`
	FfYear    int32       `json:"ff_year"`
	FpEnd     string      `json:"fp_end"`
	ReportTxt string      `json:"report_txt"`
	RptDate   string      `json:"rpt_date"`
	Fields    []interface{} `json:"fields"` // metric values; element shape varies per field set
}

// KeyFinancialMetrics holds per-period key ratios (ROE, margins, debt ratio).
type KeyFinancialMetrics struct {
	Currency    string          `json:"currency"`
	Report      string          `json:"report"`
	EmptyFields []string        `json:"empty_fields"`
	List        []KeyMetricItem `json:"list"`
}

// AIChatData holds the AI chat context embedded in analyst responses.
type AIChatData struct {
	AgentID        string `json:"agent_id"`
	HandoffAgentID string `json:"handoff_agent_id"`
	Symbol         string `json:"symbol"`
	Text           string `json:"text"`
	Type           string `json:"type"`
	WorkflowType   string `json:"workflow_type"`
}

// USConsensusEstimate holds actual vs estimated values for one metric.
type USConsensusEstimate struct {
	Actual   string `json:"actual"`
	Estimate string `json:"estimate"`
}

// USConsensusItem is one fiscal-year entry in AnalystConsensus.List.
type USConsensusItem struct {
	EBIT       USConsensusEstimate `json:"ebit"`
	EPS        USConsensusEstimate `json:"eps"`
	FiscalYear int64               `json:"fiscal_year"`
	ReportTxt  string              `json:"report_txt"`
	Revenue    USConsensusEstimate `json:"revenue"`
}

// AnalystConsensus holds analyst consensus estimates and AI analysis.
// report enum: "q1" (Q1), "qf" (quarterly), "saf" (semi-annual), "3q" (Q3), "af" (annual)
type AnalystConsensus struct {
	AISummary  string            `json:"ai_summary"`
	AIChatData AIChatData        `json:"aichat_data"`
	Currency   string            `json:"currency"`
	Report     string            `json:"report"`
	List       []USConsensusItem `json:"list"`
	OptReports []string          `json:"opt_reports"`
	H5Data     interface{}       `json:"h5_data"`
}

// FiscalYearDividend holds one fiscal-year row in ETFDividendInfo.
type FiscalYearDividend struct {
	Dividend        string `json:"dividend"`
	DividendYield   string `json:"dividend_yield"`
	FiscalYear      string `json:"fiscal_year"`
	Currency        string `json:"currency"`
	FiscalYearRange string `json:"fiscal_year_range"`
}

// ETFDividendInfo holds ETF dividend history.
type ETFDividendInfo struct {
	DividendTTM      string               `json:"dividend_ttm"`
	DividendYieldTTM string               `json:"dividend_yield_ttm"`
	DividendFreq     string               `json:"dividend_frequency"`
	Currency         string               `json:"currency"`
	FiscalYearInfo   []FiscalYearDividend `json:"fiscal_year_info"`
}

// USDividendItem is a single dividend payment record.
type USDividendItem struct {
	Dividend     string `json:"dividend"`
	DividendType string `json:"dividend_type"`
	ExDate       string `json:"ex_date"`
	PaymentDate  string `json:"payment_date"`
	RecordDate   string `json:"record_date"`
}

// USRecentDividend holds the trailing-12-month dividend summary.
type USRecentDividend struct {
	DividendTTM      string `json:"dividend_ttm"`
	DividendYieldTTM string `json:"dividend_yield_ttm"`
	Payouts          string `json:"payouts"`
	Currency         string `json:"currency"`
}

// USDividendHistoryItem is one fiscal-year row in the dividend history or payout-ratio table.
type USDividendHistoryItem struct {
	FiscalYear                string `json:"fiscal_year"`
	FiscalYearRange           string `json:"fiscal_year_range"`
	TotalShareholderYield     string `json:"total_shareholder_yield"`
	Dividend                  string `json:"dividend"`
	DividendYield             string `json:"dividend_yield"`
	DividendGrowthRate        string `json:"dividend_growth_rate"`
	DividendPayoutRatio       string `json:"dividend_payout_ratio"`
	DividendToCashflowRatio   string `json:"dividend_to_cashflow_ratio"`
	NetBuyback                string `json:"net_buyback"`
	NetBuybackYield           string `json:"net_buyback_yield"`
	NetBuybackGrowthRate      string `json:"net_buyback_growth_rate"`
	NetBuybackPayoutRatio     string `json:"net_buyback_payout_ratio"`
	NetBuybackToCashflowRatio string `json:"net_buyback_to_cashflow_ratio"`
	Currency                  string `json:"currency"`
}

// USDividendPayoutRecord is one actual dividend payment event.
type USDividendPayoutRecord struct {
	Dividend      string `json:"dividend"`
	DividendType  string `json:"dividend_type"`
	Currency      string `json:"currency"`
	ExDate        string `json:"ex_date"`
	PaymentDate   string `json:"payment_date"`
	RecordDate    string `json:"record_date"`
	Title         string `json:"title"`
	StartTimeUnix string `json:"start_time_unix"`
}

// USCompanyDividends holds historical dividend data for a US stock.
type USCompanyDividends struct {
	RecentDividends       USRecentDividend         `json:"recent_dividends"`
	DividendHistory       []USDividendHistoryItem  `json:"dividend_history"`
	PayoutRatios          []USDividendHistoryItem  `json:"payout_ratios"`
	DividendPayoutHistory []USDividendPayoutRecord `json:"dividend_payout_history"`
}

// ETFFile is a single document in the ETF file list.
type ETFFile struct {
	FileName   string `json:"file_name"`
	FilePath   string `json:"file_path"`
	UpdateDate string `json:"update_date"`
	Code       string `json:"code"`
	Format     string `json:"format"`
}

// ETFFilesResponse holds the ETF document list.
type ETFFilesResponse struct {
	Files []ETFFile `json:"files"`
}
