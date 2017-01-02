package cbor

import (
	"io"
	"time"
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
	return
}

func (e *Emitter) EmitBool(v bool) (err error) {
	return
}

func (e *Emitter) EmitInt(v int64, _ int) (err error) {
	return
}

func (e *Emitter) EmitUint(v uint64, _ int) (err error) {
	return
}

func (e *Emitter) EmitFloat(v float64, bitSize int) (err error) {
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
