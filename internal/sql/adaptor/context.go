package adaptor

import (
	"context"
	"database/sql"
)

// Executor can run SQL with context, typically *sql.DB or *sql.Tx.
type Executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

type tCtxExecutor struct{}

// WithExecutor binds e into ctx for nested Exec/Query/Tx reuse.
func WithExecutor(ctx context.Context, e Executor) context.Context {
	return context.WithValue(ctx, tCtxExecutor{}, e)
}

// ExecutorFrom returns the Executor bound by WithExecutor, if any.
func ExecutorFrom(ctx context.Context) Executor {
	e, ok := ctx.Value(tCtxExecutor{}).(Executor)
	if ok {
		return e
	}
	return nil
}
