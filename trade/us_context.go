package trade

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/longbridge/openapi-go/internal/counter"
)

// QueryUSOrders queries the paginated US order list.
//
// Path: POST /v1/orders/query
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) QueryUSOrders(ctx context.Context, req *QueryUSOrdersRequest) (*QueryUSOrdersResponse, error) {
	if err := c.opts.httpClient.CheckRegion("/v1/orders/query", "US"); err != nil {
		return nil, err
	}
	var resp QueryUSOrdersResponse
	if err := c.opts.httpClient.Post(ctx, "/v1/orders/query", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// USOrderDetail returns the detail for a single US order.
// When isAttached is true, attached take-profit/stop-loss sub-orders are
// included in the response.
//
// Path: GET /v3/orders/{order_id}
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USOrderDetail(ctx context.Context, orderID string, isAttached bool) (*USOrderDetailResponse, error) {
	path := fmt.Sprintf("/v3/orders/%s", orderID)
	if err := c.opts.httpClient.CheckRegion(path, "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("order_id_str", orderID)
	if isAttached {
		q.Set("is_attached", "true")
	}
	var raw json.RawMessage
	if err := c.opts.httpClient.Get(ctx, path, q, &raw); err != nil {
		return nil, err
	}
	var out USOrderDetailResponse
	if err := json.Unmarshal(raw, &out.Raw); err != nil {
		return nil, err
	}
	if isAttached {
		var wrapper struct {
			AttachedOrders []AttachedOrder `json:"attached_orders"`
		}
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			return nil, err
		}
		out.AttachedOrders = wrapper.AttachedOrders
	}
	return &out, nil
}

// USAssetOverview returns the full US account asset snapshot, including stock,
// option, multi-leg, and crypto positions together with purchasing power.
//
// counter_id fields in crypto positions are converted to user-facing symbols.
// asset_timestamp is converted from Unix seconds to time.Time.
//
// Path: GET /v1/us/assets/overview
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USAssetOverview(ctx context.Context) (*USAssetOverview, error) {
	if err := c.opts.httpClient.CheckRegion("/v1/us/assets/overview", "US"); err != nil {
		return nil, err
	}
	var raw usRawAssetOverview
	if err := c.opts.httpClient.Get(ctx, "/v1/us/assets/overview", nil, &raw); err != nil {
		return nil, err
	}
	cryptoList := make([]USCryptoEntry, 0, len(raw.CryptoList))
	for _, e := range raw.CryptoList {
		cryptoList = append(cryptoList, USCryptoEntry{
			AssetType:      e.AssetType,
			AverageCost:    e.AverageCost,
			Symbol:         counter.IDToSymbol(e.CounterID),
			Currency:       e.Currency,
			IndustrySymbol: counter.IDToSymbol(e.IndustryCounterID),
			IndustryName:   e.IndustryName,
		})
	}
	var ts time.Time
	if secs, err := strconv.ParseInt(raw.AssetTimestamp, 10, 64); err == nil {
		ts = time.Unix(secs, 0).UTC()
	}
	return &USAssetOverview{
		AccountType:    raw.AccountType,
		AssetTimestamp: ts,
		CashBuyPower:   raw.CashBuyPower,
		CashList:       raw.CashList,
		CryptoList:     cryptoList,
	}, nil
}

// USRealizedPL returns realized profit-and-loss for the US account.
//
// currency is required (e.g. "USD").
// category filters by asset type: "ALL", "STOCK", "OPTION", or "CRYPTO";
// pass nil for all categories.
//
// Path: GET /v1/us/assets/pl/realized
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USRealizedPL(ctx context.Context, currency string, category *string) (*USRealizedPL, error) {
	if err := c.opts.httpClient.CheckRegion("/v1/us/assets/pl/realized", "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("currency", currency)
	if category != nil {
		q.Set("category", *category)
	}
	var resp USRealizedPL
	if err := c.opts.httpClient.Get(ctx, "/v1/us/assets/pl/realized", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
