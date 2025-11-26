package builder_test

import (
	"context"
	"fmt"

	. "github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	. "github.com/xoctopus/sqlx/testdata"
)

var (
	tUser = T(
		"t_user",
		C("f_id", WithColDefOf(context.Background(), uint64(0), `,autoinc`)),
		C("f_name", WithColDefOf(context.Background(), "", `,width=128,default=''`)),
		C("f_org_id", WithColDefOf(context.Background(), uint64(0), ``)),
	)
	tOrg = T(
		"t_org",
		C("f_org_id", WithColDefOf(context.Background(), uint64(0), `,autoinc`)),
		C("f_name", WithColDefOf(context.Background(), "", `,width=128,default=''`)),
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
	f := Select(ColsOf(C("F_a"), C("F_b"))).
		From(
			T("t_x"),
			Where(Where(CT[int]("F_a").AsCond(Eq(1)))),
			GroupBy(C("F_a"), C("F_b")).
				Having(CT[int]("F_a").AsCond(Eq(1))),
			Comment("group multi columns"),
		)
	Print(context.Background(), f)

	// Output:
	// -- group multi columns
	// SELECT f_a,f_b FROM t_x WHERE f_a = ? GROUP BY f_a,f_b HAVING f_a = ?
	// [1 1]
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
				CastC[int64](tUser.C("f_org_id")).AsCond(EqCol(CastC[int64](tOrg.C("f_org_id")))),
				CastC[string](tOrg.C("f_name")).AsCond(Neq("abc")),
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

	f = Select(
		frag.Compose(
			",",
			Alias(C("f_id").Of(tUser), "f_user_id"),
			Alias(C("f_org_id").Of(tOrg), "f_org_id"),
		),
	).From(
		tUser,
		RightJoin(tOrg).On(
			CastC[int64](tUser.C("f_org_id")).AsCond(EqCol(CastC[int64](tOrg.C("f_org_id")))),
		),
		Comment("right join"),
	)
	Print(context.Background(), f)

	f = Select(
		frag.Compose(
			",",
			Alias(C("f_id").Of(tUser), "f_user_id"),
			Alias(C("f_org_id").Of(tOrg), "f_org_id"),
		),
	).From(
		tUser,
		FullJoin(tOrg).On(
			CastC[int64](tUser.C("f_org_id")).AsCond(EqCol(CastC[int64](tOrg.C("f_org_id")))),
		),
		Comment("full join"),
	)
	Print(context.Background(), f)

	f = Select(
		frag.Compose(
			",",
			Alias(C("f_id").Of(tUser), "f_user_id"),
			Alias(C("f_org_id").Of(tOrg), "f_org_id"),
		),
	).From(
		tUser,
		InnerJoin(tOrg).On(
			CastC[int64](tUser.C("f_org_id")).AsCond(EqCol(CastC[int64](tOrg.C("f_org_id")))),
		),
		Comment("inner join"),
	)
	Print(context.Background(), f)

	f = Select(nil).
		From(tUser, CrossJoin(tOrg), Comment("cross join"))
	Print(context.Background(), f)

	f = Select(
		AutoAlias(
			tUser.C("f_id"),
			tOrg.C("f_name"),
		),
	).From(
		tUser,
		FullJoin(tOrg).On(
			CastC[int](tUser.C("f_org_id")).AsCond(EqCol(CastC[int](tOrg.C("f_org_id")))),
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
	// -- right join
	// SELECT t_user.f_id AS f_user_id,t_org.f_org_id AS f_org_id FROM t_user RIGHT JOIN t_org ON t_user.f_org_id = t_org.f_org_id
	// []
	// -- full join
	// SELECT t_user.f_id AS f_user_id,t_org.f_org_id AS f_org_id FROM t_user FULL JOIN t_org ON t_user.f_org_id = t_org.f_org_id
	// []
	// -- inner join
	// SELECT t_user.f_id AS f_user_id,t_org.f_org_id AS f_org_id FROM t_user INNER JOIN t_org ON t_user.f_org_id = t_org.f_org_id
	// []
	// -- cross join
	// SELECT * FROM t_user CROSS JOIN t_org
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
			Comment("limit with offset", "comment line2"),
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
	// -- comment line2
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
		Where(CT[int]("F_id").AsCond(In(1, 2, 3))),
		OrderBy(
			Order(tUser.C("f_id")),
			AscOrder(tUser.C("f_org_id"), NullsFirst()),
			DescOrder(tUser.C("f_name"), NullsLast()),
		),
		Comment("select with orders"),
	)
	Print(context.Background(), f)

	f = Select(
		frag.Compose(",",
			DistinctOn(tUser.C("f_org_id")),
			tUser.C("f_id"),
			tUser.C("f_name"),
		),
	).From(
		tUser,
		OrderBy(Order(C("f_created_at"))),
		Comment("select with distinct on"),
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
	// -- select with orders
	// SELECT * FROM t_user WHERE f_id IN (?,?,?) ORDER BY (f_id),(f_org_id) ASC NULLS FIRST,(f_name) DESC NULLS LAST
	// [1 2 3]
	// -- select with distinct on
	// SELECT DISTINCT ON (f_org_id),f_id,f_name FROM t_user ORDER BY (f_created_at)
	// []
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

	f = Insert().
		Into(
			T("t_user"),
			OnConflict(ColsOf(C("f_id"))).DoNothing(),
			Comment("insert on conflict do nothing"),
		).
		Values(ColsOf(C("f_id"), C("f_name")), 1, "saito")
	Print(context.Background(), f)

	f = Insert().
		Into(
			T("t_user"),
			OnConflict(ColsOf(C("f_id"))).
				DoUpdateSet(ColumnsAndValues(ColsOf(C("f_name")), "saito")),
			Comment("insert on conflict do update"),
		).
		Values(ColsOf(C("f_id"), C("f_name")), 1, "saito")
	Print(context.Background(), f)

	f = Insert().
		Into(
			T("t_user"),
			OnDuplicate(ColumnsAndValues(ColsOf(C("f_id")), C("f_id"))),
			Comment("insert on duplicate key do update"),
		).
		Values(ColsOf(C("f_id"), C("f_name")), 1, "saito")
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
	// -- insert on conflict do nothing
	// INSERT INTO t_user (f_id,f_name) VALUES (?,?) ON CONFLICT (f_id) DO NOTHING
	// [1 saito]
	// -- insert on conflict do update
	// INSERT INTO t_user (f_id,f_name) VALUES (?,?) ON CONFLICT (f_id) DO UPDATE SET f_name = ?
	// [1 saito saito]
	// -- insert on duplicate key do update
	// INSERT INTO t_user (f_id,f_name) VALUES (?,?) ON DUPLICATE KEY UPDATE f_id = f_id
	// [1 saito]
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
			CastC[int](t_stu.C("f_score")).AssignBy(Value(100)),
			CastC[string](t_stu.C("f_class")).AssignBy(AsValue(CastC[string](t_class.C("f_class")))),
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
			CastC[int](t_stu.C("f_score")).AssignBy(Value(100)),
			CastC[string](t_stu.C("f_class")).AssignBy(AsValue(CastC[string](t_class.C("f_class")))),
		)
	Print(context.Background(), f)

	sub := Select(f_class.Of(t_class)).From(
		t_class,
		Where(CastC[int](f_stu_id.Of(t_class)).AsCond(EqCol(CastC[int](f_id.Of(t_stu))))),
		Limit(1),
	)
	f = Update(t_stu).
		Set(
			CastC[int](t_stu.C("f_score")).AssignBy(Value(100)),
			ColumnsAndValues(f_class.Of(t_stu), sub),
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

func Example_test() {
	var f frag.Fragment

	f = Select(nil).From(
		TUser,
		Where(
			And(
				TUser.UserID.AsCond(Eq[UserID](100)),
				TUser.OrgID.Fragment("# = ? + 1", TOrg.OrgID),
			),
		),
		LeftJoin(TOrg).On(
			TUser.OrgID.AsCond(EqCol[OrgID](TOrg.OrgID)),
		),
		Limit(100).Offset(200),
	)
	Print(context.Background(), f)

	f = Update(TUser).
		Set(
			TUser.Nickname.AssignBy(Value("new_name")),
		).
		Where(
			TUser.UserID.AsCond(Eq[UserID](100)),
		)
	Print(context.Background(), f)

	f = Insert().Into(TOrg).Values(
		ColsOf(TOrg.OrgID, TOrg.Name, TOrg.Belonged, TOrg.Manager),
		100, "org_name", 101, 102,
	)
	Print(context.Background(), f)

	f = Delete().
		From(
			TUser,
			Where(TUser.UserID.AsCond(Eq[UserID](100))),
		)
	Print(context.Background(), f)

	// Output:
	// SELECT * FROM users LEFT JOIN t_org ON users.f_org_id = t_org.f_org_id WHERE (users.f_user_id = ?) AND (users.f_org_id = t_org.f_org_id + 1) LIMIT 100 OFFSET 200
	// [100]
	// UPDATE users SET f_nick_name = ? WHERE f_user_id = ?
	// [new_name 100]
	// INSERT INTO t_org (f_org_id,f_name,f_belongs,manager) VALUES (?,?,?,?)
	// [100 org_name 101 102]
	// DELETE FROM users WHERE f_user_id = ?
	// [100]
}
