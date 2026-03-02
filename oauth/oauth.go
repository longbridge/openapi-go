// Package oauth provides OAuth 2.0 authentication support for LongPort OpenAPI.
//
// This package implements the OAuth 2.0 authorization code flow to obtain
// access tokens for API authentication.
//
// # Example
//
//	o := oauth.New("your-client-id").
//	    OnOpenURL(func(url string) {
//	        // Open the URL however you like, e.g. print it or launch a browser
//	        fmt.Println("Please visit:", url)
//	    })
//	token, err := o.Authorize(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	cfg := config.FromOAuth(o.ClientID(), token.AccessToken)
package oauth

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

const (
	authTimeout         = 5 * time.Minute
	oauthBaseURL        = "https://openapi.longportapp.com"
	defaultCallbackPort = 60355
)

// OAuthToken holds an OAuth 2.0 access token with expiration and refresh info.
type OAuthToken struct {
	// AccessToken is the Bearer token for API authentication.
	AccessToken string `json:"access_token"`
	// RefreshToken can be used to obtain a new access token.
	RefreshToken string `json:"refresh_token,omitempty"`
	// ExpiresAt is the Unix timestamp (seconds) when the token expires.
	ExpiresAt int64 `json:"expires_at"`
}

// IsExpired reports whether the token has expired.
func (t *OAuthToken) IsExpired() bool {
	return time.Now().Unix() >= t.ExpiresAt
}

// ExpiresSoon reports whether the token will expire within 1 hour.
func (t *OAuthToken) ExpiresSoon() bool {
	return t.ExpiresAt-time.Now().Unix() < 3600
}

// OAuth is the OAuth 2.0 client for LongPort OpenAPI.
type OAuth struct {
	clientID     string
	callbackPort int
	baseURL      string
	openURL      func(string)
}

// New creates a new OAuth client with the given client ID.
//
// The client ID is obtained from the LongPort developer portal.
func New(clientID string) *OAuth {
	return &OAuth{
		clientID:     clientID,
		callbackPort: defaultCallbackPort,
		baseURL:      oauthBaseURL,
	}
}

// NewWithBaseURL creates a new OAuth client with a custom base URL.
// This is primarily intended for testing.
func NewWithBaseURL(clientID, baseURL string) *OAuth {
	return &OAuth{
		clientID:     clientID,
		callbackPort: defaultCallbackPort,
		baseURL:      baseURL,
	}
}

// WithCallbackPort sets the local callback port (default: 60355).
func (o *OAuth) WithCallbackPort(port int) *OAuth {
	o.callbackPort = port
	return o
}

// OnOpenURL sets a callback to handle the authorization URL.
//
// The callback receives the authorization URL as a string. Use this to
// open the URL in a browser, print it, or handle it in any other way
// appropriate for your application.
//
// If not set, the URL is silently discarded (useful for testing or when
// you retrieve the URL through other means).
func (o *OAuth) OnOpenURL(f func(string)) *OAuth {
	o.openURL = f
	return o
}

// ClientID returns the OAuth 2.0 client ID.
func (o *OAuth) ClientID() string {
	return o.clientID
}

