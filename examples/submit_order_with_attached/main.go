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
	"github.com/longbridge/openapi-go/trade"
	"github.com/shopspring/decimal"
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
		return
	}
	defer tradeContext.Close()
	ctx := context.Background()

	tradeContext.OnTrade(func(ev *trade.PushEvent) {
		log.Printf("order event: %+v\n", ev)
	})

	// Submit a bracket order: buy 700.HK at 12.00, with a take-profit at 13.00
	// and a stop-loss at 11.00 attached to it.
	order := &trade.SubmitOrder{
		Symbol:            "700.HK",
		OrderType:         trade.OrderTypeLO,
		Side:              trade.OrderSideBuy,
		SubmittedQuantity: 200,
		TimeInForce:       trade.TimeTypeDay,
		SubmittedPrice:    decimal.NewFromFloat(12),
		AttachedParams: &trade.SubmitAttachedParams{
			AttachedOrderType: trade.AttachedOrderTypeBracket,
			ProfitTakerPrice:  decimal.NewFromFloat(13),
			StopLossPrice:     decimal.NewFromFloat(11),
		},
	}
	orderId, err := tradeContext.SubmitOrder(ctx, order)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("orderId: %v\n", orderId)

	detail, err := tradeContext.OrderDetail(ctx, orderId)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, attached := range detail.AttachedOrders {
		fmt.Printf("attached order: %+v\n", attached)

		// The attached sub-order has its own order ID; query/cancel it directly
		// via the *Attached variants instead of the parent order's ID.
		attachedDetail, err := tradeContext.OrderDetailAttached(ctx, attached.OrderId)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("attached order detail: %+v\n", attachedDetail)
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
