package sqlite

/*
import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"

	"github.com/xoctopus/x/codex"
	"modernc.org/sqlite"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/loggingdriver"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/errors"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func init() {
	// adaptor.Register(&litecli{}, "sqlite3")
}

func Open(ctx context.Context, dsn *url.URL) (adaptor.Adaptor, error) {
	return (&litecli{dsn: dsn}).Open(ctx, dsn)
}

type litecli struct {
	dialect
	adaptor.DB

	dsn *url.URL
}

func (d litecli) DriverName() string {
	return "sqlite"
}

func (d litecli) Endpoint() string {
	return d.dsn.Path
}

func (d litecli) Dialect() adaptor.Dialect {
	return d.dialect
}

func (d litecli) Connector() driver.DriverContext {
	return loggingdriver.Wrap(
		&sqlite.Driver{},
		d.DriverName(),
		loggingdriver.WithInterpolator(loggingdriver.OrderedInterpolator),
		loggingdriver.WithErrorLeveler(
			func(err error) int {
				if d.IsConflictError(err) {
					return 0
				}
				return 1
			},
		),
	)
}

func (d litecli) Open(ctx context.Context, dsn *url.URL) (adaptor.Adaptor, error) {
	if dsn.Scheme != d.DriverName() {
		return nil, fmt.Errorf("invalid dsn schema, expect '%s' but got '%s'", d.DriverName(), dsn)
	}

	c, err := d.Connector().OpenConnector(dsn.Path)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(c)
	db.SetMaxOpenConns(1)

	a := &litecli{
		DB: adaptor.Wrap(db, func(err error) error {
			if d.IsConflictError(err) {
				return codex.Errorf(errors.CONFLICT, "%v", err)
			}
			return err
		}),
	}

	query := dsn.Query()
	if !query.Has("journal_mode") {
		query.Set("journal_mode", "WAL")
	}
	if !query.Has("busy_timeout") {
		query.Set("busy_timeout", "5000")
	}
	if !query.Has("synchronous") {
		query.Set("synchronous", "NORMAL")
	}

	_, err = a.Exec(ctx, frag.Query("PRAGMA journal_mode = ?", frag.Lit(query.Get("journal_mode"))))
	if err != nil {
		return nil, err
	}

	_, err = a.Exec(ctx, frag.Query("PRAGMA busy_timeout = ?", frag.Lit(query.Get("busy_timeout"))))
	if err != nil {
		return nil, err
	}

	_, err = a.Exec(ctx, frag.Query("PRAGMA synchronous = ?", frag.Lit(query.Get("synchronous"))))
	if err != nil {
		return nil, err
	}

	return a, err
}

func (d litecli) Catalog(ctx context.Context) (builder.Catalog, error) {
	return ScanCatalog(ctx, d)
}
*/
