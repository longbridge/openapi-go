package trade

import (
	"context"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestHistoryExecutionsMissingSide(t *testing.T) {
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.URL.Path != "/v1/trade/execution/history" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"data":{"has_more":false,"trades":[
			{"order_id":"1","trade_id":"t1","symbol":"AAPL.US","trade_done_at":"1790000000","quantity":"1","price":"200"},
			{"order_id":"2","trade_id":"t2","symbol":"AAPL.US","trade_done_at":"1790000001","quantity":"1","price":"201","side":"Sell"}
		]}}`))
	}))
	defer srv.Close()

	tc := newTestTradeContext(t, srv.URL, "us_")
	trades, err := tc.HistoryExecutions(context.Background(), &GetHistoryExecutions{})
	if err != nil {
		t.Fatalf("HistoryExecutions: %v", err)
	}
	if len(trades) != 2 {
		t.Fatalf("got %d trades, want 2", len(trades))
	}
	if trades[0].Side != OrderSideUnknown {
		t.Errorf("missing side = %q, want OrderSideUnknown", trades[0].Side)
	}
	if trades[1].Side != OrderSideSell {
		t.Errorf("side = %q, want OrderSideSell", trades[1].Side)
	}
}
