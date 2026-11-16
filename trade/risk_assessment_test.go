package trade

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Real test-env response (2026-11-16).
const riskAssessmentTestEnvBody = `{"approvals":{"car":"pass","cka":"pass","is_corporate_client":"false","risk_audit":"pass"},` +
	`"deadlines":{},"pi":0,"pi_end_at":"0","re_authentication":"","risk":"A5",` +
	`"risk_assessment_at":"2026-10-14","score":0,"strategy_risk_levels":[]}`

func TestRiskAssessmentUnmarshal(t *testing.T) {
	cases := []struct {
		name string
		data string
		want RiskAssessment
	}{
		{
			name: "test env response",
			data: riskAssessmentTestEnvBody,
			want: RiskAssessment{
				Risk:               "A5",
				Approvals:          map[string]string{"car": "pass", "cka": "pass", "is_corporate_client": "false", "risk_audit": "pass"},
				Deadlines:          map[string]int64{},
				RiskAssessmentAt:   "2026-10-14",
				StrategyRiskLevels: []string{},
			},
		},
		{
			name: "ints as strings or numbers, non-string approvals",
			data: `{"risk":"R1","pi":"1","pi_end_at":1790000000,"score":"42",` +
				`"deadlines":{"pi":"1790000000","risk":1790000001,"x":null},` +
				`"approvals":{"flag":true,"n":3,"z":null},"strategy_risk_levels":null}`,
			want: RiskAssessment{
				Risk:               "R1",
				Pi:                 1,
				PiEndAt:            1790000000,
				Score:              42,
				Deadlines:          map[string]int64{"pi": 1790000000, "risk": 1790000001, "x": 0},
				Approvals:          map[string]string{"flag": "true", "n": "3", "z": ""},
				StrategyRiskLevels: []string{},
			},
		},
		{
			name: "floats and out-of-range ints become 0, non-string string fields",
			data: `{"risk":5,"re_authentication":true,"risk_assessment_at":null,"pi":1.5,"score":3000000000,` +
				`"pi_end_at":1.79e9,"deadlines":{"a":-2,"b":-2.7,"c":"x","d":18446744073709551615},` +
				`"approvals":{"o":{"b": 1, "a": 2}},"strategy_risk_levels":["S1",null,3]}`,
			want: RiskAssessment{
				Risk:               "5",
				ReAuthentication:   "true",
				Deadlines:          map[string]int64{"a": -2, "b": 0, "c": 0, "d": 0},
				Approvals:          map[string]string{"o": `{"b":1,"a":2}`},
				StrategyRiskLevels: []string{"S1", "", "3"},
			},
		},
		{
			name: "null maps",
			data: `{"approvals":null,"deadlines":null}`,
			want: RiskAssessment{Approvals: map[string]string{}, Deadlines: map[string]int64{}, StrategyRiskLevels: []string{}},
		},
		{
			name: "wrong-shaped maps and lists become empty",
			data: `{"approvals":[],"deadlines":"x","strategy_risk_levels":{"a":1}}`,
			want: RiskAssessment{Approvals: map[string]string{}, Deadlines: map[string]int64{}, StrategyRiskLevels: []string{}},
		},
		{
			name: "scalar strategy list",
			data: `{"strategy_risk_levels":"S1"}`,
			want: RiskAssessment{Approvals: map[string]string{}, Deadlines: map[string]int64{}, StrategyRiskLevels: []string{}},
		},
		{
			name: "empty object",
			data: `{}`,
			want: RiskAssessment{Approvals: map[string]string{}, Deadlines: map[string]int64{}, StrategyRiskLevels: []string{}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got RiskAssessment
			if err := json.Unmarshal([]byte(tc.data), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestRiskAssessmentRequest(t *testing.T) {
	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Method != nethttp.MethodGet || r.URL.Path != "/v1/members/me/risk-assessment" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"data":` + riskAssessmentTestEnvBody + `}`))
	}))
	defer srv.Close()

	tc := newTestTradeContext(t, srv.URL, "hk_")
	got, err := tc.RiskAssessment(context.Background())
	if err != nil {
		t.Fatalf("RiskAssessment: %v", err)
	}
	if got.Risk != "A5" || got.RiskAssessmentAt != "2026-10-14" || got.Approvals["risk_audit"] != "pass" || len(got.Approvals) != 4 {
		t.Errorf("unexpected result %+v", got)
	}
}
