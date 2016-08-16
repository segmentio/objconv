package objconv

import "reflect"

// StructField represents a single field of a struct and carries information
// useful to the algorithms of the objconv package.
type StructField struct {
	// The index of the field in the structure.
	Index []int

	// The name of the field in the structure.
	Name string

	// The type of the field in the structure.
	Type reflect.Type

	// The offset of the field within the structure.
	Offset uintptr

	// The struct tags that were set on the field.
	Tags StructTag

	// Set to true when the field is anonymous, false otherwise.
	Anonymous bool

	// Set to true when the field is exported, false otherwise.
	Exported bool
}

func makeStructField(f reflect.StructField) StructField {
	return StructField{
		Index:     f.Index,
		Name:      f.Name,
		Offset:    f.Offset,
		Type:      f.Type,
		Anonymous: f.Anonymous,
		Exported:  len(f.PkgPath) == 0,
		Tags:      ParseStructTag(f.Tag),
	}
}
