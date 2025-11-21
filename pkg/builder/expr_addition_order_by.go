package builder

import (
	"context"
	"slices"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type OrderAddition interface {
	Addition
	frag.Fragment

	Mode() string
}

func OrderBy(os ...OrderAddition) Addition {
	final := make([]OrderAddition, 0, len(os))
	for i := range os {
		o := os[i]
		if !frag.IsNil(o) {
			final = append(final, o)
		}
	}
	return &orders{orders: final}
}

func Order(by frag.Fragment, ex ...frag.Fragment) OrderAddition {
	return &order{by: by, ex: ex}
}

func AscOrder(by frag.Fragment, ex ...frag.Fragment) OrderAddition {
	return &order{by: by, mode: "ASC", ex: ex}
}

func DescOrder(by frag.Fragment, ex ...frag.Fragment) OrderAddition {
	return &order{by: by, mode: "DESC", ex: ex}
}

type order struct {
	by   frag.Fragment
	mode string
	ex   []frag.Fragment
}

func (o *order) Type() AdditionType {
	return addition_ORDER_BY
}

func (o *order) Mode() string {
	return o.mode
}

func (o *order) IsNil() bool {
	return o == nil || frag.IsNil(o.by)
}

func (o *order) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		for q, args := range frag.Block(o.by).Frag(ctx) {
			yield(q, args)
		}
		if o.mode != "" {
			yield(" "+o.mode, nil)
		}
		for _, x := range o.ex {
			if !frag.IsNil(x) {
				for q, args := range x.Frag(ctx) {
					yield(q, args)
				}
			}
		}
	}
}

type orders struct {
	orders []OrderAddition
}

func (o *orders) Type() AdditionType {
	return addition_ORDER_BY
}

func (o *orders) IsNil() bool {
	return o == nil || len(o.orders) == 0
}

func (o *orders) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		yield("ORDER BY ", nil)

		iter := frag.ComposeSeq(",", frag.NonNil(slices.Values(o.orders))).Frag(ctx)
		for q, args := range iter {
			yield(q, args)
		}
	}
}

// NullsFirst mysql/sqlite unsupported
// Alternative: ORDER BY `col` IS NULL
func NullsFirst() frag.Fragment {
	return frag.Lit(" NULLS FIRST")
}

// NullsLast mysql/sqlite unsupported
// Alternative: ORDER BY `col` IS NOT NULL
func NullsLast() frag.Fragment {
	return frag.Lit(" NULLS LAST")
}

// DistinctOn mysql/sqlite unsupported
// Alternative:
/*
SELECT `_id`, `score` FROM (
	SELECT *, ROW_NUMBER() OVER (PARTITION BY `_id` ORDER BY `score` ASE) AS `row_no`
	FROM `table`
) AS grouped
WHERE `row_no` = 1
*/
func DistinctOn(on ...frag.Fragment) frag.Fragment {
	return &distinct{on: on}
}

type distinct struct {
	on []frag.Fragment
}

func (d *distinct) IsNil() bool {
	return d == nil || len(d.on) == 0
}

func (d *distinct) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		yield("DISTINCT ON (", nil)

		iter := frag.ComposeSeq(",", frag.NonNil(slices.Values(d.on))).Frag(ctx)
		for q, args := range iter {
			yield(q, args)
		}

		yield(")", nil)
	}
}
