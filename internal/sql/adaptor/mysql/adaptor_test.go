package mysql_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/xoctopus/logx"
	"github.com/xoctopus/x/misc/must"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func Context(t testing.TB) context.Context {
	t.Helper()
	return logx.WithLogger(context.Background(), logx.Std(logx.NewHandler()))
}

func NewAdaptor(t testing.TB) adaptor.Adaptor {
	dsn, err := url.Parse("mysql://root@localhost:13306/test")
	Expect(t, err, Succeed())

	a, err := mysql.Open(Context(t), dsn)
	Expect(t, err, Succeed())

	t.Cleanup(func() {
		_ = a.Close()
	})
	return a
}

func TestOpen_Hack(t *testing.T) {
	// hack.Check(t)
	t.Run("FailedToAuth", func(t *testing.T) {
		dsn, err := url.Parse("mysql://user:pass@localhost:13306/test")
		Expect(t, err, Succeed())

		_, err = mysql.Open(Context(t), dsn)
		Expect(t, mysql.IsUnknownDatabaseError(err), BeFalse())

		ue := mysql.UnwrapError(err)
		Expect(t, any(ue), NotBeNil[any]())
	})
	t.Run("InvalidSchema", func(t *testing.T) {
		_, err := mysql.Open(context.Background(), &url.URL{Scheme: "not_mysql"})
		Expect(t, err, ErrorContains("invalid dsn schema"))
	})
	t.Run("NeedCreateDatabase", func(t *testing.T) {
		dsn := must.NoErrorV(url.Parse("mysql://root@localhost:13306/invalid.db"))
		_, err := mysql.Open(Context(t), dsn)
		Expect(t, err, Failed())
		ue := mysql.UnwrapError(err)
		Expect(t, ue.Number, Equal(uint16(1064))) // caused by invalid database name

		t.Run("Success", func(t *testing.T) {
			dsn = must.NoErrorV(url.Parse("mysql://root@localhost:13306/fresh"))
			d, err := mysql.Open(Context(t), dsn)
			Expect(t, err, Succeed())

			Expect(t, d.Endpoint(), Equal("mysql://localhost:13306"))

			dialect := d.Dialect()
			_, err = d.Exec(Context(t), dialect.SwitchSchema("mysql"))
			Expect(t, err, Succeed())
			_, err = d.Exec(Context(t), frag.Query("DROP DATABASE ?", frag.Lit("fresh")))
			Expect(t, err, Succeed())
			_ = d.Close()
		})
	})
}
