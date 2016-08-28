package objconv

import (
	"bytes"
	"io"
	"reflect"
	"sort"
	"sync"
	"time"
)

// Encode encodes the value to the output in the specified format.
//
// The format must be a string describing the content type of the data
// (like json, resp, ...).
func Encode(out io.Writer, format string, value interface{}) (err error) {
	var emitter Emitter

	if emitter, err = GetEmitter(format); err == nil {
		err = NewEncoder(EncoderConfig{
			Output:      out,
			Emitter:     emitter,
			SortMapKeys: true,
		}).Encode(value)
	}

	return
}

// EncodeBytes returns a byte slice containing a representation of the value in
// the specified format.
//
// The format must be a string describing the content type of the data
// (like json, resp, ...).
func EncodeBytes(format string, value interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := Encode(buf, format, value)
	return buf.Bytes(), err
}

// EncodeString returns a byte slice containing a representation of the value in
// the specified format.
//
// The format must be a string describing the content type of the data
// (like json, resp, ...).
func EncodeString(format string, value interface{}) (string, error) {
	b, err := EncodeBytes(format, value)
	return string(b), err
}

// EncodeLength returns the length of the encoded value in format.
//
// The format must be a string describing the content type of the data
// (like json, resp, ...).
func EncodeLength(format string, value interface{}) (n int, err error) {
	c := &counter{}
	err = Encode(c, format, value)
	n = c.n
	return
}

// An Encoder encodes and writes values to an output streams in.
type Encoder interface {
	// Encode writes v to its output.
	Encode(v interface{}) error
}

// EncoderFunc is an adapater to allow use of ordinary functions as encoders.
type EncoderFunc func(interface{}) error

// Encode calls f(v).
func (f EncoderFunc) Encode(v interface{}) error { return f(v) }

// EncoderConfig carries the configuration for creating an encoder.
type EncoderConfig struct {
	// Output is the data stream where the encoder write.
	Output io.Writer

	// Emitter defines the format used by the encoder.
	Emitter Emitter

	// Tag sets the name of the tag used when encoding struct fields.
	Tag string

	// SortMapKeys controls whether the encoder will sort map keys or not.
	SortMapKeys bool
}

// A StreamEncoder encodes and writes a stream of values to an output stream.
type StreamEncoder interface {
	Encoder

	// Open explicitly tells the encoder to start the stream, setting the number
	// of values to n.
	//
	// Depending on the actual format that the stream is encoding to, n may or
	// may not have to be accurate, some formats also support passing a negative
	// value to indicate that the number of elements is unknown.
	Open(n int) error

	// Close terminates the stream encoder.
	Close() error
}

// NewEncoder returns a new encoder configured with config.
func NewEncoder(config EncoderConfig) Encoder {
	config = setEncoderConfigDefault(config)
	return &encoder{
		sort: config.SortMapKeys,
		e:    config.Emitter,
		t:    config.Tag,
		w:    Writer{W: config.Output},
	}
}

// NewStreamEncoder returns a new stream encoder configured with config.
func NewStreamEncoder(config EncoderConfig) StreamEncoder {
	config = setEncoderConfigDefault(config)
	return &streamEncoder{
		encoder: encoder{
			sort: config.SortMapKeys,
			e:    config.Emitter,
			t:    config.Tag,
			w:    Writer{W: config.Output},
		},
	}
}

func setEncoderConfigDefault(config EncoderConfig) EncoderConfig {
	if config.Output == nil {
		panic("objconv.NewEncoder: config.Output is nil")
	}

	if config.Emitter == nil {
		panic("objconv.NewEncoder: config.Emitter is nil")
	}

	if len(config.Tag) == 0 {
		config.Tag = "objconv"
	}

	return config
}

type encoder struct {
	sort bool

	e Emitter
	t string
	w Writer
}

func (e *encoder) Encode(v interface{}) error {
	e.encodeBegin()
	e.encode(v)
	e.encodeEnd()
	return e.w.e
}

