package grid

import (
	"net/url"
	"strconv"
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

func (p params) AddOptInt(key string, val int64) {
	if val != 0 {
		p.AddInt(key, val)
	}
}

func (p params) Values() url.Values {
	vals := url.Values{}
	for k, v := range p {
		vals.Add(k, v)
	}
	return vals
}
