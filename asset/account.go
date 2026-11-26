package asset

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/longbridge/openapi-go/asset/jsontypes"
)

// The account endpoints below are served by the AP (Asia-Pacific) data center
// only: US data-center accounts get an error upstream, so these calls return
// *http.RegionRestrictedError for US sessions without sending the request.
// There is no US equivalent. They need the account.read OAuth scope.

const (
	pathBankCards         = "/v2/account/bank-cards"
	pathEddaStatus        = "/v1/account/edda-status"
	pathStockTransfersIn  = "/v1/account/stock-transfers/in"
	pathStockTransfersOut = "/v1/account/stock-transfers/out"
	pathRewards           = "/v1/account/rewards"
	pathCoupons           = "/v1/account/cards"
	pathCouponOrders      = "/v1/account/card-orders"
	pathDeposits          = "/v1/account/deposits"
	pathWithdrawals       = "/v1/account/withdrawals"
)

// couponIDPattern matches a coupon card id: 32 hexadecimal characters.
var couponIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)

// get checks the data center and sends a GET request.
func (c *AssetContext) get(ctx context.Context, path string, q url.Values, resp interface{}) error {
	if err := c.httpClient.CheckRegionCtx(ctx, path, "AP"); err != nil {
		return err
	}
	return c.httpClient.Get(ctx, path, q, resp)
}

// BankCards returns the linked bank cards with their verification, binding
// and authorization (eDDA, bank-securities link, FPS) status.
//
// Path: GET /v2/account/bank-cards. AP accounts only.
//
// Example:
//
//	cards, err := actx.BankCards(context.Background())
//	for _, card := range cards {
//	    for _, auth := range card.AuthList {
//	        if auth.IsEDDA() { fmt.Println(card.Bank, auth.Currency, auth.Status) }
//	    }
//	}
func (c *AssetContext) BankCards(ctx context.Context) ([]*BankCard, error) {
	q := url.Values{}
	q.Set("from_app_page", "cards")
	var resp jsontypes.BankCardsResponse
	if err := c.get(ctx, pathBankCards, q, &resp); err != nil {
		return nil, err
	}
	cards := make([]*BankCard, 0, len(resp.List))
	for i := range resp.List {
		cards = append(cards, toBankCard(&resp.List[i]))
	}
	return cards, nil
}

// EddaApplications returns the eDDA authorization applications of one bank
// card in one currency. cardID is BankCard.ID; currency is e.g. "HKD".
//
// Path: GET /v1/account/edda-status. AP accounts only.
func (c *AssetContext) EddaApplications(ctx context.Context, cardID int64, currency string) ([]*EddaApplication, error) {
	if cardID <= 0 {
		return nil, fmt.Errorf("asset: EddaApplications requires the numeric bank card id (BankCard.ID)")
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return nil, fmt.Errorf("asset: EddaApplications requires a currency")
	}
	q := url.Values{}
	q.Set("card_id", strconv.FormatInt(cardID, 10))
	q.Set("currency", currency)
	var resp json.RawMessage
	if err := c.get(ctx, pathEddaStatus, q, &resp); err != nil {
		return nil, err
	}
	return decodeEddaApplications(resp)
}

// StockTransfersIn returns one page of stock transfer-in applications
// (another broker → this account), newest first.
//
// Path: GET /v1/account/stock-transfers/in. AP accounts only.
func (c *AssetContext) StockTransfersIn(ctx context.Context, params *GetStockTransfers) (*StockTransferInList, error) {
	var resp jsontypes.TransferRecordsResponse
	if err := c.get(ctx, pathStockTransfersIn, params.Values(), &resp); err != nil {
		return nil, err
	}
	out := &StockTransferInList{
		Records: make([]*StockTransferIn, 0, len(resp.Records)),
		Total:   int64(resp.Total),
	}
	for i := range resp.Records {
		out.Records = append(out.Records, toStockTransferIn(&resp.Records[i]))
	}
	return out, nil
}

// StockTransfersOut returns one page of stock transfer-out applications (this
// account → another broker), newest first.
//
// Path: GET /v1/account/stock-transfers/out. AP accounts only.
func (c *AssetContext) StockTransfersOut(ctx context.Context, params *GetStockTransfers) (*StockTransferOutList, error) {
	var resp jsontypes.TransferRecordsResponse
	if err := c.get(ctx, pathStockTransfersOut, params.Values(), &resp); err != nil {
		return nil, err
	}
	out := &StockTransferOutList{Records: make([]*StockTransferOut, 0, len(resp.Records))}
	for i := range resp.Records {
		out.Records = append(out.Records, toStockTransferOut(&resp.Records[i]))
	}
	return out, nil
}

