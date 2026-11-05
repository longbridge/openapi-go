// Package oauth provides OAuth 2.0 authentication support for Longbridge OpenAPI.
//
// Use OAuth by setting it on config (like the Rust SDK). The token is stored
// internally and in ~/.longbridge/openapi/tokens/<client_id>; it is refreshed
// automatically when expired. Do not use or expose OAuthToken. Revoke revokes
// the token on the server and Clear drops it locally (memory and file).
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
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/oauth2"
)

const (
	authTimeout         = 5 * time.Minute
	oauthBaseURL        = "https://openapi.longbridge.com"
	oauthBaseURLStaging = "https://openapi-global.longbridge.xyz"
	defaultCallbackPort = 60355
	expiresSoonSecs     = 3600
	tokenDir            = ".longbridge/openapi/tokens"
	tokenRequestTimeout = 30 * time.Second
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

	// refreshSem (capacity 1) serializes token refreshes with each other and
	// with Revoke, so Revoke never reads a refresh token that a refresh is
	// about to rotate. It is a channel rather than a mutex so that waiting
	// honors ctx. Lock order: refreshSem before mu.
	refreshSem     chan struct{}
	refreshSemOnce sync.Once
	// refreshWaiters counts AccessToken callers on the refresh path (tests).
	refreshWaiters atomic.Int32

	mu    sync.Mutex
	token *oauthToken
	// gen is bumped by Clear so a refresh that was in flight when the token
	// was cleared does not store (or persist) it again.
	gen uint64
}

// New creates a new OAuth client with the given client ID.
func New(clientID string) *OAuth {
	base := oauthBaseURL
	if os.Getenv("LONGBRIDGE_ENV") == "staging" {
		base = oauthBaseURLStaging
	}
	return &OAuth{
		clientID:     clientID,
		callbackPort: defaultCallbackPort,
		baseURL:      base,
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

// acquireRefresh takes refreshSem, or returns ctx.Err() if ctx is done first.
func (o *OAuth) acquireRefresh(ctx context.Context) error {
	o.refreshSemOnce.Do(func() { o.refreshSem = make(chan struct{}, 1) })
	select {
	case o.refreshSem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (o *OAuth) releaseRefresh() {
	<-o.refreshSem
}

// Build loads a token from ~/.longbridge/openapi/tokens/<client_id> or runs the
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
		// A refresh of the previous in-memory token that is still in flight
		// must not overwrite the token just loaded.
		o.gen++
		return nil
	}
	// The lock is released during network I/O; gen detects a Clear (or
	// another Build) in the meantime, so a cleared token is not stored again.
	gen := o.gen
	if loaded != nil && loaded.RefreshToken != "" {
		o.mu.Unlock()
		if err := o.acquireRefresh(ctx); err != nil {
			o.mu.Lock()
			return err
		}
		o.mu.Lock()
		// A refresh that finished while we waited has already rotated the
		// refresh token loaded above: use its result instead of refreshing
		// with the stale one.
		if o.gen != gen {
			o.releaseRefresh()
			return o.replacementOrCleared()
		}
		o.mu.Unlock()
		refreshed, err := o.refreshToken(ctx, loaded)
		o.mu.Lock()
		o.releaseRefresh()
		if o.gen != gen {
			return o.replacementOrCleared()
		}
		if err == nil {
			o.storeToken(path, refreshed)
			return nil
		}
	}

	// No valid token: run authorization flow.
	o.mu.Unlock()
	tok, err := o.authorizeFlow(ctx)
	o.mu.Lock()
	if err != nil {
		return err
	}
	if o.gen != gen {
		return o.replacementOrCleared()
	}
	o.storeToken(path, tok)
	return nil
}

// replacementOrCleared decides Build's result when the token changed (gen
// bumped) while Build released mu: the token was cleared, or replaced by a
// concurrent refresh / Build. It keeps a valid replacement and never stores
// the token Build obtained, so a cleared token is not resurrected. The caller
// must hold o.mu.
func (o *OAuth) replacementOrCleared() error {
	if o.token != nil && !o.token.isExpired() {
		return nil
	}
	return ErrTokenCleared
}

var (
	// ErrNoToken is returned when there is no usable token: Build has not
	// been called, the token was cleared, or it has expired (or expires
	// within an hour) and there is no refresh token.
	ErrNoToken = errors.New("oauth: no valid token; call Build(ctx) first or re-run authorization")
	// ErrTokenCleared is returned when Clear ran while a refresh or
	// authorization was in flight; its result was discarded.
	ErrTokenCleared = errors.New("oauth: token was cleared while it was being obtained")
)

// RevokeError is returned by Revoke when the server rejects the revocation.
type RevokeError struct {
	// StatusCode is the HTTP status of the response.
	StatusCode int
	// Code is the OAuth error code (e.g. "unsupported_token_type"), if any.
	Code string
	// Description is the OAuth error_description, if any.
	Description string
}

func (e *RevokeError) Error() string {
	switch {
	case e.Code != "" && e.Description != "":
		return "oauth: revoke token: " + e.Code + ": " + e.Description
	case e.Code != "":
		return "oauth: revoke token: " + e.Code
	case e.Description != "":
		return fmt.Sprintf("oauth: revoke token: server returned HTTP %d: %s", e.StatusCode, e.Description)
	default:
		return fmt.Sprintf("oauth: revoke token: server returned HTTP %d", e.StatusCode)
	}
}

// withTokenHTTPClient returns ctx carrying the *http.Client used for token
// requests: the one already stored under oauth2.HTTPClient (unless nil), or
// one with a 30-second timeout. Refreshes and revocations block other callers (see
// refreshSem), so they must not hang forever.
func withTokenHTTPClient(ctx context.Context) context.Context {
	if c, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok && c != nil {
		return ctx
	}
	// Keep http.DefaultClient's transport (what oauth2 used before), so a
	// proxy or mTLS transport configured there still applies.
	return context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
		Transport: http.DefaultClient.Transport,
		Timeout:   tokenRequestTimeout,
	})
}

