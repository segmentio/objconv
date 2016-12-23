package resp

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/objconv"
)

var (
	crlfBytes  = [...]byte{'\r', '\n'}
	nullBytes  = [...]byte{'$', '-', '1', '\r', '\n'}
	trueBytes  = [...]byte{'+', 't', 'r', 'u', 'e', '\r', '\n'}
	falseBytes = [...]byte{'+', 'f', 'a', 'l', 's', 'e', '\r', '\n'}
)

// Emitter implements a RESP emitter that satisfies the objconv.Emitter
// interface.
type Emitter struct {
	w io.Writer
	s []byte
	a [128]byte
}

func NewEmitter(w io.Writer) *Emitter {
	e := &Emitter{w: w}
	e.s = e.a[:0]
	return e
}

func (e *Emitter) Reset(w io.Writer) {
	e.w = w
}

func (e *Emitter) EmitNil() (err error) {
	_, err = e.w.Write(nullBytes[:])
	return
}

func (e *Emitter) EmitBool(v bool) (err error) {
	if v {
		_, err = e.w.Write(trueBytes[:])
	} else {
		_, err = e.w.Write(falseBytes[:])
	}
	return
}

func (e *Emitter) EmitInt(v int64, _ int) (err error) {
	s := e.s[:0]

	s = append(s, ':')
	s = appendInt(s, v)
	s = appendCRLF(s)

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitUint(v uint64, _ int) (err error) {
	if v > objconv.Int64Max {
		return fmt.Errorf("objconv/resp: %d overflows the maximum integer value of %d", v, objconv.Int64Max)
	}

	s := e.s[:0]

	s = append(s, ':')
	s = appendUint(s, v)
	s = appendCRLF(s)

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitFloat(v float64, bitSize int) (err error) {
	s := e.s[:0]

	s = append(s, '+')
	s = appendFloat(s, v, bitSize)
	s = appendCRLF(s)

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitString(v string) (err error) {
	s := e.s[:0]

	if indexCRLF(v) < 0 {
		s = append(s, '+')
		s = append(s, v...)
		s = appendCRLF(s)
	} else {
		s = append(s, '$')
		s = appendUint(s, uint64(len(v)))
		s = appendCRLF(s)
		s = append(s, v...)
		s = appendCRLF(s)
	}

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitBytes(v []byte) (err error) {
	s := e.s[:0]

	s = append(s, '$')
	s = appendUint(s, uint64(len(v)))
	s = appendCRLF(s)

	if (len(v) + 2) <= (cap(s) - len(s)) { // if it fits in the buffer
		s = append(s, v...)
		s = appendCRLF(s)
		e.s = s[:0]

		_, err = e.w.Write(s)
		return
	}

	e.s = s[:0]

	if _, err = e.w.Write(s); err != nil {
		return
	}

	if _, err = e.w.Write(v); err != nil {
		return
	}

	_, err = e.w.Write(crlfBytes[:])
	return
}

func (e *Emitter) EmitTime(v time.Time) (err error) {
	s := e.s[:0]

	s = append(s, '+')
	s = v.AppendFormat(s, time.RFC3339Nano)
	s = appendCRLF(s)

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitDuration(v time.Duration) (err error) {
	s := e.s[:0]

	s = append(s, '+')
	s = objconv.AppendDuration(s, v)
	s = appendCRLF(s)

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitError(v error) (err error) {
	x := v.Error()
	s := e.s[:0]

	if i := indexCRLF(x); i >= 0 {
		x = x[:i] // only keep the first line
	}

	s = append(s, '-')
	s = append(s, x...)
	s = appendCRLF(s)

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitArrayBegin(n int) (err error) {
	s := e.s[:0]

	s = append(s, '*')
	s = appendUint(s, uint64(n))
	s = appendCRLF(s)

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitArrayEnd() (err error) {
	return
}

func (e *Emitter) EmitArrayNext() (err error) {
	return
}

func (e *Emitter) EmitMapBegin(n int) (err error) {
	s := e.s[:0]

	s = append(s, '*')
	s = appendUint(s, 2*uint64(n))
	s = appendCRLF(s)

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitMapEnd() (err error) {
	return
}

func (e *Emitter) EmitMapValue() (err error) {
	return
}

func (e *Emitter) EmitMapNext() (err error) {
	return
}

func appendInt(b []byte, v int64) []byte {
	return strconv.AppendInt(b, v, 10)
}

func appendUint(b []byte, v uint64) []byte {
	return strconv.AppendUint(b, v, 10)
}

func appendFloat(b []byte, v float64, bitSize int) []byte {
	return strconv.AppendFloat(b, v, 'g', -1, bitSize)
}

func appendCRLF(b []byte) []byte {
	return append(b, '\r', '\n')
}

func indexCRLF(s string) int {
	for i, n := 0, len(s); i != n; i++ {
		j := strings.IndexByte(s[i:], '\r')

		if j < 0 {
			break
		}

		if j++; j == n {
			break
		}

		if s[j] == '\n' {
			return j - 1
		}
	}
	return -1
}
