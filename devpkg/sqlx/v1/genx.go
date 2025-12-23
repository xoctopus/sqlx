package sqlx

import (
	"bytes"
	_ "embed"
	"go/types"
	"strings"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

var (
	//go:embed tpls/sqlx.go_model.tpl
	tplModel []byte
	//go:embed tpls/sqlx.go_model_create.tpl
	tplCreate []byte
	//go:embed tpls/sqlx.go_model_datalist.tpl
	tplList []byte
	//go:embed tpls/sqlx.go_model_fetch_by_unique.tpl
	tplFetchByUnique []byte
	//go:embed tpls/sqlx.go_model_update_by_unique.tpl
	tplUpdateByUnique []byte
	//go:embed tpls/sqlx.go_model_delete_by_unique.tpl
	tplDeleteByUnique []byte
	//go:embed tpls/sqlx.go_model_mark_delete_by_unique.tpl
	tplSoftDeleteByUnique []byte
)

type uniqueAction struct {
	name string
	code []byte
}

var byUniqueActions = []*uniqueAction{
	{name: "Fetch", code: tplFetchByUnique},
	{name: "Update", code: tplUpdateByUnique},
	{name: "Delete", code: tplDeleteByUnique},
}

var (
	// typesPath   = "github.com/xoctopus/sqlx/pkg/types"
	_frag    = "github.com/xoctopus/sqlx/pkg/frag"
	_builder = "github.com/xoctopus/sqlx/pkg/builder"
	_modeled = "github.com/xoctopus/sqlx/pkg/builder/modeled"
	_helper  = "github.com/xoctopus/sqlx/pkg/helper"
	_session = "github.com/xoctopus/sqlx/pkg/session"
	_codex   = "github.com/xoctopus/x/codex"
	_errors  = "github.com/xoctopus/sqlx/pkg/errors"
)

func init() {
	genx.Register(&g{})
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
	actions := byUniqueActions

	if m.HasSoftDeletion() {
		actions = append(byUniqueActions, &uniqueAction{
			name: "MarkDeletion", code: tplSoftDeleteByUnique,
		})
	}

	ss := []s.Snippet{
		s.Template(
			bytes.NewReader(tplModel),
			s.Arg(ctx, "T", m.Ident()),
			s.Arg(ctx, "TableName", m.TableName()),
			s.Arg(ctx, "TableDesc", m.TableDesc()),
			s.Arg(ctx, "PrimaryColList", m.PrimaryColList()),
			s.Arg(ctx, "IndexList", m.IndexList(false)),
			s.Arg(ctx, "UniqueIndexList", m.IndexList(true)),
			s.Arg(ctx, "ModeledKeyDefList", m.ModeledKeyDefList(ctx)),
			s.Arg(ctx, "ModeledColDefList", m.ModeledColDefList(ctx)),
			s.Arg(ctx, "ModeledKeyInitList", m.ModeledKeyInitList(ctx)),
			s.Arg(ctx, "ModeledColInitList", m.ModeledColInitList(ctx)),
			s.Arg(ctx, "modeled.Table", m.ModeledTable(ctx)),
			s.Arg(ctx, "modeled.M", m.ModeledM(ctx)),
			s.Arg(ctx, "Register", m.Register()),
			s.ArgExposeUnsafe(ctx, "reflect", "ValueOf"),
			s.ArgExposeUnsafe(ctx, _builder, "Model").WithName("builder.Model"),
			s.ArgExposeUnsafe(ctx, _builder, "Col").WithName("builder.Col"),
			s.ArgExposeUnsafe(ctx, _builder, "ColsOf").WithName("builder.ColsOf"),
			s.ArgExposeUnsafe(ctx, _builder, "Assignment").WithName("builder.Assignment"),
			s.ArgExposeUnsafe(ctx, _builder, "GetColDef").WithName("builder.GetColDef"),
			s.ArgExposeUnsafe(ctx, _builder, "ColumnsAndValues").WithName("builder.ColumnsAndValues"),
		),
		s.Template(
			bytes.NewReader(tplCreate),
			s.Arg(ctx, "T", m.Ident()),
			s.Arg(ctx, "CreationMarker", m.CreationMarker()),
			s.Arg(ctx, "CreateComment", m.CommentOf("Create")),
			s.ArgExposeUnsafe(ctx, "context", "Context"),
			s.ArgExposeUnsafe(ctx, _helper, "ColumnsAndValuesForInsertion").WithName("helper.ColumnsAndValuesForInsertion"),
			s.ArgExposeUnsafe(ctx, _session, "For").WithName("session.For"),
			s.ArgExposeUnsafe(ctx, _builder, "Comment").WithName("builder.Comment"),
			s.ArgExposeUnsafe(ctx, _builder, "Insert").WithName("builder.Insert"),
		),
		s.Template(
			bytes.NewReader(tplList),
			s.Arg(ctx, "T", m.Ident()),
			s.Arg(ctx, "ListComment", m.CommentOf("List")),
			s.Arg(ctx, "SoftDeletionCondition", m.SoftDeletionCondition(ctx)),
			s.ArgExposeUnsafe(ctx, "context", "Context"),
			s.ArgExposeUnsafe(ctx, _builder, "SqlCondition").WithName("builder.SqlCondition"),
			s.ArgExposeUnsafe(ctx, _builder, "Additions").WithName("builder.Additions"),
			s.ArgExposeUnsafe(ctx, _builder, "Comment").WithName("builder.Comment"),
			s.ArgExposeUnsafe(ctx, _builder, "Select").WithName("builder.Select"),
			s.ArgExposeUnsafe(ctx, _builder, "Col").WithName("builder.Col"),
			s.ArgExposeUnsafe(ctx, _builder, "ColsOf").WithName("builder.ColsOf"),
			s.ArgExposeUnsafe(ctx, _builder, "And").WithName("builder.And"),
			s.ArgExposeUnsafe(ctx, _builder, "Where").WithName("builder.Where"),
			s.ArgExposeUnsafe(ctx, _frag, "Fragment").WithName("frag.Fragment"),
			s.ArgExposeUnsafe(ctx, _helper, "Scan").WithName("helper.Scan"),
			s.ArgExposeUnsafe(ctx, _session, "For").WithName("session.For"),
		),
	}

	for _, action := range actions {
		for _, names := range m.UniqueNames() {
			suffix := strings.Join(names, "And")
			args := []*s.TArg{
				s.Arg(ctx, "T", m.Ident()),
				s.Arg(ctx, "UniqueSuffix", s.Block(suffix)),
				s.Arg(ctx, "UniqueConds", m.UniqueConditions(ctx, names)),
				s.Arg(ctx, "UniqueFields", m.UniqueFields(names)),
				s.Arg(ctx, "SoftDeletionCondition", m.SoftDeletionCondition(ctx)),
				s.Arg(ctx, action.name+"Comment", m.CommentOf(action.name+"By"+suffix)),
				s.Arg(ctx, "ModificationMarker", m.ModificationMarker()),
				s.Arg(ctx, "DeletionMarker", m.DeletionMarker()),
				s.ArgExposeUnsafe(ctx, "context", "Context"),
				s.ArgExposeUnsafe(ctx, _frag, "Fragment").WithName("frag.Fragment"),
				s.ArgExposeUnsafe(ctx, _builder, "Select").WithName("builder.Select"),
				s.ArgExposeUnsafe(ctx, _builder, "Update").WithName("builder.Update"),
				s.ArgExposeUnsafe(ctx, _builder, "Delete").WithName("builder.Delete"),
				s.ArgExposeUnsafe(ctx, _builder, "Where").WithName("builder.Where"),
				s.ArgExposeUnsafe(ctx, _builder, "Limit").WithName("builder.Limit"),
				s.ArgExposeUnsafe(ctx, _builder, "Comment").WithName("builder.Comment"),
				s.ArgExposeUnsafe(ctx, _builder, "And").WithName("builder.And"),
				s.ArgExposeUnsafe(ctx, _builder, "Col").WithName("builder.Col"),
				s.ArgExposeUnsafe(ctx, _builder, "CC").WithName("builder.CC"),
				s.ArgExposeUnsafe(ctx, _builder, "Eq").WithName("builder.Eq"),
				s.ArgExposeUnsafe(ctx, _builder, "Neq").WithName("builder.Neq"),
				s.ArgExposeUnsafe(ctx, _session, "For").WithName("session.For"),
				s.ArgExposeUnsafe(ctx, _helper, "Scan").WithName("helper.Scan"),
			}

			if action.name == "Update" {
				args = append(
					args,
					s.Arg(ctx, "codex.New", s.ExposeUnsafe(ctx, _codex, "New")),
					s.Arg(ctx, "errors.NOTFOUND", s.ExposeUnsafe(ctx, _errors, "NOTFOUND")),
				)
			}
			if action.name == "MarkDeletion" {
				args = append(
					args,
					s.ArgExposeUnsafe(ctx, "database/sql/driver", "Value").WithName("driver.Value"),
				)
			}

			ss = append(
				ss,
				s.Template(bytes.NewReader(action.code), args...),
			)
		}
	}

	c.Render(s.Snippets(s.NewLine(1), ss...))
	return nil
}
