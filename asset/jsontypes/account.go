package jsontypes

import "encoding/json"

// The types below are the JSON shapes of the account money-movement and
// reward endpoints. Upstream fields that carry personal data (ID-document and
// proof-document storage keys, postal addresses, contact people, coupon
// secrets, delivery details, receipt images) are deliberately not modelled,
// so they are dropped while decoding.

// BankCardAuth is one entry of a bank card's auth_list.
type BankCardAuth struct {
	FlowID                  FlexString `json:"flow_id"`
	Name                    FlexString `json:"name"`
	Currency                FlexString `json:"currency"`
	Status                  FlexInt32  `json:"status"`
	Message                 FlexString `json:"message"`
	IsSupport               FlexBool   `json:"is_support"`
	MaxAmount               FlexString `json:"max_amount"`
	MaxAmountUnlimited      FlexBool   `json:"max_amount_unlimited"`
	DailyMaxCount           FlexInt64  `json:"daily_max_count"`
	DailyMaxCountUnlimited  FlexBool   `json:"daily_max_count_unlimited"`
	ExpiryDate              FlexString `json:"expiry_date"`
	LongTermValid           FlexBool   `json:"long_term_valid"`
	DocReviewing            FlexBool   `json:"doc_reviewing"`
	NoPriorDepositReviewing FlexBool   `json:"no_prior_deposit_reviewing"`
}

// BankCard is one card of GET /v2/account/bank-cards.
type BankCard struct {
	ID             FlexInt64      `json:"id"`
	BankID         FlexString     `json:"bank_id"`
	Bank           FlexString     `json:"bank"`
	BankEn         FlexString     `json:"bank_en"`
	BankCode       FlexString     `json:"bank_code"`
	Account        FlexString     `json:"account"`
	AccountType    FlexString     `json:"account_type"`
	Country        FlexString     `json:"country"`
	CountryDisplay FlexString     `json:"country_display"`
	Region         FlexString     `json:"region"`
	RegionName     FlexString     `json:"region_name"`
	SwiftCode      FlexString     `json:"swift_code"`
	Remark         FlexString     `json:"remark"`
	Status         FlexInt32      `json:"status"`
	WitnessState   FlexInt32      `json:"witness_state"`
	SupportWitness FlexBool       `json:"support_witness"`
	BindState      FlexInt32      `json:"bind_state"`
	BindTip        FlexString     `json:"bind_tip"`
	CreatedAt      FlexString     `json:"created_at"`
	WithdrawFee    FlexString     `json:"withdraw_fee"`
	AuthList       []BankCardAuth `json:"auth_list"`
}

// BankCardsResponse is the response of GET /v2/account/bank-cards.
type BankCardsResponse struct {
	List []BankCard `json:"list"`
}

// EddaApplication is one eDDA authorization application.
type EddaApplication struct {
	ID        FlexString `json:"id"`
	FlowID    FlexString `json:"flow_id"`
	Name      FlexString `json:"name"`
	Currency  FlexString `json:"currency"`
	Status    FlexInt32  `json:"status"`
	Message   FlexString `json:"message"`
	CreatedAt FlexString `json:"created_at"`
	UpdatedAt FlexString `json:"updated_at"`
}

// EddaStatusResponse is the response of GET /v1/account/edda-status. Current
// builds return {list[]}; older builds return a single application object, so
// the raw body is kept and decoded by the caller.
type EddaStatusResponse = json.RawMessage

// TransferStock is one security of a stock transfer record.
type TransferStock struct {
	Code             FlexString `json:"code"`
	Name             FlexString `json:"name"`
	NameEn           FlexString `json:"name_en"`
	Market           FlexString `json:"market"`
	CounterID        FlexString `json:"counter_id"`
	Amount           FlexString `json:"amount"`
	Currency         FlexString `json:"currency"`
	Cost             FlexString `json:"cost"`
	State            FlexInt32  `json:"state"`
	StreetAmount     FlexString `json:"street_amount"`
	NomineeAmount    FlexString `json:"nominee_amount"`
	OwnAmount        FlexString `json:"own_amount"`
	CertificateCount FlexString `json:"certificate_count"`
	PositionSide     FlexString `json:"position_side"`
}

