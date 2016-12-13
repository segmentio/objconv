package objconv

import (
	"io"
	"reflect"
	"sync"
	"time"
)

// EncoderConfig carries the different configuration options that can be set
// when instantiating an instance of Encoder.
type EncoderConfig struct {
	// Emitter defines the format used by the encoder.
	Emitter Emitter

	// SortMapKeys controls whether the encoder will sort map keys or not.
	SortMapKeys bool
}

// An Encoder implements the high-level encoding algorithm that inspect encoded
// values and drive the use of an Emitter to create a serialized representation
// of the data.
//
// Instances of Encoder are not safe for use by multiple goroutines.
type Encoder struct {
	e    Emitter
	sort bool
	key  bool
}

// NewEncoder returns a new encoder that outputs values to emitter.
//
// Encoders created by this function use the default encoder configuration,
// which is equivalent to using a zero-value EncoderConfig with only the Emitter
// field set.
//
// The function panics if emitter is nil.
func NewEncoder(emitter Emitter) *Encoder {
	return NewEncoderWith(EncoderConfig{
		Emitter: emitter,
	})
}

// NewEncoderWith returns a new encoder configured with config.
//
// The function panics if config.Emitter is nil.
func NewEncoderWith(config EncoderConfig) *Encoder {
	if config.Emitter == nil {
		panic("objconv.NewEncoder: the emitter is nil")
	}
	return &Encoder{
		e:    config.Emitter,
		sort: config.SortMapKeys,
	}
}

// Encode encodes the generic value v.
func (e *Encoder) Encode(v interface{}) error {
	if err := e.encodeMapValueMaybe(); err != nil {
		return err
	}

	if v == nil {
		return e.encodeNil()
	}

	switch x := v.(type) {
	case string:
		return e.encodeString(x)

	case []byte:
		return e.encodeBytes(x)

	case bool:
		return e.encodeBool(x)

	case int:
		return e.encodeInt(int64(x))

	case int8:
		return e.encodeInt(int64(x))

	case int16:
		return e.encodeInt(int64(x))

	case int32:
		return e.encodeInt(int64(x))

	case int64:
		return e.encodeInt(x)

	case uint:
		return e.encodeUint(uint64(x))

	case uint8:
		return e.encodeUint(uint64(x))

	case uint16:
		return e.encodeUint(uint64(x))

	case uint32:
		return e.encodeUint(uint64(x))

	case uint64:
		return e.encodeUint(uint64(x))

	case uintptr:
		return e.encodeUint(uint64(x))

	case float32:
		return e.encodeFloat(float64(x))

	case float64:
		return e.encodeFloat(x)

	case time.Time:
		return e.encodeTime(x)

	case time.Duration:
		return e.encodeDuration(x)

	case error:
		return e.encodeError(x)

	default:
		return e.encodeValue(reflect.ValueOf(v))
	}
}

var (
	timeType     = reflect.TypeOf(time.Time{})
	durationType = reflect.TypeOf(time.Duration(0))
)

func (e *Encoder) encodeValue(v reflect.Value) error {
	t := v.Type()

	switch k := t.Kind(); k {
	case reflect.Bool:
		return e.encodeBool(v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.encodeInt(v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return e.encodeUint(v.Uint())

	case reflect.Float32, reflect.Float64:
		return e.encodeFloat(v.Float())

	case reflect.String:
		return e.encodeString(v.String())

	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			if k == reflect.Array {
				v = v.Slice(0, v.Len())
			}
			return e.encodeBytes(v.Bytes())
		}
		return e.encodeValueArray(v)

	case reflect.Map:
		switch {
		case v.Len() == 0:
			return e.EncodeMap(0, nil)
		case e.sort:
			return e.encodeSortedValueMap(v)
		default:
			return e.encodeUnsortedValueMap(v)
		}

	case reflect.Struct:
		return e.encodeValueStruct(v)

	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return e.encodeNil()
		}
		return e.Encode(v.Elem().Interface())

	default:
		return &UnsupportedTypeError{t}
	}
}

func (e *Encoder) encodeNil() error { return e.e.EmitNil() }

func (e *Encoder) encodeBool(v bool) error { return e.e.EmitBool(v) }

func (e *Encoder) encodeInt(v int64) error { return e.e.EmitInt(v) }

func (e *Encoder) encodeUint(v uint64) error { return e.e.EmitUint(v) }

