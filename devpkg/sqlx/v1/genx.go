package sqlx

import (
	"bytes"
	"database/sql/driver"
	_ "embed"
	"go/types"
	"reflect"
	"strings"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
	"github.com/xoctopus/typx/pkg/typx"
)

var (
	//go:embed sqlx.go_model.tpl
	tplModel []byte
	//go:embed sqlx.go_model_create.tpl
	tplCreate []byte
	//go:embed sqlx.go_model_datalist.tpl
	tplList []byte
	//go:embed sqlx.go_model_fetch_by_unique.tpl
	tplFetchByUnique []byte
	//go:embed sqlx.go_model_update_by_unique.tpl
	tplUpdateByUnique []byte
	//go:embed sqlx.go_model_delete_by_unique.tpl
	tplDeleteByUnique []byte
	//go:embed sqlx.go_model_soft_delete_by_unique.tpl
	tplSoftDeleteByUnique []byte
)

var byUniqueActions = map[string][]byte{
	"Fetch":  tplFetchByUnique,
	"Update": tplUpdateByUnique,
	"Delete": tplDeleteByUnique,
}

func init() {
	// genx.Register(&g{})
}

type g struct{}

func (x *g) New(c genx.Context) genx.Generator {
	return &g{}
}

func (x *g) Identifier() string {
	return "model"
}

