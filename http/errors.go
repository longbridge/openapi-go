package http

import (
	"errors"
	"fmt"
	nhttp "net/http"
	"time"
)

// CodeGatewayRateLimited is the code the unified gateway returns when it
// rejects a request for exceeding the rate limit.
const CodeGatewayRateLimited = 459101

// Retry policy for rate-limited calls, matching the Rust SDK: up to 5
// retries, backing off from 100ms and doubling each time.
const (
	retryCount        = 5
	retryInitialDelay = 100 * time.Millisecond
	retryFactor       = 2
)

type ApiError struct {
	HttpStatus int
	Code       int
	Message    string
	TraceID    string
}

func (ae *ApiError) Error() string {
	return fmt.Sprintf("longbridge openapi error, httpStatus:%d code:%d message:%s trace:%s", ae.HttpStatus, ae.Code, ae.Message, ae.TraceID)
}

func NewError(httpStatus int, resp *apiResponse) error {
	return &ApiError{
		HttpStatus: httpStatus,
		Code:       resp.Code,
		Message:    resp.Message,
		TraceID:    resp.TraceID,
	}
}

// isRateLimited reports whether err is a rate-limit rejection, either an
// HTTP 429 or the gateway's CodeGatewayRateLimited.
func isRateLimited(err error) bool {
	var ae *ApiError
	if !errors.As(err, &ae) {
		return false
	}
	return ae.HttpStatus == nhttp.StatusTooManyRequests || ae.Code == CodeGatewayRateLimited
}
