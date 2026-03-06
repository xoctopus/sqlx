package builder

import (
	"context"
	"slices"

	"github.com/xoctopus/x/iterx"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type OnConflictAddition interface {
	Addition

	DoNothing() OnConflictAddition
	DoUpdateSet(...Assignment) OnConflictAddition
}

func OnConflict(cols ColIter) OnConflictAddition {
	return &onconflict{cols: cols}
}

type onconflict struct {
	cols        ColIter
	nothing     bool
	assignments []Assignment
}

func (o onconflict) Type() AdditionType {
	return addition_ON_CONFLICT
}

func (o onconflict) DoNothing() OnConflictAddition {
	o.nothing = true
	return &o
}

func (o onconflict) DoUpdateSet(assignments ...Assignment) OnConflictAddition {
	// o.assigns = iterx.Map(
	// 	iterx.FilterSlice(assignments, func(a Assignment) bool {
	// 		if !frag.IsNil(a) {
	// 			o.assignments = append(o.assignments, a)
	// 			return true
	// 		}
	// 		return false
	// 	}),
	// 	func(a Assignment) frag.Fragment {
	// 		return a
	// 	},
	// )
	o.nothing = false
	o.assignments = append(o.assignments, assignments...)
	return &o
}

func (o *onconflict) IsNil() bool {
	return o == nil || o.cols == nil || (!o.nothing && len(o.assignments) == 0)
}

func (o *onconflict) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		if !yield("ON CONFLICT ", nil) {
			return
		}

		for q, args := range frag.Block(
			frag.ComposeSeq(
				",",
				iterx.Map(
					o.cols.Cols(),
					func(c Col) frag.Fragment { return c },
				),
			),
		).Frag(ctx) {
			if !yield(q, args) {
				return
			}
		}

		if !yield(" DO ", nil) {
			return
		}
		if o.nothing {
			if !yield("NOTHING", nil) {
				return
			}
			return
		}

		if !yield("UPDATE SET ", nil) {
			return
		}

		frags := iterx.Map(
			slices.Values(o.assignments),
			func(a Assignment) frag.Fragment {
				if frag.IsNil(a) {
					return nil
				}
				return a
			},
		)
		for q, args := range frag.ComposeSeq(", ", frag.NonNil(frags)).Frag(ctx) {
			if !yield(q, args) {
				return
			}
		}
	}
}
