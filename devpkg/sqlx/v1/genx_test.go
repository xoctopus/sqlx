package sqlx_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/xoctopus/genx/pkg/genx"
	"github.com/xoctopus/x/misc/must"

	_ "github.com/xoctopus/sqlx/devpkg/sqlx/v1"
)

func Example() {
	root := filepath.Join(must.NoErrorV(os.Getwd()), "..", "..", "..")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{
			filepath.Join(root, "testdata"),
			filepath.Join(root, "testdata", "v2"),
			filepath.Join(root, "example", "models"),
		},
	})

	if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
		panic(err)
	}

	// Output:
}
