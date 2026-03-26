package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	nhttp "net/http"
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
	if c.opts.OAuthClient != nil {
		token, err := c.opts.OAuthClient.AccessToken(ctx)
		if err != nil {
			return err
		}
		appKey = c.opts.OAuthClient.ClientID()
		accessToken = "Bearer " + token
		appSecret = ""
	}

	// set headers
	req.Header.Add("accept-language", string(c.opts.Language))
	req.Header.Add("x-api-key", appKey)
	req.Header.Add("authorization", accessToken)
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
	req.Close = true
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

	cli := &nhttp.Client{Timeout: opts.Timeout}

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
