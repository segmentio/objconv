package cbor

import (
	"bytes"
	"reflect"
	"sync"
	"testing"

	"github.com/segmentio/objconv"
)

var (
	cborDecodeTests [][]byte
	cborDecodeOnce  sync.Once
)

func initDecodeTests() {
	cborDecodeOnce.Do(func() {
		cborDecodeTests = make([][]byte, len(cborTests))

		for i, v := range cborTests {
			var err error

			if cborDecodeTests[i], _ = Marshal(v); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkUnmarshal(b *testing.B) {
	initDecodeTests()

	for i, test := range cborTests {
		var t reflect.Type

		if test == nil {
			t = reflect.TypeOf((*interface{})(nil)).Elem()
		} else {
			t = reflect.TypeOf(test)
		}

		v := reflect.New(t).Interface()
		s := cborDecodeTests[i]

		b.Run(testName(test), func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				if err := Unmarshal(s, v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(s)))
		})
	}
}

func BenchmarkDecoder(b *testing.B) {
	initDecodeTests()

	r := bytes.NewReader(nil)
	p := NewParser(nil)
	d := objconv.NewDecoder(p)

	for i, test := range cborTests {
		var t reflect.Type

		if test == nil {
			t = reflect.TypeOf((*interface{})(nil)).Elem()
		} else {
			t = reflect.TypeOf(test)
		}

		v := reflect.New(t).Interface()
		s := cborDecodeTests[i]

		b.Run(testName(test), func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				r.Reset(s)
				p.Reset(r)

				if err := d.Decode(v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(s)))
		})
	}
}
