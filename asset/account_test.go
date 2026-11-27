package asset

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/longbridge/openapi-go/asset/jsontypes"
	"github.com/longbridge/openapi-go/http"
)

func newTestAssetContext(t *testing.T, srvURL, prefix string) *AssetContext {
	t.Helper()
	cli, err := http.New(
		http.WithURL(srvURL),
		http.WithAppKey(prefix+"key"),
		http.WithAppSecret(prefix+"secret"),
		http.WithAccessToken(prefix+"token"),
	)
	if err != nil {
		t.Fatalf("http.New: %v", err)
	}
	return &AssetContext{httpClient: cli}
}

// serve answers every request with data and records the last request.
func serve(t *testing.T, data string, got **nethttp.Request) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		*got = r
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"","data":` + data + `}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestParseTimePlaceholders(t *testing.T) {
	tests := []struct {
		in   string
		want int64 // 0 = zero time
	}{
		{"", 0},
		{"0", 0},
		{"-2000000000", 0},
		{"1970-01-01T00:00:00Z", 0},
		{"1906-06-07T18:35:44Z", 0},
		{"253402214400", 0}, // 9999-12-31 sentinel
		{"8061747303", 0},   // past 2100
		{"garbage", 0},
		{"1792600000", 1792600000},
		{"1798732800000", 1798732800}, // milliseconds
		{"2026-11-25T09:59:06Z", 1795600746},
		{"2026-11-25T17:59:06+08:00", 1795600746},
	}
	for _, tt := range tests {
		got := parseTime(jsontypes.FlexString(tt.in))
		if tt.want == 0 {
			if !got.IsZero() {
				t.Errorf("parseTime(%q) = %v, want zero", tt.in, got)
			}
			continue
		}
		if got.Unix() != tt.want || got.Location() != time.UTC {
			t.Errorf("parseTime(%q) = %v, want unix %d UTC", tt.in, got, tt.want)
		}
	}
}

func TestBankCards(t *testing.T) {
	var req *nethttp.Request
	srv := serve(t, `{"show_add_card_button": true, "list": [{
		"id": "1024", "bank_id": 12, "bank": "HSBC", "bank_en": "HSBC", "bank_code": "004",
		"account": "012345678901", "name": "CHAN TAI MAN", "name_en": "CHAN TAI MAN",
		"account_type": 1, "country": "Hong Kong", "country_display": "Hong Kong",
		"region": "HK", "region_name": "Hong Kong", "address": "1 Queen's Road",
		"swift_code": "HSBCHKHHHKH", "status": 1, "witness_state": 3, "support_witness": true,
		"bind_state": 1, "bind_tip": "", "created_at": "1792600000",
		"picture_urls": ["id/front.jpg"], "payee_street": "x", "payee_city": "y",
		"withdraw_fee": "0", "logo": "https://x/logo.png",
		"auth_list": [
			{"flow_id": "hk_edda", "name": "eDDA", "currency": "HKD", "status": 3,
			 "message": "Name mismatch", "is_support": 1, "max_amount": "0",
			 "max_amount_unlimited": false, "daily_max_count": "5",
			 "expiry_date": 253402214400, "long_term_valid": true, "custom": "789",
			 "id": "0", "created_at": "0", "need_sms": 0},
			{"flow_id": "hk_fps", "currency": "USD", "status": 1, "expiry_date": "1830000000",
			 "long_term_valid": false, "max_amount": "50000"}
		]
	}]}`, &req)
	actx := newTestAssetContext(t, srv.URL, "")
	cards, err := actx.BankCards(context.Background())
	if err != nil {
		t.Fatalf("BankCards: %v", err)
	}
	if req.URL.Path != "/v2/account/bank-cards" || req.URL.Query().Get("from_app_page") != "cards" {
		t.Errorf("request = %s", req.URL)
	}
	if len(cards) != 1 {
		t.Fatalf("got %d cards", len(cards))
	}
	c := cards[0]
	if c.ID != "1024" || c.BankID != "12" || c.Account != "012345678901" || c.Country != "Hong Kong" ||
		c.AccountType != "1" || c.Status != BankCardStatusUnderReview ||
		c.WitnessState != BankCardWitnessStateWitnessed || c.BindState != BankCardBindStateBound ||
		!c.SupportWitness || c.CreatedAt.Unix() != 1792600000 {
		t.Errorf("card = %+v", c)
	}
	if len(c.AuthList) != 2 {
		t.Fatalf("auth list = %+v", c.AuthList)
	}
	edda := c.AuthList[0]
	if !edda.IsEDDA() || edda.Status != AuthorizationStatusFailed || edda.Message != "Name mismatch" ||
		!edda.IsSupport || edda.MaxAmount != "" || edda.DailyMaxCount != 5 ||
		!edda.LongTermValid || !edda.ExpiryDate.IsZero() {
		t.Errorf("edda auth = %+v", edda)
	}
	fps := c.AuthList[1]
	if fps.IsEDDA() || fps.Status != AuthorizationStatusAuthorized || fps.ExpiryDate.Unix() != 1830000000 || fps.MaxAmount != "50000" {
		t.Errorf("fps auth = %+v", fps)
	}
	// Sensitive upstream fields must not leak into the public type.
	b, _ := json.Marshal(c)
	for _, leaked := range []string{"CHAN TAI MAN", "Queen", "front.jpg", "payee", "logo"} {
		if strings.Contains(string(b), leaked) {
			t.Errorf("marshalled card contains %q: %s", leaked, b)
		}
	}
}

func TestEddaApplicationsShapes(t *testing.T) {
	var req *nethttp.Request
	for _, tt := range []struct {
		name string
		data string
		want int
	}{
		{"list", `{"list": [{"id": 77, "flow_id": "hk_sdda", "currency": "HKD", "status": 2, "created_at": 1792600000}, {"id": "78", "status": 3, "message": "Rejected by bank"}]}`, 2},
		{"legacy object", `{"id": "77", "status": 1, "message": ""}`, 1},
		{"empty", `{}`, 0},
		{"null", `null`, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			srv := serve(t, tt.data, &req)
			apps, err := newTestAssetContext(t, srv.URL, "").EddaApplications(context.Background(), " 1024 ", " hkd ")
			if err != nil {
				t.Fatalf("EddaApplications: %v", err)
			}
			if q := req.URL.Query(); req.URL.Path != "/v1/account/edda-status" || q.Get("card_id") != "1024" || q.Get("currency") != "HKD" {
				t.Errorf("request = %s", req.URL)
			}
			if apps == nil || len(apps) != tt.want {
				t.Fatalf("got %+v, want %d applications", apps, tt.want)
			}
			if tt.name == "list" {
				if apps[0].ID != "77" || apps[0].Status != AuthorizationStatusAuthorizing || apps[0].CreatedAt.Unix() != 1792600000 ||
					apps[1].Status != AuthorizationStatusFailed || apps[1].Message != "Rejected by bank" {
					t.Errorf("apps = %+v %+v", apps[0], apps[1])
				}
			}
			if tt.name == "legacy object" && (apps[0].ID != "77" || apps[0].Status != AuthorizationStatusAuthorized) {
				t.Errorf("app = %+v", apps[0])
			}
		})
	}
	actx := newTestAssetContext(t, "http://127.0.0.1:0", "")
	for _, id := range []string{"", " ", "abc", "12a", "-1", "1.5", "0123456789abcdef0123456789abcdef"} {
		if _, err := actx.EddaApplications(context.Background(), id, "HKD"); err == nil {
			t.Errorf("card id %q must be rejected", id)
		}
	}
	if _, err := actx.EddaApplications(context.Background(), "1", ""); err == nil {
		t.Error("empty currency must be rejected")
	}
}

