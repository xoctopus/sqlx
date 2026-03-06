package adaptor

import (
	"context"
	"database/sql"
)

type Executor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

type ctxExecutor struct{}

func WithExecutor(ctx context.Context, e Executor) context.Context {
	return context.WithValue(ctx, ctxExecutor{}, e)
}

func ExecutorFrom(ctx context.Context) Executor {
	e, ok := ctx.Value(ctxExecutor{}).(Executor)
	if ok {
		return e
	}
	return nil
}
