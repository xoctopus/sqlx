package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type GroupByAddition interface {
	Addition

	Having(cond frag.Fragment) GroupByAddition
}

func GroupBy(groups ...frag.Fragment) GroupByAddition {
	return &groupby{groups: groups}
}

type groupby struct {
	groups []frag.Fragment
	having frag.Fragment
}

func (g *groupby) Type() AdditionType {
	return addition_GROUP_BY
}

func (g *groupby) Having(cond frag.Fragment) GroupByAddition {
	g.having = cond
	return g
}

func (g *groupby) IsNil() bool {
	return g == nil || len(g.groups) == 0
}

func (g *groupby) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		if !yield("GROUP BY ", nil) {
			return
		}

		for i, group := range g.groups {
			if i > 0 {
				if !yield(",", nil) {
					return
				}
			}
			for q, args := range group.Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}
		}
		if !(frag.IsNil(g.having)) {
			if !yield(" HAVING ", nil) {
				return
			}
			for q, args := range g.having.Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}
		}
	}
}
