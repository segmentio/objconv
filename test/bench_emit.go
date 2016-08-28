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

func BenchmarkEncodeStringEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, "")
}

func BenchmarkEncodeStringShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, "Hello World!")
}

func BenchmarkEncodeStringLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, longString)
}

func BenchmarkEncodeBytesEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, bytesEmpty)
}

func BenchmarkEncodeBytesShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, bytesShort)
}

func BenchmarkEncodeBytesLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, bytesLong)
}

func BenchmarkEncodeSliceInterfaceEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceInterfaceEmpty)
}

func BenchmarkEncodeSliceInterfaceShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceInterfaceShort)
}

func BenchmarkEncodeSliceInterfaceLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceInterfaceLong)
}

func BenchmarkEncodeSliceStringEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStringEmpty)
}

func BenchmarkEncodeSliceStringShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStringShort)
}

func BenchmarkEncodeSliceStringLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStringLong)
}

func BenchmarkEncodeSliceBytesEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceBytesEmpty)
}

func BenchmarkEncodeSliceBytesShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceBytesShort)
}

func BenchmarkEncodeSliceBytesLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceBytesLong)
}

func BenchmarkEncodeSliceStructEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStructEmpty)
}

func BenchmarkEncodeSliceStructShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStructShort)
}

func BenchmarkEncodeSliceStructLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStructLong)
}

func BenchmarkEncodeStructEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, &structEmpty)
}

func BenchmarkEncodeStructShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, &structShort)
}

func BenchmarkEncodeStructLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, &structLong)
}

func BenchmarkEncode(b *testing.B, e objconv.Encoder, v interface{}) {
	for i := 0; i != b.N; i++ {
		e.Encode(v)
	}
}

const longString = `Package json implements encoding and decoding of JSON objects as defined in RFC 4627. The mapping between JSON objects and Go values is described in the documentation for the Marshal and Unmarshal functions.`

var (
	bytesEmpty = []byte{}
	bytesShort = []byte("Hello World!")
	bytesLong  = []byte(`Package json implements encoding and decoding of JSON objects as defined in RFC 4627. The mapping between JSON objects and Go values is described in the documentation for the Marshal and Unmarshal functions.`)

	sliceInterfaceEmpty = []interface{}{}
	sliceInterfaceShort = []interface{}{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil} // 10 items
	sliceInterfaceLong  = []interface{}{
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	} // 100 items

	sliceStringEmpty = []string{}
	sliceStringShort = []string{"", "", "", "", "", "", "", "", "", ""} // 10 items
	sliceStringLong  = []string{
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
	} // 100 items

	sliceBytesEmpty = [][]byte{}
	sliceBytesShort = [][]byte{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil} // 10 items
	sliceBytesLong  = [][]byte{
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	} // 100 items

	sliceStructEmpty = []struct{}{}
	sliceStructShort = []struct{}{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}} // 10 items
	sliceStructLong  = []struct{}{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	} // 100 items

	structEmpty = struct{}{}
	structShort = struct {
		A interface{}
		B interface{}
		C interface{}
	}{}
	structLong = struct {
		A interface{}
		B interface{}
		C interface{}
		D interface{}
		E interface{}
		F interface{}
		G interface{}
		H interface{}
		I interface{}
		J interface{}
		K interface{}
		L interface{}
		M interface{}
		N interface{}
		O interface{}
		P interface{}
		Q interface{}
		R interface{}
		S interface{}
		T interface{}
		U interface{}
		V interface{}
		W interface{}
		X interface{}
		Y interface{}
		Z interface{}
	}{}
)