func (e *Encoder) encodeFloat(v float64) error { return e.e.EmitFloat(v) }

func (e *Encoder) encodeString(v string) error { return e.e.EmitString(v) }

func (e *Encoder) encodeBytes(v []byte) error { return e.e.EmitBytes(v) }

func (e *Encoder) encodeTime(v time.Time) error { return e.e.EmitTime(v) }

func (e *Encoder) encodeDuration(v time.Duration) error { return e.e.EmitDuration(v) }

func (e *Encoder) encodeError(v error) error { return e.e.EmitError(v) }

func (e *Encoder) encodeValueArray(v reflect.Value) error {
	i := 0
	return e.EncodeArray(v.Len(), func(e *Encoder) (err error) {
		err = e.Encode(v.Index(i).Interface())
		i++
		return
	})
}

// EncodeArray provides the implementation of the array encoding algorithm,
// where n is the number of elements in the array, and f a function called to
// encode each element.
//
// The n argument can be set to a negative value to indicate that the program
// doesn't know how many elements it will output to the array. Be mindful that
// not all emitters support encoding arrays of unknown lengths.
//
// The f function is called to encode each element of the array.
func (e *Encoder) EncodeArray(n int, f func(*Encoder) error) (err error) {
	if err = e.encodeMapValueMaybe(); err != nil {
		return
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

func (e *Encoder) encodeArrayBegin(n int) error { return e.e.EmitArrayBegin(n) }

func (e *Encoder) encodeArrayEnd() error { return e.e.EmitArrayEnd() }

func (e *Encoder) encodeArrayNext() error { return e.e.EmitArrayNext() }

func (e *Encoder) encodeSortedValueMap(v reflect.Value) error {
	t := v.Type().Key()
	k := v.MapKeys()
	sortValues(t, k)
	return e.encodeValueMap(v, k)
}

func (e *Encoder) encodeUnsortedValueMap(v reflect.Value) error {
	return e.encodeValueMap(v, v.MapKeys())
}

func (e *Encoder) encodeValueMap(v reflect.Value, keys []reflect.Value) error {
	i := 0
	return e.EncodeMap(v.Len(), func(e *Encoder) error {
		k := keys[i]
		i++
		if err := e.Encode(k.Interface()); err != nil {
			return err
		}
		return e.Encode(v.MapIndex(k).Interface())
	})
}

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
func (e *Encoder) EncodeMap(n int, f func(*Encoder) error) (err error) {
	if err = e.encodeMapValueMaybe(); err != nil {
		return
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
		switch err = f(e); err {
		case nil:
		case End:
			break encodeMap
		default:
			return
		}
	}

	return e.encodeMapEnd()
}

func (e *Encoder) encodeMapBegin(n int) error { return e.e.EmitMapBegin(n) }

func (e *Encoder) encodeMapEnd() error { return e.e.EmitMapEnd() }

func (e *Encoder) encodeMapValue() error { return e.e.EmitMapValue() }

func (e *Encoder) encodeMapNext() error { return e.e.EmitMapNext() }

func (e *Encoder) encodeMapValueMaybe() (err error) {
	if e.key {
		e.key = false
		err = e.encodeMapValue()
	}
	return
}

func (e *Encoder) encodeValueStruct(v reflect.Value) (err error) {
	s := LookupStruct(v.Type())
	n := 0

	for _, f := range s.Fields {
		if !omit(f, v.FieldByIndex(f.Index)) {
			n++
		}
	}

	if err = e.encodeMapBegin(n); err != nil {
		return
	}
	n = 0

	for _, f := range s.Fields {
		if fv := v.FieldByIndex(f.Index); !omit(f, fv) {
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
			if err = e.Encode(fv.Interface()); err != nil {
				return
			}
			n++
		}
	}

	return e.encodeMapEnd()
}

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

// ValueEncoder is the interface that can be implemented by types that wish to
// provide their own encoding algorithms.
//
// The EncodeValue method is called when the value is found by an encoding
// algorithm.
type ValueEncoder interface {
	EncodeValue(*Encoder) error
}

// ValueEncoderFunc allows the use of regular functions or methods as value
// encoders.
type ValueEncoderFunc func(*Encoder) error

// EncodeValue calls f(e).
func (f ValueEncoderFunc) EncodeValue(e *Encoder) error { return f(e) }

var (
	stringKeysPool = sync.Pool{
		New: func() interface{} { return make([]string, 0, 20) },
	}
)
