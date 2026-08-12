package grid

import (
	"github.com/shopspring/decimal"
)

// TriggerPriceType is how grid trigger thresholds are interpreted (wire: int32).
type TriggerPriceType int32

const (
	// TriggerPriceTypeUnknown is the unknown / unset value.
	TriggerPriceTypeUnknown TriggerPriceType = 0
	// TriggerPriceTypeSpread triggers by an absolute price spread.
	TriggerPriceTypeSpread TriggerPriceType = 1
	// TriggerPriceTypePercent triggers by percent.
	TriggerPriceTypePercent TriggerPriceType = 2
)

// GridTimeInForce is the time in force for a grid order (wire: int32).
// Values other than the named ones are preserved verbatim.
type GridTimeInForce int32

const (
	// GridTimeInForceDay is a day order.
	GridTimeInForceDay GridTimeInForce = 0
	// GridTimeInForceGoodTilCanceled is a good-til-canceled order.
	GridTimeInForceGoodTilCanceled GridTimeInForce = 1
	// GridTimeInForceGoodTilDate is a good-til-date order.
	GridTimeInForceGoodTilDate GridTimeInForce = 6
)

// GridLimitEvent is the action taken when a grid boundary is reached
// (wire: int32).
type GridLimitEvent int32

const (
	// GridLimitEventUnknown is the unknown / unset value.
	GridLimitEventUnknown GridLimitEvent = 0
	// GridLimitEventIgnore keeps the grid running when the boundary is reached.
	GridLimitEventIgnore GridLimitEvent = 1
	// GridLimitEventCloseAtLast closes the position at the last price.
	GridLimitEventCloseAtLast GridLimitEvent = 2
)

// GridTrigger expresses a grid's up/down trigger thresholds. Percent and
// spread are mutually exclusive; use PercentTrigger or SpreadTrigger to build
// one so the choice is explicit (instead of four independent optional fields).
type GridTrigger struct {
	priceType TriggerPriceType
	up        decimal.Decimal
	down      decimal.Decimal
}

// PercentTrigger builds a percent-based grid trigger.
func PercentTrigger(up, down decimal.Decimal) GridTrigger {
	return GridTrigger{priceType: TriggerPriceTypePercent, up: up, down: down}
}

// SpreadTrigger builds an absolute-spread-based grid trigger.
func SpreadTrigger(up, down decimal.Decimal) GridTrigger {
	return GridTrigger{priceType: TriggerPriceTypeSpread, up: up, down: down}
}

// NewGridTradeRule creates a rule with the fields a valid grid order requires.
// The gateway still validates business rules, but this makes the minimum field
// set visible in the signature instead of leaving all fields optional. Chain
// the setter methods to populate the optional fields.
func NewGridTradeRule(
	basePrice decimal.Decimal,
	upperPrice decimal.Decimal,
	lowerPrice decimal.Decimal,
	trigger GridTrigger,
	quantity decimal.Decimal,
	upperQuantity decimal.Decimal,
	lowerQuantity decimal.Decimal,
	timeInForce GridTimeInForce,
) *GridTradeRule {
	rule := &GridTradeRule{
		SubmittedBasePrice: basePrice,
		UpperLimitPrice:    upperPrice,
		LowerLimitPrice:    lowerPrice,
		TriggerPriceType:   trigger.priceType,
		TriggerQuantity:    quantity,
		UpperLimitQuantity: upperQuantity,
		LowerLimitQuantity: lowerQuantity,
		TimeInForce:        timeInForce,
	}
	switch trigger.priceType {
	case TriggerPriceTypePercent:
		rule.TriggerPercentUp = trigger.up
		rule.TriggerPercentDown = trigger.down
	case TriggerPriceTypeSpread:
		rule.TriggerSpreadUp = trigger.up
		rule.TriggerSpreadDown = trigger.down
	}
	return rule
}

// LimitEvents sets the actions taken at the upper / lower bounds.
func (r *GridTradeRule) LimitEvents(upper, lower GridLimitEvent) *GridTradeRule {
	r.UpperLimitEvent = upper
	r.LowerLimitEvent = lower
	return r
}

// Depths sets the sell / buy order-book depths (0 = use the order type).
func (r *GridTradeRule) Depths(sell, buy int32) *GridTradeRule {
	r.TriggerSellDepth = sell
	r.TriggerBuyDepth = buy
	return r
}

// OrderTypes sets the sell / buy order types (GMO / GLO / GTG).
func (r *GridTradeRule) OrderTypes(up, down string) *GridTradeRule {
	r.GridOrderTypeUp = up
	r.GridOrderTypeDown = down
	return r
}

// MultipleTriggerFlag allows a single grid level to trigger multiple times.
func (r *GridTradeRule) MultipleTriggerFlag(value bool) *GridTradeRule {
	r.MultipleTrigger = value
	return r
}

// SupportShortsellFlag allows short selling.
func (r *GridTradeRule) SupportShortsellFlag(value bool) *GridTradeRule {
	r.SupportShortsell = value
	return r
}

// RthFlag sets the regular-trading-hours flag (0 / 1 / 2).
func (r *GridTradeRule) RthFlag(value int32) *GridTradeRule {
	r.Rth = value
	return r
}

// ExpireTimeAt sets the expiry time (unix seconds), used with a GTD
// time-in-force.
func (r *GridTradeRule) ExpireTimeAt(unixSeconds int64) *GridTradeRule {
	r.ExpireTime = unixSeconds
	return r
}
