# 包地图

本文描述 `sqlx` 代码分层和功能描述

## 公开包

| pkg               | 场景                     |
|-------------------|--------------------------|
| `adaptors`        | 底层驱动侧载注册         |
| `builder`         | 表结构与 SQL 语句构建    |
| `builder/modeled` | 模型泛型的表 / 列作用域  |
| `errors`          | SQL 领域错误             |
| `frag`            | SQL 片段与参数展开       |
| `frag/testutil`   | Fragment 测试 Matcher    |
| `session`         | 目录与数据库执行绑定     |
| `migrator`        | 基于目录的结构迁移       |
| `types`           | 常见 SQL 类型封装        |
| `types/sqlops`    | 操作时间字段 embed       |
| `types/sqltime`   | 时间底层表示             |
| `helper`          | 插入展开与 Scan 辅助     |

使用 `go doc <package>` 获取更详细的包职能和关键接口

## 内部包 (勿作为公开 API)

业务与 skill 默认只依赖 `pkg/*`. 下列包供实现使用, 外部一般不直接 import:

| pkg                                            | 职责                     |
|------------------------------------------------|--------------------------|
| `internal/sql/adaptor`                         | Adaptor / Dialect / 注册 |
| `internal/sql/adaptor/mysql` (postgres/sqlite) | 各驱动实现               |
| `internal/sql/loggingdriver`                   | 插值与日志 driver 包装   |
| `internal/sql/scanner`                         | Rows 扫描 (供 `helper`)  |
| `internal/def`                                 | 列 / 键定义解析          |
| `internal/structs`                             | 结构体与表列映射         |
| `internal/diff`                                | Catalog diff / 迁移模式  |
| `pkg/migrator/internal`                        | 迁移 SQL meta 生成       |
| `pkg/migrator/models`                          | 迁移元数据表模型         |

使用 `go doc <package>` 获取更详细的包职能和关键接口