func TestStockTransfersIn(t *testing.T) {
	var req *nethttp.Request
	// Shape of a real canary transfer-in record, 2026-11-25.
	srv := serve(t, `{"total": 1, "records": [{
		"id": "7053", "order_id": "TI202611250001", "market": "US", "market_symbol": "US",
		"broker": "富途证券", "broker_id": "9", "broker_en": "Futu Securities", "code": "FUTU",
		"account": "11222333444555", "name": "CHAN TAI MAN", "contact": "Mr Lee",
		"contact_phone": "12345678", "contact_email": "a@b.c", "user_email": "me@x.y",
		"state": "已提交", "state_code": 0, "app_state": "noticed", "can_modify": 1,
		"created_at": 1795600746, "updated_at": "1795600746", "process_at": "0",
		"exec_at": "1970-01-01T00:00:00Z", "cancel_at": -1990000000,
		"attachments": "[\"ai/proof.pdf\"]", "cash": [],
		"stocks": [{"code": "AAPL", "name": "苹果", "market": "US", "counter_id": "ST/US/AAPL",
			"amount": 100, "currency": "USD", "cost": "180.5", "state": 0,
			"usi_attachments": ["k"], "is_reality": 0, "street_amount": "0"}]
	}]}`, &req)
	state := int32(TransferInStateSubmitted)
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.FixedZone("HKT", 8*3600))
	list, err := newTestAssetContext(t, srv.URL, "ap_").StockTransfersIn(context.Background(), &GetStockTransfers{
		Market: "us", State: &state, Start: start,
	})
	if err != nil {
		t.Fatalf("StockTransfersIn: %v", err)
	}
	q := req.URL.Query()
	want := url.Values{"market": {"US"}, "state": {"0"}, "start": {"1767196800"}, "page": {"1"}, "size": {"20"}}
	if req.URL.Path != "/v1/account/stock-transfers/in" || q.Encode() != want.Encode() {
		t.Errorf("request = %s", req.URL)
	}
	if list.Total != 1 || len(list.Records) != 1 {
		t.Fatalf("list = %+v", list)
	}
	r := list.Records[0]
	if r.ID != "7053" || r.BrokerEn != "Futu Securities" || r.BrokerCode != "FUTU" || r.Account != "11222333444555" ||
		r.StateCode != TransferInStateSubmitted || !r.CanModify || !r.UserNotifiedBroker() ||
		r.CreatedAt.Unix() != 1795600746 || r.UpdatedAt.Unix() != 1795600746 {
		t.Errorf("record = %+v", r)
	}
	if !r.ProcessAt.IsZero() || !r.ExecAt.IsZero() {
		t.Errorf("placeholder times must be zero: process %v exec %v", r.ProcessAt, r.ExecAt)
	}
	if len(r.Stocks) != 1 || r.Stocks[0].Name != "苹果" || r.Stocks[0].Amount != "100" || r.Stocks[0].StateCode != TransferStockStateProcessing {
		t.Errorf("stocks = %+v", r.Stocks)
	}
	b, _ := json.Marshal(r)
	for _, leaked := range []string{"CHAN TAI MAN", "Mr Lee", "12345678\"", "a@b.c", "me@x.y", "proof.pdf", "usi_attachments"} {
		if strings.Contains(string(b), leaked) {
			t.Errorf("marshalled record contains %q: %s", leaked, b)
		}
	}
}

func TestStockTransfersOut(t *testing.T) {
	var req *nethttp.Request
	srv := serve(t, `{"records": [{
		"id": 88, "name_en": "Futu Securities", "broker": "富途证券", "state_code": 0, "can_modify": 0,
		"app_state": "default", "created_at": "1795600746", "agree_at": "1906-06-07T18:35:44Z",
		"refuse_at": "0", "cancel_at": "1970-01-01T00:00:00Z", "stocks": null
	}, {"id": "89", "state_code": 7, "refuse_at": 1795700000}]}`, &req)
	end := time.Date(2026, 1, 31, 23, 59, 59, 0, time.FixedZone("HKT", 8*3600))
	list, err := newTestAssetContext(t, srv.URL, "").StockTransfersOut(context.Background(), &GetStockTransfers{End: end, Page: 2, Size: 50})
	if err != nil {
		t.Fatalf("StockTransfersOut: %v", err)
	}
	q := req.URL.Query()
	if req.URL.Path != "/v1/account/stock-transfers/out" || q.Get("end") != "1769875199" || q.Get("page") != "2" || q.Get("size") != "50" || q.Has("start") || q.Has("state") {
		t.Errorf("request = %s", req.URL)
	}
	if len(list.Records) != 2 {
		t.Fatalf("list = %+v", list)
	}
	r := list.Records[0]
	if r.ID != "88" || r.BrokerEn != "Futu Securities" || !r.Cancellable() || r.UserNotifiedBroker() ||
		!r.AgreeAt.IsZero() || !r.RefuseAt.IsZero() || !r.CancelAt.IsZero() || r.Stocks == nil {
		t.Errorf("record 0 = %+v", r)
	}
	r = list.Records[1]
	if r.StateCode != TransferOutStateRejected || r.Cancellable() || r.RefuseAt.Unix() != 1795700000 {
		t.Errorf("record 1 = %+v", r)
	}
}

func TestRewardRecords(t *testing.T) {
	var req *nethttp.Request
	srv := serve(t, `{"total": 4, "list": [
		{"reward_type": 1, "prize": {"record_id": "r1", "activity_title": "Open account", "title": "Coffee",
			"status": 12, "deliver_status": 8, "prize_type": 0, "send_time": 1790000000,
			"receive_deadline": "0", "coupon_secret": {"secret": "ABC-123"}, "coupon_secret_id": "s1",
			"cash_card_code": "CC-9", "express_info": {"phone": "12345678", "address": "Somewhere"},
			"addresses": [{"detail": "Somewhere"}]}},
		{"reward_type": 0, "cash": {"record_id": "r2", "status": 12, "amount": 38, "currency": "HKD", "send_time": "0"}},
		{"card": {"record_id": "r3", "card_type": "commission", "card_status": 0, "status": 12}},
		{"reward_type": 3},
		{"unknown": {}}
	]}`, &req)
	rt := RewardTypeGift
	list, err := newTestAssetContext(t, srv.URL, "").RewardRecords(context.Background(), &GetRewardRecords{RewardType: &rt, PageSize: 50})
	if err != nil {
		t.Fatalf("RewardRecords: %v", err)
	}
	q := req.URL.Query()
	if req.URL.Path != "/v1/account/rewards" || q.Get("current") != "1" || q.Get("page_size") != "50" || q.Get("reward_type") != "1" {
		t.Errorf("request = %s", req.URL)
	}
	if list.Total != 4 || len(list.Records) != 4 {
		t.Fatalf("list = %+v", list)
	}
	prize := list.Records[0]
	if prize.RewardType != RewardTypeGift || prize.DetailKind != RewardDetailKindPrize || prize.Status != RewardStatusIssued ||
		prize.DeliverStatus != RewardDeliverStatusAddressMissing ||
		prize.PrizeType == nil || *prize.PrizeType != PrizeTypePhysical || prize.CardStatus != CouponStatusUnknown ||
		prize.SendTime.Unix() != 1790000000 || !prize.ReceiveDeadline.IsZero() {
		t.Errorf("prize = %+v", prize)
	}
	cash := list.Records[1]
	if cash.RewardType != RewardTypeCash || cash.Amount != "38" || !cash.SendTime.IsZero() ||
		cash.PrizeType != nil || cash.DeliverStatus != RewardDeliverStatusUnknown || cash.CardStatus != CouponStatusUnknown {
		t.Errorf("cash = %+v", cash)
	}
	if coupon := list.Records[2]; coupon.RewardType != RewardTypeCoupon || coupon.DetailKind != RewardDetailKindCard || coupon.CardType != "commission" ||
		coupon.CardStatus != CouponStatusUnused || coupon.PrizeType != nil || coupon.DeliverStatus != RewardDeliverStatusUnknown {
		t.Errorf("coupon = %+v", coupon)
	}
	if task := list.Records[3]; task.RewardType != RewardTypeTask || task.DetailKind != RewardDetailKindNone ||
		task.PrizeType != nil || task.DeliverStatus != RewardDeliverStatusUnknown || task.CardStatus != CouponStatusUnknown {
		t.Errorf("task = %+v", task)
	}
	b, _ := json.Marshal(list)
	for _, leaked := range []string{"ABC-123", "s1", "CC-9", "12345678", "Somewhere"} {
		if strings.Contains(string(b), leaked) {
			t.Errorf("marshalled rewards contain %q", leaked)
		}
	}
}