// TransferRecord is one record of GET /v1/account/stock-transfers/in and /out.
// The holder name at the other broker (record-level "name"), contact details
// and attachments are not modelled.
type TransferRecord struct {
	ID        FlexString      `json:"id"`
	OrderID   FlexString      `json:"order_id"`
	Market    FlexString      `json:"market"`
	Broker    FlexString      `json:"broker"`
	BrokerID  FlexString      `json:"broker_id"`
	BrokerEn  FlexString      `json:"broker_en"`
	NameEn    FlexString      `json:"name_en"`
	Code      FlexString      `json:"code"`
	Account   FlexString      `json:"account"`
	State     FlexString      `json:"state"`
	StateCode FlexInt32       `json:"state_code"`
	AppState  FlexString      `json:"app_state"`
	CanModify FlexBool        `json:"can_modify"`
	CreatedAt FlexString      `json:"created_at"`
	UpdatedAt FlexString      `json:"updated_at"`
	ProcessAt FlexString      `json:"process_at"`
	ExecAt    FlexString      `json:"exec_at"`
	AgreeAt   FlexString      `json:"agree_at"`
	RefuseAt  FlexString      `json:"refuse_at"`
	CancelAt  FlexString      `json:"cancel_at"`
	Remark    FlexString      `json:"remark"`
	Reason    FlexString      `json:"reason"`
	Stocks    []TransferStock `json:"stocks"`
}

// TransferRecordsResponse is the response of the stock transfer record lists.
type TransferRecordsResponse struct {
	Records []TransferRecord `json:"records"`
	Total   FlexInt64        `json:"total"`
}

// RewardDetail is the cash{}, card{} or prize{} object of a reward item.
// Coupon secrets, cash card codes, delivery details and addresses are not
// modelled.
type RewardDetail struct {
	RecordID        FlexString `json:"record_id"`
	ActivityTitle   FlexString `json:"activity_title"`
	Title           FlexString `json:"title"`
	Description     FlexString `json:"description"`
	CardType        FlexString `json:"card_type"`
	CardStatus      FlexInt32  `json:"card_status"`
	PrizeType       FlexInt32  `json:"prize_type"`
	Amount          FlexString `json:"amount"`
	Currency        FlexString `json:"currency"`
	Status          FlexInt32  `json:"status"`
	DeliverStatus   FlexInt32  `json:"deliver_status"`
	ExitStatus      FlexInt32  `json:"exit_status"`
	SendTime        FlexString `json:"send_time"`
	StartTime       FlexString `json:"start_time"`
	EndTime         FlexString `json:"end_time"`
	EffectiveTime   FlexString `json:"effective_time"`
	ReceiveDeadline FlexString `json:"receive_deadline"`
}

// RewardItem is one item of GET /v1/account/rewards.
type RewardItem struct {
	RewardType *FlexInt32    `json:"reward_type"`
	Cash       *RewardDetail `json:"cash"`
	Card       *RewardDetail `json:"card"`
	Prize      *RewardDetail `json:"prize"`
}

// RewardsResponse is the response of GET /v1/account/rewards.
type RewardsResponse struct {
	List  []RewardItem `json:"list"`
	Total FlexInt64    `json:"total"`
}

// CouponCommission is the commission{} object of a coupon card.
type CouponCommission struct {
	Markets               FlexStrings `json:"markets"`
	RebateRate            FlexString  `json:"rebate_rate"`
	BzTypes               []FlexInt32 `json:"bz_types"`
	CommissionBalance     FlexString  `json:"commission_balance"`
	CommissionDuration    FlexString  `json:"commission_duration"`
	Unlimited             FlexBool    `json:"unlimited"`
	CommissionExpiredTime FlexString  `json:"commission_expired_time"`
	UsageCount            FlexInt64   `json:"usage_count"`
	LaveUsageCount        FlexInt64   `json:"lave_usage_count"`
	SettlementModel       FlexString  `json:"settlement_model"`
}

