package asset

import "time"

// Timestamps in the account types below are time.Time in UTC. Upstream fills
// unset times with placeholders (0, "0", a negative value, a 1970 or 1906
// date, or a far-future sentinel such as 9999-12-31); those decode to the zero
// time.Time, so check IsZero() before using one. encoding/json marshals a zero
// value as "0001-01-01T00:00:00Z".
//
// Every integer status code type below has an XxxUnknown constant (-1): a code
// that upstream omitted, sent as null or sent as something that is not an
// integer decodes to it, so it can never be mistaken for a real state (0 is a
// real state in most of these tables). An undocumented integer code is kept
// verbatim.

// BankCardStatus is the verification status of a bank card.
// Undocumented integer codes are kept verbatim.
type BankCardStatus int32

const (
	// BankCardStatusUnknown means upstream sent no usable status.
	BankCardStatusUnknown BankCardStatus = -1
	// BankCardStatusUnverified means the card has not been verified.
	BankCardStatusUnverified BankCardStatus = 0
	// BankCardStatusUnderReview means the card is under manual review.
	BankCardStatusUnderReview BankCardStatus = 1
	// BankCardStatusVerified means the card is verified.
	BankCardStatusVerified BankCardStatus = 2
)

// BankCardWitnessState is the witness (identity verification) state of a bank
// card. Undocumented integer codes are kept verbatim.
type BankCardWitnessState int32

const (
	// BankCardWitnessStateUnknown means upstream sent no usable state.
	BankCardWitnessStateUnknown BankCardWitnessState = -1
	// BankCardWitnessStateNotRequired means no witness is required.
	BankCardWitnessStateNotRequired BankCardWitnessState = 1
	// BankCardWitnessStatePending means the witness is pending.
	BankCardWitnessStatePending BankCardWitnessState = 2
	// BankCardWitnessStateWitnessed means the card has been witnessed.
	BankCardWitnessStateWitnessed BankCardWitnessState = 3
)

// BankCardBindState is the binding state of a bank card.
// Undocumented integer codes are kept verbatim.
type BankCardBindState int32

const (
	// BankCardBindStateUnknown means upstream sent no usable state.
	BankCardBindStateUnknown BankCardBindState = -1
	// BankCardBindStateBinding means the card is being bound.
	BankCardBindStateBinding BankCardBindState = 0
	// BankCardBindStateBound means the card is bound.
	BankCardBindStateBound BankCardBindState = 1
	// BankCardBindStateFailed means binding failed.
	BankCardBindStateFailed BankCardBindState = 2
)

// AuthorizationStatus is the status of a bank card authorization (eDDA,
// bank-securities link, FPS) and of an eDDA application.
// Undocumented integer codes are kept verbatim.
type AuthorizationStatus int32

const (
	// AuthorizationStatusUnknown means upstream sent no usable status.
	AuthorizationStatusUnknown AuthorizationStatus = -1
	// AuthorizationStatusUnauthorized means not authorized.
	AuthorizationStatusUnauthorized AuthorizationStatus = 0
	// AuthorizationStatusAuthorized means authorized.
	AuthorizationStatusAuthorized AuthorizationStatus = 1
	// AuthorizationStatusAuthorizing means the authorization is in progress.
	AuthorizationStatusAuthorizing AuthorizationStatus = 2
	// AuthorizationStatusFailed means the authorization failed; see the
	// failure reason in Message.
	AuthorizationStatusFailed AuthorizationStatus = 3
)

// Authorization flows found in BankCardAuth.FlowID.
const (
	// AuthFlowHKEDDA is the bank-side eDDA flow.
	AuthFlowHKEDDA = "hk_edda"
	// AuthFlowHKSDDA is the in-app eDDA flow.
	AuthFlowHKSDDA = "hk_sdda"
	// AuthFlowHKBankSec is the bank-securities link flow.
	AuthFlowHKBankSec = "hk_bank_sec"
	// AuthFlowHKFPS is the FPS flow.
	AuthFlowHKFPS = "hk_fps"
)

// BankCardAuth is one authorization of a bank card (eDDA, bank-securities
// link, FPS), in one currency.
type BankCardAuth struct {
	// FlowID identifies the flow, e.g. AuthFlowHKEDDA; see IsEDDA.
	FlowID   string              `json:"flow_id"`
	Name     string              `json:"name"`
	Currency string              `json:"currency"`
	Status   AuthorizationStatus `json:"status"`
	// Message is the failure reason, if any.
	Message   string `json:"message"`
	IsSupport bool   `json:"is_support"`
	// MaxAmount is the per-transaction limit ("" when unset).
	MaxAmount          string `json:"max_amount"`
	MaxAmountUnlimited bool   `json:"max_amount_unlimited"`
	// DailyMaxCount is the daily transaction count limit (0 when unset).
	DailyMaxCount          int64 `json:"daily_max_count"`
	DailyMaxCountUnlimited bool  `json:"daily_max_count_unlimited"`
	// ExpiryDate is when the authorization expires; zero when unset or when
	// LongTermValid is true (upstream then sends a 9999-12-31 sentinel).
	ExpiryDate              time.Time `json:"expiry_date"`
	LongTermValid           bool      `json:"long_term_valid"`
	DocReviewing            bool      `json:"doc_reviewing"`
	NoPriorDepositReviewing bool      `json:"no_prior_deposit_reviewing"`
}

// IsEDDA reports whether the authorization is an eDDA flow (hk_edda or
// hk_sdda).
func (a *BankCardAuth) IsEDDA() bool {
	return a.FlowID == AuthFlowHKEDDA || a.FlowID == AuthFlowHKSDDA
}

