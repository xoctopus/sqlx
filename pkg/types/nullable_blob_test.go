package types_test

import (
	"database/sql/driver"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/types"
)

func TestBlob(t *testing.T) {
	b := types.Blob{}

	Expect(t, b.DBType("postgres"), Equal("bytea"))
	Expect(t, b.DBType("mysql"), Equal("blob"))

	v, err := b.Value()
	Expect(t, err, Succeed())
	Expect(t, v, Equal[driver.Value](nil))

	b.Set([]byte("hello"))
	v, err = b.Value()
	Expect(t, err, Succeed())
	Expect(t, v, Equal[driver.Value]([]byte("hello")))

	Expect(t, b.Scan([]byte("hello")), Succeed())
	Expect(t, b, Equal(types.Blob([]byte("hello"))))

	Expect(t, b.Scan("hello"), Succeed())
	Expect(t, b, Equal(types.Blob([]byte("hello"))))

	Expect(t, b.Scan(nil), Succeed())
	Expect(t, b, Equal(types.Blob(nil)))

	Expect(t, b.Scan(1), Failed())
}