func TestCoupons(t *testing.T) {
	var req *nethttp.Request
	// Wire shape of production coupon cards (synthetic values): every
	// per-type object key is present, null for the types that do not apply.
	srv := serve(t, `{"total": 3, "expired_count": 5, "tags": [{"name": "all"}],
		"top_list": [{"card_id": "0123456789abcdef0123456789abcdef", "card_type": "commission",
			"start_time": 1750313704, "end_time": "8061747303", "card_status": 4, "activate_time": "0",
			"commission": {"markets": ["US"], "bz_types": [0, 2, 3], "commission_expired_time": "0",
				"rebate_rate": 1, "usage_count": "3", "lave_usage_count": 7, "unlimited": 2,
				"commission_duration": "0", "superposition": false},
			"markets": ["US"], "rebate_rate": "1", "quote_duration": "0",
			"available_trade_session": "夜盘、盘中+盘前盘后", "cover": "https://x/c.png",
			"coffee_coupon": null, "quote": null, "stock": null, "cash": null}],
		"list": [{"card_id": "ffffffffffffffffffffffffffffffff", "card_type": "coffee_coupon", "end_time": 1798732800000,
			"card_status": 0, "commission": null, "quote": null, "stock": null,
			"coffee_coupon": {"coupon_channel": "chan", "coupon_id": "", "coupon_no": "NO-1",
				"is_available_all_store": false, "period_type": "day", "settled_store_id": "0",
				"store_id_list": ["101", 102], "store_list": [{"store_id": "101", "store_name": "Store A"}],
				"threshold": "1", "use_limit": "period"}}],
		"bottom_list": [{"card_id": "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "end_time": "0", "card_status": 2}]
	}`, &req)
	list, err := newTestAssetContext(t, srv.URL, "").Coupons(context.Background(), &GetCoupons{CardType: "commission"})
	if err != nil {
		t.Fatalf("Coupons: %v", err)
	}
	q := req.URL.Query()
	if req.URL.Path != "/v1/account/cards" || q.Get("card_status") != "-1" || q.Get("card_type") != "commission" || q.Get("current") != "1" || q.Get("page_size") != "20" {
		t.Errorf("request = %s", req.URL)
	}
	if list.Total != 3 || list.ExpiredCount != 5 || len(list.TopList) != 1 || len(list.List) != 1 || len(list.BottomList) != 1 {
		t.Fatalf("list = %+v", list)
	}
	top := list.TopList[0]
	if !top.NoExpiry || !top.EndTime.IsZero() || !top.ActivateTime.IsZero() || top.StartTime.Unix() != 1750313704 ||
		top.CardStatus != CouponStatusInUse || top.Commission == nil {
		t.Fatalf("top = %+v", top)
	}
	cm := top.Commission
	if len(cm.BzTypes) != 3 || cm.BzTypes[1] != CouponTargetTypeWarrants || cm.Markets[0] != "US" || cm.RebateRate != "1" ||
		!cm.CommissionExpiredTime.IsZero() || cm.UsageCount != 3 || cm.LaveUsageCount != 7 ||
		cm.Unlimited != CommissionUnlimitedTypeLifetime || top.CoffeeCoupon != nil || top.Quote != nil || top.Cash != nil {
		t.Errorf("commission = %+v", cm)
	}
	coffee := list.List[0]
	if coffee.NoExpiry || coffee.EndTime.Unix() != 1798732800 || coffee.Commission != nil || coffee.CoffeeCoupon == nil {
		t.Fatalf("coffee coupon = %+v", coffee)
	}
	if cc := coffee.CoffeeCoupon; cc.CouponChannel != "chan" || cc.CouponNo != "NO-1" || cc.IsAvailableAllStore ||
		cc.PeriodType != "day" || cc.SettledStoreID != 0 || len(cc.StoreIDList) != 2 || cc.StoreIDList[1] != 102 ||
		len(cc.StoreList) != 1 || cc.StoreList[0].StoreID != 101 || cc.StoreList[0].StoreName != "Store A" ||
		cc.Threshold != 1 || cc.UseLimit != "period" {
		t.Errorf("coffee terms = %+v", cc)
	}
	if expired := list.BottomList[0]; expired.NoExpiry || !expired.EndTime.IsZero() || expired.CardStatus != CouponStatusExpired {
		t.Errorf("expired = %+v", expired)
	}
	b, _ := json.Marshal(list)
	for _, dropped := range []string{"c.png", `"cover"`, `"quote":`, `"stock":`, "platform_fee_deduct", "quote_duration"} {
		if strings.Contains(string(b), dropped) {
			t.Errorf("marshalled coupons contain %q: %s", dropped, b)
		}
	}

	if _, err = newTestAssetContext(t, srv.URL, "").Coupons(context.Background(), &GetCoupons{Status: CouponStatusFilterExpired, Page: 2}); err != nil {
		t.Fatalf("Coupons: %v", err)
	}
	if q := req.URL.Query(); q.Get("card_status") != "2" || q.Get("current") != "2" || q.Has("card_type") {
		t.Errorf("request = %s", req.URL)
	}
}

