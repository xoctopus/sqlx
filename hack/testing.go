package hack

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/xoctopus/confx/pkg/types"
	"github.com/xoctopus/logx"
	"github.com/xoctopus/x/misc/retry"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	_ "github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/session"
)

func Check(t testing.TB) {
	if os.Getenv("HACK_TEST") != "true" {
		t.Skip("should depend on postgres/mysql")
	}
}

func Context(t testing.TB) context.Context {
	t.Helper()
	return logx.With(context.Background(), logx.Std(logx.NewHandler()))
}

func NewAdaptor(t testing.TB, dsn string) adaptor.Adaptor {
	Check(t)

	_, err := url.Parse(dsn)
	Expect(t, err, Succeed())

	var a adaptor.Adaptor

	err = (&retry.Retry{Repeats: 10, Interval: 3 * time.Second}).
		Do(func() (err error) {
			a, err = adaptor.Open(Context(t), dsn)
			return err
		})
	Expect(t, err, Succeed())

	t.Cleanup(func() { _ = a.Close() })
	return a
}

func NewAdaptorError(t testing.TB, dsn string) error {
	Check(t)

	_, err := url.Parse(dsn)
	Expect(t, err, Succeed())

	return (&retry.Retry{Repeats: 5, Interval: 3 * time.Second}).
		Do(func() (err error) {
			_, err = adaptor.Open(Context(t), dsn)
			return err
		})
}

func WithSession(ctx context.Context, t testing.TB, dsn string, catalogs ...builder.Catalog) context.Context {
	Check(t)

	_, err := url.Parse(dsn)
	Expect(t, err, Succeed())

	ep := session.Endpoint{
		Endpoint: types.Endpoint[session.EndpointOption]{Address: dsn},
	}
	ep.ApplyCatalog("sqlx.hack", catalogs...)

	err = (&retry.Retry{Repeats: 10, Interval: 3 * time.Second}).
		Do(func() (err error) {
			err = ep.Init(ctx)
			return err
		})
	Expect(t, err, Succeed())

	t.Cleanup(func() { _ = ep.Close() })

	return session.With(ctx, ep.Session())
}
