package http

import "testing"

func TestDcRegionFromCredential(t *testing.T) {
	tests := []struct {
		credential string
		want       dcRegion
	}{
		{"us_eyJabc", dcRegionUs},
		{"ap_eyJabc", dcRegionAp},
		{"hk_eyJabc", dcRegionAp},
		{"eyJabc", dcRegionAp},
		{"", dcRegionAp},
		{"Bearer us_eyJabc", dcRegionUs},
		{"Bearer ap_eyJabc", dcRegionAp},
		{"Bearer eyJabc", dcRegionAp},
	}
	for _, tc := range tests {
		got := dcRegionFromCredential(tc.credential)
		if got != tc.want {
			t.Errorf("dcRegionFromCredential(%q) = %v, want %v", tc.credential, got, tc.want)
		}
	}
}

func TestDcRegionFromCredentials(t *testing.T) {
	if dcRegionFromCredentials("ap_key", "us_secret", "ap_token") != dcRegionUs {
		t.Error("expected Us when any credential is us_")
	}
	if dcRegionFromCredentials("ap_key", "ap_secret", "ap_token") != dcRegionAp {
		t.Error("expected Ap when all credentials are ap_")
	}
	if dcRegionFromCredentials() != dcRegionAp {
		t.Error("expected Ap for empty credentials")
	}
}

func TestDcRegionAsStr(t *testing.T) {
	if dcRegionUs.asStr() != "us" {
		t.Error("Us.asStr() should be 'us'")
	}
	if dcRegionAp.asStr() != "ap" {
		t.Error("Ap.asStr() should be 'ap'")
	}
}

func TestStripBearerPrefix(t *testing.T) {
	// Region prefixes are kept as-is; only "Bearer " is stripped.
	tests := []struct {
		input string
		want  string
	}{
		{"us_m_eyJabc", "us_m_eyJabc"},
		{"hk_m_eyJabc", "hk_m_eyJabc"},
		{"ap_m_eyJabc", "ap_m_eyJabc"},
		{"eyJabc", "eyJabc"},
		{"Bearer us_m_eyJabc", "us_m_eyJabc"},
		{"Bearer hk_m_eyJabc", "hk_m_eyJabc"},
		{"Bearer eyJabc", "eyJabc"},
	}
	for _, tc := range tests {
		got := stripBearerPrefix(tc.input)
		if got != tc.want {
			t.Errorf("stripBearerPrefix(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestRegionRestrictedErrorHint(t *testing.T) {
	const apPrefix = "is only available in the AP data center and is not supported for your US-region account"
	tests := []struct {
		path, required, current, want string
	}{
		{"/v1/asset/account", "AP", "US", "this API (/v1/asset/account) " + apPrefix + "; use USCashAndBuyPower or USAssetOverview instead"},
		{"/v1/asset/cashflow", "AP", "US", "this API (/v1/asset/cashflow) " + apPrefix + "; use USTransactionHistory instead"},
		{"/v1/asset/stock", "AP", "US", "this API (/v1/asset/stock) " + apPrefix + "; use USAssetOverview instead"},
		{"/v1/asset/exchange_rates", "AP", "US", "this API (/v1/asset/exchange_rates) " + apPrefix + "; use USExchangeRates instead"},
		{"/v1/trade/order", "AP", "US", "this API (/v1/trade/order) " + apPrefix + "; use USOrderDetail instead"},
		{"/v1/asset/fund", "AP", "US", "this API (/v1/asset/fund) " + apPrefix + "; there is no US equivalent"},
		{"/v1/risk/margin-ratio", "AP", "US", "this API (/v1/risk/margin-ratio) " + apPrefix + "; there is no US equivalent; see MarginStatus in USCashAndBuyPower instead"},
		{"/v1/us/assets/overview", "US", "AP", "this API (/v1/us/assets/overview) is only available in the US data center and is not supported for your AP-region account; this API is for US data-center accounts only"},
		{"quote/brokers", "AP", "US", "this API (quote/brokers) " + apPrefix},
	}
	for _, tt := range tests {
		err := &RegionRestrictedError{Path: tt.path, Required: tt.required, Current: tt.current}
		if got := err.Error(); got != tt.want {
			t.Errorf("Error() for %s:\n got  %q\n want %q", tt.path, got, tt.want)
		}
	}
}