func TestCouponTypeTerms(t *testing.T) {
	var req *nethttp.Request
	// Synthetic values in the server's wire format: int64 as strings, enums
	// as names or numbers, "0" for unset times.
	srv := serve(t, `{"list": [{"card_id": "c", "subscription_id": "SUB",
		"ai_benefit": {"benefit_days": 30, "benefit_start": "1767196800", "benefit_end": "0", "benefit_tag": "tag"},
		"cash": {"receive_duration": "7"},
		"commission": {"bz_types": ["StockTrade", "HkStockWarrant", 9], "unlimited": "ByRatePer",
			"amount_limit": "500", "superposition": true, "commission_duration": "30"},
		"crypto": {"fallback": {"enabled": true, "counter_id": "C1"}, "optional_counters": [{"counter_id": "C2"}]},
		"finance_interest": {"discount_mode": 1, "discount_rate": 0.5, "superposition": true},
		"fund_deduct": {"max_risk_level": "3", "min_order_amount": "100", "order_types": ["buy"],
			"underlying_type": "fund", "underlying_values": ["F1"], "underlying_names": ["Fund"]},
		"fund_interest": {},
		"fundamental": {"investor_type": "INVESTOR_TYPE_PROFESSIONAL", "markets": ["US"]},
		"ipo": {"subscription_id": "S1"},
		"ipo_commission": {"bz_types": [1, "CashIPO", 0, "_"], "counter_id": "C3", "counter_name": "IPO", "financing_limit": "10"},
		"ipo_interest": {"discount_rate": "0.8"},
		"platform_fee_deduct": {"bz_types": ["UsStockOptions", "4"], "markets": ["US"], "usage_count": "1",
			"lave_usage_count": "4", "superposition": true, "counter_id": "C6", "counter_name": "P"},
		"quote": {"current_end_time": "1767196800", "investor_type": 1, "member_investor_type": "2",
			"quote_duration": "30", "quote_tags": ["L2"], "manual_activate": true, "inapplicable": false,
			"pair_spu_no": "SPU", "quote_market": "US", "superposition": false},
		"stock": {"amount": 12.5, "counter_id": "C4", "counter_name": "Stock", "redemption_type": "shares_number",
			"send_way": "account_transfer", "optional_counters": [{"counter_id": "C5", "shares": "3"}],
			"fallback": {"enabled": true, "type": "cash", "counter_id": "C7", "shares": "1"}},
		"stock_fragment": {"amounts": {"a": "1", "b": 2}, "counter_id": "C8", "counter_name": "F", "market": "US",
			"markets": ["US"], "underlying_type": "ETF"}}]}`, &req)
	list, err := newTestAssetContext(t, srv.URL, "").Coupons(context.Background(), nil)
	if err != nil {
		t.Fatalf("Coupons: %v", err)
	}
	c := list.List[0]
	if c.SubscriptionID != "SUB" {
		t.Errorf("subscription id = %q", c.SubscriptionID)
	}
	if a := c.AiBenefit; a == nil || a.BenefitDays != 30 || a.BenefitStart.Unix() != 1767196800 || !a.BenefitEnd.IsZero() || a.BenefitTag != "tag" {
		t.Errorf("ai benefit = %+v", a)
	}
	if c.Cash == nil || c.Cash.ReceiveDuration != 7 {
		t.Errorf("cash = %+v", c.Cash)
	}
	if m := c.Commission; m == nil || len(m.BzTypes) != 3 || m.BzTypes[0] != CouponTargetTypeStocks ||
		m.BzTypes[1] != CouponTargetTypeWarrants || m.BzTypes[2] != 9 || m.Unlimited != CommissionUnlimitedTypeByRatePer ||
		m.AmountLimit != 500 || !m.Superposition || m.CommissionDuration != 30 {
		t.Errorf("commission = %+v", m)
	}
	if m := c.Crypto; m == nil || !m.Fallback.Enabled || m.Fallback.CounterID != "C1" || len(m.OptionalCounters) != 1 || m.OptionalCounters[0].CounterID != "C2" {
		t.Errorf("crypto = %+v", m)
	}
	if m := c.FinanceInterest; m == nil || m.DiscountMode != 1 || m.DiscountRate != "0.5" || !m.Superposition {
		t.Errorf("finance interest = %+v", m)
	}
	if m := c.FundDeduct; m == nil || m.MaxRiskLevel != 3 || m.MinOrderAmount != "100" || m.OrderTypes[0] != "buy" ||
		m.UnderlyingType != "fund" || m.UnderlyingValues[0] != "F1" || m.UnderlyingNames[0] != "Fund" {
		t.Errorf("fund deduct = %+v", m)
	}
	if m := c.Fundamental; m == nil || m.InvestorType != QuoteInvestorTypeProfessional || m.Markets[0] != "US" {
		t.Errorf("fundamental = %+v", m)
	}
	if c.Ipo == nil || c.Ipo.SubscriptionID != "S1" {
		t.Errorf("ipo = %+v", c.Ipo)
	}
	if m := c.IpoCommission; m == nil || len(m.BzTypes) != 4 || m.BzTypes[0] != IpoCouponTargetTypeFinancing ||
		m.BzTypes[1] != IpoCouponTargetTypeCash || m.BzTypes[2] != IpoCouponTargetTypeUnknown ||
		m.BzTypes[3] != IpoCouponTargetTypeUnknown || m.CounterID != "C3" || m.FinancingLimit != "10" {
		t.Errorf("ipo commission = %+v", m)
	}
	if c.IpoInterest == nil || c.IpoInterest.DiscountRate != "0.8" {
		t.Errorf("ipo interest = %+v", c.IpoInterest)
	}
	if m := c.PlatformFeeDeduct; m == nil || len(m.BzTypes) != 2 || m.BzTypes[0] != CouponTargetTypeOptions ||
		m.BzTypes[1] != CouponTargetTypeLOF || m.UsageCount != 1 || m.LaveUsageCount != 4 || !m.Superposition || m.CounterID != "C6" {
		t.Errorf("platform fee deduct = %+v", m)
	}
	if m := c.Quote; m == nil || m.CurrentEndTime.Unix() != 1767196800 || m.InvestorType != QuoteInvestorTypeNonProfessional ||
		m.MemberInvestorType != 2 || m.QuoteDuration != 30 || m.QuoteTags[0] != "L2" || !m.ManualActivate ||
		m.PairSpuNo != "SPU" || m.QuoteMarket != "US" {
		t.Errorf("quote = %+v", m)
	}
	if m := c.Stock; m == nil || m.Amount != "12.5" || m.CounterID != "C4" || m.RedemptionType != "shares_number" ||
		m.SendWay != "account_transfer" || len(m.OptionalCounters) != 1 || m.OptionalCounters[0].Shares != "3" ||
		!m.Fallback.Enabled || m.Fallback.Type != "cash" || m.Fallback.CounterID != "C7" || m.Fallback.Shares != "1" {
		t.Errorf("stock = %+v", m)
	}
	if m := c.StockFragment; m == nil || m.UnderlyingType != CouponTargetTypeETF || m.Amounts["a"] != "1" ||
		m.Amounts["b"] != "2" || m.Market != "US" || m.CounterName != "F" {
		t.Errorf("stock fragment = %+v", m)
	}
}

func TestCouponOrders(t *testing.T) {
	var req *nethttp.Request
	srv := serve(t, `{"orders": [{"order_id": "1179", "state": "已成交", "state_val": "FilledStatus",
		"executed_price": 180.5, "executed_quantity": "100", "currency": "USD",
		"settlement_currency": "HKD", "rebate": "1.2", "rebate_dt": "2026-11-20", "state_color": "#000"}]}`, &req)
	actx := newTestAssetContext(t, srv.URL, "")
	const id = "0123456789abcdef0123456789abcdef"
	orders, err := actx.CouponOrders(context.Background(), id, CouponOrderTypeSettled)
	if err != nil {
		t.Fatalf("CouponOrders: %v", err)
	}
	if q := req.URL.Query(); req.URL.Path != "/v1/account/card-orders" || q.Get("card_id") != id || q.Get("order_type") != "settled" {
		t.Errorf("request = %s", req.URL)
	}
	if len(orders) != 1 || orders[0].OrderID != "1179" || orders[0].ExecutedPrice != "180.5" || orders[0].StateVal != "FilledStatus" {
		t.Errorf("orders = %+v", orders)
	}
	if _, err := actx.CouponOrders(context.Background(), "1024", CouponOrderTypeLocked); err == nil {
		t.Error("a bank card id must be rejected")
	}
	if _, err := actx.CouponOrders(context.Background(), id, "all"); err == nil {
		t.Error("an unknown order type must be rejected")
	}
}

