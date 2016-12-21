package objconv

import (
	"encoding"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// An Encoder implements the high-level encoding algorithm that inspect encoded
// values and drive the use of an Emitter to create a serialized representation
// of the data.
//
// Instances of Encoder are not safe for use by multiple goroutines.
type Encoder struct {
	Emitter     Emitter // the emitter used by this encoder
	SortMapKeys bool    // whether the map keys should be sorted
	key         bool
}

// NewEncoder returns a new encoder that outputs values to emitter.
//
// Encoders created by this function use the default encoder configuration,
// which is equivalent to using a zero-value EncoderConfig with only the Emitter
// field set.
//
// The function panics if emitter is nil.
func NewEncoder(emitter Emitter) *Encoder {
	if emitter == nil {
		panic("objconv: the emitter is nil")
	}
	return &Encoder{
		Emitter: emitter,
	}
}

// Encode encodes the generic value v.
func (e Encoder) Encode(v interface{}) (err error) {
	if e.key {
		if e.key, err = false, e.encodeMapValue(); err != nil {
			return
		}
	}

	if v == nil {
		return e.encodeNil()
	}

	return e.encodeValue(reflect.ValueOf(v))
}

func (e Encoder) encodeValue(v reflect.Value) error {
	return encodeFuncOf(v.Type())(e, v)
}

func (e Encoder) encodeValueNil(v reflect.Value) error {
	return e.encodeNil()
}

func (e Encoder) encodeValueBool(v reflect.Value) error {
	return e.encodeBool(v.Bool())
}

func (e Encoder) encodeValueInt(v reflect.Value) error {
	return e.encodeInt(v.Int())
}

func (e Encoder) encodeValueUint(v reflect.Value) error {
	return e.encodeUint(v.Uint())
}

func (e Encoder) encodeValueFloat(v reflect.Value) error {
	return e.encodeFloat(v.Float())
}

func (e Encoder) encodeValueString(v reflect.Value) error {
	return e.encodeString(v.String())
}

func (e Encoder) encodeValueBytes(v reflect.Value) error {
	return e.encodeBytes(v.Bytes())
}

func (e Encoder) encodeValueTime(v reflect.Value) error {
	return e.encodeTime(v.Interface().(time.Time))
}

func (e Encoder) encodeValueDuration(v reflect.Value) error {
	return e.encodeDuration(v.Interface().(time.Duration))
}

func (e Encoder) encodeValueError(v reflect.Value) error {
	return e.encodeError(v.Interface().(error))
}

func (e Encoder) encodeValueArray(v reflect.Value) error {
	return e.encodeValueArrayWith(v, encodeFuncOf(v.Type().Elem()))
}

func (e Encoder) encodeValueArrayWith(v reflect.Value, f encodeFunc) error {
	i := 0
	return e.EncodeArray(v.Len(), func(e Encoder) (err error) {
		err = f(e, v.Index(i))
		i++
		return
	})
}

func (e Encoder) encodeValueMap(v reflect.Value) error {
	t := v.Type()
	kf := encodeFuncOf(t.Key())
	vf := encodeFuncOf(t.Elem())
	return e.encodeValueMapWith(v, kf, vf)
}

func (e Encoder) encodeValueMapWith(v reflect.Value, kf encodeFunc, vf encodeFunc) error {
	keys := v.MapKeys()
	if e.SortMapKeys {
		sortValues(v.Type().Key(), keys)
	}
	i := 0
	return e.EncodeMap(v.Len(), func(ke Encoder, ve Encoder) (err error) {
		k := keys[i]
		v := v.MapIndex(k)
		if err = kf(e, k); err != nil {
			return
		}
		if err = e.encodeMapValue(); err != nil {
			return
		}
		if err = vf(e, v); err != nil {
			return
		}
		i++
		return
	})
}

func (e Encoder) encodeValueStruct(v reflect.Value) error {
	return e.encodeValueStructWith(v, LookupStruct(v.Type()))
}

func (e Encoder) encodeValueStructWith(v reflect.Value, s *Struct) (err error) {
	n := 0

	for i := range s.Fields {
		f := &s.Fields[i]
		if !f.omit(v.FieldByIndex(f.Index)) {
			n++
		}
	}

	if err = e.encodeMapBegin(n); err != nil {
		return
	}
	n = 0

	for i := range s.Fields {
		f := &s.Fields[i]
		if fv := v.FieldByIndex(f.Index); !f.omit(fv) {
			if n != 0 {
				if err = e.encodeMapNext(); err != nil {
					return
				}
			}
			if err = e.encodeString(f.Name); err != nil {
				return
			}
			if err = e.encodeMapValue(); err != nil {
				return
			}
			if err = f.encode(e, fv); err != nil {
				return
			}
			n++
		}
	}

	return e.encodeMapEnd()
}

func (e Encoder) encodeValuePointer(v reflect.Value) error {
	return e.encodeValuePointerWith(v, encodeFuncOf(v.Type().Elem()))
}

func (e Encoder) encodeValuePointerWith(v reflect.Value, f encodeFunc) error {
	if v.IsNil() {
		return e.encodeNil()
	}
	return f(e, v.Elem())
}

func (e Encoder) encodeValueEncoder(v reflect.Value) error {
	return v.Interface().(ValueEncoder).EncodeValue(e)
}

func (e Encoder) encodeValueTextMarshaler(v reflect.Value) error {
	b, err := v.Interface().(encoding.TextMarshaler).MarshalText()
	if err == nil {
		err = e.encodeString(stringNoCopy(b))
	}
	return err
}

func (e Encoder) encodeValueUnsupported(v reflect.Value) error {
	return fmt.Errorf("objconv: the encoder doesn't support values of type %s", v.Type())
}

func (e Encoder) encodeNil() error { return e.Emitter.EmitNil() }

func (e Encoder) encodeBool(v bool) error { return e.Emitter.EmitBool(v) }

func (e Encoder) encodeInt(v int64) error { return e.Emitter.EmitInt(v) }

func (e Encoder) encodeUint(v uint64) error { return e.Emitter.EmitUint(v) }

func (e Encoder) encodeFloat(v float64) error { return e.Emitter.EmitFloat(v) }

func (e Encoder) encodeString(v string) error { return e.Emitter.EmitString(v) }

func (e Encoder) encodeBytes(v []byte) error { return e.Emitter.EmitBytes(v) }

func (e Encoder) encodeTime(v time.Time) error { return e.Emitter.EmitTime(v) }

func (e Encoder) encodeDuration(v time.Duration) error { return e.Emitter.EmitDuration(v) }

func (e Encoder) encodeError(v error) error { return e.Emitter.EmitError(v) }

// EncodeArray provides the implementation of the array encoding algorithm,
// where n is the number of elements in the array, and f a function called to
// encode each element.
//
// The n argument can be set to a negative value to indicate that the program
// doesn't know how many elements it will output to the array. Be mindful that
// not all emitters support encoding arrays of unknown lengths.
//
// The f function is called to encode each element of the array.
func (e Encoder) EncodeArray(n int, f func(Encoder) error) (err error) {
	if e.key {
		if e.key, err = false, e.encodeMapValue(); err != nil {
			return
		}
	}

	if err = e.encodeArrayBegin(n); err != nil {
		return
	}

encodeArray:
	for i := 0; n < 0 || i < n; i++ {
		if i != 0 {
			if e.encodeArrayNext(); err != nil {
				return
			}
		}
		switch err = f(e); err {
		case nil:
		case End:
			break encodeArray
		default:
			return
		}
	}

	return e.encodeArrayEnd()
}

func (e Encoder) encodeArrayBegin(n int) error { return e.Emitter.EmitArrayBegin(n) }

func (e Encoder) encodeArrayEnd() error { return e.Emitter.EmitArrayEnd() }

func (e Encoder) encodeArrayNext() error { return e.Emitter.EmitArrayNext() }

// EncodeMap provides the implementation of the map encoding algorithm, where n
// is the number of elements in the map, and f a function called to encode each
// element.
//
// The n argument can be set to a negative value to indicate that the program
// doesn't know how many elements it will output to the map. Be mindful that not
// all emitters support encoding maps of unknown length.
//
// The f function is called to encode each element of the map, it is expected to
// encode two values, the first one being the key, follow by the associated value.
// The first encoder must be used to encode the key, the second for the value.
func (e Encoder) EncodeMap(n int, f func(Encoder, Encoder) error) (err error) {
	if e.key {
		if e.key, err = false, e.encodeMapValue(); err != nil {
			return
		}
	}

	if err = e.encodeMapBegin(n); err != nil {
		return
	}

encodeMap:
	for i := 0; n < 0 || i < n; i++ {
		if i != 0 {
			if err = e.encodeMapNext(); err != nil {
				return
			}
		}
		e.key = true
		err = f(
			Encoder{Emitter: e.Emitter, SortMapKeys: e.SortMapKeys},
			Encoder{Emitter: e.Emitter, SortMapKeys: e.SortMapKeys, key: true},
		)
		// Because internal calls don't use the exported methods they may not
		// reset this flag to false when expected, forcing the value here.
		e.key = false

		switch err {
		case nil:
		case End:
			break encodeMap
		default:
			return
		}
	}

	return e.encodeMapEnd()
}

func (e Encoder) encodeMapBegin(n int) error { return e.Emitter.EmitMapBegin(n) }

func (e Encoder) encodeMapEnd() error { return e.Emitter.EmitMapEnd() }

func (e Encoder) encodeMapValue() error { return e.Emitter.EmitMapValue() }

func (e Encoder) encodeMapNext() error { return e.Emitter.EmitMapNext() }

/*
// A StreamEncoder encodes and writes a stream of values to an output stream.
//
// Instances of StreamEncoder are not safe for use by multiple goroutines.
type StreamEncoder struct {
	enc    Encoder
	err    error
	len    int
	off    int
	opened bool
	closed bool
}

// NewStreamEncoder returns a new stream encoder that outputs to emitter.
func NewStreamEncoder(emitter Emitter) *StreamEncoder {
	return NewStreamEncoderWith(EncoderConfig{
		Emitter: emitter,
	})
}

// NewStreamEncoder returns a new stream encoder configured with config.
func NewStreamEncoderWith(config EncoderConfig) *StreamEncoder {
	if config.Emitter == nil {
		panic("objconv.NewStreamEncoder: the emitter is nil")
	}
	return &StreamEncoder{
		enc: Encoder{
			e:    config.Emitter,
			sort: config.SortMapKeys,
		},
	}
}

// Open explicitly tells the encoder to start the stream, setting the number
// of values to n.
//
// Depending on the actual format that the stream is encoding to, n may or
// may not have to be accurate, some formats also support passing a negative
// value to indicate that the number of elements is unknown.
func (e *StreamEncoder) Open(n int) error {
	if err := e.err; err != nil {
		return err
	}

	if e.closed {
		return io.ErrClosedPipe
	}

	if !e.opened {
		e.len = n
		e.opened = true
		e.err = e.enc.encodeArrayBegin(n)
	}

	return e.err
}

// Close terminates the stream encoder.
func (e *StreamEncoder) Close() error {
	if err := e.Open(-1); err != nil {
		return err
	}
	if !e.closed {
		e.closed = true
		e.err = e.enc.encodeArrayEnd()
	}
	return e.err
}

// Encode writes v to the stream, encoding it based on the emitter configured
// on e.
func (e *StreamEncoder) Encode(v interface{}) error {
	if err := e.Open(-1); err != nil {
		return err
	}

	if e.off != 0 {
		e.err = e.enc.encodeArrayNext()
	}

	if e.err == nil {
		e.err = e.enc.Encode(v)

		if e.off++; e.len >= 0 && e.off >= e.len {
			e.Close()
		}
	}

	return e.err
}
*/

// ValueEncoder is the interface that can be implemented by types that wish to
// provide their own encoding algorithms.
//
// The EncodeValue method is called when the value is found by an encoding
// algorithm.
type ValueEncoder interface {
	EncodeValue(Encoder) error
}

// ValueEncoderFunc allows the use of regular functions or methods as value
// encoders.
type ValueEncoderFunc func(Encoder) error

// EncodeValue calls f(e).
func (f ValueEncoderFunc) EncodeValue(e Encoder) error { return f(e) }

var (
	stringKeysPool = sync.Pool{
		New: func() interface{} { return make([]string, 0, 20) },
	}
)

// encodeFuncOpts is used to configure how the encodeFuncOf behaves.
type encodeFuncOpts struct {
	recurse bool
	structs map[reflect.Type]*Struct
}

// encodeFunc is the prototype of functions that encode values.
type encodeFunc func(Encoder, reflect.Value) error

// encodeFuncOf returns an encoder function for t.
func encodeFuncOf(t reflect.Type) encodeFunc {
	return makeEncodeFunc(t, encodeFuncOpts{})
}

func makeEncodeFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	switch {
	case t == timeType:
		return Encoder.encodeValueTime

	case t == durationType:
		return Encoder.encodeValueDuration

	case t.Implements(valueEncoderInterface):
		return Encoder.encodeValueEncoder

	case t.Implements(textMarshalerInterface):
		return Encoder.encodeValueTextMarshaler

	case t.Implements(errorInterface):
		return Encoder.encodeValueError
	}

	switch t.Kind() {
	case reflect.Bool:
		return Encoder.encodeValueBool

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Encoder.encodeValueInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Encoder.encodeValueUint

	case reflect.Float32, reflect.Float64:
		return Encoder.encodeValueFloat

	case reflect.String:
		return Encoder.encodeValueString

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return Encoder.encodeValueBytes
		}
		return makeEncodeArrayFunc(t, opts)

	case reflect.Array:
		return makeEncodeArrayFunc(t, opts)

	case reflect.Map:
		return makeEncodeMapFunc(t, opts)

	case reflect.Struct:
		return makeEncodeStructFunc(t, opts)

	case reflect.Ptr, reflect.Interface:
		return makeEncodePtrFunc(t, opts)

	default:
		return Encoder.encodeValueUnsupported
	}
}

func makeEncodeArrayFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	if !opts.recurse {
		return Encoder.encodeValueArray
	}
	f := makeEncodeFunc(t.Elem(), opts)
	return func(e Encoder, v reflect.Value) error {
		return e.encodeValueArrayWith(v, f)
	}
}

func makeEncodeMapFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	if !opts.recurse {
		return Encoder.encodeValueMap
	}
	kf := makeEncodeFunc(t.Key(), opts)
	vf := makeEncodeFunc(t.Elem(), opts)
	return func(e Encoder, v reflect.Value) error {
		return e.encodeValueMapWith(v, kf, vf)
	}
}

func makeEncodeStructFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	if !opts.recurse {
		return Encoder.encodeValueStruct
	}
	s := newStruct(t, opts.structs)
	return func(e Encoder, v reflect.Value) error {
		return e.encodeValueStructWith(v, s)
	}
}

func makeEncodePtrFunc(t reflect.Type, opts encodeFuncOpts) encodeFunc {
	if !opts.recurse {
		return Encoder.encodeValuePointer
	}
	f := makeEncodeFunc(t.Elem(), opts)
	return func(e Encoder, v reflect.Value) error {
		return e.encodeValuePointerWith(v, f)
	}
}
