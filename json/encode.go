package json

import (
	"bytes"
	"io"

	"github.com/segmentio/objconv"
)

// NewEncoder returns a new JSON encoder that writes to w.
func NewEncoder(w io.Writer) *objconv.Encoder {
	return objconv.NewEncoder(NewEmitter(w))
}

// NewStreamEncoder returns a new JSON stream encoder that writes to w.
func NewStreamEncoder(w io.Writer) *objconv.StreamEncoder {
	return objconv.NewStreamEncoder(NewEmitter(w))
}

// Marshal writes the JSON representation of v to a byte slice returned in b.
func Marshal(v interface{}) (b []byte, err error) {
	buf := &bytes.Buffer{}
	buf.Grow(1024)

	if err = NewEncoder(buf).Encode(v); err == nil {
		b = buf.Bytes()
	}

	return
}
