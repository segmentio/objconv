package objconv

import (
	"bytes"
	"io"
	"reflect"
	"time"
)

// Encode encodes the value to the output in the specified format.
//
// The format must be a string describing the content type of the data
// (like json, resp, ...).
func Encode(out io.Writer, format string, value interface{}) (err error) {
	defer func() { err = convertPanicToError(recover()) }()
	return NewEncoder(EncoderConfig{
		Output:  out,
		Emitter: NewEmitter(format),
	}).Encode(value)
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
		w: Writer{W: config.Output},
		e: config.Emitter,
		t: config.Tag,
	}
}

// NewStreamEncoder returns a new stream encoder configured with config.
func NewStreamEncoder(config EncoderConfig) StreamEncoder {
	config = setEncoderConfigDefault(config)
	return &streamEncoder{
		encoder: encoder{
			w: Writer{W: config.Output},
			e: config.Emitter,
			t: config.Tag,
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
	w Writer
	e Emitter
	t string
}

func (e *encoder) Encode(v interface{}) (err error) {
	defer func() { err = convertPanicToError(recover()) }()
	e.encodeBegin(&e.w)
	e.encode(&e.w, v)
	e.encodeEnd(&e.w)
	return
}

func (e *encoder) encode(w *Writer, v interface{}) {
	if v == nil {
		e.encodeNil(w)
		return
	}

	switch x := v.(type) {
	case string:
		e.encodeString(w, x)

	case []byte:
		e.encodeBytes(w, x)

	case bool:
		e.encodeBool(w, x)

	case int:
		e.encodeInt(w, x)

	case int8:
		e.encodeInt8(w, x)

	case int16:
		e.encodeInt16(w, x)

	case int32:
		e.encodeInt32(w, x)

	case int64:
		e.encodeInt64(w, x)

	case uint:
		e.encodeUint(w, x)

	case uint8:
		e.encodeUint8(w, x)

	case uint16:
		e.encodeUint16(w, x)

	case uint32:
		e.encodeUint32(w, x)

	case uint64:
		e.encodeUint64(w, x)

	case uintptr:
		e.encodeUintptr(w, x)

	case float32:
		e.encodeFloat32(w, x)

	case float64:
		e.encodeFloat64(w, x)

	case time.Time:
		e.encodeTime(w, x)

	case time.Duration:
		e.encodeDuration(w, x)

	case Array:
		e.encodeArray(w, x)

	case Map:
		e.encodeMap(w, x)

	case error:
		e.encodeError(w, x)

	case []rune:
		e.encodeString(w, string(x))

	default:
		e.encodeValue(w, reflect.ValueOf(v))
	}
}

func (e *encoder) encodeValue(w *Writer, v reflect.Value) {
	t := v.Type()
	k := t.Kind()

	switch k {
	case reflect.Bool:
		e.encodeBool(w, v.Bool())

	case reflect.Int:
		e.encodeInt(w, int(v.Int()))

	case reflect.Int8:
		e.encodeInt8(w, int8(v.Int()))

	case reflect.Int16:
		e.encodeInt16(w, int16(v.Int()))

	case reflect.Int32:
		e.encodeInt32(w, int32(v.Int()))

	case reflect.Int64:
		e.encodeInt64(w, v.Int())

	case reflect.Uint:
		e.encodeUint(w, uint(v.Uint()))

	case reflect.Uint8:
		e.encodeUint8(w, uint8(v.Uint()))

	case reflect.Uint16:
		e.encodeUint16(w, uint16(v.Uint()))

	case reflect.Uint32:
		e.encodeUint32(w, uint32(v.Uint()))

	case reflect.Uint64:
		e.encodeUint64(w, v.Uint())

	case reflect.Uintptr:
		e.encodeUintptr(w, uintptr(v.Uint()))

	case reflect.Float32:
		e.encodeFloat32(w, float32(v.Float()))

	case reflect.Float64:
		e.encodeFloat64(w, float64(v.Float()))

	case reflect.String:
		e.encodeString(w, v.String())

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			e.encodeBytes(w, v.Bytes())
		} else {
			e.encodeArray(w, ArraySlice(v))
		}

	case reflect.Array:
		e.encodeArray(w, ArraySlice(v))

	case reflect.Map:
		e.encodeMap(w, MapMap(v))

	case reflect.Struct:
		e.encodeMap(w, newMapStruct(e.t, v))

	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			e.encodeNil(w)
		} else {
			e.encode(w, v.Elem().Interface())
		}

	default:
		panic(&UnsupportedTypeError{t})
	}
}

func (e *encoder) encodeBegin(w *Writer) { e.e.EmitBegin(w) }

func (e *encoder) encodeEnd(w *Writer) { e.e.EmitEnd(w) }

func (e *encoder) encodeNil(w *Writer) { e.e.EmitNil(w) }

func (e *encoder) encodeBool(w *Writer, v bool) { e.e.EmitBool(w, v) }

func (e *encoder) encodeInt(w *Writer, v int) { e.e.EmitInt(w, v) }

func (e *encoder) encodeInt8(w *Writer, v int8) { e.e.EmitInt8(w, v) }

func (e *encoder) encodeInt16(w *Writer, v int16) { e.e.EmitInt16(w, v) }

func (e *encoder) encodeInt32(w *Writer, v int32) { e.e.EmitInt32(w, v) }

func (e *encoder) encodeInt64(w *Writer, v int64) { e.e.EmitInt64(w, v) }

func (e *encoder) encodeUint(w *Writer, v uint) { e.e.EmitUint(w, v) }

func (e *encoder) encodeUint8(w *Writer, v uint8) { e.e.EmitUint8(w, v) }

