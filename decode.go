package objconv

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

// Decode decodes the content from the reader into the value.
//
// The format must be a string describing the content type of the data
// (like json, resp, ...).
func Decode(in io.Reader, format string, value interface{}) (err error) {
	defer func() { err = convertPanicToError(recover()) }()
	return NewDecoder(DecoderConfig{
		Input:  in,
		Parser: NewParser(format),
	}).Decode(value)
}

// DecodeBytes decodes the content from the byte slice into the value.
//
// The format must be a string describing the content type of the data
// (like json, resp, ...).
func DecodeBytes(in []byte, format string, value interface{}) (err error) {
	return Decode(bytes.NewReader(in), format, value)
}

// DecodeString decodes the content from the string into the value.
//
// The format must be a string describing the content type of the data
// (like json, resp, ...).
func DecodeString(in string, format string, value interface{}) (err error) {
	return Decode(strings.NewReader(in), format, value)
}

// A Decoder reads and decodes values from an input stream.
type Decoder interface {
	// Decode reads the next value from the its input and stores it in the value
	// pointed to by v.
	Decode(v interface{}) error
}

// DecoderFunc is an adapter to allow use of ordinary functions as decoders.
type DecoderFunc func(interface{}) error

// Decode calls f(v).
func (f DecoderFunc) Decode(v interface{}) error { return f(v) }

// DecoderConfig carries the configuration for creating an encoder.
type DecoderConfig struct {
	// Input is the data stream that the decoder reads from.
	Input io.Reader

	// Parser defines the format used by the decoder.
	Parser Parser

	// Tag sets the name of the tag used when decoding struct fields.
	Tag string
}

// A StreamDecoder reads and decodes a stream of values from an input stream.
type StreamDecoder interface {
	Decoder

	// Len returns the expected number of elements returned from the stream.
	//
	// Depending on the actual format that the stream is decoding this value
	// may or may not be accurate, some formats may also return a negative
	// value to indicate that the number of elements is unknown.
	//
	// Note: the value returned by this method will not be meaningful until the
	// first value was decoded from the stream.
	Len() int

	// Error returns the last error encountered by the decoder.
	Error() error

	// Enoder returns a stream encoder that can be used to re-encode the values
	// read from the decoder.
	//
	// This is useful because the stream decoder abstracts the underlying types
	// of the data it reads, the application cannot tell whether it's reading
	// from a sequence or a single value.
	// If it needs to re-encode the values with the same type that they had
	// before decoding the application needs to use an encoder returned by this
	// method.
	//
	// Note: the encoder returned by this method will be nil until the first
	// value was decoded from the stream.
	Encoder(EncoderConfig) StreamEncoder
}

// NewDecoder returns a new decoder configured with config.
func NewDecoder(config DecoderConfig) Decoder {
	config = setDecoderConfigDefaults(config)
	return &decoder{
		r: config.Input,
		p: config.Parser,
		t: config.Tag,
	}
}

// NewStreamDecoder returns a new stream decoder configured with config.
func NewStreamDecoder(config DecoderConfig) StreamDecoder {
	config = setDecoderConfigDefaults(config)
	return &streamDecoder{
		decoder: decoder{
			r: config.Input,
			p: config.Parser,
			t: config.Tag,
		},
	}
}

func setDecoderConfigDefaults(config DecoderConfig) DecoderConfig {
	if config.Input == nil {
		panic("objconv.NewDecoder: config.Input is nil")
	}

	if config.Parser == nil {
		panic("objconv.NewDecoder: config.Parser is nil")
	}

	if len(config.Tag) == 0 {
		config.Tag = "objconv"
	}

	return config
}

type decoder struct {
	r io.Reader
	p Parser
	t string
}

func (d *decoder) Decode(v interface{}) (err error) {
	defer func() { err = convertPanicToError(recover()) }()
	d.decode(NewReader(d.r), v)
	return
}

func (d *decoder) parse(r *Reader, v interface{}) (interface{}, reflect.Value) {
	to := reflect.ValueOf(v).Elem()
	return d.p.Parse(r, to.Interface()), to
}

func (d *decoder) decode(r *Reader, v interface{}) {
	from, to := d.parse(r, v)
	d.decodeValue(r, from, to)
}

