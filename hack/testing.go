package hack

import (
	"context"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/xoctopus/logx"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	_ "github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
)

var once sync.Once

func Check(t testing.TB) {
	if os.Getenv("HACK_TEST") != "true" {
		t.Skip("should depend on postgres/mysql")
	}
	once.Do(func() {
		// time.Sleep(time.Second * 5) // to wait dependencies ready
	})
}

func Context(t testing.TB) context.Context {
	t.Helper()
	return logx.With(context.Background(), logx.Std(logx.NewHandler()))
}

func NewAdaptor(t testing.TB, dsn string) adaptor.Adaptor {
	Check(t)

	_, err := url.Parse(dsn)
	Expect(t, err, Succeed())

	a, err := adaptor.Open(Context(t), dsn)
	Expect(t, err, Succeed())

	t.Cleanup(func() {
		_ = a.Close()
	})
	return a
}
