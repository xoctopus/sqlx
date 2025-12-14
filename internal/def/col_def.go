package def

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
)

func ParseColDef(t typx.Type, tag reflect.StructTag) *ColumnDef {
	d := &ColumnDef{
		Type: typx.Deref(t),
		Tag:  tag,
	}

	flag := reflectx.ParseTag(tag).Get("db")
	if flag == nil {
		return d
	}
	d.ParseDBTag(flag)
	return d
}

// ColumnDef describes source and database model
type ColumnDef struct {
	Type       typx.Type
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

func (d *ColumnDef) ParseDBTag(flag *reflectx.Flag) {
	for o := range flag.Options() {
		switch strings.ToLower(o.Key()) {
		case "null":
			d.Null = true
		case "autoinc":
			d.AutoInc = true
		case "default":
			ov := o.Value()
			d.Default = &ov
		case "width":
			ov := o.Unquoted()
			v, err := strconv.ParseUint(ov, 10, 64)
			must.NoErrorF(err, "invalid width value: %s", ov)
			d.Width = v
		case "precision":
			ov := o.Unquoted()
			v, err := strconv.ParseUint(ov, 10, 64)
			must.NoErrorF(err, "invalid precision value: %s", ov)
			d.Precision = v
		case "onupdate":
			ov := o.Value()
			must.BeTrueF(len(ov) > 0, "missing onupdate value")
			d.OnUpdate = &ov
		case "deprecated":
			// TODO more deprecated actions?
			d.Deprecated = &DeprecatedActions{RenameTo: o.Value()}
		}
	}
}

type DeprecatedActions struct {
	RenameTo string `name:"rename"`
}
