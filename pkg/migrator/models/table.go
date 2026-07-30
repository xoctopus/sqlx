package models

// Table stores table-level SQL meta documentation.
// +genx:model
// @attr TableName=sql_meta_table
// @attr Register=MetaCatalog
// @def uidx ui_column Model
type Table struct {
	// Model is the table name.
	Model string `db:"f_tab,width=64"`
	// TabType is the Go type of the model.
	TabType string `db:"f_tab_type,width=1024"`
	// Comment is the table documentation text.
	Comment string `db:"f_doc,width=1024"`
}
