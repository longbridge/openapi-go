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
		SubmittedAt:       o.SubmittedAt.Time(),
		UpdatedAt:         o.UpdatedAt.Time(),
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
		MonitorPrice:             o.MonitorPrice,
		SubmittedAmount:          o.SubmittedAmount,
		PlatformDeductionsStatus: o.PlatformDeductionsStatus,
		PloyID:                   o.PloyID,
		PloyType:                 o.PloyType,
		TickerSize:               o.TickerSize,
		CurrentMillisecond:       o.CurrentMillisecond,
		OrgID:                    o.OrgID,
		Tag:                      o.Tag,
		ForceOnlyRTH:      o.ForceOnlyRTH,
		DeductionsStatus:  o.DeductionsStatus,
		FreeStatus:        o.FreeStatus,
		Trend:             o.Trend,
	}
}

func convertStockList(entries []USStockEntry) []USStockEntry {
	for i := range entries {
		entries[i].FullSymbol = counter.IDToSymbol(entries[i].CounterID)
	}
	return entries
}

// QueryUSOrders queries the US order list.
//
// Path: POST /v1/us/orders/query
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) QueryUSOrders(ctx context.Context, req *GetUSHistoryOrders) (*QueryUSOrdersResponse, error) {
	if err := c.opts.httpClient.CheckRegion("/v1/us/orders/query", "US"); err != nil {
		return nil, err
	}
	if req == nil {
		req = &GetUSHistoryOrders{}
	}
	now := time.Now().Unix()

	// Build action from Side
	action := int32(0)
	switch req.Side {
	case OrderSideBuy:
		action = 1
	case OrderSideSell:
		action = 2
	}

	// Build counter_ids from Symbol
	counterIDs := []string{}
	if req.Symbol != "" {
		counterIDs = append(counterIDs, counter.SymbolToID(req.Symbol))
	}

	// Apply defaults
	page, limit := req.Page, req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	startAt := float64(req.StartAt)
	if startAt == 0 {
		startAt = float64(now - 90*24*3600)
	}
	endAt := float64(req.EndAt)
	if endAt == 0 {
		endAt = float64(now)
	}

	body := map[string]interface{}{
		"account_channel": "",
		"action":          action,
		"start_at":        startAt,
		"end_at":          endAt,
		"counter_ids":     counterIDs,
		"security_types":  []string{},
		"query_type":      req.QueryType,
		"page":            page,
		"limit":           limit,
		"query_version":   float64(now),
	}

	var raw usRawQueryUSOrdersResponse
	if err := c.opts.httpClient.Post(ctx, "/v1/us/orders/query", body, &raw); err != nil {
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
// Path: GET /v1/us/orders/{order_id}
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USOrderDetail(ctx context.Context, orderID string) (*USOrderDetailResponse, error) {
	path := fmt.Sprintf("/v1/us/orders/%s", orderID)
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
		AccountType:       raw.AccountType,
		AssetTimestamp:    ts,
		CashBuyPower:      raw.CashBuyPower,
		OvernightBuyPower: raw.OvernightBuyPower,
		Currency:          raw.Currency,
		CashList:          raw.CashList,
		StockList:         convertStockList(raw.StockList),
		OptionList:        raw.OptionList,
		CryptoList:        cryptoList,
		MultiLeg:          raw.MultiLeg,
	}, nil
}

// USRealizedPL returns realized profit-and-loss for the US account.
//
// Path: GET /v1/us/assets/pl/realized
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USRealizedPL(ctx context.Context, req *GetUSRealizedPL) (*USRealizedPLResponse, error) {
	if err := c.opts.httpClient.CheckRegion("/v1/us/assets/pl/realized", "US"); err != nil {
		return nil, err
	}
	if req == nil {
		req = &GetUSRealizedPL{Currency: "USD"}
	}
	if req.Currency == "" {
		req.Currency = "USD"
	}
	q := url.Values{}
	q.Set("currency", req.Currency)
	if req.Category != "" {
		q.Set("category", req.Category)
	}
	var resp USRealizedPLResponse
	if err := c.opts.httpClient.Get(ctx, "/v1/us/assets/pl/realized", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
