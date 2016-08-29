package objconv

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

const longString = `Package json implements encoding and decoding of JSON objects as defined in RFC 4627. The mapping between JSON objects and Go values is described in the documentation for the Marshal and Unmarshal functions.`

func TestReadLine(t *testing.T) {
	tests := []struct {
		data  string
		size  int // size of the buffer after reading lines
		lines []string
	}{
		{
			data:  "Hello World!\n",
			size:  1024,
			lines: []string{"Hello World!"},
		},
		{
			data: "0) " + longString +
				"1) " + longString +
				"2) " + longString +
				"3) " + longString +
				"4) " + longString +
				"5) " + longString +
				"6) " + longString +
				"7) " + longString +
				"8) " + longString +
				"9) " + longString + "\n",
			size: 4096,
			lines: []string{
				"0) " + longString +
					"1) " + longString +
					"2) " + longString +
					"3) " + longString +
					"4) " + longString +
					"5) " + longString +
					"6) " + longString +
					"7) " + longString +
					"8) " + longString +
					"9) " + longString,
			},
		},
	}

testLoop:
	for i, test := range tests {
		var lines []string
		var r = NewReader(strings.NewReader(test.data))

	readLines:
		for {
			switch line, err := r.ReadLine(LF); err {
			case nil:
				lines = append(lines, string(line))
			case io.EOF:
				break readLines
			default:
				t.Errorf("[%d] %s", i, err)
				continue testLoop
			}
		}

		if !reflect.DeepEqual(lines, test.lines) {
			t.Errorf("[%d] bad lines: %v", i, lines)
		}

		if len(r.b) != test.size {
			t.Errorf("[%d] bad buffer size: %d != %d", i, test.size, len(r.b))
		}
	}
}

func TestWriteString(t *testing.T) {
	b := &bytes.Buffer{}
	w := Writer{w: b}

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
	w := Writer{w: errorWriter{io.ErrUnexpectedEOF}}
	_, err := w.WriteString("Hello World!")

	if err != io.ErrUnexpectedEOF {
		t.Error("bad error:", err)
	}
}

type errorWriter struct{ error }

func (w errorWriter) Write(b []byte) (int, error) { return 0, w.error }
