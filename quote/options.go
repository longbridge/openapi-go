package quote

import (
	"github.com/longbridge/openapi-go"
	"github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/log"
	"github.com/longbridge/openapi-go/longbridge"
	protocol "github.com/longbridge/openapi-protocol/go"
)

const (
	DefaultQuoteUrl = "wss://openapi-quote.longbridge.com/v2"
)

// Options for quote context
type Options struct {
	quoteURL           string
	httpClient         *http.Client
	lbOpts             *longbridge.Options
	logLevel           string
	logger             log.Logger
	enableOvernight    bool
	language           openapi.Language
	reconnectCallbacks []func(resubFlag bool)
}

// Option for quote context
type Option func(*Options)

// WithQuoteURL to set url for quote context
func WithQuoteURL(url string) Option {
	return func(o *Options) {
		if url != "" {
			o.quoteURL = url
		}
	}
}

// WithHttpClient to set http client for quote context
func WithHttpClient(client *http.Client) Option {
	return func(o *Options) {
		if client != nil {
			o.httpClient = client
		}
	}
}

// WithLbOptions to set longbridge options for quote context
func WithLbOptions(opts *longbridge.Options) Option {
	return func(o *Options) {
		if opts != nil {
			o.lbOpts = opts
		}
	}
}

// WithLogLevel to set log level for quote context
func WithLogLevel(level string) Option {
	return func(o *Options) {
		if level != "" {
			o.logLevel = level
		}
	}
}

// WithLogger to set logger for quote context
func WithLogger(l log.Logger) Option {
	return func(o *Options) {
		if l != nil {
			o.logger = l
		}
	}
}

// WithEnableOvernight to set enable overnight for quote context
func WithEnableOvernight(enable bool) Option {
	return func(o *Options) {
		o.enableOvernight = enable
	}
}

// WithLanguage to set language for quote context
func WithLanguage(language openapi.Language) Option {
	return func(o *Options) {
		if language != "" {
			o.language = language
		}
	}
}

// OnReconnect to set reconnect callbacks for quote context
func OnReconnect(fn func(successResub bool)) Option {
	return func(o *Options) {
		o.reconnectCallbacks = append(o.reconnectCallbacks, fn)
	}
}

func newOptions(opt ...Option) *Options {
	opts := Options{
		quoteURL: DefaultQuoteUrl,
		lbOpts:   longbridge.NewOptions(),
		logger:   &protocol.DefaultLogger{},
		language: openapi.LanguageEN,
	}
	for _, o := range opt {
		o(&opts)
	}
	return &opts
}
