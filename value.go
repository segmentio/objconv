package objconv

import (
	"encoding"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

// Type is an enumeration that represent all the base types supported by the
// emitters and parsers.
type Type int

const (
	Unknown Type = iota
	Nil
	Bool
	Int
	Uint
	Float
	String
	Bytes
	Time
	Duration
	Error
	Array
	Map
)

// String returns a human readable representation of the type.
func (t Type) String() string {
	switch t {
	case Nil:
		return "nil"
	case Bool:
		return "bool"
	case Int:
		return "int"
	case Uint:
		return "uint"
	case Float:
		return "float"
	case String:
		return "string"
	case Bytes:
		return "bytes"
	case Time:
		return "time"
	case Duration:
		return "duration"
	case Error:
		return "error"
	case Array:
		return "array"
	case Map:
		return "map"
	default:
		return "<type>"
	}
}

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

var (
	zeroCache = make(map[reflect.Type]reflect.Value)
	zeroMutex sync.RWMutex
)

// zeroValueOf and the related cache is used to keep the zero values so they
// don't need to be reallocated everytime they're used.
func zeroValueOf(t reflect.Type) reflect.Value {
	zeroMutex.RLock()
	v, ok := zeroCache[t]
	zeroMutex.RUnlock()

	if !ok {
		v = reflect.Zero(t)
		zeroMutex.Lock()
		zeroCache[t] = v
		zeroMutex.Unlock()
	}

	return v
}

var (
	timeType                 = reflect.TypeOf(time.Time{})
	timePtrType              = reflect.TypeOf((*time.Time)(nil))
	durationType             = reflect.TypeOf(time.Duration(0))
	durationPtrType          = reflect.TypeOf((*time.Duration)(nil))
	errorInterface           = reflect.TypeOf((*error)(nil)).Elem()
	valueEncoderInterface    = reflect.TypeOf((*ValueEncoder)(nil)).Elem()
	valueDecoderInterface    = reflect.TypeOf((*ValueDecoder)(nil)).Elem()
	textMarshalerInterface   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshalerInterface = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	emptyInterface           = reflect.TypeOf((*interface{})(nil)).Elem()
	emptyInterfacePtr        = reflect.TypeOf((*interface{})(nil))
)

func stringNoCopy(b []byte) string {
	n := len(b)
	if n == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&b[0])),
		Len:  n,
	}))
}
