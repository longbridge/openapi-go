package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	nhttp "net/http"
	"time"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/log"
)

type apiResponse struct {
	Code    int
	Message string
	Data    json.RawMessage
	TraceID string
}

type otpResponse struct {
	Otp string
}

// Client is a http client to access Longbridge REST OpenAPI
type Client struct {
	opts       *Options
	httpClient *nhttp.Client
}

// RequestOptions use to set additional information for the request
type RequestOptions struct {
	// Request Header
	Header nhttp.Header
	body   interface{}
}

// RequestOption use to set addition info to request
type RequestOption func(*RequestOptions)

// WithHeader set request header
func WithHeader(h nhttp.Header) RequestOption {
	return func(o *RequestOptions) {
		o.Header = h
	}
}

// WithBody to set playload
func WithBody(v interface{}) RequestOption {
	return func(o *RequestOptions) {
		if v != nil {
			o.body = v
		}
	}
}

// region returns the data center region derived from the client's credentials.
func (c *Client) region() dcRegion {
	if c.opts.OAuthClient != nil {
		// OAuth token inspection is async; cannot determine synchronously.
		// Default to AP (the gateway default) so AP-only guards don't fire.
		return dcRegionAp
	}
	return dcRegionFromCredentials(c.opts.AppKey, c.opts.AppSecret, c.opts.AccessToken)
}

// IsUS reports whether the configured credentials belong to the US data center.
func (c *Client) IsUS() bool {
	return c.region() == dcRegionUs
}

// CheckRegion verifies that the current session's region matches required.
// Returns *RegionRestrictedError if it does not; nil otherwise.
// required must be "US" or "AP".
func (c *Client) CheckRegion(path, required string) error {
	var req dcRegion
	if required == "US" {
		req = dcRegionUs
	} else {
		req = dcRegionAp
	}
	cur := c.region()
	if cur.allows(req) {
		return nil
	}
	return &RegionRestrictedError{Path: path, Required: req.display(), Current: cur.display()}
}

// Get sends Get request with queryParams
func (c *Client) Get(ctx context.Context, path string, queryParams url.Values, resp interface{}, ropts ...RequestOption) error {
	return c.Call(ctx, "GET", path, queryParams, nil, resp, ropts...)
}

// Post sends Post request with json body
func (c *Client) Post(ctx context.Context, path string, body interface{}, resp interface{}, ropts ...RequestOption) error {
	return c.Call(ctx, "POST", path, nil, body, resp, ropts...)
}

// Put sends Put request with json body
func (c *Client) Put(ctx context.Context, path string, body interface{}, resp interface{}, ropts ...RequestOption) error {
	return c.Call(ctx, "PUT", path, nil, body, resp, ropts...)
}

// Delete sends Delete request with queryParams
func (c *Client) Delete(ctx context.Context, path string, queryParams interface{}, resp interface{}, ropts ...RequestOption) error {
	return c.Call(ctx, "DELETE", path, queryParams, nil, resp, ropts...)
}

// GetOTP to get one time password
// Reference: https://open.longbridge.com/en/docs/socket-token-api
func (c *Client) GetOTP(ctx context.Context, ropts ...RequestOption) (string, error) {
	res := &otpResponse{}
	err := c.Get(ctx, "/v1/socket/token", nil, res, ropts...)
	if err != nil {
		return "", err
	}
	return res.Otp, nil
}

func (c *Client) GetOTPV2(ctx context.Context, ropts ...RequestOption) (string, error) {
	res := &otpResponse{}
	err := c.Get(ctx, "/v2/socket/token", nil, res, ropts...)
	if err != nil {
		return "", err
	}
	return res.Otp, nil
}

