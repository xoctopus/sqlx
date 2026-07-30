package loggingdriver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/xoctopus/logx"
	. "github.com/xoctopus/x/testx"
)

type stubDriver struct {
	openDSN string
	openErr error
	conn    driver.Conn
}

func (d *stubDriver) Open(name string) (driver.Conn, error) {
	d.openDSN = name
	if d.openErr != nil {
		return nil, d.openErr
	}
	if d.conn == nil {
		d.conn = &stubConn{}
	}
	return d.conn, nil
}

type stubConn struct {
	lastQ    string
	lastArgs []driver.NamedValue
	queryErr error
	execErr  error
	beginErr error
	closed   bool
}

func (c *stubConn) Prepare(string) (driver.Stmt, error) { return nil, nil }
func (c *stubConn) Close() error                        { c.closed = true; return nil }
func (c *stubConn) Begin() (driver.Tx, error)           { return &stubTx{}, nil }

func (c *stubConn) QueryContext(_ context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	c.lastQ = q
	c.lastArgs = args
	if c.queryErr != nil {
		return nil, c.queryErr
	}
	return &stubRows{}, nil
}

func (c *stubConn) ExecContext(_ context.Context, q string, args []driver.NamedValue) (driver.Result, error) {
	c.lastQ = q
	c.lastArgs = args
	if c.execErr != nil {
		return nil, c.execErr
	}
	return stubResult{}, nil
}

func (c *stubConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	if c.beginErr != nil {
		return nil, c.beginErr
	}
	return &stubTx{}, nil
}

type stubRows struct{}

func (stubRows) Columns() []string { return nil }
func (stubRows) Close() error      { return nil }
func (stubRows) Next([]driver.Value) error {
	return io.EOF
}

type stubResult struct{}

func (stubResult) LastInsertId() (int64, error) { return 1, nil }
func (stubResult) RowsAffected() (int64, error) { return 1, nil }

type stubTx struct {
	committed  bool
	rolledBack bool
	commitErr  error
	rbErr      error
}

func (t *stubTx) Commit() error {
	if t.commitErr != nil {
		return t.commitErr
	}
	t.committed = true
	return nil
}

func (t *stubTx) Rollback() error {
	if t.rbErr != nil {
		return t.rbErr
	}
	t.rolledBack = true
	return nil
}

func discardCtx() context.Context {
	return logx.With(context.Background(), logx.Discard())
}

func openWrapped(t *testing.T, d *stubDriver, name, dsn string, opts ...DriverOptionApplier) driver.Conn {
	t.Helper()
	dc := Wrap(d, name, opts...)
	c, err := dc.OpenConnector(dsn)
	Expect(t, err, Succeed())
	conn, err := c.Connect(discardCtx())
	Expect(t, err, Succeed())
	return conn
}

func TestOpenConnector(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		d := &stubDriver{}
		dc := Wrap(d, "mysql")
		c, err := dc.OpenConnector("mysql://u:p@localhost:3306/db")
		Expect(t, err, Succeed())
		conn, err := c.Connect(discardCtx())
		Expect(t, err, Succeed())
		Expect(t, d.openDSN, Equal("mysql://u:p@localhost:3306/db"))
		Expect(t, conn.Close(), Succeed())
		Expect(t, d.conn.(*stubConn).closed, BeTrue())
	})

	t.Run("InvalidDSN", func(t *testing.T) {
		dc := Wrap(&stubDriver{}, "mysql")
		_, err := dc.OpenConnector("://bad")
		Expect(t, err, Failed())
	})

	t.Run("ReadonlySuffix", func(t *testing.T) {
		d := &stubDriver{}
		dc := Wrap(d, "mysql")
		c, err := dc.OpenConnector("mysql://u@h/db?_ro=true&x=1")
		Expect(t, err, Succeed())
		conn, err := c.Connect(discardCtx())
		Expect(t, err, Succeed())
		Expect(t, conn.(*connection).name, Equal("mysql::ro"))
		Expect(t, d.openDSN, ContainsSubString("x=1"))
		Expect(t, d.openDSN, Not(ContainsSubString("_ro=")))
	})

	t.Run("DSNParser", func(t *testing.T) {
		d := &stubDriver{}
		dc := Wrap(d, "mysql", WithDsnParser(func(dsn string) (string, error) {
			return "parsed:" + dsn, nil
		}))
		c, err := dc.OpenConnector("mysql://u@h/db")
		Expect(t, err, Succeed())
		_, err = c.Connect(discardCtx())
		Expect(t, err, Succeed())
		Expect(t, d.openDSN, Equal("parsed:mysql://u@h/db"))
	})

	t.Run("DSNParserError", func(t *testing.T) {
		dc := Wrap(&stubDriver{}, "mysql", WithDsnParser(func(string) (string, error) {
			return "", fmt.Errorf("parse boom")
		}))
		_, err := dc.OpenConnector("mysql://u@h/db")
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("parse boom"))
	})

	t.Run("OpenError", func(t *testing.T) {
		d := &stubDriver{openErr: fmt.Errorf("dial boom")}
		dc := Wrap(d, "mysql")
		c, err := dc.OpenConnector("mysql://u@h/db")
		Expect(t, err, Succeed())
		_, err = c.Connect(discardCtx())
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("failed to open connection"))
		Expect(t, err, ErrorContains("dial boom"))
	})
}

