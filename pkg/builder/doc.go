// Package builder builds typed SQL fragments for tables, statements, and clauses.
//
// Fragments implement [frag.Fragment] and can be collected into query text and args.
// See Example* tests in this package for end-to-end renderings.
//
// # Definitions
//
//   - [Col] / [C] / [CT] / [CC]: a column. C creates by name; CT/CC add a Go type
//     for typed conditions and assignments. ColumnsOf / Columns / ColsOf group columns.
//   - [Key] / [PK] / [UK] / [K]: primary key and indexes on a table.
//   - [Table] / [T]: a table with columns and keys. T builds from
//     definitions; TFrom scans a [Model] (struct with db tags and optional
//     PrimaryKey / Indexes / UniqueIndexes methods) into a Table.
//   - [TFrom]: scan a go type (struct pointer) to Table
//   - [Catalog] / [NewCatalog]: a named registry of tables (Add / T / Tables).
//   - [Model]: table-backed struct type used with TFrom and generated models.
//
// Example table setup (from ExampleJoin):
//
//	tUser := T(
//		"t_user",
//		C("f_id", WithColDefOf(uint64(0), `,autoinc`)),
//		C("f_name", WithColDefOf("", `,width=128,default=''`)),
//		C("f_org_id", WithColDefOf(uint64(0), ``)),
//	)
//
// # Statements
//
//   - [Select]: SELECT ... FROM table with additions
//     (modifiers like SQL_CALC_FOUND_ROWS, DISTINCT ON, FOR UPDATE / FOR SHARE).
//
//   - [Insert]: INSERT [IGNORE] INTO ... VALUES / SELECT, with OnConflict / OnDuplicate.
//
//   - [Update] / [UpdateIgnore]: UPDATE ... SET ... [FROM|JOIN] WHERE, optional Returning.
//
//   - [Delete]: DELETE FROM ... WHERE.
//
//   - [With] / [WithRecursive]: CTE / WITH statements.
//
//     Select(nil).From(tUser, Where(CT[int]("F_id").AsCond(Eq(1))), Limit(10))
//     // SELECT * FROM t_user WHERE f_id = ? LIMIT 10
//
//     Insert().Into(T("t_user")).Values(Columns("f_a", "f_b"), 1, 2)
//     // INSERT INTO t_user (f_a,f_b) VALUES (?,?)
//
//     Delete().From(T("t_x"), Where(CT[int]("F_a").AsCond(Eq(1))))
//     // DELETE FROM t_x WHERE f_a = ?
//
//     Update(t).Set(...).Where(..., Returning(nil))
//     // UPDATE ... SET ... WHERE ... RETURNING *
//
// # Additions (clauses)
//
// Attach to From / Into / Where via [Additions]. Types sort automatically when rendered.
//
//   - [Comment]: leading SQL comments (-- ...)
//
//   - [Where]: WHERE condition
//
//   - [Join] / [LeftJoin] / [RightJoin] / [FullJoin] / [InnerJoin] / [CrossJoin]:
//     JOIN ... ON / USING
//
//   - [GroupBy]: GROUP BY ... [HAVING ...]
//
//   - [OrderBy] / [Order] / [AscOrder] / [DescOrder]: ORDER BY (NullsFirst / NullsLast)
//
//   - [Limit]: LIMIT n [.Offset(m)]; negative count omits LIMIT
//
//   - [OnConflict]: INSERT ... ON CONFLICT (...) DO NOTHING / DO UPDATE SET
//
//   - [OnDuplicate]: INSERT ... ON DUPLICATE KEY UPDATE (MySQL)
//
//   - [Returning]: RETURNING ...
//
//   - [ForUpdate] / [ForUpdateSkipLocked] / [ForUpdateNoWait]: FOR UPDATE
//
//   - [ForShare] / [ForShareSkipLocked] / [ForShareNoWait]: FOR SHARE
//
// e.g.
//
//	Select(nil).
//	From(
//		T("t_x"),
//		Where(CT[int]("F_a").AsCond(Eq(1))),
//		GroupBy(C("F_a")).Having(CT[int]("F_a").AsCond(Eq(1))),
//		Limit(20).Offset(100),
//		Comment("limit with offset"),
//	)
//
//	Insert().Into(
//		T("t_user"),
//		OnConflict(ColsOf(C("f_id"))).DoNothing(),
//	).Values(ColsOf(C("f_id"), C("f_name")), 1, "saito")
//
// # Conditions
//
// Build predicates with typed column helpers, then compose:
//
//   - Col.AsCond with Eq / Neq / Lt / Gt / Gte / In / Like / LLike / RLike / EqCol / ...
//
//   - [And] / [Or] / [Xor]: compose conditions (nil entries skipped)
//
//   - [AsCond]: wrap any fragment as a condition
//
//   - [Exists]: EXISTS (subquery)
//
// e.g.:
//
//	Xor(
//		Or(
//			And(
//				CT[int]("a").AsCond(Lt(1)),
//				CT[string]("b").AsCond(LLike[string]("text")),
//			),
//			CT[int]("a").AsCond(Eq(2)),
//		),
//		CT[string]("b").AsCond(RLike[string]("g")),
//	)
//	// (((a < ?) AND (b LIKE ?)) OR (a = ?)) XOR (b LIKE ?)
//
// # Functions
//
// Aggregate and misc SQL functions via [Func] helpers:
//
//   - [Count]: COUNT(1) when empty, else COUNT(args); Func("COUNT") → COUNT(*)
//
//   - [Avg] / [Min] / [Max] / [Sum]
//
//   - [AnyValue] / [First] / [Last] / [Distinct]
//
//     Count()       // COUNT(1)
//     Count(C("a")) // COUNT(a)
//     Avg(C("a"))   // AVG(a)
//
// # Assignments
//
// [ColumnsAndValues] and typed AssignBy / Value / AsValue / Inc build SET and
// INSERT value lists (see ExampleAssignment, ExampleUpdate, ExampleInsert).
package builder
