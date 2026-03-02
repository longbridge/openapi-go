package config

import (
	"path"
	"strings"
)

type Options struct {
	tp       ConfigType
	filePath string

	appKey      *string
	appSecret   *string
	accessToken *string
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

// WithOAuth configures the client to use OAuth 2.0 authentication.
//
// This sets the AppKey to clientID and AccessToken to "Bearer <accessToken>",
// and clears AppSecret (not required for OAuth 2.0).
func WithOAuth(clientID, accessToken string) Option {
	return func(o *Options) {
		o.appKey = &clientID
		secret := ""
		o.appSecret = &secret
		if !strings.HasPrefix(accessToken, "Bearer ") {
			accessToken = "Bearer " + accessToken
		}
		o.accessToken = &accessToken
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
