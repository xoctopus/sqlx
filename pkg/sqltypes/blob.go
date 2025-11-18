package sqltypes

import (
	"database/sql/driver"
	"strings"
)

type Blob []byte

func (v *Blob) Set(x []byte) {
	*v = Blob(x)
}

func (Blob) DBType(driver string) string {
	if strings.HasPrefix(driver, "postgres") {
		return "bytea"
	}
	return "blob"
}

func (v Blob) Value() (driver.Value, error) {
	if len(v) == 0 {
		return nil, nil
	}
	return []byte(v), nil
}

func (v *Blob) Scan(src any) error {
	switch x := src.(type) {
	case string:
		*v = Blob(x)
	case []byte:
		*v = Blob(x)
	}
	return nil
}
