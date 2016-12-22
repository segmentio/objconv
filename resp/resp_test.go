package resp

import (
	"errors"
	"time"
)

var respTests = []struct {
	v interface{}
	s string
}{
	{nil, "$-1\r\n"},

	{true, "+true\r\n"},
	{false, "+false\r\n"},

	{0, ":0\r\n"},
	{-1, ":-1\r\n"},
	{42, ":42\r\n"},

	{0.0, "+0\r\n"},
	{0.5, "+0.5\r\n"},

	{"", "+\r\n"},
	{"Hello World!", "+Hello World!\r\n"},
	{"Hello\nWorld!", "+Hello\nWorld!\r\n"},
	{"Hello\r\nWorld!", "$13\r\nHello\r\nWorld!\r\n"},

	{[]byte(nil), "$0\r\n\r\n"},
	{[]byte("Hello World!"), "$12\r\nHello World!\r\n"},

	{errors.New(""), "-\r\n"},
	{errors.New("oops"), "-oops\r\n"},
	{errors.New("A\r\nB\r\nC\r\n"), "-A\r\n"},

	{time.Date(2016, 12, 20, 0, 20, 1, 0, time.UTC), "+2016-12-20T00:20:01Z\r\n"},
	{time.Second, "+1s\r\n"},

	{[]int{}, "*0\r\n"},
	{[]int{1, 2, 3}, "*3\r\n:1\r\n:2\r\n:3\r\n"},

	{struct{}{}, "*0\r\n"},
	{struct{ A int }{42}, "*2\r\n+A\r\n:42\r\n"},
}
