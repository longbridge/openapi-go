package counter

import "testing"

func TestValidateSymbol(t *testing.T) {
	tests := []struct {
		symbol  string
		wantErr bool
	}{
		{"AAPL.US", false},
		{"SPY.US", false},
		{"BTCUSD.BKKT", false},
		{"AXTI", true},        // no dot
		{"", true},            // empty
		{".US", true},         // empty code
		{"AAPL.", true},       // empty market
	}
	for _, tt := range tests {
		err := ValidateSymbol(tt.symbol)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateSymbol(%q) error = %v, wantErr %v", tt.symbol, err, tt.wantErr)
		}
	}
}
