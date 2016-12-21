package json

import (
	"bytes"
	"errors"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/segmentio/objconv"
)

var jsonTests = []struct {
	v interface{}
	s string
}{
	{nil, `null`},
	{true, `true`},
	{false, `false`},

	{0, `0`},
	{-1, `-1`},

	{uint(1), `1`},
	{uint(42), `42`},

	{float32(0.5), `0.5`},
	{float64(1.234), `1.234`},

	{"", `""`},
	{"Hello World!", `"Hello World!"`},
	{"Hello\\World!", `"Hello\\World!"`},
	{"Hello\"World!", `"Hello\"World!"`},
	{"Hello/World!", `"Hello\/World!"`},
	{"Hello\bWorld!", `"Hello\bWorld!"`},
	{"Hello\fWorld!", `"Hello\fWorld!"`},
	{"Hello\nWorld!", `"Hello\nWorld!"`},
	{"Hello\rWorld!", `"Hello\rWorld!"`},
	{"Hello\tWorld!", `"Hello\tWorld!"`},

	{[]byte(""), `""`},
	{[]byte("Hello World!"), `"Hello World!"`},

	{errors.New("error"), `"error"`},

	{time.Date(2016, 12, 20, 0, 20, 1, 0, time.UTC), `"2016-12-20T00:20:01Z"`},
	{time.Second, `"1s"`},

	{[]int{}, `[]`},
	{[]int{1, 2, 3}, `[1,2,3]`},

	{map[string]int{}, `{}`},
	{map[string]int{"answer": 42}, `{"answer":42}`},
	{map[string]string{}, `{}`},
	{map[string]string{"hello": "world"}, `{"hello":"world"}`},
	{map[string]interface{}{}, `{}`},
	{map[string]interface{}{"hello": "world"}, `{"hello":"world"}`},

	{struct{}{}, `{}`},
	{struct {
		A int `objconv:"a"`
		B int `objconv:"-"`
		C int `objconv:",omitempty"`
		D int `objconv:",omitzero"`
		E int
	}{A: 1, E: 42}, `{"a":1,"E":42}`},
}

func TestMarshal(t *testing.T) {
	for _, test := range jsonTests {
		t.Run(test.s, func(t *testing.T) {
			b, err := Marshal(test.v)

			if err != nil {
				t.Error(err)
			}

			if s := string(b); s != test.s {
				t.Error(s)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	for _, test := range jsonTests {
		t.Run(test.s, func(t *testing.T) {
			var typ reflect.Type

			if test.v == nil {
				typ = reflect.TypeOf((*interface{})(nil)).Elem()
			} else {
				typ = reflect.TypeOf(test.v)
			}

			val := reflect.New(typ)
			err := Unmarshal([]byte(test.s), val.Interface())

			if err != nil {
				t.Error(err)
			}

			v1 := test.v
			v2 := val.Elem().Interface()

			if !reflect.DeepEqual(v1, v2) {
				t.Error(v2)
			}
		})
	}
}

func TestStreamEncoder(t *testing.T) {
	buf := &bytes.Buffer{}
	enc := NewStreamEncoder(buf)

	for i := 0; i != 10; i++ {
		if err := enc.Encode(i); err != nil {
			t.Error(err)
		}
	}

	if err := enc.Close(); err != nil {
		t.Error(err)
	}

	if s := buf.String(); s != `[0,1,2,3,4,5,6,7,8,9]` {
		t.Error(s)
	}
}

func TestStreamDecoder(t *testing.T) {
	json := `[0,1,2,3,4,5,6,7,8,9]`

	dec := NewStreamDecoder(strings.NewReader(json))
	cnt := 0
	val := 0

	for dec.Decode(&val) == nil {
		if val != cnt {
			t.Error(val, "!=", cnt)
		}
		cnt++
	}

	if cnt != 10 {
		t.Error(cnt)
	}

	if err := dec.Err(); err != nil {
		t.Error(err)
	}
}

func BenchmarkObjconvEncoder(b *testing.B) {
	e := NewEncoder(ioutil.Discard)

	for _, test := range jsonTests {
		b.Run(test.s, func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				if err := e.Encode(test.v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(test.s)))
		})
	}
}

func BenchmarkObjconvDecoder(b *testing.B) {
	r := strings.NewReader("")
	p := NewParser(nil)
	d := objconv.NewDecoder(p)

	for _, test := range jsonTests {
		var t reflect.Type

		if test.v == nil {
			t = reflect.TypeOf((*interface{})(nil)).Elem()
		} else {
			t = reflect.TypeOf(test.v)
		}

		v := reflect.New(t).Interface()

		b.Run(test.s, func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				r.Reset(test.s)
				p.Reset(r)

				if err := d.Decode(v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(test.s)))
		})
	}
}