func (e *encoder) encode(v interface{}) {
	if v == nil {
		e.encodeNil()
		return
	}

	switch x := v.(type) {
	case string:
		e.encodeString(x)

	case []byte:
		e.encodeBytes(x)

	case bool:
		e.encodeBool(x)

	case int:
		e.encodeInt(x)

	case int8:
		e.encodeInt8(x)

	case int16:
		e.encodeInt16(x)

	case int32:
		e.encodeInt32(x)

	case int64:
		e.encodeInt64(x)

	case uint:
		e.encodeUint(x)

	case uint8:
		e.encodeUint8(x)

	case uint16:
		e.encodeUint16(x)

	case uint32:
		e.encodeUint32(x)

	case uint64:
		e.encodeUint64(x)

	case uintptr:
		e.encodeUintptr(x)

	case float32:
		e.encodeFloat32(x)

	case float64:
		e.encodeFloat64(x)

	case time.Time:
		e.encodeTime(x)

	case time.Duration:
		e.encodeDuration(x)

	case []interface{}:
		e.encodeSliceInterface(x)

	case []string:
		e.encodeSliceString(x)

	case [][]byte:
		e.encodeSliceBytes(x)

	case map[string]string:
		e.encodeMapStringString(x)

	case map[string]interface{}:
		e.encodeMapStringInterface(x)

	case MapSlice:
		e.encodeMapSlice(x)

	case Array:
		e.encodeArray(x)

	case Map:
		e.encodeMap(x)

	case error:
		e.encodeError(x)

	case []rune:
		e.encodeString(string(x))

	default:
		e.encodeValue(reflect.ValueOf(v))
	}
}

func (e *encoder) encodeValue(v reflect.Value) {
	t := v.Type()
	k := t.Kind()

	switch k {
	case reflect.Bool:
		e.encodeBool(v.Bool())

	case reflect.Int:
		e.encodeInt(int(v.Int()))

	case reflect.Int8:
		e.encodeInt8(int8(v.Int()))

	case reflect.Int16:
		e.encodeInt16(int16(v.Int()))

	case reflect.Int32:
		e.encodeInt32(int32(v.Int()))

	case reflect.Int64:
		e.encodeInt64(v.Int())

	case reflect.Uint:
		e.encodeUint(uint(v.Uint()))

	case reflect.Uint8:
		e.encodeUint8(uint8(v.Uint()))

	case reflect.Uint16:
		e.encodeUint16(uint16(v.Uint()))

	case reflect.Uint32:
		e.encodeUint32(uint32(v.Uint()))

	case reflect.Uint64:
		e.encodeUint64(v.Uint())

	case reflect.Uintptr:
		e.encodeUintptr(uintptr(v.Uint()))

	case reflect.Float32:
		e.encodeFloat32(float32(v.Float()))

	case reflect.Float64:
		e.encodeFloat64(float64(v.Float()))

	case reflect.String:
		e.encodeString(v.String())

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			e.encodeBytes(v.Bytes())
		} else {
			e.encodeSliceValue(v)
		}

	case reflect.Array:
		e.encodeSliceValue(v)

	case reflect.Map:
		if e.sort {
			e.encodeMap(SortedMap(v))
		} else {
			e.encodeMap(UnsortedMap(v))
		}

	case reflect.Struct:
		e.encodeStruct(v)

	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			e.encodeNil()
		} else {
			e.encode(v.Elem().Interface())
		}

	default:
		e.w.e = &UnsupportedTypeError{t}
	}
}

func (e *encoder) encodeBegin() { e.e.EmitBegin(&e.w) }

func (e *encoder) encodeEnd() { e.e.EmitEnd(&e.w) }

func (e *encoder) encodeNil() { e.e.EmitNil(&e.w) }

func (e *encoder) encodeBool(v bool) { e.e.EmitBool(&e.w, v) }

func (e *encoder) encodeInt(v int) { e.e.EmitInt(&e.w, v) }

func (e *encoder) encodeInt8(v int8) { e.e.EmitInt8(&e.w, v) }

func (e *encoder) encodeInt16(v int16) { e.e.EmitInt16(&e.w, v) }

func (e *encoder) encodeInt32(v int32) { e.e.EmitInt32(&e.w, v) }

func (e *encoder) encodeInt64(v int64) { e.e.EmitInt64(&e.w, v) }

func (e *encoder) encodeUint(v uint) { e.e.EmitUint(&e.w, v) }

func (e *encoder) encodeUint8(v uint8) { e.e.EmitUint8(&e.w, v) }

func (e *encoder) encodeUint16(v uint16) { e.e.EmitUint16(&e.w, v) }

func (e *encoder) encodeUint32(v uint32) { e.e.EmitUint32(&e.w, v) }

func (e *encoder) encodeUint64(v uint64) { e.e.EmitUint64(&e.w, v) }

func (e *encoder) encodeUintptr(v uintptr) { e.e.EmitUintptr(&e.w, v) }