// BankCard is a linked bank card. The account holder's name, postal address
// and ID-document references returned upstream are not modelled.
type BankCard struct {
	// ID is the bank card id, a string of digits (use it with
	// EddaApplications).
	ID       string `json:"id"`
	BankID   string `json:"bank_id"`
	Bank     string `json:"bank"`
	BankEn   string `json:"bank_en"`
	BankCode string `json:"bank_code"`
	// Account is the full account number, as returned upstream.
	Account     string `json:"account"`
	AccountType string `json:"account_type"`
	// Country is the bank's country or region name (e.g. "Hong Kong"), not a
	// currency.
	Country        string               `json:"country"`
	CountryDisplay string               `json:"country_display"`
	Region         string               `json:"region"`
	RegionName     string               `json:"region_name"`
	SwiftCode      string               `json:"swift_code"`
	Remark         string               `json:"remark"`
	Status         BankCardStatus       `json:"status"`
	WitnessState   BankCardWitnessState `json:"witness_state"`
	SupportWitness bool                 `json:"support_witness"`
	BindState      BankCardBindState    `json:"bind_state"`
	BindTip        string               `json:"bind_tip"`
	// CreatedAt is when the card was added (zero when unset).
	CreatedAt   time.Time      `json:"created_at"`
	WithdrawFee string         `json:"withdraw_fee"`
	AuthList    []BankCardAuth `json:"auth_list"`
}

// EddaApplication is one eDDA authorization application of a bank card.
type EddaApplication struct {
	ID       string              `json:"id"`
	FlowID   string              `json:"flow_id"`
	Name     string              `json:"name"`
	Currency string              `json:"currency"`
	Status   AuthorizationStatus `json:"status"`
	// Message is the failure reason, if any.
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TransferInState is the state_code of a stock transfer-in record.
// Undocumented integer codes are kept verbatim.
type TransferInState int32

const (
	// TransferInStateUnknown means upstream sent no usable state code.
	TransferInStateUnknown TransferInState = -1
	// TransferInStateSubmitted means the application was submitted.
	TransferInStateSubmitted TransferInState = 0
	// TransferInStateCancelled means the application was cancelled.
	TransferInStateCancelled TransferInState = 1
	// TransferInStateCompleted means the transfer completed.
	TransferInStateCompleted TransferInState = 3
	// TransferInStateProcessing means the transfer is in progress.
	TransferInStateProcessing TransferInState = 5
	// TransferInStatePartiallyCompleted means some securities were transferred.
	TransferInStatePartiallyCompleted TransferInState = 7
	// TransferInStateFailed means the transfer failed.
	TransferInStateFailed TransferInState = 8
)

// TransferOutState is the state_code of a stock transfer-out record.
// Undocumented integer codes are kept verbatim.
type TransferOutState int32

const (
	// TransferOutStateUnknown means upstream sent no usable state code.
	TransferOutStateUnknown TransferOutState = -1
	// TransferOutStateSubmitted means the application was submitted; only
	// this state can be cancelled.
	TransferOutStateSubmitted TransferOutState = 0
	// TransferOutStateProcessing means the application is being processed.
	TransferOutStateProcessing TransferOutState = 1
	// TransferOutStateExecuting means the transfer is executing.
	TransferOutStateExecuting TransferOutState = 2
	// TransferOutStateCompleted means the transfer completed.
	TransferOutStateCompleted TransferOutState = 3
	// TransferOutStateCancelled means the application was cancelled.
	TransferOutStateCancelled TransferOutState = 4
	// TransferOutStatePartiallyCompleted means some securities were
	// transferred.
	TransferOutStatePartiallyCompleted TransferOutState = 5
	// TransferOutStateFailed means the transfer failed.
	TransferOutStateFailed TransferOutState = 6
	// TransferOutStateRejected means the application was rejected.
	TransferOutStateRejected TransferOutState = 7
)

// TransferStockState is the state of one security in a stock transfer.
// Undocumented integer codes are kept verbatim.
type TransferStockState int32

const (
	// TransferStockStateUnknown means upstream sent no usable state.
	TransferStockStateUnknown TransferStockState = -1
	// TransferStockStateProcessing means the security is being transferred.
	TransferStockStateProcessing TransferStockState = 0
	// TransferStockStateSucceeded means the security was transferred.
	TransferStockStateSucceeded TransferStockState = 1
	// TransferStockStateFailed means the security failed to transfer.
	TransferStockStateFailed TransferStockState = 2
)

// TransferStock is one security of a stock transfer record. Proof-document
// references are not modelled.
type TransferStock struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	NameEn    string `json:"name_en"`
	Market    string `json:"market"`
	CounterID string `json:"counter_id"`
	// Amount is the quantity, as returned upstream.
	Amount           string             `json:"amount"`
	Currency         string             `json:"currency"`
	Cost             string             `json:"cost"`
	StateCode        TransferStockState `json:"state_code"`
	StreetAmount     string             `json:"street_amount"`
	NomineeAmount    string             `json:"nominee_amount"`
	OwnAmount        string             `json:"own_amount"`
	CertificateCount string             `json:"certificate_count"`
	PositionSide     string             `json:"position_side"`
}

// AppStateNoticed is the StockTransferIn/Out.AppState value set when the user
// tapped "I have notified the broker" in the app (the user's own
// confirmation, not a Longbridge action). Other records carry "default".
const AppStateNoticed = "noticed"

