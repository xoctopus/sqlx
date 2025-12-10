package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-sql-driver/mysql"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/loggingdriver"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func init() {
	c := &mycli{}
	adaptor.Register(c, c.DriverName())
}

func Open(ctx context.Context, dsn *url.URL) (adaptor.Adaptor, error) {
	return (&mycli{}).Open(ctx, dsn)
}

type mycli struct {
	dialect
	adaptor.DB

	database string
}

func (d *mycli) Dialect() adaptor.Dialect {
	return &d.dialect
}

func (d *mycli) DriverName() string {
	return "mysql"
}

func (d *mycli) Connector() driver.DriverContext {
	return loggingdriver.Wrap(
		mysql.MySQLDriver{},
		d.DriverName(),
		loggingdriver.WithErrorLeveler(
			func(err error) int {
				var e *mysql.MySQLError // duplicate entry
				if errors.As(err, &e) && e.Number == 1062 {
					return 0
				}
				return 1
			},
		),
		loggingdriver.WithDsnParser(ParseDSN),
	)
}

// Open return
// dsn: mysql://[user[:password]@][addr]/database[?param1=value1&paramN=valueN]
func (d *mycli) Open(ctx context.Context, dsn *url.URL) (adaptor.Adaptor, error) {
	if dsn.Scheme != d.DriverName() {
		return nil, fmt.Errorf("invalid dsn schema, expect '%s' but got '%s'", d.DriverName(), dsn)
	}

	database := DatabaseNameFromDSN(dsn)
	conn, err := d.Connector().OpenConnector(dsn.String())
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(conn)

	if err = db.PingContext(ctx); err != nil {
		if IsUnknownDatabase(err) {
			if err = d.CreateDatabase(ctx, *dsn, database); err != nil {
				return nil, err
			}
			return d.Open(ctx, dsn)
		}
		return nil, err
	}

	return &mycli{
		dialect:  dialect{},
		DB:       adaptor.Wrap(db, func(err error) error { return nil }),
		database: "",
	}, nil
}

func (d *mycli) Catalog(ctx context.Context) (builder.Catalog, error) {
	return ScanCatalog(ctx, d, d.database)
}

func (d *mycli) CreateDatabase(ctx context.Context, dsn url.URL, database string) error {
	dsn.Path = "/mysql"
	dsn.RawPath = "mysql"

	a, err := d.Open(ctx, &dsn)
	if err != nil {
		return err
	}
	defer a.Close()

	_, err = a.Exec(context.Background(), frag.Query("CREATE DATABASE ?", frag.Lit(database)))
	return err
}

func DatabaseNameFromDSN(u *url.URL) string {
	database := u.Path
	if len(database) > 0 && database[0] == '/' {
		database = database[1:]
	}
	return database
}

func IsUnknownDatabase(err error) bool {
	var e *mysql.MySQLError
	if errors.As(err, &e) && e.Number == 1049 {
		return true
	}
	return false
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
