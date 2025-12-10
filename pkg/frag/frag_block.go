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

func BlockWithoutBrackets(seq iter.Seq[Fragment]) Fragment {
	if seq == nil {
		return nil
	}

	return &block{seq: seq, sep: ","}
}

type block struct {
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
				yield("(", nil)
			}
			if i > 0 {
				yield(f.sep, nil)
			}
			for query, args := range frag.Frag(ctx) {
				yield(query, args)
				i++
			}
			if f.bracketed && i > 0 {
				yield(")", nil)
			}
		}
	}
}
