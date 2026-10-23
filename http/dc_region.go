package http

import (
	"fmt"
	"strings"
)

// dcRegionHeader is the HTTP header that selects the data center serving a request.
// An absent header is treated as the AP (Asia-Pacific) data center by the API gateway.
const dcRegionHeader = "x-dc-region"

// dcRegion identifies which data center the API gateway should route a request to.
// Independent of the CN/global host selection: that picks the *.longbridge.cn vs
// *.longbridge.com host, while this selects which data center (us/ap) the gateway
// sources data from.
type dcRegion int

const (
	dcRegionAp dcRegion = iota // Asia-Pacific (gateway default)
	dcRegionUs                 // US data center
)

// dcRegionFromCredential derives the DC region from a single credential's prefix.
//
// Longbridge credentials are prefixed with their data center: "us_…" for the US
// data center, "ap_…" for Asia-Pacific. A "us_" prefix maps to dcRegionUs;
// everything else — including "ap_"-prefixed and unprefixed — maps to dcRegionAp,
// matching the gateway default. A leading "Bearer " is tolerated.
func dcRegionFromCredential(credential string) dcRegion {
	credential = strings.TrimPrefix(credential, "Bearer ")
	if strings.HasPrefix(credential, "us_") {
		return dcRegionUs
	}
	return dcRegionAp
}

// dcRegionFromCredentials returns dcRegionUs if any credential carries the "us_" prefix.
// Used for API-key auth where app_key, app_secret, and access_token all carry the prefix.
func dcRegionFromCredentials(credentials ...string) dcRegion {
	for _, c := range credentials {
		if dcRegionFromCredential(c) == dcRegionUs {
			return dcRegionUs
		}
	}
	return dcRegionAp
}

// asStr returns the x-dc-region header value for this region ("us" or "ap").
func (r dcRegion) asStr() string {
	if r == dcRegionUs {
		return "us"
	}
	return "ap"
}

// display returns the human-readable uppercase name of the region ("US" or "AP").
func (r dcRegion) display() string {
	if r == dcRegionUs {
		return "US"
	}
	return "AP"
}

// allows reports whether this region satisfies a restriction to required.
func (r dcRegion) allows(required dcRegion) bool {
	return r == required
}

// RegionRestrictedError is returned when an API is called from the wrong data center.
// It mirrors the Rust SDK's HttpClientError::DcRegionRestricted.
type RegionRestrictedError struct {
	// Path is the API path or WebSocket command that was called.
	Path string
	// Required is the data center the API requires ("US" or "AP").
	Required string
	// Current is the data center of the current session ("US" or "AP").
	Current string
}

// Error returns the region message followed, when one is known, by a hint
// naming the API to call instead, so callers (including LLM agents) can recover
// without parsing the path.
func (e *RegionRestrictedError) Error() string {
	msg := fmt.Sprintf(
		"this API (%s) is only available in the %s data center and is not supported for your %s-region account",
		e.Path, e.Required, e.Current,
	)
	if hint := regionAlternative(e.Path, e.Required); hint != "" {
		msg += "; " + hint
	}
	return msg
}

// Hints for paths that have no counterpart in the other data center.
const (
	noUSEquivalent = "there is no US equivalent"
	noAPEquivalent = "there is no AP equivalent"
)

// regionAlternative returns the recovery hint for a region-restricted API
// path, or "" when there is none. required is compared case-insensitively.
//
// The hint text matches the Rust SDK, whose message is shared by every
// binding and by MCP servers with different method and tool names, so the
// alternative is named by HTTP endpoint rather than by method.
func regionAlternative(path, required string) string {
	if strings.EqualFold(required, "US") {
		return usOnlyAlternative(path)
	}
	return apOnlyAlternative(path)
}

// apOnlyAlternative is the hint for a US-region account calling an AP-only path.
func apOnlyAlternative(path string) string {
	switch path {
	case "/v1/asset/account":
		return "use the US endpoint GET /v1/us/assets/cash-and-buy-power or GET /v1/us/assets/overview (TradeContext US methods)"
	case "/v1/asset/cashflow":
		return "use the US endpoint POST /v1/us/assets/transaction-history (TradeContext US methods)"
	case "/v1/asset/stock":
		return "use the US endpoint GET /v1/us/assets/overview (TradeContext US methods)"
	case "/v1/asset/exchange_rates":
		return "the US equivalent is on TradeContext (GET /v1/us/assets/exchange-rates)"
	case "/v1/trade/order":
		return "use the US endpoint GET /v1/us/orders/{order_id} (TradeContext US methods)"
	case "/v1/risk/margin-ratio":
		return "there is no US equivalent; account-level margin is at GET /v1/us/assets/cash-and-buy-power (margin_status)"
	case "/v1/asset/fund", "quote/brokers", "/v1/quote/operatings":
		return noUSEquivalent
	}
	if strings.HasPrefix(path, "/v1/quote/broker-holding") || strings.HasPrefix(path, "/v1/dailycoins/") {
		return noUSEquivalent
	}
	return ""
}

// usOnlyAlternative is the hint for an AP-region account calling a US-only path.
func usOnlyAlternative(path string) string {
	switch path {
	case "/v1/us/orders/query":
		return "use the AP endpoint GET /v1/trade/order/today or GET /v1/trade/order/history"
	case "/v1/us/assets/overview":
		return "use the AP endpoints GET /v1/asset/account and GET /v1/asset/stock"
	case "/v1/us/assets/transaction-history":
		return "use the AP endpoint GET /v1/asset/cashflow"
	case "/v1/us/assets/exchange-rates":
		return "use the AP endpoint GET /v1/asset/exchange_rates"
	case "/v1/us/assets/cash-and-buy-power":
		return "use the AP endpoint GET /v1/asset/account"
	case "/v1/us/assets/available-holding":
		return "use the AP endpoint GET /v1/asset/stock (available_quantity)"
	}
	switch {
	case strings.HasPrefix(path, "/v1/us/orders/"):
		return "use the AP endpoint GET /v1/trade/order"
	case strings.HasPrefix(path, "/v1/us/"):
		return noAPEquivalent
	}
	return ""
}

// stripBearerPrefix strips any leading "Bearer " from a credential.
//
// Region prefixes (hk_m_, us_m_, ap_m_, …) are routing metadata consumed by
// dcRegionFromCredential to derive the x-dc-region header. The gateway
// accepts the full prefixed token and routes by the header, so no region
// prefix is stripped — only "Bearer " is removed.
func stripBearerPrefix(credential string) string {
	return strings.TrimPrefix(credential, "Bearer ")
}
