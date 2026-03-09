package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/longbridge/openapi-go/config"
	lbhttp "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/oauth"
)

func main() {
	o := oauth.New("your-client-id").
		OnOpenURL(func(url string) { fmt.Println("Open this URL to authorize:", url) })
	if err := o.Build(context.Background()); err != nil {
		log.Fatal(err)
	}
	conf, err := config.New(config.WithOAuthClient(o))
	if err != nil {
		log.Fatal(err)
	}
	client, err := lbhttp.NewFromCfg(conf)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	var data json.RawMessage
	if err := client.Get(ctx, "/v1/trade/execution/today", nil, &data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
