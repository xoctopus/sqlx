package builder

import (
	"context"
	"slices"
	"strings"

	"github.com/xoctopus/x/iterx"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type AssignmentMarker interface {
	asAssignment()
}

// ColumnsAndValues returns Assigment cols should be a Col or Cols
func ColumnsAndValues(cols frag.Fragment, values ...any) Assignment {
	count := 1
	if x, ok := cols.(interface{ Len() int }); ok {
		count = x.Len()
	}

	return &assignment{
		cols:   cols,
		count:  count,
		values: values,
	}
}

type Assignment interface {
	frag.Fragment

	AssignmentMarker
}

type assignment struct {
	cols   frag.Fragment
	count  int
	values []any
}

func (a *assignment) asAssignment() {}

func (a *assignment) IsNil() bool {
	return a == nil || frag.IsNil(a.cols) || len(a.values) == 0
}

func (a *assignment) Frag(ctx context.Context) frag.Iter {
	usev := HasToggle(ctx, TOGGLE__ASSIGNMENTS)

	return func(yield func(string, []any) bool) {
		if usev || len(a.values) > 1 {
			// (f_a,f_b...)
			for q, args := range frag.Block(a.cols).Frag(WithoutToggles(ctx, TOGGLE__MULTI_TABLE)) {
				yield(q, args)
			}
			values := a.values

			if len(values) == 1 {
				if stmt, ok := values[0].(SelectStatement); ok {
					yield(" ", nil)
					for q, args := range stmt.Frag(ctx) {
						yield(q, args)
					}
					return
				}
			}

			yield(" VALUES ", nil)
			frags := iterx.Map(
				slices.Chunk(values, a.count),
				func(values []any) frag.Fragment {
					return frag.Query(
						"("+strings.Repeat(",?", len(values))[1:]+")", // (?,?,...)
						values...,
					)
				},
			)
			for q, args := range frag.BlockWithoutBrackets(frags).Frag(ctx) {
				yield(q, args)
			}
			return
		}
		for q, args := range a.cols.Frag(WithoutToggles(ctx, TOGGLE__MULTI_TABLE)) {
			yield(q, args)
		}

		value := a.values[0]
		if stmt, ok := value.(SelectStatement); ok {
			value = frag.Block(stmt)
		}

		for q, args := range frag.Query(" = ?", value).Frag(ctx) {
			yield(q, args)
		}
	}
}

type Assignments []Assignment

func (as Assignments) asAssignment() {}

func (as Assignments) IsNil() bool {
	if len(as) == 0 {
		return true
	}
	for i := range as {
		a := as[i]
		if !frag.IsNil(a) {
			return false
		}
	}
	return true
}

func (as Assignments) Frag(ctx context.Context) frag.Iter {
	return frag.ComposeSeq(", ", frag.NonNil(slices.Values(as))).Frag(ctx)
}
