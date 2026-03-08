module github.com/xoctopus/sqlx

go 1.26.0

tool github.com/xoctopus/sqlx/internal/cmd/example

require (
	github.com/xoctopus/genx v0.2.0
	github.com/xoctopus/logx v0.3.0
	github.com/xoctopus/pkgx v0.4.0
	github.com/xoctopus/typx v0.4.0
	github.com/xoctopus/x v0.4.0
)

// drivers
require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/go-sql-driver/mysql v1.9.3
)

// extended datatypes
require github.com/shopspring/decimal v1.4.0

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/tools v0.42.0 // indirect
)
