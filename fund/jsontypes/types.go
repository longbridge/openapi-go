// Package jsontypes contains the raw JSON wire types for the fund (mutual
// fund) API. These types match the exact JSON field names returned by the
// Longbridge API. Use the parent fund package for the idiomatic Go types.
//
// Numeric fields that arrive as strings on the wire are kept as string to
// preserve the exact server formatting; unix-second timestamps are kept as
// int64. Server-defined "any" structures are kept as json.RawMessage.
package jsontypes

import "encoding/json"

// ----- request bodies -----

// GetFundsBody is the request body for the fund list (POST /v1/fund/funds).
type GetFundsBody struct {
	Filter       json.RawMessage `json:"filter,omitempty"`
	QuickIds     []int64         `json:"quick_ids,omitempty"`
	TimeInterval []string        `json:"time_interval,omitempty"`
}

// ValidateFundOrderBody is the request body for validating a fund order.
type ValidateFundOrderBody struct {
	Symbol         string `json:"symbol"`
	Action         string `json:"action"`
	Currency       string `json:"currency"`
	Amount         string `json:"amount,omitempty"`
	Units          string `json:"units,omitempty"`
	DividendOption *int32 `json:"dividend_option,omitempty"`
	FundSource     *int32 `json:"fund_source,omitempty"`
	AccountChannel string `json:"account_channel,omitempty"`
}

// SubmitFundOrderBody is the request body for submitting a fund order.
type SubmitFundOrderBody struct {
	Symbol         string `json:"symbol"`
	Action         string `json:"action"`
	Currency       string `json:"currency"`
	Amount         string `json:"amount,omitempty"`
	Units          string `json:"units,omitempty"`
	DividendOption *int32 `json:"dividend_option,omitempty"`
	Fee            string `json:"fee,omitempty"`
	IsSellAll      *bool  `json:"is_sell_all,omitempty"`
	Remark         string `json:"remark,omitempty"`
	TradeMethod    *int32 `json:"trade_method,omitempty"`
}

// CancelFundOrderBody is the request body for cancelling a fund order.
type CancelFundOrderBody struct {
	UtId int64 `json:"ut_id"`
}

// ----- response envelopes -----

// HotFundsResponse wraps the hot-fund list.
type HotFundsResponse struct {
	List []*HotFund `json:"list"`
}

// FundsResponse wraps the fund list.
type FundsResponse struct {
	Funds []*FundBrief `json:"funds"`
}

// AnnualReturnsResponse wraps the annual-return list.
type AnnualReturnsResponse struct {
	List []*FundAnnualReturn `json:"list"`
}

// QuarterlyReturnsResponse wraps the quarterly-return list.
type QuarterlyReturnsResponse struct {
	List []*FundQuarterlyReturn `json:"list"`
}

// PerformanceResponse wraps the fund performance list.
type PerformanceResponse struct {
	Value []*FundPerformance `json:"value"`
}

// NavResponse wraps the fund latest net-value list.
type NavResponse struct {
	Value []*FundNavValue `json:"value"`
}

// NavHistoryResponse wraps the fund historical net-value list.
type NavHistoryResponse struct {
	HistoryValue []*FundNavValue `json:"history_value"`
}

// StockHoldingsResponse wraps the fund reverse stock-holdings list.
type StockHoldingsResponse struct {
	Lists []*FundStockHolding `json:"lists"`
}

// PositionPerformanceResponse wraps the held-fund performance list.
type PositionPerformanceResponse struct {
	Value []*FundPositionPerformance `json:"value"`
}

// PositionNavResponse wraps the held-fund net-value list.
type PositionNavResponse struct {
	HistoryValue []*FundPositionNav `json:"history_value"`
}

// OrdersResponse wraps the fund orders list.
type OrdersResponse struct {
	Orders []*FundOrder `json:"orders"`
}

// TransactionsResponse wraps the fund transactions list.
type TransactionsResponse struct {
	List []*FundTransaction `json:"list"`
}

// ----- response models -----

