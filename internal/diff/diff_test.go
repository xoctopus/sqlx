package diff_test

import (
	"context"
	"testing"

	"github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/hack"
	"github.com/xoctopus/sqlx/internal/diff"
	"github.com/xoctopus/sqlx/pkg/builder"
	"github.com/xoctopus/sqlx/pkg/frag/testutil"
	"github.com/xoctopus/sqlx/testdata"
	"github.com/xoctopus/sqlx/testdata/v2"
)

func TestDiff_mysql(t *testing.T) {
	d := hack.NewAdaptor(t, "mysql://root@localhost:13306/test")
	cbg := context.Background()

	t.Run("InitUserV1", func(t *testing.T) {
		actions := diff.Diff(cbg, d.Dialect(), nil, builder.TFrom(&testdata.User{}))

		// q, _ := frag.Collect(cbg, actions)
		// t.Log(q)
		testx.Expect(t, actions, testutil.BeFragment(`CREATE TABLE IF NOT EXISTS t_user (
	f_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	f_user_id BIGINT UNSIGNED NOT NULL,
	f_org_id BIGINT UNSIGNED NOT NULL,
	f_name VARCHAR(127) NOT NULL,
	f_real_name TEXT NOT NULL,
	f_username VARCHAR(255) NOT NULL,
	f_nick_name VARCHAR(127) NOT NULL,
	f_age INT NOT NULL,
	f_gender TINYINT NOT NULL,
	f_asset DECIMAL(32,4) NOT NULL,
	f_created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	f_updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	f_deleted_at DATETIME NOT NULL DEFAULT '0001-01-01 00:00:00',
	PRIMARY KEY (f_id)
);
CREATE UNIQUE INDEX ui_name ON t_user (f_name,f_deleted_at);
CREATE UNIQUE INDEX ui_user_id ON t_user (f_user_id,f_deleted_at);
CREATE INDEX i_age ON t_user (f_age);
CREATE INDEX i_nickname ON t_user (f_nick_name,f_deleted_at);`,
		))
	})

	t.Run("InitUserV2", func(t *testing.T) {
		actions := diff.Diff(cbg, d.Dialect(), nil, builder.TFrom(&v2.User{}))
		// q, _ := frag.Collect(cbg, actions)
		// t.Log(q)
		testx.Expect(t, actions, testutil.BeFragment(`CREATE TABLE IF NOT EXISTS t_user (
	f_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	f_user_id BIGINT UNSIGNED NOT NULL,
	f_org_id BIGINT UNSIGNED NOT NULL,
	f_real_name VARCHAR(255) NOT NULL DEFAULT '',
	f_age TINYINT NOT NULL DEFAULT 0,
	f_gender TINYINT NOT NULL,
	f_desc TEXT NOT NULL,
	f_created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	f_updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (f_id)
);
CREATE UNIQUE INDEX ui_name ON t_user (f_real_name);
CREATE UNIQUE INDEX ui_user_id ON t_user (f_user_id);
CREATE INDEX i_gender ON t_user (f_gender);`,
		))
	})

	t.Run("DiffUserAndUserV2", func(t *testing.T) {
		actions := diff.Diff(cbg, d.Dialect(), builder.TFrom(&testdata.User{}), builder.TFrom(&v2.User{}))

		// q, _ := frag.Collect(cbg, actions)
		// t.Log(q)
		testx.Expect(t, actions, testutil.BeFragment(`ALTER TABLE t_user DROP INDEX i_age;
ALTER TABLE t_user DROP INDEX i_nickname;
ALTER TABLE t_user DROP INDEX ui_name;
ALTER TABLE t_user DROP INDEX ui_user_id;
ALTER TABLE t_user DROP COLUMN f_nick_name;
ALTER TABLE t_user DROP COLUMN f_real_name;
ALTER TABLE t_user DROP COLUMN f_username;
ALTER TABLE t_user RENAME COLUMN f_name TO f_real_name;
ALTER TABLE t_user MODIFY COLUMN f_age TINYINT NOT NULL DEFAULT 0; /* from INT NOT NULL */
ALTER TABLE t_user MODIFY COLUMN f_real_name VARCHAR(255) NOT NULL DEFAULT ''; /* from TEXT NOT NULL */
ALTER TABLE t_user ADD COLUMN f_desc TEXT NOT NULL;
CREATE INDEX i_gender ON t_user (f_gender);
CREATE UNIQUE INDEX ui_name ON t_user (f_real_name);
CREATE UNIQUE INDEX ui_user_id ON t_user (f_user_id);`))
	})
}
