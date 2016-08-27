package objconv

import (
	"bytes"
	"io"
	"testing"
)

func TestWriter(t *testing.T) {
	b := &bytes.Buffer{}
	w := Writer{W: b}

	n, e := w.WriteString("Hello World!")

	if n != 12 {
		t.Errorf("invalid number of bytes written: %d", n)
	}

	if e != nil {
		t.Errorf("unexpected error returned: %s", e)
	}

	if s := b.String(); s != "Hello World!" {
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
