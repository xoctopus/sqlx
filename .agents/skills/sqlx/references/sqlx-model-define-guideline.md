# 模型定义规范

本文描述以下规范和建议
- 模型定义
  + 代码定义模型的规范和建议
  + 针对已有表结构定义模型代码
- 模型字段标注

## 数据库模型定义

### 顶层约定

- 一个数据库表在代码里是一个结构体
- 结构体的每个字段对应一个表的字段
- **强烈建议** 区分数据库 **自增主键** 和 **业务关联键**. 要严格区分什么字段是数据库生成的, 什么字段是业务生成的.
- **强烈建议** 类型复用. 组合稳定的模型定义.

以上约束参考 `github.com/xoctopus/sqlx/example/models` 的定义实现

举例:
- 统一的自增字段 `types.Serial`
- 业务/资源的唯一主键 `models.RelOrder` `models.RelProduct` `models.RelUser`
- 业务/资源的相对静态数据段 `models.OrderMeta` `models.ProductMeta`
- 业务/资源的经常变更数据段 `models.OrderState` `models.ShipmentState`
- 业务/资源的记录操作变更段 包括 `types`, `sqlops`, 包内时间相关定义

文档查阅:

- `go doc github.com/xoctopus/sqlx/types`
- `go doc github.com/xoctopus/sqlx/types/sqlops`

### 数据类型

| 数据类型     | 说明                                                               |
|--------------|--------------------------------------------------------------------|
| 自增主键     | 参考  `types.Serial` 或自行实现.                                   |
| 业务主键     | 底层为 `整型` 或 `字符串`. 生成可以依赖 `uuid` `雪花ID` 等额外方案 |
| 记录操作时间 | 参考 `types`, `sqlops` 和 `sqltime`                                |
| JSON数据     | 参考 `types.JSONArray`, `types.JSONObject`                         |
| 枚举数据     | 联动 `.agents/skills/genx/references/generator-enumx.spec.md`      |
| 文本/流数据  | 参考 `types.Text` `types.Blob`                                     |
| 浮点数据     | 参考 `types.Decimal`                                               |
| 布尔数据     | 参考 `types.Bool`                                                  |

无论是定义新的模型还是从已有数据库模型上定义, 都需要遵循上述规则.

## 模型字段标注

举例: `db:"<列名>,width=<宽度>,precision=<精度>,default=<默认值>,null,autoinc,onupdate=<更新子句>"`

1. 列名 `db:"f_name"`
2. 自增 `db:"f_id,autoinc"`
3. `VARCHAR` 宽度: `db:"name,width=128"` => `VARCHAR(128)`
4. `DECIMAL` 精度/宽度: `db:"name,width=22,precision=4"` => `DECIMAL(22,4)`
5. `DATETIME` 精度: `db:"created_at,precision=3"` => `DATETIME(3)`
6. 默认值(函数): `db:"created_at,default=CURRENT_TIMESTAMP(3)"` => `DEFAULT CURRENT_TIMESTAMP(3)`
7. 默认值(值): `db:"created_at,default='1970-01-01 00:00:00'"` => `DEFAULT '1970-01-01 00:00:00'`
8. `ON UPDATE` 触发更新: `db:"updated_at,onupdate=CURRENT_TIMESTAMP(3)"` => `ON UPDATE CURRENT_TIMESTAMP(3)`
9. `NULL` 是否允许为空: 如果没有该标签 => `NOT NULL`; 否则不会有非空约束 `db:"...,null"`
