# 生成器集成指南

本文描述
- 模型代码生成器接入
- 模型代码生成标注规范

## 接入

按照 `github.com/xoctopus/sqlx/devpkg/sqlx/v1` 约定接入
- `go doc github.com/xoctopus/sqlx/devpkg/sqlx/v1`


## 模型代码生成标注规范

| 规范          | 说明                                                               |
|:--------------|:-------------------------------------------------------------------|
| 注释规范      | 联动 `.agents/skills/genx/references/generator-docx.spec.md`       |
| 生成器指令    | `// +genx:model`                                                   |
| 定义表名      | `// @model TableName=<表名字符串>`                                 |
| 注册模型目录  | `// @model Register=<目录变量名>`                                  |
| 定义主键      | `// @model pk=<字段名>`                                            |
| 定义一般索引  | `// @model idx=<索引名[,选项]>;<字段名[,选项]>[;<字段名[,选项]>]`  |
| 定义唯一索引  | `// @model uidx=<索引名[,选项]>;<字段名[,选项]>[;<字段名[,选项]>]` |
| 关联标注      | `// @model rel=<模型类型>.<字段名>`.                               |

补充:
- 索引定义描述: `go doc github.com/xoctopus/sqlx/internal/def.ParseKeyDef`
- 关联标注并不会创建外键约束但是会通过类型复用标注关联关系. 参考 `example/models.RelProduct`
- 以上约束参见 `github.com/xoctopus/sqlx/example/models` 的定义实现

