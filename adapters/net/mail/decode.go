package mail

import (
	"errors"
	"net/mail"
	"reflect"

	"github.com/segmentio/objconv"
)

// DecodeAddress decodes a mail.Address value into to from a string
// representation using d.
func DecodeAddress(d objconv.Decoder, to reflect.Value) (err error) {
	var a *mail.Address
	var s string

	if err = d.Decode(&s); err != nil {
		return
	}

	if a, err = mail.ParseAddress(s); err != nil {
		err = errors.New("objconv: bad email address: " + err.Error())
		return
	}

	if to.IsValid() {
		to.Set(reflect.ValueOf(*a))
	}
	return
}