// StockTransferIn is a stock transfer-in application (another broker → this
// account). The account holder's name and the contact details at the other
// broker, and proof-document references, are not modelled.
type StockTransferIn struct {
	ID       string `json:"id"`
	OrderID  string `json:"order_id"`
	Market   string `json:"market"`
	Broker   string `json:"broker"`
	BrokerID string `json:"broker_id"`
	BrokerEn string `json:"broker_en"`
	// BrokerCode is the other broker's code (upstream "code").
	BrokerCode string `json:"broker_code"`
	// Account is the full account number at the other broker.
	Account string `json:"account"`
	// State is the upstream state text; StateCode the typed code.
	State     string          `json:"state"`
	StateCode TransferInState `json:"state_code"`
	// AppState is AppStateNoticed or "default"; see UserNotifiedBroker.
	AppState string `json:"app_state"`
	// CanModify reports whether the application can still be cancelled or
	// changed.
	CanModify bool      `json:"can_modify"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ProcessAt time.Time `json:"process_at"`
	// ExecAt is when the transfer executed; zero while it has not.
	ExecAt time.Time       `json:"exec_at"`
	Remark string          `json:"remark"`
	Reason string          `json:"reason"`
	Stocks []TransferStock `json:"stocks"`
}

// UserNotifiedBroker reports whether the user confirmed in the app that the
// other broker was notified.
func (r *StockTransferIn) UserNotifiedBroker() bool { return r.AppState == AppStateNoticed }

// StockTransferInList is one page of stock transfer-in records.
type StockTransferInList struct {
	Records []*StockTransferIn `json:"records"`
	Total   int64              `json:"total"`
}

// StockTransferOut is a stock transfer-out application (this account →
// another broker). The account holder's name and contact details at the other
// broker, and proof-document references, are not modelled.
type StockTransferOut struct {
	ID       string `json:"id"`
	OrderID  string `json:"order_id"`
	Market   string `json:"market"`
	Broker   string `json:"broker"`
	BrokerID string `json:"broker_id"`
	// BrokerEn is the other broker's English name (upstream "name_en").
	BrokerEn string `json:"broker_en"`
	// BrokerCode is the other broker's code (upstream "code").
	BrokerCode string `json:"broker_code"`
	// Account is the full account number at the other broker.
	Account string `json:"account"`
	// State is the upstream state text; StateCode the typed code.
	State     string           `json:"state"`
	StateCode TransferOutState `json:"state_code"`
	// AppState is AppStateNoticed or "default"; see UserNotifiedBroker.
	AppState  string    `json:"app_state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ProcessAt time.Time `json:"process_at"`
	// ExecAt is when the transfer executed; zero while it has not.
	ExecAt   time.Time       `json:"exec_at"`
	AgreeAt  time.Time       `json:"agree_at"`
	RefuseAt time.Time       `json:"refuse_at"`
	CancelAt time.Time       `json:"cancel_at"`
	Remark   string          `json:"remark"`
	Reason   string          `json:"reason"`
	Stocks   []TransferStock `json:"stocks"`
}

// UserNotifiedBroker reports whether the user confirmed in the app that the
// other broker was notified.
func (r *StockTransferOut) UserNotifiedBroker() bool { return r.AppState == AppStateNoticed }

// Cancellable reports whether the application can still be cancelled. Upstream
// can_modify is 0 even for a cancellable transfer-out record, so only the
// submitted state counts.
func (r *StockTransferOut) Cancellable() bool { return r.StateCode == TransferOutStateSubmitted }

// StockTransferOutList is one page of stock transfer-out records. Upstream
// returns no total: a page shorter than the requested size is the last one.
type StockTransferOutList struct {
	Records []*StockTransferOut `json:"records"`
}

// GetStockTransfers contains the filters of StockTransfersIn and
// StockTransfersOut. All fields are optional.
type GetStockTransfers struct {
	// Market filters by market, e.g. "HK" or "US".
	Market string
	// State filters by state code: a TransferInState for StockTransfersIn or a
	// TransferOutState for StockTransfersOut. Nil means all states.
	State *int32
	// Start and End bound the application time (inclusive; zero means
	// unbounded) and are sent as unix seconds (start / end). Upstream filters on unix seconds, so only the instant
	// matters, not the time.Time's location. A user's calendar dates are Hong
	// Kong dates: pass 00:00:00 and 23:59:59 in UTC+8, e.g.
	// time.Date(2026, 11, 1, 0, 0, 0, 0, time.FixedZone("HKT", 8*3600)).
	Start time.Time
	End   time.Time
	// Page is the 1-based page number; 0 defaults to 1.
	Page int32
	// Size is the page size; 0 defaults to 20 and values above 50 are sent
	// as 50.
	Size int32
}

// RewardType is the type of a reward record.
// Undocumented integer codes are kept verbatim.
type RewardType int32

const (
	// RewardTypeUnknown means upstream sent no usable reward type and the
	// detail object did not identify one either.
	RewardTypeUnknown RewardType = -1
	// RewardTypeCash is a cash reward.
	RewardTypeCash RewardType = 0
	// RewardTypeGift is a gift (prize).
	RewardTypeGift RewardType = 1
	// RewardTypeCoupon is a coupon card.
	RewardTypeCoupon RewardType = 2
	// RewardTypeTask is a task reward.
	RewardTypeTask RewardType = 3
)

// RewardStatus is the status of a reward record.
// Undocumented integer codes are kept verbatim.
type RewardStatus int32

const (
	// RewardStatusUnknown means upstream sent no usable status.
	RewardStatusUnknown RewardStatus = -1
	// RewardStatusUpcoming means the reward is not issued yet.
	RewardStatusUpcoming RewardStatus = 0
	// RewardStatusCancelled means the reward was cancelled.
	RewardStatusCancelled RewardStatus = 1
	// RewardStatusExitedActivity means the user left the activity.
	RewardStatusExitedActivity RewardStatus = 2
	// RewardStatusUnclaimed means the reward has not been claimed.
	RewardStatusUnclaimed RewardStatus = 8
	// RewardStatusIssuing means the reward is being issued.
	RewardStatusIssuing RewardStatus = 10
	// RewardStatusIssued means the reward was issued.
	RewardStatusIssued RewardStatus = 12
	// RewardStatusFailed means issuing failed.
	RewardStatusFailed RewardStatus = 13
	// RewardStatusExchanged means the reward was exchanged.
	RewardStatusExchanged RewardStatus = 16
	// RewardStatusExpired means the reward expired.
	RewardStatusExpired RewardStatus = 17
	// RewardStatusRefunded means the reward was refunded.
	RewardStatusRefunded RewardStatus = 18
)

// RewardDeliverStatus is the delivery status of an issued prize. It is only
// meaningful when the record is a prize whose Status is RewardStatusIssued.
// Undocumented integer codes are kept verbatim.
type RewardDeliverStatus int32

