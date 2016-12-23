package objconv

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"time"
)

// A Decoder implements the algorithms for building data structures from their
// serialized forms.
//
// Decoders are not safe for use by multiple goroutines.
type Decoder struct {
	// Parser to use to load values.
	Parser Parser

	// DecodeMap can be set to a function used to decode maps when there is no
	// destination type (like when decoding to an empty interface for example).
	DecodeMapFunc func(Decoder, Decoder) error

	off int // offset of the value when decoding a map
}

// NewDecoder returns a decoder object that uses p, will panic if p is nil.
func NewDecoder(p Parser) *Decoder {
	if p == nil {
		panic("objconv: the parser is nil")
	}
	return &Decoder{Parser: p}
}

// Decode expects v to be a pointer to a value in which the decoder will load
// the next parsed data.
//
// The method panics if v is neither a pointer type nor implements the
// ValueDecoder interface, or if v is a nil pointer.
func (d Decoder) Decode(v interface{}) (err error) {
	to := reflect.ValueOf(v)

	switch {
	case to.Kind() != reflect.Ptr:
		panic("objconv.Decoder.Decode: v must be a pointer")

	case to.IsNil():
		panic("objconv.Decoder.Decode: v cannot be a nil pointer")
	}

	if d.off != 0 {
		if d.off, err = 0, d.decodeMapValue(d.off-1); err != nil {
			return
		}
	}

	_, err = d.decodeValue(to.Elem())
	return
}

func (d Decoder) decodeValue(to reflect.Value) (Type, error) {
	return decodeFuncOf(to.Type())(d, to)
}

