// Package content provides a client for the Longbridge Content OpenAPI.
// It covers community topics and replies.
package content

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/content/jsontypes"
	httplib "github.com/longportapp/openapi-go/http"
)

// ContentContext is a client for the Longbridge Content OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	cctx, err := content.NewFromCfg(conf)
//	topics, err := cctx.MyTopics(context.Background(), &content.ListMyTopicsOptions{Size: 20})
type ContentContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a ContentContext from a *config.Config.
func NewFromCfg(c *config.Config) (*ContentContext, error) {
	cli, err := httplib.NewFromCfg(c)
	if err != nil {
		return nil, err
	}
	return &ContentContext{httpClient: cli}, nil
}

// TopicDetail returns the full details of a topic by ID.
//
// Reference: GET /v1/content/topics/{id}
//
// Example:
//
//	topic, err := cctx.TopicDetail(ctx, "6993508780031016960")
func (c *ContentContext) TopicDetail(ctx context.Context, id string) (*OwnedTopic, error) {
	resp := &jsontypes.TopicDetailResponse{}
	if err := c.httpClient.Get(ctx, "/v1/content/topics/"+id, nil, resp); err != nil {
		return nil, err
	}
	return convertOwnedTopic(&resp.Item)
}

// MyTopics returns topics created by the currently authenticated user.
//
// Reference: GET /v1/content/topics/mine
//
// Example:
//
//	topics, err := cctx.MyTopics(ctx, &content.ListMyTopicsOptions{
//	  Page: 1, Size: 50, TopicType: "article",
//	})
func (c *ContentContext) MyTopics(ctx context.Context, opts *ListMyTopicsOptions) ([]*OwnedTopic, error) {
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
// Reference: POST /v1/content/topics
//
// Two content types are supported:
//   - "post" (default): plain text; Markdown is NOT rendered.
//   - "article": Markdown body (server converts to HTML); Title is required.
//
// Permission: user must hold a funded Longbridge account (HTTP 403 otherwise).
//
// Stock symbols mentioned in Body (e.g. "700.HK", "TSLA.US") are automatically
// recognized and linked by the platform. Use Tickers to associate additional
// symbols not explicitly mentioned in the body.
// WARNING: do not abuse symbol linking for unrelated stocks — content moderation
// may restrict publishing or suspend the account.
//
// Rate limit: max 3 topics/min and 10/24h per user (HTTP 429 on excess).
//
// Example:
//
//	id, err := cctx.CreateTopic(ctx, &content.CreateTopicOptions{
//	  Body:      "Bullish on 700.HK today",
//	  Tickers:   []string{"700.HK"},
//	  Hashtags:  []string{"hongkong"},
//	})
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
// Reference: GET /v1/content/topics/{topic_id}/comments
//
// Each reply includes author info, body (plain text), likes/replies counts, and a
// ReplyToID field: "0" indicates a top-level reply; any other value is a nested reply.
//
// Example:
//
//	replies, err := cctx.ListTopicReplies(ctx, "6993508780031016960",
//	  &content.ListTopicRepliesOptions{Page: 1, Size: 20},
//	)
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
// Reference: POST /v1/content/topics/{topic_id}/comments
//
// Permission: user must hold a funded Longbridge account (HTTP 403 otherwise).
//
// Body is plain text only — Markdown is not rendered.
// Stock symbols mentioned in Body (e.g. "700.HK", "TSLA.US") are automatically
// recognized and linked by the platform.
// WARNING: do not abuse symbol linking for unrelated stocks — content moderation
// may restrict publishing or suspend the account.
//
// Rate limit per user per topic: first 3 replies have no wait; subsequent replies
// require incrementally longer intervals (3 s → 5 s → 8 s → 13 s → 21 s → 34 s → 55 s cap).
// Exceeding the limit returns HTTP 429.
//
// To post a top-level reply, leave ReplyToID empty or set it to "0".
// To nest under an existing reply, set ReplyToID to that reply's ID.
//
// Example:
//
//	reply, err := cctx.CreateTopicReply(ctx, "6993508780031016960",
//	  &content.CreateReplyOptions{Body: "Great post!"},
//	)
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
