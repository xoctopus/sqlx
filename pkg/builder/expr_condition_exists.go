package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func Exists(f frag.Fragment) SqlCondition {
	return AsCond(&exists{f: f})
}

type exists struct {
	f frag.Fragment
}

func (e *exists) IsNil() bool {
	return e.f.IsNil()
}

func (e *exists) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		yield("EXISTS ", nil)
		for q, args := range frag.Block(e.f).Frag(ctx) {
			yield(q, args)
		}
	}
}
