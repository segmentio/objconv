package resp

import (
	"testing"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/test"
)

var (
	nilReader  = test.NewBenchmarkReader("$-1\r\n")
	nilDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  nilReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeNil(b *testing.B) { test.BenchmarkDecode(b, nilDecoder, nilReader) }

var (
	trueReader  = test.NewBenchmarkReader(":1\r\n")
	trueDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  trueReader,
		Parser: &Parser{},
	})

	falseReader  = test.NewBenchmarkReader(":0\r\n")
	falseDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  falseReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeBoolTrue(b *testing.B)  { test.BenchmarkDecode(b, trueDecoder, trueReader) }
func BenchmarkDecodeBoolFalse(b *testing.B) { test.BenchmarkDecode(b, falseDecoder, falseReader) }

var (
	intReader  = test.NewBenchmarkReader(":42\r\n")
	intDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  intReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeInt(b *testing.B) { test.BenchmarkDecode(b, intDecoder, intReader) }

var (
	stringReader  = test.NewBenchmarkReader("+Hello World!\r\n")
	stringDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  stringReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeString(b *testing.B) { test.BenchmarkDecode(b, stringDecoder, stringReader) }

var (
	arrayReader  = test.NewBenchmarkReader("*3\r\n:1\r\n:2\r\n:3\r\n")
	arrayDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  arrayReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeArray(b *testing.B) { test.BenchmarkDecode(b, arrayDecoder, arrayReader) }

var (
	mapReader  = test.NewBenchmarkReader("*6\r\n+A\r\n:1\r\n+B\r\n:2\r\n+C\r\n:3\r\n")
	mapDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  mapReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeMap(b *testing.B) { test.BenchmarkDecode(b, mapDecoder, mapReader) }