// Call will send request with signature to http server
func (c *Client) Call(ctx context.Context, method, path string, queryParams interface{}, body interface{}, resp interface{}, ropts ...RequestOption) (err error) {
	var (
		br       io.Reader
		bb       []byte
		httpResp *nhttp.Response
		rb       []byte
	)

	ro := &RequestOptions{}
	for _, opt := range ropts {
		opt(ro)
	}

	if body == nil && ro.body != nil {
		body = ro.body
	}

	if body != nil {
		bb, err = json.Marshal(body)
		if err != nil {
			return err
		}
		br = bytes.NewBuffer(bb)
	}

	req, err := nhttp.NewRequestWithContext(ctx, method, c.opts.URL+path, br)
	if err != nil {
		return err
	}

	appKey := c.opts.AppKey
	accessToken := c.opts.AccessToken
	appSecret := c.opts.AppSecret
	var region dcRegion
	if c.opts.OAuthClient != nil {
		token, err := c.opts.OAuthClient.AccessToken(ctx)
		if err != nil {
			return err
		}
		// Derive DC region from the token prefix ("us_" → US, otherwise AP),
		// then strip the prefix so only the bare token is sent to the gateway.
		region = dcRegionFromCredential(token)
		appKey = c.opts.OAuthClient.ClientID()
		accessToken = "Bearer " + stripBearerPrefix(token)
		appSecret = ""
	} else {
		// API-key auth: any of the three credentials may carry the region prefix.
		region = dcRegionFromCredentials(appKey, appSecret, accessToken)
		accessToken = stripBearerPrefix(accessToken)
		appKey = stripBearerPrefix(appKey)
	}

	// set headers
	req.Header.Add("accept-language", string(c.opts.Language))
	req.Header.Add("x-api-key", appKey)
	req.Header.Add("authorization", accessToken)
	req.Header.Add(dcRegionHeader, region.asStr())
	for k, v := range c.opts.ExtraHeaders {
		req.Header.Set(k, v)
	}
	if ro.Header != nil {
		for k, v := range ro.Header {
			req.Header[k] = v
		}
	}
	if len(bb) != 0 {
		req.Header.Add("content-type", "application/json; charset=utf-8")
	}

	// set query params
	if queryParams != nil {
		vals, ok := queryParams.(url.Values)
		if !ok {
			if vals, err = query.Values(queryParams); err != nil {
				return
			}
		}
		req.URL.RawQuery = vals.Encode()
	}
	// set signature (no-op for OAuth when appSecret is empty)
	signature(req, appSecret, bb)

	log.Debugf("http call method:%v url:%v body:%v", req.Method, req.URL, string(bb))
	httpResp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}
	log.Debugf("http call response headers:%v", httpResp.Header)
	defer httpResp.Body.Close()

	if rb, err = io.ReadAll(httpResp.Body); err != nil {
		return err
	}
	log.Debugf("http call response body:%s", rb)

	apiResp := &apiResponse{}

	if v := httpResp.Header.Get("x-trace-id"); v != "" {
		apiResp.TraceID = v
	}

	if isJSON(httpResp.Header.Get("content-type")) {
		if err = jsonUnmarshal(bytes.NewReader(rb), apiResp); err != nil {
			return err
		}
	} else {
		apiResp.Message = string(rb)
	}

	if httpResp.StatusCode != nhttp.StatusOK || apiResp.Code != 0 {
		return NewError(httpResp.StatusCode, apiResp)
	}

	if resp == nil {
		return
	}

	if err = jsonUnmarshal(bytes.NewReader(apiResp.Data), resp); err != nil {
		return err
	}
	return nil
}

func isJSON(ct string) bool {
	return strings.Contains(ct, "application/json")
}

func jsonUnmarshal(r io.Reader, v interface{}) error {
	d := json.NewDecoder(r)
	d.UseNumber()
	return d.Decode(v)
}

// New create http client to call Longbridge REST OpenAPI
func New(opt ...Option) (*Client, error) {
	opts := newOptions(opt...)
	if opts.URL == "" {
		return nil, errors.New("http url is empty")
	}

	cli := &nhttp.Client{
		Timeout: opts.Timeout,
		Transport: &nhttp.Transport{
			IdleConnTimeout: 60 * time.Second,
		},
	}

	if opts.Client != nil {
		cli = opts.Client
	}

	client := &Client{
		opts:       opts,
		httpClient: cli,
	}
	return client, nil
}

// NewFromCfg init longbridge http client from *config.Config
func NewFromCfg(c *config.Config) (*Client, error) {
	url := c.HttpURL
	if url == "" {
		url = DefaultHttpUrl
	}
	opts := []Option{
		WithTimeout(c.HTTPTimeout),
		WithClient(c.Client),
		WithURL(url),
		WithLanguage(c.Language),
		WithExtraHeaders(c.ExtraHeaders),
	}
	if c.OAuthClient != nil {
		opts = append(opts, WithOAuthClient(c.OAuthClient))
	} else {
		opts = append(opts,
			WithAccessToken(c.AccessToken),
			WithAppKey(c.AppKey),
			WithAppSecret(c.AppSecret),
		)
	}
	return New(opts...)
}
