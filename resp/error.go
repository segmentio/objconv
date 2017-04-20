package resp

// The Error type represents redis errors.
type Error struct {
	raw string
}

// NewError returns a new redis error.
func NewError(s string) *Error {
	return &Error{
		raw: s,
	}
}

// Error satsifies the error interface.
func (e *Error) Error() string {
	return e.raw
}
