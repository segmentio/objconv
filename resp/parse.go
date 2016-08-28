package resp

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/segmentio/objconv"
)

// Parser implements a RESP parser that satisfies the objconv.Parser interface.
type Parser struct{}

// Parse returns the next value read from r.
func (p *Parser) Parse(r *objconv.Reader, hint interface{}) interface{} {
	switch b, _ := r.ReadByte(); b {
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
	line := string(p.readLine(r))

	if v, err := strconv.ParseInt(line, 10, 64); err != nil {
		// This is an extension to the standard RESP specs so we can support
		// uint, uint64 and uintptr.
		u, err := strconv.ParseUint(line, 10, 64)
		objconv.Assertf(err == nil, "RESP parser expected an integer but found '%#v'", line)
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
	s := string(p.readLine(r))

	if hint != nil {
		switch hint.(type) {
		case string, []byte, []rune:
			// fast path

		case time.Time:
			v, err := time.Parse(time.RFC3339Nano, s)
			objconv.AssertErr(err)
			return v

		case time.Duration:
			v, err := time.ParseDuration(s)
			objconv.AssertErr(err)
			return v

		default:
			switch reflect.TypeOf(hint).Kind() {
			case reflect.Float32, reflect.Float64:
				v, err := strconv.ParseFloat(s, 64)
				objconv.AssertErr(err)
				return v
			}
		}
	}

	return s
}

func (p *Parser) parseBulk(r *objconv.Reader, hint interface{}) interface{} {
	n := p.readInt(r)
	if n < 0 {
		return nil
	}

	b := make([]byte, int(n)+2)
	r.ReadFull(b)
	assertHasCRLF(b)
	b = b[:len(b)-2]

	if hint != nil {
		switch hint.(type) {
		case string, []byte, []rune:
			// fast path

		case time.Time:
			v, err := time.Parse(time.RFC3339Nano, string(b))
			objconv.AssertErr(err)
			return v

		case time.Duration:
			v, err := time.ParseDuration(string(b))
			objconv.AssertErr(err)
			return v

		default:
			switch reflect.TypeOf(hint).Kind() {
			case reflect.Float32, reflect.Float64:
				v, err := strconv.ParseFloat(string(b), 64)
				objconv.AssertErr(err)
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
			objconv.Assert((n&1) == 0, "RESP parser requires arrays to have an even number of elements to be decoded as maps or structs")
			return objconv.NewFixedMapParser(p, int(n)/2)
		}
	}

	return objconv.ArrayParserLen(int(n), objconv.NewArrayParser(p))
}

func (p *Parser) readInt(r *objconv.Reader) int64 {
	line := p.readLine(r)
	v, _ := strconv.ParseInt(string(line), 10, 64)
	objconv.Assertf(err == nil, "RESP parser expected an integer but found '%c'", line[0])
	return v
}

func (p *Parser) readLine(r *objconv.Reader) []byte { return r.ReadLine(objconv.CRLF) }

func init() {
	f := func() objconv.Parser { return &Parser{} }
	objconv.RegisterParser("resp", f)
	objconv.RegisterParser("application/resp", f)
}

func assertHasCRLF(b []byte) {
	objconv.Assert(
		bytes.HasSuffix(b, []byte(objconv.CRLF)),
		"RESP parser expected a CRLF sequence at the end of a bulk string",
	)
}
