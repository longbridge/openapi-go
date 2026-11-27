package asset

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/longbridge/openapi-go/asset/jsontypes"
)

// Placeholder bounds for upstream timestamps: anything before minValidUnix
// (2000-01-01, which catches 0, negative values and the 1906 / 1970
// placeholders) or from maxValidUnix (2100-01-01, which catches the
// 9999-12-31 "long-term" sentinel) on is treated as unset.
const (
	minValidUnix = 946684800
	maxValidUnix = 4102444800
	// msThreshold separates unix milliseconds from unix seconds.
	msThreshold = 1_000_000_000_000
)

// rawUnix parses an upstream timestamp (unix seconds or milliseconds, as a
// number or a string, or an RFC 3339 string) into unix seconds. ok is false
// when the value is empty or unparsable.
func rawUnix(s jsontypes.FlexString) (sec int64, ok bool) {
	t := strings.TrimSpace(string(s))
	if t == "" {
		return 0, false
	}
	if n, err := strconv.ParseInt(t, 10, 64); err == nil {
		if n > msThreshold {
			n /= 1000
		}
		return n, true
	}
	if parsed, err := time.Parse(time.RFC3339, t); err == nil {
		return parsed.Unix(), true
	}
	return 0, false
}

// parseTime converts an upstream timestamp to UTC time.Time; placeholders and
// unparsable values give the zero time.
func parseTime(s jsontypes.FlexString) time.Time {
	sec, ok := rawUnix(s)
	if !ok || sec < minValidUnix || sec >= maxValidUnix {
		return time.Time{}
	}
	return time.Unix(sec, 0).UTC()
}

// unknownCode is what a missing, null or non-integer status code decodes to;
// every code type has a matching XxxUnknown constant. Any integer upstream
// sends, documented or not, is kept as is.
const unknownCode int32 = -1

// code returns the decoded status code, or unknownCode when upstream sent no
// usable value.
func code(c jsontypes.FlexCode) int32 {
	if !c.Valid {
		return unknownCode
	}
	return c.Value
}

// rawObject returns v, or nil when it is empty or JSON null.
func rawObject(v json.RawMessage) json.RawMessage {
	t := bytes.TrimSpace(v)
	if len(t) == 0 || bytes.Equal(t, []byte("null")) {
		return nil
	}
	return v
}

func toBankCard(c *jsontypes.BankCard) *BankCard {
	card := &BankCard{
		ID:             string(c.ID),
		BankID:         string(c.BankID),
		Bank:           string(c.Bank),
		BankEn:         string(c.BankEn),
		BankCode:       string(c.BankCode),
		Account:        string(c.Account),
		AccountType:    string(c.AccountType),
		Country:        string(c.Country),
		CountryDisplay: string(c.CountryDisplay),
		Region:         string(c.Region),
		RegionName:     string(c.RegionName),
		SwiftCode:      string(c.SwiftCode),
		Remark:         string(c.Remark),
		Status:         BankCardStatus(code(c.Status)),
		WitnessState:   BankCardWitnessState(code(c.WitnessState)),
		SupportWitness: bool(c.SupportWitness),
		BindState:      BankCardBindState(code(c.BindState)),
		BindTip:        string(c.BindTip),
		CreatedAt:      parseTime(c.CreatedAt),
		WithdrawFee:    string(c.WithdrawFee),
		AuthList:       make([]BankCardAuth, 0, len(c.AuthList)),
	}
	for i := range c.AuthList {
		a := &c.AuthList[i]
		maxAmount := string(a.MaxAmount)
		if maxAmount == "0" {
			maxAmount = ""
		}
		auth := BankCardAuth{
			FlowID:                  string(a.FlowID),
			Name:                    string(a.Name),
			Currency:                string(a.Currency),
			Status:                  AuthorizationStatus(code(a.Status)),
			Message:                 string(a.Message),
			IsSupport:               bool(a.IsSupport),
			MaxAmount:               maxAmount,
			MaxAmountUnlimited:      bool(a.MaxAmountUnlimited),
			DailyMaxCount:           int64(a.DailyMaxCount),
			DailyMaxCountUnlimited:  bool(a.DailyMaxCountUnlimited),
			LongTermValid:           bool(a.LongTermValid),
			DocReviewing:            bool(a.DocReviewing),
			NoPriorDepositReviewing: bool(a.NoPriorDepositReviewing),
		}
		if !auth.LongTermValid {
			auth.ExpiryDate = parseTime(a.ExpiryDate)
		}
		card.AuthList = append(card.AuthList, auth)
	}
	return card
}

