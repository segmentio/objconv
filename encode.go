package objconv

import (
	"encoding"
	"fmt"
	"io"
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
	SortMapKeys bool    // whether map keys should be sorted
	key         bool
}

// NewEncoder returns a new encoder that outputs values to e.
//
// Encoders created by this function use the default encoder configuration,
// which is equivalent to using a zero-value EncoderConfig with only the Emitter
// field set.
//
// The function panics if e is nil.
func NewEncoder(e Emitter) *Encoder {
	if e == nil {
		panic("objconv: the emitter is nil")
	}
	return &Encoder{Emitter: e}
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
	return e.encodeInt(v.Int(), 0)
}

func (e Encoder) encodeValueInt8(v reflect.Value) error {
	return e.encodeInt(v.Int(), 8)
}

func (e Encoder) encodeValueInt16(v reflect.Value) error {
	return e.encodeInt(v.Int(), 16)
}

func (e Encoder) encodeValueInt32(v reflect.Value) error {
	return e.encodeInt(v.Int(), 32)
}

func (e Encoder) encodeValueInt64(v reflect.Value) error {
	return e.encodeInt(v.Int(), 64)
}

func (e Encoder) encodeValueUint(v reflect.Value) error {
	return e.encodeUint(v.Uint(), 0)
}

func (e Encoder) encodeValueUint8(v reflect.Value) error {
	return e.encodeUint(v.Uint(), 8)
}

func (e Encoder) encodeValueUint16(v reflect.Value) error {
	return e.encodeUint(v.Uint(), 16)
}

func (e Encoder) encodeValueUint32(v reflect.Value) error {
	return e.encodeUint(v.Uint(), 32)
}

func (e Encoder) encodeValueUint64(v reflect.Value) error {
	return e.encodeUint(v.Uint(), 64)
}

func (e Encoder) encodeValueUintptr(v reflect.Value) error {
	return e.encodeUint(v.Uint(), 0)
}

func (e Encoder) encodeValueFloat32(v reflect.Value) error {
	return e.encodeFloat(v.Float(), 32)
}

func (e Encoder) encodeValueFloat64(v reflect.Value) error {
	return e.encodeFloat(v.Float(), 64)
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
	t := v.Type()

	if !e.SortMapKeys {
		switch {
		case t.ConvertibleTo(mapInterfaceInterfaceType):
			return e.encodeValueMapInterfaceInterface(v.Convert(mapInterfaceInterfaceType))

		case t.ConvertibleTo(mapStringInterfaceType):
			return e.encodeValueMapStringInterface(v.Convert(mapStringInterfaceType))

		case t.ConvertibleTo(mapStringStringType):
			return e.encodeValueMapStringString(v.Convert(mapStringStringType))
		}
	}

	var k []reflect.Value
	var n = v.Len()
	var i = 0

	if n != 0 {
		k = v.MapKeys()

		if e.SortMapKeys {
			sortValues(t.Key(), k)
		}
	}

	return e.EncodeMap(n, func(ke Encoder, ve Encoder) (err error) {
		if err = kf(e, k[i]); err != nil {
			return
		}
		if err = e.encodeMapValue(); err != nil {
			return
		}
		if err = vf(e, v.MapIndex(k[i])); err != nil {
			return
		}
		i++
		return
	})
}

func (e Encoder) encodeValueMapInterfaceInterface(v reflect.Value) (err error) {
	m := v.Interface().(map[interface{}]interface{})
	n := len(m)
	i := 0

	if err = e.encodeMapBegin(n); err != nil {
		return
	}

	for k, v := range m {
		if i != 0 {
			if err = e.encodeMapNext(); err != nil {
				return
			}
		}
		if err = e.Encode(k); err != nil {
			return
		}
		if err = e.encodeMapValue(); err != nil {
			return
		}
		if err = e.Encode(v); err != nil {
			return
		}
		i++
	}

	return e.encodeMapEnd()
}

func (e Encoder) encodeValueMapStringInterface(v reflect.Value) (err error) {
	m := v.Interface().(map[string]interface{})
	n := len(m)
	i := 0

	if err = e.encodeMapBegin(n); err != nil {
		return
	}

	for k, v := range m {
		if i != 0 {
			if err = e.encodeMapNext(); err != nil {
				return
			}
		}
		if err = e.encodeString(k); err != nil {
			return
		}
		if err = e.encodeMapValue(); err != nil {
			return
		}
		if err = e.Encode(v); err != nil {
			return
		}
		i++
	}

	return e.encodeMapEnd()
}

