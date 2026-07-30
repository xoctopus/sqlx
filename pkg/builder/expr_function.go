package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

// Count builds COUNT(...); empty args become COUNT(1).
func Count(fragments ...frag.Fragment) *Function {
	if len(fragments) == 0 {
		return Func("COUNT", frag.Lit("1"))
	}
	return Func("COUNT", fragments...)
}

// Avg builds AVG(...).
func Avg(fragments ...frag.Fragment) *Function {
	return Func("AVG", fragments...)
}

// AnyValue builds ANY_VALUE(...).
func AnyValue(fragments ...frag.Fragment) *Function {
	return Func("ANY_VALUE", fragments...)
}

// Distinct builds DISTINCT(...).
func Distinct(fragments ...frag.Fragment) *Function {
	return Func("DISTINCT", fragments...)
}

// Min builds MIN(...).
func Min(fragments ...frag.Fragment) *Function {
	return Func("MIN", fragments...)
}

// Max builds MAX(...).
func Max(fragments ...frag.Fragment) *Function {
	return Func("MAX", fragments...)
}

// First builds FIRST(...).
func First(fragments ...frag.Fragment) *Function {
	return Func("FIRST", fragments...)
}

// Last builds LAST(...).
func Last(fragments ...frag.Fragment) *Function {
	return Func("LAST", fragments...)
}

// Sum builds SUM(...).
func Sum(fragments ...frag.Fragment) *Function {
	return Func("SUM", fragments...)
}

// Func builds a named SQL function call.
func Func(name string, args ...frag.Fragment) *Function {
	if name == "" {
		return nil
	}
	return &Function{
		name: name,
		args: args,
	}
}

// Function is a SQL function call fragment.
type Function struct {
	name string
	args []frag.Fragment
}

func (f *Function) IsNil() bool {
	return f == nil || f.name == ""
}

func (f *Function) Frag(ctx context.Context) frag.Iter {
	return func(yield func(string, []any) bool) {
		yield(f.name, nil)

		if len(f.args) == 0 {
			for q, args := range frag.Block(frag.Lit("*")).Frag(ctx) {
				yield(q, args)
			}
			return
		}

		iter := frag.Block(frag.Compose(",", f.args...)).Frag(ctx)
		for q, args := range iter {
			yield(q, args)
		}
	}
}
