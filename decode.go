package objconv

import (
	"bytes"
	"errors"
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
	var parser Parser

	if parser, err = GetParser(format); err == nil {
		err = NewDecoder(DecoderConfig{
			Input:  in,
			Parser: parser,
		}).Decode(value)
	}

	return
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
	d := makeDecoder(config.Parser, config.Input)
	return &d
}

// NewStreamDecoder returns a new stream decoder configured with config.
func NewStreamDecoder(config DecoderConfig) StreamDecoder {
	config = setDecoderConfigDefaults(config)
	return &streamDecoder{decoder: makeDecoder(config.Parser, config.Input)}
}

func setDecoderConfigDefaults(config DecoderConfig) DecoderConfig {
	if config.Input == nil {
		panic("objconv.NewDecoder: config.Input is nil")
	}

	if config.Parser == nil {
		panic("objconv.NewDecoder: config.Parser is nil")
	}

	return config
}

type decoder struct {
	p Parser
	r *Reader
	Reader
}

func makeDecoder(p Parser, r io.Reader) (d decoder) {
	d.p = p
	d.r = &d.Reader

	// Use the reader directly if it's already an instance of Reader.
	switch x := r.(type) {
	case *Reader:
		d.r = x
	default:
		d.r.r = r
	}

	return
}

func (d *decoder) Decode(v interface{}) (err error) {
	defer func() { err = convertPanicToError(recover()) }()
	d.decode(d.r, v)
	return
}

func (d *decoder) parse(r *Reader, v interface{}) (interface{}, reflect.Value) {
	to := reflect.ValueOf(v).Elem()
	v, err := d.p.Parse(r, to.Interface())
	if err != nil {
		panic(err)
	}
	return v, to
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
	switch to.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map:
		if !to.IsNil() {
			to.Set(reflect.Zero(to.Type()))
		}

	case reflect.Bool:
		to.SetBool(false)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		to.SetInt(0)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		to.SetUint(0)

	case reflect.Float32, reflect.Float64:
		to.SetFloat(0)

	case reflect.String:
		to.SetString("")

	default:
		to.Set(reflect.Zero(to.Type()))
	}
}

