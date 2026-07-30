package mysql_test

import (
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/hack"
	"github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func TestOpen_Hack(t *testing.T) {
	t.Run("FailedToAuth", func(t *testing.T) {
		err := hack.NewInvalidAdaptor(t, "mysql://user1:pass@localhost:13306/test")
		Expect(t, err, Failed())
		Expect(t, mysql.IsUnknownDatabaseError(err), BeFalse())

		ue := mysql.UnwrapError(err)
		Expect(t, any(ue), NotBeNil[any]())
	})
	t.Run("InvalidSchema", func(t *testing.T) {
		err := hack.NewInvalidAdaptor(t, "invalid://user1:pass@localhost:13306/test")
		Expect(t, err, ErrorContains("missing adaptor"))
	})
	t.Run("NeedCreateDatabase", func(t *testing.T) {
		err := hack.NewInvalidAdaptor(t, "mysql://root@localhost:13306/invalid.db")
		Expect(t, err, Failed())
		ue := mysql.UnwrapError(err)
		Expect(t, ue.Number, Equal(uint16(1064))) // caused by invalid database name
	})
	t.Run("Success", func(t *testing.T) {
		t.Run("NoOption", func(t *testing.T) {
			dsn := "mysql://root@localhost:13306/fresh"
			d := hack.NewAdaptor(t, dsn)

			dialect := d.Dialect()
			_, err := d.Exec(hack.Context(t), dialect.SwitchSchema("mysql"))
			Expect(t, err, Succeed())
			_, err = d.Exec(hack.Context(t), frag.Query("DROP DATABASE ?", frag.Lit("fresh")))
			Expect(t, err, Succeed())
		})
		t.Run("HasOptions", func(t *testing.T) {
			dsn := "mysql://root@localhost:13306/fresh2?interpolateParams=true"
			d := hack.NewAdaptor(t, dsn)

			dialect := d.Dialect()
			_, err := d.Exec(hack.Context(t), dialect.SwitchSchema("mysql"))
			Expect(t, err, Succeed())
			_, err = d.Exec(hack.Context(t), frag.Query("DROP DATABASE ?", frag.Lit("fresh2")))
			Expect(t, err, Succeed())
		})
	})
}
