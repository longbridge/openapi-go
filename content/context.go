// Package content provides ContentContext for accessing Longbridge content APIs
package content

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/content/jsontypes"
)

// ContentContext is a client for interacting with Longbridge Content OpenAPI
type ContentContext struct {
	httpClient *http.Client
}

// NewFromCfg returns a ContentContext with config.Config.
func NewFromCfg(cfg *config.Config) (*ContentContext, error) {
	httpClient, err := http.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &ContentContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a ContentContext configured from environment variables.
func NewFromEnv() (*ContentContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// Topics returns the discussion topics list for a symbol.
// Reference: https://open.longbridge.com/en/docs/quote/security/topics
func (c *ContentContext) Topics(ctx context.Context, symbol string) (items []*TopicItem, err error) {
	var resp jsontypes.TopicList
	err = c.httpClient.Get(ctx, fmt.Sprintf("/v1/content/%s/topics", symbol), url.Values{}, &resp)
	if err != nil {
		return
	}
	items = make([]*TopicItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &TopicItem{
			Id:            item.Id,
			Title:         item.Title,
			Description:   item.Description,
			Url:           item.Url,
			PublishedAt:   time.Unix(item.PublishedAt, 0).UTC(),
			CommentsCount: item.CommentsCount,
			LikesCount:    item.LikesCount,
			SharesCount:   item.SharesCount,
		})
	}
	return
}

// News returns the news list for a symbol.
// Reference: https://open.longbridge.com/en/docs/quote/security/news
func (c *ContentContext) News(ctx context.Context, symbol string) (items []*NewsItem, err error) {
	var resp jsontypes.NewsList
	err = c.httpClient.Get(ctx, fmt.Sprintf("/v1/content/%s/news", symbol), url.Values{}, &resp)
	if err != nil {
		return
	}
	items = make([]*NewsItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &NewsItem{
			Id:            item.Id,
			Title:         item.Title,
			Description:   item.Description,
			Url:           item.Url,
			PublishedAt:   time.Unix(item.PublishedAt, 0).UTC(),
			CommentsCount: item.CommentsCount,
			LikesCount:    item.LikesCount,
			SharesCount:   item.SharesCount,
		})
	}
	return
}
