package url

import (
	"net/url"
	"reflect"

	"github.com/segmentio/objconv"
)

func init() {
	objconv.Install(reflect.TypeOf(url.URL{}), URLAdapter())
}

// URLAdapter returns the adapter to encode and decode url.URL values.
func URLAdapter() objconv.Adapter {
	return objconv.Adapter{
		Encode: encodeURL,
		Decode: decodeURL,
	}
}
