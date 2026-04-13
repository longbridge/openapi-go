package ipo

// SubscriptionMethod represents the IPO subscription method/channel (认购方式).
type SubscriptionMethod int32

// SubscriptionStatus represents the order status of an IPO subscription.
type SubscriptionStatus int32

// SubmitOrderResult is the result returned by SubmitOrder and AmendOrder.
type SubmitOrderResult struct {
	ID                int64
	Name              string
	Symbol            string
	Market            string
	SubQty            string // 申购数量
	SubAmount         string // 申购金额
	HandingFee        string // 申购手续费
	WithdrawAble      bool   // 能否撤单
	LotWinQty         string // 中签数量
	Status            SubscriptionStatus
	Currency          string
	FinancingAmount   string
	NeedToPay         string             // 需补缴
	IpoStatus         int32
	Method            SubscriptionMethod // 认购方式
	FinancingInterest string             // 融资利息
	FinanceFeeRate    string             // 融资利率
	IpoPrice          string
	BorrowAmount      string
	BorrowInterest    string
	TotalAmount       string // 合计金额
	CurrentAmount     string // 当前币种现金（港币现金）
	RefundAmount      string // 需退款金额
	LuckAt            int64
	IpoDate           int64  // 上市时间
	MartBegin         int64  // 暗盘开始时间
	MartEnd           int64  // 暗盘结束时间
	RejectText        string // 拒绝理由
	ModifyEnable      bool   // 是否能够改单
	BatchID           int64
	BorrowEnable      bool   // 是否支持借币
	LuckyFee          string
}

// OrderListItem is a summary item in the IPO order list.
type OrderListItem struct {
	ID              int64
	Name            string
	Symbol          string
	Status          SubscriptionStatus
	SubQty          string
	Method          SubscriptionMethod
	CreatedAt       int64
	Currency        string
	LotWinQty       string
	SubAmount       string
	CashAmount      string
	FinancingAmount string
	FinancingRatio  string
	HandingFee      string
}

// OrderDetail is the full detail of an IPO subscription order.
type OrderDetail struct {
	ID                int64
	Name              string
	Code              string
	Market            string
	SubQty            string // 申购数量
	SubAmount         string // 申购金额
	HandingFee        string // 申购手续费
	WithdrawAble      bool   // 能否撤单
	LotWinQty         string // 中签数量
	Status            SubscriptionStatus
	Currency          string
	FinancingAmount   string
	NeedToPay         string             // 需补缴
	IpoStatus         int32
	Method            SubscriptionMethod // 认购方式
	FinancingInterest string             // 融资利息
	FinanceFeeRate    string             // 融资利率
	Explanation       string             // 需补缴文案
	IpoPrice          string
	BorrowAmount      string
	BorrowInterest    string
	Symbol            string
	TotalAmount       string // 合计金额
	CurrentAmount     string // 当前币种现金（港币现金）
	RefundAmount      string // 需退款金额
	LuckAt            int64
	RejectText        string // 拒绝理由
	Multiple          string
	ModifyEnable      bool   // 是否能够改单
	LuckyFee          string
	Channel           string
	BatchID           int64
}

// MarginItem is a financing scheme for an IPO (融资方案).
type MarginItem struct {
	Method           SubscriptionMethod
	HandingFee       string
	Currency         string
	Deadline         int64
	Multiple         string
	FinancingFeeRate string
	RemainAmount     string
	BatchID          int64
	IsEnded          bool
	FinancingRatio   string
	StartAt          int64
	DisplayName      string
	MinCash          string
}

// PaymentListItem is one selectable lot quantity and its subscription amount.
type PaymentListItem struct {
	Number    int32  // lot number option
	SubQty    string // subscription quantity
	SubAmount string
	Currency  string
}

// BuyingPower contains the user's available funds for IPO subscription (打新购买力).
type BuyingPower struct {
	TotalAmount   string
	CurrentAmount string // 当前币种现金（港币现金）
	MmfAmount     string
	CashAmount    string
}

// SubmitOrderRequest contains parameters for submitting an IPO subscription order.
type SubmitOrderRequest struct {
	Symbol  string             // required
	SubQty  string             // required, subscription quantity
	BatchID int64              // required
	Method  SubscriptionMethod // required, subscription method
}

// AmendOrderRequest contains parameters for modifying an IPO subscription order.
type AmendOrderRequest struct {
	Symbol  string             // required
	SubQty  string             // required
	BatchID int64              // required
	Method  SubscriptionMethod // required
	OrderID int64              // required
}

// WithdrawOrderRequest contains parameters for cancelling an IPO subscription.
type WithdrawOrderRequest struct {
	OrderID int64 // required
}

// FetchOrderListRequest contains parameters for querying IPO order list.
type FetchOrderListRequest struct {
	Symbol   string // optional
	Page     int32  // optional
	PageSize int32  // optional
}

// FetchOrderDetailRequest contains parameters for querying IPO order detail.
type FetchOrderDetailRequest struct {
	OrderID int64 // required
}

// FetchMarginListRequest contains parameters for querying financing schemes.
type FetchMarginListRequest struct {
	Symbol string // required
}

// FetchIpoPaymentListRequest contains parameters for querying lot quantity options.
type FetchIpoPaymentListRequest struct {
	Symbol string // required
}

// FetchBuyLimitRequest contains parameters for querying buying power.
type FetchBuyLimitRequest struct {
	Symbol   string // required
	Currency string // required
}
