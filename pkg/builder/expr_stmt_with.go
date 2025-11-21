package builder

import (
	"context"
	"strings"

	"github.com/xoctopus/x/iterx"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type SubQuery func(Table) frag.Fragment

func WithRecursive(t Table, q SubQuery) *WithStmt {
	return With(t, q, "RECURSIVE")
}

func With(t Table, q SubQuery, modifiers ...string) *WithStmt {
	return (&WithStmt{modifiers: modifiers}).With(t, q)
}

type WithStmt struct {
	modifiers []string
	tables    []Table
	asList    []SubQuery
	stmt      func(...Table) frag.Fragment
}

func (w WithStmt) With(t Table, q SubQuery) *WithStmt {
	w.tables = append(w.tables, t)
	w.asList = append(w.asList, q)
	return &w
}

func (w WithStmt) Exec(stmt func(...Table) frag.Fragment) *WithStmt {
	w.stmt = stmt
	return &w
}

func (w *WithStmt) IsNil() bool {
	return w == nil || len(w.tables) == 0 || len(w.asList) == 0 || w.stmt == nil
}

func (w *WithStmt) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		yield("WITH", nil)

		if len(w.modifiers) > 0 {
			yield(" "+strings.Join(w.modifiers, " "), nil)
		}

		for i, t := range w.tables {
			if i > 0 {
				yield(",", nil)
			}
			yield(" ", nil)

			for q, args := range t.Frag(ctx) {
				yield(q, args)
			}

			iter := frag.Block(
				frag.ComposeSeq(
					",",
					iterx.Map(t.Cols(), func(c Col) frag.Fragment { return c }),
				),
			).Frag(ctx)
			for q, args := range iter {
				yield(q, args)
			}

			yield(" AS ", nil)

			iter = frag.Block(w.asList[i](t)).Frag(ctx)
			for q, args := range iter {
				yield(q, args)
			}
		}

		for q, args := range w.stmt(w.tables...).Frag(ctx) {
			yield(q, args)
		}
	}
}
