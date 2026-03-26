package main

import (
	"context"
	"fmt"
	"log"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/oauth"
	"github.com/longbridge/openapi-go/trade"
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
	tradeContext, err := trade.NewFromCfg(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer tradeContext.Close()
	ctx := context.Background()
	orders, err := tradeContext.TodayOrders(ctx, &trade.GetTodayOrders{})
	if err != nil {
		log.Fatal(err)
	}
	for _, o := range orders {
		fmt.Printf("%+v\n", o)
	}
}
