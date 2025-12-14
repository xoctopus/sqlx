package sqlx

import (
	"context"
	"go/token"
	"go/types"
	"reflect"
	"strings"

	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
	"github.com/xoctopus/pkgx/pkg/pkgx"
	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/stringsx"

	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/internal/structs"
)

func NewModel(g genx.Context, t types.Type) *Model {
	n, ok := t.(*types.Named)
	if !ok {
		return nil
	}
	x := typx.NewTType(t)
	if x.Kind() != reflect.Struct {
		return nil
	}
	e := g.Package().TypeNames().ElementByName(n.Obj().Name())
	must.BeTrueF(e != nil, "expect %s lookup in package %s", n.Obj().Name(), g.Package().Path())

	m := &Model{
		ctx:    g,
		typ:    x,
		fields: structs.FieldsFor(x),
		t:      s.IdentTT(g.Context(), t),
		name:   "t_" + stringsx.LowerSnakeCase(n.Obj().Name()),
	}

	fm := make(map[string]*structs.Field)
	for _, f := range m.fields {
		fm[f.FieldName] = f
		if f.ColumnDef.Comment == "" {
			doc := g.Package().DocOf(token.Pos(typx.PosOfStructField(f.Field)))
			if doc != nil {
				f.ColumnDef.Comment = strings.Join(doc.Desc(), " ")
			}
		}
	}

	for _, line := range e.Doc().Desc() {
		if strings.HasPrefix(line, "@def ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "@def "))
			k := def.ParseKeyDef(line)
			must.BeTrueF(k != nil, "failed to parse @def: %s", line)
			for _, o := range k.Options {
				_, exists := fm[o.Name]
				must.BeTrueF(exists, "field def found: %s", o.Name)
			}
			switch k.Kind {
			case def.KEY_KIND__PRIMARY:
				m.primary = k
			case def.KEY_KIND__INDEX:
				m.indexes = append(m.indexes, k)
			case def.KEY_KIND__UNIQUE_INDEX:
				m.uniques = append(m.uniques, k)
			}
			continue
		}
		if strings.HasPrefix(line, "@attr ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "@attr "))
			if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
				switch parts[0] {
				case "TableName":
					m.name = strings.TrimSpace(parts[1])
				}
			}
			continue
		}
		m.desc = append(m.desc, strings.TrimSpace(line))
	}
	return m
}

type Model struct {
	ctx    genx.Context
	typ    typx.Type
	doc    *pkgx.Doc
	fields []*structs.Field

	t    s.Snippet
	name string
	desc []string

	primary *def.KeyDefine
	indexes []*def.KeyDefine
	uniques []*def.KeyDefine
}

func (m *Model) PrimaryColList() s.Snippet {
	if m.primary != nil {
		return s.Strings(",", "\n", m.primary.OptionsStrings()...)
	}
	return s.Snippet(&s.Placeholder{})
}

func (m *Model) IndexList(unique bool) s.Snippet {
	indexes := m.indexes
	if unique {
		indexes = m.uniques
	}

	if len(indexes) > 0 {
		ss := make([]s.Snippet, 0, len(indexes))
		for _, i := range indexes {
			ss = append(ss, s.BlockF("%q: {", i.Name))
			ss = append(ss, s.Strings(",", "\n", i.OptionsNames()...))
			ss = append(ss, s.Block("},"))
		}
		return s.Snippets(s.NewLine(1), ss...)
	}
	return s.Snippet(&s.Placeholder{})
}

func (m *Model) ModeledKeyDefList(ctx context.Context) s.Snippet {
	keyTyp := s.Expose(ctx, "github.com/xoctopus/sqlx/pkg/builder/modeled", "Key", m.t)

	ss := make([]s.Snippet, 0, len(m.indexes)+len(m.uniques)+1)
	if m.primary != nil {
		ss = append(ss, s.Compose(s.Block("Primary "), keyTyp))
	}
	for _, i := range m.indexes {
		ss = append(ss, s.Compose(s.Block(stringsx.UpperCamelCase(i.Name)+" "), keyTyp))
	}
	for _, i := range m.uniques {
		ss = append(ss, s.Compose(s.Block(stringsx.UpperCamelCase(i.Name)+" "), keyTyp))
	}
	return s.Snippets(s.NewLine(1), ss...)
}

func (m *Model) ModeledColDefList(ctx context.Context) s.Snippet {
	ss := []s.Snippet{s.NewLine(1)}
	for _, f := range m.fields {
		if f.ColumnDef.Comment != "" {
			ss = append(ss, s.Comments(f.ColumnDef.Comment))
		}
		ss = append(ss, s.Compose(s.Block(f.FieldName+" "), m.ModeledTCol(ctx, f.Type)))
	}
	return s.Snippets(s.NewLine(1), ss...)
}

func (m *Model) ModeledKeyInitList(_ context.Context) s.Snippet {
	ss := make([]s.Snippet, 0, len(m.indexes)+len(m.uniques)+1)
	if m.primary != nil {
		ss = append(ss, s.Compose(s.BlockF("Primary: m.MK(%q),", m.primary.Name)))
	}
	for _, i := range m.indexes {
		f := stringsx.UpperCamelCase(i.Name)
		ss = append(ss, s.Compose(s.BlockF("%s: m.MK(%q),", f, i.Name)))
	}
	for _, i := range m.uniques {
		f := stringsx.UpperCamelCase(i.Name)
		ss = append(ss, s.Compose(s.BlockF("%s: m.MK(%q),", f, i.Name)))
	}
	return s.Snippets(s.NewLine(1), ss...)
}

func (m *Model) ModeledColInitList(ctx context.Context) s.Snippet {
	ss := make([]s.Snippet, 0, len(m.fields))
	for _, f := range m.fields {
		ss = append(ss,
			s.Compose(
				s.Block(f.FieldName+": "),
				m.ModeledCT(ctx, f.Type),
				s.BlockF("(m.C(%q)),", f.FieldName),
			),
		)
	}
	return s.Snippets(s.NewLine(1), ss...)
}

func (m *Model) T() s.Snippet { return m.t }

func (m *Model) ModeledM(ctx context.Context) s.Snippet {
	return s.Expose(
		ctx,
		"github.com/xoctopus/sqlx/pkg/builder/modeled",
		"M", m.t,
	)
}

func (m *Model) ModeledTable(ctx context.Context) s.Snippet {
	return s.Expose(
		ctx,
		"github.com/xoctopus/sqlx/pkg/builder/modeled",
		"Table", m.t,
	)
}

func (m *Model) ModeledTCol(ctx context.Context, t typx.Type) s.Snippet {
	return s.Expose(
		ctx,
		"github.com/xoctopus/sqlx/pkg/builder/modeled",
		"TCol",
		m.t, s.Ident(ctx, t),
	)
}

func (m *Model) ModeledCT(ctx context.Context, t typx.Type) s.Snippet {
	return s.Expose(
		ctx,
		"github.com/xoctopus/sqlx/pkg/builder/modeled",
		"CT",
		m.t, s.Ident(ctx, t),
	)
}

func (m *Model) TableName() string { return m.name }

func (m *Model) TableDesc() []string { return m.desc }
