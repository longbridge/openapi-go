// Package asset provides a client for the Longbridge Asset OpenAPI.
// It covers account asset queries and downloads.
package asset

import (
	"context"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/asset/jsontypes"
	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// AssetContext is a client for the Longbridge Asset API.
//
// Example:
//
//	conf, err := config.New()
//	sctx, err := asset.NewFromCfg(conf)
//	items, err := sctx.Statements(context.Background(), &asset.GetStatementList{
//	    StatementType: asset.StatementTypeDaily,
//	    PageSize: 5,
//	})
type AssetContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a AssetContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*AssetContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &AssetContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a AssetContext configured from environment variables.
func NewFromEnv() (*AssetContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// Statements returns the list of available statements.
//
// Reference: GET /v1/statement/list
func (c *AssetContext) Statements(ctx context.Context, params *GetStatementList) (items []*StatementItem, err error) {
	var resp jsontypes.StatementListResponse
	err = c.httpClient.Get(ctx, "/v1/statement/list", params.Values(), &resp)
	if err != nil {
		return
	}
	items = make([]*StatementItem, 0, len(resp.List))
	for _, item := range resp.List {
		items = append(items, &StatementItem{
			Date:    item.Dt,
			FileKey: item.FileKey,
		})
	}
	return
}

// StatementDownloadURL returns a presigned URL to download the statement JSON.
//
// Reference: GET /v1/statement/download
func (c *AssetContext) StatementDownloadURL(ctx context.Context, params *GetStatementDownloadURL) (downloadURL string, err error) {
	var resp jsontypes.StatementDownloadURLResponse
	err = c.httpClient.Get(ctx, "/v1/statement/download", params.Values(), &resp)
	if err != nil {
		return
	}
	downloadURL = resp.URL
	return
}
