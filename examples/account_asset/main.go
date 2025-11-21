package main

import (
	"context"
	"fmt"
	"log"

	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/trade"
)

func main() {
	// create trade context from environment variables
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	tradeContext, err := trade.NewFromCfg(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer tradeContext.Close()
	ctx := context.Background()
	// Get AccountBalance infomation
	ab, err := tradeContext.AccountBalance(ctx, &trade.GetAccountBalance{Currency: trade.CurrencyHKD})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", ab[0])
}
