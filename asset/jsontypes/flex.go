package jsontypes

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

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

// FlexInt64 is an int64 field that also accepts a numeric string; null or an
// unparsable value (including a float) decodes to 0.
type FlexInt64 int64

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexInt64) UnmarshalJSON(data []byte) error {
	var s FlexString
	if err := s.UnmarshalJSON(data); err != nil {
		return err
	}
	v, err := strconv.ParseInt(strings.TrimSpace(string(s)), 10, 64)
	if err != nil {
		v = 0
	}
	*f = FlexInt64(v)
	return nil
}

// FlexCode is an integer status code that also accepts a numeric string.
// Valid is false when the field is missing, null, or not an integer that fits
// in int32 (e.g. "", "abc", 1.5), so a missing or garbled code can never be
// mistaken for a real code such as 0.
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
	v, err := strconv.ParseInt(strings.TrimSpace(string(s)), 10, 32)
	if err != nil {
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
// that is not an array decodes to nil.
type FlexStrings []string

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexStrings) UnmarshalJSON(data []byte) error {
	var items []FlexString
	if err := json.Unmarshal(data, &items); err != nil {
		*f = nil
		return nil
	}
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = string(item)
	}
	*f = out
	return nil
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

// FlexInt64s is an int64 list whose items are decoded as FlexInt64; a value
// that is not an array decodes to nil.
type FlexInt64s []int64

// UnmarshalJSON implements json.Unmarshaler.
func (f *FlexInt64s) UnmarshalJSON(data []byte) error {
	var items []FlexInt64
	if err := json.Unmarshal(data, &items); err != nil {
		*f = nil
		return nil
	}
	out := make([]int64, len(items))
	for i, item := range items {
		out[i] = int64(item)
	}
	*f = out
	return nil
}
