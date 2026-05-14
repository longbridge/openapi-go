package calendar

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"github.com/longbridge/openapi-go/calendar/jsontypes"
	"github.com/longbridge/openapi-go/config"
	httplib "github.com/longbridge/openapi-go/http"
	"github.com/longbridge/openapi-go/internal/util"
)

// CalendarContext is a client for financial calendar events.
type CalendarContext struct {
	httpClient *httplib.Client
}

// NewFromCfg creates a CalendarContext from a Config.
func NewFromCfg(cfg *config.Config) (*CalendarContext, error) {
	httpClient, err := httplib.NewFromCfg(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "create http client error")
	}
	return &CalendarContext{httpClient: httpClient}, nil
}

// NewFromEnv creates a CalendarContext from environment variables.
func NewFromEnv() (*CalendarContext, error) {
	cfg, err := config.NewFormEnv()
	if err != nil {
		return nil, errors.Wrap(err, "load config from env error")
	}
	return NewFromCfg(cfg)
}

// FinanceCalendar returns financial calendar events for the given date range and category.
// start and end are date strings in YYYY-MM-DD format.
// market is optional (pass empty string to omit).
func (c *CalendarContext) FinanceCalendar(ctx context.Context, category CalendarCategory, start, end string, market string) (*CalendarEventsResponse, error) {
	values := url.Values{}
	values.Add("date", start)
	values.Add("date_end", end)
	values.Add("types[]", string(category))
	if market != "" {
		values.Add("markets[]", market)
	}

	var resp jsontypes.CalendarEventsResponse
	if err := c.httpClient.Get(ctx, "/v1/quote/finance_calendar", values, &resp); err != nil {
		return nil, err
	}

	result := &CalendarEventsResponse{Date: resp.Date}
	for _, g := range resp.List {
		group := &CalendarDateGroup{Date: g.Date, Count: g.Count}
		for _, info := range g.Infos {
			event := &CalendarEventInfo{
				Symbol:              util.CounterIDToSymbol(info.CounterId),
				Market:              info.Market,
				Content:             info.Content,
				CounterName:         info.CounterName,
				DateType:            info.DateType,
				Date:                info.Date,
				ChartUid:            info.ChartUid,
				EventType:           info.EventType,
				Datetime:            info.Datetime,
				Icon:                info.Icon,
				Star:                info.Star,
				Id:                  info.Id,
				FinancialMarketTime: info.FinancialMarketTime,
				Currency:            info.Currency,
				ActivityType:        info.ActivityType,
			}
			for _, kv := range info.DataKv {
				event.DataKv = append(event.DataKv, &CalendarDataKv{
					Key:       kv.Key,
					Value:     kv.Value,
					ValueType: kv.ValueType,
					ValueRaw:  kv.ValueRaw,
				})
			}
			group.Infos = append(group.Infos, event)
		}
		result.List = append(result.List, group)
	}
	return result, nil
}