// decodeEddaApplications accepts both reply shapes of
// GET /v1/account/edda-status: {list[]} and, from older builds, a single
// {id, status, message} object.
func decodeEddaApplications(data json.RawMessage) ([]*EddaApplication, error) {
	out := []*EddaApplication{}
	if rawObject(data) == nil {
		return out, nil
	}
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}
	var apps []jsontypes.EddaApplication
	if list, ok := probe["list"]; ok {
		if rawObject(list) != nil {
			if err := json.Unmarshal(list, &apps); err != nil {
				return nil, err
			}
		}
	} else if _, ok := probe["status"]; ok {
		var app jsontypes.EddaApplication
		if err := json.Unmarshal(data, &app); err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	for i := range apps {
		a := &apps[i]
		out = append(out, &EddaApplication{
			ID:        string(a.ID),
			FlowID:    string(a.FlowID),
			Name:      string(a.Name),
			Currency:  string(a.Currency),
			Status:    AuthorizationStatus(code(a.Status)),
			Message:   string(a.Message),
			CreatedAt: parseTime(a.CreatedAt),
			UpdatedAt: parseTime(a.UpdatedAt),
		})
	}
	return out, nil
}

func toTransferStocks(stocks []jsontypes.TransferStock) []TransferStock {
	out := make([]TransferStock, 0, len(stocks))
	for i := range stocks {
		s := &stocks[i]
		out = append(out, TransferStock{
			Code:             string(s.Code),
			Name:             string(s.Name),
			NameEn:           string(s.NameEn),
			Market:           string(s.Market),
			CounterID:        string(s.CounterID),
			Amount:           string(s.Amount),
			Currency:         string(s.Currency),
			Cost:             string(s.Cost),
			StateCode:        TransferStockState(code(s.State)),
			StreetAmount:     string(s.StreetAmount),
			NomineeAmount:    string(s.NomineeAmount),
			OwnAmount:        string(s.OwnAmount),
			CertificateCount: string(s.CertificateCount),
			PositionSide:     string(s.PositionSide),
		})
	}
	return out
}

// brokerEn returns the broker's English name: transfer-out records name it
// name_en instead of broker_en.
func brokerEn(r *jsontypes.TransferRecord) string {
	if r.BrokerEn != "" {
		return string(r.BrokerEn)
	}
	return string(r.NameEn)
}

func toStockTransferIn(r *jsontypes.TransferRecord) *StockTransferIn {
	return &StockTransferIn{
		ID:         string(r.ID),
		OrderID:    string(r.OrderID),
		Market:     string(r.Market),
		Broker:     string(r.Broker),
		BrokerID:   string(r.BrokerID),
		BrokerEn:   brokerEn(r),
		BrokerCode: string(r.Code),
		Account:    string(r.Account),
		State:      string(r.State),
		StateCode:  TransferInState(code(r.StateCode)),
		AppState:   string(r.AppState),
		CanModify:  bool(r.CanModify),
		CreatedAt:  parseTime(r.CreatedAt),
		UpdatedAt:  parseTime(r.UpdatedAt),
		ProcessAt:  parseTime(r.ProcessAt),
		ExecAt:     parseTime(r.ExecAt),
		Remark:     string(r.Remark),
		Reason:     string(r.Reason),
		Stocks:     toTransferStocks(r.Stocks),
	}
}

