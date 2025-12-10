package modeled

import (
	"iter"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

type Col[M Model] interface {
	ModelNewer[M]
	builder.Col

	ComputedBy(frag.Fragment) Col[M]
}

func CastC[M Model](c builder.Col) Col[M] {
	return &col[M]{Col: c}
}

type col[M Model] struct {
	Newer[M]
	builder.Col
}

func (c *col[M]) Unwrap() builder.Col {
	return c.Col
}

func (c *col[M]) ComputedBy(f frag.Fragment) Col[M] {
	return CastC[M](builder.CC[any](c, builder.WithColComputed(f)))
}

type TCol[M Model, T any] interface {
	ModelNewer[M]
	builder.TCol[T]

	ComputedBy(frag.Fragment) Col[M]
	TypedComputedBy(frag.Fragment) TCol[M, T]
}

func CT[M Model, T any](c builder.Col) TCol[M, T] {
	return &tcol[M, T]{TCol: builder.CC[T](c)}
}

type tcol[M Model, T any] struct {
	Newer[M]
	builder.TCol[T]
}

func (c *tcol[M, T]) Unwrap() builder.Col {
	return c.TCol
}

func (c *tcol[M, T]) ComputedBy(f frag.Fragment) Col[M] {
	return CastC[M](builder.CC[any](c, builder.WithColComputed(f)))
}

func (c *tcol[M, T]) TypedComputedBy(f frag.Fragment) TCol[M, T] {
	return CT[M, T](builder.CC[T](c, builder.WithColComputed(f)))
}

type ColIter[M Model] interface {
	builder.ColIter
	MCols() iter.Seq[Col[M]]
}
