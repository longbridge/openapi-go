package trade

// QueryUSOrdersRequest is the request body for QueryUSOrders.
type QueryUSOrdersRequest struct {
	AccountChannel string   `json:"account_channel"`
	Action         int32    `json:"action"`
	StartAt        float64  `json:"start_at"`
	EndAt          float64  `json:"end_at"`
	CounterIDs     []string `json:"counter_ids"`
	SecurityTypes  []string `json:"security_types"`
	QueryType      int32    `json:"query_type"`
	Page           int32    `json:"page"`
	Limit          int32    `json:"limit"`
	QueryVersion   float64  `json:"query_version"`
}

// QueryUSOrdersResponse holds the paged list of US orders.
type QueryUSOrdersResponse struct {
	Orders     []map[string]interface{} `json:"orders"`
	TotalCount int32                    `json:"total_count"`
}

// AttachedOrder is a take-profit or stop-loss sub-order attached to a US order.
type AttachedOrder struct {
	OrderID      string `json:"order_id"`
	Type         string `json:"type"`
	Side         string `json:"side"`
	Price        string `json:"price"`
	TrailAmount  string `json:"trail_amount"`
	TrailPercent string `json:"trail_percent"`
	Status       string `json:"status"`
}

// USOrderDetailResponse is the response for USOrderDetail.
// The raw order detail fields are passed through as-is; attached_orders is
// populated only when isAttached=true.
type USOrderDetailResponse struct {
	Raw            map[string]interface{} `json:"-"`
	AttachedOrders []AttachedOrder        `json:"attached_orders,omitempty"`
}

// ── USAssetOverview ────────────────────────────────────────────────────────

// USCashEntry is one currency cash entry in USAssetOverview.
type USCashEntry struct {
	Currency      string `json:"currency"`
	FrozenBuyCash string `json:"frozen_buy_cash"`
	Outstanding   string `json:"outstanding"`
	SettledCash   string `json:"settled_cash"`
	TotalAmount   string `json:"total_amount"`
	TotalCash     string `json:"total_cash"`
}

// USCryptoEntry is one cryptocurrency holding in USAssetOverview.
type USCryptoEntry struct {
	AssetType         string `json:"asset_type"`
	AverageCost       string `json:"average_cost"`
	CounterID         string `json:"counter_id"`
	Currency          string `json:"currency"`
	IndustryCounterID string `json:"industry_counter_id"`
	IndustryName      string `json:"industry_name"`
	// Additional fields passed through without type assertion.
	Extra map[string]interface{} `json:"-"`
}

// USAssetOverview is the US account asset snapshot.
// Field names match the actual API response from /v1/us/assets/overview.
type USAssetOverview struct {
	AccountType    string        `json:"account_type"`
	AssetTimestamp string        `json:"asset_timestamp"`
	CashBuyPower   string        `json:"cash_buy_power"`
	CashList       []USCashEntry `json:"cash_list"`
	CryptoList     []USCryptoEntry `json:"crypto_list"`
	// The full response may contain additional fields (stock positions, option
	// positions, etc.) that are preserved here for forward compatibility.
	Extra map[string]interface{} `json:"-"`
}

// ── USRealizedPL ───────────────────────────────────────────────────────────

// USRealizedPLMetric is one time-period metric within a USRealizedPLEntry.
type USRealizedPLMetric struct {
	Amount string `json:"amount"`
	Period int32  `json:"period"`
	Rate   string `json:"rate"`
}

// USRealizedPLEntry is one asset-category entry in USRealizedPL.
// Category: 0=all, 1=stock, 2=option, 3=crypto (server-defined values).
type USRealizedPLEntry struct {
	Category int32                `json:"category"`
	Currency string               `json:"currency"`
	Metrics  []USRealizedPLMetric `json:"metrics"`
}

// USRealizedPL is the response for GET /v1/us/assets/pl/realized.
type USRealizedPL struct {
	RealizedPLList []USRealizedPLEntry `json:"realized_pl_list"`
}
