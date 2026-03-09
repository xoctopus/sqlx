package session

import (
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"

	"github.com/xoctopus/sqlx/pkg/builder"
)

var catalogs = syncx.NewXmap[string, string]()

func Register(session string, cats ...builder.Catalog) {
	for _, cat := range cats {
		for t := range cat.Tables() {
			k := session + "." + t.TableName()
			loaded, ok := catalogs.LoadOrStore(k, session)
			must.BeTrueF(!ok, "model %s already registered to %s", loaded)
		}
	}
}
