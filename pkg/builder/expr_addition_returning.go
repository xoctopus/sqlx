package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type ReturningAddition interface {
	Addition
}

func Returning(p frag.Fragment) ReturningAddition {
	return &returning{p: p}
}

type returning struct {
	p frag.Fragment
}

func (r *returning) Type() AdditionType {
	return addition_RETURNING
}

func (r *returning) IsNil() bool {
	return r == nil || frag.IsNil(r.p)
}

func (r *returning) Frag(ctx context.Context) frag.Iter {
	return frag.Query("RETURNING ?", r.p).Frag(WithToggles(ctx, TOGGLE__IN_PROJECT))
}
