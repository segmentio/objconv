package resp

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/bytesconv"
)

// Parser implements a RESP parser that satisfies the objconv.Parser interface.
type Parser struct{}

// Parse returns the next value read from r.
func (p *Parser) Parse(r *objconv.Reader, hint interface{}) interface{} {
	switch b := p.readByte(r); b {
	case ':':
		return p.parseInt(r, hint)

	case '+':
		return p.parseString(r, hint)

	case '-':
		return p.parseError(r, hint)

	case '$':
		return p.parseBulk(r, hint)

	case '*':
		return p.parseArray(r, hint)

	default:
		panic(fmt.Sprintf("RESP parser found an invalid character '%c'", b))
	}

}

func (p *Parser) parseInt(r *objconv.Reader, hint interface{}) interface{} {
	line := p.readLine(r)

	if v, err := bytesconv.ParseInt(line, 10, 64); err != nil {
		// This is an extension to the standard RESP specs so we can support
		// uint, uint64 and uintptr.
		u, err := bytesconv.ParseUint(line, 10, 64)
		objconv.Check(err)
		if hint != nil && reflect.TypeOf(hint).Kind() == reflect.Bool {
			return v != 0
		}
		return u
	} else {
		if hint != nil && reflect.TypeOf(hint).Kind() == reflect.Bool {
			return v != 0
		}
		return v
	}
}

func (p *Parser) parseString(r *objconv.Reader, hint interface{}) interface{} {
	s := p.readLine(r)

	if hint != nil {
		switch hint.(type) {
		case string, []byte, []rune:
			// fast path

		case time.Time:
			v, err := time.Parse(time.RFC3339Nano, string(s))
			objconv.Check(err)
			return v

		case time.Duration:
			v, err := time.ParseDuration(string(s))
			objconv.Check(err)
			return v

		default:
			switch reflect.TypeOf(hint).Kind() {
			case reflect.Float32, reflect.Float64:
				v, err := bytesconv.ParseFloat(s, 64)
				objconv.Check(err)
				return v
			}
		}
	}

	return string(s)
}

func (p *Parser) parseBulk(r *objconv.Reader, hint interface{}) interface{} {
	n := p.readInt(r)
	if n < 0 {
		return nil
	}

	b := make([]byte, int(n)+2)
	r.ReadFull(b)

	if !bytes.HasSuffix(b, crlf[:]) {
		panic("RESP parser expected a CRLF sequence at the end of a bulk string")
	}

	b = b[:len(b)-2]

	if hint != nil {
		switch hint.(type) {
		case string, []byte, []rune:
			// fast path

		case time.Time:
			v, err := time.Parse(time.RFC3339Nano, string(b))
			objconv.Check(err)
			return v

		case time.Duration:
			v, err := time.ParseDuration(string(b))
			objconv.Check(err)
			return v

		default:
			switch reflect.TypeOf(hint).Kind() {
			case reflect.Float32, reflect.Float64:
				v, err := bytesconv.ParseFloat(b, 64)
				objconv.Check(err)
				return v
			}
		}
	}

	return b
}

func (p *Parser) parseError(r *objconv.Reader, hint interface{}) interface{} {
	return NewError(string(p.readLine(r)))
}

func (p *Parser) parseArray(r *objconv.Reader, hint interface{}) interface{} {
	n := p.readInt(r)
	if n < 0 {
		return nil
	}

	if hint != nil {
		switch reflect.TypeOf(hint).Kind() {
		case reflect.Map, reflect.Struct:
			if (n & 1) != 0 {
				panic("RESP parser requires arrays to have an even number of elements to be decoded as maps or structs")
			}
			return objconv.NewFixedMapParser(p, int(n)/2)
		}
	}

	return objconv.ArrayParserLen(int(n), objconv.NewArrayParser(p))
}

func (p *Parser) readByte(r *objconv.Reader) byte {
	b, err := r.ReadByte()
	objconv.Check(err)
	return b
}

func (p *Parser) readInt(r *objconv.Reader) int64 {
	line := p.readLine(r)
	v, err := bytesconv.ParseInt(line, 10, 64)
	objconv.Check(err)
	return v
}

func (p *Parser) readLine(r *objconv.Reader) []byte {
	line, err := r.ReadLine(objconv.CRLF)
	objconv.Check(err)
	return line
}

func init() {
	f := func() objconv.Parser { return &Parser{} }
	objconv.RegisterParser("resp", f)
	objconv.RegisterParser("application/resp", f)
}
