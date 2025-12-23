package session

import (
	"context"
	"fmt"

	"github.com/xoctopus/x/contextx"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/builder"
)

type OptionFunc func(*option)

type option struct {
	ReadOnly bool
}

func ReadOnly() OptionFunc {
	return func(o *option) {
		o.ReadOnly = true
	}
}

type Session interface {
	// Database physically endpoint
	Database() string
	// Schema logically
	Schema() string
	// T picks table from session
	T(any) builder.Table
	// Tx exec query
	Tx(context.Context, func(context.Context) error) error
	// Adaptor returns session adaptor
	Adaptor(...OptionFunc) adaptor.Adaptor
}

func New(a adaptor.Adaptor, name string) Session {
	return &session{
		database: a.Endpoint(),
		schema:   a.Schema(),
		name:     name,
		a:        a,
	}
}

func NewRO(rw adaptor.Adaptor, ro adaptor.Adaptor, name string) Session {
	return &session{
		database: ro.Endpoint(),
		schema:   ro.Schema(),
		name:     name,
		a:        rw,
		ro:       ro,
	}
}

func For(ctx context.Context, m any) Session {
	switch x := m.(type) {
	case interface{ Unwrap() builder.Model }:
		return For(ctx, x)
	case string:
		return FromContext(ctx, x)
	case builder.Model:
		// if t, ok := catalogs.Load(x.TableName()); ok {
		// 	x.TableName()
		// }
		// return FromContext(ctx)
	}
	return nil
}

func FromContext(ctx context.Context, name string) Session {
	s, ok := ctxSession.From(ctx)
	if ok && s.Schema() == name {
		return s
	}
	panic(fmt.Errorf("missing session '%s'", name))
}

type session struct {
	database string
	schema   string

	name string
	a    adaptor.Adaptor
	ro   adaptor.Adaptor
}

func (s *session) Schema() string {
	return s.schema
}

func (s *session) Database() string {
	return s.database
}

func (s *session) Name() string {
	return s.name
}

func (s *session) T(m any) builder.Table {
	switch x := m.(type) {
	case builder.WithTable:
		return x.T()
	case builder.Table:
		return x
	default:
		return builder.TFrom(m)
	}
}

func (s *session) Tx(ctx context.Context, exec func(context.Context) error) error {
	return s.a.Tx(ctx, exec)
}

func (s *session) Adaptor(options ...OptionFunc) adaptor.Adaptor {
	opt := &option{}
	for _, o := range options {
		o(opt)
	}

	if opt.ReadOnly {
		return s.ro
	}
	return s.a
}

var ctxSession = contextx.NewT[Session]()