const (
	// RewardDeliverStatusUnknown means upstream sent no usable delivery
	// status for a prize.
	RewardDeliverStatusUnknown RewardDeliverStatus = -1
	// RewardDeliverStatusNotActivated means the prize is not activated.
	RewardDeliverStatusNotActivated RewardDeliverStatus = 0
	// RewardDeliverStatusActivatedOrShipping means the prize is activated or
	// being shipped.
	RewardDeliverStatusActivatedOrShipping RewardDeliverStatus = 1
	// RewardDeliverStatusConsumed means the prize was consumed.
	RewardDeliverStatusConsumed RewardDeliverStatus = 2
	// RewardDeliverStatusEVoucherActivating means an e-voucher is activating.
	RewardDeliverStatusEVoucherActivating RewardDeliverStatus = 5
	// RewardDeliverStatusEVoucherExpired means an e-voucher expired before
	// activation.
	RewardDeliverStatusEVoucherExpired RewardDeliverStatus = 6
	// RewardDeliverStatusEVoucherInvalid means an e-voucher is invalid.
	RewardDeliverStatusEVoucherInvalid RewardDeliverStatus = 7
	// RewardDeliverStatusAddressMissing means no delivery address was given.
	RewardDeliverStatusAddressMissing RewardDeliverStatus = 8
	// RewardDeliverStatusAddressFilled means the delivery address was given.
	RewardDeliverStatusAddressFilled RewardDeliverStatus = 9
	// RewardDeliverStatusShipping means the prize is being shipped.
	RewardDeliverStatusShipping RewardDeliverStatus = 10
	// RewardDeliverStatusAddressFilledAlt is another code the app also shows
	// as "address filled".
	RewardDeliverStatusAddressFilledAlt RewardDeliverStatus = 11
)

// PrizeTypePhysical is the RewardRecord.PrizeType of physical goods (the only
// prize type with a delivery address and shipping states). Other prize types
// are kept verbatim.
const PrizeTypePhysical int32 = 0

// RewardDetailKind names which upstream detail object a reward record came
// from.
type RewardDetailKind string

const (
	// RewardDetailKindNone means the item carried no known detail object
	// (e.g. task rewards); only RewardType is set.
	RewardDetailKindNone RewardDetailKind = ""
	// RewardDetailKindCash is the cash{} object.
	RewardDetailKindCash RewardDetailKind = "cash"
	// RewardDetailKindCard is the card{} (coupon) object.
	RewardDetailKindCard RewardDetailKind = "card"
	// RewardDetailKindPrize is the prize{} (gift) object.
	RewardDetailKindPrize RewardDetailKind = "prize"
)

