package test

import (
	"io/ioutil"
	"testing"

	"github.com/segmentio/objconv"
)

func BenchmarkEmitBoolTrue(b *testing.B, e objconv.Emitter) {
	BenchmarkEmitBool(b, e, true)
}

func BenchmarkEmitBoolFalse(b *testing.B, e objconv.Emitter) {
	BenchmarkEmitBool(b, e, false)
}

func BenchmarkEmitBool(b *testing.B, e objconv.Emitter, v bool) {
	enc := objconv.NewEncoder(objconv.EncoderConfig{
		Output:  ioutil.Discard,
		Emitter: e,
	})

	b.ResetTimer()

	for i := 0; i != b.N; i++ {
		enc.Encode(v)
	}
}
