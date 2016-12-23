package msgpack

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

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

	// uint8, uint16, uint32, uint64
	uint8(objconv.Uint8Max),
	uint16(objconv.Uint16Max),
	uint32(objconv.Uint32Max),
	uint64(objconv.Uint64Max),

	// float32
	float32(0),
	float32(objconv.Float32IntMin),
	float32(objconv.Float32IntMax),

	// float64
	float64(0),
	float64(objconv.Float64IntMin),
	float64(objconv.Float64IntMax),

	// fixstr
	"",
	"Hello World!",

	// str8, str16, str32
	strings.Repeat("A", 32),
	strings.Repeat("A", objconv.Uint8Max+1),
	strings.Repeat("A", objconv.Uint16Max+1),

	// bin8, bin16, bin32
	[]byte(""),
	[]byte("Hello World!"),
	bytes.Repeat([]byte("A"), objconv.Uint8Max+1),
	bytes.Repeat([]byte("A"), objconv.Uint16Max+1),

	// duration
	time.Nanosecond,
	time.Microsecond,
	time.Millisecond,
	time.Second,
	time.Minute,
	time.Hour,

	// time
	time.Time{},
	time.Now(),

	// error
	errors.New(""),
	errors.New("Hello World!"),
	errors.New(strings.Repeat("A", objconv.Uint8Max+1)),
	errors.New(strings.Repeat("A", objconv.Uint16Max+1)),

	// fixarray
	[]int{},
	[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},

	// array16, array32
	make([]int, objconv.Uint8Max+1),
	make([]int, objconv.Uint16Max+1),

	// fixmap
	makeMap(0),
	makeMap(15),

	// map16, map32
	makeMap(objconv.Uint8Max + 1),
	makeMap(objconv.Uint16Max + 1),
}

func makeMap(n int) map[int]int {
	m := make(map[int]int, n)
	for i := 0; i != n; i++ {
		m[i] = i
	}
	return m
}

func testName(v interface{}) string {
	s := fmt.Sprintf("%v", v)
	if len(s) > 20 {
		s = s[:20] + "..."
	}
	return s
}

func TestMsgpack(t *testing.T) {
	for _, test := range msgpackTests {
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
