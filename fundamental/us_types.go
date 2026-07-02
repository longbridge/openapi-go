package fundamental

// RankTag is an industry-rank label returned by CompanyOverview.
type RankTag struct {
	Name     string `json:"name"`
	Chg      string `json:"chg"`
	RankType string `json:"rank_type"`
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

// FinancialOverview is the US financial overview (revenue, net income, EPS, cash flow).
// The server defines the exact inner fields; callers receive the raw map.
type FinancialOverview map[string]interface{}

// FinancialPeriod holds one reporting period's financial data.
type FinancialPeriod struct {
	Date   string                 `json:"date"`
	Values map[string]interface{} `json:"values"`
}

// FinancialStatement is the US financial statement (IS/BS/CF).
type FinancialStatement struct {
	Revenue   string            `json:"revenue"`
	NetIncome string            `json:"net_income"`
	NetMargin string            `json:"net_margin"`
	Periods   []FinancialPeriod `json:"periods"`
	Currency  string            `json:"currency"`
}

// KeyFinancialMetrics holds per-period key ratios (ROE, margins, leverage).
// The server defines the exact inner fields; callers receive the raw map.
type KeyFinancialMetrics map[string]interface{}

// AnalystConsensus holds per-period EPS and revenue forecasts.
// The server defines the exact inner fields; callers receive the raw map.
type AnalystConsensus map[string]interface{}

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
