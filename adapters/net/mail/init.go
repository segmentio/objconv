package mail

import (
	"net/mail"
	"reflect"

	"github.com/segmentio/objconv"
)

func init() {
	objconv.Install(reflect.TypeOf(mail.Address{}), AddressAdapter())
}

// AddressAdapter returns the adapter to encode and decode mail.Address values.
func AddressAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeAddress,
		Decode: decodeAddress,
	}
}
