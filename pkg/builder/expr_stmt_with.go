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
		if !yield("WITH", nil) {
			return
		}

		if len(w.modifiers) > 0 {
			if !yield(" "+strings.Join(w.modifiers, " "), nil) {
				return
			}
		}

		for i, t := range w.tables {
			if i > 0 {
				if !yield(",", nil) {
					return
				}
			}
			if !yield(" ", nil) {
				return
			}

			for q, args := range t.Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}

			iter := frag.Block(
				frag.ComposeSeq(
					",",
					iterx.Map(t.Cols(), func(c Col) frag.Fragment { return c }),
				),
			).Frag(ctx)
			for q, args := range iter {
				if !yield(q, args) {
					return
				}
			}

			if !yield(" AS ", nil) {
				return
			}

			iter = frag.Block(w.asList[i](t)).Frag(ctx)
			for q, args := range iter {
				if !yield(q, args) {
					return
				}
			}
		}

		for q, args := range w.stmt(w.tables...).Frag(ctx) {
			if !yield(q, args) {
				return
			}
		}
	}
}
