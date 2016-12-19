package objconv

import (
	"errors"
	"fmt"
	"reflect"
)

// ArrayLengthError is returned when decoding into an array produced an invalid
// number of elements.
type ArrayLengthError struct {
	// Type of the array.
	Type reflect.Type

	// Length that didn't match the number of elements in the array.
	Length int
}

// Error satisfies the error interface.
func (e *ArrayLengthError) Error() string {
	return fmt.Sprintf("objconv: length = %d: array length mismatch with %s", e.Length, e.Type)
}

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
	From Type

	// To is the destination type where the value is decoded.
	To Type
}

// Error satsifies the error interface.
func (e *TypeConversionError) Error() string {
	return "objconv: cannot convert " + e.From.String() + " to " + e.To.String()
}

var (
	// End is expected to be returned to indicate that a function has completed
	// its owrk, this is usually employed in generic algorithms.
	End = errors.New("end")

	// This error value is used as a building block for reflection and is never
	// returned by the package.
	errBase = errors.New("")
)
