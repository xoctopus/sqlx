package models

// TableColumn stores column-level SQL meta documentation.
// +genx:model
// @attr TableName=sql_meta_table_column
// @attr Register=MetaCatalog
// @def uidx ui_column Model;Col
type TableColumn struct {
	// Model is the table name.
	Model string `db:"f_tab,width=64"`
	// Col is the column name.
	Col string `db:"f_col,width=64"`
	// ColType is the column datatype string.
	ColType string `db:"f_col_typ,width=1024"`
	// Field is the struct field name.
	Field string `db:"f_field,width=64"`
	// Rel is the column relation text.
	Rel string `db:"f_rel,width=128,defualt=''"`
	// Comment is the column documentation text.
	Comment string `db:"f_doc,width=1024"`
}
