package resp

import (
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/objconv"
)

// Emitter implements a RESP emitter that satisfies the objconv.Emitter
// interface.
type Emitter struct {
	// EmitBulkOnly controls whether the emitter is allowed to emit simple
	// strings or should only emit bulk strings.
	EmitBulkStringsOnly bool
}

func (f *Emitter) EmitBegin(w *objconv.Writer) {}

func (f *Emitter) EmitEnd(w *objconv.Writer) {}

func (f *Emitter) EmitNil(w *objconv.Writer) {
	w.WriteString("$-1\r\n")
}

func (f *Emitter) EmitBool(w *objconv.Writer, v bool) {
	if v {
		f.EmitInt(w, 1)
	} else {
		f.EmitInt(w, 0)
	}
}

func (f *Emitter) EmitInt(w *objconv.Writer, v int) { f.EmitInt64(w, int64(v)) }

func (f *Emitter) EmitInt8(w *objconv.Writer, v int8) { f.EmitInt64(w, int64(v)) }

func (f *Emitter) EmitInt16(w *objconv.Writer, v int16) { f.EmitInt64(w, int64(v)) }

func (f *Emitter) EmitInt32(w *objconv.Writer, v int32) { f.EmitInt64(w, int64(v)) }

func (f *Emitter) EmitInt64(w *objconv.Writer, v int64) {
	var a [64]byte
	w.WriteByte(':')
	w.Write(strconv.AppendInt(a[:0], v, 10))
	f.crlf(w)
}

func (f *Emitter) EmitUint(w *objconv.Writer, v uint) { f.EmitUint64(w, uint64(v)) }

func (f *Emitter) EmitUint8(w *objconv.Writer, v uint8) { f.EmitUint64(w, uint64(v)) }

func (f *Emitter) EmitUint16(w *objconv.Writer, v uint16) { f.EmitUint64(w, uint64(v)) }

func (f *Emitter) EmitUint32(w *objconv.Writer, v uint32) { f.EmitUint64(w, uint64(v)) }

func (f *Emitter) EmitUint64(w *objconv.Writer, v uint64) {
	var a [64]byte
	w.WriteByte(':')
	w.Write(strconv.AppendUint(a[:0], v, 10))
	f.crlf(w)
}

func (f *Emitter) EmitUintptr(w *objconv.Writer, v uintptr) { f.EmitUint64(w, uint64(v)) }

func (f *Emitter) EmitFloat32(w *objconv.Writer, v float32) { f.formatFloat(w, float64(v), 32) }

func (f *Emitter) EmitFloat64(w *objconv.Writer, v float64) { f.formatFloat(w, v, 64) }

func (f *Emitter) formatFloat(w *objconv.Writer, v float64, p int) {
	var a [64]byte
	w.WriteByte('+')
	w.Write(strconv.AppendFloat(a[:0], v, 'g', -1, p))
	f.crlf(w)
}

func (f *Emitter) EmitString(w *objconv.Writer, v string) {
	if f.EmitBulkStringsOnly || len(v) > 100 || strings.IndexByte(v, '\r') >= 0 || strings.IndexByte(v, '\n') >= 0 {
		f.formatBulkString(w, v)
	} else {
		f.formatSimpleString(w, v)
	}
}

func (f *Emitter) EmitBytes(w *objconv.Writer, v []byte) {
	f.formatBulkLength(w, len(v))
	w.Write(v)
	f.crlf(w)
}

func (f *Emitter) EmitTime(w *objconv.Writer, v time.Time) {
	f.EmitString(w, v.Format(time.RFC3339Nano))
}

func (f *Emitter) EmitDuration(w *objconv.Writer, v time.Duration) {
	f.EmitString(w, v.String())
}

func (f *Emitter) EmitError(w *objconv.Writer, v error) {
	w.WriteByte('-')
	w.WriteString(v.Error())
	f.crlf(w)
}

func (f *Emitter) EmitArrayBegin(w *objconv.Writer, n int) {
	var a [64]byte
	w.WriteByte('*')
	w.Write(strconv.AppendInt(a[:0], int64(n), 10))
	f.crlf(w)
}

func (f *Emitter) EmitArrayEnd(w *objconv.Writer) {}

func (f *Emitter) EmitArrayNext(w *objconv.Writer) {}

func (f *Emitter) EmitMapBegin(w *objconv.Writer, n int) { f.EmitArrayBegin(w, n+n) }

func (f *Emitter) EmitMapEnd(w *objconv.Writer) {}

func (f *Emitter) EmitMapValue(w *objconv.Writer) {}

func (f *Emitter) EmitMapNext(w *objconv.Writer) {}

func (f *Emitter) formatBulkLength(w *objconv.Writer, n int) {
	var a [64]byte
	w.WriteByte('$')
	w.Write(strconv.AppendInt(a[:0], int64(n), 10))
	f.crlf(w)
}

func (f *Emitter) formatBulkString(w *objconv.Writer, v string) {
	f.formatBulkLength(w, len(v))
	w.WriteString(v)
	f.crlf(w)
}

func (f *Emitter) formatSimpleString(w *objconv.Writer, v string) {
	w.WriteByte('+')
	w.WriteString(v)
	f.crlf(w)
}

func (f *Emitter) crlf(w *objconv.Writer) { w.WriteString(string(objconv.CRLF)) }

func init() {
	f := func() objconv.Emitter { return &Emitter{} }
	objconv.RegisterEmitter("resp", f)
	objconv.RegisterEmitter("application/resp", f)
}
