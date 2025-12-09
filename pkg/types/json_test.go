package types_test

import (
	"database/sql/driver"
	"errors"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/pkg/types"
)

func TestJSON(t *testing.T) {
	t.Run("JSONArray", func(t *testing.T) {
		var arr = types.JSONArrayOf[int](nil)
		t.Run("Scan", func(t *testing.T) {
			for _, src := range []any{[]byte(`null`), "", []byte(``), nil} {
				Expect(t, arr.Scan(src), Succeed())
				Expect(t, arr.IsZero(), BeTrue())
				Expect(t, arr.Get(), HaveLen[[]int](0))
			}
			Expect(t, arr.Scan(100), Failed())
			Expect(t, arr.Scan(`[1,2,3]`), Succeed())
			Expect(t, arr.Get(), Equal[[]int]([]int{1, 2, 3}))

			Expect(t, arr.UnmarshalJSON([]byte(`[`)), Failed())
		})
		t.Run("Value", func(t *testing.T) {
			arr.Set(nil)
			v, _ := arr.Value()
			Expect(t, v, Equal[driver.Value](""))
			arr.Set([]int{1})
			v, _ = arr.Value()
			Expect(t, v, Equal[driver.Value]("[1]"))
		})
		t.Run("DatabaseType", func(t *testing.T) {
			Expect(t, arr.DBType(), Equal("TEXT"))
			arr.WithDBType("JSON")
			Expect(t, arr.DBType(), Equal("JSON"))
		})
	})

	t.Run("JSONObject", func(t *testing.T) {
		type T struct {
			A int `json:"a"`
		}
		type M map[int]int
		var (
			objT = types.JSONObjectOf[T](nil)
			objM = types.JSONObjectOf[M](nil)
		)
		t.Run("Scan", func(t *testing.T) {
			for _, src := range []any{[]byte(``), []byte(`null`), "", nil} {
				Expect(t, objT.Scan(src), Succeed())
				Expect(t, objT.IsZero(), BeTrue())
				Expect(t, objT.Get(), BeNil[*T]())
				Expect(t, objM.Scan(src), Succeed())
				Expect(t, objM.IsZero(), BeTrue())
				Expect(t, objM.Get(), BeNil[*M]())
			}
			Expect(t, objT.Scan(100), Failed())
			Expect(t, objT.Scan(`{}`), Succeed())
			Expect(t, objT.Get(), Equal[*T](&T{}))
			Expect(t, objT.UnmarshalJSON([]byte(`{`)), Failed())
		})
		t.Run("Value", func(t *testing.T) {
			objT.Set(nil)
			v, _ := objT.Value()
			Expect(t, v, Equal[driver.Value](""))
			objM.Set(nil)
			v, _ = objT.Value()
			Expect(t, v, Equal[driver.Value](""))

			objT.Set(&T{1})
			v, _ = objT.Value()
			Expect(t, v, Equal[driver.Value](`{"a":1}`))
			objM.Set(&M{1: 1})
			v, _ = objM.Value()
			Expect(t, v, Equal[driver.Value](`{"1":1}`))
		})
		ExpectPanic[error](t, func() { types.JSONObjectOf(new(int)) })
	})

	_, err := types.DriverJSONValue(MustFailedJsonArshaler{})
	Expect(t, err, Failed())
}

type MustFailedJsonArshaler struct{}

func (MustFailedJsonArshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New("brick")
}
