package http

import (
	"context"
	"errors"
	"io"
	nhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestClient(t *testing.T, url string) *Client {
	t.Helper()
	c, err := New(
		WithURL(url),
		WithAppKey("test_app_key"),
		WithAppSecret("test_app_secret"),
		WithAccessToken("test_access_token"),
	)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	return c
}

func TestCallSSESuccess(t *testing.T) {
	srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
		if got := r.Header.Get("accept"); got != "text/event-stream" {
			t.Errorf("Accept header = %q, want text/event-stream", got)
		}
		if got := r.Header.Get("x-api-key"); got != "test_app_key" {
			t.Errorf("x-api-key header = %q", got)
		}
		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("event: message\ndata: {\"hello\":\"world\"}\n\n"))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	rc, err := c.CallSSE(context.Background(), "POST", "/v1/stream", map[string]string{"query": "hi"})
	if err != nil {
		t.Fatalf("CallSSE error: %v", err)
	}
	defer rc.Close()

	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if !strings.Contains(string(b), `{"hello":"world"}`) {
		t.Errorf("body = %s", b)
	}
}

func TestCallSSEErrorResponse(t *testing.T) {
	srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"code":123,"message":"bad request"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := c.CallSSE(context.Background(), "POST", "/v1/stream", map[string]string{"query": "hi"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("expected *ApiError, got %T: %v", err, err)
	}
	if apiErr.Code != 123 || apiErr.Message != "bad request" {
		t.Errorf("unexpected ApiError: %+v", apiErr)
	}
}

func TestCallWithRequestTimeout(t *testing.T) {
	srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"code":0,"message":"success","data":{}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	var resp map[string]interface{}
	err := c.Call(context.Background(), "GET", "/v1/slow", nil, nil, &resp, WithRequestTimeout(20*time.Millisecond))
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded, got: %v", err)
	}
}

func TestCallSuccess(t *testing.T) {
	srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"code":0,"message":"success","data":{"ok":true}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	var resp struct {
		OK bool `json:"ok"`
	}
	if err := c.Call(context.Background(), "GET", "/v1/ok", nil, nil, &resp); err != nil {
		t.Fatalf("Call error: %v", err)
	}
	if !resp.OK {
		t.Errorf("resp.OK = false, want true")
	}
}

func TestAPIResponseUnmarshalCode(t *testing.T) {
	cases := []struct {
		body string
		want int
	}{
		{`{"code":0,"message":"ok"}`, 0},
		{`{"code":301600,"message":"bad"}`, 301600},
		{`{"code":"459101","message":"request too many"}`, 459101},
		{`{"code":"abc"}`, 500},
		{`{"code":1.5}`, 500},
		{`{"code":true}`, 500},
		{`{"code":null}`, 0},
		{`{"message":"no code"}`, 0},
	}
	for _, tc := range cases {
		var r apiResponse
		if err := jsonUnmarshal(strings.NewReader(tc.body), &r); err != nil {
			t.Fatalf("%s: unmarshal error: %v", tc.body, err)
		}
		if r.Code != tc.want {
			t.Errorf("%s: code = %d, want %d", tc.body, r.Code, tc.want)
		}
	}

	var r apiResponse
	if err := jsonUnmarshal(strings.NewReader(`{"code":"1","message":"m","data":{"a":1}}`), &r); err != nil {
		t.Fatal(err)
	}
	if r.Message != "m" || string(r.Data) != `{"a":1}` {
		t.Errorf("other fields not decoded: %+v", r)
	}
}

func rateLimitedHandler(limited int, status int, ct, body string, calls *int) nhttp.HandlerFunc {
	return func(w nhttp.ResponseWriter, r *nhttp.Request) {
		*calls++
		if *calls <= limited {
			if ct != "" {
				w.Header().Set("content-type", ct)
			}
			w.WriteHeader(status)
			_, _ = io.WriteString(w, body)
			return
		}
		w.Header().Set("content-type", "application/json")
		_, _ = io.WriteString(w, `{"code":0,"data":{"ok":true}}`)
	}
}

