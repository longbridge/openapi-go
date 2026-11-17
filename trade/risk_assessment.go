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
	// Risk is the risk level of the latest questionnaire (e.g. "A5"). The
	// server does not document the scale, so do not infer an ordering.
	Risk string `json:"risk"`
	// Pi is the professional investor (PI) status code. The server does not
	// document the code values; treat it as opaque.
	Pi int32 `json:"pi"`
	// Approvals are the approval states keyed by item (e.g. "car", "cka",
	// "risk_audit", "is_corporate_client"). Values are strings, not booleans:
	// "is_corporate_client" arrives as "true"/"false". Never nil.
	Approvals map[string]string `json:"approvals"`
	// ReAuthentication is the re-authentication flag, as returned by the
	// server (empty when not set).
	ReAuthentication string `json:"re_authentication"`
	// PiEndAt is the professional investor status end time, as returned by
	// the server (0 when not set). Raw integer: the unit (seconds or
	// milliseconds) is not documented, so it is not converted to a time.Time.
	PiEndAt int64 `json:"pi_end_at"`
	// Deadlines are the deadlines keyed by item, as returned by the server.
	// Raw integers with the same undocumented unit as PiEndAt. Never nil.
	Deadlines map[string]int64 `json:"deadlines"`
	// RiskAssessmentAt is the date of the latest risk assessment, as returned
	// by the server (e.g. "2026-10-14"). Kept as the server's string rather
	// than a time.Time: the "yyyy-mm-dd" form is not documented to be
	// guaranteed. May be empty; note that lenient decoding also renders a null
	// or non-string value as "", so emptiness does not by itself mean "never
	// assessed".
	RiskAssessmentAt string `json:"risk_assessment_at"`
	// Score is the questionnaire score.
	Score int32 `json:"score"`
	// StrategyRiskLevels are the strategy risk levels the member is allowed to
	// trade. Never nil.
	StrategyRiskLevels []string `json:"strategy_risk_levels"`
}

// UnmarshalJSON decodes the server response leniently, following the same rules
// as the Rust SDK: integers (including map values) are accepted as JSON
// integers or integer strings, and anything else (null, floats, out-of-range
// values) becomes 0; string fields, approval values and strategy risk levels
// keep a JSON string as-is, turn null into "", and render any other value as
// its compact JSON text (non-integer numbers keep their original text, e.g.
// "1.50", where Rust normalizes them); a map or list field holding anything
// other than an object / array (e.g. null, or [] for a map) decodes to an empty,
// non-nil value.
func (r *RiskAssessment) UnmarshalJSON(data []byte) error {
	var raw struct {
		Risk               json.RawMessage     `json:"risk"`
		Pi                 jsontypes.FlexInt32 `json:"pi"`
		Approvals          json.RawMessage     `json:"approvals"`
		ReAuthentication   json.RawMessage     `json:"re_authentication"`
		PiEndAt            jsontypes.FlexInt64 `json:"pi_end_at"`
		Deadlines          json.RawMessage     `json:"deadlines"`
		RiskAssessmentAt   json.RawMessage     `json:"risk_assessment_at"`
		Score              jsontypes.FlexInt32 `json:"score"`
		StrategyRiskLevels json.RawMessage     `json:"strategy_risk_levels"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	// A wrong-shaped value fails these inner decodes and leaves them empty.
	var approvalsRaw map[string]json.RawMessage
	_ = json.Unmarshal(raw.Approvals, &approvalsRaw)
	approvals := make(map[string]string, len(approvalsRaw))
	for k, v := range approvalsRaw {
		approvals[k] = rawJSONText(v)
	}
	var deadlinesRaw map[string]jsontypes.FlexInt64
	_ = json.Unmarshal(raw.Deadlines, &deadlinesRaw)
	deadlines := make(map[string]int64, len(deadlinesRaw))
	for k, v := range deadlinesRaw {
		deadlines[k] = int64(v)
	}
	var levelsRaw []json.RawMessage
	_ = json.Unmarshal(raw.StrategyRiskLevels, &levelsRaw)
	strategyRiskLevels := make([]string, len(levelsRaw))
	for i, v := range levelsRaw {
		strategyRiskLevels[i] = rawJSONText(v)
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
