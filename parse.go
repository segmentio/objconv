package objconv

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// The Parser interface must be implemented by types that provide decoding of a
// specific format (like json, resp, ...).
//
// Parsers are not expected to be safe for use by multiple goroutines.
type Parser interface {
	// ParseType is called by a decoder to ask the parser what is the type of
	// the next value that can be parsed.
	//
	// ParseType must be idempotent, it must be possible to call it multiple
	// without actually changing the state of the parser.
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
	//
	// The string is returned as a byte slice because it is expected to be
	// pointing at an internal memory buffer, the decoder will make a copy of
	// the value. This design allows more memory allocation optimizations.
	ParseString() ([]byte, error)

	// ParseBool parses a byte array value.
	//
	// The returned byte slice is expected to be pointing at an internal memory
	// buffer, the decoder will make a copy of the value. This design allows more
	// memory allocation optimizations.
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
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the array.
	ParseArrayEnd(int) error

	// ParseArrayNext is called by the array-decoding algorithm between each
	// value parsed in the array.
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the array.
	//
	// If the ParseArrayBegin method returned a negative value this method
	// should return objconv.End to indicated that there is no more elements to
	// parse in the array. In this case the method is also called right before
	// decoding the first element ot handle the case where the array is empty
	// and the end-of-array marker can be read right away.
	ParseArrayNext(int) error

	// ParseMapBegin is called by the map-decoding algorithm when it starts.
	//
	// The method should return the length of the map being decoded, or a
	// negative value if it is unknown (some formats like json don't keep track
	// of the length of the map).
	ParseMapBegin() (int, error)

	// ParseMapEnd is called by the map-decoding algorithm when it completes.
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the map.
	ParseMapEnd(int) error

	// ParseMapValue is called by the map-decoding algorithm after parsing a key
	// but before parsing the associated value.
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the map.
	ParseMapValue(int) error

	// ParseMapNext is called by the map-decoding algorithm between each
	// value parsed in the map.
	//
	// The method receives the iteration counter as argument, which indicates
	// how many values were decoded from the map.
	//
	// If the ParseMapBegin method returned a negative value this method should
	// return objconv.End to indicated that there is no more elements to parse
	// in the map. In this case the method is also called right before decoding
	// the first element ot handle the case where the array is empty and the
	// end-of-map marker can be read right away.
	ParseMapNext(int) error
}

// The bytesDecoder interface may optionnaly be implemented by a Parser to
// provide an extra step in decoding a byte slice. This is sometimes necessary
// if the associated Emitter has transformed bytes slices because the format is
// not capable of representing binary data.
type bytesDecoder interface {
	// DecodeBytes is called when the destination variable for a string or a
	// byte slice is a byte slice, allowing the parser to apply a transformation
	// before the value is stored.
	DecodeBytes([]byte) ([]byte, error)
}

// ValueParser is parser that uses "natural" in-memory representation of data
// structures.
//
// This is mainly useful for testing the decoder algorithms.
type ValueParser struct {
	stack []reflect.Value
	ctx   []valueParserContext
}

type valueParserContext struct {
	value  reflect.Value
	keys   []reflect.Value
	fields []structField
}

// NewValueParser creates a new parser that exposes the value v.
func NewValueParser(v interface{}) *ValueParser {
	return &ValueParser{
		stack: []reflect.Value{reflect.ValueOf(v)},
	}
}

func (p *ValueParser) ParseType() (Type, error) {
	v := p.value()

	if !v.IsValid() {
		return Nil, nil
	}

	switch v.Interface().(type) {
	case time.Time:
		return Time, nil

	case time.Duration:
		return Duration, nil

	case error:
		return Error, nil
	}

	switch v.Kind() {
	case reflect.Bool:
		return Bool, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Uint, nil

	case reflect.Float32, reflect.Float64:
		return Float, nil

	case reflect.String:
		return String, nil

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return Bytes, nil
		}
		return Array, nil

	case reflect.Array:
		return Array, nil

	case reflect.Map:
		return Map, nil

	case reflect.Struct:
		return Map, nil

	case reflect.Interface:
		if v.IsNil() {
			return Nil, nil
		}
	}

	return Nil, errors.New("objconv: unsupported type found in value parser: " + v.Type().String())
}

