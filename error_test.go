package objconv

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestUnsupportedTypeError(t *testing.T) {
	e := &UnsupportedTypeError{Type: reflect.TypeOf(1)}
	s := e.Error()

	if s != "objconv: unsupported type: int" {
		t.Errorf("invalid error message for unsupported type: %#v", s)
	}
}

func TestConvertPanicToError(t *testing.T) {
	tests := []struct {
		v interface{}
		e error
	}{
		{
			v: nil,
			e: nil,
		},
		{
			v: io.EOF,
			e: io.EOF,
		},
		{
			v: "Hello World",
			e: fmt.Errorf("objconv: Hello World"),
		},
	}

	for _, test := range tests {
		if e := convertPanicToError(test.v); !reflect.DeepEqual(e, test.e) {
			t.Errorf("convertPanicToError(%#v): %#v != %#v", test.v, test.e, e)
		}
	}
}
