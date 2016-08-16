package resp

import "testing"

func TestIsErrorPrefix(t *testing.T) {
	tests := []struct {
		s string
		x bool
	}{
		{
			s: "",
			x: false,
		},
		{
			s: "abc",
			x: false,
		},
		{
			s: "ABc",
			x: false,
		},
		{
			s: "ABC",
			x: true,
		},
	}

	for _, test := range tests {
		if x := isErrorPrefix(test.s); x != test.x {
			t.Errorf("isErrorPrefix(%#v): %v != %v", test.s, test.x, x)
		}
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		s string
		t string
		m string
	}{
		{
			s: "",
			t: "",
			m: "",
		},
		{
			s: "Hello World!",
			t: "",
			m: "Hello World!",
		},
		{
			s: "ERR",
			t: "ERR",
			m: "",
		},
		{
			s: "ERR Hello World!",
			t: "ERR",
			m: "Hello World!",
		},
	}

	for _, test := range tests {
		e := NewError(test.s)

		if e.Type != test.t {
			t.Errorf("%s: invalid error type: %#v != %#v", test.s, test.t, e.Type)
		}

		if e.Message != test.m {
			t.Errorf("%s: invalid error message: %#v != %#v", test.s, test.m, e.Message)
		}

		if s := e.Error(); s != test.s {
			t.Errorf("%s: invalid error string: %#v != %#v", test.s, test.s, s)
		}
	}
}
