package trade

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/trade/jsontypes"
)

func marshalMap(t *testing.T, v interface{}) map[string]interface{} {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return m
}

// TestDelayedUntilSerialize verifies delayed_until reaches the wire for submit, multi-leg
// submit and replace, and is omitted when zero.
func TestDelayedUntilSerialize(t *testing.T) {
	var submit jsontypes.SubmitOrder
	if err := util.Copy(&submit, &SubmitOrder{Symbol: "700.HK", OrderType: OrderTypeLO, Side: OrderSideBuy, SubmittedQuantity: 100, SubmittedPrice: decimal.RequireFromString("300"), TimeInForce: TimeTypeDay, DelayedUntil: 1767225600}); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if got := marshalMap(t, submit)["delayed_until"]; got != float64(1767225600) {
		t.Errorf("submit delayed_until = %v", got)
	}

	var multileg jsontypes.SubmitMultiLegOrder
	if err := util.Copy(&multileg, &SubmitMultiLegOrder{Side: OrderSideBuy, OrderType: OrderTypeMO, SubmittedQuantity: decimal.RequireFromString("1"), Strategy: MultiLegStrategyStraddle, DelayedUntil: 1767225600}); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if got := marshalMap(t, multileg)["delayed_until"]; got != float64(1767225600) {
		t.Errorf("multileg delayed_until = %v", got)
	}

	var replace jsontypes.ReplaceOrder
	if err := util.Copy(&replace, &ReplaceOrder{OrderId: "1", Quantity: 100, DelayedUntil: ClearDelayedUntil}); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if got := marshalMap(t, replace)["delayed_until"]; got != float64(1) {
		t.Errorf("replace delayed_until = %v", got)
	}

	var regular jsontypes.ReplaceOrder
	if err := util.Copy(&regular, &ReplaceOrder{OrderId: "1", Quantity: 100}); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if _, ok := marshalMap(t, regular)["delayed_until"]; ok {
		t.Error("delayed_until should be omitted when zero")
	}
}

// TestDelayedUntilDeserialize verifies the order query responses carry delayed_until.
func TestDelayedUntilDeserialize(t *testing.T) {
	var wire jsontypes.Order
	if err := json.Unmarshal([]byte(`{"order_id":"1","status":"DelayedNotReported","delayed_until":"1767225600"}`), &wire); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	var order Order
	if err := util.Copy(&order, &wire); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if order.DelayedUntil != "1767225600" || order.Status != OrderDelayedNotReported {
		t.Errorf("order = %+v", order)
	}

	var wireDetail jsontypes.OrderDetail
	if err := json.Unmarshal([]byte(`{"order_id":"1","delayed_until":"0"}`), &wireDetail); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	var detail OrderDetail
	if err := util.Copy(&detail, &wireDetail); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if detail.DelayedUntil != "0" {
		t.Errorf("detail.DelayedUntil = %q", detail.DelayedUntil)
	}
}
