module github.com/xoctopus/sqlx

go 1.25.5

tool github.com/xoctopus/sqlx/internal/cmd/example

require (
	github.com/xoctopus/confx v0.2.14
	github.com/xoctopus/genx v0.1.16
	github.com/xoctopus/logx v0.1.1
	github.com/xoctopus/pkgx v0.1.10
	github.com/xoctopus/typx v0.3.4
	github.com/xoctopus/x v0.2.11
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
	github.com/go-think/openssl v1.20.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/tools v0.40.0 // indirect
)
