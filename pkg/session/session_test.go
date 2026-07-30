package session_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/session"
)

type recordingAdaptor struct {
	name    string
	execN   int
	queryN  int
	txN     int
	execErr error
	lastF   frag.Fragment
}

func (a *recordingAdaptor) D() *sql.DB   { return nil }
func (a *recordingAdaptor) Close() error { return nil }
func (a *recordingAdaptor) DriverName() string {
	if a.name != "" {
		return a.name
	}
	return "rw"
}
func (a *recordingAdaptor) Schema() string           { return "test" }
func (a *recordingAdaptor) Dialect() adaptor.Dialect { return nil }
func (a *recordingAdaptor) Catalog(context.Context) (builder.Catalog, error) {
	return nil, nil
}

func (a *recordingAdaptor) Tx(ctx context.Context, fn func(context.Context) error) error {
	a.txN++
	return fn(ctx)
}

func (a *recordingAdaptor) Exec(_ context.Context, f frag.Fragment) (sql.Result, error) {
	a.execN++
	a.lastF = f
	if a.execErr != nil {
		return nil, a.execErr
	}
	return execResult{}, nil
}

func (a *recordingAdaptor) Query(_ context.Context, f frag.Fragment) (*sql.Rows, error) {
	a.queryN++
	a.lastF = f
	return nil, nil
}

type execResult struct{}

func (execResult) LastInsertId() (int64, error) { return 0, nil }
func (execResult) RowsAffected() (int64, error) { return 0, nil }

type withTable struct {
	tab builder.Table
}

func (w withTable) T() builder.Table { return w.tab }

type sessionModel struct {
	ID int64 `db:"f_id"`
}

func (sessionModel) TableName() string { return "t_session_model" }

func TestSession(t *testing.T) {
	rw := &recordingAdaptor{name: "rw"}
	ro := &recordingAdaptor{name: "ro"}

	t.Run("New", func(t *testing.T) {
		s := session.New(rw, "main")
		Expect(t, s.Name(), Equal("main"))
		Expect(t, s.Adaptor(), Be[adaptor.Adaptor](rw))
		Expect(t, s.Adaptor(session.ReadOnly()), BeNil[adaptor.Adaptor]())
	})

	t.Run("NewReadonly", func(t *testing.T) {
		s := session.NewReadonly(rw, ro, "ro-session")
		Expect(t, s.Name(), Equal("ro-session"))
		Expect(t, s.Adaptor(), Be[adaptor.Adaptor](rw))
		Expect(t, s.Adaptor(session.ReadOnly()), Be[adaptor.Adaptor](ro))
	})

	t.Run("T", func(t *testing.T) {
		s := session.New(rw, "main")
		tab := builder.T("t_direct")

		Expect(t, s.T(tab).TableName(), Equal("t_direct"))
		Expect(t, s.T(withTable{tab: tab}).TableName(), Equal("t_direct"))
		Expect(t, s.T(&sessionModel{}).TableName(), Equal("t_session_model"))
	})

	t.Run("TxExecQuery", func(t *testing.T) {
		rw.execN, rw.queryN, rw.txN = 0, 0, 0
		s := session.New(rw, "main")
		f := frag.Query("SELECT 1")

		Expect(t, s.Tx(context.Background(), func(ctx context.Context) error {
			Expect(t, session.InTx(ctx), BeFalse())
			return nil
		}), Succeed())
		Expect(t, rw.txN, Equal(1))

		_, err := s.Exec(context.Background(), f)
		Expect(t, err, Succeed())
		Expect(t, rw.execN, Equal(1))
		Expect(t, rw.lastF, Be[frag.Fragment](f))

		_, err = s.Query(context.Background(), f)
		Expect(t, err, Succeed())
		Expect(t, rw.queryN, Equal(1))

		rw.execErr = fmt.Errorf("exec boom")
		_, err = s.Exec(context.Background(), f)
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("exec boom"))
		rw.execErr = nil
	})
}
