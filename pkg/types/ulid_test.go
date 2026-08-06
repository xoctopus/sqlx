package types_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/types"
)

func TestULID(t *testing.T) {
	u := types.MakeULID()

	t.Run("DBType", func(t *testing.T) {
		Expect(t, u.DBType("postgres"), Equal("uuid"))
		Expect(t, u.DBType("Postgres"), Equal("uuid"))
		Expect(t, u.DBType("duckdb"), Equal("uuid"))
		Expect(t, u.DBType("DuckDB"), Equal("uuid"))
		Expect(t, u.DBType("sqlite"), Equal("blob"))
		Expect(t, u.DBType("SQLite"), Equal("blob"))
		Expect(t, u.DBType("mysql"), Equal("binary"))
		Expect(t, u.DBType(""), Equal("binary"))
	})

	t.Run("DBFixedWidth", func(t *testing.T) {
		Expect(t, u.DBFixedWidth("postgres"), BeNil[*uint64]())
		Expect(t, u.DBFixedWidth("duckdb"), BeNil[*uint64]())
		Expect(t, u.DBFixedWidth("sqlite"), BeNil[*uint64]())
		Expect(t, u.DBFixedWidth("mysql"), Equal(new(uint64(16))))
		Expect(t, u.DBFixedWidth(""), Equal(new(uint64(16))))
	})

	t.Log(u.BigInt())
	t.Log(u.String())
}
