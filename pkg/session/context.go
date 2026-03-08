package session

import (
	"context"
	"fmt"

	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/pkg/builder"
)

type ModelWrapper interface {
	Unwrap() builder.Table
}

// For retrieves Session by session name or table
func For(ctx context.Context, m any) (Session, bool) {
	switch x := m.(type) {
	case interface{ Unwrap() builder.Model }:
		return For(ctx, x.Unwrap().TableName())
	case interface{ Unwrap() builder.Table }:
		return For(ctx, x.Unwrap().TableName())
	case interface{ TableName() string }:
		return From(ctx, x.TableName())
	case string:
		return From(ctx, x)
	default:
		return nil, false
	}
}

func MustFor(ctx context.Context, m any) Session {
	switch x := m.(type) {
	case interface{ Unwrap() builder.Model }:
		return MustFor(ctx, x.Unwrap().TableName())
	case interface{ Unwrap() builder.Table }:
		return MustFor(ctx, x.Unwrap().TableName())
	case interface{ TableName() string }:
		return Must(ctx, x.TableName())
	case string:
		return Must(ctx, x)
	default:
		panic(fmt.Errorf("missing session invalid type: %T", x))
	}
}

type tSessionKey struct {
	name string
}

// From retrieve Session from ctx by Session.Name
func From(ctx context.Context, name string) (Session, bool) {
	s, ok := ctx.Value(tSessionKey{name}).(Session)
	// must.BeTrueF(ok, "missing session: %s", name)
	return s, ok
}

func Must(ctx context.Context, name string) Session {
	s, ok := From(ctx, name)
	must.BeTrueF(ok, "missing session for: %s", name)
	return s
}

// With injects Session
func With(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, tSessionKey{name: session.Name()}, session)
}

// Carry returns context carrier
func Carry(session Session) contextx.Carrier {
	return func(ctx context.Context) context.Context {
		return With(ctx, session)
	}
}
