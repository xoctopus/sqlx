package loggingdriver

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"net/url"
	"strconv"

	"github.com/xoctopus/logx"
)

type connector struct {
	d    driver.Driver
	dsn  string
	name string
}

func (c *connector) OpenConnector(dsn string) (driver.Connector, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	name := c.name
	if q.Get("_ro") == "true" {
		name += "_ro"
		q.Del("_ro")
		u.RawQuery = q.Encode()
	}

	return &connector{
		d:    c.d,
		name: name,
		dsn:  u.String(),
	}, nil
}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	return c.Open(c.dsn)
}

func (c *connector) Driver() driver.Driver {
	return c
}

func (c *connector) Open(dsn string) (driver.Conn, error) {
	conn, err := c.d.Open(c.dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w. %s", err, dsn)
	}
	return &connection{
		Conn: conn,
		name: c.name,
	}, nil
}

type connection struct {
	driver.Conn
	name        string
	placeholder string
}

func (c *connection) Prepare(q string) (driver.Stmt, error) {
	panic("forbidden") // to universe dialects
}

func (c *connection) Close() error {
	return c.Conn.Close()
}

func (c *connection) QueryContext(ctx context.Context, q string, args []driver.NamedValue) (rows driver.Rows, err error) {
	_, log := logx.Enter(ctx)
	span := Span()

	defer func() {
		microseconds := span().Microseconds()
		printer := Interpolator(q, args)
		log = log.With("driver", c.name, "query", printer, "cost[µs]", microseconds)
		if err != nil {
			log.Error(fmt.Errorf("query failed: %w", err))
		} else {
			log.Debug("")
		}
		log.End()
	}()

	rows, err = c.Conn.(driver.QueryerContext).QueryContext(ctx, c.prepare(q), args)
	return rows, err
}

func (c *connection) ExecContext(ctx context.Context, q string, args []driver.NamedValue) (res driver.Result, err error) {
	_, log := logx.Enter(ctx)
	span := Span()

	defer func() {
		microseconds := span().Microseconds()
		printer := Interpolator(q, args)
		log = log.With("driver", c.name, "query", printer, "cost[µs]", microseconds)
		if err != nil {
			log.Error(fmt.Errorf("exec failed: %w", err))
		} else {
			log.Debug("")
		}
		log.End()
	}()

	res, err = c.Conn.(driver.ExecerContext).ExecContext(ctx, c.prepare(q), args)
	return res, err
}

func (c *connection) prepare(q string) string {
	if len(c.placeholder) == 0 {
		return q
	}

	b := bytes.NewBuffer(nil)
	for i := range q {
		switch v := q[i]; v {
		case '?':
			b.WriteString(c.placeholder)
			b.WriteString(strconv.FormatInt(int64(i+1), 10))
			i++
		default:
			b.WriteByte(v)
		}
	}

	return b.String()
}

func (c *connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	log := logx.FromContext(ctx)

	log.Debug("=========== Transaction Begin     ===========")
	// don't pass ctx into real driver to avoid connect discount
	tx, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if err != nil {
		log.Error(fmt.Errorf("failed to begin transaction: %w", err))
		return nil, err
	}

	return &transaction{tx: tx, log: log}, nil
}

type transaction struct {
	log logx.Logger
	tx  driver.Tx
}

func (tx *transaction) Commit() error {
	if err := tx.tx.Commit(); err != nil {
		tx.log.Debug("failed to commit transaction: %s", err)
		return err
	}
	tx.log.Debug("=========== Transaction Committed ===========")
	return nil
}

func (tx *transaction) Rollback() error {
	if err := tx.tx.Rollback(); err != nil {
		tx.log.Debug("failed to rollback transaction: %s", err)
		return err
	}
	tx.log.Debug("=========== Transaction Rollback  ===========")
	return nil
}
