package quote

import (
	"context"
	"strings"
	"time"
	_ "time/tzdata" // New York rules must resolve even without a system tz database

	quotev1 "github.com/longbridge/openapi-protobufs/gen/go/quote"
)

// Keeps HTTP-transport quote data aligned with the WebSocket's US overnight
// rule.
//
// The quote WebSocket only returns the US overnight session when the
// connection opted in (`need_over_night_quote`, i.e. WithEnableOvernight),
// and applies that before count / range caps. The REST endpoints always
// return overnight data and offer no opt-out, so with TransportHTTP and
// overnight disabled the SDK drops overnight bars client-side and, for capped
// candlestick windows, pages further so both transports yield the same bars.

const (
	// maxHistoryCandlesticks is the server-side cap on a candlestick request's count.
	maxHistoryCandlesticks = 1000
	// maxRefillRounds bounds the extra requests one call may issue while
	// topping up a window.
	maxRefillRounds = 8
)

var newYork = mustLoadLocation("America/New_York")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic("quote: load " + name + ": " + err.Error())
	}
	return loc
}

// alignsOvernight reports whether this context must drop US overnight data to
// match what the WebSocket returns.
func (c *core) alignsOvernight() bool {
	return c.transport == TransportHTTP && !c.enableOvernight
}

// isUSEquitySymbol is true for US equities / ETFs (`AAPL.US`), whose overnight
// session is what the WebSocket gates. Indices (`.VIX.US`) trade on their own
// schedule and are left alone.
func isUSEquitySymbol(symbol string) bool {
	return strings.HasSuffix(symbol, ".US") && !strings.HasPrefix(symbol, ".")
}

// isUSOvernight reports whether a unix timestamp falls in the US overnight
// session (20:00–04:00 New York).
func isUSOvernight(unix int64) bool {
	h := time.Unix(unix, 0).In(newYork).Hour()
	return h < 4 || h >= 20
}

func isOvernightCandlestick(candlestick *quotev1.Candlestick) bool {
	return candlestick.GetTradeSession() == quotev1.TradeSession_OVERNIGHT_TRADE
}

// newYorkLocal returns the New York wall-clock time of a unix timestamp, as
// used by the offset query.
func newYorkLocal(unix int64) time.Time {
	return time.Unix(unix, 0).In(newYork)
}

// newYorkMidnight returns midnight New York of the given `YYYYMMDD` date, or
// zero if it does not parse.
func newYorkMidnight(yyyymmdd string) time.Time {
	t, err := time.ParseInLocation("20060102", yyyymmdd, newYork)
	if err != nil {
		return time.Time{}
	}
	return t
}

// rawEdge is the earliest (`!forward`) or latest (`forward`) timestamp in a
// raw page, or 0 for an empty page.
func rawEdge(candlesticks []*quotev1.Candlestick, forward bool) int64 {
	var edge int64
	for _, candlestick := range candlesticks {
		ts := candlestick.GetTimestamp()
		if edge == 0 || (forward && ts > edge) || (!forward && ts < edge) {
			edge = ts
		}
	}
	return edge
}

// dropOvernight returns the non-overnight candlesticks of a raw page and how
// many were dropped.
func dropOvernight(candlesticks []*quotev1.Candlestick) ([]*quotev1.Candlestick, int) {
	kept := make([]*quotev1.Candlestick, 0, len(candlesticks))
	for _, candlestick := range candlesticks {
		if !isOvernightCandlestick(candlestick) {
			kept = append(kept, candlestick)
		}
	}
	return kept, len(candlesticks) - len(kept)
}

// candlestickFetch fetches `count` raw candlesticks from `anchor` (New York
// wall-clock time, inclusive) in the direction fixed by the caller.
type candlestickFetch func(ctx context.Context, anchor time.Time, count int32) ([]*quotev1.Candlestick, error)

// refillCandlesticks tops up a capped candlestick window after overnight bars
// were dropped.
//
// `candlesticks` is the already-filtered first page (ascending) and `cursor`
// the edge of its raw window (overnight bars included). Pages are fetched
// beyond the cursor until `target` bars are collected, the data runs out,
// `bound` (midnight New York of a date-range query's start, zero for none;
// backward queries only) is crossed, or maxRefillRounds is reached.
func refillCandlesticks(
	ctx context.Context,
	forward bool,
	candlesticks []*quotev1.Candlestick,
	cursor int64,
	target int,
	bound time.Time,
	fetch candlestickFetch,
) ([]*quotev1.Candlestick, error) {
	for round := 0; round < maxRefillRounds; round++ {
		if len(candlesticks) >= target || cursor == 0 {
			break
		}
		need := target - len(candlesticks)
		// Over-fetch: roughly a third of a US day is overnight, and the anchor
		// bar itself comes back and is discarded.
		count := need + need/2 + 1
		if count > maxHistoryCandlesticks {
			count = maxHistoryCandlesticks
		}
		raw, err := fetch(ctx, newYorkLocal(cursor), int32(count))
		if err != nil {
			return nil, err
		}
		exhausted := len(raw) < count
		nextCursor := rawEdge(raw, forward)

		more, _ := dropOvernight(raw)
		filtered := more[:0]
		crossedBound := false
		for _, candlestick := range more {
			ts := candlestick.GetTimestamp()
			if (forward && ts <= cursor) || (!forward && ts >= cursor) {
				continue
			}
			if !forward && !bound.IsZero() && ts < bound.Unix() {
				crossedBound = true
				continue
			}
			filtered = append(filtered, candlestick)
		}

		if forward {
			if len(filtered) > need {
				filtered = filtered[:need]
			}
			candlesticks = append(candlesticks, filtered...)
		} else {
			if len(filtered) > need {
				filtered = filtered[len(filtered)-need:]
			}
			candlesticks = append(append([]*quotev1.Candlestick{}, filtered...), candlesticks...)
		}

		if exhausted || crossedBound || nextCursor == 0 || nextCursor == cursor {
			break
		}
		cursor = nextCursor
	}
	return candlesticks, nil
}
