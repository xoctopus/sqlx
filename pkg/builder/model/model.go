package model

/*
import (
	"iter"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/internal"
)

type M[Model internal.Model] struct{}

func (M[Model]) Model() *Model {
	return new(Model)
}

type ModelNewer[Model internal.Model] internal.ModelNewer[Model]

type Table[Model internal.Model] interface {
	builder.Table
	ModelNewer[Model]
	ColIter[Model]
	KeyIter[Model]
}

type ColIter[Model internal.Model] interface {
	builder.ColIter
	Cols() iter.Seq[Col[Model]]
}

type KeyIter[Model internal.Model] interface {
	Keys() iter.Seq[Key[Model]]
}

type table[Model internal.Model] struct {
	M[Model]
	builder.Table
}

func (t *table[Model]) MKeys() iter.Seq {}
*/
