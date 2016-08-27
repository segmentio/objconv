package objconv

import (
	"bytes"
	"io"
	"unicode/utf8"
)

// EOL represent end-of-line terminaators.
type EOL string

const (
	// LF is a simple new line sequence.
	LF EOL = "\n"

	// CRLF is a two-byte new line sequence.
	CRLF EOL = "\r\n"
)

// Reader implements the io.Reader interface but reports errors through panics
// instead of returning them.
//
// It's not safe to use the reader concurrently from multiple goroutines.
type Reader struct {
	R io.Reader
	n int
	c [4]byte   // ReadByte and ReadRune buffer
	b [100]byte // ReadLine buffer
}

// NewReader returns a Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	switch v := r.(type) {
	case *Reader:
		return v
	default:
		return &Reader{R: r}
	}
}

// Read reads bytes into b from r, panics if there was an error.
//
// The method returns an error to satisfy the io.Read interface, it will always
// be nil and can be ignored.
func (r *Reader) Read(b []byte) (n int, err error) {
	n, err = r.R.Read(b)

	if n > 0 {
		err = nil
		r.n += n
	} else if err != nil {
		if err == io.EOF && r.n != 0 {
			err = io.ErrUnexpectedEOF
		}
		panic(err)
	}

	return
}

// ReadByte reads a byte from r, panics if there was an error.
//
// The method returns an error to satisfy the io.ByteReader interface, it will
// always be nil and can be ignored.
func (r *Reader) ReadByte() (c byte, err error) {
	r.Read(r.c[:1])
	c = r.c[0]
	return
}

// ReadRune reads a rune from r, panics if there was an error.
//
// The method returns an error to satisfy the io.RuneReader interface, it will
// always be nil and can be ignored.
func (r *Reader) ReadRune() (c rune, n int, err error) {
	b, _ := r.ReadByte()

	if (b & 0x80) == 0 {
		c = rune(b)
		n = 1
		return

	} else if (b & 0xE0) == 0xC0 {
		n = 2

	} else if (b & 0xF0) == 0xE0 {
		n = 3

	} else {
		n = 4
	}

	r.c[0] = b
	r.ReadFull(r.c[1:n])
	c, _ = utf8.DecodeRune(r.c[:n])
	return
}

// ReadLine reads a line ending with eol from r, panics if there was an error.
func (r *Reader) ReadLine(eol EOL) (line []byte) {
	line = r.b[:0]
	end := []byte(eol)

	for !bytes.HasSuffix(line, end) {
		b, _ := r.ReadByte()
		line = append(line, b)
	}

	return line[:len(line)-len(eol)]
}

// ReadFull fills up b with data read from r, panics if there was an error.
func (r *Reader) ReadFull(b []byte) {
	io.ReadFull(r, b)
}

// Writer implements the io.Writer interface but reports errors through panics
// instead of returning them.
//
// It's not safe to use the writer concurrently from multiple goroutines.
type Writer struct {
	W io.Writer
	b []byte   // WriteString buffer
	c [64]byte // WriteByte, WriteRune buffer, also used for short strings
}

// NewWriter returns a Writer that reads from r.
func NewWriter(w io.Writer) *Writer {
	switch x := w.(type) {
	case *Writer:
		return x
	default:
		return &Writer{W: w}
	}
}

// Write writes b to w, panics if there was an error.
func (w *Writer) Write(b []byte) (n int, err error) {
	if n, err = w.W.Write(b); err != nil {
		panic(err)
	}
	return
}

// WriteByte writes b to w, panics if there was an error.
func (w *Writer) WriteByte(b byte) (err error) {
	switch x := w.W.(type) {
	case io.ByteWriter:
		if err = x.WriteByte(b); err != nil {
			panic(err)
		}
	default:
		w.c[0] = b
		_, err = w.Write(w.c[:1])
	}
	return
}

// WriteRune writes r to w, panics if there was an error.
func (w *Writer) WriteRune(r rune) (n int, err error) {
	switch x := w.W.(type) {
	case interface {
		WriteRune(rune) (int, error)
	}:
		if n, err = x.WriteRune(r); err != nil {
			panic(err)
		}
	default:
		n = utf8.EncodeRune(w.c[:], r)
		n, err = w.Write(w.c[:n])
	}
	return
}

// WriteString writes s to w, panics if there was an error.
func (w *Writer) WriteString(s string) (n int, err error) {
	switch x := w.W.(type) {
	case interface {
		WriteString(string) (int, error)
	}:
		if n, err = x.WriteString(s); err != nil {
			panic(err)
		}

	default:
		var b []byte

		if n = len(s); n < len(w.c) {
			b = w.c[:n]
			copy(b, s)
		} else {
			if w.b == nil {
				w.b = make([]byte, 0, 512)
			}
			w.b = append(w.b[:0], s...)
			b = w.b
		}

		n, err = w.Write(b)
	}

	return
}

type counter struct{ n int }

func (c *counter) Write(b []byte) (n int, err error) {
	n = len(b)
	c.n += n
	return
}

func (c *counter) WriteByte(b byte) (err error) {
	c.n++
	return
}

func (c *counter) WriteRune(r rune) (n int, err error) {
	n = utf8.RuneLen(r)
	c.n += n
	return
}

func (c *counter) WriteString(s string) (n int, err error) {
	n = len(s)
	c.n += n
	return
}