func TestConnection(t *testing.T) {
	t.Run("PreparePanics", func(t *testing.T) {
		conn := openWrapped(t, &stubDriver{}, "mysql", "mysql://u@h/db")
		ExpectPanic[string](t, func() {
			_, _ = conn.Prepare("SELECT 1")
		}, Equal("forbidden"))
	})

	t.Run("QueryContext", func(t *testing.T) {
		d := &stubDriver{}
		conn := openWrapped(t, d, "mysql", "mysql://u@h/db")
		rows, err := conn.(driver.QueryerContext).QueryContext(
			discardCtx(),
			"SELECT ?",
			[]driver.NamedValue{{Ordinal: 1, Value: 1}},
		)
		Expect(t, err, Succeed())
		Expect(t, rows, NotBeNil[driver.Rows]())
		Expect(t, d.conn.(*stubConn).lastQ, Equal("SELECT ?"))
		Expect(t, len(d.conn.(*stubConn).lastArgs), Equal(1))
	})

	t.Run("QueryContextError", func(t *testing.T) {
		d := &stubDriver{conn: &stubConn{queryErr: fmt.Errorf("query boom")}}
		conn := openWrapped(t, d, "mysql", "mysql://u@h/db")
		_, err := conn.(driver.QueryerContext).QueryContext(discardCtx(), "SELECT 1", nil)
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("query boom"))
	})

	t.Run("ExecContext", func(t *testing.T) {
		d := &stubDriver{}
		conn := openWrapped(t, d, "mysql", "mysql://u@h/db")
		res, err := conn.(driver.ExecerContext).ExecContext(
			discardCtx(),
			"UPDATE t SET v = ?",
			[]driver.NamedValue{{Ordinal: 1, Value: "x"}},
		)
		Expect(t, err, Succeed())
		Expect(t, res, NotBeNil[driver.Result]())
		Expect(t, d.conn.(*stubConn).lastQ, Equal("UPDATE t SET v = ?"))
	})

	t.Run("ExecContextError", func(t *testing.T) {
		d := &stubDriver{conn: &stubConn{execErr: fmt.Errorf("exec boom")}}
		conn := openWrapped(t, d, "mysql", "mysql://u@h/db")
		_, err := conn.(driver.ExecerContext).ExecContext(discardCtx(), "DELETE FROM t", nil)
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("exec boom"))
	})

	t.Run("Interpolator", func(t *testing.T) {
		d := &stubDriver{}
		conn := openWrapped(t, d, "mysql", "mysql://u@h/db", WithInterpolator(func(q string, args []driver.NamedValue) (string, []driver.NamedValue) {
			return "SELECT 1", nil
		}))
		_, err := conn.(driver.QueryerContext).QueryContext(
			discardCtx(),
			"SELECT ?",
			[]driver.NamedValue{{Value: 1}},
		)
		Expect(t, err, Succeed())
		Expect(t, d.conn.(*stubConn).lastQ, Equal("SELECT 1"))
		Expect(t, d.conn.(*stubConn).lastArgs, BeNil[[]driver.NamedValue]())
	})

	t.Run("ErrorLevel", func(t *testing.T) {
		conn := openWrapped(t, &stubDriver{}, "mysql", "mysql://u@h/db", WithErrorLeveler(func(error) int {
			return 0
		}))
		Expect(t, conn.(*connection).ErrorLevel(fmt.Errorf("x")), Equal(0))

		conn = openWrapped(t, &stubDriver{}, "mysql", "mysql://u@h/db")
		Expect(t, conn.(*connection).ErrorLevel(fmt.Errorf("x")), Equal(1))
	})
}

func TestBeginTx(t *testing.T) {
	t.Run("BeginError", func(t *testing.T) {
		d := &stubDriver{conn: &stubConn{beginErr: fmt.Errorf("begin boom")}}
		conn := openWrapped(t, d, "mysql", "mysql://u@h/db")
		_, err := conn.(driver.ConnBeginTx).BeginTx(discardCtx(), driver.TxOptions{})
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("begin boom"))
	})

	t.Run("Commit", func(t *testing.T) {
		stx := &stubTx{}
		d := &stubDriver{conn: &beginConn{tx: stx}}
		conn := openWrapped(t, d, "mysql", "mysql://u@h/db")
		tx, err := conn.(driver.ConnBeginTx).BeginTx(discardCtx(), driver.TxOptions{})
		Expect(t, err, Succeed())
		Expect(t, tx.Commit(), Succeed())
		Expect(t, stx.committed, BeTrue())
	})

	t.Run("Rollback", func(t *testing.T) {
		stx := &stubTx{}
		d := &stubDriver{conn: &beginConn{tx: stx}}
		conn := openWrapped(t, d, "mysql", "mysql://u@h/db")
		tx, err := conn.(driver.ConnBeginTx).BeginTx(discardCtx(), driver.TxOptions{})
		Expect(t, err, Succeed())
		Expect(t, tx.Rollback(), Succeed())
		Expect(t, stx.rolledBack, BeTrue())
	})

	t.Run("CommitError", func(t *testing.T) {
		stx := &stubTx{commitErr: fmt.Errorf("commit boom")}
		d := &stubDriver{conn: &beginConn{tx: stx}}
		conn := openWrapped(t, d, "mysql", "mysql://u@h/db")
		tx, err := conn.(driver.ConnBeginTx).BeginTx(discardCtx(), driver.TxOptions{})
		Expect(t, err, Succeed())
		err = tx.Commit()
		Expect(t, err, Failed())
		Expect(t, err, ErrorContains("commit boom"))
	})
}

type beginConn struct {
	stubConn
	tx *stubTx
}

func (c *beginConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return c.tx, nil
}

func TestCost(t *testing.T) {
	span := Cost()
	time.Sleep(2 * time.Millisecond)
	Expect(t, span() >= 2*time.Millisecond, BeTrue())
}
