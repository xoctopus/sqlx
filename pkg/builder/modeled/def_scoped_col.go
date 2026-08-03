package modeled

import (
	"iter"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

// Col is a [builder.Col] bound to model type M.
type Col[M Model] interface {
	ModelNewer[M]
	builder.Col

	// ComputedBy returns a copy with a computed expression.
	ComputedBy(frag.Fragment) Col[M]
}

// CastC wraps a [builder.Col] as a model-scoped [Col].
func CastC[M Model](c builder.Col) Col[M] {
	return &col[M]{Col: c}
}

type col[M Model] struct {
	Newer[M]
	builder.Col
}

// Unwrap returns the underlying [builder.Col].
func (c *col[M]) Unwrap() builder.Col {
	return c.Col
}

// ComputedBy returns a copy with a computed expression.
func (c *col[M]) ComputedBy(f frag.Fragment) Col[M] {
	return CastC[M](builder.CC[any](c, builder.WithColComputed(f)))
}

// TCol is a typed [builder.TCol] bound to model type M and value type T.
type TCol[M Model, T any] interface {
	ModelNewer[M]
	builder.TCol[T]

	// ComputedBy returns an untyped model-scoped column with a computed expression.
	ComputedBy(frag.Fragment) Col[M]
	// TypedComputedBy returns a typed model-scoped column with a computed expression.
	TypedComputedBy(frag.Fragment) TCol[M, T]
}

// CT wraps a [builder.Col] as a typed model-scoped [TCol].
func CT[M Model, T any](c builder.Col) TCol[M, T] {
	return &tcol[M, T]{TCol: builder.CC[T](c)}
}

type tcol[M Model, T any] struct {
	Newer[M]
	builder.TCol[T]
}

// Unwrap returns the underlying [builder.Col].
func (c *tcol[M, T]) Unwrap() builder.Col {
	return c.TCol
}

// ComputedBy returns an untyped model-scoped column with a computed expression.
func (c *tcol[M, T]) ComputedBy(f frag.Fragment) Col[M] {
	return CastC[M](builder.CC[any](c, builder.WithColComputed(f)))
}

// TypedComputedBy returns a typed model-scoped column with a computed expression.
func (c *tcol[M, T]) TypedComputedBy(f frag.Fragment) TCol[M, T] {
	return CT[M, T](builder.CC[T](c, builder.WithColComputed(f)))
}

// ColIter iterates model-scoped columns.
type ColIter[M Model] interface {
	builder.ColIter
	MCols() iter.Seq[Col[M]]
}
