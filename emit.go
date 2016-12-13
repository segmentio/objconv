package objconv

import "time"

// The Emitter interface must be implemented by types that provide encoding
// of a specific format (like json, resp, ...).
type Emitter interface {
	// EmitNil writes a nil value to the writer.
	EmitNil() error

	// EmitBool writes a bool value to the writer.
	EmitBool(bool) error

	// EmitInt writes a int value to the writer.
	EmitInt(int) error

	// EmitInt8 writes a int8 value to the writer.
	EmitInt8(int8) error

	// EmitInt16 writes a int16 value to the writer.
	EmitInt16(int16) error

	// EmitInt32 writes a int32 value to the writer.
	EmitInt32(int32) error

	// EmitInt64 writes a int64 value to the writer.
	EmitInt64(int64) error

	// EmitUint writes a uint value to the writer.
	EmitUint(uint) error

	// EmitUint8 writes a uint8 value to the writer.
	EmitUint8(uint8) error

	// EmitUint16 writes a uint16 value to the writer.
	EmitUint16(uint16) error

	// EmitUint32 writes a uint32 value to the writer.
	EmitUint32(uint32) error

	// EmitUint64 writes a uint64 value to the writer.
	EmitUint64(uint64) error

	// EmitUintptr writes a uintptr value to the writer.
	EmitUintptr(uintptr) error

	// EmitFloat32 writes a float32 value to the writer.
	EmitFloat32(float32) error

	// EmitFloat64 writes a float64 value to the writer.
	EmitFloat64(float64) error

	// EmitString writes a string value to the writer.
	EmitString(string) error

	// EmitBytes writes a []byte value to the writer.
	EmitBytes([]byte) error

	// EmitTime writes a time.Time value to the writer.
	EmitTime(time.Time) error

	// EmitDuration writes a time.Duration value to the writer.
	EmitDuration(time.Duration) error

	// EmitError writes an error value to the writer.
	EmitError(error) error

	// EmitArrayBegin writes the beginning of an array value to the writer.
	// The method receives the length of the array.
	EmitArrayBegin(int) error

	// EmitArrayEnd writes the end of an array value to the writer.
	EmitArrayEnd() error

	// EmitArrayNext is called after each array value except to the last one.
	EmitArrayNext() error

	// EmitMapBegin writes the beginning of a map value to the writer.
	// The method receives the length of the map.
	EmitMapBegin(int) error

	// EmitMapEnd writes the end of a map value to the writer.
	EmitMapEnd() error

	// EmitMapValue is called after each map key was written.
	EmitMapValue() error

	// EmitMapNext is called after each map value was written except the last one.
	EmitMapNext() error
}

// ValueEmitter is a special kind of emitter, instead of serializing the values
// it receives it builds an in-memory representation of the data.
//
// This is useful for testing the high-level API of the package without actually
// having to generate a serialized representation.
type ValueEmitter struct {
	stack []interface{}
	marks []int
}

// Value returns the value built in the emitter.
func (e *ValueEmitter) Value() interface{} { return e.stack[0] }

func (e *ValueEmitter) EmitNil() error { return e.push(nil) }

func (e *ValueEmitter) EmitBool(v bool) error { return e.push(v) }

func (e *ValueEmitter) EmitInt(v int) error { return e.push(v) }

func (e *ValueEmitter) EmitInt8(v int8) error { return e.push(v) }

func (e *ValueEmitter) EmitInt16(v int16) error { return e.push(v) }

func (e *ValueEmitter) EmitInt32(v int32) error { return e.push(v) }

func (e *ValueEmitter) EmitInt64(v int64) error { return e.push(v) }

func (e *ValueEmitter) EmitUint(v uint) error { return e.push(v) }

func (e *ValueEmitter) EmitUint8(v uint8) error { return e.push(v) }

func (e *ValueEmitter) EmitUint16(v uint16) error { return e.push(v) }

func (e *ValueEmitter) EmitUint32(v uint32) error { return e.push(v) }

func (e *ValueEmitter) EmitUint64(v uint64) error { return e.push(v) }

func (e *ValueEmitter) EmitUintptr(v uintptr) error { return e.push(v) }

func (e *ValueEmitter) EmitFloat32(v float32) error { return e.push(v) }

