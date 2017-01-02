package cbor

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/segmentio/objconv"
)

type Parser struct {
	r io.Reader // reader to load bytes from
	i int       // offset of the first unread byte in b
	j int       // offset + 1 of the last unread byte in b
	s []byte    // string buffer
	b [240]byte // read buffer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{r: r}
}

func (p *Parser) Reset(r io.Reader) {
	p.r = r
	p.i = 0
	p.j = 0
}

func (p *Parser) Buffered() io.Reader {
	return bytes.NewReader(p.b[p.i:p.j])
}

func (p *Parser) ParseType() (typ objconv.Type, err error) {
	var s []byte

	if s, err = p.peek(1); err != nil {
		return
	}

	switch m, b := majorTypeOf(s[0]); m {
	case MajorType0:
		typ = objconv.Uint

	case MajorType1:
		typ = objconv.Int

	case MajorType2:
		typ = objconv.Bytes

	case MajorType3:
		typ = objconv.String

	case MajorType4:
		typ = objconv.Array

	case MajorType5:
		typ = objconv.Map

	case MajorType6:
		// TODO:
		err = fmt.Errorf("objconv/cbor: tags are not supported yet")

	default:
		switch b {
		case Null, Undefined:
			typ = objconv.Nil

		case False, True:
			typ = objconv.Bool

		case Float16, Float32, Float64:
			typ = objconv.Float

		case Extension:
			typ = objconv.Uint

		default:
			err = fmt.Errorf("objconv/cbor: unexpected value in major type 7: %d", b)
		}
	}

	return
}

func (p *Parser) parseType7() (b byte, err error) {
	var s []byte

	if s, err = p.peek(1); err != nil {
		return
	}

	if _, b = majorTypeOf(s[0]); b != Extension {
		p.i++
	} else {
		if s, err = p.peek(2); err != nil {
			return
		}
		if b = s[1]; b < 32 {
			err = fmt.Errorf("objconv/cbor: invalid extended simple value in major type 7: %d", b)
			return
		}
		p.i += 2
	}

	return
}

func (p *Parser) ParseNil() (err error) {
	_, err = p.parseType7()
	return
}

func (p *Parser) ParseBool() (v bool, err error) {
	var b byte

	if b, err = p.parseType7(); err != nil {
		return
	}

	v = b == True
	return
}

func (p *Parser) ParseInt() (v int64, err error) {
	return
}

func (p *Parser) ParseUint() (v uint64, err error) {
	var s []byte
	var m byte
	var b byte
	var n int

	if s, err = p.peek(1); err != nil {
		return
	}

	if m, b = majorTypeOf(s[0]); m != MajorType0 { // m == MajorType7 && b == Extension
		if b, err = p.parseType7(); err != nil {
			return
		}
		v = uint64(b)
		return
	}

	if b <= 23 {
		v = uint64(b)
		return
	}

	switch b {
	case Uint8:
		n = 2
	case Uint16:
		n = 3
	case Uint32:
		n = 5
	default:
		n = 9
	}

	if s, err = p.peek(n); err != nil {
		return
	}

	switch b {
	case Uint8:
		v = uint64(s[1])
	case Uint16:
		v = uint64(getUint16(s[1:]))
	case Uint32:
		v = uint64(getUint32(s[1:]))
	default:
		v = getUint64(s[1:])
	}

	p.i += n
	return
}

func (p *Parser) ParseFloat() (v float64, err error) {
	var s []byte
	var n int
	var b byte

	if b, err = p.parseType7(); err != nil {
		return
	}

	switch b {
	case Float16:
		n = 2
	case Float32:
		n = 4
	default:
		n = 8
	}

	if s, err = p.peek(n); err != nil {
		return
	}
	p.i += n

	switch b {
	case Float16:
		v = float64(math.Float32frombits(f16tof32bits(getUint16(s))))
	case Float32:
		v = float64(math.Float32frombits(getUint32(s)))
	default:
		v = math.Float64frombits(getUint64(s))
	}

	return
}

func (p *Parser) ParseString() (v []byte, err error) {
	return
}

func (p *Parser) ParseBytes() (v []byte, err error) {
	return
}

func (p *Parser) ParseTime() (v time.Time, err error) {
	return
}

func (p *Parser) ParseDuration() (v time.Duration, err error) {
	return
}

func (p *Parser) ParseError() (v error, err error) {
	return
}

func (p *Parser) ParseArrayBegin() (n int, err error) {
	return
}

func (p *Parser) ParseArrayEnd(n int) (err error) {
	return
}

func (p *Parser) ParseArrayNext(n int) (err error) {
	return
}

func (p *Parser) ParseMapBegin() (n int, err error) {
	return
}

func (p *Parser) ParseMapEnd(n int) (err error) {
	return
}

func (p *Parser) ParseMapValue(n int) (err error) {
	return
}

func (p *Parser) ParseMapNext(n int) (err error) {
	return
}

func (p *Parser) read(n int) (b []byte, err error) {
	if n <= (p.j - p.i) { // check if the string is already buffered
		b = p.b[p.i : p.i+n]
		p.i += n
		return
	}

	if n <= len(p.b) { // check if the string can be loaded in the read buffer
		if b, err = p.peek(n); err != nil {
			return
		}
		p.i += n
		return
	}

	if cap(p.s) < n {
		p.s = make([]byte, n, align(n, 1024))
	} else {
		p.s = p.s[:n]
	}

	copy(p.s, p.b[p.i:p.j])
	n = p.j - p.i
	p.i = 0
	p.j = 0

	if _, err = io.ReadFull(p.r, p.s[n:]); err != nil {
		return
	}

	b = p.s
	return
}

func (p *Parser) peek(n int) (b []byte, err error) {
	for (p.i + n) > p.j {
		if err = p.fill(); err != nil {
			return
		}
	}
	b = p.b[p.i : p.i+n]
	return
}

func (p *Parser) fill() (err error) {
	n := p.j - p.i
	copy(p.b[:], p.b[p.i:p.j])
	p.i = 0
	p.j = n

	if n, err = p.r.Read(p.b[n:]); n > 0 {
		err = nil
		p.j += n
	} else if err != nil {
		return
	} else {
		err = io.ErrNoProgress
		return
	}

	return
}
