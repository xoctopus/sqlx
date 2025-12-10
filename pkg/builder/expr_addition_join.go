package builder

import (
	"context"
	"slices"
	"strings"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type JoinAddition interface {
	Addition

	On(condition frag.Fragment) JoinAddition
	Using(cols ...Col) JoinAddition
}

func Join(t frag.Fragment, methods ...string) JoinAddition {
	return &join{
		method: strings.Join(methods, " "),
		target: t,
	}
}

func InnerJoin(t frag.Fragment) JoinAddition {
	return Join(t, "INNER")
}

func LeftJoin(t frag.Fragment) JoinAddition {
	return Join(t, "LEFT")
}

// RightJoin sqlite unsupported
func RightJoin(t frag.Fragment) JoinAddition {
	return Join(t, "RIGHT")
}

// FullJoin mysql/sqlite unsupported
func FullJoin(t frag.Fragment) JoinAddition {
	return Join(t, "FULL")
}

func CrossJoin(t frag.Fragment) JoinAddition {
	return Join(t, "CROSS")
}

type join struct {
	method string
	target frag.Fragment
	cond   frag.Fragment
	cols   []Col
}

func (j *join) Type() AdditionType {
	return addition_JOIN
}

func (j *join) On(cond frag.Fragment) JoinAddition {
	j.cond = cond
	return j
}

func (j *join) Using(cols ...Col) JoinAddition {
	j.cols = cols
	return j
}

func (j *join) IsNil() bool {
	return j == nil ||
		frag.IsNil(j.target) ||
		(j.method != "CROSS" && len(j.cols) == 0 && frag.IsNil(j.cond))
}

func (j *join) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		method := "JOIN "
		if j.method != "" {
			method = j.method + " " + method
		}
		yield(method, nil)

		for q, args := range j.target.Frag(ctx) {
			yield(q, args)
		}

		if !frag.IsNil(j.cond) {
			yield(" ON ", nil)
			for q, args := range j.cond.Frag(ctx) {
				yield(q, args)
			}
		}

		if len(j.cols) > 0 {
			yield(" USING (", nil)
			cols := frag.ComposeSeq(",", frag.NonNil(slices.Values(j.cols)))
			for q, args := range cols.Frag(TrimToggles(ctx, TOGGLE__MULTI_TABLE)) {
				yield(q, args)
			}
			yield(")", nil)
		}
	}
}
