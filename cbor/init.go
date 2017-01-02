package cbor

import (
	"io"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/mimetype"
)

func init() {
	codec := mimetype.Codec{
		NewEmitter: func(w io.Writer) objconv.Emitter { return NewEmitter(w) },
		NewParser:  func(r io.Reader) objconv.Parser { return NewParser(r) },
	}

	for _, name := range [...]string{
		"application/cbor",
		"cbor",
	} {
		mimetype.Register(name, codec)
	}
}
