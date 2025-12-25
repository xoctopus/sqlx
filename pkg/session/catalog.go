package session

import (
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/syncx"

	"github.com/xoctopus/sqlx/pkg/builder"
)

var (
	catalogs = syncx.NewXmap[string, string]()
	sessions = syncx.NewXmap[string, struct{}]()
)

func register(session string, catalog builder.Catalog) {
	_, ok := sessions.LoadOrStore(session, struct{}{})
	must.BeTrueF(!ok, "session already registered: %s", session)

	for t := range catalog.Tables() {
		registered, ok := catalogs.LoadOrStore(t.TableName(), session)
		must.BeTrueF(!ok, "model %s already registered to %s", registered)
	}
}
