---
name: SQLx
description:
  - 说明如何在其他项目中引入并使用 `github.com/xoctopus/sqlx` 的主要包与集成入口
  - 当 agent 需要把 sqlx 作为依赖接入宿主项目时使用.
---

# SQLx

- 包地图: [references/sqlx-package-map.spec.md](references/sqlx-package-map.spec.md),
- 选型: [sqlx-howto-guideline.md](references/sqlx-howto-guideline.md)
- 代码生成: [sqlx-generator.spec.md](references/sqlx-generator.spec.md)
- 模型定义: [sqlx-model-define-guideline.md](references/sqlx-model-define-guideline.md)
- 基于目录的表迁移: `migrator.Migrate`

## 构建SQL片段

所有语句都是 `frag.Fragment`. 条件用 typed 列:

- 原始片段

```go
f := frag.Query("SELECT * FROM t_user WHERE f_id = ?", 1)
```

- 手动拼装

参见 `github.com/xoctopus/sqlx/pkg/builder/expr_test.go`

- 已生成模型中间代码

拿 `example/models.User` 举例

```go
builder.And(
	TUser.UserID.AsCond(builder.Eq(uid)),
	TUser.Name.AsCond(builder.Like("%x%")),
)
```

- 更多用法. 参考 `go doc github.com/xoctopus/sqlx/pkg/builder`


## 执行SQL

### 打开连接并注入 Session:

```go
a, err := session.Open(ctx, dsn) // 或 adaptor endpoint
s := session.New(a, "main")
// s := session.NewReadonly(rw, ro, "main")

ctx = session.With(ctx, s)
ctx = session.WithModel(ctx, &User{}, s) // 按表名绑定, 供 MustFor(ctx, TUser)
```

### 按模型取 `Session` 后 Exec / Query / Tx:

```go
s := session.MustFor(ctx, TUser)

// 写
_, err = s.Exec(ctx, builder.Insert().Into(TUser).Values(...))
// 或 s.Adaptor().Exec(ctx, f) — 生成代码惯用写法

// 读
rows, err := s.Query(ctx, builder.Select(nil).From(TUser, ...))
defer rows.Close()
err = helper.Scan(ctx, rows, &dst) // dst: *Model / *[]Model / 标量

// 一步查询+扫描
err = helper.QueryAndScan(ctx, s.Adaptor(), f, &dst)

// 事务 (ctx 内可嵌套感知 InTx)
err = s.Tx(ctx, func(ctx context.Context) error {
	if err := m.UpdateByID(ctx); err != nil {
		return err
	}
	return m.FetchByID(ctx)
})

// 只读连接 (NewReadonly 时)
ro := s.Adaptor(session.ReadOnly())
```

### 原始片段也可直接执行:

```go
f := frag.Query("SELECT * FROM t_user WHERE f_id = ?", 1)
_, err = s.Exec(ctx, f)
// 调试: q, args := frag.Collect(ctx, f)
```
