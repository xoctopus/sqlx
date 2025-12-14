package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"

	"github.com/go-sql-driver/mysql"
	"github.com/xoctopus/x/codex"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/loggingdriver"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/errors"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func init() {
	adaptor.Register(&mycli{})
}

func Open(ctx context.Context, dsn *url.URL) (adaptor.Adaptor, error) {
	return (&mycli{dsn: dsn}).Open(ctx, dsn)
}

type mycli struct {
	dialect
	adaptor.DB

	database string
	dsn      *url.URL
}

func (d *mycli) Dialect() adaptor.Dialect {
	return d.dialect
}

func (d *mycli) DriverName() string {
	return "mysql"
}

func (d *mycli) Endpoint() string {
	return (&url.URL{
		Scheme: d.dsn.Scheme,
		Host:   d.dsn.Host,
	}).String()
}

func (d *mycli) Connector() driver.DriverContext {
	return loggingdriver.Wrap(
		mysql.MySQLDriver{},
		d.DriverName(),
		loggingdriver.WithDsnParser(ParseDSN),
		loggingdriver.WithInterpolator(loggingdriver.DefaultInterpolate),
	)
}

// Open return
// dsn: mysql://[user[:password]@][addr]/database[?param1=value1&paramN=valueN]
func (d *mycli) Open(ctx context.Context, dsn *url.URL) (a adaptor.Adaptor, err error) {
	if dsn.Scheme != d.DriverName() {
		return nil, fmt.Errorf("invalid dsn schema, expect '%s' but got '%s'", d.DriverName(), dsn)
	}

	var (
		database = adaptor.DatabaseNameFromDSN(dsn)
		conn     driver.Connector
	)

	conn, err = d.Connector().OpenConnector(dsn.String())
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(conn)

	a = &mycli{
		DB: adaptor.Wrap(db, func(err error) error {
			if d.IsConflictError(err) {
				return codex.Errorf(errors.CONFLICT, "%v", err)
			}
			return err
		}),
		database: database,
		dsn:      d.dsn,
	}

	defer func() {
		if err != nil {
			_ = db.Close()
		}
	}()

	if err = db.PingContext(ctx); err != nil {
		if d.IsUnknownDatabaseError(err) {
			if err = d.CreateDatabase(ctx, *dsn, database); err != nil {
				return nil, err
			}
			return d.Open(ctx, dsn)
		}
		return nil, err
	}

	return a, nil
}

func (d *mycli) Catalog(ctx context.Context) (builder.Catalog, error) {
	return ScanCatalog(ctx, d, d.database)
}

func (d *mycli) CreateDatabase(ctx context.Context, dsn url.URL, database string) error {
	dsn.Path = "/mysql"

	a, err := d.Open(ctx, &dsn)
	if err != nil {
		return err
	}
	defer a.Close()

	_, err = a.Exec(ctx, frag.Query("CREATE DATABASE ?", frag.Lit(database)))
	return err
}

func ParseDSN(dsn string) (string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	dsn = fmt.Sprintf("%s@tcp(%s)%s", u.User, u.Host, u.Path)
	if q := u.Query(); len(q) > 0 {
		dsn += "?" + q.Encode()
	}
	return dsn, nil
}