// FundNavValue is a fund net-asset-value data point (latest / historical).
type FundNavValue struct {
	Change              string `json:"change"`
	ChangePercent       string `json:"change_percent"`
	ChangePercentFormat string `json:"change_percent_format"`
	CounterId           string `json:"counter_id"`
	CounterName         string `json:"counter_name"`
	Currency            string `json:"currency"`
	DateFormat          string `json:"date_format"`
	Isin                string `json:"isin"`
	LastUpdateTime      int64  `json:"last_update_time"`
	Value               string `json:"value"`
	ValueFormat         string `json:"value_format"`
}

// FundPerformancePoint is a recent performance point used by the hot-fund list.
type FundPerformancePoint struct {
	Date     string `json:"date"`
	LastDone string `json:"last_done"`
}

// HotFund is a hot-selling fund entry.
type HotFund struct {
	AssetClass        int32                   `json:"asset_class"`
	AssetClassName    string                  `json:"asset_class_name"`
	CounterId         string                  `json:"counter_id"`
	Currency          string                  `json:"currency"`
	EarningRate       string                  `json:"earning_rate"`
	FundPerformances  []*FundPerformancePoint `json:"fund_performances"`
	Name              string                  `json:"name"`
	PurchaseAmount    string                  `json:"purchase_amount"`
	RecommendationTxt string                  `json:"recommendation_text"`
	RiskLevel         int32                   `json:"risk_level"`
	RiskLevelName     string                  `json:"risk_level_name"`
	TimeInterval      string                  `json:"time_interval"`
}

// FundBrief is a fund entry in the fund list.
type FundBrief struct {
	AssetClass        int32  `json:"asset_class"`
	AssetClassName    string `json:"asset_class_name"`
	Code              string `json:"code"`
	CounterId         string `json:"counter_id"`
	Currency          string `json:"currency"`
	Description       string `json:"description"`
	EarningRate       string `json:"earning_rate"`
	Holding           bool   `json:"holding"`
	Isin              string `json:"isin"`
	Name              string `json:"name"`
	Product           string `json:"product"`
	PurchaseAmount    string `json:"purchase_amount"`
	RecommendationTxt string `json:"recommendation_text"`
	RiskLevel         int32  `json:"risk_level"`
	RiskLevelName     string `json:"risk_level_name"`
	TimeInterval      string `json:"time_interval"`
	UnitValue         string `json:"unit_value"`
}

// FundFilters is the fund list filter options.
type FundFilters struct {
	AssetClass           []json.RawMessage `json:"asset_class"`
	Company              []json.RawMessage `json:"company"`
	Currency             []json.RawMessage `json:"currency"`
	IndustryCategoryName []json.RawMessage `json:"industry_category_name"`
	RiskLevel            []json.RawMessage `json:"risk_level"`
}

// FundAssetAllocationItem is a single holding inside a fund's asset allocation.
type FundAssetAllocationItem struct {
	Code          string `json:"code"`
	CounterId     string `json:"counter_id"`
	Name          string `json:"name"`
	PositionRatio string `json:"position_ratio"`
}

// FundAssetAllocation is a fund's asset allocation.
type FundAssetAllocation struct {
	AssetType  int32                      `json:"asset_type"`
	Lists      []*FundAssetAllocationItem `json:"lists"`
	ReportDate string                     `json:"report_date"`
}

