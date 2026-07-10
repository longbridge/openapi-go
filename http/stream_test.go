package http

import (
	"context"
	"fmt"
	nhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(t *testing.T, url string) *Client {
	t.Helper()
	c, err := New(
		WithURL(url),
		WithAppKey("key"),
		WithAppSecret("secret"),
		WithAccessToken("token"),
	)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	return c
}

func sseTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
		if got := r.Header.Get("accept"); got != "text/event-stream" {
			t.Errorf("accept header = %q, want text/event-stream", got)
		}
		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(nhttp.StatusOK)
		fl := w.(nhttp.Flusher)
		// two data lines -> joined with \n
		fmt.Fprint(w, "data: hello\ndata: world\n\n")
		fl.Flush()
		fmt.Fprint(w, ": heartbeat\n\n") // comment-only event, skipped
		fl.Flush()
		fmt.Fprint(w, "event: msg\r\ndata: {\"a\":1}\r\n\r\n") // CRLF + non-data field
		fl.Flush()
		fmt.Fprint(w, "data: trailing") // no terminating blank line -> flushed on EOF
	}))
}

func checkEvents(t *testing.T, got []string) {
	t.Helper()
	want := []string{"hello\nworld", `{"a":1}`, "trailing"}
	if len(got) != len(want) {
		t.Fatalf("got %d events %q, want %d %q", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("event[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestStreamNext(t *testing.T) {
	srv := sseTestServer(t)
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	stream, err := c.Stream(context.Background(), "GET", "/v1/stream", nil, nil)
	if err != nil {
		t.Fatalf("Stream() error: %v", err)
	}
	defer stream.Close()

	var got []string
	for stream.Next() {
		got = append(got, string(stream.Current()))
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
	checkEvents(t, got)

	if stream.Next() {
		t.Error("Next() after exhaustion = true, want false")
	}
}

func TestStreamEventsIterator(t *testing.T) {
	srv := sseTestServer(t)
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	stream, err := c.Stream(context.Background(), "GET", "/v1/stream", nil, nil)
	if err != nil {
		t.Fatalf("Stream() error: %v", err)
	}
	defer stream.Close()

	var got []string
	for data := range stream.Events() {
		got = append(got, string(data))
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
	checkEvents(t, got)
}

func TestStreamEventsIteratorBreak(t *testing.T) {
	srv := sseTestServer(t)
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	stream, err := c.Stream(context.Background(), "GET", "/v1/stream", nil, nil)
	if err != nil {
		t.Fatalf("Stream() error: %v", err)
	}
	defer stream.Close()

	var got []string
	for data := range stream.Events() {
		got = append(got, string(data))
		break
	}
	if len(got) != 1 || got[0] != "hello\nworld" {
		t.Fatalf("got %q, want just the first event", got)
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
}

func TestStreamNon200(t *testing.T) {
	srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("x-trace-id", "trace-123")
		w.WriteHeader(nhttp.StatusBadRequest)
		fmt.Fprint(w, `{"code":40001,"message":"bad request"}`)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	stream, err := c.Stream(context.Background(), "GET", "/v1/stream", nil, nil)
	if err == nil {
		stream.Close()
		t.Fatal("Stream() error = nil, want *ApiError")
	}
	ae, ok := err.(*ApiError)
	if !ok {
		t.Fatalf("Stream() error type = %T, want *ApiError", err)
	}
	if ae.HttpStatus != nhttp.StatusBadRequest || ae.Code != 40001 || ae.TraceID != "trace-123" {
		t.Errorf("unexpected ApiError: %+v", ae)
	}
}

func TestStreamContextCancel(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {
		w.Header().Set("content-type", "text/event-stream")
		w.WriteHeader(nhttp.StatusOK)
		w.(nhttp.Flusher).Flush()
		<-release // hold the connection open until the test is done
	}))
	defer srv.Close()
	defer close(release)

	c := newTestClient(t, srv.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	stream, err := c.Stream(ctx, "GET", "/v1/stream", nil, nil)
	if err != nil {
		t.Fatalf("Stream() error: %v", err)
	}
	defer stream.Close()

	done := make(chan struct{})
	go func() {
		for stream.Next() {
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("stream did not terminate after context timeout")
	}

	if stream.Err() == nil {
		t.Fatal("Err() = nil, want context deadline / timeout error")
	}
}

func TestStreamConnectError(t *testing.T) {
	// point at a closed server so the request fails to connect
	srv := httptest.NewServer(nhttp.HandlerFunc(func(w nhttp.ResponseWriter, r *nhttp.Request) {}))
	url := srv.URL
	srv.Close()

	c := newTestClient(t, url)
	if _, err := c.Stream(context.Background(), "GET", "/v1/stream", nil, nil); err == nil {
		t.Fatal("Stream() error = nil, want connection error")
	}
}

func TestStreamClose(t *testing.T) {
	srv := sseTestServer(t)
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	stream, err := c.Stream(context.Background(), "GET", "/v1/stream", nil, nil)
	if err != nil {
		t.Fatalf("Stream() error: %v", err)
	}
	if err := stream.Close(); err != nil {
		t.Fatalf("Close() error: %v", err)
	}
	if err := stream.Close(); err != nil {
		t.Fatalf("second Close() error: %v", err)
	}
	if stream.Next() {
		t.Error("Next() after Close = true, want false")
	}
	if err := stream.Err(); err != nil {
		t.Errorf("Err() after Close = %v, want nil", err)
	}
}
