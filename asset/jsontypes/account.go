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
	Status                  FlexCode   `json:"status"`
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
	ID             FlexString     `json:"id"`
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
	Status         FlexCode       `json:"status"`
	WitnessState   FlexCode       `json:"witness_state"`
	SupportWitness FlexBool       `json:"support_witness"`
	BindState      FlexCode       `json:"bind_state"`
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
	Status    FlexCode   `json:"status"`
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
	State            FlexCode   `json:"state"`
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
	StateCode FlexCode        `json:"state_code"`
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
	CardStatus      FlexCode   `json:"card_status"`
	PrizeType       FlexCode   `json:"prize_type"`
	Amount          FlexString `json:"amount"`
	Currency        FlexString `json:"currency"`
	Status          FlexCode   `json:"status"`
	DeliverStatus   FlexCode   `json:"deliver_status"`
	ExitStatus      FlexCode   `json:"exit_status"`
	SendTime        FlexString `json:"send_time"`
	StartTime       FlexString `json:"start_time"`
	EndTime         FlexString `json:"end_time"`
	EffectiveTime   FlexString `json:"effective_time"`
	ReceiveDeadline FlexString `json:"receive_deadline"`
}

// RewardItem is one item of GET /v1/account/rewards.
type RewardItem struct {
	RewardType FlexCode      `json:"reward_type"`
	Cash       *RewardDetail `json:"cash"`
	Card       *RewardDetail `json:"card"`
	Prize      *RewardDetail `json:"prize"`
}

// RewardsResponse is the response of GET /v1/account/rewards.
type RewardsResponse struct {
	List  []RewardItem `json:"list"`
	Total FlexInt64    `json:"total"`
}

// CouponCommission is the commission{} object of a coupon card. BzTypes and
// Unlimited are proto enums, sent as a name or a number.
type CouponCommission struct {
	Markets               FlexStrings  `json:"markets"`
	RebateRate            FlexString   `json:"rebate_rate"`
	BzTypes               []FlexString `json:"bz_types"`
	CommissionBalance     FlexString   `json:"commission_balance"`
	CommissionDuration    FlexInt64    `json:"commission_duration"`
	Unlimited             FlexString   `json:"unlimited"`
	CommissionExpiredTime FlexString   `json:"commission_expired_time"`
	UsageCount            FlexInt64    `json:"usage_count"`
	LaveUsageCount        FlexInt64    `json:"lave_usage_count"`
	SettlementModel       FlexString   `json:"settlement_model"`
	AmountLimit           FlexInt64    `json:"amount_limit"`
	Superposition         FlexBool     `json:"superposition"`
}

// CouponAiBenefit is the ai_benefit{} object of a coupon card.
type CouponAiBenefit struct {
	BenefitDays  FlexInt64  `json:"benefit_days"`
	BenefitStart FlexString `json:"benefit_start"`
	BenefitEnd   FlexString `json:"benefit_end"`
	BenefitTag   FlexString `json:"benefit_tag"`
}

// CouponCash is the cash{} object of a coupon card.
type CouponCash struct {
	ReceiveDuration FlexInt64 `json:"receive_duration"`
}

// CouponCoffeeStore is one store_list entry of a coffee coupon.
type CouponCoffeeStore struct {
	StoreID   FlexInt64  `json:"store_id"`
	StoreName FlexString `json:"store_name"`
}

// CouponCoffee is the coffee_coupon{} object of a coupon card.
type CouponCoffee struct {
	CouponChannel       FlexString          `json:"coupon_channel"`
	CouponID            FlexString          `json:"coupon_id"`
	CouponNo            FlexString          `json:"coupon_no"`
	IsAvailableAllStore FlexBool            `json:"is_available_all_store"`
	PeriodType          FlexString          `json:"period_type"`
	SettledStoreID      FlexInt64           `json:"settled_store_id"`
	StoreIDList         FlexInt64s          `json:"store_id_list"`
	StoreList           []CouponCoffeeStore `json:"store_list"`
	Threshold           FlexInt64           `json:"threshold"`
	UseLimit            FlexString          `json:"use_limit"`
}

