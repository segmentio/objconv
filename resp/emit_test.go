package resp

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
		{nil, "$-1\r\n"},

		// bool
		{false, ":0\r\n"},
		{true, ":1\r\n"},

		// int
		{int(42), ":42\r\n"},
		{int8(42), ":42\r\n"},
		{int16(42), ":42\r\n"},
		{int32(42), ":42\r\n"},
		{int64(42), ":42\r\n"},

		// uint
		{uint(42), ":42\r\n"},
		{uint8(42), ":42\r\n"},
		{uint16(42), ":42\r\n"},
		{uint32(42), ":42\r\n"},
		{uint64(42), ":42\r\n"},
		{uintptr(42), ":42\r\n"},

		// float
		{float32(0.5), "$3\r\n0.5\r\n"},
		{float64(0.5), "$3\r\n0.5\r\n"},

		// string
		{"", "+\r\n"},
		{"Hello World!", "+Hello World!\r\n"},
		{"Hello\nWorld!", "$12\r\nHello\nWorld!\r\n"},

		// []byte
		{[]byte(""), "$0\r\n\r\n"},
		{[]byte("Hello World!"), "$12\r\nHello World!\r\n"},
		{[]byte("Hello\nWorld!"), "$12\r\nHello\nWorld!\r\n"},

		// []rune
		{[]rune(""), "+\r\n"},
		{[]rune("Hello World!"), "+Hello World!\r\n"},
		{[]rune("Hello\nWorld!"), "$12\r\nHello\nWorld!\r\n"},

		// time
		{time.Unix(0, 0).In(time.UTC), "+1970-01-01T00:00:00Z\r\n"},
		{time.Second, "+1s\r\n"},

		// error
		{io.EOF, "-EOF\r\n"},

		// slice
		{[]int{}, "*0\r\n"},
		{[]int{1, 2, 3}, "*3\r\n:1\r\n:2\r\n:3\r\n"},

		// map
		{map[int]int{}, "*0\r\n"},
		{map[int]int{1: 0, 2: 0, 3: 0}, "*6\r\n:1\r\n:0\r\n:2\r\n:0\r\n:3\r\n:0\r\n"},

		// struct
		{struct{}{}, "*0\r\n"},
		{struct{ A int }{42}, "*2\r\n+A\r\n:42\r\n"},
		{
			v: struct {
				A int `objconv:"a"`
				B int `objconv:"b,omitempty"`
				C int
			}{1, 0, 42},
			s: "*4\r\n+a\r\n:1\r\n+C\r\n:42\r\n",
		},
	}

	for _, test := range tests {
		if s, err := objconv.EncodeString("resp", test.v); err != nil {
			t.Errorf("%#v: %s", test.v, err)
		} else if s != test.s {
			t.Errorf("%#v:\n- %#v\n- %#v", test.v, test.s, s)
		}
	}
}

func TestCodec(t *testing.T) {
	test.Codec(t, &Emitter{}, &Parser{})
}

func TestStreamCodec(t *testing.T) {
	test.Codec(t, &Emitter{}, &Parser{})
}
