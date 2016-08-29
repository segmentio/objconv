package resp

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/bytesconv"
)

// Parser implements a RESP parser that satisfies the objconv.Parser interface.
type Parser struct{}

// Parse returns the next value read from r.
func (p *Parser) Parse(r *objconv.Reader, hint interface{}) (v interface{}, err error) {
	var b byte

	if b, err = r.ReadByte(); err == nil {
		switch b {
		case ':':
			v, err = p.parseInt(r, hint)

		case '+':
			v, err = p.parseString(r, hint)

		case '-':
			v, err = p.parseError(r, hint)

		case '$':
			v, err = p.parseBulk(r, hint)

		case '*':
			v, err = p.parseArray(r, hint)

		default:
			err = fmt.Errorf("RESP parser found an invalid character '%c'", b)
		}
	}

	return
}

func (p *Parser) parseInt(r *objconv.Reader, hint interface{}) (v interface{}, err error) {
	var line []byte
	var i int64
	var u uint64

	if line, err = readLine(r); err != nil {
		return
	}

	if i, err = bytesconv.ParseInt(line, 10, 64); err != nil {
		// This is an extension to the standard RESP specs so we can support
		// uint, uint64 and uintptr.
		if u, err = bytesconv.ParseUint(line, 10, 64); err != nil {
			return
		}
		v = u
	} else {
		v = i
	}

	if hint != nil {
		if reflect.TypeOf(hint).Kind() == reflect.Bool {
			switch x := v.(type) {
			case int64:
				v = x != 0
			case uint64:
				v = x != 0
			}
		}
	}

	return
}

func (p *Parser) parseString(r *objconv.Reader, hint interface{}) (v interface{}, err error) {
	var line []byte

	if line, err = readLine(r); err != nil {
		return
	}

	if hint != nil {
		switch hint.(type) {
		case string, []byte, []rune:
			// fast path

		case time.Time:
			return time.Parse(time.RFC3339Nano, string(line))

		case time.Duration:
			return time.ParseDuration(string(line))

		default:
			switch reflect.TypeOf(hint).Kind() {
			case reflect.Float32, reflect.Float64:
				return bytesconv.ParseFloat(line, 64)
			}
		}
	}

	v = string(line)
	return
}

func (p *Parser) parseBulk(r *objconv.Reader, hint interface{}) (v interface{}, err error) {
	var b []byte
	var n int64

	if n, err = readInt(r); err != nil {
		return
	}

	if n < 0 {
		return // v = nil
	}

	b = make([]byte, int(n)+2)

	if _, err = r.ReadFull(b); err != nil {
		return
	}

	if !bytes.HasSuffix(b, crlf[:]) {
		err = errBulkMissingCRLF
	}

	b = b[:len(b)-2]

	if hint != nil {
		switch hint.(type) {
		case string, []byte, []rune:
			// fast path

		case time.Time:
			return time.Parse(time.RFC3339Nano, string(b))

		case time.Duration:
			return time.ParseDuration(string(b))

		default:
			switch reflect.TypeOf(hint).Kind() {
			case reflect.Float32, reflect.Float64:
				return bytesconv.ParseFloat(b, 64)
			}
		}
	}

	v = b
	return
}

func (p *Parser) parseError(r *objconv.Reader, hint interface{}) (v interface{}, err error) {
	var line []byte

	if line, err = readLine(r); err == nil {
		v = NewError(string(line))
	}

	return
}

func (p *Parser) parseArray(r *objconv.Reader, hint interface{}) (v interface{}, err error) {
	var n int64

	if n, err = readInt(r); err != nil {
		return
	}

	if n < 0 {
		return // v = nil
	}

	if hint != nil {
		switch reflect.TypeOf(hint).Kind() {
		case reflect.Map, reflect.Struct:
			if (n & 1) != 0 {
				err = errMapRequiresEvenArrayElements
			}
			v = objconv.NewFixedMapParser(p, int(n)/2)
			return
		}
	}

	v = objconv.ArrayParserLen(int(n), objconv.NewArrayParser(p))
	return
}

func readInt(r *objconv.Reader) (v int64, err error) {
	var line []byte

	if line, err = readLine(r); err == nil {
		v, err = bytesconv.ParseInt(line, 10, 64)
	}

	return
}

func readLine(r *objconv.Reader) (v []byte, err error) { return r.ReadLine(objconv.CRLF) }

func init() {
	f := func() objconv.Parser { return &Parser{} }
	objconv.RegisterParser("resp", f)
	objconv.RegisterParser("application/resp", f)
}

var (
	errBulkMissingCRLF              = errors.New("RESP parser expected a CRLF sequence at the end of a bulk string")
	errMapRequiresEvenArrayElements = errors.New("RESP parser requires arrays to have an even number of elements to be decoded as maps or structs")
)
