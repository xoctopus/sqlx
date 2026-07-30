package adaptor

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type recordingExecutor struct {
	execN   int
	queryN  int
	lastQ   string
	lastA   []any
	execErr error
	rows    *sql.Rows
}

func (e *recordingExecutor) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	e.execN++
	e.lastQ = query
	e.lastA = args
	if e.execErr != nil {
		return nil, e.execErr
	}
	return sqlmock.NewResult(1, 1), nil
}

func (e *recordingExecutor) QueryContext(_ context.Context, query string, args ...any) (*sql.Rows, error) {
	e.queryN++
	e.lastQ = query
	e.lastA = args
	return e.rows, nil
}

func TestWrap(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	Expect(t, err, Succeed())
	t.Cleanup(func() {
		mock.ExpectClose()
		_ = sqlDB.Close()
	})

	wrapped := Wrap(sqlDB, nil)
	Expect(t, wrapped.D(), Be(sqlDB))

	wrapped = Wrap(sqlDB, func(err error) error {
		return fmt.Errorf("wrapped: %w", err)
	})
	Expect(t, wrapped.D(), Be(sqlDB))
}

func TestDBExec(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	Expect(t, err, Succeed())
	t.Cleanup(func() {
		mock.ExpectClose()
		_ = sqlDB.Close()
	})
	d := Wrap(sqlDB, nil)

	t.Run("NilFragment", func(t *testing.T) {
		result, err := d.Exec(context.Background(), nil)
		Expect(t, err, Succeed())
		Expect(t, result, BeNil[sql.Result]())
	})

	t.Run("ViaDB", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO t").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
		result, err := d.Exec(context.Background(), frag.Query("INSERT INTO t VALUES (?)", 1))
		Expect(t, err, Succeed())
		Expect(t, result, NotBeNil[sql.Result]())
		Expect(t, mock.ExpectationsWereMet(), Succeed())
	})

	t.Run("ViaExecutor", func(t *testing.T) {
		exec := &recordingExecutor{}
		ctx := WithExecutor(context.Background(), exec)
		result, err := d.Exec(ctx, frag.Query("UPDATE t SET v = ?", 2))
		Expect(t, err, Succeed())
		Expect(t, result, NotBeNil[sql.Result]())
		Expect(t, exec.execN, Equal(1))
		Expect(t, exec.lastQ, Equal("UPDATE t SET v = ?"))
		Expect(t, exec.lastA, Equal([]any{2}))
	})

	t.Run("ExecutorError", func(t *testing.T) {
		exec := &recordingExecutor{execErr: fmt.Errorf("exec boom")}
		ctx := WithExecutor(context.Background(), exec)
		_, err := d.Exec(ctx, frag.Query("DELETE FROM t"))
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("exec boom"))
	})
}

func TestDBQuery(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	Expect(t, err, Succeed())
	t.Cleanup(func() {
		mock.ExpectClose()
		_ = sqlDB.Close()
	})
	d := Wrap(sqlDB, nil)

	t.Run("NilFragment", func(t *testing.T) {
		rows, err := d.Query(context.Background(), nil)
		Expect(t, err, Succeed())
		Expect(t, rows, BeNil[*sql.Rows]())
	})

	t.Run("ViaDB", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM t").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		rows, err := d.Query(context.Background(), frag.Query("SELECT id FROM t"))
		Expect(t, err, Succeed())
		Expect(t, rows, NotBeNil[*sql.Rows]())
		_ = rows.Close()
		Expect(t, mock.ExpectationsWereMet(), Succeed())
	})

	t.Run("ViaExecutor", func(t *testing.T) {
		exec := &recordingExecutor{}
		ctx := WithExecutor(context.Background(), exec)
		_, err := d.Query(ctx, frag.Query("SELECT ?", "x"))
		Expect(t, err, Succeed())
		Expect(t, exec.queryN, Equal(1))
		Expect(t, exec.lastQ, Equal("SELECT ?"))
		Expect(t, exec.lastA, Equal([]any{"x"}))
	})
}

func TestDBTx(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	Expect(t, err, Succeed())
	t.Cleanup(func() {
		mock.ExpectClose()
		_ = sqlDB.Close()
	})
	d := Wrap(sqlDB, nil)

	t.Run("Commit", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()
		var sawTx bool
		err := d.Tx(context.Background(), func(ctx context.Context) error {
			exec := ExecutorFrom(ctx)
			Expect(t, exec, NotBeNil[Executor]())
			_, ok := exec.(*sql.Tx)
			Expect(t, ok, BeTrue())
			sawTx = true
			return nil
		})
		Expect(t, err, Succeed())
		Expect(t, sawTx, BeTrue())
		Expect(t, mock.ExpectationsWereMet(), Succeed())
	})

	t.Run("RollbackOnError", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()
		err := d.Tx(context.Background(), func(ctx context.Context) error {
			return fmt.Errorf("tx boom")
		})
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("tx boom"))
		Expect(t, mock.ExpectationsWereMet(), Succeed())
	})

	t.Run("NestedReuseTx", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()
		err := d.Tx(context.Background(), func(ctx context.Context) error {
			return d.Tx(ctx, func(ctx context.Context) error {
				Expect(t, ExecutorFrom(ctx), NotBeNil[Executor]())
				return nil
			})
		})
		Expect(t, err, Succeed())
		Expect(t, mock.ExpectationsWereMet(), Succeed())
	})

	t.Run("PanicError", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()
		err := d.Tx(context.Background(), func(ctx context.Context) error {
			panic(fmt.Errorf("panic boom"))
		})
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("cause"))
		Expect(t, err, ErrorContains("panic boom"))
		Expect(t, mock.ExpectationsWereMet(), Succeed())
	})

	t.Run("PanicNonError", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()
		ExpectPanic[error](t, func() {
			_ = d.Tx(context.Background(), func(ctx context.Context) error {
				panic("raw panic")
			})
		}, ErrorContains("raw panic"))
		Expect(t, mock.ExpectationsWereMet(), Succeed())
	})
}
