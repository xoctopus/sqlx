package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ArrayAsListElement interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string
}

type ArrayAsList[T ArrayAsListElement] []T

var (
	_ driver.Valuer = (*ArrayAsList[int])(nil)
	_ sql.Scanner   = (*ArrayAsList[int])(nil)
)

func (aa ArrayAsList[T]) DBType(driver string) string {
	return "text"
}

func (aa ArrayAsList[T]) Value() (driver.Value, error) {
	return strings.Join(aa.Elements(), ","), nil
}

func (aa ArrayAsList[T]) String() string {
	return strings.Join(aa.Elements(), ",")
}

func (aa *ArrayAsList[T]) Scan(v any) error {
	var (
		x    ArrayAsList[T]
		err  error
		kind = reflect.TypeFor[T]().Kind()
	)

	switch src := v.(type) {
	case []byte:
		x = make(ArrayAsList[T], 0)
		err = x.AppendString(string(src))
	case string:
		x = make(ArrayAsList[T], 0)
		err = x.AppendString(src)
	default:
		err = fmt.Errorf("cannot scan type %T into ArrayAsList[%s]", v, kind)
	}
	if err == nil {
		*aa = x
	}
	return err
}

func (aa ArrayAsList[T]) Elements() []string {
	elements := make([]string, 0, len(aa))
	for _, e := range aa {
		elements = append(elements, fmt.Sprintf("%v", e))
	}
	return elements
}

func (aa *ArrayAsList[T]) Append(values ...T) {
	*aa = append(*aa, values...)
}

func (aa *ArrayAsList[T]) AppendString(values ...string) error {
	typed := make([]T, 0, len(values))
	for _, v := range values {
		for _, part := range strings.Split(v, ",") {
			tv := new(T)
			switch reflect.TypeFor[T]().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				u, err := strconv.ParseInt(part, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert %s to int: %w", part, err)
				}
				reflect.ValueOf(tv).Elem().SetInt(u)
				typed = append(typed, *tv)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				u, err := strconv.ParseUint(part, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert %s to uint: %w", part, err)
				}
				reflect.ValueOf(tv).Elem().SetUint(u)
				typed = append(typed, *tv)
			case reflect.String:
				reflect.ValueOf(tv).Elem().SetString(part)
				typed = append(typed, *tv)
			case reflect.Float32, reflect.Float64:
				u, err := strconv.ParseFloat(part, 64)
				if err != nil {
					return fmt.Errorf("failed to convert %s to float: %w", part, err)
				}
				reflect.ValueOf(tv).Elem().SetFloat(u)
				typed = append(typed, *tv)
			default:
				return fmt.Errorf("cannot append %s to ArrayAsList[%s]", v, reflect.TypeFor[T]().Kind())
			}
		}
	}
	aa.Append(typed...)
	return nil
}

func (aa *ArrayAsList[T]) UnmarshalJSON(data []byte) error {
	s, err := string(data), error(nil)
	if len(data) > 0 && data[0] == '"' {
		s, err = strconv.Unquote(s)
		if err != nil {
			return err
		}
	}
	x := make(ArrayAsList[T], 0)
	if err = x.AppendString(s); err != nil {
		return err
	}
	*aa = x
	return nil
}

func (aa *ArrayAsList[T]) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(aa.String())), nil
}
