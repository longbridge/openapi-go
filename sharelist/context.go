// Package sharelist provides a client for the Longbridge Sharelist API.
package sharelist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// SharelistContext is a client for managing share-lists.
type SharelistContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a SharelistContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*SharelistContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &SharelistContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a SharelistContext configured from environment variables.
func NewFromEnv() (*SharelistContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// ListSharelists returns share-lists, optionally limited by count.
//
// Reference: GET /v1/sharelists
func (c *SharelistContext) ListSharelists(ctx context.Context, count *int) (json.RawMessage, error) {
	values := url.Values{}
	if count != nil {
		values.Add("count", fmt.Sprintf("%d", *count))
	}
	var resp json.RawMessage
	err := c.httpClient.Get(ctx, "/v1/sharelists", values, &resp)
	return resp, err
}

// SharelistDetail returns detail for a specific share-list.
//
// Reference: GET /v1/sharelists/{id}
func (c *SharelistContext) SharelistDetail(ctx context.Context, id string) (json.RawMessage, error) {
	var resp json.RawMessage
	err := c.httpClient.Get(ctx, fmt.Sprintf("/v1/sharelists/%s", id), nil, &resp)
	return resp, err
}

// CreateSharelistOptions contains parameters for creating a share-list.
type CreateSharelistOptions struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CreateSharelist creates a new share-list.
//
// Reference: POST /v1/sharelists
func (c *SharelistContext) CreateSharelist(ctx context.Context, opts *CreateSharelistOptions) (json.RawMessage, error) {
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, "/v1/sharelists", opts, &resp)
	return resp, err
}

// DeleteSharelist deletes a share-list by ID.
//
// Reference: DELETE /v1/sharelists/{id}
func (c *SharelistContext) DeleteSharelist(ctx context.Context, id string) error {
	return c.httpClient.Delete(ctx, fmt.Sprintf("/v1/sharelists/%s", id), nil, nil)
}

// AddSharelistItems adds symbols to a share-list.
//
// Reference: POST /v1/sharelists/{id}/items
func (c *SharelistContext) AddSharelistItems(ctx context.Context, id string, symbols []string) (json.RawMessage, error) {
	body := map[string]interface{}{"symbols": symbols}
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, fmt.Sprintf("/v1/sharelists/%s/items", id), body, &resp)
	return resp, err
}

// RemoveSharelistItems removes symbols from a share-list.
//
// Reference: DELETE /v1/sharelists/{id}/items
func (c *SharelistContext) RemoveSharelistItems(ctx context.Context, id string, symbols []string) error {
	values := url.Values{}
	values.Add("symbols", strings.Join(symbols, ","))
	return c.httpClient.Delete(ctx, fmt.Sprintf("/v1/sharelists/%s/items", id), values, nil)
}

// SortSharelistItems reorders items in a share-list.
//
// Reference: POST /v1/sharelists/{id}/items/sort
func (c *SharelistContext) SortSharelistItems(ctx context.Context, id string, symbols []string) (json.RawMessage, error) {
	body := map[string]interface{}{"symbols": symbols}
	var resp json.RawMessage
	err := c.httpClient.Post(ctx, fmt.Sprintf("/v1/sharelists/%s/items/sort", id), body, &resp)
	return resp, err
}

// PopularSharelists returns popular share-lists.
//
// Reference: GET /v1/sharelists/popular
func (c *SharelistContext) PopularSharelists(ctx context.Context, count *int) (json.RawMessage, error) {
	values := url.Values{}
	if count != nil {
		values.Add("count", fmt.Sprintf("%d", *count))
	}
	var resp json.RawMessage
	err := c.httpClient.Get(ctx, "/v1/sharelists/popular", values, &resp)
	return resp, err
}
