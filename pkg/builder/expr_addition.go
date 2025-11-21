package builder

import (
	"cmp"
	"context"
	"slices"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type Addition interface {
	frag.Fragment

	Type() AdditionType
}

type AdditionType int8

const (
	addition_JOIN AdditionType = iota + 1
	addition_WHERE
	addition_GROUP_BY
	addition_ORDER_BY
	addition_ON_CONFLICT
	addition_LIMIT
	addition_RETURNING
	addition_FOR_UPDATE
	addition_COMBINATION
	addition_COMMENT = 127
)

type Additions []Addition

func (as Additions) IsNil() bool {
	return len(as) == 0
}

func (as Additions) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		additions := slices.SortedFunc(slices.Values(as), func(a, b Addition) int {
			return cmp.Compare(a.Type(), b.Type())
		})

		for _, a := range additions {
			if HasToggle(ctx, TOGGLE__SKIP_COMMENTS) && a.Type() == addition_COMMENT ||
				HasToggle(ctx, TOGGLE__SKIP_JOIN) && a.Type() == addition_JOIN {
				continue
			}

			if frag.IsNil(a) {
				continue
			}

			if a.Type() != addition_COMMENT {
				yield(" ", nil)
			}

			for q, args := range a.Frag(ctx) {
				yield(q, args)
			}
		}
	}
}

func AsAddition(t AdditionType, f frag.Fragment) Addition {
	return &addition{Fragment: f, T: t}
}

func ExtractAdditions(t AdditionType, additions ...Addition) (filtered Additions) {
	for _, a := range additions {
		if frag.IsNil(a) {
			continue
		}
		if a.Type() == t {
			filtered = append(filtered, a)
		}
	}
	return
}

type addition struct {
	frag.Fragment

	T AdditionType
}

func (a *addition) Type() AdditionType {
	return a.T
}

func (a *addition) IsNil() bool {
	return a == nil || frag.IsNil(a.Fragment)
}
