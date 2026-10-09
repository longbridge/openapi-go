package trade

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/trade/jsontypes"
)

// TestEstimateMultiLegRequestSerialize verifies the domain request is copied into
// the wire body with the decimal quantity/price rendered as strings and each leg
// carrying its symbol.
func TestEstimateMultiLegRequestSerialize(t *testing.T) {
	params := &GetEstimateMultiLegAvailableQuantity{
		Side:              OrderSideBuy,
		OrderType:         OrderTypeLO,
		SubmittedQuantity: decimal.RequireFromString("1"),
		Strategy:          MultiLegStrategyVerticalCallSpread,
		SubmittedPrice:    decimal.RequireFromString("1.5"),
		Legs: []EstimateMultiLegOrderLeg{
			{Symbol: "QQQ260731C764000.US"},
			{Symbol: "QQQ260731C767000.US"},
		},
	}

	var body jsontypes.EstimateMultiLegAvailableQuantity
	if err := util.Copy(&body, params); err != nil {
		t.Fatalf("copy: %v", err)
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"side":"Buy","order_type":"LO","submitted_quantity":"1","strategy":"VerticalCallSpread","legs":[{"symbol":"QQQ260731C764000.US"},{"symbol":"QQQ260731C767000.US"}],"submitted_price":"1.5"}`
	if string(b) != want {
		t.Fatalf("wire mismatch\n got: %s\nwant: %s", string(b), want)
	}
}

// TestEstimateMultiLegResponseDeserialize verifies the wire response decodes and
// copies into the domain type with the margin fields parsed to decimals.
func TestEstimateMultiLegResponseDeserialize(t *testing.T) {
	raw := `{"max_open_qty":"10","unit_margin":"1500.00","initial_margin_change":"15000.00","maintenance_margin_change":"12000.00"}`
	var wire jsontypes.EstimateMultiLegAvailableQuantityResponse
	if err := json.Unmarshal([]byte(raw), &wire); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	var resp EstimateMultiLegAvailableQuantityResponse
	if err := util.Copy(&resp, wire); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if !resp.MaxOpenQty.Equal(decimal.RequireFromString("10")) {
		t.Fatalf("max_open_qty mismatch: %s", resp.MaxOpenQty)
	}
	if !resp.UnitMargin.Equal(decimal.RequireFromString("1500.00")) {
		t.Fatalf("unit_margin mismatch: %s", resp.UnitMargin)
	}
	if !resp.MaintenanceMarginChange.Equal(decimal.RequireFromString("12000.00")) {
		t.Fatalf("maintenance_margin_change mismatch: %s", resp.MaintenanceMarginChange)
	}
}
