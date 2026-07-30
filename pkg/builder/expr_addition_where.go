package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

// Where builds a WHERE addition.
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
		if !yield("WHERE ", nil) {
			return
		}
		for q, args := range w.condition.Frag(ctx) {
			if !yield(q, args) {
				return
			}
		}
	}
}
