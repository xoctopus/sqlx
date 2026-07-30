package frag

import (
	"context"
	"iter"
	"slices"
)

// Compose joins frags with sep, skipping nil fragments.
func Compose(sep string, frags ...Fragment) Fragment {
	return &compose{sep: sep, seq: NonNil(slices.Values(frags))}
}

// ComposeSeq joins an iterator of fragments with sep, skipping nil fragments.
func ComposeSeq(sep string, frags iter.Seq[Fragment]) Fragment {
	return &compose{sep: sep, seq: NonNil(frags)}
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
				if !yield(f.sep, nil) {
					return
				}
			}
			for q, args := range frag.Frag(ctx) {
				if !yield(q, args) {
					return
				}
				i++
			}
		}
	}
}