func (d Decoder) decodeValueNil(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueNilFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueNilFromType(t Type, to reflect.Value) (err error) {
	switch t {
	case Nil:
		err = d.decodeNil()
	default:
		err = typeConversionError(t, Nil)
	}

	to.Set(zeroValueOf(to.Type()))
	return
}

func (d Decoder) decodeValueBool(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueBoolFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueBoolFromType(t Type, to reflect.Value) (err error) {
	var v bool

	switch t {
	case Nil:
		err = d.decodeNil()

	case Bool:
		v, err = d.decodeBool()

	default:
		err = typeConversionError(t, Bool)
	}

	if err != nil {
		return
	}

	to.SetBool(v)
	return
}

func (d Decoder) decodeValueInt(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueIntFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueIntFromType(t Type, to reflect.Value) (err error) {
	var i int64
	var u uint64

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
		err = typeConversionError(t, Int)
	}

	if err != nil {
		return
	}

	to.SetInt(i)
	return
}

func (d Decoder) decodeValueUint(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueUintFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueUintFromType(t Type, to reflect.Value) (err error) {
	var i int64
	var u uint64

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
			err = checkUint64Bounds(u, uint64(UintMax), t)
		case reflect.Uint8:
			err = checkUint64Bounds(u, Uint8Max, t)
		case reflect.Uint16:
			err = checkUint64Bounds(u, Uint16Max, t)
		case reflect.Uint32:
			err = checkUint64Bounds(u, Uint32Max, t)
		}

	default:
		err = typeConversionError(t, Uint)
	}

	if err != nil {
		return
	}

	to.SetUint(u)
	return
}

func (d Decoder) decodeValueFloat(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueFloatFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueFloatFromType(t Type, to reflect.Value) (err error) {
	var i int64
	var u uint64
	var f float64

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
		err = typeConversionError(t, Float)
	}

	if err != nil {
		return
	}

	to.SetFloat(f)
	return
}

func (d Decoder) decodeValueString(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueStringFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueStringFromType(t Type, to reflect.Value) (err error) {
	var b []byte

	switch t {
	case Nil:
		err = d.decodeNil()

	case String:
		b, err = d.decodeString()

	case Bytes:
		b, err = d.decodeBytes()

	default:
		err = typeConversionError(t, String)
	}

	if err != nil {
		return
	}

	to.SetString(string(b))
	return
}

func (d Decoder) decodeValueBytes(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueBytesFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueBytesFromType(t Type, to reflect.Value) (err error) {
	var b []byte
	var v []byte

	switch t {
	case Nil:
		err = d.decodeNil()

	case String:
		b, err = d.decodeString()

	case Bytes:
		b, err = d.decodeBytes()

	default:
		err = typeConversionError(t, String)
	}

	if err != nil {
		return
	}

	if b != nil {
		v = make([]byte, len(b))
		copy(v, b)
	}

	to.SetBytes(v)
	return
}

func (d Decoder) decodeValueTime(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueTimeFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueTimeFromType(t Type, to reflect.Value) (err error) {
	var s []byte
	var v time.Time

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
		v, err = time.Parse(time.RFC3339Nano, string(s))
	}

	*(to.Addr().Interface().(*time.Time)) = v
	return
}

func (d Decoder) decodeValueDuration(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueDurationFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueDurationFromType(t Type, to reflect.Value) (err error) {
	var s []byte
	var v time.Duration

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

	to.SetInt(int64(v))
	return
}

func (d Decoder) decodeValueError(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueErrorFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueErrorFromType(t Type, to reflect.Value) (err error) {
	var s []byte
	var v error

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

func (d Decoder) decodeValueSlice(to reflect.Value) (t Type, err error) {
	return d.decodeValueSliceWith(to, decodeFuncOf(to.Type().Elem()))
}

func (d Decoder) decodeValueSliceWith(to reflect.Value, f decodeFunc) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueSliceFromTypeWith(t, to, f)
	}
	return
}

func (d Decoder) decodeValueSliceFromType(typ Type, to reflect.Value) (err error) {
	return d.decodeValueSliceFromTypeWith(typ, to, decodeFuncOf(to.Type().Elem()))
}

func (d Decoder) decodeValueSliceFromTypeWith(typ Type, to reflect.Value, f decodeFunc) (err error) {
	t := to.Type()

	if typ == Nil {
		to.Set(zeroValueOf(t))
		return
	}

	s := reflect.MakeSlice(t, 0, 0)
	i := 0
	n := 0

	if err = d.decodeArrayFromType(typ, func(d Decoder) (err error) {
		if i == n {
			if n *= 5; n == 0 {
				n = 10
			}
			sc := reflect.MakeSlice(t, n, n)
			reflect.Copy(sc, s)
			s = sc
		}
		if _, err = f(d, s.Index(i)); err != nil {
			return
		}
		i++
		return
	}); err != nil {
		return
	}

	if i != n {
		s = s.Slice(0, i)
	}

	to.Set(s)
	return
}

func (d Decoder) decodeValueArray(to reflect.Value) (t Type, err error) {
	return d.decodeValueArrayWith(to, decodeFuncOf(to.Type().Elem()))
}

func (d Decoder) decodeValueArrayWith(to reflect.Value, f decodeFunc) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueArrayFromTypeWith(t, to, f)
	}
	return
}

func (d Decoder) decodeValueArrayFromType(typ Type, to reflect.Value) (err error) {
	return d.decodeValueArrayFromTypeWith(typ, to, decodeFuncOf(to.Type().Elem()))
}

func (d Decoder) decodeValueArrayFromTypeWith(typ Type, to reflect.Value, f decodeFunc) (err error) {
	n := to.Len()       // len(to)
	t := to.Type()      // [...]T
	e := t.Elem()       // T
	z := zeroValueOf(e) // T{}

	for i := 0; i != n; i++ {
		to.Index(i).Set(z) // reset to the zero-value
	}

	i := 0

	if err = d.decodeArrayFromType(typ, func(d Decoder) (err error) {
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
		err = fmt.Errorf("objconv: array length mismatch, expected %d but only %d elements were decoded", n, i)
	}

	return
}

func (d Decoder) decodeValueMap(to reflect.Value) (Type, error) {
	t := to.Type()
	return d.decodeValueMapWith(to, decodeFuncOf(t.Key()), decodeFuncOf(t.Elem()))
}

func (d Decoder) decodeValueMapWith(to reflect.Value, kf decodeFunc, vf decodeFunc) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueMapFromTypeWith(t, to, kf, vf)
	}
	return
}

func (d Decoder) decodeValueMapFromType(typ Type, to reflect.Value) (err error) {
	t := to.Type()
	return d.decodeValueMapFromTypeWith(typ, to, decodeFuncOf(t.Key()), decodeFuncOf(t.Elem()))
}

func (d Decoder) decodeValueMapFromTypeWith(typ Type, to reflect.Value, kf decodeFunc, vf decodeFunc) (err error) {
	t := to.Type() // map[K]V

	if typ == Nil {
		to.Set(zeroValueOf(t))
		return
	}

	switch t {
	case mapInterfaceInterfaceType:
		return d.decodeValueMapInterfaceInterface(typ, to)

	case mapStringInterfaceType:
		return d.decodeValueMapStringInterface(typ, to)

	case mapStringStringType:
		return d.decodeValueMapStringString(typ, to)
	}

	m := reflect.MakeMap(t) // make(map[K]V)

	kt := t.Key()                // K
	kz := zeroValueOf(kt)        // K{}
	kv := reflect.New(kt).Elem() // &K{}

	vt := t.Elem()               // V
	vz := zeroValueOf(vt)        // V{}
	vv := reflect.New(vt).Elem() // &V{}

	f := d.DecodeMapFunc
	if f == nil {
		f = func(kd Decoder, vd Decoder) (err error) {
			kv.Set(kz) // reset the key to its zero-value
			vv.Set(vz) // reset the value to its zero-value
			if _, err = kf(d, kv); err != nil {
				return
			}
			if err = d.decodeMapValue(vd.off - 1); err != nil {
				return
			}
			if _, err = vf(d, vv); err != nil {
				return
			}
			m.SetMapIndex(kv, vv)
			return
		}
	}

	if err = d.decodeMapFromType(typ, f); err != nil {
		return
	}

	to.Set(m)
	return
}

func (d Decoder) decodeValueMapInterfaceInterface(typ Type, to reflect.Value) error {
	m := to.Interface().(map[interface{}]interface{})

	if m == nil {
		m = make(map[interface{}]interface{})
		to.Set(reflect.ValueOf(m))
	}

	for k := range m {
		delete(m, k)
	}

	return d.decodeMapFromType(typ, func(kd Decoder, vd Decoder) (err error) {
		var k interface{}
		var v interface{}

		if err = kd.Decode(&k); err != nil {
			return
		}
		if err = vd.Decode(&v); err != nil {
			return
		}

		m[k] = v
		return
	})
}

func (d Decoder) decodeValueMapStringInterface(typ Type, to reflect.Value) (err error) {
	m := to.Interface().(map[string]interface{})

	if m == nil {
		m = make(map[string]interface{})
		to.Set(reflect.ValueOf(m))
	}

	for k := range m {
		delete(m, k)
	}

	return d.decodeMapFromType(typ, func(kd Decoder, vd Decoder) (err error) {
		var b []byte
		var k string
		var v interface{}

		if b, err = d.decodeTypeAndString(); err != nil {
			return
		}
		k = string(b)

		if err = vd.Decode(&v); err != nil {
			return
		}

		m[k] = v
		return
	})
}

func (d Decoder) decodeValueMapStringString(typ Type, to reflect.Value) (err error) {
	m := to.Interface().(map[string]string)

	if m == nil {
		m = make(map[string]string)
		to.Set(reflect.ValueOf(m))
	}

	for k := range m {
		delete(m, k)
	}

	return d.decodeMapFromType(typ, func(kd Decoder, vd Decoder) (err error) {
		var b []byte
		var k string
		var v string

		if b, err = d.decodeTypeAndString(); err != nil {
			return
		}
		k = string(b)

		if err = d.decodeMapValue(vd.off - 1); err != nil {
			return
		}

		if b, err = d.decodeTypeAndString(); err != nil {
			return
		}
		v = string(b)

		m[k] = v
		return
	})
}

func (d Decoder) decodeValueStruct(to reflect.Value) (Type, error) {
	return d.decodeValueStructWith(to, LookupStruct(to.Type()))
}

func (d Decoder) decodeValueStructWith(to reflect.Value, s *Struct) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueStructFromTypeWith(t, to, s)
	}
	return
}

