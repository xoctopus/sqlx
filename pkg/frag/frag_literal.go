package frag

import "context"

// Lit creates a literal SQL snippet with no arguments.
func Lit(q string) Fragment {
	return literal(q)
}

// Empty is an empty fragment (IsNil reports true).
func Empty() Fragment {
	return literal("")
}

type literal string

func (l literal) IsNil() bool {
	return len(l) == 0
}

func (l literal) Frag(ctx context.Context) Iter {
	return func(yield func(string, []any) bool) {
		if !yield(string(l), nil) {
			return
		}
	}
}
