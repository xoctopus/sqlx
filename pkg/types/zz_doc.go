// Package types provides common SQL column types and model field embeds.
//
// # Basic time types (package sqltime)
//
//   - sqltime.Datetime: calendar date-time stored as a SQL datetime type
//     (MySQL DATETIME, PostgreSQL TIMESTAMP).
//   - sqltime.Timestamp: integer unix timestamp at second precision.
//   - sqltime.TimestampMilli: integer unix timestamp at millisecond precision.
//
// # Operation-time embeds over sqltime.Datetime
//
// Combinations of CreatedAt / UpdatedAt / DeletedAt for record lifecycle.
// Soft-deletion variants (CreationModificationDeletion*) implement SoftDeletion.
// Aliases: OperationDatetime, OperationDatetimePrecise.
//
// Second precision:
//
//   - CreationDatetime
//   - CreationModificationDatetime
//   - CreationModificationDeletionDatetime
//   - sqlops.CreationDatetime
//   - sqlops.CreationModificationDatetime
//   - sqlops.CreationModificationDeletionDatetime
//
// Millisecond precision (Datetime with precision=3):
//
//   - CreationDatetimeMilli
//   - CreationModificationDatetimePrecise
//   - CreationModificationDeletionDatetimePrecise
//   - sqlops.CreationDatetimeMilli
//   - sqlops.CreationModificationDatetimePrecise
//   - sqlops.CreationModificationDeletionDatetimePrecise
//
// Column names differ: types uses an `f_` prefix (e.g. f_created_at);
// sqlops uses unprefixed names (e.g. created_at).
//
// # Operation-time embeds over sqltime.Timestamp / TimestampMilli
//
// Second precision (sqltime.Timestamp):
//
//   - CreationTime
//   - CreationModificationTime
//   - CreationModificationDeletionTime
//   - sqlops.CreationTime
//   - sqlops.CreationModificationTime
//   - sqlops.CreationModificationDeletionTime
//
// Millisecond precision (sqltime.TimestampMilli):
//
//   - CreationTimePrecise
//   - CreationModificationTimePrecise
//   - CreationModificationDeletionTimePrecise
//   - sqlops.CreationTimePrecise
//   - sqlops.CreationModificationTimePrecise
//   - sqlops.CreationModificationDeletionTimePrecise
//
// Soft-deletion variants implement SoftDeletion.
// Aliases: OperationTime, OperationTimePrecise.
// Same f_ vs unprefixed column naming as the Datetime embeds above.
//
// # Bool
//
// types.Bool is a three-state boolean (TRUE / FALSE / UNKNOWN ) for avoiding
// ambiguity with Go's two-state bool. Use Boolean(b) to convert from bool.
//
// # Text
//
// types.Text wraps SQL TEXT. An empty string is stored as NULL.
//
// # Blob
//
// types.Blob wraps SQL BLOB (BYTEA on PostgreSQL). An empty slice is stored as NULL.
//
// # Decimal
//
// types.Decimal is a high-precision decimal for monetary and exact numeric columns.
//
// # JSON
//
// types.JSONArray and types.JSONObject marshal Go values to/from JSON columns.
//
// # Other
//
// AutoIncID is an embed for an auto-increment primary key (f_id).
// ID is a generic unsigned identifier type.
//
// Marker interfaces CreationMarker, ModificationMarker, DeletionMarker, and
// SoftDeletion describe create / update / soft-delete behavior used by model
// code generation.
// +genx:doc
package types
