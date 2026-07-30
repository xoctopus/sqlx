package session

import (
	"context"

	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/pkg/builder"
)

// For retrieves a Session by session name (string) or model (builder.Model).
func For(ctx context.Context, m any) (Session, bool) {
	switch x := m.(type) {
	case string:
		return From(ctx, x)
	case builder.Model:
		return FromByModel(ctx, x)
	default:
		return nil, false
	}
}

// MustFor is like For but panics when the session is missing.
func MustFor(ctx context.Context, m any) Session {
	s, ok := For(ctx, m)
	must.BeTrueF(ok, "missing session for %T", m)
	return s
}

type tSession struct {
	name string
}

// From retrieves a Session bound under the given name.
func From(ctx context.Context, name string) (Session, bool) {
	s, ok := ctx.Value(tSession{name}).(Session)
	return s, ok
}

// FromByModel retrieves a Session by schema-qualified key, falling back to TableName.
func FromByModel(ctx context.Context, m builder.Model) (Session, bool) {
	s, ok := From(ctx, keyOf(m))
	if !ok {
		s, ok = From(ctx, m.TableName())
	}
	return s, ok
}

// Must is like From but panics when the session is missing.
func Must(ctx context.Context, name string) Session {
	s, ok := From(ctx, name)
	must.BeTrueF(ok, "missing session for: %s", name)
	return s
}

// MustByModel is like FromByModel but panics when the session is missing.
func MustByModel(ctx context.Context, m builder.Model) Session {
	s, ok := FromByModel(ctx, m)
	must.BeTrueF(ok, "missing session for: %T", m)
	return s
}

// With binds s into ctx under s.Name().
func With(ctx context.Context, s Session) context.Context {
	return context.WithValue(ctx, tSession{name: s.Name()}, s)
}

// WithModel binds s into ctx under m.TableName().
func WithModel(ctx context.Context, m builder.Model, s Session) context.Context {
	v := tSession{name: m.TableName()}
	return context.WithValue(ctx, v, s)
}

// WithSchemaModel binds s into ctx under the schema-qualified model key.
func WithSchemaModel(ctx context.Context, m builder.Model, s Session) context.Context {
	v := tSession{name: keyOf(m)}
	return context.WithValue(ctx, v, s)
}

// CarryModel returns a Carrier that applies WithModel.
func CarryModel(m builder.Model, s Session) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return WithModel(ctx, m, s)
	}
}

// CarrySchemaModel returns a Carrier that applies WithSchemaModel.
func CarrySchemaModel(m builder.Model, s Session) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return WithSchemaModel(ctx, m, s)
	}
}

func keyOf(m any) string {
	k := ""
	if x, ok := m.(builder.Model); ok {
		k = x.TableName()
	}
	if x, ok := m.(builder.HasSchema); ok {
		k = x.Schema() + "." + k
	}
	return k
}
