package util

import (
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var opt = copier.Option{
	IgnoreEmpty: true,
	DeepCopy:    false,
	Converters: []copier.TypeConverter{
		{
			SrcType: copier.String,
			DstType: decimal.Decimal{},
			Fn: func(src interface{}) (interface{}, error) {
				value, ok := src.(string)

				if !ok {
					return nil, errors.New("convert string to decimal, but src type not matching")
				}

				return decimal.NewFromString(value)
			},
		},
		{
			SrcType: time.Time{},
			DstType: copier.String,
			Fn: func(src interface{}) (interface{}, error) {
				value, ok := src.(time.Time)

				if !ok {
					return nil, errors.New("convert time to string, but src type not matching")
				}
				return FormatDate(&value), nil
			},
		},
		{
			SrcType: copier.String,
			DstType: int64(0),
			Fn: func(src interface{}) (interface{}, error) {
				value, ok := src.(string)

				if !ok {
					return nil, errors.New("convert string to int64, but src type not matching")
				}

				return strconv.ParseInt(value, 10, 64)
			},
		},
		{
			SrcType: int64(0),
			DstType: time.Time{},
			Fn: func(src interface{}) (interface{}, error) {
				value, ok := src.(int64)

				if !ok {
					return nil, errors.New("convert int64 to time.Time, but src type not matching")
				}

				if value < 0 {
					return nil, errors.New("convert int64 to time.Time, but src value is less than 0")
				}

				// Check if it's seconds (most common case)
				if value <= 4102444800 { // ~2106-02-07
					return time.Unix(value, 0), nil
				}

				// Check if it's milliseconds
				if value <= 4102444800000 { // ~2286-11-20
					return time.Unix(0, value*int64(time.Millisecond)), nil
				}

				// Check if it's microseconds
				if value <= 4102444800000000 { // ~2262-04-11
					return time.Unix(0, value*int64(time.Microsecond)), nil
				}

				// Check if it's nanoseconds
				if value <= 4102444800000000000 { // ~2262-04-11
					return time.Unix(0, value), nil
				}

				return nil, errors.New("convert int64 to time.Time, but src value is not valid")
			},
		},
	},
}

func Copy(toValue interface{}, fromValue interface{}) error {
	return copier.CopyWithOption(toValue, fromValue, opt)
}
