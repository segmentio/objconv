package json

import (
	"errors"
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
		{`null`, nil},

		// bool
		{`false`, false},
		{`true`, true},

		// int
		{`42`, int(42)},
		{`42`, int(42)},
		{`42`, int8(42)},
		{`42`, int16(42)},
		{`42`, int32(42)},
		{`42`, int64(42)},

		// uint
		{`42`, uint(42)},
		{`42`, uint8(42)},
		{`42`, uint16(42)},
		{`42`, uint32(42)},
		{`42`, uint64(42)},
		{`42`, uintptr(42)},

		// float
		{`42`, float32(42)},
		{`42`, float64(42)},
		{`0.5`, float32(0.5)},
		{`0.5`, float64(0.5)},
		{`42`, float32(42)},
		{`42`, float64(42)},
		{`0.5`, float32(0.5)},
		{`0.5`, float64(0.5)},
		{`1e9`, float32(1e9)},
		{`1e9`, float64(1e9)},

		// time
		{`"1970-01-01T00:00:00Z"`, time.Unix(0, 0).In(time.UTC)},
		{`"1970-01-01T00:00:00Z"`, time.Unix(0, 0).In(time.UTC)},
		{`"1s"`, time.Second},
		{`"1s"`, time.Second},

		// string
		{`""`, ""},
		{`"Hello World!"`, "Hello World!"},
		{`"Hello World!"`, "Hello World!"},
		{`"Hello\"World!"`, "Hello\"World!"},
		{`"Hello\/World!"`, "Hello/World!"},
		{`"Hello\\World!"`, "Hello\\World!"},
		{`"Hello\bWorld!"`, "Hello\bWorld!"},
		{`"Hello\fWorld!"`, "Hello\fWorld!"},
		{`"Hello\nWorld!"`, "Hello\nWorld!"},
		{`"Hello\rWorld!"`, "Hello\rWorld!"},
		{`"Hello\tWorld!"`, "Hello\tWorld!"},

		// []byte
		{`""`, []byte("")},
		{`"Hello World!"`, []byte("Hello World!")},
		{`"Hello World!"`, []byte("Hello World!")},

		// []rune
		{`""`, []rune("")},
		{`"Hello World!"`, []rune("Hello World!")},
		{`"Hello World!"`, []rune("Hello World!")},

		// error
		{`""`, errors.New("")},
		{`"oops"`, errors.New("oops")},

		// slice
		{`[]`, []int{}},
		{`[1,2,3]`, []int{1, 2, 3}},
		{`[1,2,3]`, []interface{}{int64(1), int64(2), int64(3)}},

		// map
		{`{}`, map[int]int{}},
		{`{"A":1,"B":2,"C":3}`, map[string]int{"A": 1, "B": 2, "C": 3}},

		// struct
		{`{}`, struct{}{}},
		{`{"A":1,"B":2,"C":3}`, struct {
			A int
			B int
			C int
		}{1, 2, 3}},
		{`{"a":1,"b":2}`, struct {
			A int `objconv:"a"`
			B int `objconv:"b"`
			C int `objconv:"c"`
		}{1, 2, 0}},
	}

	for _, test := range tests {
		z := zero(test.v)
		v := z.Interface()

		if err := objconv.DecodeString(test.s, "json", v); err != nil {
			t.Errorf("%#v: %s", test.s, err)
		} else if v = z.Elem().Interface(); !reflect.DeepEqual(v, test.v) {
			t.Errorf("%#v:\n- %#v\n- %#v", test.s, test.v, v)
		}
	}
}

func zero(v interface{}) reflect.Value { return test.NewZero(v) }
