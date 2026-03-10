// Package oauth provides OAuth 2.0 authentication support for Longbridge OpenAPI.
//
// Use OAuth by setting it on config (like the Rust SDK). The token is stored
// internally and in ~/.longbridge-openapi/tokens/<client_id>; it is refreshed
// automatically when expired. Do not use or expose OAuthToken.
//
// # Example
//
//	o := oauth.New("your-client-id").
//	    OnOpenURL(func(url string) { fmt.Println("Please visit:", url) })
//	if err := o.Build(context.Background()); err != nil {
//	    log.Fatal(err)
//	}
//	cfg, _ := config.New(config.WithOAuthClient(o))
//	tctx, _ := trade.NewFromCfg(cfg)
package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

const (
	authTimeout         = 5 * time.Minute
	oauthBaseURL        = "https://openapi.longbridge.com"
	defaultCallbackPort = 60355
	expiresSoonSecs     = 3600
	tokenDir            = ".longbridge-openapi/tokens"
)

// oauthToken is the internal token type (not exported).
type oauthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresAt    int64  `json:"expires_at"`
}

func (t *oauthToken) isExpired() bool {
	return t == nil || time.Now().Unix() >= t.ExpiresAt
}

func (t *oauthToken) expiresSoon() bool {
	return t == nil || t.ExpiresAt-time.Now().Unix() < expiresSoonSecs
}

// OAuth is the OAuth 2.0 client for Longbridge OpenAPI.
// Use New and Build, then pass to config via config.WithOAuthClient(o).
// AccessToken(ctx) is used by the HTTP client to get a valid token (auto-refresh).
type OAuth struct {
	clientID     string
	callbackPort int
	baseURL      string
	openURL      func(string)

	mu    sync.Mutex
	token *oauthToken
}

// New creates a new OAuth client with the given client ID.
func New(clientID string) *OAuth {
	return &OAuth{
		clientID:     clientID,
		callbackPort: defaultCallbackPort,
		baseURL:      oauthBaseURL,
	}
}

// NewWithBaseURL creates a new OAuth client with a custom base URL (e.g. for tests).
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

// OnOpenURL sets a callback for the authorization URL (e.g. open in browser or print).
func (o *OAuth) OnOpenURL(f func(string)) *OAuth {
	o.openURL = f
	return o
}

// ClientID returns the OAuth 2.0 client ID.
func (o *OAuth) ClientID() string {
	return o.clientID
}

// Build loads a token from ~/.longbridge-openapi/tokens/<client_id> or runs the
// authorization flow, then stores the token in memory and on disk. Call this
// once before passing the OAuth to config. If a token exists and is expired,
// it is refreshed automatically.
func (o *OAuth) Build(ctx context.Context) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	path, err := o.tokenPath()
	if err != nil {
		return err
	}

	loaded, _ := loadTokenFromPath(path)
	if loaded != nil && !loaded.isExpired() {
		o.token = loaded
		return nil
	}
	if loaded != nil && loaded.RefreshToken != "" {
		o.mu.Unlock()
		refreshed, err := o.refreshToken(ctx, loaded)
		o.mu.Lock()
		if err == nil {
			_ = saveTokenToPath(path, refreshed)
			o.token = refreshed
			return nil
		}
	}

	// No valid token: run authorization flow.
	o.mu.Unlock()
	tok, err := o.authorizeFlow(ctx)
	if err != nil {
		return err
	}
	o.mu.Lock()
	o.token = tok
	_ = saveTokenToPath(path, tok)
	return nil
}

