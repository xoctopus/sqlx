package modeled

import (
	"github.com/xoctopus/sqlx/internal"
	"github.com/xoctopus/sqlx/pkg/builder"
)

// Newer embeds Model() for scoped table / column / key wrappers.
type Newer[M internal.Model] struct{}

// Model returns a new zero value of M.
func (m *Newer[M]) Model() *M {
	return new(M)
}

type (
	// Model is a table-backed struct type.
	Model = internal.Model
	// ModelNewer constructs a model instance of type M.
	ModelNewer[M Model] internal.ModelNewer[M]
)

// OrderAddition is a model-scoped order clause (same as [builder.OrderAddition]).
type OrderAddition[M Model] interface {
	builder.OrderAddition
}
