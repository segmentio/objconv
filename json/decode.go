package json

import (
	"bytes"
	"io"
	"sync"

	"github.com/segmentio/objconv"
)

// NewDecoder returns a new JSON deocder that parses values from r.
func NewDecoder(r io.Reader) *objconv.Decoder {
	return &objconv.Decoder{
		Parser: NewParser(r),
	}
}

// Unmarshal decodes a JSON representation of v from b.
func Unmarshal(b []byte, v interface{}) (err error) {
	// Get a reader from the pool, this saves a memory allocation because Go
	// fails to realize that the reader doesn't escape.
	r := readerPool.Get().(*bytes.Reader)
	r.Reset(b)

	// Get a parser from the pool, this saves a memory allocation because Go
	// fails to realize that the parser doesn't escape.
	p := parserPool.Get().(*Parser)
	p.Reset(r)

	err = (objconv.Decoder{Parser: p}).Decode(v)

	p.Reset(nil)
	parserPool.Put(p)

	r.Reset(nil)
	readerPool.Put(r)
	return
}

var readerPool = sync.Pool{
	New: func() interface{} { return bytes.NewReader(nil) },
}

var parserPool = sync.Pool{
	New: func() interface{} { return &Parser{} },
}
