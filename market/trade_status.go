package market

import "encoding/json"

// TradeStatus is a market trading status code returned by /v1/quote/market-status.
type TradeStatus int32

const (
	// TradeStatusUnknown is an unknown or unsupported market trading status.
	TradeStatusUnknown TradeStatus = -1
	// TradeStatusNoRegisterQuote means quote is not registered.
	TradeStatusNoRegisterQuote TradeStatus = 0
	// TradeStatusClean is clearing before the market opens.
	TradeStatusClean TradeStatus = 101
	// TradeStatusOpenBid is the opening auction.
	TradeStatusOpenBid TradeStatus = 102
	// TradeStatusMorningClosing is the morning break, currently used by VIX indexes.
	TradeStatusMorningClosing TradeStatus = 103
	// TradeStatusTrading is regular trading.
	TradeStatusTrading TradeStatus = 105
	// TradeStatusNoonClosing is the midday break.
	TradeStatusNoonClosing TradeStatus = 106
	// TradeStatusCloseBid is the closing auction.
	TradeStatusCloseBid TradeStatus = 107
	// TradeStatusClosing means the market is closed.
	TradeStatusClosing TradeStatus = 108
	// TradeStatusDarkWait is dark trading waiting to open.
	TradeStatusDarkWait TradeStatus = 110
	// TradeStatusDarkTrading is dark trading.
	TradeStatusDarkTrading TradeStatus = 111
	// TradeStatusDarkClosing is dark trading closed.
	TradeStatusDarkClosing TradeStatus = 112
	// TradeStatusAfterFix is after-hours fixed-price trading.
	TradeStatusAfterFix TradeStatus = 120
	// TradeStatusHalfClosing is half-day market closed.
	TradeStatusHalfClosing TradeStatus = 121
	// TradeStatusNotOpened means the exchange is waiting to open under special conditions.
	TradeStatusNotOpened TradeStatus = 122
	// TradeStatusRealtimeQuote is a temporary intraday break.
	TradeStatusRealtimeQuote TradeStatus = 123
	// TradeStatusUSPrev is US pre-market.
	TradeStatusUSPrev TradeStatus = 201
	// TradeStatusUSTrading is US regular trading.
	TradeStatusUSTrading TradeStatus = 202
	// TradeStatusUSAfter is US post-market.
	TradeStatusUSAfter TradeStatus = 203
	// TradeStatusUSClosing is US closed.
	TradeStatusUSClosing TradeStatus = 204
	// TradeStatusUSStop is US halted.
	TradeStatusUSStop TradeStatus = 205
	// TradeStatusUSClean is US clearing plus pre-market.
	TradeStatusUSClean TradeStatus = 206
	// TradeStatusUSNight is US overnight trading.
	TradeStatusUSNight TradeStatus = 207
	// TradeStatusUSPrevMarketClean is a US pre-market clearing alias returned by the quote engine.
	TradeStatusUSPrevMarketClean TradeStatus = 209
	// TradeStatusUSAfterMarketClean is a US post-market clearing alias returned by the quote engine.
	TradeStatusUSAfterMarketClean TradeStatus = 210
	// TradeStatusRefresh is stock refresh. It is deprecated in the status definition.
	TradeStatusRefresh TradeStatus = 1000
	// TradeStatusDelist is delisted.
	TradeStatusDelist TradeStatus = 1001
	// TradeStatusPrepare is preparing to list.
	TradeStatusPrepare TradeStatus = 1002
	// TradeStatusCodeChange is code changed.
	TradeStatusCodeChange TradeStatus = 1003
	// TradeStatusStop is halted.
	TradeStatusStop TradeStatus = 1004
	// TradeStatusWillOpen is waiting to open, typically for a US IPO auction.
	TradeStatusWillOpen TradeStatus = 1005
	// TradeStatusCommonSuspend is split or merge suspended.
	TradeStatusCommonSuspend TradeStatus = 1006
	// TradeStatusExpire is expired.
	TradeStatusExpire TradeStatus = 1007
	// TradeStatusNoQuote means no quote data.
	TradeStatusNoQuote TradeStatus = 1008
	// TradeStatusUnited is not listed. The historical variant name is kept for compatibility.
	TradeStatusUnited TradeStatus = 1009
	// TradeStatusTradingHalt is terminated trading, usually for warrants.
	TradeStatusTradingHalt TradeStatus = 1010
	// TradeStatusWaitListing is waiting to list, usually for new warrants.
	TradeStatusWaitListing TradeStatus = 1011
	// TradeStatusFuse is fuse.
	TradeStatusFuse TradeStatus = 2001
)

