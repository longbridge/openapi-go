package alert

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// AlertCondition represents the trigger condition for a price alert.
type AlertCondition int

const (
	// AlertConditionPriceRise triggers when the price rises to a given value.
	AlertConditionPriceRise AlertCondition = 1
	// AlertConditionPriceFall triggers when the price falls to a given value.
	AlertConditionPriceFall AlertCondition = 2
	// AlertConditionPercentRise triggers when the percentage change rises to a given value.
	AlertConditionPercentRise AlertCondition = 3
	// AlertConditionPercentFall triggers when the percentage change falls to a given value.
	AlertConditionPercentFall AlertCondition = 4
)

// AlertFrequency controls how often a triggered alert fires.
type AlertFrequency int

const (
	// AlertFrequencyDaily fires at most once per day.
	AlertFrequencyDaily AlertFrequency = 1
	// AlertFrequencyEveryTime fires every time the condition is met.
	AlertFrequencyEveryTime AlertFrequency = 2
	// AlertFrequencyOnce fires exactly once.
	AlertFrequencyOnce AlertFrequency = 3
)

// AlertList is the top-level response containing groups of alerts per security.
type AlertList struct {
	// Lists holds alert groups, one per security.
	Lists []*AlertSymbolGroup
}

// AlertSymbolGroup holds all price alerts for a single security.
type AlertSymbolGroup struct {
	// Symbol is the security identifier (e.g. "700.HK").
	Symbol string
	// Code is the short code of the security.
	Code string
	// Market is the market identifier (e.g. "HK").
	Market string
	// Name is the display name of the security.
	Name string
	// Price is the current price, if available.
	Price *decimal.Decimal
	// Chg is the price change, if available.
	Chg *decimal.Decimal
	// PChg is the percentage change, if available.
	PChg *decimal.Decimal
	// Product is the product type string (may be empty).
	Product string
	// Indicators is the list of individual alert items for this security.
	Indicators []*AlertItem
}

// AlertItem is a single price-alert configuration.
type AlertItem struct {
	// ID is the unique identifier for this alert.
	ID string
	// IndicatorID is the condition type code (matches AlertCondition values).
	IndicatorID string
	// Enabled indicates whether the alert is currently active.
	Enabled bool
	// Frequency controls how often the alert fires.
	Frequency int
	// Scope is an internal scope value.
	Scope int
	// Text is the human-readable description of the alert trigger (e.g. "价格涨到 600").
	Text string
	// State tracks the current trigger state.
	State []int
	// ValueMap holds the threshold values as raw JSON (e.g. {"price":"600"}).
	ValueMap json.RawMessage
}
