package types

import (
	"fmt"
)

// Boolean converts a Go bool to Bool (TRUE/FALSE).
func Boolean(b bool) Bool {
	if b {
		return TRUE
	}
	return FALSE
}

// Bool a ternary boolean state UNKNOWN, TRUE or FALSE
type Bool int

const (
	_     Bool = iota
	TRUE       // true
	FALSE      // false
)

func (v Bool) Bool() bool {
	return v == TRUE
}

func (Bool) OpenAPISchemaType() []string {
	return []string{"boolean"}
}

func (v Bool) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

func (v *Bool) UnmarshalJSON(data []byte) (err error) {
	return v.UnmarshalText(data)
}

func (v Bool) MarshalText() ([]byte, error) {
	switch v {
	case FALSE:
		return []byte("false"), nil
	case TRUE:
		return []byte("true"), nil
	default:
		return []byte("null"), nil
	}
}

func (v *Bool) UnmarshalText(data []byte) error {
	switch string(data) {
	case "false", `"false"`:
		*v = FALSE
		return nil
	case "true", `"true"`:
		*v = TRUE
		return nil
	case "null":
		*v = 0
		return nil
	default:
		return fmt.Errorf("invalid boolean value: %q", string(data))
	}
}
