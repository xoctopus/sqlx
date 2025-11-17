package frag

import (
	"context"
)

func Single(arg any) Fragment {
	return &single{arg: arg}
}

type single struct {
	arg any
}

func (f *single) IsNil() bool { return false }

func (f *single) Frag(_ context.Context) Iter {
	return func(yield func(string, []any) bool) {
		yield("?", []any{f.arg})
	}
}
