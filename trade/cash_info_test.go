package trade

import (
	"encoding/json"
	"testing"

	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/trade/jsontypes"
)

func TestCashInfoRedemptionCash(t *testing.T) {
	var resp jsontypes.AccountBalances
	body := `{"list":[{"currency":"HKD","cash_infos":[
		{"currency":"HKD","withdraw_cash":"1","available_cash":"2","frozen_cash":"3","settling_cash":"4","redemption_cash":"1500.25"},
		{"currency":"USD","withdraw_cash":"1","available_cash":"2","frozen_cash":"3","settling_cash":"4"}
	]}]}`
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatal(err)
	}
	var accounts []*AccountBalance
	if err := util.Copy(&accounts, resp.List); err != nil {
		t.Fatal(err)
	}
	infos := accounts[0].CashInfos
	if got := infos[0].RedemptionCash; got == nil || got.String() != "1500.25" {
		t.Fatalf("redemption_cash: got %v, want 1500.25", got)
	}
	if got := infos[1].SettlingCash; got == nil || got.String() != "4" {
		t.Fatalf("settling_cash: got %v, want 4", got)
	}
	// Older responses omit redemption_cash; it must not fail the call.
	if got := infos[1].RedemptionCash; got != nil && !got.IsZero() {
		t.Fatalf("missing redemption_cash: got %v, want nil or 0", got)
	}
}