// FundDetail is a fund detail.
type FundDetail struct {
	AdditionalPurchaseAmount   string              `json:"additional_purchase_amount"`
	AffirmDay                  int32               `json:"affirm_day"`
	AmountAffirmDay            string              `json:"amount_affirm_day"`
	AssetAllocation            FundAssetAllocation `json:"asset_allocation"`
	AssetClass                 int32               `json:"asset_class"`
	AssetClassName             string              `json:"asset_class_name"`
	BillPurchaseRate           string              `json:"bill_purchase_rate"`
	Channel                    string              `json:"channel"`
	ClosePeriod                string              `json:"close_period"`
	Code                       string              `json:"code"`
	Currency                   string              `json:"currency"`
	CutOffTime                 string              `json:"cut_off_time"`
	Derivatives                bool                `json:"derivatives"`
	DoneDay                    int32               `json:"done_day"`
	ExcessReturnFee            string              `json:"excess_return_fee"`
	GstRate                    string              `json:"gst_rate"`
	Introduce                  string              `json:"introduce"`
	IsCashPlus                 bool                `json:"is_cash_plus"`
	IsComplex                  bool                `json:"is_complex"`
	IsNewCashPlus              bool                `json:"is_new_cash_plus"`
	IsYinghebao                bool                `json:"is_yinghebao"`
	Isin                       string              `json:"isin"`
	ManageRate                 string              `json:"manage_rate"`
	Manager                    string              `json:"manager"`
	MinHoldCash                string              `json:"min_hold_cash"`
	MinHoldShare               string              `json:"min_hold_share"`
	MinSellShare               string              `json:"min_sell_share"`
	MonthRaiseDay              string              `json:"month_raise_day"`
	Name                       string              `json:"name"`
	NavDeadline                string              `json:"nav_deadline"`
	NoLoad                     bool                `json:"no_load"`
	OpenDate                   string              `json:"open_date"`
	OpenPeriod                 string              `json:"open_period"`
	Product                    string              `json:"product"`
	ProductInformationLocals   string              `json:"product_information_locals"`
	Profile                    string              `json:"profile"`
	Purchasable                int32               `json:"purchasable"`
	PurchaseAffirmDay          string              `json:"purchase_affirm_day"`
	PurchaseAmount             string              `json:"purchase_amount"`
	PurchaseRate               string              `json:"purchase_rate"`
	Rating                     int32               `json:"rating"`
	Redeemable                 int32               `json:"redeemable"`
	RedemptionAdvanceDay       string              `json:"redemption_advance_day"`
	RedemptionAmount           string              `json:"redemption_amount"`
	RedemptionClosePeriodShows string              `json:"redemption_close_period_shows"`
	RedemptionDoneDay          string              `json:"redemption_done_day"`
	RedemptionOpenDayShows     string              `json:"redemption_open_day_shows"`
	RiskLevel                  int32               `json:"risk_level"`
	RiskLevelName              string              `json:"risk_level_name"`
	VerifyStatus               int32               `json:"verify_status"`
	VirtualCurrency            bool                `json:"virtual_currency"`
	YearToDateYield            string              `json:"year_to_date_yield"`
	YtdYieldType               int32               `json:"ytd_yield_type"`
}

// FundAnalysis is a fund analysis (level 1).
type FundAnalysis struct {
	ActualPeriod  int32           `json:"actual_period"`
	CostLevel     json.RawMessage `json:"cost_level"`
	ReturnAbility json.RawMessage `json:"return_ability"`
	RiskAbility   json.RawMessage `json:"risk_ability"`
	UpdatedAt     string          `json:"updated_at"`
	ValueForMoney json.RawMessage `json:"value_for_money"`
	Visible       bool            `json:"visible"`
}

// FundAnalysisDetail is a fund analysis detail (level 2).
type FundAnalysisDetail struct {
	ActualPeriod     int32           `json:"actual_period"`
	AvailablePeriods []int32         `json:"available_periods"`
	CostLevel        json.RawMessage `json:"cost_level"`
	ReturnAbility    json.RawMessage `json:"return_ability"`
	RiskAbility      json.RawMessage `json:"risk_ability"`
	UpdatedAt        string          `json:"updated_at"`
	ValueForMoney    json.RawMessage `json:"value_for_money"`
	Visible          bool            `json:"visible"`
}

// FundTrendContrast is a benchmark contrast series in a fund trend chart.
type FundTrendContrast struct {
	BenchmarkName string            `json:"benchmark_name"`
	Performances  []json.RawMessage `json:"performances"`
}

