package def

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/xoctopus/typx/pkg/typx"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/reflectx"
)

const (
	// DefineFromCatalog scan from database
	DefineFromCatalog = 1
	// DefineFromUser type implements builder.ColDef interface
	DefineFromUser = 2
	// DefineFromGoType parse depends on Go type
	DefineFromGoType = 3
)

func ParseColDef(t typx.Type, tag reflect.StructTag) *ColumnDef {
	d := &ColumnDef{
		Type: typx.Deref(t),
		Tag:  tag,
	}

	flag := reflectx.ParseTag(
		tag,
		reflectx.WithOptionSplitter('='),
		reflectx.WithExpectFlags("db"),
	).Get("db")
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
	IsUnsigned *bool
	Precision  *uint64
	Width      *uint64
	Default    *string
	OnUpdate   *string
	Null       bool
	AutoInc    bool
	Comment    string
	Desc       []string
	Relation   []string
	Deprecated *DeprecatedActions
	DefineFrom int
}

func (d *ColumnDef) ParseDBTag(flag *reflectx.Flag) {
	for o := range flag.Options() {
		switch strings.ToLower(o.Key()) {
		case "null":
			d.Null = true
		case "autoinc":
			d.AutoInc = true
		case "default":
			// NOTE: here is no special processing is performed during parsing
			// default values; instead, the raw content is retrieved in its entirety.
			// eg:
			// - For a `varchar` with tag: `default='content'`, `'content'` will be retrieved.
			// - For a `datetime` with tag: `default=CURRENT_TIMESTAMP` or `default='2001-01-01'`.
			// consider whether to wrap the default value in single quotes to
			// prevent unintended errors.
			d.Default = new(o.Value())
		case "width":
			ov := o.Unquoted()
			v, err := strconv.ParseUint(ov, 10, 64)
			must.NoErrorF(err, "invalid width value: %s", ov)
			d.Width = &v
		case "precision":
			ov := o.Unquoted()
			v, err := strconv.ParseUint(ov, 10, 64)
			must.NoErrorF(err, "invalid precision value: %s", ov)
			d.Precision = &v
		case "onupdate":
			ov := o.Value()
			must.BeTrueF(len(ov) > 0, "missing onupdate value")
			d.OnUpdate = &ov
		case "unsigned":
			// NOTE: Go's float32 and float64 types are always signed by default.
			// Marking these as 'unsigned' in the database mapping carries an
			// overflow or data-integrity risk, as Go cannot natively enforce
			// unsigned constraints for floating-point numbers at the type level.
			d.IsUnsigned = new(true)
		case "deprecated":
			// TODO more deprecated actions?
			d.Deprecated = &DeprecatedActions{RenameTo: o.Value()}
		}
	}
}

type DeprecatedActions struct {
	RenameTo string `name:"rename"`
}
