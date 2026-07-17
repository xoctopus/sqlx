package models

// Enumeration
// +genx:model
// @attr TableName=sql_meta_enumeration
// @attr Register=MetaCatalog
type Enumeration struct {
	Model string `db:"f_tab,width=64"`
	Col   string `db:"f_col,width=64"`
	Enum  string `db:"f_enum,width=255"`
	Value string `db:"f_value,width=64"`
	Kind  string `db:"f_kind,width=16"`
	Key   string `db:"f_key,width=64"`
	Text  string `db:"f_text,width=128"`
}