// RewardRecords returns one page of reward-center records (cash, gifts,
// coupons, task rewards). The endpoint has no date filter: to find records in
// a period, page through and filter on SendTime.
//
// Path: GET /v1/account/rewards. AP accounts only.
func (c *AssetContext) RewardRecords(ctx context.Context, params *GetRewardRecords) (*RewardRecordList, error) {
	var resp jsontypes.RewardsResponse
	if err := c.get(ctx, pathRewards, params.Values(), &resp); err != nil {
		return nil, err
	}
	out := &RewardRecordList{
		Records: make([]*RewardRecord, 0, len(resp.List)),
		Total:   int64(resp.Total),
	}
	for i := range resp.List {
		if rec := toRewardRecord(&resp.List[i]); rec != nil {
			out.Records = append(out.Records, rec)
		}
	}
	return out, nil
}

// Coupons returns one page of coupon cards. TopList and BottomList are pinned
// sections that upstream repeats on every page.
//
// Path: GET /v1/account/cards. AP accounts only.
func (c *AssetContext) Coupons(ctx context.Context, params *GetCoupons) (*CouponList, error) {
	var resp jsontypes.CouponsResponse
	if err := c.get(ctx, pathCoupons, params.Values(), &resp); err != nil {
		return nil, err
	}
	return &CouponList{
		TopList:      toCoupons(resp.TopList),
		List:         toCoupons(resp.List),
		BottomList:   toCoupons(resp.BottomList),
		Total:        int64(resp.Total),
		ExpiredCount: int64(resp.ExpiredCount),
	}, nil
}

// CouponOrders returns the orders a coupon card was applied to: those it is
// locked to (CouponOrderTypeLocked) or whose rebate was settled
// (CouponOrderTypeSettled). cardID is the 32-hex Coupon.CardID.
//
// The gateway rate-limits this route per member tightly; fetch the two order
// types one after the other rather than concurrently.
//
// Path: GET /v1/account/card-orders. AP accounts only.
func (c *AssetContext) CouponOrders(ctx context.Context, cardID string, orderType CouponOrderType) ([]*CouponOrder, error) {
	cardID = strings.TrimSpace(cardID)
	if !couponIDPattern.MatchString(cardID) {
		return nil, fmt.Errorf("asset: CouponOrders requires the 32-character coupon id (Coupon.CardID), got %q", cardID)
	}
	if orderType != CouponOrderTypeLocked && orderType != CouponOrderTypeSettled {
		return nil, fmt.Errorf("asset: CouponOrders order type must be %q or %q, got %q", CouponOrderTypeLocked, CouponOrderTypeSettled, orderType)
	}
	q := url.Values{}
	q.Set("card_id", cardID)
	q.Set("order_type", string(orderType))
	var resp jsontypes.CouponOrdersResponse
	if err := c.get(ctx, pathCouponOrders, q, &resp); err != nil {
		return nil, err
	}
	orders := make([]*CouponOrder, 0, len(resp.Orders))
	for i := range resp.Orders {
		orders = append(orders, toCouponOrder(&resp.Orders[i]))
	}
	return orders, nil
}

// Deposits returns one page of deposit records.
//
// Path: GET /v1/account/deposits. AP accounts only.
func (c *AssetContext) Deposits(ctx context.Context, params *GetDeposits) (*DepositList, error) {
	var resp jsontypes.DepositsResponse
	if err := c.get(ctx, pathDeposits, params.Values(), &resp); err != nil {
		return nil, err
	}
	out := &DepositList{
		Items: make([]*Deposit, 0, len(resp.Items)),
		Total: int64(resp.Total),
	}
	for i := range resp.Items {
		out.Items = append(out.Items, toDeposit(&resp.Items[i]))
	}
	return out, nil
}

// Withdrawals returns one page of withdrawal records.
//
// Path: GET /v1/account/withdrawals. AP accounts only.
func (c *AssetContext) Withdrawals(ctx context.Context, params *GetWithdrawals) (*WithdrawalList, error) {
	var resp jsontypes.WithdrawalsResponse
	if err := c.get(ctx, pathWithdrawals, params.Values(), &resp); err != nil {
		return nil, err
	}
	out := &WithdrawalList{
		List:  make([]*Withdrawal, 0, len(resp.List)),
		Total: int64(resp.Total),
	}
	for i := range resp.List {
		out.List = append(out.List, toWithdrawal(&resp.List[i]))
	}
	return out, nil
}

// setUnix sets key to t as unix seconds when t is not zero.
func setUnix(v url.Values, key string, t time.Time) {
	if !t.IsZero() {
		v.Set(key, strconv.FormatInt(t.Unix(), 10))
	}
}

// setPaging sets the page number and size, defaulting to 1 and 20.
func setPaging(v url.Values, pageKey string, page int32, sizeKey string, size int32) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	v.Set(pageKey, strconv.Itoa(int(page)))
	v.Set(sizeKey, strconv.Itoa(int(size)))
}
