package grid

import (
	"github.com/shopspring/decimal"
)

// GridOrders is the response for the grid orders list request.
type GridOrders struct {
	GridOrder []*GridOrder
	HasMore   bool
}

// GridTriggerHistory is the response for the grid trigger history request.
type GridTriggerHistory struct {
	TriggerOrders []*TriggerOrder
	HasMore       bool
}

// GridOrder is a grid trading order (element of the list / by-ids responses).
type GridOrder struct {
	OrderId              string
	Symbol               string
	StockName            string
	Market               string
	Status               string
	GridStatus           string
	SubmittedBasePrice   *decimal.Decimal
	CurrentBasePrice     *decimal.Decimal
	PreTriggerBasePrice  *decimal.Decimal
	PostTriggerBasePrice *decimal.Decimal
	UpperLimitPrice      *decimal.Decimal
	LowerLimitPrice      *decimal.Decimal
	TriggerPriceType     TriggerPriceType // Spread / Percent
	TriggerSpreadUp      *decimal.Decimal
	TriggerSpreadDown    *decimal.Decimal
	TriggerPercentUp     *decimal.Decimal
	TriggerPercentDown   *decimal.Decimal
	PullbackPercent      *decimal.Decimal
	PullbackSpread       *decimal.Decimal
	ReboundPercent       *decimal.Decimal
	ReboundSpread        *decimal.Decimal
	TriggerSellOrderType string
	TriggerBuyOrderType  string
	TriggerSellDepth     int32
	TriggerBuyDepth      int32
	TriggerQuantity      *decimal.Decimal
	TriggerSellQuantity  *decimal.Decimal
	TriggerBuyQuantity   *decimal.Decimal
	UpperLimitQuantity   *decimal.Decimal
	LowerLimitQuantity   *decimal.Decimal
	UpperLimitEvent      GridLimitEvent
	LowerLimitEvent      GridLimitEvent
	MultipleTrigger      bool
	TriggerTimes         int32
	TotalBuyQuantity     *decimal.Decimal
	TotalSellQuantity    *decimal.Decimal
	TotalProfitBalance   *decimal.Decimal
	SettlementCurrency   string
	TimeInForce          GridTimeInForce // Day / GTC / GTD
	Gtd                  string
	CreatedAt            string
	Rth                  int32
	SupportShortsell     bool
	GridOrderTypeUp      string // GMO / GLO / GTG
	GridOrderTypeDown    string // GMO / GLO / GTG
}

// GridOrderSubOrder is a triggered sub-order carried in the grid order detail.
type GridOrderSubOrder struct {
	Id          string
	Price       *decimal.Decimal
	OrderType   string
	Quantity    *decimal.Decimal
	ExecutedQty *decimal.Decimal
	Action      int32
	Status      string
	SubmittedAt string
	Rth         int32
}

// GridOrderHistory is a grid order lifecycle-history entry.
type GridOrderHistory struct {
	HistoryId     string
	CreatedAt     string
	Status        string
	SuspendReason string
	Reason        string
}

// GridOrderDetail is the detail of a grid trading order.
type GridOrderDetail struct {
	OrderId             string
	Symbol              string
	StockName           string
	Status              string
	GridStatus          string
	SuspendReason       string
	SleepingReason      string
	SubmittedBasePrice  *decimal.Decimal
	CurrentBasePrice    *decimal.Decimal
	UpperLimitPrice     *decimal.Decimal
	LowerLimitPrice     *decimal.Decimal
	TriggerPriceType    TriggerPriceType
	TriggerSpreadUp     *decimal.Decimal
	TriggerSpreadDown   *decimal.Decimal
	TriggerPercentUp    *decimal.Decimal
	TriggerPercentDown  *decimal.Decimal
	PullbackPercent     *decimal.Decimal
	PullbackSpread      *decimal.Decimal
	ReboundPercent      *decimal.Decimal
	ReboundSpread       *decimal.Decimal
	MultipleTrigger     bool
	TimeInForce         GridTimeInForce
	TriggerQuantity     *decimal.Decimal
	TriggerSellQuantity *decimal.Decimal
	TriggerBuyQuantity  *decimal.Decimal
	UpperLimitQuantity  *decimal.Decimal
	LowerLimitQuantity  *decimal.Decimal
	UpperLimitEvent     GridLimitEvent
	LowerLimitEvent     GridLimitEvent
	TriggerSellDepth    int32
	TriggerBuyDepth     int32
	CreatedAt           string
	UpdatedAt           string
	SettlementCurrency  string
	ExpireTime          string
	Gtd                 string
	GridSubOrders       []*GridOrderSubOrder
	SubHasMore          bool
	GridOrderHistory    []*GridOrderHistory
	HistoryHasMore      bool
	SupportShortsell    bool
	Rth                 int32
	GridOrderTypeUp     string
	GridOrderTypeDown   string
}

// TriggerOrder is a grid trigger-history entry (one triggered order).
type TriggerOrder struct {
	Id            string
	Status        string
	Name          string
	Symbol        string
	Price         *decimal.Decimal
	Quantity      *decimal.Decimal
	ExecutedPrice *decimal.Decimal
	ExecutedQty   *decimal.Decimal
	SubmittedAt   string
	Action        int32
	OrderType     string
	TriggerPrice  *decimal.Decimal
	Msg           string
	Currency      string
	LastDone      *decimal.Decimal
	UpdatedAt     string
	TimeInForce   GridTimeInForce
	Gtd           string
	TriggerAt     string
	TriggerStatus int32
}

// GridBidSize is a price-step (bid-size) rule entry from the order-info response.
type GridBidSize struct {
	StrProceed *decimal.Decimal
	EndProceed *decimal.Decimal
	BidSize    *decimal.Decimal
}

// GridChannelInfo is the channel / authorization info nested in order-info.
type GridChannelInfo struct {
	StrategyGranted    bool
	SupportRth         bool
	Currency           string
	SettlementCurrency []string
}

// GridOrderInfo is the /v1/orders/info response used by the grid order window.
type GridOrderInfo struct {
	Name         string
	LastDone     *decimal.Decimal
	LotSize      *decimal.Decimal
	BuyLotSize   *decimal.Decimal
	SellLotSize  *decimal.Decimal
	BidSizes     []*GridBidSize
	ChannelInfos GridChannelInfo
}
