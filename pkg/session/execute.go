package session

import (
	"context"
	"database/sql"

	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
)

func InTx(ctx context.Context) bool {
	executor := adaptor.ExecutorFrom(ctx)
	if _, ok := executor.(*sql.Tx); ok {
		return true
	}
	return false
}
