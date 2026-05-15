package portfolio

import "github.com/shopspring/decimal"

// ── exchange_rate ─────────────────────────────────────────────────

// ExchangeRates is the response for ExchangeRate.
type ExchangeRates struct {
	// List of exchange rates.
	Exchanges []ExchangeRate
}

// ExchangeRate is one currency exchange rate.
type ExchangeRate struct {
	// Average rate (base_currency / other_currency).
	AverageRate float64
	// Base currency, e.g. "USD".
	BaseCurrency string
	// Bid rate.
	BidRate float64
	// Offer rate.
	OfferRate float64
	// Other currency, e.g. "HKD".
	OtherCurrency string
}

// ── profit_analysis ───────────────────────────────────────────────

// ProfitAnalysis is the combined response for ProfitAnalysis.
// It merges the summary and per-security sublist.
type ProfitAnalysis struct {
	// Summary overview.
	Summary ProfitAnalysisSummary
	// Per-security breakdown.
	Sublist ProfitAnalysisSublist
}

// ProfitAnalysisSummary is the account-level P&L summary.
type ProfitAnalysisSummary struct {
	// Account currency.
	Currency string
	// Current total asset value.
	CurrentTotalAsset *decimal.Decimal
	// Query start date string.
	StartDate string
	// Query end date string.
	EndDate string
	// Start time (unix timestamp string).
	StartTime string
	// End time (unix timestamp string).
	EndTime string
	// Ending asset value.
	EndingAssetValue *decimal.Decimal
	// Initial asset value.
	InitialAssetValue *decimal.Decimal
	// Total invested amount.
	InvestAmount *decimal.Decimal
	// Whether any trades occurred.
	IsTraded bool
	// Total profit/loss.
	SumProfit *decimal.Decimal
	// Total profit/loss rate.
	SumProfitRate *decimal.Decimal
	// Per-asset-type breakdown.
	Profits ProfitSummaryBreakdown
}

// ProfitSummaryBreakdown is the P&L breakdown by asset type.
type ProfitSummaryBreakdown struct {
	// Stock P&L.
	Stock *decimal.Decimal
	// Fund P&L.
	Fund *decimal.Decimal
	// Crypto P&L.
	Crypto *decimal.Decimal
	// Money market fund P&L.
	Mmf *decimal.Decimal
	// Other P&L.
	Other *decimal.Decimal
	// Cumulative transaction amount.
	CumulativeTransactionAmount *decimal.Decimal
	// Total number of orders.
	TradeOrderNum string
	// Total number of traded securities.
	TradeStockNum string
	// IPO P&L.
	Ipo *decimal.Decimal
	// IPO hits.
	IpoHit int32
	// IPO subscriptions.
	IpoSubscription int32
	// Per-category summary info.
	SummaryInfo []ProfitSummaryInfo
}

// ProfitSummaryInfo holds P&L info for one asset category.
type ProfitSummaryInfo struct {
	// Asset type.
	AssetType AssetType
	// Security with the maximum profit.
	ProfitMax string
	// Name of the max-profit security.
	ProfitMaxName string
	// Security with the maximum loss.
	LossMax string
	// Name of the max-loss security.
	LossMaxName string
}

// ProfitAnalysisSublist is the per-security P&L breakdown.
type ProfitAnalysisSublist struct {
	// Start time (unix timestamp string).
	Start string
	// End time (unix timestamp string).
	End string
	// Start date string.
	StartDate string
	// End date string.
	EndDate string
	// Last updated time (unix timestamp string).
	UpdatedAt string
	// Last updated date string.
	UpdatedDate string
	// Per-security items.
	Items []ProfitAnalysisItem
}

// ProfitAnalysisItem is the P&L for one security.
type ProfitAnalysisItem struct {
	// Security name.
	Name string
	// Market.
	Market string
	// Whether still holding.
	IsHolding bool
	// Profit/loss amount.
	Profit *decimal.Decimal
	// Profit/loss rate.
	ProfitRate *decimal.Decimal
	// Number of completed trades.
	ClearanceTimes int64
	// Asset type.
	ItemType AssetType
	// Currency.
	Currency string
	// Security symbol (converted from counter_id).
	Symbol string
	// Holding period display string.
	HoldingPeriod string
	// Ticker code.
	SecurityCode string
	// ISIN (for funds).
	Isin string
	// Underlying stock P&L.
	UnderlyingProfit *decimal.Decimal
	// Derivatives P&L.
	DerivativesProfit *decimal.Decimal
	// P&L in order currency.
	OrderProfit *decimal.Decimal
}

// ── profit_analysis_detail ────────────────────────────────────────

