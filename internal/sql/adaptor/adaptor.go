package adaptor

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/xoctopus/x/syncx"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

// DB executes SQL fragments against a database connection pool.
type DB interface {
	// D returns the underlying *sql.DB.
	D() *sql.DB

	Exec(context.Context, frag.Fragment) (sql.Result, error)
	Query(context.Context, frag.Fragment) (*sql.Rows, error)
	Tx(context.Context, func(context.Context) error) error
	Close() error
}

// Connector opens an Adaptor from a parsed DSN URL.
type Connector interface {
	Open(context.Context, *url.URL) (Adaptor, error)
}

// Adaptor is a driver-specific database access handle.
type Adaptor interface {
	DB

	// DriverName returns the registered driver scheme (for example "mysql").
	DriverName() string
	// Schema returns the current database/schema name.
	Schema() string

	Dialect() Dialect
	Catalog(context.Context) (builder.Catalog, error)
}

// Dialect renders dialect-specific DDL/DML fragments and classifies driver errors.
type Dialect interface {
	CreateSchema(string) frag.Fragment
	SwitchSchema(string) frag.Fragment

	CreateTableIfNotExists(t builder.Table) []frag.Fragment
	DropTable(t builder.Table) frag.Fragment
	TruncateTable(t builder.Table) frag.Fragment

	AddColumn(builder.Col) frag.Fragment
	DropColumn(builder.Col) frag.Fragment
	RenameColumn(builder.Col, builder.Col) frag.Fragment
	ModifyColumn(builder.Col, builder.Col) frag.Fragment

	AddIndex(key builder.Key) frag.Fragment
	DropIndex(key builder.Key) frag.Fragment

	DBType(builder.ColumnDef) frag.Fragment
	ColDefine(dd *builder.ColumnDef)
	IsUnknownDatabaseError(error) bool
	IsConflictError(err error) bool
}

var adaptors = syncx.NewXmap[string, Adaptor]()

// Register stores an Adaptor under its DriverName and optional aliases.
func Register(a Adaptor, aliases ...string) {
	adaptors.Store(a.DriverName(), a)
	for _, alias := range aliases {
		adaptors.Store(alias, a)
	}
}

// Open parses dsn and opens an Adaptor whose scheme matches a registered driver.
func Open(ctx context.Context, dsn string) (Adaptor, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	for driver, adaptor := range adaptors.Range {
		if driver == u.Scheme {
			return adaptor.(Connector).Open(ctx, u)
		}
	}
	return nil, fmt.Errorf("missing adaptor: %s", dsn)
}
