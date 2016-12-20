package json

import (
	"io"
	"strconv"
	"time"
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
	b [128]byte
	s []byte

	// W is the writer where the emitter outputs the JSON representation of the
	// values it encodes.
	W io.Writer
}

func NewEmitter(w io.Writer) *Emitter {
	return &Emitter{W: w}
}

func (e *Emitter) EmitNil() (err error) {
	_, err = e.W.Write(nullBytes[:])
	return
}

func (e *Emitter) EmitBool(v bool) (err error) {
	if v {
		_, err = e.W.Write(trueBytes[:])
	} else {
		_, err = e.W.Write(falseBytes[:])
	}
	return
}

func (e *Emitter) EmitInt(v int64) (err error) {
	_, err = e.W.Write(strconv.AppendInt(e.b[:0], v, 10))
	return
}

func (e *Emitter) EmitUint(v uint64) (err error) {
	_, err = e.W.Write(strconv.AppendUint(e.b[:0], v, 10))
	return
}

func (e *Emitter) EmitFloat(v float64) (err error) {
	_, err = e.W.Write(strconv.AppendFloat(e.b[:0], v, 'g', -1, 64))
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

	_, err = e.W.Write(s)
	return
}

func (e *Emitter) EmitBytes(v []byte) error {
	return e.EmitString(string(v))
}

func (e *Emitter) EmitTime(v time.Time) error {
	return e.EmitString(v.Format(time.RFC3339Nano))
}

func (e *Emitter) EmitDuration(v time.Duration) error {
	return e.EmitString(v.String())
}

func (e *Emitter) EmitError(v error) error {
	return e.EmitString(v.Error())
}

func (e *Emitter) EmitArrayBegin(_ int) (err error) {
	_, err = e.W.Write(arrayOpen[:])
	return
}

func (e *Emitter) EmitArrayEnd() (err error) {
	_, err = e.W.Write(arrayClose[:])
	return
}

func (e *Emitter) EmitArrayNext() (err error) {
	_, err = e.W.Write(comma[:])
	return
}

func (e *Emitter) EmitMapBegin(_ int) (err error) {
	_, err = e.W.Write(mapOpen[:])
	return
}

func (e *Emitter) EmitMapEnd() (err error) {
	_, err = e.W.Write(mapClose[:])
	return
}

func (e *Emitter) EmitMapValue() (err error) {
	_, err = e.W.Write(column[:])
	return
}

func (e *Emitter) EmitMapNext() (err error) {
	_, err = e.W.Write(comma[:])
	return
}
