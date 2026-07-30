package internal_test

import (
	"context"
	"slices"
	"testing"

	. "github.com/xoctopus/x/testx"

	"github.com/xoctopus/sqlx/example/models"
	"github.com/xoctopus/sqlx/hack"
	"github.com/xoctopus/sqlx/pkg/migrator/internal"
)

func TestGenerateDocuments(t *testing.T) {
	a := hack.NewAdaptor(t, "mysql://root@localhost:13306/test_sql_meta?multiStatements=true&interpolateParams=true")

	fragments := slices.Concat(
		internal.GenerateTableDocuments(t.Context(), a, models.Catalog),
		internal.GenerateTableColumnDocuments(t.Context(), a, models.Catalog),
		internal.GenerateTableEnumerationDocument(t.Context(), a, models.Catalog),
	)

	err := a.Tx(t.Context(), func(ctx context.Context) error {
		for _, f := range fragments {
			if _, err := a.Exec(ctx, f); err != nil {
				return err
			}
		}
		return nil
	})
	Expect(t, err, Succeed())
}
