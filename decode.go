package objconv

/*
// DecoderConfig carries the configuration for creating an encoder.
type DecoderConfig struct {
	// Input is the data stream that the decoder reads from.
	Input io.Reader

	// Parser defines the format used by the decoder.
	Parser Parser
}

type Decoder struct {
	p  Parser
	r  *Reader
	r2 Reader
}

// NewDecoder returns a new decoder configured with config.
func NewDecoder(config DecoderConfig) *Decoder {
	if config.Input == nil {
		panic("objconv.NewDecoder: the input is nil")
	}
	if config.Parser == nil {
		panic("objconv.NewDecoder: the parser is nil")
	}
	d := makeDecoder(config.Parser, config.Input)
	return &d
}

func makeDecoder(p Parser, r io.Reader) (d Decoder) {
	d.p = p
	d.r = &d.r2

	// Use the reader directly if it's already an instance of Reader.
	switch x := r.(type) {
	case *Reader:
		d.r = x
	default:
		d.r.r = r
	}

	return
}

func (d *Decoder) Decode(v interface{}) (err error) {
	return
}
*/

/*
	return d.decode(d.r, v)
}

func (d *Decoder) parse(r *Reader, v interface{}) (interface{}, reflect.Value, error) {
	to := reflect.ValueOf(v).Elem()
	v, err := d.p.Parse(r, to.Interface())
	return v, to, err
}

func (d *Decoder) decode(r *Reader, v interface{}) error {
	from, to, err := d.parse(r, v)
	if err == nil {
		err = d.decodeValue(r, from, to)
	}
	return err
}

func (d *Decoder) decodeValue(r *Reader, v interface{}, to reflect.Value) error {
	switch x := v.(type) {
	case bool:
		return d.decodeBool(x, to)

	case int64:
		return d.decodeInt(x, to)

	case uint64:
		return d.decodeUint(x, to)

	case float64:
		return d.decodeFloat(x, to)

	case string:
		return d.decodeString(x, to)

	case []byte:
		return d.decodeBytes(x, to)

	case time.Time:
		return d.decodeTime(x, to)

	case time.Duration:
		return d.decodeDuration(x, to)

	case error:
		return d.decodeError(x, to)

	case ArrayParser:
		return d.decodeArray(r, x, to)

	case MapParser:
		return d.decodeMap(r, x, to)

	default:
		if x == nil {
			return d.decodeNil(to)
		} else {
			return fmt.Errorf("the parser produced an unsupported value of type %T, this is a bug", x)
		}
	}
}

func (d *Decoder) decodeNil(to reflect.Value) (err error) {
	switch to.Kind() {
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

	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map:
		if !to.IsNil() {
			to.Set(reflect.Zero(to.Type()))
		}

	default:
		to.Set(reflect.Zero(to.Type()))
	}

	return
}

func (d *Decoder) decodeBool(v bool, to reflect.Value) (err error) {
	switch to.Kind() {
	case reflect.Bool:
		to.SetBool(v)
	default:
		err = setValue(to, reflect.ValueOf(v))
	}
	return
}

func (d *Decoder) decodeInt(v int64, to reflect.Value) (err error) {
	switch to.Kind() {
	case reflect.Int:
		if _, err = convertInt64ToInt(v); err == nil {
			to.SetInt(v)
		}

	case reflect.Int8:
		if _, err = convertInt64ToInt(v); err == nil {
			to.SetInt(v)
		}

	case reflect.Int16:
		if _, err = convertInt64ToInt16(v); err == nil {
			to.SetInt(v)
		}

	case reflect.Int32:
		if _, err = convertInt64ToInt32(v); err == nil {
			to.SetInt(v)
		}

	case reflect.Int64:
		to.SetInt(v)

	case reflect.Uint:
		if _, err = convertInt64ToUint(v); err == nil {
			to.SetUint(uint64(v))
		}

	case reflect.Uint8:
		if _, err = convertInt64ToUint8(v); err == nil {
			to.SetUint(uint64(v))
		}

	case reflect.Uint16:
		if _, err = convertInt64ToUint16(v); err == nil {
			to.SetUint(uint64(v))
		}

	case reflect.Uint32:
		if _, err = convertInt64ToUint32(v); err == nil {
			to.SetUint(uint64(v))
		}

	case reflect.Uint64:
		if _, err = convertInt64ToUint64(v); err == nil {
			to.SetUint(uint64(v))
		}

	case reflect.Uintptr:
		if _, err = convertInt64ToUintptr(v); err == nil {
			to.SetUint(uint64(v))
		}

	case reflect.Float32:
		if _, err = convertInt64ToFloat32(v); err == nil {
			to.SetFloat(float64(v))
		}

	case reflect.Float64:
		if _, err = convertInt64ToFloat64(v); err == nil {
			to.SetFloat(float64(v))
		}

	default:
		err = setValue(to, reflect.ValueOf(v))
	}

	return
}

func (d *Decoder) decodeUint(v uint64, to reflect.Value) (err error) {
	switch to.Kind() {
	case reflect.Int:
		if _, err = convertUint64ToInt(v); err == nil {
			to.SetInt(int64(v))
		}

	case reflect.Int8:
		if _, err = convertUint64ToInt(v); err == nil {
			to.SetInt(int64(v))
		}

	case reflect.Int16:
		if _, err = convertUint64ToInt16(v); err == nil {
			to.SetInt(int64(v))
		}

	case reflect.Int32:
		if _, err = convertUint64ToInt32(v); err == nil {
			to.SetInt(int64(v))
		}

	case reflect.Int64:
		if _, err = convertUint64ToInt64(v); err == nil {
			to.SetInt(int64(v))
		}

	case reflect.Uint:
		if _, err = convertUint64ToUint(v); err == nil {
			to.SetUint(v)
		}

	case reflect.Uint8:
		if _, err = convertUint64ToUint8(v); err == nil {
			to.SetUint(v)
		}

	case reflect.Uint16:
		if _, err = convertUint64ToUint16(v); err == nil {
			to.SetUint(v)
		}

	case reflect.Uint32:
		if _, err = convertUint64ToUint32(v); err == nil {
			to.SetUint(v)
		}

	case reflect.Uint64:
		to.SetUint(v)

	case reflect.Uintptr:
		if _, err = convertUint64ToUintptr(v); err == nil {
			to.SetUint(v)
		}

	case reflect.Float32:
		if _, err = convertUint64ToFloat32(v); err == nil {
			to.SetUint(v)
		}

	case reflect.Float64:
		if _, err = convertUint64ToFloat64(v); err == nil {
			to.SetUint(v)
		}

	default:
		err = setValue(to, reflect.ValueOf(v))
	}

	return
}

func (d *Decoder) decodeFloat(v float64, to reflect.Value) (err error) {
	switch to.Kind() {
	case reflect.Float32, reflect.Float64:
		to.SetFloat(v)
	default:
		err = setValue(to, reflect.ValueOf(v))
	}
	return
}

func (d *Decoder) decodeString(v string, to reflect.Value) (err error) {
	switch to.Kind() {
	case reflect.String:
		to.SetString(v)

	case reflect.Slice:
		err = d.decodeStringToSlice(v, to)

	default:
		err = d.decodeStringToOther(v, to)
	}
	return
}

func (d *Decoder) decodeStringToSlice(v string, to reflect.Value) (err error) {
	switch to.Type().Elem().Kind() {
	case reflect.Uint8: // []byte
		to.SetBytes([]byte(v))

	case reflect.Int32: // []rune
		to.Set(reflect.ValueOf([]rune(v)))

	default:
		err = setValue(to, reflect.ValueOf(v))
	}
	return
}

func (d *Decoder) decodeStringToOther(v string, to reflect.Value) (err error) {
	switch to.Interface().(type) {
	case time.Time:
		var t time.Time
		if t, err = time.Parse(time.RFC3339Nano, v); err == nil {
			to.Set(reflect.ValueOf(t))
		}

	case time.Duration:
		var d time.Duration
		if d, err = time.ParseDuration(v); err == nil {
			to.Set(reflect.ValueOf(d))
		}

	case error:
		to.Set(reflect.ValueOf(errors.New(v)))

	default:
		err = setValue(to, reflect.ValueOf(v))
	}
	return
}

func (d *Decoder) decodeBytes(v []byte, to reflect.Value) (err error) {
	switch to.Kind() {
	case reflect.Slice:
		err = d.decodeBytesToSlice(v, to)

	case reflect.String:
		to.SetString(string(v))

	default:
		err = setValue(to, reflect.ValueOf(v))
	}
	return
}

func (d *Decoder) decodeBytesToSlice(v []byte, to reflect.Value) (err error) {
	switch to.Type().Elem().Kind() {
	case reflect.Uint8: // []byte
		to.SetBytes(v)

	case reflect.Int32: // []rune
		to.Set(reflect.ValueOf([]rune(string(v))))

	default:
		err = setValue(to, reflect.ValueOf(v))
	}
	return
}

func (d *Decoder) decodeDuration(v time.Duration, to reflect.Value) error {
	return setValue(to, reflect.ValueOf(v))
}

func (d *Decoder) decodeTime(v time.Time, to reflect.Value) error {
	return setValue(to, reflect.ValueOf(v))
}

func (d *Decoder) decodeError(v error, to reflect.Value) error {
	return setValue(to, reflect.ValueOf(v))
}

func (d *Decoder) decodeArray(r *Reader, v ArrayParser, to reflect.Value) (err error) {
	switch t := to.Type(); t.Kind() {
	case reflect.Slice:
		err = d.decodeArrayToSlice(r, v, to, t)
	default:
		err = d.decodeArrayToInterface(r, v, to)
	}
	return
}

func (d *Decoder) decodeArrayToSlice(r *Reader, v ArrayParser, to reflect.Value, t reflect.Type) (err error) {
	n := v.Len()

	if n < 0 {
		n = 20
	}

	s := reflect.MakeSlice(t, 0, n)
	z := reflect.Zero(t.Elem())
	h := z.Interface()

	for i := 0; true; i++ {
		x, e := v.Parse(r, h)

		if e == io.EOF {
			break
		}

		if e != nil {
			err = e
			return
		}

		s = reflect.Append(s, z)

		if err = d.decodeValue(r, x, s.Index(i)); err != nil {
			return
		}
	}

	to.Set(s)
	return
}

func (d *Decoder) decodeArrayToInterface(r *Reader, v ArrayParser, to reflect.Value) (err error) {
	n := v.Len()

	if n < 0 {
		n = 20
	}

	s := make([]interface{}, 0, n)

	for i := 0; true; i++ {
		x, e := v.Parse(r, nil)

		if e == io.EOF {
			break
		}

		if e != nil {
			err = e
			return
		}

		switch s = append(s, nil); x.(type) {
		case ArrayParser, MapParser:
			err = d.decodeValue(r, x, reflect.ValueOf(&s[i]).Elem())
		default:
			s[i] = x
		}

		if err != nil {
			return
		}
	}

	to.Set(reflect.ValueOf(s))
	return
}

func (d *Decoder) decodeMap(r *Reader, v MapParser, to reflect.Value) (err error) {
	switch t := to.Type(); t.Kind() {
	case reflect.Struct:
		err = d.decodeMapToStruct(r, v, to, t)

	case reflect.Map:
		err = d.decodeMapToMap(r, v, to, t)

	default:
		var m map[interface{}]interface{}
		var mv = reflect.ValueOf(&m).Elem()
		if err = d.decodeMapToMap(r, v, mv, reflect.TypeOf(m)); err != nil {
			err = setValue(to, mv)
		}
	}
	return
}

func (d *Decoder) decodeMapToMap(r *Reader, v MapParser, to reflect.Value, t reflect.Type) (err error) {
	m := reflect.MakeMap(t)

	kt := t.Key()
	vt := t.Elem()

	ke := reflect.New(kt).Elem()
	ve := reflect.New(vt).Elem()

	ki := ke.Interface()
	vi := ve.Interface()

	for {
		x, e := v.ParseKey(r, ki)
		if e == io.EOF {
			break
		}
		if e != nil {
			err = e
			return
		}
		if err = d.decodeValue(r, x, ke); err != nil {
			return
		}

		x, e = v.ParseValue(r, vi)
		if e == io.EOF {
			break
		}
		if e != nil {
			err = e
			return
		}
		if err = d.decodeValue(r, x, ve); err != nil {
			return
		}

		m.SetMapIndex(ke, ve)
	}

	to.Set(m)
	return
}

func (d *Decoder) decodeMapToStruct(r *Reader, v MapParser, to reflect.Value, t reflect.Type) (err error) {
	s := LookupStruct(t)
	f := ""

	fv := reflect.ValueOf(&f).Elem()
	fi := fv.Interface()

	for {
		x, e := v.ParseKey(r, fi)
		if e == io.EOF {
			break
		}
		if e != nil {
			err = e
			return
		}
		if err = d.decodeValue(r, x, fv); err != nil {
			return
		}

		if sf := s.FieldByName[f]; sf != nil {
			fv := to.FieldByIndex(sf.Index)
			x, e := v.ParseValue(r, fv.Interface())
			if e != nil {
				err = e
				return
			}
			if err = d.decodeValue(r, x, fv); err != nil {
				return
			}
		}

		f = ""
	}

	return
}
*/