func TestCallRetriesRateLimited(t *testing.T) {
	cases := []struct {
		name   string
		status int
		ct     string
		body   string
	}{
		{"gateway 429", nhttp.StatusTooManyRequests, "application/json", `{"code":"459101","message":"request too many"}`},
		{"429 only", nhttp.StatusTooManyRequests, "application/json", `{"code":1,"message":"slow down"}`},
		{"429 empty json body", nhttp.StatusTooManyRequests, "application/json", ``},
		{"429 truncated json", nhttp.StatusTooManyRequests, "application/json", `{"code":`},
		{"429 mistyped json", nhttp.StatusTooManyRequests, "application/json", `{"message":42}`},
		{"429 plain text", nhttp.StatusTooManyRequests, "text/plain", `too many`},
		{"code 459101 only", nhttp.StatusOK, "application/json", `{"code":"459101","message":"request too many"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			srv := httptest.NewServer(rateLimitedHandler(1, tc.status, tc.ct, tc.body, &calls))
			defer srv.Close()

			var resp struct{ Ok bool }
			if err := newTestClient(t, srv.URL).Call(context.Background(), "POST", "/v1/x", nil, map[string]int{"n": 1}, &resp); err != nil {
				t.Fatalf("Call error: %v", err)
			}
			if calls != 2 || !resp.Ok {
				t.Errorf("calls = %d, resp = %+v; want 2 calls and ok", calls, resp)
			}
		})
	}
}

func TestCallNoRetryOnOtherError(t *testing.T) {
	cases := []struct {
		name     string
		status   int
		body     string
		wantCode int
	}{
		{"business error", nhttp.StatusOK, `{"code":301600,"message":"bad"}`, 301600},
		{"server error", nhttp.StatusInternalServerError, `{"code":`, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
				calls++
				w.Header().Set("content-type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer srv.Close()

			err := newTestClient(t, srv.URL).Call(context.Background(), "GET", "/v1/x", nil, nil, nil)
			var ae *ApiError
			if !errors.As(err, &ae) || ae.HttpStatus != tc.status || ae.Code != tc.wantCode {
				t.Fatalf("err = %v, want ApiError status %d code %d", err, tc.status, tc.wantCode)
			}
			if calls != 1 {
				t.Errorf("calls = %d, want 1", calls)
			}
		})
	}
}

func TestCallMalformedErrorBodyKeepsStatus(t *testing.T) {
	srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("x-trace-id", "trace-1")
		w.WriteHeader(nhttp.StatusBadGateway)
		_, _ = io.WriteString(w, `<html>bad gateway</html>`)
	}))
	defer srv.Close()

	err := newTestClient(t, srv.URL).Call(context.Background(), "GET", "/v1/x", nil, nil, nil)
	var ae *ApiError
	if !errors.As(err, &ae) {
		t.Fatalf("err = %v, want ApiError", err)
	}
	if ae.HttpStatus != nhttp.StatusBadGateway || ae.Message != `<html>bad gateway</html>` || ae.TraceID != "trace-1" {
		t.Errorf("ApiError = %+v", ae)
	}
}

// cancelOnClose cancels a context once a response body is closed, i.e.
// right after Call has read and decoded a response, before any backoff.
type cancelOnClose struct {
	cancel context.CancelFunc
}

func (t cancelOnClose) RoundTrip(r *nhttp.Request) (*nhttp.Response, error) {
	resp, err := nhttp.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	resp.Body = cancelBody{ReadCloser: resp.Body, cancel: t.cancel}
	return resp, nil
}

type cancelBody struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (b cancelBody) Close() error {
	b.cancel()
	return b.ReadCloser.Close()
}

func TestCallRetryStopsOnContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	calls := 0
	srv := httptest.NewServer(rateLimitedHandler(100, nhttp.StatusTooManyRequests, "application/json", `{"code":"459101","message":"request too many"}`, &calls))
	defer srv.Close()

	c, err := New(
		WithURL(srv.URL),
		WithAppKey("test_app_key"),
		WithAppSecret("test_app_secret"),
		WithAccessToken("test_access_token"),
		WithClient(&nhttp.Client{Transport: cancelOnClose{cancel: cancel}}),
	)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	err = c.Call(ctx, "GET", "/v1/x", nil, nil, nil)
	var ae *ApiError
	if !errors.As(err, &ae) || ae.Code != CodeGatewayRateLimited || ae.HttpStatus != nhttp.StatusTooManyRequests {
		t.Fatalf("err = %v, want rate-limited ApiError", err)
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}
}

func TestCallNoRetryWhenContextAlreadyDone(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(rateLimitedHandler(100, nhttp.StatusTooManyRequests, "application/json", `{"code":"459101"}`, &calls))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := newTestClient(t, srv.URL).Call(ctx, "GET", "/v1/x", nil, nil, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if calls != 0 {
		t.Errorf("calls = %d, want 0", calls)
	}
}

func TestParseRetryAfter(t *testing.T) {
	cases := []struct {
		in   string
		want time.Duration
	}{
		{"", 0},
		{"2", 2 * time.Second},
		{"1.5s", 1500 * time.Millisecond},
		{"300ms", 300 * time.Millisecond},
		{"1m0s", time.Minute},
		{"-1s", 0},
		{"soon", 0},
		{time.Now().Add(-time.Minute).UTC().Format(nhttp.TimeFormat), 0},
	}
	for _, tc := range cases {
		if got := parseRetryAfter(tc.in); got != tc.want {
			t.Errorf("parseRetryAfter(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}

	if got := parseRetryAfter(time.Now().Add(time.Hour).UTC().Format(nhttp.TimeFormat)); got < 59*time.Minute || got > time.Hour {
		t.Errorf("parseRetryAfter(HTTP date in 1h) = %v", got)
	}
}

func retryAfterHandler(retryAfter string, calls *int) nhttp.HandlerFunc {
	return func(w nhttp.ResponseWriter, r *nhttp.Request) {
		*calls++
		w.Header().Set("content-type", "application/json")
		if *calls == 1 {
			w.Header().Set("retry-after", retryAfter)
			w.WriteHeader(nhttp.StatusTooManyRequests)
			_, _ = io.WriteString(w, `{"code":429002,"message":"rate limit"}`)
			return
		}
		_, _ = io.WriteString(w, `{"code":0,"data":{}}`)
	}
}

func TestCallRetryHonorsRetryAfter(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(retryAfterHandler("300ms", &calls))
	defer srv.Close()

	start := time.Now()
	if err := newTestClient(t, srv.URL).Call(context.Background(), "GET", "/v1/x", nil, nil, nil); err != nil {
		t.Fatalf("Call error: %v", err)
	}
	if calls != 2 {
		t.Errorf("calls = %d, want 2", calls)
	}
	// The backoff alone would wait 100ms.
	if elapsed := time.Since(start); elapsed < 300*time.Millisecond {
		t.Errorf("elapsed = %v, want >= 300ms", elapsed)
	}
}

func TestCallNoRetryWhenRetryAfterTooLong(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(retryAfterHandler("1m0s", &calls))
	defer srv.Close()

	err := newTestClient(t, srv.URL).Call(context.Background(), "GET", "/v1/x", nil, nil, nil)
	var ae *ApiError
	if !errors.As(err, &ae) || ae.HttpStatus != nhttp.StatusTooManyRequests || ae.RetryAfter != time.Minute {
		t.Fatalf("err = %v, want 429 ApiError with RetryAfter 1m", err)
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}
}
