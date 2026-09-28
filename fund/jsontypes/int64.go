package jsontypes

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// Int64 is an int64 that tolerates the fund backend encoding int64 values as
// JSON strings (e.g. "1790229760"). It accepts a JSON number, a quoted numeric
// string, null, and an empty string (decoded as 0). This mirrors the Rust
// core's serde `int64_str` / `int64_str_empty_is_none` handling.
type Int64 int64

// UnmarshalJSON implements json.Unmarshaler.
func (i *Int64) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	// null -> 0 (leave the field at its zero value).
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*i = 0
		return nil
	}

	// Strip surrounding quotes if the value arrived as a JSON string.
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
		data = bytes.TrimSpace(data)
		// Empty string -> 0 (mirrors int64_str_empty_is_none).
		if len(data) == 0 {
			*i = 0
			return nil
		}
	}

	v, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*i = Int64(v)
	return nil
}

// MarshalJSON implements json.Marshaler, emitting a plain JSON number.
func (i Int64) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, int64(i), 10), nil
}

// compile-time interface checks.
var (
	_ json.Unmarshaler = (*Int64)(nil)
	_ json.Marshaler   = Int64(0)
)
