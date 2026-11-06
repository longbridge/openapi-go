package forex

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/forex/jsontypes"
	"github.com/longbridge/openapi-go/internal/util"
)

// TestQuoteDeserialize verifies a quote's wire JSON is decoded and copied into
// the domain type, with rate parsed to a decimal and expire_at tolerated as a
// JSON number.
func TestQuoteDeserialize(t *testing.T) {
	raw := `{"ok":true,"quote_id":"Q1","rate":"7.7812","expire_at":1790000000000,"ccy_pair":"USD/HKD"}`
	var wire jsontypes.ForexQuote
	if err := json.Unmarshal([]byte(raw), &wire); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	quote := &ForexQuote{}
	if err := util.Copy(quote, &wire); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if quote.QuoteID != "Q1" || quote.CcyPair != "USD/HKD" {
		t.Fatalf("string fields mismatch: %+v", quote)
	}
	if !quote.Rate.Equal(decimal.RequireFromString("7.7812")) {
		t.Fatalf("rate mismatch: %s", quote.Rate)
	}
	if quote.ExpireAt != 1790000000000 {
		t.Fatalf("expire_at mismatch: %d", quote.ExpireAt)
	}
}

// TestOrderProcessingEmptyAmounts verifies an in-progress order (state 0) with
// empty amount strings decodes to the Processing state and nil decimal pointers.
func TestOrderProcessingEmptyAmounts(t *testing.T) {
	raw := `{"state":0,"rate":"","from_amount":"","to_amount":"","fail_reason":""}`
	var wire jsontypes.ForexOrderDetail
	if err := json.Unmarshal([]byte(raw), &wire); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	detail := &ForexOrderDetail{}
	if err := util.Copy(detail, &wire); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if detail.State != ForexOrderProcessing {
		t.Fatalf("state mismatch: %v", detail.State)
	}
	if detail.Rate != nil || detail.FromAmount != nil || detail.ToAmount != nil {
		t.Fatalf("expected nil amounts, got rate=%v from=%v to=%v", detail.Rate, detail.FromAmount, detail.ToAmount)
	}
}

// TestOrderSuccess verifies a filled order (state 1) decodes to Success with
// the amounts parsed to decimals.
func TestOrderSuccess(t *testing.T) {
	raw := `{"state":1,"rate":"7.7812","from_amount":"1000","to_amount":"7781.2","fail_reason":""}`
	var wire jsontypes.ForexOrderDetail
	if err := json.Unmarshal([]byte(raw), &wire); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	detail := &ForexOrderDetail{}
	if err := util.Copy(detail, &wire); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if detail.State != ForexOrderSuccess {
		t.Fatalf("state mismatch: %v", detail.State)
	}
	if detail.ToAmount == nil || !detail.ToAmount.Equal(decimal.RequireFromString("7781.2")) {
		t.Fatalf("to_amount mismatch: %v", detail.ToAmount)
	}
}

// TestQuoteBodySerialize verifies the quote request body renders the decimal
// amount as a string and omits target_amount when unset.
func TestQuoteBodySerialize(t *testing.T) {
	amount := decimal.RequireFromString("1000.00")
	body := jsontypes.GetForexQuoteBody{From: "USD", To: "HKD"}
	if amount != (decimal.Decimal{}) {
		body.Amount = amount.String()
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"from":"USD","to":"HKD","amount":"1000"}`
	if string(b) != want {
		t.Fatalf("wire mismatch\n got: %s\nwant: %s", string(b), want)
	}
}
