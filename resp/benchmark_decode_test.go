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
