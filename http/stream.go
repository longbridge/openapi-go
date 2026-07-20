package http

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"iter"
	nhttp "net/http"

	"github.com/longbridge/openapi-go/log"
)

// Stream represents a server-sent events (SSE) response. It is a pull-based
// reader in the style of bufio.Scanner / sql.Rows: no background goroutine,
// events are read from the connection on demand.
//
//	stream, err := client.Stream(ctx, "GET", "/v1/xxx", nil, nil)
//	if err != nil {
//		return err // request/connection errors and non-200 responses
//	}
//	defer stream.Close()
//
//	for stream.Next() {
//		handle(stream.Current())
//	}
//	if err := stream.Err(); err != nil {
//		// mid-stream read failure, timeout, or cancellation
//	}
//
// Or equivalently, range over Events():
//
//	for data := range stream.Events() {
//		handle(data)
//	}
//	if err := stream.Err(); err != nil { ... }
//
// A Stream is not safe for concurrent use. To abort a stream from another
// goroutine, cancel the context passed to Client.Stream.
type Stream struct {
	resp   *nhttp.Response
	reader *bufio.Reader
	cur    []byte
	err    error
	closed bool
}

// Stream sends a request expecting a server-sent events (SSE) response. It
// sets the "Accept: text/event-stream" request header and blocks until the
// response headers arrive. Request build errors, connection failures and
// non-200 responses (as *ApiError) are returned immediately; errors occurring
// after that — mid-stream read failures, timeouts, cancellation — are
// reported by Stream.Err.
//
// The stream's lifetime is bounded by ctx: the client-wide timeout does not
// apply, so pass a context with a deadline or cancellation to limit it.
// Always Close the returned stream.
func (c *Client) Stream(ctx context.Context, method, path string, queryParams interface{}, body interface{}, ropts ...RequestOption) (*Stream, error) {
	req, bb, err := c.newRequest(ctx, method, path, queryParams, body, ropts...)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "text/event-stream")

	log.Debugf("http stream method:%v url:%v body:%v", req.Method, req.URL, string(bb))
	resp, err := c.streamHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != nhttp.StatusOK {
		defer resp.Body.Close()
		rb, _ := io.ReadAll(resp.Body)
		apiResp := &apiResponse{TraceID: resp.Header.Get("x-trace-id")}
		if isJSON(resp.Header.Get("content-type")) {
			_ = jsonUnmarshal(bytes.NewReader(rb), apiResp)
		} else {
			apiResp.Message = string(rb)
		}
		return nil, NewError(resp.StatusCode, apiResp)
	}

	return &Stream{
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}, nil
}

// streamHTTPClient returns an *http.Client suitable for a long-lived stream.
// It is a shallow copy of the configured client with Timeout cleared, since a
// non-zero http.Client.Timeout caps the whole request (headers + body read)
// and would abort the stream. Callers control the lifetime via the context.
func (c *Client) streamHTTPClient() *nhttp.Client {
	cli := *c.httpClient
	cli.Timeout = 0
	return &cli
}

// Next reads the next event from the stream. It returns false when the
// stream ends — because the server closed it, an error occurred, or the
// stream was closed. After Next returns false, check Err.
func (s *Stream) Next() bool {
	if s.closed || s.err != nil {
		return false
	}

	var data bytes.Buffer
	for {
		line, err := s.reader.ReadBytes('\n')
		if len(line) > 0 {
			// strip the trailing line terminator (\n or \r\n)
			line = bytes.TrimRight(line, "\r\n")
			switch {
			case len(line) == 0:
				// blank line marks the end of an event
				if data.Len() > 0 {
					s.cur = data.Bytes()
					return true
				}
			case bytes.HasPrefix(line, []byte(":")):
				// comment / heartbeat, ignore
			case bytes.HasPrefix(line, []byte("data:")):
				val := bytes.TrimPrefix(line[len("data:"):], []byte(" "))
				if data.Len() > 0 {
					data.WriteByte('\n')
				}
				data.Write(val)
			default:
				// other fields (event/id/retry) are not part of the payload
			}
		}

		if err != nil {
			if err == io.EOF {
				// flush a trailing event that had no terminating blank line
				if data.Len() > 0 {
					s.cur = data.Bytes()
					return true
				}
			} else if !s.closed {
				s.err = err
			}
			return false
		}
	}
}

// Current returns the payload (the concatenated "data" fields) of the event
// read by the last successful call to Next.
func (s *Stream) Current() []byte {
	return s.cur
}

// Events returns an iterator over the remaining events of the stream,
// yielding each event's payload:
//
//	for data := range stream.Events() { ... }
//
// It is a convenience wrapper around Next/Current; after the loop, check Err.
func (s *Stream) Events() iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for s.Next() {
			if !yield(s.cur) {
				return
			}
		}
	}
}

// Err returns the error, if any, that terminated the stream. Mid-stream read
// failures, timeouts and cancellation are all reported here. It returns nil
// for a stream that ended normally or was closed by Close.
func (s *Stream) Err() error {
	return s.err
}

// Close releases the underlying connection. It is safe to call multiple
// times. Always Close a stream, even one that was fully consumed.
func (s *Stream) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	return s.resp.Body.Close()
}
