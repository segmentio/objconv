package json

import (
	"io"

	"github.com/segmentio/objconv"
)

func init() {
	codec := objconv.Codec{
		NewEmitter: func(w io.Writer) objconv.Emitter { return NewEmitter(w) },
		NewParser:  func(r io.Reader) objconv.Parser { return NewParser(r) },
	}

	for _, name := range [...]string{
		"application/json",
		"text/json",
		"json",
	} {
		objconv.Register(name, codec)
	}
}
