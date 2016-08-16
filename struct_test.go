package objconv

import (
	"reflect"
	"testing"
)

func TestStructIter(t *testing.T) {
	type T struct {
		A int `test:"a"`
		B int
		C int `test:",omitempty"`
		D int `test:"-"`
		e int
	}

	v := T{}
	s := LookupStruct(reflect.TypeOf(v))
	i := 0

	for it := s.Iter("test", v, FilterUnexported|FilterAnonymous|FilterOmitempty|FilterSkipped); true; {
		if name, val, ok := it.Next(); !ok {
			break
		} else {
			switch i {
			case 0:
				if name != "a" {
					t.Error("invalid field returned by struct iterator:", name, val)
				}

			case 1:
				if name != "B" {
					t.Error("invalid field returned by struct iterator:", name, val)
				}

			default:
				t.Error("too many value returned by the struct iterator")
			}
			i++
		}
	}
}

func TestStructSetter(t *testing.T) {
	type T struct {
		A int `test:"a"`
		B int
		C int `test:",omitempty"`
		D int `test:"-"`
		e int
	}

	v := T{}
	s := LookupStruct(reflect.TypeOf(v)).Setter("test", &v)

	s["a"].SetInt(42)

	if v.A != 42 {
		t.Error("struct setter did not modify the original value")
	}
}
