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
				Index:     []int{0},
				Name:      "A",
				Type:      reflect.TypeOf(0),
				Offset:    0,
				Tags:      StructTag{},
				Anonymous: false,
				Exported:  true,
			},
		},

		{
			s: reflect.TypeOf(struct{ a int }{}).Field(0),
			f: StructField{
				Index:     []int{0},
				Name:      "a",
				Type:      reflect.TypeOf(0),
				Offset:    0,
				Tags:      StructTag{},
				Anonymous: false,
				Exported:  false,
			},
		},

		{
			s: reflect.TypeOf(B{}).Field(0),
			f: StructField{
				Index:     []int{0},
				Name:      "A",
				Type:      reflect.TypeOf(A{}),
				Offset:    0,
				Tags:      StructTag{},
				Anonymous: true,
				Exported:  true,
			},
		},

		{
			s: reflect.TypeOf(b{}).Field(0),
			f: StructField{
				Index:     []int{0},
				Name:      "a",
				Type:      reflect.TypeOf(a{}),
				Offset:    0,
				Tags:      StructTag{},
				Anonymous: true,
				Exported:  false,
			},
		},
	}

	for _, test := range tests {
		if f := makeStructField(test.s); !reflect.DeepEqual(test.f, f) {
			t.Errorf("%#v != %#v", test.f, f)
		}
	}
}
