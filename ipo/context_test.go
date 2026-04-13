package ipo_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/longbridgeapp/assert"

	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/ipo"
)

// newTestServer starts a mock HTTP server and returns an IPOContext wired to it.
// handler receives all requests and writes the JSON response body (without envelope).
func newTestServer(t *testing.T, handler http.HandlerFunc) (*ipo.IPOContext, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}))
	httpClient, err := httplib.New(
		httplib.WithURL(srv.URL),
		httplib.WithAppKey("test-key"),
		httplib.WithAppSecret("test-secret"),
		httplib.WithAccessToken("test-token"),
	)
	if err != nil {
		t.Fatal(err)
	}
	return ipo.NewWithHTTPClient(httpClient), srv
}

// envelope wraps data in the standard Longbridge API response envelope.
func envelope(t *testing.T, data interface{}) []byte {
	t.Helper()
	inner, _ := json.Marshal(data)
	out, _ := json.Marshal(map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    json.RawMessage(inner),
	})
	return out
}

func TestSubmitOrder(t *testing.T) {
	ctx, srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(envelope(t, map[string]interface{}{
			"id":          int64(100001),
			"name":        "Alibaba",
			"symbol":      "9988.HK",
			"market":      "HK",
			"sub_qty":     "1000",
			"sub_amount":  "88000.00",
			"handing_fee": "100.00",
			"status":      1,
			"currency":    "HKD",
			"method":      2,
			"batch_id":    int64(12345),
			"modify_enable": true,
			"withdraw_able": false,
		}))
	})
	defer srv.Close()

	result, err := ctx.SubmitOrder(context.Background(), &ipo.SubmitOrderRequest{
		Symbol:  "9988.HK",
		SubQty:  "1000",
		BatchID: 12345,
		Method:  2,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(100001), result.ID)
	assert.Equal(t, "9988.HK", result.Symbol)
	assert.Equal(t, "HK", result.Market)
	assert.Equal(t, "1000", result.SubQty)
	assert.Equal(t, "88000.00", result.SubAmount)
	assert.Equal(t, "100.00", result.HandingFee)
	assert.Equal(t, ipo.SubscriptionStatus(1), result.Status)
	assert.Equal(t, ipo.SubscriptionMethod(2), result.Method)
	assert.Equal(t, int64(12345), result.BatchID)
	assert.True(t, result.ModifyEnable)
	assert.False(t, result.WithdrawAble)
}

func TestAmendOrder(t *testing.T) {
	ctx, srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(envelope(t, map[string]interface{}{
			"id":      int64(100001),
			"symbol":  "9988.HK",
			"sub_qty": "2000",
			"status":  1,
			"method":  2,
		}))
	})
	defer srv.Close()

	result, err := ctx.AmendOrder(context.Background(), &ipo.AmendOrderRequest{
		Symbol:  "9988.HK",
		SubQty:  "2000",
		BatchID: 12345,
		Method:  2,
		OrderID: 100001,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(100001), result.ID)
	assert.Equal(t, "2000", result.SubQty)
	assert.Equal(t, ipo.SubscriptionStatus(1), result.Status)
	assert.Equal(t, ipo.SubscriptionMethod(2), result.Method)
}

func TestWithdrawOrder(t *testing.T) {
	ctx, srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(envelope(t, map[string]interface{}{"result": true}))
	})
	defer srv.Close()

	ok, err := ctx.WithdrawOrder(context.Background(), &ipo.WithdrawOrderRequest{
		OrderID: 100001,
	})
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestFetchOrderList(t *testing.T) {
	ctx, srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(envelope(t, map[string]interface{}{
			"list": []map[string]interface{}{
				{
					"id":               int64(100001),
					"name":             "Alibaba",
					"symbol":           "9988.HK",
					"status":           2,
					"sub_qty":          "1000",
					"method":           1,
					"created_at":       int64(1700000000),
					"currency":         "HKD",
					"lot_win_qty":      "0",
					"sub_amount":       "88000.00",
					"cash_amount":      "8800.00",
					"financing_amount": "79200.00",
					"financing_ratio":  "0.9",
					"handing_fee":      "100.00",
				},
			},
			"has_more": true,
		}))
	})
	defer srv.Close()

	orders, hasMore, err := ctx.FetchOrderList(context.Background(), &ipo.FetchOrderListRequest{
		Symbol:   "9988.HK",
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.True(t, hasMore)
	assert.Equal(t, 1, len(orders))

	o := orders[0]
	assert.Equal(t, int64(100001), o.ID)
	assert.Equal(t, "9988.HK", o.Symbol)
	assert.Equal(t, ipo.SubscriptionStatus(2), o.Status)
	assert.Equal(t, ipo.SubscriptionMethod(1), o.Method)
	assert.Equal(t, "1000", o.SubQty)
	assert.Equal(t, int64(1700000000), o.CreatedAt)
	assert.Equal(t, "HKD", o.Currency)
	assert.Equal(t, "88000.00", o.SubAmount)
	assert.Equal(t, "8800.00", o.CashAmount)
	assert.Equal(t, "79200.00", o.FinancingAmount)
	assert.Equal(t, "100.00", o.HandingFee)
}

func TestFetchOrderDetail(t *testing.T) {
	ctx, srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(envelope(t, map[string]interface{}{
			"id":                int64(100001),
			"name":              "Alibaba",
			"code":              "9988",
			"symbol":            "9988.HK",
			"market":            "HK",
			"sub_qty":           "1000",
			"sub_amount":        "88000.00",
			"handing_fee":       "100.00",
			"status":            2,
			"currency":          "HKD",
			"method":            1,
			"total_amount":      "88100.00",
			"need_to_pay":       "0",
			"financing_amount":  "0",
			"ipo_price":         "88.00",
			"modify_enable":     true,
			"withdraw_able":     false,
			"luck_at":           int64(0),
			"batch_id":          int64(12345),
			"channel":           "PUBLIC_OFFER",
		}))
	})
	defer srv.Close()

	detail, err := ctx.FetchOrderDetail(context.Background(), &ipo.FetchOrderDetailRequest{
		OrderID: 100001,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(100001), detail.ID)
	assert.Equal(t, "9988", detail.Code)
	assert.Equal(t, "9988.HK", detail.Symbol)
	assert.Equal(t, "HK", detail.Market)
	assert.Equal(t, "1000", detail.SubQty)
	assert.Equal(t, "88000.00", detail.SubAmount)
	assert.Equal(t, "88100.00", detail.TotalAmount)
	assert.Equal(t, "88.00", detail.IpoPrice)
	assert.Equal(t, ipo.SubscriptionStatus(2), detail.Status)
	assert.Equal(t, ipo.SubscriptionMethod(1), detail.Method)
	assert.Equal(t, int64(12345), detail.BatchID)
	assert.Equal(t, "PUBLIC_OFFER", detail.Channel)
	assert.True(t, detail.ModifyEnable)
	assert.False(t, detail.WithdrawAble)
}

