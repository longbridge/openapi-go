package trade

import (
	"net/url"
	"strconv"
	"time"

	"github.com/longbridge/openapi-go/internal/util"
	"github.com/shopspring/decimal"
)

type params map[string]string

func (p params) Add(key string, val string) {
	if len(val) > 0 {
		p[key] = val
	}
}

func (p params) AddInt(key string, val int64) {
	p[key] = strconv.FormatInt(val, 10)
}

func (p params) AddUInt(key string, val uint64) {
	p[key] = strconv.FormatUint(val, 10)
}

func (p params) AddDate(key string, val time.Time) {
	if !val.IsZero() {
		p[key] = util.FormatDate(&val)
	}
}

func (p params) AddOptInt(key string, val int64) {
	if val != 0 {
		p.AddInt(key, val)
	}
}

func (p params) AddPoniterInt(key string, val *int64) {
	if val != nil {
		p.AddInt(key, *val)
	}
}

func (p params) AddOptInt32(key string, val *int32) {
	if val != nil {
		p.AddInt(key, int64(*val))
	}
}

func (p params) AddOptUint(key string, val *uint64) {
	if val != nil {
		p[key] = strconv.FormatUint(*val, 10)
	}
}

func (p params) AddOptDecimal(key string, val decimal.Decimal) {
	if !val.IsZero() {
		p[key] = val.String()
	}
}

func (p params) Values() url.Values {
	vals := url.Values{}
	for k, v := range p {
		vals.Add(k, v)
	}
	return vals
}
