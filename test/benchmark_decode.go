package test

import (
	"bytes"
	"testing"

	"github.com/segmentio/objconv"
)

type BenchmarkReader struct {
	*bytes.Reader
	b []byte
}

func NewBenchmarkReader(b []byte) *BenchmarkReader {
	return &BenchmarkReader{bytes.NewReader(b), b}
}

func (r *BenchmarkReader) Reset() { r.Reader.Reset(r.b) }

func BenchmarkDecode(b *testing.B, d objconv.Decoder, r *BenchmarkReader) {
	for i := 0; i != b.N; i++ {
		var v interface{}
		d.Decode(&v)
		r.Reset()
	}
}