// TradeStatusFromCode converts a raw market trading status code to TradeStatus.
func TradeStatusFromCode(code int32) TradeStatus {
	status := TradeStatus(code)
	switch status {
	case TradeStatusUnknown,
		TradeStatusNoRegisterQuote,
		TradeStatusClean,
		TradeStatusOpenBid,
		TradeStatusMorningClosing,
		TradeStatusTrading,
		TradeStatusNoonClosing,
		TradeStatusCloseBid,
		TradeStatusClosing,
		TradeStatusDarkWait,
		TradeStatusDarkTrading,
		TradeStatusDarkClosing,
		TradeStatusAfterFix,
		TradeStatusHalfClosing,
		TradeStatusNotOpened,
		TradeStatusRealtimeQuote,
		TradeStatusUSPrev,
		TradeStatusUSTrading,
		TradeStatusUSAfter,
		TradeStatusUSClosing,
		TradeStatusUSStop,
		TradeStatusUSClean,
		TradeStatusUSNight,
		TradeStatusUSPrevMarketClean,
		TradeStatusUSAfterMarketClean,
		TradeStatusRefresh,
		TradeStatusDelist,
		TradeStatusPrepare,
		TradeStatusCodeChange,
		TradeStatusStop,
		TradeStatusWillOpen,
		TradeStatusCommonSuspend,
		TradeStatusExpire,
		TradeStatusNoQuote,
		TradeStatusUnited,
		TradeStatusTradingHalt,
		TradeStatusWaitListing,
		TradeStatusFuse:
		return status
	default:
		return TradeStatusUnknown
	}
}

// UnmarshalJSON decodes numeric market trading status codes.
func (s *TradeStatus) UnmarshalJSON(data []byte) error {
	var code int32
	if err := json.Unmarshal(data, &code); err != nil {
		return err
	}
	*s = TradeStatusFromCode(code)
	return nil
}

// Code returns the raw numeric status code.
func (s TradeStatus) Code() int32 {
	return int32(s)
}

// String returns the full English status name.
func (s TradeStatus) String() string {
	return s.Name()
}

// Label returns a simplified label for key display states.
func (s TradeStatus) Label() string {
	status := s.Normalize()
	switch status {
	case TradeStatusUSPrev,
		TradeStatusUSTrading,
		TradeStatusUSAfter,
		TradeStatusUSNight,
		TradeStatusUSClosing,
		TradeStatusTrading,
		TradeStatusClosing:
		return status.Name()
	default:
		return ""
	}
}

