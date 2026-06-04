package main

import (
	"context"
	"fmt"
	"log"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/counter"
	"github.com/longbridge/openapi-go/fundamental"
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
	fctx, err := fundamental.NewFromCfg(conf)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	const symbol = "QQQ.US"
	fmt.Printf("%s -> counter_id %s (is_etf=%v)\n",
		symbol, counter.SymbolToCounterID(symbol), counter.IsETF(symbol))

	resp, err := fctx.EtfAssetAllocation(ctx, symbol)
	if err != nil {
		log.Fatal(err)
	}
	for _, group := range resp.Info {
		fmt.Printf("group: type=%d report_date=%s items=%d\n",
			group.AssetType, group.ReportDate, len(group.Lists))
		for _, item := range group.Lists {
			fmt.Printf("  %s (%s) ratio=%s\n", item.Name, item.Symbol, item.PositionRatio)
		}
	}
}
