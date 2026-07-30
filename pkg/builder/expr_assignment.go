package builder

import (
	"context"
	"slices"
	"strings"

	"github.com/xoctopus/x/iterx"

	"github.com/xoctopus/sqlx/pkg/frag"
)

// AssignmentMarker marks assignment fragments.
type AssignmentMarker interface {
	asAssignment()
}

// ColumnsAndValues returns Assigment cols should be a Col or Cols
func ColumnsAndValues(cols frag.Fragment, values ...any) Assignment {
	count := 1
	if x, ok := cols.(interface{ Len() int }); ok {
		count = x.Len()
	}

	a := &assignment{
		cols:   cols,
		count:  count,
		values: values,
	}

	if count > 0 {
		a.flatten = make([]frag.Fragment, 0, count)
		if x, ok := cols.(ColIter); ok {
			for f := range x.Cols() {
				a.flatten = append(a.flatten, f)
			}
		}
		// must.BeTrueF(
		// 	len(a.flatten) == a.count && len(values)%a.count == 0,
		// 	"unmatched count of columns and values columns:%d values:%d",
		// 	len(a.flatten), len(values),
		// )
	}

	return a
}

// Assignment is a SET/VALUES assignment fragment.
type Assignment interface {
	frag.Fragment

	AssignmentMarker
}

type assignment struct {
	cols    frag.Fragment
	flatten []frag.Fragment
	count   int
	values  []any
}

func (a *assignment) asAssignment() {}

func (a *assignment) IsNil() bool {
	return a == nil || frag.IsNil(a.cols) || len(a.values) == 0
}

func (a *assignment) Frag(ctx context.Context) frag.Iter {
	toggled := HasToggle(ctx, TOGGLE__TUPLE_ASSIGNMENTS)

	return func(yield func(string, []any) bool) {
		// tuple mode (f_a,f_b...) VALUES
		if toggled {
			for q, args := range frag.Block(a.cols).Frag(TrimToggles(ctx, TOGGLE__MULTI_TABLE)) {
				if !yield(q, args) {
					return
				}
			}
			values := a.values

			if len(values) == 1 {
				if stmt, ok := values[0].(SelectStatement); ok {
					if !yield(" ", nil) {
						return
					}
					for q, args := range stmt.Frag(ctx) {
						if !yield(q, args) {
							return
						}
					}
					return
				}
			}

			if !yield("\nVALUES", nil) {
				return
			}
			frags := iterx.Map(
				slices.Chunk(values, a.count),
				func(values []any) frag.Fragment {
					return frag.Query(
						"\n("+strings.Repeat(",?", len(values))[1:]+")", // (?,?,...)
						values...,
					)
				},
			)
			for q, args := range frag.BlockWithoutBrackets(frags).Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}
			return
		}

		if a.count <= 1 {
			for q, args := range a.cols.Frag(TrimToggles(ctx, TOGGLE__MULTI_TABLE)) {
				if !yield(q, args) {
					return
				}
			}

			value := a.values[0]
			if stmt, ok := value.(SelectStatement); ok {
				value = frag.Block(stmt)
			}

			for q, args := range frag.Query(" = ?", value).Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}
			return
		}

		if len(a.values) == a.count {
			var flatten Assignments
			for i := 0; i < a.count; i++ {
				flatten = append(flatten, &assignment{
					cols:   a.flatten[i],
					count:  1,
					values: []any{a.values[i]},
				})
			}

			for q, args := range flatten.Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}
			return
		}

		// TODO invalid or unmatch
	}
}

// Assignments is a list of Assignment values.
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
