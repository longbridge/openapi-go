package quote

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	quotev1 "github.com/longbridge/openapi-protobufs/gen/go/quote"
)

// fakeHistory serves offset pages from `data` and records the requests.
func fakeHistory(data []*quotev1.Candlestick, calls *[]*quotev1.SecurityHistoryCandlestickRequest) historyFetch {
	return func(_ context.Context, req *quotev1.SecurityHistoryCandlestickRequest) ([]*quotev1.Candlestick, error) {
		*calls = append(*calls, req)
		oq := req.GetOffsetRequest()
		anchor, err := time.ParseInLocation("200601021504", oq.GetDate()+oq.GetMinute(), newYork)
		if err != nil {
			return nil, err
		}
		forward := oq.GetDirection() == quotev1.Direction_FORWARD
		calls2 := 0
		return testServer(data, forward, &calls2)(context.Background(), anchor, oq.GetCount())
	}
}

func offsetReq(symbol string, forward bool, count int32, session CandlestickTradeSession) *quotev1.SecurityHistoryCandlestickRequest {
	dir := quotev1.Direction_BACKWARD
	if forward {
		dir = quotev1.Direction_FORWARD
	}
	return &quotev1.SecurityHistoryCandlestickRequest{
		Symbol:        symbol,
		Period:        quotev1.Period_ONE_MINUTE,
		QueryType:     quotev1.HistoryCandlestickQueryType_QUERY_BY_OFFSET,
		OffsetRequest: &quotev1.SecurityHistoryCandlestickRequest_OffsetQuery{Direction: dir, Count: count},
		TradeSession:  int32(session),
	}
}

func dateReq(symbol, start, end string) *quotev1.SecurityHistoryCandlestickRequest {
	return &quotev1.SecurityHistoryCandlestickRequest{
		Symbol:       symbol,
		Period:       quotev1.Period_ONE_MINUTE,
		QueryType:    quotev1.HistoryCandlestickQueryType_QUERY_BY_DATE,
		DateRequest:  &quotev1.SecurityHistoryCandlestickRequest_DateQuery{StartDate: start, EndDate: end},
		TradeSession: int32(CandlestickTradeSessionAll),
	}
}

func TestAlignCandlesticksEarlyExits(t *testing.T) {
	data := testTimeline(testDay, 24*60)
	page := data[19*60 : 21*60] // 60 post + 60 overnight
	var calls []*quotev1.SecurityHistoryCandlestickRequest
	fetch := fakeHistory(data, &calls)

	// Not a US equity: untouched, overnight bars included.
	out, err := alignCandlesticks(context.Background(), offsetReq(".VIX.US", false, 120, CandlestickTradeSessionAll), page, fetch)
	if err != nil || len(out) != 120 || len(calls) != 0 {
		t.Fatalf("index: len=%d calls=%d err=%v", len(out), len(calls), err)
	}
	// Normal session only: untouched.
	out, _ = alignCandlesticks(context.Background(), offsetReq("AAPL.US", false, 120, CandlestickTradeSessionNormal), page, fetch)
	if len(out) != 120 || len(calls) != 0 {
		t.Fatalf("normal: len=%d calls=%d", len(out), len(calls))
	}
	// Nothing dropped: no top-up request.
	out, _ = alignCandlesticks(context.Background(), offsetReq("AAPL.US", false, 60, CandlestickTradeSessionAll), data[18*60:19*60], fetch)
	if len(out) != 60 || len(calls) != 0 {
		t.Fatalf("no-drop: len=%d calls=%d", len(out), len(calls))
	}
}

