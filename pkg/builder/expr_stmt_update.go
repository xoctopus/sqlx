package builder

import (
	"context"
	"iter"
	"slices"

	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/pkg/frag"
)

var ErrUpdateNeedLimitation = any(nil)

func Update(t Table, modifiers ...string) *StmtUpdate {
	return &StmtUpdate{table: t, modifiers: modifiers}
}

func UpdateIgnore(t Table) *StmtUpdate {
	return Update(t, "IGNORE")
}

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
				yield(q, args)
			}
			yield("\n", nil)
		}

		yield("UPDATE", nil)

		for i := range s.modifiers {
			yield(" "+s.modifiers[i], nil)
		}

		yield(" ", nil)

		for q, args := range s.table.Frag(ctx) {
			yield(q, args)
		}

		joins := ExtractAdditions(addition_JOIN, s.additions...)
		if !frag.IsNil(joins) {
			for q, args := range joins.Frag(ctx) {
				yield(q, args)
			}
		}

		if assignments := s.assignments; assignments != nil {
			yield(" SET ", nil)

			for q, args := range frag.ComposeSeq(", ", frag.NonNil(assignments)).Frag(ctx) {
				yield(q, args)
			}
		}

		if s.from != nil {
			yield(" FROM ", nil)
			for q, args := range s.from.Frag(ctx) {
				yield(q, args)
			}
		}

		for q, args := range s.additions.Frag(WithToggles(ctx, TOGGLE__SKIP_COMMENTS, TOGGLE__SKIP_JOIN)) {
			yield(q, args)
		}
	}
}
