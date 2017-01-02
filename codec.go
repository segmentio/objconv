package objconv

import (
	"io"
)

// A Codec is a factory for encoder and decoders that work on byte streams.
type Codec struct {
	NewEmitter func(io.Writer) Emitter
	NewParser  func(io.Reader) Parser
}

// NewEncoder returns a new encoder that outputs to w.
func (c Codec) NewEncoder(w io.Writer) *Encoder {
	return NewEncoder(c.NewEmitter(w))
}

// NewDecoder returns a new decoder that takes input from r.
func (c Codec) NewDecoder(r io.Reader) *Decoder {
	return NewDecoder(c.NewParser(r))
}

// NewStreamEncoder returns a new stream encoder that outputs to w.
func (c Codec) NewStreamEncoder(w io.Writer) *StreamEncoder {
	return NewStreamEncoder(c.NewEmitter(w))
}

// NewStreamDecoder returns a new stream decoder that takes input from r.
func (c Codec) NewStreamDecoder(r io.Reader) *StreamDecoder {
	return NewStreamDecoder(c.NewParser(r))
}
