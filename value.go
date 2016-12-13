package objconv

import (
	"reflect"
	"unsafe"
)

// Type is an enumeration that define the type of a Value.
type Type int

const (
	NilType Type = iota
	BoolType
	IntType
	UintType
	FloatType
	StringType
	BytesType
	TimeType
	DurationType
	ArrayType
	MapType
)

func (t Type) String() string {
	switch t {
	case NilType:
		return "nil"
	case BoolType:
		return "bool"
	case IntType:
		return "int"
	case UintType:
		return "uint"
	case FloatType:
		return "float"
	case StringType:
		return "string"
	case BytesType:
		return "bytes"
	case TimeType:
		return "time"
	case DurationType:
		return "duration"
	case ArrayType:
		return "array"
	case MapType:
		return "map"
	default:
		return "<type>"
	}
}

/*
type Value struct {
	t Type
	x bool
	i int64
	n uint64
	f float64
	s string
	b []byte
	a Array
	m Map
}

func (v *Value) Type() Type {
	return v.t
}

func (v *Value) SetBool(x bool) {
	v.t = BoolType
	v.x = x
}

func (v *Value) SetInt(i int64) {
	v.t = IntType
	v.i = i
}

func (v *Value) SetUint(u uint64) {
	v.t = UintType
	v.u = u
}

func (v *Value) SetFloat(f float64) {
	v.t = FloatType
	v.f = f
}

func (v *Value) SetString(s string) {
	v.t = StringType
	v.s = s
}

func (v *Value) SetBytes(b []byte) {
	v.t = BytesType
	v.b = b
}

func (v *Value) SetTime(t time.Time) {

}

func (v *Value) SetDuration(d time.Duration) {

}

func (v *Value) SetArray(a Array) {
	v.t = ArrayType
	v.a = a
}

func (v *Value) SetMap(m Map) {
	v.t = MapType
	v.m = m
}

func (v *Value) check(t Type) {
	if v.t != t {
		panic("objconv: type mismatch while reading value: " + t.String() + " != " + v.t.String())
	}
}
*/

// IsEmptyValue returns true if the value given as argument would be considered
// empty by the standard library packages, and therefore not serialized if
// `omitempty` is set on a struct field with this value.
func IsEmptyValue(v interface{}) bool {
	return isEmptyValue(reflect.ValueOf(v))
}

// Based on https://golang.org/src/encoding/json/encode.go?h=isEmpty
func isEmptyValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true // nil interface{}
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
	case reflect.Interface, reflect.Ptr, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.UnsafePointer:
		return unsafe.Pointer(v.Pointer()) == nil
	}
	return false
}

// IsZeroValue returns true if the value given as argument is the zero-value of
// the type of v.
func IsZeroValue(v interface{}) bool {
	return isZeroValue(reflect.ValueOf(v))
}

func isZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true // nil interface{}
	}
	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.Len() == 0
	case reflect.UnsafePointer:
		return unsafe.Pointer(v.Pointer()) == nil
	case reflect.Array:
		return isZeroArray(v)
	case reflect.Struct:
		return isZeroStruct(v)
	}
	return false
}

func isZeroArray(v reflect.Value) bool {
	for i, n := 0, v.Len(); i != n; i++ {
		if !isZeroValue(v.Index(i)) {
			return false
		}
	}
	return true
}

func isZeroStruct(v reflect.Value) bool {
	s := LookupStruct(v.Type())

	for _, f := range s.Fields {
		if !isZeroValue(v.FieldByIndex(f.Index)) {
			return false
		}
	}

	return true
}

func setValue(v1 reflect.Value, v2 reflect.Value) (err error) {
	t1 := v1.Type()
	t2 := v2.Type()

	switch {
	case t2.AssignableTo(t1):
		v1.Set(v2)

	case t2.ConvertibleTo(t1):
		v1.Set(v2.Convert(t1))

	default:
		err = &TypeConversionError{
			From: t2,
			To:   t1,
		}
	}

	return
}
