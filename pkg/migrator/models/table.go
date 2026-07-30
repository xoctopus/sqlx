package models

// Table stores table-level SQL meta documentation.
// +genx:model
// @model TableName=sql_meta_table
// @model Register=MetaCatalog
// @model uidx=ui_column;Model
type Table struct {
	// Model table name.
	Model string `db:"f_tab,width=64"`
	// TabType Go type of the model.
	TabType string `db:"f_tab_type,width=1024"`
	// Comment table documentation text.
	Comment string `db:"f_doc,width=1024"`
}
