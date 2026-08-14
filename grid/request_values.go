package grid

import (
	"net/url"
)

func (r *GetGridOrders) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	p := &params{}
	p.AddOptInt("page", int64(r.Page))
	p.AddOptInt("limit", int64(r.Limit))
	p.Add("market", string(r.Market))
	p.Add("status", r.Status)
	p.Add("symbol", r.Symbol)
	p.Add("sort_by", r.SortBy)
	p.Add("sort_order", r.SortOrder)
	return p.Values()
}

func (r *GetGridOrderDetail) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("order_id", r.OrderId)
	p.Add("history_id", r.HistoryId)
	p.AddOptInt("limit", int64(r.Limit))
	return p.Values()
}

func (r *GetGridTriggerHistory) Values() url.Values {
	if r == nil {
		return url.Values{}
	}
	p := &params{}
	p.Add("grid_order_id", r.GridOrderId)
	p.AddOptInt("page", int64(r.Page))
	p.AddOptInt("limit", int64(r.Limit))
	return p.Values()
}
