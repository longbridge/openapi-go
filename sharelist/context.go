package sharelist

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/sharelist/jsontypes"
)

// SharelistContext is a client for sharelist (watchlist collections) management.
type SharelistContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a SharelistContext from a Config.
func NewFromCfg(cfg *config.Config) (*SharelistContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &SharelistContext{httpClient: httpClient}, nil
}

// NewFromEnv creates a SharelistContext from environment variables.
func NewFromEnv() (*SharelistContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

func convertSharelistInfo(info *jsontypes.SharelistInfo) *SharelistInfo {
	si := &SharelistInfo{
		Id:               info.Id,
		Name:             info.Name,
		Description:      info.Description,
		Cover:            info.Cover,
		SubscribersCount: info.SubscribersCount,
		CreatedAt:        info.CreatedAt,
		EditedAt:         info.EditedAt,
		ThisYearChg:      info.ThisYearChg,
		Subscribed:       info.Subscribed,
		Chg:              info.Chg,
		SharelistType:    info.SharelistType,
		IndustryCode:     info.IndustryCode,
	}
	for _, s := range info.Stocks {
		si.Stocks = append(si.Stocks, &SharelistStock{
			Symbol:                  util.CounterIDToSymbol(s.CounterId),
			Name:                    s.Name,
			Market:                  s.Market,
			Code:                    s.Code,
			Intro:                   s.Intro,
			UnreadChangeLogCategory: s.UnreadChangeLogCategory,
			Change:                  s.Change,
			LastDone:                s.LastDone,
			TradeStatus:             s.TradeStatus,
			Latency:                 s.Latency,
		})
	}
	return si
}

// List returns the user's own sharelists and subscribed sharelists.
func (c *SharelistContext) List(ctx context.Context, count int32) (*SharelistList, error) {
	values := url.Values{}
	values.Add("size", strconv.FormatInt(int64(count), 10))
	values.Add("self", "true")
	values.Add("subscription", "true")
	var resp jsontypes.SharelistList
	if err := c.httpClient.Get(ctx, "/v1/sharelists", values, &resp); err != nil {
		return nil, err
	}
	result := &SharelistList{TailMark: resp.TailMark}
	for _, info := range resp.Sharelists {
		result.Sharelists = append(result.Sharelists, convertSharelistInfo(info))
	}
	for _, info := range resp.SubscribedSharelists {
		result.SubscribedSharelists = append(result.SubscribedSharelists, convertSharelistInfo(info))
	}
	return result, nil
}

// Detail returns the detail of a sharelist by ID, including constituent securities.
func (c *SharelistContext) Detail(ctx context.Context, id int64) (*SharelistDetail, error) {
	values := url.Values{}
	values.Add("constituent", "true")
	values.Add("quote", "true")
	values.Add("subscription", "true")
	var resp jsontypes.SharelistDetail
	if err := c.httpClient.Get(ctx, fmt.Sprintf("/v1/sharelists/%d", id), values, &resp); err != nil {
		return nil, err
	}
	result := &SharelistDetail{}
	if resp.Sharelist != nil {
		result.Sharelist = convertSharelistInfo(resp.Sharelist)
	}
	if resp.Scopes != nil {
		result.Scopes = &SharelistScopes{
			Subscription: resp.Scopes.Subscription,
			IsSelf:       resp.Scopes.IsSelf,
		}
	}
	return result, nil
}

// Popular returns the most popular sharelists.
func (c *SharelistContext) Popular(ctx context.Context, count int32) ([]*SharelistInfo, error) {
	values := url.Values{}
	values.Add("size", strconv.FormatInt(int64(count), 10))
	var resp jsontypes.SharelistList
	if err := c.httpClient.Get(ctx, "/v1/sharelists/popular", values, &resp); err != nil {
		return nil, err
	}
	var result []*SharelistInfo
	for _, info := range resp.Sharelists {
		result = append(result, convertSharelistInfo(info))
	}
	return result, nil
}

// Create creates a new sharelist.
func (c *SharelistContext) Create(ctx context.Context, name, description string) (*SharelistInfo, error) {
	body := map[string]interface{}{
		"name":        name,
		"description": description,
	}
	var resp jsontypes.SharelistDetail
	if err := c.httpClient.Post(ctx, "/v1/sharelists", body, &resp); err != nil {
		return nil, err
	}
	if resp.Sharelist == nil {
		return nil, nil
	}
	return convertSharelistInfo(resp.Sharelist), nil
}

// Delete removes a sharelist by ID.
func (c *SharelistContext) Delete(ctx context.Context, id int64) error {
	return c.httpClient.Delete(ctx, fmt.Sprintf("/v1/sharelists/%d", id), nil, nil)
}

// AddSecurities adds securities to a sharelist.
func (c *SharelistContext) AddSecurities(ctx context.Context, id int64, symbols []string) error {
	counterIds := make([]string, len(symbols))
	for i, s := range symbols {
		counterIds[i] = util.SymbolToCounterID(s)
	}
	return c.httpClient.Post(ctx, fmt.Sprintf("/v1/sharelists/%d/items", id), map[string]interface{}{
		"counter_ids": strings.Join(counterIds, ","),
	}, nil)
}

// RemoveSecurities removes securities from a sharelist.
func (c *SharelistContext) RemoveSecurities(ctx context.Context, id int64, symbols []string) error {
	counterIds := make([]string, len(symbols))
	for i, s := range symbols {
		counterIds[i] = util.SymbolToCounterID(s)
	}
	values := url.Values{}
	values.Add("counter_ids", strings.Join(counterIds, ","))
	return c.httpClient.Delete(ctx, fmt.Sprintf("/v1/sharelists/%d/items", id), values, nil)
}

// SortSecurities reorders the securities in a sharelist.
func (c *SharelistContext) SortSecurities(ctx context.Context, id int64, symbols []string) error {
	counterIds := make([]string, len(symbols))
	for i, s := range symbols {
		counterIds[i] = util.SymbolToCounterID(s)
	}
	return c.httpClient.Post(ctx, fmt.Sprintf("/v1/sharelists/%d/items/sort", id), map[string]interface{}{
		"counter_ids": strings.Join(counterIds, ","),
	}, nil)
}
