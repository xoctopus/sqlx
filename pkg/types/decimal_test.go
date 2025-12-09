package types_test

import (
	"database/sql/driver"
	"testing"

	"github.com/shopspring/decimal"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/types"
)

func TestDecimal(t *testing.T) {
	v := types.AsDecimal(decimal.NewFromInt(100))
	Expect(t, v.DBType("any"), Equal("decimal"))

	Expect(t, v.Scan([]byte("111111.11")), Succeed())
	fv, _ := v.Float64()
	Expect(t, fv, Equal(111111.11))

	dv, err := v.Value()
	Expect(t, err, Succeed())
	Expect(t, dv, Equal[driver.Value]("111111.11"))
}
