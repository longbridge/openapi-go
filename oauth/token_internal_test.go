package oauth

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/longbridgeapp/assert"
)

func Test_oauthToken_IsExpired(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		tok := &oauthToken{
			AccessToken: "test",
			ExpiresAt:   time.Now().Unix() + 7200,
		}
		assert.False(t, tok.isExpired())
	})
	t.Run("expired", func(t *testing.T) {
		tok := &oauthToken{
			AccessToken: "test",
			ExpiresAt:   time.Now().Unix() - 1,
		}
		assert.True(t, tok.isExpired())
	})
}

func Test_oauthToken_ExpiresSoon(t *testing.T) {
	t.Run("expires soon (30 min)", func(t *testing.T) {
		tok := &oauthToken{
			AccessToken: "test",
			ExpiresAt:   time.Now().Unix() + 1800,
		}
		assert.True(t, tok.expiresSoon())
	})
	t.Run("not expires soon (2 hours)", func(t *testing.T) {
		tok := &oauthToken{
			AccessToken: "test",
			ExpiresAt:   time.Now().Unix() + 7200,
		}
		assert.False(t, tok.expiresSoon())
	})
}

func Test_oauthToken_JSONRoundtrip(t *testing.T) {
	tok := &oauthToken{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		ExpiresAt:    1234567890,
	}
	data, err := json.Marshal(tok)
	assert.NoError(t, err)
	var decoded oauthToken
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, tok.AccessToken, decoded.AccessToken)
	assert.Equal(t, tok.RefreshToken, decoded.RefreshToken)
	assert.Equal(t, tok.ExpiresAt, decoded.ExpiresAt)
}
