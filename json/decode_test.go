package json

import (
	"reflect"
	"strings"
	"testing"

	"github.com/segmentio/objconv"
)

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

func TestUnicode(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{`"\u2022"`, "•"},
		{`"\uDC00D800"`, "�"},
	}

	for _, test := range tests {
		t.Run(test.out, func(t *testing.T) {
			var s string

			if err := Unmarshal([]byte(test.in), &s); err != nil {
				t.Error(err)
			}

			if s != test.out {
				t.Error(s)
			}
		})
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

func BenchmarkUnmarshal(b *testing.B) {
	for _, test := range jsonTests {
		var t reflect.Type

		if test.v == nil {
			t = reflect.TypeOf((*interface{})(nil)).Elem()
		} else {
			t = reflect.TypeOf(test.v)
		}

		v := reflect.New(t).Interface()
		s := []byte(test.s)

		b.Run(test.s, func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				if err := Unmarshal(s, v); err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(test.s)))
		})
	}
}

func BenchmarkDecoder(b *testing.B) {
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
