package sqlx_test

import (
	"context"
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/xoctopus/genx/pkg/genx"
	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/contextx"
	"github.com/xoctopus/x/misc/must"
)

func Example() {
	cwd := must.NoErrorV(os.Getwd())

	entry := filepath.Join(cwd, "..", "..", "..", "testdata")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{
			entry,
			filepath.Join(entry, "v2"),
		},
	})

	if err := ctx.Execute(context.Background(), genx.Get()...); err != nil {
		panic(err)
	}

	// Output:
}

type mock struct{}

func (x *mock) New(c genx.Context) genx.Generator {
	return &mock{}
}

func (x *mock) Identifier() string {
	return "model"
}

func (x *mock) Generate(c genx.Context, t types.Type) error {
	target := c.PackageByPath("github.com/xoctopus/sqlx/pkg/types").
		TypeNames().ElementByName("SoftDeletion").Type()

	if typx.NewTType(t).Implements(target) {
		fmt.Println(t.String())
	}
	return nil
}

func TestImpl(t *testing.T) {
	genx.Register(&mock{})

	cwd := must.NoErrorV(os.Getwd())

	entry := filepath.Join(cwd, "..", "..", "..", "example", "models")
	root := filepath.Join(cwd, "..", "..", "..")

	ctx := genx.NewContext(&genx.Args{
		Entrypoint: []string{entry, filepath.Join(root, "pkg", "types")},
	})

	if err := ctx.Execute(
		contextx.Compose(
		// pkgx.CtxWorkdir.Carry(),
		)(context.Background()),
		genx.Get()...,
	); err != nil {
		panic(err)
	}
}
