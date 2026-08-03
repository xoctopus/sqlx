// Package modeled provides [builder] Table, Col, and Key variants parameterized
// by a concrete [Model], so the compiler knows which model a fragment belongs to.
//
// # Why modeled scope
//
// [builder.Table], [builder.Col], and [builder.Key] are model-agnostic: the
// same column name can belong to any table. modeled tells the compiler which
// Model owns a Table, Col, or Key.
//
// e.g.
//
//	Table[User]: modeled table by User model.
//	TCol[User, uint64]: column from User model and its type is uint64.
//	Key[User]: index defined on User model.
//
// # Usage
//
//   - M / CastT builds a scoped table from a model
//   - CastC builds a model scoped column
//   - CT builds a typed model scoped column
//   - CastK builds a model scoped index
package modeled
