package frag

import (
	"bytes"
	"context"
	"iter"
	"slices"
	"strings"
)

// Iter will yield a SQL fragment with a raw query(string with holder) and values
// eg:
//
//	query = INSERT INTO t_table (f_x,f_y,f_z) VALUES (?,?,?)
//	args  = 1, '2', 3.0
type Iter = iter.Seq2[string, []any]

// IsNil returns if a Fragment is nil
func IsNil(e Fragment) bool {
	return e == nil || e.IsNil()
}

// NonNil filter nil fragments
func NonNil[F Fragment](seq iter.Seq[F]) iter.Seq[Fragment] {
	return func(yield func(Fragment) bool) {
		for frag := range seq {
			if IsNil(frag) {
				continue
			}
			yield(frag)
		}
	}
}

type Func func(ctx context.Context) Iter

func (f Func) IsNil() bool {
	return f == nil
}

func (f Func) Frag(ctx context.Context) Iter {
	return f(ctx)
}

// Fragment defines an interface to present a sql fragment
type Fragment interface {
	IsNil() bool
	Frag(ctx context.Context) Iter
}

func Collect(ctx context.Context, f Fragment) (string, []any) {
	if IsNil(f) {
		return "", nil
	}

	var (
		query = bytes.NewBuffer(nil)
		args  = make([]any, 0)
	)

	for q, x := range f.Frag(ctx) {
		if len(q) > 0 {
			query.WriteString(strings.TrimPrefix(q, "\n"))
		}
		if len(args) > 0 {
			args = slices.Concat(args, x)
		}
	}
	return query.String(), args
}
