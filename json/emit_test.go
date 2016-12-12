package json

import (
	"io"
	"testing"
	"time"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/test"
)

func TestEmit(t *testing.T) {
	tests := []struct {
		v interface{}
		s string
	}{
		// nil
		{nil, "null"},

		// bool
		{false, "false"},
		{true, "true"},

		// int
		{int(42), "42"},
		{int8(42), "42"},
		{int16(42), "42"},
		{int32(42), "42"},
		{int64(42), "42"},

		// uint
		{uint(42), "42"},
		{uint8(42), "42"},
		{uint16(42), "42"},
		{uint32(42), "42"},
		{uint64(42), "42"},
		{uintptr(42), "42"},

		// float
		{float32(0.5), "0.5"},
		{float64(0.5), "0.5"},

		// string
		{"", `""`},
		{"Hello World!", `"Hello World!"`},
		{"Hello\nWorld!", `"Hello\nWorld!"`},

		// []byte
		{[]byte(""), `""`},
		{[]byte("Hello World!"), `"Hello World!"`},
		{[]byte("Hello\nWorld!"), `"Hello\nWorld!"`},

		// []rune
		{[]rune(""), `""`},
		{[]rune("Hello World!"), `"Hello World!"`},
		{[]rune("Hello\nWorld!"), `"Hello\nWorld!"`},

		// time
		{time.Unix(0, 0).In(time.UTC), `"1970-01-01T00:00:00Z"`},
		{time.Second, `"1s"`},

		// error
		{io.EOF, `"EOF"`},

		// slice
		{[]int{}, `[]`},
		{[]int{1, 2, 3}, `[1,2,3]`},

		// map
		{map[string]int{}, `{}`},
		{map[string]int{"A": 1, "B": 2, "C": 3}, `{"A":1,"B":2,"C":3}`},

		// struct
		{struct{}{}, `{}`},
		{struct{ A int }{42}, `{"A":42}`},
		{
			v: struct {
				A int `objconv:"a"`
				B int `objconv:"b,omitempty"`
				C int
			}{1, 0, 42},
			s: `{"a":1,"C":42}`,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			if s, err := objconv.EncodeString("json", test.v); err != nil {
				t.Errorf("%#v: %s", test.v, err)
			} else if s != test.s {
				t.Errorf("%#v:\n- %#v\n- %#v", test.v, test.s, s)
			}
		})
	}
}

func TestCodec(t *testing.T) {
	test.Codec(t, &Emitter{}, &Parser{})
}

func TestStreamCodec(t *testing.T) {
	test.Codec(t, &Emitter{}, &Parser{})
}
