package sqltypes

import (
	"github.com/pkg/errors"
)

type Bool int

const (
	_     Bool = iota
	TRUE       // true
	FALSE      // false
)

func Boolean(b bool) Bool {
	if b {
		return TRUE
	}
	return FALSE
}

func (v Bool) Bool() bool {
	return v == TRUE
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
		return errors.Errorf("invalid boolean value: %q", string(data))
	}
}
