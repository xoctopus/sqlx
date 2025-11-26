package testdata

// NotStruct
// +genx:model
type NotStruct string

// NotNamed
// +genx:model
type NotNamed = NotStruct

// NoIndexDef
// +genx:model
type NoIndexDef struct {
}
