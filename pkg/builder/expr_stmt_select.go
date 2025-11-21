package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type SelectStatement interface {
	frag.Fragment

	asSelection()
}

func Select(f frag.Fragment, modifiers ...frag.Fragment) *StmtSelect {
	return &StmtSelect{
		projects:  f,
		modifiers: modifiers,
	}
}

type StmtSelect struct {
	SelectStatement

	table     frag.Fragment
	modifiers []frag.Fragment
	projects  frag.Fragment
	additions Additions
}

func (s *StmtSelect) IsNil() bool {
	return s == nil
}

func (s StmtSelect) From(t frag.Fragment, additions ...Addition) *StmtSelect {
	s.table = t
	s.additions = append(s.additions, additions...)
	return &s
}

func (s *StmtSelect) Frag(ctx context.Context) frag.Iter {
	for i := range s.additions {
		a := s.additions[i]
		if frag.IsNil(a) {
			continue
		}
		if a.Type() == addition_JOIN {
			ctx = WithToggles(ctx, TOGGLE__MULTI_TABLE)
		}
	}

	return func(yield func(string, []any) bool) {
		comments := ExtractAdditions(addition_COMMENT, s.additions...)
		if !frag.IsNil(comments) {
			for q, args := range comments.Frag(ctx) {
				yield(q, args)
			}
			yield("\n", nil)
		}

		yield("SELECT", nil)

		for _, m := range s.modifiers {
			for q, args := range m.Frag(ctx) {
				yield(" "+q, args)
			}
		}
		yield(" ", nil)

		projects := s.projects
		if frag.IsNil(s.projects) {
			projects = frag.Lit("*")
		}

		for q, args := range projects.Frag(WithToggles(ctx, TOGGLE__IN_PROJECT)) {
			yield(q, args)
		}

		if !frag.IsNil(s.table) {
			yield(" FROM ", nil)

			for q, args := range s.table.Frag(ctx) {
				yield(q, args)
			}
		}

		for q, args := range s.additions.Frag(WithToggles(ctx, TOGGLE__SKIP_COMMENTS)) {
			yield(q, args)
		}
	}
}
