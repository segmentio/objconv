package url

import (
	"net/url"
	"reflect"

	"github.com/segmentio/objconv"
)

func init() {
	objconv.Install(reflect.TypeOf(url.URL{}), objconv.Adapter{
		Encode: EncodeURL,
		Decode: DecodeURL,
	})
}
