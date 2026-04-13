package ipo

import (
	"net/url"
	"strconv"
)


// Values converts FetchOrderListRequest to URL query parameters.
func (req *FetchOrderListRequest) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	v := url.Values{}
	if req.Symbol != "" {
		v.Set("symbol", req.Symbol)
	}
	if req.Page != 0 {
		v.Set("page", strconv.Itoa(int(req.Page)))
	}
	if req.PageSize != 0 {
		v.Set("page_size", strconv.Itoa(int(req.PageSize)))
	}
	return v
}

// Values converts FetchOrderDetailRequest to URL query parameters.
func (req *FetchOrderDetailRequest) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	v := url.Values{}
	v.Set("order_id", strconv.FormatInt(req.OrderID, 10))
	return v
}

// Values converts FetchMarginListRequest to URL query parameters.
func (req *FetchMarginListRequest) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	v := url.Values{}
	v.Set("symbol", req.Symbol)
	return v
}

// Values converts FetchIpoPaymentListRequest to URL query parameters.
func (req *FetchIpoPaymentListRequest) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	v := url.Values{}
	v.Set("symbol", req.Symbol)
	return v
}

// Values converts FetchBuyLimitRequest to URL query parameters.
func (req *FetchBuyLimitRequest) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	v := url.Values{}
	v.Set("symbol", req.Symbol)
	v.Set("currency", req.Currency)
	return v
}
