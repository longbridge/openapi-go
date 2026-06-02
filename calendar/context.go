package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/longbridge/openapi-go/calendar/jsontypes"
	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
)

// CalendarContext is a client for the Longbridge Financial Calendar OpenAPI.
//
// Example:
//
//	conf, err := config.NewFromEnv()
//	cctx, err := calendar.NewFromCfg(conf)
//	resp, err := cctx.FinanceCalendar(ctx, calendar.CalendarCategoryReport, "2025-05-01", "2025-05-31", nil)
type CalendarContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a CalendarContext from a *config.Config.
func NewFromCfg(cfg *config.Config) (*CalendarContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &CalendarContext{httpClient: httpClient}, nil
}

// NewFromEnv returns a CalendarContext configured from environment variables.
func NewFromEnv() (*CalendarContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// FinanceCalendar retrieves financial calendar events for a date range.
//
// category selects the event type (earnings, dividend, split, IPO, etc.).
// start and end are date strings in "YYYY-MM-DD" format.
// market is optional; pass nil or an empty string to retrieve all markets.
//
// The endpoint is paginated via next_date; this method follows the cursor
// automatically and returns the complete merged result.
//
// Reference: GET /v1/quote/finance_calendar
func (c *CalendarContext) FinanceCalendar(
	ctx context.Context,
	category CalendarCategory,
	start string,
	end string,
	market *string,
) (*CalendarEventsResponse, error) {
	cursor := start
	// seen deduplicates events across pages by ID.
	seen := make(map[string]struct{})
	var merged *CalendarEventsResponse

	for {
		params := url.Values{}
		params.Set("date", cursor)
		params.Set("date_end", end)
		params.Set("types[]", category.String())
		if market != nil && *market != "" {
			params.Set("markets[]", *market)
		}

		var raw jsontypes.CalendarEventsResponse
		if err := c.httpClient.Get(ctx, "/v1/quote/finance_calendar", params, &raw); err != nil {
			return nil, fmt.Errorf("calendar: finance_calendar: %w", err)
		}

		page, err := convertCalendarEventsResponse(&raw)
		if err != nil {
			return nil, err
		}

		if merged == nil {
			merged = &CalendarEventsResponse{Date: page.Date}
		}

		for _, grp := range page.List {
			var deduped []CalendarEventInfo
			for _, info := range grp.Infos {
				if _, dup := seen[info.ID]; !dup {
					seen[info.ID] = struct{}{}
					deduped = append(deduped, info)
				}
			}
			if len(deduped) > 0 {
				merged.List = append(merged.List, CalendarDateGroup{
					Date:  grp.Date,
					Count: int32(len(deduped)),
					Infos: deduped,
				})
			}
		}

		if raw.NextDate == "" {
			break
		}
		cursor = raw.NextDate
	}

	if merged == nil {
		merged = &CalendarEventsResponse{}
	}
	return merged, nil
}

// --- internal converters ---

func convertCalendarEventsResponse(raw *jsontypes.CalendarEventsResponse) (*CalendarEventsResponse, error) {
	resp := &CalendarEventsResponse{
		Date: raw.Date,
		List: make([]CalendarDateGroup, 0, len(raw.List)),
	}
	for _, rg := range raw.List {
		grp, err := convertCalendarDateGroup(&rg)
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, *grp)
	}
	return resp, nil
}

func convertCalendarDateGroup(raw *jsontypes.CalendarDateGroup) (*CalendarDateGroup, error) {
	grp := &CalendarDateGroup{
		Date:  raw.Date,
		Count: raw.Count,
		Infos: make([]CalendarEventInfo, 0, len(raw.Infos)),
	}
	for _, ri := range raw.Infos {
		info, err := convertCalendarEventInfo(&ri)
		if err != nil {
			return nil, err
		}
		grp.Infos = append(grp.Infos, *info)
	}
	return grp, nil
}

func convertCalendarEventInfo(raw *jsontypes.CalendarEventInfo) (*CalendarEventInfo, error) {
	info := &CalendarEventInfo{
		Symbol:              raw.Symbol,
		Market:              raw.Market,
		Content:             raw.Content,
		CounterName:         raw.CounterName,
		DateType:            raw.DateType,
		Date:                raw.Date,
		ChartUID:            raw.ChartUID,
		EventType:           raw.EventType,
		Datetime:            raw.Datetime,
		Icon:                raw.Icon,
		Star:                raw.Star,
		ID:                  raw.ID,
		FinancialMarketTime: raw.FinancialMarketTime,
		Currency:            raw.Currency,
		ActivityType:        raw.ActivityType,
	}

	// Preserve raw JSON blobs for Live and Ext
	if raw.Live != nil {
		b, err := raw.Live.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("calendar: marshal live: %w", err)
		}
		rm := json.RawMessage(b)
		info.Live = &rm
	}
	if raw.Ext != nil {
		b, err := raw.Ext.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("calendar: marshal ext: %w", err)
		}
		rm := json.RawMessage(b)
		info.Ext = &rm
	}

	info.DataKV = make([]CalendarDataKv, 0, len(raw.DataKV))
	for _, rkv := range raw.DataKV {
		kv, err := convertCalendarDataKv(&rkv)
		if err != nil {
			return nil, err
		}
		info.DataKV = append(info.DataKV, *kv)
	}
	return info, nil
}

func convertCalendarDataKv(raw *jsontypes.CalendarDataKv) (*CalendarDataKv, error) {
	kv := &CalendarDataKv{
		Key:       raw.Key,
		Value:     raw.Value,
		ValueType: raw.ValueType,
	}
	if raw.ValueRaw != "" {
		d, err := decimal.NewFromString(raw.ValueRaw)
		if err == nil {
			kv.ValueRaw = &d
		}
		// non-numeric value_raw strings are silently ignored (matches Rust behaviour)
	}
	return kv, nil
}
