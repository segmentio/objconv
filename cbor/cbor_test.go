package cbor

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/segmentio/objconv"
)

var cborTests = []interface{}{
	// constants
	nil,
	true,
	false,

	// positive integer
	uint(0),
	1,
	23,
	24,
	objconv.Uint8Max,
	objconv.Uint8Max + 1,
	objconv.Uint16Max,
	objconv.Uint16Max + 1,
	objconv.Uint32Max,
	objconv.Uint32Max + 1,

	// negative integer
	-1,
	objconv.Int8Min,
	objconv.Int8Min - 1,
	objconv.Int16Min,
	objconv.Int16Min - 1,
	objconv.Int32Min,
	objconv.Int32Min - 1,

	// float
	float32(0.5),
	float64(0.5),

	// string
	"",
	"Hello World!",

	// bytes
	[]byte(""),
	[]byte("Hello World!"),

	// duration
	time.Nanosecond,
	time.Microsecond,
	time.Millisecond,
	time.Second,
	time.Minute,
	time.Hour,

	// time
	time.Unix(0, 0),
	time.Unix(1, 42),
	time.Unix(17179869184, 999999999),

	// error
	errors.New(""),
	errors.New("Hello World!"),

	// array
	[]int{},
	[]int{1, 2, 3},

	// map
	map[int]int{},
	map[int]int{1: 21, 2: 42, 3: 84},

	// struct
	struct{}{},
	struct{ A int }{42},
	struct{ A, B, C int }{1, 2, 3},
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

func TestMajorType(t *testing.T) {
	m, b := majorType(majorByte(MajorType7, 24))

	if m != MajorType7 {
		t.Error("bad major type:", m)
	}

	if b != 24 {
		t.Error("bad info value:", b)
	}
}

func TestStream(t *testing.T) {
	r, w := io.Pipe()

	e := NewStreamEncoder(w)
	d := NewStreamDecoder(r)

	go func() {
		defer e.Close()

		for i := 0; i != 100; i++ {
			if err := e.Encode(i); err != nil {
				t.Error(err)
			}
		}
	}()

	var i int
	for j := 0; d.Decode(&i) == nil; j++ {
		if i != j {
			t.Errorf("%d != %d", i, j)
		}
	}

	if err := d.Err(); err != nil {
		t.Error(err)
	}
}
