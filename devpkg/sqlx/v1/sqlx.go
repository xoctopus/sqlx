package sqlx

import (
	"bytes"
	"context"
	"fmt"
	"go/token"
	"go/types"
	"reflect"
	"strings"

	"github.com/xoctopus/genx/pkg/docx"
	"github.com/xoctopus/genx/pkg/genx"
	s "github.com/xoctopus/genx/pkg/snippet"
	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/slicex"
	"github.com/xoctopus/x/stringsx"

	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/internal/structs"
)

func NewModel(g genx.Context, t types.Type) *Model {
	x := typx.NewTType(t)
	if x.Kind() != reflect.Struct {
		return nil
	}

	e := g.Package().TypeNames().ElementByName(x.Name())
	must.NotNilF(e, "expect %s lookup in package %s", x.Name(), g.Package().Path())
	must.BeTrueF(types.Identical(e.Type(), t), "")

	doc := docx.Parse(e.Doc())

	m := &Model{
		typ:    x,
		ptr:    typx.NewTType(types.NewPointer(e.Type())),
		fields: structs.FieldsFor(x),
		ident:  s.IdentTT(g.Context(), t),
		attrs: map[Attr]string{
			AttrTableName: "t_" + stringsx.LowerSnakeCase(x.Name()),
		},
		desc: doc.Descriptions(e.Name()),
	}

	if p := g.PackageByPath("github.com/xoctopus/sqlx/pkg/types"); p != nil {
		typenames := p.TypeNames()
		if tn := typenames.ElementByName("CreationMarker"); tn != nil {
			m.tCreationMarker = tn.Type()
		}
		if tn := typenames.ElementByName("ModificationMarker"); tn != nil {
			m.tModificationMarker = tn.Type()
		}
		if tn := typenames.ElementByName("DeletionMarker"); tn != nil {
			m.tDeletionMarker = tn.Type()
		}
		if tn := typenames.ElementByName("SoftDeletion"); tn != nil {
			m.tSoftDeletion = tn.Type()
		}
	}

	fm := make(map[string]*structs.Field)
	for _, f := range m.fields {
		fm[f.FieldName] = f
		if f.ColumnDef.Comment == "" {
			e.GetFieldDocByName(f.FieldName)
			d := docx.Parse(g.Packages().DocByPos(token.Pos(typx.PosOfStructField(f.Field))))
			f.ColumnDef.Comment = d.Description(f.FieldName, " ")
			if annotations, ok := d.AnnotationsByName("rel"); ok {
				f.ColumnDef.Relation = slicex.Mapping(
					annotations,
					func(a docx.Annotation) string { return a.Text() },
				)
			}
		}
	}

	annotations, _ := doc.AnnotationsByName(identifier)
	for _, anno := range annotations {
		switch k := anno.Key(); k {
		case string(AttrTableName), string(AttrRegister):
			m.attrs[Attr(k)] = anno.Value()
		default:
			keydef := def.ParseKeyDef(strings.ToLower(k), anno.Value())
			must.BeTrueF(keydef != nil, "failed to parse %s @def: %s", m.typ.Name(), anno.Text())
			for _, o := range keydef.Options {
				_, exists := fm[o.Name]
				must.BeTrueF(exists, "failed to parse %s field def not found: %s", m.typ.Name(), o.Name)
			}
			switch kind := keydef.Kind; kind {
			case def.KEY_KIND__PRIMARY:
				if m.primary == nil {
					m.primary = keydef
				}
			case def.KEY_KIND__INDEX:
				m.indexes = append(m.indexes, keydef)
			case def.KEY_KIND__UNIQUE_INDEX:
				m.uniques = append(m.uniques, keydef)
			}
		}
	}

	return m
}

