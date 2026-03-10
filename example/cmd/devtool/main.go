package main

import (
	"context"
	"os"
	"path/filepath"

	_ "github.com/xoctopus/genx/devpkg/docx"
	_ "github.com/xoctopus/genx/devpkg/enumx"
	"github.com/xoctopus/genx/pkg/genx"
	"github.com/xoctopus/x/misc/must"

	_ "github.com/xoctopus/sqlx/devpkg/sqlx/v1"
)

func main() {
	cwd := must.NoErrorV(os.Getwd())

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{
			filepath.Join(cwd, "example", "models"),
			filepath.Join(cwd, "example", "enums"),
			filepath.Join(cwd, "testdata"),
			filepath.Join(cwd, "testdata", "v2"),
		},
	})

	if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
		panic(err)
	}
}
