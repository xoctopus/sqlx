package types

import (
	"database/sql"
	"database/sql/driver"

	"github.com/xoctopus/sqlx/internal"
)

// DBValue can convert between RDB and Go values and describe the RDB datatype.
type DBValue interface {
	driver.Valuer
	sql.Scanner
	DBType(driver string) string
}

// DBTypeAdapter allows overriding the SQL datatype string.
type DBTypeAdapter interface {
	WithDBType(driver string)
}

// CreationMarker marks a record as created.
type CreationMarker interface {
	MarkCreatedAt()
}

// ModificationMarker marks a record as modified.
type ModificationMarker interface {
	MarkModifiedAt()
}

// DeletionMarker marks a record as deleted.
type DeletionMarker interface {
	MarkDeletedAt()
}

// SoftDeletion describes soft-deletion field metadata.
type SoftDeletion interface {
	// SoftDeletion returns soft deletion field name, modification fields if any,
	// and the default value of the deletion field.
	SoftDeletion() (deletion string, modifications []string, v driver.Value)
}

// HasSoftDeletion reports whether M implements SoftDeletion.
func HasSoftDeletion[M internal.Model]() bool {
	_, ok := any(new(M)).(SoftDeletion)
	return ok
}
