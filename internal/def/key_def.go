package def

import "strings"

// ParseKeyDef parses key define
// eg:
//
//	| Kind         | Name[,Using]       | Field[,Option]             |
//	| :---         | :---               | :----                      |
//	| idx          | idx_name,BTREE     | Name                       |
//	| index        | idx_name,GIST      | Geo,gist_trgm_ops          |
//	| unique_index | idx_name           | OrgID,NULLS,FIRST;MemberID |
//	| u_idx        | idx_name           | OrgID;MemberID,NULLS,FIRST |
//	| primary      |                    | ID                         |
//	| pk           |                    | ID                         |
func ParseKeyDef(def string) *KeyDefine {
	parts := strings.Fields(def)

	d := &KeyDefine{}
	switch parts[0] {
	case "idx", "index":
		if len(parts) != 3 {
			return nil
		}
		d.Kind = KEY_KIND__INDEX
		d.Name, d.Using = ResolveIndexNameAndUsing(parts[1])
	case "unique_index", "u_idx", "uidx", "ui":
		if len(parts) != 3 {
			return nil
		}
		d.Kind = KEY_KIND__UNIQUE_INDEX
		d.Name, d.Using = ResolveIndexNameAndUsing(parts[1])
	case "primary", "pk":
		if len(parts) != 2 {
			return nil
		}
		d.Kind = KEY_KIND__PRIMARY
		d.Name = "primary"
	default:
		return nil
	}
	if d.Name == "" && d.Kind != KEY_KIND__PRIMARY {
		return nil
	}

	d.Options = ResolveKeyColumnOptions(parts[len(parts)-1])
	if len(d.Options) == 0 {
		return nil
	}

	return d
}

func ResolveIndexNameAndUsing(s string) (name string, using string) {
	parts := strings.Split(s, ",")
	name = parts[0]
	if len(parts[1:]) > 0 {
		using = parts[1]
	}
	return
}

type KeyKind int8

const (
	KEY_KIND__INDEX KeyKind = iota + 1
	KEY_KIND__UNIQUE_INDEX
	KEY_KIND__PRIMARY
)

type KeyDefine struct {
	Kind    KeyKind
	Name    string
	Using   string
	Comment string
	Options []KeyColumnOption
}

func (d *KeyDefine) OptionsFieldNames() []string {
	names := make([]string, len(d.Options))
	for i, opt := range d.Options {
		names[i] = opt.FieldName
	}
	return names
}

func (d *KeyDefine) OptionsStrings() []string {
	ss := make([]string, len(d.Options))
	for i, opt := range d.Options {
		ss[i] = opt.String()
	}
	return ss
}

func ResolveKeyColumnOptions(s string) (options []KeyColumnOption) {
	fields := strings.Split(s, ";")
	for _, field := range fields {
		if parts := strings.Split(field, ","); len(parts) > 0 {
			option := KeyColumnOption{
				FieldName: parts[0],
				Options:   parts[1:],
			}
			if option.FieldName == "" {
				continue
			}
			options = append(options, option)
		}
	}
	return
}

func ResolveKeyColumnOptionsFromStrings(ss ...string) (options []KeyColumnOption) {
	for _, s := range ss {
		options = append(options, ResolveKeyColumnOptions(s)...)
	}
	return
}

func KeyColumnOptionByNames(names ...string) []KeyColumnOption {
	options := make([]KeyColumnOption, len(names))
	for i := range names {
		options[i].FieldName = names[i]
	}
	return options
}

type KeyColumnOption struct {
	FieldName string
	Options   []string
}

func (o *KeyColumnOption) String() string {
	if len(o.Options) == 0 {
		return o.FieldName
	}
	return o.FieldName + "," + strings.Join(o.Options, ",")
}
