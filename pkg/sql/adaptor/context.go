package adaptor

import (
	"context"
	"database/sql"

	"github.com/xoctopus/sqlx/pkg/frag"
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

type ExecutorX interface {
	Exec(context.Context, frag.Fragment) (sql.Result, error)
	Query(context.Context, frag.Fragment) (*sql.Rows, error)
	Tx(context.Context, func(context.Context) error) error
}

type tCtxExecutorX struct{}

func WithExecutorX(ctx context.Context, e ExecutorX) context.Context {
	return context.WithValue(ctx, tCtxExecutorX{}, e)
}

func ExecutorXFrom(ctx context.Context) ExecutorX {
	e, ok := ctx.Value(tCtxExecutorX{}).(ExecutorX)
	if ok {
		return e
	}
	return nil
}
