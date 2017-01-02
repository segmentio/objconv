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

	// This stack is used to cache lengths for variable size arrays and maps.
	stack []int

	// sback is used as the initial backing array for the stack slice to avoid
	// dynamic memory allocations for the most common use cases.
	sback [8]int
}

func NewEmitter(w io.Writer) *Emitter {
	e := &Emitter{w: w}
	e.stack = e.sback[:0]
	return e
}

func (e *Emitter) Reset(w io.Writer) {
	e.w = w
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
		return e.EmitUint(uint64(v), 0)
	}

	// TODO
	return
}

func (e *Emitter) EmitUint(v uint64, _ int) (err error) {
	var n int

	switch {
	case v <= 23:
		n = 1
		e.b[0] = majorByte(MajorType0, byte(v))

	case v <= objconv.Uint8Max:
		n = 2
		e.b[0] = majorByte(MajorType0, Uint8)
		e.b[1] = uint8(v)

	case v <= objconv.Uint16Max:
		n = 3
		e.b[0] = majorByte(MajorType0, Uint16)
		putUint16(e.b[1:], uint16(v))

	case v <= objconv.Uint32Max:
		n = 5
		e.b[0] = majorByte(MajorType0, Uint32)
		putUint32(e.b[1:], uint32(v))

	default:
		n = 9
		e.b[0] = majorByte(MajorType0, Uint64)
		putUint64(e.b[1:], v)
	}

	_, err = e.w.Write(e.b[:n])
	return
}

func (e *Emitter) EmitFloat(v float64, bitSize int) (err error) {
	if bitSize == 32 {
		e.b[0] = majorByte(MajorType7, Float32)
		putUint32(e.b[1:], math.Float32bits(float32(v)))
	} else {
		e.b[0] = majorByte(MajorType7, Float64)
		putUint64(e.b[1:], math.Float64bits(v))
	}
	_, err = e.w.Write(e.b[:1+(bitSize/8)])
	return
}

func (e *Emitter) EmitString(v string) (err error) {
	return
}

func (e *Emitter) EmitBytes(v []byte) (err error) {
	return
}

func (e *Emitter) EmitTime(v time.Time) (err error) {
	return
}

func (e *Emitter) EmitDuration(v time.Duration) (err error) {
	return
}

func (e *Emitter) EmitError(v error) (err error) {
	return
}

func (e *Emitter) EmitArrayBegin(n int) (err error) {
	return
}

func (e *Emitter) EmitArrayEnd() (err error) {
	return
}

func (e *Emitter) EmitArrayNext() (err error) {
	return
}

func (e *Emitter) EmitMapBegin(n int) (err error) {
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
