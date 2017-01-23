package mail

import (
	"net/mail"
	"reflect"

	"github.com/segmentio/objconv"
)

func encodeAddress(e objconv.Encoder, v reflect.Value) error {
	a := v.Interface().(mail.Address)
	return e.Emitter.EmitString(a.String())
}
