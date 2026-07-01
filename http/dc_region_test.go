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

func TestStripRegionPrefix(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"us_eyJabc", "eyJabc"},
		{"ap_eyJabc", "eyJabc"},
		{"hk_eyJabc", "hk_eyJabc"}, // unknown prefix — pass through
		{"eyJabc", "eyJabc"},
		{"Bearer us_eyJabc", "eyJabc"},
		{"Bearer ap_eyJabc", "eyJabc"},
		{"Bearer eyJabc", "eyJabc"},
	}
	for _, tc := range tests {
		got := stripRegionPrefix(tc.input)
		if got != tc.want {
			t.Errorf("stripRegionPrefix(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
