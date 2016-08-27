package resp

import (
	"testing"

	"github.com/segmentio/objconv/test"
)

func BenchmarkEmitBoolTrue(b *testing.B) { test.BenchmarkEmitBoolTrue(b, &Emitter{}) }

func BenchmarkEmitBoolFalse(b *testing.B) { test.BenchmarkEmitBoolFalse(b, &Emitter{}) }