func TestDeposits(t *testing.T) {
	var req *nethttp.Request
	srv := serve(t, `{"total": 20, "items": [{"id": 501, "created_at": "1795600746", "updated_at": 0,
		"currency": "HKD", "amount": "10000.00", "receivable_amount": "0", "state": "处理中", "state_code": 0,
		"sub_state_code": 2, "bank_operation_type_id": 3, "bank_operation_type_name": "FPS",
		"reference_no": "FPS123", "can_cancel": true, "user_bank_swift_code": "HSBCHKHH",
		"vouchers": [{"key": "voucher/a.jpg", "img_url": "https://x/a.jpg"},
			{"key": "voucher/b.jpg", "img_url": "https://x/b.jpg"}], "disable_link": false}]}`, &req)
	list, err := newTestAssetContext(t, srv.URL, "").Deposits(context.Background(), &GetDeposits{
		States:     []DepositState{DepositStateProcessing, DepositStateCredited},
		Currencies: []string{"hkd", " USD"},
		Start:      time.Unix(1767196800, 0),
		End:        time.Unix(1769875199, 0),
		Size:       50,
	})
	if err != nil {
		t.Fatalf("Deposits: %v", err)
	}
	want := url.Values{"page": {"1"}, "size": {"20"}, "states": {"0,1"}, "currencies": {"HKD,USD"}, "from_date": {"1767196800"}, "to_date": {"1769875199"}}
	if req.URL.Path != "/v1/account/deposits" || req.URL.Query().Encode() != want.Encode() {
		t.Errorf("request = %s", req.URL)
	}
	d := list.Items[0]
	if list.HasMore || d.ID != "501" || d.CreatedAt.Unix() != 1795600746 || !d.UpdatedAt.IsZero() ||
		d.StateCode != DepositStateProcessing || d.BankOperationTypeID != "3" || !d.CanCancel || d.ReceivableAmount != "0" ||
		d.VoucherCount != 2 {
		t.Errorf("deposit = %+v", d)
	}
	b, _ := json.Marshal(d)
	for _, leaked := range []string{"voucher/", "img_url", ".jpg", `"key"`, `"vouchers"`} {
		if strings.Contains(string(b), leaked) {
			t.Errorf("marshalled deposit leaks voucher %q: %s", leaked, b)
		}
	}
	if !strings.Contains(string(b), `"voucher_count":2`) {
		t.Errorf("marshalled deposit lacks voucher_count: %s", b)
	}
}

func TestDepositVoucherCount(t *testing.T) {
	var req *nethttp.Request
	srv := serve(t, `{"items": [
		{"id": "1", "vouchers": []},
		{"id": "2"},
		{"id": "3", "vouchers": null},
		{"id": "4", "vouchers": "voucher/x.jpg"},
		{"id": "5", "vouchers": {"key": "voucher/y.jpg"}},
		{"id": "6", "vouchers": [{"key": "voucher/a.jpg"}, {"key": "voucher/b.jpg"}]}]}`, &req)
	list, err := newTestAssetContext(t, srv.URL, "").Deposits(context.Background(), nil)
	if err != nil {
		t.Fatalf("Deposits: %v", err)
	}
	want := []int64{0, 0, 0, 0, 0, 2}
	if len(list.Items) != len(want) {
		t.Fatalf("items = %d", len(list.Items))
	}
	for i, d := range list.Items {
		if d.VoucherCount != want[i] {
			t.Errorf("deposit %s: voucher_count = %d, want %d", d.ID, d.VoucherCount, want[i])
		}
	}
}

func TestWithdrawals(t *testing.T) {
	var req *nethttp.Request
	srv := serve(t, `{"total": 1, "list": [{"id": "9", "order_id": "W1", "created_at": 1795600746,
		"currency": "HKD", "amount": "500", "fee": "0", "settlement_amount": "500", "state": 4,
		"raw_state": 7, "state_description": "已汇款", "bank_account": "012345678901",
		"display_bank_account": "尾号 8901", "account_channel": "lb"}]}`, &req)
	list, err := newTestAssetContext(t, srv.URL, "").Withdrawals(context.Background(), &GetWithdrawals{Start: time.Unix(1767196800, 0)})
	if err != nil {
		t.Fatalf("Withdrawals: %v", err)
	}
	if q := req.URL.Query(); req.URL.Path != "/v1/account/withdrawals" || q.Get("begin") != "1767196800" || q.Has("end") || q.Get("page") != "1" || q.Get("size") != "20" {
		t.Errorf("request = %s", req.URL)
	}
	w := list.List[0]
	if list.Total != 1 || list.HasMore || w.StateCode != WithdrawalStateRemitted || w.StateDescription != "已汇款" || w.RawState != "7" || w.BankAccount != "012345678901" || w.CreatedAt.Unix() != 1795600746 || w.AccountChannel != "lb" {
		t.Errorf("withdrawal = %+v", w)
	}
}

func TestAccountEndpointsAreAPOnly(t *testing.T) {
	called := false
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		called = true
	}))
	defer srv.Close()
	actx := newTestAssetContext(t, srv.URL, "us_")
	ctx := context.Background()
	calls := map[string]func() error{
		"BankCards":         func() error { _, err := actx.BankCards(ctx); return err },
		"EddaApplications":  func() error { _, err := actx.EddaApplications(ctx, "1", "HKD"); return err },
		"StockTransfersIn":  func() error { _, err := actx.StockTransfersIn(ctx, nil); return err },
		"StockTransfersOut": func() error { _, err := actx.StockTransfersOut(ctx, nil); return err },
		"RewardRecords":     func() error { _, err := actx.RewardRecords(ctx, nil); return err },
		"Coupons":           func() error { _, err := actx.Coupons(ctx, nil); return err },
		"CouponOrders": func() error {
			_, err := actx.CouponOrders(ctx, "0123456789abcdef0123456789abcdef", CouponOrderTypeLocked)
			return err
		},
		"Deposits":    func() error { _, err := actx.Deposits(ctx, nil); return err },
		"Withdrawals": func() error { _, err := actx.Withdrawals(ctx, nil); return err },
	}
	for name, call := range calls {
		var rr *http.RegionRestrictedError
		if err := call(); !errors.As(err, &rr) || rr.Required != "AP" {
			t.Errorf("%s: err = %v, want *http.RegionRestrictedError for AP", name, err)
		}
	}
	if called {
		t.Error("a US session must not send the request")
	}
}

func TestNilParamsUseDefaults(t *testing.T) {
	for _, tt := range []struct {
		name string
		got  url.Values
		want string
	}{
		{"transfers", (*GetStockTransfers)(nil).Values(), "page=1&size=20"},
		{"rewards", (*GetRewardRecords)(nil).Values(), "current=1&page_size=20"},
		{"coupons", (*GetCoupons)(nil).Values(), "card_status=-1&current=1&page_size=20"},
		{"deposits", (*GetDeposits)(nil).Values(), "page=1&size=20"},
		{"withdrawals", (*GetWithdrawals)(nil).Values(), "page=1&size=20"},
	} {
		if got := tt.got.Encode(); got != tt.want {
			t.Errorf("%s: %s, want %s", tt.name, got, tt.want)
		}
	}
}

