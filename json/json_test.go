package json

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

func TestUnicode(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{`"\u2022"`, "•"},
		{`"\uDC00D800"`, "�"},
	}

	for _, test := range tests {
		t.Run(test.out, func(t *testing.T) {
			var s string

			if err := Unmarshal([]byte(test.in), &s); err != nil {
				t.Error(err)
			}

			if s != test.out {
				t.Error(s)
			}
		})
	}
}