func (e Encoder) encodeValueMapStringString(v reflect.Value) (err error) {
	m := v.Interface().(map[string]string)
	n := len(m)
	i := 0

	if err = e.encodeMapBegin(n); err != nil {
		return
	}

	for k, v := range m {
		if i != 0 {
			if err = e.encodeMapNext(); err != nil {
				return
			}
		}
		if err = e.encodeString(k); err != nil {
			return
		}
		if err = e.encodeMapValue(); err != nil {
			return
		}
		if err = e.encodeString(v); err != nil {
			return
		}
		i++
	}

	return e.encodeMapEnd()
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

func (e Encoder) encodeValueInterface(v reflect.Value) error {
	if v.IsNil() {
		return e.encodeNil()
	}
	return e.encodeValue(v.Elem())
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

func (e Encoder) encodeInt(v int64, bs int) error { return e.Emitter.EmitInt(v, bs) }

func (e Encoder) encodeUint(v uint64, bs int) error { return e.Emitter.EmitUint(v, bs) }

func (e Encoder) encodeFloat(v float64, bs int) error { return e.Emitter.EmitFloat(v, bs) }

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

// A StreamEncoder encodes and writes a stream of values to an output stream.
//
// Instances of StreamEncoder are not safe for use by multiple goroutines.
type StreamEncoder struct {
	Emitter     Emitter // the emiiter used by this encoder
	SortMapKeys bool    // whether map keys should be sorted

	err     error
	max     int
	cnt     int
	opened  bool
	closed  bool
	oneshot bool
}

// NewStreamEncoder returns a new stream encoder that outputs to e.
//
// The function panics if e is nil.
func NewStreamEncoder(e Emitter) *StreamEncoder {
	if e == nil {
		panic("objconv.NewStreamEncoder: the emitter is nil")
	}
	return &StreamEncoder{Emitter: e}
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
		e.max = n
		e.opened = true

		if !e.oneshot {
			e.err = e.encoder().encodeArrayBegin(n)
		}
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

		if !e.oneshot {
			e.err = e.encoder().encodeArrayEnd()
		}
	}

	return e.err
}

// Encode writes v to the stream, encoding it based on the emitter configured
// on e.
func (e *StreamEncoder) Encode(v interface{}) error {
	if err := e.Open(-1); err != nil {
		return err
	}

	if e.max >= 0 && e.cnt >= e.max {
		return fmt.Errorf("objconv: too many values sent to a stream encoder exceed the configured limit of %d", e.max)
	}

	enc := e.encoder()

	if !e.oneshot && e.cnt != 0 {
		e.err = enc.encodeArrayNext()
	}

	if e.err == nil {
		e.err = enc.Encode(v)

		if e.cnt++; e.max >= 0 && e.cnt >= e.max {
			e.Close()
		}
	}

	return e.err
}

func (e *StreamEncoder) encoder() Encoder {
	return Encoder{
		Emitter:     e.Emitter,
		SortMapKeys: e.SortMapKeys,
	}
}

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
	switch t {
	case boolType:
		return Encoder.encodeValueBool

	case stringType:
		return Encoder.encodeValueString

	case bytesType:
		return Encoder.encodeValueBytes

	case timeType:
		return Encoder.encodeValueTime

	case durationType:
		return Encoder.encodeValueDuration

	case emptyInterface:
		return Encoder.encodeValueInterface

	case intType:
		return Encoder.encodeValueInt

	case int8Type:
		return Encoder.encodeValueInt8

	case int16Type:
		return Encoder.encodeValueInt16

	case int32Type:
		return Encoder.encodeValueInt32

	case int64Type:
		return Encoder.encodeValueInt64

	case uintType:
		return Encoder.encodeValueUint

	case uint8Type:
		return Encoder.encodeValueUint8

	case uint16Type:
		return Encoder.encodeValueUint16

	case uint32Type:
		return Encoder.encodeValueUint32

	case uint64Type:
		return Encoder.encodeValueUint64

	case uintptrType:
		return Encoder.encodeValueUintptr

	case float32Type:
		return Encoder.encodeValueFloat32

	case float64Type:
		return Encoder.encodeValueFloat64
	}

	switch {
	case t.Implements(valueEncoderInterface):
		return Encoder.encodeValueEncoder

	case t.Implements(textMarshalerInterface):
		return Encoder.encodeValueTextMarshaler

	case t.Implements(errorInterface):
		return Encoder.encodeValueError
	}

	switch t.Kind() {
	case reflect.Struct:
		return makeEncodeStructFunc(t, opts)

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return Encoder.encodeValueBytes
		}
		return makeEncodeArrayFunc(t, opts)

	case reflect.Map:
		return makeEncodeMapFunc(t, opts)

	case reflect.Ptr:
		return makeEncodePtrFunc(t, opts)

	case reflect.Array:
		return makeEncodeArrayFunc(t, opts)

	case reflect.String:
		return Encoder.encodeValueString

	case reflect.Bool:
		return Encoder.encodeValueBool

	case reflect.Int:
		return Encoder.encodeValueInt

	case reflect.Int8:
		return Encoder.encodeValueInt8

	case reflect.Int16:
		return Encoder.encodeValueInt16

	case reflect.Int32:
		return Encoder.encodeValueInt32

	case reflect.Int64:
		return Encoder.encodeValueInt64

	case reflect.Uint:
		return Encoder.encodeValueUint

	case reflect.Uint8:
		return Encoder.encodeValueUint8

	case reflect.Uint16:
		return Encoder.encodeValueUint16

	case reflect.Uint32:
		return Encoder.encodeValueUint32

	case reflect.Uint64:
		return Encoder.encodeValueUint64

	case reflect.Uintptr:
		return Encoder.encodeValueUintptr

	case reflect.Float32:
		return Encoder.encodeValueFloat32

	case reflect.Float64:
		return Encoder.encodeValueFloat64

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
