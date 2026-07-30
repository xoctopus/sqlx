package models

// TableColumn stores column-level SQL meta documentation.
// +genx:model
// @model TableName=sql_meta_table_column
// @model Register=MetaCatalog
// @model uidx=ui_column;Model;Col
type TableColumn struct {
	// Model table name.
	Model string `db:"f_tab,width=64"`
	// Col column name.
	Col string `db:"f_col,width=64"`
	// ColType column datatype string.
	ColType string `db:"f_col_typ,width=1024"`
	// Field struct field name.
	Field string `db:"f_field,width=64"`
	// Rel column relation.
	Rel string `db:"f_rel,width=128,defualt=''"`
	// Comment column documentation text.
	Comment string `db:"f_doc,width=1024"`
}