func (e *ValueEmitter) EmitFloat64(v float64) error { return e.push(v) }

func (e *ValueEmitter) EmitString(v string) error { return e.push(v) }

func (e *ValueEmitter) EmitBytes(v []byte) error { return e.push(v) }

func (e *ValueEmitter) EmitTime(v time.Time) error { return e.push(v) }

func (e *ValueEmitter) EmitDuration(v time.Duration) error { return e.push(v) }

func (e *ValueEmitter) EmitError(v error) error { return e.push(v) }

func (e *ValueEmitter) EmitArrayBegin(v int) error { return e.pushMark() }

func (e *ValueEmitter) EmitArrayEnd() error {
	v := e.pop(e.popMark())
	a := make([]interface{}, len(v))
	copy(a, v)
	return e.push(a)
}

func (e *ValueEmitter) EmitArrayNext() error { return nil }

func (e *ValueEmitter) EmitMapBegin(v int) error { return e.pushMark() }

func (e *ValueEmitter) EmitMapEnd() error {
	v := e.pop(e.popMark())
	n := len(v)
	m := make(map[interface{}]interface{}, n/2)

	for i := 0; i != n; i += 2 {
		m[v[i]] = v[i+1]
	}

	return e.push(m)
}

func (e *ValueEmitter) EmitMapValue() error { return nil }

func (e *ValueEmitter) EmitMapNext() error { return nil }

func (e *ValueEmitter) push(v interface{}) error {
	e.stack = append(e.stack, v)
	return nil
}

func (e *ValueEmitter) pop(n int) []interface{} {
	v := e.stack[n:]
	e.stack = e.stack[:n]
	return v
}

func (e *ValueEmitter) pushMark() error {
	e.marks = append(e.marks, len(e.stack))
	return nil
}

func (e *ValueEmitter) popMark() int {
	n := len(e.marks) - 1
	m := e.marks[n]
	e.marks = e.marks[:n]
	return m
}

// DiscardEmitter is a special emitter that outputs nothing and simply discards
// the values.
//
// This emitter is mostly useful to benchmark the encoder, but it can also be
// used to disable an encoder output if necessary.
type DiscardEmitter struct{}

func (e DiscardEmitter) EmitNil() error { return nil }

func (e DiscardEmitter) EmitBool(v bool) error { return nil }

func (e DiscardEmitter) EmitInt(v int) error { return nil }

func (e DiscardEmitter) EmitInt8(v int8) error { return nil }

func (e DiscardEmitter) EmitInt16(v int16) error { return nil }

func (e DiscardEmitter) EmitInt32(v int32) error { return nil }

func (e DiscardEmitter) EmitInt64(v int64) error { return nil }

func (e DiscardEmitter) EmitUint(v uint) error { return nil }

func (e DiscardEmitter) EmitUint8(v uint8) error { return nil }

func (e DiscardEmitter) EmitUint16(v uint16) error { return nil }

func (e DiscardEmitter) EmitUint32(v uint32) error { return nil }

func (e DiscardEmitter) EmitUint64(v uint64) error { return nil }

func (e DiscardEmitter) EmitUintptr(v uintptr) error { return nil }

func (e DiscardEmitter) EmitFloat32(v float32) error { return nil }

func (e DiscardEmitter) EmitFloat64(v float64) error { return nil }

func (e DiscardEmitter) EmitString(v string) error { return nil }

func (e DiscardEmitter) EmitBytes(v []byte) error { return nil }

func (e DiscardEmitter) EmitTime(v time.Time) error { return nil }

func (e DiscardEmitter) EmitDuration(v time.Duration) error { return nil }

func (e DiscardEmitter) EmitError(v error) error { return nil }

func (e DiscardEmitter) EmitArrayBegin(v int) error { return nil }

func (e DiscardEmitter) EmitArrayEnd() error { return nil }

func (e DiscardEmitter) EmitArrayNext() error { return nil }

func (e DiscardEmitter) EmitMapBegin(v int) error { return nil }

func (e DiscardEmitter) EmitMapEnd() error { return nil }

func (e DiscardEmitter) EmitMapDiscard() error { return nil }

func (e DiscardEmitter) EmitMapNext() error { return nil }

func (e DiscardEmitter) EmitMapValue() error { return nil }
