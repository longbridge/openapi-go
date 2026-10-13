package quote

import (
	"context"
	"testing"
	"time"

	quotev1 "github.com/longbridge/openapi-protobufs/gen/go/quote"
)

// One-minute bars for a New York day from `start` for `n` minutes, with the
// trade session derived from the clock.
func testTimeline(start time.Time, n int) []*quotev1.Candlestick {
	out := make([]*quotev1.Candlestick, 0, n)
	for i := 0; i < n; i++ {
		ts := start.Add(time.Duration(i) * time.Minute)
		local := ts.In(newYork)
		var session quotev1.TradeSession
		switch h := local.Hour(); {
		case h < 4 || h >= 20:
			session = quotev1.TradeSession_OVERNIGHT_TRADE
		case h < 9 || (h == 9 && local.Minute() < 30):
			session = quotev1.TradeSession_PRE_TRADE
		case h < 16:
			session = quotev1.TradeSession_NORMAL_TRADE
		default:
			session = quotev1.TradeSession_POST_TRADE
		}
		out = append(out, &quotev1.Candlestick{
			Close: "1", Open: "1", Low: "1", High: "1", Volume: 1, Turnover: "1",
			Timestamp:    ts.Unix(),
			TradeSession: session,
		})
	}
	return out
}

// A fake server: `count` bars from the anchor (inclusive) in `forward`
// direction, ascending.
func testServer(data []*quotev1.Candlestick, forward bool, calls *int) candlestickFetch {
	return func(_ context.Context, anchor time.Time, count int32) ([]*quotev1.Candlestick, error) {
		*calls++
		at := anchor.Unix()
		var page []*quotev1.Candlestick
		if forward {
			for _, c := range data {
				if c.Timestamp >= at && len(page) < int(count) {
					page = append(page, c)
				}
			}
		} else {
			for i := len(data) - 1; i >= 0 && len(page) < int(count); i-- {
				if data[i].Timestamp <= at {
					page = append(page, data[i])
				}
			}
			for i, j := 0, len(page)-1; i < j; i, j = i+1, j-1 {
				page[i], page[j] = page[j], page[i]
			}
		}
		return page, nil
	}
}

func hhmm(c *quotev1.Candlestick) string {
	return time.Unix(c.Timestamp, 0).In(newYork).Format("15:04")
}

func assertSortedUnique(t *testing.T, cs []*quotev1.Candlestick) {
	t.Helper()
	for i := 1; i < len(cs); i++ {
		if cs[i-1].Timestamp >= cs[i].Timestamp {
			t.Fatalf("not ascending/unique at %d", i)
		}
	}
}

var testDay = time.Date(2026, 10, 9, 0, 0, 0, 0, newYork) // EDT

func TestRefillBackwardAfterDroppingOvernight(t *testing.T) {
	data := testTimeline(testDay, 24*60)
	// Raw page 19:00–20:39 NY: 60 post + 40 overnight.
	end := 20*60 + 40
	page := data[end-100 : end]
	cursor := rawEdge(page, false)
	kept, dropped := dropOvernight(page)
	if len(kept) != 60 || dropped != 40 {
		t.Fatalf("kept=%d dropped=%d", len(kept), dropped)
	}
	calls := 0
	out, err := refillCandlesticks(context.Background(), false, kept, cursor, 100, time.Time{}, testServer(data, false, &calls))
	if err != nil {
		t.Fatal(err)
	}
	assertSortedUnique(t, out)
	if len(out) != 100 || hhmm(out[0]) != "18:20" || hhmm(out[99]) != "19:59" {
		t.Fatalf("len=%d first=%s last=%s", len(out), hhmm(out[0]), hhmm(out[len(out)-1]))
	}
	if calls != 1 {
		t.Fatalf("calls=%d", calls)
	}
}

func TestRefillForwardAfterDroppingOvernight(t *testing.T) {
	data := testTimeline(testDay, 24*60)
	// First 100 raw bars from 03:00 NY: 60 overnight + 40 pre.
	page := data[3*60 : 3*60+100]
	cursor := rawEdge(page, true)
	kept, dropped := dropOvernight(page)
	if len(kept) != 40 || dropped != 60 {
		t.Fatalf("kept=%d dropped=%d", len(kept), dropped)
	}
	calls := 0
	out, err := refillCandlesticks(context.Background(), true, kept, cursor, 100, time.Time{}, testServer(data, true, &calls))
	if err != nil {
		t.Fatal(err)
	}
	assertSortedUnique(t, out)
	if len(out) != 100 || hhmm(out[0]) != "04:00" || hhmm(out[99]) != "05:39" {
		t.Fatalf("len=%d first=%s last=%s", len(out), hhmm(out[0]), hhmm(out[len(out)-1]))
	}
	if calls != 1 {
		t.Fatalf("calls=%d", calls)
	}
}

