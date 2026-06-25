package jsontypes

import (
	"encoding/json"

	"github.com/shopspring/decimal"
)

// AlertList is the top-level response from GET /v1/notify/reminders.
type AlertList struct {
	Lists []*AlertSymbolGroup `json:"lists"`
}

// AlertSymbolGroup holds all alert items for a single security.
type AlertSymbolGroup struct {
	// Symbol is the security identifier (e.g. "700.HK"), decoded from counter_id.
	Symbol     string           `json:"counter_id"`
	Code       string           `json:"code"`
	Market     string           `json:"market"`
	Name       string           `json:"name"`
	Price      *decimal.Decimal `json:"price"`
	Chg        *decimal.Decimal `json:"chg"`
	PChg       *decimal.Decimal `json:"p_chg"`
	Product    string           `json:"product"`
	Indicators []*AlertItem     `json:"indicators"`
}

// AlertItem is a single price-alert configuration.
type AlertItem struct {
	ID          string          `json:"id"`
	IndicatorID string          `json:"indicator_id"`
	Enabled     bool            `json:"enabled"`
	Frequency   int             `json:"frequency"`
	Scope       int             `json:"scope"`
	Text        string          `json:"text"`
	State       []int           `json:"state"`
	ValueMap    json.RawMessage `json:"value_map"`
}
