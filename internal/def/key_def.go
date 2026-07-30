package def

import "strings"

// ParseKeyDef parses key define
// eg:
//
//	| Kind | Name[,Using]       | Field[,Option]                |
//	| :--- | :---               | :----                         |
//	| idx  | idx_name,BTREE     | Name                          |
//	| uidx | idx_name           | OrgID;f_member_id,NULLS,FIRST |
//	| pk   |                    | ID                            |
func ParseKeyDef(k, def string) *KeyDefine {
	def = strings.TrimSpace(def)
	if len(def) == 0 {
		return nil
	}

	parts := strings.Split(def, ";")

	d := &KeyDefine{}
	switch k {
	case "idx":
		if len(parts) < 2 {
			return nil
		}
		d.Kind = KEY_KIND__INDEX
		d.Name, d.Using = ResolveIndexNameAndUsing(parts[0])
		parts = parts[1:]
	case "uidx":
		if len(parts) < 2 {
			return nil
		}
		d.Kind = KEY_KIND__UNIQUE_INDEX
		d.Name, d.Using = ResolveIndexNameAndUsing(parts[0])
		parts = parts[1:]
	case "pk":
		d.Kind = KEY_KIND__PRIMARY
		d.Name = "primary"
	default:
		return nil
	}

	if d.Name == "" {
		return nil
	}

	d.Options = ResolveKeyColumnOptionsFromStrings(parts...)
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

func (d *KeyDefine) OptionsNames() []string {
	names := make([]string, len(d.Options))
	for i, opt := range d.Options {
		names[i] = opt.Name
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

func ResolveKeyColumnOptions(s string) KeyColumnOption {
	if parts := strings.Split(s, ","); len(s) > 0 && len(parts) > 0 {
		if len(parts) > 0 {
			return KeyColumnOption{
				Name:    parts[0],
				Options: parts[1:],
			}
		}
	}
	return KeyColumnOption{}
}

func ResolveKeyColumnOptionsFromStrings(ss ...string) (options []KeyColumnOption) {
	for _, s := range ss {
		if option := ResolveKeyColumnOptions(s); len(option.Name) > 0 {
			options = append(options, option)
		}
	}
	return
}

func KeyColumnOptionByNames(names ...string) []KeyColumnOption {
	options := make([]KeyColumnOption, len(names))
	for i := range names {
		options[i].Name = names[i]
	}
	return options
}

type KeyColumnOption struct {
	Name    string // maybe column name or field name
	Options []string
}

func (o *KeyColumnOption) String() string {
	if len(o.Options) == 0 {
		return o.Name
	}
	return o.Name + "," + strings.Join(o.Options, ",")
}