func (p *ValueParser) ParseNil() (err error) {
	return
}

func (p *ValueParser) ParseBool() (v bool, err error) {
	v = p.value().Bool()
	return
}

func (p *ValueParser) ParseInt() (v int64, err error) {
	v = p.value().Int()
	return
}

func (p *ValueParser) ParseUint() (v uint64, err error) {
	v = p.value().Uint()
	return
}

func (p *ValueParser) ParseFloat() (v float64, err error) {
	v = p.value().Float()
	return
}

func (p *ValueParser) ParseString() (v []byte, err error) {
	v = []byte(p.value().String())
	return
}

func (p *ValueParser) ParseBytes() (v []byte, err error) {
	v = p.value().Bytes()
	return
}

func (p *ValueParser) ParseTime() (v time.Time, err error) {
	v = p.value().Interface().(time.Time)
	return
}

func (p *ValueParser) ParseDuration() (v time.Duration, err error) {
	v = p.value().Interface().(time.Duration)
	return
}

func (p *ValueParser) ParseError() (v error, err error) {
	v = p.value().Interface().(error)
	return
}

func (p *ValueParser) ParseArrayBegin() (n int, err error) {
	v := p.value()
	n = v.Len()
	p.pushContext(valueParserContext{value: v})

	if n != 0 {
		p.push(v.Index(0))
	}

	return
}

func (p *ValueParser) ParseArrayEnd(n int) (err error) {
	if n != 0 {
		p.pop()
	}
	p.popContext()
	return
}

func (p *ValueParser) ParseArrayNext(n int) (err error) {
	ctx := p.context()
	p.pop()
	p.push(ctx.value.Index(n))
	return
}

func (p *ValueParser) ParseMapBegin() (n int, err error) {
	v := p.value()

	if v.Kind() == reflect.Map {
		n = v.Len()
		k := v.MapKeys()
		p.pushContext(valueParserContext{value: v, keys: k})
		if n != 0 {
			p.push(k[0])
		}
	} else {
		c := valueParserContext{value: v}
		s := structCache.lookup(v.Type())

		for _, f := range s.fields {
			if !f.omit(v.FieldByIndex(f.index)) {
				c.fields = append(c.fields, f)
				n++
			}
		}

		p.pushContext(c)
		if n != 0 {
			p.push(reflect.ValueOf(c.fields[0].name))
		}
	}

	return
}

func (p *ValueParser) ParseMapEnd(n int) (err error) {
	if n != 0 {
		p.pop()
	}
	p.popContext()
	return
}

func (p *ValueParser) ParseMapValue(n int) (err error) {
	ctx := p.context()
	p.pop()

	if ctx.keys != nil {
		p.push(ctx.value.MapIndex(ctx.keys[n]))
	} else {
		p.push(ctx.value.FieldByIndex(ctx.fields[n].index))
	}

	return
}

func (p *ValueParser) ParseMapNext(n int) (err error) {
	ctx := p.context()
	p.pop()

	if ctx.keys != nil {
		p.push(ctx.keys[n])
	} else {
		p.push(reflect.ValueOf(ctx.fields[n].name))
	}

	return
}

func (p *ValueParser) value() reflect.Value {
	v := p.stack[len(p.stack)-1]

	if !v.IsValid() {
		return v
	}

	switch v.Interface().(type) {
	case error:
		return v
	}

dereference:
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if !v.IsNil() {
			v = v.Elem()
			goto dereference
		}
	}

	return v
}

func (p *ValueParser) push(v reflect.Value) {
	p.stack = append(p.stack, v)
}

func (p *ValueParser) pop() {
	p.stack = p.stack[:len(p.stack)-1]
}

func (p *ValueParser) pushContext(ctx valueParserContext) {
	p.ctx = append(p.ctx, ctx)
}

func (p *ValueParser) popContext() {
	p.ctx = p.ctx[:len(p.ctx)-1]
}

func (p *ValueParser) context() *valueParserContext {
	return &p.ctx[len(p.ctx)-1]
}

