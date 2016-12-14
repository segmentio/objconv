package objconv

import (
	"errors"
	"reflect"
	"time"
)

// DecoderConfig carries the configuration for instantiating decoders.
type DecoderConfig struct {
	// Parser defines the format used by the decoder.
	Parser Parser

	// TimeFormat is used to parse time values from strings, defaults to
	// time.RFC33339Nano.
	TimeFormat string

	// TimeLocation is the location used for parsing time values from strings if
	// none is explicitly set, defaults to time.Local.
	TimeLocation *time.Location
}

// A Decoder implements the algorithms for building data structures from their
// serialized forms.
//
// Decoders are not safe for use by multiple goroutines.
type Decoder struct {
	p   Parser
	key bool

	timeFormat   string
	timeLocation *time.Location
}

// NewDecoder returns a new decoder object that uses parser to deserialize data
// structures.
//
// The function panics if parser is nil.
func NewDecoder(parser Parser) *Decoder {
	return NewDecoderWith(DecoderConfig{
		Parser: parser,
	})
}

// NewDecoderWith returns a new decoder configured with config.
//
// The function panics if config.Parser is nil.
func NewDecoderWith(config DecoderConfig) *Decoder {
	if config.Parser == nil {
		panic("objconv.NewDecoder: the parser is nil")
	}
	if len(config.TimeFormat) == 0 {
		config.TimeFormat = time.RFC3339Nano
	}
	if config.TimeLocation == nil {
		config.TimeLocation = time.Local
	}
	return &Decoder{
		p:            config.Parser,
		timeFormat:   config.TimeFormat,
		timeLocation: config.TimeLocation,
	}
}

// Decode expects v to be a pointer to a value in which the decoder will load
// the next parsed data.
//
// The method panics if v is neither a pointer type nor implements the
// ValueDecoder interface, or if v is a nil pointer.
func (d *Decoder) Decode(v interface{}) error {
	to := reflect.ValueOf(v)

	switch {
	case to.Kind() != reflect.Ptr:
		panic("objconv.(*Decoder).Decode: v must be a pointer")

	case to.IsNil():
		panic("objconv.(*Decoder).Decode: v cannot be a nil pointer")
	}

	if err := d.decodeMapValueMaybe(); err != nil {
		return err
	}

	_, err := d.decodeValue(to)
	return err
}

func (d *Decoder) decodeValue(to reflect.Value) (Type, error) {
	return decodeFuncOf(to.Type())(d, to)
}

func (d *Decoder) decodeValueBool(to reflect.Value) (t Type, err error) {
	var v bool

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()

	case Bool:
		v, err = d.decodeBool()

	default:
		err = &TypeConversionError{
			From: t,
			To:   Bool,
		}
	}

	if err != nil {
		return
	}

	to.SetBool(v)
	return
}

func (d *Decoder) decodeValueInt(to reflect.Value) (t Type, err error) {
	var i int64
	var u uint64

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()

	case Int:
		if i, err = d.decodeInt(); err != nil {
			return
		}

		switch t := to.Type(); t.Kind() {
		case reflect.Int:
			err = checkInt64Bounds(i, int64(IntMin), uint64(IntMax), t)
		case reflect.Int8:
			err = checkInt64Bounds(i, Int8Min, Int8Max, t)
		case reflect.Int16:
			err = checkInt64Bounds(i, Int16Min, Int16Max, t)
		case reflect.Int32:
			err = checkInt64Bounds(i, Int32Min, Int32Max, t)
		}

	case Uint:
		if u, err = d.decodeUint(); err != nil {
			return
		}

		switch t := to.Type(); t.Kind() {
		case reflect.Int:
			err = checkUint64Bounds(u, uint64(IntMax), t)
		case reflect.Int8:
			err = checkUint64Bounds(u, Int8Max, t)
		case reflect.Int16:
			err = checkUint64Bounds(u, Int16Max, t)
		case reflect.Int32:
			err = checkUint64Bounds(u, Int32Max, t)
		case reflect.Int64:
			err = checkUint64Bounds(u, Int64Max, t)
		}

		i = int64(u)

	default:
		err = &TypeConversionError{
			From: t,
			To:   Int,
		}
	}

	if err != nil {
		return
	}

	to.SetInt(i)
	return
}