func toStockTransferOut(r *jsontypes.TransferRecord) *StockTransferOut {
	return &StockTransferOut{
		ID:         string(r.ID),
		OrderID:    string(r.OrderID),
		Market:     string(r.Market),
		Broker:     string(r.Broker),
		BrokerID:   string(r.BrokerID),
		BrokerEn:   brokerEn(r),
		BrokerCode: string(r.Code),
		Account:    string(r.Account),
		State:      string(r.State),
		StateCode:  TransferOutState(code(r.StateCode)),
		AppState:   string(r.AppState),
		CreatedAt:  parseTime(r.CreatedAt),
		UpdatedAt:  parseTime(r.UpdatedAt),
		ProcessAt:  parseTime(r.ProcessAt),
		ExecAt:     parseTime(r.ExecAt),
		AgreeAt:    parseTime(r.AgreeAt),
		RefuseAt:   parseTime(r.RefuseAt),
		CancelAt:   parseTime(r.CancelAt),
		Remark:     string(r.Remark),
		Reason:     string(r.Reason),
		Stocks:     toTransferStocks(r.Stocks),
	}
}

// toRewardRecord flattens one reward item. It returns nil for an item with
// neither a reward_type nor a known detail object.
func toRewardRecord(item *jsontypes.RewardItem) *RewardRecord {
	var (
		kind     RewardDetailKind
		detail   *jsontypes.RewardDetail
		inferred RewardType
	)
	switch {
	case item.Cash != nil:
		kind, detail, inferred = RewardDetailKindCash, item.Cash, RewardTypeCash
	case item.Card != nil:
		kind, detail, inferred = RewardDetailKindCard, item.Card, RewardTypeCoupon
	case item.Prize != nil:
		kind, detail, inferred = RewardDetailKindPrize, item.Prize, RewardTypeGift
	}
	if detail == nil && !item.RewardType.Valid {
		return nil
	}
	rec := &RewardRecord{
		RewardType:    inferred,
		DetailKind:    kind,
		Status:        RewardStatusUnknown,
		CardStatus:    CouponStatusUnknown,
		DeliverStatus: RewardDeliverStatusUnknown,
		ExitStatus:    -1,
	}
	if item.RewardType.Valid {
		rec.RewardType = RewardType(item.RewardType.Value)
	}
	if detail == nil {
		return rec
	}
	rec.RecordID = string(detail.RecordID)
	rec.ActivityTitle = string(detail.ActivityTitle)
	rec.Title = string(detail.Title)
	rec.Description = string(detail.Description)
	rec.CardType = string(detail.CardType)
	rec.Amount = string(detail.Amount)
	rec.Currency = string(detail.Currency)
	rec.Status = RewardStatus(code(detail.Status))
	rec.ExitStatus = code(detail.ExitStatus)
	switch kind {
	case RewardDetailKindCard:
		rec.CardStatus = CouponStatus(code(detail.CardStatus))
	case RewardDetailKindPrize:
		if detail.PrizeType.Valid {
			prizeType := detail.PrizeType.Value
			rec.PrizeType = &prizeType
		}
		rec.DeliverStatus = RewardDeliverStatus(code(detail.DeliverStatus))
	}
	rec.SendTime = parseTime(detail.SendTime)
	rec.StartTime = parseTime(detail.StartTime)
	rec.EndTime = parseTime(detail.EndTime)
	rec.EffectiveTime = parseTime(detail.EffectiveTime)
	rec.ReceiveDeadline = parseTime(detail.ReceiveDeadline)
	return rec
}