func TestHasMoreFollowsFullPage(t *testing.T) {
	item := `{"id": "1", "state_code": 1, "state": 1}`
	page := func(n int) string {
		items := make([]string, n)
		for i := range items {
			items[i] = item
		}
		return "[" + strings.Join(items, ",") + "]"
	}
	var req *nethttp.Request
	for _, tt := range []struct {
		name string
		size int32
		n    int
		want bool
	}{
		{"default size, full", 0, 20, true},
		{"default size, short", 0, 19, false},
		{"explicit size, full", 5, 5, true},
		{"size capped at 50, full", 100, 50, true},
		{"size 30, 25 records", 30, 25, false},
		{"empty", 10, 0, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// total follows the page size upstream; it must not drive HasMore.
			srv := serve(t, `{"total": 999, "items": `+page(tt.n)+`, "list": `+page(tt.n)+`}`, &req)
			actx := newTestAssetContext(t, srv.URL, "")
			deposits, err := actx.Deposits(context.Background(), &GetDeposits{Size: tt.size})
			if err != nil {
				t.Fatalf("Deposits: %v", err)
			}
			// Deposits cap the size at 20, so 25 records fill a deposit page.
			wantDeposits := tt.want || tt.n >= maxDepositPageSize
			if deposits.HasMore != wantDeposits || len(deposits.Items) != tt.n {
				t.Errorf("deposits: has_more = %v (%d items), want %v", deposits.HasMore, len(deposits.Items), wantDeposits)
			}
			wantSize := effectivePageSize(tt.size)
			if wantSize > maxDepositPageSize {
				wantSize = maxDepositPageSize
			}
			if got := req.URL.Query().Get("size"); got != strconv.Itoa(int(wantSize)) {
				t.Errorf("deposit size sent = %s, want %d", got, wantSize)
			}
			withdrawals, err := actx.Withdrawals(context.Background(), &GetWithdrawals{Size: tt.size})
			if err != nil {
				t.Fatalf("Withdrawals: %v", err)
			}
			if withdrawals.HasMore != tt.want || withdrawals.Total != 999 {
				t.Errorf("withdrawals: has_more = %v total = %d, want %v", withdrawals.HasMore, withdrawals.Total, tt.want)
			}
		})
	}
	if b, _ := json.Marshal(&DepositList{}); strings.Contains(string(b), "total") {
		t.Errorf("deposit list must not expose total: %s", b)
	}
}

func TestFlexCode(t *testing.T) {
	for _, tt := range []struct {
		in    string
		want  int32
		valid bool
	}{
		{`0`, 0, true},
		{`"0"`, 0, true},
		{`" 7 "`, 7, true},
		{`-1`, -1, true},
		{`null`, 0, false},
		{`""`, 0, false},
		{`"abc"`, 0, false},
		{`1.5`, 0, false},
		{`true`, 0, false},
		{`{}`, 0, false},
		{`[]`, 0, false},
		{`99999999999`, 0, false},
	} {
		var c jsontypes.FlexCode
		if err := json.Unmarshal([]byte(tt.in), &c); err != nil {
			t.Errorf("%s: %v", tt.in, err)
			continue
		}
		if c.Valid != tt.valid || c.Value != tt.want {
			t.Errorf("%s: got %+v, want {%d %v}", tt.in, c, tt.want, tt.valid)
		}
	}
	// A missing field stays invalid.
	var v struct {
		Code jsontypes.FlexCode `json:"code"`
	}
	if err := json.Unmarshal([]byte(`{}`), &v); err != nil || v.Code.Valid {
		t.Errorf("missing field: %+v %v", v.Code, err)
	}
}

// TestMissingAndGarbageCodesAreUnknown checks that a status code upstream
// omitted, sent as null or sent as garbage never decodes to a real state (0 is
// a real state in most code tables).
func TestMissingAndGarbageCodesAreUnknown(t *testing.T) {
	for _, tt := range []struct {
		name string
		code string // JSON value, or "" to omit the field
	}{
		{"missing", ""},
		{"null", "null"},
		{"empty string", `""`},
		{"garbage string", `"n/a"`},
		{"float", "1.5"},
		{"object", "{}"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			field := func(name string) string {
				if tt.code == "" {
					return ""
				}
				return `"` + name + `": ` + tt.code + `,`
			}
			var req *nethttp.Request
			ctx := context.Background()

			srv := serve(t, `{"list": [{`+field("status")+field("witness_state")+field("bind_state")+
				`"id": "1", "auth_list": [{`+field("status")+`"flow_id": "hk_edda"}]}]}`, &req)
			actx := newTestAssetContext(t, srv.URL, "")
			cards, err := actx.BankCards(ctx)
			if err != nil {
				t.Fatalf("BankCards: %v", err)
			}
			c := cards[0]
			if c.Status != BankCardStatusUnknown || c.WitnessState != BankCardWitnessStateUnknown ||
				c.BindState != BankCardBindStateUnknown || c.AuthList[0].Status != AuthorizationStatusUnknown {
				t.Errorf("card = %+v", c)
			}

			srv = serve(t, `{"list": [{`+field("status")+`"id": "9"}]}`, &req)
			apps, err := newTestAssetContext(t, srv.URL, "").EddaApplications(ctx, "1", "HKD")
			if err != nil {
				t.Fatalf("EddaApplications: %v", err)
			}
			if apps[0].Status != AuthorizationStatusUnknown {
				t.Errorf("edda = %+v", apps[0])
			}

			srv = serve(t, `{"records": [{`+field("state_code")+`"id": "1", "stocks": [{`+field("state")+`"code": "AAPL"}]}]}`, &req)
			actx = newTestAssetContext(t, srv.URL, "")
			in, err := actx.StockTransfersIn(ctx, nil)
			if err != nil {
				t.Fatalf("StockTransfersIn: %v", err)
			}
			if in.Records[0].StateCode != TransferInStateUnknown || in.Records[0].Stocks[0].StateCode != TransferStockStateUnknown {
				t.Errorf("transfer in = %+v", in.Records[0])
			}
			out, err := actx.StockTransfersOut(ctx, nil)
			if err != nil {
				t.Fatalf("StockTransfersOut: %v", err)
			}
			if out.Records[0].StateCode != TransferOutStateUnknown || out.Records[0].Cancellable() {
				t.Errorf("transfer out = %+v", out.Records[0])
			}

			srv = serve(t, `{"list": [
				{`+field("reward_type")+`"prize": {`+field("status")+field("prize_type")+field("deliver_status")+field("exit_status")+`"record_id": "p"}},
				{"reward_type": 2, "card": {`+field("card_status")+`"record_id": "c"}},
				{`+field("reward_type")+`"other": {}}
			]}`, &req)
			rewards, err := newTestAssetContext(t, srv.URL, "").RewardRecords(ctx, nil)
			if err != nil {
				t.Fatalf("RewardRecords: %v", err)
			}
			if len(rewards.Records) != 2 {
				t.Fatalf("an item with neither a usable reward_type nor a detail must be skipped: %+v", rewards.Records)
			}
			prize := rewards.Records[0]
			if prize.RewardType != RewardTypeGift || prize.Status != RewardStatusUnknown || prize.PrizeType != nil ||
				prize.DeliverStatus != RewardDeliverStatusUnknown || prize.ExitStatus != -1 {
				t.Errorf("prize = %+v", prize)
			}
			if card := rewards.Records[1]; card.CardStatus != CouponStatusUnknown {
				t.Errorf("card reward = %+v", card)
			}

			srv = serve(t, `{"list": [{`+field("card_status")+`"card_id": "x", "commission": {"bz_types": [0, `+orNull(tt.code)+`]}}]}`, &req)
			coupons, err := newTestAssetContext(t, srv.URL, "").Coupons(ctx, nil)
			if err != nil {
				t.Fatalf("Coupons: %v", err)
			}
			cp := coupons.List[0]
			// An unusable bz_types entry is skipped rather than kept as Unknown.
			if cp.CardStatus != CouponStatusUnknown || len(cp.Commission.BzTypes) != 1 ||
				cp.Commission.BzTypes[0] != CouponTargetTypeStocks {
				t.Errorf("coupon = %+v %+v", cp, cp.Commission)
			}

			srv = serve(t, `{"items": [{`+field("state_code")+field("sub_state_code")+`"id": "1"}], "list": [{`+field("state")+`"id": "2"}]}`, &req)
			actx = newTestAssetContext(t, srv.URL, "")
			deposits, err := actx.Deposits(ctx, nil)
			if err != nil {
				t.Fatalf("Deposits: %v", err)
			}
			if d := deposits.Items[0]; d.StateCode != DepositStateUnknown || d.SubStateCode != -1 {
				t.Errorf("deposit = %+v", d)
			}
			withdrawals, err := actx.Withdrawals(ctx, nil)
			if err != nil {
				t.Fatalf("Withdrawals: %v", err)
			}
			if w := withdrawals.List[0]; w.StateCode != WithdrawalStateUnknown {
				t.Errorf("withdrawal = %+v", w)
			}
		})
	}
}

