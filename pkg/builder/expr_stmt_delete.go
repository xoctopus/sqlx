package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

// Delete starts a DELETE statement builder.
func Delete() *StmtDelete {
	return &StmtDelete{}
}

// StmtDelete builds a DELETE statement.
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
				if !yield(q, args) {
					return
				}
			}
			if !yield("\n", nil) {
				return
			}
		}
		if !yield("DELETE FROM ", nil) {
			return
		}

		for q, args := range s.table.Frag(ctx) {
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
