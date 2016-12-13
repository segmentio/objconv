package objconv

import (
	"errors"
	"fmt"
	"reflect"
)

// OutOfBoundsError is returned when decoding numeric values that do not fit
// into their destination type, for example decoding the value 1000 into a int8
// would return this error.
type OutOfBoundsError struct {
	// Value is the decoded value.
	Value interface{}

	// Type of the destination.
	Type reflect.Type
}

// Error satisfies the error interface.
func (e *OutOfBoundsError) Error() string {
	return "objconv: value out of bounds for " + e.Type.String() + ": " + fmt.Sprint(e.Value)
}

// UnsupportedFormatError is returned when the format specified for an encoding
// or decoding operation is not recognized.
type UnsupportedFormatError struct {
	// Format is the name of the unsupported format.
	Format string
}

// Error satisfies the error interface.
func (e *UnsupportedFormatError) Error() string {
	return "objconv: unsupported format: " + e.Format
}

// UnsupportedTypeError is returned by encoding functions when attempting to
// encode an unsupported value type.
type UnsupportedTypeError struct {
	// Type is the unsupported type.
	Type reflect.Type
}

// Error satisfies the error interface.
func (e *UnsupportedTypeError) Error() string {
	return "objconv: unsupported type: " + e.Type.String()
}

// TypeConversionError is returned by decoding functions when a a type mismatch
// occurs between the decoded value and its destination.
type TypeConversionError struct {
	// From is the type of the value being decoded.
	From reflect.Type

	// To is the destination type where the value is decoded.
	To reflect.Type
}

// Error satsifies the error interface.
func (e *TypeConversionError) Error() string {
	return "objconv: type mismatch between " + e.From.String() + " and " + e.To.String()
}

var (
	// End is expected to be returned to indicate that a function has completed
	// its owrk, this is usually employed in generic algorithms.
	End = errors.New("end")
)
