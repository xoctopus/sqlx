package helper_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/xoctopus/x/testx"
	"github.com/xoctopus/x/testx/bdd"

	"github.com/xoctopus/sqlx/pkg/builder"
	sqlerrs "github.com/xoctopus/sqlx/pkg/errors"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/helper"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
)

type M struct {
	ID int    `db:"id,autoinc"`
	V  string `db:"v,width=32"`
}

func (M) TableName() string {
	return "helper_test"
}

func TestCVsForInsertion(t *testing.T) {
	cols, vals := helper.CVsForInsertion(
		M{ID: 1, V: "1"},
		M{ID: 2, V: "2"},
		M{ID: 3, V: "3"},
	)
	Expect(t, cols.Len(), Equal(1))
	Expect(t, vals, Equal([]any{"1", "2", "3"}))

	cols, vals = helper.CVsForInsertion(
		&M{ID: 1, V: "1"},
		&M{ID: 2, V: "2"},
		&M{ID: 3, V: "3"},
	)
	Expect(t, cols.Len(), Equal(1))
	Expect(t, vals, Equal([]any{"1", "2", "3"}))

	cols, vals = helper.CVsForInsertion[M]()
	Expect(t, cols, BeNil[builder.Cols]())
	Expect(t, len(vals), Equal(0))
}

type stubAdaptor struct {
	db *sql.DB
}

func (a *stubAdaptor) D() *sql.DB { return a.db }

func (a *stubAdaptor) Close() error { return a.db.Close() }

func (a *stubAdaptor) Exec(context.Context, frag.Fragment) (sql.Result, error) {
	panic("not implemented")
}

func (a *stubAdaptor) Query(ctx context.Context, f frag.Fragment) (*sql.Rows, error) {
	if frag.IsNil(f) {
		return nil, nil
	}
	q, args := frag.Collect(ctx, f)
	return a.db.QueryContext(ctx, q, args...)
}

func (a *stubAdaptor) Tx(context.Context, func(context.Context) error) error {
	panic("not implemented")
}

func (a *stubAdaptor) DriverName() string { return "sqlmock" }

func (a *stubAdaptor) Schema() string { return "test" }

func (a *stubAdaptor) Dialect() adaptor.Dialect { return nil }

func (a *stubAdaptor) Catalog(context.Context) (builder.Catalog, error) {
	return nil, nil
}

func TestScan(t *testing.T) {
	bdd.From(t).Given("rows with one record", func(t bdd.T) {
		db, mock, _ := sqlmock.New()
		t.Cleanup(func() { _ = db.Close() })

		_ = mock.ExpectQuery("SELECT id, v FROM helper_test").WillReturnRows(
			mock.NewRows([]string{"id", "v"}).AddRow(1, "one"),
		)

		t.When("scanning into a struct pointer", func(t bdd.T) {
			rows, qerr := db.Query("SELECT id, v FROM helper_test")
			dst := &M{}
			err := helper.Scan(context.Background(), rows, dst)

			t.Then(
				"it should succeed and fill fields",
				bdd.Succeed(qerr),
				bdd.Succeed(err),
				bdd.Equal(dst.ID, 1),
				bdd.Equal(dst.V, "one"),
			)
		})
	})

	bdd.From(t).Given("rows with no records", func(t bdd.T) {
		db, mock, _ := sqlmock.New()
		t.Cleanup(func() { _ = db.Close() })

		_ = mock.ExpectQuery("SELECT id, v FROM helper_test").WillReturnRows(
			mock.NewRows([]string{"id", "v"}),
		)

		t.When("scanning into a struct pointer", func(t bdd.T) {
			rows, qerr := db.Query("SELECT id, v FROM helper_test")
			err := helper.Scan(context.Background(), rows, &M{})

			t.Then(
				"it should return NOTFOUND",
				bdd.Succeed(qerr),
				bdd.BeTrue(sqlerrs.IsErrNotFound(err)),
			)
		})
	})

	bdd.From(t).Given("rows is nil", func(t bdd.T) {
		t.When("scanning", func(t bdd.T) {
			err := helper.Scan(context.Background(), nil, &M{})
			t.Then("it should succeed", bdd.Succeed(err))
		})
	})
}

func TestQueryAndScan(t *testing.T) {
	bdd.From(t).Given("adaptor query returns one record", func(t bdd.T) {
		db, mock, _ := sqlmock.New()
		t.Cleanup(func() { _ = db.Close() })

		_ = mock.ExpectQuery("SELECT id, v FROM helper_test").WillReturnRows(
			mock.NewRows([]string{"id", "v"}).AddRow(2, "two"),
		)
		a := &stubAdaptor{db: db}
		f := frag.Query("SELECT id, v FROM helper_test")

		t.When("query and scan into a struct pointer", func(t bdd.T) {
			dst := &M{}
			err := helper.QueryAndScan(context.Background(), a, f, dst)

			t.Then(
				"it should succeed and fill fields",
				bdd.Succeed(err),
				bdd.Equal(dst.ID, 2),
				bdd.Equal(dst.V, "two"),
			)
		})
	})

	bdd.From(t).Given("adaptor query returns no records", func(t bdd.T) {
		db, mock, _ := sqlmock.New()
		t.Cleanup(func() { _ = db.Close() })

		_ = mock.ExpectQuery("SELECT id, v FROM helper_test").WillReturnRows(
			mock.NewRows([]string{"id", "v"}),
		)
		a := &stubAdaptor{db: db}
		f := frag.Query("SELECT id, v FROM helper_test")

		t.When("query and scan into a struct pointer", func(t bdd.T) {
			err := helper.QueryAndScan(context.Background(), a, f, &M{})
			t.Then("it should return NOTFOUND", bdd.BeTrue(sqlerrs.IsErrNotFound(err)))
		})
	})

	bdd.From(t).Given("adaptor query fails", func(t bdd.T) {
		db, mock, _ := sqlmock.New()
		t.Cleanup(func() { _ = db.Close() })

		_ = mock.ExpectQuery("SELECT id, v FROM helper_test").WillReturnError(fmt.Errorf("boom"))
		a := &stubAdaptor{db: db}
		f := frag.Query("SELECT id, v FROM helper_test")

		t.When("query and scan", func(t bdd.T) {
			err := helper.QueryAndScan(context.Background(), a, f, &M{})
			t.Then(
				"it should return the query error",
				bdd.Failed(err),
				bdd.ErrorContains(err, "boom"),
			)
		})
	})

	bdd.From(t).Given("adaptor query returns rows and dst is nil", func(t bdd.T) {
		db, mock, _ := sqlmock.New()
		t.Cleanup(func() { _ = db.Close() })

		_ = mock.ExpectQuery("SELECT id, v FROM helper_test").WillReturnRows(
			mock.NewRows([]string{"id", "v"}).AddRow(3, "three"),
		)
		a := &stubAdaptor{db: db}
		f := frag.Query("SELECT id, v FROM helper_test")

		t.When("query and scan with nil dst", func(t bdd.T) {
			err := helper.QueryAndScan(context.Background(), a, f, nil)
			t.Then("it should succeed without scanning", bdd.Succeed(err))
		})
	})
}
