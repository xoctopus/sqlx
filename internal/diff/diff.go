package diff

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/xoctopus/x/misc/must"

	"github.com/xoctopus/sqlx/internal/sql/adaptor"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
)

type action struct {
	typ   ActType
	tar   string
	frags []frag.Fragment
}

func (a *action) IsNil() bool {
	return len(a.frags) == 0
}

func (a *action) Frag(ctx context.Context) frag.Iter {
	return frag.Compose("\n", a.frags...).Frag(ctx)
}

type differ struct {
	dialect adaptor.Dialect
	actions []*action

	indexes map[string]struct{}
	columns map[string]ActType
}

func (d *differ) IsNil() bool {
	return len(d.actions) == 0
}

func (d *differ) Frag(ctx context.Context) frag.Iter {
	actions := slices.SortedFunc(slices.Values(d.actions), func(a *action, b *action) int {
		ret := cmp.Compare(a.typ, b.typ)
		if ret == 0 {
			return cmp.Compare(a.tar, b.tar)
		}
		return ret
	})

	return frag.ComposeSeq("\n", frag.NonNil(slices.Values(actions))).Frag(ctx)
}

func (d *differ) do(typ ActType, target string, frags ...frag.Fragment) {
	if len(frags) > 0 {
		if typ == ACT_DROP_IDX || typ == ACT_CREATE_IDX {
			name := fmt.Sprintf("%d/%s", typ, target)
			if _, ok := d.indexes[name]; ok {
				return
			}
			d.indexes[name] = struct{}{}
		}

		d.actions = append(d.actions, &action{
			typ:   typ,
			tar:   target,
			frags: frags,
		})
	}
}

// Diff curr and next table for migration
// curr from catalog scanning
// next from current tables define from code ast
func Diff(ctx context.Context, d adaptor.Dialect, curr, next builder.Table) frag.Fragment {
	dd := &differ{
		dialect: d,
		indexes: make(map[string]struct{}),
		columns: make(map[string]ActType),
	}

	// create new table
	if curr == nil {
		dd.do(ACT_CREATE_TABLE, next.TableName(), d.CreateTableIfNotExists(next)...)
		return dd
	}

	if mode, ok := CtxMode.From(ctx); ok && mode.Is(MODE_CREATE_TABLE) {
		return nil
	}

	// diff columns
	for nextc := range next.Cols() {
		currc := curr.C(nextc.Name())
		if !frag.IsNil(nextc) && !frag.IsNil(currc) {
			if def := builder.GetColDef(nextc); def.Deprecated != nil {
				if renameto := def.Deprecated.RenameTo; renameto != "" {
					// rename column
					if c := curr.C(renameto); c != nil {
						dd.columns[c.Name()] = ACT_DROP_COL
						dd.do(ACT_DROP_COL, c.Name(), d.DropColumn(c))
					}
					renamedc := next.C(renameto)
					must.BeTrueF(!frag.IsNil(renamedc), "column '%s' is not declared", renameto)
					// ALTER TABLE @table RENAME nextc TO renamedc
					dd.columns[renameto] = ACT_RENAME_COL
					dd.do(ACT_RENAME_COL, nextc.Name(), d.RenameColumn(nextc, renamedc))
					curr.(builder.ColsManager).AddCol(renamedc)
					continue
				}
				// else {
				// drop deprecated column
				// dd.do(ACT_DROP_COL, nextc.Name(), d.DropColumn(nextc))
				// }
				continue
			}
			// diff datatype and modify
			typCurr, _ := frag.Collect(ctx, d.DBType(builder.GetColDef(currc)))
			typNext, _ := frag.Collect(ctx, d.DBType(builder.GetColDef(nextc)))
			if !strings.EqualFold(typCurr, typNext) {
				dd.columns[nextc.Name()] = ACT_MODIFY_COL
				dd.do(ACT_MODIFY_COL, nextc.Name(), d.ModifyColumn(nextc, currc))
			}
			dd.columns[nextc.Name()] = ACT_KEEP_COL
			continue
		}
	}

	// drop deprecated/tmp column if mark with prefix dd_tmp__
	for currc := range curr.Cols() {
		if strings.HasPrefix(currc.Name(), "dd_tmp__") && frag.IsNil(next.C(currc.Name())) {
			dd.do(ACT_DROP_COL, currc.Name(), d.DropColumn(currc))
		}
		// if frag.IsNil(next.C(currc.Name())) {
		// 	dd.do(ACT_DROP_COL, currc.Name(), d.DropColumn(currc))
		// }
	}

	// create new columns and drop deprecated column
	for nextc := range next.Cols() {
		if def := builder.GetColDef(nextc); def.Deprecated == nil {
			if _, changed := dd.columns[nextc.Name()]; !changed {
				dd.do(ACT_CREATE_COL, nextc.Name(), d.AddColumn(nextc))
			}
		} else {
			// it MUST be declared to drop real data column
			if def.Deprecated.RenameTo == "" && !frag.IsNil(curr.C(nextc.Name())) {
				dd.do(ACT_DROP_COL, nextc.Name(), d.DropColumn(nextc))
			}
		}
	}

	for nextk := range next.Keys() {
		name := nextk.Name()
		if nextk.IsPrimary() {
			if currk := curr.K(nextk.Name()); currk == nil {
				dd.do(ACT_CREATE_IDX, nextk.Name(), d.AddIndex(nextk))
			}
			// won't modify primary key
			continue
		}

		for c := range nextk.Cols() {
			if act, ok := dd.columns[c.Name()]; ok && act == ACT_MODIFY_COL {
				// reindex when column modified
				dd.do(ACT_DROP_IDX, nextk.Name(), d.DropIndex(nextk))
				dd.do(ACT_CREATE_IDX, nextk.Name(), d.AddIndex(nextk))
			}
		}

		if currk := curr.K(name); currk == nil {
			dd.do(ACT_CREATE_IDX, nextk.Name(), d.AddIndex(nextk))
		} else {
			if !nextk.IsPrimary() {
				nextkdef, _ := frag.Collect(
					context.Background(), builder.ColsIterOf(
						slices.Values(slices.SortedFunc(
							nextk.Cols(),
							func(c1, c2 builder.Col) int {
								return cmp.Compare(c1.Name(), c2.Name())
							},
						)),
					),
				)
				currkdef, _ := frag.Collect(
					context.Background(), builder.ColsIterOf(
						slices.Values(slices.SortedFunc(
							currk.Cols(),
							func(c1, c2 builder.Col) int {
								return cmp.Compare(c1.Name(), c2.Name())
							},
						)),
					),
				)
				if !strings.EqualFold(nextkdef, currkdef) {
					dd.do(ACT_DROP_IDX, nextk.Name(), d.DropIndex(nextk))
					dd.do(ACT_CREATE_IDX, nextk.Name(), d.AddIndex(nextk))
				}
			}
		}
	}

	for currk := range curr.Keys() {
		dropped := false
		for c := range currk.Cols() {
			if act, ok := dd.columns[c.Name()]; ok && act == ACT_DROP_COL {
				dropped = true
				break
			}
		}
		// when column dropped
		if dropped {
			dd.do(ACT_DROP_IDX, currk.Name(), d.DropIndex(currk))
			continue
		}
		// when index dropped
		if next.K(currk.Name()) == nil {
			dd.do(ACT_DROP_IDX, currk.Name(), d.DropIndex(currk))
		}
	}
	return dd
}
