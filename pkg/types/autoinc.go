package types

type AutoIncID struct {
	ID uint64 `db:"f_id,autoinc" json:"-"`
}

type AutoIncPrimary struct {
	ID uint64 `db:"f_id,autoinc" json:"-"`
}

func (AutoIncPrimary) PrimaryKey() []string {
	return []string{"ID"}
}

type ID uint64
