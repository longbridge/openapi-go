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
	const usPrefix = "is only available in the US data center and is not supported for your AP-region account"
	ap := func(path, hint string) [4]string {
		want := "this API (" + path + ") " + apPrefix
		if hint != "" {
			want += "; " + hint
		}
		return [4]string{path, "AP", "US", want}
	}
	us := func(path, hint string) [4]string {
		want := "this API (" + path + ") " + usPrefix
		if hint != "" {
			want += "; " + hint
		}
		return [4]string{path, "US", "AP", want}
	}
	tests := [][4]string{
		ap("/v1/asset/account", "use the US endpoint GET /v1/us/assets/cash-and-buy-power or GET /v1/us/assets/overview (TradeContext US methods)"),
		ap("/v1/asset/cashflow", "use the US endpoint POST /v1/us/assets/transaction-history (TradeContext US methods)"),
		ap("/v1/asset/stock", "use the US endpoint GET /v1/us/assets/overview (TradeContext US methods)"),
		ap("/v1/asset/exchange_rates", "use the US endpoint GET /v1/us/assets/exchange-rates (TradeContext US methods)"),
		ap("/v1/trade/order", "use the US endpoint GET /v1/us/orders/{order_id} (TradeContext US methods)"),
		ap("/v1/risk/margin-ratio", "there is no US equivalent; account-level margin is at GET /v1/us/assets/cash-and-buy-power (margin_status)"),
		ap("/v3/trade/execution/all", "use GET /v1/trade/execution/history"),
		ap("/v1/asset/fund", "there is no US equivalent"),
		ap("quote/brokers", "there is no US equivalent"),
		ap("/v1/quote/operatings", "there is no US equivalent"),
		ap("/v1/quote/broker-holding", "there is no US equivalent"),
		ap("/v1/quote/broker-holding/detail", "there is no US equivalent"),
		ap("/v1/quote/broker-holding/daily", "there is no US equivalent"),
		ap("/v1/dailycoins/query", "there is no US equivalent"),
		ap("/v1/dailycoins/toggle", "there is no US equivalent"),
		ap("/v1/unknown", ""),
		us("/v1/us/orders/123", "use the AP endpoint GET /v1/trade/order"),
		us("/v1/us/orders/query", "use the AP endpoint GET /v1/trade/order/today or GET /v1/trade/order/history"),
		us("/v1/us/assets/overview", "use the AP endpoints GET /v1/asset/account and GET /v1/asset/stock"),
		us("/v1/us/assets/transaction-history", "use the AP endpoint GET /v1/asset/cashflow"),
		us("/v1/us/assets/exchange-rates", "use the AP endpoint GET /v1/asset/exchange_rates"),
		us("/v1/us/assets/cash-and-buy-power", "use the AP endpoint GET /v1/asset/account"),
		us("/v1/us/assets/available-holding", "use the AP endpoint GET /v1/asset/stock (available_quantity)"),
		us("/v1/us/assets/pl/realized", "there is no AP equivalent"),
		us("/v1/us/gemini/crypto-overview", "there is no AP equivalent"),
		us("/v1/us/stock-info/etf-files", "there is no AP equivalent"),
		us("/v1/us/stock-info/company-overview", "use the AP endpoint GET /v1/quote/comp-overview"),
		us("/v1/us/stock-info/valuation-overview", "use the AP endpoint GET /v1/quote/valuation"),
		us("/v1/us/stock-info/finn-overview", "use the AP endpoint GET /v1/quote/financial-reports"),
		us("/v1/us/quote/financials/statements", "use the AP endpoint GET /v1/quote/financial-reports"),
		us("/v1/us/stock-info/fin-keyfactor", "use the AP endpoint GET /v1/quote/financials/earnings-snapshot"),
		us("/v1/us/stock-info/fin-consensus", "use the AP endpoint GET /v1/quote/financial-consensus-detail"),
		us("/v1/us/stock-info/etf-dividend-info", "use the AP endpoint GET /v1/quote/dividends"),
		us("/v1/us/stock-info/company-dividends", "use the AP endpoint GET /v1/quote/dividends"),
		us("/v1/unknown", ""),
	}
	for _, tt := range tests {
		err := &RegionRestrictedError{Path: tt[0], Required: tt[1], Current: tt[2]}
		if got := err.Error(); got != tt[3] {
			t.Errorf("Error() for %s:\n got  %q\n want %q", tt[0], got, tt[3])
		}
	}
}

func TestRegionAlternativeCaseInsensitive(t *testing.T) {
	for _, required := range []string{"US", "us", "Us"} {
		if got := regionAlternative("/v1/us/assets/cash-and-buy-power", required); got != "use the AP endpoint GET /v1/asset/account" {
			t.Errorf("regionAlternative(required=%q) = %q", required, got)
		}
	}
	if err := checkRegion("/v1/us/assets/overview", "us", dcRegionUs); err != nil {
		t.Errorf("checkRegion with lowercase \"us\" on a US session = %v, want nil", err)
	}
	if err := checkRegion("/v1/asset/account", "ap", dcRegionAp); err != nil {
		t.Errorf("checkRegion with lowercase \"ap\" on an AP session = %v, want nil", err)
	}
}
