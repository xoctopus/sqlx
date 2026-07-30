package models

// Enumeration stores enum-value SQL meta documentation.
// +genx:model
// @model TableName=sql_meta_enumeration
// @model Register=MetaCatalog
type Enumeration struct {
	// Model table name.
	Model string `db:"f_tab,width=64"`
	// Col column name.
	Col string `db:"f_col,width=64"`
	// EnumType Go enum type name.
	EnumType string `db:"f_enum_type,width=1024"`
	// Value enum numeric value.
	Value string `db:"f_value,width=64"`
	// Kind underlying kind of the enum value.
	Kind string `db:"f_kind,width=16"`
	// Key enum constant name.
	Key string `db:"f_key,width=64"`
	// Text enum display text.
	Text string `db:"f_text,width=128"`
}
