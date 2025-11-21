package types

import (
	"database/sql"
	"database/sql/driver"
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
