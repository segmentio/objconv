package json

import (
	"bytes"
	"fmt"
	"io"

	"github.com/segmentio/objconv"
	"github.com/segmentio/objconv/bytesconv"
)

// Parser implements a JSON parser that satisfies the objconv.Parser interface.
type Parser struct {
	b [64]byte // buffer
}

// Parse returns the next value read from r.
func (p *Parser) Parse(r *objconv.Reader, hint interface{}) (v interface{}, err error) {
	var b byte

	if b, err = peekByte(r); err != nil {
		return
	}

	switch b {
	case 'n':
		return p.parseNull(r)

	case 't':
		return p.parseTrue(r)

	case 'f':
		return p.parseFalse(r)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', '.', 'e', 'E':
		return p.parseNumber(r)

	case '"':
		return p.parseString(r)

	case '[':
		return p.parseArray(r)

	case '{':
		return p.parseMap(r)

	default:
		err = fmt.Errorf("JSON parser found an invalid character '%c'", b)
		return
	}
}

func (p *Parser) parseNull(r *objconv.Reader) (v interface{}, err error) {
	err = p.parseWord(r, nullBytes[:])
	return
}

func (p *Parser) parseTrue(r *objconv.Reader) (v interface{}, err error) {
	if err = p.parseWord(r, trueBytes[:]); err == nil {
		v = true
	}
	return
}

func (p *Parser) parseFalse(r *objconv.Reader) (v interface{}, err error) {
	if err = p.parseWord(r, falseBytes[:]); err == nil {
		v = false
	}
	return
}

func (p *Parser) parseNumber(r *objconv.Reader) (v interface{}, err error) {
	a := p.b[:0]
	f := false
readNumber:
	for {
		var b byte

		if b, err = r.ReadByte(); err != nil {
			if err == io.EOF && len(a) != 0 {
				err = nil
				break readNumber
			}
			return
		}

		switch b {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-':
			a = append(a, b)

		case '.', 'e', 'E':
			a = append(a, b)
			f = true // float

		default:
			if err = r.UnreadByte(); err != nil {
				return
			}
			break readNumber
		}
	}

	if f {
		v, err = bytesconv.ParseFloat(a, 64)

	} else if a[0] == '-' || a[0] == '+' {
		v, err = bytesconv.ParseInt(a, 10, 64)

	} else {
		var u uint64

		if u, err = bytesconv.ParseUint(a, 10, 64); err == nil {
			if u < objconv.Int64Max {
				v = int64(u)
			} else {
				v = u
			}
		}
	}

	return
}

func (p *Parser) parseString(r *objconv.Reader) (v interface{}, err error) {
	if err = p.parseWord(r, quote[:]); err == nil {
		escape := false

		// use the parser's buffer to avoid an extra memory allocation for short strings
		s := p.b[:0]

		for {
			var b byte

			switch b, err = r.ReadByte(); err {
			case nil:
			case io.EOF:
				err = io.ErrUnexpectedEOF
				return
			default:
				return
			}

			if escape {
				escape = false

				switch b {
				case '"', '/', '\\':
					// b = b

				case 'b':
					b = '\b'

				case 'f':
					b = '\f'

				case 'n':
					b = '\n'

				case 'r':
					b = '\r'

				case 't':
					b = '\t'
				}

			} else {
				if b == '"' {
					break
				}

				if b == '\\' {
					escape = true
					continue
				}
			}

			s = append(s, b)
		}

		v = string(s)
	}
	return
}

func (p *Parser) parseArray(r *objconv.Reader) (v interface{}, err error) {
	if err = p.parseWord(r, arrayOpen[:]); err == nil {
		v = &arrayParser{p: p}
	}
	return
}

func (p *Parser) parseMap(r *objconv.Reader) (v interface{}, err error) {
	if err = p.parseWord(r, mapOpen[:]); err == nil {
		v = &mapParser{p: p}
	}
	return
}

func (p *Parser) parseWord(r *objconv.Reader, word []byte) (err error) {
	n := len(word)

	if _, err = r.ReadFull(p.b[:n]); err == nil {
		if !bytes.Equal(p.b[:n], word) {
			err = fmt.Errorf("JSON parser expected '%s' but found '%s'", string(word), string(p.b[:n]))
		}
	}

	return
}

func peekByte(r *objconv.Reader) (b byte, err error) {
	if b, err = readByte(r); err == nil {
		err = r.UnreadByte()
	}
	return
}

func readByte(r *objconv.Reader) (b byte, err error) {
	if err = skipSpaces(r); err == nil {
		b, err = r.ReadByte()
	}
	return
}

func skipSpaces(r *objconv.Reader) (err error) {
	for {
		var b byte

		if b, err = r.ReadByte(); err != nil {
			return
		}

		switch b {
		case ' ', '\b', '\f', '\n', '\r', '\t':
		default:
			err = r.UnreadByte()
			return
		}
	}
}

func init() {
	f := func() objconv.Parser { return &Parser{} }
	objconv.RegisterParser("json", f)
	objconv.RegisterParser("text/json", f)
	objconv.RegisterParser("application/json", f)
}

type arrayParser struct {
	p *Parser
	n int
}

func (a *arrayParser) Len() int { return -1 }

func (a *arrayParser) Parse(r *objconv.Reader, hint interface{}) (v interface{}, err error) {
	var b byte

	switch b, err = readByte(r); err {
	case nil:
	case io.EOF:
		err = io.ErrUnexpectedEOF
		return
	default:
		return
	}

	if b == ']' {
		err = io.EOF
		return
	}

	if a.n == 0 {
		r.UnreadByte()

	} else if b != ',' {
		err = fmt.Errorf("JSON parser expected ',' but found '%c'", b)
		return
	}

	a.n++
	return a.p.Parse(r, hint)
}

type mapParser struct {
	p *Parser
	n int
}

func (a *mapParser) Len() int { return -1 }

func (a *mapParser) ParseKey(r *objconv.Reader, hint interface{}) (k interface{}, err error) {
	var b byte

	switch b, err = readByte(r); err {
	case nil:
	case io.EOF:
		err = io.ErrUnexpectedEOF
		return
	default:
		return
	}

	if b == '}' {
		err = io.EOF
		return
	}

	if a.n == 0 {
		r.UnreadByte()

	} else if b != ',' {
		err = fmt.Errorf("JSON parser expected ',' but found '%c'", b)
		return
	}

	return a.p.Parse(r, hint)
}

func (a *mapParser) ParseValue(r *objconv.Reader, hint interface{}) (v interface{}, err error) {
	var b byte

	switch b, err = readByte(r); err {
	case nil:
	case io.EOF:
		err = io.ErrUnexpectedEOF
		return
	default:
		return
	}

	if b != ':' {
		err = fmt.Errorf("JSON parser expected ':' but found '%c'", b)
		return
	}

	a.n++
	return a.p.Parse(r, hint)
}
