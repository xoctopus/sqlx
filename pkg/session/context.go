package session

import (
	"context"

	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/pkg/builder"
)

// For retrieves Session by session name or table
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

func MustFor(ctx context.Context, m any) Session {
	s, ok := For(ctx, m)
	must.BeTrueF(ok, "missing session for %T", m)
	return s
}

type tSession struct {
	name string
}

func From(ctx context.Context, name string) (Session, bool) {
	s, ok := ctx.Value(tSession{name}).(Session)
	return s, ok
}

func FromByModel(ctx context.Context, m builder.Model) (Session, bool) {
	s, ok := From(ctx, keyOf(m))
	if !ok {
		s, ok = From(ctx, m.TableName())
	}
	return s, ok
}

func Must(ctx context.Context, name string) Session {
	s, ok := From(ctx, name)
	must.BeTrueF(ok, "missing session for: %s", name)
	return s
}

func MustByModel(ctx context.Context, m builder.Model) Session {
	s, ok := FromByModel(ctx, m)
	must.BeTrueF(ok, "missing session for: %T", m)
	return s
}

func With(ctx context.Context, s Session) context.Context {
	return context.WithValue(ctx, tSession{name: s.Name()}, s)
}

func WithModel(ctx context.Context, m builder.Model, s Session) context.Context {
	v := tSession{name: m.TableName()}
	return context.WithValue(ctx, v, s)
}

func WithSchemaModel(ctx context.Context, m builder.Model, s Session) context.Context {
	v := tSession{name: keyOf(m)}
	return context.WithValue(ctx, v, s)
}

func CarryModel(m builder.Model, s Session) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return WithModel(ctx, m, s)
	}
}

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
