package models

// Enumeration stores enum-value SQL meta documentation.
// +genx:model
// @attr TableName=sql_meta_enumeration
// @attr Register=MetaCatalog
type Enumeration struct {
	// Model is the table name.
	Model string `db:"f_tab,width=64"`
	// Col is the column name.
	Col string `db:"f_col,width=64"`
	// EnumType is the Go enum type name.
	EnumType string `db:"f_enum_type,width=1024"`
	// Value is the enum numeric value.
	Value string `db:"f_value,width=64"`
	// Kind is the underlying kind of the enum value.
	Kind string `db:"f_kind,width=16"`
	// Key is the enum constant name.
	Key string `db:"f_key,width=64"`
	// Text is the enum display text.
	Text string `db:"f_text,width=128"`
}
