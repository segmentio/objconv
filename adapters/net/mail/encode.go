package mail

import (
	"net/mail"
	"reflect"

	"github.com/segmentio/objconv"
)

// EncodeAddress encodes the mail.Address value in v as a string using e.
func EncodeAddress(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(mail.Address)
	return e.Emitter.EmitString(a.String())
}
