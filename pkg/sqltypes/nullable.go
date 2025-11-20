package sqltypes

import (
	"database/sql/driver"
	"strings"
)

type Text string

func (v *Text) Set(str string) {
	*v = Text(str)
}

func (Text) DBType(driver string) string {
	return "text"
}

func (v Text) Value() (driver.Value, error) {
	if v == "" {
		return nil, nil
	}
	return string(v), nil
}

func (v *Text) Scan(src any) error {
	switch x := src.(type) {
	case string:
		*v = Text(x)
	case []byte:
		*v = Text(x)
	}
	return nil
}

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
