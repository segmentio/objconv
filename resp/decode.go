package resp

import (
	"bytes"
	"io"
	"sync"

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
	u := unmarshalerPool.Get().(*unmarshaler)
	u.reset(b)

	err := (objconv.Decoder{Parser: u}).Decode(v)

	u.reset(nil)
	unmarshalerPool.Put(u)
	return err
}

var unmarshalerPool = sync.Pool{
	New: func() interface{} { return newUnmarshaler() },
}

type unmarshaler struct {
	Parser
	b bytes.Buffer
}

func newUnmarshaler() *unmarshaler {
	u := &unmarshaler{}
	u.s = u.c[:0]
	u.r = &u.b
	return u
}

func (u *unmarshaler) reset(b []byte) {
	v := bytes.NewBuffer(b)
	u.b = *v
}
