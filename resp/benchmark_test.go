package resp

import (
	"io/ioutil"
	"testing"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/test"
)

var (
	benchmarkEncoder = objconv.NewEncoder(objconv.EncoderConfig{
		Output:  ioutil.Discard,
		Emitter: &Emitter{},
	})
)

func BenchmarkEncodeNil(b *testing.B) { test.BenchmarkEncodeNil(b, benchmarkEncoder) }

func BenchmarkEncodeBoolTrue(b *testing.B)  { test.BenchmarkEncodeBoolTrue(b, benchmarkEncoder) }
func BenchmarkEncodeBoolFalse(b *testing.B) { test.BenchmarkEncodeBoolFalse(b, benchmarkEncoder) }

func BenchmarkEncodeIntZero(b *testing.B)  { test.BenchmarkEncodeIntZero(b, benchmarkEncoder) }
func BenchmarkEncodeIntShort(b *testing.B) { test.BenchmarkEncodeIntShort(b, benchmarkEncoder) }
func BenchmarkEncodeIntLong(b *testing.B)  { test.BenchmarkEncodeIntLong(b, benchmarkEncoder) }

func BenchmarkEncodeUintZero(b *testing.B)  { test.BenchmarkEncodeUintZero(b, benchmarkEncoder) }
func BenchmarkEncodeUintShort(b *testing.B) { test.BenchmarkEncodeUintShort(b, benchmarkEncoder) }
func BenchmarkEncodeUintLong(b *testing.B)  { test.BenchmarkEncodeUintLong(b, benchmarkEncoder) }

func BenchmarkEncodeFloatZero(b *testing.B)  { test.BenchmarkEncodeFloatZero(b, benchmarkEncoder) }
func BenchmarkEncodeFloatShort(b *testing.B) { test.BenchmarkEncodeFloatShort(b, benchmarkEncoder) }
func BenchmarkEncodeFloatLong(b *testing.B)  { test.BenchmarkEncodeFloatLong(b, benchmarkEncoder) }

func BenchmarkEncodeStringEmpty(b *testing.B) { test.BenchmarkEncodeStringEmpty(b, benchmarkEncoder) }
func BenchmarkEncodeStringShort(b *testing.B) { test.BenchmarkEncodeStringShort(b, benchmarkEncoder) }
func BenchmarkEncodeStringLong(b *testing.B)  { test.BenchmarkEncodeStringLong(b, benchmarkEncoder) }

func BenchmarkEncodeBytesEmpty(b *testing.B) { test.BenchmarkEncodeBytesEmpty(b, benchmarkEncoder) }
func BenchmarkEncodeBytesShort(b *testing.B) { test.BenchmarkEncodeBytesShort(b, benchmarkEncoder) }
func BenchmarkEncodeBytesLong(b *testing.B)  { test.BenchmarkEncodeBytesLong(b, benchmarkEncoder) }

func BenchmarkEncodeSliceInterfaceEmpty(b *testing.B) {
	test.BenchmarkEncodeSliceInterfaceEmpty(b, benchmarkEncoder)
}
func BenchmarkEncodeSliceInterfaceShort(b *testing.B) {
	test.BenchmarkEncodeSliceInterfaceShort(b, benchmarkEncoder)
}
func BenchmarkEncodeSliceInterfaceLong(b *testing.B) {
	test.BenchmarkEncodeSliceInterfaceLong(b, benchmarkEncoder)
}

func BenchmarkEncodeSliceStringEmpty(b *testing.B) {
	test.BenchmarkEncodeSliceStringEmpty(b, benchmarkEncoder)
}
func BenchmarkEncodeSliceStringShort(b *testing.B) {
	test.BenchmarkEncodeSliceStringShort(b, benchmarkEncoder)
}
func BenchmarkEncodeSliceStringLong(b *testing.B) {
	test.BenchmarkEncodeSliceStringLong(b, benchmarkEncoder)
}

func BenchmarkEncodeSliceBytesEmpty(b *testing.B) {
	test.BenchmarkEncodeSliceBytesEmpty(b, benchmarkEncoder)
}
func BenchmarkEncodeSliceBytesShort(b *testing.B) {
	test.BenchmarkEncodeSliceBytesShort(b, benchmarkEncoder)
}
func BenchmarkEncodeSliceBytesLong(b *testing.B) {
	test.BenchmarkEncodeSliceBytesLong(b, benchmarkEncoder)
}

func BenchmarkEncodeSliceStructEmpty(b *testing.B) {
	test.BenchmarkEncodeSliceStructEmpty(b, benchmarkEncoder)
}
func BenchmarkEncodeSliceStructShort(b *testing.B) {
	test.BenchmarkEncodeSliceStructShort(b, benchmarkEncoder)
}
func BenchmarkEncodeSliceStructLong(b *testing.B) {
	test.BenchmarkEncodeSliceStructLong(b, benchmarkEncoder)
}
