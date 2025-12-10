package frag_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"slices"
	"testing"

	"github.com/shopspring/decimal"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	. "github.com/xoctopus/sqlx/pkg/frag/testutil"
)

type Fragment = frag.Fragment

func TestFragment(t *testing.T) {
	t.Run("CollectEmpty", func(t *testing.T) {
		q, args := frag.Collect(nil, nil)
		Expect(t, q, Equal(""))
		Expect(t, args, HaveLen[[]any](0))
	})
	t.Run("Const", func(t *testing.T) {
		Expect(t, frag.Lit(""), BeFragment(""))
		Expect(t, frag.Lit("SELECT 1"), BeFragment("SELECT 1"))
	})
	t.Run("Flatten", func(t *testing.T) {
		val := []any{1, 2, 3}
		seq := slices.Values(val)
		t.Run("Values", func(t *testing.T) {
			Expect(t, frag.Query(`IN (?,?,?)`, val...), BeFragment("IN (?,?,?)", val...))
		})
		t.Run("Valuer", func(t *testing.T) {
			Expect(t, frag.Query(`f_balance = ?`, decimal.NewFromFloat(100.001)), BeFragment("f_balance = ?", decimal.NewFromFloat(100.001)))
		})
		t.Run("Seq", func(t *testing.T) {
			Expect(t, frag.Query(`IN (?)`, seq), BeFragment("IN (?,?,?)", 1, 2, 3))
			Expect(t, frag.Query(`IN (?)`, slices.Values([]int{1, 2, 3})), BeFragment("IN (?,?,?)", 1, 2, 3))
			sub := slices.Values([]any{frag.Query("SELECT 1,2")})
			Expect(t, frag.Query(`IN (?)`, sub), BeFragment("IN (SELECT 1,2)"))
		})
		t.Run("Slice", func(t *testing.T) {
			type String string

			Expect(t, frag.Query(`IN (?)`, val), BeFragment("IN (?,?,?)", 1, 2, 3))
			Expect(t, frag.Query(`IN (?)`, []any{0, 1}), BeFragment("IN (?,?)", 0, 1))
			Expect(t, frag.Query(`IN (?)`, []bool{true, false}), BeFragment("IN (?,?)", true, false))
			Expect(t, frag.Query(`IN (?)`, []string{"a", "b"}), BeFragment("IN (?,?)", "a", "b"))
			Expect(t, frag.Query(`IN (?)`, []float32{0.1, 0.2}), BeFragment("IN (?,?)", float32(0.1), float32(0.2)))
			Expect(t, frag.Query(`IN (?)`, []float64{0.1, 0.2}), BeFragment("IN (?,?)", 0.1, 0.2))
			Expect(t, frag.Query(`IN (?)`, []int{0, 1}), BeFragment("IN (?,?)", 0, 1))
			Expect(t, frag.Query(`IN (?)`, []int8{3, 4}), BeFragment("IN (?,?)", int8(3), int8(4)))
			Expect(t, frag.Query(`IN (?)`, []int16{0, 1}), BeFragment("IN (?,?)", int16(0), int16(1)))
			Expect(t, frag.Query(`IN (?)`, []int32{0, 1}), BeFragment("IN (?,?)", int32(0), int32(1)))
			Expect(t, frag.Query(`IN (?)`, []int64{0, 1}), BeFragment("IN (?,?)", int64(0), int64(1)))
			Expect(t, frag.Query(`IN (?)`, []uint{5, 6}), BeFragment("IN (?,?)", uint(5), uint(6)))
			Expect(t, frag.Query(`IN (?)`, []uint16{7, 8}), BeFragment("IN (?,?)", uint16(7), uint16(8)))
			Expect(t, frag.Query(`IN (?)`, []uint32{0, 1}), BeFragment("IN (?,?)", uint32(0), uint32(1)))
			Expect(t, frag.Query(`IN (?)`, []uint64{0, 1}), BeFragment("IN (?,?)", uint64(0), uint64(1)))
			Expect(t, frag.Query(`IN (?)`, []String{"123", "456"}), BeFragment("IN (?,?)", String("123"), String("456")))

			t.Run("Bytes", func(t *testing.T) {
				Expect(t, frag.Query(`IN (?)`, []byte("123")), BeFragment("IN (?)", []byte("123")))
			})
		})
		t.Run("Composed", func(t *testing.T) {
			Expect(
				t,
				frag.Query(
					`DO UPDATE SET f_name = ?`,
					[]any{frag.Query("EXCLUDED.?", frag.Lit("f_name"))},
				),
				BeFragment("DO UPDATE SET f_name = EXCLUDED.f_name"),
			)
		})
		t.Run("HasSub", func(t *testing.T) {
			Expect(
				t,
				frag.Query(`#ID = ?`, frag.Query("#ID+?", 1)),
				BeFragment("#ID = #ID+?", 1),
			)
		})
		t.Run("CustomValueArg", func(t *testing.T) {
			Expect(
				t,
				frag.Query(`#Point = ?`, Point{1, 1}),
				BeFragment("#Point = ST_GeomFromText(?)", Point{1, 1}),
			)
		})
		t.Run("NamedArg", func(t *testing.T) {
			t.Run("WithNamedArg", func(t *testing.T) {
				Expect(
					t,
					frag.Query(
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
					frag.Query(
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
				frag.Query(
					`CREATE TABLE IF NOT EXISTS @table @columns`,
					frag.NamedArgs{
						"table": frag.Query("t"),
						"columns": frag.Block(
							frag.Compose(
								", ",
								frag.Query(
									"@column @datatype",
									frag.NamedArgs{
										"column":   frag.Query("f_id"),
										"datatype": frag.Query("bigint"),
									},
								),
								nil,
								frag.Query(
									"@column @datatype",
									frag.NamedArgs{
										"column":   frag.Query("f_name"),
										"datatype": frag.Query("varchar(255)"),
									},
								),
							),
						),
					}),
				BeFragment(`CREATE TABLE IF NOT EXISTS t (f_id bigint, f_name varchar(255))`),
			)
		})
	})
	t.Run("EmptyArgs", func(t *testing.T) {
		Expect(t, frag.Query("f_a = f_b?", frag.Empty()), BeFragment("f_a = f_b"))
		Expect(t, frag.Query("f_a = f_b?", []any{}), BeFragment("f_a = f_b"))
		Expect(t, frag.Arg(nil).IsNil(), BeFalse())
		Expect(t, frag.Values[int](slices.Values([]int{1})).IsNil(), BeFalse())
		ExpectPanic[error](t, func() { frag.ArgIter(context.Background(), func() {}) })
	})
}

type Point struct{ X, Y float64 }

func (Point) ValueEx() string {
	return `ST_GeomFromText(?)`
}

func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%v %v)", p.X, p.Y), nil
}

func TestBlock(t *testing.T) {
	Expect(t, frag.Block(frag.Lit("?,?")), BeFragment("(?,?)"))
	Expect(t, frag.Block(nil), BeFragment(""))
	Expect(t, frag.BlockWithoutBrackets(
		slices.Values([]Fragment{
			builder.C("f_id"),
			builder.C("f_name"),
		}),
	), BeFragment("f_id,f_name"))
	Expect(t, frag.BlockWithoutBrackets(nil), BeFragment(""))
}

func TestCompose(t *testing.T) {
	Expect(
		t,
		frag.Compose(",", builder.C("f_id"), builder.C("f_name")),
		BeFragment("f_id,f_name"),
	)

	Expect(
		t,
		frag.ComposeSeq(",", slices.Values([]Fragment{builder.C("f_id"), builder.C("f_name"), nil})),
		BeFragment("f_id,f_name"),
	)
}

func TestFunc(t *testing.T) {
	f := frag.ArgIterFunc(100)
	Expect(t, f, BeFragment("?", 100))
}