func TestAlignCandlesticksOffsetTargetIsMinOfCountAndPage(t *testing.T) {
	data := testTimeline(testDay, 24*60)
	page := data[19*60 : 21*60] // 60 post + 60 overnight, count asked 5000
	var calls []*quotev1.SecurityHistoryCandlestickRequest
	out, err := alignCandlesticks(context.Background(), offsetReq("AAPL.US", false, 5000, CandlestickTradeSessionAll), page, fakeHistory(data, &calls))
	if err != nil {
		t.Fatal(err)
	}
	// Target is the page size (120), not 5000: topped up to 120 bars 18:00–19:59.
	if len(out) != 120 || hhmm(out[0]) != "18:00" || hhmm(out[119]) != "19:59" {
		t.Fatalf("len=%d first=%s last=%s", len(out), hhmm(out[0]), hhmm(out[len(out)-1]))
	}
	for _, c := range calls {
		if c.GetTradeSession() != int32(CandlestickTradeSessionAll) || c.GetSymbol() != "AAPL.US" ||
			c.GetQueryType() != quotev1.HistoryCandlestickQueryType_QUERY_BY_OFFSET || c.GetDateRequest() != nil ||
			c.GetOffsetRequest().GetCount() > maxHistoryCandlesticks {
			t.Fatalf("bad top-up request: %+v", c)
		}
	}
}

func TestAlignCandlesticksDateRange(t *testing.T) {
	data := testTimeline(testDay, 24*60)
	var calls []*quotev1.SecurityHistoryCandlestickRequest
	fetch := fakeHistory(data, &calls)

	// Page reaches midnight of start: not capped → only dropping, no request.
	out, err := alignCandlesticks(context.Background(), dateReq("AAPL.US", "20261009", "20261009"), data, fetch)
	if err != nil || len(out) != 16*60 || len(calls) != 0 {
		t.Fatalf("uncapped: len=%d calls=%d err=%v", len(out), len(calls), err)
	}
	// Page starts at 20:00 of the start day: capped at 240 → refill back to
	// 240 non-overnight bars, bounded by midnight of the start day.
	out, err = alignCandlesticks(context.Background(), dateReq("AAPL.US", "20261009", "20261010"), data[20*60:], fetch)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 240 || hhmm(out[0]) != "16:00" || hhmm(out[239]) != "19:59" {
		t.Fatalf("capped: len=%d first=%s last=%s", len(out), hhmm(out[0]), hhmm(out[len(out)-1]))
	}
	// No start date and a page shorter than the cap: not treated as capped.
	calls = nil
	out, _ = alignCandlesticks(context.Background(), dateReq("AAPL.US", "", "20261010"), data[20*60:], fetch)
	if len(out) != 0 || len(calls) != 0 {
		t.Fatalf("no-start short page: len=%d calls=%d", len(out), len(calls))
	}
}

func TestRefillSmallNeedInsideOvernightBlock(t *testing.T) {
	// Backward: raw page 03:50–05:29 (10 overnight + 90 pre); the 10 missing
	// bars are on the far side of the 8-hour overnight block.
	data := testTimeline(testDay.AddDate(0, 0, -1), 2*24*60)
	start := 24*60 + 3*60 + 50
	page := data[start : start+100]
	kept, dropped := dropOvernight(page)
	if len(kept) != 90 || dropped != 10 {
		t.Fatalf("kept=%d dropped=%d", len(kept), dropped)
	}
	calls := 0
	out, err := refillCandlesticks(context.Background(), false, kept, rawEdge(page, false), 100, time.Time{}, testServer(data, false, &calls))
	if err != nil {
		t.Fatal(err)
	}
	assertSortedUnique(t, out)
	if len(out) != 100 || hhmm(out[0]) != "19:50" || hhmm(out[9]) != "19:59" || hhmm(out[10]) != "04:00" || calls > 3 {
		t.Fatalf("backward: len=%d first=%s calls=%d", len(out), hhmm(out[0]), calls)
	}

	// Forward: raw page 19:40–21:19 (20 post + 80 overnight); next bars at
	// 04:00 the following day.
	data = testTimeline(testDay, 2*24*60)
	start = 19*60 + 40
	page = data[start : start+100]
	kept, dropped = dropOvernight(page)
	if len(kept) != 20 || dropped != 80 {
		t.Fatalf("kept=%d dropped=%d", len(kept), dropped)
	}
	calls = 0
	out, err = refillCandlesticks(context.Background(), true, kept, rawEdge(page, true), 100, time.Time{}, testServer(data, true, &calls))
	if err != nil {
		t.Fatal(err)
	}
	assertSortedUnique(t, out)
	if len(out) != 100 || hhmm(out[19]) != "19:59" || hhmm(out[20]) != "04:00" || hhmm(out[99]) != "05:19" || calls > 3 {
		t.Fatalf("forward: len=%d last=%s calls=%d", len(out), hhmm(out[len(out)-1]), calls)
	}
}