func (d Decoder) decodeValueStructFromType(typ Type, to reflect.Value) (err error) {
	return d.decodeValueStructFromTypeWith(typ, to, LookupStruct(to.Type()))
}

func (d Decoder) decodeValueStructFromTypeWith(typ Type, to reflect.Value, s *Struct) (err error) {
	if err = d.decodeMapFromType(typ, func(kd Decoder, vd Decoder) (err error) {
		var b []byte

		if b, err = d.decodeTypeAndString(); err != nil {
			return
		}

		if err = d.decodeMapValue(vd.off - 1); err != nil {
			return
		}

		f := s.FieldsByName[string(b)]
		if f == nil {
			var v interface{} // discard
			return d.Decode(&v)
		}

		_, err = f.decode(d, to.FieldByIndex(f.Index))
		return
	}); err != nil {
		to.Set(zeroValueOf(to.Type()))
	}
	return
}

func (d Decoder) decodeValuePointer(to reflect.Value) (Type, error) {
	return d.decodeValuePointerWith(to, decodeFuncOf(to.Type().Elem()))
}

func (d Decoder) decodeValuePointerWith(to reflect.Value, f decodeFunc) (typ Type, err error) {
	var t = to.Type()
	var v reflect.Value

	if to.IsNil() {
		v = reflect.New(t.Elem())
	} else {
		v = to
	}

	if typ, err = f(d, v.Elem()); err != nil {
		return
	}

	if to.CanSet() {
		switch {
		case typ == Nil:
			to.Set(zeroValueOf(t))

		case to.IsNil():
			to.Set(v)
		}
	}

	return
}

