package sqlx

import (
	"bytes"
	"context"
	"fmt"
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

var (
	path = "github.com/xoctopus/sqlx/pkg/types"
	pkg  = pkgx.NewPackages(context.Background(), path).Package(path)

	tCreationMarker     = pkg.TypeNames().ElementByName("CreationMarker").Type()
	tModificationMarker = pkg.TypeNames().ElementByName("ModificationMarker").Type()
	tSoftDeletion       = pkg.TypeNames().ElementByName("SoftDeletion").Type()
)

func NewModel(g genx.Context, t types.Type) *Model {
	x := typx.NewTType(t)
	if x.Kind() != reflect.Struct {
		return nil
	}

	e := g.Package().TypeNames().ElementByName(x.Name())
	must.NotNilF(e, "expect %s lookup in package %s", x.Name(), g.Package().Path())
	must.BeTrueF(types.Identical(e.Type(), t), "")

	m := &Model{
		typ:    x,
		ptr:    typx.NewTType(types.NewPointer(e.Type())),
		fields: structs.FieldsFor(x),
		ident:  s.IdentTT(g.Context(), t),
		attrs: map[Attr]string{
			AttrTableName: "t_" + stringsx.LowerSnakeCase(x.Name()),
		},
	}

	fm := make(map[string]*structs.Field)
	for _, f := range m.fields {
		fm[f.FieldName] = f
		if f.ColumnDef.Comment == "" {
			doc := g.Packages().DocOf(token.Pos(typx.PosOfStructField(f.Field)))
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
				if attr := HasAttr(parts[0]); attr != "" {
					m.attrs[attr] = parts[1]
				}
			}
			continue
		}
		m.desc = append(m.desc, strings.TrimSpace(line))
	}
	return m
}

type Model struct {
	typ    typx.Type
	ptr    typx.Type
	doc    *pkgx.Doc
	fields []*structs.Field

	ident s.Snippet
	attrs map[Attr]string
	desc  []string

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
	keyTyp := s.Expose(ctx, "github.com/xoctopus/sqlx/pkg/builder/modeled", "Key", m.ident)

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

func (m *Model) Ident() s.Snippet { return m.ident }

func (m *Model) ModeledM(ctx context.Context) s.Snippet {
	return s.Expose(
		ctx,
		"github.com/xoctopus/sqlx/pkg/builder/modeled",
		"M", m.ident,
	)
}

func (m *Model) ModeledTable(ctx context.Context) s.Snippet {
	return s.Expose(
		ctx,
		"github.com/xoctopus/sqlx/pkg/builder/modeled",
		"Table", m.ident,
	)
}

func (m *Model) ModeledTCol(ctx context.Context, t typx.Type) s.Snippet {
	return s.Expose(
		ctx,
		"github.com/xoctopus/sqlx/pkg/builder/modeled",
		"TCol",
		m.ident, s.Ident(ctx, t),
	)
}

func (m *Model) ModeledCT(ctx context.Context, t typx.Type) s.Snippet {
	return s.Expose(
		ctx,
		"github.com/xoctopus/sqlx/pkg/builder/modeled",
		"CT",
		m.ident, s.Ident(ctx, t),
	)
}

func (m *Model) TableName() s.Snippet {
	return s.BlockRaw(m.attrs[AttrTableName])
}

func (m *Model) TableDesc() s.Snippet {
	if len(m.desc) > 0 {
		return s.Strings(",", "\n", m.desc...)
	}
	return nil
}

func (m *Model) Register() s.Snippet {
	if register, _ := m.attrs[AttrRegister]; len(register) > 0 {
		return s.BlockF("%s.Add(T%s)", register, m.typ.Name())
	}
	return nil
}

func (m *Model) CreationMarker() s.Snippet {
	if m.typ.Implements(tCreationMarker) || m.ptr.Implements(tCreationMarker) {
		return s.Block("m.MarkCreatedAt()")
	}
	return nil
}

func (m *Model) ModificationMarker() s.Snippet {
	if m.typ.Implements(tModificationMarker) || m.ptr.Implements(tModificationMarker) {
		return s.Block("m.MarkModifiedAt()")
	}
	return nil
}

func (m *Model) CommentOf(ref string) s.Snippet {
	return s.BlockF("\"%s.%s\"", m.typ.Name(), ref)
}

func (m *Model) UniqueNames() [][]string {
	suffixes := make([][]string, 0)

	uniques := append([]*def.KeyDefine{m.primary}, m.uniques...)
	for _, i := range uniques {
		if i != nil {
			suffixes = append(suffixes, i.OptionsNames())
		}
	}
	return suffixes
}

func (m *Model) UniqueFields(names []string) s.Snippet {
	fields := make([]string, len(names))
	for i, name := range names {
		fields[i] = m.typ.Name() + "." + name
	}
	return s.Block(strings.Join(fields, " and "))
}

func (m *Model) UniqueConditions(ctx context.Context, names []string) s.Snippet {
	code := bytes.NewBufferString(`
@def T
@def builder.Eq
--FetchByUniqueConds
`)
	for i, name := range names {
		if i > 0 {
			code.WriteString("\n")
		}
		code.WriteString(fmt.Sprintf("\t\tT#T#.%s.AsCond(#builder.Eq#(m.%s)),", name, name))
	}
	return s.Template(
		code,
		s.Arg(ctx, "T", m.Ident()),
		s.Arg(ctx, "builder.Eq", s.ExposeUnsafe(ctx, "github.com/xoctopus/sqlx/pkg/builder", "Eq")),
	)
}