func (x *g) Generate(c genx.Context, t types.Type) error {
	m := NewModel(c, t)
	if m == nil {
		return nil
	}
	ctx := c.Context()

	if typx.NewTType(t).Implements(tSoftDeletion) {
		byUniqueActions["SoftDelete"] = tplSoftDeleteByUnique
	}

	var (
		typesPath   = "github.com/xoctopus/sqlx/pkg/types"
		fragPath    = "github.com/xoctopus/sqlx/pkg/frag"
		builderPath = "github.com/xoctopus/sqlx/pkg/builder"
		sessionPath = "github.com/xoctopus/sqlx/pkg/session"
		helperPath  = "github.com/xoctopus/sqlx/pkg/helper"
		errorsPath  = "github.com/xoctopus/sqlx/pkg/errors"
		codexPath   = "github.com/xoctopus/x/codex"
	)

	args := []*s.TArg{
		// context.Context
		s.ArgExpose(ctx, "context", "Context"),
		// reflect.ValueOf
		s.ArgExpose(ctx, "reflect", "ValueOf"),
		// driver.Value
		s.Arg(ctx, "driver.Value", s.IdentRT(ctx, reflect.TypeFor[driver.Value]())),

		// model ident
		s.Arg(ctx, "T", m.Ident()),

		s.Arg(ctx, "TableName", m.TableName()),
		s.Arg(ctx, "TableDesc", m.TableDesc()),
		// primary column list
		s.Arg(ctx, "PrimaryColList", m.PrimaryColList()),
		// index list
		s.Arg(ctx, "IndexList", m.IndexList(false)),
		// unique index list
		s.Arg(ctx, "UniqueIndexList", m.IndexList(true)),
		// modeled key define list
		s.Arg(ctx, "ModeledKeyDefList", m.ModeledKeyDefList(ctx)),
		// modeled column define list
		s.Arg(ctx, "ModeledColDefList", m.ModeledColDefList(ctx)),
		// modeled key init list
		s.Arg(ctx, "ModeledKeyInitList", m.ModeledKeyInitList(ctx)),
		// modeled column init list
		s.Arg(ctx, "ModeledColInitList", m.ModeledColInitList(ctx)),
		// modeled.Table[T]
		s.Arg(ctx, "ModeledTable", m.ModeledTable(ctx)),
		// modeled.M[T]()
		s.Arg(ctx, "ModeledM", m.ModeledM(ctx)),
		// Register
		s.Arg(ctx, "Register", m.Register()),

		// CreationMarker
		s.Arg(ctx, "CreationMarker", m.CreationMarker()),
		// ModificationMarker
		s.Arg(ctx, "ModificationMarker", m.ModificationMarker()),

		// CreateComment
		s.Arg(ctx, "CreateComment", m.CommentOf("Create")),
		// ListComment
		s.Arg(ctx, "ListComment", m.CommentOf("List")),

		// builder.Model
		s.ArgExpose(ctx, builderPath, "Model").WithName("builder.Model"),
		// builder.SqlCondition
		s.ArgExpose(ctx, builderPath, "SqlCondition").WithName("builder.SqlCondition"),
		// builder.Additions
		s.ArgExpose(ctx, builderPath, "Additions").WithName("builder.Additions"),
		// builder.CC
		s.Arg(ctx, "builder.CC", s.ExposeUnsafe(ctx, builderPath, "CC")),
		// builder.Eq
		s.Arg(ctx, "builder.Eq", s.ExposeUnsafe(ctx, builderPath, "Eq")),
		// builder.And
		s.ArgExpose(ctx, builderPath, "And").WithName("builder.And"),
		// builder.Col
		s.ArgExpose(ctx, builderPath, "Col").WithName("builder.Col"),
		// builder.Cols
		s.ArgExpose(ctx, builderPath, "Cols").WithName("builder.Cols"),
		// builder.ColsOf
		s.ArgExpose(ctx, builderPath, "ColsOf").WithName("builder.ColsOf"),
		// builder.ColsIterOf
		s.ArgExpose(ctx, builderPath, "ColsIterOf").WithName("builder.ColsIterOf"),
		// builder.GetColOf
		s.ArgExpose(ctx, builderPath, "GetColDef").WithName("builder.GetColDef"),
		// builder.Insert
		s.ArgExpose(ctx, builderPath, "Insert").WithName("builder.Insert"),
		// builder.Select
		s.ArgExpose(ctx, builderPath, "Select").WithName("builder.Select"),
		// builder.Update
		s.ArgExpose(ctx, builderPath, "Update").WithName("builder.Update"),
		// builder.Delete
		s.ArgExpose(ctx, builderPath, "Delete").WithName("builder.Delete"),
		// builder.Where
		s.ArgExpose(ctx, builderPath, "Where").WithName("builder.Where"),
		// builder.Comment
		s.ArgExpose(ctx, builderPath, "Comment").WithName("builder.Comment"),
		// builder.Limit
		s.ArgExpose(ctx, builderPath, "Limit").WithName("builder.Limit"),

		// helper.Scan
		s.ArgExpose(ctx, helperPath, "Scan").WithName("helper.Scan"),
		// helper.ColumnsAndValues
		s.ArgExpose(ctx, helperPath, "ColumnsAndValuesForInsertion").WithName("helper.ColumnsAndValuesForInsertion"),

		// types.SoftDeletion
		s.ArgExpose(ctx, typesPath, "SoftDeletion").WithName("types.SoftDeletion"),

		// session.For
		s.ArgExpose(ctx, sessionPath, "For").WithName("session.For"),

		// frag.Fragment
		s.ArgExpose(ctx, fragPath, "Fragment").WithName("frag.Fragment"),

		// codex.New
		s.Arg(ctx, "codex.New", s.ExposeUnsafe(ctx, codexPath, "New")),

		// errors.NOTFOUND
		s.Arg(ctx, "errors.NOTFOUND", s.ExposeUnsafe(ctx, errorsPath, "NOTFOUND")),
	}

	ss := []s.Snippet{
		s.Template(bytes.NewReader(tplModel), args...),
		s.Template(bytes.NewReader(tplCreate), args...),
		s.Template(bytes.NewReader(tplList), args...),
	}

	for action, tpl := range byUniqueActions {
		for _, names := range m.UniqueNames() {
			suffix := strings.Join(names, "And")
			ss = append(
				ss,
				s.Template(
					bytes.NewReader(tpl),
					append(
						args,
						s.Arg(ctx, "UniqueSuffix", s.Block(suffix)),
						s.Arg(ctx, "UniqueConds", m.UniqueConditions(ctx, names)),
						s.Arg(ctx, "UniqueFields", m.UniqueFields(names)),
						s.Arg(ctx, action+"Comment", m.CommentOf(action+"By"+suffix)),
					)...,
				),
			)
		}
	}

	c.Render(s.Snippets(s.NewLine(1), ss...))
	return nil
}
