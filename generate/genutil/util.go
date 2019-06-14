package genutil

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/objutil"
)

// SkipValue properly skips a value in a parser based on the type
func SkipValue(p objconv.Parser) (err error) {
	var typ objconv.Type

	typ, err = p.ParseType()
	if err != nil {
		return err
	}

	switch typ {
	case objconv.Nil:
		return p.ParseNil()
	case objconv.Bool:
		_, err = p.ParseBool()
	case objconv.Int:
		_, err = p.ParseInt()
	case objconv.Uint:
		_, err = p.ParseUint()
	case objconv.Float:
		_, err = p.ParseFloat()
	case objconv.String:
		_, err = p.ParseString()
	case objconv.Bytes:
		_, err = p.ParseBytes()
	case objconv.Time:
		_, err = p.ParseTime()
	case objconv.Duration:
		_, err = p.ParseDuration()
	case objconv.Error:
		_, err = p.ParseError()
	case objconv.Array:
		var count int
		count, err = p.ParseArrayBegin()
		if err != nil {
			return
		}
		var i int
		for i = 0; count < 0 || i < count; i++ {
			if count < 0 || i != 0 {
				err = p.ParseArrayNext(i)
				if err == objconv.End {
					break
				}
				if err != nil {
					return
				}
			}
			if err = SkipValue(p); err != nil {
				return
			}
		}
		err = p.ParseArrayEnd(i)
	case objconv.Map:
		var count int
		count, err = p.ParseMapBegin()
		if err != nil {
			return err
		}
		var i int
		for i = 0; count < 0 || i < count; i++ {
			if count < 0 || i != 0 {
				err = p.ParseMapNext(i)
				if err == objconv.End {
					break
				}
				if err != nil {
					return
				}
			}
			if err = SkipValue(p); err != nil {
				return
			}
			if err = p.ParseMapValue(i); err != nil {
				return
			}
			if err = SkipValue(p); err != nil {
				return
			}
		}
		err = p.ParseMapEnd(i)
	default:
		err = fmt.Errorf("Unknown type: %v", typ)
	}
	return
}

// ParseInt parses a value into an int64
func ParseInt(p objconv.Parser) (int64, error) {
	typ, err := p.ParseType()
	if err != nil {
		return 0, err
	}
	var v int64
	switch typ {
	case objconv.Int:
		v, err = p.ParseInt()
	case objconv.Uint:
		var uv uint64
		uv, err = p.ParseUint()
		if uv > objutil.Int64Max {
			err = fmt.Errorf("objconv: %d does not fit in int64", uv)
		}
		v = int64(uv)
	case objconv.Float:
		var fv float64
		fv, err = p.ParseFloat()
		v = int64(fv)
		if float64(fv) != fv {
			err = fmt.Errorf("objconv: %f does not fit in int64", fv)
		}
	default:
		err = fmt.Errorf("objconv: cannot decode %s into int64", typ.String())
	}
	return v, err
}

// ParseUint parses a value into a uint64
func ParseUint(p objconv.Parser) (uint64, error) {
	typ, err := p.ParseType()
	if err != nil {
		return 0, err
	}
	var v uint64
	switch typ {
	case objconv.Int:
		var iv int64
		iv, err = p.ParseInt()
		if iv < 0 {
			err = fmt.Errorf("objconv: %d does not fit in uint64", iv)
		}
	case objconv.Uint:
		v, err = p.ParseUint()
	case objconv.Float:
		var fv float64
		fv, err = p.ParseFloat()
		v = uint64(fv)
		if float64(v) != fv {
			err = fmt.Errorf("objconv: %f does not fit in uint64", fv)
		}
	default:
		err = fmt.Errorf("objconv: cannot decode %s into uint64", typ.String())
	}
	return v, err
}

// ParseString parses a value into a string
func ParseString(p objconv.Parser) ([]byte, error) {
	typ, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	var b []byte
	switch typ {
	case objconv.String:
		b, err = p.ParseString()
	case objconv.Bytes:
		b, err = p.ParseBytes()
	default:
		err = fmt.Errorf("objconv: cannot decode %s into string", typ.String())
	}
	return b, err
}

// ParseFloat parses a value into a float64
func ParseFloat(p objconv.Parser) (float64, error) {
	typ, err := p.ParseType()
	if err != nil {
		return 0, err
	}
	var v float64
	switch typ {
	case objconv.Int:
		var iv int64
		iv, err = p.ParseInt()
		v = float64(iv)
	case objconv.Uint:
		var uv uint64
		uv, err = p.ParseUint()
		v = float64(uv)
	case objconv.Float:
		v, err = p.ParseFloat()
	default:
		err = fmt.Errorf("objconv: cannot decode %s into float64", typ.String())
	}

	return v, err
}

// ParseBool parses a value into a bool
func ParseBool(p objconv.Parser) (bool, error) {
	typ, err := p.ParseType()
	if err != nil {
		return false, err
	}
	if typ != objconv.Bool {
		return false, fmt.Errorf("objconv: cannot decode %s into bool", typ.String())
	}
	return p.ParseBool()
}

// ParseTime parses a value into a time.Time
func ParseTime(p objconv.Parser) (time.Time, error) {
	typ, err := p.ParseType()
	if err != nil {
		return time.Time{}, err
	}
	if typ == objconv.Time {
		return p.ParseTime()
	}
	if typ != objconv.String {
		return time.Time{}, fmt.Errorf("objconv: cannot decode %s into time.Time", typ.String())
	}
	b, err := p.ParseString()
	if err != nil {
		return time.Time{}, err
	}

	tstr := *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&b[0])),
		Len:  len(b),
	}))

	return time.ParseInLocation(time.RFC3339Nano, tstr, time.UTC)
}
