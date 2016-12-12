package objconv

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestIsEmptyTrue(t *testing.T) {
	tests := []interface{}{
		nil,

		false,

		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),

		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		uintptr(0),

		float32(0),
		float64(0),

		"",
		[]byte(nil),
		[]int{},
		[0]int{},

		(map[string]int)(nil),
		map[string]int{},

		(*int)(nil),
		unsafe.Pointer(nil),
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test), func(t *testing.T) {
			if !IsEmptyValue(test) {
				t.Errorf("%T, %#v should be an empty value", test, test)
			}
		})
	}
}
