package trade

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/longbridge/openapi-go/trade/jsontypes"
)

// GetUSAvailableHolding is the request for USAvailableHolding.
type GetUSAvailableHolding struct {
	// Symbol is the security symbol, e.g. "AAPL.US", an option symbol such as
	// "AAPL261120C345000.US", or a crypto pair such as "BTCUSD.BKKT". Required.
	Symbol string
	// EquityType is the holding type: "STOCK", "OPTION" or "CRYPTO".
	// Empty defaults to "STOCK".
	EquityType string
	// PositionSide is "LONG" or "SHORT". Empty defaults to "LONG".
	PositionSide string
}

// USAvailableHolding is the response for GET /v1/us/assets/available-holding.
type USAvailableHolding struct {
	// AvailableQty is the quantity available to sell (long) or buy to cover
	// (short); negative for a short position.
	AvailableQty string `json:"available_qty"`
}

// USBuyLimitInfo is the buying power in USCashAndBuyPower.
type USBuyLimitInfo struct {
	Currency               string `json:"currency"`
	CashBuyPower           string `json:"cash_buy_power"`
	OvernightBuyPower      string `json:"overnight_buy_power"`
	DayTradeBuyPower       string `json:"day_trade_buy_power"`
	InstantBuyPower        string `json:"instant_buy_power"`
	IntradayMarginBuyPower string `json:"intraday_margin_buy_power"`
	OptionBuyPower         string `json:"option_buy_power"`
	CryptoBuyPower         string `json:"crypto_buy_power"`
	// MME is the margin maintenance excess used for buying power.
	MME string `json:"mme"`
}

// USCashDetail is one cash currency entry in USCashAndBuyPower.
type USCashDetail struct {
	Currency      string `json:"currency"`
	TotalAmount   string `json:"total_amount"`
	AvailableCash string `json:"available_cash"`
	WithdrawCash  string `json:"withdraw_cash"`
	FrozenCash    string `json:"frozen_cash"`
	SettledCash   string `json:"settled_cash"`
	Outstanding   string `json:"outstanding"`
	// ExternalAmount is cash held for external use (e.g. IPO or fund subscriptions).
	ExternalAmount string `json:"external_amount"`
}

// USMarginStatus is the margin account status in USCashAndBuyPower.
type USMarginStatus struct {
	Currency string `json:"currency"`
	// NLV is the net liquidation value.
	NLV                     string `json:"nlv"`
	MarginEquity            string `json:"margin_equity"`
	MarginMaintenanceExcess string `json:"margin_maintenance_excess"`
	// MMR is the maintenance margin requirement.
	MMR        string `json:"mmr"`
	RegTExcess string `json:"reg_t_excess"`
	// SMA is the special memorandum account.
	SMA string `json:"sma"`
}

// USCashAndBuyPower is the response for GET /v1/us/assets/cash-and-buy-power.
type USCashAndBuyPower struct {
	BuyLimitInfo USBuyLimitInfo `json:"buy_limit_info"`
	CashList     []USCashDetail `json:"cash_list"`
	MarginStatus USMarginStatus `json:"margin_status"`
}

// GetUSTransactionHistory is the request for USTransactionHistory.
type GetUSTransactionHistory struct {
	// StartAt is the start time (inclusive). Required.
	StartAt time.Time
	// EndAt is the end time (inclusive). Required.
	EndAt time.Time
	// Page is the 1-based page number; 0 defaults to 1.
	Page int32
	// Size is the page size; 0 defaults to 20.
	Size int32
	// Types filters by transaction type; empty means all. Values:
	// "TT_DEPOSIT_AND_WITHDRAWAL", "TT_STOCKS_AND_ETFS_TRADING",
	// "TT_OPTIONS_TRADING", "TT_CRYPTO_TRADING", "TT_CORPORATE_ACTION",
	// "TT_TRANSFER_IN_AND_OUT", "TT_CRYPTO_TRANSFER", "TT_MARGIN_FINANCING",
	// "TT_OTHER".
	Types []string
	// Currency filters by currency, e.g. "USD"; empty means all.
	Currency string
}

