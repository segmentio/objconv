package yaml

import (
	"testing"

	"github.com/segmentio/objconv/objtests"
)

func TestCodec(t *testing.T) {
	objtests.TestCodec(t, codec)
}

func BenchmarkCodec(b *testing.B) {
	objtests.BenchmarkCodec(b, codec)
}
