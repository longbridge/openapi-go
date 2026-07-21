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
