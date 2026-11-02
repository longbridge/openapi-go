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

func rateLimitedHandler(limited int, calls *int) nhttp.HandlerFunc {
	return func(w nhttp.ResponseWriter, r *nhttp.Request) {
		*calls++
		w.Header().Set("content-type", "application/json")
		if *calls <= limited {
			w.WriteHeader(nhttp.StatusTooManyRequests)
			_, _ = io.WriteString(w, `{"code":"459101","message":"request too many"}`)
			return
		}
		_, _ = io.WriteString(w, `{"code":0,"data":{"ok":true}}`)
	}
}

func TestCallRetriesRateLimited(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(rateLimitedHandler(1, &calls))
	defer srv.Close()

	var resp struct{ Ok bool }
	if err := newTestClient(t, srv.URL).Call(context.Background(), "POST", "/v1/x", nil, nil, &resp); err != nil {
		t.Fatalf("Call error: %v", err)
	}
	if calls != 2 || !resp.Ok {
		t.Errorf("calls = %d, resp = %+v; want 2 calls and ok", calls, resp)
	}
}

func TestCallNoRetryOnOtherError(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
		calls++
		w.Header().Set("content-type", "application/json")
		_, _ = io.WriteString(w, `{"code":301600,"message":"bad"}`)
	}))
	defer srv.Close()

	err := newTestClient(t, srv.URL).Call(context.Background(), "GET", "/v1/x", nil, nil, nil)
	var ae *ApiError
	if !errors.As(err, &ae) || ae.Code != 301600 {
		t.Fatalf("err = %v, want ApiError code 301600", err)
	}
	if calls != 1 {
		t.Errorf("calls = %d, want 1", calls)
	}
}

func TestCallRetryStopsOnContextDone(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(rateLimitedHandler(100, &calls))
	defer srv.Close()

	err := newTestClient(t, srv.URL).Call(context.Background(), "GET", "/v1/x", nil, nil, nil,
		WithRequestTimeout(250*time.Millisecond))
	var ae *ApiError
	if !errors.As(err, &ae) || ae.Code != CodeGatewayRateLimited || ae.HttpStatus != nhttp.StatusTooManyRequests {
		t.Fatalf("err = %v, want rate-limited ApiError", err)
	}
	// 0ms, 100ms, 300ms: the third attempt is past the 250ms deadline.
	if calls != 2 {
		t.Errorf("calls = %d, want 2", calls)
	}
}
