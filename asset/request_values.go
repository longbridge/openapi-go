package asset

import (
	"net/url"
	"strconv"
	"strings"
)

// Values converts GetStatementList to URL query parameters.
func (req *GetStatementList) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	v := url.Values{}
	if req.StatementType != 0 {
		v.Set("statement_type", strconv.Itoa(int(req.StatementType)))
	}
	if req.Page != 0 {
		v.Set("page", strconv.Itoa(int(req.Page)))
	}
	if req.PageSize != 0 {
		v.Set("page_size", strconv.Itoa(int(req.PageSize)))
	}
	return v
}

// Values converts GetStatementDownloadURL to URL query parameters.
func (req *GetStatementDownloadURL) Values() url.Values {
	if req == nil {
		return url.Values{}
	}
	v := url.Values{}
	v.Set("file_key", req.FileKey)
	return v
}

// Values converts GetStockTransfers to URL query parameters. Page defaults to
// 1 and Size to 20.
func (req *GetStockTransfers) Values() url.Values {
	v := url.Values{}
	if req == nil {
		setPaging(v, "page", 0, "size", 0)
		return v
	}
	setPaging(v, "page", req.Page, "size", req.Size)
	if m := strings.ToUpper(strings.TrimSpace(req.Market)); m != "" {
		v.Set("market", m)
	}
	if req.State != nil {
		v.Set("state", strconv.Itoa(int(*req.State)))
	}
	setUnix(v, "start", req.Start)
	setUnix(v, "end", req.End)
	return v
}

// Values converts GetRewardRecords to URL query parameters. Page defaults to
// 1 and PageSize to 20.
func (req *GetRewardRecords) Values() url.Values {
	v := url.Values{}
	if req == nil {
		setPaging(v, "current", 0, "page_size", 0)
		return v
	}
	setPaging(v, "current", req.Page, "page_size", req.PageSize)
	if req.RewardType != nil {
		v.Set("reward_type", strconv.Itoa(int(*req.RewardType)))
	}
	return v
}

// Values converts GetCoupons to URL query parameters. Status defaults to
// CouponStatusFilterValid, Page to 1 and PageSize to 20.
func (req *GetCoupons) Values() url.Values {
	v := url.Values{}
	if req == nil {
		req = &GetCoupons{}
	}
	setPaging(v, "current", req.Page, "page_size", req.PageSize)
	status := req.Status
	if status == 0 {
		status = CouponStatusFilterValid
	}
	v.Set("card_status", strconv.Itoa(int(status)))
	if t := strings.TrimSpace(req.CardType); t != "" {
		v.Set("card_type", t)
	}
	return v
}

// Values converts GetDeposits to URL query parameters. States are sent as
// comma-separated numeric codes. Page defaults to 1 and Size to 20.
func (req *GetDeposits) Values() url.Values {
	v := url.Values{}
	if req == nil {
		setPaging(v, "page", 0, "size", 0)
		return v
	}
	setPaging(v, "page", req.Page, "size", req.Size)
	if len(req.States) > 0 {
		states := make([]string, 0, len(req.States))
		for _, s := range req.States {
			states = append(states, strconv.Itoa(int(s)))
		}
		v.Set("states", strings.Join(states, ","))
	}
	if len(req.Currencies) > 0 {
		currencies := make([]string, 0, len(req.Currencies))
		for _, c := range req.Currencies {
			if c = strings.ToUpper(strings.TrimSpace(c)); c != "" {
				currencies = append(currencies, c)
			}
		}
		if len(currencies) > 0 {
			v.Set("currencies", strings.Join(currencies, ","))
		}
	}
	setUnix(v, "from_date", req.Start)
	setUnix(v, "to_date", req.End)
	return v
}

// Values converts GetWithdrawals to URL query parameters. Page defaults to 1
// and Size to 20.
func (req *GetWithdrawals) Values() url.Values {
	v := url.Values{}
	if req == nil {
		setPaging(v, "page", 0, "size", 0)
		return v
	}
	setPaging(v, "page", req.Page, "size", req.Size)
	setUnix(v, "begin", req.Start)
	setUnix(v, "end", req.End)
	return v
}
