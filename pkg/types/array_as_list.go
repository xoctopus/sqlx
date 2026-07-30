package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/xoctopus/x/misc/must"
)

type TElement interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string
}

func parseElement[T TElement](s string) (T, error) {
	v := new(T)
	switch t := reflect.TypeFor[T](); t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, t.Bits())
		if err != nil {
			return *v, err
		}
		reflect.ValueOf(v).Elem().SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, t.Bits())
		if err != nil {
			return *v, err
		}
		reflect.ValueOf(v).Elem().SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(s, t.Bits())
		if err != nil {
			return *v, err
		}
		reflect.ValueOf(v).Elem().SetFloat(i)
	default:
		must.BeTrueF(t.Kind() == reflect.String, "unknown element type: %s", t)
		reflect.ValueOf(v).Elem().SetString(s)
	}
	return *v, nil
}

func ParseArrayAsList[T TElement](s string) (ArrayAsList[T], error) {
	list := make(ArrayAsList[T], 0)
	for part := range strings.SplitSeq(s, ",") {
		if part = strings.TrimSpace(part); len(part) == 0 {
			return nil, fmt.Errorf("empty content is not accepted")
		}
		v, err := parseElement[T](part)
		if err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}

type ArrayAsList[T TElement] []T

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
	switch src := v.(type) {
	case nil:
		*aa = ArrayAsList[T]{}
		return nil
	case []byte:
		return aa.UnmarshalText(src)
	case string:
		return aa.UnmarshalText([]byte(src))
	default:
		return fmt.Errorf("cannot scan type %T into ArrayAsList[%s]", v, reflect.TypeFor[T]())
	}
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

func (aa *ArrayAsList[T]) UnmarshalText(data []byte) error {
	x, err := ParseArrayAsList[T](string(data))
	if err != nil {
		return err
	}
	*aa = x
	return nil
}

func (aa *ArrayAsList[T]) MarshalText() ([]byte, error) {
	return []byte(aa.String()), nil
}

func (aa *ArrayAsList[T]) UnmarshalJSON(data []byte) error {
	s := string(data)
	if len(s) == 0 || s == "null" {
		*aa = make(ArrayAsList[T], 0)
		return nil
	}

	s, err := strconv.Unquote(s)
	if err != nil {
		return err
	}
	l, err := ParseArrayAsList[T](s)
	if err != nil {
		return err
	}
	*aa = l
	return nil
}

func (aa *ArrayAsList[T]) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(aa.String())), nil
}
