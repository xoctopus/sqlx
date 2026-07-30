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

	// WithDatatypeDesc provides a driver-specific datatype override.
	WithDatatypeDesc interface {
		DBType(driver string) string
	}

	// Newer creates a new model instance.
	Newer interface {
		New() Model
	}

	// Model is a table-backed model.
	Model = internal.Model
)
