package types_test

import (
	"database/sql/driver"
	"testing"

	. "github.com/xoctopus/x/testx"

	. "github.com/xoctopus/sqlx/pkg/types"
)

func TestArrayAsList(t *testing.T) {
	type ID uint64

	t.Run("Parse", func(t *testing.T) {
		t.Run("Int", func(t *testing.T) {
			list, err := ParseArrayAsList[int]("1,2,3")
			Expect(t, err, Succeed())
			Expect(t, list, Equal(ArrayAsList[int]{1, 2, 3}))

			list, err = ParseArrayAsList[int](" 1 , 2 , 3 ")
			Expect(t, err, Succeed())
			Expect(t, list, Equal(ArrayAsList[int]{1, 2, 3}))

			_, err = ParseArrayAsList[int]("")
			Expect(t, err, Failed())
			_, err = ParseArrayAsList[int]("1,,2")
			Expect(t, err, Failed())
			_, err = ParseArrayAsList[int]("1,2,")
			Expect(t, err, Failed())
			_, err = ParseArrayAsList[int]("1,a,3")
			Expect(t, err, Failed())
		})

		t.Run("Int8Overflow", func(t *testing.T) {
			_, err := ParseArrayAsList[int8]("256")
			Expect(t, err, Failed())
		})

		t.Run("Uint", func(t *testing.T) {
			list, err := ParseArrayAsList[ID]("1,2,3")
			Expect(t, err, Succeed())
			Expect(t, list, Equal(ArrayAsList[ID]{1, 2, 3}))

			_, err = ParseArrayAsList[ID]("xxx")
			Expect(t, err, Failed())
		})

		t.Run("Float", func(t *testing.T) {
			list, err := ParseArrayAsList[float64]("1.5,2.25")
			Expect(t, err, Succeed())
			Expect(t, list, Equal(ArrayAsList[float64]{1.5, 2.25}))

			_, err = ParseArrayAsList[float64]("xxx")
			Expect(t, err, Failed())
		})

		t.Run("String", func(t *testing.T) {
			list, err := ParseArrayAsList[string]("a,b,c")
			Expect(t, err, Succeed())
			Expect(t, list, Equal(ArrayAsList[string]{"a", "b", "c"}))

			list, err = ParseArrayAsList[string](" a , b ")
			Expect(t, err, Succeed())
			Expect(t, list, Equal(ArrayAsList[string]{"a", "b"}))

			_, err = ParseArrayAsList[string]("")
			Expect(t, err, Failed())
			_, err = ParseArrayAsList[string]("a,,b")
			Expect(t, err, Failed())
		})
	})

	t.Run("DBType", func(t *testing.T) {
		Expect(t, ArrayAsList[int]{}.DBType("mysql"), Equal("text"))
		Expect(t, ArrayAsList[string]{}.DBType("postgres"), Equal("text"))
	})

	t.Run("AppendElementsStringValue", func(t *testing.T) {
		aa := ArrayAsList[int]{}
		aa.Append(1, 2)
		aa.Append(3)

		Expect(t, aa.Elements(), Equal([]string{"1", "2", "3"}))
		Expect(t, aa.String(), Equal("1,2,3"))

		v, err := aa.Value()
		Expect(t, err, Succeed())
		Expect(t, v, Equal[driver.Value]("1,2,3"))

		empty := ArrayAsList[string]{}
		Expect(t, empty.String(), Equal(""))
		v, err = empty.Value()
		Expect(t, err, Succeed())
		Expect(t, v, Equal[driver.Value](""))
	})

	t.Run("Scan", func(t *testing.T) {
		aa := ArrayAsList[int]{9}

		Expect(t, aa.Scan("1,2,3"), Succeed())
		Expect(t, aa, Equal(ArrayAsList[int]{1, 2, 3}))

		Expect(t, aa.Scan([]byte("4, 5")), Succeed())
		Expect(t, aa, Equal(ArrayAsList[int]{4, 5}))

		Expect(t, aa.Scan(nil), Succeed())
		Expect(t, aa, Equal(ArrayAsList[int]{}))

		Expect(t, aa.Scan("1,,2"), Failed())
		Expect(t, aa.Scan(float64(1)), Failed())
	})

	t.Run("Text", func(t *testing.T) {
		aa := ArrayAsList[string]{"a", "b"}

		data, err := aa.MarshalText()
		Expect(t, err, Succeed())
		Expect(t, data, Equal([]byte("a,b")))

		var decoded ArrayAsList[string]
		Expect(t, decoded.UnmarshalText(data), Succeed())
		Expect(t, decoded, Equal(aa))

		Expect(t, decoded.UnmarshalText([]byte("")), Failed())
		Expect(t, decoded.UnmarshalText([]byte("a,,b")), Failed())
	})

	t.Run("JSON", func(t *testing.T) {
		aa := ArrayAsList[ID]{1, 2, 3}

		data, err := aa.MarshalJSON()
		Expect(t, err, Succeed())
		Expect(t, data, Equal([]byte(`"1,2,3"`)))

		var decoded ArrayAsList[ID]
		Expect(t, decoded.UnmarshalJSON(data), Succeed())
		Expect(t, decoded, Equal(aa))

		Expect(t, decoded.UnmarshalJSON([]byte("null")), Succeed())
		Expect(t, decoded, Equal(ArrayAsList[ID]{}))

		Expect(t, decoded.UnmarshalJSON([]byte(nil)), Succeed())
		Expect(t, decoded, Equal(ArrayAsList[ID]{}))

		Expect(t, decoded.UnmarshalJSON([]byte(`"1,2,3"`)), Succeed())
		Expect(t, decoded, Equal(aa))

		Expect(t, decoded.UnmarshalJSON([]byte(`[1,2]`)), Failed())
		Expect(t, decoded.UnmarshalJSON([]byte(`"1,,2"`)), Failed())
	})

	t.Run("RoundTrip", func(t *testing.T) {
		src := ArrayAsList[int]{7, 8, 9}
		v, err := src.Value()
		Expect(t, err, Succeed())

		var dst ArrayAsList[int]
		Expect(t, dst.Scan(v), Succeed())
		Expect(t, dst, Equal(src))

		text, err := src.MarshalText()
		Expect(t, err, Succeed())
		Expect(t, dst.UnmarshalText(text), Succeed())
		Expect(t, dst, Equal(src))

		js, err := src.MarshalJSON()
		Expect(t, err, Succeed())
		Expect(t, dst.UnmarshalJSON(js), Succeed())
		Expect(t, dst, Equal(src))
	})
}
