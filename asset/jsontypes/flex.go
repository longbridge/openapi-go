package jsontypes

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

// ParseIntegral parses s as an integer that fits in bitSize bits. Besides
// plain integer text it accepts an integral float ("6.0", "2e1"), which some
// JSON encoders emit for integer fields. ok is false for empty, non-numeric,
// fractional or out-of-range text.
func ParseIntegral(s string, bitSize int) (v int64, ok bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if n, err := strconv.ParseInt(s, 10, bitSize); err == nil {
		return n, true
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsInf(f, 0) || math.IsNaN(f) || f != math.Trunc(f) {
		return 0, false
	}
	limit := math.Ldexp(1, bitSize-1)
	if f < -limit || f >= limit {
		return 0, false
	}
	return int64(f), true
}

// FlexString is a string field that also accepts a JSON number (kept as its
// literal text) or a boolean (as "true" / "false"); null, objects and arrays
// decode to "". It keeps a server-side type change from failing the whole
// response.
type FlexString string

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexString) UnmarshalJSON(data []byte) error {
	var v interface{}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return err
	}
	switch v := v.(type) {
	case string:
		*f = FlexString(v)
	case json.Number:
		*f = FlexString(v.String())
	case bool:
		*f = FlexString(strconv.FormatBool(v))
	default:
		*f = ""
	}
	return nil
}

// FlexInt64 is an int64 field that also accepts a numeric string or an
// integral float ("6.0"); null or an unparsable value (including a fractional
// number) decodes to 0.
type FlexInt64 int64

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexInt64) UnmarshalJSON(data []byte) error {
	var s FlexString
	if err := s.UnmarshalJSON(data); err != nil {
		return err
	}
	v, _ := ParseIntegral(string(s), 64)
	*f = FlexInt64(v)
	return nil
}

// FlexCode is an integer status code that also accepts a numeric string or an
// integral float ("2.0"). Valid is false when the field is missing, null, or
// not an integer that fits in int32 (e.g. "", "abc", 1.5), so a missing or
// garbled code can never be mistaken for a real code such as 0.
type FlexCode struct {
	Value int32
	Valid bool
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexCode) UnmarshalJSON(data []byte) error {
	*f = FlexCode{}
	var s FlexString
	if err := s.UnmarshalJSON(data); err != nil {
		return err
	}
	v, ok := ParseIntegral(string(s), 32)
	if !ok {
		return nil
	}
	*f = FlexCode{Value: int32(v), Valid: true}
	return nil
}

// FlexBool is a bool field that also accepts a number (non-zero is true) or a
// string ("true" / "1" and other non-zero numbers are true); anything else
// decodes to false.
type FlexBool bool

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexBool) UnmarshalJSON(data []byte) error {
	var s FlexString
	if err := s.UnmarshalJSON(data); err != nil {
		return err
	}
	t := strings.ToLower(strings.TrimSpace(string(s)))
	if t == "true" {
		*f = true
		return nil
	}
	n, err := strconv.ParseFloat(t, 64)
	*f = FlexBool(err == nil && n != 0)
	return nil
}

// FlexStrings is a string list whose items are decoded as FlexString; a value
// that is not an array decodes to nil, and null, object and array items are
// skipped.
type FlexStrings []string

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexStrings) UnmarshalJSON(data []byte) error {
	*f = nil
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil || items == nil {
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if !isScalar(item) {
			continue
		}
		var s FlexString
		if err := s.UnmarshalJSON(item); err != nil {
			continue
		}
		out = append(out, string(s))
	}
	*f = out
	return nil
}

// isScalar reports whether a raw JSON value is a string, number or boolean.
func isScalar(raw json.RawMessage) bool {
	t := bytes.TrimSpace(raw)
	if len(t) == 0 {
		return false
	}
	switch t[0] {
	case '{', '[', 'n':
		return false
	}
	return true
}

// FlexStringMap is a string-to-string object whose values are decoded as
// FlexString; a value that is not an object decodes to nil.
type FlexStringMap map[string]string

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexStringMap) UnmarshalJSON(data []byte) error {
	var items map[string]FlexString
	if err := json.Unmarshal(data, &items); err != nil || items == nil {
		*f = nil
		return nil
	}
	out := make(map[string]string, len(items))
	for k, v := range items {
		out[k] = string(v)
	}
	*f = out
	return nil
}

// FlexInt64s is an int64 list whose items are integers, numeric strings or
// integral floats; a value that is not an array decodes to nil, and items
// that are not integers are skipped.
type FlexInt64s []int64

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexInt64s) UnmarshalJSON(data []byte) error {
	*f = nil
	var items FlexStrings
	if err := items.UnmarshalJSON(data); err != nil || items == nil {
		return nil
	}
	out := make([]int64, 0, len(items))
	for _, item := range items {
		if v, ok := ParseIntegral(item, 64); ok {
			out = append(out, v)
		}
	}
	*f = out
	return nil
}

// Object is an optional nested JSON object. V is nil when the value is
// missing, null, not a JSON object, or an object that does not decode into T,
// so one malformed nested object never fails the enclosing response.
type Object[T any] struct {
	V *T
}

// UnmarshalJSON implements json.Unmarshaler.
func (o *Object[T]) UnmarshalJSON(data []byte) error {
	o.V = nil
	t := bytes.TrimSpace(data)
	if len(t) == 0 || t[0] != '{' {
		return nil
	}
	v := new(T)
	if err := json.Unmarshal(t, v); err != nil {
		return nil
	}
	o.V = v
	return nil
}

// List is a JSON array of objects decoded leniently: a value that is not an
// array decodes to nil, and elements that are null, not JSON objects or do
// not decode into T are skipped, so one malformed element never fails the
// enclosing response.
type List[T any] []T

// UnmarshalJSON implements json.Unmarshaler.
func (l *List[T]) UnmarshalJSON(data []byte) error {
	*l = nil
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil || items == nil {
		return nil
	}
	out := make([]T, 0, len(items))
	for _, item := range items {
		t := bytes.TrimSpace(item)
		if len(t) == 0 || t[0] != '{' {
			continue
		}
		var v T
		if err := json.Unmarshal(t, &v); err != nil {
			continue
		}
		out = append(out, v)
	}
	*l = out
	return nil
}
