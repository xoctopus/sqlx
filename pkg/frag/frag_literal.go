package frag

import "context"

func Lit(q string) Fragment {
	return literal(q)
}

func Empty() Fragment {
	return literal("")
}

type literal string

func (l literal) IsNil() bool {
	return len(l) == 0
}

func (l literal) Frag(ctx context.Context) Iter {
	return func(yield func(string, []any) bool) {
		yield(string(l), nil)
	}
}
