package objconv

import (
	"reflect"
	"time"
)

// The Parser interface must be implemented by types that provide decoding of a
// specific format (like json, resp, ...).
type Parser interface {
	// ParseType is called by a decoder to ask the parser what is the type of
	// the next value that can be parsed.
	ParseType() (Type, error)

	// ParseNil parses a nil value.
	ParseNil() error

	// ParseBool parses a boolean value.
	ParseBool() (bool, error)

	// ParseInt parses an integer value.
	ParseInt() (int64, error)

	// ParseBool parses an unsigned integer value.
	ParseUint() (uint64, error)

	// ParseBool parses a floating point value.
	ParseFloat() (float64, error)

	// ParseBool parses a string value.
	ParseString() (string, error)

	// ParseBool parses a byte array value.
	ParseBytes() ([]byte, error)

	// ParseBool parses a time value.
	ParseTime() (time.Time, error)

	// ParseBool parses a duration value.
	ParseDuration() (time.Duration, error)

	// ParseError parses an error value.
	ParseError() (error, error)

	// ParseArrayBegin is called by the array-decoding algorithm when it starts.
	//
	// The method should return the length of the array being decoded, or a
	// negative value if it is unknown (some formats like json don't keep track
	// of the length of the array).
	ParseArrayBegin() (int, error)

	// ParseArrayEnd is called by the array-decoding algorithm when it
	// completes.
	ParseArrayEnd() error

	// ParseArrayNext is called by the array-decoding algorithm between each
	// value parsed in the array.
	//
	// If the ParseArrayBegin method returned a negative value this method
	// should return objconv.End to indicated that there is no more elements to
	// parse in the array.
	ParseArrayNext() error

	// ParseMapBegin is called by the map-decoding algorithm when it starts.
	//
	// The method should return the length of the map being decoded, or a
	// negative value if it is unknown (some formats like json don't keep track
	// of the length of the map).
	ParseMapBegin() (int, error)

	// ParseMapEnd is called by the map-decoding algorithm when it completes.
	ParseMapEnd() error

	// ParseMapValue is called by the map-decoding algorithm after parsing a key
	// but before parsing the associated value.
	ParseMapValue() error

	// ParseMapNext is called by the map-decoding algorithm between each
	// value parsed in the map.
	//
	// If the ParseMapBegin method returned a negative value this method should
	// return objconv.End to indicated that there is no more elements to parse
	// in the map.
	ParseMapNext() error
}

type ValueParser struct {
	stack []valueParserContext
}

type valueParserContext struct {
	value reflect.Value
	index int
	keys  []reflect.Value
}

func NewValueParser(v interface{}) *ValueParser {
	return &ValueParser{
		stack: []valueParserContext{
			makeValueParserContext(reflect.ValueOf(v)),
		},
	}
}

func (p *ValueParser) ParseType() (t Type, err error) {
	return
}

func (p *ValueParser) ParseNil() (err error) {
	return
}

func (p *ValueParser) ParseBool() (v bool, err error) {
	return
}

func (p *ValueParser) ParseInt() (v int64, err error) {
	return
}

func (p *ValueParser) ParseUint() (v uint64, err error) {
	return
}

func (p *ValueParser) ParseFloat() (v float64, err error) {
	return
}

func (p *ValueParser) ParseString() (v string, err error) {
	return
}

func (p *ValueParser) ParseBytes() (v []byte, err error) {
	return
}

func (p *ValueParser) ParseTime() (v time.Time, err error) {
	return
}

func (p *ValueParser) ParseDuration() (v time.Duration, err error) {
	return
}

func (p *ValueParser) ParseError() (v error, err error) {
	return
}

func (p *ValueParser) ParseArrayBegin() (n int, err error) {
	return
}

func (p *ValueParser) ParseArrayEnd() (err error) {
	return
}

func (p *ValueParser) ParseArrayNext() (err error) {
	return
}

func (p *ValueParser) ParseMapBegin() (n int, err error) {
	return
}

func (p *ValueParser) ParseMapEnd() (err error) {
	return
}

func (p *ValueParser) ParseMapValue() (err error) {
	return
}

func (p *ValueParser) ParseMapNext() (err error) {
	return
}

func makeValueParserContext(v reflect.Value) (ctx valueParserContext) {
	return
}
