package models

// Table
// +genx:model
// @attr TableName=sql_meta_table
// @attr Register=MetaCatalog
// @def uidx ui_column Model
type Table struct {
	Model   string `db:"f_tab,width=64"`
	TabType string `db:"f_tab_type,width=1024"`
	Comment string `db:"f_doc,width=1024"`
}
