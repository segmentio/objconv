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

func BenchmarkDecode(b *testing.B, p objconv.Parser, r *BenchmarkReader) {
	var v interface{}
	var d = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  r,
		Parser: p,
	})

	b.ResetTimer()

	for i := 0; i != b.N; i++ {
		d.Decode(&v)
		r.Reset()
		v = nil
	}
}
