package builder

import "github.com/xoctopus/sqlx/internal"

type (
	// WithTableDesc provides table description lines.
	WithTableDesc interface {
		TableDesc() []string
	}

	// WithPrimaryKey provides primary key field names.
	WithPrimaryKey interface {
		PrimaryKey() []string
	}

	// WithUniqueIndexes provides unique index definitions.
	WithUniqueIndexes interface {
		// UniqueIndexes returns field name => []string{index name, options...}
		UniqueIndexes() map[string][]string
	}

	// WithIndexes provides non-unique index definitions.
	WithIndexes interface {
		// Indexes returns field name => []string{index name, options...}
		Indexes() map[string][]string
	}

	// WithColumnComment provides per-field column comments.
	WithColumnComment interface {
		ColumnComment() map[string]string
	}

	// WithColumnDesc provides per-field column descriptions.
	WithColumnDesc interface {
		ColumnDesc() map[string][]string
	}

	// WithColumnRel provides per-field column relations.
	WithColumnRel interface {
		ColumnRel() map[string][]string
	}

	// WithDatatypeDesc lets a custom domain type declare its SQL datatype so the
	// Go type and column definition stay aligned across dialects.
	//
	// Typical uses:
	//   - Fixed domain scale, e.g. Money as DECIMAL(22,4):
	//
	//	type Money struct{ decimal.Decimal }
	//	func (Money) DBType(string) string             { return "decimal" }
	//	func (Money) DBFixedWidth(string) *uint64      { return new(uint64(22)) }
	//	func (Money) DBFixedPrecision(string) *uint64  { return new(uint64(4)) }
	//
	//   - Dialect-specific storage, e.g. types.ULID: uuid on postgres/duckdb,
	//     blob on sqlite, binary(16) on mysql.
	//
	// Pair with [WithFixedWidthDesc] / [WithFixedPrecisionDesc] when needed.
	WithDatatypeDesc interface {
		DBType(driver string) string
	}

	// WithFixedWidthDesc lets a custom domain type fix its storage width.
	// A non-nil return overrides the db tag width. See [WithDatatypeDesc].
	WithFixedWidthDesc interface {
		DBFixedWidth(driver string) *uint64
	}

	// WithFixedPrecisionDesc lets a custom domain type fix its storage precision.
	// A non-nil return overrides the db tag precision. See [WithDatatypeDesc].
	WithFixedPrecisionDesc interface {
		DBFixedPrecision(driver string) *uint64
	}

	// Newer creates a new model instance.
	Newer interface {
		New() Model
	}

	// Model is a table-backed model.
	Model = internal.Model
)
