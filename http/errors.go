package http

import (
	"errors"
	"fmt"
	nhttp "net/http"
	"strconv"
	"strings"
	"time"
)

// CodeGatewayRateLimited is the code the unified gateway returns when it
// rejects a request for exceeding the rate limit.
const CodeGatewayRateLimited = 459101

// Retry policy for rate-limited calls, matching the Rust SDK: up to 5
// retries, backing off from 100ms and doubling each time. A Retry-After
// response header replaces the backoff delay; one longer than
// maxRetryAfter is not waited out, the error is returned instead.
const (
	retryCount        = 5
	retryInitialDelay = 100 * time.Millisecond
	retryFactor       = 2
	maxRetryAfter     = 10 * time.Second
)

type ApiError struct {
	HttpStatus int
	Code       int
	Message    string
	TraceID    string
	// RetryAfter is the wait the server asked for via the Retry-After
	// header, or 0 if it sent none.
	RetryAfter time.Duration
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
		RetryAfter: resp.RetryAfter,
	}
}

// rateLimitError returns err as an *ApiError if it is a rate-limit
// rejection, either an HTTP 429 or the gateway's CodeGatewayRateLimited,
// and nil otherwise.
func rateLimitError(err error) *ApiError {
	var ae *ApiError
	if !errors.As(err, &ae) {
		return nil
	}
	if ae.HttpStatus == nhttp.StatusTooManyRequests || ae.Code == CodeGatewayRateLimited {
		return ae
	}
	return nil
}

// parseRetryAfter parses a Retry-After header value: delay seconds or an
// HTTP date (RFC 9110), or a Go duration such as "1.5s", which is what the
// OpenAPI gateway sends. It returns 0 for an empty, invalid or past value.
func parseRetryAfter(v string) time.Duration {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	var d time.Duration
	if secs, err := strconv.Atoi(v); err == nil {
		d = time.Duration(secs) * time.Second
	} else if pd, err := time.ParseDuration(v); err == nil {
		d = pd
	} else if t, err := nhttp.ParseTime(v); err == nil {
		d = time.Until(t)
	}
	if d < 0 {
		return 0
	}
	return d
}
