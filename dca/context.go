// Package dca provides a client for the Longbridge Dollar-Cost Averaging (DCA) API.
package dca

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// DcaContext is a client for managing DCA (Dollar-Cost Averaging) plans.
type DcaContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a DcaContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*DcaContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &DcaContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a DcaContext configured from environment variables.
func NewFromEnv() (*DcaContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// ListDcaPlans returns DCA plans, optionally filtered by status.
//
// Reference: GET /v1/dailycoins/query
func (c *DcaContext) ListDcaPlans(ctx context.Context, status *string) (json.RawMessage, error) {
	values := url.Values{}
	if status != nil {
		values.Add("status", *status)
	}
	var resp json.RawMessage
	err := c.httpClient.Get(ctx, "/v1/dailycoins/query", values, &resp)
	return resp, err
}

// CreateDcaPlanOptions contains parameters for creating a DCA plan.
type CreateDcaPlanOptions struct {
	Symbol    string `json:"symbol"`
	Amount    string `json:"amount"`
	Frequency string `json:"frequency"`
}

// CreateDcaPlan creates a new DCA plan.
//
// Reference: POST /v1/dailycoins/create
func (c *DcaContext) CreateDcaPlan(ctx context.Context, opts *CreateDcaPlanOptions) (json.RawMessage, error) {
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, "/v1/dailycoins/create", opts, &resp)
	return resp, err
}

// UpdateDcaPlanOptions contains parameters for updating a DCA plan.
type UpdateDcaPlanOptions struct {
	PlanID string                 `json:"plan_id"`
	Extra  map[string]interface{} `json:"extra,omitempty"`
}

// UpdateDcaPlan updates an existing DCA plan.
//
// Reference: POST /v1/dailycoins/update
func (c *DcaContext) UpdateDcaPlan(ctx context.Context, opts *UpdateDcaPlanOptions) (json.RawMessage, error) {
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, "/v1/dailycoins/update", opts, &resp)
	return resp, err
}

// PauseDcaPlan pauses a DCA plan.
//
// Reference: POST /v1/dailycoins/toggle
func (c *DcaContext) PauseDcaPlan(ctx context.Context, planID string) (json.RawMessage, error) {
	body := map[string]string{"plan_id": planID, "action": "pause"}
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, "/v1/dailycoins/toggle", body, &resp)
	return resp, err
}

// ResumeDcaPlan resumes a DCA plan.
//
// Reference: POST /v1/dailycoins/toggle
func (c *DcaContext) ResumeDcaPlan(ctx context.Context, planID string) (json.RawMessage, error) {
	body := map[string]string{"plan_id": planID, "action": "resume"}
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, "/v1/dailycoins/toggle", body, &resp)
	return resp, err
}

// StopDcaPlan stops a DCA plan.
//
// Reference: POST /v1/dailycoins/toggle
func (c *DcaContext) StopDcaPlan(ctx context.Context, planID string) (json.RawMessage, error) {
	body := map[string]string{"plan_id": planID, "action": "stop"}
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, "/v1/dailycoins/toggle", body, &resp)
	return resp, err
}

// DcaHistoryOptions contains parameters for querying DCA history.
type DcaHistoryOptions struct {
	PlanID    *string
	StartDate *string
	EndDate   *string
}

// DcaHistory returns execution history for DCA plans.
//
// Reference: GET /v1/dailycoins/query-records
func (c *DcaContext) DcaHistory(ctx context.Context, opts *DcaHistoryOptions) (json.RawMessage, error) {
	values := url.Values{}
	if opts != nil {
		if opts.PlanID != nil {
			values.Add("plan_id", *opts.PlanID)
		}
		if opts.StartDate != nil {
			values.Add("start_date", *opts.StartDate)
		}
		if opts.EndDate != nil {
			values.Add("end_date", *opts.EndDate)
		}
	}
	var resp json.RawMessage
	err := c.httpClient.Get(ctx, "/v1/dailycoins/query-records", values, &resp)
	return resp, err
}

// DcaStatistics returns DCA statistics, optionally filtered by symbol.
//
// Reference: GET /v1/dailycoins/statistic
func (c *DcaContext) DcaStatistics(ctx context.Context, symbol *string) (json.RawMessage, error) {
	values := url.Values{}
	if symbol != nil {
		values.Add("symbol", *symbol)
	}
	var resp json.RawMessage
	err := c.httpClient.Get(ctx, "/v1/dailycoins/statistic", values, &resp)
	return resp, err
}

// CheckDcaSupportOptions contains symbols to check for DCA support.
type CheckDcaSupportOptions struct {
	Symbols []string `json:"symbols"`
}

// CheckDcaSupport checks DCA support for a list of symbols.
//
// Reference: POST /v1/dailycoins/batch-check-support
func (c *DcaContext) CheckDcaSupport(ctx context.Context, opts *CheckDcaSupportOptions) (json.RawMessage, error) {
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, "/v1/dailycoins/batch-check-support", opts, &resp)
	return resp, err
}