func (d Decoder) decodeValueDecoder(to reflect.Value) (Type, error) {
	return Bool /* just needs to not be Nil */, to.Interface().(ValueDecoder).DecodeValue(d)
}

func (d Decoder) decodeValueTextUnmarshaler(to reflect.Value) (t Type, err error) {
	var b []byte
	var v = reflect.ValueOf(&b).Elem()

	if t, err = d.decodeValueBytes(v); err != nil {
		return
	}

	err = to.Interface().(encoding.TextUnmarshaler).UnmarshalText(b)
	return
}

func (d Decoder) decodeValueInterface(to reflect.Value) (t Type, err error) {
	if t, err = d.decodeType(); err == nil {
		err = d.decodeValueInterfaceFromType(t, to)
	}
	return
}

func (d Decoder) decodeValueInterfaceFromType(t Type, to reflect.Value) (err error) {
	switch t {
	case Nil:
		err = d.decodeValueInterfaceFromNil(to)
	case Bool:
		err = d.decodeValueInterfaceFrom(boolType, t, to, Decoder.decodeValueBoolFromType)
	case Int:
		err = d.decodeValueInterfaceFrom(int64Type, t, to, Decoder.decodeValueIntFromType)
	case Uint:
		err = d.decodeValueInterfaceFrom(uint64Type, t, to, Decoder.decodeValueUintFromType)
	case Float:
		err = d.decodeValueInterfaceFrom(float64Type, t, to, Decoder.decodeValueFloatFromType)
	case String:
		err = d.decodeValueInterfaceFrom(stringType, t, to, Decoder.decodeValueStringFromType)
	case Bytes:
		err = d.decodeValueInterfaceFrom(bytesType, t, to, Decoder.decodeValueBytesFromType)
	case Time:
		err = d.decodeValueInterfaceFrom(timeType, t, to, Decoder.decodeValueTimeFromType)
	case Duration:
		err = d.decodeValueInterfaceFrom(durationType, t, to, Decoder.decodeValueDurationFromType)
	case Error:
		err = d.decodeValueInterfaceFrom(errorInterface, t, to, Decoder.decodeValueErrorFromType)
	case Array:
		err = d.decodeValueInterfaceFrom(sliceInterfaceType, t, to, Decoder.decodeValueSliceFromType)
	case Map:
		err = d.decodeValueInterfaceFrom(mapInterfaceInterfaceType, t, to, Decoder.decodeValueMapFromType)
	default:
		panic("objconv: parser returned an unsupported value type: " + t.String())
	}
	return
}