func toCoupon(c *jsontypes.Coupon) *Coupon {
	coupon := &Coupon{
		CardID:                string(c.CardID),
		CardType:              string(c.CardType),
		Title:                 string(c.Title),
		TitleSec:              string(c.TitleSec),
		Description:           string(c.Description),
		Detail:                string(c.Detail),
		CardStatus:            CouponStatus(code(c.CardStatus)),
		StartTime:             parseTime(c.StartTime),
		EndTime:               parseTime(c.EndTime),
		ActivateTime:          parseTime(c.ActivateTime),
		Amount:                string(c.Amount),
		Currency:              string(c.Currency),
		AvailableTradeSession: string(c.AvailableTradeSession),
		TradeSessions:         []string(c.TradeSessions),
	}
	if sec, ok := rawUnix(c.EndTime); ok && sec >= maxValidUnix {
		coupon.NoExpiry = true
	}
	coupon.SubscriptionID = string(c.SubscriptionID)
	if m := c.Commission.V; m != nil {
		coupon.Commission = &CouponCommission{
			Markets:               []string(m.Markets),
			RebateRate:            string(m.RebateRate),
			BzTypes:               targetTypes(m.BzTypes),
			CommissionBalance:     string(m.CommissionBalance),
			CommissionDuration:    int64(m.CommissionDuration),
			Unlimited:             CommissionUnlimitedType(enumCode(m.Unlimited, unlimitedNames)),
			CommissionExpiredTime: parseTime(m.CommissionExpiredTime),
			UsageCount:            int64(m.UsageCount),
			LaveUsageCount:        int64(m.LaveUsageCount),
			SettlementModel:       string(m.SettlementModel),
			AmountLimit:           int64(m.AmountLimit),
			Superposition:         bool(m.Superposition),
		}
	}
	if m := c.AiBenefit.V; m != nil {
		coupon.AiBenefit = &CouponAiBenefit{
			BenefitDays:  int32(m.BenefitDays),
			BenefitStart: parseTime(m.BenefitStart),
			BenefitEnd:   parseTime(m.BenefitEnd),
			BenefitTag:   string(m.BenefitTag),
		}
	}
	if m := c.Cash.V; m != nil {
		coupon.Cash = &CouponCash{ReceiveDuration: int64(m.ReceiveDuration)}
	}
	if m := c.CoffeeCoupon.V; m != nil {
		coffee := &CouponCoffee{
			CouponChannel:       string(m.CouponChannel),
			CouponID:            string(m.CouponID),
			CouponNo:            string(m.CouponNo),
			IsAvailableAllStore: bool(m.IsAvailableAllStore),
			PeriodType:          string(m.PeriodType),
			SettledStoreID:      int64(m.SettledStoreID),
			StoreIDList:         []int64(m.StoreIDList),
			StoreList:           make([]CouponCoffeeStore, 0, len(m.StoreList)),
			Threshold:           int64(m.Threshold),
			UseLimit:            string(m.UseLimit),
		}
		for _, st := range m.StoreList {
			coffee.StoreList = append(coffee.StoreList, CouponCoffeeStore{
				StoreID:   int64(st.StoreID),
				StoreName: string(st.StoreName),
			})
		}
		coupon.CoffeeCoupon = coffee
	}
	if m := c.Crypto.V; m != nil {
		crypto := &CouponCrypto{OptionalCounters: make([]CouponCryptoCounter, 0, len(m.OptionalCounters))}
		if f := m.Fallback.V; f != nil {
			crypto.Fallback = CouponCryptoFallback{Enabled: bool(f.Enabled), CounterID: string(f.CounterID)}
		}
		for _, o := range m.OptionalCounters {
			crypto.OptionalCounters = append(crypto.OptionalCounters, CouponCryptoCounter{CounterID: string(o.CounterID)})
		}
		coupon.Crypto = crypto
	}
	if m := c.FinanceInterest.V; m != nil {
		coupon.FinanceInterest = &CouponFinanceInterest{
			DiscountMode:  int32(m.DiscountMode),
			DiscountRate:  string(m.DiscountRate),
			Superposition: bool(m.Superposition),
		}
	}
	if m := c.FundDeduct.V; m != nil {
		coupon.FundDeduct = &CouponFundDeduct{
			MaxRiskLevel:     int32(m.MaxRiskLevel),
			MinOrderAmount:   string(m.MinOrderAmount),
			OrderTypes:       []string(m.OrderTypes),
			UnderlyingType:   string(m.UnderlyingType),
			UnderlyingValues: []string(m.UnderlyingValues),
			UnderlyingNames:  []string(m.UnderlyingNames),
		}
	}
	coupon.Fundamental = toCouponQuote(c.Fundamental.V)
	coupon.Quote = toCouponQuote(c.Quote.V)
	if m := c.Ipo.V; m != nil {
		coupon.Ipo = &CouponIpo{SubscriptionID: string(m.SubscriptionID)}
	}
	if m := c.IpoCommission.V; m != nil {
		ipo := &CouponIpoCommission{
			BzTypes:        make([]IpoCouponTargetType, 0, len(m.BzTypes)),
			CounterID:      string(m.CounterID),
			CounterName:    string(m.CounterName),
			FinancingLimit: string(m.FinancingLimit),
		}
		for _, t := range m.BzTypes {
			v, ok := namedCode(t, ipoTargetNames)
			if !ok {
				continue
			}
			if v == 0 {
				// The proto's unset value ("_").
				v = unknownCode
			}
			ipo.BzTypes = append(ipo.BzTypes, IpoCouponTargetType(v))
		}
		coupon.IpoCommission = ipo
	}
	if m := c.IpoInterest.V; m != nil {
		coupon.IpoInterest = &CouponIpoInterest{DiscountRate: string(m.DiscountRate)}
	}
	if m := c.PlatformFeeDeduct.V; m != nil {
		coupon.PlatformFeeDeduct = &CouponPlatformFeeDeduct{
			BzTypes:        targetTypes(m.BzTypes),
			CounterID:      string(m.CounterID),
			CounterName:    string(m.CounterName),
			Markets:        []string(m.Markets),
			UsageCount:     int64(m.UsageCount),
			LaveUsageCount: int64(m.LaveUsageCount),
			Superposition:  bool(m.Superposition),
		}
	}
	if m := c.Stock.V; m != nil {
		stock := &CouponStock{
			CounterID:        string(m.CounterID),
			CounterName:      string(m.CounterName),
			Amount:           string(m.Amount),
			RedemptionType:   string(m.RedemptionType),
			SendWay:          string(m.SendWay),
			OptionalCounters: make([]CouponStockCounter, 0, len(m.OptionalCounters)),
		}
		for _, o := range m.OptionalCounters {
			stock.OptionalCounters = append(stock.OptionalCounters, CouponStockCounter{
				CounterID: string(o.CounterID),
				Shares:    string(o.Shares),
			})
		}
		if f := m.Fallback.V; f != nil {
			stock.Fallback = CouponStockFallback{
				Enabled:   bool(f.Enabled),
				Type:      string(f.Type),
				CounterID: string(f.CounterID),
				Shares:    string(f.Shares),
			}
		}
		coupon.Stock = stock
	}
	if m := c.StockFragment.V; m != nil {
		coupon.StockFragment = &CouponStockFragment{
			CounterID:      string(m.CounterID),
			CounterName:    string(m.CounterName),
			Market:         string(m.Market),
			Markets:        []string(m.Markets),
			UnderlyingType: CouponTargetType(enumCode(m.UnderlyingType, targetTypeNames)),
			Amounts:        map[string]string(m.Amounts),
		}
	}
	return coupon
}

