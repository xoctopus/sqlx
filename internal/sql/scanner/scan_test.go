package scanner_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/xoctopus/pkgx/pkg/pkgx"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/sql/scanner"
	"github.com/xoctopus/sqlx/pkg/errors"
)

type T struct {
	I int    `db:"f_i"`
	S string `db:"f_s"`
}

type Any string

type T2 T

func (t *T2) ColumnReceivers() map[string]any {
	return map[string]any{
		"f_i": &t.I,
		"f_s": &t.S,
	}
}

type TDataList struct {
	Data []T
}

func (*TDataList) New() any {
	return &T{}
}

func (l *TDataList) Next(v any) error {
	t := v.(*T)
	l.Data = append(l.Data, *t)
	return nil
}

var ctx = pkgx.CtxLoadTests.With(context.Background(), true)

func BenchmarkScan(b *testing.B) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := "SELECT f_i,f_s from t"
	b.Run("ScanToStruct", func(b *testing.B) {
		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "4")

		_ = mock.ExpectQuery(query).WillReturnRows(mockRows)

		target := &T{}

		for _ = range b.N {
			rows, _ := db.Query(query)
			_ = scanner.Scan(ctx, rows, target)
		}

		b.Log(target)
	})

	b.Run("ScanToStructWithColumnReceivers", func(b *testing.B) {
		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "4")

		_ = mock.ExpectQuery(query).WillReturnRows(mockRows)

		target := &T2{}

		for _ = range b.N {
			rows, _ := db.Query(query)
			_ = scanner.Scan(ctx, rows, target)
		}

		b.Log(target)
	})
}

func TestScan(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := "SELECT f_i,f_s from t"

	t.Run("ScanToStruct", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "4")
		_ = mock.ExpectQuery(query).WillReturnRows(mockRows)

		target := &T{}
		rows, _ := db.Query(query)
		err := scanner.Scan(ctx, rows, target)
		Expect(t, err, Succeed())
		Expect(t, target, Equal(&T{I: 2, S: "4"}))
	})

	t.Run("ScanToStructWithColumnReceivers", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "4")
		_ = mock.ExpectQuery(query).WillReturnRows(mockRows)

		target := &T2{}
		rows, _ := db.Query(query)
		err := scanner.Scan(ctx, rows, target)
		Expect(t, err, Succeed())
		Expect(t, target, Equal(&T2{I: 2, S: "4"}))
	})

	t.Run("ScanToStructNoRecords", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		_ = mock.ExpectQuery(query).WillReturnRows(mockRows)

		target := &T{}
		rows, err := db.Query(query)
		Expect(t, err, Succeed())

		err = scanner.Scan(ctx, rows, target)
		Expect(t, errors.IsErrNotFound(err), Be(true))
	})

	t.Run("ScanCount", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"count(1)"})
		mockRows.AddRow(10)
		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		count := 0
		rows, err := db.Query("SELECT count(1) from t")
		Expect(t, err, BeNil[error]())

		err = scanner.Scan(ctx, rows, &count)
		Expect(t, err, BeNil[error]())
		Expect(t, count, Equal(10))
	})

	t.Run("ScanCountBadReceiver", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"count(1)"})
		mockRows.AddRow(10)
		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		v := Any("")
		rows, err := db.Query("SELECT count(1) from t")
		Expect(t, err, Be[error](nil))

		err = scanner.Scan(ctx, rows, &v)
		Expect(t, err, Not(Be[error](nil)))
	})

	t.Run("ScanToSlice", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "2")
		mockRows.AddRow(3, "3")
		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		list := make([]T, 0)
		rows, err := db.Query("SELECT f_i,f_b from t")
		Expect(t, err, BeNil[error]())

		err = scanner.Scan(ctx, rows, &list)

		Expect(t, err, BeNil[error]())
		Expect(t, list, Equal([]T{
			{I: 2, S: "2"},
			{I: 3, S: "3"},
		}))
	})

	t.Run("ScanToIterator", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "2")
		mockRows.AddRow(3, "3")

		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		rows, err := db.Query("SELECT f_i,f_b from t")
		Expect(t, err, Be[error](nil))

		list := TDataList{}

		err = scanner.Scan(ctx, rows, &list)

		Expect(t, err, Be[error](nil))
		Expect(t, list.Data, Equal([]T{
			{I: 2, S: "2"},
			{I: 3, S: "3"},
		}))
	})
}