// FundTrend is a fund trend chart.
type FundTrend struct {
	ActualPeriod                int32             `json:"actual_period"`
	AvailablePeriods            []int32           `json:"available_periods"`
	CategoryAveragePerformances []json.RawMessage `json:"category_average_performances"`
	ContrastPerformances        FundTrendContrast `json:"contrast_performances"`
	FundPerformances            []json.RawMessage `json:"fund_performances"`
}

// FundNamedContrast is a named contrast performance series.
type FundNamedContrast struct {
	Name         string            `json:"name"`
	Performances []json.RawMessage `json:"performances"`
}

// FundPerformanceComparison is a fund performance comparison.
type FundPerformanceComparison struct {
	ContrastPerformances []*FundNamedContrast `json:"contrast_performances"`
	FundPerformances     []json.RawMessage    `json:"fund_performances"`
}

// FundAnnualReturn is a fund annual return entry.
type FundAnnualReturn struct {
	ChangePercent string `json:"change_percent"`
	Year          int32  `json:"year"`
}

// FundQuarterlyReturn is a fund quarterly return entry.
type FundQuarterlyReturn struct {
	ChangePercent string `json:"change_percent"`
	Quarter       int32  `json:"quarter"`
	Year          int32  `json:"year"`
}

// FundPerformance is a fund's detailed performance figures.
type FundPerformance struct {
	AnnualizedReturnFive         string `json:"annualized_return_five"`
	AnnualizedReturnOne          string `json:"annualized_return_one"`
	AnnualizedReturnTen          string `json:"annualized_return_ten"`
	AnnualizedReturnThree        string `json:"annualized_return_three"`
	AnnualizedReturnTwo          string `json:"annualized_return_two"`
	CounterId                    string `json:"counter_id"`
	FundName                     string `json:"fund_name"`
	PerformanceRankFiveYears     int32  `json:"performance_rank_five_years"`
	PerformanceRankOneDay        int32  `json:"performance_rank_one_day"`
	PerformanceRankOneMonth      int32  `json:"performance_rank_one_month"`
	PerformanceRankOneWeek       int32  `json:"performance_rank_one_week"`
	PerformanceRankOneYear       int32  `json:"performance_rank_one_year"`
	PerformanceRankSixMonths     int32  `json:"performance_rank_six_months"`
	PerformanceRankTenYears      int32  `json:"performance_rank_ten_years"`
	PerformanceRankThreeMonths   int32  `json:"performance_rank_three_months"`
	PerformanceRankThreeYears    int32  `json:"performance_rank_three_years"`
	PerformanceRankTwoYears      int32  `json:"performance_rank_two_years"`
	PerformanceRankYtd           int32  `json:"performance_rank_ytd"`
	PerformanceReturnFiveYears   string `json:"performance_return_five_years"`
	PerformanceReturnOneDay      string `json:"performance_return_one_day"`
	PerformanceReturnOneMonth    string `json:"performance_return_one_month"`
	PerformanceReturnOneWeek     string `json:"performance_return_one_week"`
	PerformanceReturnOneYear     string `json:"performance_return_one_year"`
	PerformanceReturnSixMonths   string `json:"performance_return_six_months"`
	PerformanceReturnTenYears    string `json:"performance_return_ten_years"`
	PerformanceReturnThreeMonths string `json:"performance_return_three_months"`
	PerformanceReturnThreeYears  string `json:"performance_return_three_years"`
	PerformanceReturnTwoYears    string `json:"performance_return_two_years"`
	PerformanceReturnYtd         string `json:"performance_return_ytd"`
	PerformanceTotalFiveYears    int32  `json:"performance_total_five_years"`
	PerformanceTotalOneDay       int32  `json:"performance_total_one_day"`
	PerformanceTotalOneMonth     int32  `json:"performance_total_one_month"`
	PerformanceTotalOneWeek      int32  `json:"performance_total_one_week"`
	PerformanceTotalOneYear      int32  `json:"performance_total_one_year"`
	PerformanceTotalSixMonths    int32  `json:"performance_total_six_months"`
	PerformanceTotalTenYears     int32  `json:"performance_total_ten_years"`
	PerformanceTotalThreeMonths  int32  `json:"performance_total_three_months"`
	PerformanceTotalThreeYears   int32  `json:"performance_total_three_years"`
	PerformanceTotalTwoYears     int32  `json:"performance_total_two_years"`
	PerformanceTotalYtd          int32  `json:"performance_total_ytd"`
	SevenDaysAnnualized          string `json:"seven_days_annualized"`
	TenThousandPrice             string `json:"ten_thousand_price"`
	UpdateTime                   int64  `json:"update_time"`
}

