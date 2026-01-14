package hack

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/xoctopus/confx/pkg/types"
	"github.com/xoctopus/logx"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	_ "github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/session"
)

var once sync.Once

func Check(t testing.TB) {
	if os.Getenv("HACK_TEST") != "true" {
		t.Skip("should depend on postgres/mysql")
	}
	once.Do(func() {
		fmt.Println("HACK_TEST ENABLED")
		time.Sleep(time.Second * 5) // to wait dependencies ready
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

func WithSession(ctx context.Context, t testing.TB, dsn string, catalogs ...builder.Catalog) context.Context {
	Check(t)

	_, err := url.Parse(dsn)
	Expect(t, err, Succeed())

	ep := session.Endpoint{
		Endpoint: types.Endpoint[session.EndpointOption]{Address: dsn},
	}

	ep.ApplyCatalog("sqlx.hack", catalogs...)
	Expect(t, ep.Init(ctx), Succeed())

	return session.With(ctx, ep.Session())
}
