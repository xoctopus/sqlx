// Package loggingdriver wraps a database/sql/driver.Driver to log SQL execution
// and optionally interpolate placeholders for readable logs.
//
// Adaptor implementations (for example mysql) call [Wrap] around the native
// driver before sql.OpenDB, so every Query / Exec / transaction is timed and
// recorded through logx with driver name, interpolated query text, and cost.
//
// # Wrap
//
//	dc := loggingdriver.Wrap(mysql.MySQLDriver{}, "mysql",
//		loggingdriver.WithDsnParser(ParseDSN),
//		loggingdriver.WithErrorLeveler(ErrorLevel),
//		loggingdriver.WithInterpolator(loggingdriver.DefaultInterpolate),
//	)
//	conn, err := dc.OpenConnector(dsn)
//
// [Wrap] returns a [driver.DriverContext]. OpenConnector may rewrite the DSN
// via [WithDsnParser], and marks read-only pools when the DSN query has
// `_ro=true` (driver name becomes `name::ro` for logs).
//
// # Logging behavior
//
//   - QueryContext / ExecContext: optional interpolate, then log cost_ms and
//     query; failures use Error or Warn based on [WithErrorLeveler] (level > 0
//     → Error, else Warn); success → Debug
//   - BeginTx / Commit / Rollback: transaction boundary Debug lines
//   - Prepare is deliberately unimplemented (panics) so dialects stay on
//     context Exec/Query paths
//
// # Interpolation
//
//   - [DefaultInterpolate]: substitute `?` with literal SQL values (for logs
//     or when the native driver does not interpolate); returns empty remaining args
//   - [OrderedInterpolator]: rewrite `?` to `$1`, `$2`, ... (PostgreSQL-style)
//   - [Interpolate] / [NewPrinter]: core substitution helpers used by logging
//
// Migrator dry-run output also uses [DefaultInterpolate] to print a full SQL
// script from collected fragments.
package loggingdriver
