package trade

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func convertUSOrder(o usRawOrder) USOrder {
	return USOrder{
		OrderID:                  o.ID,
		AAID:                     o.AAID,
		AccountChannel:           o.AccountChannel,
		Action:                   o.Action,
		Symbol:                   o.Symbol,
		UnderlyingSymbol:         o.UnderlyingSymbol,
		Code:                     o.Code,
		Name:                     o.Name,
		SecurityType:             o.SecurityType,
		Currency:                 o.Currency,
		TradeCurrency:            o.TradeCurrency,
		OrderType:                o.OrderType,
		Status:                   o.Status,
		Price:                    o.Price,
		Quantity:                 o.Quantity,
		ExecutedQty:              o.ExecutedQty,
		ExecutedPrice:            o.ExecutedPrice,
		ExecutedAmount:           o.ExecutedAmount,
		OperateDirection:         o.OperateDirection,
		TimeInForce:              o.TimeInForce,
		GTD:                      o.GTD,
		SubmittedAt:              o.SubmittedAt.Time(),
		UpdatedAt:                o.UpdatedAt.Time(),
		Msg:                      o.Msg,
		Report:                   o.Report,
		ContractDirection:        o.ContractDirection,
		ContractDueDate:          o.ContractDueDate,
		StrikePrice:              o.StrikePrice,
		TailingAmount:            o.TailingAmount,
		TailingPercent:           o.TailingPercent,
		TriggerPrice:             o.TriggerPrice,
		TriggerStatus:            o.TriggerStatus,
		TriggerAt:                o.TriggerAt,
		TriggerExchange:          o.TriggerExchange,
		TriggerLastDone:          o.TriggerLastDone,
		TriggerCount:             o.TriggerCount,
		LotSize:                  o.LotSize,
		LimitOffset:              o.LimitOffset,
		LimitDepthLevel:          o.LimitDepthLevel,
		MarketPrice:              o.MarketPrice,
		LastDone:                 o.LastDone,
		MonitorPrice:             o.MonitorPrice,
		SubmittedAmount:          o.SubmittedAmount,
		PlatformDeductionsStatus: o.PlatformDeductionsStatus,
		PloyID:                   o.PloyID,
		PloyType:                 o.PloyType,
		TickerSize:               o.TickerSize,
		CurrentMillisecond:       o.CurrentMillisecond,
		OrgID:                    o.OrgID,
		Tag:                      o.Tag,
		ForceOnlyRTH:             o.ForceOnlyRTH,
		DeductionsStatus:         o.DeductionsStatus,
		FreeStatus:               o.FreeStatus,
		Trend:                    o.Trend,
	}
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

	// Build symbols from Symbol
	symbols := []string{}
	if req.Symbol != "" {
		symbols = append(symbols, req.Symbol)
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
		"symbols":         symbols,
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
		StockList:         raw.StockList,
		OptionList:        raw.OptionList,
		CryptoList:        raw.CryptoList,
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

// USCashAndBuyPower returns cash, buying power and margin status for the US
// account. currency is optional (e.g. "USD"); empty uses the server default.
//
// Path: GET /v1/us/assets/cash-and-buy-power
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USCashAndBuyPower(ctx context.Context, currency string) (*USCashAndBuyPower, error) {
	const path = "/v1/us/assets/cash-and-buy-power"
	if err := c.opts.httpClient.CheckRegion(path, "US"); err != nil {
		return nil, err
	}
	q := url.Values{}
	if currency != "" {
		q.Set("currency", currency)
	}
	var resp USCashAndBuyPower
	if err := c.opts.httpClient.Get(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// USAvailableHolding returns the quantity of a US holding that is available to
// trade. EquityType defaults to "STOCK" and PositionSide to "LONG".
//
// Path: GET /v1/us/assets/available-holding
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USAvailableHolding(ctx context.Context, req *GetUSAvailableHolding) (*USAvailableHolding, error) {
	const path = "/v1/us/assets/available-holding"
	if err := c.opts.httpClient.CheckRegion(path, "US"); err != nil {
		return nil, err
	}
	if req == nil || req.Symbol == "" {
		return nil, fmt.Errorf("trade: USAvailableHolding requires a symbol")
	}
	equityType := req.EquityType
	if equityType == "" {
		equityType = "STOCK"
	}
	positionSide := req.PositionSide
	if positionSide == "" {
		positionSide = "LONG"
	}
	q := url.Values{}
	q.Set("symbol", req.Symbol)
	q.Set("equity_type", equityType)
	q.Set("position_side", positionSide)
	var resp USAvailableHolding
	if err := c.opts.httpClient.Get(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// USTransactionHistory returns one page of the US account transaction history
// (deposits, withdrawals, trades, corporate actions, transfers and so on).
// Page defaults to 1 and Size to 20.
//
// Path: POST /v1/us/assets/transaction-history
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USTransactionHistory(ctx context.Context, req *GetUSTransactionHistory) (*USTransactionHistory, error) {
	const path = "/v1/us/assets/transaction-history"
	if err := c.opts.httpClient.CheckRegion(path, "US"); err != nil {
		return nil, err
	}
	if req == nil || req.StartAt.IsZero() || req.EndAt.IsZero() {
		return nil, fmt.Errorf("trade: USTransactionHistory requires StartAt and EndAt")
	}
	body := map[string]interface{}{
		"start": req.StartAt.Unix(),
		"end":   req.EndAt.Unix(),
		"page":  defaultInt32(req.Page, 1),
		"size":  defaultInt32(req.Size, 20),
	}
	if len(req.Types) > 0 {
		body["types"] = req.Types
	}
	if req.Currency != "" {
		body["currency"] = req.Currency
	}
	var resp USTransactionHistory
	if err := c.opts.httpClient.Post(ctx, path, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// USExchangeRates returns the exchange rates used by the US account. The
// upstream two-level map is flattened into USExchangeRates.Exchanges; each entry
// keeps the two upstream map keys as UniversalKey and SpecKey.
//
// Path: GET /v1/us/assets/exchange-rates
// US token required; returns *http.RegionRestrictedError for non-US credentials.
func (c *TradeContext) USExchangeRates(ctx context.Context) (*USExchangeRates, error) {
	const path = "/v1/us/assets/exchange-rates"
	if err := c.opts.httpClient.CheckRegion(path, "US"); err != nil {
		return nil, err
	}
	var raw usRawExchangeRates
	if err := c.opts.httpClient.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	return raw.flatten(), nil
}

func defaultInt32(v, def int32) int32 {
	if v > 0 {
		return v
	}
	return def
}
