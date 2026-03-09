package oauth_test

import (
	"context"
	"testing"

	"github.com/longbridgeapp/assert"
	"github.com/longbridge/openapi-go/oauth"
)

func TestOAuth_New(t *testing.T) {
	o := oauth.New("test-client-id")
	assert.Equal(t, "test-client-id", o.ClientID())
}

func TestOAuth_WithCallbackPort(t *testing.T) {
	o := oauth.New("test-client-id").WithCallbackPort(8080)
	assert.Equal(t, "test-client-id", o.ClientID())
}

func TestOAuth_AccessToken_ErrorsWithoutBuild(t *testing.T) {
	o := oauth.New("test-client-id")
	_, err := o.AccessToken(context.Background())
	assert.True(t, err != nil)
	assert.Contains(t, err.Error(), "no valid token")
}
