# 选型

本文描述 SQLx 框架的选型

## 选型

| 目标                   | 组合                                            |
|------------------------|-------------------------------------------------|
| 减少手写 SQL           | `frag` + `builder`                              |
| 从已有模型扫表信息     | `builder.TFrom` / `builder.NewCatalog` + `Add`  |
| 执行 SQL               | `session` + `frag` + `builder`                  |
| 按模型拼 insert / scan | `helper` + `builder/modeled`                    |
| 数据库结构迁移         | `builder` + `session`/`adaptor` + `migrator`    |
| 模型时间 / 软删字段    | `types` + `types/sqlops` (+ `sqltime` 如需底层) |
| 判定业务 SQL 错误      | `errors.IsErrNotFound` / `IsErrConflict` / ...  |

## 关键入口

- 表定义: `builder.T` / `builder.TFrom` / `builder.C` / `builder.Key*`
- 语句: `builder.Select` / `Insert` / `Update` / `Delete` / `With` / `Where` / `Join` / Limit 等 Addition
- 片段: `frag.Query` / `Literal` / `Compose` / `Args`
- 会话: `session.New` / `NewReadonly`, `Session.Exec` / `Query` / `Tx`, `session.Register`
- 迁移: `migrator.Migrate`, `migrator.DIFF_MODE_DRY_RUN` / `DIFF_MODE_CREATE_TABLE`
- 扫描: `helper.Scan` / `QueryAndScan`, `helper.CVsForInsertion`


