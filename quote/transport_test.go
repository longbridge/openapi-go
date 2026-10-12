package quote

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	quotev1 "github.com/longbridge/openapi-protobufs/gen/go/quote"
	"github.com/longbridge/openapi-protocol/go/client"

	lbhttp "github.com/longbridge/openapi-go/http"
)

func TestParseTransport(t *testing.T) {
	for in, want := range map[string]Transport{"ws": TransportWebSocket, "WS": TransportWebSocket, " http ": TransportHTTP} {
		got, ok := ParseTransport(in)
		if !ok || got != want {
			t.Fatalf("ParseTransport(%q) = %q, %v", in, got, ok)
		}
	}
	for _, bad := range []string{"grpc", "websocket"} {
		if _, ok := ParseTransport(bad); ok {
			t.Fatalf("%q accepted", bad)
		}
	}
}

// TestHTTPPathsCoverPullCommands guards the command → REST mapping: every
// pull command the core issues must be listed, and no push/subscription
// command may be.
func TestHTTPPathsCoverPullCommands(t *testing.T) {
	pull := []quotev1.Command{
		quotev1.Command_QueryMarketTradePeriod, quotev1.Command_QueryMarketTradeDay,
		quotev1.Command_QuerySecurityStaticInfo, quotev1.Command_QuerySecurityQuote,
		quotev1.Command_QueryOptionQuote, quotev1.Command_QueryWarrantQuote,
		quotev1.Command_QueryDepth, quotev1.Command_QueryBrokers,
		quotev1.Command_QueryParticipantBrokerIds, quotev1.Command_QueryTrade,
		quotev1.Command_QueryIntraday, quotev1.Command_QueryCandlestick,
		quotev1.Command_QueryOptionChainDate, quotev1.Command_QueryWarrantIssuerInfo,
		quotev1.Command_QueryWarrantFilterList, quotev1.Command_QueryCapitalFlowIntraday,
		quotev1.Command_QueryCapitalFlowDistribution, quotev1.Command_QuerySecurityCalcIndex,
		quotev1.Command_QueryHistoryCandlestick,
	}
	if len(httpPaths) != len(pull) {
		t.Fatalf("httpPaths has %d entries, want %d", len(httpPaths), len(pull))
	}
	for _, cmd := range pull {
		if _, ok := httpPaths[cmd]; !ok {
			t.Fatalf("%v has no REST path", cmd)
		}
	}
	for _, cmd := range []quotev1.Command{
		quotev1.Command_Subscribe, quotev1.Command_Unsubscribe, quotev1.Command_Subscription,
		quotev1.Command_QueryUserQuoteProfile,
	} {
		if _, ok := httpPaths[cmd]; ok {
			t.Fatalf("%v must stay on the WebSocket", cmd)
		}
	}
}

// TestDecodeGatewayProtoJSON replays the gateway's proto-JSON shape (proto
// field names, int64 as strings, int enums, null sub-messages, an extra
// `volume_str`) into the protobuf response types.
func TestDecodeGatewayProtoJSON(t *testing.T) {
	quotes := json.RawMessage(`{"secu_quote":[{"high":"425.800","last_done":"424.800","low":"414.800","open":"415.000",
		"over_night_quote":null,"post_market_quote":null,"pre_market_quote":{"last_done":"331.400","timestamp":"1791552600","volume":"1262677","turnover":"420496915.751","high":"340.590","low":"330.750","prev_close":"340.420"},
		"prev_close":"411.400","symbol":"700.HK","timestamp":"1791533288","trade_status":0,"turnover":"8610387928.000","volume":"20422500","volume_str":""}]}`)
	var q quotev1.SecurityQuoteResponse
	if err := protoJSONUnmarshal.Unmarshal(quotes, &q); err != nil {
		t.Fatal(err)
	}
	got := q.GetSecuQuote()[0]
	if got.GetSymbol() != "700.HK" || got.GetTimestamp() != 1791533288 || got.GetVolume() != 20422500 ||
		got.GetTradeStatus() != quotev1.TradeStatus_NORMAL || got.GetOverNightQuote() != nil ||
		got.GetPreMarketQuote().GetVolume() != 1262677 {
		t.Fatalf("decoded quote mismatch: %+v", got)
	}

	depth := json.RawMessage(`{"ask":[{"order_num":"272","position":1,"price":"425.000","volume":"908200","volume_str":""}],
		"bid":[{"order_num":"0","position":1,"price":"","volume":"0","volume_str":""}],"symbol":"700.HK"}`)
	var d quotev1.SecurityDepthResponse
	if err := protoJSONUnmarshal.Unmarshal(depth, &d); err != nil {
		t.Fatal(err)
	}
	if d.GetAsk()[0].GetOrderNum() != 272 || d.GetAsk()[0].GetPosition() != 1 || d.GetBid()[0].GetPrice() != "" {
		t.Fatalf("decoded depth mismatch: %+v", &d)
	}

	// Request encoding uses proto field names and omits absent sub-messages.
	req := &quotev1.SecurityHistoryCandlestickRequest{
		Symbol: "AAPL.US", Period: quotev1.Period_ONE_MINUTE,
		QueryType:    quotev1.HistoryCandlestickQueryType_QUERY_BY_DATE,
		DateRequest:  &quotev1.SecurityHistoryCandlestickRequest_DateQuery{StartDate: "20261008", EndDate: "20261009"},
		TradeSession: int32(CandlestickTradeSessionAll),
	}
	b, err := protoJSONMarshal.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if _, has := m["offset_request"]; has {
		t.Fatalf("absent offset_request must be omitted: %s", b)
	}
	if _, has := m["date_request"]; !has {
		t.Fatalf("date_request missing: %s", b)
	}
	if string(m["trade_session"]) != "100" || string(m["symbol"]) != `"AAPL.US"` {
		t.Fatalf("unexpected encoding: %s", b)
	}
}

