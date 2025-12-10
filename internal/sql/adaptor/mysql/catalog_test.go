package mysql_test

import (
	"slices"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/hack"
	"github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
	"github.com/xoctopus/sqlx/pkg/builder"
)

func TestScanCatalog_Hack(t *testing.T) {
	hack.Check(t)

	catalog, err := mysql.ScanCatalog(Context(t), NewAdaptor(t), "test")
	Expect(t, err, Succeed())
	Expect(t, catalog, NotBeNil[builder.Catalog]())

	tUsers := catalog.T("users")
	Expect(t, tUsers.C("f_id"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_user_id"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_name"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_email"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_created_at"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_updated_at"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_deleted_at"), NotBeNil[builder.Col]())

	k := tUsers.K("ui_email")
	cols := slices.Collect(k.Cols())
	Expect(t, k.IsUnique(), BeTrue())
	Expect(t, cols, HaveLen[[]builder.Col](2))

	k = tUsers.K("i_name")
	Expect(t, k.IsUnique(), BeFalse())
	Expect(t, slices.Collect(k.Cols()), HaveLen[[]builder.Col](2))
}