func TestRefillPageOfOnlyOvernightBarsMakesProgress(t *testing.T) {
	data := testTimeline(testDay, 24*60)
	page := data[21*60 : 21*60+100] // 21:00–22:39, all overnight
	cursor := rawEdge(page, false)
	kept, dropped := dropOvernight(page)
	if len(kept) != 0 || dropped != 100 {
		t.Fatalf("kept=%d dropped=%d", len(kept), dropped)
	}
	calls := 0
	out, err := refillCandlesticks(context.Background(), false, kept, cursor, 100, time.Time{}, testServer(data, false, &calls))
	if err != nil {
		t.Fatal(err)
	}
	assertSortedUnique(t, out)
	if len(out) != 100 || hhmm(out[0]) != "18:20" || hhmm(out[99]) != "19:59" {
		t.Fatalf("len=%d first=%s last=%s", len(out), hhmm(out[0]), hhmm(out[len(out)-1]))
	}
	if calls > 3 {
		t.Fatalf("took %d rounds", calls)
	}
}

func TestRefillStopsAtBoundAndWhenDataRunsOut(t *testing.T) {
	data := testTimeline(testDay, 24*60) // only this day exists
	page := data[:100]                   // 00:00–01:39, all overnight
	cursor := rawEdge(page, false)
	kept, _ := dropOvernight(page)

	calls := 0
	out, err := refillCandlesticks(context.Background(), false, kept, cursor, 100, testDay, testServer(data, false, &calls))
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 || calls != 1 {
		t.Fatalf("bound: len=%d calls=%d", len(out), calls)
	}

	calls = 0
	out, err = refillCandlesticks(context.Background(), false, kept, cursor, 100, time.Time{}, testServer(data, false, &calls))
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 || calls != 1 {
		t.Fatalf("exhausted: len=%d calls=%d", len(out), calls)
	}
}

func TestRefillRoundLimit(t *testing.T) {
	// A pathological server that always returns a full page of overnight bars.
	calls := 0
	ts := testDay.Unix()
	fetch := func(_ context.Context, _ time.Time, count int32) ([]*quotev1.Candlestick, error) {
		calls++
		ts -= 60 * int64(count)
		page := make([]*quotev1.Candlestick, 0, count)
		for i := int32(0); i < count; i++ {
			page = append(page, &quotev1.Candlestick{
				Timestamp:    ts + int64(i)*60,
				TradeSession: quotev1.TradeSession_OVERNIGHT_TRADE,
			})
		}
		return page, nil
	}
	out, err := refillCandlesticks(context.Background(), false, nil, testDay.Unix(), 50, time.Time{}, fetch)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 || calls != maxRefillRounds {
		t.Fatalf("len=%d calls=%d", len(out), calls)
	}
}

func TestUSOvernightWindow(t *testing.T) {
	for _, day := range []time.Time{
		time.Date(2026, 10, 9, 0, 0, 0, 0, newYork), // EDT
		time.Date(2026, 1, 9, 0, 0, 0, 0, newYork),  // EST
	} {
		at := func(h, m int) int64 { return day.Add(time.Duration(h*60+m) * time.Minute).Unix() }
		for _, tc := range []struct {
			h, m int
			want bool
		}{{0, 0, true}, {3, 59, true}, {4, 0, false}, {9, 30, false}, {16, 0, false}, {19, 59, false}, {20, 0, true}, {23, 59, true}} {
			if got := isUSOvernight(at(tc.h, tc.m)); got != tc.want {
				t.Fatalf("%s %02d:%02d overnight=%v want %v", day.Format("2006-01-02"), tc.h, tc.m, got, tc.want)
			}
		}
	}
	for sym, want := range map[string]bool{"aapl.us": true, ".US": false, "AAPL.US": true, "AAPL251017C340000.US": true, ".VIX.US": false, "700.HK": false} {
		if isUSEquitySymbol(sym) != want {
			t.Fatalf("isUSEquitySymbol(%q) != %v", sym, want)
		}
	}
}
