// Package grid provides a client for the Longbridge grid trading OpenAPI.
// It supports submitting, replacing, querying, and controlling
// (cancel / suspend / restart) grid trading orders, as well as recording the
// strategy risk-disclosure questionnaire and fetching order-window info.
package grid

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/grid/jsontypes"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
)

// GridContext is a client for the Longbridge grid trading OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	gctx, err := grid.NewFromCfg(conf)
//	orders, err := gctx.GridOrders(ctx, &grid.GetGridOrders{})
type GridContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a GridContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*GridContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &GridContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a GridContext configured from environment variables.
func NewFromEnv() (*GridContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// Submit submits a grid trading order.
// Reference: https://open.longbridge.com/en/docs/trade/gridtrading/submit
func (c *GridContext) Submit(ctx context.Context, params *SubmitGridOrder) (orderId string, err error) {
	var rule jsontypes.GridTradeRule
	if err = util.Copy(&rule, params.GridTradingRule); err != nil {
		return
	}
	body := jsontypes.SubmitGridOrder{
		Symbol:             params.Symbol,
		SettlementCurrency: params.SettlementCurrency,
		GridTradingRule:    rule,
	}
	resp := &jsontypes.SubmitGridOrderResponse{}
	err = c.httpClient.Post(ctx, "/v1/gridtrading/submit", body, resp)
	if err != nil {
		return
	}
	return resp.OrderId, nil
}

// Replace replaces (modifies) a grid trading order.
// Reference: https://open.longbridge.com/en/docs/trade/gridtrading/replace
func (c *GridContext) Replace(ctx context.Context, params *ReplaceGridOrder) (err error) {
	var rule jsontypes.GridTradeRule
	if err = util.Copy(&rule, params.GridTradingRule); err != nil {
		return
	}
	body := jsontypes.ReplaceGridOrder{
		OrderId:         params.OrderId,
		GridTradingRule: rule,
	}
	return c.httpClient.Post(ctx, "/v1/gridtrading/replace", body, nil)
}

// List returns a paged list of grid trading orders.
// Reference: https://open.longbridge.com/en/docs/trade/gridtrading/list
func (c *GridContext) List(ctx context.Context, params *GetGridOrders) (orders *GridOrders, err error) {
	resp := &jsontypes.GridOrders{}
	err = c.httpClient.Get(ctx, "/v1/gridtrading/list", params.Values(), resp)
	if err != nil {
		return
	}
	orders = &GridOrders{HasMore: resp.HasMore}
	err = util.Copy(&orders.GridOrder, resp.GridOrder)
	return
}

// ListByIds queries grid trading orders by their IDs.
// Reference: https://open.longbridge.com/en/docs/trade/gridtrading/list
func (c *GridContext) ListByIds(ctx context.Context, orderIds []string) (orders []*GridOrder, err error) {
	body := jsontypes.GridOrderIdsBody{OrderIds: orderIds}
	resp := &jsontypes.GridOrders{}
	err = c.httpClient.Post(ctx, "/v1/gridtrading/list", body, resp)
	if err != nil {
		return
	}
	err = util.Copy(&orders, resp.GridOrder)
	return
}

// Detail returns the detail (and paged history) of a grid trading order.
// Reference: https://open.longbridge.com/en/docs/trade/gridtrading/detail
func (c *GridContext) Detail(ctx context.Context, params *GetGridOrderDetail) (detail *GridOrderDetail, err error) {
	resp := &jsontypes.GridOrderDetail{}
	err = c.httpClient.Get(ctx, "/v1/gridtrading/detail", params.Values(), resp)
	if err != nil {
		return
	}
	detail = &GridOrderDetail{}
	err = util.Copy(detail, resp)
	return
}

// TriggerHistory returns the trigger history of a grid trading order.
// Reference: https://open.longbridge.com/en/docs/trade/gridtrading/trigger_history
func (c *GridContext) TriggerHistory(ctx context.Context, params *GetGridTriggerHistory) (history *GridTriggerHistory, err error) {
	resp := &jsontypes.GridTriggerHistory{}
	err = c.httpClient.Get(ctx, "/v1/gridtrading/trigger_history_list", params.Values(), resp)
	if err != nil {
		return
	}
	history = &GridTriggerHistory{HasMore: resp.HasMore}
	err = util.Copy(&history.TriggerOrders, resp.TriggerOrders)
	return
}

// Cancel cancels a grid trading order.
// Reference: https://open.longbridge.com/en/docs/trade/gridtrading/cancel
func (c *GridContext) Cancel(ctx context.Context, orderId string) (err error) {
	return c.gridAction(ctx, "/v1/gridtrading/cancel", orderId)
}

// Suspend suspends a grid trading order.
// Reference: https://open.longbridge.com/en/docs/trade/gridtrading/suspend
func (c *GridContext) Suspend(ctx context.Context, orderId string) (err error) {
	return c.gridAction(ctx, "/v1/gridtrading/suspend", orderId)
}

// Restart restarts a grid trading order.
// Reference: https://open.longbridge.com/en/docs/trade/gridtrading/restart
func (c *GridContext) Restart(ctx context.Context, orderId string) (err error) {
	return c.gridAction(ctx, "/v1/gridtrading/restart", orderId)
}

// gridAction is the shared body for the cancel / suspend / restart grid actions.
func (c *GridContext) gridAction(ctx context.Context, path string, orderId string) (err error) {
	body := jsontypes.GridOrderIdBody{OrderId: orderId}
	return c.httpClient.Post(ctx, path, body, nil)
}

// SubmitStrategyQuestionnaire records the user's consent to the strategy
// risk disclosure required before using grid trading. The body sent is
// {"type": "strategy", "items": {"agree": "true"}}.
func (c *GridContext) SubmitStrategyQuestionnaire(ctx context.Context) (err error) {
	body := jsontypes.SubmitStrategyQuestionnaire{
		Type:  "strategy",
		Items: map[string]string{"agree": "true"},
	}
	return c.httpClient.Post(ctx, "/v1/record/questionnaire", body, nil)
}

// OrderInfo returns order info used by the grid order window (lot size,
// authorization flag, settlement currency, etc.).
func (c *GridContext) OrderInfo(ctx context.Context, symbol string) (info *GridOrderInfo, err error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	resp := &jsontypes.GridOrderInfo{}
	err = c.httpClient.Get(ctx, "/v1/orders/info", values, resp)
	if err != nil {
		return
	}
	info = &GridOrderInfo{}
	err = util.Copy(info, resp)
	return
}
