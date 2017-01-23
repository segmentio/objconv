package objtests

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/segmentio/objconv"
	_ "github.com/segmentio/objconv/adapters"
)

// TestValues is an array of all the values used by the TestCodec suite.
var TestValues = [...]interface{}{
	// constants
	nil,
	false,
	true,

	// int
	0,
	1,
	23,
	24,
	127,
	-1,
	-10,
	-31,
	-32,
	objconv.Int8Min,
	objconv.Int8Max + 1,
	objconv.Int8Min - 1,
	objconv.Int16Max,
	objconv.Int16Min,
	objconv.Int16Max + 1,
	objconv.Int16Min - 1,
	objconv.Int32Max,
	objconv.Int32Min,
	objconv.Int32Max + 1,
	objconv.Int32Min - 1,
	int64(objconv.Int64Max),
	int64(objconv.Int64Min),

	// uint
	uint(0),
	uint(1),
	uint8(objconv.Uint8Max),
	uint16(objconv.Uint8Max) + 1,
	uint16(objconv.Uint16Max),
	uint32(objconv.Uint16Max) + 1,
	uint32(objconv.Uint32Max),
	uint64(objconv.Uint32Max) + 1,

	// float
	float32(0),
	float32(objconv.Float32IntMin),
	float32(objconv.Float32IntMax),
	float64(0),
	float64(0.5),

	// string
	"",
	"Hello World!",
	"Hello\"World!",
	"Hello\\World!",
	"Hello\nWorld!",
	"Hello\rWorld!",
	"Hello\tWorld!",
	"Hello\bWorld!",
	"Hello\fWorld!",
	"你好",
	strings.Repeat("A", 32),
	strings.Repeat("A", objconv.Uint8Max+1),
	strings.Repeat("A", objconv.Uint16Max+1),

	// bytes
	[]byte(""),
	[]byte("Hello World!"),
	bytes.Repeat([]byte("A"), objconv.Uint8Max+1),
	bytes.Repeat([]byte("A"), objconv.Uint16Max+1),

	// duration
	time.Nanosecond,
	time.Microsecond,
	time.Millisecond,
	time.Second,
	time.Minute,
	time.Hour,

	// time
	time.Unix(0, 0).In(time.UTC),
	time.Unix(1, 42).In(time.UTC),
	time.Unix(17179869184, 999999999).In(time.UTC),
	time.Date(2016, 12, 20, 0, 20, 1, 0, time.UTC),

	// error
	errors.New(""),
	errors.New("Hello World!"),
	errors.New(strings.Repeat("A", objconv.Uint8Max+1)),
	errors.New(strings.Repeat("A", objconv.Uint16Max+1)),

	// array
	[]int{},
	[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	make([]int, objconv.Uint8Max+1),
	make([]int, objconv.Uint16Max+1),

	// map
	makeMap(0),
	makeMap(15),
	makeMap(objconv.Uint8Max + 1),
	makeMap(objconv.Uint16Max + 1),

	// struct
	struct{}{},
	struct{ A int }{42},
	struct{ A, B, C int }{1, 2, 3},
	struct {
		A int
		T time.Time
		S string
	}{42, time.Date(2016, 12, 20, 0, 20, 1, 0, time.UTC), "Hello World!"},

	// net
	net.TCPAddr{
		IP:   net.ParseIP("::1"),
		Port: 4242,
		Zone: "zone",
	},
	net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: 4242,
		Zone: "zone",
	},
	net.IPAddr{
		IP:   net.ParseIP("::1"),
		Zone: "zone",
	},
	net.IPv4(127, 0, 0, 1),

	// url
	parseURL("http://localhost:4242/hello/world?answer=42#question"),
	parseQuery("answer=42&message=Hello+World"),

	// mail
	parseEmail("git@github.com"),
	parseEmailList("Alice <alice@example.com>, Bob <bob@example.com>, Eve <eve@example.com>"),
}

func makeMap(n int) map[string]string {
	m := make(map[string]string, n)
	for i := 0; i != n; i++ {
		m[strconv.Itoa(i)] = "A"
	}
	return m
}

func testName(v interface{}) string {
	s := fmt.Sprintf("%T:%v", v, v)
	if len(s) > 42 {
		s = s[:42] + "..."
	}
	return s
}

// TestCodec implements a test suite for validating that a codec properly
// support encoding and decoding values of different types. The function also
// tests that the codec behaves properly when used with stream encoders and
// decoders.
func TestCodec(t *testing.T, codec objconv.Codec) {
	t.Run("Values", func(t *testing.T) { testCodecValues(t, codec) })
	t.Run("Stream", func(t *testing.T) { testCodecStream(t, codec) })
}

func newValue(model interface{}) reflect.Value {
	if model == nil {
		return reflect.New(reflect.TypeOf(&model).Elem())
	} else {
		return reflect.New(reflect.TypeOf(model))
	}
}

