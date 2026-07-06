package http

import "strings"

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

// stripRegionPrefix strips any leading "Bearer " from a credential.
//
// Region prefixes (hk_m_, us_m_, ap_m_, …) are routing metadata consumed by
// dcRegionFromCredential to derive the x-dc-region header. The gateway
// accepts the full prefixed token and routes by the header, so no region
// prefix is stripped — only "Bearer " is removed.
func stripRegionPrefix(credential string) string {
	return strings.TrimPrefix(credential, "Bearer ")
}
