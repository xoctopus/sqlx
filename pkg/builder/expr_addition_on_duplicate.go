package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func OnDuplicate(assignments ...Assignment) Addition {
	return &onduplicate{
		assignments: assignments,
	}
}

type onduplicate struct {
	assignments Assignments
}

func (o onduplicate) Type() AdditionType {
	return addition_ON_DUPLICATE
}

func (o *onduplicate) IsNil() bool {
	return o == nil || frag.IsNil(o.assignments)
}

func (o *onduplicate) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		yield("ON DUPLICATE KEY UPDATE ", nil)
		for q, args := range o.assignments.Frag(ctx) {
			yield(q, args)
		}
	}
}
