// Package session binds a named database [Session] into context for Exec /
// Query / Tx, and looks it up by session name or model.
//
// Generated model methods call [MustFor](ctx, TUser) (or similar) to obtain the
// Session that should run against that table. Callers create a Session with
// [New] / [NewReadonly], inject it with [With] / [WithModel] /
// [WithSchemaModel], then pass the context down.
//
// This package blank-imports [github.com/xoctopus/sqlx/pkg/adaptors] so
// registered drivers are available when opening connections.
//
// # Session
//
//	a, err := adaptor.Open(ctx, dsn)
//	s := session.New(a, "main")
//	// or session.NewReadonly(rw, ro, "main")
//
//	ctx = session.With(ctx, s)
//	// and/or session.WithModel(ctx, &User{}, s)
//
//	_, err = session.MustFor(ctx, TUser).Exec(ctx, stmt)
//	// generated code typically uses:
//	// session.MustFor(ctx, TUser).Adaptor().Exec(...)
//
// [Session] exposes Name, T (resolve table), Adaptor (optional [ReadOnly]),
// and Exec / Query / Tx on the read-write adaptor. [InTx] reports whether ctx
// already carries an active *sql.Tx.
//
// # Context lookup
//
//   - [With]: bind under s.Name()
//   - [WithModel]: bind under m.TableName()
//   - [WithSchemaModel]: bind under schema.TableName
//   - [From] / [Must]: lookup by name
//   - [FromByModel] / [MustByModel]: lookup by model key, fallback TableName
//   - [For] / [MustFor]: name (string) or [builder.Model]
//   - [CarryModel] / [CarrySchemaModel]: contextx.Carrier helpers
//
// # Catalog Register
//
// [Register] records tables from one or more [builder.Catalog] values and
// panics if the same table key is registered twice (TableName, or
// Schema.TableName). This is separate from `@model Register=Catalog`, which
// only adds generated tables into a Catalog variable at init.
package session