// CouponCryptoCounter is one optional_counters entry of a crypto coupon.
type CouponCryptoCounter struct {
	CounterID FlexString `json:"counter_id"`
}

// CouponCryptoFallback is the fallback{} object of a crypto coupon.
type CouponCryptoFallback struct {
	Enabled   FlexBool   `json:"enabled"`
	CounterID FlexString `json:"counter_id"`
}

// CouponCrypto is the crypto{} object of a coupon card.
type CouponCrypto struct {
	Fallback         *CouponCryptoFallback `json:"fallback"`
	OptionalCounters []CouponCryptoCounter `json:"optional_counters"`
}

// CouponFinanceInterest is the finance_interest{} object of a coupon card.
type CouponFinanceInterest struct {
	DiscountMode  FlexInt64  `json:"discount_mode"`
	DiscountRate  FlexString `json:"discount_rate"`
	Superposition FlexBool   `json:"superposition"`
}

// CouponFundDeduct is the fund_deduct{} object of a coupon card.
type CouponFundDeduct struct {
	MaxRiskLevel     FlexInt64   `json:"max_risk_level"`
	MinOrderAmount   FlexString  `json:"min_order_amount"`
	OrderTypes       FlexStrings `json:"order_types"`
	UnderlyingType   FlexString  `json:"underlying_type"`
	UnderlyingValues FlexStrings `json:"underlying_values"`
	UnderlyingNames  FlexStrings `json:"underlying_names"`
}

// CouponQuote is the quote{} / fundamental{} object of a coupon card.
// InvestorType is a proto enum, sent as a name or a number.
type CouponQuote struct {
	CurrentEndTime     FlexString  `json:"current_end_time"`
	Inapplicable       FlexBool    `json:"inapplicable"`
	InvestorType       FlexString  `json:"investor_type"`
	ManualActivate     FlexBool    `json:"manual_activate"`
	Markets            FlexStrings `json:"markets"`
	MemberInvestorType FlexInt64   `json:"member_investor_type"`
	PairSpuNo          FlexString  `json:"pair_spu_no"`
	QuoteDuration      FlexInt64   `json:"quote_duration"`
	QuoteMarket        FlexString  `json:"quote_market"`
	QuoteTags          FlexStrings `json:"quote_tags"`
	Superposition      FlexBool    `json:"superposition"`
}

// CouponIpo is the ipo{} object of a coupon card.
type CouponIpo struct {
	SubscriptionID FlexString `json:"subscription_id"`
}

// CouponIpoCommission is the ipo_commission{} object of a coupon card.
type CouponIpoCommission struct {
	BzTypes        []FlexString `json:"bz_types"`
	CounterID      FlexString   `json:"counter_id"`
	CounterName    FlexString   `json:"counter_name"`
	FinancingLimit FlexString   `json:"financing_limit"`
}

// CouponIpoInterest is the ipo_interest{} object of a coupon card.
type CouponIpoInterest struct {
	DiscountRate FlexString `json:"discount_rate"`
}

// CouponPlatformFeeDeduct is the platform_fee_deduct{} object of a coupon
// card.
type CouponPlatformFeeDeduct struct {
	BzTypes        []FlexString `json:"bz_types"`
	CounterID      FlexString   `json:"counter_id"`
	CounterName    FlexString   `json:"counter_name"`
	Markets        FlexStrings  `json:"markets"`
	UsageCount     FlexInt64    `json:"usage_count"`
	LaveUsageCount FlexInt64    `json:"lave_usage_count"`
	Superposition  FlexBool     `json:"superposition"`
}

// CouponStockCounter is one optional_counters entry of a stock coupon.
type CouponStockCounter struct {
	CounterID FlexString `json:"counter_id"`
	Shares    FlexString `json:"shares"`
}

// CouponStockFallback is the fallback{} object of a stock coupon.
type CouponStockFallback struct {
	Enabled   FlexBool   `json:"enabled"`
	Type      FlexString `json:"type"`
	CounterID FlexString `json:"counter_id"`
	Shares    FlexString `json:"shares"`
}