// Proto enum names the coupon terms may carry instead of the number.
var (
	targetTypeNames = map[string]int32{
		"StockTrade": 0, "UsStockOptions": 1, "HkStockWarrant": 2, "ETF": 3, "LOF": 4,
	}
	ipoTargetNames = map[string]int32{"_": 0, "FinancingIPO": 1, "CashIPO": 2}
	unlimitedNames = map[string]int32{
		"Proportionate": 0, "Full": 1, "Lifetime": 2, "ByAmountPer": 6, "ByRatePer": 7,
	}
	investorTypeNames = map[string]int32{
		"INVESTOR_TYPE_ANY":              0,
		"INVESTOR_TYPE_NON_PROFESSIONAL": 1,
		"INVESTOR_TYPE_PROFESSIONAL":     2,
	}
)

// namedCode decodes a proto enum sent as its name or as its number (a JSON
// number, a numeric string or an integral float). A number is kept as is even
// when it is not a documented value; ok is false for an empty or
// unrecognised value.
func namedCode(s string, names map[string]int32) (int32, bool) {
	t := strings.TrimSpace(s)
	if v, ok := names[t]; ok {
		return v, true
	}
	if v, ok := jsontypes.ParseIntegral(t, 32); ok {
		return int32(v), true
	}
	return 0, false
}

