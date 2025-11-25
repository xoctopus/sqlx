package builder

import (
	"context"
	"reflect"

	"github.com/xoctopus/typex"

	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/internal/structs"
)

func scan(ctx context.Context, m any) Table {
	t := typex.Deref(typex.NewTType(ctx, reflect.TypeOf(m)))

	var (
		comments = map[string]string{}
		descs    = map[string][]string{}
		rels     = map[string][]string{}
	)

	// column comment
	if x, ok := m.(WithColumnComment); ok {
		comments = x.ColumnComment()
	}

	// column descriptions
	if x, ok := m.(WithColumnDesc); ok {
		descs = x.ColumnDesc()
	}

	// column rel: eg: OrgID string => t_org.f_org_id
	if x, ok := m.(WithColumnRel); ok {
		rels = x.ColumnRel()
	}

	tab := &table{
		cs: &columns{},
		ks: &keys{},
	}

	for _, f := range structs.FieldsFor(ctx, t) {
		c := &column[any]{
			fname: f.FieldName,
			name:  f.Name,
			def:   f.ColumnDef,
		}
		if text, ok := comments[c.fname]; ok {
			c.def.Comment = text
		}
		if lines, ok := descs[c.fname]; ok {
			c.def.Desc = lines
		}
		if rel, ok := rels[c.fname]; ok {
			c.def.Relation = rel
		}
		tab.cs.(ColsManager).AddCol(c.Of(tab))
	}

	if x, ok := m.(WithTableDesc); ok {
		tab.desc = x.TableDesc()
	}

	if x, ok := m.(WithPrimaryKey); ok {
		tab.ks.(KeysManager).AddKey(
			(&key{
				kind:    def.KEY_KIND__PRIMARY,
				name:    "primary",
				unique:  true,
				options: def.KeyColumnOptionByNames(x.PrimaryKey()...),
			}).Of(tab))
	}

	if x, ok := m.(WithUniqueIndexes); ok {
		for idx, fields := range x.UniqueIndexes() {
			name, method := def.ResolveIndexNameAndUsing(idx)
			tab.ks.(KeysManager).AddKey((&key{
				name:    name,
				method:  method,
				unique:  true,
				options: def.ResolveKeyColumnOptionsFromStrings(fields...),
			}).Of(tab))
		}
	}

	if x, ok := m.(WithIndexes); ok {
		for idx, fields := range x.Indexes() {
			name, method := def.ResolveIndexNameAndUsing(idx)
			tab.ks.(KeysManager).AddKey((&key{
				name:    name,
				method:  method,
				unique:  false,
				options: def.ResolveKeyColumnOptionsFromStrings(fields...),
			}).Of(tab))
		}
	}
	return tab
}
