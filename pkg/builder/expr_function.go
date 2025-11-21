package builder

import (
	"context"

	"github.com/xoctopus/sqlx/pkg/frag"
)

func Count(fragments ...frag.Fragment) *Function {
	if len(fragments) == 0 {
		return Func("COUNT", frag.Lit("1"))
	}
	return Func("COUNT", fragments...)
}

func Avg(fragments ...frag.Fragment) *Function {
	return Func("AVG", fragments...)
}

func AnyValue(fragments ...frag.Fragment) *Function {
	return Func("ANY_VALUE", fragments...)
}

func Distinct(fragments ...frag.Fragment) *Function {
	return Func("DISTINCT", fragments...)
}

func Min(fragments ...frag.Fragment) *Function {
	return Func("MIN", fragments...)
}

func Max(fragments ...frag.Fragment) *Function {
	return Func("MAX", fragments...)
}

func First(fragments ...frag.Fragment) *Function {
	return Func("FIRST", fragments...)
}

func Last(fragments ...frag.Fragment) *Function {
	return Func("LAST", fragments...)
}

func Sum(fragments ...frag.Fragment) *Function {
	return Func("SUM", fragments...)
}

func Func(name string, args ...frag.Fragment) *Function {
	if name == "" {
		return nil
	}
	return &Function{
		name: name,
		args: args,
	}
}

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
