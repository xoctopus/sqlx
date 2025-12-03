package structs_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/xoctopus/pkgx/pkg/pkgx"
	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/ptrx"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/internal/structs"
	"github.com/xoctopus/sqlx/pkg/types"
)

type RefOrg struct {
	OrgID uint64 `db:"f_org_id"`
}

type RefUser struct {
	UserID uint64 `db:"f_user_id"`
}

type User struct {
	types.AutoIncID

	RefUser
	RefOrg

	UserData
	Ignored any `db:"-" json:"org_info"`
}

func (User) TableName() string { return "t_user" }

type UserData struct {
	Name       string `db:"f_name,width=128"`
	Age        int    `db:"f_age"`
	unexported any
}

func ExampleFieldsFor() {
	ctx := contextx.Compose(
		def.CtxTagKey.Carry("db"),
		pkgx.CtxLoadTests.Carry(true),
	)(context.Background())

	fields := structs.FieldsFor(ctx, typx.NewRType(reflect.TypeFor[User]()))

	for _, f := range fields {
		fmt.Printf("%-10s %-6s %v\n", f.Name, f.Type.Name(), f.Loc)
	}

	seq := structs.FieldsSeqFor(ctx, typx.NewRType(reflect.TypeFor[User]()))
	for f := range seq {
		fmt.Printf("%-10s %-6s %v\n", f.Name, f.Type.Name(), f.Loc)
	}

	// Output:
	// f_id       uint64 [0 0]
	// f_user_id  uint64 [1 0]
	// f_org_id   uint64 [2 0]
	// f_name     string [3 0]
	// f_age      int    [3 1]
	// f_id       uint64 [0 0]
	// f_user_id  uint64 [1 0]
	// f_org_id   uint64 [2 0]
	// f_name     string [3 0]
	// f_age      int    [3 1]
}

type M struct {
	Sub
	*PtrSub
	F6 *string  `db:"f_f6"`
	F7 **string `db:"f_f7,deprecated='migrated'"`
}

func (M) TableName() string { return "t_m" }

type Sub struct {
	SubSub
	F2 string `db:"f_f2"`
}

type PtrSub struct {
	F3 []string          `db:"f_f3"`
	F4 map[string]string `db:"f_f4"`
	F5 *string           `db:"f_f5"`
}

type SubSub struct {
	SubSubSub
}

type SubSubSub struct {
	F1 string `db:"f_f1"`
}

var V = &M{
	Sub: Sub{
		SubSub: SubSub{
			SubSubSub{
				F1: "f1",
			},
		},
		F2: "f2",
	},
	PtrSub: &PtrSub{
		F3: []string{"f3"},
		F4: map[string]string{"f4": "f4"},
		F5: nil,
	},
	F6: nil,
	F7: ptrx.Ptr(ptrx.Ptr("f7")),
}

func TestField_Value(t *testing.T) {
	ctx := contextx.Compose(
		def.CtxTagKey.Carry("db"),
		pkgx.CtxLoadTests.Carry(true),
	)(context.Background())

	v := reflect.ValueOf(V).Elem()

	fields := structs.FieldsFor(ctx, typx.NewRType(reflect.TypeFor[M]()))
	Expect(t, fields, HaveLen[[]*structs.Field](7))
	Expect(t, fields[0].Name, Equal("f_f1"))
	Expect(t, fields[0].Value(v), Equal[any](V.F1))

	Expect(t, fields[1].Name, Equal("f_f2"))
	Expect(t, fields[1].Value(v), Equal[any](V.F2))

	Expect(t, fields[2].Name, Equal("f_f3"))
	Expect(t, fields[2].Value(v), Equal[any](V.F3))

	Expect(t, fields[3].Name, Equal("f_f4"))
	Expect(t, fields[3].Value(v), Equal[any](V.F4))

	Expect(t, fields[4].Name, Equal("f_f5"))
	Expect(t, fields[4].Value(v), Equal[any](V.F5))

	Expect(t, fields[5].Name, Equal("f_f6"))
	Expect(t, fields[5].Value(v), Equal[any](V.F6))

	Expect(t, fields[6].Name, Equal("f_f7"))
	Expect(t, fields[6].Value(v), Equal[any](V.F7))
}

func TestTableFields(t *testing.T) {
	ctx := contextx.Compose(
		def.CtxTagKey.Carry("db"),
		pkgx.CtxLoadTests.Carry(true),
	)(context.Background())

	fields := structs.TableFields(ctx, V)

	Expect(t, fields, HaveLen[[]*structs.TableField](6))
	Expect(t, fields[0].Field.Name, Equal("f_f1"))
	Expect(t, fields[0].Value.Interface(), Equal[any](V.F1))

	Expect(t, fields[1].Field.Name, Equal("f_f2"))
	Expect(t, fields[1].Value.Interface(), Equal[any](V.F2))

	Expect(t, fields[2].Field.Name, Equal("f_f3"))
	Expect(t, fields[2].Value.Interface(), Equal[any](V.F3))

	Expect(t, fields[3].Field.Name, Equal("f_f4"))
	Expect(t, fields[3].Value.Interface(), Equal[any](V.F4))

	Expect(t, fields[4].Field.Name, Equal("f_f5"))
	Expect(t, fields[4].Value.Interface(), Equal[any](V.F5))

	Expect(t, fields[5].Field.Name, Equal("f_f6"))
	Expect(t, fields[5].Value.Interface(), Equal[any](V.F6))
}
