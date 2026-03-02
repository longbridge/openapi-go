package oauth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/longbridgeapp/assert"
	"github.com/longportapp/openapi-go/oauth"
)

func TestOAuthToken_IsExpired(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		tok := &oauth.OAuthToken{
			AccessToken: "test",
			ExpiresAt:   time.Now().Unix() + 7200,
		}
		assert.False(t, tok.IsExpired())
	})

	t.Run("expired", func(t *testing.T) {
		tok := &oauth.OAuthToken{
			AccessToken: "test",
			ExpiresAt:   time.Now().Unix() - 1,
		}
		assert.True(t, tok.IsExpired())
	})
}

func TestOAuthToken_ExpiresSoon(t *testing.T) {
	t.Run("expires soon (30 min)", func(t *testing.T) {
		tok := &oauth.OAuthToken{
			AccessToken: "test",
			ExpiresAt:   time.Now().Unix() + 1800,
		}
		assert.True(t, tok.ExpiresSoon())
	})

	t.Run("not expires soon (2 hours)", func(t *testing.T) {
		tok := &oauth.OAuthToken{
			AccessToken: "test",
			ExpiresAt:   time.Now().Unix() + 7200,
		}
		assert.False(t, tok.ExpiresSoon())
	})
}

func TestOAuthToken_JSONSerialization(t *testing.T) {
	tok := &oauth.OAuthToken{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		ExpiresAt:    1234567890,
	}

	data, err := json.Marshal(tok)
	assert.NoError(t, err)

	var decoded oauth.OAuthToken
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, tok.AccessToken, decoded.AccessToken)
	assert.Equal(t, tok.RefreshToken, decoded.RefreshToken)
	assert.Equal(t, tok.ExpiresAt, decoded.ExpiresAt)
}

func TestOAuthToken_JSONSerialization_NoRefresh(t *testing.T) {
	tok := &oauth.OAuthToken{
		AccessToken: "test_access_token",
		ExpiresAt:   1234567890,
	}

	data, err := json.Marshal(tok)
	assert.NoError(t, err)

	// refresh_token should be omitted from JSON when empty
	assert.False(t, strings.Contains(string(data), "refresh_token"))

	var decoded oauth.OAuthToken
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "", decoded.RefreshToken)
}

func TestOAuth_New(t *testing.T) {
	o := oauth.New("test-client-id")
	assert.Equal(t, "test-client-id", o.ClientID())
}

func TestOAuth_WithCallbackPort(t *testing.T) {
	o := oauth.New("test-client-id").WithCallbackPort(8080)
	assert.Equal(t, "test-client-id", o.ClientID())
}

func TestOAuth_Refresh(t *testing.T) {
	// Start a test HTTP server to mock the token endpoint
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/oauth2/token", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		r.ParseForm()
		assert.Equal(t, "refresh_token", r.FormValue("grant_type"))
		assert.Equal(t, "old-refresh-token", r.FormValue("refresh_token"))
		assert.Equal(t, "test-client-id", r.FormValue("client_id"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  "new-access-token",
			"refresh_token": "new-refresh-token",
			"expires_in":    3600,
			"token_type":    "Bearer",
		})
	}))
	defer srv.Close()

	// Override the base URL via the exported helper for tests
	o := oauth.NewWithBaseURL("test-client-id", srv.URL)
	ctx := t.Context()
	tok, err := o.Refresh(ctx, "old-refresh-token")
	assert.NoError(t, err)
	assert.Equal(t, "new-access-token", tok.AccessToken)
	assert.Equal(t, "new-refresh-token", tok.RefreshToken)
	assert.False(t, tok.IsExpired())
}

func TestOAuth_Refresh_PreservesRefreshToken(t *testing.T) {
	// Server returns no refresh_token — client should keep the old one
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "new-access-token",
			"expires_in":   3600,
		})
	}))
	defer srv.Close()

	o := oauth.NewWithBaseURL("test-client-id", srv.URL)
	tok, err := o.Refresh(t.Context(), "original-refresh-token")
	assert.NoError(t, err)
	assert.Equal(t, "new-access-token", tok.AccessToken)
	assert.Equal(t, "original-refresh-token", tok.RefreshToken)
}
