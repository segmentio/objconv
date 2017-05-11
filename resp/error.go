package resp

// The Error type represents redis errors.
type Error string

// NewError returns a new redis error.
func NewError(s string) *Error {
	e := Error(s)
	return &e
}

// Error satsifies the error interface.
func (e *Error) Error() string {
	return string(*e)
}
