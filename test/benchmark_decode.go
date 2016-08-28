package test

import (
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

func BenchmarkDecode(b *testing.B, d objconv.Decoder, r *BenchmarkReader) {
	for i := 0; i != b.N; i++ {
		var v interface{}
		d.Decode(&v)
		r.Reset()
	}
}
