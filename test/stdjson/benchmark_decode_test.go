package json

import (
	"encoding/json"
	"testing"

	"github.com/segmentio/objconv/test"
)

var (
	nilReader  = test.NewBenchmarkReader(`null`)
	nilDecoder = json.NewDecoder(nilReader)
)

func BenchmarkDecodeNil(b *testing.B) { test.BenchmarkDecode(b, nilDecoder, nilReader, nil) }

var (
	trueReader  = test.NewBenchmarkReader(`true`)
	trueDecoder = json.NewDecoder(trueReader)

	falseReader  = test.NewBenchmarkReader(`false`)
	falseDecoder = json.NewDecoder(falseReader)
)

func BenchmarkDecodeBoolTrue(b *testing.B)  { test.BenchmarkDecode(b, trueDecoder, trueReader, false) }
func BenchmarkDecodeBoolFalse(b *testing.B) { test.BenchmarkDecode(b, falseDecoder, falseReader, false) }

var (
	intReader  = test.NewBenchmarkReader(`42`)
	intDecoder = json.NewDecoder(intReader)
)

func BenchmarkDecodeInt(b *testing.B) { test.BenchmarkDecode(b, intDecoder, intReader, 0) }

var (
	floatReader  = test.NewBenchmarkReader("1.234")
	floatDecoder = json.NewDecoder(floatReader)
)

func BenchmarkDecodeFloat(b *testing.B) { test.BenchmarkDecode(b, floatDecoder, floatReader, 0.0) }

var (
	stringReader  = test.NewBenchmarkReader(`"Hello World!"`)
	stringDecoder = json.NewDecoder(stringReader)
)

func BenchmarkDecodeString(b *testing.B) { test.BenchmarkDecode(b, stringDecoder, stringReader, "") }

var (
	arrayReader  = test.NewBenchmarkReader(`[1,2,3]`)
	arrayDecoder = json.NewDecoder(arrayReader)
)

func BenchmarkDecodeArray(b *testing.B) {
	test.BenchmarkDecode(b, arrayDecoder, arrayReader, ([]interface{})(nil))
}

var (
	mapReader  = test.NewBenchmarkReader(`{"A":1,"B":2,"C":3}`)
	mapDecoder = json.NewDecoder(mapReader)
)

func BenchmarkDecodeMap(b *testing.B) {
	test.BenchmarkDecode(b, mapDecoder, mapReader, (map[string]interface{})(nil))
}

var (
	structReader  = test.NewBenchmarkReader(`{"A":1,"B":2,"C":3}`)
	structDecoder = json.NewDecoder(structReader)
)

func BenchmarkDecodeStruct(b *testing.B) {
	test.BenchmarkDecode(b, structDecoder, structReader, struct {
		A int
		B int
		C int
	}{})
}
