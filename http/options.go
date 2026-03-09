package http

import (
	"net/http"
	"time"

	"github.com/longbridge/openapi-go"
	"github.com/longbridge/openapi-go/oauth"
)

// DefaultHttpUrl
const DefaultHttpUrl = "https://openapi.longportapp.com"

// DefaultTimeout
const DefaultTimeout = 15 * time.Second

// Options for http client
type Options struct {
	URL         string
	AppKey      string
	AppSecret   string
	AccessToken string
	Timeout     time.Duration
	Client      *http.Client
	Language    openapi.Language

	// OAuthClient when set: token and app key are taken from it (auto-refresh).
	OAuthClient *oauth.OAuth
}

// Option for http client
type Option func(*Options)

// WithClient use custom *http.Client
func WithClient(cli *http.Client) Option {
	return func(opts *Options) {
		if cli != nil {
			opts.Client = cli
		}
	}
}

// WithURL to set url
func WithURL(url string) Option {
	return func(opts *Options) {
		if url != "" {
			opts.URL = url
		}
	}
}

// WithAppKey to set app key
func WithAppKey(appKey string) Option {
	return func(opts *Options) {
		if appKey != "" {
			opts.AppKey = appKey
		}
	}
}

// WithAppSecret to set app secret
func WithAppSecret(appSecret string) Option {
	return func(opts *Options) {
		if appSecret != "" {
			opts.AppSecret = appSecret
		}
	}
}

// WithAccessToken to set access token
func WithAccessToken(accessToken string) Option {
	return func(opts *Options) {
		if accessToken != "" {
			opts.AccessToken = accessToken
		}
	}
}

// WithTimeout to set http client timeout. Worked when Options.Client is not set
func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		if timeout > 0 {
			opts.Timeout = timeout
		}
	}
}

// WithLanguage to set language
func WithLanguage(language openapi.Language) Option {
	return func(opts *Options) {
		if language != "" {
			opts.Language = language
		}
	}
}

// WithOAuthClient sets the OAuth client for token auto-refresh.
func WithOAuthClient(o *oauth.OAuth) Option {
	return func(opts *Options) {
		opts.OAuthClient = o
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		Timeout: DefaultTimeout,
		URL:     DefaultHttpUrl,
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
