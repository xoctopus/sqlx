package builder

import (
	"context"
	"fmt"

	"github.com/xoctopus/sqlx/pkg/frag"
)

type LimitAddition interface {
	Addition

	Offset(offset int64) LimitAddition
}

func Limit(count int64) LimitAddition {
	return &limit{count: count}
}

type limit struct {
	count  int64
	offset int64
}

func (l *limit) Type() AdditionType {
	return addition_LIMIT
}

func (l *limit) Offset(offset int64) LimitAddition {
	l.offset = offset
	return l
}

func (l *limit) IsNil() bool {
	return l == nil || l.count <= 0
}

func (l *limit) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		yield(fmt.Sprintf("LIMIT %d", l.count), nil)
		if l.offset > 0 {
			yield(fmt.Sprintf(" OFFSET %d", l.offset), nil)
		}
	}
}
