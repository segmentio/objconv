package mail

import (
	"net/mail"
	"reflect"

	"github.com/segmentio/objconv"
)

func init() {
	objconv.Install(reflect.TypeOf(mail.Address{}), AddressAdapter())
	objconv.Install(reflect.TypeOf(([]*mail.Address)(nil)), AddressListAdapter())
}

// AddressAdapter returns the adapter to encode and decode mail.Address values.
func AddressAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeAddress,
		Decode: decodeAddress,
	}
}

// AddressListAdapter returns the adapter to encode and decode []*mail.Address
// values.
func AddressListAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeAddressList,
		Decode: decodeAddressList,
	}
}
