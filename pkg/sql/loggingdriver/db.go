package loggingdriver

import (
	"context"
	"database/sql/driver"
	"fmt"
	"net/url"
	"time"

	"github.com/xoctopus/logx"
)

func Wrap(d driver.Driver, name string, opts ...DriverOptionApplier) driver.DriverContext {
	c := &connector{d: d, name: name}
	for _, applier := range opts {
		applier(&c.options)
	}

	return c
}

type connector struct {
	d    driver.Driver
	name string
	dsn  string
	options
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

	dsn = u.String()
	if c.DSNParser != nil {
		dsn, err = c.DSNParser(dsn)
		if err != nil {
			return nil, err
		}
	}

	return &connector{
		d:    c.d,
		name: name,
		dsn:  dsn,
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
		Conn:         conn,
		name:         c.name,
		level:        c.ErrorLevel,
		interpolator: c.ValueHolderReplacer,
	}, nil
}

type connection struct {
	driver.Conn
	name         string
	level        func(error) int
	interpolator Interpolator
}

func (c *connection) Prepare(q string) (driver.Stmt, error) {
	panic("forbidden") // to universe dialects
}

func (c *connection) Close() error {
	return c.Conn.Close()
}

func (c *connection) ErrorLevel(err error) int {
	if c.level != nil {
		return c.level(err)
	}
	return 1
}

func (c *connection) QueryContext(ctx context.Context, q string, args []driver.NamedValue) (rows driver.Rows, err error) {
	_, log := logx.Enter(ctx)
	span := Cost()

	if c.interpolator == nil {
		c.interpolator = DefaultInterpolate
	}
	q, args = c.interpolator(q, args)

	defer func() {
		millis := span().Milliseconds()
		printer := NewPrinter(q, args)
		log = log.With("driver", c.name, "query", printer.String(), "cost_ms", millis)
		if err != nil {
			if c.ErrorLevel(err) > 0 {
				log.Error(fmt.Errorf("query failed: %w", err))
			} else {
				log.Warn(fmt.Errorf("query failed: %w", err))
			}
		} else {
			log.Debug("")
		}
		log.End()
	}()

	// mysql set InterpolateParams default to false.
	rows, err = c.Conn.(driver.QueryerContext).QueryContext(ctx, q, args)
	return rows, err
}

func (c *connection) ExecContext(ctx context.Context, q string, args []driver.NamedValue) (res driver.Result, err error) {
	_, log := logx.Enter(ctx)
	span := Cost()

	if c.interpolator == nil {
		c.interpolator = DefaultInterpolate
	}
	q, args = c.interpolator(q, args)

	defer func() {
		millis := span().Milliseconds()
		printer := NewPrinter(q, args)
		log = log.With("driver", c.name, "query", printer.String(), "cost_ms", millis)
		if err != nil {
			if c.ErrorLevel(err) > 0 {
				log.Error(fmt.Errorf("exec failed: %w", err))
			} else {
				log.Warn(fmt.Errorf("exec failed: %w", err))
			}
		} else {
			log.Debug("")
		}
		log.End()
	}()

	res, err = c.Conn.(driver.ExecerContext).ExecContext(ctx, q, args)
	return res, err
}

func (c *connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	log := logx.From(ctx)

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

func Cost() func() time.Duration {
	ts := time.Now()
	return func() time.Duration {
		return time.Since(ts)
	}
}
