package objconv

import (
	"reflect"
	"testing"
)

func TestParseTag(t *testing.T) {
	tests := []struct {
		tag string
		res Tag
	}{
		{
			tag: "",
			res: Tag{},
		},
		{
			tag: "hello",
			res: Tag{Name: "hello"},
		},
		{
			tag: ",omitempty",
			res: Tag{Omitempty: true},
		},
		{
			tag: "-",
			res: Tag{Name: "-", Skip: true},
		},
		{
			tag: "hello,omitempty",
			res: Tag{Name: "hello", Omitempty: true},
		},
		{
			tag: "-,omitempty",
			res: Tag{Name: "-", Omitempty: true, Skip: true},
		},
	}

	for _, test := range tests {
		if res := ParseTag(test.tag); res != test.res {
			t.Errorf("%s: %#v != %#v", test.tag, test.res, res)
		}
	}
}

func TestParseStructTag(t *testing.T) {
	tests := []struct {
		val interface{}
		res StructTag
	}{
		{
			val: struct{ F int }{},
			res: StructTag{},
		},
		{
			val: struct {
				F int `json:"f"`
			}{},
			res: StructTag{"json": Tag{Name: "f"}},
		},
		{
			val: struct {
				F int `json:"-"`
			}{},
			res: StructTag{"json": Tag{Name: "-", Skip: true}},
		},
		{
			val: struct {
				F int `json:",omitempty` // missing closing "
			}{},
			res: StructTag{"json": Tag{Omitempty: true}},
		},
		{
			val: struct {
				F int `json:f` // missing opening "
			}{},
			res: StructTag{"json": Tag{}, "f": Tag{}},
		},
	}

	for _, test := range tests {
		if res := ParseStructTag(reflect.TypeOf(test.val).Field(0).Tag); !reflect.DeepEqual(res, test.res) {
			t.Errorf("%s: %#v != %#v", test.val, test.res, res)
		}
	}
}