func (d *Decoder) decodeValueUint(to reflect.Value) (t Type, err error) {
	var i int64
	var u uint64

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()

	case Int:
		if i, err = d.decodeInt(); err != nil {
			return
		}

		switch t := to.Type(); t.Kind() {
		case reflect.Uint:
			err = checkInt64Bounds(i, 0, uint64(UintMax), t)
		case reflect.Uint8:
			err = checkInt64Bounds(i, 0, Uint8Max, t)
		case reflect.Uint16:
			err = checkInt64Bounds(i, 0, Uint16Max, t)
		case reflect.Uint32:
			err = checkInt64Bounds(i, 0, Uint32Max, t)
		case reflect.Uint64:
			err = checkInt64Bounds(i, 0, Uint64Max, t)
		}

		u = uint64(i)

	case Uint:
		if u, err = d.decodeUint(); err != nil {
			return
		}

		switch t := to.Type(); t.Kind() {
		case reflect.Uint:
			err = checkUint64Bounds(u, uint64(IntMax), t)
		case reflect.Uint8:
			err = checkUint64Bounds(u, Int8Max, t)
		case reflect.Uint16:
			err = checkUint64Bounds(u, Int16Max, t)
		case reflect.Uint32:
			err = checkUint64Bounds(u, Int32Max, t)
		}

	default:
		err = &TypeConversionError{
			From: t,
			To:   Uint,
		}
	}

	if err != nil {
		return
	}

	to.SetUint(u)
	return
}

func (d *Decoder) decodeValueFloat(to reflect.Value) (t Type, err error) {
	var i int64
	var u uint64
	var f float64

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()

	case Int:
		i, err = d.decodeInt()
		f = float64(i)

	case Uint:
		u, err = d.decodeUint()
		f = float64(u)

	case Float:
		f, err = d.decodeFloat()

	default:
		err = &TypeConversionError{
			From: t,
			To:   Float,
		}
	}

	if err != nil {
		return
	}

	to.SetFloat(f)
	return
}

func (d *Decoder) decodeValueString(to reflect.Value) (t Type, err error) {
	var b []byte

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()

	case String:
		b, err = d.decodeString()

	case Bytes:
		b, err = d.decodeBytes()

	default:
		err = &TypeConversionError{
			From: t,
			To:   String,
		}
	}

	if err != nil {
		return
	}

	to.SetString(string(b))
	return
}

func (d *Decoder) decodeValueBytes(to reflect.Value) (t Type, err error) {
	var b []byte
	var v []byte

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()

	case String:
		b, err = d.decodeString()

	case Bytes:
		b, err = d.decodeBytes()

	default:
		err = &TypeConversionError{
			From: t,
			To:   String,
		}
	}

	if err != nil {
		return
	}

	if n := len(b); n != 0 {
		v = make([]byte, n)
		copy(v, b)
	}

	to.SetBytes(v)
	return
}

func (d *Decoder) decodeValueTime(to reflect.Value) (t Type, err error) {
	var s []byte
	var v time.Time

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()

	case String:
		s, err = d.decodeString()

	case Bytes:
		s, err = d.decodeBytes()

	case Time:
		v, err = d.decodeTime()
	}

	if err != nil {
		return
	}

	if t == String || t == Bytes {
		v, err = time.ParseInLocation(d.timeFormat, string(s), d.timeLocation)
	}

	to.Set(reflect.ValueOf(v))
	return
}

func (d *Decoder) decodeValueDuration(to reflect.Value) (t Type, err error) {
	var s []byte
	var v time.Duration

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()

	case String:
		s, err = d.decodeString()

	case Bytes:
		s, err = d.decodeBytes()

	case Duration:
		v, err = d.decodeDuration()
	}

	if err != nil {
		return
	}

	if t == String || t == Bytes {
		v, err = time.ParseDuration(string(s))
	}

	to.Set(reflect.ValueOf(v))
	return
}

func (d *Decoder) decodeValueError(to reflect.Value) (t Type, err error) {
	var s []byte
	var v error

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()

	case String:
		s, err = d.decodeString()

	case Bytes:
		s, err = d.decodeBytes()

	case Error:
		v, err = d.decodeError()
	}

	if err != nil {
		return
	}

	if t == String || t == Bytes {
		v = errors.New(string(s))
	}

	to.Set(reflect.ValueOf(v))
	return
}

func (d *Decoder) decodeValueSlice(to reflect.Value) (typ Type, err error) {
	t := to.Type()                   // []T
	e := t.Elem()                    // T
	z := reflect.Zero(e)             // T{}
	v := reflect.New(e).Elem()       // &T{}
	s := reflect.MakeSlice(t, 0, 50) // make([]T, 0, 50)
	f := decodeFuncOf(e)

	if typ, err = d.DecodeArray(func(d *Decoder) (err error) {
		v.Set(z) // reset to the zero-value
		if _, err = f(d, v); err != nil {
			return
		}
		s = reflect.Append(s, v)
		return
	}); err != nil {
		return
	}

	if typ != Nil {
		to.Set(s)
	} else {
		to.Set(reflect.Zero(t))
	}

	return
}

