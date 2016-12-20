package objconv

import (
	"io"
	"strconv"
	"unicode/utf8"
)

// Formatter is a utility type that makes it easier to build clear and efficient
// emitters.
//
// The type exposes methods to efficiently write values to an io.Writer. Using a
// Formatter is usually faster than using routines like io.WriteString or the
// functions from strconv because it manages an internal buffer to avoid memory
// allocations.
type Formatter struct {
	W io.Writer // the writer to output formatted values to
	b [64]byte  // the internal buffer used for formating values
	e error     // cache of errors returned by the writer
}

// Err returns the error caught by the formatter when writing bytes to its
// output, or nil if everything went fine.
func (f *Formatter) Err() error {
	return f.e
}

// AppendInt writes i in base to the formatter's output.
func (f *Formatter) AppendInt(i int64, base int) {
	if f.e == nil {
		_, f.e = f.Write(strconv.AppendInt(f.b[:0], i, base))
	}
}

// AppendUint writes u in base to the formatter's output.
func (f *Formatter) AppendUint(u uint64, base int) {
	if f.e == nil {
		_, f.e = f.Write(strconv.AppendUint(f.b[:0], u, base))
	}
}

// AppendByte writes b to the formatter's output.
func (f *Formatter) AppendByte(b byte) {
	if f.e == nil {
		f.e = f.WriteByte(b)
	}
}

// AppendRune writes r to the formatter's output.
func (f *Formatter) AppendRune(r rune) {
	if f.e == nil {
		_, f.e = f.WriteRune(r)
	}
}

// AppendString writes s to the formatter's output.
func (f *Formatter) AppendString(s string) {
	if f.e == nil {
		_, f.e = f.WriteString(s)
	}
}

// WriteByte writes b to the formatter's output.
//
// This method doesn't update the formatter's internal error cache.
func (f *Formatter) WriteByte(b byte) (err error) {
	_, err = f.Write(append(f.b[:0], b))
	return
}

// WriteRune writes r to the formatter's output.
//
// This method doesn't update the formatter's internal error cache.
func (f *Formatter) WriteRune(r rune) (n int, err error) {
	n = utf8.EncodeRune(f.b[:], r)
	return f.Write(f.b[:n])
}

// WriteString writes s to the formatter's output.
//
// This method doesn't update the formatter's internal error cache.
func (f *Formatter) WriteString(s string) (n int, err error) {
	for err == nil && len(s) != 0 {
		n1 := len(s)
		n2 := len(f.b)
		n3 := 0

		if n1 > n2 {
			n1 = n2
		}

		copy(f.b[:n1], s[:n1])

		n3, err = f.Write(f.b[:n1])
		n += n3
		s = s[n1:]
	}
	return
}

// Write writes b to the formatter's output.
//
// This method doesn't update the formatter's internal error cache.
func (f *Formatter) Write(b []byte) (n int, err error) {
	return f.W.Write(b)
}