// orNull returns code, or null when it is empty (an array entry cannot be
// omitted).
func orNull(code string) string {
	if code == "" {
		return "null"
	}
	return code
}

// couponNullTerms is every per-type terms key of a coupon card set to null, as
// the gateway sends them for the types that do not apply.
const couponNullTerms = `"ai_benefit": null, "cash": null, "coffee_coupon": null, "crypto": null,
	"finance_interest": null, "fund_deduct": null, "fund_interest": null, "fundamental": null, "ipo": null,
	"ipo_commission": null, "ipo_interest": null, "platform_fee_deduct": null, "quote": null, "stock": null,
	"stock_fragment": null`

func TestCouponGatewayWireShape(t *testing.T) {
	var req *nethttp.Request
	// Mirrors a real gateway card (synthetic values): all per-type keys
	// present and null except the card's own, int64 values as strings, enums
	// as numbers, "0" for unset times.
	srv := serve(t, `{"total": "1", "expired_count": "0", "list": [{
		"card_id": "0123456789abcdef0123456789abcdef", "card_type": "commission", "title": "T", "title_sec": "",
		"description": "", "detail": "", "card_status": 1, "start_time": "1767196800", "end_time": "1798732800",
		"activate_time": "0", "amount": "50", "currency": "HKD", "member_id": "1000001",
		"available_trade_session": "", "trade_sessions": [], "subscription_id": "", "markets": ["HK"],
		"bz_types": [0], "rebate_rate": "", "commission_balance": "", "commission_duration": "0",
		"quote_duration": "0", "receive_duration": "0", "style": "", "cover": "", "background": "",
		"invisible": false, "new_flag": true, "account_channel": "lb",
		`+strings.Replace(couponNullTerms, `"ai_benefit": null, `, "", 1)+`,
		"ai_benefit": null,
		"commission": {"markets": ["HK"], "rebate_rate": "", "bz_types": [0, 2], "commission_balance": "50",
			"commission_duration": "2592000", "unlimited": 6, "commission_expired_time": "0", "usage_count": "5",
			"lave_usage_count": "4", "settlement_model": "", "amount_limit": "0", "superposition": false}}]}`, &req)
	list, err := newTestAssetContext(t, srv.URL, "").Coupons(context.Background(), nil)
	if err != nil {
		t.Fatalf("Coupons: %v", err)
	}
	if list.Total != 1 || len(list.List) != 1 {
		t.Fatalf("list = %+v", list)
	}
	c := list.List[0]
	if c.CardStatus != CouponStatus(1) || c.StartTime.Unix() != 1767196800 || !c.ActivateTime.IsZero() || c.Amount != "50" {
		t.Errorf("card = %+v", c)
	}
	if c.AiBenefit != nil || c.Cash != nil || c.CoffeeCoupon != nil || c.Crypto != nil || c.FinanceInterest != nil ||
		c.FundDeduct != nil || c.Fundamental != nil || c.Ipo != nil || c.IpoCommission != nil || c.IpoInterest != nil ||
		c.PlatformFeeDeduct != nil || c.Quote != nil || c.Stock != nil || c.StockFragment != nil {
		t.Errorf("null per-type terms must stay nil: %+v", c)
	}
	m := c.Commission
	if m == nil || m.Unlimited != CommissionUnlimitedTypeByAmountPer || m.CommissionDuration != 2592000 ||
		m.UsageCount != 5 || m.LaveUsageCount != 4 || len(m.BzTypes) != 2 || m.BzTypes[1] != CouponTargetTypeWarrants ||
		m.Markets[0] != "HK" || !m.CommissionExpiredTime.IsZero() {
		t.Errorf("commission = %+v", m)
	}
}

// TestCouponMalformedShapesDoNotFailPage checks that odd per-type objects and
// nested lists are dropped instead of failing the whole page.
func TestCouponMalformedShapesDoNotFailPage(t *testing.T) {
	var req *nethttp.Request
	srv := serve(t, `{"top_list": "oops", "bottom_list": [null, 1, "x"], "list": [
		{"card_id": "a", "quote": "", "stock": [], "cash": 0, "crypto": "x", "fundamental": true,
			"commission": {"bz_types": 1, "markets": "HK", "unlimited": "", "usage_count": 2.0}},
		{"card_id": "b", "commission": {"bz_types": ["", null, {}, "ETF", 2.0, 9, "Nope"], "unlimited": 5},
			"coffee_coupon": {"store_list": "x", "store_id_list": ["1", null, "x", 2.0]},
			"stock": {"optional_counters": "x", "fallback": ""},
			"crypto": {"optional_counters": [null, 1, {"counter_id": "C1"}], "fallback": []},
			"ipo_commission": {"bz_types": [0, "_", 9, "CashIPO", "bad"]},
			"stock_fragment": {"amounts": "x", "underlying_type": null},
			"quote": {"investor_type": 7, "markets": [null, "US", {}]}},
		"not a card",
		{"card_id": "c", "card_status": "2.0"}]}`, &req)
	list, err := newTestAssetContext(t, srv.URL, "").Coupons(context.Background(), nil)
	if err != nil {
		t.Fatalf("Coupons: %v", err)
	}
	if len(list.TopList) != 0 || len(list.BottomList) != 0 || len(list.List) != 3 {
		t.Fatalf("list = %+v", list)
	}
	a := list.List[0]
	if a.Quote != nil || a.Stock != nil || a.Cash != nil || a.Crypto != nil || a.Fundamental != nil || a.Commission == nil {
		t.Errorf("non-object terms must be nil: %+v", a)
	}
	if m := a.Commission; len(m.BzTypes) != 0 || len(m.Markets) != 0 || m.Unlimited != CommissionUnlimitedTypeUnknown || m.UsageCount != 2 {
		t.Errorf("commission a = %+v", m)
	}
	b := list.List[1]
	if m := b.Commission; len(m.BzTypes) != 3 || m.BzTypes[0] != CouponTargetTypeETF ||
		m.BzTypes[1] != CouponTargetTypeWarrants || m.BzTypes[2] != 9 || m.Unlimited != 5 {
		t.Errorf("commission b = %+v", m)
	}
	if m := b.CoffeeCoupon; m == nil || len(m.StoreList) != 0 || len(m.StoreIDList) != 2 || m.StoreIDList[1] != 2 {
		t.Errorf("coffee = %+v", m)
	}
	if m := b.Stock; m == nil || len(m.OptionalCounters) != 0 || m.Fallback.Enabled {
		t.Errorf("stock = %+v", m)
	}
	if m := b.Crypto; m == nil || len(m.OptionalCounters) != 1 || m.OptionalCounters[0].CounterID != "C1" {
		t.Errorf("crypto = %+v", m)
	}
	if m := b.IpoCommission; m == nil || len(m.BzTypes) != 4 || m.BzTypes[0] != IpoCouponTargetTypeUnknown ||
		m.BzTypes[1] != IpoCouponTargetTypeUnknown || m.BzTypes[2] != 9 || m.BzTypes[3] != IpoCouponTargetTypeCash {
		t.Errorf("ipo commission = %+v", m)
	}
	if m := b.StockFragment; m == nil || m.Amounts != nil || m.UnderlyingType != CouponTargetTypeUnknown {
		t.Errorf("stock fragment = %+v", m)
	}
	if m := b.Quote; m == nil || m.InvestorType != 7 || len(m.Markets) != 1 || m.Markets[0] != "US" {
		t.Errorf("quote = %+v", m)
	}
	if c := list.List[2]; c.CardID != "c" || c.CardStatus != CouponStatusExpired {
		t.Errorf("integral float card_status = %+v", c)
	}
}

