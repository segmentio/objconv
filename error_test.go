package objconv

import (
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
