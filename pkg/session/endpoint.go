package session

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xoctopus/confx/pkg/types"
	"github.com/xoctopus/x/flagx"

	"github.com/xoctopus/sqlx/internal/diff"
	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	_ "github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/migrator"
)

type Endpoint struct {
	types.Endpoint[EndpointOption]
	Readonly types.Endpoint[EndpointOption]

	name    string
	catalog builder.Catalog
	db      adaptor.Adaptor
	ro      adaptor.Adaptor
}

// ApplyCatalog should do before endpoint initialization
func (d *Endpoint) ApplyCatalog(name string, catalogs ...builder.Catalog) {
	d.name = name
	d.catalog = builder.NewCatalog()

	for _, catalog := range catalogs {
		for table := range catalog.Tables() {
			d.catalog.Add(table)
		}
	}
}

func (d *Endpoint) Init(ctx context.Context) error {
	if d.db != nil {
		return nil
	}

	if err := d.Endpoint.Init(); err != nil {
		return fmt.Errorf("failed to init main endpoint: %w", err)
	}
	main := d.Endpoint
	db, err := adaptor.Open(ctx, main.String())
	if err != nil {
		return err
	}

	d.db = db

	if !d.Readonly.IsZero() {
		if !d.Readonly.IsZero() {
			if err = d.Readonly.Init(); err != nil {
				return fmt.Errorf("failed to init readonly endpoint: %w", err)
			}
		}
		// readonly endpoint
		ro := d.Readonly
		// reuse main configurations
		if ro.Auth.IsZero() {
			ro.Auth = main.Auth
		}
		ro.AddOption("_ro", "true")
		db, err = adaptor.Open(ctx, ro.String())
		if err != nil {
			return err
		}
		d.ro = db
	}

	register(d.Name(), d.catalog)

	if v := d.LivenessCheck(ctx); !v.Reachable {
		return fmt.Errorf("failed to probe liveness: %s", v.Message)
	}

	return nil
}

func (d *Endpoint) LivenessCheck(ctx context.Context) (v types.LivenessData) {
	var db adaptor.Adaptor
	if d.db != nil {
		db = d.db
	}
	if d.ro != nil {
		db = d.ro
	}
	if db == nil {
		v.Message = "connection lost"
		return
	}

	span := types.Cost()
	_, err := db.Query(ctx, frag.Query("SELECT 1"))
	cost := span()
	if err != nil {
		v.Message = err.Error()
		return
	}
	v.Reachable = true
	v.RTT = types.Duration(cost)
	return
}

func (d *Endpoint) Name() string {
	return d.name
}

func (d *Endpoint) Session() Session {
	if d.ro != nil {
		return NewRO(d.db, d.ro, d.Name())
	}
	return New(d.db, d.Name())
}

func (d *Endpoint) Catalog() builder.Catalog {
	return d.catalog
}

func (d *Endpoint) Run(ctx context.Context) error {
	o := d.Endpoint.Option

	if o.AutoMigration {
		f := flagx.NewFlag[diff.Mode]()
		if o.DryRun {
			f.With(diff.MODE_DRY_RUN)
		}
		if o.CreateTableOnly {
			f.With(diff.MODE_CREATE_TABLE)
		}
		ctx = diff.CtxMode.With(ctx, f)
		q, err := migrator.Migrate(ctx, d.db, d.catalog)
		fmt.Println(q)

		if dir, ok := migrator.CtxOutput.From(ctx); ok {
			filename := filepath.Join(
				dir,
				d.name+"_"+time.Now().Format("20060102_150405.000")+".sql",
			)
			_ = os.WriteFile(filename, []byte(q), 0o666)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Endpoint) Close() error {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			return err
		}
	}
	if d.ro != nil {
		if err := d.ro.Close(); err != nil {
			return err
		}
	}
	return nil
}
