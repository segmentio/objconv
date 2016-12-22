package json

import (
	"io"
	"strconv"
	"time"

	"github.com/segmentio/objconv"
)

var (
	nullBytes  = [...]byte{'n', 'u', 'l', 'l'}
	trueBytes  = [...]byte{'t', 'r', 'u', 'e'}
	falseBytes = [...]byte{'f', 'a', 'l', 's', 'e'}

	arrayOpen  = [...]byte{'['}
	arrayClose = [...]byte{']'}

	mapOpen  = [...]byte{'{'}
	mapClose = [...]byte{'}'}

	comma  = [...]byte{','}
	column = [...]byte{':'}
)

// Emitter implements a JSON emitter that satisfies the objconv.Emitter
// interface.
type Emitter struct {
	w io.Writer
	s []byte
	b [128]byte
}

func NewEmitter(w io.Writer) *Emitter {
	return &Emitter{w: w}
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

func (e *Emitter) EmitInt(v int64) (err error) {
	_, err = e.w.Write(strconv.AppendInt(e.b[:0], v, 10))
	return
}

func (e *Emitter) EmitUint(v uint64) (err error) {
	_, err = e.w.Write(strconv.AppendUint(e.b[:0], v, 10))
	return
}

func (e *Emitter) EmitFloat(v float64) (err error) {
	_, err = e.w.Write(strconv.AppendFloat(e.b[:0], v, 'g', -1, 64))
	return
}

func (e *Emitter) EmitString(v string) (err error) {
	i := 0
	j := 0
	n := len(v)

	if e.s == nil {
		e.s = e.b[:0]
	}

	s := e.s[:0]
	s = append(s, '"')

	for j != n {
		b := v[j]
		j++

		switch b {
		case '"', '\\', '/':
			// b = b

		case '\b':
			b = 'b'

		case '\f':
			b = 'f'

		case '\n':
			b = 'n'

		case '\r':
			b = 'r'

		case '\t':
			b = 't'

		default:
			continue
		}

		s = append(s, v[i:j-1]...)
		s = append(s, '\\', b)
		i = j
	}

	s = append(s, v[i:j]...)
	s = append(s, '"')
	e.s = s[:0] // in case the buffer was reallocated

	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitBytes(v []byte) error {
	return e.EmitString(string(v))
}

func (e *Emitter) EmitTime(v time.Time) error {
	return e.EmitString(string(v.AppendFormat(e.b[:0], time.RFC3339Nano)))
}

func (e *Emitter) EmitDuration(v time.Duration) error {
	return e.EmitString(string(objconv.AppendDuration(e.b[:0], v)))
}

func (e *Emitter) EmitError(v error) error {
	return e.EmitString(v.Error())
}

func (e *Emitter) EmitArrayBegin(_ int) (err error) {
	_, err = e.w.Write(arrayOpen[:])
	return
}

func (e *Emitter) EmitArrayEnd() (err error) {
	_, err = e.w.Write(arrayClose[:])
	return
}

func (e *Emitter) EmitArrayNext() (err error) {
	_, err = e.w.Write(comma[:])
	return
}

func (e *Emitter) EmitMapBegin(_ int) (err error) {
	_, err = e.w.Write(mapOpen[:])
	return
}

func (e *Emitter) EmitMapEnd() (err error) {
	_, err = e.w.Write(mapClose[:])
	return
}

func (e *Emitter) EmitMapValue() (err error) {
	_, err = e.w.Write(column[:])
	return
}

func (e *Emitter) EmitMapNext() (err error) {
	_, err = e.w.Write(comma[:])
	return
}