// enumCode is namedCode for a single enum field: a missing or unrecognised
// value is unknownCode.
func enumCode(s jsontypes.FlexString, names map[string]int32) int32 {
	if v, ok := namedCode(string(s), names); ok {
		return v
	}
	return unknownCode
}

// targetTypes decodes a list of CouponTargetType entries, skipping entries
// that are neither a known name nor an integer.
func targetTypes(list []string) []CouponTargetType {
	out := make([]CouponTargetType, 0, len(list))
	for _, t := range list {
		if v, ok := namedCode(t, targetTypeNames); ok {
			out = append(out, CouponTargetType(v))
		}
	}
	return out
}

func toCouponQuote(m *jsontypes.CouponQuote) *CouponQuote {
	if m == nil {
		return nil
	}
	return &CouponQuote{
		CurrentEndTime:     parseTime(m.CurrentEndTime),
		Inapplicable:       bool(m.Inapplicable),
		InvestorType:       QuoteInvestorType(enumCode(m.InvestorType, investorTypeNames)),
		ManualActivate:     bool(m.ManualActivate),
		Markets:            []string(m.Markets),
		MemberInvestorType: int32(m.MemberInvestorType),
		PairSpuNo:          string(m.PairSpuNo),
		QuoteDuration:      int64(m.QuoteDuration),
		QuoteMarket:        string(m.QuoteMarket),
		QuoteTags:          []string(m.QuoteTags),
		Superposition:      bool(m.Superposition),
	}
}

func toCoupons(list jsontypes.List[jsontypes.Coupon]) []*Coupon {
	out := make([]*Coupon, 0, len(list))
	for i := range list {
		out = append(out, toCoupon(&list[i]))
	}
	return out
}

func toCouponOrder(o *jsontypes.CouponOrder) *CouponOrder {
	return &CouponOrder{
		OrderID:            string(o.OrderID),
		State:              string(o.State),
		StateVal:           string(o.StateVal),
		ExecutedPrice:      string(o.ExecutedPrice),
		ExecutedQuantity:   string(o.ExecutedQuantity),
		Currency:           string(o.Currency),
		SettlementCurrency: string(o.SettlementCurrency),
		Rebate:             string(o.Rebate),
		RebateDt:           string(o.RebateDt),
	}
}

func toDeposit(d *jsontypes.Deposit) *Deposit {
	return &Deposit{
		ID:                    string(d.ID),
		CreatedAt:             parseTime(d.CreatedAt),
		UpdatedAt:             parseTime(d.UpdatedAt),
		Currency:              string(d.Currency),
		Amount:                string(d.Amount),
		ReceivableAmount:      string(d.ReceivableAmount),
		State:                 string(d.State),
		StateCode:             DepositState(code(d.StateCode)),
		SubStateCode:          code(d.SubStateCode),
		BankOperationTypeID:   string(d.BankOperationTypeID),
		BankOperationTypeName: string(d.BankOperationTypeName),
		ReferenceNo:           string(d.ReferenceNo),
		CanCancel:             bool(d.CanCancel),
		UserBankSwiftCode:     string(d.UserBankSwiftCode),
		VoucherCount:          int64(d.Vouchers),
	}
}

func toWithdrawal(w *jsontypes.Withdrawal) *Withdrawal {
	return &Withdrawal{
		ID:                 string(w.ID),
		OrderID:            string(w.OrderID),
		CreatedAt:          parseTime(w.CreatedAt),
		Currency:           string(w.Currency),
		Amount:             string(w.Amount),
		Fee:                string(w.Fee),
		SettlementAmount:   string(w.SettlementAmount),
		StateCode:          WithdrawalState(code(w.State)),
		RawState:           string(w.RawState),
		StateDescription:   string(w.StateDescription),
		BankAccount:        string(w.BankAccount),
		DisplayBankAccount: string(w.DisplayBankAccount),
		UserBankSwiftCode:  string(w.UserBankSwiftCode),
		AccountChannel:     string(w.AccountChannel),
	}
}
