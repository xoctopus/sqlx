// Package def parses column and index definitions used across sqlx.
//
// [structs.FieldsFor] attaches a [ColumnDef] to each model field; [builder]
// and dialects consume it for DDL. The genx model generator uses [ParseKeyDef]
// to turn `@model pk|idx|uidx=...` annotations into [KeyDefine] values.
//
// # ColumnDef
//
//	d := def.ParseColDef(t, `db:"name,width=128,default='',null"`)
//
// [ParseColDef] reads the `db` struct tag options into [ColumnDef]:
//
//   - null / autoinc / unsigned
//   - width / precision
//   - default (raw string, including quotes or SQL functions)
//   - onupdate
//   - deprecated (optional rename target for migrations)
//
// [DefineFromCatalog] / [DefineFromUser] / [DefineFromGoType] record where a
// definition originated when filled later by catalog scan or dialect.
//
// # KeyDefine
//
// [ParseKeyDef] parses `@model` index annotation values. See that function's
// comment and ExampleParseKeyDef for annotation → parse mappings.
//
//	pk=<field>[;<field>...]
//	idx=<name[,method]>;<field[,options...]>[;...]
//	uidx=<name[,method]>;<field[,options...]>[;...]
//
// Fields are separated by `;`, options by `,`. [KeyColumnOption] holds a field
// (or column) name plus per-column options such as NULLS,FIRST.
// Helpers: [ResolveIndexNameAndUsing], [ResolveKeyColumnOptions],
// [KeyColumnOptionByNames].
package def
