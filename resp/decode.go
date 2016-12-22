package resp

import (
	"bytes"
	"io"

	"github.com/segmentio/objconv"
)

// NewDecoder returns a new RESP decoder that parses values from r.
func NewDecoder(r io.Reader) *objconv.Decoder {
	return objconv.NewDecoder(NewParser(r))
}

// NewStreamDecoder returns a new RESP stream decoder that parses values from r.
func NewStreamDecoder(r io.Reader) *objconv.StreamDecoder {
	return objconv.NewStreamDecoder(NewParser(r))
}

// Unmarshal decodes a RESP representation of v from b.
func Unmarshal(b []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(b)).Decode(v)
}
