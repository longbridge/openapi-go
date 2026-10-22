package trade

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/longbridge/openapi-go/http"
)

func newTestTradeContext(t *testing.T, url, prefix string) *TradeContext {
	t.Helper()
	cli, err := http.New(
		http.WithURL(url),
		http.WithAppKey(prefix+"key"),
		http.WithAppSecret(prefix+"secret"),
		http.WithAccessToken(prefix+"token"),
	)
	if err != nil {
		t.Fatalf("http.New: %v", err)
	}
	return &TradeContext{opts: &Options{httpClient: cli}}
}

func TestUSExchangeRatesFlatten(t *testing.T) {
	data := `{
		"universal_exchanges": {
			"2": {"spec_exchanges": {"USD": {"exchanges": [
				{"base_currency": "USD", "other_currency": "HKD", "bid_rate": 7.8, "offer_rate": "7.85", "average_rate": null}
			]}}},
			"1": {"spec_exchanges": {
				"b": {"exchanges": [{"base_currency": "USD", "other_currency": "CNH", "bid_rate": "7.1"}]},
				"a": {"exchanges": null}
			}}
		}
	}`
	var raw usRawExchangeRates
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	got := raw.flatten().Exchanges
	want := []USExchangeRate{
		{UniversalKey: "1", SpecKey: "b", BaseCurrency: "USD", OtherCurrency: "CNH", BidRate: "7.1"},
		{UniversalKey: "2", SpecKey: "USD", BaseCurrency: "USD", OtherCurrency: "HKD", BidRate: "7.8", OfferRate: "7.85"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d entries, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("entry %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestUSTransactionHistoryUnmarshal(t *testing.T) {
	data := `{
		"current_page": "2",
		"items": [
			{"title": "Buy", "amount": "-100", "currency": "USD", "symbol": "AAPL.US", "create_at": "1792600000", "done_at": 1792600060},
			{"title": "Deposit", "amount": "500", "currency": "", "symbol": "BTCUSD.BKKT", "counter_id": "VA/BKKT/BTCUSD", "done_at": "0", "remark": null}
		]
	}`
	var resp USTransactionHistory
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.CurrentPage != 2 || len(resp.Items) != 2 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.Items[0].Symbol != "AAPL.US" || resp.Items[0].CreateAt.Unix() != 1792600000 || resp.Items[0].DoneAt.Unix() != 1792600060 {
		t.Errorf("item 0 = %+v", resp.Items[0])
	}
	if resp.Items[1].Symbol != "BTCUSD.BKKT" || !resp.Items[1].DoneAt.IsZero() {
		t.Errorf("item 1 = %+v", resp.Items[1])
	}
}

// TestUSAssetProdShapes decodes the field names, JSON types and map keys that
// US prod returns.
func TestUSAssetProdShapes(t *testing.T) {
	var hist USTransactionHistory
	if err := json.Unmarshal([]byte(`{"current_page": 1, "items": [{
		"amount": "0.5", "biz_code": "LVAVNSCR", "counter_id": "VA/BKKT/BTCUSD",
		"create_at": "1791478565", "currency": "", "done_at": "1791478565",
		"equity_type": "CRYPTO", "market": "BKKT", "position_side": "LONG",
		"quantity": "0.5", "symbol": "BTCUSD.BKKT", "title": "Cryptocurrency Deposit: BTCUSD"
	}]}`), &hist); err != nil {
		t.Fatalf("unmarshal transaction history: %v", err)
	}
	if hist.CurrentPage != 1 || hist.Items[0].Symbol != "BTCUSD.BKKT" || hist.Items[0].DoneAt.Unix() != 1791478565 {
		t.Errorf("transaction history = %+v", hist)
	}

	var raw usRawExchangeRates
	if err := json.Unmarshal([]byte(`{"universal_exchanges": {"us_lb_apex": {"spec_exchanges": {
		"reference": {"exchanges": [{"average_rate": 1, "base_currency": "USD", "bid_rate": 1, "offer_rate": 1, "other_currency": "USD"}]},
		"am": {"exchanges": [{"average_rate": 1, "base_currency": "USD", "bid_rate": 1, "offer_rate": 1, "other_currency": "USD"}]}
	}}}}`), &raw); err != nil {
		t.Fatalf("unmarshal exchange rates: %v", err)
	}
	rates := raw.flatten().Exchanges
	if len(rates) != 2 || rates[0].UniversalKey != "us_lb_apex" || rates[0].SpecKey != "am" || rates[1].SpecKey != "reference" || rates[0].AverageRate != "1" {
		t.Errorf("exchange rates = %+v", rates)
	}

	var overview usRawAssetOverview
	if err := json.Unmarshal([]byte(`{"account_type": "1", "asset_timestamp": "1792661597",
		"stock_list": [{"symbol": "MULL", "full_symbol": "MULL.US", "counter_id": "ETF/US/MULL", "quantity": "10"}],
		"crypto_list": [{"asset_type": "CRYPTO", "average_cost": "", "counter_id": "VA/BKKT/ETHUSD",
			"currency": "USD", "last_done": "2500", "market": "BKKT", "market_price": "2500",
			"name": "Ethereum", "open_position_time": "0", "position_side": "LONG",
			"quantity": "0.1", "symbol": "ETHUSD", "today_pl": "", "trade_status": "105"}],
		"option_list": [], "multi_leg": {"holdings": []}}`), &overview); err != nil {
		t.Fatalf("unmarshal asset overview: %v", err)
	}
	if overview.StockList[0].FullSymbol != "MULL.US" {
		t.Errorf("stock entry = %+v", overview.StockList[0])
	}
	c := overview.CryptoList[0]
	if c.Symbol != "ETHUSD" || c.Market != "BKKT" || c.Quantity != "0.1" || c.Name != "Ethereum" || c.MarketPrice != "2500" {
		t.Errorf("crypto entry = %+v", c)
	}
}

func TestUSTransactionHistoryRequest(t *testing.T) {
	var body map[string]interface{}
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Method != nethttp.MethodPost || r.URL.Path != "/v1/us/assets/transaction-history" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"data":{"current_page":1,"items":[]}}`))
	}))
	defer srv.Close()

	tc := newTestTradeContext(t, srv.URL, "us_")
	start := time.Unix(1790000000, 0)
	end := time.Unix(1790086400, 0)
	resp, err := tc.USTransactionHistory(context.Background(), &GetUSTransactionHistory{StartAt: start, EndAt: end})
	if err != nil {
		t.Fatalf("USTransactionHistory: %v", err)
	}
	if resp.CurrentPage != 1 {
		t.Errorf("current_page = %d", resp.CurrentPage)
	}
	if body["start"] != float64(1790000000) || body["end"] != float64(1790086400) || body["page"] != float64(1) || body["size"] != float64(20) {
		t.Errorf("unexpected body: %v", body)
	}
	if _, ok := body["types"]; ok {
		t.Errorf("empty types should be omitted: %v", body)
	}
}

func TestUSAssetRegionRestrictions(t *testing.T) {
	ctx := context.Background()
	var rre *http.RegionRestrictedError

	ap := newTestTradeContext(t, "http://127.0.0.1:1", "ap_")
	if _, err := ap.USCashAndBuyPower(ctx, ""); !errors.As(err, &rre) {
		t.Errorf("USCashAndBuyPower with AP token: got %v, want RegionRestrictedError", err)
	}
	if _, err := ap.USAvailableHolding(ctx, &GetUSAvailableHolding{Symbol: "AAPL.US"}); !errors.As(err, &rre) {
		t.Errorf("USAvailableHolding with AP token: got %v, want RegionRestrictedError", err)
	}
	if _, err := ap.USTransactionHistory(ctx, &GetUSTransactionHistory{StartAt: time.Unix(1, 0), EndAt: time.Unix(2, 0)}); !errors.As(err, &rre) {
		t.Errorf("USTransactionHistory with AP token: got %v, want RegionRestrictedError", err)
	}
	if _, err := ap.USExchangeRates(ctx); !errors.As(err, &rre) {
		t.Errorf("USExchangeRates with AP token: got %v, want RegionRestrictedError", err)
	}

	us := newTestTradeContext(t, "http://127.0.0.1:1", "us_")
	if _, err := us.FundPositions(ctx, nil); !errors.As(err, &rre) {
		t.Errorf("FundPositions with US token: got %v, want RegionRestrictedError", err)
	}
	if _, err := us.MarginRatio(ctx, "AAPL.US"); !errors.As(err, &rre) {
		t.Errorf("MarginRatio with US token: got %v, want RegionRestrictedError", err)
	}
	if _, err := us.OrderDetail(ctx, "1"); !errors.As(err, &rre) {
		t.Errorf("OrderDetail with US token: got %v, want RegionRestrictedError", err)
	}
	if _, err := us.OrderDetailAttached(ctx, "1"); !errors.As(err, &rre) {
		t.Errorf("OrderDetailAttached with US token: got %v, want RegionRestrictedError", err)
	}
}
