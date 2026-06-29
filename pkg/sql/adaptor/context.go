package adaptor

import (
	"context"
	"database/sql"
)

type Executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

type tCtxExecutor struct{}

func WithExecutor(ctx context.Context, e Executor) context.Context {
	return context.WithValue(ctx, tCtxExecutor{}, e)
}

func ExecutorFrom(ctx context.Context) Executor {
	e, ok := ctx.Value(tCtxExecutor{}).(Executor)
	if ok {
		return e
	}
	return nil
}
