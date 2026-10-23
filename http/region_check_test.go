package http

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/longbridge/openapi-go/oauth"
)

// newOAuthTestClient returns a Client backed by an OAuth session whose stored
// access token is accessToken. The token is written to a temporary HOME so
// oauth.Build loads it without running the authorization flow.
func newOAuthTestClient(t *testing.T, accessToken string) *Client {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	const clientID = "test-client"
	dir := filepath.Join(home, ".longbridge", "openapi", "tokens")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("mkdir token dir: %v", err)
	}
	data, err := json.Marshal(map[string]any{
		"access_token": accessToken,
		"expires_at":   time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		t.Fatalf("marshal token: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, clientID), data, 0o600); err != nil {
		t.Fatalf("write token: %v", err)
	}

	o := oauth.New(clientID)
	if err := o.Build(context.Background()); err != nil {
		t.Fatalf("oauth Build: %v", err)
	}
	c, err := New(WithURL("http://127.0.0.1:1"), WithOAuthClient(o))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestCheckRegionCtxOAuthUSToken(t *testing.T) {
	ctx := context.Background()
	c := newOAuthTestClient(t, "us_oauth_token")

	if err := c.CheckRegionCtx(ctx, "/v1/us/assets/overview", "US"); err != nil {
		t.Errorf("US check with us_ OAuth token: got %v, want nil", err)
	}
	var rre *RegionRestrictedError
	if err := c.CheckRegionCtx(ctx, "/v1/asset/account", "AP"); !errors.As(err, &rre) {
		t.Errorf("AP check with us_ OAuth token: got %v, want RegionRestrictedError", err)
	}
	isUS, err := c.IsUSCtx(ctx)
	if err != nil || !isUS {
		t.Errorf("IsUSCtx with us_ OAuth token = %v, %v; want true, nil", isUS, err)
	}
}

func TestCheckRegionCtxOAuthAPToken(t *testing.T) {
	ctx := context.Background()
	c := newOAuthTestClient(t, "plain_oauth_token")

	if err := c.CheckRegionCtx(ctx, "/v1/asset/account", "AP"); err != nil {
		t.Errorf("AP check with unprefixed OAuth token: got %v, want nil", err)
	}
	var rre *RegionRestrictedError
	if err := c.CheckRegionCtx(ctx, "/v1/us/assets/overview", "US"); !errors.As(err, &rre) {
		t.Errorf("US check with unprefixed OAuth token: got %v, want RegionRestrictedError", err)
	}
	isUS, err := c.IsUSCtx(ctx)
	if err != nil || isUS {
		t.Errorf("IsUSCtx with unprefixed OAuth token = %v, %v; want false, nil", isUS, err)
	}
}

func TestCheckRegionCtxOAuthTokenError(t *testing.T) {
	// An OAuth client that was never built has no token: the check must
	// surface the token error rather than guess a region.
	c, err := New(WithURL("http://127.0.0.1:1"), WithOAuthClient(oauth.New("unbuilt")))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	err = c.CheckRegionCtx(context.Background(), "/v1/asset/account", "AP")
	var rre *RegionRestrictedError
	if err == nil || errors.As(err, &rre) {
		t.Errorf("CheckRegionCtx without a token: got %v, want token error", err)
	}
}

func TestCheckRegionCtxAPIKey(t *testing.T) {
	ctx := context.Background()
	us, err := New(WithURL("http://127.0.0.1:1"), WithAppKey("us_key"), WithAppSecret("us_secret"), WithAccessToken("us_token"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := us.CheckRegionCtx(ctx, "/v1/us/assets/overview", "US"); err != nil {
		t.Errorf("US check with us_ API key: got %v, want nil", err)
	}
	var rre *RegionRestrictedError
	if err := us.CheckRegionCtx(ctx, "/v1/asset/account", "AP"); !errors.As(err, &rre) {
		t.Errorf("AP check with us_ API key: got %v, want RegionRestrictedError", err)
	}
}
