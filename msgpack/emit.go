package msgpack

import (
	"io"
	"time"
	"unsafe"

	"github.com/segmentio/objconv"
)

// Emitter implements a MessagePack emitter that satisfies the objconv.Emitter
// interface.
type Emitter struct {
	w io.Writer
	b [240]byte
}

func NewEmitter(w io.Writer) *Emitter {
	return &Emitter{w: w}
}

func (e *Emitter) Reset(w io.Writer) {
	e.w = w
}

func (e *Emitter) EmitNil() (err error) {
	e.b[0] = Nil
	_, err = e.w.Write(e.b[:1])
	return
}

func (e *Emitter) EmitBool(v bool) (err error) {
	if v {
		e.b[0] = True
	} else {
		e.b[0] = False
	}
	_, err = e.w.Write(e.b[:1])
	return
}

func (e *Emitter) EmitInt(v int64) (err error) {
	n := 0

	if v >= 0 {
		switch {
		case v <= objconv.Int8Max:
			e.b[0] = byte(v) | PositiveFixintTag
			n = 1

		case v <= objconv.Int16Max:
			e.b[0] = Int16
			putUint16(e.b[1:], uint16(v))
			n = 3

		case v <= objconv.Int32Max:
			e.b[0] = Int32
			putUint32(e.b[1:], uint32(v))
			n = 5

		default:
			e.b[0] = Int64
			putUint64(e.b[1:], uint64(v))
			n = 9
		}

	} else {
		switch {
		case v >= -31:
			e.b[0] = byte(v) | NegativeFixintTag
			n = 1

		case v >= objconv.Int8Min:
			e.b[0] = Int8
			e.b[1] = byte(v)
			n = 2

		case v >= objconv.Int16Min:
			e.b[0] = Int16
			putUint16(e.b[1:], uint16(v))
			n = 3

		case v >= objconv.Int32Min:
			e.b[0] = Int32
			putUint32(e.b[1:], uint32(v))
			n = 5

		default:
			e.b[0] = Int64
			putUint64(e.b[1:], uint64(v))
			n = 9
		}
	}

	_, err = e.w.Write(e.b[:n])
	return
}

func (e *Emitter) EmitUint(v uint64) (err error) {
	n := 0

	switch {
	case v <= objconv.Uint8Max:
		e.b[0] = Uint8
		e.b[1] = byte(v)
		n = 2

	case v <= objconv.Uint16Max:
		e.b[0] = Uint16
		putUint16(e.b[1:], uint16(v))
		n = 3

	case v <= objconv.Uint32Max:
		e.b[0] = Uint32
		putUint32(e.b[1:], uint32(v))
		n = 5

	default:
		e.b[0] = Uint64
		putUint64(e.b[1:], v)
		n = 9
	}

	_, err = e.w.Write(e.b[:n])
	return
}

func (e *Emitter) EmitFloat(v float64) (err error) {
	e.b[0] = Float64
	putUint64(e.b[1:], *((*uint64)(unsafe.Pointer(&v))))
	_, err = e.w.Write(e.b[:9])
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