func (d Decoder) decodeValueInterfaceFromNil(to reflect.Value) (err error) {
	if err = d.decodeNil(); err == nil {
		to.Set(zeroValueOf(to.Type()))
	}
	return
}

func (d Decoder) decodeValueInterfaceFrom(from reflect.Type, t Type, to reflect.Value, decode func(Decoder, Type, reflect.Value) error) (err error) {
	v := reflect.New(from).Elem()

	if err = decode(d, t, v); err == nil {
		to.Set(v)
	}

	return
}

func (d Decoder) decodeValueUnsupported(to reflect.Value) (Type, error) {
	return Nil, fmt.Errorf("objconv: the decoder doesn't support values of type %s", to.Type())
}

func (d Decoder) decodeTypeAndString() (b []byte, err error) {
	var t Type

	if t, err = d.decodeType(); err == nil {
		// This algorithm is the same than the one used in
		// decodeValueStringWithType, and should be kept in sync.
		switch t {
		case Nil:
			err = d.decodeNil()
		case String:
			b, err = d.decodeString()
		case Bytes:
			b, err = d.decodeBytes()
		default:
			err = typeConversionError(t, String)
		}
	}

	return
}

func (d Decoder) decodeType() (Type, error) { return d.Parser.ParseType() }

func (d Decoder) decodeNil() error { return d.Parser.ParseNil() }

func (d Decoder) decodeBool() (bool, error) { return d.Parser.ParseBool() }

func (d Decoder) decodeInt() (int64, error) { return d.Parser.ParseInt() }

func (d Decoder) decodeUint() (uint64, error) { return d.Parser.ParseUint() }

func (d Decoder) decodeFloat() (float64, error) { return d.Parser.ParseFloat() }

func (d Decoder) decodeString() ([]byte, error) { return d.Parser.ParseString() }

func (d Decoder) decodeBytes() ([]byte, error) { return d.Parser.ParseBytes() }

func (d Decoder) decodeTime() (time.Time, error) { return d.Parser.ParseTime() }

func (d Decoder) decodeDuration() (time.Duration, error) { return d.Parser.ParseDuration() }

func (d Decoder) decodeError() (error, error) { return d.Parser.ParseError() }

// DecodeArray provides the implementation of the algorithm for decoding arrays,
// where f is called to decode each element of the array.
//
// The method returns the underlying type of the value returned by the parser.
func (d Decoder) DecodeArray(f func(Decoder) error) (t Type, err error) {
	if t, err = d.decodeType(); err != nil {
		return
	}

	if d.off != 0 {
		if d.off, err = 0, d.decodeMapValue(d.off-1); err != nil {
			return
		}
	}

	err = d.decodeArrayFromType(t, f)
	return
}

func (d Decoder) decodeArrayFromType(t Type, f func(Decoder) error) (err error) {
	var n int

	switch t {
	case Nil:
		err = d.decodeNil()
		return

	case Array:
		n, err = d.decodeArrayBegin()

	default:
		err = typeConversionError(t, Array)
	}

	if err != nil {
		return
	}

	i := 0

	for n < 0 || i < n {
		if n < 0 || i != 0 {
			if err = d.decodeArrayNext(i); err != nil {
				if err == End {
					err = nil
					break
				}
				return
			}
		}
		if err = f(d); err != nil {
			return
		}
		i++
	}

	err = d.decodeArrayEnd(i)
	return
}

func (d Decoder) decodeArrayBegin() (int, error) { return d.Parser.ParseArrayBegin() }

func (d Decoder) decodeArrayEnd(n int) error { return d.Parser.ParseArrayEnd(n) }

func (d Decoder) decodeArrayNext(n int) error { return d.Parser.ParseArrayNext(n) }

