package quote

import (
	"encoding/json"
	"testing"

	"github.com/longbridge/openapi-go/quote/jsontypes"
)

// TestOptionChainContractMapping verifies the raw option_chain_list rows map onto the typed
// domain contract, including the empty-string special cases for option_type / standard_attr.
func TestOptionChainContractMapping(t *testing.T) {
	const raw = `{"list":[
		{"symbol":"AAPL230120C130000.US","expiry_date":"20230120","strike_price":"130","direction":"C","option_type":"","standard_attr":"","days_to_expiry":5},
		{"symbol":"BABA2261218C10000.US","expiry_date":"20261218","strike_price":"100","direction":"P","option_type":"W","standard_attr":"old","days_to_expiry":-1}
	]}`
	var resp jsontypes.OptionChainContractList
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.List) != 2 {
		t.Fatalf("list = %d", len(resp.List))
	}

	// Reproduce the mapping the context method performs.
	a := resp.List[0]
	if OptionDirection(a.Direction) != OptionDirectionCall {
		t.Fatalf("direction = %q", a.Direction)
	}
	if OptionExpiryCycleType(a.OptionType) != OptionExpiryCycleTypeMonthly {
		t.Fatalf("empty option_type should be monthly, got %q", a.OptionType)
	}
	if OptionStandardAttr(a.StandardAttr) != OptionStandardAttrNormal {
		t.Fatalf("empty standard_attr should be normal, got %q", a.StandardAttr)
	}

	b := resp.List[1]
	if OptionExpiryCycleType(b.OptionType) != OptionExpiryCycleTypeWeekly {
		t.Fatalf("option_type = %q", b.OptionType)
	}
	if OptionStandardAttr(b.StandardAttr) != OptionStandardAttrOld {
		t.Fatalf("standard_attr = %q", b.StandardAttr)
	}
	if b.DaysToExpiry != -1 {
		t.Fatalf("days_to_expiry = %d", b.DaysToExpiry)
	}
}
