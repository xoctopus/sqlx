package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func Delete() *StmtDelete {
	return &StmtDelete{}
}

type StmtDelete struct {
	table     Table
	additions Additions
}

func (s *StmtDelete) IsNil() bool {
	return s == nil || frag.IsNil(s.table)
}

func (s StmtDelete) From(t Table, additions ...Addition) *StmtDelete {
	s.table = t
	s.additions = append(s.additions, additions...)
	return &s
}

func (s *StmtDelete) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		comments := ExtractAdditions(addition_COMMENT, s.additions...)
		if !frag.IsNil(comments) {
			for q, args := range comments.Frag(ctx) {
				yield(q, args)
			}
			yield("\n", nil)
		}
		yield("DELETE FROM ", nil)

		for q, args := range s.table.Frag(ctx) {
			yield(q, args)
		}

		for q, args := range s.additions.Frag(WithToggles(ctx, TOGGLE__SKIP_COMMENTS)) {
			yield(q, args)
		}
	}
}
