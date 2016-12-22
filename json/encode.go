package json

import (
	"bytes"
	"io"
	"sync"

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
	buf := bufferPool.Get().(*bytes.Buffer)
	emt := emitterPool.Get().(*Emitter)
	emt.Reset(buf)

	enc := objconv.Encoder{
		Emitter: emt,
	}

	if err = enc.Encode(v); err == nil {
		b, *buf = buf.Bytes(), bytes.Buffer{}
	}

	emt.Reset(nil)
	emitterPool.Put(emt)
	bufferPool.Put(buf)
	return
}

var bufferPool = sync.Pool{
	New: func() interface{} { return &bytes.Buffer{} },
}

var emitterPool = sync.Pool{
	New: func() interface{} { return &Emitter{} },
}
