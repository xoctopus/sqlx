package session

import (
	"github.com/xoctopus/x/syncx"

	"github.com/xoctopus/sqlx/pkg/builder"
)

var catalogs = syncx.NewXmap[string, string]()

func RegisterCatalog(database string, catalog builder.Catalog) {
	for t := range catalog.Tables() {
		catalogs.Store(database+"."+t.TableName(), database)
	}
}
