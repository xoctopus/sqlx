package frag

import "context"

func Func(f func(context.Context) Iter) Fragment {
	return _func(f)
}

type _func func(context.Context) Iter

func (f _func) IsNil() bool {
	return f == nil
}

func (f _func) Frag(ctx context.Context) Iter {
	return f(ctx)
}
