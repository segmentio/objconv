package objconv

import "reflect"

// StructField represents a single field of a struct and carries information
// useful to the algorithms of the objconv package.
type StructField struct {
	// The index of the field in the structure.
	Index []int

	// The name of the field in the structure.
	Name string

	// Omitempty is set to true when the field should be omitted if it has an
	// empty value.
	Omitempty bool

	// Omitzero is set to true when the field should be omitted if it has a zero
	// value.
	Omitzero bool
}

func makeStructField(f reflect.StructField) StructField {
	t := ParseTag(f.Tag.Get("objconv"))
	s := StructField{
		Index:     f.Index,
		Name:      f.Name,
		Omitempty: t.Omitempty,
		Omitzero:  t.Omitzero,
	}

	if len(t.Name) != 0 {
		s.Name = t.Name
	}

	return s
}

func omit(f StructField, v reflect.Value) bool {
	return (f.Omitempty && isEmptyValue(v)) || (f.Omitzero && isZeroValue(v))
}
