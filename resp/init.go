package resp

import (
	"io"

	"github.com/segmentio/objconv"
)

// Codec for the RESP format.
var Codec = objconv.Codec{
	NewEmitter: func(w io.Writer) objconv.Emitter { return NewEmitter(w) },
	NewParser:  func(r io.Reader) objconv.Parser { return NewParser(r) },
}

func init() {
	for _, name := range [...]string{
		"application/resp",
		"text/resp",
		"resp",
	} {
		objconv.Register(name, Codec)
	}
}
