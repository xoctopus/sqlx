package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/xoctopus/genx/pkg/genx"
	"github.com/xoctopus/x/misc/must"

	_ "github.com/xoctopus/sqlx/devpkg/sqlx/v1"
)

func main() {
	entry := filepath.Join(must.NoErrorV(os.Getwd()), "testdata")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{entry},
	})

	if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
		panic(err)
	}
}
