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

type DB interface {
	Exec(context.Context, frag.Fragment) (sql.Result, error)
	Query(context.Context, frag.Fragment) (*sql.Rows, error)
	Tx(context.Context, func(context.Context) error) error
	Close() error
}

type Connector interface {
	Open(context.Context, *url.URL) (Adaptor, error)
}

type Adaptor interface {
	DB

	DriverName() string
	Dialect() Dialect
	Catalog(context.Context) (builder.Catalog, error)
}

type Dialect interface {
	CreateTableIsNotExists(t builder.Table) []frag.Fragment
	DropTable(t builder.Table) frag.Fragment
	TruncateTable(t builder.Table) frag.Fragment
	AddColumn(builder.Col) frag.Fragment
	DropColumn(builder.Col) frag.Fragment
	RenameColumn(builder.Col, builder.Col) frag.Fragment
	ModifyColumn(builder.Col, builder.Col) frag.Fragment
	AddIndex(key builder.Key) frag.Fragment
	DropIndex(key builder.Key) frag.Fragment
	DBType(builder.ColumnDef) frag.Fragment
}

var adaptors = syncx.NewXmap[string, Adaptor]()

func Register(a Adaptor, aliases ...string) {
	adaptors.Store(a.DriverName(), a)
	for _, alias := range aliases {
		adaptors.Store(alias, a)
	}
}

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
	return nil, fmt.Errorf("missing adaptor for %s", u.Scheme)
}
