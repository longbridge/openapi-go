// Package ipo provides a client for the Longbridge IPO OpenAPI.
// It covers IPO subscription submission, modification, cancellation, and queries.
package ipo

import (
	"context"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/ipo/jsontypes"
)

// IPOContext is a client for the Longbridge IPO API.
//
// Example:
//
//	conf, err := config.New()
//	ictx, err := ipo.NewFromCfg(conf)
//	result, err := ictx.SubmitOrder(context.Background(), &ipo.SubmitOrderRequest{
//	    Symbol:  "9988.HK",
//	    SubQty:  "1000",
//	    BatchID: 12345,
//	    Method:  1,
//	})
type IPOContext struct {
	httpClient *httplib.Client
}

// NewWithHTTPClient creates an IPOContext with a pre-built http.Client (useful for testing).
func NewWithHTTPClient(httpClient *httplib.Client) *IPOContext {
	return &IPOContext{httpClient: httpClient}
}

// NewFromCfg creates an IPOContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*IPOContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &IPOContext{httpClient: httpClient}, nil
}

// NewFromEnv returns an IPOContext configured from environment variables.
func NewFromEnv() (*IPOContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// SubmitOrder submits a new IPO subscription order.
//
// Reference: POST /v1/openapi/ipo/submit
func (c *IPOContext) SubmitOrder(ctx context.Context, req *SubmitOrderRequest) (*SubmitOrderResult, error) {
	body := map[string]interface{}{
		"symbol":   req.Symbol,
		"sub_qty":  req.SubQty,
		"batch_id": req.BatchID,
		"method":   int32(req.Method),
	}
	var resp jsontypes.SubmitOrderResponse
	if err := c.httpClient.Post(ctx, "/v1/openapi/ipo/submit", body, &resp); err != nil {
		return nil, err
	}
	return convertSubmitOrderResult(&resp), nil
}

// AmendOrder modifies an existing IPO subscription order.
//
// Reference: POST /v1/openapi/ipo/amend
func (c *IPOContext) AmendOrder(ctx context.Context, req *AmendOrderRequest) (*SubmitOrderResult, error) {
	body := map[string]interface{}{
		"symbol":   req.Symbol,
		"sub_qty":  req.SubQty,
		"batch_id": req.BatchID,
		"method":   int32(req.Method),
		"order_id": req.OrderID,
	}
	var resp jsontypes.SubmitOrderResponse
	if err := c.httpClient.Post(ctx, "/v1/openapi/ipo/amend", body, &resp); err != nil {
		return nil, err
	}
	return convertSubmitOrderResult(&resp), nil
}

// WithdrawOrder cancels an IPO subscription order.
//
// Reference: POST /v1/openapi/ipo/withdraw
func (c *IPOContext) WithdrawOrder(ctx context.Context, req *WithdrawOrderRequest) (bool, error) {
	body := map[string]interface{}{
		"order_id": req.OrderID,
	}
	var resp jsontypes.WithdrawOrderResponse
	if err := c.httpClient.Post(ctx, "/v1/openapi/ipo/withdraw", body, &resp); err != nil {
		return false, err
	}
	return resp.Result, nil
}

// FetchOrderList returns the IPO subscription order list.
//
// Reference: GET /v1/openapi/ipo/order/list
func (c *IPOContext) FetchOrderList(ctx context.Context, req *FetchOrderListRequest) (orders []*OrderListItem, hasMore bool, err error) {
	var resp jsontypes.FetchOrderListResponse
	if err = c.httpClient.Get(ctx, "/v1/openapi/ipo/order/list", req.Values(), &resp); err != nil {
		return
	}
	hasMore = resp.HasMore
	orders = make([]*OrderListItem, 0, len(resp.List))
	for _, item := range resp.List {
		orders = append(orders, convertOrderListItem(item))
	}
	return
}

// FetchOrderDetail returns the full detail of an IPO subscription order.
//
// Reference: GET /v1/openapi/ipo/order
func (c *IPOContext) FetchOrderDetail(ctx context.Context, req *FetchOrderDetailRequest) (*OrderDetail, error) {
	var resp jsontypes.FetchOrderDetailResponse
	if err := c.httpClient.Get(ctx, "/v1/openapi/ipo/order", req.Values(), &resp); err != nil {
		return nil, err
	}
	return convertOrderDetail(&resp), nil
}

// FetchMarginList returns the financing scheme list for an IPO (融资方案列表).
//
// Reference: GET /v1/openapi/ipo/margin/list
func (c *IPOContext) FetchMarginList(ctx context.Context, req *FetchMarginListRequest) ([]*MarginItem, error) {
	var resp jsontypes.FetchMarginListResponse
	if err := c.httpClient.Get(ctx, "/v1/openapi/ipo/margin/list", req.Values(), &resp); err != nil {
		return nil, err
	}
	items := make([]*MarginItem, 0, len(resp.List))
	for _, item := range resp.List {
		items = append(items, convertMarginItem(item))
	}
	return items, nil
}

// FetchIpoPaymentList returns available lot quantity options and their subscription amounts.
//
// Reference: GET /v1/openapi/ipo/payable/list
func (c *IPOContext) FetchIpoPaymentList(ctx context.Context, req *FetchIpoPaymentListRequest) ([]*PaymentListItem, error) {
	var resp jsontypes.FetchIpoPaymentListResponse
	if err := c.httpClient.Get(ctx, "/v1/openapi/ipo/payable/list", req.Values(), &resp); err != nil {
		return nil, err
	}
	items := make([]*PaymentListItem, 0, len(resp.List))
	for _, item := range resp.List {
		items = append(items, convertPaymentListItem(item))
	}
	return items, nil
}

// FetchBuyLimit returns the user's buying power for IPO subscription (打新购买力).
//
// Reference: GET /v1/openapi/ipo/buylimit
func (c *IPOContext) FetchBuyLimit(ctx context.Context, req *FetchBuyLimitRequest) (*BuyingPower, error) {
	var resp jsontypes.FetchBuyLimitResponse
	if err := c.httpClient.Get(ctx, "/v1/openapi/ipo/buylimit", req.Values(), &resp); err != nil {
		return nil, err
	}
	return &BuyingPower{
		TotalAmount:   resp.TotalAmount,
		CurrentAmount: resp.CurrentAmount,
		MmfAmount:     resp.MmfAmount,
		CashAmount:    resp.CashAmount,
	}, nil
}

// --- internal converters ---

func convertSubmitOrderResult(j *jsontypes.SubmitOrderResponse) *SubmitOrderResult {
	return &SubmitOrderResult{
		ID:                j.ID,
		Name:              j.Name,
		Symbol:            j.Symbol,
		Market:            j.Market,
		SubQty:            j.SubQty,
		SubAmount:         j.SubAmount,
		HandingFee:        j.HandingFee,
		WithdrawAble:      j.WithdrawAble,
		LotWinQty:         j.LotWinQty,
		Status:            SubscriptionStatus(j.Status),
		Currency:          j.Currency,
		FinancingAmount:   j.FinancingAmount,
		NeedToPay:         j.NeedToPay,
		IpoStatus:         j.IpoStatus,
		Method:            SubscriptionMethod(j.Method),
		FinancingInterest: j.FinancingInterest,
		FinanceFeeRate:    j.FinanceFeeRate,
		IpoPrice:          j.IpoPrice,
		BorrowAmount:      j.BorrowAmount,
		BorrowInterest:    j.BorrowInterest,
		TotalAmount:       j.TotalAmount,
		CurrentAmount:     j.CurrentAmount,
		RefundAmount:      j.RefundAmount,
		LuckAt:            j.LuckAt,
		IpoDate:           j.IpoDate,
		MartBegin:         j.MartBegin,
		MartEnd:           j.MartEnd,
		RejectText:        j.RejectText,
		ModifyEnable:      j.ModifyEnable,
		BatchID:           j.BatchID,
		BorrowEnable:      j.BorrowEnable,
		LuckyFee:          j.LuckyFee,
	}
}

func convertOrderListItem(j *jsontypes.OrderListItem) *OrderListItem {
	return &OrderListItem{
		ID:              j.ID,
		Name:            j.Name,
		Symbol:          j.Symbol,
		Status:          SubscriptionStatus(j.Status),
		SubQty:          j.SubQty,
		Method:          SubscriptionMethod(j.Method),
		CreatedAt:       j.CreatedAt,
		Currency:        j.Currency,
		LotWinQty:       j.LotWinQty,
		SubAmount:       j.SubAmount,
		CashAmount:      j.CashAmount,
		FinancingAmount: j.FinancingAmount,
		FinancingRatio:  j.FinancingRatio,
		HandingFee:      j.HandingFee,
	}
}

func convertOrderDetail(j *jsontypes.FetchOrderDetailResponse) *OrderDetail {
	return &OrderDetail{
		ID:                j.ID,
		Name:              j.Name,
		Code:              j.Code,
		Market:            j.Market,
		SubQty:            j.SubQty,
		SubAmount:         j.SubAmount,
		HandingFee:        j.HandingFee,
		WithdrawAble:      j.WithdrawAble,
		LotWinQty:         j.LotWinQty,
		Status:            SubscriptionStatus(j.Status),
		Currency:          j.Currency,
		FinancingAmount:   j.FinancingAmount,
		NeedToPay:         j.NeedToPay,
		IpoStatus:         j.IpoStatus,
		Method:            SubscriptionMethod(j.Method),
		FinancingInterest: j.FinancingInterest,
		FinanceFeeRate:    j.FinanceFeeRate,
		Explanation:       j.Explanation,
		IpoPrice:          j.IpoPrice,
		BorrowAmount:      j.BorrowAmount,
		BorrowInterest:    j.BorrowInterest,
		Symbol:            j.Symbol,
		TotalAmount:       j.TotalAmount,
		CurrentAmount:     j.CurrentAmount,
		RefundAmount:      j.RefundAmount,
		LuckAt:            j.LuckAt,
		RejectText:        j.RejectText,
		Multiple:          j.Multiple,
		ModifyEnable:      j.ModifyEnable,
		LuckyFee:          j.LuckyFee,
		Channel:           j.Channel,
		BatchID:           j.BatchID,
	}
}

func convertMarginItem(j *jsontypes.MarginItem) *MarginItem {
	return &MarginItem{
		Method:           SubscriptionMethod(j.Method),
		HandingFee:       j.HandingFee,
		Currency:         j.Currency,
		Deadline:         j.Deadline,
		Multiple:         j.Multiple,
		FinancingFeeRate: j.FinancingFeeRate,
		RemainAmount:     j.RemainAmount,
		BatchID:          j.BatchID,
		IsEnded:          j.IsEnded,
		FinancingRatio:   j.FinancingRatio,
		StartAt:          j.StartAt,
		DisplayName:      j.DisplayName,
		MinCash:          j.MinCash,
	}
}

func convertPaymentListItem(j *jsontypes.PaymentListItem) *PaymentListItem {
	return &PaymentListItem{
		Number:    j.Number,
		SubQty:    j.SubQty,
		SubAmount: j.SubAmount,
		Currency:  j.Currency,
	}
}
