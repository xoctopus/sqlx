package frag

import (
	"bytes"
	"context"
	"fmt"
	"iter"
	"slices"
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
			if !yield(frag) {
				return
			}
		}
	}
}

// Fragment is a SQL fragment that yields query pieces and arguments.
type Fragment interface {
	// IsNil reports whether the fragment is empty.
	IsNil() bool
	// Frag yields SQL text chunks and their bound arguments.
	Frag(ctx context.Context) Iter
}

// Collect flattens f into a query string and argument list.
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
			// query.WriteString(strings.TrimPrefix(q, "\n"))
			query.WriteString(q)
		}
		if len(x) > 0 {
			args = slices.Concat(args, x)
		}
	}
	return query.String(), args
}

// Stringify returns a debug representation of f as "query | args".
func Stringify(f Fragment) string {
	q, args := Collect(context.Background(), f)
	return fmt.Sprintf("%s | %v", q, args)
}
