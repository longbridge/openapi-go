package trade

import (
	"time"

	"github.com/longbridge/openapi-go"
	"github.com/shopspring/decimal"
)

type GetHistoryExecutions struct {
	Symbol  string    // optional
	StartAt time.Time // optional
	EndAt   time.Time // optional
}

type GetTodayExecutions struct {
	Symbol  string // optional
	OrderId string // optional
}

type GetHistoryOrders struct {
	Symbol  string         // optional
	Status  []OrderStatus  // optional
	Side    OrderSide      // optional
	Market  openapi.Market // optional
	StartAt int64          // optional
	EndAt   int64          // optional
}

type GetTodayOrders struct {
	Symbol string         // optional
	Status []OrderStatus  // optional
	Side   OrderSide      // optional
	Market openapi.Market // optional
}

type GetFundPositions struct {
	Symbols []string // optional
}

type GetStockPositions struct {
	Symbols []string // optional
}

type GetCashFlow struct {
	StartAt      int64 // start timestamp , required
	EndAt        int64 // end timestamp, required
	BusinessType BalanceType
	Symbol       string
	Page         int64
	Size         int64
}

type ReplaceOrder struct {
	OrderId         string          // required
	Quantity        uint64          // required
	Price           decimal.Decimal // LO / ELO / ALO / ODD / LIT Order Required
	TriggerPrice    decimal.Decimal // LIT / MIT Order Required
	LimitOffset     decimal.Decimal // TSLPAMT / TSLPPCT Order Required
	TrailingAmount  decimal.Decimal // TSLPAMT / TSMAMT Order Required
	TrailingPercent decimal.Decimal // TSLPPCT / TSMAPCT Order Required
	Remark          string
}

type SubmitOrder struct {
	Symbol            string          // required
	OrderType         OrderType       // required
	Side              OrderSide       // required
	SubmittedQuantity uint64          // required
	SubmittedPrice    decimal.Decimal // LO / ELO / ALO / ODD / LIT Order Required
	TriggerPrice      decimal.Decimal // LIT / MIT Order Required
	LimitOffset       decimal.Decimal // TSLPAMT / TSLPPCT Order Required
	TrailingAmount    decimal.Decimal // TSLPAMT / TSMAMT Order Required
	TrailingPercent   decimal.Decimal // TSLPPCT / TSMAPCT Order Required
	ExpireDate        *time.Time      // required when time_in_force is GTD
	OutsideRTH        OutsideRTH
	Remark            string
	TimeInForce       TimeType // required
}

type GetEstimateMaxPurchaseQuantity struct {
	Symbol    string
	OrderType OrderType
	Price     decimal.Decimal
	Currency  string
	OrderId   string
	Side      OrderSide
}

type GetAccountBalance struct {
	Currency Currency // optional
}