// Authorize starts the OAuth 2.0 authorization code flow.
//
// It will:
//  1. Start a local HTTP server to receive the callback
//  2. Invoke the OnOpenURL callback with the authorization URL, so the
//     caller can open it in a browser or handle it in any other way
//  3. Wait for the user to authorize and receive the authorization code
//  4. Exchange the code for an access token
//
// The context can be used to cancel the flow before it completes.
func (o *OAuth) Authorize(ctx context.Context) (*OAuthToken, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", o.callbackPort))
	if err != nil {
		return nil, fmt.Errorf("oauth: failed to bind callback server: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", port)

	cfg := o.oauth2Config(redirectURI)
	state := oauth2.GenerateVerifier() // cryptographically random

	authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Invoke caller-supplied callback with the authorization URL
	if o.openURL != nil {
		o.openURL(authURL)
	}

	type result struct {
		code  string
		state string
		err   error
	}
	ch := make(chan result, 1)
	go func() {
		code, st, err := waitForCallback(listener)
		ch <- result{code, st, err}
	}()

	timeoutCtx, cancel := context.WithTimeout(ctx, authTimeout)
	defer cancel()

	var res result
	select {
	case res = <-ch:
	case <-timeoutCtx.Done():
		return nil, fmt.Errorf("oauth: authorization timeout — no response within 5 minutes")
	}
	if res.err != nil {
		return nil, res.err
	}
	if res.state != state {
		return nil, fmt.Errorf("oauth: CSRF state mismatch")
	}

	t, err := cfg.Exchange(ctx, res.code)
	if err != nil {
		return nil, fmt.Errorf("oauth: failed to exchange code for token: %w", err)
	}
	return tokenFromOAuth2(t), nil
}

// Refresh obtains a new access token using a refresh token.
func (o *OAuth) Refresh(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	cfg := o.oauth2Config("")
	src := cfg.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	t, err := src.Token()
	if err != nil {
		return nil, fmt.Errorf("oauth: failed to refresh token: %w", err)
	}
	tok := tokenFromOAuth2(t)
	if tok.RefreshToken == "" {
		tok.RefreshToken = refreshToken
	}
	return tok, nil
}

// ── internal helpers ──────────────────────────────────────────────────────────

func (o *OAuth) oauth2Config(redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID: o.clientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:   o.baseURL + "/oauth2/authorize",
			TokenURL:  o.baseURL + "/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams, // send client_id in POST body, not Basic Auth
		},
		RedirectURL: redirectURI,
	}
}

func tokenFromOAuth2(t *oauth2.Token) *OAuthToken {
	expiresAt := t.Expiry.Unix()
	if t.Expiry.IsZero() {
		expiresAt = time.Now().Unix() + 3600
	}
	rt, _ := t.Extra("refresh_token").(string)
	if rt == "" {
		rt = t.RefreshToken
	}
	return &OAuthToken{
		AccessToken:  t.AccessToken,
		RefreshToken: rt,
		ExpiresAt:    expiresAt,
	}
}

const callbackStyle = `<style>html{font-family:system-ui,-apple-system,BlinkMacSystemFont,sans-serif;` +
	`font-size:16px;color:#e0e0e0;background:#202020;padding:2rem;text-align:center;}</style>`

func waitForCallback(listener net.Listener) (code, state string, err error) {
	conn, err := listener.Accept()
	if err != nil {
		return "", "", fmt.Errorf("oauth: callback server error: %w", err)
	}
	defer conn.Close()

	// Read only the request line; we don't need headers.
	requestLine, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("oauth: failed to read callback request: %w", err)
	}

	// e.g. "GET /callback?code=xxx&state=yyy HTTP/1.1"
	parts := splitFields(requestLine)
	if len(parts) < 2 {
		writeHTML(conn, 400, "Bad Request", "<h1>Bad Request</h1>")
		return "", "", fmt.Errorf("oauth: malformed callback request")
	}

	u, err := url.Parse("http://localhost" + parts[1])
	if err != nil {
		writeHTML(conn, 400, "Bad Request", "<h1>Invalid URL</h1>")
		return "", "", fmt.Errorf("oauth: invalid callback URL: %w", err)
	}

	q := u.Query()
	if oauthErr := q.Get("error"); oauthErr != "" {
		writeHTML(conn, 400, "Authorization Failed",
			fmt.Sprintf("<h1>Authorization Failed</h1><p>Error: %s</p>", oauthErr))
		return "", "", fmt.Errorf("oauth: authorization failed: %s", oauthErr)
	}

	code, state = q.Get("code"), q.Get("state")
	if code == "" || state == "" {
		writeHTML(conn, 400, "Missing Parameters",
			"<h1>Missing Parameters</h1><p>Authorization code or state not received</p>")
		return "", "", fmt.Errorf("oauth: missing code or state in callback")
	}

	writeHTML(conn, 200, "OK",
		"<h1>✓ Authorization Successful!</h1><p>You can close this window and return to the terminal.</p>")
	return code, state, nil
}

func writeHTML(conn net.Conn, status int, statusText, body string) {
	html := fmt.Sprintf("<html><body>%s%s</body></html>", callbackStyle, body)
	fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s",
		status, statusText, html)
}

func splitFields(s string) []string {
	var fields []string
	start := -1
	for i, c := range s {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			if start >= 0 {
				fields = append(fields, s[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		fields = append(fields, s[start:])
	}
	return fields
}
