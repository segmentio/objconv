package objconv

import (
	"bytes"
	"io"
	"unicode/utf8"
)

// EOL represent end-of-line terminators.
type EOL string

const (
	// LF is a simple new line sequence.
	LF EOL = "\n"

	// CRLF is a two-byte new line sequence.
	CRLF EOL = "\r\n"
)

// Reader implements the io.Reader interface.
//
// It's not safe to use the reader concurrently from multiple goroutines.
type Reader struct {
	r io.Reader
	c [utf8.UTFMax]byte // ReadRune and EOL buffer
	b []byte            // general purpose buffer
	i int               // start offset in b
	j int               // end offset in b
}

// NewReader returns a Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	switch v := r.(type) {
	case *Reader:
		return v
	default:
		return &Reader{r: r}
	}
}

// Buffered returns the bytes that are currently buffered by the readered and
// were not consumed yet.
//
// The returned by slice is valid until the next call to one of the read
// methods of the reader.
func (r *Reader) Buffered() []byte { return r.b[r.i : r.j-r.i] }

// Reset clears all buffered data in the reader, forcing it to load bytes from
// its backend reader on the next call to one of the read methods.
func (r *Reader) Reset() {
	r.i = 0
	r.j = 0
}

// Read reads bytes into b from r.
func (r *Reader) Read(b []byte) (n int, err error) {
	if r.i == r.j {
		// The reader has nothing buffered and the destination is greater than
		// the reader's internal buffer, bypass buffering and load directly into
		// b.
		if r.b != nil && len(b) >= len(r.b) {
			return r.r.Read(b)
		}

		// We need more data, buffering more!
		if _, err = r.load(); err != nil {
			return
		}
	}

	n1 := r.j - r.i
	n2 := len(b)

	if n = n1; n > n2 {
		n = n2
	}

	copy(b[:n], r.b[r.i:r.i+n])
	r.i += n
	return
}

// ReadByte reads a byte from r.
func (r *Reader) ReadByte() (c byte, err error) {
	if r.i == r.j {
		if _, err = r.load(); err != nil {
			return
		}
	}
	c = r.b[r.i]
	r.i++
	return
}

// ReadRune reads a rune from r.
func (r *Reader) ReadRune() (c rune, n int, err error) {
	var b byte

	if b, err = r.ReadByte(); err != nil {
		return
	}

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

	if _, err = r.ReadFull(r.c[1:n]); err != nil {
		return
	}

	c, _ = utf8.DecodeRune(r.c[:n])
	return
}

// ReadLine reads a line ending with eol from r.
//
// The returned byte slice points to the reader's internal buffer and is valid
// until the next call to one of the reader functions.
func (r *Reader) ReadLine(eol EOL) (line []byte, err error) {
	suff := append(r.c[:0], eol...)

	for {
		if off := bytes.Index(r.b[r.i:r.j], suff); off >= 0 {
			line = r.b[r.i : r.i+off]
			r.i += off + len(suff)
			break
		}
		// There was no EOL in the current buffer, we need to load more data
		// on the next iteration.
		if _, err = r.load(); err != nil {
			return
		}
	}

	return
}

// ReadFull fills up b with data read from r.
func (r *Reader) ReadFull(b []byte) (n int, err error) { return io.ReadFull(r, b) }

func (r *Reader) load() (n int, err error) {
	if r.b == nil {
		// Lazy allocation of the reader's internal buffer.
		r.b = make([]byte, 1024)
	} else if r.j >= (len(r.b) / 2) {
		// Double the size of the reader's internal buffer and copy any bytes
		// that may still be in the old buffer.
		b := make([]byte, 2*len(r.b))
		copy(b, r.b[r.i:r.j])
		r.j -= r.i
		r.i = 0
		r.b = b
	}

	if n, err = r.r.Read(r.b[r.j:]); err != nil && n != 0 {
		err = nil
	}

	r.j += n
	return
}

// Writer implements the io.Writer interface.
//
// The writer is used to optimize output operations throughout the objconv
// package and sub-packages. It wraps around a generic io.Writer and capture
// the errors it would return, allowing the encoder to detect whether an
// error occurred or not after serializing values. It makes it easier to write
// emitters since they don't need to do error checking in their implementations.
//
// It's not safe to use the writer concurrently from multiple goroutines.
type Writer struct {
	w io.Writer
	e error
	b [64]byte // buffer
}

// NewWriter returns a Writer that reads from r.
func NewWriter(w io.Writer) *Writer {
	switch x := w.(type) {
	case *Writer:
		return x
	default:
		return &Writer{w: w}
	}
}

// Write writes b to w.
func (w *Writer) Write(b []byte) (n int, err error) {
	if err = w.e; err == nil {
		n, err = w.w.Write(b)
		w.e = err
	}
	return
}

// WriteByte writes b to w.
func (w *Writer) WriteByte(b byte) (err error) {
	w.b[0] = b
	_, err = w.Write(w.b[:1])
	return
}

// WriteRune writes r to w.
func (w *Writer) WriteRune(r rune) (n int, err error) {
	return w.Write(w.b[:utf8.EncodeRune(w.b[:], r)])
}

// WriteString writes s to w.
func (w *Writer) WriteString(s string) (n int, err error) {
	// Writes the string by chunks of len(w.b) at a time.
	for len(s) != 0 {
		n1 := len(w.b)
		n2 := len(s)
		n3 := n1
		n4 := 0

		if n3 > n2 {
			n3 = n2
		}

		copy(w.b[:], s[:n3])
		n4, err = w.Write(w.b[:n3])
		n += n4
		s = s[n3:]
	}
	return
}

type runeWriter interface {
	WriteRune(rune) (int, error)
}

type stringWriter interface {
	WriteString(string) (int, error)
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
