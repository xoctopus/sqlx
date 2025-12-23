@def T
@def ModeledKeyInitList
@def ModeledColInitList
@def modeled.M
@def Register
--Init
var T#T# *t#T#

func init() {
	m := #modeled.M#()

	T#T# = &t#T#{
		Table: m,

		I: i#T#{
			#ModeledKeyInitList#
		},

		#ModeledColInitList#
	}

	#Register#
}

@def T
@def ModeledKeyDefList
--iTable
// i#T# includes all modeled indexes of #T#
type i#T# struct {
	#ModeledKeyDefList#
}

@def T
@def modeled.Table
@def ModeledColDefList
--tTable
// t#T# includes modeled table, indexes and column list.
type t#T# struct {
	#modeled.Table#
	I i#T#

	#ModeledColDefList#
}

@def T
@def builder.Model
--New
// New creates a new #T#
func (t *t#T#) New() #builder.Model# {
	return &#T#{}
}

@def T
@def builder.Col
@def builder.ColsOf
@def builder.Assignment
@def reflect.ValueOf
@def builder.GetColDef
@def builder.ColumnsAndValues
// AssignmentFor returns assignment by m with expects columns
func (t *t#T#) AssignmentFor(m *#T#, expects ...#builder.Col#) #builder.Assignment# {
	cols := t.Pick()
	if len(expects) > 0 {
		cols = #builder.ColsOf#(expects...)
	}
	vals := make([]any, 0, cols.Len())
	rv := #reflect.ValueOf#(m).Elem()

	for c := range cols.Cols() {
		if !#builder.GetColDef#(c).AutoInc {
			vals = append(vals, rv.FieldByName(c.FieldName()).Interface())
		}
	}

	return #builder.ColumnsAndValues#(cols, vals...)
}


@def T
@def TableName
--TableName
// TableName returns database table name of #T#
func (m #T#) TableName() string {
	return #TableName#
}

@def T
@def TableDesc
--TableDesc
// TableDesc returns descriptions of #T#
func (m #T#) TableDesc() []string {
	return []string{
		#TableDesc#
	}
}

@def T
@def PrimaryColList
--PrimaryKey
// PrimaryKey returns column list of #T#'s primary key
func (m #T#) PrimaryKey() []string {
	return []string{
		#PrimaryColList#
	}
}

@def T
@def IndexList
--Indexes
// Indexes returns index list of #T#
func (m #T#) Indexes() map[string][]string {
	return map[string][]string{
		#IndexList#
	}
}

@def T
@def UniqueIndexList
--UniqueIndexes
// UniqueIndexes returns unique index list of #T#
func (m #T#) UniqueIndexes() map[string][]string {
	return map[string][]string{
		#UniqueIndexList#
	}
}

