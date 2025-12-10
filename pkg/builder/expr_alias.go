package builder

import (
	"context"
	"slices"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func Alias(f frag.Fragment, name string) frag.Fragment {
	return &alias{
		name:     name,
		Fragment: f,
	}
}

type alias struct {
	frag.Fragment

	name string
}

func (a *alias) IsNil() bool {
	return a == nil || a.name == "" || frag.IsNil(a.Fragment)
}

func (a *alias) Frag(ctx context.Context) frag.Iter {
	return frag.Query(
		"? AS ?", a.Fragment, frag.Lit(a.name),
	).Frag(TrimToggles(ctx, TOGGLE__AUTO_ALIAS))
}

func AutoAlias(columns ...frag.Fragment) frag.Fragment {
	return &autoalias{
		columns: slices.Collect(frag.NonNil(slices.Values(columns))),
	}
}

type autoalias struct {
	columns []frag.Fragment
}

func (a *autoalias) IsNil() bool {
	return a == nil || len(a.columns) == 0
}

func (a *autoalias) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		for q, args := range frag.Compose(", ", a.columns...).Frag(WithToggles(ctx, TOGGLE__AUTO_ALIAS)) {
			yield(q, args)
		}
	}
}
