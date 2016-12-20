package json

import (
	"bytes"
	"io"

	"github.com/segmentio/objconv"
)

func NewDecoder(r io.Reader) *objconv.Decoder {
	return objconv.NewDecoder(NewParser(r))
}

func Unmarshal(b []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(b)).Decode(v)
}
