package session

import (
	"context"
	"fmt"
	"net/url"

	"github.com/xoctopus/datatypex"
	"github.com/xoctopus/x/flagx"

	"github.com/xoctopus/sqlx/internal/diff"
	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	_ "github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/migrator"
)

type Database struct {
	Endpoint datatypex.Endpoint
	Readonly datatypex.Endpoint

	AutoMigration   bool
	DryRun          bool
	CreateTableOnly bool

	name    string
	catalog builder.Catalog

	db adaptor.Adaptor
	ro adaptor.Adaptor
}

func (d *Database) SetDefault() {}

// ApplyCatalog should do before endpoint initialization
func (d *Database) ApplyCatalog(name string, catalogs ...builder.Catalog) {
	d.name = name
	d.catalog = builder.NewCatalog()

	for _, catalog := range catalogs {
		for table := range catalog.Tables() {
			d.catalog.Add(table)
		}
	}
}

func (d *Database) Init(ctx context.Context) error {
	if d.db != nil {
		return nil
	}

	main := d.Endpoint
	db, err := adaptor.Open(ctx, main.String())
	if err != nil {
		return err
	}

	d.db = db

	if !d.Readonly.IsZero() {
		// readonly endpoint
		ro := d.Readonly
		// reuse main configurations
		if ro.Username == "" {
			ro.Username = main.Username
			ro.Password = main.Password
		}
		if ro.Param == nil {
			ro.Param = url.Values{}
		}
		ro.Param.Set("_ro", "true")
		db, err = adaptor.Open(ctx, ro.String())
		if err != nil {
			return err
		}
		d.ro = db
	}

	register(d.Name(), d.catalog)

	return nil
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) Session() Session {
	if d.ro != nil {
		return NewRO(d.db, d.ro, d.Name())
	}
	return New(d.db, d.Name())
}

func (d *Database) Catalog() builder.Catalog {
	return d.catalog
}

func (d *Database) Run(ctx context.Context) error {
	if d.AutoMigration {
		f := flagx.NewFlag[diff.Mode]()
		if d.DryRun {
			f.With(diff.MODE_DRY_RUN)
		}
		if d.CreateTableOnly {
			f.With(diff.MODE_CREATE_TABLE)
		}
		ctx = diff.CtxMode.With(ctx, f)
		q, err := migrator.Migrate(ctx, d.db, d.catalog)
		if err != nil {
			return err
		}
		fmt.Println(q)
	}
	return nil
}
