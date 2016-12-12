package resp

import (
	"reflect"
	"testing"
	"time"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/test"
)

func TestParser(t *testing.T) {
	tests := []struct {
		s string
		v interface{}
	}{
		// nil
		{"$-1\r\n", nil},
		{"*-1\r\n", nil},

		// bool
		{":0\r\n", false},
		{":1\r\n", true},
		{":2\r\n", true},
		{":18446744073709551615\r\n", true},

		// int
		{":42\r\n", int(42)},
		{":42\r\n", int(42)},
		{":42\r\n", int8(42)},
		{":42\r\n", int16(42)},
		{":42\r\n", int32(42)},
		{":42\r\n", int64(42)},

		// uint
		{":42\r\n", uint(42)},
		{":42\r\n", uint8(42)},
		{":42\r\n", uint16(42)},
		{":42\r\n", uint32(42)},
		{":42\r\n", uint64(42)},
		{":42\r\n", uintptr(42)},

		// float
		{":42\r\n", float32(42)},
		{":42\r\n", float64(42)},
		{"+0.5\r\n", float32(0.5)},
		{"+0.5\r\n", float64(0.5)},
		{"$2\r\n42\r\n", float32(42)},
		{"$2\r\n42\r\n", float64(42)},
		{"$3\r\n0.5\r\n", float32(0.5)},
		{"$3\r\n0.5\r\n", float64(0.5)},

		// time
		{"$20\r\n1970-01-01T00:00:00Z\r\n", time.Unix(0, 0).In(time.UTC)},
		{"+1970-01-01T00:00:00Z\r\n", time.Unix(0, 0).In(time.UTC)},
		{"$2\r\n1s\r\n", time.Second},
		{"+1s\r\n", time.Second},

		// string
		{"+\r\n", ""},
		{"+Hello World!\r\n", "Hello World!"},
		{"$0\r\n\r\n", ""},
		{"$12\r\nHello World!\r\n", "Hello World!"},

		// []byte
		{"+\r\n", []byte("")},
		{"+Hello World!\r\n", []byte("Hello World!")},
		{"$0\r\n\r\n", []byte("")},
		{"$12\r\nHello World!\r\n", []byte("Hello World!")},

		// []rune
		{"+\r\n", []rune("")},
		{"+Hello World!\r\n", []rune("Hello World!")},
		{"$0\r\n\r\n", []rune("")},
		{"$12\r\nHello World!\r\n", []rune("Hello World!")},

		// error
		{"-\r\n", NewError("")},
		{"-oops\r\n", NewError("oops")},
		{"-ERR oops\r\n", NewError("ERR oops")},

		// slice
		{"*0\r\n", []int{}},
		{"*3\r\n:1\r\n:2\r\n:3\r\n", []int{1, 2, 3}},
		{"*3\r\n:1\r\n:2\r\n:3\r\n", []interface{}{int64(1), int64(2), int64(3)}},

		// map
		{"*0\r\n", map[int]int{}},
		{"*6\r\n+A\r\n:1\r\n+B\r\n:2\r\n+C\r\n:3\r\n", map[string]int{"A": 1, "B": 2, "C": 3}},

		// struct
		{"*0\r\n", struct{}{}},
		{"*6\r\n+A\r\n:1\r\n+B\r\n:2\r\n+C\r\n:3\r\n", struct {
			A int
			B int
			C int
		}{1, 2, 3}},
		{"*4\r\n+a\r\n:1\r\n+b\r\n:2\r\n\r\n", struct {
			A int `objconv:"a"`
			B int `objconv:"b"`
			C int `objconv:"c"`
		}{1, 2, 0}},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			z := zero(test.v)
			v := z.Interface()

			if err := objconv.DecodeString(test.s, "resp", v); err != nil {
				t.Errorf("%#v: %s", test.s, err)
			} else if v = z.Elem().Interface(); !reflect.DeepEqual(v, test.v) {
				t.Errorf("%#v:\n- %#v\n- %#v", test.s, test.v, v)
			}
		})
	}
}

func zero(v interface{}) reflect.Value { return test.NewZero(v) }
