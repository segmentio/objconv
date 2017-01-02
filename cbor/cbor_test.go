package cbor

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/segmentio/objconv"
)

var cborTests = []interface{}{
	// constants
	nil,
	true,
	false,

	// positive integer
	0,
	1,
	23,
	24,
	objconv.Uint8Max,
	objconv.Uint8Max + 1,
	objconv.Uint16Max,
	objconv.Uint16Max + 1,
	objconv.Uint32Max,
	objconv.Uint32Max + 1,

	// float
	float32(0.5),
	float64(0.5),
}

func makeMap(n int) map[string]string {
	m := make(map[string]string, n)
	for i := 0; i != n; i++ {
		m[strconv.Itoa(i)] = "A"
	}
	return m
}

func testName(v interface{}) string {
	s := fmt.Sprintf("%T:%v", v, v)
	if len(s) > 20 {
		s = s[:20] + "..."
	}
	return s
}

func TestCBOR(t *testing.T) {
	for _, test := range cborTests {
		t.Run(testName(test), func(t *testing.T) {
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

func TestMajorTypeOf(t *testing.T) {
	m, b := majorTypeOf(majorByte(MajorType7, 24))

	if m != MajorType7 {
		t.Error("bad major type:", m)
	}

	if b != 24 {
		t.Error("bad info value:", b)
	}
}
