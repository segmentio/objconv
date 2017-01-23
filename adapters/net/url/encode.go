package url

import (
	"net/url"
	"reflect"

	"github.com/segmentio/objconv"
)

func encodeURL(e objconv.Encoder, v reflect.Value) error {
	u := v.Interface().(url.URL)
	return e.Emitter.EmitString(u.String())
}
