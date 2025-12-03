package frag

import (
	"context"
	"database/sql/driver"
	"iter"
	"reflect"
	"slices"
	"strings"

	"github.com/xoctopus/x/reflectx"
)

type Values[T any] iter.Seq[T]

func (vs Values[T]) IsNil() bool {
	return vs == nil
}

func (vs Values[T]) Frag(_ context.Context) Iter {
	return func(yield func(string, []any) bool) {
		i := 0
		for v := range vs {
			if i == 0 {
				yield("?", []any{v})
			} else {
				yield(",?", []any{v})
			}
			i++
		}
	}
}

// CustomValueArg supports replacing ? to some sql snippet
// eg: ? => UPPER(?)
type CustomValueArg interface {
	ValueEx() string
}

func ArgIter(ctx context.Context, v any) Iter {
	switch x := v.(type) {
	case CustomValueArg:
		return func(yield func(string, []any) bool) {
			yield(x.ValueEx(), []any{x})
		}
	case Fragment:
		if !IsNil(x) {
			return func(yield func(string, []any) bool) {
				for query, args := range x.Frag(ctx) {
					yield(query, args)
				}
			}
		}
		return Empty().Frag(ctx)
	case driver.Valuer:
		return Arg(x).Frag(ctx)
	case iter.Seq[any]:
		if f := Values[any](x); !IsNil(f) {
			return f.Frag(ctx)
		}
		return Empty().Frag(ctx)
	case []any:
		if len(x) > 0 {
			return Query(strings.Repeat(",?", len(x))[1:], x...).Frag(ctx)
		}
		return Empty().Frag(ctx)
	default:
		asArg := func(yield func(string, []any) bool) {
			yield("?", []any{x})
		}
		tpe := reflect.TypeOf(x)
		switch tpe.Kind() {
		case reflect.Slice:
			if !reflectx.IsBytes(tpe) {
				return ArgsIter(ctx, x)
			}
			return asArg
		case reflect.Func:
			if tpe.CanSeq() {
				rv := reflect.ValueOf(x)
				return Values[any](func(yield func(any) bool) {
					for xx := range rv.Seq() {
						yield(xx.Interface())
					}
				}).Frag(ctx)
			}
			return asArg
		default:
			return asArg
		}
	}
}

func ArgsIter(ctx context.Context, v any) Iter {
	switch x := v.(type) {
	case []bool:
		return Values[bool](slices.Values(x)).Frag(ctx)
	case []string:
		return Values[string](slices.Values(x)).Frag(ctx)
	case []float32:
		return Values[float32](slices.Values(x)).Frag(ctx)
	case []float64:
		return Values[float64](slices.Values(x)).Frag(ctx)
	case []int:
		return Values[int](slices.Values(x)).Frag(ctx)
	case []int8:
		return Values[int8](slices.Values(x)).Frag(ctx)
	case []int16:
		return Values[int16](slices.Values(x)).Frag(ctx)
	case []int32:
		return Values[int32](slices.Values(x)).Frag(ctx)
	case []int64:
		return Values[int64](slices.Values(x)).Frag(ctx)
	case []uint:
		return Values[uint](slices.Values(x)).Frag(ctx)
	case []uint8:
		return Values[uint8](slices.Values(x)).Frag(ctx)
	case []uint16:
		return Values[uint16](slices.Values(x)).Frag(ctx)
	case []uint32:
		return Values[uint32](slices.Values(x)).Frag(ctx)
	case []uint64:
		return Values[uint64](slices.Values(x)).Frag(ctx)
	case []any:
		return Values[any](slices.Values(x)).Frag(ctx)
	}

	rv := reflect.ValueOf(v)
	return Values[any](func(yield func(any) bool) {
		for i := 0; i < rv.Len(); i++ {
			yield(rv.Index(i).Interface())
		}
	}).Frag(ctx)
}

// Arg presents a single argument fragment
func Arg(v any) Fragment {
	return &argument{v: v}
}

type argument struct {
	v any
}

func (f *argument) IsNil() bool { return false }

func (f *argument) Frag(_ context.Context) Iter {
	return func(yield func(string, []any) bool) {
		yield("?", []any{f.v})
	}
}
