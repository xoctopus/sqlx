package migrator_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/xoctopus/x/flagx"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/migrator"
)

type execResult struct{}

func (execResult) LastInsertId() (int64, error) { return 0, nil }
func (execResult) RowsAffected() (int64, error) { return 0, nil }

type stubDialect struct{}

func (stubDialect) CreateSchema(string) frag.Fragment { return frag.Query("") }
func (stubDialect) SwitchSchema(string) frag.Fragment { return frag.Query("") }
func (stubDialect) TruncateTable(builder.Table) frag.Fragment {
	return frag.Query("")
}
func (stubDialect) AddColumn(builder.Col) frag.Fragment  { return frag.Query("") }
func (stubDialect) DropColumn(builder.Col) frag.Fragment { return frag.Query("") }
func (stubDialect) RenameColumn(builder.Col, builder.Col) frag.Fragment {
	return frag.Query("")
}
func (stubDialect) ModifyColumn(builder.Col, builder.Col) frag.Fragment {
	return frag.Query("")
}
func (stubDialect) AddIndex(builder.Key) frag.Fragment  { return frag.Query("") }
func (stubDialect) DropIndex(builder.Key) frag.Fragment { return frag.Query("") }
func (stubDialect) DBType(builder.ColumnDef) frag.Fragment {
	return frag.Lit("INT")
}
func (stubDialect) ColDefine(*builder.ColumnDef)      {}
func (stubDialect) IsUnknownDatabaseError(error) bool { return false }
func (stubDialect) IsConflictError(error) bool        { return false }

func (stubDialect) DropTable(t builder.Table) frag.Fragment {
	return frag.Query("DROP TABLE IF EXISTS ?;", t)
}

func (stubDialect) CreateTableIfNotExists(t builder.Table) []frag.Fragment {
	return []frag.Fragment{frag.Query("CREATE TABLE IF NOT EXISTS ? ();", t)}
}

type stubAdaptor struct {
	dialect    adaptor.Dialect
	catalog    builder.Catalog
	catalogErr error
	execErr    error
	execN      int
	txN        int
}

func (a *stubAdaptor) D() *sql.DB   { return nil }
func (a *stubAdaptor) Close() error { return nil }
func (a *stubAdaptor) Query(context.Context, frag.Fragment) (*sql.Rows, error) {
	return nil, fmt.Errorf("not implemented")
}
func (a *stubAdaptor) DriverName() string { return "stub" }
func (a *stubAdaptor) Schema() string     { return "test" }
func (a *stubAdaptor) Dialect() adaptor.Dialect {
	if a.dialect != nil {
		return a.dialect
	}
	return stubDialect{}
}
func (a *stubAdaptor) Catalog(context.Context) (builder.Catalog, error) {
	if a.catalogErr != nil {
		return nil, a.catalogErr
	}
	if a.catalog != nil {
		return a.catalog, nil
	}
	return builder.NewCatalog(), nil
}

func (a *stubAdaptor) Tx(ctx context.Context, fn func(context.Context) error) error {
	a.txN++
	return fn(ctx)
}

func (a *stubAdaptor) Exec(ctx context.Context, f frag.Fragment) (sql.Result, error) {
	a.execN++
	if a.execErr != nil {
		return nil, a.execErr
	}
	return execResult{}, nil
}

func demoCatalog() builder.Catalog {
	id := builder.C(
		"f_id",
		builder.WithColFieldName("ID"),
		builder.WithColDefOf(int64(0), ",autoinc"),
	)
	tab := builder.T(
		"t_demo",
		id,
		builder.PK(builder.ColsOf(id)),
	)
	cat := builder.NewCatalog()
	cat.Add(tab)
	return cat
}

func dryRunCtx() context.Context {
	mode := flagx.NewFlag[migrator.Mode]()
	mode.With(migrator.DIFF_MODE_DRY_RUN)
	return migrator.CtxMode.With(context.Background(), mode)
}

func TestMigrate(t *testing.T) {
	t.Run("CatalogError", func(t *testing.T) {
		a := &stubAdaptor{catalogErr: fmt.Errorf("catalog boom")}
		q, err := migrator.Migrate(context.Background(), a, demoCatalog())
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("catalog boom"))
		Expect(t, q, Equal(""))
		Expect(t, a.txN, Equal(0))
		Expect(t, a.execN, Equal(0))
	})

	t.Run("DryRunCreateTable", func(t *testing.T) {
		a := &stubAdaptor{}
		q, err := migrator.Migrate(dryRunCtx(), a, demoCatalog())
		Expect(t, err, Succeed())
		Expect(t, q, ContainsSubString("CREATE TABLE IF NOT EXISTS t_demo"))
		Expect(t, q, ContainsSubString("DROP TABLE IF EXISTS sql_meta_table"))
		Expect(t, q, ContainsSubString("DROP TABLE IF EXISTS sql_meta_table_column"))
		Expect(t, q, ContainsSubString("DROP TABLE IF EXISTS sql_meta_enumeration"))
		Expect(t, a.txN, Equal(0))
		Expect(t, a.execN, Equal(0))
	})

	t.Run("DryRunNoSchemaChange", func(t *testing.T) {
		next := demoCatalog()
		a := &stubAdaptor{catalog: next}
		q, err := migrator.Migrate(dryRunCtx(), a, next)
		Expect(t, err, Succeed())
		Expect(t, q, Not(ContainsSubString("CREATE TABLE IF NOT EXISTS t_demo")))
		Expect(t, q, ContainsSubString("sql_meta_table"))
		Expect(t, a.txN, Equal(0))
		Expect(t, a.execN, Equal(0))
	})

	t.Run("Apply", func(t *testing.T) {
		a := &stubAdaptor{}
		q, err := migrator.Migrate(context.Background(), a, demoCatalog())
		Expect(t, err, Succeed())
		Expect(t, q, ContainsSubString("CREATE TABLE IF NOT EXISTS t_demo"))
		Expect(t, a.txN, Equal(1))
		Expect(t, a.execN > 0, BeTrue())
	})

	t.Run("ExecFailed", func(t *testing.T) {
		a := &stubAdaptor{execErr: fmt.Errorf("exec boom")}
		q, err := migrator.Migrate(context.Background(), a, demoCatalog())
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("migrate failed"))
		Expect(t, err, ErrorContains("exec boom"))
		Expect(t, q, ContainsSubString("CREATE TABLE IF NOT EXISTS t_demo"))
		Expect(t, a.txN, Equal(1))
	})

	t.Run("EmptyNextStillWritesMeta", func(t *testing.T) {
		a := &stubAdaptor{}
		q, err := migrator.Migrate(dryRunCtx(), a, builder.NewCatalog())
		Expect(t, err, Succeed())
		Expect(t, q, ContainsSubString("sql_meta_table"))
		Expect(t, a.txN, Equal(0))
	})
}
