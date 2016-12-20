package objconv

import (
	"bytes"
	"testing"
)

func TestFormatAppendInt(t *testing.T) {
	tests := []struct {
		i int64
		b int
		s string
	}{
		{0, 10, "0"},
		{1, 10, "1"},
		{42, 16, "2a"},
	}

	buf := &bytes.Buffer{}
	buf.Grow(128)

	for _, test := range tests {
		buf.Reset()

		t.Run(test.s, func(t *testing.T) {
			f := Formatter{Out: buf}
			f.AppendInt(test.i, test.b)

			if f.Err != nil {
				t.Error(f.Err)
			}

			if s := buf.String(); s != test.s {
				t.Error(s)
			}
		})
	}
}

func TestFormatAppendUint(t *testing.T) {
	tests := []struct {
		i uint64
		b int
		s string
	}{
		{0, 10, "0"},
		{1, 10, "1"},
		{42, 16, "2a"},
	}

	buf := &bytes.Buffer{}
	buf.Grow(128)

	for _, test := range tests {
		buf.Reset()

		t.Run(test.s, func(t *testing.T) {
			f := Formatter{Out: buf}
			f.AppendUint(test.i, test.b)

			if f.Err != nil {
				t.Error(f.Err)
			}

			if s := buf.String(); s != test.s {
				t.Error(s)
			}
		})
	}
}

func TestFormatAppendByte(t *testing.T) {
	b := &bytes.Buffer{}
	f := Formatter{Out: b}
	f.AppendByte('A')

	if f.Err != nil {
		t.Error(f.Err)
	}

	if s := b.String(); s != "A" {
		t.Error(s)
	}
}

func TestFormatAppendRune(t *testing.T) {
	b := &bytes.Buffer{}
	f := Formatter{Out: b}
	f.AppendRune('好')

	if f.Err != nil {
		t.Error(f.Err)
	}

	if s := b.String(); s != "好" {
		t.Error(s)
	}
}

func TestFormatAppendString(t *testing.T) {
	b := &bytes.Buffer{}
	f := Formatter{Out: b}
	f.AppendString("你好")

	if f.Err != nil {
		t.Error(f.Err)
	}

	if s := b.String(); s != "你好" {
		t.Error(s)
	}
}