func (e *encoder) encodeFloat32(v float32) { e.e.EmitFloat32(&e.w, v) }

func (e *encoder) encodeFloat64(v float64) { e.e.EmitFloat64(&e.w, v) }

func (e *encoder) encodeString(v string) { e.e.EmitString(&e.w, v) }

func (e *encoder) encodeBytes(v []byte) { e.e.EmitBytes(&e.w, v) }

func (e *encoder) encodeTime(v time.Time) { e.e.EmitTime(&e.w, v) }

func (e *encoder) encodeDuration(v time.Duration) { e.e.EmitDuration(&e.w, v) }

func (e *encoder) encodeError(v error) { e.e.EmitError(&e.w, v) }

func (e *encoder) encodeSliceInterface(v []interface{}) {
	n := len(v)
	e.encodeArrayBegin(n)
	if n != 0 {
		e.encode(v[0])
		for i := 1; i != n; i++ {
			e.encodeArrayNext()
			e.encode(v[i])
		}
	}
	e.encodeArrayEnd()
}

func (e *encoder) encodeSliceString(v []string) {
	n := len(v)
	e.encodeArrayBegin(n)
	if n != 0 {
		e.encode(v[0])
		for i := 1; i != n; i++ {
			e.encodeArrayNext()
			e.encodeString(v[i])
		}
	}
	e.encodeArrayEnd()
}

func (e *encoder) encodeSliceBytes(v [][]byte) {
	n := len(v)
	e.encodeArrayBegin(n)
	if n != 0 {
		e.encode(v[0])
		for i := 1; i != n; i++ {
			e.encodeArrayNext()
			e.encodeBytes(v[i])
		}
	}
	e.encodeArrayEnd()
}

func (e *encoder) encodeSliceValue(v reflect.Value) {
	n := v.Len()
	e.encodeArrayBegin(n)
	if n != 0 {
		e.encode(v.Index(0).Interface())
		for i := 1; i != n; i++ {
			e.encodeArrayNext()
			e.encode(v.Index(i).Interface())
		}
	}
	e.encodeArrayEnd()
}

func (e *encoder) encodeArray(v Array) {
	e.encodeArrayBegin(v.Len())
	it := v.Iter()
	for i := 0; true; i++ {
		if v, ok := it.Next(); !ok {
			break
		} else {
			if i != 0 {
				e.encodeArrayNext()
			}
			e.encode(v)
		}
	}
	e.encodeArrayEnd()
}

func (e *encoder) encodeArrayBegin(n int) { e.e.EmitArrayBegin(&e.w, n) }

func (e *encoder) encodeArrayEnd() { e.e.EmitArrayEnd(&e.w) }

func (e *encoder) encodeArrayNext() { e.e.EmitArrayNext(&e.w) }

func (e *encoder) encodeMapStringString(v map[string]string) {
	n := len(v)
	e.encodeMapBegin(n)

	if n != 0 {
		if e.sort {
			keys := stringKeysPool.Get().([]string)

			for x := range v {
				keys = append(keys, x)
			}

			sort.Strings(keys)

			for i, k := range keys {
				if i != 0 {
					e.encodeMapNext()
				}
				e.encodeString(k)
				e.encodeMapValue()
				e.encodeString(v[k])
			}

			stringKeysPool.Put(keys[:0])
		} else {
			i := 0
			for k, v := range v {
				if i != 0 {
					e.encodeMapNext()
				}
				e.encodeString(k)
				e.encodeMapValue()
				e.encodeString(v)
				i++
			}
		}
	}

	e.encodeMapEnd()
}

func (e *encoder) encodeMapStringInterface(v map[string]interface{}) {
	n := len(v)
	e.encodeMapBegin(n)

	if n != 0 {
		if e.sort {
			keys := stringKeysPool.Get().([]string)

			for x := range v {
				keys = append(keys, x)
			}

			sort.Strings(keys)

			for i, k := range keys {
				if i != 0 {
					e.encodeMapNext()
				}
				e.encodeString(k)
				e.encodeMapValue()
				e.encode(v[k])
			}

			stringKeysPool.Put(keys[:0])
		} else {
			i := 0
			for k, v := range v {
				if i != 0 {
					e.encodeMapNext()
				}
				e.encodeString(k)
				e.encodeMapValue()
				e.encode(v)
				i++
			}
		}
	}

	e.encodeMapEnd()
}