func (d *Decoder) decodeValueArray(to reflect.Value) (typ Type, err error) {
	n := to.Len()        // len(to)
	t := to.Type()       // [...]T
	e := t.Elem()        // T
	z := reflect.Zero(e) // T{}
	f := decodeFuncOf(e)

	for i := 0; i != n; i++ {
		to.Index(i).Set(z) // reset to the zero-value
	}

	i := 0

	if typ, err = d.DecodeArray(func(d *Decoder) (err error) {
		if i < n {
			if _, err = f(d, to.Index(i)); err != nil {
				return
			}
		}
		i++
		return
	}); err != nil {
		return
	}

	if (typ != Nil) && (i != n) {
		typ, err = Nil, &ArrayLengthError{
			Type:   t,
			Length: n,
		}
	}

	return
}

func (d *Decoder) decodeValueMap(to reflect.Value) (typ Type, err error) {
	t := to.Type()          // map[K]V
	m := reflect.MakeMap(t) // make(map[K]V)

	kt := t.Key()                // K
	kz := reflect.Zero(kt)       // K{}
	kv := reflect.New(kt).Elem() // &K{}
	kf := decodeFuncOf(kt)

	vt := t.Elem()               // V
	vz := reflect.Zero(vt)       // V{}
	vv := reflect.New(vt).Elem() // &V{}
	vf := decodeFuncOf(vt)

	if typ, err = d.DecodeMap(func(d *Decoder) (err error) {
		kv.Set(kz) // reset the key to its zero-value
		vv.Set(vz) // reset the value to its zero-value
		if _, err = kf(d, kv); err != nil {
			return
		}
		if err = d.decodeMapValue(); err != nil {
			return
		}
		if _, err = vf(d, vv); err != nil {
			return
		}
		m.SetMapIndex(kv, vv)
		return
	}); err != nil {
		return
	}

	if typ != Nil {
		to.Set(m)
	} else {
		to.Set(reflect.Zero(t))
	}

	return
}

func (d *Decoder) decodeValueStruct(to reflect.Value) (typ Type, err error) {
	t := to.Type()
	s := LookupStruct(t)

	if typ, err = d.DecodeMap(func(d *Decoder) (err error) {
		var b []byte

		if b, err = d.decodeString(); err != nil {
			return
		}
		if err = d.decodeMapValue(); err != nil {
			return
		}

		f, ok := s.FieldByName[string(b)]
		if !ok {
			var v interface{} // discard
			return d.Decode(&v)
		}

		_, err = f.decode(d, to.FieldByIndex(f.Index))
		return
	}); err != nil {
		to.Set(reflect.Zero(t))
	}

	return
}

func (d *Decoder) decodeValuePointer(to reflect.Value) (typ Type, err error) {
	var t = to.Type()
	var v reflect.Value

	if to.IsNil() {
		v = reflect.New(t.Elem())
	} else {
		v = to
	}

	if typ, err = d.decodeValue(v.Elem()); err != nil {
		return
	}

	if to.CanSet() {
		switch {
		case typ == Nil:
			to.Set(reflect.Zero(t))

		case to.IsNil():
			to.Set(v)
		}
	}

	return
}

func (d *Decoder) decodeValueDecoder(to reflect.Value) (typ Type, err error) {
	typ = Bool // just needs to not be Nil
	err = to.Interface().(ValueDecoder).DecodeValue(d)
	return
}

func (d *Decoder) decodeValueUnsupported(to reflect.Value) (Type, error) {
	return Nil, &UnsupportedTypeError{Type: to.Type()}
}

func (d *Decoder) decodeType() (Type, error) { return d.p.ParseType() }

func (d *Decoder) decodeNil() error { return d.p.ParseNil() }

func (d *Decoder) decodeBool() (bool, error) { return d.p.ParseBool() }

func (d *Decoder) decodeInt() (int64, error) { return d.p.ParseInt() }

func (d *Decoder) decodeUint() (uint64, error) { return d.p.ParseUint() }

func (d *Decoder) decodeFloat() (float64, error) { return d.p.ParseFloat() }

func (d *Decoder) decodeString() ([]byte, error) { return d.p.ParseString() }

func (d *Decoder) decodeBytes() ([]byte, error) { return d.p.ParseBytes() }

func (d *Decoder) decodeTime() (time.Time, error) { return d.p.ParseTime() }

func (d *Decoder) decodeDuration() (time.Duration, error) { return d.p.ParseDuration() }

func (d *Decoder) decodeError() (error, error) { return d.p.ParseError() }

