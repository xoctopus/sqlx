package types

// AutoIncID is an auto-increment primary key field.
type AutoIncID struct {
	ID uint64 `db:"f_id,autoinc" json:"-"`
}

// ID is a generic unsigned identifier.
type ID uint64
