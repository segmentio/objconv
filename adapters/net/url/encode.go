package url

import (
	"net/url"
	"reflect"

	"github.com/segmentio/objconv"
)

// EncodeURL encodes the url.URL value in v as a string using e.
func EncodeURL(e objconv.Encoder, v reflect.Value) error {
	u := v.Interface().(url.URL)
	return e.Emitter.EmitString(u.String())
}
