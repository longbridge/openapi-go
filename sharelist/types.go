package sharelist

import (
	"time"

	"github.com/shopspring/decimal"
)

// SharelistType represents the kind of sharelist.
type SharelistType int32

const (
	// SharelistTypeRegular is a regular user sharelist.
	SharelistTypeRegular SharelistType = 0
	// SharelistTypeOfficial is an officially curated sharelist.
	SharelistTypeOfficial SharelistType = 3
	// SharelistTypeIndustry is an industry-themed sharelist.
	SharelistTypeIndustry SharelistType = 4
)

// SharelistList is the result of the List and Popular methods.
type SharelistList struct {
	// Sharelists contains the user's own sharelists (and, for Popular, the
	// trending sharelists).
	Sharelists []SharelistInfo
	// SubscribedSharelists contains sharelists the user is subscribed to.
	// This field may be empty when returned by Popular.
	SubscribedSharelists []SharelistInfo
	// TailMark is the pagination cursor for the subscribed list.
	TailMark string
}

// SharelistDetail is the result of the Detail method.
type SharelistDetail struct {
	// Sharelist holds the sharelist metadata and constituents.
	Sharelist SharelistInfo
	// Scopes holds subscription/ownership status for the authenticated user.
	Scopes SharelistScopes
}

// SharelistInfo holds metadata and constituent stocks for a sharelist.
type SharelistInfo struct {
	// ID is the sharelist identifier.
	ID int64
	// Name is the display name of the sharelist.
	Name string
	// Description is a short description.
	Description string
	// Cover is the URL of the cover image.
	Cover string
	// SubscribersCount is the number of subscribers.
	SubscribersCount int64
	// CreatedAt is when the sharelist was created.
	CreatedAt time.Time
	// EditedAt is when the constituent list was last edited.
	EditedAt time.Time
	// ThisYearChg is the YTD change percentage; nil when not available.
	ThisYearChg *decimal.Decimal
	// Stocks holds the constituent securities.
	Stocks []SharelistStock
	// Subscribed is true when the authenticated user is subscribed.
	Subscribed bool
	// Chg is the day change percentage; nil when not available.
	Chg *decimal.Decimal
	// SharelistType classifies the sharelist (regular / official / industry).
	SharelistType SharelistType
	// IndustryCode is populated for industry-type sharelists.
	IndustryCode string
}

// SharelistStock describes a security within a sharelist.
type SharelistStock struct {
	// Symbol is the security identifier, e.g. "TSLA.US" or "700.HK".
	// It is converted from the wire-level counter_id field.
	Symbol string
	// Name is the display name of the security.
	Name string
	// Market is the exchange market code, e.g. "HK" or "US".
	Market string
	// Code is the ticker code.
	Code string
	// Intro is a brief description.
	Intro string
	// UnreadChangeLogCategory is the unread change log category.
	UnreadChangeLogCategory string
	// Change is the day change percentage; nil when not available.
	Change *decimal.Decimal
	// LastDone is the latest price; nil when not available.
	LastDone *decimal.Decimal
	// TradeStatus is the trade status code; nil when not available.
	TradeStatus *int32
	// Latency indicates a delayed quote when true; nil when not available.
	Latency *bool
}

// SharelistScopes holds subscription and ownership status for the current user.
type SharelistScopes struct {
	// Subscription is true when the authenticated user is subscribed.
	Subscription bool
	// IsSelf is true when the authenticated user is the creator.
	IsSelf bool
}
