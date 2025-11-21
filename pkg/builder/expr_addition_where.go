package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func Where(f frag.Fragment) Addition {
	switch x := f.(type) {
	case *where:
		return x
	default:
		return &where{condition: AsCond(x)}
	}
}

type where struct {
	condition SqlCondition
}

func (w *where) Type() AdditionType {
	return addition_WHERE
}

func (w *where) IsNil() bool {
	return w == nil || frag.IsNil(w.condition)
}

func (w *where) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		yield("WHERE ", nil)
		for q, args := range w.condition.Frag(ctx) {
			yield(q, args)
		}
	}
}
