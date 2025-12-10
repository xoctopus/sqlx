package builder

import "github.com/xoctopus/sqlx/internal"

type (
	WithTableDesc interface {
		TableDesc() []string
	}

	WithPrimaryKey interface {
		PrimaryKey() []string
	}

	WithUniqueIndexes interface {
		// UniqueIndexes returns field name => []string{index name, options...}
		UniqueIndexes() map[string][]string
	}

	WithIndexes interface {
		// Indexes returns field name => []string{index name, options...}
		Indexes() map[string][]string
	}

	WithColumnComment interface {
		ColumnComment() map[string]string
	}

	WithColumnDesc interface {
		ColumnDesc() map[string][]string
	}

	WithColumnRel interface {
		ColumnRel() map[string][]string
	}

	WithDatatypeDesc interface {
		DBType(driver string) string
	}

	Model = internal.Model
)