// RewardRecord is one reward record, flattened from the item's cash{},
// card{} or prize{} object. Coupon secrets, cash card codes, delivery details
// and addresses are not modelled.
type RewardRecord struct {
	// RewardType is upstream reward_type, or inferred from DetailKind when
	// upstream omits it.
	RewardType RewardType `json:"reward_type"`
	// DetailKind is the detail object the fields below came from.
	DetailKind    RewardDetailKind `json:"detail_kind"`
	RecordID      string           `json:"record_id"`
	ActivityTitle string           `json:"activity_title"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	// CardType is set for coupon rewards.
	CardType string `json:"card_type"`
	// CardStatus applies only to coupon rewards (DetailKind
	// RewardDetailKindCard). It is CouponStatusUnknown for every other reward
	// and when upstream sent no usable status (as in the Rust SDK).
	CardStatus CouponStatus `json:"card_status"`
	// PrizeType is set only for prizes (DetailKind RewardDetailKindPrize) that
	// carry an integer prize_type; PrizeTypePhysical is physical goods. nil for
	// every other reward, and for a prize without a usable prize_type.
	PrizeType *int32 `json:"prize_type,omitempty"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
	// Status is RewardStatusUnknown when upstream sent no usable status.
	Status RewardStatus `json:"status"`
	// DeliverStatus applies only to prizes, and only once the prize is
	// RewardStatusIssued. It is RewardDeliverStatusUnknown for every other
	// reward and when upstream sent no usable delivery status (as in the Rust
	// SDK).
	DeliverStatus RewardDeliverStatus `json:"deliver_status"`
	// ExitStatus is the upstream exit_status code, -1 when absent.
	ExitStatus      int32     `json:"exit_status"`
	SendTime        time.Time `json:"send_time"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	EffectiveTime   time.Time `json:"effective_time"`
	ReceiveDeadline time.Time `json:"receive_deadline"`
}

// RewardRecordList is one page of reward records.
type RewardRecordList struct {
	Records []*RewardRecord `json:"records"`
	Total   int64           `json:"total"`
}

// GetRewardRecords contains the parameters of RewardRecords.
type GetRewardRecords struct {
	// RewardType filters by type; nil means all types.
	RewardType *RewardType
	// Page is the 1-based page number; 0 defaults to 1.
	Page int32
	// PageSize is the page size; 0 defaults to 20 and values above 50 are
	// sent as 50.
	PageSize int32
}

// CouponStatus is the status of a coupon card.
// Undocumented integer codes are kept verbatim.
type CouponStatus int32

const (
	// CouponStatusUnknown means upstream sent no usable status.
	CouponStatusUnknown CouponStatus = -1
	// CouponStatusUnused means the coupon is unused.
	CouponStatusUnused CouponStatus = 0
	// CouponStatusUsed means the coupon was used.
	CouponStatusUsed CouponStatus = 1
	// CouponStatusExpired means the coupon expired.
	CouponStatusExpired CouponStatus = 2
	// CouponStatusRedeeming means the coupon is being redeemed.
	CouponStatusRedeeming CouponStatus = 3
	// CouponStatusInUse means the coupon is in use.
	CouponStatusInUse CouponStatus = 4
	// CouponStatusProcessing means the coupon is being processed.
	CouponStatusProcessing CouponStatus = 5
)

// CouponStatusFilter selects coupons by validity in GetCoupons.
type CouponStatusFilter int32

const (
	// CouponStatusFilterValid selects coupons that have not expired (default).
	CouponStatusFilterValid CouponStatusFilter = -1
	// CouponStatusFilterExpired selects expired coupons.
	CouponStatusFilterExpired CouponStatusFilter = 2
)

// Coupon card types accepted by GetCoupons.CardType (not exhaustive).
const (
	CouponCardTypeCommission        = "commission"
	CouponCardTypePlatformFeeDeduct = "platform_fee_deduct"
	CouponCardTypeStock             = "stock"
	CouponCardTypeQuote             = "quote"
	CouponCardTypeCash              = "cash"
)

// CouponTargetType is an instrument type a coupon applies to (the server's
// bz_types / underlying_type, sent as a proto enum name such as
// "HkStockWarrant" or as its number). Undocumented integer codes are kept
// verbatim.
type CouponTargetType int32

const (
	// CouponTargetTypeUnknown is an entry upstream sent that is neither an
	// integer nor a known name.
	CouponTargetTypeUnknown CouponTargetType = -1
	// CouponTargetTypeStocks is stocks (StockTrade).
	CouponTargetTypeStocks CouponTargetType = 0
	// CouponTargetTypeOptions is US stock options (UsStockOptions).
	CouponTargetTypeOptions CouponTargetType = 1
	// CouponTargetTypeWarrants is HK warrants (HkStockWarrant).
	CouponTargetTypeWarrants CouponTargetType = 2
	// CouponTargetTypeETF is ETFs.
	CouponTargetTypeETF CouponTargetType = 3
	// CouponTargetTypeLOF is LOFs.
	CouponTargetTypeLOF CouponTargetType = 4
)

// IpoCouponTargetType is an IPO subscription type an IPO commission coupon
// applies to (sent as a proto enum name or number). The proto's unset value
// 0 decodes to IpoCouponTargetTypeUnknown.
type IpoCouponTargetType int32

const (
	// IpoCouponTargetTypeUnknown is an unset or unrecognised entry.
	IpoCouponTargetTypeUnknown IpoCouponTargetType = -1
	// IpoCouponTargetTypeFinancing is a financing (margin) IPO subscription
	// (FinancingIPO).
	IpoCouponTargetTypeFinancing IpoCouponTargetType = 1
	// IpoCouponTargetTypeCash is a cash IPO subscription (CashIPO).
	IpoCouponTargetTypeCash IpoCouponTargetType = 2
)

// CommissionUnlimitedType is the server's `unlimited` code of a commission
// coupon (a proto enum sent as a name or a number). Undocumented integer
// codes are kept verbatim.
type CommissionUnlimitedType int32

const (
	// CommissionUnlimitedTypeUnknown means upstream sent no usable value.
	CommissionUnlimitedTypeUnknown CommissionUnlimitedType = -1
	// CommissionUnlimitedTypeProportionate is `Proportionate`.
	CommissionUnlimitedTypeProportionate CommissionUnlimitedType = 0
	// CommissionUnlimitedTypeFull is `Full`.
	CommissionUnlimitedTypeFull CommissionUnlimitedType = 1
	// CommissionUnlimitedTypeLifetime is `Lifetime`.
	CommissionUnlimitedTypeLifetime CommissionUnlimitedType = 2
	// CommissionUnlimitedTypeByAmountPer is `ByAmountPer`.
	CommissionUnlimitedTypeByAmountPer CommissionUnlimitedType = 6
	// CommissionUnlimitedTypeByRatePer is `ByRatePer`.
	CommissionUnlimitedTypeByRatePer CommissionUnlimitedType = 7
)

// QuoteInvestorType is the investor type a quote coupon applies to (a proto
// enum sent as a name or a number). Undocumented integer codes are kept
// verbatim.
type QuoteInvestorType int32

const (
	// QuoteInvestorTypeUnknown means upstream sent no usable value.
	QuoteInvestorTypeUnknown QuoteInvestorType = -1
	// QuoteInvestorTypeAny is any investor (INVESTOR_TYPE_ANY).
	QuoteInvestorTypeAny QuoteInvestorType = 0
	// QuoteInvestorTypeNonProfessional is non-professional investors
	// (INVESTOR_TYPE_NON_PROFESSIONAL).
	QuoteInvestorTypeNonProfessional QuoteInvestorType = 1
	// QuoteInvestorTypeProfessional is professional investors
	// (INVESTOR_TYPE_PROFESSIONAL).
	QuoteInvestorTypeProfessional QuoteInvestorType = 2
)

// CouponCommission is the commission terms of a commission coupon.
type CouponCommission struct {
	Markets    []string `json:"markets"`
	RebateRate string   `json:"rebate_rate"`
	// BzTypes are the instrument types the coupon applies to.
	BzTypes            []CouponTargetType `json:"bz_types"`
	CommissionBalance  string             `json:"commission_balance"`
	CommissionDuration int64              `json:"commission_duration"`
	// Unlimited is the server's `unlimited` code.
	Unlimited CommissionUnlimitedType `json:"unlimited"`
	// CommissionExpiredTime is zero when unset.
	CommissionExpiredTime time.Time `json:"commission_expired_time"`
	UsageCount            int64     `json:"usage_count"`
	LaveUsageCount        int64     `json:"lave_usage_count"`
	SettlementModel       string    `json:"settlement_model"`
	// AmountLimit is the amount limit, as returned by the server.
	AmountLimit int64 `json:"amount_limit"`
	// Superposition reports whether the coupon can be combined with others.
	Superposition bool `json:"superposition"`
}

// CouponAiBenefit is the terms of an AI benefit coupon.
type CouponAiBenefit struct {
	// BenefitDays is the benefit days, as returned by the server.
	BenefitDays int32 `json:"benefit_days"`
	// BenefitStart and BenefitEnd are zero when unset.
	BenefitStart time.Time `json:"benefit_start"`
	BenefitEnd   time.Time `json:"benefit_end"`
	BenefitTag   string    `json:"benefit_tag"`
}

// CouponCash is the terms of a cash coupon.
type CouponCash struct {
	// ReceiveDuration is the receive duration, as returned by the server.
	ReceiveDuration int64 `json:"receive_duration"`
}

// CouponCoffeeStore is a store of a coffee coupon.
type CouponCoffeeStore struct {
	StoreID   int64  `json:"store_id"`
	StoreName string `json:"store_name"`
}

// CouponCoffee is the terms of a coffee coupon. String fields are as
// returned by the server.
type CouponCoffee struct {
	CouponChannel       string `json:"coupon_channel"`
	CouponID            string `json:"coupon_id"`
	CouponNo            string `json:"coupon_no"`
	IsAvailableAllStore bool   `json:"is_available_all_store"`
	// PeriodType is e.g. "day".
	PeriodType string `json:"period_type"`
	// SettledStoreID is 0 when the coupon was not settled at a store.
	SettledStoreID int64               `json:"settled_store_id"`
	StoreIDList    []int64             `json:"store_id_list"`
	StoreList      []CouponCoffeeStore `json:"store_list"`
	Threshold      int64               `json:"threshold"`
	// UseLimit is "total" or "period".
	UseLimit string `json:"use_limit"`
}

// CouponCryptoCounter is a counter (security) of a crypto coupon.
type CouponCryptoCounter struct {
	CounterID string `json:"counter_id"`
}

// CouponCryptoFallback is the fallback of a crypto coupon.
type CouponCryptoFallback struct {
	Enabled   bool   `json:"enabled"`
	CounterID string `json:"counter_id"`
}

// CouponCrypto is the terms of a crypto coupon.
type CouponCrypto struct {
	Fallback         CouponCryptoFallback  `json:"fallback"`
	OptionalCounters []CouponCryptoCounter `json:"optional_counters"`
}

// CouponFinanceInterest is the terms of a finance interest coupon.
type CouponFinanceInterest struct {
	// DiscountMode is the discount mode, as returned by the server.
	DiscountMode int32 `json:"discount_mode"`
	// DiscountRate is the discount rate as decimal text.
	DiscountRate  string `json:"discount_rate"`
	Superposition bool   `json:"superposition"`
}

// CouponFundDeduct is the terms of a fund deduction coupon. String fields
// are as returned by the server.
type CouponFundDeduct struct {
	MaxRiskLevel     int32    `json:"max_risk_level"`
	MinOrderAmount   string   `json:"min_order_amount"`
	OrderTypes       []string `json:"order_types"`
	UnderlyingType   string   `json:"underlying_type"`
	UnderlyingValues []string `json:"underlying_values"`
	UnderlyingNames  []string `json:"underlying_names"`
}

// CouponQuote is the terms of a quote (market data) or fundamental-data
// coupon.
type CouponQuote struct {
	// CurrentEndTime is zero when unset.
	CurrentEndTime time.Time         `json:"current_end_time"`
	Inapplicable   bool              `json:"inapplicable"`
	InvestorType   QuoteInvestorType `json:"investor_type"`
	ManualActivate bool              `json:"manual_activate"`
	Markets        []string          `json:"markets"`
	// MemberInvestorType is the member's investor type, as returned by the
	// server.
	MemberInvestorType int32  `json:"member_investor_type"`
	PairSpuNo          string `json:"pair_spu_no"`
	// QuoteDuration is the quote duration, as returned by the server.
	QuoteDuration int64    `json:"quote_duration"`
	QuoteMarket   string   `json:"quote_market"`
	QuoteTags     []string `json:"quote_tags"`
	Superposition bool     `json:"superposition"`
}

// CouponIpo is the terms of an IPO coupon.
type CouponIpo struct {
	SubscriptionID string `json:"subscription_id"`
}

// CouponIpoCommission is the terms of an IPO commission coupon.
type CouponIpoCommission struct {
	BzTypes     []IpoCouponTargetType `json:"bz_types"`
	CounterID   string                `json:"counter_id"`
	CounterName string                `json:"counter_name"`
	// FinancingLimit is the financing limit, as returned by the server.
	FinancingLimit string `json:"financing_limit"`
}

// CouponIpoInterest is the terms of an IPO interest coupon.
type CouponIpoInterest struct {
	// DiscountRate is the discount rate, as returned by the server.
	DiscountRate string `json:"discount_rate"`
}

// CouponPlatformFeeDeduct is the terms of a platform fee deduction coupon.
type CouponPlatformFeeDeduct struct {
	BzTypes        []CouponTargetType `json:"bz_types"`
	CounterID      string             `json:"counter_id"`
	CounterName    string             `json:"counter_name"`
	Markets        []string           `json:"markets"`
	UsageCount     int64              `json:"usage_count"`
	LaveUsageCount int64              `json:"lave_usage_count"`
	Superposition  bool               `json:"superposition"`
}

// CouponStockCounter is an optional counter of a stock coupon.
type CouponStockCounter struct {
	CounterID string `json:"counter_id"`
	// Shares is as returned by the server.
	Shares string `json:"shares"`
}

// CouponStockFallback is the fallback of a stock coupon. String fields are
// as returned by the server.
type CouponStockFallback struct {
	Enabled bool `json:"enabled"`
	// Type is the server's `type`.
	Type      string `json:"type"`
	CounterID string `json:"counter_id"`
	Shares    string `json:"shares"`
}

// CouponStock is the terms of a stock coupon.
type CouponStock struct {
	CounterID   string `json:"counter_id"`
	CounterName string `json:"counter_name"`
	// Amount is decimal text.
	Amount string `json:"amount"`
	// RedemptionType is "by_card_value" or "shares_number".
	RedemptionType string `json:"redemption_type"`
	// SendWay is "card_exchange" or "account_transfer".
	SendWay          string               `json:"send_way"`
	OptionalCounters []CouponStockCounter `json:"optional_counters"`
	Fallback         CouponStockFallback  `json:"fallback"`
}

// CouponStockFragment is the terms of a stock fragment coupon.
type CouponStockFragment struct {
	CounterID      string           `json:"counter_id"`
	CounterName    string           `json:"counter_name"`
	Market         string           `json:"market"`
	Markets        []string         `json:"markets"`
	UnderlyingType CouponTargetType `json:"underlying_type"`
	// Amounts is as returned by the server.
	Amounts map[string]string `json:"amounts"`
}

// Coupon is a coupon card. The type-specific terms are in the pointer field
// named after the card type (Commission, Quote, Stock, …); the fields of the
// other types are nil. The fund_interest object (an arbitrary object in the
// proto), deprecated top-level duplicates of the commission terms and app
// styling are not modelled.
type Coupon struct {
	// CardID is the 32-hex coupon id (use it with CouponOrders).
	CardID      string       `json:"card_id"`
	CardType    string       `json:"card_type"`
	Title       string       `json:"title"`
	TitleSec    string       `json:"title_sec"`
	Description string       `json:"description"`
	Detail      string       `json:"detail"`
	CardStatus  CouponStatus `json:"card_status"`
	StartTime   time.Time    `json:"start_time"`
	// EndTime is when the coupon expires; zero when unset or when NoExpiry is
	// true. Upstream sends seconds or milliseconds; both are handled.
	EndTime time.Time `json:"end_time"`
	// NoExpiry reports a coupon that never expires (upstream end time past
	// 2100).
	NoExpiry     bool      `json:"no_expiry"`
	ActivateTime time.Time `json:"activate_time"`
	Amount       string    `json:"amount"`
	Currency     string    `json:"currency"`
	// AvailableTradeSession is the upstream display text of the sessions the
	// coupon applies to.
	AvailableTradeSession string   `json:"available_trade_session"`
	TradeSessions         []string `json:"trade_sessions"`
	// SubscriptionID is as returned by the server.
	SubscriptionID    string                   `json:"subscription_id,omitempty"`
	AiBenefit         *CouponAiBenefit         `json:"ai_benefit,omitempty"`
	Cash              *CouponCash              `json:"cash,omitempty"`
	CoffeeCoupon      *CouponCoffee            `json:"coffee_coupon,omitempty"`
	Commission        *CouponCommission        `json:"commission,omitempty"`
	Crypto            *CouponCrypto            `json:"crypto,omitempty"`
	FinanceInterest   *CouponFinanceInterest   `json:"finance_interest,omitempty"`
	FundDeduct        *CouponFundDeduct        `json:"fund_deduct,omitempty"`
	Fundamental       *CouponQuote             `json:"fundamental,omitempty"`
	Ipo               *CouponIpo               `json:"ipo,omitempty"`
	IpoCommission     *CouponIpoCommission     `json:"ipo_commission,omitempty"`
	IpoInterest       *CouponIpoInterest       `json:"ipo_interest,omitempty"`
	PlatformFeeDeduct *CouponPlatformFeeDeduct `json:"platform_fee_deduct,omitempty"`
	Quote             *CouponQuote             `json:"quote,omitempty"`
	Stock             *CouponStock             `json:"stock,omitempty"`
	StockFragment     *CouponStockFragment     `json:"stock_fragment,omitempty"`
}

// CouponList is one page of coupon cards. Upstream splits the list into pinned
// sections: TopList and BottomList are repeated on every page.
type CouponList struct {
	TopList      []*Coupon `json:"top_list"`
	List         []*Coupon `json:"list"`
	BottomList   []*Coupon `json:"bottom_list"`
	Total        int64     `json:"total"`
	ExpiredCount int64     `json:"expired_count"`
}

// GetCoupons contains the parameters of Coupons.
type GetCoupons struct {
	// Status selects valid or expired coupons; 0 defaults to
	// CouponStatusFilterValid.
	Status CouponStatusFilter
	// CardType filters by type, e.g. CouponCardTypeCommission; empty means all.
	CardType string
	// Page is the 1-based page number; 0 defaults to 1.
	Page int32
	// PageSize is the page size; 0 defaults to 20 and values above 50 are
	// sent as 50.
	PageSize int32
}

// CouponOrderType selects which orders CouponOrders returns.
type CouponOrderType string

const (
	// CouponOrderTypeLocked is orders the coupon is locked to (in progress).
	CouponOrderTypeLocked CouponOrderType = "locked"
	// CouponOrderTypeSettled is orders whose rebate was settled.
	CouponOrderTypeSettled CouponOrderType = "settled"
)

// CouponOrder is an order a coupon card was applied to.
type CouponOrder struct {
	OrderID string `json:"order_id"`
	// State is the upstream display text; StateVal the status name (e.g.
	// "FilledStatus").
	State              string `json:"state"`
	StateVal           string `json:"state_val"`
	ExecutedPrice      string `json:"executed_price"`
	ExecutedQuantity   string `json:"executed_quantity"`
	Currency           string `json:"currency"`
	SettlementCurrency string `json:"settlement_currency"`
	Rebate             string `json:"rebate"`
	// RebateDt is the rebate date, as returned upstream.
	RebateDt string `json:"rebate_dt"`
}

// DepositState is the state_code of a deposit.
// Undocumented integer codes are kept verbatim.
type DepositState int32

const (
	// DepositStateUnknown means upstream sent no usable state code.
	DepositStateUnknown DepositState = -1
	// DepositStateProcessing means the deposit is being processed.
	DepositStateProcessing DepositState = 0
	// DepositStateCredited means the deposit was credited.
	DepositStateCredited DepositState = 1
	// DepositStateFailed means the deposit failed.
	DepositStateFailed DepositState = 2
	// DepositStateCancelled means the deposit was cancelled.
	DepositStateCancelled DepositState = 3
	// DepositStateInvalidated means the deposit was invalidated.
	DepositStateInvalidated DepositState = 4
	// DepositStateRejected means the deposit was rejected.
	DepositStateRejected DepositState = 5
)

// Deposit is a deposit record. Receipt images are not modelled.
type Deposit struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Currency  string    `json:"currency"`
	// Amount is the notified (applied) amount.
	Amount string `json:"amount"`
	// ReceivableAmount is the amount received ("0" when not received yet).
	ReceivableAmount string `json:"receivable_amount"`
	// State is the upstream state text; StateCode the typed code.
	State     string       `json:"state"`
	StateCode DepositState `json:"state_code"`
	// SubStateCode is the upstream sub_state_code, -1 when absent.
	SubStateCode          int32  `json:"sub_state_code"`
	BankOperationTypeID   string `json:"bank_operation_type_id"`
	BankOperationTypeName string `json:"bank_operation_type_name"`
	ReferenceNo           string `json:"reference_no"`
	CanCancel             bool   `json:"can_cancel"`
	UserBankSwiftCode     string `json:"user_bank_swift_code"`
}

// DepositList is one page of deposits. Upstream also returns a total, but it
// follows the page size rather than the record count, so it is not exposed.
type DepositList struct {
	Items []*Deposit `json:"items"`
	// HasMore reports that a full page came back, so the next page may hold
	// more records.
	HasMore bool `json:"has_more"`
}

// GetDeposits contains the parameters of Deposits. All fields are optional.
type GetDeposits struct {
	// States filters by state; empty means all.
	States []DepositState
	// Currencies filters by currency, e.g. "HKD"; empty means all.
	Currencies []string
	// Start and End bound the deposit time (inclusive; zero means unbounded)
	// and are sent as unix seconds (from_date / to_date). Upstream filters on unix seconds, so only the instant
	// matters, not the time.Time's location. A user's calendar dates are Hong
	// Kong dates: pass 00:00:00 and 23:59:59 in UTC+8, e.g.
	// time.Date(2026, 11, 1, 0, 0, 0, 0, time.FixedZone("HKT", 8*3600)).
	Start time.Time
	End   time.Time
	// Page is the 1-based page number; 0 defaults to 1.
	Page int32
	// Size is the page size; 0 defaults to 20 and values above 50 are sent
	// as 50.
	Size int32
}

// WithdrawalState is the state of a withdrawal.
// Undocumented integer codes are kept verbatim.
type WithdrawalState int32

const (
	// WithdrawalStateUnknown means upstream sent no usable state.
	WithdrawalStateUnknown WithdrawalState = -1
	// WithdrawalStateSubmitted means the withdrawal was submitted.
	WithdrawalStateSubmitted WithdrawalState = 0
	// WithdrawalStateProcessing means the withdrawal is being processed.
	WithdrawalStateProcessing WithdrawalState = 1
	// WithdrawalStateRejected means the withdrawal was rejected.
	WithdrawalStateRejected WithdrawalState = 2
	// WithdrawalStateCancelled means the withdrawal was cancelled.
	WithdrawalStateCancelled WithdrawalState = 3
	// WithdrawalStateRemitted means the money was paid out.
	WithdrawalStateRemitted WithdrawalState = 4
	// WithdrawalStateFailed means the withdrawal failed.
	WithdrawalStateFailed WithdrawalState = 5
	// WithdrawalStateRefunded means the withdrawal was refunded.
	WithdrawalStateRefunded WithdrawalState = 6
	// WithdrawalStateApprovedAwaitingPayout means approved, awaiting payout.
	WithdrawalStateApprovedAwaitingPayout WithdrawalState = 7
	// WithdrawalStateLargeAmountReview means under large-amount review.
	WithdrawalStateLargeAmountReview WithdrawalState = 8
	// WithdrawalStateChangePendingApproval means a change awaits approval.
	WithdrawalStateChangePendingApproval WithdrawalState = 9
	// WithdrawalStateDeletionPendingApproval means a deletion awaits approval.
	WithdrawalStateDeletionPendingApproval WithdrawalState = 10
	// WithdrawalStateDeleted means the withdrawal was deleted.
	WithdrawalStateDeleted WithdrawalState = 11
	// WithdrawalStateRejectionPendingApproval means a rejection awaits
	// approval.
	WithdrawalStateRejectionPendingApproval WithdrawalState = 12
	// WithdrawalStateAwaitingConfirmation means awaiting confirmation.
	WithdrawalStateAwaitingConfirmation WithdrawalState = 15
)

// Withdrawal is a withdrawal record.
type Withdrawal struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
	Currency  string    `json:"currency"`
	// Amount is the applied amount.
	Amount string `json:"amount"`
	Fee    string `json:"fee"`
	// SettlementAmount is the amount paid out.
	SettlementAmount string `json:"settlement_amount"`
	// StateCode is the typed state (upstream "state"); StateDescription is
	// the display text.
	StateCode        WithdrawalState `json:"state_code"`
	RawState         string          `json:"raw_state"`
	StateDescription string          `json:"state_description"`
	// BankAccount is the full receiving account number, as returned upstream.
	BankAccount        string `json:"bank_account"`
	DisplayBankAccount string `json:"display_bank_account"`
	UserBankSwiftCode  string `json:"user_bank_swift_code"`
	AccountChannel     string `json:"account_channel"`
}

// WithdrawalList is one page of withdrawals.
type WithdrawalList struct {
	List  []*Withdrawal `json:"list"`
	Total int64         `json:"total"`
	// HasMore reports that a full page came back, so the next page may hold
	// more records.
	HasMore bool `json:"has_more"`
}

// GetWithdrawals contains the parameters of Withdrawals. All fields are
// optional.
type GetWithdrawals struct {
	// Start and End bound the withdrawal time (inclusive; zero means
	// unbounded) and are sent as unix seconds (begin / end). Upstream filters on unix seconds, so only the instant
	// matters, not the time.Time's location. A user's calendar dates are Hong
	// Kong dates: pass 00:00:00 and 23:59:59 in UTC+8, e.g.
	// time.Date(2026, 11, 1, 0, 0, 0, 0, time.FixedZone("HKT", 8*3600)).
	Start time.Time
	End   time.Time
	// Page is the 1-based page number; 0 defaults to 1.
	Page int32
	// Size is the page size; 0 defaults to 20 and values above 50 are sent
	// as 50.
	Size int32
}
