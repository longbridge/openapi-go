package quote

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	quotev1 "github.com/longbridge/openapi-protobufs/gen/go/quote"
	protocol "github.com/longbridge/openapi-protocol/go"
	"github.com/longbridge/openapi-protocol/go/client"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Transport selects how the quote pull (request/response) APIs reach the
// server. Subscriptions, push events and the APIs without a REST equivalent
// always use the WebSocket connection.
type Transport string

const (
	// TransportWebSocket sends pull requests over the quote WebSocket
	// connection (default).
	TransportWebSocket Transport = "ws"
	// TransportHTTP sends pull requests over HTTP (`POST /quote/*`) where the
	// API has a REST equivalent, falling back to the WebSocket for the rest.
	// The WebSocket is only dialed when something actually needs it, so a
	// context that only uses HTTP-backed APIs never opens a connection.
	TransportHTTP Transport = "http"
)

// ParseTransport parses a transport name (`ws` or `http`, case-insensitive). Unknown names return false.
func ParseTransport(s string) (Transport, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "ws":
		return TransportWebSocket, true
	case "http":
		return TransportHTTP, true
	}
	return "", false
}

// httpPaths maps a WebSocket command to the REST endpoint that serves the same
// request and response messages.
var httpPaths = map[quotev1.Command]string{
	quotev1.Command_QueryMarketTradePeriod:       "/quote/markets/trading-sessions",
	quotev1.Command_QueryMarketTradeDay:          "/quote/markets/trading-days",
	quotev1.Command_QuerySecurityStaticInfo:      "/quote/static-info",
	quotev1.Command_QuerySecurityQuote:           "/quote/quotes",
	quotev1.Command_QueryOptionQuote:             "/quote/options/quotes",
	quotev1.Command_QueryWarrantQuote:            "/quote/warrants/quotes",
	quotev1.Command_QueryDepth:                   "/quote/depth",
	quotev1.Command_QueryBrokers:                 "/quote/brokers",
	quotev1.Command_QueryParticipantBrokerIds:    "/quote/participants",
	quotev1.Command_QueryTrade:                   "/quote/trades",
	quotev1.Command_QueryIntraday:                "/quote/intraday",
	quotev1.Command_QueryCandlestick:             "/quote/candlesticks",
	quotev1.Command_QueryOptionChainDate:         "/quote/options/expiry-dates",
	quotev1.Command_QueryWarrantIssuerInfo:       "/quote/warrants/issuers",
	quotev1.Command_QueryWarrantFilterList:       "/quote/warrants",
	quotev1.Command_QueryCapitalFlowIntraday:     "/quote/capital-flow",
	quotev1.Command_QueryCapitalFlowDistribution: "/quote/capital-distribution",
	quotev1.Command_QuerySecurityCalcIndex:       "/quote/calc-indexes",
	quotev1.Command_QueryHistoryCandlestick:      "/quote/history-candlesticks",
}

// ErrClosed is returned when the WebSocket is needed after Close.
var ErrClosed = errors.New("quote: context closed")

var (
	protoJSONMarshal   = protojson.MarshalOptions{UseProtoNames: true}
	protoJSONUnmarshal = protojson.UnmarshalOptions{DiscardUnknown: true}
)

// do sends `req` as command `cmd` and decodes the response into `ret`, over
// HTTP when the transport is HTTP and the command has a REST equivalent,
// otherwise over the WebSocket. `req` may be nil for body-less commands.
func (c *core) do(ctx context.Context, cmd quotev1.Command, req, ret proto.Message) error {
	if c.transport == TransportHTTP {
		if path, ok := httpPaths[cmd]; ok {
			return c.httpDo(ctx, path, req, ret)
		}
	}
	return c.wsDo(ctx, cmd, req, ret)
}

// wsDo sends a command over the WebSocket, dialing it first if needed.
func (c *core) wsDo(ctx context.Context, cmd quotev1.Command, req, ret proto.Message) error {
	cl, err := c.ws(ctx)
	if err != nil {
		return err
	}
	r := &client.Request{Cmd: uint32(cmd)}
	if req != nil {
		r.Body = req
	}
	var res *protocol.Packet
	if res, err = cl.Do(ctx, r); err != nil {
		return err
	}
	if ret == nil {
		return nil
	}
	return res.Unmarshal(ret)
}

// httpDo posts `req` as proto-JSON to `path` and decodes the proto-JSON
// response into `ret`.
func (c *core) httpDo(ctx context.Context, path string, req, ret proto.Message) error {
	body := json.RawMessage("{}")
	if req != nil {
		b, err := protoJSONMarshal.Marshal(req)
		if err != nil {
			return errors.Wrap(err, "encode quote request")
		}
		body = b
	}
	var raw json.RawMessage
	if err := c.httpClient.Post(ctx, path, body, &raw); err != nil {
		return err
	}
	// http.Client leaves an absent `data` as an empty RawMessage (see Call).
	// An absent or `null` `data` is an empty response.
	if ret == nil || len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil
	}
	normalized, err := normalizeProtoJSON(raw, ret.ProtoReflect().Descriptor())
	if err != nil {
		return errors.Wrapf(err, "decode quote response from %s", path)
	}
	if err := protoJSONUnmarshal.Unmarshal(normalized, ret); err != nil {
		return errors.Wrapf(err, "decode quote response from %s: %s", path, excerpt(raw))
	}
	return nil
}

