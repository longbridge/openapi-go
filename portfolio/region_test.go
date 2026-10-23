package portfolio

import (
	"context"
	"errors"
	"testing"

	httplib "github.com/longbridge/openapi-go/http"
)

func TestExchangeRateRegionRestriction(t *testing.T) {
	cli, err := httplib.New(
		httplib.WithURL("http://127.0.0.1:1"),
		httplib.WithAppKey("us_key"),
		httplib.WithAppSecret("us_secret"),
		httplib.WithAccessToken("us_token"),
	)
	if err != nil {
		t.Fatalf("http.New: %v", err)
	}
	pctx := &PortfolioContext{httpClient: cli}
	var rre *httplib.RegionRestrictedError
	if _, err := pctx.ExchangeRate(context.Background()); !errors.As(err, &rre) {
		t.Errorf("ExchangeRate with US token: got %v, want RegionRestrictedError", err)
	}
}
