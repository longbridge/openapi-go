package jsontypes

type ExchangeRates struct {
	Exchanges []*ExchangeRate `json:"exchanges"`
}

type ExchangeRate struct {
	AverageRate   float64 `json:"average_rate"`
	BaseCurrency  string  `json:"base_currency"`
	BidRate       float64 `json:"bid_rate"`
	OfferRate     float64 `json:"offer_rate"`
	OtherCurrency string  `json:"other_currency"`
}

type ProfitAnalysisSummary struct {
	Currency          string                 `json:"currency"`
	CurrentTotalAsset string                 `json:"current_total_asset"`
	StartDate         string                 `json:"start_date"`
	EndDate           string                 `json:"end_date"`
	StartTime         string                 `json:"start_time"`
	EndTime           string                 `json:"end_time"`
	EndingAssetValue  string                 `json:"ending_asset_value"`
	InitialAssetValue string                 `json:"initial_asset_value"`
	InvestAmount      string                 `json:"invest_amount"`
	IsTraded          bool                   `json:"is_traded"`
	SumProfit         string                 `json:"sum_profit"`
	SumProfitRate     string                 `json:"sum_profit_rate"`
	Profits           *ProfitSummaryBreakdown `json:"profits"`
}

type ProfitSummaryBreakdown struct {
	Stock                      string               `json:"stock"`
	Fund                       string               `json:"fund"`
	Crypto                     string               `json:"crypto"`
	Mmf                        string               `json:"mmf"`
	Other                      string               `json:"other"`
	CumulativeTransactionAmount string              `json:"cumulative_transaction_amount"`
	TradeOrderNum              string               `json:"trade_order_num"`
	TradeStockNum              string               `json:"trade_stock_num"`
	Ipo                        string               `json:"ipo"`
	IpoHit                     int32                `json:"ipo_hit"`
	IpoSubscription            int32                `json:"ipo_subscription"`
	SummaryInfo                []*ProfitSummaryInfo `json:"summary_info"`
}

type ProfitSummaryInfo struct {
	AssetType      string `json:"asset_type"`
	ProfitMax      string `json:"profit_max"`
	ProfitMaxName  string `json:"profit_max_name"`
	LossMax        string `json:"loss_max"`
	LossMaxName    string `json:"loss_max_name"`
}

type ProfitAnalysisSublist struct {
	Start       string                 `json:"start"`
	End         string                 `json:"end"`
	StartDate   string                 `json:"start_date"`
	EndDate     string                 `json:"end_date"`
	UpdatedAt   string                 `json:"updated_at"`
	UpdatedDate string                 `json:"updated_date"`
	Items       []*ProfitAnalysisItem  `json:"items"`
}

type ProfitAnalysisItem struct {
	Name              string `json:"name"`
	Market            string `json:"market"`
	IsHolding         bool   `json:"is_holding"`
	Profit            string `json:"profit"`
	ProfitRate        string `json:"profit_rate"`
	ClearanceTimes    int64  `json:"clearance_times"`
	ItemType          string `json:"type"`
	Currency          string `json:"currency"`
	CounterId         string `json:"counter_id"`
	HoldingPeriod     string `json:"holding_period"`
	SecurityCode      string `json:"security_code"`
	Isin              string `json:"isin"`
	UnderlyingProfit  string `json:"underlying_profit"`
	DerivativesProfit string `json:"derivatives_profit"`
	OrderProfit       string `json:"order_profit"`
}

type ProfitAnalysisDetail struct {
	Profit                string        `json:"profit"`
	UnderlyingDetails     *ProfitDetails `json:"underlying_details"`
	DerivativePnlDetails  *ProfitDetails `json:"derivative_pnl_details"`
	Name                  string        `json:"name"`
	UpdatedAt             string        `json:"updated_at"`
	UpdatedDate           string        `json:"updated_date"`
	Currency              string        `json:"currency"`
	DefaultTag            int32         `json:"default_tag"`
	Start                 string        `json:"start"`
	End                   string        `json:"end"`
	StartDate             string        `json:"start_date"`
	EndDate               string        `json:"end_date"`
}

type ProfitDetails struct {
	HoldingValue                string               `json:"holding_value"`
	Profit                      string               `json:"profit"`
	CumulativeCreditedAmount    string               `json:"cumulative_credited_amount"`
	CreditedDetails             []*ProfitDetailEntry `json:"credited_details"`
	CumulativeDebitedAmount     string               `json:"cumulative_debited_amount"`
	DebitedDetails              []*ProfitDetailEntry `json:"debited_details"`
	CumulativeFeeAmount         string               `json:"cumulative_fee_amount"`
	FeeDetails                  []*ProfitDetailEntry `json:"fee_details"`
	ShortHoldingValue           string               `json:"short_holding_value"`
	LongHoldingValue            string               `json:"long_holding_value"`
	HoldingValueAtBeginning     string               `json:"holding_value_at_beginning"`
	HoldingValueAtEnding        string               `json:"holding_value_at_ending"`
}

type ProfitDetailEntry struct {
	Describe string `json:"describe"`
	Amount   string `json:"amount"`
}

type ProfitAnalysisByMarket struct {
	Profit     string                      `json:"profit"`
	HasMore    bool                        `json:"has_more"`
	StockItems []*ProfitAnalysisByMarketItem `json:"stock_items"`
}

type ProfitAnalysisByMarketItem struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Market string `json:"market"`
	Profit string `json:"profit"`
}

type ProfitAnalysisFlows struct {
	FlowsList []*FlowItem `json:"flows_list"`
	HasMore   bool        `json:"has_more"`
}

type FlowItem struct {
	ExecutedDate      string  `json:"executed_date"`
	ExecutedTimestamp string  `json:"executed_timestamp"`
	Code              string  `json:"code"`
	Direction         string  `json:"direction"`
	ExecutedQuantity  string  `json:"executed_quantity"`
	ExecutedPrice     string  `json:"executed_price"`
	ExecutedCost      string  `json:"executed_cost"`
	Describe          string  `json:"describe"`
}