// CouponStock is the stock{} object of a coupon card.
type CouponStock struct {
	CounterID        FlexString           `json:"counter_id"`
	CounterName      FlexString           `json:"counter_name"`
	Amount           FlexString           `json:"amount"`
	RedemptionType   FlexString           `json:"redemption_type"`
	SendWay          FlexString           `json:"send_way"`
	OptionalCounters []CouponStockCounter `json:"optional_counters"`
	Fallback         *CouponStockFallback `json:"fallback"`
}

// CouponStockFragment is the stock_fragment{} object of a coupon card.
// UnderlyingType is a proto enum, sent as a name or a number.
type CouponStockFragment struct {
	CounterID      FlexString    `json:"counter_id"`
	CounterName    FlexString    `json:"counter_name"`
	Market         FlexString    `json:"market"`
	Markets        FlexStrings   `json:"markets"`
	UnderlyingType FlexString    `json:"underlying_type"`
	Amounts        FlexStringMap `json:"amounts"`
}

// Coupon is one coupon card of GET /v1/account/cards. Every per-type terms
// object is modelled (upstream sends null for the types that do not apply)
// except fund_interest, which the proto describes as an arbitrary object;
// app styling and the deprecated top-level duplicates of the commission
// terms are dropped.
type Coupon struct {
	CardID                FlexString               `json:"card_id"`
	CardType              FlexString               `json:"card_type"`
	Title                 FlexString               `json:"title"`
	TitleSec              FlexString               `json:"title_sec"`
	Description           FlexString               `json:"description"`
	Detail                FlexString               `json:"detail"`
	CardStatus            FlexCode                 `json:"card_status"`
	StartTime             FlexString               `json:"start_time"`
	EndTime               FlexString               `json:"end_time"`
	ActivateTime          FlexString               `json:"activate_time"`
	Amount                FlexString               `json:"amount"`
	Currency              FlexString               `json:"currency"`
	AvailableTradeSession FlexString               `json:"available_trade_session"`
	TradeSessions         FlexStrings              `json:"trade_sessions"`
	SubscriptionID        FlexString               `json:"subscription_id"`
	AiBenefit             *CouponAiBenefit         `json:"ai_benefit"`
	Cash                  *CouponCash              `json:"cash"`
	CoffeeCoupon          *CouponCoffee            `json:"coffee_coupon"`
	Commission            *CouponCommission        `json:"commission"`
	Crypto                *CouponCrypto            `json:"crypto"`
	FinanceInterest       *CouponFinanceInterest   `json:"finance_interest"`
	FundDeduct            *CouponFundDeduct        `json:"fund_deduct"`
	Fundamental           *CouponQuote             `json:"fundamental"`
	Ipo                   *CouponIpo               `json:"ipo"`
	IpoCommission         *CouponIpoCommission     `json:"ipo_commission"`
	IpoInterest           *CouponIpoInterest       `json:"ipo_interest"`
	PlatformFeeDeduct     *CouponPlatformFeeDeduct `json:"platform_fee_deduct"`
	Quote                 *CouponQuote             `json:"quote"`
	Stock                 *CouponStock             `json:"stock"`
	StockFragment         *CouponStockFragment     `json:"stock_fragment"`
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
	StateCode             FlexCode   `json:"state_code"`
	SubStateCode          FlexCode   `json:"sub_state_code"`
	BankOperationTypeID   FlexString `json:"bank_operation_type_id"`
	BankOperationTypeName FlexString `json:"bank_operation_type_name"`
	ReferenceNo           FlexString `json:"reference_no"`
	CanCancel             FlexBool   `json:"can_cancel"`
	UserBankSwiftCode     FlexString `json:"user_bank_swift_code"`
}

// DepositsResponse is the response of GET /v1/account/deposits. Its total is
// not modelled: upstream caps it at the page size.
type DepositsResponse struct {
	Items []Deposit `json:"items"`
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
	State              FlexCode   `json:"state"`
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
