package models

// TableColumn
// +genx:model
// @attr TableName=sql_meta_table_column
// @attr Register=MetaCatalog
// @def uidx ui_column Model;Col
type TableColumn struct {
	Model   string `db:"f_tab,width=64"`
	Col     string `db:"f_col,width=64"`
	ColType string `db:"f_col_typ,width=1024"`
	Field   string `db:"f_field,width=64"`
	Rel     string `db:"f_rel,width=128,defualt=''"`
	Comment string `db:"f_doc,width=1024"`
}
