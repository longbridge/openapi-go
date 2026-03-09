package config_test

import (
	"strings"
	"testing"
	"time"

	"github.com/longbridgeapp/assert"
	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/oauth"
)

var expectedConfig = &config.Config{
	HttpURL:         "http://test",
	HTTPTimeout:     12 * time.Second,
	AppKey:          "test_app_key",
	AppSecret:       "test_app_secret",
	AccessToken:     "test_access_token",
	TradeUrl:        "http://trade_test",
	QuoteUrl:        "http://quote_test",
	EnableOvernight: true,
	AuthTimeout:     12 * time.Second,
	Timeout:         12 * time.Second,
	WriteQueueSize:  12,
	ReadQueueSize:   12,
	ReadBufferSize:  12,
	MinGzipSize:     12,
	Region:          "hk",
}

func Test_withConfigKey(t *testing.T) {
	var c, err = config.New(config.WithConfigKey("appKey", "appSecret", "accessToken"))
	assert.NoError(t, err)
	assert.Equal(t, "appKey", c.AppKey)
	assert.Equal(t, "appSecret", c.AppSecret)
	assert.Equal(t, "accessToken", c.AccessToken)
}

func Test_YamlConfig(t *testing.T) {
	c, err := config.New(config.WithFilePath("./testdata/test_config.yaml"))
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, c)
}

func Test_TomlConfig(t *testing.T) {
	c, err := config.New(config.WithFilePath("./testdata/test_config.toml"))
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, c)
}

func Test_FromOAuth(t *testing.T) {
	c := config.FromOAuth("my-client-id", "my-access-token")
	assert.Equal(t, "my-client-id", c.AppKey)
	assert.Equal(t, "Bearer my-access-token", c.AccessToken)
	assert.Equal(t, "", c.AppSecret)
	assert.True(t, strings.HasPrefix(c.AccessToken, "Bearer "))
}

func Test_FromOAuth_WithBearerPrefix(t *testing.T) {
	c := config.FromOAuth("my-client-id", "Bearer already-prefixed")
	assert.Equal(t, "Bearer already-prefixed", c.AccessToken)
	assert.True(t, strings.HasPrefix(c.AccessToken, "Bearer "))
}

func Test_LegacyMode_NoBearerToken(t *testing.T) {
	c, err := config.New(config.WithConfigKey("appKey", "appSecret", "accessToken"))
	assert.NoError(t, err)
	assert.True(t, c.OAuthClient == nil && !strings.HasPrefix(c.AccessToken, "Bearer "))
}

func Test_WithOAuth(t *testing.T) {
	c, err := config.New(config.WithOAuth("my-client-id", "my-access-token"))
	assert.NoError(t, err)
	assert.Equal(t, "my-client-id", c.AppKey)
	assert.Equal(t, "Bearer my-access-token", c.AccessToken)
	assert.Equal(t, "", c.AppSecret)
	assert.True(t, strings.HasPrefix(c.AccessToken, "Bearer "))
}

func Test_WithOAuthClient(t *testing.T) {
	o := oauth.New("my-client-id")
	c, err := config.New(config.WithOAuthClient(o))
	assert.NoError(t, err)
	assert.True(t, c.OAuthClient != nil)
}

