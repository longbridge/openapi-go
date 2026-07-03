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

func (e *RegionRestrictedError) Error() string {
	return fmt.Sprintf(
		"this API (%s) is only available in the %s data center and is not supported for your %s-region account",
		e.Path, e.Required, e.Current,
	)
}

// stripRegionPrefix strips the routing prefix from a credential so only the
// bare JWT is sent in the Authorization header.
//
// Access tokens are prefixed with their data-center identifier
// (e.g. "us_m_eyJ…", "hk_m_eyJ…"). The prefix is routing metadata consumed
// by dcRegionFromCredential; the gateway validates only the bare JWT (starting
// with "eyJ"). Sending the full prefixed string causes JWT header decode
// failure on the server side.
//
// Stripping order:
//  1. Remove any leading "Bearer " OAuth wrapper.
//  2. Remove any region prefix before the JWT start ("eyJ").
//
// App keys (hex strings, no "eyJ") are returned unchanged.
func stripRegionPrefix(credential string) string {
	credential = strings.TrimPrefix(credential, "Bearer ")
	if idx := strings.Index(credential, "eyJ"); idx > 0 {
		credential = credential[idx:]
	}
	return credential
}
