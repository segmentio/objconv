package resp

import (
	"bytes"
	"io"
	"reflect"
	"time"
	"unsafe"

	"github.com/segmentio/objconv"
)

type Parser struct {
	r io.Reader // reader to load bytes from
	s []byte    // buffer used for building strings
	i int       // offset of the first byte in b
	j int       // offset of the last byte in b
	b [128]byte // buffer where bytes are loaded from the reader
	c [128]byte // initial backend array for s
}

func NewParser(r io.Reader) *Parser {
	p := &Parser{r: r}
	p.s = p.c[:0]
	return p
}

func (p *Parser) Reset(r io.Reader) {
	p.r = r
	p.i = 0
	p.j = 0
}

func (p *Parser) Buffered() io.Reader {
	return bytes.NewReader(p.b[p.i:p.j])
}

func (p *Parser) ParseType() (t objconv.Type, err error) {

	return
}

func (p *Parser) ParseNil() (err error) {

	return
}

func (p *Parser) ParseBool() (v bool, err error) {

	return
}

func (p *Parser) ParseInt() (v int64, err error) {

	return
}

func (p *Parser) ParseUint() (v uint64, err error) {

	return
}

func (p *Parser) ParseFloat() (v float64, err error) {

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

func stringNoCopy(b []byte) string {
	n := len(b)
	if n == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&b[0])),
		Len:  n,
	}))
}
