package url

import (
	"errors"
	"net/url"
	"reflect"

	"github.com/segmentio/objconv"
)

// DecodeURL decodes a url.URL value into to from a string representation
// using d.
func DecodeURL(d objconv.Decoder, to reflect.Value) (err error) {
	var u *url.URL
	var s string

	if err = d.Decode(&s); err != nil {
		return
	}

	if u, err = url.Parse(s); err != nil {
		err = errors.New("objconv: bad URL: " + err.Error())
		return
	}

	if to.IsValid() {
		to.Set(reflect.ValueOf(*u))
	}
	return
}