// storeToken replaces the token in memory and on disk. It bumps gen so that a
// refresh of the previous token that is still in flight is discarded. The
// caller must hold o.mu.
func (o *OAuth) storeToken(path string, tok *oauthToken) {
	o.token = tok
	o.gen++
	_ = saveTokenToPath(path, tok)
}

// AccessToken returns a valid access token, refreshing it if expired or soon to expire.
// The HTTP client calls this when OAuth is set on config.
func (o *OAuth) AccessToken(ctx context.Context) (string, error) {
	o.mu.Lock()
	if tok := o.token; tok != nil && !tok.isExpired() && !tok.expiresSoon() {
		accessToken := tok.AccessToken
		o.mu.Unlock()
		return accessToken, nil
	}
	o.mu.Unlock()

	// Only one refresh at a time; it also waits for a Revoke in flight.
	o.refreshWaiters.Add(1)
	defer o.refreshWaiters.Add(-1)
	if err := o.acquireRefresh(ctx); err != nil {
		return "", err
	}
	defer o.releaseRefresh()

	o.mu.Lock()
	tok := o.token
	needRefresh := tok == nil || tok.isExpired() || tok.expiresSoon()
	if needRefresh && tok != nil && tok.RefreshToken != "" {
		// Copy so we can refresh without holding the lock (network I/O).
		tokCopy := *tok
		gen := o.gen
		o.mu.Unlock()
		refreshed, err := o.refreshToken(ctx, &tokCopy)
		o.mu.Lock()
		// Check gen before err, like Build: a Clear or replacement during the
		// refresh decides the result even if the refresh itself failed.
		if o.gen != gen {
			// The token was cleared or replaced while refreshing: never
			// store this result; use the replacement if there is one.
			err := o.replacementOrCleared()
			var accessToken string
			if err == nil {
				accessToken = o.token.AccessToken
			}
			o.mu.Unlock()
			return accessToken, err
		}
		if err != nil {
			o.mu.Unlock()
			return "", err
		}
		path, _ := o.tokenPath()
		o.storeToken(path, refreshed)
		o.mu.Unlock()
		return refreshed.AccessToken, nil
	}
	if needRefresh {
		o.mu.Unlock()
		return "", ErrNoToken
	}
	accessToken := tok.AccessToken
	o.mu.Unlock()
	return accessToken, nil
}

