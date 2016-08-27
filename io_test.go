package objconv

import (
	"bytes"
	"io"
	"testing"
)

func TestWriteString(t *testing.T) {
	const longString = `Package json implements encoding and decoding of JSON objects as defined in RFC 4627. The mapping between JSON objects and Go values is described in the documentation for the Marshal and Unmarshal functions.`

	b := &bytes.Buffer{}
	w := Writer{W: b}

	n, e := w.WriteString(longString)

	if n != len(longString) {
		t.Errorf("invalid number of bytes written: %d", n)
	}

	if e != nil {
		t.Errorf("unexpected error returned: %s", e)
	}

	if s := b.String(); s != longString {
		t.Errorf("invalid buffer content: %#v", s)
	}
}

func TestWriterError(t *testing.T) {
	w := Writer{W: errorWriter{io.ErrUnexpectedEOF}}
	_, err := w.WriteString("Hello World!")

	if err != io.ErrUnexpectedEOF {
		t.Error("bad error:", err)
	}
}

type errorWriter struct{ error }

func (w errorWriter) Write(b []byte) (int, error) { return 0, w.error }
