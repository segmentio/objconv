package json

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestMarshal(t *testing.T) {
	for _, test := range jsonTests {
		t.Run(test.s, func(t *testing.T) {
			b, err := Marshal(test.v)

			if err != nil {
				t.Error(err)
			}

			if s := string(b); s != test.s {
				t.Error(s)
			}
		})
	}
}

func TestStreamEncoder(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := NewStreamEncoder(buf)

	for i := 0; i != 10; i++ {
		if err := enc.Encode(i); err != nil {
			t.Error(err)
		}
	}

	if err := enc.Close(); err != nil {
		t.Error(err)
	}

	if s := buf.String(); s != `[0,1,2,3,4,5,6,7,8,9]` {
		t.Error(s)
	}
}

func BenchmarkMarshal(b *testing.B) {
	for _, test := range jsonTests {
		b.Run(test.s, func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				if _, err := Marshal(test.v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEncoder(b *testing.B) {
	e := NewEncoder(ioutil.Discard)

	for _, test := range jsonTests {
		b.Run(test.s, func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				if err := e.Encode(test.v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(test.s)))
		})
	}
}
