package msgpack

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/segmentio/objconv"
)

var msgpackTests = []interface{}{
	// constants
	nil,
	false,
	true,

	// positive fixint
	0,
	1,
	127,

	// negative fixint
	-1,
	-10,
	-31,

	// int8
	-32,
	objconv.Int8Min,

	// int16
	objconv.Int8Max + 1,
	objconv.Int8Min - 1,
	objconv.Int16Max,
	objconv.Int16Min,

	// int32
	objconv.Int16Max + 1,
	objconv.Int16Min - 1,
	objconv.Int32Max,
	objconv.Int32Min,

	// int64
	objconv.Int32Max + 1,
	objconv.Int32Min - 1,
	int64(objconv.Int64Max),
	int64(objconv.Int64Min),

	// uint8
	uint8(objconv.Uint8Max),

	// uint16
	uint16(objconv.Uint16Max),

	// uint32
	uint32(objconv.Uint32Max),

	// uint64
	uint64(objconv.Uint64Max),

	// float32
	float32(0),
	float32(objconv.Float32IntMin),
	float32(objconv.Float32IntMax),

	// float64
	float64(0),
	float64(objconv.Float64IntMin),
	float64(objconv.Float64IntMax),
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
