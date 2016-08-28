package objconv

import (
	"sync"
	"time"
)

// The Emitter interface must be implemented by types that provide encoding
// of a specific format (like json, resp, ...).
type Emitter interface {
	// EmitBegin is called when an encoding operation begins.
	EmitBegin(*Writer)

	// EmitEnd is called when an encoding operation ends.
	EmitEnd(*Writer)

	// EmitNil writes a nil value to the writer.
	EmitNil(*Writer)

	// EmitBool writes a bool value to the writer.
	EmitBool(*Writer, bool)

	// EmitInt writes a int value to the writer.
	EmitInt(*Writer, int)

	// EmitInt8 writes a int8 value to the writer.
	EmitInt8(*Writer, int8)

	// EmitInt16 writes a int16 value to the writer.
	EmitInt16(*Writer, int16)

	// EmitInt32 writes a int32 value to the writer.
	EmitInt32(*Writer, int32)

	// EmitInt64 writes a int64 value to the writer.
	EmitInt64(*Writer, int64)

	// EmitUint writes a uint value to the writer.
	EmitUint(*Writer, uint)

	// EmitUint8 writes a uint8 value to the writer.
	EmitUint8(*Writer, uint8)

	// EmitUint16 writes a uint16 value to the writer.
	EmitUint16(*Writer, uint16)

	// EmitUint32 writes a uint32 value to the writer.
	EmitUint32(*Writer, uint32)

	// EmitUint64 writes a uint64 value to the writer.
	EmitUint64(*Writer, uint64)

	// EmitUintptr writes a uintptr value to the writer.
	EmitUintptr(*Writer, uintptr)

	// EmitFloat32 writes a float32 value to the writer.
	EmitFloat32(*Writer, float32)

	// EmitFloat64 writes a float64 value to the writer.
	EmitFloat64(*Writer, float64)

	// EmitString writes a string value to the writer.
	EmitString(*Writer, string)

	// EmitBytes writes a []byte value to the writer.
	EmitBytes(*Writer, []byte)

	// EmitTime writes a time.Time value to the writer.
	EmitTime(*Writer, time.Time)

	// EmitDuration writes a time.Duration value to the writer.
	EmitDuration(*Writer, time.Duration)

	// EmitError writes an error value to the writer.
	EmitError(*Writer, error)

	// EmitArrayBegin writes the beginning of an array value to the writer.
	// The method receives the length of the array.
	EmitArrayBegin(*Writer, int)

	// EmitArrayEnd writes the end of an array value to the writer.
	EmitArrayEnd(*Writer)

	// EmitArrayNext is called after each array value except to the last one.
	EmitArrayNext(*Writer)

	// EmitMapBegin writes the beginning of a map value to the writer.
	// The method receives the length of the map.
	EmitMapBegin(*Writer, int)

	// EmitMapEnd writes the end of a map value to the writer.
	EmitMapEnd(*Writer)

	// EmitMapValue is called after each map key was written.
	EmitMapValue(*Writer)

	// EmitMapNext is called after each map value was written except the last one.
	EmitMapNext(*Writer)
}

// RegisterEmitter adds a new emitter factory under the given name.
func RegisterEmitter(format string, factory func() Emitter) {
	emitterMutex.Lock()
	emitterStore[format] = factory
	emitterMutex.Unlock()
}

// UnregisterEmitter removes the emitter registered under the given name.
func UnregisterEmitter(format string) {
	emitterMutex.Lock()
	delete(emitterStore, format)
	emitterMutex.Unlock()
}

// GetEmitter returns a new emitter for the given format, or an error if no emitter
// was registered for that format prior to the call.
func GetEmitter(format string) (p Emitter, err error) {
	emitterMutex.RLock()
	if f := emitterStore[format]; f == nil {
		err = &UnsupportedFormatError{format}
	} else {
		p = f()
	}
	emitterMutex.RUnlock()
	return
}

// NewEmitter returns a new emitter for the given format, or panics if not emitter
// was registered for that format prior to the call.
func NewEmitter(format string) Emitter {
	if e, err := GetEmitter(format); err != nil {
		panic(err)
	} else {
		return e
	}
}

var (
	emitterMutex sync.RWMutex
	emitterStore = map[string](func() Emitter){}
)
