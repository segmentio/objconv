package resp

import (
	"testing"

	"github.com/segmentio/objconv/test"
)

var (
	benchmarkParser = &Parser{}

	nilReader = test.NewBenchmarkReader([]byte("$-1\r\n"))
)

func BenchmarkDecodeNil(b *testing.B) { test.BenchmarkDecode(b, benchmarkParser, nilReader) }
