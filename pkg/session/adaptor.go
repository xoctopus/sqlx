package session

import (
	"context"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	_ "github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
)

type Adaptor = adaptor.Adaptor

func Open(ctx context.Context, endpoint string) (Adaptor, error) {
	return adaptor.Open(ctx, endpoint)
}
