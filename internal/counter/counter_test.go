package counter

import "testing"

func TestIDToSymbol(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"ST/US/TSLA", "TSLA.US"},
		{"ST/HK/700", "700.HK"},
		{"ST/SH/600519", "600519.SH"},
		{"ST/SZ/000001", "000001.SZ"},
		{"IX/HK/HSI", "HSI.HK"},
		// non-slash formats returned unchanged
		{"AAPL.US", "AAPL.US"},
		{"", ""},
	}
	for _, c := range cases {
		got := IDToSymbol(c.in)
		if got != c.want {
			t.Errorf("IDToSymbol(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
