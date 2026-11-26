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
// every code type has a matching XxxUnknown constant.
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
	if m := c.Commission; m != nil {
		commission := &CouponCommission{
			Markets:               []string(m.Markets),
			RebateRate:            string(m.RebateRate),
			BzTypes:               make([]CouponTargetType, 0, len(m.BzTypes)),
			CommissionBalance:     string(m.CommissionBalance),
			CommissionDuration:    string(m.CommissionDuration),
			Unlimited:             bool(m.Unlimited),
			CommissionExpiredTime: parseTime(m.CommissionExpiredTime),
			UsageCount:            int64(m.UsageCount),
			LaveUsageCount:        int64(m.LaveUsageCount),
			SettlementModel:       string(m.SettlementModel),
		}
		for _, t := range m.BzTypes {
			commission.BzTypes = append(commission.BzTypes, CouponTargetType(code(t)))
		}
		coupon.Commission = commission
	}
	return coupon
}

func toCoupons(list []jsontypes.Coupon) []*Coupon {
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
