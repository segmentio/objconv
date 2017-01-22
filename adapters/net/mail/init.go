package mail

import (
	"net/mail"
	"reflect"

	"github.com/segmentio/objconv"
)

func init() {
	objconv.Install(reflect.TypeOf(mail.Address{}), objconv.Adapter{
		Encode: EncodeAddress,
		Decode: DecodeAddress,
	})
}