func (e *encoder) encodeMapSlice(v MapSlice) {
	n := len(v)
	e.encodeMapBegin(n)
	if n != 0 {
		e.encode(v[0].Key)
		e.encodeMapValue()
		e.encode(v[0].Value)
		for i := 1; i != n; i++ {
			e.encodeMapNext()
			e.encode(v[i].Key)
			e.encodeMapValue()
			e.encode(v[i].Value)
		}
	}
	e.encodeMapEnd()
}

func (e *encoder) encodeMap(v Map) {
	e.encodeMapBegin(v.Len())
	it := v.Iter()

	for i := 0; true; i++ {
		if item, ok := it.Next(); !ok {
			break
		} else {
			if i != 0 {
				e.encodeMapNext()
			}
			e.encode(item.Key)
			e.encodeMapValue()
			e.encode(item.Value)
		}
	}

	e.encodeMapEnd()
}

func (e *encoder) encodeMapBegin(n int) { e.e.EmitMapBegin(&e.w, n) }

func (e *encoder) encodeMapEnd() { e.e.EmitMapEnd(&e.w) }

func (e *encoder) encodeMapValue() { e.e.EmitMapValue(&e.w) }

func (e *encoder) encodeMapNext() { e.e.EmitMapNext(&e.w) }

func (e *encoder) encodeStruct(v reflect.Value) {
	f := structFieldPool.Get().([]structField)
	s := LookupStruct(v.Type())

	it := s.IterValue(e.t, v, FilterUnexported|FilterAnonymous|FilterSkipped|FilterOmitempty)

	for {
		if k, v, ok := it.NextValue(); !ok {
			break
		} else {
			f = append(f, structField{name: k, value: v.Interface()})
		}
	}

	n := len(f)
	e.encodeMapBegin(n)
	if n != 0 {
		e.encodeString(f[0].name)
		e.encodeMapValue()
		e.encode(f[0].value)
		for i := 1; i != n; i++ {
			e.encodeMapNext()
			e.encodeString(f[i].name)
			e.encodeMapValue()
			e.encode(f[i].value)
		}
	}
	e.encodeMapEnd()
	structFieldPool.Put(f[:0])
}

type streamEncoder struct {
	encoder
	len    int
	off    int
	opened bool
	closed bool
}

func (e *streamEncoder) Open(n int) (err error) {
	return e.open(n)
}

func (e *streamEncoder) Close() (err error) {
	if err = e.open(-1); err != nil {
		return
	}
	return e.close()
}

func (e *streamEncoder) Encode(v interface{}) (err error) {
	if err = e.open(-1); err != nil {
		return
	}

	if e.off != 0 {
		e.encodeArrayNext()
	}

	e.encode(v)

	if e.off++; e.len >= 0 && e.off >= e.len {
		e.close()
	}

	return e.w.e
}

func (e *streamEncoder) open(n int) (err error) {
	if err = e.w.e; err != nil {
		return
	}

	if e.closed {
		err = io.ErrClosedPipe
		return
	}

	if !e.opened {
		e.len = n
		e.opened = true
		e.encodeBegin()
		e.encodeArrayBegin(n)
	}

	return e.w.e
}

func (e *streamEncoder) close() (err error) {
	if !e.closed {
		e.closed = true
		e.encodeArrayEnd()
		e.encodeEnd()
	}
	return e.w.e
}

type nonstreamEncoder struct {
	encoder
	opened bool
	closed bool
}

func (e *nonstreamEncoder) Open(n int) (err error) {
	return e.open()
}

func (e *nonstreamEncoder) Close() (err error) {
	if err = e.open(); err != nil {
		return
	}
	return e.close()
}

func (e *nonstreamEncoder) Encode(v interface{}) (err error) {
	if err = e.open(); err != nil {
		return
	}
	e.encode(v)
	return e.close()
}

func (e *nonstreamEncoder) open() (err error) {
	if err = e.w.e; err != nil {
		return
	}

	if e.closed {
		err = io.ErrClosedPipe
		return
	}

	if !e.opened {
		e.opened = true
		e.encodeBegin()
	}

	return e.w.e
}

func (e *nonstreamEncoder) close() (err error) {
	if !e.closed {
		e.closed = true
		e.encodeEnd()
	}
	return e.w.e
}

type structField struct {
	name  string
	value interface{}
}

var (
	structFieldPool = sync.Pool{
		New: func() interface{} { return make([]structField, 0, 20) },
	}

	stringKeysPool = sync.Pool{
		New: func() interface{} { return make([]string, 0, 20) },
	}
)