// FundHolding is a single fund holding (top-10 holdings).
type FundHolding struct {
	BondType           string `json:"bond_type"`
	BondTypeName       string `json:"bond_type_name"`
	CountryName        string `json:"country_name"`
	HoldingType        string `json:"holding_type"`
	IndustryName       string `json:"industry_name"`
	MarketValue        string `json:"market_value"`
	MaturityDate       string `json:"maturity_date"`
	Name               string `json:"name"`
	ShareChange        string `json:"share_change"`
	ShareChangePercent string `json:"share_change_percent"`
	Shares             string `json:"shares"`
	Weighting          string `json:"weighting"`
}

// FundHoldings is a fund's top-10 holdings.
type FundHoldings struct {
	Holdings   []*FundHolding `json:"holdings"`
	ReportDate string         `json:"report_date"`
	Weighting  string         `json:"weighting"`
}

// FundStockHolding is a stock held by the fund (reverse lookup).
type FundStockHolding struct {
	Code          string `json:"code"`
	CounterId     string `json:"counter_id"`
	Currency      string `json:"currency"`
	Name          string `json:"name"`
	PositionRatio string `json:"position_ratio"`
	ReportDate    string `json:"report_date"`
}

// FundPosition is a single fund position held by the user.
type FundPosition struct {
	Amount           string `json:"amount"`
	CounterId        string `json:"counter_id"`
	Currency         string `json:"currency"`
	FreezeUnits      string `json:"freeze_units"`
	HoldingProfit    string `json:"holding_profit"`
	HoldingUnits     string `json:"holding_units"`
	Name             string `json:"name"`
	RecentProfit     string `json:"recent_profit"`
	RecentTradingDay int64  `json:"recent_trading_day"`
	SumRecentProfit  string `json:"sum_recent_profit"`
}

// FundPositions is the user's fund positions overview.
type FundPositions struct {
	AccountChannel          string          `json:"account_channel"`
	List                    []*FundPosition `json:"list"`
	PendingBuyOrders        string          `json:"pending_buy_orders"`
	RecentTradingDay        int64           `json:"recent_trading_day"`
	SoldPendingCreditOrders string          `json:"sold_pending_credit_orders"`
}

// FundDatedValue is a dated value point.
type FundDatedValue struct {
	Date  int64  `json:"date"`
	Value string `json:"value"`
}

// FundUnitValue is a fund unit-value point (position view).
type FundUnitValue struct {
	Date            int64  `json:"date"`
	DayIncreaseRate string `json:"day_increase_rate"`
	TotalValue      string `json:"total_value"`
	UnitValue       string `json:"unit_value"`
}

// FundPositionDetailValues holds the detail values of a single fund position.
type FundPositionDetailValues struct {
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	HoldingCost             string `json:"holding_cost"`
	HoldingProfit           string `json:"holding_profit"`
	HoldingProfitRate       string `json:"holding_profit_rate"`
	HoldingUnits            string `json:"holding_units"`
	HoldingValue            string `json:"holding_value"`
	PendingBuyValue         string `json:"pending_buy_value"`
	PendingSellValue        string `json:"pending_sell_value"`
	ProfitAmountAccumTd     string `json:"profit_amount_accum_td"`
	ProfitAmountAccumTdRate string `json:"profit_amount_accum_td_rate"`
	RecentProfit            string `json:"recent_profit"`
	RecentTradingday        int64  `json:"recent_tradingday"`
	RecentUnitValue         string `json:"recent_unit_value"`
	SoldPendingConfirmUnits string `json:"sold_pending_confirm_units"`
}

