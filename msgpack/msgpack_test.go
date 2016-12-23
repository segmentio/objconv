package msgpack

import (
	"fmt"
	"reflect"
	"testing"
)

var msgpackTests = []interface{}{
	nil,
}

func TestMsgpack(t *testing.T) {
	for _, test := range msgpackTests {
		t.Run(fmt.Sprintf("%#v", test), func(t *testing.T) {
			var typ reflect.Type

			if test == nil {
				typ = reflect.TypeOf((*interface{})(nil)).Elem()
			} else {
				typ = reflect.TypeOf(test)
			}

			val := reflect.New(typ)
			b, err := Marshal(test)

			if err != nil {
				t.Error(err)
				return
			}

			if err := Unmarshal(b, val.Interface()); err != nil {
				t.Error(err)
				return
			}

			x1 := test
			x2 := val.Elem().Interface()

			if !reflect.DeepEqual(x1, x2) {
				t.Errorf("%#v", x2)
			}
		})
	}
}