func testCodecValues(t *testing.T, codec objconv.Codec) {
	b := &bytes.Buffer{}
	b.Grow(1024)

	for _, v1 := range TestValues {
		t.Run(testName(v1), func(t *testing.T) {
			b.Reset()
			e := objconv.NewEncoder(codec.NewEmitter(b))
			d := objconv.NewDecoder(codec.NewParser(b))
			v2 := newValue(v1)

			if err := e.Encode(v1); err != nil {
				t.Error(err)
				return
			}

			if err := d.Decode(v2.Interface()); err != nil {
				t.Error(err)
				return
			}

			x1 := v1
			x2 := v2.Elem().Interface()

			if !reflect.DeepEqual(x1, x2) {
				t.Errorf("%#v", x2)
			}
		})
	}
}

func testCodecStream(t *testing.T, codec objconv.Codec) {
	r, w := io.Pipe()
	defer r.Close()

	e := objconv.NewStreamEncoder(codec.NewEmitter(w))
	d := objconv.NewStreamDecoder(codec.NewParser(r))

	go func() {
		defer w.Close()
		defer e.Close()

		for _, v := range TestValues {
			if err := e.Encode(v); err != nil {
				if err != io.ErrClosedPipe {
					t.Error(err)
				}
				return
			}
		}
	}()

	for _, v1 := range TestValues {
		v2 := newValue(v1)

		if err := d.Decode(v2.Interface()); err != nil {
			return
		}

		x1 := v1
		x2 := v2.Elem().Interface()

		if !reflect.DeepEqual(x1, x2) {
			t.Errorf("%#v", x2)
		}
	}

	var v interface{}
	if err := d.Decode(&v); err == nil {
		t.Error("too many values decoded from the stream")
	}

	if err := d.Err(); err != nil {
		t.Error(err)
	}
}

type counter struct {
	n int
}

func (c *counter) Write(b []byte) (n int, err error) {
	n = len(b)
	c.n += n
	return
}

// BenchmarkCodec implements a benchmark suite for codecs, making it easy to get
// comparable performance results for various formats.
func BenchmarkCodec(b *testing.B, codec objconv.Codec) {
	b.Run("Encoder", func(b *testing.B) { benchmarkEncoder(b, codec) })
	b.Run("Decoder", func(b *testing.B) { benchmarkDecoder(b, codec) })
	b.Run("StreamEncoder", func(b *testing.B) { benchmarkStreamEncoder(b, codec) })
	b.Run("StreamDecoder", func(b *testing.B) { benchmarkStreamDecoder(b, codec) })
}

func benchmarkEncoder(b *testing.B, codec objconv.Codec) {
	for _, v := range TestValues {
		b.Run(testName(v), func(b *testing.B) {
			c := &counter{}
			e := objconv.NewEncoder(codec.NewEmitter(c))

			for i := 0; i != b.N; i++ {
				e.Encode(v)
			}

			b.SetBytes(int64(c.n / b.N))
		})
	}
}

func benchmarkDecoder(b *testing.B, codec objconv.Codec) {
	a := &bytes.Buffer{}
	a.Grow(1024)

	for _, v := range TestValues {
		e := objconv.NewEncoder(codec.NewEmitter(a))
		e.Encode(v)

		s := a.Bytes()
		r := bytes.NewReader(s)

		b.Run(testName(v), func(b *testing.B) {
			d := objconv.NewDecoder(codec.NewParser(r))

			for i := 0; i != b.N; i++ {
				var x interface{}
				d.Decode(&x)
				r.Reset(s)
			}

			b.SetBytes(int64(len(s)))
		})

		a.Reset()
	}
}

func benchmarkStreamEncoder(b *testing.B, codec objconv.Codec) {
	for _, v := range TestValues {
		b.Run(testName(v), func(b *testing.B) {
			c := &counter{}
			e := objconv.NewStreamEncoder(codec.NewEmitter(c))

			for i := 0; i != b.N; i++ {
				e.Encode(v)
			}

			e.Close()
			b.SetBytes(int64(c.n / b.N))
		})
	}
}

func benchmarkStreamDecoder(b *testing.B, codec objconv.Codec) {
	a := &bytes.Buffer{}
	a.Grow(131072)

	for _, v := range TestValues {
		e := objconv.NewStreamEncoder(codec.NewEmitter(a))
		e.Encode(v)
		e.Close()

		s := a.Bytes()
		r := bytes.NewReader(s)

		b.Run(testName(v), func(b *testing.B) {
			d := objconv.NewStreamDecoder(codec.NewParser(r))

			for i := 0; i != b.N; i++ {
				var x interface{}
				d.Decode(&x)
				r.Reset(s)
			}

			b.SetBytes(int64(len(s)))
		})

		a.Reset()
	}
}

func parseURL(s string) url.URL {
	u, _ := url.Parse(s)
	return *u
}

func parseQuery(s string) url.Values {
	v, _ := url.ParseQuery(s)
	return v
}

func parseEmail(s string) mail.Address {
	a, _ := mail.ParseAddress(s)
	return *a
}

func parseEmailList(s string) []*mail.Address {
	l, _ := mail.ParseAddressList(s)
	return l
}