// excerpt returns the start of a payload for error messages.
func excerpt(b []byte) string {
	const max = 256
	if len(b) <= max {
		return string(b)
	}
	return string(b[:max]) + "…"
}

// normalizeProtoJSON rewrites the few gateway encodings that protojson
// rejects, using the message descriptor to find the affected fields: an empty
// string on a numeric / enum field is dropped (→ default), a quoted number on
// an enum field is unquoted, `"true"` / `"false"` / `1` / `0` on a bool field
// become bools, a number on a string field becomes a string, and `null`
// elements are removed from lists. Everything else, including unknown fields
// and camelCase keys, is left to protojson.
//
// It deliberately decodes, walks and re-encodes the whole payload (about 3x
// the CPU of protojson alone on a 1000-row warrant page), which is small next
// to the HTTP round trip.
func normalizeProtoJSON(raw json.RawMessage, md protoreflect.MessageDescriptor) (json.RawMessage, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	normalizeValue(v, md)
	return json.Marshal(v)
}

func normalizeValue(v interface{}, md protoreflect.MessageDescriptor) {
	obj, ok := v.(map[string]interface{})
	if !ok {
		return
	}
	fields := md.Fields()
	for key, val := range obj {
		// Same lookup order as protojson: JSON name first, then proto name.
		fd := fields.ByJSONName(key)
		if fd == nil {
			fd = fields.ByName(protoreflect.Name(key))
		}
		if fd == nil {
			continue
		}
		switch {
		case fd.IsMap():
			continue
		case fd.IsList():
			list, ok := val.([]interface{})
			if !ok {
				continue
			}
			out := list[:0]
			for _, el := range list {
				if fd.Kind() == protoreflect.MessageKind {
					if el == nil {
						continue
					}
					normalizeValue(el, fd.Message())
				} else {
					var drop bool
					el, drop = normalizeScalar(el, fd.Kind())
					if drop {
						continue
					}
				}
				out = append(out, el)
			}
			obj[key] = out
		case fd.Kind() == protoreflect.MessageKind:
			normalizeValue(val, fd.Message())
		default:
			if nv, drop := normalizeScalar(val, fd.Kind()); drop {
				delete(obj, key)
			} else {
				obj[key] = nv
			}
		}
	}
}

// normalizeScalar coerces a scalar JSON value for a field of `kind`; `drop`
// asks the caller to omit the value (so protojson applies the default).
func normalizeScalar(v interface{}, kind protoreflect.Kind) (out interface{}, drop bool) {
	switch kind {
	case protoreflect.BoolKind:
		if s, ok := v.(string); ok {
			switch s {
			case "true":
				return true, false
			case "false":
				return false, false
			case "":
				return nil, true
			}
		}
		if n, ok := v.(json.Number); ok {
			switch n.String() {
			case "1":
				return true, false
			case "0":
				return false, false
			}
		}
	case protoreflect.StringKind:
		if n, ok := v.(json.Number); ok {
			return n.String(), false
		}
	case protoreflect.MessageKind, protoreflect.GroupKind, protoreflect.BytesKind:
	default: // numeric and enum kinds
		if s, ok := v.(string); ok && s == "" {
			return nil, true
		}
		// protojson accepts quoted numbers for integer kinds but not for
		// enums; the gateway sends enums as numbers.
		if s, ok := v.(string); ok && kind == protoreflect.EnumKind {
			if _, err := json.Number(s).Int64(); err == nil {
				return json.Number(s), false
			}
		}
		if v == nil {
			return nil, true
		}
	}
	return v, false
}

// ws returns the WebSocket client, dialing and authenticating it on first
// use. With TransportWebSocket this happens in newCore, so later calls only
// return the existing connection. Nothing that itself needs the connection is
// done while wsMu is held (the mutex is not reentrant).
func (c *core) ws(ctx context.Context) (client.Client, error) {
	c.wsMu.Lock()
	defer c.wsMu.Unlock()
	if c.closed {
		return nil, ErrClosed
	}
	if c.client != nil {
		return c.client, nil
	}
	// The connection outlives this call; don't tie it to the caller's ctx.
	cl, err := c.dial(context.WithoutCancel(ctx))
	if err != nil {
		return nil, err
	}
	c.client = cl
	for _, apply := range c.pendingHandlers {
		apply(cl)
	}
	c.pendingHandlers = nil
	return cl, nil
}

// connected reports whether the WebSocket has been dialed.
func (c *core) connected() bool {
	c.wsMu.Lock()
	defer c.wsMu.Unlock()
	return c.client != nil && !c.closed
}

// withClient runs `apply` against the WebSocket client now, or as soon as it
// is dialed when the transport is HTTP and no connection exists yet. It does
// not dial by itself.
func (c *core) withClient(apply func(client.Client)) {
	c.wsMu.Lock()
	defer c.wsMu.Unlock()
	if c.closed {
		return
	}
	if c.client != nil {
		apply(c.client)
		return
	}
	c.pendingHandlers = append(c.pendingHandlers, apply)
}
