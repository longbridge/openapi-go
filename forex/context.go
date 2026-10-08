// Package forex provides a client for the Longbridge forex (currency exchange)
// OpenAPI. It lets a caller convert cash between two currencies inside their
// securities account. The flow is three steps and fully asynchronous:
//
//  1. Quote    - get an exchange rate and a quote_id valid for a short window.
//  2. SubmitOrder - place the conversion; a successful call only means accepted.
//  3. Order    - poll by client_order_id until the order reaches a terminal state.
package forex

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/forex/jsontypes"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
)

// ForexContext is a client for the Longbridge forex (currency exchange) OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	fctx, err := forex.NewFromCfg(conf)
//	quote, err := fctx.Quote(ctx, &forex.GetForexQuote{From: "USD", To: "HKD", Amount: &amount})
type ForexContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a ForexContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*ForexContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &ForexContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a ForexContext configured from environment variables.
func NewFromEnv() (*ForexContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// Quote returns a forex quote. Provide at most one of params.Amount /
// params.TargetAmount; the server validates the choice.
func (c *ForexContext) Quote(ctx context.Context, params *GetForexQuote) (quote *ForexQuote, err error) {
	body := jsontypes.GetForexQuoteBody{
		From: params.From,
		To:   params.To,
	}
	if params.Amount != nil {
		body.Amount = params.Amount.String()
	}
	if params.TargetAmount != nil {
		body.TargetAmount = params.TargetAmount.String()
	}
	resp := &jsontypes.ForexQuote{}
	if err = c.httpClient.Post(ctx, "/v1/forex/quote", body, resp); err != nil {
		return
	}
	quote = &ForexQuote{}
	err = util.Copy(quote, resp)
	return
}

// SubmitOrder submits a forex order.
//
// A nil error only means the order was accepted; conversion is asynchronous —
// poll Order with the same ClientOrderID for the final state. Reusing a
// ClientOrderID yields at most one order; on timeout / 5xx, do not change it:
// query first, and only re-quote + resubmit (same ClientOrderID) if no order is
// found. Each quote_id can be submitted successfully only once (reuse -> 603006).
func (c *ForexContext) SubmitOrder(ctx context.Context, params *SubmitForexOrder) (err error) {
	body := jsontypes.SubmitForexOrderBody{
		QuoteID:       params.QuoteID,
		ClientOrderID: params.ClientOrderID,
	}
	return c.httpClient.Post(ctx, "/v1/forex/order", body, nil)
}

// Order queries a forex order by clientOrderID.
func (c *ForexContext) Order(ctx context.Context, clientOrderID string) (detail *ForexOrderDetail, err error) {
	vals := url.Values{}
	vals.Set("client_order_id", clientOrderID)
	resp := &jsontypes.ForexOrderDetail{}
	if err = c.httpClient.Get(ctx, "/v1/forex/order", vals, resp); err != nil {
		return
	}
	detail = &ForexOrderDetail{}
	err = util.Copy(detail, resp)
	return
}
