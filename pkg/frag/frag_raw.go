package frag

import (
	"context"
)

type Raw string

func (r Raw) IsNil() bool { return len(r) == 0 }

func (r Raw) Frag(_ context.Context) Iter {
	return func(yield func(string, []any) bool) {
		yield(string(r), nil)
	}
}
