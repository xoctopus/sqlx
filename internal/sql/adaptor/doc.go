// Package adaptor defines the database driver abstraction used by sqlx.
//
// An [Adaptor] is a connected, driver-specific handle: it executes
// [frag.Fragment] values ([DB]), exposes a [Dialect] for DDL/DML rendering and
// error classification, and can scan the live schema into a [builder.Catalog].
// Higher layers such as migrator and session depend on this package; public
// callers usually reach it via github.com/xoctopus/sqlx/pkg/adaptors or
// session.New.
//
// # Interfaces
//
//   - [DB]: Exec / Query / Tx over Fragments (backed by *sql.DB or a ctx Executor)
//   - [Adaptor]: DB + DriverName / Schema / Dialect / Catalog
//   - [Connector]: Open from a parsed DSN URL
//   - [Dialect]: CreateTable / AddColumn / AddIndex / DBType / conflict detection, etc.
//
// # Registration and Open
//
// Driver implementations call [Register] in init (optionally with aliases).
// [Open] parses a DSN (`scheme://...`), looks up the scheme, and delegates to
// that driver's [Connector.Open].
//
//	a, err := adaptor.Open(ctx, "mysql://user:pass@tcp(localhost:3306)/db")
//
// Blank-import github.com/xoctopus/sqlx/pkg/adaptors (or a concrete driver
// package such as mysql) so Register runs before Open.
//
// # Execution helpers
//
//   - [Wrap]: adapt *sql.DB to [DB]
//   - [WithExecutor] / [ExecutorFrom]: bind *sql.DB or *sql.Tx into context so
//     nested Exec/Query reuse the same executor (used by [DB.Tx])
//   - [DatabaseNameFromDSN]: extract database/schema name from a DSN URL
//
// # Driver packages
//
//   - mysql: full Adaptor + Dialect (registered on import)
//   - postgres / sqlite: undone
package adaptor
