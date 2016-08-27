package test

import (
	"math"
	"testing"

	"github.com/segmentio/objconv"
)

func BenchmarkEncodeNil(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, nil)
}

func BenchmarkEncodeBoolTrue(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, true)
}

func BenchmarkEncodeBoolFalse(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, false)
}

func BenchmarkEncodeIntZero(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, int64(0))
}

func BenchmarkEncodeIntShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, int64(100))
}

func BenchmarkEncodeIntLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, int64(math.MaxInt64))
}

func BenchmarkEncodeUintZero(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, uint64(0))
}

func BenchmarkEncodeUintShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, uint64(100))
}

func BenchmarkEncodeUintLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, uint64(math.MaxUint64))
}

func BenchmarkEncodeFloatZero(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, float64(0))
}

func BenchmarkEncodeFloatShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, float64(1.234))
}

func BenchmarkEncodeFloatLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, float64(math.MaxFloat64))
}

func BenchmarkEncodeStringZero(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, "")
}

func BenchmarkEncodeStringShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, "Hello World!")
}

func BenchmarkEncodeStringLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, longString)
}

func BenchmarkEncodeBytesZero(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, zeroBytes)
}

func BenchmarkEncodeBytesShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, shortBytes)
}

func BenchmarkEncodeBytesLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, longBytes)
}

func BenchmarkEncode(b *testing.B, e objconv.Encoder, v interface{}) {
	for i := 0; i != b.N; i++ {
		e.Encode(v)
	}
}

const longString = `Package json implements encoding and decoding of JSON objects as defined in RFC 4627. The mapping between JSON objects and Go values is described in the documentation for the Marshal and Unmarshal functions.`

var (
	zeroBytes  = []byte{}
	shortBytes = []byte("Hello World!")
	longBytes  = []byte(`Package json implements encoding and decoding of JSON objects as defined in RFC 4627. The mapping between JSON objects and Go values is described in the documentation for the Marshal and Unmarshal functions.`)
)
