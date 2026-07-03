package trade

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/longbridge/openapi-go/internal/counter"
)

func convertUSOrder(o usRawOrder) USOrder {
	return USOrder{
		OrderID:           o.ID,
		AAID:              o.AAID,
		AccountChannel:    o.AccountChannel,
		Action:            o.Action,
		Symbol:            counter.IDToSymbol(o.CounterID),
		UnderlyingSymbol:  counter.IDToSymbol(o.UnderlyingCounterID),
		Code:              o.Code,
		Name:              o.Name,
		SecurityType:      o.SecurityType,
		Currency:          o.Currency,
		TradeCurrency:     o.TradeCurrency,
		OrderType:         o.OrderType,
		Status:            o.Status,
		Price:             o.Price,
		Quantity:          o.Quantity,
		ExecutedQty:       o.ExecutedQty,
		ExecutedPrice:     o.ExecutedPrice,
		ExecutedAmount:    o.ExecutedAmount,
		OperateDirection:  o.OperateDirection,
		TimeInForce:       o.TimeInForce,
		GTD:               o.GTD,
		SubmittedAt:       unixStrToTime(o.SubmittedAt),
		UpdatedAt:         unixStrToTime(o.UpdatedAt),
		Msg:               o.Msg,
		Report:            o.Report,
		ContractDirection: o.ContractDirection,
		ContractDueDate:   o.ContractDueDate,
		StrikePrice:       o.StrikePrice,
		TailingAmount:     o.TailingAmount,
		TailingPercent:    o.TailingPercent,
		TriggerPrice:      o.TriggerPrice,
		TriggerStatus:     o.TriggerStatus,
		TriggerAt:         o.TriggerAt,
		TriggerExchange:   o.TriggerExchange,
		TriggerLastDone:   o.TriggerLastDone,
		TriggerCount:      o.TriggerCount,
		LotSize:           o.LotSize,
		LimitOffset:       o.LimitOffset,
		LimitDepthLevel:   o.LimitDepthLevel,
		MarketPrice:       o.MarketPrice,
		LastDone:          o.LastDone,
		OrgID:             o.OrgID,
		Tag:               o.Tag,
		ForceOnlyRTH:      o.ForceOnlyRTH,
		DeductionsStatus:  o.DeductionsStatus,
		FreeStatus:        o.FreeStatus,
		Trend:             o.Trend,
	}
}

func unixStrToTime(s string) time.Time {
	if s == "" || s == "0" {
		return time.Time{}
	}
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(n, 0).UTC()
	}
	return time.Time{}
}

// QueryUSOrders queries the paginated US order list.
//
// Path: POST /v1/orders/query
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) QueryUSOrders(ctx context.Context, req *QueryUSOrdersRequest) (*QueryUSOrdersResponse, error) {
	if err := c.opts.httpClient.CheckRegion("/v1/orders/query", "US"); err != nil {
		return nil, err
	}
	var raw usRawQueryUSOrdersResponse
	if err := c.opts.httpClient.Post(ctx, "/v1/orders/query", req, &raw); err != nil {
		return nil, err
	}
	orders := make([]USOrder, 0, len(raw.Orders))
	for _, o := range raw.Orders {
		orders = append(orders, convertUSOrder(o))
	}
	return &QueryUSOrdersResponse{Orders: orders, TotalCount: raw.TotalCount}, nil
}

// USOrderDetail returns the detail for a single US order.
//
// Path: GET /v1/orders/{order_id}
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USOrderDetail(ctx context.Context, orderID string) (*USOrderDetailResponse, error) {
	path := fmt.Sprintf("/v1/orders/%s", orderID)
	if err := c.opts.httpClient.CheckRegion(path, "US"); err != nil {
		return nil, err
	}
	var out USOrderDetailResponse
	if err := c.opts.httpClient.Get(ctx, path, nil, &out); err != nil {
		return nil, err
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
			AssetType:    e.AssetType,
			AverageCost:  e.AverageCost,
			Symbol:       counter.IDToSymbol(e.CounterID),
			Currency:     e.Currency,
			IndustryName: e.IndustryName,
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
