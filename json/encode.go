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
	m := marshalerPool.Get().(*marshaler)
	e := objconv.Encoder{
		Emitter: m,
	}

	if err = e.Encode(v); err == nil {
		b = m.bytes()
	} else {
		m.reset()
	}

	marshalerPool.Put(m)
	return
}

var marshalerPool = sync.Pool{
	New: func() interface{} { return newMarshaler() },
}

type marshaler struct {
	Emitter
	b bytes.Buffer
}

func newMarshaler() *marshaler {
	m := &marshaler{}
	m.w = &m.b
	return m
}

func (m *marshaler) bytes() (b []byte) {
	b, m.b = m.b.Bytes(), bytes.Buffer{}
	return
}

func (m *marshaler) reset() {
	m.b.Reset()
}
