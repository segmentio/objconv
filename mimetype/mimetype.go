package mimetype

import (
	"io"

	"github.com/segmentio/objconv"
)

// The global registry to which packages add their codecs.
var reg Registry

// Register adds a codec for a mimetype to the global registry.
func Register(mimetype string, codec Codec) {
	reg.Register(mimetype, codec)
}

// Unregister removes the codec for a mimetype from the global registry.
func Unregister(mimetype string) {
	reg.Unregister(mimetype)
}

// Lookup returns the codec associated with mimetype, ok is set to true or false
// based on whether a codec was found.
func Lookup(mimetype string) (Codec, bool) {
	return reg.Lookup(mimetype)
}

// NewEncoder returns a new encoder for mimetype that outputs to w.
//
// The function returns nil if non codec was registered for mimetype.
func NewEncoder(mimetype string, w io.Writer) *objconv.Encoder {
	return reg.NewEncoder(mimetype, w)
}

// NewDecoder returns a new encoder for mimetype that takes input from r.
//
// The function returns nil if non codec was registered for mimetype.
func NewDecoder(mimetype string, r io.Reader) *objconv.Decoder {
	return reg.NewDecoder(mimetype, r)
}

// NewStreamEncoder returns a new encoder for mimetype that outputs to w.
//
// The function returns nil if non codec was registered for mimetype.
func NewStreamEncoder(mimetype string, w io.Writer) *objconv.StreamEncoder {
	return reg.NewStreamEncoder(mimetype, w)
}

// NewStreamDecoder returns a new encoder for mimetype that takes input from r.
//
// The function returns nil if non codec was registered for mimetype.
func NewStreamDecoder(mimetype string, r io.Reader) *objconv.StreamDecoder {
	return reg.NewStreamDecoder(mimetype, r)
}