// TestHTTPDoEndToEnd drives httpDo against a fake gateway: the request path
// and body, absent / null `data`, normalisation, and API errors.
func TestHTTPDoEndToEnd(t *testing.T) {
	var gotPath, gotBody string
	reply := ""
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		gotPath = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("content-type", "application/json")
		_, _ = io.WriteString(w, reply)
	}))
	defer srv.Close()
	hc, err := lbhttp.New(lbhttp.WithURL(srv.URL), lbhttp.WithAppKey("k"), lbhttp.WithAppSecret("s"), lbhttp.WithAccessToken("t"))
	if err != nil {
		t.Fatal(err)
	}
	c := &core{transport: TransportHTTP, httpClient: hc}
	ctx := context.Background()

	reply = `{"code":0,"data":{"ask":[{"position":1,"price":"1.5","volume":"","order_num":"3"}]}}`
	var d quotev1.SecurityDepthResponse
	if err := c.do(ctx, quotev1.Command_QueryDepth, &quotev1.SecurityRequest{Symbol: "700.HK"}, &d); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/quote/depth" || gotBody != `{"symbol":"700.HK"}` || d.GetAsk()[0].GetOrderNum() != 3 {
		t.Fatalf("path=%s body=%s depth=%+v", gotPath, gotBody, &d)
	}

	for _, empty := range []string{`{"code":0,"message":"ok"}`, `{"code":0,"data":null}`} {
		reply = empty
		var p quotev1.ParticipantBrokerIdsResponse
		if err := c.do(ctx, quotev1.Command_QueryParticipantBrokerIds, nil, &p); err != nil {
			t.Fatalf("%s: %v", empty, err)
		}
		if gotBody != `{}` || len(p.GetParticipantBrokerNumbers()) != 0 {
			t.Fatalf("%s: body=%s resp=%+v", empty, gotBody, &p)
		}
	}

	reply = `{"code":301600,"message":"invalid symbol"}`
	var q quotev1.SecurityQuoteResponse
	err = c.do(ctx, quotev1.Command_QuerySecurityQuote, &quotev1.MultiSecurityRequest{Symbol: []string{"X"}}, &q)
	var apiErr *lbhttp.ApiError
	if !errors.As(err, &apiErr) || apiErr.Code != 301600 {
		t.Fatalf("want ApiError 301600, got %v", err)
	}
}

// TestClosedContextDoesNotDial checks that nothing re-dials the WebSocket
// after Close on the HTTP transport, and that queued handlers are dropped.
func TestClosedContextDoesNotDial(t *testing.T) {
	c := &core{transport: TransportHTTP}
	c.withClient(func(client.Client) { t.Fatal("handler applied") })
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ws(context.Background()); !errors.Is(err, ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
	if len(c.pendingHandlers) != 0 {
		t.Fatal("pending handlers kept after Close")
	}
	subs, err := c.Subscriptions(context.Background())
	if err != nil || len(subs) != 0 {
		t.Fatalf("subscriptions without a connection: %v %v", subs, err)
	}
	if err := c.Subscribe(context.Background(), []string{"700.HK"}, []SubType{SubTypeQuote}, false); !errors.Is(err, ErrClosed) {
		t.Fatalf("Subscribe after Close must fail, got %v", err)
	}
}