// DecodeArray provides the implementation of the algorithm for decoding arrays,
// where f is called to decode each element of the array.
//
// The method returns the underlying type of the value returned by the parser.
func (d *Decoder) DecodeArray(f func(*Decoder) error) (t Type, err error) {
	var n int

	if err = d.decodeMapValueMaybe(); err != nil {
		return
	}

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()
		return

	case Array:
		n, err = d.decodeArrayBegin()

	default:
		err = &TypeConversionError{
			From: t,
			To:   Array,
		}
	}

	if err != nil {
		return
	}

decodeArray:
	for i := 0; n < 0 || i < n; i++ {
		if i != 0 {
			switch err = d.decodeArrayNext(); err {
			case nil:
			case End:
				break decodeArray
			default:
				return
			}
		}
		if err = f(d); err != nil {
			return
		}
	}

	err = d.decodeArrayEnd()
	return
}

func (d *Decoder) decodeArrayBegin() (int, error) { return d.p.ParseArrayBegin() }

func (d *Decoder) decodeArrayEnd() error { return d.p.ParseArrayEnd() }

func (d *Decoder) decodeArrayNext() error { return d.p.ParseArrayNext() }

// DecodeMap provides the implementation of the algorithm for decoding maps,
// where f is called to decode each pair of key and value.
//
// The function f is expected to decode two values from the map, the first one
// being the key and the second the associated value.
//
// The method returns the underlying type of the value returned by the parser.
func (d *Decoder) DecodeMap(f func(*Decoder) error) (t Type, err error) {
	var n int

	if err = d.decodeMapValueMaybe(); err != nil {
		return
	}

	if t, err = d.decodeType(); err != nil {
		return
	}

	switch t {
	case Nil:
		err = d.decodeNil()
		return

	case Map:
		n, err = d.decodeMapBegin()

	default:
		err = &TypeConversionError{
			From: t,
			To:   Map,
		}
	}

	if err != nil {
		return
	}

decodeMap:
	for i := 0; n < 0 || i < n; i++ {
		if i != 0 {
			switch err = d.decodeMapNext(); err {
			case nil:
			case End:
				break decodeMap
			default:
				return
			}
		}

		d.key = true
		err = f(d)
		// Because internal calls don't use the exported methods they may not
		// reset this flag to false when expected, forcing the value here.
		d.key = false

		if err != nil {
			return
		}
	}

	err = d.decodeMapEnd()
	return
}

func (d *Decoder) decodeMapBegin() (int, error) { return d.p.ParseMapBegin() }

func (d *Decoder) decodeMapEnd() error { return d.p.ParseMapEnd() }

func (d *Decoder) decodeMapValue() error { return d.p.ParseMapValue() }

func (d *Decoder) decodeMapNext() error { return d.p.ParseMapNext() }

func (d *Decoder) decodeMapValueMaybe() (err error) {
	if d.key {
		d.key = false
		err = d.decodeMapValue()
	}
	return
}

func decodeFuncOf(t reflect.Type) func(*Decoder, reflect.Value) (Type, error) {
	switch {
	case t == timeType:
		return (*Decoder).decodeValueTime

	case t == durationType:
		return (*Decoder).decodeValueDuration

	case t.Implements(valueDecoderInterface):
		return (*Decoder).decodeValueDecoder

	case t.Implements(errorInterface):
		return (*Decoder).decodeValueError
	}

	switch t.Kind() {
	case reflect.Bool:
		return (*Decoder).decodeValueBool

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return (*Decoder).decodeValueInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return (*Decoder).decodeValueUint

	case reflect.Float32, reflect.Float64:
		return (*Decoder).decodeValueFloat

	case reflect.String:
		return (*Decoder).decodeValueString

	case reflect.Slice:
		switch {
		case t.Elem().Kind() == reflect.Uint8:
			return (*Decoder).decodeValueBytes
		default:
			return (*Decoder).decodeValueSlice
		}

	case reflect.Array:
		return (*Decoder).decodeValueArray

	case reflect.Map:
		return (*Decoder).decodeValueMap

	case reflect.Struct:
		return (*Decoder).decodeValueStruct

	case reflect.Ptr:
		return (*Decoder).decodeValuePointer

	default:
		return (*Decoder).decodeValueUnsupported
	}
}

// ValueDecoder is the interface that can be implemented by types that wish to
// provide their own decoding algorithms.
//
// The DecodeValue method is called when the value is found by a decoding
// algorithm.
type ValueDecoder interface {
	DecodeValue(*Decoder) error
}

// ValueDecoderFunc allos the use of regular functions or methods as value
// decoders.
type ValueDecoderFunc func(*Decoder) error

// DecodeValue calls f(d).
func (f ValueDecoderFunc) DecodeValue(d *Decoder) error { return f(d) }