// Name returns the full English status name.
func (s TradeStatus) Name() string {
	switch s.Normalize() {
	case TradeStatusUnknown, TradeStatusNoRegisterQuote:
		return "Unknown"
	case TradeStatusOpenBid:
		return "Open Bid"
	case TradeStatusMorningClosing:
		return "Morning Break"
	case TradeStatusTrading, TradeStatusUSTrading, TradeStatusUSAfterMarketClean:
		return "Trading"
	case TradeStatusNoonClosing:
		return "Mid-Day Break"
	case TradeStatusCloseBid:
		return "Close Bid"
	case TradeStatusClosing, TradeStatusClean, TradeStatusHalfClosing, TradeStatusUSClosing, TradeStatusUSPrevMarketClean:
		return "Closed"
	case TradeStatusDarkWait:
		return "Dark Wait"
	case TradeStatusDarkTrading:
		return "Dark Trading"
	case TradeStatusDarkClosing:
		return "Closing"
	case TradeStatusAfterFix:
		return "After Fix"
	case TradeStatusNotOpened:
		return "Not Open"
	case TradeStatusRealtimeQuote:
		return "Temporary Break"
	case TradeStatusUSPrev, TradeStatusUSClean:
		return "Pre-Market"
	case TradeStatusUSAfter:
		return "Post-Market"
	case TradeStatusUSStop, TradeStatusStop:
		return "Stop"
	case TradeStatusUSNight:
		return "Overnight"
	case TradeStatusRefresh:
		return "Refresh"
	case TradeStatusDelist:
		return "Delist"
	case TradeStatusPrepare:
		return "Prepare"
	case TradeStatusCodeChange:
		return "Code Change"
	case TradeStatusWillOpen:
		return "Will Open"
	case TradeStatusCommonSuspend:
		return "Common Suspend"
	case TradeStatusExpire:
		return "Expire"
	case TradeStatusNoQuote:
		return "No Quote"
	case TradeStatusUnited:
		return "Not Listed"
	case TradeStatusTradingHalt:
		return "Terminated"
	case TradeStatusWaitListing:
		return "Wait Listing"
	case TradeStatusFuse:
		return "Fuse"
	default:
		return "Unknown"
	}
}

// IsUSMarket reports whether this is a US market status.
func (s TradeStatus) IsUSMarket() bool {
	return s.Code() >= 200 && s.Code() < 300
}

// IsUSPrePost reports whether this is a US pre/post-market status.
func (s TradeStatus) IsUSPrePost() bool {
	return s.IsUSPreMarket() || s.IsUSPostMarket()
}

// IsUSNight reports whether this is a US overnight status.
func (s TradeStatus) IsUSNight() bool {
	return s == TradeStatusUSNight
}

// IsUSClosing reports whether this is a US closed status.
func (s TradeStatus) IsUSClosing() bool {
	return s == TradeStatusUSClosing || s == TradeStatusUSPrevMarketClean
}

// IsClosing reports whether this is a closed status.
func (s TradeStatus) IsClosing() bool {
	return s == TradeStatusUSClosing ||
		s == TradeStatusUSPrevMarketClean ||
		s == TradeStatusClosing ||
		s == TradeStatusHalfClosing
}

// IsUSPreMarket reports whether this is a US pre-market status.
func (s TradeStatus) IsUSPreMarket() bool {
	return s == TradeStatusUSPrev || s == TradeStatusUSClean
}

// IsUSPostMarket reports whether this is a US post-market status.
func (s TradeStatus) IsUSPostMarket() bool {
	return s == TradeStatusUSAfter
}

// IsTrading reports whether this is a trading status.
func (s TradeStatus) IsTrading() bool {
	return s == TradeStatusTrading ||
		s == TradeStatusUSTrading ||
		s == TradeStatusUSAfterMarketClean
}

// IsDark reports whether this is a dark-pool status.
func (s TradeStatus) IsDark() bool {
	return s == TradeStatusDarkWait ||
		s == TradeStatusDarkTrading ||
		s == TradeStatusDarkClosing
}

// AllowTrading reports whether this status allows trading.
func (s TradeStatus) AllowTrading() bool {
	return s == TradeStatusOpenBid ||
		s == TradeStatusTrading ||
		s == TradeStatusCloseBid ||
		s == TradeStatusNotOpened ||
		s == TradeStatusNoonClosing ||
		s == TradeStatusUSTrading ||
		s == TradeStatusUSAfterMarketClean
}

// Normalize maps quote-engine aliases to their display-equivalent status.
func (s TradeStatus) Normalize() TradeStatus {
	switch s {
	case TradeStatusClean:
		return TradeStatusClosing
	case TradeStatusUSPrevMarketClean:
		return TradeStatusUSClosing
	case TradeStatusUSClean:
		return TradeStatusUSPrev
	case TradeStatusUSAfterMarketClean:
		return TradeStatusUSTrading
	default:
		return s
	}
}

// IsSpecial reports whether this is a special non-regular status.
func (s TradeStatus) IsSpecial() bool {
	return s.Code() < 100 || s == TradeStatusUSStop || s.Code() >= 1000
}
