package objconv

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

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
		t.Run(fmt.Sprint(test.v), func(t *testing.T) {
			if e := convertPanicToError(test.v); !reflect.DeepEqual(e, test.e) {
				t.Errorf("convertPanicToError(%#v): %#v != %#v", test.v, test.e, e)
			}
		})
	}
}
