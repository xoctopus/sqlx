package types

// AutoIncID is an auto-increment primary key field.
//
// Deprecated: use [Serial] instead. uint64 conflicts with database/sql's
// int64-centric Value model (see Serial).
type AutoIncID struct {
	ID uint64 `db:"f_id,autoinc" json:"-"`
}

// ID is a generic unsigned identifier.
type ID uint64

// Serial is an auto-increment primary key embed (column f_id).
//
// ID is int64 rather than uint64 for two reasons:
//
//  1. SQL integer types used for identity columns are signed in the SQL
//     standard and in common dialects (INTEGER / BIGINT; PostgreSQL SERIAL /
//     BIGSERIAL map to signed integers). Unsigned 64-bit integers are not a
//     portable auto-increment type across databases.
//
//  2. database/sql treats integer driver.Value as int64. Passing a uint64 with
//     the high bit set (value >= 1<<63) fails with
//     "uint64 values with high bit set are not supported". Result.LastInsertId
//     also returns int64. Using int64 keeps insert args, scan targets, and
//     last-insert-id in one representation without conversion edge cases.
//
// ref:
//   - https://pkg.go.dev/database/sql/driver#Value
//   - https://www.postgresql.org/docs/current/datatype-numeric.html
type Serial struct {
	// ID 自增主键
	ID int64 `db:"f_id,autoinc" json:"-"`
}
