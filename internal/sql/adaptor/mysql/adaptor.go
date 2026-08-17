package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"net/url"
	"slices"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/xoctopus/x/codex"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/internal/sql/loggingdriver"
	"github.com/xoctopus/sqlx/pkg/builder"
	sqlerr "github.com/xoctopus/sqlx/pkg/errors"
	"github.com/xoctopus/sqlx/pkg/frag"
)

func init() {
	// HACK: temporarily treat TiDB as MySQL.
	// TiDB speaks the MySQL wire protocol and is largely compatible with the
	// mysql Dialect/Catalog, so register "tidb" as a DSN scheme alias instead
	// of maintaining a separate adaptor. DriverName() remains "mysql" so type
	// mapping and dialect behavior stay on the mysql path.
	// TODO: remove once a dedicated tidb adaptor exists, or if TiDB-specific
	// DDL/catalog differences need first-class handling.
	aliases := []string{"tidb"}
	adaptor.Register(&mycli{aliases: aliases}, aliases...)
}

type mycli struct {
	dialect
	adaptor.DB

	aliases []string
	schema  string
	dsn     *url.URL
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
	name := d.DriverName()
	if d.dsn != nil && d.dsn.Scheme != "" {
		name = d.dsn.Scheme
	}
	return loggingdriver.Wrap(mysql.MySQLDriver{}, name, options...)
}

// Open returns mysql adaptor.Adaptor
// dsn: mysql://[user[:password]@][addr]/database[?param1=value1&paramN=valueN]
func (d *mycli) Open(ctx context.Context, dsn *url.URL) (adaptor.Adaptor, error) {
	scheme := dsn.Scheme
	names := append(d.aliases, d.DriverName())
	must.BeTrueF(
		slices.Contains(names, scheme),
		"invalid dsn schema, expect %v but got '%s'", names, dsn,
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
		schema:  database,
		dsn:     dsn,
		aliases: d.aliases,
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
	conf := mysql.NewConfig()
	conf.User = u.User.Username()
	conf.Passwd = pass
	conf.Net = "tcp"
	conf.Addr = u.Host
	conf.DBName = strings.TrimPrefix(u.Path, "/")
	conf.Params = make(map[string]string)
	for k, vs := range u.Query() {
		conf.Params[k] = vs[0]
	}
	return conf.FormatDSN(), nil
}
