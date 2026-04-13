package jsontypes

// SubmitOrderResponse is the JSON response for submit/amend IPO order.
type SubmitOrderResponse struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	Market            string `json:"market"`
	SubQty            string `json:"sub_qty"`
	SubAmount         string `json:"sub_amount"`
	HandingFee        string `json:"handing_fee"`
	WithdrawAble      bool   `json:"withdraw_able"`
	LotWinQty         string `json:"lot_win_qty"`
	Status            int32  `json:"status"`
	Currency          string `json:"currency"`
	FinancingAmount   string `json:"financing_amount"`
	NeedToPay         string `json:"need_to_pay"`
	IpoStatus         int32  `json:"ipo_status"`
	Method            int32  `json:"method"`
	FinancingInterest string `json:"financing_interest"`
	FinanceFeeRate    string `json:"finance_fee_rate"`
	IpoPrice          string `json:"ipo_price"`
	BorrowAmount      string `json:"borrow_amount"`
	BorrowInterest    string `json:"borrow_interest"`
	TotalAmount       string `json:"total_amount"`
	CurrentAmount     string `json:"current_amount"`
	RefundAmount      string `json:"refund_amount"`
	LuckAt            int64  `json:"luck_at"`
	IpoDate           int64  `json:"ipo_date"`
	MartBegin         int64  `json:"mart_begin"`
	MartEnd           int64  `json:"mart_end"`
	RejectText        string `json:"reject_text"`
	ModifyEnable      bool   `json:"modify_enable"`
	BatchID           int64  `json:"batch_id"`
	BorrowEnable      bool   `json:"borrow_enable"`
	LuckyFee          string `json:"lucky_fee"`
}

// WithdrawOrderResponse is the JSON response for withdraw IPO order.
type WithdrawOrderResponse struct {
	Result bool `json:"result"`
}

// OrderListItem is one item in the IPO order list.
type OrderListItem struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Status          int32  `json:"status"`
	SubQty          string `json:"sub_qty"`
	Method          int32  `json:"method"`
	CreatedAt       int64  `json:"created_at"`
	Currency        string `json:"currency"`
	LotWinQty       string `json:"lot_win_qty"`
	SubAmount       string `json:"sub_amount"`
	CashAmount      string `json:"cash_amount"`
	FinancingAmount string `json:"financing_amount"`
	FinancingRatio  string `json:"financing_ratio"`
	HandingFee      string `json:"handing_fee"`
}

// FetchOrderListResponse is the JSON response for order list.
type FetchOrderListResponse struct {
	List    []*OrderListItem `json:"list"`
	HasMore bool             `json:"has_more"`
}

// FetchOrderDetailResponse is the JSON response for order detail.
type FetchOrderDetailResponse struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Code              string `json:"code"`
	Market            string `json:"market"`
	SubQty            string `json:"sub_qty"`
	SubAmount         string `json:"sub_amount"`
	HandingFee        string `json:"handing_fee"`
	WithdrawAble      bool   `json:"withdraw_able"`
	LotWinQty         string `json:"lot_win_qty"`
	Status            int32  `json:"status"`
	Currency          string `json:"currency"`
	FinancingAmount   string `json:"financing_amount"`
	NeedToPay         string `json:"need_to_pay"`
	IpoStatus         int32  `json:"ipo_status"`
	Method            int32  `json:"method"`
	FinancingInterest string `json:"financing_interest"`
	FinanceFeeRate    string `json:"finance_fee_rate"`
	Explanation       string `json:"explanation"`
	IpoPrice          string `json:"ipo_price"`
	BorrowAmount      string `json:"borrow_amount"`
	BorrowInterest    string `json:"borrow_interest"`
	Symbol            string `json:"symbol"`
	TotalAmount       string `json:"total_amount"`
	CurrentAmount     string `json:"current_amount"`
	RefundAmount      string `json:"refund_amount"`
	LuckAt            int64  `json:"luck_at"`
	RejectText        string `json:"reject_text"`
	Multiple          string `json:"multiple"`
	ModifyEnable      bool   `json:"modify_enable"`
	LuckyFee          string `json:"lucky_fee"`
	Channel           string `json:"channel"`
	BatchID           int64  `json:"batch_id"`
}

// MarginItem is one financing scheme item.
type MarginItem struct {
	Method           int32  `json:"method"`
	HandingFee       string `json:"handing_fee"`
	Currency         string `json:"currency"`
	Deadline         int64  `json:"deadline"`
	Multiple         string `json:"multiple"`
	FinancingFeeRate string `json:"financing_fee_rate"`
	RemainAmount     string `json:"remain_amount"`
	BatchID          int64  `json:"batch_id"`
	IsEnded          bool   `json:"is_ended"`
	FinancingRatio   string `json:"financing_ratio"`
	StartAt          int64  `json:"start_at"`
	DisplayName      string `json:"display_name"`
	MinCash          string `json:"min_cash"`
}

// FetchMarginListResponse is the JSON response for margin/financing scheme list.
type FetchMarginListResponse struct {
	List []*MarginItem `json:"list"`
}

// PaymentListItem is one item in the IPO payment/lot list.
type PaymentListItem struct {
	Number    int32  `json:"number"`
	SubQty    string `json:"sub_qty"`
	SubAmount string `json:"sub_amount"`
	Currency  string `json:"currency"`
}

// FetchIpoPaymentListResponse is the JSON response for IPO payment list.
type FetchIpoPaymentListResponse struct {
	List []*PaymentListItem `json:"list"`
}

// FetchBuyLimitResponse is the JSON response for buying power.
type FetchBuyLimitResponse struct {
	TotalAmount   string `json:"total_amount"`
	CurrentAmount string `json:"current_amount"`
	MmfAmount     string `json:"mmf_amount"`
	CashAmount    string `json:"cash_amount"`
}
