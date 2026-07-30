package session

import (
	"context"
	"database/sql"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	"github.com/xoctopus/sqlx/pkg/sql/adaptor"
	_ "github.com/xoctopus/sqlx/pkg/sql/adaptor/mysql"
	_ "github.com/xoctopus/sqlx/pkg/sql/adaptor/postgres"
	_ "github.com/xoctopus/sqlx/pkg/sql/adaptor/sqlite"
)

// Session is a named database access handle.
type Session interface {
	// Name returns the logical session name.
	Name() string
	// T resolves a table from a model, Table, or WithTable value.
	T(any) builder.Table
	// Adaptor returns the underlying adaptor; ReadOnly() selects the RO adaptor when set.
	Adaptor(...AdaptorOptionApplier) adaptor.Adaptor

	// Tx runs fn inside a transaction on the RW adaptor.
	Tx(context.Context, func(context.Context) error) error
	// Exec executes f on the RW adaptor.
	Exec(context.Context, frag.Fragment) (sql.Result, error)
	// Query runs f on the RW adaptor.
	Query(context.Context, frag.Fragment) (*sql.Rows, error)
}

// New creates a Session backed by a single adaptor.
func New(a adaptor.Adaptor, name string) Session {
	return &session{
		name: name,
		a:    a,
	}
}

// NewReadonly creates a Session with separate read-write and read-only adaptors.
func NewReadonly(rw adaptor.Adaptor, ro adaptor.Adaptor, name string) Session {
	return &session{
		name: name,
		a:    rw,
		ro:   ro,
	}
}

type session struct {
	name string

	a  adaptor.Adaptor
	ro adaptor.Adaptor
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

func (s *session) Exec(ctx context.Context, f frag.Fragment) (sql.Result, error) {
	return s.a.Exec(ctx, f)
}

func (s *session) Query(ctx context.Context, f frag.Fragment) (*sql.Rows, error) {
	return s.a.Query(ctx, f)
}

func (s *session) Adaptor(appliers ...AdaptorOptionApplier) adaptor.Adaptor {
	opt := &AdaptorOption{}
	for _, apply := range appliers {
		apply(opt)
	}

	if opt.ReadOnly {
		return s.ro
	}
	return s.a
}

// AdaptorOption configures Adaptor selection.
type AdaptorOption struct {
	// ReadOnly selects the read-only adaptor when true.
	ReadOnly bool
}

// AdaptorOptionApplier mutates AdaptorOption.
type AdaptorOptionApplier func(*AdaptorOption)

// ReadOnly selects the read-only adaptor from NewReadonly.
func ReadOnly() AdaptorOptionApplier {
	return func(o *AdaptorOption) {
		o.ReadOnly = true
	}
}
