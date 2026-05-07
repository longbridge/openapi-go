// Package alert provides a client for the Longbridge Alert (price reminder) API.
package alert

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// AlertContext is a client for managing price alerts (reminders).
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

// ListAlerts returns all price alerts.
//
// Reference: GET /v1/notify/reminders
func (c *AlertContext) ListAlerts(ctx context.Context) (json.RawMessage, error) {
	var resp json.RawMessage
	err := c.httpClient.Get(ctx, "/v1/notify/reminders", nil, &resp)
	return resp, err
}

// AddAlertOptions contains parameters for creating a new price alert.
type AddAlertOptions struct {
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Direction string `json:"direction"`
	Remark    string `json:"remark,omitempty"`
}

// AddAlert creates a new price alert.
//
// Reference: POST /v1/notify/reminders
func (c *AlertContext) AddAlert(ctx context.Context, opts *AddAlertOptions) (json.RawMessage, error) {
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, "/v1/notify/reminders", opts, &resp)
	return resp, err
}

// DeleteAlerts deletes price alerts by IDs.
//
// Reference: DELETE /v1/notify/reminders
func (c *AlertContext) DeleteAlerts(ctx context.Context, ids []string) error {
	values := url.Values{}
	values.Add("ids", strings.Join(ids, ","))
	return c.httpClient.Delete(ctx, "/v1/notify/reminders", values, nil)
}

// EnableAlert enables a price alert by ID.
//
// Reference: PUT /v1/notify/reminders
func (c *AlertContext) EnableAlert(ctx context.Context, id string) (json.RawMessage, error) {
	body := map[string]interface{}{
		"id":      id,
		"enabled": true,
	}
	var resp json.RawMessage
	err := c.httpClient.Put(ctx, "/v1/notify/reminders", body, &resp)
	return resp, err
}

// DisableAlert disables a price alert by ID.
//
// Reference: PUT /v1/notify/reminders
func (c *AlertContext) DisableAlert(ctx context.Context, id string) (json.RawMessage, error) {
	body := map[string]interface{}{
		"id":      id,
		"enabled": false,
	}
	var resp json.RawMessage
	err := c.httpClient.Put(ctx, "/v1/notify/reminders", body, &resp)
	return resp, err
}
