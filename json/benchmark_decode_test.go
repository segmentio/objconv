package json

import (
	"testing"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/test"
)

var (
	nilReader  = test.NewBenchmarkReader(`null`)
	nilDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  nilReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeNil(b *testing.B) { test.BenchmarkDecode(b, nilDecoder, nilReader, nil) }

var (
	trueReader  = test.NewBenchmarkReader(`true`)
	trueDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  trueReader,
		Parser: &Parser{},
	})

	falseReader  = test.NewBenchmarkReader(`false`)
	falseDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  falseReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeBoolTrue(b *testing.B)  { test.BenchmarkDecode(b, trueDecoder, trueReader, false) }
func BenchmarkDecodeBoolFalse(b *testing.B) { test.BenchmarkDecode(b, falseDecoder, falseReader, false) }

var (
	intReader  = test.NewBenchmarkReader(`42`)
	intDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  intReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeInt(b *testing.B) { test.BenchmarkDecode(b, intDecoder, intReader, 0) }

var (
	floatReader  = test.NewBenchmarkReader(`1.234`)
	floatDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  floatReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeFloat(b *testing.B) { test.BenchmarkDecode(b, floatDecoder, floatReader, 0.0) }

var (
	stringReader  = test.NewBenchmarkReader(`"Hello World!"`)
	stringDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  stringReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeString(b *testing.B) { test.BenchmarkDecode(b, stringDecoder, stringReader, "") }

var (
	arrayReader  = test.NewBenchmarkReader(`[1,2,3]`)
	arrayDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  arrayReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeArray(b *testing.B) {
	test.BenchmarkDecode(b, arrayDecoder, arrayReader, ([]interface{})(nil))
}

var (
	mapReader  = test.NewBenchmarkReader(`{"A":1,"B":2,"C":3}`)
	mapDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  mapReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeMap(b *testing.B) {
	test.BenchmarkDecode(b, mapDecoder, mapReader, (map[string]interface{})(nil))
}

var (
	structReader  = test.NewBenchmarkReader(`{"A":1,"B":2,"C":3}`)
	structDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  structReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeStruct(b *testing.B) {
	test.BenchmarkDecode(b, structDecoder, structReader, struct {
		A int
		B int
		C int
	}{})
}
