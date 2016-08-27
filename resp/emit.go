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
	// EmitBulkStringsOnly forces the emitter to only output bulk strings.
	EmitBulkStringsOnly bool

	// general purpose buffer
	b [32]byte

	// bulk length buffer
	c [32]byte
}

var (
	nilBulk = [...]byte{'$', '-', '1', '\r', '\n'}

	boolBulkTrue  = [...]byte{'$', '1', '\r', '\n', '1', '\r', '\n'}
	boolBulkFalse = [...]byte{'$', '1', '\r', '\n', '0', '\r', '\n'}

	boolIntTrue  = [...]byte{':', '1', '\r', '\n'}
	boolIntFalse = [...]byte{':', '0', '\r', '\n'}

	crlf = [...]byte{'\r', '\n'}
)

func (e *Emitter) EmitBegin(w *objconv.Writer) {}

func (e *Emitter) EmitEnd(w *objconv.Writer) {}

func (e *Emitter) EmitNil(w *objconv.Writer) { w.Write(nilBulk[:]) }

func (e *Emitter) EmitBool(w *objconv.Writer, v bool) {
	if e.EmitBulkStringsOnly {
		if v {
			w.Write(boolBulkTrue[:])
		} else {
			w.Write(boolBulkFalse[:])
		}
	} else {
		if v {
			w.Write(boolIntTrue[:])
		} else {
			w.Write(boolIntFalse[:])
		}
	}
}

func (e *Emitter) EmitInt(w *objconv.Writer, v int) { e.EmitInt64(w, int64(v)) }

func (e *Emitter) EmitInt8(w *objconv.Writer, v int8) { e.EmitInt64(w, int64(v)) }

func (e *Emitter) EmitInt16(w *objconv.Writer, v int16) { e.EmitInt64(w, int64(v)) }

func (e *Emitter) EmitInt32(w *objconv.Writer, v int32) { e.EmitInt64(w, int64(v)) }

func (e *Emitter) EmitInt64(w *objconv.Writer, v int64) {
	if b := strconv.AppendInt(e.b[:0], v, 10); e.EmitBulkStringsOnly {
		e.emitBulkString(w, string(b))
	} else {
		w.WriteByte(':')
		w.Write(b)
		e.crlf(w)
	}
}

func (e *Emitter) EmitUint(w *objconv.Writer, v uint) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitUint8(w *objconv.Writer, v uint8) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitUint16(w *objconv.Writer, v uint16) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitUint32(w *objconv.Writer, v uint32) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitUint64(w *objconv.Writer, v uint64) {
	if b := strconv.AppendUint(e.b[:0], v, 10); e.EmitBulkStringsOnly {
		e.emitBulkString(w, string(b))
	} else {
		w.WriteByte(':')
		w.Write(b)
		e.crlf(w)
	}
}

func (e *Emitter) EmitUintptr(w *objconv.Writer, v uintptr) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitFloat32(w *objconv.Writer, v float32) { e.emitFloat(w, float64(v), 32) }

func (e *Emitter) EmitFloat64(w *objconv.Writer, v float64) { e.emitFloat(w, v, 64) }

func (e *Emitter) emitFloat(w *objconv.Writer, v float64, p int) {
	e.EmitBytes(w, strconv.AppendFloat(e.b[:0], v, 'g', -1, p))
}

func (e *Emitter) EmitString(w *objconv.Writer, v string) {
	if e.EmitBulkStringsOnly || len(v) > 100 || strings.IndexByte(v, '\r') >= 0 || strings.IndexByte(v, '\n') >= 0 {
		e.emitBulkString(w, v)
	} else {
		e.emitSimpleString(w, v)
	}
}

func (e *Emitter) EmitBytes(w *objconv.Writer, v []byte) {
	e.emitBulkLength(w, len(v))
	w.Write(v)
	e.crlf(w)
}

func (e *Emitter) EmitTime(w *objconv.Writer, v time.Time) {
	e.EmitString(w, v.Format(time.RFC3339Nano))
}

func (e *Emitter) EmitDuration(w *objconv.Writer, v time.Duration) {
	e.EmitString(w, v.String())
}

func (e *Emitter) EmitError(w *objconv.Writer, v error) {
	if s := v.Error(); e.EmitBulkStringsOnly {
		e.emitBulkString(w, s)
	} else {
		w.WriteByte('-')
		w.WriteString(s)
		e.crlf(w)
	}
}

func (e *Emitter) EmitArrayBegin(w *objconv.Writer, n int) {
	w.WriteByte('*')
	w.Write(strconv.AppendInt(e.b[:0], int64(n), 10))
	e.crlf(w)
}

func (e *Emitter) EmitArrayEnd(w *objconv.Writer) {}

func (e *Emitter) EmitArrayNext(w *objconv.Writer) {}

func (e *Emitter) EmitMapBegin(w *objconv.Writer, n int) { e.EmitArrayBegin(w, n+n) }

func (e *Emitter) EmitMapEnd(w *objconv.Writer) {}

func (e *Emitter) EmitMapValue(w *objconv.Writer) {}

func (e *Emitter) EmitMapNext(w *objconv.Writer) {}

func (e *Emitter) emitBulkLength(w *objconv.Writer, n int) {
	w.WriteByte('$')
	w.Write(strconv.AppendInt(e.c[:0], int64(n), 10))
	e.crlf(w)
}

func (e *Emitter) emitBulkString(w *objconv.Writer, v string) {
	e.emitBulkLength(w, len(v))
	w.WriteString(v)
	e.crlf(w)
}

func (e *Emitter) emitSimpleString(w *objconv.Writer, v string) {
	w.WriteByte('+')
	w.WriteString(v)
	e.crlf(w)
}

func (e *Emitter) emitBulkBytes(w *objconv.Writer, v []byte) {
	e.emitBulkLength(w, len(v))
	w.Write(v)
	e.crlf(w)
}

func (e *Emitter) emitSimpleBytes(w *objconv.Writer, v []byte) {
	w.WriteByte('+')
	w.Write(v)
	e.crlf(w)
}

func (e *Emitter) crlf(w *objconv.Writer) { w.Write(crlf[:]) }

func init() {
	f := func() objconv.Emitter { return &Emitter{} }
	objconv.RegisterEmitter("resp", f)
	objconv.RegisterEmitter("application/resp", f)
}
