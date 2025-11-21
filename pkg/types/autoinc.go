package types

type AutoIncID struct {
	ID uint64 `db:"f_id,autoinc" json:"-"`
}
