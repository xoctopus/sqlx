package types

import (
	"bytes"
	"database/sql/driver"
	"reflect"

	"github.com/go-json-experiment/json"
	jsonv1 "github.com/go-json-experiment/json/v1"
	"github.com/pkg/errors"
	"github.com/xoctopus/x/misc/must"
)

var options = []json.Options{
	jsonv1.OmitEmptyWithLegacySemantics(true),
}

// ScanJSONValue scan database input to value
func ScanJSONValue(src, dst any) error {
	switch x := src.(type) {
	case []byte:
		if len(x) == 0 {
			x = []byte("null")
		}
		return json.Unmarshal(x, dst)
	case string:
		if len(x) == 0 {
			x = "null"
		}
		return json.Unmarshal([]byte(x), dst)
	case nil:
		return nil
	default:
		return errors.Errorf("cannot sql.Scan from `%T` to `%T", src, dst)
	}
}

// DriverJSONValue unmarshal input to driver.Value
func DriverJSONValue(v any) (driver.Value, error) {
	if x, ok := v.(interface{ IsZero() bool }); ok && x.IsZero() {
		return "", nil
	}

	data, err := json.Marshal(v, options...)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type JSONDBType struct {
	typ string
}

func (t JSONDBType) DBType() string {
	if t.typ != "" {
		return t.typ
	}
	return "TEXT"
}

func (t *JSONDBType) WithDBType(typ string) {
	t.typ = typ
}

func JSONArrayOf[T any](s []T) *JSONArray[T] {
	return &JSONArray[T]{v: s}
}

type JSONArray[T any] struct {
	JSONDBType
	v []T
}

func (v *JSONArray[T]) Set(x []T) {
	v.v = x
}

func (v *JSONArray[T]) Get() []T {
	return v.v
}

func (v JSONArray[T]) IsZero() bool {
	return len(v.v) == 0
}

func (v JSONArray[T]) Value() (driver.Value, error) {
	return DriverJSONValue(v)
}

func (v *JSONArray[T]) Scan(src any) error {
	return ScanJSONValue(src, v)
}

func (v *JSONArray[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		v.v = nil
		return nil
	}

	x := make([]T, 0)
	if err := json.Unmarshal(data, &x, options...); err != nil {
		return err
	}
	v.v = x
	return nil
}

func (v *JSONArray[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.v, options...)
}

func JSONObjectOf[T any](v *T) JSONObject[T] {
	k := reflect.TypeFor[T]().Kind()
	must.BeTrueF(
		k == reflect.Struct || k == reflect.Map,
		"it shouldn't be a slice but a struct",
	)
	return JSONObject[T]{v: v}
}

type JSONObject[T any] struct {
	JSONDBType
	v *T
}

func (v *JSONObject[T]) Set(x *T) {
	v.v = x
}

func (v *JSONObject[T]) Get() *T {
	return v.v
}

func (v JSONObject[T]) IsZero() bool {
	return v.v == nil
}

func (v JSONObject[T]) Value() (driver.Value, error) {
	return DriverJSONValue(v)
}

func (v *JSONObject[T]) Scan(src any) error {
	return ScanJSONValue(src, v)
}

func (v *JSONObject[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		v.v = nil
		return nil
	}

	x := new(T)
	if err := json.Unmarshal(data, x, options...); err != nil {
		return err
	}
	v.v = x
	return nil
}

func (v JSONObject[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.v, options...)
}

func (v *JSONObject[T]) OneOf() []any {
	return []any{new(T)}
}
