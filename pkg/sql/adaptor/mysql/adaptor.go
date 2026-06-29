package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"

	"github.com/go-sql-driver/mysql"
	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/pkg/builder"
	sqlerr "github.com/xoctopus/sqlx/pkg/errors"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/sql/loggingdriver"
)

func init() {
	adaptor.Register(&mycli{})
}

type mycli struct {
	dialect
	adaptor.DB

	schema string
	dsn    *url.URL
}

func (d *mycli) Dialect() adaptor.Dialect {
	return d.dialect
}

func (d *mycli) DriverName() string {
	return "mysql"
}

func (d *mycli) Schema() string {
	return d.schema
}

func (d *mycli) Connector() driver.DriverContext {
	options := []loggingdriver.DriverOptionApplier{
		loggingdriver.WithDsnParser(ParseDSN),
		loggingdriver.WithErrorLeveler(ErrorLevel),
	}

	if v := d.dsn.Query().Get("interpolateParams"); len(v) == 0 || v == "false" {
		options = append(
			options,
			loggingdriver.WithInterpolator(loggingdriver.DefaultInterpolate),
		)
	}

	return loggingdriver.Wrap(mysql.MySQLDriver{}, d.DriverName(), options...)
}

// Open returns mysql adaptor.Adaptor
// dsn: mysql://[user[:password]@][addr]/database[?param1=value1&paramN=valueN]
func (d *mycli) Open(ctx context.Context, dsn *url.URL) (adaptor.Adaptor, error) {
	must.BeTrueF(
		dsn.Scheme == d.DriverName(),
		"invalid dsn schema, expect '%s' but got '%s'", d.DriverName(), dsn,
	)
	database := adaptor.DatabaseNameFromDSN(dsn)
	d.dsn = dsn

	conn, err := d.Connector().OpenConnector(d.dsn.String())
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(conn)

	if err = db.PingContext(ctx); err != nil {
		// always do closing if it needs creating database
		defer func() {
			println("closing: " + dsn.String())
			_ = db.Close()
		}()
		if d.IsUnknownDatabaseError(err) {
			if err = d.CreateDatabase(ctx, *dsn, database); err != nil {
				return nil, err
			}
			return d.Open(ctx, dsn)
		}
		return nil, err
	}

	return &mycli{
		DB: adaptor.Wrap(db, func(err error) error {
			if d.IsConflictError(err) {
				return codex.Errorf(sqlerr.CONFLICT, "%v", err)
			}
			return err
		}),
		schema: database,
		dsn:    dsn,
	}, nil
}

func (d *mycli) Catalog(ctx context.Context) (builder.Catalog, error) {
	return ScanCatalog(ctx, d, d.schema)
}

func (d *mycli) CreateDatabase(ctx context.Context, dsn url.URL, database string) error {
	dsn.Path = "/mysql"

	a, err := d.Open(ctx, &dsn)
	if err != nil {
		return err
	}
	defer func() {
		println("closing: " + dsn.String())
		_ = a.Close()
	}()

	_, err = a.Exec(ctx, frag.Query("CREATE DATABASE ?", frag.Lit(database)))
	return err
}

func ParseDSN(dsn string) (string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	pass, _ := u.User.Password()
	dsn = fmt.Sprintf("%s:%s@tcp(%s)%s", u.User.Username(), pass, u.Host, u.Path)
	if q := u.Query(); len(q) > 0 {
		dsn += "?" + q.Encode()
	}
	return dsn, nil
}
