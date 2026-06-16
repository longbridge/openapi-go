package market

import (
	"encoding/json"
	"testing"
)

func TestTradeStatusJSONDeserializesNumericCodes(t *testing.T) {
	var status TradeStatus

	if err := json.Unmarshal([]byte("202"), &status); err != nil {
		t.Fatalf("unmarshal known status: %v", err)
	}
	if status != TradeStatusUSTrading {
		t.Fatalf("known status = %v, want %v", status, TradeStatusUSTrading)
	}

	if err := json.Unmarshal([]byte("456"), &status); err != nil {
		t.Fatalf("unmarshal unknown status: %v", err)
	}
	if status != TradeStatusUnknown {
		t.Fatalf("unknown status = %v, want %v", status, TradeStatusUnknown)
	}

	data, err := json.Marshal(TradeStatusUSClean)
	if err != nil {
		t.Fatalf("marshal status: %v", err)
	}
	if string(data) != "206" {
		t.Fatalf("marshal status = %s, want 206", data)
	}
}

func TestTradeStatusHelpersMatchOpenAPIDefinition(t *testing.T) {
	cases := []struct {
		code int32
		want TradeStatus
		name string
	}{
		{101, TradeStatusClean, "Closed"},
		{123, TradeStatusRealtimeQuote, "Temporary Break"},
		{202, TradeStatusUSTrading, "Trading"},
		{1009, TradeStatusUnited, "Not Listed"},
		{1010, TradeStatusTradingHalt, "Terminated"},
		{2001, TradeStatusFuse, "Fuse"},
	}

	for _, tc := range cases {
		status := TradeStatusFromCode(tc.code)
		if status != tc.want {
			t.Fatalf("TradeStatusFromCode(%d) = %v, want %v", tc.code, status, tc.want)
		}
		if status.Code() != tc.code {
			t.Fatalf("TradeStatusFromCode(%d).Code() = %d, want %d", tc.code, status.Code(), tc.code)
		}
		if status.Name() != tc.name {
			t.Fatalf("TradeStatusFromCode(%d).Name() = %q, want %q", tc.code, status.Name(), tc.name)
		}
	}

	if got := TradeStatusFromCode(456); got != TradeStatusUnknown {
		t.Fatalf("TradeStatusFromCode(456) = %v, want %v", got, TradeStatusUnknown)
	}
}

func TestTradeStatusNormalizesEngineAliases(t *testing.T) {
	cases := []struct {
		status TradeStatus
		want   TradeStatus
	}{
		{TradeStatusClean, TradeStatusClosing},
		{TradeStatusUSClean, TradeStatusUSPrev},
		{TradeStatusUSPrevMarketClean, TradeStatusUSClosing},
		{TradeStatusUSAfterMarketClean, TradeStatusUSTrading},
	}

	for _, tc := range cases {
		if got := tc.status.Normalize(); got != tc.want {
			t.Fatalf("%v.Normalize() = %v, want %v", tc.status, got, tc.want)
		}
	}
}

func TestTradeStatusLabelMatchesEngineDisplay(t *testing.T) {
	cases := []struct {
		status TradeStatus
		want   string
	}{
		{TradeStatusUSPrev, "Pre-Market"},
		{TradeStatusUSClean, "Pre-Market"},
		{TradeStatusUSAfter, "Post-Market"},
		{TradeStatusUSClosing, "Closed"},
		{TradeStatusUSAfterMarketClean, "Trading"},
		{TradeStatusUSTrading, "Trading"},
		{TradeStatusTrading, "Trading"},
		{TradeStatusClean, "Closed"},
		{TradeStatusOpenBid, ""},
		{TradeStatusNoonClosing, ""},
	}

	for _, tc := range cases {
		if got := tc.status.Label(); got != tc.want {
			t.Fatalf("%v.Label() = %q, want %q", tc.status, got, tc.want)
		}
	}
}

func TestMarketTimeItemUsesTradeStatusType(t *testing.T) {
	item := MarketTimeItem{
		TradeStatus:      TradeStatusUSTrading,
		DelayTradeStatus: TradeStatusUSClosing,
	}

	if item.TradeStatus != TradeStatusUSTrading {
		t.Fatalf("TradeStatus = %v, want %v", item.TradeStatus, TradeStatusUSTrading)
	}
	if item.DelayTradeStatus != TradeStatusUSClosing {
		t.Fatalf("DelayTradeStatus = %v, want %v", item.DelayTradeStatus, TradeStatusUSClosing)
	}
}
