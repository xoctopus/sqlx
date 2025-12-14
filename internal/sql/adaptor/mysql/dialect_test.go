package mysql_test

import (
	"context"
	"database/sql/driver"
	"errors"
	"reflect"
	"testing"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/xoctopus/typx/pkg/typx"
	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/hack"
	"github.com/xoctopus/sqlx/internal/def"
	"github.com/xoctopus/sqlx/internal/sql/adaptor/mysql"
	"github.com/xoctopus/sqlx/internal/sql/loggingdriver"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag"
	. "github.com/xoctopus/sqlx/pkg/frag/testutil"
	"github.com/xoctopus/sqlx/pkg/types"
)

func TestDialect_Hack(t *testing.T) {
	hack.Check(t)

	d := NewAdaptor(t).Dialect()
	Expect(t, d.CreateSchema("any"), BeFragment("CREATE DATABASE IF NOT EXISTS any;"))
	Expect(t, d.SwitchSchema("any"), BeFragment("USE any;"))
	Expect(t, d.DropTable(builder.T("t_table")), BeFragment("DROP TABLE IF EXISTS t_table;"))
	Expect(t, d.TruncateTable(builder.T("t_table")), BeFragment("TRUNCATE TABLE t_table;"))

	cIDDef := def.ParseColDef(typx.NewRType(reflect.TypeFor[types.ID]()), `db:",autoinc"`)
	cIDDef.Comment = "PK ID"

	tab := builder.T("demo")
	tab.(builder.ColsManager).AddCol(
		builder.C("f_id", builder.WithColDef(cIDDef)),
		builder.C("f_name", builder.WithColDefOf("", `db:",width=255,default=''"`)),
		builder.C("f_org_id", builder.WithColDefOf(types.ID(0), `db:",default=0"`)),
		builder.C("f_bool", builder.WithColDefOf(new(bool), "")),
		builder.C("f_tinyint", builder.WithColDefOf(int8(0), "")),
		builder.C("f_tinyint_unsigned", builder.WithColDefOf(uint8(0), "")),
		builder.C("f_smallint", builder.WithColDefOf(int16(0), "")),
		builder.C("f_smallint_unsigned", builder.WithColDefOf(uint16(0), "")),
		builder.C("f_int", builder.WithColDefOf(int(0), "")),
		builder.C("f_int_unsigned", builder.WithColDefOf(uint32(0), "")),
		builder.C("f_bigint", builder.WithColDefOf(int64(0), "")),
		builder.C("f_bigint_unsigned", builder.WithColDefOf(uint64(0), "")),
		builder.C("f_time", builder.WithColDefOf(time.Time{}, "")),
		builder.C("f_desc", builder.WithColDefOf("", `db:",default=('')"`)),
		builder.C("f_float", builder.WithColDefOf(float32(0), "")),
		builder.C("f_double", builder.WithColDefOf(float64(0), "")),
		builder.C("f_created_at", builder.WithColDefOf(types.Datetime{}, `db:",precision=3,default=CURRENT_TIMESTAMP(3)"`)),
		builder.C("f_updated_at", builder.WithColDefOf(types.Datetime{}, `db:",precision=3,default=CURRENT_TIMESTAMP(3),onupdate=CURRENT_TIMESTAMP(3)"`)),
		builder.C("f_deleted_at", builder.WithColDefOf(types.Datetime{}, `db:",precision=3,default='0001-01-01 00:00:00'"`)),
		builder.C("f_deprecated", builder.WithColDefOf("", `db:",deprecated"`)),
	)
	tab.(builder.KeysManager).AddKey(
		builder.PK(builder.ColsOf(tab.C("f_id"))),
		builder.UK("ui_name", builder.ColsOf(tab.C("f_name"), tab.C("f_deleted_at"))),
		builder.K("i_org_id", builder.ColsOf(tab.C("f_org_id")), builder.WithKeyMethod("BTREE")),
	)

	q, args := frag.Collect(context.Background(), frag.Compose("\n", d.CreateTableIfNotExists(tab)...))
	named := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		named[i].Value = arg
	}

	q, _ = loggingdriver.DefaultInterpolate(q, named)
	Expect(t, q, Equal(`CREATE TABLE IF NOT EXISTS demo (
    f_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'PK ID',
    f_name VARCHAR(255) NOT NULL DEFAULT '',
    f_org_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    f_bool BOOLEAN NOT NULL,
    f_tinyint TINYINT NOT NULL,
    f_tinyint_unsigned TINYINT UNSIGNED NOT NULL,
    f_smallint SMALLINT NOT NULL,
    f_smallint_unsigned SMALLINT UNSIGNED NOT NULL,
    f_int INT NOT NULL,
    f_int_unsigned INT UNSIGNED NOT NULL,
    f_bigint BIGINT NOT NULL,
    f_bigint_unsigned BIGINT UNSIGNED NOT NULL,
    f_time DATETIME NOT NULL,
    f_desc TEXT NOT NULL DEFAULT (''),
    f_float FLOAT NOT NULL,
    f_double DOUBLE PRECISION NOT NULL,
    f_created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    f_updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    f_deleted_at DATETIME(3) NOT NULL DEFAULT '0001-01-01 00:00:00',
    PRIMARY KEY (f_id)
);
CREATE UNIQUE INDEX ui_name ON demo (f_name,f_deleted_at);
CREATE INDEX i_org_id ON demo (f_org_id) USING BTREE;`))

	c1 := builder.C("f_appended", builder.WithColDefOf("", `db:",width=255,default=''"`)).Of(tab)
	Expect(t, d.AddColumn(c1), BeFragment("ALTER TABLE demo ADD COLUMN f_appended VARCHAR(255) NOT NULL DEFAULT '';"))
	c2 := builder.C("f_renamed").Of(tab)
	Expect(t, d.RenameColumn(c1, c2), BeFragment("ALTER TABLE demo RENAME COLUMN f_appended TO f_renamed;"))
	Expect(t, d.DropColumn(c2), BeFragment("ALTER TABLE demo DROP COLUMN f_renamed;"))

	pk := builder.PK(builder.ColsOf(tab.C("f_id"))).Of(tab)
	Expect(t, d.DropIndex(pk), BeFragment("ALTER TABLE demo DROP PRIMARY KEY;"))
	Expect(t, d.AddIndex(pk), BeFragment("ALTER TABLE demo ADD PRIMARY KEY (f_id);"))

	uk := builder.UK("ui_idx", builder.ColsOf(tab.C("f_id"), tab.C("f_name"))).Of(tab)
	Expect(t, d.DropIndex(uk), BeFragment("ALTER TABLE demo DROP INDEX ui_idx;"))
	Expect(t, d.AddIndex(uk), BeFragment("CREATE UNIQUE INDEX ui_idx ON demo (f_id,f_name);"))

	err := &mysqldriver.MySQLError{Number: 1049}
	Expect(t, d.IsUnknownDatabaseError(err), BeTrue())
	err.Number = 1062
	Expect(t, d.IsConflictError(err), BeTrue())

	Expect(t, mysql.UnwrapError(err), NotBeNil[*mysqldriver.MySQLError]())
	Expect(t, mysql.UnwrapError(errors.New("any")), BeNil[*mysqldriver.MySQLError]())
}
