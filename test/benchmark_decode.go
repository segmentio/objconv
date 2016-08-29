package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/segmentio/objconv"
)

type BenchmarkReader struct {
	*strings.Reader
	s string
}

func NewBenchmarkReader(s string) *BenchmarkReader {
	return &BenchmarkReader{strings.NewReader(s), s}
}

func (r *BenchmarkReader) Reset() { r.Reader.Reset(r.s) }

func BenchmarkDecode(b *testing.B, d objconv.Decoder, r *BenchmarkReader, v interface{}) {
	var z interface{}

	if v == nil {
		z = &v
	} else {
		z = reflect.New(reflect.TypeOf(v)).Interface()
	}

	b.ResetTimer()

	for i := 0; i != b.N; i++ {
		d.Decode(z)
		r.Reset()
	}
}
