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

func BenchmarkEncodeStringZero(b *testing.B)  { test.BenchmarkEncodeStringZero(b, benchmarkEncoder) }
func BenchmarkEncodeStringShort(b *testing.B) { test.BenchmarkEncodeStringShort(b, benchmarkEncoder) }
func BenchmarkEncodeStringLong(b *testing.B)  { test.BenchmarkEncodeStringLong(b, benchmarkEncoder) }

func BenchmarkEncodeBytesZero(b *testing.B)  { test.BenchmarkEncodeBytesZero(b, benchmarkEncoder) }
func BenchmarkEncodeBytesShort(b *testing.B) { test.BenchmarkEncodeBytesShort(b, benchmarkEncoder) }
func BenchmarkEncodeBytesLong(b *testing.B)  { test.BenchmarkEncodeBytesLong(b, benchmarkEncoder) }