// FundPositionDetail is the detail of a single fund position.
type FundPositionDetail struct {
	DetailValues FundPositionDetailValues `json:"detail_values"`
	SumProfit    []*FundDatedValue        `json:"sum_profit"`
	UtValue      []*FundUnitValue         `json:"ut_value"`
}

// FundPositionPerformance holds the performance figures for a held fund.
type FundPositionPerformance struct {
	AnnualizedReturnFive         string `json:"annualized_return_five"`
	AnnualizedReturnOne          string `json:"annualized_return_one"`
	AnnualizedReturnTen          string `json:"annualized_return_ten"`
	AnnualizedReturnThree        string `json:"annualized_return_three"`
	AnnualizedReturnTwo          string `json:"annualized_return_two"`
	CounterId                    string `json:"counter_id"`
	FundName                     string `json:"fund_name"`
	PerformanceReturnFiveYears   string `json:"performance_return_five_years"`
	PerformanceReturnOneDay      string `json:"performance_return_one_day"`
	PerformanceReturnOneMonth    string `json:"performance_return_one_month"`
	PerformanceReturnOneWeek     string `json:"performance_return_one_week"`
	PerformanceReturnOneYear     string `json:"performance_return_one_year"`
	PerformanceReturnSixMonths   string `json:"performance_return_six_months"`
	PerformanceReturnTenYears    string `json:"performance_return_ten_years"`
	PerformanceReturnThreeMonths string `json:"performance_return_three_months"`
	PerformanceReturnThreeYears  string `json:"performance_return_three_years"`
	PerformanceReturnTwoYears    string `json:"performance_return_two_years"`
	PerformanceReturnYtd         string `json:"performance_return_ytd"`
	UpdateTime                   int64  `json:"update_time"`
}

// FundPositionProfits is the user's cumulative profit for a held fund.
type FundPositionProfits struct {
	Currency       string            `json:"currency"`
	HistoryValue   []*FundDatedValue `json:"history_value"`
	LastUpdateTime int64             `json:"last_update_time"`
	SumProfit      string            `json:"sum_profit"`
}

// FundPositionNav is a held-fund net-value point (position view).
type FundPositionNav struct {
	Change         string `json:"change"`
	ChangePercent  string `json:"change_percent"`
	CounterId      string `json:"counter_id"`
	CounterName    string `json:"counter_name"`
	LastUpdateTime int64  `json:"last_update_time"`
	Value          string `json:"value"`
}

// FundDividend is a cash dividend record for a held fund.
type FundDividend struct {
	Amount    string `json:"amount"`
	CounterId string `json:"counter_id"`
	Currency  string `json:"currency"`
	Date      int64  `json:"date"`
	DivMethod string `json:"div_method"`
	Name      string `json:"name"`
}

// FundDividends is the user's dividend records for a held fund.
type FundDividends struct {
	Currency     string          `json:"currency"`
	DivCashInfos []*FundDividend `json:"div_cash_infos"`
	LastestDate  int64           `json:"lastest_date"`
	TotalDivCash string          `json:"total_div_cash"`
}

// FundOrder is a fund order (list view).
type FundOrder struct {
	Action      string `json:"action"`
	Amount      string `json:"amount"`
	CounterId   string `json:"counter_id"`
	CreatedAt   int64  `json:"created_at"`
	Currency    string `json:"currency"`
	FundName    string `json:"fund_name"`
	Id          int64  `json:"id"`
	IsAuto      bool   `json:"is_auto"`
	NetWorth    string `json:"net_worth"`
	ProductType string `json:"product_type"`
	State       string `json:"state"`
	StateDesc   string `json:"state_desc"`
	Units       string `json:"units"`
}

// FundOrderKeyword is a keyword block in a fund order detail.
type FundOrderKeyword struct {
	Content      string `json:"content"`
	Group        string `json:"group"`
	Key          string `json:"key"`
	LineStrategy string `json:"line_strategy"`
	Title        string `json:"title"`
}

