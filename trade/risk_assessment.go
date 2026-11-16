package trade

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/longbridge/openapi-go/trade/jsontypes"
)

// RiskAssessment is the risk assessment (suitability questionnaire) result of
// the current member. It mirrors the server response of
// GET /v1/members/me/risk-assessment field-for-field.
type RiskAssessment struct {
	// Risk is the risk level of the latest questionnaire (e.g. "A5").
	Risk string `json:"risk"`
	// Pi is the professional investor (PI) status code.
	Pi int32 `json:"pi"`
	// Approvals are the approval states keyed by item (e.g. "car", "cka",
	// "risk_audit", "is_corporate_client"). Never nil.
	Approvals map[string]string `json:"approvals"`
	// ReAuthentication is the re-authentication flag, as returned by the
	// server (empty when not set).
	ReAuthentication string `json:"re_authentication"`
	// PiEndAt is the professional investor status end time, as returned by
	// the server (0 when not set).
	PiEndAt int64 `json:"pi_end_at"`
	// Deadlines are the deadlines keyed by item, as returned by the server.
	// Never nil.
	Deadlines map[string]int64 `json:"deadlines"`
	// RiskAssessmentAt is the date of the latest risk assessment, as returned
	// by the server (e.g. "2026-10-14").
	RiskAssessmentAt string `json:"risk_assessment_at"`
	// Score is the questionnaire score.
	Score int32 `json:"score"`
	// StrategyRiskLevels are the strategy risk levels the member is allowed to
	// trade.
	StrategyRiskLevels []string `json:"strategy_risk_levels"`
}

// UnmarshalJSON decodes the server response leniently, matching the Rust SDK:
// integers (including map values) are accepted as JSON numbers or numeric
// strings, and anything else (null, floats, out-of-range values) becomes 0;
// string fields, approval values and strategy risk levels keep a JSON string
// as-is, turn null into "", and render any other value as its compact JSON
// text. Null or missing maps decode to empty, non-nil maps.
func (r *RiskAssessment) UnmarshalJSON(data []byte) error {
	var raw struct {
		Risk               json.RawMessage                `json:"risk"`
		Pi                 jsontypes.FlexInt32            `json:"pi"`
		Approvals          map[string]json.RawMessage     `json:"approvals"`
		ReAuthentication   json.RawMessage                `json:"re_authentication"`
		PiEndAt            jsontypes.FlexInt64            `json:"pi_end_at"`
		Deadlines          map[string]jsontypes.FlexInt64 `json:"deadlines"`
		RiskAssessmentAt   json.RawMessage                `json:"risk_assessment_at"`
		Score              jsontypes.FlexInt32            `json:"score"`
		StrategyRiskLevels []json.RawMessage              `json:"strategy_risk_levels"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	approvals := make(map[string]string, len(raw.Approvals))
	for k, v := range raw.Approvals {
		approvals[k] = rawJSONText(v)
	}
	var strategyRiskLevels []string
	if raw.StrategyRiskLevels != nil {
		strategyRiskLevels = make([]string, len(raw.StrategyRiskLevels))
		for i, v := range raw.StrategyRiskLevels {
			strategyRiskLevels[i] = rawJSONText(v)
		}
	}
	deadlines := make(map[string]int64, len(raw.Deadlines))
	for k, v := range raw.Deadlines {
		deadlines[k] = int64(v)
	}
	*r = RiskAssessment{
		Risk:               rawJSONText(raw.Risk),
		Pi:                 int32(raw.Pi),
		Approvals:          approvals,
		ReAuthentication:   rawJSONText(raw.ReAuthentication),
		PiEndAt:            int64(raw.PiEndAt),
		Deadlines:          deadlines,
		RiskAssessmentAt:   rawJSONText(raw.RiskAssessmentAt),
		Score:              int32(raw.Score),
		StrategyRiskLevels: strategyRiskLevels,
	}
	return nil
}

// rawJSONText returns a JSON string's value, "" for null or an absent value,
// and the compacted JSON text of any other value.
func rawJSONText(v json.RawMessage) string {
	t := bytes.TrimSpace(v)
	if len(t) == 0 || bytes.Equal(t, []byte("null")) {
		return ""
	}
	var s string
	if err := json.Unmarshal(t, &s); err == nil {
		return s
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, t); err != nil {
		return string(t)
	}
	return buf.String()
}

// RiskAssessment returns the risk assessment (suitability questionnaire)
// result of the current member.
//
// Path: GET /v1/members/me/risk-assessment
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	tctx, err := trade.NewFromCfg(conf)
//	ra, err := tctx.RiskAssessment(context.Background())
func (c *TradeContext) RiskAssessment(ctx context.Context) (*RiskAssessment, error) {
	var resp RiskAssessment
	if err := c.opts.httpClient.Get(ctx, "/v1/members/me/risk-assessment", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
