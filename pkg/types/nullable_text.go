package types

import (
	"database/sql/driver"
	"fmt"
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
	case nil:
		*v = Text("")
	default:
		return fmt.Errorf("cannot sql.Scan() %T to Text", src)
	}
	return nil
}
