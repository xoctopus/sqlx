package sqlx

import (
	"bytes"
	_ "embed"
	"go/types"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
)

//go:embed sqlx.go_model.tpl
var tplModel []byte

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

	c.Render(s.Template(
		bytes.NewReader(tplModel),

		// model type name
		s.Arg(ctx, "T", m.T()),
		// table name
		s.Arg(ctx, "TableName", s.BlockRaw(m.TableName())),
		// table desc
		s.Arg(ctx, "TableDesc", s.Strings(",", "\n", m.TableDesc()...)),
		// context.Background
		s.ArgExpose(ctx, "context", "Background"),
		// builder.Model
		s.Arg(ctx, "BuilderModel", s.Expose(ctx, "github.com/xoctopus/sqlx/pkg/builder", "Model")),
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
	))
	return nil
}
