package mimetype

import (
	"io"
	"sync"

	"github.com/segmentio/objconv"
)

// A Registry associates mime types to codecs.
//
// It is safe to use a registry concurrently from multiple goroutines.
type Registry struct {
	mutex  sync.RWMutex
	codecs map[string]Codec
}

// Register adds a codec for a mimetype to r.
func (reg *Registry) Register(mimetype string, codec Codec) {
	defer reg.mutex.Unlock()
	reg.mutex.Lock()

	if reg.codecs == nil {
		reg.codecs = make(map[string]Codec)
	}

	reg.codecs[mimetype] = codec
}

// Unregister removes the codec for a mimetype from r.
func (reg *Registry) Unregister(mimetype string) {
	defer reg.mutex.Unlock()
	reg.mutex.Lock()

	delete(reg.codecs, mimetype)
}

// Lookup returns the codec associated with mimetype, ok is set to true or false
// based on whether a codec was found.
func (reg *Registry) Lookup(mimetype string) (codec Codec, ok bool) {
	reg.mutex.RLock()
	codec, ok = reg.codecs[mimetype]
	reg.mutex.RUnlock()
	return
}

// NewEncoder returns a new encoder for mimetype that outputs to w.
//
// The function returns nil if non codec was registered for mimetype.
func (reg *Registry) NewEncoder(mimetype string, w io.Writer) *objconv.Encoder {
	codec, ok := reg.Lookup(mimetype)
	if !ok {
		return nil
	}
	return codec.NewEncoder(w)
}

// NewDecoder returns a new encoder for mimetype that takes input from r.
//
// The function returns nil if non codec was registered for mimetype.
func (reg *Registry) NewDecoder(mimetype string, r io.Reader) *objconv.Decoder {
	codec, ok := reg.Lookup(mimetype)
	if !ok {
		return nil
	}
	return codec.NewDecoder(r)
}

// NewStreamEncoder returns a new encoder for mimetype that outputs to w.
//
// The function returns nil if non codec was registered for mimetype.
func (reg *Registry) NewStreamEncoder(mimetype string, w io.Writer) *objconv.StreamEncoder {
	codec, ok := reg.Lookup(mimetype)
	if !ok {
		return nil
	}
	return codec.NewStreamEncoder(w)
}

// NewStreamDecoder returns a new encoder for mimetype that takes input from r.
//
// The function returns nil if non codec was registered for mimetype.
func (reg *Registry) NewStreamDecoder(mimetype string, r io.Reader) *objconv.StreamDecoder {
	codec, ok := reg.Lookup(mimetype)
	if !ok {
		return nil
	}
	return codec.NewStreamDecoder(r)
}
