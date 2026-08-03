// Package sqlx provides a generator for generating sqlx model related implementations.
//
// # How to Integrate
//
//	import _ "github.com/xoctopus/sqlx/devpkg/sqlx/v1"
//
//	ctx := genx.NewContext(&genx.Args{
//		Entrypoint: []string{entry},
//	})
//	_ = ctx.Execute(context.Background(), genx.Get()...)
//
// # Code Definition Conventions
//
//  1. Add `// +genx:model` comment to your model type.
//  2. Use `@model` annotations in its comment to configure the generator
//     (e.g., `// @model TableName=t_product`, `// @model Register=Catalog`,
//     `// @model pk=ID`, `// @model uidx=...`, `// @model idx=...`).
//  3. Map columns with `db` struct tags
//     (e.g., `db:"name,width=128"`, `db:"f_id,autoinc"`).
//  4. Soft relations use `// @model rel=<ModelType>.<Field>` on the related
//     field comment (does not create foreign keys).
//
// Index fields are separated by `;`; options by `,`
// (e.g., `// @model idx=i_name,BTREE;Name,NULLS,FIRST;DeletedAt`).
//
// Prefer composing reusable embeds such as `types.AutoIncID`, `Rel*`, `Meta`,
// `State`, and operation-time types from `pkg/types` / `pkg/types/sqlops`.
// See `github.com/xoctopus/sqlx/example/models`.
//
// Example:
//
//	var Catalog = builder.NewCatalog()
//
//	// Product 商品
//	// +genx:model
//	// @model TableName=t_product
//	// @model Register=Catalog
//	// @model pk=ID
//	// @model uidx=ui_product_id;ProductID;DeletedAt
//	// @model idx=i_product_name;Name
//	type Product struct {
//		types.AutoIncID
//		RelProduct
//		ProductMeta
//		types.CreationModificationDeletionTime
//	}
//
//	type RelProduct struct {
//		// @model rel=Product.ProductID
//		ProductID ProductID `db:"product_id"`
//	}
package sqlx
