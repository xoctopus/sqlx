package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func Insert(modifiers ...string) *StmtInsert {
	return &StmtInsert{modifiers: modifiers}
}

type StmtInsert struct {
	table       Table
	modifiers   []string
	assignments Assignments
	additions   Additions
}

func (s StmtInsert) Into(t Table, additions ...Addition) *StmtInsert {
	s.table = t
	s.additions = append(s.additions, additions...)
	return &s
}

func (s StmtInsert) Values(cols Cols, values ...any) *StmtInsert {
	s.assignments = Assignments{ColumnsAndValues(cols, values...)}
	return &s
}

func (s *StmtInsert) IsNil() bool {
	return s == nil || s.table == nil || len(s.assignments) == 0
}

func (s *StmtInsert) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		comments := ExtractAdditions(addition_COMMENT, s.additions...)
		if !frag.IsNil(comments) {
			for q, args := range comments.Frag(ctx) {
				yield(q, args)
			}
			yield("\n", nil)
		}

		yield("INSERT", nil)

		for i := range s.modifiers {
			yield(" "+s.modifiers[i], nil)
		}

		yield(" INTO ", nil)

		for q, args := range s.table.Frag(ctx) {
			yield(q, args)
		}

		yield(" ", nil)

		for q, args := range s.assignments.Frag(WithToggles(ctx, TOGGLE__ASSIGNMENTS)) {
			yield(q, args)
		}

		for q, args := range s.additions.Frag(WithToggles(ctx, TOGGLE__SKIP_COMMENTS)) {
			yield(q, args)
		}
	}
}
