// Package sharelist provides a client for the Longbridge community sharelist API.
// It covers listing, creating, deleting, and managing securities in sharelists.
//
// Example:
//
//	cfg, err := config.NewFormEnv()
//	sctx, err := sharelist.NewFromCfg(cfg)
//	result, err := sctx.List(context.Background(), 20)
package sharelist

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/counter"
	"github.com/longbridge/openapi-go/sharelist/jsontypes"
)

// SharelistContext is a client for the Longbridge Sharelist API.
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

// List returns up to count of the user's own and subscribed sharelists.
//
// Path: GET /v1/sharelists
func (c *SharelistContext) List(ctx context.Context, count uint32) (*SharelistList, error) {
	params := url.Values{}
	params.Set("size", strconv.FormatUint(uint64(count), 10))
	params.Set("self", "true")
	params.Set("subscription", "true")

	var resp jsontypes.SharelistList
	if err := c.httpClient.Get(ctx, "/v1/sharelists", params, &resp); err != nil {
		return nil, err
	}
	return convertSharelistList(&resp)
}

// Detail returns the full information for a sharelist, including constituents
// and quotes.
//
// Path: GET /v1/sharelists/{id}
func (c *SharelistContext) Detail(ctx context.Context, id int64) (*SharelistDetail, error) {
	params := url.Values{}
	params.Set("constituent", "true")
	params.Set("quote", "true")
	params.Set("subscription", "true")

	path := fmt.Sprintf("/v1/sharelists/%d", id)
	var resp jsontypes.SharelistDetail
	if err := c.httpClient.Get(ctx, path, params, &resp); err != nil {
		return nil, err
	}
	info, err := convertSharelistInfo(&resp.Sharelist)
	if err != nil {
		return nil, err
	}
	return &SharelistDetail{
		Sharelist: *info,
		Scopes: SharelistScopes{
			Subscription: resp.Scopes.Subscription,
			IsSelf:       resp.Scopes.IsSelf,
		},
	}, nil
}

// Popular returns up to count trending sharelists.
//
// Path: GET /v1/sharelists/popular
func (c *SharelistContext) Popular(ctx context.Context, count uint32) (*SharelistList, error) {
	params := url.Values{}
	params.Set("size", strconv.FormatUint(uint64(count), 10))

	var resp jsontypes.SharelistList
	if err := c.httpClient.Get(ctx, "/v1/sharelists/popular", params, &resp); err != nil {
		return nil, err
	}
	return convertSharelistList(&resp)
}

// Create creates a new sharelist with the given name and optional description.
// When description is empty the name is used as the description.
//
// Path: POST /v1/sharelists
func (c *SharelistContext) Create(ctx context.Context, name string, description string) error {
	if description == "" {
		description = name
	}
	body := map[string]interface{}{
		"name":        name,
		"description": description,
		"cover":       "https://pub.pbkrs.com/files/202107/kaJSk6BsvPt6NJ3Q/sharelist_v1.png",
	}
	var resp interface{}
	return c.httpClient.Post(ctx, "/v1/sharelists", body, &resp)
}

// Delete deletes a sharelist by ID.
//
// Path: DELETE /v1/sharelists/{id}
func (c *SharelistContext) Delete(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/v1/sharelists/%d", id)
	var resp interface{}
	return c.httpClient.Delete(ctx, path, nil, &resp, httplib.WithBody(map[string]interface{}{}))
}

// AddSecurities adds one or more securities to a sharelist.
// Symbols should be in "CODE.MARKET" format, e.g. "TSLA.US" or "700.HK".
//
// Path: POST /v1/sharelists/{id}/items
func (c *SharelistContext) AddSecurities(ctx context.Context, id int64, symbols []string) error {
	counterIDs := symbolsToCounterIDs(symbols)
	path := fmt.Sprintf("/v1/sharelists/%d/items", id)
	body := map[string]interface{}{
		"counter_ids": counterIDs,
	}
	var resp interface{}
	return c.httpClient.Post(ctx, path, body, &resp)
}

// RemoveSecurities removes one or more securities from a sharelist.
// Symbols should be in "CODE.MARKET" format, e.g. "TSLA.US" or "700.HK".
//
// Path: DELETE /v1/sharelists/{id}/items
func (c *SharelistContext) RemoveSecurities(ctx context.Context, id int64, symbols []string) error {
	counterIDs := symbolsToCounterIDs(symbols)
	path := fmt.Sprintf("/v1/sharelists/%d/items", id)
	var resp interface{}
	return c.httpClient.Delete(ctx, path, nil, &resp, httplib.WithBody(map[string]interface{}{
		"counter_ids": counterIDs,
	}))
}

