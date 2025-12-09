package def

import (
	"context"

	"github.com/xoctopus/x/contextx"
)

var CtxTagKey = contextx.NewT[string]()

func WithModelTagKey(ctx context.Context, key string) context.Context {
	return CtxTagKey.With(ctx, key)
}

func ModelTagKeyFrom(ctx context.Context) string {
	key, _ := CtxTagKey.From(ctx)
	if key != "" {
		return key
	}
	return "db"
}