// Revoke revokes the current token on the server (RFC 7009,
// POST /oauth2/revoke): the refresh token when present, otherwise the access
// token. It only contacts the server; the in-memory and cached tokens are left
// untouched, so call Clear afterwards to drop them locally.
//
// The request uses the *http.Client stored in ctx under oauth2.HTTPClient (as
// token refreshes do), or a client with a 30-second timeout otherwise (token
// refreshes use the same default). A rejection is returned as *RevokeError;
// with no token it returns an error wrapping ErrNoToken. It is
// serialized with token refreshes (in AccessToken and Build), so it never
// revokes a refresh token that a concurrent refresh is rotating; while it
// runs, AccessToken calls that need a refresh wait for it. It is not
// serialized with a browser authorization still running in Build, which
// stores a new token afterwards.
func (o *OAuth) Revoke(ctx context.Context) error {
	if err := o.acquireRefresh(ctx); err != nil {
		return err
	}
	defer o.releaseRefresh()

	o.mu.Lock()
	tok := o.token
	o.mu.Unlock()
	if tok == nil {
		return fmt.Errorf("%w (nothing to revoke)", ErrNoToken)
	}

	form := url.Values{"client_id": {o.clientID}}
	if tok.RefreshToken != "" {
		form.Set("token", tok.RefreshToken)
		form.Set("token_type_hint", "refresh_token")
	} else {
		form.Set("token", tok.AccessToken)
		form.Set("token_type_hint", "access_token")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		o.baseURL+"/oauth2/revoke", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("oauth: revoke token: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := withTokenHTTPClient(ctx).Value(oauth2.HTTPClient).(*http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("oauth: revoke token: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newRevokeError(resp.StatusCode, body)
	}
	return nil
}

// newRevokeError builds the error for a failed revocation response, keeping
// the server's error and error_description when the body carries them.
func newRevokeError(status int, body []byte) *RevokeError {
	var e struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	_ = json.Unmarshal(body, &e)
	return &RevokeError{StatusCode: status, Code: e.Error, Description: e.ErrorDescription}
}

// Clear drops the token locally: it deletes the cached token file
// (~/.longbridge/openapi/tokens/<client_id>; a missing file is not an error)
// and then clears the in-memory token. It does not contact the server, so call
// Revoke first to disconnect the application completely. After Clear,
// AccessToken fails until Build is called again.
//
// The file is deleted first: if that fails the error is returned and the
// in-memory token is kept, so memory and disk never disagree about whether
// the client is logged in.
func (o *OAuth) Clear() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	path, err := o.tokenPath()
	if err != nil {
		return err
	}
	if err := deleteTokenAtPath(path); err != nil {
		return err
	}
	o.token = nil
	o.gen++
	return nil
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
	src := cfg.TokenSource(withTokenHTTPClient(ctx), &oauth2.Token{RefreshToken: tok.RefreshToken})
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
	authURL, state, verifier := buildAuthorizeURL(cfg)

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

	t, err := cfg.Exchange(withTokenHTTPClient(ctx), res.code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, fmt.Errorf("oauth: failed to exchange code for token: %w", err)
	}
	return tokenFromOAuth2(t), nil
}

// buildAuthorizeURL returns the authorization URL with a fresh CSRF state and
// PKCE challenge (RFC 7636, S256). The server requires PKCE for public clients
// and ignores it for clients that don't need it; the returned verifier must be
// sent with the token exchange.
func buildAuthorizeURL(cfg *oauth2.Config) (authURL, state, verifier string) {
	state = oauth2.GenerateVerifier()
	verifier = oauth2.GenerateVerifier()
	authURL = cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	return authURL, state, verifier
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

// deleteTokenAtPath removes the token file. A missing file is not an error.
func deleteTokenAtPath(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("oauth: failed to delete token file %s: %w", path, err)
	}
	return nil
}

const callbackStyle = `<style>html{font-family:system-ui,-apple-system,BlinkMacSystemFont,sans-serif;` +
	`font-size:16px;color:#e0e0e0;background:#202020;padding:2rem;text-align:center;}</style>`

func writeHTML(w http.ResponseWriter, status int, body string) {
	html := fmt.Sprintf("<html><body>%s%s</body></html>", callbackStyle, body)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, html)
}
