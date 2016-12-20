package json

import (
	"errors"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	date := time.Date(2016, 12, 20, 0, 20, 1, 0, time.UTC)

	tests := []struct {
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

		{date, `"2016-12-20T00:20:01Z"`},
		{time.Second, `"1s"`},

		{[]int{}, `[]`},
		{[]int{1, 2, 3}, `[1,2,3]`},

		{map[string]int{}, `{}`},
		{map[string]int{"answer": 42}, `{"answer":42}`},

		{struct{}{}, `{}`},
		{struct {
			A int `objconv:"a"`
			B int `objconv:"-"`
			C int `objconv:",omitempty"`
			D int `objconv:",omitzero"`
			E int
		}{A: 1, E: 42}, `{"a":1,"E":42}`},
	}

	for _, test := range tests {
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