func (d *decoder) decodeValue(r *Reader, v interface{}, to reflect.Value) {
	switch x := v.(type) {
	case bool:
		d.decodeBool(x, to)

	case int64:
		d.decodeInt(x, to)

	case uint64:
		d.decodeUint(x, to)

	case float64:
		d.decodeFloat(x, to)

	case string:
		d.decodeString(x, to)

	case []byte:
		d.decodeBytes(x, to)

	case time.Time:
		d.decodeTime(x, to)

	case time.Duration:
		d.decodeDuration(x, to)

	case error:
		d.decodeError(x, to)

	case ArrayParser:
		d.decodeArray(r, x, to)

	case MapParser:
		d.decodeMap(r, x, to)

	default:
		if x == nil {
			d.decodeNil(to)
		} else {
			panic(fmt.Sprintf("the parser produced an unsupported value of type %T, this is a bug", x))
		}
	}
}

func (d *decoder) decodeNil(to reflect.Value) {
	to.Set(reflect.Zero(to.Type()))
}

func (d *decoder) decodeBool(v bool, to reflect.Value) {
	to.SetBool(v)
}

func (d *decoder) decodeInt(v int64, to reflect.Value) {
	switch to.Kind() {
	case reflect.Int:
		to.SetInt(int64(convertInt64ToInt(v)))

	case reflect.Int8:
		to.SetInt(int64(convertInt64ToInt8(v)))

	case reflect.Int16:
		to.SetInt(int64(convertInt64ToInt16(v)))

	case reflect.Int32:
		to.SetInt(int64(convertInt64ToInt32(v)))

	case reflect.Int64:
		to.SetInt(v)

	case reflect.Uint:
		to.SetUint(uint64(convertInt64ToUint(v)))

	case reflect.Uint8:
		to.SetUint(uint64(convertInt64ToUint8(v)))

	case reflect.Uint16:
		to.SetUint(uint64(convertInt64ToUint16(v)))

	case reflect.Uint32:
		to.SetUint(uint64(convertInt64ToUint32(v)))

	case reflect.Uint64:
		to.SetUint(convertInt64ToUint64(v))

	case reflect.Uintptr:
		to.SetUint(uint64(convertInt64ToUintptr(v)))

	case reflect.Float32:
		to.SetFloat(float64(convertInt64ToFloat32(v)))

	case reflect.Float64:
		to.SetFloat(convertInt64ToFloat64(v))

	default:
		to.Set(reflect.ValueOf(v))
	}
}

func (d *decoder) decodeUint(v uint64, to reflect.Value) {
	switch to.Kind() {
	case reflect.Int:
		to.SetInt(int64(convertUint64ToInt(v)))

	case reflect.Int8:
		to.SetInt(int64(convertUint64ToInt8(v)))

	case reflect.Int16:
		to.SetInt(int64(convertUint64ToInt16(v)))

	case reflect.Int32:
		to.SetInt(int64(convertUint64ToInt32(v)))

	case reflect.Int64:
		to.SetInt(convertUint64ToInt64(v))

	case reflect.Uint:
		to.SetUint(uint64(convertUint64ToUint(v)))

	case reflect.Uint8:
		to.SetUint(uint64(convertUint64ToUint8(v)))

	case reflect.Uint16:
		to.SetUint(uint64(convertUint64ToUint16(v)))

	case reflect.Uint32:
		to.SetUint(uint64(convertUint64ToUint32(v)))

	case reflect.Uint64:
		to.SetUint(v)

	case reflect.Uintptr:
		to.SetUint(uint64(convertUint64ToUintptr(v)))

	case reflect.Float32:
		to.SetFloat(float64(convertUint64ToFloat32(v)))

	case reflect.Float64:
		to.SetFloat(convertUint64ToFloat64(v))

	default:
		to.Set(reflect.ValueOf(v))
	}
}

func (d *decoder) decodeFloat(v float64, to reflect.Value) {
	switch to.Kind() {
	case reflect.Float32, reflect.Float64:
		to.SetFloat(v)

	default:
		to.Set(reflect.ValueOf(v))
	}
}

func (d *decoder) decodeString(v string, to reflect.Value) {
	switch to.Kind() {
	case reflect.Slice:
		d.decodeStringToSlice(v, to)

	case reflect.String:
		to.SetString(v)

	default:
		to.Set(reflect.ValueOf(v))
	}
}

func (d *decoder) decodeStringToSlice(v string, to reflect.Value) {
	switch to.Type().Elem().Kind() {
	case reflect.Uint8: // []byte
		to.SetBytes([]byte(v))

	case reflect.Int32: // []rune
		to.Set(reflect.ValueOf([]rune(string(v))))

	default:
		to.SetString(v)
	}
}

func (d *decoder) decodeBytes(v []byte, to reflect.Value) {
	switch to.Kind() {
	case reflect.Slice:
		d.decodeBytesToSlice(v, to)

	case reflect.String:
		to.SetString(string(v))

	default:
		to.Set(reflect.ValueOf(v))
	}
}

