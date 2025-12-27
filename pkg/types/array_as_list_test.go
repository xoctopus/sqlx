package types_test

import (
	"database/sql/driver"
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/sqlx/pkg/types"
)

func TestArrayAsList(t *testing.T) {
	aa := ArrayAsList[string]{}
	Expect(t, aa.DBType("any"), Equal("text"))

	Expect(t, aa.AppendString("1"), Succeed())
	Expect(t, aa.AppendString("2"), Succeed())

	v, err := aa.Value()
	Expect(t, err, Succeed())
	Expect(t, v, Equal[driver.Value]("1,2"))

	aa2 := ArrayAsList[int]{}
	aa2.Append(1, 2)
	v, err = aa2.Value()
	Expect(t, err, Succeed())
	Expect(t, v, Equal[driver.Value]("1,2"))

	Expect(t, aa.Scan("1,2,3"), Succeed())
	Expect(t, aa.String(), Equal("1,2,3"))

	Expect(t, aa.Scan([]byte("1,2,3,4")), Succeed())
	Expect(t, aa.String(), Equal("1,2,3,4"))

	Expect(t, aa.Scan(float32(0)), Failed())

	type ID uint64

	aa3 := ArrayAsList[ID]{}
	Expect(t, aa3.AppendString("1,2,3,4,5"), Succeed())
	Expect(t, aa3.String(), Equal("1,2,3,4,5"))

	data, err := aa3.MarshalJSON()
	Expect(t, err, Succeed())
	Expect(t, data, Equal([]byte(`"1,2,3,4,5"`)))

	aa4 := ArrayAsList[ID]{}
	err = aa4.UnmarshalJSON(data)
	Expect(t, err, Succeed())
	Expect(t, aa4, Equal(aa3))
}