// FundOrderStage is a processing stage in a fund order detail.
type FundOrderStage struct {
	Desc     string `json:"desc"`
	Key      string `json:"key"`
	Link     string `json:"link"`
	LinkText string `json:"link_text"`
	Progress string `json:"progress"`
	Stage    string `json:"stage"`
}

// FundOrderInfo is the full information of a fund order.
type FundOrderInfo struct {
	Aaid           int64  `json:"aaid"`
	AccountChannel string `json:"account_channel"`
	Action         string `json:"action"`
	Amount         string `json:"amount"`
	Channel        string `json:"channel"`
	CounterId      string `json:"counter_id"`
	CreatedAt      int64  `json:"created_at"`
	Currency       string `json:"currency"`
	DividendOption string `json:"dividend_option"`
	EqAt           int64  `json:"eq_at"`
	Fee            string `json:"fee"`
	FundName       string `json:"fund_name"`
	FundSource     string `json:"fund_source"`
	Histories      string `json:"histories"`
	Id             int64  `json:"id"`
	Message        string `json:"message"`
	NetWorth       string `json:"net_worth"`
	PriceAt        int64  `json:"price_at"`
	ProcessedAt    int64  `json:"processed_at"`
	ProductType    string `json:"product_type"`
	Repurchaseable bool   `json:"repurchaseable"`
	SaleProceeds   string `json:"sale_proceeds"`
	SalesCharge    string `json:"sales_charge"`
	SalesPrice     string `json:"sales_price"`
	SalesUnit      string `json:"sales_unit"`
	State          string `json:"state"`
	StateDesc      string `json:"state_desc"`
	Status         int32  `json:"status"`
	StatusEx       int32  `json:"status_ex"`
	TDescription   string `json:"t_description"`
	TimePartition  string `json:"time_partition"`
	TotalAmount    string `json:"total_amount"`
	TransactionAt  int64  `json:"transaction_at"`
	Units          string `json:"units"`
	WithdrawAt     int64  `json:"withdraw_at"`
	Withdrawable   bool   `json:"withdrawable"`
}

// FundOrderDetail is a fund order detail.
type FundOrderDetail struct {
	Keywords []*FundOrderKeyword `json:"keywords"`
	Order    FundOrderInfo       `json:"order"`
	Stages   []*FundOrderStage   `json:"stages"`
}

// FundTransaction is a fund transaction / cash-flow record.
type FundTransaction struct {
	Amount              string `json:"amount"`
	Category            string `json:"category"`
	CreatedAt           int64  `json:"created_at"`
	Currency            string `json:"currency"`
	Description         string `json:"description"`
	DetailCreatedAt     int64  `json:"detail_created_at"`
	DetailType          string `json:"detail_type"`
	DoneAt              int64  `json:"done_at"`
	QuantityDescription string `json:"quantity_description"`
	RedirectPage        string `json:"redirect_page"`
	RedirectPageV2      string `json:"redirect_page_v2"`
	RefNo               string `json:"ref_no"`
	StockQuantity       string `json:"stock_quantity"`
	TxType              string `json:"tx_type"`
	TypeName            string `json:"type_name"`
}

// FundOrderValidation is the result of validating a fund order.
type FundOrderValidation struct {
	AuthToken     string `json:"auth_token"`
	EvalAddress   string `json:"eval_address"`
	FundRiskLevel int32  `json:"fund_risk_level"`
	Msg           string `json:"msg"`
	UserPi        int32  `json:"user_pi"`
	UserRiskLevel int32  `json:"user_risk_level"`
}

// FundOrderSubmitResponse is the result of submitting a fund order.
type FundOrderSubmitResponse struct {
	Action    string `json:"action"`
	Amount    string `json:"amount"`
	CounterId string `json:"counter_id"`
	CreatedAt int64  `json:"created_at"`
	FundName  string `json:"fund_name"`
	Id        int64  `json:"id"`
	Msg       string `json:"msg"`
	Status    int32  `json:"status"`
	Units     string `json:"units"`
}
