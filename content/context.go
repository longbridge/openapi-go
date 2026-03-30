// Package content provides a client for the Longbridge Content OpenAPI.
// It covers community topics, replies, and related data.
package content

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/content/jsontypes"
	httplib "github.com/longbridge/openapi-go/http"
)

// ContentContext is a client for the Longbridge Content OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	cctx, err := content.NewFromCfg(conf)
//	topics, err := cctx.MyTopics(context.Background(), &content.MyTopicsOptions{Size: 20})
type ContentContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a ContentContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*ContentContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
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

// TopicDetail returns the full details of a topic by ID.
//
// See: https://open.longbridge.com/docs/api?op=topic_detail
func (c *ContentContext) TopicDetail(ctx context.Context, id string) (*OwnedTopic, error) {
	resp := &jsontypes.TopicDetailResponse{}
	if err := c.httpClient.Get(ctx, "/v1/content/topics/"+id, nil, resp); err != nil {
		return nil, err
	}
	return convertOwnedTopic(&resp.Item)
}

// MyTopics returns topics created by the currently authenticated user.
//
// See: https://open.longbridge.com/docs/api?op=list_my_topics
func (c *ContentContext) MyTopics(ctx context.Context, opts *MyTopicsOptions) ([]*OwnedTopic, error) {
	resp := &struct {
		Items []*jsontypes.OwnedTopic `json:"items"`
	}{}
	if err := c.httpClient.Get(ctx, "/v1/content/topics/mine", opts.values(), resp); err != nil {
		return nil, err
	}
	out := make([]*OwnedTopic, 0, len(resp.Items))
	for _, item := range resp.Items {
		t, err := convertOwnedTopic(item)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

// CreateTopic publishes a new community topic and returns the new topic ID.
//
// See: https://open.longbridge.com/docs/api?op=create_topic
func (c *ContentContext) CreateTopic(ctx context.Context, opts *CreateTopicOptions) (string, error) {
	body := map[string]interface{}{
		"body": opts.Body,
	}
	if opts.Title != "" {
		body["title"] = opts.Title
	}
	if opts.TopicType != "" {
		body["topic_type"] = opts.TopicType
	}
	if len(opts.Tickers) > 0 {
		body["tickers"] = opts.Tickers
	}
	if len(opts.Hashtags) > 0 {
		body["hashtags"] = opts.Hashtags
	}

	resp := &struct {
		ID string `json:"id"`
	}{}
	wrapper := &struct {
		Item json.RawMessage `json:"item"`
	}{}
	if err := c.httpClient.Post(ctx, "/v1/content/topics", body, wrapper); err != nil {
		return "", err
	}
	if err := json.Unmarshal(wrapper.Item, resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

// ListTopicReplies returns a paginated list of replies for a topic.
//
// See: https://open.longbridge.com/docs/api?op=list_topic_replies
func (c *ContentContext) ListTopicReplies(ctx context.Context, topicID string, opts *ListTopicRepliesOptions) ([]*TopicReply, error) {
	resp := &jsontypes.TopicRepliesResponse{}
	path := fmt.Sprintf("/v1/content/topics/%s/comments", topicID)
	if err := c.httpClient.Get(ctx, path, opts.values(), resp); err != nil {
		return nil, err
	}
	out := make([]*TopicReply, 0, len(resp.Items))
	for _, item := range resp.Items {
		r, err := convertTopicReply(item)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

// CreateTopicReply posts a reply to a topic and returns the created reply.
//
// See: https://open.longbridge.com/docs/api?op=create_topic_reply
func (c *ContentContext) CreateTopicReply(ctx context.Context, topicID string, opts *CreateReplyOptions) (*TopicReply, error) {
	body := map[string]interface{}{
		"body": opts.Body,
	}
	if opts.ReplyToID != "" && opts.ReplyToID != "0" {
		body["reply_to_id"] = opts.ReplyToID
	}

	resp := &jsontypes.TopicReplyResponse{}
	path := fmt.Sprintf("/v1/content/topics/%s/comments", topicID)
	if err := c.httpClient.Post(ctx, path, body, resp); err != nil {
		return nil, err
	}
	return convertTopicReply(&resp.Item)
}

// --- internal converters ---

func convertOwnedTopic(j *jsontypes.OwnedTopic) (*OwnedTopic, error) {
	createdAt, err := parseUnixTimestamp(j.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("content: parse created_at: %w", err)
	}
	updatedAt, err := parseUnixTimestamp(j.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("content: parse updated_at: %w", err)
	}
	t := &OwnedTopic{
		ID:          j.ID,
		Title:       j.Title,
		Description: j.Description,
		Body:        j.Body,
		Author: TopicAuthor{
			MemberID: j.Author.MemberID,
			Name:     j.Author.Name,
			Avatar:   j.Author.Avatar,
		},
		Tickers:   j.Tickers,
		Hashtags:  j.Hashtags,
		TopicType: j.TopicType,
		DetailURL: j.DetailURL,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	t.LikesCount, _ = strconv.ParseInt(j.LikesCount, 10, 64)
	t.CommentsCount, _ = strconv.ParseInt(j.CommentsCount, 10, 64)
	t.ViewsCount, _ = strconv.ParseInt(j.ViewsCount, 10, 64)
	t.SharesCount, _ = strconv.ParseInt(j.SharesCount, 10, 64)
	for _, img := range j.Images {
		t.Images = append(t.Images, TopicImage{URL: img.URL, Sm: img.Sm, Lg: img.Lg})
	}
	return t, nil
}

func convertTopicReply(j *jsontypes.TopicReply) (*TopicReply, error) {
	createdAt, err := parseUnixTimestamp(j.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("content: parse created_at: %w", err)
	}
	r := &TopicReply{
		ID:        j.ID,
		TopicID:   j.TopicID,
		Body:      j.Body,
		ReplyToID: j.ReplyToID,
		Author: TopicAuthor{
			MemberID: j.Author.MemberID,
			Name:     j.Author.Name,
			Avatar:   j.Author.Avatar,
		},
		CreatedAt: createdAt,
	}
	r.LikesCount, _ = strconv.ParseInt(j.LikesCount, 10, 64)
	r.CommentsCount, _ = strconv.ParseInt(j.CommentsCount, 10, 64)
	for _, img := range j.Images {
		r.Images = append(r.Images, TopicImage{URL: img.URL, Sm: img.Sm, Lg: img.Lg})
	}
	return r, nil
}

func parseUnixTimestamp(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0), nil
}
