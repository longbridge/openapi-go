package fund

import "encoding/json"

// GetFunds is the request for the fund list (POST /v1/fund/funds).
type GetFunds struct {
	Filter       json.RawMessage // optional, server-defined filter object
	QuickIds     []int64         // optional, quick-filter ids
	TimeInterval []string        // optional, earning-rate time intervals
}

// GetFundAnalysis is the request for the fund analysis / trend / comparison
// endpoints.
type GetFundAnalysis struct {
	Period int32 // optional, analysis period (0 = omit)
}

// FundPage is a paging option (page / size).
type FundPage struct {
	Page int32 // optional
	Size int32 // optional
}

// FundNavRange is a net-value range option (relative months / years before now).
type FundNavRange struct {
	MonthBefore int32 // optional, number of months before now
	YearBefore  int32 // optional, number of years before now
}

// GetFundHoldings is the request for the fund holdings endpoint.
type GetFundHoldings struct {
	Scene int32 // optional
}

// GetFundStockHoldings is the request for the fund stock-holdings (reverse)
// endpoint.
type GetFundStockHoldings struct {
	Limit int32 // optional, maximum number of stocks to return
}

// GetFundPositions is the request for the fund positions overview endpoint.
type GetFundPositions struct {
	AccountChannel string // optional
	Aaid           int64  // optional, account id
}

// GetFundPosition is the request for a single fund position detail.
type GetFundPosition struct {
	AccountChannel string // optional
	Aaid           int64  // optional, account id
	Start          string // optional, range start
	End            string // optional, range end
}

// GetFundPositionProfits is the request for a single fund position
// cumulative-profit list.
type GetFundPositionProfits struct {
	AccountChannel string // optional
	Aaid           int64  // optional, account id
	Start          string // optional, range start
	End            string // optional, range end
	Page           int32  // optional
	Size           int32  // optional
}

// GetFundPositionDividends is the request for a single fund position dividend
// list.
type GetFundPositionDividends struct {
	AccountChannel string // optional
	Aaid           int64  // optional, account id
	Currency       string // optional
	Start          int64  // optional, range start (unix seconds)
	End            int64  // optional, range end (unix seconds)
	Page           int32  // optional
	Size           int32  // optional
}

// GetFundOrders is the request for the fund orders list.
type GetFundOrders struct {
	Symbols  []string // optional, filter by fund symbols (query key "symbol")
	Actions  string   // optional, comma-separated actions
	States   string   // optional, comma-separated states
	Currency string   // optional
	Start    int64    // optional, range start (unix seconds)
	End      int64    // optional, range end (unix seconds)
	Page     int32    // optional
	Size     int32    // optional
}

// GetFundTransactions is the request for the fund transactions (cash-flow) list.
type GetFundTransactions struct {
	AccountChannel string // optional
	BusinessType   string // optional
	Category       string // optional
	Currencies     string // optional, comma-separated
	Start          int64  // optional, range start (unix seconds)
	End            int64  // optional, range end (unix seconds)
	Page           int32  // optional
	Size           int32  // optional
}

// ValidateFundOrder is the request for validating a fund order before
// submitting. Symbol, Action and Currency are required; the remaining fields
// are optional (nil pointers are omitted).
type ValidateFundOrder struct {
	Symbol         string // required, fund symbol
	Action         string // required, buy / sell
	Currency       string // required
	Amount         string // optional, for amount-based orders
	Units          string // optional, for unit-based orders
	DividendOption *int32 // optional
	FundSource     *int32 // optional
	AccountChannel string // optional
}

// SubmitFundOrder is the request for submitting a fund order (buy / sell).
// Symbol, Action and Currency are required; the remaining fields are optional
// (nil pointers are omitted).
type SubmitFundOrder struct {
	Symbol         string // required, fund symbol
	Action         string // required, buy / sell
	Currency       string // required
	Amount         string // optional, for amount-based orders
	Units          string // optional, for unit-based orders
	DividendOption *int32 // optional
	Fee            string // optional
	IsSellAll      *bool  // optional
	Remark         string // optional
	TradeMethod    *int32 // optional
}
