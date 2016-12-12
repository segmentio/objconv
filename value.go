package objconv

import (
	"reflect"
	"unsafe"
)

// IsEmptyValue returns true if the value given as argument would be considered
// empty by the standard library packages, and therefore not serialized if
// `omitempty` is set on a struct field with this value.
func IsEmptyValue(v interface{}) bool {
	return isEmptyValue(reflect.ValueOf(v))
}

// Based on https://golang.org/src/encoding/json/encode.go?h=isEmpty
func isEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true // nil empty interface
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.UnsafePointer:
		return unsafe.Pointer(v.Pointer()) == nil
	}
	return false
}
