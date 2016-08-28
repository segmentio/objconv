package resp

import (
	"testing"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/test"
)

var (
	nilReader  = test.NewBenchmarkReader([]byte("$-1\r\n"))
	nilDecoder = objconv.NewDecoder(objconv.DecoderConfig{
		Input:  nilReader,
		Parser: &Parser{},
	})
)

func BenchmarkDecodeNil(b *testing.B) { test.BenchmarkDecode(b, nilDecoder, nilReader) }
