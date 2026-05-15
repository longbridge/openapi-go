// Package jsontypes contains raw JSON response structs for the Longbridge
// Sharelist API. These types mirror the wire format exactly; callers should
// use the public types in the parent sharelist package instead.
package jsontypes

// SharelistList is the response for the list and popular endpoints.
type SharelistList struct {
	Sharelists           []SharelistInfo `json:"sharelists"`
	SubscribedSharelists []SharelistInfo `json:"subscribed_sharelists"`
	TailMark             string          `json:"tail_mark"`
}

// SharelistDetail is the response for the detail endpoint.
type SharelistDetail struct {
	Sharelist SharelistInfo   `json:"sharelist"`
	Scopes    SharelistScopes `json:"scopes"`
}

// SharelistInfo represents a sharelist with metadata and constituent stocks.
type SharelistInfo struct {
	// ID may be returned as a string or integer by the API.
	ID interface{} `json:"id"`
	// Name of the sharelist.
	Name string `json:"name"`
	// Description of the sharelist.
	Description string `json:"description"`
	// Cover image URL.
	Cover string `json:"cover"`
	// Number of subscribers.
	SubscribersCount int64 `json:"subscribers_count"`
	// Creation time as Unix timestamp string.
	CreatedAt string `json:"created_at"`
	// Last stock edit time as Unix timestamp string.
	EditedAt string `json:"edited_at"`
	// YTD change percentage (may be "" when absent).
	ThisYearChg string `json:"this_year_chg"`
	// Creator info (kept as raw JSON to match the Rust implementation).
	Creator interface{} `json:"creator"`
	// Constituent stocks.
	Stocks []SharelistStock `json:"stocks"`
	// Whether the current user is subscribed.
	Subscribed bool `json:"subscribed"`
	// Day change percentage (may be "" when absent).
	Chg string `json:"chg"`
	// Sharelist type: 0=regular, 3=official, 4=industry.
	SharelistType int32 `json:"sharelist_type"`
	// Industry code (for industry sharelists).
	IndustryCode string `json:"industry_code"`
}

// SharelistStock describes a security in a sharelist.
type SharelistStock struct {
	// CounterID is the raw counter_id from the API (e.g. "ST/US/TSLA").
	// Use Symbol in the public type instead.
	CounterID string `json:"counter_id"`
	// Name is the security display name.
	Name string `json:"name"`
	// Market, e.g. "HK".
	Market string `json:"market"`
	// Ticker code.
	Code string `json:"code"`
	// Brief description.
	Intro string `json:"intro"`
	// Unread change log category.
	UnreadChangeLogCategory string `json:"unread_change_log_category"`
	// Day change percentage (may be "" when absent).
	Change string `json:"change"`
	// Latest price (may be "" when absent).
	LastDone string `json:"last_done"`
	// Trade status code.
	TradeStatus *int32 `json:"trade_status"`
	// Whether delayed quote.
	Latency *bool `json:"latency"`
}

// SharelistScopes holds subscription/ownership flags for the current user.
type SharelistScopes struct {
	// Subscription indicates whether the current user is subscribed.
	Subscription bool `json:"subscription"`
	// IsSelf indicates whether the current user is the creator.
	IsSelf bool `json:"self"`
}
