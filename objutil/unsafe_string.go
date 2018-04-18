package objutil

import (
	"reflect"
	"unsafe"
)

// UnsafeString returns a string that is only safe to use under the following conditions:
// - b points to data on the heap
// - the bytes pointed to by b will not be modified while the returned string exists
// - the returned string will not be stored past the lifetime of b
func UnsafeString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&b[0])),
		Len:  len(b),
	}))
}
