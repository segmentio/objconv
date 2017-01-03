package cbor

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

func TestMajorType(t *testing.T) {
	m, b := majorType(majorByte(majorType7, 24))

	if m != majorType7 {
		t.Error("bad major type:", m)
	}

	if b != 24 {
		t.Error("bad info value:", b)
	}
}
