package alert

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/alert/jsontypes"
	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// AlertContext is a client for price alert management.
type AlertContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates an AlertContext from a Config.
func NewFromCfg(cfg *config.Config) (*AlertContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &AlertContext{httpClient: httpClient}, nil
}

// NewFromEnv creates an AlertContext from environment variables.
func NewFromEnv() (*AlertContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

func parseDecimalOpt(s string) *decimal.Decimal {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return &d
}

// List returns all price alerts for the current user.
func (c *AlertContext) List(ctx context.Context) (*AlertList, error) {
	var resp jsontypes.AlertList
	if err := c.httpClient.Get(ctx, "/v1/notify/reminders", nil, &resp); err != nil {
		return nil, err
	}
	result := &AlertList{}
	for _, g := range resp.Lists {
		group := &AlertSymbolGroup{
			Symbol:  jsontypes.CounterIDToSymbol(g.CounterId),
			Code:    g.Code,
			Market:  g.Market,
			Name:    g.Name,
			Price:   parseDecimalOpt(g.Price),
			Chg:     parseDecimalOpt(g.Chg),
			PChg:    parseDecimalOpt(g.PChg),
			Product: g.Product,
		}
		for _, item := range g.Items {
			group.Items = append(group.Items, &AlertItem{
				Id:          item.Id,
				IndicatorId: item.IndicatorId,
				Enabled:     item.Enabled,
				Frequency:   item.Frequency,
				Scope:       item.Scope,
				Text:        item.Text,
				State:       item.State,
				ValueMap:    item.ValueMap,
			})
		}
		result.Lists = append(result.Lists, group)
	}
	return result, nil
}

// Add creates a new price alert.
func (c *AlertContext) Add(ctx context.Context, symbol string, condition AlertCondition, triggerValue string, frequency AlertFrequency) error {
	counterID := jsontypes.SymbolToCounterID(symbol)
	var key string
	switch condition {
	case AlertConditionPriceRise, AlertConditionPriceFall:
		key = "price"
	default:
		key = "chg"
	}
	body := map[string]interface{}{
		"counter_id":   counterID,
		"indicator_id": strconv.Itoa(int(condition)),
		"value_map":    map[string]string{key: triggerValue},
		"frequency":    int(frequency),
		"enabled":      true,
		"scope":        0,
		"state":        []int{1},
	}
	return c.httpClient.Post(ctx, "/v1/notify/reminders", body, nil)
}

// Update modifies an existing alert. Pass the AlertItem obtained from List.
func (c *AlertContext) Update(ctx context.Context, item *AlertItem) error {
	body := map[string]interface{}{
		"id":           item.Id,
		"indicator_id": item.IndicatorId,
		"frequency":    item.Frequency,
		"scope":        item.Scope,
		"state":        item.State,
		"value_map":    item.ValueMap,
		"enabled":      item.Enabled,
	}
	return c.httpClient.Post(ctx, "/v1/notify/reminders", body, nil)
}

// Delete removes price alerts by their IDs.
func (c *AlertContext) Delete(ctx context.Context, alertIds []string) error {
	return c.httpClient.Delete(ctx, "/v1/notify/reminders", nil, nil, httplib.WithBody(map[string]interface{}{
		"ids": alertIds,
	}))
}