// AccessToken returns a valid access token, refreshing it if expired or soon to expire.
// The HTTP client calls this when OAuth is set on config.
func (o *OAuth) AccessToken(ctx context.Context) (string, error) {
	o.mu.Lock()
	tok := o.token
	needRefresh := tok == nil || tok.isExpired() || tok.expiresSoon()
	if needRefresh && tok != nil && tok.RefreshToken != "" {
		// Copy so we can refresh without holding the lock (network I/O).
		tokCopy := *tok
		o.mu.Unlock()
		refreshed, err := o.refreshToken(ctx, &tokCopy)
		if err != nil {
			return "", err
		}
		o.mu.Lock()
		o.token = refreshed
		path, _ := o.tokenPath()
		_ = saveTokenToPath(path, refreshed)
		o.mu.Unlock()
		return refreshed.AccessToken, nil
	}
	if needRefresh {
		o.mu.Unlock()
		return "", fmt.Errorf("oauth: no valid token; call Build(ctx) first or re-run authorization")
	}
	accessToken := tok.AccessToken
	o.mu.Unlock()
	return accessToken, nil
}

func (o *OAuth) tokenPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("oauth: cannot get home dir: %w", err)
	}
	return filepath.Join(home, tokenDir, o.clientID), nil
}

func (o *OAuth) refreshToken(ctx context.Context, tok *oauthToken) (*oauthToken, error) {
	cfg := o.oauth2Config("")
	src := cfg.TokenSource(ctx, &oauth2.Token{RefreshToken: tok.RefreshToken})
	t, err := src.Token()
	if err != nil {
		return nil, fmt.Errorf("oauth: refresh token: %w", err)
	}
	out := tokenFromOAuth2(t)
	if out.RefreshToken == "" {
		out.RefreshToken = tok.RefreshToken
	}
	return out, nil
}

func (o *OAuth) authorizeFlow(ctx context.Context) (*oauthToken, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", o.callbackPort))
	if err != nil {
		return nil, fmt.Errorf("oauth: failed to bind callback server: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", port)

	cfg := o.oauth2Config(redirectURI)
	state := oauth2.GenerateVerifier()
	authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)

	if o.openURL != nil {
		o.openURL(authURL)
	}

	type result struct {
		code  string
		state string
		err   error
	}
	ch := make(chan result, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if oauthErr := q.Get("error"); oauthErr != "" {
			writeHTML(w, http.StatusBadRequest,
				fmt.Sprintf("<h1>Authorization Failed</h1><p>Error: %s</p>", oauthErr))
			ch <- result{err: fmt.Errorf("oauth: authorization failed: %s", oauthErr)}
			return
		}
		code, st := q.Get("code"), q.Get("state")
		if code == "" || st == "" {
			writeHTML(w, http.StatusBadRequest,
				"<h1>Missing Parameters</h1><p>Authorization code or state not received</p>")
			ch <- result{err: fmt.Errorf("oauth: missing code or state in callback")}
			return
		}
		writeHTML(w, http.StatusOK,
			"<h1>✓ Authorization Successful!</h1><p>You can close this window and return to the terminal.</p>")
		ch <- result{code: code, state: st}
	})

	srv := &http.Server{Handler: mux}
	go srv.Serve(listener) //nolint:errcheck

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

func (o *OAuth) oauth2Config(redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID: o.clientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:   o.baseURL + "/oauth2/authorize",
			TokenURL:  o.baseURL + "/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: redirectURI,
	}
}

func tokenFromOAuth2(t *oauth2.Token) *oauthToken {
	expiresAt := t.Expiry.Unix()
	if t.Expiry.IsZero() {
		expiresAt = time.Now().Unix() + 3600
	}
	rt, _ := t.Extra("refresh_token").(string)
	if rt == "" {
		rt = t.RefreshToken
	}
	return &oauthToken{
		AccessToken:  t.AccessToken,
		RefreshToken: rt,
		ExpiresAt:    expiresAt,
	}
}

func loadTokenFromPath(path string) (*oauthToken, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t oauthToken
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func saveTokenToPath(path string, t *oauthToken) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

const callbackStyle = `<style>html{font-family:system-ui,-apple-system,BlinkMacSystemFont,sans-serif;` +
	`font-size:16px;color:#e0e0e0;background:#202020;padding:2rem;text-align:center;}</style>`

func writeHTML(w http.ResponseWriter, status int, body string) {
	html := fmt.Sprintf("<html><body>%s%s</body></html>", callbackStyle, body)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, html)
}
