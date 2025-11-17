package frag

import (
	"context"
	"iter"
	"slices"
)

func Compose(sep string, frags ...Fragment) Fragment {
	return &compose{sep: sep, seq: NonNil(slices.Values(frags))}
}

type compose struct {
	seq iter.Seq[Fragment]
	sep string
}

func (f *compose) IsNil() bool { return f.seq == nil }

func (f *compose) Frag(ctx context.Context) Iter {
	return func(yield func(string, []any) bool) {
		i := 0
		for frag := range NonNil(f.seq) {
			if i > 0 {
				yield(f.sep, nil)
			}
			for query, args := range frag.Frag(ctx) {
				yield(query, args)
				i++
			}
		}
	}
}
