package mysql_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/hack"
	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func TestOpen_Hack(t *testing.T) {
	hack.Check(t)

	t.Run("FailedToAuth", func(t *testing.T) {
		_, err := adaptor.Open(hack.Context(t), "mysql://user:pass@localhost:13306/test")
		Expect(t, mysql.IsUnknownDatabaseError(err), BeFalse())

		ue := mysql.UnwrapError(err)
		Expect(t, any(ue), NotBeNil[any]())
	})
	t.Run("InvalidSchema", func(t *testing.T) {
		_, err := adaptor.Open(hack.Context(t), "invalid://user:pass@localhost:13306/test")
		Expect(t, err, ErrorContains("missing adaptor"))
	})
	t.Run("NeedCreateDatabase", func(t *testing.T) {
		_, err := adaptor.Open(hack.Context(t), "mysql://root@localhost:13306/invalid.db")
		Expect(t, err, Failed())
		ue := mysql.UnwrapError(err)
		Expect(t, ue.Number, Equal(uint16(1064))) // caused by invalid database name

		t.Run("Success", func(t *testing.T) {
			d, err := adaptor.Open(hack.Context(t), "mysql://root@localhost:13306/fresh")
			Expect(t, err, Succeed())

			t.Cleanup(func() {
				_ = d.Close()
			})

			Expect(t, d.Endpoint(), Equal("mysql://localhost:13306"))
			dialect := d.Dialect()
			_, err = d.Exec(hack.Context(t), dialect.SwitchSchema("mysql"))
			Expect(t, err, Succeed())
			_, err = d.Exec(hack.Context(t), frag.Query("DROP DATABASE ?", frag.Lit("fresh")))
			Expect(t, err, Succeed())
		})
	})
}