func TestNormalizeProtoJSON(t *testing.T) {
	// Encodings protojson rejects on its own: "" on int64, "true" on bool,
	// a number on a string field, null inside a scalar list; plus the usual
	// int64-as-string, camelCase keys and unknown fields, which pass through.
	raw := json.RawMessage(`{"secuQuote":[{"symbol":"700.HK","last_done":424.8,"volume":"","timestamp":"1791533288",
		"trade_status":"","over_night_quote":null,"pre_market_quote":{"volume":"","turnover":7},"volume_str":""}]}`)
	var q quotev1.SecurityQuoteResponse
	normalized, err := normalizeProtoJSON(raw, q.ProtoReflect().Descriptor())
	if err != nil {
		t.Fatal(err)
	}
	if err := protoJSONUnmarshal.Unmarshal(normalized, &q); err != nil {
		t.Fatalf("%v\n%s", err, normalized)
	}
	got := q.GetSecuQuote()[0]
	if got.GetLastDone() != "424.8" || got.GetVolume() != 0 || got.GetTimestamp() != 1791533288 ||
		got.GetTradeStatus() != quotev1.TradeStatus_NORMAL || got.GetPreMarketQuote().GetTurnover() != "7" {
		t.Fatalf("decoded: %+v", got)
	}

	raw = json.RawMessage(`{"position":"1","broker_ids":["1",2,null,"-3",""]}`)
	var b quotev1.Brokers
	normalized, err = normalizeProtoJSON(raw, b.ProtoReflect().Descriptor())
	if err != nil {
		t.Fatal(err)
	}
	if err := protoJSONUnmarshal.Unmarshal(normalized, &b); err != nil {
		t.Fatalf("%v\n%s", err, normalized)
	}
	if len(b.GetBrokerIds()) != 3 || b.GetBrokerIds()[2] != -3 {
		t.Fatalf("broker ids: %v", b.GetBrokerIds())
	}

	raw = json.RawMessage(`{"price":"1","standard":"true"}`)
	var s quotev1.StrikePriceInfo
	normalized, err = normalizeProtoJSON(raw, s.ProtoReflect().Descriptor())
	if err != nil {
		t.Fatal(err)
	}
	if err := protoJSONUnmarshal.Unmarshal(normalized, &s); err != nil || !s.GetStandard() {
		t.Fatalf("bool: %v %+v", err, &s)
	}

	// Enum as a quoted number, bool as 0/1, "false", null in a message list.
	raw = json.RawMessage(`{"secu_quote":[null,{"symbol":"A","trade_status":"0"}]}`)
	var q2 quotev1.SecurityQuoteResponse
	normalized, _ = normalizeProtoJSON(raw, q2.ProtoReflect().Descriptor())
	if err := protoJSONUnmarshal.Unmarshal(normalized, &q2); err != nil || len(q2.GetSecuQuote()) != 1 {
		t.Fatalf("enum/null list: %v %s", err, normalized)
	}
	for in, want := range map[string]bool{`1`: true, `0`: false, `"false"`: false} {
		var s2 quotev1.StrikePriceInfo
		normalized, _ = normalizeProtoJSON(json.RawMessage(`{"standard":`+in+`}`), s2.ProtoReflect().Descriptor())
		if err := protoJSONUnmarshal.Unmarshal(normalized, &s2); err != nil || s2.GetStandard() != want {
			t.Fatalf("bool %s: %v %v", in, err, s2.GetStandard())
		}
	}

	// Still rejected: a non-numeric string on an int field.
	raw = json.RawMessage(`{"ask":[{"volume":"abc"}]}`)
	var d quotev1.SecurityDepthResponse
	normalized, _ = normalizeProtoJSON(raw, d.ProtoReflect().Descriptor())
	if err := protoJSONUnmarshal.Unmarshal(normalized, &d); err == nil {
		t.Fatal("garbage int accepted")
	}
}
