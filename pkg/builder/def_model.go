package builder

type (
	WithTableDesc interface {
		TableDesc() []string
	}

	WithPrimaryKey interface {
		PrimaryKey() []string
	}

	WithUniqueIndexes interface {
		UniqueIndexes() map[string][]string
	}

	WithIndexes interface {
		Indexes() map[string][]string
	}

	WithColumnComments interface {
		Comments() map[string]string
	}

	WithColumnDesc interface {
		ColumnDesc() map[string][]string
	}

	WithColumnRel interface {
		ColumnRel() map[string][]string
	}
)
