package util_test

import (
	"testing"

	"github.com/longbridgeapp/assert"
	"github.com/longbridge/openapi-go/internal/util"

	"github.com/shopspring/decimal"
)

func TestCopy(t *testing.T) {
	fv := []*From{{
		Price:  "0.1",
		Price1: "a",
		Price2: "1.1",
		Price3: "1a",
		Num:    "",
		Num1:   "12",
	}}
	price, err := decimal.NewFromString("0.1")
	assert.NoError(t, err)
	price2 := decimal.NewFromFloat(1.1)
	expectedTv := To{
		Price:  &price,
		Price1: &decimal.Decimal{},
		Price2: price2,
		Price3: decimal.Decimal{},
		Num1: 12,
	}
	var tv To
	err = util.Copy(&tv, fv[0])
	assert.NoError(t, err)
	assert.Equal(t, expectedTv, tv)
}

type From struct {
	Price  string
	Price1 string
	Price2 string
	Price3 string
	Num    string
	Num1   string
}

type To struct {
	Price  *decimal.Decimal
	Price1 *decimal.Decimal
	Price2 decimal.Decimal
	Price3 decimal.Decimal
	Num    int64
	Num1   int64
}
