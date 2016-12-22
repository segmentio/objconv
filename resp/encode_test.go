package resp

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	for _, test := range respTests {
		t.Run(strings.Replace(test.s, "\r\n", "", -1), func(t *testing.T) {
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

func BenchmarkEncoder(b *testing.B) {
	e := NewEncoder(ioutil.Discard)

	for _, test := range respTests {
		b.Run(strings.Replace(test.s, "\r\n", "", -1), func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				if err := e.Encode(test.v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(test.s)))
		})
	}
}

func BenchmarkMarshal(b *testing.B) {
	for _, test := range respTests {
		b.Run(strings.Replace(test.s, "\r\n", "", -1), func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				if _, err := Marshal(test.v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(test.s)))
		})
	}
}
