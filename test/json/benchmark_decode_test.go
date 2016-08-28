package json

import (
	"encoding/json"
	"testing"

	"github.com/segmentio/objconv/test"
)

var (
	nilReader  = test.NewBenchmarkReader([]byte("null"))
	nilDecoder = json.NewDecoder(nilReader)
)

func BenchmarkDecodeNil(b *testing.B) { test.BenchmarkDecode(b, nilDecoder, nilReader) }
