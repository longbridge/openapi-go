package fundamental

// RankTag is an industry-rank label returned by CompanyOverview.
type RankTag struct {
	Name     string `json:"name"`
	Chg      string `json:"chg"`
	RankType int32 `json:"rank_type"`
}

// USCompanyOverview is the US company summary snapshot.
type USCompanyOverview struct {
	Intro       string    `json:"intro"`
	MarketCap   string    `json:"market_cap"`
	CcySymbol   string    `json:"ccy_symbol"`
	TopRankTags []RankTag `json:"top_rank_tags"`
	DetailURL   string    `json:"detail_url"`
}

// ValuationIndicator holds the current valuation metric detail.
type ValuationIndicator struct {
	Circle     string `json:"circle"`
	Part       string `json:"part"`
	Metric     string `json:"metric"`
	MetricType string `json:"metric_type"`
	Desc       string `json:"desc"`
	CcySymbol  string `json:"ccy_symbol"`
}

// ValuationOverview is the US valuation snapshot.
type ValuationOverview struct {
	Indicator        string             `json:"indicator"`
	CurrentIndicator ValuationIndicator `json:"current_indicator"`
	Range            int32              `json:"range"`
	Date             string             `json:"date"`
	AISummary        string             `json:"ai_summary"`
}

// USReportPeriod identifies a reporting period (quarter/half/annual).
type USReportPeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	ReportTxt string `json:"report_txt"`
}

// FinancialISItem is one income-statement entry in FinancialOverview.
type FinancialISItem struct {
	Revenue   string                `json:"revenue"`
	NetIncome string                `json:"net_income"`
	NetMargin string                `json:"net_margin"`
	Report    USReportPeriod `json:"report"`
}

// FinancialBSItem is one balance-sheet entry in FinancialOverview.
type FinancialBSItem struct {
	DebtAssetsRatio  string                `json:"debt_assets_ratio"`
	TotalAssets      string                `json:"total_assets"`
	TotalLiabilities string                `json:"total_liabilities"`
	Report           USReportPeriod `json:"report"`
}

// FinancialCFItem is one cash-flow entry in FinancialOverview.
type FinancialCFItem struct {
	Operating string                `json:"operating"`
	Investing string                `json:"investing"`
	Financing string                `json:"financing"`
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

// FinancialStatement is the US financial statement (IS/BS/CF).
type FinancialStatement struct {
	Revenue   string `json:"revenue"`
	NetIncome string `json:"net_income"`
	NetMargin string `json:"net_margin"`
	Periods   []struct {
		Date   string      `json:"date"`
		Values interface{} `json:"values"`
	} `json:"periods"`
	Currency string `json:"currency"`
}

// KeyMetricItem is one period entry in KeyFinancialMetrics.
type KeyMetricItem struct {
	FfPeriod  string      `json:"ff_period"`
	FfYear    int32       `json:"ff_year"`
	FpEnd     string      `json:"fp_end"`
	ReportTxt string      `json:"report_txt"`
	RptDate   string      `json:"rpt_date"`
	Fields    interface{} `json:"fields"` // metric values; shape varies per field set
}

// KeyFinancialMetrics holds per-period key ratios (ROE, margins, debt ratio).
type KeyFinancialMetrics struct {
	Currency    string          `json:"currency"`
	Report      string          `json:"report"`
	EmptyFields []string        `json:"empty_fields"`
	List        []KeyMetricItem `json:"list"`
}

// AIChatData holds the AI chat context embedded in AnalystConsensus.
type AIChatData struct {
	AgentID        string `json:"agent_id"`
	HandoffAgentID string `json:"handoff_agent_id"`
	Symbol         string `json:"symbol"`
	Text           string `json:"text"`
	Type           string `json:"type"`
	WorkflowType   string `json:"workflow_type"`
}

// AnalystConsensus holds analyst consensus estimates and AI analysis.
type AnalystConsensus struct {
	AISummary  string      `json:"ai_summary"`
	AIChatData AIChatData  `json:"aichat_data"`
	Currency   string      `json:"currency"`
	Report     string      `json:"report"`
	List       interface{} `json:"list"`       // consensus detail; shape TBD from production
	OptReports interface{} `json:"opt_reports"` // option consensus; shape TBD from production
	H5Data     interface{} `json:"h5_data"`
}

// FiscalYearDividend holds dividend records for one fiscal year.
type FiscalYearDividend struct {
	Year          string                   `json:"year"`
	TotalDividend string                   `json:"total_dividend"`
	Records       []map[string]interface{} `json:"records"`
}

// ETFDividendInfo holds ETF dividend history.
type ETFDividendInfo struct {
	DividendTTM      string               `json:"dividend_ttm"`
	DividendYieldTTM string               `json:"dividend_yield_ttm"`
	DividendFreq     string               `json:"dividend_frequency"`
	Currency         string               `json:"currency"`
	FiscalYearInfo   []FiscalYearDividend `json:"fiscal_year_info"`
}

// USDividendItem is a single dividend payment record for a US stock.
type USDividendItem struct {
	Dividend     string `json:"dividend"`
	DividendType string `json:"dividend_type"`
	ExDate       string `json:"ex_date"`
	PaymentDate  string `json:"payment_date"`
	RecordDate   string `json:"record_date"`
}

// USCompanyDividends holds historical dividend data for a US stock.
type USCompanyDividends struct {
	DividendTTM      string           `json:"dividend_ttm"`
	DividendYieldTTM string           `json:"dividend_yield_ttm"`
	Payouts          string           `json:"payouts"`
	Currency         string           `json:"currency"`
	Items            []USDividendItem `json:"items"`
}

// ETFFile is a single document in the ETF file list.
type ETFFile struct {
	Name     string `json:"name"`
	FileType string `json:"file_type"`
	URL      string `json:"url"`
}

// ETFFilesResponse holds the ETF document list.
type ETFFilesResponse struct {
	Files []ETFFile `json:"files"`
}
