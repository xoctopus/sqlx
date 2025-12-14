package postgres

/*
import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/xoctopus/x/codex"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/loggingdriver"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/errors"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func init() {
	// adaptor.Register(&pgcli{}, "pg", "postgresql")
}

func Open(ctx context.Context, dsn *url.URL) (adaptor.Adaptor, error) {
	return (&pgcli{dsn: dsn}).Open(ctx, dsn)
}

type pgcli struct {
	dialect
	adaptor.DB

	database string
	dsn      *url.URL

	pool *pgxpool.Pool
	once sync.Once
	perr error
}

func (d *pgcli) Dialect() adaptor.Dialect {
	return d.dialect
}

func (d *pgcli) DriverName() string {
	return "postgres"
}

func (d *pgcli) Endpoint() string {
	return (&url.URL{
		Scheme: d.dsn.Scheme,
		Host:   d.dsn.Host,
		Path:   d.dsn.Path,
	}).String()
}

func (d *pgcli) Connector() driver.DriverContext {
	return loggingdriver.Wrap(
		stdlib.GetPoolConnector(d.pool).Driver(),
		d.DriverName(),
		loggingdriver.WithErrorLeveler(
			func(err error) int {
				if d.IsConflictError(err) {
					return 0
				}
				return 1
			},
		),
		loggingdriver.WithInterpolator(loggingdriver.OrderedInterpolator),
	)
}

func (d *pgcli) Open(ctx context.Context, dsn *url.URL) (adaptor.Adaptor, error) {
	if dsn.Scheme != d.DriverName() {
		return nil, fmt.Errorf("invalid dsn schema, expect '%s' but got '%s'", d.DriverName(), dsn.Scheme)
	}

	pconn := url.Values{}
	ppool := url.Values{}

	for k, vs := range dsn.Query() {
		if _, ok := extras[k]; ok {
			pconn[k] = vs
		}
		ppool[k] = vs
	}

	d.once.Do(func() {
		if !ppool.Has("pool_max_conns") {
			ppool.Set("pool_max_conns", "10")
		}
		if !ppool.Has("pool_max_conn_lifetime") {
			ppool.Set("pool_max_conn_lifetime", "1h")
		}
		dsn.RawQuery = ppool.Encode()

		pool, err := pgxpool.New(ctx, dsn.String())
		if err != nil {
			d.perr = err
			return
		}
		d.pool = pool
	})

	database := adaptor.DatabaseNameFromDSN(dsn)
	dsn.RawQuery = pconn.Encode()

	c, err := d.Connector().OpenConnector(dsn.String())
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(c)
	if err = db.PingContext(ctx); err != nil {
		if d.IsUnknownDatabaseError(err) {
			if err = d.CreateDatabase(ctx, database, *dsn); err != nil {
				return nil, err
			}
			return d.Open(ctx, dsn)
		}
		return nil, err
	}

	return &pgcli{
		database: database,
		DB: adaptor.Wrap(db, func(err error) error {
			if d.IsConflictError(err) {
				return codex.Errorf(errors.CONFLICT, "%v", err)
			}
			return err
		}),
	}, nil
}

func (d *pgcli) Catalog(ctx context.Context) (builder.Catalog, error) {
	return ScanCatalog(ctx, d)
}

func (d *pgcli) CreateDatabase(ctx context.Context, database string, dsn url.URL) error {
	dsn.Path = "/postgres"

	a, err := d.Open(ctx, &dsn)
	if err != nil {
		return err
	}
	defer a.Close()

	_, err = a.Exec(context.Background(), frag.Query("CREATE DATABASE ?;", frag.Lit(database)))
	return err
}

// extras parameters for non-connection
var extras = map[string]struct{}{
	"host":                 {},
	"port":                 {},
	"database":             {},
	"user":                 {},
	"password":             {},
	"passfile":             {},
	"connect_timeout":      {},
	"sslmode":              {},
	"sslkey":               {},
	"sslcert":              {},
	"sslrootcert":          {},
	"sslpassword":          {},
	"sslsni":               {},
	"krbspn":               {},
	"krbsrvname":           {},
	"target_session_attrs": {},
	"service":              {},
	"servicefile":          {},
	"_ro":                  {},
}
*/
