package frag_test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"slices"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/frag"
	. "github.com/xoctopus/sqlx/pkg/frag/testutil"
)

type Fragment = frag.Fragment

func TestFragment(t *testing.T) {
	t.Run("Const", func(t *testing.T) {
		Expect[Fragment](t, frag.Raw(""), BeFragment(""))
		Expect[Fragment](t, frag.Raw("SELECT 1"), BeFragment("SELECT 1"))
	})
	t.Run("Flatten", func(t *testing.T) {
		val := []any{1, 2, 3}
		seq := slices.Values(val)
		t.Run("Values", func(t *testing.T) {
			Expect[Fragment](t, frag.Pair(`IN (?,?,?)`, val...), BeFragment("IN (?,?,?)", val...))
		})
		t.Run("Seq", func(t *testing.T) {
			Expect[Fragment](t, frag.Pair(`IN (?)`, seq), BeFragment("IN (?,?,?)", 1, 2, 3))
		})
		t.Run("Slice", func(t *testing.T) {
			Expect[Fragment](t, frag.Pair(`IN (?)`, val), BeFragment("IN (?,?,?)", 1, 2, 3))
		})
		t.Run("Composed", func(t *testing.T) {
			Expect[Fragment](
				t,
				frag.Pair(
					`DO UPDATE SET f_name = ?`,
					[]any{frag.Pair("EXCLUDED.?", frag.Raw("f_name"))},
				),
				BeFragment("DO UPDATE SET f_name = EXCLUDED.f_name"),
			)
		})
		t.Run("HasSub", func(t *testing.T) {
			Expect[Fragment](
				t,
				frag.Pair(`#ID = ?`, frag.Pair("#ID+?", 1)),
				BeFragment("#ID = #ID+?", 1),
			)
		})
		t.Run("CustomValueArg", func(t *testing.T) {
			Expect[Fragment](
				t,
				frag.Pair(`#Point = ?`, Point{1, 1}),
				BeFragment("#Point = ST_GeomFromText(?)", Point{1, 1}),
			)
		})
		t.Run("NamedArg", func(t *testing.T) {
			t.Run("WithNamedArg", func(t *testing.T) {
				Expect[Fragment](
					t,
					frag.Pair(
						`time > @min AND time < @max`,
						sql.Named("min", 10),
						sql.Named("max", 20),
					),
					BeFragment("time > ? AND time < ?", 10, 20),
				)
			})
			t.Run("WithNamedArgSet", func(t *testing.T) {
				Expect[Fragment](
					t,
					frag.Pair(
						`time > @min AND time < @max`,
						frag.NamedArgs{
							"min": 10,
							"max": 20,
						},
					),
					BeFragment("time > ? AND time < ?", 10, 20),
				)
			})
		})
		t.Run("Embedded", func(t *testing.T) {
			Expect(t,
				frag.Pair(
					`CREATE TABLE IF NOT EXISTS @table @columns`,
					frag.NamedArgs{
						"table": frag.Pair("t"),
						"columns": frag.Block(
							frag.Compose(
								", ",
								frag.Pair(
									"@column @datatype",
									frag.NamedArgs{
										"column":   frag.Pair("f_id"),
										"datatype": frag.Pair("bigint"),
									},
								),
								nil,
								frag.Pair(
									"@column @datatype",
									frag.NamedArgs{
										"column":   frag.Pair("f_name"),
										"datatype": frag.Pair("varchar(255)"),
									},
								),
							),
						),
					}),
				BeFragment(`CREATE TABLE IF NOT EXISTS t (f_id bigint, f_name varchar(255))`),
			)
		})
	})

}

type Point struct{ X, Y float64 }

func (Point) ValueEx() string {
	return `ST_GeomFromText(?)`
}

func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%v %v)", p.X, p.Y), nil
}