func TestFetchMarginList(t *testing.T) {
	ctx, srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(envelope(t, map[string]interface{}{
			"list": []map[string]interface{}{
				{
					"method":             1,
					"handing_fee":        "100.00",
					"currency":           "HKD",
					"financing_fee_rate": "0",
					"remain_amount":      "0",
					"batch_id":           int64(12345),
					"is_ended":           false,
					"display_name":       "Cash",
					"min_cash":           "10000.00",
					"start_at":           int64(1699900000),
					"deadline":           int64(1700000000),
					"multiple":           "1",
					"financing_ratio":    "0",
				},
				{
					"method":             2,
					"handing_fee":        "150.00",
					"currency":           "HKD",
					"financing_fee_rate": "0.06",
					"remain_amount":      "200000.00",
					"batch_id":           int64(12345),
					"is_ended":           false,
					"display_name":       "10x Margin",
					"min_cash":           "5000.00",
					"start_at":           int64(1699900000),
					"deadline":           int64(1700000000),
					"multiple":           "10",
					"financing_ratio":    "0.9",
				},
			},
		}))
	})
	defer srv.Close()

	margins, err := ctx.FetchMarginList(context.Background(), &ipo.FetchMarginListRequest{
		Symbol: "9988.HK",
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(margins))

	cash := margins[0]
	assert.Equal(t, ipo.SubscriptionMethod(1), cash.Method)
	assert.Equal(t, "Cash", cash.DisplayName)
	assert.Equal(t, "HKD", cash.Currency)
	assert.Equal(t, "100.00", cash.HandingFee)
	assert.Equal(t, "10000.00", cash.MinCash)
	assert.Equal(t, int64(12345), cash.BatchID)
	assert.Equal(t, int64(1700000000), cash.Deadline)
	assert.False(t, cash.IsEnded)

	margin := margins[1]
	assert.Equal(t, ipo.SubscriptionMethod(2), margin.Method)
	assert.Equal(t, "0.06", margin.FinancingFeeRate)
	assert.Equal(t, "200000.00", margin.RemainAmount)
	assert.Equal(t, "10", margin.Multiple)
}

func TestFetchIpoPaymentList(t *testing.T) {
	ctx, srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(envelope(t, map[string]interface{}{
			"list": []map[string]interface{}{
				{"number": 1, "sub_qty": "500", "sub_amount": "44000.00", "currency": "HKD"},
				{"number": 2, "sub_qty": "1000", "sub_amount": "88000.00", "currency": "HKD"},
				{"number": 3, "sub_qty": "2000", "sub_amount": "176000.00", "currency": "HKD"},
			},
		}))
	})
	defer srv.Close()

	list, err := ctx.FetchIpoPaymentList(context.Background(), &ipo.FetchIpoPaymentListRequest{
		Symbol: "9988.HK",
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(list))

	assert.Equal(t, int32(1), list[0].Number)
	assert.Equal(t, "500", list[0].SubQty)
	assert.Equal(t, "44000.00", list[0].SubAmount)
	assert.Equal(t, "HKD", list[0].Currency)

	assert.Equal(t, int32(3), list[2].Number)
	assert.Equal(t, "2000", list[2].SubQty)
	assert.Equal(t, "176000.00", list[2].SubAmount)
}

func TestFetchBuyLimit(t *testing.T) {
	ctx, srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(envelope(t, map[string]interface{}{
			"total_amount":   "500000.00",
			"current_amount": "300000.00",
			"mmf_amount":     "100000.00",
			"cash_amount":    "200000.00",
		}))
	})
	defer srv.Close()

	bp, err := ctx.FetchBuyLimit(context.Background(), &ipo.FetchBuyLimitRequest{
		Symbol:   "9988.HK",
		Currency: "HKD",
	})
	assert.NoError(t, err)
	assert.Equal(t, "500000.00", bp.TotalAmount)
	assert.Equal(t, "300000.00", bp.CurrentAmount)
	assert.Equal(t, "100000.00", bp.MmfAmount)
	assert.Equal(t, "200000.00", bp.CashAmount)
}
