package msgpack

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
		"application/msgpack",
		"msgpack",
	} {
		mimetype.Register(name, codec)
	}
}
