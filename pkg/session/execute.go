package session

import (
	"context"
	"database/sql"

	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
	"github.com/xoctopus/x/misc/must"
)

func InTx(ctx context.Context) bool {
	executor := adaptor.ExecutorFrom(ctx)
	if _, ok := executor.(*sql.Tx); ok {
		return true
	}
	return false
}

func ExecutorFor(ctx context.Context, m any) (adaptor.ExecutorX, bool) {
	if executor := adaptor.ExecutorXFrom(ctx); executor != nil {
		return executor, true
	}
	if s, ok := For(ctx, m); ok {
		return s.Adaptor(), true
	}
	return nil, false
}

func MustExecutorFor(ctx context.Context, m any) adaptor.ExecutorX {
	e, ok := ExecutorFor(ctx, m)
	must.BeTrueF(ok, "missing executor for %T", m)
	return e
}