// Coupon is one coupon card of GET /v1/account/cards. The coffee_coupon
// object (which can carry redemption codes) and app styling are not modelled.
type Coupon struct {
	CardID                FlexString        `json:"card_id"`
	CardType              FlexString        `json:"card_type"`
	Title                 FlexString        `json:"title"`
	TitleSec              FlexString        `json:"title_sec"`
	Description           FlexString        `json:"description"`
	Detail                FlexString        `json:"detail"`
	CardStatus            FlexInt32         `json:"card_status"`
	StartTime             FlexString        `json:"start_time"`
	EndTime               FlexString        `json:"end_time"`
	ActivateTime          FlexString        `json:"activate_time"`
	Amount                FlexString        `json:"amount"`
	Currency              FlexString        `json:"currency"`
	AvailableTradeSession FlexString        `json:"available_trade_session"`
	TradeSessions         FlexStrings       `json:"trade_sessions"`
	Commission            *CouponCommission `json:"commission"`
	Quote                 json.RawMessage   `json:"quote"`
	Stock                 json.RawMessage   `json:"stock"`
	PlatformFeeDeduct     json.RawMessage   `json:"platform_fee_deduct"`
}

// CouponsResponse is the response of GET /v1/account/cards.
type CouponsResponse struct {
	TopList      []Coupon  `json:"top_list"`
	List         []Coupon  `json:"list"`
	BottomList   []Coupon  `json:"bottom_list"`
	Total        FlexInt64 `json:"total"`
	ExpiredCount FlexInt64 `json:"expired_count"`
}

// CouponOrder is one order of GET /v1/account/card-orders.
type CouponOrder struct {
	OrderID            FlexString `json:"order_id"`
	State              FlexString `json:"state"`
	StateVal           FlexString `json:"state_val"`
	ExecutedPrice      FlexString `json:"executed_price"`
	ExecutedQuantity   FlexString `json:"executed_quantity"`
	Currency           FlexString `json:"currency"`
	SettlementCurrency FlexString `json:"settlement_currency"`
	Rebate             FlexString `json:"rebate"`
	RebateDt           FlexString `json:"rebate_dt"`
}

// CouponOrdersResponse is the response of GET /v1/account/card-orders.
type CouponOrdersResponse struct {
	Orders []CouponOrder `json:"orders"`
}

// Deposit is one item of GET /v1/account/deposits. Receipt images (vouchers)
// are not modelled.
type Deposit struct {
	ID                    FlexString `json:"id"`
	CreatedAt             FlexString `json:"created_at"`
	UpdatedAt             FlexString `json:"updated_at"`
	Currency              FlexString `json:"currency"`
	Amount                FlexString `json:"amount"`
	ReceivableAmount      FlexString `json:"receivable_amount"`
	State                 FlexString `json:"state"`
	StateCode             FlexInt32  `json:"state_code"`
	SubStateCode          FlexInt32  `json:"sub_state_code"`
	BankOperationTypeID   FlexString `json:"bank_operation_type_id"`
	BankOperationTypeName FlexString `json:"bank_operation_type_name"`
	ReferenceNo           FlexString `json:"reference_no"`
	CanCancel             FlexBool   `json:"can_cancel"`
	UserBankSwiftCode     FlexString `json:"user_bank_swift_code"`
}

// DepositsResponse is the response of GET /v1/account/deposits.
type DepositsResponse struct {
	Items []Deposit `json:"items"`
	Total FlexInt64 `json:"total"`
}

// Withdrawal is one item of GET /v1/account/withdrawals.
type Withdrawal struct {
	ID                 FlexString `json:"id"`
	OrderID            FlexString `json:"order_id"`
	CreatedAt          FlexString `json:"created_at"`
	Currency           FlexString `json:"currency"`
	Amount             FlexString `json:"amount"`
	Fee                FlexString `json:"fee"`
	SettlementAmount   FlexString `json:"settlement_amount"`
	State              FlexInt32  `json:"state"`
	RawState           FlexString `json:"raw_state"`
	StateDescription   FlexString `json:"state_description"`
	BankAccount        FlexString `json:"bank_account"`
	DisplayBankAccount FlexString `json:"display_bank_account"`
	UserBankSwiftCode  FlexString `json:"user_bank_swift_code"`
	AccountChannel     FlexString `json:"account_channel"`
}

// WithdrawalsResponse is the response of GET /v1/account/withdrawals.
type WithdrawalsResponse struct {
	List  []Withdrawal `json:"list"`
	Total FlexInt64    `json:"total"`
}
