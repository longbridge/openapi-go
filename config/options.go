package config

import (
	"path"

	"github.com/longbridge/openapi-go/oauth"
)

type Options struct {
	tp       ConfigType
	filePath string

	appKey      *string
	appSecret   *string
	accessToken *string
	oauthClient *oauth.OAuth
}

type Option func(*Options)

// WithFilePath config path
func WithFilePath(filePath string) Option {
	return func(o *Options) {
		if filePath != "" {
			o.filePath = filePath
			fileSuffix := path.Ext(filePath)
			if fileSuffix != "" {
				o.tp = ConfigType(fileSuffix)
			}
		}
	}
}

// WithConfigKey config appKey, appSecret, accessToken
func WithConfigKey(appKey string, appSecret string, accessToken string) Option {
	return func(o *Options) {
		o.appKey = &appKey
		o.appSecret = &appSecret
		o.accessToken = &accessToken
	}
}

// WithOAuthClient configures the client to use OAuth 2.0 with the given OAuth
// client. The token is refreshed automatically. Call oauth.Build(ctx) before
// creating config. Usage (like Rust SDK):
//
//	o := oauth.New("client-id").OnOpenURL(func(url string) { ... })
//	if err := o.Build(ctx); err != nil { ... }
//	cfg, _ := config.New(config.WithOAuthClient(o))
func WithOAuthClient(o *oauth.OAuth) Option {
	return func(opts *Options) {
		opts.oauthClient = o
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		tp: ConfigTypeEnv,
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
