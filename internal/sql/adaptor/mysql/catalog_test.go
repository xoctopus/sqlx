package mysql_test

import (
	"slices"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/hack"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag/testutil"
)

func TestScanCatalog_Hack(t *testing.T) {
	hack.Check(t)

	d := NewAdaptor(t)

	catalog, err := d.Catalog(Context(t))

	Expect(t, err, Succeed())
	Expect(t, catalog, NotBeNil[builder.Catalog]())

	tUsers := catalog.T("users")
	Expect(t, tUsers.C("f_id"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_user_id"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_name"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_email"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_token"), NotBeNil[builder.Col]())
	Expect(t, tUsers.C("f_balance"), NotBeNil[builder.Col]())
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

	datatype := d.Dialect().DBType(tUsers.C("f_id").(builder.ColDef).Def())
	Expect(t, datatype, testutil.BeFragment("bigint unsigned NOT NULL AUTO_INCREMENT"))

	datatype = d.Dialect().DBType(tUsers.C("f_user_id").(builder.ColDef).Def())
	Expect(t, datatype, testutil.BeFragment("bigint unsigned NOT NULL"))

	datatype = d.Dialect().DBType(tUsers.C("f_name").(builder.ColDef).Def())
	Expect(t, datatype, testutil.BeFragment("varchar(64) NOT NULL"))

	datatype = d.Dialect().DBType(tUsers.C("f_email").(builder.ColDef).Def())
	Expect(t, datatype, testutil.BeFragment("varchar(255) NOT NULL"))

	datatype = d.Dialect().DBType(tUsers.C("f_token").(builder.ColDef).Def())
	Expect(t, datatype, testutil.BeFragment("varbinary(1024)"))

	datatype = d.Dialect().DBType(tUsers.C("f_balance").(builder.ColDef).Def())
	Expect(t, datatype, testutil.BeFragment("decimal(22,4) NOT NULL"))

	datatype = d.Dialect().DBType(tUsers.C("f_created_at").(builder.ColDef).Def())
	Expect(t, datatype, testutil.BeFragment("bigint unsigned NOT NULL"))

	datatype = d.Dialect().DBType(tUsers.C("f_updated_at").(builder.ColDef).Def())
	Expect(t, datatype, testutil.BeFragment("bigint unsigned NOT NULL"))

	datatype = d.Dialect().DBType(tUsers.C("f_deleted_at").(builder.ColDef).Def())
	Expect(t, datatype, testutil.BeFragment("bigint unsigned NOT NULL DEFAULT 0"))
}
