package resp

import (
	"strings"
	"unicode"
)

// Error represents a redis error returned by redis servers.
type Error struct {
	// Type of the error, often something like "ERR" or "WRONGTYPE".
	Type string

	// Message is the human readable error message.
	Message string
}

// NewError creates a new error value from a string which is expected to be
// formatted following redis error conventions ("<TYPE> <message>").
func NewError(s string) *Error {
	var t string
	var m = s
	var i = strings.IndexByte(s, ' ')
	var j = i + 1

	if i < 0 {
		i = len(s)
		j = i
	}

	if i > 0 && isErrorPrefix(s[:i]) {
		t, m = s[:i], s[j:]
	}

	return &Error{
		Type:    t,
		Message: m,
	}
}

// Error satisfies the error interface.
func (e *Error) Error() string {
	if len(e.Message) == 0 {
		return e.Type
	}

	if len(e.Type) == 0 {
		return e.Message
	}

	return e.Type + " " + e.Message
}

func isErrorPrefix(s string) bool {
	for _, c := range s {
		if !unicode.IsUpper(c) {
			return false
		}
	}
	return len(s) != 0
}
