package types

import (
	"database/sql/driver"
	"fmt"
)

type Blob []byte

func (v *Blob) Set(x []byte) {
	*v = Blob(x)
}

func (Blob) DBType(driver string) string {
	switch driver {
	case "postgres", "pg":
		return "bytea"
	default:
		return "blob"
	}
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
	case nil:
		*v = Blob(nil)
	default:
		return fmt.Errorf("cannot sql.Scan() %T to Blob", src)
	}
	return nil
}