func (d *decoder) decodeBool(v bool, to reflect.Value) {
	switch to.Kind() {
	case reflect.Bool:
		to.SetBool(v)

	default:
		to.Set(reflect.ValueOf(v))
	}
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
	case reflect.String:
		to.SetString(v)

	case reflect.Slice:
		d.decodeStringToSlice(v, to)

	default:
		d.decodeStringToOther(v, to)
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

func (d *decoder) decodeStringToOther(v string, to reflect.Value) {
	switch to.Interface().(type) {
	case time.Time:
		if t, err := time.Parse(time.RFC3339Nano, v); err != nil {
			panic(err)
		} else {
			to.Set(reflect.ValueOf(t))
		}

	case time.Duration:
		if d, err := time.ParseDuration(v); err != nil {
			panic(err)
		} else {
			to.Set(reflect.ValueOf(d))
		}

	case error:
		to.Set(reflect.ValueOf(errors.New(v)))

	default:
		to.Set(reflect.ValueOf(v))
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
	switch t := to.Type(); t.Kind() {
	case reflect.Slice:
		d.decodeArrayToSlice(r, v, to, t)
	default:
		d.decodeArrayToInterface(r, v, to)
	}
}

func (d *decoder) decodeArrayToSlice(r *Reader, v ArrayParser, to reflect.Value, t reflect.Type) {
	n := v.Len()

	if n < 0 {
		n = 20
	}

	s := reflect.MakeSlice(t, 0, n)
	z := reflect.Zero(t.Elem())
	h := z.Interface()

decodeLoop:
	for i := 0; true; i++ {
		switch x, err := v.Parse(r, h); err {
		case nil:
			s = reflect.Append(s, z)
			d.decodeValue(r, x, s.Index(i))
		case io.EOF:
			break decodeLoop
		default:
			panic(err)
		}
	}

	to.Set(s)
}

func (d *decoder) decodeArrayToInterface(r *Reader, v ArrayParser, to reflect.Value) {
	n := v.Len()

	if n < 0 {
		n = 20
	}

	s := make([]interface{}, 0, n)

decodeLoop:
	for i := 0; true; i++ {
		switch x, err := v.Parse(r, nil); err {
		case nil:
			switch s = append(s, nil); x.(type) {
			case ArrayParser, MapParser:
				d.decodeValue(r, x, reflect.ValueOf(&s[i]).Elem())
			default:
				s[i] = x
			}
		case io.EOF:
			break decodeLoop
		default:
			panic(err)
		}
	}

	to.Set(reflect.ValueOf(s))
}

func (d *decoder) decodeMap(r *Reader, v MapParser, to reflect.Value) {
	switch t := to.Type(); t.Kind() {
	case reflect.Struct:
		d.decodeMapToStruct(r, v, to, t)

	case reflect.Map:
		d.decodeMapToMap(r, v, to, t)

	default:
		var m map[interface{}]interface{}
		var mv = reflect.ValueOf(&m).Elem()
		d.decodeMapToMap(r, v, mv, reflect.TypeOf(m))
		to.Set(mv)
	}
}

func (d *decoder) decodeMapToMap(r *Reader, v MapParser, to reflect.Value, t reflect.Type) {
	m := reflect.MakeMap(t)

	kt := t.Key()
	vt := t.Elem()

	ke := reflect.New(kt).Elem()
	ve := reflect.New(vt).Elem()

	ki := ke.Interface()
	vi := ve.Interface()

decodeLoop:
	for {
		switch x, err := v.ParseKey(r, ki); err {
		case nil:
			d.decodeValue(r, x, ke)
		case io.EOF:
			break decodeLoop
		default:
			panic(err)
		}

		switch x, err := v.ParseValue(r, vi); err {
		case nil:
			d.decodeValue(r, x, ve)
		case io.EOF:
			panic(io.ErrUnexpectedEOF)
		default:
			panic(err)
		}

		m.SetMapIndex(ke, ve)
	}

	to.Set(m)
}

func (d *decoder) decodeMapToStruct(r *Reader, v MapParser, to reflect.Value, t reflect.Type) {
	s := LookupStruct(t)
	f := ""

	fv := reflect.ValueOf(&f).Elem()
	fi := fv.Interface()

decodeLoop:
	for {
		switch x, err := v.ParseKey(r, fi); err {
		case nil:
			d.decodeValue(r, x, fv)
		case io.EOF:
			break decodeLoop
		default:
			panic(err)
		}

		if sf := s.FieldByName[f]; sf != nil {
			fv := to.FieldByIndex(sf.Index)
			switch x, err := v.ParseValue(r, fv.Interface()); err {
			case nil:
				d.decodeValue(r, x, fv)
			default:
				panic(err)
			}
		}

		f = ""
	}
}

type streamDecoder struct {
	decoder
	parser ArrayParser
	err    error
	count  int
	array  bool
}

func (d *streamDecoder) Decode(v interface{}) (err error) {
	if err = d.err; err != nil {
		return
	}

	defer func() { err = d.convertPanicToError(recover()) }()

	ve := reflect.ValueOf(v).Elem()
	vi := ve.Interface()

	if d.parser == nil {
		from, _ := d.parse(d.r, v)
		switch a := from.(type) {
		case ArrayParser:
			d.parser = a
			d.array = true
		default:
			d.parser = ArrayParserLen(1, ArrayParserFunc(func(r *Reader, hint interface{}) (interface{}, error) {
				return a, nil
			}))
		}
	}

	if x, err := d.parser.Parse(d.r, vi); err != nil {
		panic(err)
	} else {
		d.decodeValue(d.r, x, ve)
	}

	d.count++
	return
}

func (d *streamDecoder) Len() (n int) {
	defer func() { d.convertPanicToError(recover()) }()
	n = -1

	if d.parser == nil {
		// The length was requested but we have no idea what value the program
		// will try to decode from the stream so calling d.parse will likely
		// produce a value of a mismatching type.
		// We need to be able to roll back to the stream after reading what the
		// length was.
		// To achieve this behavior in a seamless way we use an intermediary
		// buffer that records everything the parser consumes, then rebuild the
		// base reader to include the recorded bytes.
		// This is kind of tricky but it works like a charm!
		var z interface{}

		b := &bytes.Buffer{}
		r := d.r.r

		d.r.r = io.TeeReader(r, b)

		v, _ := d.parse(d.r, &z)
		switch a := v.(type) {
		case ArrayParser:
			n = a.Len()
		default:
			n = 1
		}

		d.r.r = io.MultiReader(b, r)
		d.r.Reset()
		return
	}

	if n = d.parser.Len(); n > 0 {
		n -= d.count
	}

	return
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
			sort: config.SortMapKeys,
			e:    config.Emitter,
			w:    Writer{w: config.Output},
		},
	}
}

func (d *streamDecoder) convertPanicToError(v interface{}) (err error) {
	if err = convertPanicToError(v); err != nil {
		d.err = err
	}
	return
}
