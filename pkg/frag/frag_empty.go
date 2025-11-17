package frag

import (
	"context"
)

func Empty() Fragment { return &empty{} }

type empty struct{}

func (empty) IsNil() bool {
	return true
}

func (empty) Frag(_ context.Context) Iter {
	return func(yield func(string, []any) bool) {
		return
	}
}
