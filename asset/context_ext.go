package asset

import (
	"context"
	"encoding/json"
)

// ExchangeRates returns all exchange rates.
//
// Reference: GET /v1/asset/exchange_rates
func (c *AssetContext) ExchangeRates(ctx context.Context) (json.RawMessage, error) {
	var resp json.RawMessage
	err := c.httpClient.Get(ctx, "/v1/asset/exchange_rates", nil, &resp)
	return resp, err
}
