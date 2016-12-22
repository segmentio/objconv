package objconv

import "testing"

var parseIntTests = []struct {
	v int64
	s string
}{
	{0, "0"},
	{1, "1"},
	{-1, "-1"},
	{1234567890, "1234567890"},
	{-1234567890, "-1234567890"},
}

func TestParseInt(t *testing.T) {
	for _, test := range parseIntTests {
		t.Run(test.s, func(t *testing.T) {
			v, err := ParseInt([]byte(test.s))

			if err != nil {
				t.Error(err)
			}

			if v != test.v {
				t.Error(v)
			}
		})
	}
}

func BenchmarkParseInt(b *testing.B) {
	for _, test := range parseIntTests {
		b.Run(test.s, func(b *testing.B) {
			s := []byte(test.s)

			for i := 0; i != b.N; i++ {
				ParseInt(s)
			}
		})
	}
}

var parseUintHexTests = []struct {
	v uint64
	s string
}{
	{0x0, "0"},
	{0x1, "1"},
	{0xA, "a"},
	{0xA, "A"},
	{0x10, "10"},
	{0xABCDEF, "abcdef"},
	{0xABCDEF, "ABCDEF"},
}

func TestParseUintHex(t *testing.T) {
	for _, test := range parseUintHexTests {
		t.Run(test.s, func(t *testing.T) {
			v, err := ParseUintHex([]byte(test.s))

			if err != nil {
				t.Error(err)
			}

			if v != test.v {
				t.Error(v)
			}
		})
	}
}

func BenchmarkParseUintHex(b *testing.B) {
	for _, test := range parseUintHexTests {
		b.Run(test.s, func(b *testing.B) {
			s := []byte(test.s)

			for i := 0; i != b.N; i++ {
				ParseUintHex(s)
			}
		})
	}
}
