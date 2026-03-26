package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/longbridge/openapi-go/config"
	"github.com/longbridge/openapi-go/oauth"
	"github.com/longbridge/openapi-go/quote"
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
	quoteContext, err := quote.NewFromCfg(conf)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer quoteContext.Close()
	ctx := context.Background()
	quotes, err := quoteContext.Quote(ctx, []string{"700.HK", "AAPL.US", "TSLA.US", "NFLX.US"})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("quotes: %+v\n", quotes[0])

	warrants, err := quoteContext.WarrantList(ctx, "700.HK", quote.WarrantFilter{
		SortBy:     quote.WarrantVolume,
		SortOrder:  quote.WarrantAsc,
		SortOffset: 0,
		SortCount:  10,
	}, quote.WarrantEN)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("warrants: %+v\n", warrants[0])

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