// USTransaction is one record in USTransactionHistory.
type USTransaction struct {
	Title       string `json:"title"`
	BizCode     string `json:"biz_code"`
	BizCodeName string `json:"biz_code_name"`
	// Amount is negative for outflows.
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	// Symbol is the related security symbol, if any (e.g. "AAPL.US").
	Symbol              string `json:"symbol"`
	Market              string `json:"market"`
	EquityType          string `json:"equity_type"`
	PositionSide        string `json:"position_side"`
	Quantity            string `json:"quantity"`
	QuantityDescription string `json:"quantity_description"`
	RefNo               string `json:"ref_no"`
	Remark              string `json:"remark"`
	// CreateAt is the creation time; zero when not provided.
	CreateAt time.Time `json:"create_at"`
	// DoneAt is the completion time; zero when not provided.
	DoneAt time.Time `json:"done_at"`
}

// UnmarshalJSON reads create_at / done_at as unix seconds (number or string)
// and accepts the symbol under either `symbol` or `counter_id`.
func (t *USTransaction) UnmarshalJSON(data []byte) error {
	type plain USTransaction
	var raw struct {
		plain
		CounterID string        `json:"counter_id"`
		CreateAt  unixTimestamp `json:"create_at"`
		DoneAt    unixTimestamp `json:"done_at"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*t = USTransaction(raw.plain)
	if t.Symbol == "" {
		t.Symbol = raw.CounterID
	}
	t.CreateAt = raw.CreateAt.Time()
	t.DoneAt = raw.DoneAt.Time()
	return nil
}

// USTransactionHistory is the response for POST /v1/us/assets/transaction-history.
type USTransactionHistory struct {
	// CurrentPage is the current page number.
	CurrentPage int32 `json:"current_page"`
	// Items are the records on this page. The API returns no total count: a
	// page shorter than the requested size is the last one.
	Items []USTransaction `json:"items"`
}

// UnmarshalJSON accepts current_page as a JSON number or a numeric string.
func (h *USTransactionHistory) UnmarshalJSON(data []byte) error {
	var raw struct {
		CurrentPage jsontypes.FlexInt32 `json:"current_page"`
		Items       []USTransaction     `json:"items"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	h.CurrentPage = int32(raw.CurrentPage)
	h.Items = raw.Items
	return nil
}

// USExchangeRate is one exchange rate in USExchangeRates.
type USExchangeRate struct {
	// UniversalKey is the account channel the rates apply to, e.g. "us_lb_apex"
	// (outer `universal_exchanges` key).
	UniversalKey string `json:"universal_key"`
	// SpecKey is the rate type, e.g. "am" or "reference" (inner
	// `spec_exchanges` key).
	SpecKey       string `json:"spec_key"`
	BaseCurrency  string `json:"base_currency"`
	OtherCurrency string `json:"other_currency"`
	// BidRate, OfferRate and AverageRate accept a JSON number or string
	// upstream; empty when not provided.
	BidRate     string `json:"bid_rate"`
	OfferRate   string `json:"offer_rate"`
	AverageRate string `json:"average_rate"`
}

// USExchangeRates is the response for GET /v1/us/assets/exchange-rates.
//
// The upstream response is a two-level map
// (universal_exchanges.<key>.spec_exchanges.<key>.exchanges[]); it is
// flattened into Exchanges, sorted by the two upstream map keys, which are
// kept on each entry as UniversalKey and SpecKey.
type USExchangeRates struct {
	Exchanges []USExchangeRate `json:"exchanges"`
}

type usRawExchangeRate struct {
	BaseCurrency  string               `json:"base_currency"`
	OtherCurrency string               `json:"other_currency"`
	BidRate       jsontypes.FlexString `json:"bid_rate"`
	OfferRate     jsontypes.FlexString `json:"offer_rate"`
	AverageRate   jsontypes.FlexString `json:"average_rate"`
}

type usRawExchangeRates struct {
	UniversalExchanges map[string]struct {
		SpecExchanges map[string]struct {
			Exchanges []usRawExchangeRate `json:"exchanges"`
		} `json:"spec_exchanges"`
	} `json:"universal_exchanges"`
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (r usRawExchangeRates) flatten() *USExchangeRates {
	out := &USExchangeRates{Exchanges: []USExchangeRate{}}
	for _, uk := range sortedKeys(r.UniversalExchanges) {
		spec := r.UniversalExchanges[uk].SpecExchanges
		for _, sk := range sortedKeys(spec) {
			for _, e := range spec[sk].Exchanges {
				out.Exchanges = append(out.Exchanges, USExchangeRate{
					UniversalKey:  uk,
					SpecKey:       sk,
					BaseCurrency:  e.BaseCurrency,
					OtherCurrency: e.OtherCurrency,
					BidRate:       string(e.BidRate),
					OfferRate:     string(e.OfferRate),
					AverageRate:   string(e.AverageRate),
				})
			}
		}
	}
	return out
}
