package cbor

import (
	"io"
	"math"
	"time"

	"github.com/segmentio/objconv"
)

// Emitter implements a MessagePack emitter that satisfies the objconv.Emitter
// interface.
type Emitter struct {
	w io.Writer
	b [240]byte

	// This stack is used to keep track of the array map lengths being parsed.
	// The sback array is the initial backend array for the stack.
	stack []int
	sback [16]int
}

func NewEmitter(w io.Writer) *Emitter {
	e := &Emitter{w: w}
	e.stack = e.sback[:0]
	return e
}

func (e *Emitter) Reset(w io.Writer) {
	e.w = w
	e.stack = e.stack[:0]
}

func (e *Emitter) EmitNil() (err error) {
	e.b[0] = majorByte(MajorType7, Null)
	_, err = e.w.Write(e.b[:1])
	return
}

func (e *Emitter) EmitBool(v bool) (err error) {
	if v {
		e.b[0] = majorByte(MajorType7, True)
	} else {
		e.b[0] = majorByte(MajorType7, False)
	}
	_, err = e.w.Write(e.b[:1])
	return
}

func (e *Emitter) EmitInt(v int64, _ int) (err error) {
	if v >= 0 {
		return e.emitUint(MajorType0, uint64(v))
	}
	return e.emitUint(MajorType1, uint64(-(v + 1)))
}

func (e *Emitter) EmitUint(v uint64, _ int) (err error) {
	return e.emitUint(MajorType0, v)
}

func (e *Emitter) EmitFloat(v float64, bitSize int) (err error) {
	n := 0

	if bitSize == 32 {
		n = 5
		e.b[0] = majorByte(MajorType7, Float32)
		putUint32(e.b[1:], math.Float32bits(float32(v)))
	} else {
		n = 9
		e.b[0] = majorByte(MajorType7, Float64)
		putUint64(e.b[1:], math.Float64bits(v))
	}

	_, err = e.w.Write(e.b[:n])
	return
}

func (e *Emitter) EmitString(v string) (err error) {
	if err = e.emitUint(MajorType3, uint64(len(v))); err != nil {
		return
	}

	for len(v) != 0 {
		n1 := len(v)
		n2 := len(e.b)

		if n1 > n2 {
			n1 = n2
		}

		copy(e.b[:], v[:n1])

		if _, err = e.w.Write(e.b[:n1]); err != nil {
			return
		}

		v = v[n1:]
	}

	return
}

func (e *Emitter) EmitBytes(v []byte) (err error) {
	if err = e.emitUint(MajorType2, uint64(len(v))); err != nil {
		return
	}
	_, err = e.w.Write(v)
	return
}

func (e *Emitter) EmitTime(v time.Time) (err error) {
	e.b[0] = majorByte(MajorType6, TagDateTime)

	if _, err = e.w.Write(e.b[:1]); err != nil {
		return
	}

	var a [64]byte
	var b = v.AppendFormat(a[:0], time.RFC3339Nano)

	if err = e.emitUint(MajorType3, uint64(len(b))); err != nil {
		return
	}

	_, err = e.w.Write(append(e.b[:0], b...))
	return
}

func (e *Emitter) EmitDuration(v time.Duration) (err error) {
	return e.EmitString(string(objconv.AppendDuration(e.b[:0], v)))
}

func (e *Emitter) EmitError(v error) (err error) {
	return e.EmitString(v.Error())
}

func (e *Emitter) EmitArrayBegin(n int) (err error) {
	e.stack = append(e.stack, n)

	if n >= 0 {
		return e.emitUint(MajorType4, uint64(n))
	}

	e.b[0] = majorByte(MajorType4, 31)
	_, err = e.w.Write(e.b[:1])
	return
}

func (e *Emitter) EmitArrayEnd() (err error) {
	i := len(e.stack) - 1
	n := e.stack[i]
	e.stack = e.stack[:i]

	if n < 0 {
		e.b[0] = 0xFF
		_, err = e.w.Write(e.b[:1])
	}
	return
}

func (e *Emitter) EmitArrayNext() (err error) {
	return
}

func (e *Emitter) EmitMapBegin(n int) (err error) {
	e.stack = append(e.stack, n)

	if n >= 0 {
		return e.emitUint(MajorType5, uint64(n))
	}

	e.b[0] = majorByte(MajorType5, 31)
	_, err = e.w.Write(e.b[:1])
	return
}

func (e *Emitter) EmitMapEnd() (err error) {
	i := len(e.stack) - 1
	n := e.stack[i]
	e.stack = e.stack[:i]

	if n < 0 {
		e.b[0] = 0xFF
		_, err = e.w.Write(e.b[:1])
	}
	return
}

func (e *Emitter) EmitMapValue() (err error) {
	return
}

func (e *Emitter) EmitMapNext() (err error) {
	return
}

func (e *Emitter) emitUint(m byte, v uint64) (err error) {
	var n int

	switch {
	case v <= 23:
		n = 1
		e.b[0] = majorByte(m, byte(v))

	case v <= objconv.Uint8Max:
		n = 2
		e.b[0] = majorByte(m, Uint8)
		e.b[1] = uint8(v)

	case v <= objconv.Uint16Max:
		n = 3
		e.b[0] = majorByte(m, Uint16)
		putUint16(e.b[1:], uint16(v))

	case v <= objconv.Uint32Max:
		n = 5
		e.b[0] = majorByte(m, Uint32)
		putUint32(e.b[1:], uint32(v))

	default:
		n = 9
		e.b[0] = majorByte(m, Uint64)
		putUint64(e.b[1:], v)
	}

	_, err = e.w.Write(e.b[:n])
	return
}