func (d *decoder) decodeBytesToSlice(v []byte, to reflect.Value) {
	switch to.Type().Elem().Kind() {
	case reflect.Int32: // []rune
		to.Set(reflect.ValueOf([]rune(string(v))))

	default:
		to.SetBytes(v)
	}
}

func (d *decoder) decodeDuration(v time.Duration, to reflect.Value) { to.Set(reflect.ValueOf(v)) }

func (d *decoder) decodeTime(v time.Time, to reflect.Value) { to.Set(reflect.ValueOf(v)) }

func (d *decoder) decodeError(v error, to reflect.Value) { to.Set(reflect.ValueOf(v)) }

func (d *decoder) decodeArray(r *Reader, v ArrayParser, to reflect.Value) {
	t := to.Type()

	switch t.Kind() {
	case reflect.Slice:
	default:
		t = reflect.TypeOf(([]interface{})(nil))
	}

	n := v.Len()
	if n < 0 {
		n = 20
	}

	s := reflect.MakeSlice(t, 0, n)
	z := reflect.Zero(t.Elem())
	h := z.Interface()

	for i := 0; true; i++ {
		if x, ok := v.Parse(r, h); !ok {
			break
		} else {
			s = reflect.Append(s, z)
			d.decodeValue(r, x, s.Index(i))
		}
	}

	to.Set(s)
}

func (d *decoder) decodeMap(r *Reader, v MapParser, to reflect.Value) {
	if t := to.Type(); t.Kind() == reflect.Struct {
		d.decodeMapToStruct(r, v, to, t)
	} else {
		d.decodeMapToMap(r, v, to, t)
	}
}

func (d *decoder) decodeMapToMap(r *Reader, v MapParser, to reflect.Value, t reflect.Type) {
	m := reflect.MakeMap(t)
	kt := t.Key()
	vt := t.Elem()

	for {
		key := reflect.New(kt).Elem()
		if x, ok := v.ParseKey(r, key.Interface()); !ok {
			break
		} else {
			d.decodeValue(r, x, key)
		}

		val := reflect.New(vt).Elem()
		d.decodeValue(r, v.ParseValue(r, val.Interface()), val)
		m.SetMapIndex(key, val)
	}

	to.Set(m)
}

func (d *decoder) decodeMapToStruct(r *Reader, v MapParser, to reflect.Value, t reflect.Type) {
	s := LookupStruct(t).SetterValue(d.t, to)

	for {
		var f string

		if x, ok := v.ParseKey(r, f); !ok {
			break
		} else {
			d.decodeValue(r, x, reflect.ValueOf(&f).Elem())
		}

		if fv, ok := s[f]; ok {
			d.decodeValue(r, v.ParseValue(r, fv.Interface()), fv)
		}
	}
}

type streamDecoder struct {
	decoder
	reader *Reader
	parser ArrayParser
	err    error
	count  int
	array  bool
}

func (d *streamDecoder) Decode(v interface{}) (err error) {
	if err = d.err; err != nil {
		return
	}

	if d.reader == nil {
		d.reader = NewReader(d.r)
	}

	defer func() { err = d.convertPanicToError(recover()) }()

	if d.parser == nil {
		from, _ := d.parse(d.reader, v)
		switch x := from.(type) {
		case ArrayParser:
			d.parser = x
			d.array = true
		default:
			d.parser = ArrayParserLen(1, ArrayParserFunc(func(r *Reader, hint interface{}) (interface{}, bool) {
				return x, true
			}))
		}
	}

	if x, ok := d.parser.Parse(d.reader, v); !ok {
		panic(io.EOF)
	} else {
		d.decodeValue(d.reader, x, reflect.ValueOf(v).Elem())
	}

	d.count++
	return
}

func (d *streamDecoder) Len() int {
	if d.parser == nil {
		return -1
	}
	n := d.parser.Len()
	if n > 0 {
		n -= d.count
	}
	return n
}

func (d *streamDecoder) Error() error {
	return d.err
}

func (d *streamDecoder) Encoder(config EncoderConfig) StreamEncoder {
	if d.parser == nil {
		return nil
	}

	if d.array {
		return NewStreamEncoder(config)
	}

	config = setEncoderConfigDefault(config)
	return &nonstreamEncoder{
		encoder: encoder{
			w: config.Output,
			e: config.Emitter,
			t: config.Tag,
		},
	}
}

func (d *streamDecoder) convertPanicToError(v interface{}) (err error) {
	if err = convertPanicToError(v); err != nil {
		d.err = err
	}
	return
}
