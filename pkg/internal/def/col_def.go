package def

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
)

func ParseColumnDef(t reflect.Type, key string, tag reflect.StructTag) *ColumnDef {
	d := &ColumnDef{
		Type: reflectx.Deref(t),
		Tag:  tag,
	}

	flag := reflectx.ParseTag(tag).Get(key)
	if flag == nil {
		return d
	}

	for o := range flag.Options() {
		ov := o.Unquoted()
		switch strings.ToLower(o.Key()) {
		case "null":
			d.Null = true
		case "autoinc":
			d.AutoInc = true
		case "default":
			must.BeTrueF(len(ov) > 0, "missing default value")
			d.Default = &ov
		case "width":
			v, err := strconv.ParseUint(ov, 10, 64)
			must.NoErrorF(err, "invalid width value: %s", ov)
			d.Width = v
		case "precision":
			v, err := strconv.ParseUint(ov, 10, 64)
			must.NoErrorF(err, "invalid precision value: %s", ov)
			d.Precision = v
		case "onupdate":
			must.BeTrueF(len(ov) > 0, "missing onupdate value")
			d.OnUpdate = &ov
		case "deprecated":
			// TODO more deprecated actions?
			d.Deprecated = &DeprecatedActions{RenameTo: ov}
		}
	}
	return d
}

// ColumnDef describes source and database model
type ColumnDef struct {
	Type       reflect.Type
	Tag        reflect.StructTag
	DataType   string
	Width      uint64
	Precision  uint64
	Default    *string
	OnUpdate   *string
	Null       bool
	AutoInc    bool
	Comment    string
	Desc       []string
	Relation   []string
	Deprecated *DeprecatedActions
}

type DeprecatedActions struct {
	RenameTo string `name:"rename"`
}

/*
TODO adapt gorm tags
func parseFieldIndexes(field *Field) (indexes []Index, err error) {
	for _, value := range strings.Split(field.Tag.Get("gorm"), ";") {
		if value != "" {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if k == "INDEX" || k == "UNIQUEINDEX" {
				var (
					name       string
					tag        = strings.Join(v[1:], ":")
					idx        = strings.IndexByte(tag, ',')
					tagSetting = strings.Join(strings.Split(tag, ",")[1:], ",")
					settings   = ParseTagSetting(tagSetting, ",")
					length, _  = strconv.Atoi(settings["LENGTH"])
				)

				if idx == -1 {
					idx = len(tag)
				}

				name = tag[0:idx]
				if name == "" {
					subName := field.Name
					const key = "COMPOSITE"
					if composite, found := settings[key]; found {
						if len(composite) == 0 || composite == key {
							err = fmt.Errorf(
								"the composite tag of %s.%s cannot be empty",
								field.Schema.Name,
								field.Name)
							return
						}
						subName = composite
					}
					name = field.Schema.namer.IndexName(
						field.Schema.Table, subName)
				}

				if (k == "UNIQUEINDEX") || settings["UNIQUE"] != "" {
					settings["CLASS"] = "UNIQUE"
				}

				priority, err := strconv.Atoi(settings["PRIORITY"])
				if err != nil {
					priority = 10
				}

				indexes = append(indexes, Index{
					Name:    name,
					Class:   settings["CLASS"],
					Type:    settings["TYPE"],
					Where:   settings["WHERE"],
					Comment: settings["COMMENT"],
					Option:  settings["OPTION"],
					Fields: []IndexOption{{
						Field:      field,
						Expression: settings["EXPRESSION"],
						Sort:       settings["SORT"],
						Collate:    settings["COLLATE"],
						Length:     length,
						Priority:   priority,
					}},
				})
			}
		}
	}

	err = nil
	return
}
*/
