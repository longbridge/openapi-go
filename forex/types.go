package forex

import "github.com/shopspring/decimal"

// ForexOrderStatus is a forex order's state.
type ForexOrderStatus int32

const (
	// ForexOrderProcessing means the order is in progress (including manual
	// review); keep polling.
	ForexOrderProcessing ForexOrderStatus = 0
	// ForexOrderSuccess means the conversion succeeded.
	ForexOrderSuccess ForexOrderStatus = 1
	// ForexOrderFailed means the conversion failed; frozen funds have been
	// returned.
	ForexOrderFailed ForexOrderStatus = 2
)

// forexOrderStatusString maps ForexOrderStatus values to a readable name.
var forexOrderStatusString = map[ForexOrderStatus]string{
	ForexOrderProcessing: "Processing",
	ForexOrderSuccess:    "Success",
	ForexOrderFailed:     "Failed",
}

// String returns a readable name for a ForexOrderStatus.
func (s ForexOrderStatus) String() string {
	if name, ok := forexOrderStatusString[s]; ok {
		return name
	}
	return "Processing"
}

// ForexQuote is a forex quote.
type ForexQuote struct {
	// QuoteID is the quote id, used when submitting the order.
	QuoteID string
	// Rate is the customer execution rate in standard currency-pair terms
	// (1 BASE = Rate QUOTE), independent of the request from/to direction.
	Rate decimal.Decimal
	// ExpireAt is the quote expiry, Unix milliseconds.
	ExpireAt int64
	// CcyPair is the standard currency pair for Rate, format BASE/QUOTE.
	CcyPair string
}

// ForexOrderDetail is a forex order's state.
type ForexOrderDetail struct {
	// State is the order state.
	State ForexOrderStatus
	// Rate is the execution rate (standard currency-pair terms); nil until filled.
	Rate *decimal.Decimal
	// FromAmount is the converted-from amount (actual value once filled); nil until filled.
	FromAmount *decimal.Decimal
	// ToAmount is the converted-to amount (actual value once filled); nil until filled.
	ToAmount *decimal.Decimal
	// FailReason is the failure reason; empty when not failed.
	FailReason string
}

// GetForexQuote are the parameters for ForexContext.Quote. Provide at most one
// of Amount / TargetAmount; the server validates the choice.
type GetForexQuote struct {
	// From is the convert-out currency (ISO 4217, e.g. USD).
	From string
	// To is the convert-in currency (e.g. HKD).
	To string
	// Amount is the convert-out amount; mutually exclusive with TargetAmount.
	Amount *decimal.Decimal
	// TargetAmount is the convert-in amount; mutually exclusive with Amount.
	TargetAmount *decimal.Decimal
}

// SubmitForexOrder are the parameters for ForexContext.SubmitOrder.
type SubmitForexOrder struct {
	// QuoteID is the quote id returned by ForexContext.Quote.
	QuoteID string
	// ClientOrderID is a caller-supplied order id, unique across all of the
	// caller's accounts.
	ClientOrderID string
}
