package resp

import (
	"errors"
	"reflect"
	"testing"
)

var respDecodeTests = []struct {
	v interface{}
	s string
}{
	{nil, "$-1\r\n"},

	{0, ":0\r\n"},
	{-1, ":-1\r\n"},
	{42, ":42\r\n"},

	{"", "+\r\n"},
	{"Hello World!", "+Hello World!\r\n"},
	{"Hello\nWorld!", "+Hello\nWorld!\r\n"},
	{"Hello\r\nWorld!", "$13\r\nHello\r\nWorld!\r\n"},

	{[]byte{}, "$0\r\n\r\n"},
	{[]byte("Hello World!"), "$12\r\nHello World!\r\n"},

	{errors.New(""), "-\r\n"},
	{errors.New("oops"), "-oops\r\n"},
	{errors.New("A"), "-A\r\n"},

	{[]int{}, "*0\r\n"},
	{[]int{1, 2, 3}, "*3\r\n:1\r\n:2\r\n:3\r\n"},
}

func TestUnmarshal(t *testing.T) {
	for _, test := range respDecodeTests {
		t.Run(testName(test.s), func(t *testing.T) {
			var typ reflect.Type

			if test.v == nil {
				typ = reflect.TypeOf((*interface{})(nil)).Elem()
			} else {
				typ = reflect.TypeOf(test.v)
			}

			val := reflect.New(typ)
			err := Unmarshal([]byte(test.s), val.Interface())

			if err != nil {
				t.Error(err)
			}

			v1 := test.v
			v2 := val.Elem().Interface()

			if !reflect.DeepEqual(v1, v2) {
				t.Error(v2)
			}
		})
	}
}