/*
type StreamDecoder struct {
	dec    Decoder
	parser ArrayParser
	err    error
	count  int
	array  bool
}

// NewStreamDecoder returns a new stream decoder configured with config.
func NewStreamDecoder(config DecoderConfig) StreamDecoder {
	if config.Input == nil {
		panic("objconv.NewStreamDecoder: the input is nil")
	}
	if config.Parse == nil {
		panic("objconv.NewStreamDecoder: the parser is nil")
	}
	return &StreamDecoder{
		dec: Decoder{
			r: Reader{r: config.Input},
		},
	}
}

func (d *StreamDecoder) Decode(v interface{}) (err error) {
	if err = d.err; err != nil {
		return
	}

	ve := reflect.ValueOf(v).Elem()
	vi := ve.Interface()

	if d.parser == nil {
		from, _, e := d.parse(d.r, v)
		if e != nil {
			d.err, err = e, e
			return
		}
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

	var x interface{}

	if x, err = d.parser.Parse(d.r, vi); err != nil {
		d.err = err
		return
	}

	if err = d.decodeValue(d.r, x, ve); err != nil {
		d.err = err
		return
	}

	d.count++
	return
}

// Len returns the expected number of elements returned from the stream.
//
// Depending on the actual format that the stream is decoding this value
// may or may not be accurate, some formats may also return a negative
// value to indicate that the number of elements is unknown.
func (d *StreamDecoder) Len() (n int) {
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

		v, _, e := d.parse(d.r, &z)
		if e != nil {
			d.err = e
			return
		}
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

// Error returns the last error encountered by the decoder.
func (d *StreamDecoder) Err() error {
	return d.err
}

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
func (d *StreamDecoder) Encoder(config EncoderConfig) *StreamEncoder {
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
*/
