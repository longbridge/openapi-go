// Package alert provides a client for the Longbridge Price Alert OpenAPI.
// It covers listing, creating, updating, and deleting price alerts.
package alert

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/alert/jsontypes"
	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// AlertContext is a client for the Longbridge Price Alert OpenAPI.
//
// Example:
//
//	conf, err := config.NewFormEnv()
//	actx, err := alert.NewFromCfg(conf)
//	list, err := actx.List(context.Background())
type AlertContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates an AlertContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*AlertContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &AlertContext{httpClient: httpClient}, nil
}

// NewFromEnv returns an AlertContext configured from environment variables.
func NewFromEnv() (*AlertContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// List returns all price alerts for the authenticated user.
//
// Path: GET /v1/notify/reminders
func (c *AlertContext) List(ctx context.Context) (*AlertList, error) {
	var resp jsontypes.AlertList
	if err := c.httpClient.Get(ctx, "/v1/notify/reminders", url.Values{}, &resp); err != nil {
		return nil, err
	}
	return convertAlertList(&resp), nil
}

// Add creates a new price alert.
//
// Path: POST /v1/notify/reminders
func (c *AlertContext) Add(
	ctx context.Context,
	symbol string,
	condition AlertCondition,
	triggerValue string,
	frequency AlertFrequency,
) error {
	var key string
	switch condition {
	case AlertConditionPriceRise, AlertConditionPriceFall:
		key = "price"
	case AlertConditionPercentRise, AlertConditionPercentFall:
		key = "chg"
	default:
		return fmt.Errorf("alert: unknown condition: %d", condition)
	}

	body := map[string]interface{}{
		"symbol":       symbol,
		"indicator_id": strconv.Itoa(int(condition)),
		"value_map":    map[string]string{key: triggerValue},
		"frequency":    int(frequency),
		"enabled":      true,
		"scope":        0,
		"state":        []int{1},
	}
	return c.httpClient.Post(ctx, "/v1/notify/reminders", body, nil)
}

// Update modifies an existing price alert in-place.
//
// Requires the AlertItem obtained from List. Set item.Enabled to true to
// re-enable or false to disable. All required fields are taken directly
// from the item — no extra round-trip needed.
//
// Path: POST /v1/notify/reminders
func (c *AlertContext) Update(ctx context.Context, item *AlertItem) error {
	body := map[string]interface{}{
		"id":           item.ID,
		"indicator_id": item.IndicatorID,
		"frequency":    item.Frequency,
		"scope":        item.Scope,
		"state":        item.State,
		"value_map":    item.ValueMap,
		"enabled":      item.Enabled,
	}
	return c.httpClient.Post(ctx, "/v1/notify/reminders", body, nil)
}

// Delete removes one or more price alerts by their IDs.
//
// Path: DELETE /v1/notify/reminders
func (c *AlertContext) Delete(ctx context.Context, alertIDs []string) error {
	body := map[string]interface{}{
		"ids": alertIDs,
	}
	return c.httpClient.Call(ctx, "DELETE", "/v1/notify/reminders", nil, body, nil)
}

// --- internal converters ---

func convertAlertList(j *jsontypes.AlertList) *AlertList {
	out := &AlertList{
		Lists: make([]*AlertSymbolGroup, 0, len(j.Lists)),
	}
	for _, g := range j.Lists {
		out.Lists = append(out.Lists, convertAlertSymbolGroup(g))
	}
	return out
}

func convertAlertSymbolGroup(j *jsontypes.AlertSymbolGroup) *AlertSymbolGroup {
	g := &AlertSymbolGroup{
		Symbol:     j.Symbol,
		Code:       j.Code,
		Market:     j.Market,
		Name:       j.Name,
		Price:      j.Price,
		Chg:        j.Chg,
		PChg:       j.PChg,
		Product:    j.Product,
		Indicators: make([]*AlertItem, 0, len(j.Indicators)),
	}
	for _, item := range j.Indicators {
		g.Indicators = append(g.Indicators, convertAlertItem(item))
	}
	return g
}

func convertAlertItem(j *jsontypes.AlertItem) *AlertItem {
	state := j.State
	if state == nil {
		state = []int{}
	}
	return &AlertItem{
		ID:          j.ID,
		IndicatorID: j.IndicatorID,
		Enabled:     j.Enabled,
		Frequency:   j.Frequency,
		Scope:       j.Scope,
		Text:        j.Text,
		State:       state,
		ValueMap:    j.ValueMap,
	}
}
