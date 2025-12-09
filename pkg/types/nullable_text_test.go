package types_test

import (
	"database/sql/driver"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/types"
)

func TestText(t *testing.T) {
	b := types.Text("")

	Expect(t, b.DBType("any"), Equal("text"))

	v, err := b.Value()
	Expect(t, err, Succeed())
	Expect(t, v, Equal[driver.Value](nil))

	b.Set("hello")
	v, err = b.Value()
	Expect(t, err, Succeed())
	Expect(t, v, Equal[driver.Value]("hello"))

	Expect(t, b.Scan([]byte("hello")), Succeed())
	Expect(t, b, Equal(types.Text("hello")))

	Expect(t, b.Scan("hello"), Succeed())
	Expect(t, b, Equal(types.Text("hello")))

	Expect(t, b.Scan(nil), Succeed())
	Expect(t, b, Equal(types.Text("")))

	Expect(t, b.Scan(1), Failed())
}
