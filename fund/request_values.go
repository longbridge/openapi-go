package fund

import "net/url"

// Values builds the query string for the fund analysis / trend / comparison
// endpoints.
func (r *GetFundAnalysis) Values() url.Values {
	p := newParams()
	if r != nil {
		p.AddOptInt("period", int64(r.Period))
	}
	return p.Values()
}

// Values builds the query string for a paged endpoint.
func (r *FundPage) Values() url.Values {
	p := newParams()
	if r != nil {
		p.AddOptInt("page", int64(r.Page))
		p.AddOptInt("size", int64(r.Size))
	}
	return p.Values()
}

// Values builds the query string for a net-value range endpoint.
func (r *FundNavRange) Values() url.Values {
	p := newParams()
	if r != nil {
		p.AddOptInt("month_before", int64(r.MonthBefore))
		p.AddOptInt("year_before", int64(r.YearBefore))
	}
	return p.Values()
}

// Values builds the query string for the fund holdings endpoint.
func (r *GetFundHoldings) Values() url.Values {
	p := newParams()
	if r != nil {
		p.AddOptInt("scene", int64(r.Scene))
	}
	return p.Values()
}

// Values builds the query string for the fund stock-holdings endpoint.
func (r *GetFundStockHoldings) Values() url.Values {
	p := newParams()
	if r != nil {
		p.AddOptInt("limit", int64(r.Limit))
	}
	return p.Values()
}

// Values builds the query string for the fund positions overview endpoint.
func (r *GetFundPositions) Values() url.Values {
	p := newParams()
	if r != nil {
		p.Add("account_channel", r.AccountChannel)
		p.AddOptInt("aaid", r.Aaid)
	}
	return p.Values()
}

// Values builds the query string for a single fund position detail.
func (r *GetFundPosition) Values() url.Values {
	p := newParams()
	if r != nil {
		p.Add("account_channel", r.AccountChannel)
		p.AddOptInt("aaid", r.Aaid)
		p.Add("start", r.Start)
		p.Add("end", r.End)
	}
	return p.Values()
}

// Values builds the query string for a single fund position profit list.
func (r *GetFundPositionProfits) Values() url.Values {
	p := newParams()
	if r != nil {
		p.Add("account_channel", r.AccountChannel)
		p.AddOptInt("aaid", r.Aaid)
		p.Add("start", r.Start)
		p.Add("end", r.End)
		p.AddOptInt("page", int64(r.Page))
		p.AddOptInt("size", int64(r.Size))
	}
	return p.Values()
}

// Values builds the query string for a single fund position dividend list.
func (r *GetFundPositionDividends) Values() url.Values {
	p := newParams()
	if r != nil {
		p.Add("account_channel", r.AccountChannel)
		p.AddOptInt("aaid", r.Aaid)
		p.Add("currency", r.Currency)
		p.AddOptInt("start", r.Start)
		p.AddOptInt("end", r.End)
		p.AddOptInt("page", int64(r.Page))
		p.AddOptInt("size", int64(r.Size))
	}
	return p.Values()
}

// Values builds the query string for the fund orders list.
func (r *GetFundOrders) Values() url.Values {
	p := newParams()
	if r != nil {
		p.AddMulti("symbol", r.Symbols)
		p.Add("actions", r.Actions)
		p.Add("states", r.States)
		p.Add("currency", r.Currency)
		p.AddOptInt("start", r.Start)
		p.AddOptInt("end", r.End)
		p.AddOptInt("page", int64(r.Page))
		p.AddOptInt("size", int64(r.Size))
	}
	return p.Values()
}

// Values builds the query string for the fund transactions list.
func (r *GetFundTransactions) Values() url.Values {
	p := newParams()
	if r != nil {
		p.Add("account_channel", r.AccountChannel)
		p.Add("business_type", r.BusinessType)
		p.Add("category", r.Category)
		p.Add("currencies", r.Currencies)
		p.AddOptInt("start", r.Start)
		p.AddOptInt("end", r.End)
		p.AddOptInt("page", int64(r.Page))
		p.AddOptInt("size", int64(r.Size))
	}
	return p.Values()
}
