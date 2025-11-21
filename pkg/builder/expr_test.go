package builder_test

import (
	"context"
	"fmt"

	. "github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

var (
	tUser = T(
		"t_user",
		C("f_id", WithColDefOf(uint64(0), `,autoinc`)),
		C("f_name", WithColDefOf("", `,width=128,default=''`)),
		C("f_org_id", WithColDefOf(uint64(0), ``)),
	)
	tOrg = T(
		"t_org",
		C("f_org_id", WithColDefOf(uint64(0), `,autoinc`)),
		C("f_name", WithColDefOf("", `,width=128,default=''`)),
	)
)

func Print(ctx context.Context, f frag.Fragment) {
	if f == nil {
		return
	}

	query, args := frag.Collect(ctx, f)
	fmt.Println(query)
	fmt.Println(args)
}

func PrintQuery(ctx context.Context, f frag.Fragment) {
	if f == nil {
		return
	}

	query, _ := frag.Collect(ctx, f)
	fmt.Println(query)
}

func ExampleGroupBy() {
	tab := T("t_x")

	f := Select(ColsOf(C("F_a"), C("F_b"))).
		From(
			tab,
			Where(Where(CT[int]("F_a").AsCond(Eq(1)))),
			GroupBy(
				C("F_a"),
				C("F_b"),
			).Having(CT[int]("F_a").AsCond(Eq(1))),
			Comment("group multi columns"),
		)
	Print(context.Background(), f)

	f = Select(nil).
		From(
			tab,
			Where(CT[int]("F_a").AsCond(Eq(1))),
			GroupBy(
				AscOrder(C("F_a")),
				DescOrder(C("F_b")),
			),
			Comment("group with order"),
		)
	Print(context.Background(), f)

	// Output:
	// -- group multi columns
	// SELECT f_a,f_b FROM t_x WHERE f_a = ? GROUP BY f_a,f_b HAVING f_a = ?
	// [1 1]
	// -- group with order
	// SELECT * FROM t_x WHERE f_a = ? GROUP BY (f_a) ASC,(f_b) DESC
	// [1]
}

func ExampleJoin() {
	f := Select(frag.Compose(", ",
		Alias(tUser.C("f_id"), "f_user_id"),
		Alias(tUser.C("f_org_id"), "f_org_id"),
		Alias(tOrg.C("f_name"), "f_org_name"),
	)).From(
		tUser,
		Join(Alias(tOrg, "t_org")).On(
			And(
				CTOf[int64](tUser, "f_org_id").AsCond(EqCol(CTOf[int64](tOrg, "f_org_id"))),
				CTOf[string](tOrg, "f_name").AsCond(Neq("abc")),
			),
		),
		Comment("join on"),
	)
	Print(context.Background(), f)

	f = Select(nil).
		From(
			tUser,
			LeftJoin(tOrg).Using(tUser.C("f_org_id")),
			Comment("join using"),
		)
	Print(context.Background(), f)

	// TODO long sql query should break to lines
	/*
		SELECT
			t_user.f_id AS f_user_id,
			t_user.f_name AS f_user_name,
			t_user.f_org_id AS f_org_id,
			t_org.f_name AS f_org_name
		FROM
			t_user
		JOIN
			t_org AS t_org
		ON
			(t_user.f_org_id = t_org.f_org_id)
			AND
			(t_org.f_name <> ?)
	*/
	f = Select(
		AutoAlias(
			tUser.C("f_id"),
			tOrg.C("f_name"),
		),
	).From(
		tUser,
		FullJoin(tOrg).On(
			CTOf[int](tUser, "f_org_id").AsCond(EqCol(CTOf[int](tOrg, "f_org_id"))),
		),
		Comment("full join + auto alias"),
	)
	Print(context.Background(), f)

	// Output:
	// -- join on
	// SELECT t_user.f_id AS f_user_id, t_user.f_org_id AS f_org_id, t_org.f_name AS f_org_name FROM t_user JOIN t_org AS t_org ON (t_user.f_org_id = t_org.f_org_id) AND (t_org.f_name <> ?)
	// [abc]
	// -- join using
	// SELECT * FROM t_user LEFT JOIN t_org USING (f_org_id)
	// []
	// -- full join + auto alias
	// SELECT t_user.f_id AS t_user__f_id, t_org.f_name AS t_org__f_name FROM t_user FULL JOIN t_org ON t_user.f_org_id = t_org.f_org_id
	// []
}

func ExampleLimit() {
	f := Select(nil).
		From(
			T("t_x"),
			Where(CT[int]("F_a").AsCond(Eq(1))),
			Limit(1),
			Comment("limit"),
		)
	Print(context.Background(), f)

	f = Select(nil).
		From(
			T("t_x"),
			Where(CT[int]("F_a").AsCond(Eq(1))),
			Limit(-1),
			Comment("without limit"),
		)
	Print(context.Background(), f)

	f = Select(nil).
		From(
			T("t_x"),
			Where(CT[int]("F_a").AsCond(Eq(1))),
			Limit(20).Offset(100),
			Comment("limit with offset"),
		)
	Print(context.Background(), f)

	// Output:
	// -- limit
	// SELECT * FROM t_x WHERE f_a = ? LIMIT 1
	// [1]
	// -- without limit
	// SELECT * FROM t_x WHERE f_a = ?
	// [1]
	// -- limit with offset
	// SELECT * FROM t_x WHERE f_a = ? LIMIT 20 OFFSET 100
	// [1]
}

func ExampleOrderBy() {
	f := Select(nil).
		From(
			T("t_x"),
			OrderBy(
				AscOrder(C("F_a")),
				DescOrder(C("F_b")),
			),
			Where(
				And(
					CT[int]("F_a").AsCond(Eq(1)),
					CT[int]("F_b").Fragment("# = ?+1", CT[int]("F_a")),
				),
			),
		)
	Print(context.Background(), f)

	// Output:
	// SELECT * FROM t_x WHERE (f_a = ?) AND (f_b = f_a+1) ORDER BY (f_a) ASC,(f_b) DESC
	// [1]
}

func ExampleCondition() {
	f := Xor(
		Or(
			And(
				nil,
				CT[int]("a").AsCond(Lt(1)),
				CT[string]("b").AsCond(LLike[string]("text")),
			),
			CT[int]("a").AsCond(Eq(2)),
		),
		CT[string]("b").AsCond(RLike[string]("g")),
	)
	Print(context.Background(), f)

	f = Xor(
		Or(
			And(
				(*Condition)(nil),
				(*Condition)(nil),
				(*Condition)(nil),
				(*Condition)(nil),
				CT[int]("c").AsCond(In(1, 2)),
				CT[int]("c").AsCond(In(3, 4)),
				CT[int]("a").AsCond(Eq(1)),
				CT[string]("b").AsCond(Like("text")),
			),
			CT[int]("a").AsCond(Eq(2)),
		),
		CT[string]("b").AsCond(Like("g")),
	)
	Print(context.Background(), f)

	// Output:
	// (((a < ?) AND (b LIKE ?)) OR (a = ?)) XOR (b LIKE ?)
	// [1 %text 2 g%]
	// (((c IN (?,?)) AND (c IN (?,?)) AND (a = ?) AND (b LIKE ?)) OR (a = ?)) XOR (b LIKE ?)
	// [1 2 3 4 1 %text% 2 %g%]
}

func ExampleFunction() {
	PrintQuery(context.Background(), Count())
	PrintQuery(context.Background(), Count(C("a")))
	PrintQuery(context.Background(), Avg(C("a")))
	PrintQuery(context.Background(), AnyValue(C("a")))
	PrintQuery(context.Background(), Min(C("a")))
	PrintQuery(context.Background(), Max(C("a")))
	PrintQuery(context.Background(), First(C("a")))
	PrintQuery(context.Background(), Last(C("a")))
	PrintQuery(context.Background(), Sum(C("a")))
	PrintQuery(context.Background(), Func(""))
	PrintQuery(context.Background(), Func("COUNT"))

	// Output:
	// COUNT(1)
	// COUNT(a)
	// AVG(a)
	// ANY_VALUE(a)
	// MIN(a)
	// MAX(a)
	// FIRST(a)
	// LAST(a)
	// SUM(a)
	//
	// COUNT(*)
}

func ExampleSelect() {
	f := Select(
		nil,
		frag.Lit("SQL_CALC_FOUND_ROWS"),
	).From(
		tUser,
		Where(CT[int]("F_created_at").AsCond(Gte(593650800))),
		Limit(10),
		Comment("select with modifier"),
	)
	Print(context.Background(), f)

	f = Select(nil).From(
		tUser,
		Where(CT[int]("F_id").AsCond(Eq(1))),
		ForUpdate(),
		Comment("select 1 row for update"),
	)
	Print(context.Background(), f)

	// Output:
	// -- select with modifier
	// SELECT SQL_CALC_FOUND_ROWS * FROM t_user WHERE f_created_at >= ? LIMIT 10
	// [593650800]
	// -- select 1 row for update
	// SELECT * FROM t_user WHERE f_id = ? FOR UPDATE
	// [1]
}

func ExampleDelete() {
	f := Delete().From(
		T("t_x"),
		Where(CT[int]("F_a").AsCond(Eq(1))),
		Comment("delete"),
	)
	Print(context.Background(), f)

	// Output:
	// -- delete
	// DELETE FROM t_x WHERE f_a = ?
	// [1]
}

func ExampleInsert() {
	f := Insert().
		Into(
			T("t_user"),
			Comment("insert"),
		).
		Values(Columns("f_a", "f_b"), 1, 2)
	Print(context.Background(), f)

	f = Insert("IGNORE").
		Into(
			T("t_user"),
			Comment("insert ignore and multi"),
		).
		Values(Columns("f_a", "f_b"), 1, 2, 1, 2, 1, 2)
	Print(context.Background(), f)

	f = Insert().
		Into(
			T("t_user_migrated"),
			Comment("insert from selection"),
		).
		Values(
			Columns("f_a", "f_b"),
			Select(Columns("f_a", "f_b")).
				From(
					T("t_user_previous"),
					Where(CT[string]("f_status").AsCond(Eq("valid"))),
				),
		)

	Print(context.Background(), f)

	// Output:
	// -- insert
	// INSERT INTO t_user (f_a,f_b) VALUES (?,?)
	// [1 2]
	// -- insert ignore and multi
	// INSERT IGNORE INTO t_user (f_a,f_b) VALUES (?,?),(?,?),(?,?)
	// [1 2 1 2 1 2]
	// -- insert from selection
	// INSERT INTO t_user_migrated (f_a,f_b) SELECT f_a,f_b FROM t_user_previous WHERE f_status = ?
	// [valid]
}

func ExampleUpdate() {
	var (
		f_id     = CT[int]("f_id")
		f_stu_id = CT[int]("f_stu_id")
		f_score  = CT[int]("f_score")
		f_class  = CT[string]("f_class")

		t_stu   = T("t_stu", f_id, f_score, f_class)
		t_class = T("t_class", f_stu_id, f_class)
	)

	f := Update(t_stu).
		Set(
			CTOf[int](t_stu, "f_score").AssignBy(Value(100)),
			CTOf[string](t_stu, f_class.Name()).AssignBy(AsValue(CTOf[string](t_class, f_class.Name()))),
		).
		From(t_class).
		Where(
			CastC[string](t_stu.C(f_id.Name())).
				AsCond(EqCol(CastC[string](t_class.C(f_stu_id.Name())))),
			Comment("update from postgres supported"),
		)
	Print(context.Background(), f)

	f = Update(t_stu).
		From(
			nil,
			Join(t_class).On(
				CastC[string](t_stu.C(f_id.Name())).
					AsCond(EqCol(CastC[string](t_class.C(f_stu_id.Name())))),
			),
			Comment("update join mysql supported"),
		).
		Set(
			CTOf[int](t_stu, "f_score").AssignBy(Value(100)),
			CTOf[string](t_stu, f_class.Name()).AssignBy(AsValue(CTOf[string](t_class, f_class.Name()))),
		)
	Print(context.Background(), f)

	sub := Select(f_class.Of(t_class)).From(
		t_class,
		Where(CastC[int](f_stu_id.Of(t_class)).AsCond(EqCol(CastC[int](f_id.Of(t_stu))))),
		Limit(1),
	)
	f = Update(t_stu).
		Set(
			CTOf[int](t_stu, "f_score").AssignBy(Value(100)),
			ColumnsAndValues(
				f_class.Of(t_stu),
				sub,
			),
		).
		Where(
			Exists(sub),
			Comment("update with sub query and exists condition sqlite supported"),
		)
	Print(context.Background(), f)

	// Output:
	// -- update from postgres supported
	// UPDATE t_stu SET f_score = ?, f_class = t_class.f_class FROM t_class WHERE t_stu.f_id = t_class.f_stu_id
	// [100]
	// -- update join mysql supported
	// UPDATE t_stu JOIN t_class ON t_stu.f_id = t_class.f_stu_id SET f_score = ?, f_class = t_class.f_class
	// [100]
	// -- update with sub query and exists condition sqlite supported
	// UPDATE t_stu SET f_score = ?, f_class = (SELECT f_class FROM t_class WHERE f_stu_id = f_id LIMIT 1) WHERE EXISTS (SELECT f_class FROM t_class WHERE f_stu_id = f_id LIMIT 1)
	// [100]
}
