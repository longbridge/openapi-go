// Package jsontypes contains the raw JSON wire types for the forex (currency
// exchange) API. These types match the exact JSON field names returned by the
// Longbridge API. Use the parent forex package for the idiomatic Go types.
//
// Decimal money fields are kept as string (the backend sends them as decimal
// strings, and the order amounts arrive empty until the order is filled); the
// parent package converts them to *decimal.Decimal. The order state is an
// integer-coded enum.
package jsontypes

// ForexQuote is a forex quote.
type ForexQuote struct {
	QuoteID  string `json:"quote_id"`
	Rate     string `json:"rate"`
	ExpireAt Int64  `json:"expire_at"`
	CcyPair  string `json:"ccy_pair"`
}

// ForexOrderDetail is a forex order's state.
type ForexOrderDetail struct {
	State      int32  `json:"state"`
	Rate       string `json:"rate"`
	FromAmount string `json:"from_amount"`
	ToAmount   string `json:"to_amount"`
	FailReason string `json:"fail_reason"`
}

// GetForexQuoteBody is the request body for getting a forex quote.
type GetForexQuoteBody struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       string `json:"amount,omitempty"`
	TargetAmount string `json:"target_amount,omitempty"`
}

// SubmitForexOrderBody is the request body for submitting a forex order.
type SubmitForexOrderBody struct {
	QuoteID       string `json:"quote_id"`
	ClientOrderID string `json:"client_order_id"`
}
