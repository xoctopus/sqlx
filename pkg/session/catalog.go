package session

import (
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"

	"github.com/xoctopus/sqlx/pkg/builder"
)

var catalogs = syncx.NewXmap[string, struct{}]()

func Register(cats ...builder.Catalog) {
	for _, cat := range cats {
		for t := range cat.Tables() {
			k := t.TableName()
			if x, ok := t.(builder.HasSchema); ok {
				k = x.Schema() + "." + k
			}

			loaded, ok := catalogs.LoadOrStore(k, struct{}{})
			must.BeTrueF(!ok, "model %s already registered to %s", loaded)
		}
	}
}