// ProfitAnalysisDetail is the response for ProfitAnalysisDetail.
type ProfitAnalysisDetail struct {
	// Total profit/loss.
	Profit *decimal.Decimal
	// Underlying stock P&L details.
	UnderlyingDetails ProfitDetails
	// Derivative P&L details.
	DerivativePnlDetails ProfitDetails
	// Security name.
	Name string
	// Last updated time (unix timestamp string).
	UpdatedAt string
	// Last updated date string.
	UpdatedDate string
	// Currency.
	Currency string
	// Default detail tab: 0 = underlying, 1 = derivative.
	DefaultTag int32
	// Query start time (unix timestamp string).
	Start string
	// Query end time (unix timestamp string).
	End string
	// Query start date string.
	StartDate string
	// Query end date string.
	EndDate string
}

// ProfitDetails is the detailed P&L breakdown for one asset class.
type ProfitDetails struct {
	// Current holding market value.
	HoldingValue *decimal.Decimal
	// Total profit/loss.
	Profit *decimal.Decimal
	// Cumulative credited amount.
	CumulativeCreditedAmount *decimal.Decimal
	// Credit detail entries.
	CreditedDetails []ProfitDetailEntry
	// Cumulative debited amount.
	CumulativeDebitedAmount *decimal.Decimal
	// Debit detail entries.
	DebitedDetails []ProfitDetailEntry
	// Cumulative fee amount.
	CumulativeFeeAmount *decimal.Decimal
	// Fee detail entries.
	FeeDetails []ProfitDetailEntry
	// Short position holding value.
	ShortHoldingValue *decimal.Decimal
	// Long position holding value.
	LongHoldingValue *decimal.Decimal
	// Opening position market value at period start.
	HoldingValueAtBeginning *decimal.Decimal
	// Closing position market value at period end.
	HoldingValueAtEnding *decimal.Decimal
}

// ProfitDetailEntry is one P&L detail line item (credit, debit, or fee).
type ProfitDetailEntry struct {
	// Description.
	Describe string
	// Amount.
	Amount *decimal.Decimal
}

// ── profit_analysis_by_market ─────────────────────────────────────

// ProfitAnalysisByMarket is the response for ProfitAnalysisByMarket.
type ProfitAnalysisByMarket struct {
	// Total P&L across all returned items.
	Profit *decimal.Decimal
	// Whether more pages are available.
	HasMore bool
	// Per-security P&L items for the requested market/page.
	StockItems []ProfitAnalysisByMarketItem
}

// ProfitAnalysisByMarketItem is one security entry in a by-market P&L response.
type ProfitAnalysisByMarketItem struct {
	// Security symbol (ticker code).
	Code string
	// Security name.
	Name string
	// Market, e.g. "HK", "US".
	Market string
	// Profit/loss amount.
	Profit *decimal.Decimal
}

// ── profit_analysis_flows ─────────────────────────────────────────

// ProfitAnalysisFlows is the response for ProfitAnalysisFlows.
type ProfitAnalysisFlows struct {
	// Paginated list of flow items.
	FlowsList []FlowItem
	// Whether there are more pages.
	HasMore bool
}

// FlowItem is one profit-analysis flow record.
type FlowItem struct {
	// Execution date string, e.g. "2024-01-15".
	ExecutedDate string
	// Execution timestamp string.
	ExecutedTimestamp string
	// Security code / ticker.
	Code string
	// Direction of the flow.
	Direction FlowDirection
	// Executed quantity.
	ExecutedQuantity *decimal.Decimal
	// Executed price.
	ExecutedPrice *decimal.Decimal
	// Executed cost.
	ExecutedCost *decimal.Decimal
	// Human-readable description.
	Describe string
}

// ── enums ─────────────────────────────────────────────────────────

// FlowDirection represents the direction of a portfolio flow.
type FlowDirection int

const (
	// FlowDirectionUnknown is an unknown flow direction.
	FlowDirectionUnknown FlowDirection = iota
	// FlowDirectionBuy is a buy flow.
	FlowDirectionBuy
	// FlowDirectionSell is a sell flow.
	FlowDirectionSell
)

// flowDirectionFromString converts a wire string to FlowDirection.
func flowDirectionFromString(s string) FlowDirection {
	switch s {
	case "buy":
		return FlowDirectionBuy
	case "sell":
		return FlowDirectionSell
	default:
		return FlowDirectionUnknown
	}
}

// AssetType represents the type of a portfolio asset.
type AssetType int

const (
	// AssetTypeUnknown is an unknown asset type.
	AssetTypeUnknown AssetType = iota
	// AssetTypeStock is a stock.
	AssetTypeStock
	// AssetTypeFund is a fund.
	AssetTypeFund
	// AssetTypeCrypto is a crypto asset.
	AssetTypeCrypto
)

// assetTypeFromString converts a wire string to AssetType.
func assetTypeFromString(s string) AssetType {
	switch s {
	case "stock":
		return AssetTypeStock
	case "fund":
		return AssetTypeFund
	case "crypto":
		return AssetTypeCrypto
	default:
		return AssetTypeUnknown
	}
}
