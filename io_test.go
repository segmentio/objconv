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

func TestWriterPanic(t *testing.T) {
	defer func() {
		if x := recover(); x != io.EOF {
			t.Errorf("expected error to be reported in panic bu %#v was found", x)
		}
	}()
	w := Writer{W: errorWriter{io.EOF}}
	w.WriteString("Hello World!")
}

type errorWriter struct{ error }

func (w errorWriter) Write(b []byte) (int, error) { return 0, w.error }