// SortSecurities reorders the securities in a sharelist.
// The symbols slice defines the desired order; each symbol should be in
// "CODE.MARKET" format.
//
// Path: POST /v1/sharelists/{id}/items/sort
func (c *SharelistContext) SortSecurities(ctx context.Context, id int64, symbols []string) error {
	counterIDs := symbolsToCounterIDs(symbols)
	path := fmt.Sprintf("/v1/sharelists/%d/items/sort", id)
	body := map[string]interface{}{
		"counter_ids": counterIDs,
	}
	var resp interface{}
	return c.httpClient.Post(ctx, path, body, &resp)
}

// --- symbol <-> counter_id helpers ---

// symbolToCounterID converts a symbol like "TSLA.US" to a counter_id like
// "ST/US/TSLA". This mirrors the Rust symbol_to_counter_id helper. Unlike the
// Rust implementation, ETF detection via an embedded CSV is not performed; all
// symbols are treated as equities (ST prefix).
func symbolToCounterID(symbol string) string {
	idx := strings.LastIndex(symbol, ".")
	if idx < 0 {
		return symbol
	}
	code := symbol[:idx]
	market := strings.ToUpper(symbol[idx+1:])
	return fmt.Sprintf("ST/%s/%s", market, code)
}

// symbolsToCounterIDs converts a slice of symbols to a comma-joined
// counter_ids string as expected by the API.
func symbolsToCounterIDs(symbols []string) string {
	ids := make([]string, 0, len(symbols))
	for _, s := range symbols {
		ids = append(ids, symbolToCounterID(s))
	}
	return strings.Join(ids, ",")
}


// --- internal converters ---

func convertSharelistList(j *jsontypes.SharelistList) (*SharelistList, error) {
	out := &SharelistList{
		TailMark: j.TailMark,
	}

	for i := range j.Sharelists {
		info, err := convertSharelistInfo(&j.Sharelists[i])
		if err != nil {
			return nil, err
		}
		out.Sharelists = append(out.Sharelists, *info)
	}
	for i := range j.SubscribedSharelists {
		info, err := convertSharelistInfo(&j.SubscribedSharelists[i])
		if err != nil {
			return nil, err
		}
		out.SubscribedSharelists = append(out.SubscribedSharelists, *info)
	}
	return out, nil
}

func convertSharelistInfo(j *jsontypes.SharelistInfo) (*SharelistInfo, error) {
	createdAt, err := parseUnixTimestamp(j.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("sharelist: parse created_at: %w", err)
	}
	editedAt, err := parseUnixTimestamp(j.EditedAt)
	if err != nil {
		return nil, fmt.Errorf("sharelist: parse edited_at: %w", err)
	}

	info := &SharelistInfo{
		ID:               parseID(j.ID),
		Name:             j.Name,
		Description:      j.Description,
		Cover:            j.Cover,
		SubscribersCount: j.SubscribersCount,
		CreatedAt:        createdAt,
		EditedAt:         editedAt,
		ThisYearChg:      parseOptionalDecimal(j.ThisYearChg),
		Subscribed:       j.Subscribed,
		Chg:              parseOptionalDecimal(j.Chg),
		SharelistType:    SharelistType(j.SharelistType),
		IndustryCode:     j.IndustryCode,
	}

	for i := range j.Stocks {
		info.Stocks = append(info.Stocks, convertSharelistStock(&j.Stocks[i]))
	}
	return info, nil
}

func convertSharelistStock(j *jsontypes.SharelistStock) SharelistStock {
	return SharelistStock{
		Symbol:                  counter.IDToSymbol(j.CounterID),
		Name:                    j.Name,
		Market:                  j.Market,
		Code:                    j.Code,
		Intro:                   j.Intro,
		UnreadChangeLogCategory: j.UnreadChangeLogCategory,
		Change:                  parseOptionalDecimal(j.Change),
		LastDone:                parseOptionalDecimal(j.LastDone),
		TradeStatus:             j.TradeStatus,
		Latency:                 j.Latency,
	}
}

// parseID handles the case where the API returns the sharelist ID as either a
// JSON number or a quoted string. Since json.Decoder.UseNumber() is active in
// the HTTP client, numbers arrive as json.Number; strings arrive as string.
func parseID(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch id := v.(type) {
	case float64:
		return int64(id)
	case string:
		n, _ := strconv.ParseInt(id, 10, 64)
		return n
	case int64:
		return id
	default:
		s := fmt.Sprintf("%v", v)
		n, _ := strconv.ParseInt(s, 10, 64)
		return n
	}
}

// parseOptionalDecimal converts an empty-or-numeric string to *decimal.Decimal.
// Returns nil when the string is empty or unparseable.
func parseOptionalDecimal(s string) *decimal.Decimal {
	if s == "" {
		return nil
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return &d
}

// parseUnixTimestamp parses a Unix timestamp string into a time.Time.
// An empty string yields the zero time without error.
func parseUnixTimestamp(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(sec, 0).UTC(), nil
}
