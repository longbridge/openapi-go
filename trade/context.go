// Package trade provide TradeContext
package trade

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/longbridge"
	"github.com/longbridge/openapi-go/trade/jsontypes"
)

// TradeContext is a client for interacting with Longbridge Trade OpenAPI.
// Longbridge Trade OpenAPI document is https://open.longbridge.com/en/docs/trade/trade-overview
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	tctx.OnTrade(func(orderEvent *trade.PushEvent) {
//	  fmt.Printf("order event: %v", orderEvent)
//	})
//	_, err := tctx.Subscribe(context.Background(), []string{"private"})
//	price := decimal.NewFromString("175.62")
//	oid, err := tctx.SubmitOrder(context.Background(), &trade.SubmitOrder{
//	  Symbol: "AAPL.US",
//	  OrderType: trade.OrderTypeLO,
//	  Side: trade.OrderSideBuy,
//	  SubmittedPrice: price,
//	  SubmittedQuantity: 2,
//	  TimeInForce: trade.TimeTypeDay,
//	})
type TradeContext struct {
	opts *Options
	core *core
}

// OnQuote set callback function which will be called when server push events.
func (c *TradeContext) OnTrade(f func(*PushEvent)) {
	c.core.SetHandler(f)
}

// Subscribe topics then the handler will receive push event.
// Reference: https://open.longbridge.com/en/docs/trade/trade-push#subscribe
func (c *TradeContext) Subscribe(ctx context.Context, topics []string) (subRes *SubResponse, err error) {
	return c.core.Subscribe(ctx, topics)
}

// Unsubscribe topics then the handler will not receive the symbol's event.
// Reference: https://open.longbridge.com/en/docs/trade/trade-push#cancel-subscribe
func (c *TradeContext) Unsubscribe(ctx context.Context, topics []string) (unsubRes *UnsubResponse, err error) {
	return c.core.Unsubscribe(ctx, topics)
}

// HistoryExecutions will return history executions.
// Reference: https://open.longbridge.com/en/docs/trade/execution/history_executions
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	trades, err := tctx.HistoryExecutions(context.Background(), &trade.GetHistoryExecutions{
//	  Symbol: "AAPL.US",
//	  StartAt: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
//	  EndAt: time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC),
//	})
func (c *TradeContext) HistoryExecutions(ctx context.Context, params *GetHistoryExecutions) (trades []*Execution, err error) {
	resp := &jsontypes.Executions{}
	err = c.opts.httpClient.Get(ctx, "/v1/trade/execution/history", params.Values(), &resp)
	if err != nil {
		return
	}
	err = util.Copy(&trades, resp.Trades)
	return
}

// TodayExecutions will return today's executions
// Reference: https://open.longbridge.com/en/docs/trade/execution/today_executions
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	trades, err := tctx.TodayExecutions(context.Background(), &trade.GetTodayExecutions{Symbol: "AAPL.US"})
func (c *TradeContext) TodayExecutions(ctx context.Context, params *GetTodayExecutions) (trades []*Execution, err error) {
	resp := &jsontypes.Executions{}
	err = c.opts.httpClient.Get(ctx, "/v1/trade/execution/today", params.Values(), resp)
	if err != nil {
		return
	}
	err = util.Copy(&trades, resp.Trades)
	return
}

// HistoryOrders will return history orders
// Reference: https://open.longbridge.com/en/docs/trade/order/history_orders
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	orders, hasMore, err := tctx.HistoryOrders(context.Background(), &trade.GetHistoryOrders{
//	  Symbol: "AAPL.US",
//	  Status: []trade.OrderStatus{trade.OrderFilledStatus},
//	})
func (c *TradeContext) HistoryOrders(ctx context.Context, params *GetHistoryOrders) (orders []*Order, hasMore bool, err error) {
	resp := &jsontypes.Orders{}
	err = c.opts.httpClient.Get(ctx, "/v1/trade/order/history", params.Values(), resp)
	if err != nil {
		return
	}
	hasMore = resp.HasMore
	err = util.Copy(&orders, resp.Orders)
	return
}

// TodayOrders will return today orders
// Reference: https://open.longbridge.com/en/docs/trade/order/today_orders
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	orders, err := tctx.TodayOrders(context.Background(), &trade.GetTodayOrders{Symbol: "AAPL.US"})
func (c *TradeContext) TodayOrders(ctx context.Context, params *GetTodayOrders) (orders []*Order, err error) {
	resp := &jsontypes.Orders{}
	err = c.opts.httpClient.Get(ctx, "/v1/trade/order/today", params.Values(), resp)
	if err != nil {
		return
	}
	err = util.Copy(&orders, resp.Orders)
	return
}

// ReplaceOrder modify quantity or price
// Reference: https://open.longbridge.com/en/docs/trade/order/replace
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	err := tctx.ReplaceOrder(context.Background(), &trade.ReplaceOrder{OrderId: "123123", Quantity: 2, Remark: "just replace the order"})
func (c *TradeContext) ReplaceOrder(ctx context.Context, params *ReplaceOrder) (err error) {
	var jsonbody jsontypes.ReplaceOrder
	err = util.Copy(&jsonbody, params)
	if err != nil {
		return
	}
	err = c.opts.httpClient.Put(ctx, "/v1/trade/order", jsonbody, nil)
	return
}

// SubmitOrder HK and US stocks, warrant and option
// Reference: https://open.longbridge.com/en/docs/trade/order/submit
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	price := decimal.NewFromString("175.62")
//	oid, err := tctx.SubmitOrder(context.Background(), &trade.SubmitOrder{
//	  Symbol: "AAPL.US",
//	  OrderType: trade.OrderTypeLO,
//	  Side: trade.OrderSideBuy,
//	  SubmittedPrice: price,
//	  SubmittedQuantity: 2,
//	  TimeInForce: trade.TimeTypeDay,
//	})

func (c *TradeContext) SubmitOrder(ctx context.Context, params *SubmitOrder) (orderId string, err error) {
	var jsonbody jsontypes.SubmitOrder
	err = util.Copy(&jsonbody, params)
	if err != nil {
		return
	}
	resp := &jsontypes.SubmitOrderResponse{}
	err = c.opts.httpClient.Post(ctx, "/v1/trade/order", jsonbody, resp)
	if err != nil {
		return
	}
	return resp.OrderId, nil
}

// WithdrawOrder to close an open order. It is same as CancelOrder function.
// Reference: https://open.longbridge.com/en/docs/trade/order/withdraw
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	err = tctx.WithdrawOrder(context.Background(), "12123123")
func (c *TradeContext) WithdrawOrder(ctx context.Context, orderId string) (err error) {
	return c.CancelOrder(ctx, orderId)
}

// CancelOrder to close an open order
// Reference: https://open.longbridge.com/en/docs/trade/order/withdraw
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	err = tctx.CancelOrder(context.Background(), "12123123")
func (c *TradeContext) CancelOrder(ctx context.Context, orderId string) (err error) {
	values := url.Values{}
	values.Add("order_id", orderId)
	err = c.opts.httpClient.Delete(ctx, "/v1/trade/order", values, nil)
	return
}

// AccountBalance to obtain the available, desirable, frozen, to-be-settled, and in-transit funds (fund purchase and redemption) information for each currency of the user.
// Reference: https://open.longbridge.com/en/docs/trade/asset/account
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	ab, err := trade.AccountBalance(context.Background(), &trade.GetAccountBalance{Currency: trade.CurrencyHKD})
func (c *TradeContext) AccountBalance(ctx context.Context, params *GetAccountBalance) (accounts []*AccountBalance, err error) {
	var resp jsontypes.AccountBalances
	err = c.opts.httpClient.Get(ctx, "/v1/asset/account", params.Values(), &resp)
	if err != nil {
		return
	}
	err = util.Copy(&accounts, resp.List)
	return
}

// CashFlow to obtain capital inflow/outflow direction, capital type, capital amount, occurrence time, associated stock code and capital flow description information.
// Reference: https://open.longbridge.com/en/docs/trade/asset/cashflow
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	start := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC).Unix()
//	end := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC).Unix()
//	flows, err := tctx.CashFlow(context.Background(), trade.GetCashFlow{
//	  StartAt: start,
//	  EndAt: end,
//	  BussinessType: trade.BalanceTypeCash,
//	})
func (c *TradeContext) CashFlow(ctx context.Context, params *GetCashFlow) (cashflows []*CashFlow, err error) {
	var resp jsontypes.CashFlows
	err = c.opts.httpClient.Get(ctx, "/v1/asset/cashflow", params.Values(), &resp)
	if err != nil {
		return
	}
	err = util.Copy(&cashflows, resp.List)
	return
}

// FundPositions to obtain fund position information including account, fund code, holding share, cost net worth, current net worth, and currency.
// Reference: https://open.longbridge.com/en/docs/trade/asset/fund
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	channels, err := tctx.FundPositions(context.Background, []string{"AAPL.US", "700.HK"})
func (c *TradeContext) FundPositions(ctx context.Context, symbols []string) (fundPositionChannels []*FundPositionChannel, err error) {
	params := &GetFundPositions{
		Symbols: symbols,
	}
	var resp jsontypes.FundPositions
	err = c.opts.httpClient.Get(ctx, "/v1/asset/fund", params.Values(), &resp)
	if err != nil {
		return
	}
	err = util.Copy(&fundPositionChannels, resp.List)
	return
}

// StockPositions to obtain stock position information including account, stock code, number of shares held, number of available shares, average position price (calculated according to account settings), and currency.
// Reference: https://open.longbridge.com/en/docs/trade/asset/stock
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	channels, err := tctx.StockPositions(context.Background(), []string{"AAPL.US"})
func (c *TradeContext) StockPositions(ctx context.Context, symbols []string) (stockPositionChannels []*StockPositionChannel, err error) {
	params := &GetStockPositions{
		Symbols: symbols,
	}
	var resp jsontypes.StockPositions
	err = c.opts.httpClient.Get(ctx, "/v1/asset/stock", params.Values(), &resp)
	if err != nil {
		return
	}
	err = util.Copy(&stockPositionChannels, resp.List)
	return
}

// MarginRatio is used to obtain the initial margin ratio, maintain the margin ratio and strengthen the margin ratio of stocks.
// Reference: https://open.longbridge.com/en/docs/trade/asset/margin_ratio
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	mr, err := tctx.MarginRatio(context.Background(), "AAPL.US")
func (c *TradeContext) MarginRatio(ctx context.Context, symbol string) (marginRatio MarginRatio, err error) {
	values := url.Values{}
	values.Add("symbol", symbol)
	var resp jsontypes.MarginRatio
	err = c.opts.httpClient.Get(ctx, "/v1/risk/margin-ratio", values, &resp)
	if err != nil {
		return
	}
	err = util.Copy(&marginRatio, resp)
	return
}

// OrderDetail is used for order detail query
// Reference: https://open.longbridge.com/en/docs/trade/order/order_detail
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	od, err := tctx.OrderDetail(context.Background(), "1123123123")
func (c *TradeContext) OrderDetail(ctx context.Context, orderId string) (orderDetail OrderDetail, err error) {
	values := url.Values{}
	values.Add("order_id", orderId)
	var resp jsontypes.OrderDetail
	err = c.opts.httpClient.Get(ctx, "/v1/trade/order", values, &resp)
	if err != nil {
		return
	}
	err = util.Copy(&orderDetail, resp)
	return
}

// EstimateMaxPurchaseQuantity is used for estimating the maximum purchase quantity for Hong Kong and US stocks, warrants, and options.
// Reference: https://open.longbridge.com/en/docs/trade/order/estimate_available_buy_limit
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	price, _ := decimal.NewFromString("175.62")
//	empqr, err := trade.EstimateMaxPurchaseQuantity(context.Background(), &trade.GetEstimateMaxPurchaseQuantity{
//	  Symbol: "AAPL.US",
//	  Price: price,
//	  OrderType: trade.OrderTypeLO,
//	  Currency: "USD",
//	  Side: trade.OrderSideBuy,
//	})
func (c *TradeContext) EstimateMaxPurchaseQuantity(ctx context.Context, params *GetEstimateMaxPurchaseQuantity) (empqr EstimateMaxPurchaseQuantityResponse, err error) {
	values := params.Values()
	var resp jsontypes.EstimateMaxPurchaseQuantityResponse
	err = c.opts.httpClient.Get(ctx, "/v1/trade/estimate/buy_limit", values, &resp)
	if err != nil {
		return
	}
	err = util.Copy(&empqr, resp)
	return
}

// Close
func (c *TradeContext) Close() error {
	return c.core.Close()
}

// Deprecated: NewFormEnv return TradeContext with environment variables, use NewFromCfg plz
func NewFormEnv() (*TradeContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, err
	}
	return NewFromCfg(cfg)
}

// NewFromCfg return TradeContext with config.Config.
func NewFromCfg(cfg *config.Config) (*TradeContext, error) {
	httpClient, err := http.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	lbOpts := longbridge.NewOptions(
		longbridge.WithAuthTimeout(cfg.AuthTimeout),
		longbridge.WithTimeout(cfg.Timeout),
		longbridge.WithReadBufferSize(cfg.ReadBufferSize),
		longbridge.WithReadQueueSize(cfg.ReadQueueSize),
		longbridge.WithWriteQueueSize(cfg.WriteQueueSize),
		longbridge.WithMinGzipSize(cfg.MinGzipSize),
	)
	return New(
		WithTradeURL(cfg.TradeUrl),
		WithHttpClient(httpClient),
		WithLbOptions(lbOpts),
		WithLogger(cfg.Logger()),
		WithLogLevel(cfg.LogLevel),
	)
}

// New return TradeContext with option.
// A connection will be created with Trade server.
func New(opt ...Option) (*TradeContext, error) {
	opts := newOptions(opt...)
	core, err := newCore(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create core")
	}
	tc := &TradeContext{
		opts: opts,
		core: core,
	}
	return tc, nil
}