func (e *encoder) encodeUint16(w *Writer, v uint16) { e.e.EmitUint16(w, v) }

func (e *encoder) encodeUint32(w *Writer, v uint32) { e.e.EmitUint32(w, v) }

func (e *encoder) encodeUint64(w *Writer, v uint64) { e.e.EmitUint64(w, v) }

func (e *encoder) encodeUintptr(w *Writer, v uintptr) { e.e.EmitUintptr(w, v) }

func (e *encoder) encodeFloat32(w *Writer, v float32) { e.e.EmitFloat32(w, v) }

func (e *encoder) encodeFloat64(w *Writer, v float64) { e.e.EmitFloat64(w, v) }

func (e *encoder) encodeString(w *Writer, v string) { e.e.EmitString(w, v) }

func (e *encoder) encodeBytes(w *Writer, v []byte) { e.e.EmitBytes(w, v) }

func (e *encoder) encodeTime(w *Writer, v time.Time) { e.e.EmitTime(w, v) }

func (e *encoder) encodeDuration(w *Writer, v time.Duration) { e.e.EmitDuration(w, v) }

func (e *encoder) encodeError(w *Writer, v error) { e.e.EmitError(w, v) }

func (e *encoder) encodeArray(w *Writer, v Array) {
	e.encodeArrayBegin(w, v.Len())
	it := v.Iter()

	for i := 0; true; i++ {
		if v, ok := it.Next(); !ok {
			break
		} else {
			if i != 0 {
				e.encodeArrayNext(w)
			}
			e.encode(w, v)
		}
	}

	e.encodeArrayEnd(w)
}

func (e *encoder) encodeArrayBegin(w *Writer, n int) { e.e.EmitArrayBegin(w, n) }

func (e *encoder) encodeArrayEnd(w *Writer) { e.e.EmitArrayEnd(w) }

func (e *encoder) encodeArrayNext(w *Writer) { e.e.EmitArrayNext(w) }

func (e *encoder) encodeMap(w *Writer, v Map) {
	e.encodeMapBegin(w, v.Len())
	it := v.Iter()

	for i := 0; true; i++ {
		if item, ok := it.Next(); !ok {
			break
		} else {
			if i != 0 {
				e.encodeMapNext(w)
			}
			e.encode(w, item.Key)
			e.encodeMapValue(w)
			e.encode(w, item.Value)
		}
	}

	e.encodeMapEnd(w)
}

func (e *encoder) encodeMapBegin(w *Writer, n int) { e.e.EmitMapBegin(w, n) }

func (e *encoder) encodeMapEnd(w *Writer) { e.e.EmitMapEnd(w) }

func (e *encoder) encodeMapValue(w *Writer) { e.e.EmitMapValue(w) }

func (e *encoder) encodeMapNext(w *Writer) { e.e.EmitMapNext(w) }

type streamEncoder struct {
	encoder
	err    error
	len    int
	off    int
	opened bool
	closed bool
}

func (e *streamEncoder) Open(n int) (err error) {
	defer func() { err = e.convertPanicToError(recover()) }()
	e.check()
	e.open(n)
	return
}

func (e *streamEncoder) Close() (err error) {
	defer func() { err = e.convertPanicToError(recover()) }()
	e.check()
	e.open(-1)
	e.close()
	return
}

func (e *streamEncoder) Encode(v interface{}) (err error) {
	defer func() { err = e.convertPanicToError(recover()) }()
	e.check()
	e.open(-1)

	if e.off != 0 {
		e.encodeArrayNext(&e.w)
	}

	e.encode(&e.w, v)

	if e.off++; e.len >= 0 && e.off >= e.len {
		e.close()
	}

	return
}

func (e *streamEncoder) check() {
	if e.err != nil {
		panic(e.err)
	}
	if e.closed {
		panic(io.ErrClosedPipe)
	}
}

func (e *streamEncoder) open(n int) {
	if !e.opened {
		e.len = n
		e.opened = true
		e.encodeBegin(&e.w)
		e.encodeArrayBegin(&e.w, n)
	}
}

func (e *streamEncoder) close() {
	if !e.closed {
		e.closed = true
		e.encodeArrayEnd(&e.w)
		e.encodeEnd(&e.w)
	}
}

func (e *streamEncoder) convertPanicToError(v interface{}) (err error) {
	if err = convertPanicToError(v); err != nil {
		e.err = err
	}
	return
}

type nonstreamEncoder struct {
	encoder
	err    error
	opened bool
	closed bool
}

func (e *nonstreamEncoder) Open(n int) (err error) {
	defer func() { err = e.convertPanicToError(recover()) }()
	e.check()
	e.open()
	return
}

func (e *nonstreamEncoder) Close() (err error) {
	defer func() { err = e.convertPanicToError(recover()) }()
	e.check()
	e.open()
	e.close()
	return
}

func (e *nonstreamEncoder) Encode(v interface{}) (err error) {
	defer func() { err = e.convertPanicToError(recover()) }()
	e.check()
	e.open()
	e.encode(&e.w, v)
	e.close()
	return
}

func (e *nonstreamEncoder) check() {
	if e.err != nil {
		panic(e.err)
	}
	if e.closed {
		panic(io.ErrClosedPipe)
	}
}

func (e *nonstreamEncoder) open() {
	if !e.opened {
		e.opened = true
		e.encodeBegin(&e.w)
	}
}

func (e *nonstreamEncoder) close() {
	if !e.closed {
		e.closed = true
		e.encodeEnd(&e.w)
	}
}

func (e *nonstreamEncoder) convertPanicToError(v interface{}) (err error) {
	if err = convertPanicToError(v); err != nil {
		e.err = err
	}
	return
}
