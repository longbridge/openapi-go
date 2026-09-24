package trade

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/trade/jsontypes"
)

// TestSubmitMultiLegOrderSerialize verifies the domain request is copied into the wire
// type with decimal legs rendered as strings and optional fields omitted when empty.
func TestSubmitMultiLegOrderSerialize(t *testing.T) {
	params := &SubmitMultiLegOrder{
		Side:              OrderSideBuy,
		OrderType:         OrderTypeLO,
		SubmittedQuantity: decimal.RequireFromString("1"),
		Strategy:          MultiLegStrategyVerticalCallSpread,
		Legs: []SubmitMultiLegOrderLeg{
			{Symbol: "QQQ260731C764000.US", RatioQuantity: decimal.RequireFromString("1")},
			{Symbol: "QQQ260731C767000.US", RatioQuantity: decimal.RequireFromString("2")},
		},
		SubmittedPrice: decimal.RequireFromString("3.5"),
	}

	var body jsontypes.SubmitMultiLegOrder
	if err := util.Copy(&body, params); err != nil {
		t.Fatalf("copy: %v", err)
	}
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	want := `{"side":"Buy","order_type":"LO","submitted_quantity":"1","strategy":"VerticalCallSpread","legs":[{"symbol":"QQQ260731C764000.US","ratio_quantity":"1"},{"symbol":"QQQ260731C767000.US","ratio_quantity":"2"}],"submitted_price":"3.5"}`
	if got != want {
		t.Fatalf("wire mismatch\n got: %s\nwant: %s", got, want)
	}
}

// TestMultiLegInfoDeserialize verifies a server order carrying multi_leg is decoded into
// the typed domain MultiLegInfo, including nested leg enums and decimals.
func TestMultiLegInfoDeserialize(t *testing.T) {
	const raw = `{
		"order_id": "1",
		"status": "FilledStatus",
		"multi_leg": {
			"strategy": "VerticalCallSpread",
			"strategy_name": "Vertical Call Spread",
			"multileg_id": "ml-1",
			"code": "VCS",
			"legs": [
				{"symbol": "QQQ260731C764000.US", "side": "Buy", "position": "LONG", "ratio_quantity": "1", "strike_price": "764", "expire_date": "2026-07-31", "contract_direction": "C"}
			]
		}
	}`
	var jt jsontypes.Order
	if err := json.Unmarshal([]byte(raw), &jt); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	var order Order
	if err := util.Copy(&order, jt); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if order.MultiLeg == nil {
		t.Fatal("MultiLeg is nil")
	}
	if order.MultiLeg.Strategy != MultiLegStrategyVerticalCallSpread {
		t.Fatalf("strategy = %q", order.MultiLeg.Strategy)
	}
	if len(order.MultiLeg.Legs) != 1 {
		t.Fatalf("legs = %d", len(order.MultiLeg.Legs))
	}
	leg := order.MultiLeg.Legs[0]
	if leg.Position != MultiLegPositionLong || leg.ContractDirection != ContractDirectionCall {
		t.Fatalf("leg enums: pos=%q dir=%q", leg.Position, leg.ContractDirection)
	}
	if !leg.RatioQuantity.Equal(decimal.RequireFromString("1")) {
		t.Fatalf("ratio = %s", leg.RatioQuantity)
	}
	if leg.StrikePrice == nil || !leg.StrikePrice.Equal(decimal.RequireFromString("764")) {
		t.Fatalf("strike = %v", leg.StrikePrice)
	}
}