// DecodeMap provides the implementation of the algorithm for decoding maps,
// where f is called to decode each pair of key and value.
//
// The function f is expected to decode two values from the map, the first one
// being the key and the second the associated value. The first decoder must be
// used to decode the key, the second one for the value.
//
// The method returns the underlying type of the value returned by the parser.
func (d Decoder) DecodeMap(f func(Decoder, Decoder) error) (t Type, err error) {
	if t, err = d.decodeType(); err != nil {
		return
	}

	if d.off != 0 {
		if d.off, err = 0, d.decodeMapValue(d.off-1); err != nil {
			return
		}
	}

	err = d.decodeMapFromType(t, f)
	return
}

func (d Decoder) decodeMapFromType(t Type, f func(Decoder, Decoder) error) (err error) {
	var n int

	switch t {
	case Nil:
		err = d.decodeNil()
		return

	case Map:
		n, err = d.decodeMapBegin()

	default:
		err = typeConversionError(t, Map)
	}

	if err != nil {
		return
	}

	i := 0

	for n < 0 || i < n {
		if n < 0 || i != 0 {
			if err = d.decodeMapNext(i); err != nil {
				if err == End {
					err = nil
					break
				}
				return
			}
		}

		d1 := d
		d2 := d
		d2.off = i + 1

		if err = f(d1, d2); err != nil {
			return
		}

		i++
	}

	err = d.decodeMapEnd(i)
	return
}

func (d Decoder) decodeMapBegin() (int, error) { return d.Parser.ParseMapBegin() }

func (d Decoder) decodeMapEnd(n int) error { return d.Parser.ParseMapEnd(n) }

func (d Decoder) decodeMapValue(n int) error { return d.Parser.ParseMapValue(n) }

func (d Decoder) decodeMapNext(n int) error { return d.Parser.ParseMapNext(n) }

// StreamDecoder decodes values in a streaming fashion, allowing an array to be
// consumed without loading it fully in memory.
//
// Instances of StreamDecoder are not safe for use by multiple goroutines.
type StreamDecoder struct {
	// Parser to use to load values.
	Parser Parser

	// DecodeMap can be set to a function used to decode maps when there is no
	// destination type (like when decoding to an empty interface for example).
	DecodeMapFunc func(Decoder, Decoder) error

	err error
	typ Type
	cnt int
	max int
}

// NewStreamDecoder returns a new stream decoder that takes input from p.
//
// The funciton panics if p is nil.
func NewStreamDecoder(p Parser) *StreamDecoder {
	if p == nil {
		panic("objconv: the parser is nil")
	}
	return &StreamDecoder{Parser: p}
}

// Err returns the last error returned by the Decode method.
//
// The method returns nil if the stream reached its natural end.
func (d *StreamDecoder) Err() error {
	if d.err == End {
		return nil
	}
	return d.err
}

// Decodes the next value from the stream into v.
func (d *StreamDecoder) Decode(v interface{}) error {
	cnt := d.cnt
	max := d.max
	dec := Decoder{
		Parser:        d.Parser,
		DecodeMapFunc: d.DecodeMapFunc,
	}

	typ, err := d.decodeType()
	if err != nil {
		return err
	}

	switch typ {
	default:
		if max = 1; cnt == max {
			err = End
		}
	case Array:
		if cnt == 0 {
			max, err = dec.decodeArrayBegin()
		}
		if cnt == max {
			err = dec.decodeArrayEnd(cnt)
		} else if cnt != 0 {
			err = dec.decodeArrayNext(cnt)
		}
	}

	if err == nil {
		if cnt == max {
			err = End
		} else {
			switch err = dec.Decode(v); err {
			case nil:
				cnt++
			case End:
				cnt++
				max = cnt
			}
		}
	}

	d.err = err
	d.cnt = cnt
	d.max = max
	return err
}

// Encoder returns a new StreamEncoder which can be used to re-encode the stream
// decoded by d into e.
//
// The method panics if e is nil.
func (d *StreamDecoder) Encoder(e Emitter) (enc *StreamEncoder, err error) {
	var typ Type

	if typ, err = d.decodeType(); err == nil {
		enc = NewStreamEncoder(e)
		enc.oneshot = typ != Array
	}

	return
}

