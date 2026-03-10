# Longport OpenAPI SDK for Go

`Longport` provides an easy-to-use interface for invokes [`Longport OpenAPI`](https://open.longportapp.com/en/).

## Quickstart

_With Go module support , simply add the following import_

```golang
import "github.com/longbridge/openapi-go"
```

## Authentication

### 1. OAuth 2.0 (Recommended)

OAuth 2.0 is the modern authentication method that uses Bearer tokens without requiring HMAC signatures.

**Step 1: Register OAuth Client**

First, register an OAuth client to get your `client_id`:

```bash
curl -X POST https://openapi.longbridge.com/v1/oauth2/client/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Application",
    "redirect_uris": ["http://localhost:60355/callback"],
    "grant_types": ["authorization_code", "refresh_token"]
  }'
```

Response:
```json
{
  "client_id": "your-client-id-here",
  "name": "My Application",
  "redirect_uris": ["http://localhost:60355/callback"]
}
```

Save the `client_id` for use in your application.

**Step 2: Build OAuth and set on config (same usage as Rust SDK)**

The SDK stores the token in `~/.longbridge-openapi/tokens/<client_id>` and refreshes it automatically. You do not handle tokens yourself.

```golang
import (
    "context"
    "fmt"
    "log"

    "github.com/longbridge/openapi-go/config"
    "github.com/longbridge/openapi-go/oauth"
    "github.com/longbridge/openapi-go/quote"
)

func main() {
    o := oauth.New("your-client-id").
        OnOpenURL(func(url string) {
            fmt.Println("Please visit:", url)
        })
    // Load token from disk or run authorization flow; token is persisted and auto-refreshed
    if err := o.Build(context.Background()); err != nil {
        log.Fatal(err)
    }

    // Set OAuth on config (like Rust SDK: Config::from_oauth(oauth))
    cfg, err := config.New(config.WithOAuthClient(o))
    if err != nil {
        log.Fatal(err)
    }

    quoteContext, err := quote.NewFromCfg(cfg)
    // ...
}
```

**Benefits:**
- Token is stored in `~/.longbridge-openapi/tokens/<client_id>` and refreshed automatically
- No need to handle or expose tokens in your code
- Same usage pattern as the Rust SDK (OAuth set on config)
- More secure (no shared secret), no HMAC signature

### 2. Legacy API Key (Environment Variables)

For backward compatibility, you can still use the traditional API key method:

_Setting environment variables (MacOS/Linux)_

```bash
export LONGPORT_APP_KEY="App Key get from user center"
export LONGPORT_APP_SECRET="App Secret get from user center"
export LONGPORT_ACCESS_TOKEN="Access Token get from user center"
```

_Setting environment variables (Windows)_

```powershell
setx LONGPORT_APP_KEY "App Key get from user center"
setx LONGPORT_APP_SECRET "App Secret get from user center"
setx LONGPORT_ACCESS_TOKEN "Access Token get from user center"
```

## Config

### Load from env

Support init config from env, and support load env from `.env` file

```golang
import (
    "github.com/longbridge/openapi-go/config"
    "github.com/longbridge/openapi-go/trade"
    "github.com/longbridge/openapi-go/http"
)

func main() {
    c, err := config.New()

    if err != nil {
        // panic
    }

    // init http client from config
    c, err := http.NewFromCfg(c)

    // init trade context from config
    tc, err := trade.NewFromCfg(c)

    // init quote context from config
    qc, err := quote.NewFromCfg(c)
}

```

All envs is listed in the last of [README](#environment-variables)

### Load from file[yaml,toml]

#### yaml example

To load configuration from a YAML file, use the following code snippet:

```golang
conf, err := config.New(config.WithFilePath("./test.yaml"))
```

Here is an example of what the `test.yaml` file might look like:


```yaml
longport:
  app_key: xxxxx
  app_secret: xxxxx 
  access_token: xxxxx 
```

#### toml example

Similarly, to load configuration from a TOML file, use this code snippet:

```golang
conf, err := config.New(config.WithFilePath("./test.toml"))
```

And here is an example of a `test.toml` file:

```toml
[longport]
app_key = "xxxxx"
app_secret = "xxxxx"
access_token = "xxxxx"
```

### Init Config manually

Config structure as follow:

```golang
type Config struct {
    HttpURL     string        `env:"LONGPORT_HTTP_URL" yaml:"http_url" toml:"http_url"`
    HTTPTimeout time.Duration `env:"LONGPORT_HTTP_TIMEOUT" yaml:"http_timeout" toml:"http_timeout"`
    AppKey      string        `env:"LONGPORT_APP_KEY" yaml:"app_key" toml:"app_key"`
    AppSecret   string        `env:"LONGPORT_APP_SECRET" yaml:"app_secret" toml:"app_secret"`
    AccessToken string        `env:"LONGPORT_ACCESS_TOKEN" yaml:"access_token" toml:"access_token"`
    TradeUrl    string        `env:"LONGPORT_TRADE_URL" yaml:"trade_url" toml:"trade_url"`
    QuoteUrl    string        `env:"LONGPORT_QUOTE_URL" yaml:"quote_url" toml:"quote_url"`
    EnableOvernight bool          `env:"LONGPORT_ENABLE_OVERNIGHT" yaml:"enable_overnight" toml:"enable_overnight"`
    Language    openapi.Language `env:"LONGPORT_LANGUAGE" yaml:"language" toml:"language"`

    LogLevel string `env:"LONGPORT_LOG_LEVEL" yaml:"log_level" toml:"log_level"`
    // LONGPORT protocol config
    AuthTimeout    time.Duration `env:"LONGPORT_AUTH_TIMEOUT" yaml:"auth_timeout" toml:"timeout"`
    Timeout        time.Duration `env:"LONGPORT_TIMEOUT" yaml:"timeout" toml:"timeout"`
    WriteQueueSize int           `env:"LONGPORT_WRITE_QUEUE_SIZE" yaml:"write_queue_size" toml:"write_queue_size"`
    ReadQueueSize  int           `env:"LONGPORT_READ_QUEUE_SIZE" yaml:"read_queue_size" toml:"read_queue_size"`
    ReadBufferSize int           `env:"LONGPORT_READ_BUFFER_SIZE" yaml:"read_buffer_size" toml:"read_buffer_size"`
    MinGzipSize    int           `env:"LONGPORT_MIN_GZIP_SIZE" yaml:"min_gzip_size" toml:"min_gzip_size"`
    Region Region `env:"LONGPORT_REGION" yaml:"region" toml:"region"`
}

```

set config field manually

```golang
c, err := config.New()
c.AppKey = "xxx"
c.AppSecret = "xxx"
c.AccessToken = "xxx"

```

### set custom logger

Our logger interface as follow:

```golang
type Logger interface {
    SetLevel(string)
    Info(msg string)
    Error(msg string)
    Warn(msg string)
    Debug(msg string)
    Infof(msg string, args ...interface{})
    Errorf(msg string, args ...interface{})
    Warnf(msg string, args ...interface{})
    Debugf(msg string, args ...interface{})
}

```

Your can use you own logger by imply the interface

```golang
c, err := config.New()

l := newOwnLogger()

c.SetLogger(l)

```

### use custom \*(net/http).Client

the default http client is initialized simply as follow:

```golang
cli := &http.Client{Timeout: opts.Timeout}
```

we only set timeout here, you can use you own \*(net/http).Client.

```golang
c, err := config.New()

c.Client = &http.Client{
    Transport: ...
}

```

## Quote API (Get basic information of securities)

**Using OAuth 2.0 (Recommended):**

```golang
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/longbridge/openapi-go/config"
    "github.com/longbridge/openapi-go/oauth"
    "github.com/longbridge/openapi-go/quote"
)

func main() {
    o := oauth.New("your-client-id").
        OnOpenURL(func(url string) {
            fmt.Println("Please visit:", url)
        })
    if err := o.Build(context.Background()); err != nil {
        log.Fatal(err)
    }
    cfg, err := config.New(config.WithOAuthClient(o))
    if err != nil {
        log.Fatal(err)
    }
    quoteContext, err := quote.NewFromCfg(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer quoteContext.Close()
    quotes, err := quoteContext.Quote(context.Background(), []string{"700.HK", "AAPL.US", "TSLA.US", "NFLX.US"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("quotes: %v", quotes)
}
```

**Using legacy API key (environment variables):**

```golang
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/longbridge/openapi-go/quote"
    "github.com/longbridge/openapi-go/config"
)

func main() {
    conf, err := config.New()
    if err != nil {
        log.Fatal(err)
        return
    }
    // create quote context from environment variables
    quoteContext, err := quote.NewFromCfg(conf)
    if err != nil {
        log.Fatal(err)
        return
    }
    defer quoteContext.Close()
    ctx := context.Background()
    // Get basic information of securities
    quotes, err := quoteContext.Quote(ctx, []string{"700.HK", "AAPL.US", "TSLA.US", "NFLX.US"})
    if err != nil {
        log.Fatal(err)
        return
    }
    fmt.Printf("quotes: %v", quotes)
}
```

## Trade API (Submit order)

**Using OAuth 2.0 (Recommended):**

```golang
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/longbridge/openapi-go/config"
    "github.com/longbridge/openapi-go/oauth"
    "github.com/longbridge/openapi-go/trade"
    "github.com/shopspring/decimal"
)

func main() {
    o := oauth.New("your-client-id").
        OnOpenURL(func(url string) {
            fmt.Println("Please visit:", url)
        })
    if err := o.Build(context.Background()); err != nil {
        log.Fatal(err)
    }
    cfg, err := config.New(config.WithOAuthClient(o))
    if err != nil {
        log.Fatal(err)
    }
    tradeContext, err := trade.NewFromCfg(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer tradeContext.Close()
    order := &trade.SubmitOrder{
        Symbol:            "700.HK",
        OrderType:         trade.OrderTypeLO,
        Side:              trade.OrderSideBuy,
        SubmittedQuantity: 200,
        TimeInForce:       trade.TimeTypeDay,
        SubmittedPrice:    decimal.NewFromFloat(12),
    }
    orderId, err := tradeContext.SubmitOrder(context.Background(), order)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("orderId: %v\n", orderId)
}
```

**Using legacy API key (environment variables):**

```golang
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/longbridge/openapi-go/trade"
    "github.com/longbridge/openapi-go/config"
    "github.com/shopspring/decimal"
)

func main() {
    conf, err := config.New()
    if err != nil {
        log.Fatal(err)
        return
    }
    // create trade context from environment variables
    tradeContext, err := trade.NewFromCfg(conf)
    if err != nil {
        log.Fatal(err)
        return
    }
    defer tradeContext.Close()
    ctx := context.Background()
    // submit order
    order := &trade.SubmitOrder{
        Symbol:            "700.HK",
        OrderType:         trade.OrderTypeLO,
        Side:              trade.OrderSideBuy,
        SubmittedQuantity: 200,
        TimeInForce:       trade.TimeTypeDay,
        SubmittedPrice:    decimal.NewFromFloat(12),
    }
    orderId, err := tradeContext.SubmitOrder(ctx, order)
    if err != nil {
        log.Fatal(err)
        return
    }
    fmt.Printf("orderId: %v\n", orderId)
}
```

## Environment Variables

Support load env from `.env` file.

| name                      | description                                                                                           | default value                       | example | optional       |
| ------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------------- | ------- | -------------- |
| LONGPORT_REGION           | Set access region, if region equals `cn`, sdk will set httpUrl, quoteUrl, tradeUrl to China endpoints | -                                   | cn      | cn             |
| LONGPORT_HTTP_URL         | LONGPORT rest api url                                                                                 | <https://openapi.longbridge.com>    |         |                |
| LONGPORT_APP_KEY          | app key                                                                                               |                                     |         |                |
| LONGPORT_APP_SECRET       | app secret                                                                                            |                                     |         |                |
| LONGPORT_ACCESS_TOKEN     | access token                                                                                          |                                     |         |                |
| LONGPORT_TRADE_URL        | LONGPORT protocol url for trade context                                                               | wss://openapi-trade.longbridge.com  |         |                |
| LONGPORT_QUOTE_URL        | LONGPORT protocol url for quote context                                                               | wss://openapi-quote.longbridge.com  |         |                |
| LONGPORT_LOG_LEVEL        | log level                                                                                             | info                                |         |                |
| LONGPORT_AUTH_TIMEOUT     | LONGPORT protocol authorize request time out                                                          | 10 second                           | 10s     |                |
| LONGPORT_TIMEOUT          | LONGPORT protocol dial timeout                                                                        | 5 second                            | 6s      |                |
| LONGPORT_WRITE_QUEUE_SIZE | longport protocol write queue size                                                                    | 16                                  |         |                |
| LONGPORT_READ_QUEUE_SIZE  | longport protocol read queue size                                                                     | 16                                  |         |                |
| LONGPORT_READ_BUFFER_SIZE | longport protocol read buffer size                                                                    | 4096                                |         |                |
| LONGPORT_MIN_GZIP_SIZE    | longport protocol minimal gzip size                                                                   | 1024                                |         |                |
| LONGPORT_ENABLE_OVERNIGHT | enable overnight quote subscription feature                                                           | false                               |         |                |
| LONGPORT_LANGUAGE         | set user language for some information.                                                               | -                                   | en      | en,zh-CN,zh-HK |

## License

Licensed under either of

- Apache License, Version 2.0,([LICENSE-APACHE](./LICENSE-APACHE) or <http://www.apache.org/licenses/LICENSE-2.0>)
- MIT license ([LICENSE-MIT](./LICENSE-MIT) or <http://opensource.org/licenses/MIT>) at your option.
