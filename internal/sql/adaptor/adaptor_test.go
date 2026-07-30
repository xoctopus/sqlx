package adaptor

import (
	"context"
	"database/sql"
	"net/url"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

type stubAdaptor struct {
	name   string
	opened *url.URL
	openN  int
}

func (a *stubAdaptor) D() *sql.DB   { return nil }
func (a *stubAdaptor) Close() error { return nil }
func (a *stubAdaptor) Exec(context.Context, frag.Fragment) (sql.Result, error) {
	return nil, nil
}
func (a *stubAdaptor) Query(context.Context, frag.Fragment) (*sql.Rows, error) {
	return nil, nil
}
func (a *stubAdaptor) Tx(context.Context, func(context.Context) error) error {
	return nil
}
func (a *stubAdaptor) DriverName() string { return a.name }
func (a *stubAdaptor) Schema() string     { return "test" }
func (a *stubAdaptor) Dialect() Dialect   { return nil }
func (a *stubAdaptor) Catalog(context.Context) (builder.Catalog, error) {
	return nil, nil
}

func (a *stubAdaptor) Open(_ context.Context, u *url.URL) (Adaptor, error) {
	a.openN++
	a.opened = u
	return a, nil
}

func TestRegisterOpen(t *testing.T) {
	const (
		driver = "stubtestdriver"
		alias  = "stubtestalias"
	)
	a := &stubAdaptor{name: driver}
	Register(a, alias)
	t.Cleanup(func() {
		adaptors.Delete(driver)
		adaptors.Delete(alias)
	})

	t.Run("ByDriverName", func(t *testing.T) {
		a.openN = 0
		got, err := Open(context.Background(), driver+"://host/db")
		Expect(t, err, Succeed())
		Expect(t, got, Be[Adaptor](a))
		Expect(t, a.openN, Equal(1))
		Expect(t, a.opened.Scheme, Equal(driver))
		Expect(t, a.opened.Host, Equal("host"))
		Expect(t, DatabaseNameFromDSN(a.opened), Equal("db"))
	})

	t.Run("ByAlias", func(t *testing.T) {
		a.openN = 0
		got, err := Open(context.Background(), alias+"://h/db2")
		Expect(t, err, Succeed())
		Expect(t, got, Be[Adaptor](a))
		Expect(t, a.openN, Equal(1))
		Expect(t, DatabaseNameFromDSN(a.opened), Equal("db2"))
	})

	t.Run("MissingAdaptor", func(t *testing.T) {
		_, err := Open(context.Background(), "missingstubdriver://h/db")
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("missing adaptor"))
	})

	t.Run("InvalidDSN", func(t *testing.T) {
		_, err := Open(context.Background(), "://bad")
		Expect(t, err, Failed())
	})
}
