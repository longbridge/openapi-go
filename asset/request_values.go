package asset

import (
	"net/url"
	"strconv"
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