type Model struct {
	typ    typx.Type
	ptr    typx.Type
	fields []*structs.Field

	ident s.Snippet
	attrs map[Attr]string
	desc  []string

	primary *def.KeyDefine
	indexes []*def.KeyDefine
	uniques []*def.KeyDefine

	tSoftDeletion       types.Type
	tCreationMarker     types.Type
	tModificationMarker types.Type
	tDeletionMarker     types.Type
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

func (m *Model) ColumnRelList() s.Snippet {
	ss := make([]s.Snippet, 0)
	for _, f := range m.fields {
		if len(f.ColumnDef.Relation) > 0 {
			ss = append(ss, s.BlockF("%q: {", f.FieldName))
			for _, r := range f.ColumnDef.Relation {
				ss = append(ss, s.BlockF("%q,", r))
			}
			ss = append(ss, s.BlockF("},"))
		}
	}
	return s.Snippets(s.NewLine(1), ss...)
}

func (m *Model) TagMapping() s.Snippet {
	ss := make([]s.Snippet, 0, len(m.fields)*2)
	keys := make(map[string]struct{})
	for _, f := range m.fields {
		if f.Flag != nil {
			if column := f.Flag.Name(); len(column) > 0 && column != "-" {
				if _, ok := keys[column]; !ok {
					ss = append(ss, s.BlockF("%q: %q,", column, column))
					keys[column] = struct{}{}
				}
				if f.NameTag != nil {
					for _, k := range []string{"json", "name"} {
						if flag := f.NameTag.Get(k); flag != nil {
							if name := flag.Name(); len(name) > 0 && name != "-" {
								if _, ok := keys[name]; !ok {
									ss = append(ss, s.BlockF("%q: %q,", name, column))
									keys[column] = struct{}{}
								}
							}
						}
					}
				}
			}
		}
	}
	return s.Snippets(s.NewLine(1), ss...)
}

func (m *Model) ModeledKeyDefList(ctx context.Context) s.Snippet {
	keyTyp := s.Expose(ctx, _modeled, "Key", m.ident)

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
	return s.Compose(
		s.ExposeUnsafe(ctx, _modeled, "M"),
		s.Block("["), m.ident, s.Block("]"),
	)
	// return s.Expose(ctx, _modeled, "M", m.ident)
}

func (m *Model) ModeledTable(ctx context.Context) s.Snippet {
	return s.Compose(
		s.ExposeUnsafe(ctx, _modeled, "Table"),
		s.Block("["), m.ident, s.Block("]"),
	)
	// return s.Expose(ctx, _modeled, "Table", m.ident)
}

func (m *Model) ModeledTCol(ctx context.Context, t typx.Type) s.Snippet {
	return s.Compose(
		s.ExposeUnsafe(ctx, _modeled, "TCol"),
		s.Block("["), m.ident, s.Block(", "), s.Ident(ctx, t), s.Block("]"),
	)
	// return s.Expose(ctx, _modeled, "TCol", m.ident, s.Ident(ctx, t))
}

func (m *Model) ModeledCT(ctx context.Context, t typx.Type) s.Snippet {
	return s.Compose(
		s.ExposeUnsafe(ctx, _modeled, "CT"),
		s.Block("["), m.ident, s.Block(", "), s.Ident(ctx, t), s.Block("]"),
	)
	// return s.Expose(ctx, _modeled, "CT", m.ident, s.Ident(ctx, t))
}

func (m *Model) TableName() s.Snippet {
	return s.BlockRaw(m.attrs[AttrTableName])
}

func (m *Model) TableDesc() (snip s.Snippet) {
	if len(m.desc) > 0 {
		snip = s.Strings(",", "\n", m.desc...)
	}
	return snip
}

func (m *Model) Register() s.Snippet {
	if register := m.attrs[AttrRegister]; len(register) > 0 {
		return s.BlockF("%s.Add(T%s)", register, m.typ.Name())
	}
	return nil
}

func (m *Model) CreationMarker() s.Snippet {
	if m.tCreationMarker != nil &&
		(m.typ.Implements(m.tCreationMarker) || m.ptr.Implements(m.tCreationMarker)) {
		return s.Block("m.MarkCreatedAt()")
	}
	return nil
}

func (m *Model) ModificationMarker() s.Snippet {
	if m.tModificationMarker != nil &&
		(m.typ.Implements(m.tModificationMarker) || m.ptr.Implements(m.tModificationMarker)) {
		return s.Block("m.MarkModifiedAt()")
	}
	return nil
}

func (m *Model) DeletionMarker() s.Snippet {
	if m.tDeletionMarker != nil &&
		(m.typ.Implements(m.tDeletionMarker) || m.ptr.Implements(m.tModificationMarker)) {
		return s.Block("m.MarkDeletedAt()")
	}
	return nil
}

func (m *Model) HasSoftDeletion() bool {
	return m.tSoftDeletion != nil &&
		(m.typ.Implements(m.tSoftDeletion) || m.ptr.Implements(m.tSoftDeletion))
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
		_, _ = fmt.Fprintf(code, "\t\tT#T#.%s.AsCond(#builder.Eq#(m.%s)),", name, name)
	}
	return s.Template(
		code,
		s.Arg(ctx, "T", m.Ident()),
		s.Arg(ctx, "builder.Eq", s.ExposeUnsafe(ctx, "github.com/xoctopus/sqlx/pkg/builder", "Eq")),
	)
}

func (m *Model) SoftDeletionCondition(ctx context.Context) s.Snippet {
	if !m.HasSoftDeletion() {
		return nil
	}
	args := []*s.TArg{
		s.Arg(ctx, "T", m.Ident()),
		s.Arg(ctx, "frag.Fragment", s.ExposeUnsafe(ctx, _frag, "Fragment")),
		s.Arg(ctx, "builder.Eq", s.ExposeUnsafe(ctx, _builder, "Eq")),
		s.Arg(ctx, "builder.CC", s.ExposeUnsafe(ctx, _builder, "CC")),
		s.Arg(ctx, "driver.Value", s.ExposeUnsafe(ctx, "database/sql/driver", "Value")),
	}

	code := `
@def frag.Fragment
@def builder.Eq
@def builder.CC
@def driver.Value
--SoftDeletionCondition
	deletion, _, v := m.SoftDeletion()
	conds = append(
		conds,
		#builder.CC#[#driver.Value#](T#T#.C(deletion)).AsCond(#builder.Eq#(v)),
	)`
	args = append(
		args,
		s.Arg(ctx, "T", m.Ident()),
	)

	return s.Template(bytes.NewBufferString(code), args...)
}