// TestUndocumentedCodesAreKept checks the one rule for integer codes: any
// integer is kept as is, documented or not.
func TestUndocumentedCodesAreKept(t *testing.T) {
	var req *nethttp.Request
	ctx := context.Background()
	srv := serve(t, `{"list": [{"reward_type": 3, "prize": {"status": 12, "deliver_status": 11, "prize_type": 0}},
		{"reward_type": 77, "prize": {"status": 99.0, "deliver_status": "42"}}]}`, &req)
	rewards, err := newTestAssetContext(t, srv.URL, "").RewardRecords(ctx, nil)
	if err != nil {
		t.Fatalf("RewardRecords: %v", err)
	}
	if r := rewards.Records[0]; r.DeliverStatus != 11 {
		t.Errorf("deliver status 11 = %+v", r)
	}
	if r := rewards.Records[1]; r.RewardType != 77 || r.Status != 99 || r.DeliverStatus != 42 {
		t.Errorf("undocumented codes = %+v", r)
	}

	srv = serve(t, `{"items": [{"state_code": 17, "sub_state_code": "8"}], "list": [{"state": 23.0}]}`, &req)
	actx := newTestAssetContext(t, srv.URL, "")
	deposits, err := actx.Deposits(ctx, nil)
	if err != nil {
		t.Fatalf("Deposits: %v", err)
	}
	if d := deposits.Items[0]; d.StateCode != 17 || d.SubStateCode != 8 {
		t.Errorf("deposit = %+v", d)
	}
	withdrawals, err := actx.Withdrawals(ctx, nil)
	if err != nil {
		t.Fatalf("Withdrawals: %v", err)
	}
	if w := withdrawals.List[0]; w.StateCode != 23 {
		t.Errorf("withdrawal = %+v", w)
	}
}

func TestFlexIntegralFloats(t *testing.T) {
	var v struct {
		A jsontypes.FlexInt64 `json:"a"`
		B jsontypes.FlexInt64 `json:"b"`
		C jsontypes.FlexInt64 `json:"c"`
		D jsontypes.FlexCode  `json:"d"`
		E jsontypes.FlexCode  `json:"e"`
		F jsontypes.FlexCode  `json:"f"`
	}
	if err := json.Unmarshal([]byte(`{"a": "6.0", "b": 2.0, "c": 1.5, "d": 2.0, "e": "1e1", "f": 3e10}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v.A != 6 || v.B != 2 || v.C != 0 || !v.D.Valid || v.D.Value != 2 || !v.E.Valid || v.E.Value != 10 || v.F.Valid {
		t.Errorf("got %+v", v)
	}
}

// testJWT builds an unsigned JWT-shaped token with the given payload.
func testJWT(payload string) string {
	enc := func(s string) string { return base64.RawURLEncoding.EncodeToString([]byte(s)) }
	return enc(`{"alg":"none"}`) + "." + enc(payload) + ".sig"
}

func TestAccountChannelQuery(t *testing.T) {
	for _, tt := range []struct {
		name, token, want string
	}{
		{"ac claim", testJWT(`{"ac":"lb_papertrading"}`), "lb_papertrading"},
		{"sub claim", "hk_m_" + testJWT(`{"sub":"{\"member_id\":1,\"account_channel\":\"lb\"}"}`), "lb"},
		{"no claim defaults to lb", testJWT(`{"exp":1,"sub":"not json"}`), "lb"},
		{"undecodable payload", "aaa.!!!.sig", ""},
		{"not a jwt", "opaque-token", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var req *nethttp.Request
			srv := serve(t, `{"list": []}`, &req)
			cli, err := http.New(http.WithURL(srv.URL), http.WithAppKey("key"), http.WithAppSecret("secret"),
				http.WithAccessToken(tt.token))
			if err != nil {
				t.Fatalf("http.New: %v", err)
			}
			actx := &AssetContext{httpClient: cli}
			if _, err := actx.Coupons(context.Background(), nil); err != nil {
				t.Fatalf("Coupons: %v", err)
			}
			q := req.URL.Query()
			if q.Get("account_channel") != tt.want || (tt.want == "" && q.Has("account_channel")) {
				t.Errorf("account_channel = %q, want %q (%s)", q.Get("account_channel"), tt.want, req.URL)
			}
			if q.Get("page_size") != "20" {
				t.Errorf("other params lost: %s", req.URL)
			}
		})
	}
}

// TestDepositsPageSizeCappedAt20 checks that Deposits never asks for more
// than 20 records (upstream returns only 10 for a size of 21 to 50) and that
// HasMore follows the size actually sent.
func TestDepositsPageSizeCappedAt20(t *testing.T) {
	item := `{"id": "1", "state_code": 1}`
	items := make([]string, 20)
	for i := range items {
		items[i] = item
	}
	var req *nethttp.Request
	srv := serve(t, `{"items": [`+strings.Join(items, ",")+`]}`, &req)
	actx := newTestAssetContext(t, srv.URL, "")
	for _, size := range []int32{0, 20, 21, 50, 100} {
		list, err := actx.Deposits(context.Background(), &GetDeposits{Size: size})
		if err != nil {
			t.Fatalf("Deposits: %v", err)
		}
		if got := req.URL.Query().Get("size"); got != "20" {
			t.Errorf("size %d: sent %s, want 20", size, got)
		}
		if !list.HasMore {
			t.Errorf("size %d: 20 records must report HasMore", size)
		}
	}
	if _, err := actx.Deposits(context.Background(), &GetDeposits{Size: 5}); err != nil {
		t.Fatalf("Deposits: %v", err)
	}
	if got := req.URL.Query().Get("size"); got != "5" {
		t.Errorf("size 5: sent %s", got)
	}
	if _, err := actx.Deposits(context.Background(), nil); err != nil {
		t.Fatalf("Deposits: %v", err)
	}
	if got := req.URL.Query().Get("size"); got != "20" {
		t.Errorf("nil params: sent %s", got)
	}
}

func TestRewardDeliverStatusIsAddressFilled(t *testing.T) {
	for s, want := range map[RewardDeliverStatus]bool{
		RewardDeliverStatusAddressFilled:    true,
		RewardDeliverStatusAddressFilledAlt: true,
		RewardDeliverStatusAddressMissing:   false,
		RewardDeliverStatusShipping:         false,
		RewardDeliverStatusUnknown:          false,
	} {
		if got := s.IsAddressFilled(); got != want {
			t.Errorf("%d.IsAddressFilled() = %v, want %v", s, got, want)
		}
	}
}
