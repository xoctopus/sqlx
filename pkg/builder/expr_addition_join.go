package builder

import (
	"context"
	"slices"
	"strings"

	"github.com/xoctopus/sqlx/pkg/frag"
)

// JoinAddition is a JOIN clause.
type JoinAddition interface {
	Addition

	On(condition frag.Fragment) JoinAddition
	Using(cols ...Col) JoinAddition
}

// Join builds a JOIN on t with optional join methods.
func Join(t frag.Fragment, methods ...string) JoinAddition {
	return &join{
		method: strings.Join(methods, " "),
		target: t,
	}
}

// InnerJoin builds an INNER JOIN on t.
func InnerJoin(t frag.Fragment) JoinAddition {
	return Join(t, "INNER")
}

// LeftJoin builds a LEFT JOIN on t.
func LeftJoin(t frag.Fragment) JoinAddition {
	return Join(t, "LEFT")
}

// RightJoin builds a RIGHT JOIN on t. Unsupported on SQLite.
func RightJoin(t frag.Fragment) JoinAddition {
	return Join(t, "RIGHT")
}

// FullJoin builds a FULL JOIN on t. Unsupported on MySQL/SQLite.
func FullJoin(t frag.Fragment) JoinAddition {
	return Join(t, "FULL")
}

// CrossJoin builds a CROSS JOIN on t.
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
		if !yield(method, nil) {
			return
		}

		for q, args := range j.target.Frag(ctx) {
			if !yield(q, args) {
				return
			}
		}

		if !frag.IsNil(j.cond) {
			if !yield(" ON ", nil) {
				return
			}
			for q, args := range j.cond.Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}
		}

		if len(j.cols) > 0 {
			if !yield(" USING (", nil) {
				return
			}
			cols := frag.ComposeSeq(",", frag.NonNil(slices.Values(j.cols)))
			for q, args := range cols.Frag(TrimToggles(ctx, TOGGLE__MULTI_TABLE)) {
				if !yield(q, args) {
					return
				}
			}
			if !yield(")", nil) {
				return
			}
		}
	}
}
