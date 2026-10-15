package trade

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/trade/jsontypes"
)

// TestAttachedParamsDepthLevel verifies the attached-order depth levels reach the wire
// for submit and replace, and are omitted when zero.
func TestAttachedParamsDepthLevel(t *testing.T) {
	var submit jsontypes.SubmitOrder
	if err := util.Copy(&submit, &SubmitOrder{
		Symbol: "700.HK", OrderType: OrderTypeLO, Side: OrderSideBuy, SubmittedQuantity: 100,
		SubmittedPrice: decimal.RequireFromString("300"), TimeInForce: TimeTypeDay,
		AttachedParams: &SubmitAttachedParams{AttachedOrderType: AttachedOrderTypeBracket, ProfitTakerDepthLevel: 2, StopLossDepthLevel: 5},
	}); err != nil {
		t.Fatalf("copy: %v", err)
	}
	ap := marshalMap(t, submit)["attached_params"].(map[string]interface{})
	if ap["profit_taker_depth_level"] != float64(2) || ap["stop_loss_depth_level"] != float64(5) {
		t.Errorf("submit attached_params = %v", ap)
	}

	var replace jsontypes.ReplaceOrder
	if err := util.Copy(&replace, &ReplaceOrder{
		OrderId: "1", Quantity: 100,
		AttachedParams: &ReplaceAttachedParams{AttachedOrderType: AttachedOrderTypeBracket, StopLossDepthLevel: 3},
	}); err != nil {
		t.Fatalf("copy: %v", err)
	}
	ap = marshalMap(t, replace)["attached_params"].(map[string]interface{})
	if _, ok := ap["profit_taker_depth_level"]; ok {
		t.Errorf("profit_taker_depth_level should be omitted when zero: %v", ap)
	}
	if ap["stop_loss_depth_level"] != float64(3) {
		t.Errorf("replace attached_params = %v", ap)
	}
}

// TestAttachedOrderLimitDepthLevel verifies attached_orders[].limit_depth_level is read.
func TestAttachedOrderLimitDepthLevel(t *testing.T) {
	for raw, want := range map[string]int32{`3`: 3, `0`: 0, `null`: 0} {
		var wire jsontypes.OrderDetail
		body := `{"order_id":"1","attached_orders":[{"order_id":"2","limit_depth_level":` + raw + `}]}`
		if err := json.Unmarshal([]byte(body), &wire); err != nil {
			t.Fatalf("unmarshal %s: %v", raw, err)
		}
		var detail OrderDetail
		if err := util.Copy(&detail, &wire); err != nil {
			t.Fatalf("copy: %v", err)
		}
		if len(detail.AttachedOrders) != 1 || detail.AttachedOrders[0].LimitDepthLevel != want {
			t.Errorf("%s: attached orders = %+v", raw, detail.AttachedOrders)
		}
	}
}
