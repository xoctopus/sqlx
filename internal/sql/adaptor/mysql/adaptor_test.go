package mysql_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/xoctopus/logx"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/hack"
	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
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
	hack.Check(t)
	t.Run("FailedToAuth", func(t *testing.T) {
		dsn, err := url.Parse("mysql://user:pass@localhost:13306/test")
		Expect(t, err, Succeed())

		_, err = mysql.Open(context.Background(), dsn)
		Expect(t, mysql.IsUnknownDatabase(err), BeFalse()) // this should be an auth error
	})

	_ = NewAdaptor(t)
}
