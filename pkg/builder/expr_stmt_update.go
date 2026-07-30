package builder

import (
	"context"
	"iter"
	"slices"

	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/pkg/frag"
)

// ErrUpdateNeedLimitation is reserved for update-without-limit errors.
var ErrUpdateNeedLimitation = any(nil)

// Update starts an UPDATE statement on t.
func Update(t Table, modifiers ...string) *StmtUpdate {
	return &StmtUpdate{table: t, modifiers: modifiers}
}

// UpdateIgnore starts an UPDATE IGNORE statement on t.
func UpdateIgnore(t Table) *StmtUpdate {
	return Update(t, "IGNORE")
}

// StmtUpdate builds an UPDATE statement.
type StmtUpdate struct {
	table       Table
	from        Table
	modifiers   []string
	assignments iter.Seq[Assignment]
	additions   Additions
}

func (s StmtUpdate) Set(assignments ...Assignment) *StmtUpdate {
	if len(assignments) != 0 {
		s.assignments = slices.Values(assignments)
	}
	return &s
}

func (s StmtUpdate) From(from Table, additions ...Addition) *StmtUpdate {
	s.from = from
	s.additions = append(s.additions, additions...)
	return &s
}

func (s StmtUpdate) Where(cond frag.Fragment, additions ...Addition) *StmtUpdate {
	if cond != nil {
		s.additions = []Addition{Where(cond)}
	}
	s.additions = append(s.additions, additions...)
	return &s
}

func (s *StmtUpdate) IsNil() bool {
	return s == nil || frag.IsNil(s.table) || s.assignments == nil
}

func (s *StmtUpdate) Frag(ctx context.Context) frag.Iter {
	hasFrom, hasJoin := false, false
	if s.from != nil {
		ctx = WithToggles(ctx, TOGGLE__MULTI_TABLE)
		hasFrom = true
	}
	for _, a := range s.additions {
		if a.Type() == addition_JOIN {
			ctx = WithToggles(ctx, TOGGLE__MULTI_TABLE)
			hasJoin = true
		}
	}

	must.BeTrueF(
		!(hasFrom && hasJoin),
		"",
	)

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

		if !yield("UPDATE", nil) {
			return
		}

		for i := range s.modifiers {
			if !yield(" "+s.modifiers[i], nil) {
				return
			}
		}

		if !yield(" ", nil) {
			return
		}

		for q, args := range s.table.Frag(ctx) {
			if !yield(q, args) {
				return
			}
		}

		joins := ExtractAdditions(addition_JOIN, s.additions...)
		if !frag.IsNil(joins) {
			for q, args := range joins.Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}
		}

		if assignments := s.assignments; assignments != nil {
			if !yield(" SET ", nil) {
				return
			}
			for q, args := range frag.ComposeSeq(", ", frag.NonNil(assignments)).Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}
		}

		if s.from != nil {
			if !yield(" FROM ", nil) {
				return
			}
			for q, args := range s.from.Frag(ctx) {
				if !yield(q, args) {
					return
				}
			}
		}

		for q, args := range s.additions.Frag(WithToggles(ctx, TOGGLE__SKIP_COMMENTS, TOGGLE__SKIP_JOIN)) {
			if !yield(q, args) {
				return
			}
		}
	}
}
