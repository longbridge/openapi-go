package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	quoteContext.OnQuote(func(pe *quote.PushQuote) {
		bytes, _ := json.Marshal(pe)
		fmt.Println(string(bytes))
	})
	quoteContext.OnDepth(func(d *quote.PushDepth) {
		bytes, _ := json.Marshal(d)
		if d.Sequence != 0 {
			fmt.Print(time.UnixMicro(d.Sequence/1000).Format(time.RFC3339) + " ")
		}
		fmt.Println(string(bytes))
	})

	err = quoteContext.Subscribe(ctx, []string{"700.HK"}, []quote.SubType{quote.SubTypeDepth}, true)
	if err != nil {
		log.Fatal(err)
		return
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
