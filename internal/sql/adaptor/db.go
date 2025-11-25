package adaptor

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func Wrap(d *sql.DB, ew func(error) error) DB {
	_db := &db{
		DB:  d,
		err: func(err error) error { return err },
	}
	if ew != nil {
		_db.err = ew
	}
	return _db
}

type db struct {
	*sql.DB
	err func(error) error
}

func (d *db) Exec(ctx context.Context, f frag.Fragment) (sql.Result, error) {
	if frag.IsNil(f) {
		return nil, nil
	}
	q, args := frag.Collect(ctx, f)
	if exec := ExecutorFrom(ctx); exec != nil {
		result, err := exec.ExecContext(ctx, q, args...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	result, err := d.ExecContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *db) Query(ctx context.Context, f frag.Fragment) (*sql.Rows, error) {
	if frag.IsNil(f) {
		return nil, nil
	}
	q, args := frag.Collect(ctx, f)
	if exec := ExecutorFrom(ctx); exec != nil {
		result, err := exec.QueryContext(ctx, q, args...)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	result, err := d.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *db) Tx(ctx context.Context, f func(context.Context) error) (err error) {
	var (
		entry = false
		tx    *sql.Tx
	)

	if exec := ExecutorFrom(ctx); exec != nil {
		if _tx, ok := exec.(*sql.Tx); ok {
			tx = _tx
		}
	}

	if tx == nil {
		_tx, _err := d.BeginTx(ctx, nil)
		if _err != nil {
			return _err
		}
		tx = _tx
		entry = true
	}

	defer func() {
		if caught := recover(); caught != nil {
			rollback := tx.Rollback() // make sure rollback
			cause := func(caught, rollback error) error {
				if rollback == nil {
					return fmt.Errorf("cause: %w", caught)
				}
				return fmt.Errorf("caught: %w rollback: %v", caught, rollback)
			}
			switch e := caught.(type) {
			case runtime.Error:
				panic(cause(e, rollback))
			case error:
				if rollback != nil {
					err = cause(e, rollback)
				} else {
					err = e
				}
				return
			default:
				panic(cause(fmt.Errorf("%v", e), rollback))
			}
		}
		if entry {
			if err != nil {
				err = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}
	}()

	err = f(WithExecutor(ctx, tx))
	return err
}
