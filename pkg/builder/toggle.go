package builder

import (
	"context"

	"github.com/xoctopus/x/contextx"
)

type ToggleType int8

const (
	TOGGLE__MULTI_TABLE ToggleType = iota + 1
	TOGGLE__AUTO_ALIAS
	TOGGLE__ASSIGNMENTS
	TOGGLE__IN_PROJECT
	TOGGLE__SKIP_COMMENTS
	TOGGLE__SKIP_JOIN
)

type Toggles map[ToggleType]bool

// Injector?
// func (ts Toggles) Inject(ctx context.Context) context.Context {
// 	return ContextWithToggles(ctx, ts)
// }

func (ts Toggles) Merge(next Toggles) Toggles {
	final := Toggles{}
	for k, v := range ts {
		if v {
			final[k] = true
		}
	}
	for k, v := range next {
		if v {
			final[k] = true
		} else {
			delete(final, k)
		}
	}
	return final
}

func (ts Toggles) Is(key ToggleType) bool {
	if v, ok := ts[key]; ok {
		return v
	}
	return false
}

type ctxTogglesKey struct{}

func ContextWithToggles(ctx context.Context, ts Toggles) context.Context {
	return contextx.WithValue(
		ctx, ctxTogglesKey{},
		TogglesFromContext(ctx).Merge(ts),
	)
}

func WithToggles(ctx context.Context, toggles ...ToggleType) context.Context {
	for _, toggle := range toggles {
		ctx = ContextWithToggles(ctx, Toggles{toggle: true})
	}
	return ctx
}

func TrimToggles(ctx context.Context, toggles ...ToggleType) context.Context {
	for _, toggle := range toggles {
		ctx = ContextWithToggles(ctx, Toggles{toggle: false})
	}
	return ctx
}

func HasToggle(ctx context.Context, toggle ToggleType) bool {
	return TogglesFromContext(ctx).Is(toggle)
}

func TogglesFromContext(ctx context.Context) Toggles {
	if ctx == nil {
		return Toggles{}
	}
	if toggles, ok := ctx.Value(ctxTogglesKey{}).(Toggles); ok {
		return toggles
	}
	return Toggles{}
}
