// Package structs walks Go model structs into SQL column field metadata.
//
// [builder.TFrom], the genx model generator, [scanner.Scan], and helper
// insert helpers all rely on this package to discover `db`-tagged fields,
// including nested embeds such as AutoIncID / Rel* / Meta / State.
//
// # FieldsFor
//
//	fields := structs.FieldsFor(typx.NewRType(reflect.TypeFor[User]()))
//	// f_id, f_user_id, ... with Loc indexes into the struct tree
//
// [FieldsFor] / [FieldsSeqFor] walk an exported struct type (results are
// cached per type). Rules:
//
//   - skip unexported fields and `db:"-"`
//   - expand anonymous (or same-named) struct embeds that are not driver.Valuer
//   - column name from `db` tag name, else lowercased field name
//   - parse column options into [Field.ColumnDef] (width, autoinc, ...)
//   - keep json/name tags on [Field.NameTag]
//
// Each [Field] carries ColumnName, FieldName, Type, Loc (field index path),
// and ColumnDef for DDL / scan / insert.
//
// # TableFields
//
//	tfs := structs.TableFields(&user)
//
// [TableFields] / [TableFieldsSeq] bind FieldsFor to a live value: each
// [TableField] has the Field metadata, the addressable [TableField.Value],
// and TableName when v implements Model. Deprecated columns
// (`db:",deprecated=..."`) are omitted.
//
// See ExampleFieldsFor for embed flattening output.
package structs