// ParseInt parses a decimanl representation of an int64 from b.
//
// The function is equivalent to calling strconv.ParseInt(string(b), 10, 64) but
// it prevents Go from making a memory allocation for converting a byte slice to
// a string (escape analysis fails due to the error returned by strconv.ParseInt).
//
// Because it only works with base 10 the function is also significantly faster
// than strconv.ParseInt.
func ParseInt(b []byte) (int64, error) {
	var val int64

	if len(b) == 0 {
		return 0, errorInvalidUint64(b)
	}

	if b[0] == '-' {
		const max = Int64Min
		const lim = max / 10

		if b = b[1:]; len(b) == 0 {
			return 0, errorInvalidUint64(b)
		}

		for _, d := range b {
			if !(d >= '0' && d <= '9') {
				return 0, errorInvalidInt64(b)
			}

			if val < lim {
				return 0, errorOverflowInt64(b)
			}

			val *= 10
			x := int64(d - '0')

			if val < (max + x) {
				return 0, errorOverflowInt64(b)
			}

			val -= x
		}
	} else {
		const max = Int64Max
		const lim = max / 10

		for _, d := range b {
			if !(d >= '0' && d <= '9') {
				return 0, errorInvalidInt64(b)
			}
			x := int64(d - '0')

			if val > lim {
				return 0, errorOverflowInt64(b)
			}

			if val *= 10; val > (max - x) {
				return 0, errorOverflowInt64(b)
			}

			val += x
		}
	}

	return val, nil
}

// ParseUintHex parses a hexadecimanl representation of a uint64 from b.
//
// The function is equivalent to calling strconv.ParseUint(string(b), 16, 64) but
// it prevents Go from making a memory allocation for converting a byte slice to
// a string (escape analysis fails due to the error returned by strconv.ParseUint).
//
// Because it only works with base 16 the function is also significantly faster
// than strconv.ParseUint.
func ParseUintHex(b []byte) (uint64, error) {
	const max = Uint64Max
	const lim = max / 0x10
	var val uint64

	if len(b) == 0 {
		return 0, errorInvalidUint64(b)
	}

	for _, d := range b {
		var x uint64

		switch {
		case d >= '0' && d <= '9':
			x = uint64(d - '0')

		case d >= 'A' && d <= 'F':
			x = uint64(d-'A') + 0xA

		case d >= 'a' && d <= 'f':
			x = uint64(d-'a') + 0xA

		default:
			return 0, errorInvalidUint64(b)
		}

		if val > lim {
			return 0, errorOverflowUint64(b)
		}

		if val *= 0x10; val > (max - x) {
			return 0, errorOverflowUint64(b)
		}

		val += x
	}

	return val, nil
}

func parseNetAddr(s string) (ip net.IP, port int, zone string, err error) {
	var h string
	var p string

	if h, p, err = net.SplitHostPort(s); err != nil {
		h, p = s, ""
	}

	if len(h) != 0 {
		if off := strings.IndexByte(h, '%'); off >= 0 {
			h, zone = h[:off], h[off+1:]
		}
		if ip = net.ParseIP(h); ip == nil {
			err = errors.New("objconv: bad IP address: " + s)
			return
		}
	}

	if len(p) != 0 {
		if port, err = strconv.Atoi(p); err != nil || port < 0 || port > 65535 {
			err = errors.New("objconv: bad port number: " + s)
			return
		}
	}

	return
}

func errorInvalidInt64(b []byte) error {
	return fmt.Errorf("objconv: %#v is not a valid decimal representation of a signed 64 bits integer", string(b))
}

func errorOverflowInt64(b []byte) error {
	return fmt.Errorf("objconv: %#v overflows the maximum values of a signed 64 bits integer", string(b))
}

func errorInvalidUint64(b []byte) error {
	return fmt.Errorf("objconv: %#v is not a valid decimal representation of an unsigned 64 bits integer", string(b))
}

func errorOverflowUint64(b []byte) error {
	return fmt.Errorf("objconv: %#v overflows the maximum values of an unsigned 64 bits integer", string(b))
}