func (d *StreamDecoder) decodeType() (Type, error) {
	if d.err == nil && d.typ == Unknown {
		d.typ, d.err = d.Parser.ParseType()
	}
	return d.typ, d.err
}

// ValueDecoder is the interface that can be implemented by types that wish to
// provide their own decoding algorithms.
//
// The DecodeValue method is called when the value is found by a decoding
// algorithm.
type ValueDecoder interface {
	DecodeValue(Decoder) error
}

// ValueDecoderFunc allos the use of regular functions or methods as value
// decoders.
type ValueDecoderFunc func(Decoder) error

// DecodeValue calls f(d).
func (f ValueDecoderFunc) DecodeValue(d Decoder) error { return f(d) }

type decodeFuncOpts struct {
	recurse bool
	structs map[reflect.Type]*Struct
}

type decodeFunc func(Decoder, reflect.Value) (Type, error)

func decodeFuncOf(t reflect.Type) decodeFunc {
	return makeDecodeFunc(t, decodeFuncOpts{})
}

func makeDecodeFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	// fast path: check if it's a basic go type
	switch t {
	case boolType:
		return Decoder.decodeValueBool

	case stringType:
		return Decoder.decodeValueString

	case bytesType:
		return Decoder.decodeValueBytes

	case timeType:
		return Decoder.decodeValueTime

	case durationType:
		return Decoder.decodeValueDuration

	case emptyInterface:
		return Decoder.decodeValueInterface

	case intType, int8Type, int16Type, int32Type, int64Type:
		return Decoder.decodeValueInt

	case uintType, uint8Type, uint16Type, uint32Type, uint64Type, uintptrType:
		return Decoder.decodeValueUint

	case float32Type, float64Type:
		return Decoder.decodeValueFloat
	}

	// check if it implements one of the special case interfaces
	switch p := reflect.PtrTo(t); {
	case p.Implements(valueDecoderInterface):
		return Decoder.decodeValueDecoder

	case p.Implements(textUnmarshalerInterface):
		return Decoder.decodeValueTextUnmarshaler

	case t.Implements(errorInterface):
		return Decoder.decodeValueError
	}

	// check what kind is the type, potentially generate a decoder
	switch t.Kind() {
	case reflect.Struct:
		return makeDecodeStructFunc(t, opts)

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return Decoder.decodeValueBytes
		}
		return makeDecodeSliceFunc(t, opts)

	case reflect.Map:
		return makeDecodeMapFunc(t, opts)

	case reflect.Ptr:
		return makeDecodePtrFunc(t, opts)

	case reflect.Array:
		return makeDecodeArrayFunc(t, opts)

	case reflect.Bool:
		return Decoder.decodeValueBool

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Decoder.decodeValueInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return Decoder.decodeValueUint

	case reflect.Float32, reflect.Float64:
		return Decoder.decodeValueFloat

	case reflect.String:
		return Decoder.decodeValueString

	default:
		return Decoder.decodeValueUnsupported
	}
}

func makeDecodeSliceFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodeValueSlice
	}
	f := makeDecodeFunc(t.Elem(), opts)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodeValueSliceWith(v, f)
	}
}

func makeDecodeArrayFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodeValueArray
	}
	f := makeDecodeFunc(t.Elem(), opts)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodeValueArrayWith(v, f)
	}
}

func makeDecodeMapFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodeValueMap
	}
	kf := makeDecodeFunc(t.Key(), opts)
	vf := makeDecodeFunc(t.Elem(), opts)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodeValueMapWith(v, kf, vf)
	}
}

func makeDecodeStructFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodeValueStruct
	}
	s := newStruct(t, opts.structs)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodeValueStructWith(v, s)
	}
}

func makeDecodePtrFunc(t reflect.Type, opts decodeFuncOpts) decodeFunc {
	if !opts.recurse {
		return Decoder.decodeValuePointer
	}
	f := makeDecodeFunc(t.Elem(), opts)
	return func(d Decoder, v reflect.Value) (Type, error) {
		return d.decodeValuePointerWith(v, f)
	}
}
