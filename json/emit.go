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

func (e *Emitter) EmitInt(v int64) (err error) {
	_, err = e.w.Write(strconv.AppendInt(e.s[:0], v, 10))
	return
}

func (e *Emitter) EmitUint(v uint64) (err error) {
	_, err = e.w.Write(strconv.AppendUint(e.s[:0], v, 10))
	return
}

func (e *Emitter) EmitFloat(v float64) (err error) {
	_, err = e.w.Write(strconv.AppendFloat(e.s[:0], v, 'g', -1, 64))
	return
}

func (e *Emitter) EmitString(v string) (err error) {
	i := 0
	j := 0
	n := len(v)
	s := append(e.s[:0], '"')

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

func (e *Emitter) EmitBytes(v []byte) (err error) {
	return e.EmitString(string(v))
}

func (e *Emitter) EmitTime(v time.Time) (err error) {
	s := e.s[:0]

	s = append(s, '"')
	s = v.AppendFormat(s, time.RFC3339Nano)
	s = append(s, '"')

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitDuration(v time.Duration) (err error) {
	s := e.s[:0]

	s = append(s, '"')
	s = objconv.AppendDuration(s, v)
	s = append(s, '"')

	e.s = s[:0]
	_, err = e.w.Write(s)
	return
}

func (e *Emitter) EmitError(v error) (err error) {
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
