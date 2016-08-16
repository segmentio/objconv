package test

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/segmentio/objconv"
)

// Codec tests an encoder with f set as the formatter and decoder with p set as
// the parser.
// The function is intended to be used in tests.
func Codec(t *testing.T, e objconv.Emitter, p objconv.Parser) {
	type T struct {
		A int
		B int
		C int
		D int `objconv:"-"`
	}

	tests := []interface{}{
		// nil
		nil,

		// bool
		true,
		false,

		// int
		int(0),
		int(-1),
		int(42),
		int(objconv.IntMin),
		int(objconv.IntMax),

		// int8
		int8(0),
		int8(-1),
		int8(42),
		int8(objconv.Int8Min),
		int8(objconv.Int8Max),

		// int16
		int16(0),
		int16(-1),
		int16(42),
		int16(objconv.Int16Min),
		int16(objconv.Int16Max),

		// int32
		int32(0),
		int32(-1),
		int32(42),
		int32(objconv.Int32Min),
		int32(objconv.Int32Max),

		// int64
		int64(0),
		int64(-1),
		int64(42),
		int64(objconv.Int64Min),
		int64(objconv.Int64Max),

		// uint
		uint(0),
		uint(42),
		uint(objconv.UintMin),
		uint(objconv.UintMax),

		// uint8
		uint8(0),
		uint8(42),
		uint8(objconv.Uint8Min),
		uint8(objconv.Uint8Max),

		// uint16
		uint16(0),
		uint16(42),
		uint16(objconv.Uint16Min),
		uint16(objconv.Uint16Max),

		// uint32
		uint32(0),
		uint32(42),
		uint32(objconv.Uint32Min),
		uint32(objconv.Uint32Max),

		// uint64
		uint64(0),
		uint64(42),
		uint64(objconv.Uint64Min),
		uint64(objconv.Uint64Max),

		// uintptr
		uintptr(0),
		uintptr(42),
		uintptr(objconv.UintptrMin),
		uintptr(objconv.UintptrMax),

		// float32
		float32(0),
		float32(1),
		float32(42),
		float32(0.5),
		float32(1.234),
		float32(objconv.Float32IntMax),
		float32(objconv.Float32IntMin),

		// float64
		float64(0),
		float64(1),
		float64(42),
		float64(0.5),
		float64(1.234),
		float64(objconv.Float64IntMax),
		float64(objconv.Float64IntMin),

		// string
		"",
		"Hello World!",
		"师傅",
		"�",

		// []byte
		[]byte(""),
		[]byte("Hello World!"),
		[]byte("师傅"),
		[]byte("�"),

		// []rune
		[]rune(""),
		[]rune("Hello World!"),
		[]rune("师傅"),
		[]rune("�"),

		// []int
		[]int{},
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},

		// []string
		[]string{},
		[]string{""},
		[]string{"hello", "world"},
		[]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},

		// []interface{}
		[]interface{}{},
		[]interface{}{nil, int64(42), "Hello World!", []byte("Hello World!"), []interface{}{"A", "B", "C"}},

		// map[int]int
		map[int]int{},
		map[int]int{1: 1, 2: 2, 3: 3},

		// map[string]interface{}
		map[string]interface{}{},
		map[string]interface{}{"A": nil, "B": int64(42)},

		// struct
		struct{}{},
		struct{ A int }{42},
		struct {
			V []interface{} `objconv:"values,omitempty"`
		}{},
		struct {
			V []T `objconv:"values,omitempty"`
		}{[]T{{1, 2, 3, 0}, {4, 5, 6, 0}, {7, 8, 9, 0}}},
	}

	b := &bytes.Buffer{}

	for _, test := range tests {
		b.Reset()
		v := NewZero(test).Interface()
		e := objconv.NewEncoder(objconv.EncoderConfig{Output: b, Emitter: e})
		d := objconv.NewDecoder(objconv.DecoderConfig{Input: b, Parser: p})

		if err := e.Encode(test); err != nil {
			t.Errorf("%T(%v):\nencode: %s\nbuffer: %#v", test, test, err, b.String())
			continue
		}

		if err := d.Decode(v); err != nil {
			t.Errorf("%T(%v):\ndecode: %s\nbuffer: %#v", test, test, err, b.String())
			continue
		}

		if res := reflect.ValueOf(v).Elem().Interface(); !reflect.DeepEqual(test, res) {
			t.Errorf("%T(%v): mismatch:\n- exp: %#v\n- got: %#v", test, test, test, res)
		}
	}
}

// StreamCodec tests the stream encoder and decoder with e and p as emitter and
// parser.
func StreamCodec(t *testing.T, e objconv.Emitter, p objconv.Parser) {
	r, w := io.Pipe()

	enc := objconv.NewStreamEncoder(objconv.EncoderConfig{
		Output:  w,
		Emitter: e,
	})

	dec := objconv.NewStreamDecoder(objconv.DecoderConfig{
		Input:  r,
		Parser: p,
	})

	test := []interface{}{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	go func() {
		defer w.Close()

		if err := enc.Open(len(test)); err != nil {
			t.Error("opening encoder:", err)
		}

		for i, v := range test {
			if err := enc.Encode(v); err != nil {
				t.Errorf("encoding %#v (%d): %s", v, i, err)
			}
		}

		if err := enc.Close(); err != nil {
			t.Error("closing encoder:", err)
		}

		if err := enc.Close(); err != objconv.ErrClosed {
			t.Error("closed encoder did not returned ErrClosed when closed again")
		}
	}()

	for i := 0; true; i++ {
		var v interface{}

		if err := dec.Decode(&v); err != nil {
			if err != io.EOF || i != len(test) {
				t.Errorf("decoding %#v (%d): %s", test[i], i, err)
			}
			break
		}

		if !reflect.DeepEqual(v, test[i]) {
			t.Errorf("decoding %#v (%d): value mismatch: %#v", test[i], i, v)
		}
	}

	if err := dec.Error(); err != io.EOF {
		t.Error("decoder should have EOF set after consuming all the stream")
	}

	r.Close()
}
