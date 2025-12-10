package modeled

import (
	"github.com/xoctopus/sqlx/internal"
	"github.com/xoctopus/sqlx/pkg/builder"
)

type Newer[M internal.Model] struct{}

func (m *Newer[M]) Model() *M {
	return new(M)
}

type (
	Model               = internal.Model
	ModelNewer[M Model] internal.ModelNewer[M]
)

type OrderAddition[M Model] interface {
	builder.OrderAddition
}
