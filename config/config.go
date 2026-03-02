package config

import (
	"net/http"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"

	"github.com/longportapp/openapi-go"
	"github.com/longportapp/openapi-go/log"
)

type IConfig interface {
	GetConfig(opts *Options) (*Config, error)
}

var configTypeMap = map[ConfigType]IConfig{
	ConfigTypeEnv:  &EnvConfig{},
	ConfigTypeYAML: &YAMLConfig{},
	ConfigTypeTOML: &TOMLConfig{},
}

type Region string

var (
	RegionCN Region = "cn"

	cnHttpUrl  = "https://openapi.longportapp.cn"
	cnQuoteUrl = "wss://openapi-quote.longportapp.cn"
	cnTradeUrl = "wss://openapi-trade.longportapp.cn"
)

// Config store Longbridge config
type Config struct {
	// Client custom http client
	Client *http.Client

	HttpURL         string           `env:"LONGBRIDGE_HTTP_URL,LONGPORT_HTTP_URL" yaml:"http_url" toml:"http_url"`
	HTTPTimeout     time.Duration    `env:"LONGBRIDGE_HTTP_TIMEOUT,LONGPORT_HTTP_TIMEOUT" yaml:"http_timeout" toml:"http_timeout"`
	AppKey          string           `env:"LONGBRIDGE_APP_KEY,LONGPORT_APP_KEY" yaml:"app_key" toml:"app_key"`
	AppSecret       string           `env:"LONGBRIDGE_APP_SECRET,LONGPORT_APP_SECRET" yaml:"app_secret" toml:"app_secret"`
	AccessToken     string           `env:"LONGBRIDGE_ACCESS_TOKEN,LONGPORT_ACCESS_TOKEN" yaml:"access_token" toml:"access_token"`
	TradeUrl        string           `env:"LONGBRIDGE_TRADE_URL,LONGPORT_TRADE_URL" yaml:"trade_url" toml:"trade_url"`
	QuoteUrl        string           `env:"LONGBRIDGE_QUOTE_URL,LONGPORT_QUOTE_URL" yaml:"quote_url" toml:"quote_url"`
	EnableOvernight bool             `env:"LONGPORT_ENABLE_OVERNIGHT" yaml:"enable_overnight" toml:"enable_overnight"`
	Language        openapi.Language `env:"LONGPORT_LANGUAGE" yaml:"language" toml:"language"`

	LogLevel string `env:"LONGBRIDGE_LOG_LEVEL,LONGPORT_LOG_LEVEL" yaml:"log_level" toml:"log_level"`
	logger   log.Logger

	// longbridge protocol config
	AuthTimeout    time.Duration `env:"LONGBRIDGE_AUTH_TIMEOUT,LONGPORT_AUTH_TIMEOUT" yaml:"auth_timeout" toml:"auth_timeout"`
	Timeout        time.Duration `env:"LONGBRIDGE_TIMEOUT,LONGPORT_TIMEOUT" yaml:"timeout" toml:"timeout"`
	WriteQueueSize int           `env:"LONGBRIDGE_WRITE_QUEUE_SIZE,LONGPORT_WRITE_QUEUE_SIZE" yaml:"write_queue_size" toml:"write_queue_size"`
	ReadQueueSize  int           `env:"LONGBRIDGE_READ_QUEUE_SIZE,LONGPORT_READ_QUEUE_SIZE" yaml:"read_queue_size" toml:"read_queue_size"`
	ReadBufferSize int           `env:"LONGBRIDGE_READ_BUFFER_SIZE,LONGPORT_READ_BUFFER_SIZE" yaml:"read_buffer_size" toml:"read_buffer_size"`
	MinGzipSize    int           `env:"LONGBRIDGE_MIN_GZIP_SIZE,LONGPORT_MIN_GZIP_SIZE" yaml:"min_gzip_size" toml:"min_gzip_size"`
	Region         Region        `env:"LONGPORT_REGION" yaml:"region" toml:"region"`
}

// parseConfig is a config for toml/yaml
type parseConfig struct {
	Longport *Config `toml:"longport" yaml:"longport"`
}

func (c *Config) SetLogger(l log.Logger) {
	if l != nil {
		l.SetLevel(c.LogLevel)
		c.logger = l
		log.SetLogger(l)
	}
}

func (c *Config) Logger() log.Logger {
	return c.logger
}

// FromOAuth creates a Config using OAuth 2.0 authentication.
//
// OAuth 2.0 mode uses Bearer token authentication and does not require
// AppSecret. The access_token will be automatically prefixed with "Bearer "
// if not already present.
//
// # Arguments
//
//   - clientID: OAuth 2.0 client ID (used as AppKey)
//   - accessToken: OAuth 2.0 access token
func FromOAuth(clientID, accessToken string) *Config {
	// Ensure Bearer prefix
	if !strings.HasPrefix(accessToken, "Bearer ") {
		accessToken = "Bearer " + accessToken
	}
	return &Config{
		AppKey:      clientID,
		AccessToken: accessToken,
	}
}

// IsOAuth2 reports whether this config is using OAuth 2.0 mode.
// OAuth 2.0 mode is detected when the access token starts with "Bearer ".
func (c *Config) IsOAuth2() bool {
	return strings.HasPrefix(c.AccessToken, "Bearer ")
}

func New(opts ...Option) (configData *Config, err error) {
	options := newOptions(opts...)
	conf, exist := configTypeMap[options.tp]
	if !exist {
		err = errors.Errorf("config type:%+v not support", options.tp)
		return
	}
	configData, err = conf.GetConfig(options)
	if err != nil {
		err = errors.Wrapf(err, "GetConfig err")
		return
	}

	if options.appKey != nil {
		configData.AppKey = *options.appKey
	}
	if options.appSecret != nil {
		configData.AppSecret = *options.appSecret
	}
	if options.accessToken != nil {
		configData.AccessToken = *options.accessToken
	}

	if configData.Region == RegionCN {
		configData.HttpURL = cnHttpUrl
		configData.QuoteUrl = cnQuoteUrl
		configData.TradeUrl = cnTradeUrl
	}

	err = configData.check()
	if err != nil {
		err = errors.Wrapf(err, "New config check err")
		return
	}
	log.SetLevel(configData.LogLevel)
	return
}

func (c *Config) check() (err error) {
	if c.AccessToken == "" {
		err = errors.New("Don't has accessToken. Please set access token on LONGPORT_ACCESS_TOKEN env")
		return
	}
	if c.AppKey == "" {
		err = errors.New("Don't has appKey. Please set app key on LONGPORT_APP_KEY env")
		return
	}
	// OAuth 2.0 mode does not require AppSecret
	if !c.IsOAuth2() && c.AppSecret == "" {
		err = errors.New("Don't has appSecret. Please set app secret on LONGPORT_APP_SECRET env")
		return
	}
	return
}

// Deprecated: NewFormEnv to create config with enviromente variables
func NewFormEnv() (*Config, error) {
	return New()
}
