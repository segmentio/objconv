package mimetype

import (
	"io"

	"github.com/segmentio/objconv"
)

// A Codec is a factory for encoder and decoders that work on byte streams.
type Codec struct {
	NewEmitter func(io.Writer) objconv.Emitter
	NewParser  func(io.Reader) objconv.Parser
}

// NewEncoder returns a new encoder that outputs to w.
func (c Codec) NewEncoder(w io.Writer) *objconv.Encoder {
	return objconv.NewEncoder(c.NewEmitter(w))
}

// NewDecoder returns a new decoder that takes input from r.
func (c Codec) NewDecoder(r io.Reader) *objconv.Decoder {
	return objconv.NewDecoder(c.NewParser(r))
}

// NewStreamEncoder returns a new stream encoder that outputs to w.
func (c Codec) NewStreamEncoder(w io.Writer) *objconv.StreamEncoder {
	return objconv.NewStreamEncoder(c.NewEmitter(w))
}

// NewStreamDecoder returns a new stream decoder that takes input from r.
func (c Codec) NewStreamDecoder(r io.Reader) *objconv.StreamDecoder {
	return objconv.NewStreamDecoder(c.NewParser(r))
}
