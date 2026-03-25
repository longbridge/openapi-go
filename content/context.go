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

// TopicsMine returns the topics created by the current authenticated user.
// Path: GET /v1/content/topics/mine
func (c *ContentContext) TopicsMine(ctx context.Context, opts *TopicsMineOptions) (items []*OwnedTopic, err error) {
	params := url.Values{}
	if opts != nil {
		if opts.Page > 0 {
			params.Set("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.Size > 0 {
			params.Set("size", fmt.Sprintf("%d", opts.Size))
		}
		if opts.TopicType != "" {
			params.Set("topic_type", opts.TopicType)
		}
	}
	var resp jsontypes.OwnedTopicList
	err = c.httpClient.Get(ctx, "/v1/content/topics/mine", params, &resp)
	if err != nil {
		return
	}
	items = make([]*OwnedTopic, 0, len(resp.Items))
	for _, item := range resp.Items {
		mi := &OwnedTopic{
			Id:            item.Id,
			Title:         item.Title,
			Description:   item.Description,
			Body:          item.Body,
			Tickers:       item.Tickers,
			Hashtags:      item.Hashtags,
			LikesCount:    item.LikesCount,
			CommentsCount: item.CommentsCount,
			ViewsCount:    item.ViewsCount,
			SharesCount:   item.SharesCount,
			TopicType:     item.TopicType,
			License:       item.License,
			DetailUrl:     item.DetailUrl,
			CreatedAt:     time.Unix(item.CreatedAt, 0).UTC(),
			UpdatedAt:     time.Unix(item.UpdatedAt, 0).UTC(),
		}
		if item.Author != nil {
			mi.Author = &Author{
				MemberId: item.Author.MemberId,
				Name:     item.Author.Name,
				Avatar:   item.Author.Avatar,
			}
		}
		if len(item.Images) > 0 {
			mi.Images = make([]*Image, 0, len(item.Images))
			for _, img := range item.Images {
				mi.Images = append(mi.Images, &Image{Url: img.Url, Sm: img.Sm, Lg: img.Lg})
			}
		}
		items = append(items, mi)
	}
	return
}

// CreateTopic creates a new topic for the current authenticated user.
// Path: POST /v1/content/topics
func (c *ContentContext) CreateTopic(ctx context.Context, opts *CreateTopicOptions) (*OwnedTopic, error) {
	req := &jsontypes.CreateTopicRequest{
		Title:     opts.Title,
		Body:      opts.Body,
		TopicType: opts.TopicType,
		Tickers:   opts.Tickers,
		Hashtags:  opts.Hashtags,
		License:   opts.License,
	}
	var resp jsontypes.CreateTopicResponse
	err := c.httpClient.Post(ctx, "/v1/content/topics", req, &resp)
	if err != nil {
		return nil, err
	}
	item := resp.Item
	if item == nil {
		return nil, nil
	}
	result := &OwnedTopic{
		Id:            item.Id,
		Title:         item.Title,
		Description:   item.Description,
		Body:          item.Body,
		Tickers:       item.Tickers,
		Hashtags:      item.Hashtags,
		LikesCount:    item.LikesCount,
		CommentsCount: item.CommentsCount,
		ViewsCount:    item.ViewsCount,
		SharesCount:   item.SharesCount,
		TopicType:     item.TopicType,
		License:       item.License,
		DetailUrl:     item.DetailUrl,
		CreatedAt:     time.Unix(item.CreatedAt, 0).UTC(),
		UpdatedAt:     time.Unix(item.UpdatedAt, 0).UTC(),
	}
	if item.Author != nil {
		result.Author = &Author{
			MemberId: item.Author.MemberId,
			Name:     item.Author.Name,
			Avatar:   item.Author.Avatar,
		}
	}
	if len(item.Images) > 0 {
		result.Images = make([]*Image, 0, len(item.Images))
		for _, img := range item.Images {
			result.Images = append(result.Images, &Image{Url: img.Url, Sm: img.Sm, Lg: img.Lg})
		}
	}
	return result, nil
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
