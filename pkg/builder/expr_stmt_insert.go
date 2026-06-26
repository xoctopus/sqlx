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
				if !yield(q, args) {
					return
				}
			}
			if !yield("\n", nil) {
				return
			}
		}

		if !yield("INSERT", nil) {
			return
		}

		for i := range s.modifiers {
			if !yield(" "+s.modifiers[i], nil) {
				return
			}
		}

		if !yield(" INTO ", nil) {
			return
		}

		for q, args := range s.table.Frag(ctx) {
			if !yield(q, args) {
				return
			}
		}

		if !yield(" ", nil) {
			return
		}

		// must use tuple mode
		for q, args := range s.assignments.Frag(WithToggles(ctx, TOGGLE__TUPLE_ASSIGNMENTS)) {
			if !yield(q, args) {
				return
			}
		}

		for q, args := range s.additions.Frag(WithToggles(ctx, TOGGLE__SKIP_COMMENTS)) {
			if !yield(q, args) {
				return
			}
		}
	}
}
