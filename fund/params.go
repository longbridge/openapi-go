package fund

import (
	"net/url"
	"strconv"
)

// params is a small query-string builder that omits empty values.
type params struct {
	vals url.Values
}

func newParams() *params {
	return &params{vals: url.Values{}}
}

// Add sets a single string value when it is non-empty.
func (p *params) Add(key string, val string) {
	if len(val) > 0 {
		p.vals.Set(key, val)
	}
}

// AddMulti appends each non-empty value under key (repeated query parameter).
func (p *params) AddMulti(key string, vals []string) {
	for _, v := range vals {
		if len(v) > 0 {
			p.vals.Add(key, v)
		}
	}
}

// AddOptInt sets an integer value when it is non-zero.
func (p *params) AddOptInt(key string, val int64) {
	if val != 0 {
		p.vals.Set(key, strconv.FormatInt(val, 10))
	}
}

// Values returns the accumulated url.Values.
func (p *params) Values() url.Values {
	return p.vals
}
