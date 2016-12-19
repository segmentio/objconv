package objconv

import (
	"reflect"
	"testing"
)

func TestMakeStructField(t *testing.T) {
	type A struct{ A int }
	type a struct{ a int }
	type B struct{ A }
	type b struct{ a }

	tests := []struct {
		s reflect.StructField
		f StructField
	}{
		{
			s: reflect.TypeOf(A{}).Field(0),
			f: StructField{
				Index: []int{0},
				Name:  "A",
			},
		},

		{
			s: reflect.TypeOf(a{}).Field(0),
			f: StructField{
				Index: []int{0},
				Name:  "a",
			},
		},

		{
			s: reflect.TypeOf(B{}).Field(0),
			f: StructField{
				Index: []int{0},
				Name:  "A",
			},
		},

		{
			s: reflect.TypeOf(b{}).Field(0),
			f: StructField{
				Index: []int{0},
				Name:  "a",
			},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			f := makeStructField(test.s)
			f.decode = nil // function types are not comparable
			f.encode = nil

			if !reflect.DeepEqual(test.f, f) {
				t.Errorf("%#v != %#v", test.f, f)
			}
		})
	}
}
