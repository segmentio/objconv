package json

import (
	"strconv"
	"time"

	"github.com/segmentio/objconv"
)

// Emitter implements a JSON emitter that satisfies the objconv.Emitter
// interface.
type Emitter struct {
	b [32]byte // buffer
}

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
	quote  = [...]byte{'"'}
)

func (e *Emitter) EmitBegin(w *objconv.Writer) {}

func (e *Emitter) EmitEnd(w *objconv.Writer) {}

func (e *Emitter) EmitNil(w *objconv.Writer) { w.Write(nullBytes[:]) }

func (e *Emitter) EmitBool(w *objconv.Writer, v bool) {
	if v {
		w.Write(trueBytes[:])
	} else {
		w.Write(falseBytes[:])
	}
}

func (e *Emitter) EmitInt(w *objconv.Writer, v int) { e.EmitInt64(w, int64(v)) }

func (e *Emitter) EmitInt8(w *objconv.Writer, v int8) { e.EmitInt64(w, int64(v)) }

func (e *Emitter) EmitInt16(w *objconv.Writer, v int16) { e.EmitInt64(w, int64(v)) }

func (e *Emitter) EmitInt32(w *objconv.Writer, v int32) { e.EmitInt64(w, int64(v)) }

func (e *Emitter) EmitInt64(w *objconv.Writer, v int64) {
	w.Write(strconv.AppendInt(e.b[:0], v, 10))
}

func (e *Emitter) EmitUint(w *objconv.Writer, v uint) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitUint8(w *objconv.Writer, v uint8) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitUint16(w *objconv.Writer, v uint16) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitUint32(w *objconv.Writer, v uint32) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitUint64(w *objconv.Writer, v uint64) {
	w.Write(strconv.AppendUint(e.b[:0], v, 10))
}

func (e *Emitter) EmitUintptr(w *objconv.Writer, v uintptr) { e.EmitUint64(w, uint64(v)) }

func (e *Emitter) EmitFloat32(w *objconv.Writer, v float32) { e.emitFloat(w, float64(v), 32) }

func (e *Emitter) EmitFloat64(w *objconv.Writer, v float64) { e.emitFloat(w, v, 64) }

func (e *Emitter) emitFloat(w *objconv.Writer, v float64, p int) {
	w.Write(strconv.AppendFloat(e.b[:0], v, 'g', -1, p))
}

func (e *Emitter) EmitString(w *objconv.Writer, v string) {
	i := 0
	j := 0
	n := len(v)
	w.Write(quote[:])

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

		e.b[0] = '\\'
		e.b[1] = b
		w.WriteString(v[i : j-1])
		w.Write(e.b[:2])
		i = j
	}

	w.WriteString(v[i:j])
	w.Write(quote[:])
}

func (e *Emitter) EmitBytes(w *objconv.Writer, v []byte) {
	i := 0
	j := 0
	n := len(v)
	w.Write(quote[:])

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

		e.b[0] = '\\'
		e.b[1] = b
		w.Write(v[i : j-1])
		w.Write(e.b[:2])
		i = j
	}

	w.Write(v[i:j])
	w.Write(quote[:])
}

func (e *Emitter) EmitTime(w *objconv.Writer, v time.Time) {
	e.EmitString(w, v.Format(time.RFC3339Nano))
}

func (e *Emitter) EmitDuration(w *objconv.Writer, v time.Duration) {
	e.EmitString(w, v.String())
}

func (e *Emitter) EmitError(w *objconv.Writer, v error) {
	e.EmitString(w, v.Error())
}

func (e *Emitter) EmitArrayBegin(w *objconv.Writer, _ int) { w.Write(arrayOpen[:]) }

func (e *Emitter) EmitArrayEnd(w *objconv.Writer) { w.Write(arrayClose[:]) }

func (e *Emitter) EmitArrayNext(w *objconv.Writer) { w.Write(comma[:]) }

func (e *Emitter) EmitMapBegin(w *objconv.Writer, _ int) { w.Write(mapOpen[:]) }

func (e *Emitter) EmitMapEnd(w *objconv.Writer) { w.Write(mapClose[:]) }

func (e *Emitter) EmitMapValue(w *objconv.Writer) { w.Write(column[:]) }

func (e *Emitter) EmitMapNext(w *objconv.Writer) { w.Write(comma[:]) }

func init() {
	f := func() objconv.Emitter { return &Emitter{} }
	objconv.RegisterEmitter("json", f)
	objconv.RegisterEmitter("text/json", f)
	objconv.RegisterEmitter("application/json", f)
}
