// Package signal provides a client for the Longbridge strategy-signal OpenAPI:
// strategy signals and the catalyst (fact) events behind them.
package signal

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/signal/jsontypes"
)

// SignalContext is a client for the Longbridge strategy-signal OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	sctx, err := signal.NewFromCfg(conf)
//	resp, err := sctx.Signals(ctx, &signal.SignalsOptions{SymbolName: "AAPL.US"})
type SignalContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a SignalContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*SignalContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &SignalContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a SignalContext configured from environment variables.
func NewFromEnv() (*SignalContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// Signals queries signals, filtered by symbol, strategy, catalyst and time range.
//
// Path: GET /v1/signals
func (c *SignalContext) Signals(ctx context.Context, opts *SignalsOptions) (*SignalsResponse, error) {
	var raw jsontypes.SignalsResponse
	if err := c.httpClient.Get(ctx, "/v1/signals", opts.Values(), &raw); err != nil {
		return nil, fmt.Errorf("signal: signals: %w", err)
	}
	resp := &SignalsResponse{
		Total:   raw.Total,
		Signals: make([]*Signal, 0, len(raw.Signals)),
	}
	for i := range raw.Signals {
		resp.Signals = append(resp.Signals, convertSignal(&raw.Signals[i]))
	}
	return resp, nil
}

// Signal gets one signal by ID, including the full analysis in Signal.JsonData.
//
// Path: GET /v1/signals/{signal_id}
func (c *SignalContext) Signal(ctx context.Context, signalID string) (*Signal, error) {
	var raw jsontypes.SignalResponse
	path := "/v1/signals/" + url.PathEscape(signalID)
	if err := c.httpClient.Get(ctx, path, url.Values{}, &raw); err != nil {
		return nil, fmt.Errorf("signal: signal: %w", err)
	}
	return convertSignal(&raw.Signal), nil
}

// SecurityFacts lists the fact (catalyst) events for one security — anomaly
// detections, factor readings, data sources and natural-language summaries.
//
// Facts are what strategies react to: a signal names the fact that triggered it
// in Signal.KeyFactID.
//
// Path: GET /v1/facts/security_facts
func (c *SignalContext) SecurityFacts(ctx context.Context, opts *SecurityFactsOptions) ([]*SecurityFact, error) {
	var raw jsontypes.SecurityFactsResponse
	if err := c.httpClient.Get(ctx, "/v1/facts/security_facts", opts.Values(), &raw); err != nil {
		return nil, fmt.Errorf("signal: security_facts: %w", err)
	}
	facts := make([]*SecurityFact, 0, len(raw.Facts))
	for i := range raw.Facts {
		facts = append(facts, convertSecurityFact(&raw.Facts[i]))
	}
	return facts, nil
}

// --- internal converters ---

func convertSignal(j *jsontypes.Signal) *Signal {
	s := &Signal{}
	// util.Copy carries the scalar fields, the string-typed enums, and the
	// unix-millisecond CreatedAt/UpdatedAt -> time.Time conversion.
	_ = util.Copy(s, j)
	return s
}

func convertSecurityFact(j *jsontypes.SecurityFact) *SecurityFact {
	f := &SecurityFact{
		FactID:    j.FactID,
		FactType:  FactType(j.FactType),
		Direction: FactDirection(j.Direction),
		NlInfo: FactNlInfo{
			Title:      j.NlInfo.Title,
			SubTitle:   j.NlInfo.SubTitle,
			Summary:    j.NlInfo.Summary,
			InvestAnal: j.NlInfo.InvestAnal,
			EliExplain: j.NlInfo.EliExplain,
		},
	}
	// occur_time is an RFC3339 string, which util.Copy cannot map onto time.Time.
	if t, err := time.Parse(time.RFC3339, j.OccurTime); err == nil {
		f.OccurTime = t
	}
	// The nested slices carry no timestamps, so util.Copy handles them (including
	// the string-typed FactDirection / FactType enums) directly.
	_ = util.Copy(&f.SymbolsInfo, j.SymbolsInfo)
	_ = util.Copy(&f.Factors, j.Factors)
	_ = util.Copy(&f.DataSource, j.DataSource)
	return f
}
