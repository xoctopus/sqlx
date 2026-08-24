module github.com/xoctopus/sqlx

go 1.26.5

tool (
	github.com/xoctopus/sqlx/internal/cmd/example
	github.com/xoctopus/sqlx/internal/cmd/skill-install
)

require (
	// +skill:genx
	github.com/xoctopus/genx v0.3.7
	// +skill:testx
	github.com/xoctopus/x v0.5.5
)

require (
	github.com/xoctopus/logx v0.3.5
	github.com/xoctopus/pkgx v0.4.3
	github.com/xoctopus/typx v0.4.6
)

// drivers
require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/go-sql-driver/mysql v1.10.0
	github.com/jackc/pgx/v5 v5.10.0
	modernc.org/sqlite v1.57.0
)

// extended datatypes
require (
	github.com/oklog/ulid/v2 v2.1.2
	github.com/shopspring/decimal v1.4.0
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	golang.org/x/tools v0.48.0 // indirect
	modernc.org/libc v1.74.4 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)
