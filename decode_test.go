package objconv

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestDecoderDecodeType(t *testing.T) {
	date := time.Date(2016, 12, 12, 01, 01, 01, 0, time.UTC)
	err := errors.New("error")

	tests := [...]struct {
		in  interface{}
		out interface{}
	}{
		// nil -> bool
		{nil, false},

		// nil -> int
		{nil, int(0)},
		{nil, int8(0)},
		{nil, int16(0)},
		{nil, int32(0)},
		{nil, int64(0)},

		// nil -> uint
		{nil, uint(0)},
		{nil, uint8(0)},
		{nil, uint16(0)},
		{nil, uint32(0)},
		{nil, uint64(0)},
		{nil, uintptr(0)},

		// nil -> float
		{nil, float32(0)},
		{nil, float64(0)},

		// nil -> string
		{nil, ""},

		// nil -> bytes
		{nil, []byte(nil)},

		// nil -> time
		{nil, time.Time{}},

		// nil -> duration
		{nil, time.Duration(0)},

		// nil -> array
		{nil, [...]int{}},
		{nil, [...]int{0, 0, 0}},

		// nil -> slice
		{nil, []int(nil)},

		// nil -> map
		{nil, (map[int]int)(nil)},

		// nil -> struct
		{nil, struct{}{}},
		{nil, struct{ A int }{}},

		// nil -> ptr
		{nil, (*int)(nil)},

		// bool -> bool
		{false, false},
		{true, true},

		// int -> int
		{int64(1), int(1)},
		{int64(1), int8(1)},
		{int64(1), int16(1)},
		{int64(1), int32(1)},
		{int64(1), int64(1)},

		// int -> uint
		{int64(1), uint(1)},
		{int64(1), uint8(1)},
		{int64(1), uint16(1)},
		{int64(1), uint32(1)},
		{int64(1), uint64(1)},

		// int -> float
		{int64(1), float32(1)},
		{int64(1), float64(1)},

		// uint -> uint
		{uint64(1), uint(1)},
		{uint64(1), uint8(1)},
		{uint64(1), uint16(1)},
		{uint64(1), uint32(1)},
		{uint64(1), uint64(1)},
		{uint64(1), uintptr(1)},

		// uint -> int
		{uint64(1), int(1)},
		{uint64(1), int8(1)},
		{uint64(1), int16(1)},
		{uint64(1), int32(1)},
		{uint64(1), int64(1)},

		// uint -> float
		{uint64(1), float32(1)},
		{uint64(1), float64(1)},

		// float -> float
		{float64(1), float32(1)},
		{float64(1), float64(1)},

		// string -> string
		{"Hello World!", "Hello World!"},

		// string -> bytes
		{"Hello World!", []byte("Hello World!")},

		// string -> time
		{"2016-12-12T01:01:01.000Z", date},

		// string -> duration
		{"1s", time.Second},

		// string -> error
		{"error", err},

		// bytes -> bytes
		{[]byte("Hello World!"), []byte("Hello World!")},

		// bytes -> string
		{[]byte("Hello World!"), "Hello World!"},

		// bytes -> time
		{[]byte("2016-12-12T01:01:01.000Z"), date},

		// bytes -> duration
		{[]byte("1s"), time.Second},

		// bytes -> error
		{[]byte("error"), err},

		// time -> time
		{date, date},

		// duration -> duration
		{time.Second, time.Second},

		// error -> error
		{err, err},

		// array -> array
		{[...]int{}, [...]int{}},
		{[...]int{1, 2, 3}, [...]int{1, 2, 3}},

		// slice -> slice
		{[]int{}, []int{}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},

		// map -> map
		{map[int]int{}, map[int]int{}},
		{map[int]int{1: 21, 2: 42}, map[int]int{1: 21, 2: 42}},
		{map[int]map[int]int{}, map[int]map[int]int{}},
		{map[int]map[int]int{1: map[int]int{2: 3}}, map[int]map[int]int{1: map[int]int{2: 3}}},

		// map -> struct
		{map[string]int{}, struct{}{}},
		{map[string]int{"A": 42}, struct{ A int }{42}},

		// struct -> struct
		{struct{}{}, struct{}{}},
		{struct{ A int }{42}, struct{ A int }{42}},
		{struct{ A, B, C int }{1, 2, 3}, struct{ A, B, C int }{1, 2, 3}},

		// struct -> ptr
		{struct{ A int }{42}, &struct{ A int }{42}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T->%T", test.in, test.out), func(t *testing.T) {
			dec := NewDecoder(NewValueParser(test.in))
			val := reflect.New(reflect.TypeOf(test.out))

			if err := dec.Decode(val.Interface()); err != nil {
				t.Error(err)
			}

			if v := val.Elem().Interface(); !reflect.DeepEqual(v, test.out) {
				t.Errorf("%T => %#v != %v", v, v, test.out)
			}
		})
	}
}
