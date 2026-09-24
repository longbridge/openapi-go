package signal

import (
	"encoding/json"
	"testing"

	"github.com/longbridge/openapi-go/signal/jsontypes"
)

func TestConvertSignal(t *testing.T) {
	const raw = `{
		"id": "sign_1",
		"symbol": "992.HK",
		"outlook": "Strong bullish",
		"status": 1,
		"analysis_price": 12.5,
		"json_data": "{\"k\":1}",
		"created_at": 1783674041337,
		"updated_at": 1783674041337
	}`
	var j jsontypes.Signal
	if err := json.Unmarshal([]byte(raw), &j); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	s := convertSignal(&j)
	if s.Outlook != OutlookStrongBullish {
		t.Fatalf("outlook = %q", s.Outlook)
	}
	if s.Status != SignalStatusActive {
		t.Fatalf("status = %d", s.Status)
	}
	if s.AnalysisPrice != 12.5 {
		t.Fatalf("analysis_price = %v", s.AnalysisPrice)
	}
	if s.CreatedAt.IsZero() {
		t.Fatal("created_at not parsed")
	}
	// 1783674041337 ms = 2026-07-21T... — must land in 2026, proving ms (not s) handling.
	if got := s.CreatedAt.UTC().Year(); got != 2026 {
		t.Fatalf("created_at year = %d (ms not decoded as milliseconds?)", got)
	}
	if s.JsonData != `{"k":1}` {
		t.Fatalf("json_data = %q", s.JsonData)
	}
}

func TestConvertSecurityFact(t *testing.T) {
	const raw = `{
		"fact_id": "technical_rsi_14_short_1",
		"fact_type": "Technical",
		"direction": "short",
		"occur_time": "2026-07-21T10:00:00Z",
		"symbols_info": [{"symbol": "AAPL.US", "security_name": "Apple"}],
		"factors": [{
			"name": "rsi_14",
			"factor_groups": ["MOMENTUM"],
			"long_short_direction": "short",
			"anomaly_detection": {"anomaly_result": "hit", "thresholds": {"low": "1", "high": "3"}}
		}],
		"data_source": [{"source_name": "Nasdaq", "type": "Technical"}],
		"nl_info": {
			"title": "RSI overbought",
			"summary": "[{\"tag\":\"RSI\",\"value\":\"14 is high\"}]"
		}
	}`
	var j jsontypes.SecurityFact
	if err := json.Unmarshal([]byte(raw), &j); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	f := convertSecurityFact(&j)
	if f.FactType != FactTypeTechnical || f.Direction != FactDirectionShort {
		t.Fatalf("enums: type=%q dir=%q", f.FactType, f.Direction)
	}
	if f.OccurTime.IsZero() || f.OccurTime.UTC().Year() != 2026 {
		t.Fatalf("occur_time not parsed: %v", f.OccurTime)
	}
	if len(f.Factors) != 1 || f.Factors[0].LongShortDirection != FactDirectionShort {
		t.Fatalf("factor dir: %+v", f.Factors)
	}
	if f.Factors[0].AnomalyDetection.Thresholds.High != "3" {
		t.Fatalf("nested threshold: %+v", f.Factors[0].AnomalyDetection)
	}
	if len(f.DataSource) != 1 || f.DataSource[0].SourceType != FactTypeTechnical {
		t.Fatalf("data source: %+v", f.DataSource)
	}
	tags := f.NlInfo.SummaryTags()
	if len(tags) != 1 || tags[0].Tag != "RSI" || tags[0].Value != "14 is high" {
		t.Fatalf("summary tags: %+v", tags)
	}
	// A non-JSON field yields an empty slice, not a panic.
	if got := f.NlInfo.InvestAnalTags(); got != nil {
		t.Fatalf("empty invest_anal should yield nil, got %+v", got)
	}
}

func TestSignalsOptionsValues(t *testing.T) {
	limit := int32(50)
	o := &SignalsOptions{SymbolName: "AAPL.US", StrategyID: "buffett-value", Limit: &limit}
	v := o.Values()
	if v.Get("symbol_name") != "AAPL.US" || v.Get("strategy_id") != "buffett-value" || v.Get("limit") != "50" {
		t.Fatalf("values = %v", v)
	}
	if _, ok := v["offset"]; ok {
		t.Fatalf("offset should be omitted when nil: %v", v)
	}
}
