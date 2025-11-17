package frag

import (
	"context"
	"iter"
)

// Block returns a bracketed Fragment
// eg:
//
//	select ... from ... => (select ... from ...)
func Block(f Fragment) Fragment {
	if IsNil(f) {
		return Empty()
	}
	return &block{
		sep:       ",",
		seq:       func(yield func(Fragment) bool) { yield(f) },
		bracketed: true,
	}
}

func BlockWithoutBrackets(f Fragment) Fragment {
	b := Block(f)
	if x, ok := b.(*block); ok {
		x.bracketed = false
	}
	return b
}

type block struct {
	f         Fragment
	seq       iter.Seq[Fragment]
	sep       string
	bracketed bool
}

func (f *block) IsNil() bool { return f.seq == nil }

func (f *block) Frag(ctx context.Context) Iter {
	return func(yield func(string, []any) bool) {
		i := 0
		for frag := range NonNil(f.seq) {
			if f.bracketed && i == 0 {
				if !yield("(", nil) {
					return
				}
			}
			if i > 0 {
				if !yield(f.sep, nil) {
					return
				}
			}
			for query, args := range frag.Frag(ctx) {
				if !yield(query, args) {
					return
				}
				i++
			}
			if f.bracketed && i > 0 {
				if !yield(")", nil) {
					return
				}
			}
		}
	}
}
