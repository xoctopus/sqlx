package types

import (
	"database/sql"
	"database/sql/driver"

	"github.com/xoctopus/sqlx/internal"
)

// DBValue can convert between rdb value and go value with description of rdb datatype
type DBValue interface {
	driver.Valuer
	sql.Scanner
	DBType(driver string) string
}

type DBTypeAdapter interface {
	WithDBType(driver string)
}

type CreationMarker interface {
	MarkCreatedAt()
}

type ModificationMarker interface {
	MarkModifiedAt()
}

type DeletionMarker interface {
	MarkDeletedAt()
}

type SoftDeletion interface {
	// SoftDeletion returns soft deletion field name, modifications fields if exists
	// and default value of deletion field
	SoftDeletion() (deletion string, modifications []string, v driver.Value)
}

func HasSoftDeletion[M internal.Model]() bool {
	_, ok := any(new(M)).(SoftDeletion)
	return ok
}
