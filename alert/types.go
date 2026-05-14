package alert

import "github.com/shopspring/decimal"

// AlertCondition is the price alert trigger condition.
type AlertCondition int32

const (
	AlertConditionPriceRise   AlertCondition = 1
	AlertConditionPriceFall   AlertCondition = 2
	AlertConditionPercentRise AlertCondition = 3
	AlertConditionPercentFall AlertCondition = 4
)

// AlertFrequency controls how often an alert fires.
type AlertFrequency int32

const (
	AlertFrequencyDaily     AlertFrequency = 1
	AlertFrequencyEveryTime AlertFrequency = 2
	AlertFrequencyOnce      AlertFrequency = 3
)

// AlertList is the response for List.
type AlertList struct {
	Lists []*AlertSymbolGroup
}

// AlertSymbolGroup groups alert items by security.
type AlertSymbolGroup struct {
	Symbol  string
	Code    string
	Market  string
	Name    string
	Price   *decimal.Decimal
	Chg     *decimal.Decimal
	PChg    *decimal.Decimal
	Product string
	Items   []*AlertItem
}

// AlertItem is a single price alert.
type AlertItem struct {
	Id          string
	IndicatorId string
	Enabled     bool
	Frequency   int32
	Scope       int32
	Text        string
	State       []int32
	ValueMap    map[string]interface{}
}
