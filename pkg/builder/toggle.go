package builder

import (
	"context"

	"github.com/xoctopus/x/contextx"
)

// ToggleType is a SQL rendering toggle.
type ToggleType int8

const (
	// TOGGLE__MULTI_TABLE qualifies columns with table names.
	TOGGLE__MULTI_TABLE ToggleType = iota + 1
	// TOGGLE__AUTO_ALIAS aliases multi-table columns automatically.
	TOGGLE__AUTO_ALIAS
	// TOGGLE__TUPLE_ASSIGNMENTS enables tuple-style assignments.
	TOGGLE__TUPLE_ASSIGNMENTS
	// TOGGLE__IN_PROJECT marks rendering inside a projection list.
	TOGGLE__IN_PROJECT
	// TOGGLE__SKIP_COMMENTS skips comment additions.
	TOGGLE__SKIP_COMMENTS
	// TOGGLE__SKIP_JOIN skips join additions.
	TOGGLE__SKIP_JOIN
)

// Toggles is a set of enabled ToggleType values.
type Toggles map[ToggleType]bool

// Injector?
// func (ts Toggles) Inject(ctx context.Context) context.Context {
// 	return ContextWithToggles(ctx, ts)
// }

// Merge combines toggles; false values in next remove keys.
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

// Is reports whether key is enabled.
func (ts Toggles) Is(key ToggleType) bool {
	if v, ok := ts[key]; ok {
		return v
	}
	return false
}

type ctxTogglesKey struct{}

// ContextWithToggles merges ts into the context toggles.
func ContextWithToggles(ctx context.Context, ts Toggles) context.Context {
	return contextx.WithValue(
		ctx, ctxTogglesKey{},
		TogglesFromContext(ctx).Merge(ts),
	)
}

// WithToggles enables the given toggles on ctx.
func WithToggles(ctx context.Context, toggles ...ToggleType) context.Context {
	for _, toggle := range toggles {
		ctx = ContextWithToggles(ctx, Toggles{toggle: true})
	}
	return ctx
}

// TrimToggles disables the given toggles on ctx.
func TrimToggles(ctx context.Context, toggles ...ToggleType) context.Context {
	for _, toggle := range toggles {
		ctx = ContextWithToggles(ctx, Toggles{toggle: false})
	}
	return ctx
}

// HasToggle reports whether toggle is enabled in ctx.
func HasToggle(ctx context.Context, toggle ToggleType) bool {
	return TogglesFromContext(ctx).Is(toggle)
}

// TogglesFromContext returns toggles stored in ctx.
func TogglesFromContext(ctx context.Context) Toggles {
	if ctx == nil {
		return Toggles{}
	}
	if toggles, ok := ctx.Value(ctxTogglesKey{}).(Toggles); ok {
		return toggles
	}
	return Toggles{}
}
