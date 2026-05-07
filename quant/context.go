// Package quant provides a client for the Longbridge Quant script execution API.
package quant

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// QuantContext is a client for executing quantitative scripts.
type QuantContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a QuantContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*QuantContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &QuantContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a QuantContext configured from environment variables.
func NewFromEnv() (*QuantContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// RunQuantScriptOptions contains parameters for running a quant script.
type RunQuantScriptOptions struct {
	Symbol string          `json:"symbol"`
	Period *string         `json:"period,omitempty"`
	Script string          `json:"script"`
	Input  json.RawMessage `json:"input,omitempty"`
}

// RunQuantScript executes a quantitative script.
//
// Reference: POST /v1/quant/run_script
func (c *QuantContext) RunQuantScript(ctx context.Context, opts *RunQuantScriptOptions) (json.RawMessage, error) {
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, "/v1/quant/run_script", opts, &resp)
	return resp, err
}
