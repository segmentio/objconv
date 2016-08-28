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

func BenchmarkDecodeNil(b *testing.B) { test.BenchmarkDecode(b, nilDecoder, nilReader) }

var (
	trueReader  = test.NewBenchmarkReader(`true`)
	trueDecoder = json.NewDecoder(trueReader)

	falseReader  = test.NewBenchmarkReader(`false`)
	falseDecoder = json.NewDecoder(falseReader)
)

func BenchmarkDecodeBoolTrue(b *testing.B)  { test.BenchmarkDecode(b, trueDecoder, trueReader) }
func BenchmarkDecodeBoolFalse(b *testing.B) { test.BenchmarkDecode(b, falseDecoder, falseReader) }

var (
	intReader  = test.NewBenchmarkReader(`42`)
	intDecoder = json.NewDecoder(intReader)
)

func BenchmarkDecodeInt(b *testing.B) { test.BenchmarkDecode(b, intDecoder, intReader) }

var (
	stringReader  = test.NewBenchmarkReader(`"Hello World!"`)
	stringDecoder = json.NewDecoder(stringReader)
)

func BenchmarkDecodeString(b *testing.B) { test.BenchmarkDecode(b, stringDecoder, stringReader) }

var (
	arrayReader  = test.NewBenchmarkReader(`[1,2,3]`)
	arrayDecoder = json.NewDecoder(arrayReader)
)

func BenchmarkDecodeArray(b *testing.B) { test.BenchmarkDecode(b, arrayDecoder, arrayReader) }

var (
	mapReader  = test.NewBenchmarkReader(`{"A":1,"B":2,"C":3}`)
	mapDecoder = json.NewDecoder(mapReader)
)

func BenchmarkDecodeMap(b *testing.B) { test.BenchmarkDecode(b, mapDecoder, mapReader) }
