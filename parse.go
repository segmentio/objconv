package objconv

import (
	"io"
	"sync"
)

// ArrayParser is an interface representing a parser of array values.
type ArrayParser interface {
	// Len returns the number of elements in the array.
	Len() int

	// Parse returns the next element of the array.
	//
	// The method is expected to return io.EOF if there was no more elements to
	// read, or any other error to indicate that something went wrong.
	Parse(r *Reader, hint interface{}) (interface{}, error)
}

// ArrayParserFunc is an adapter to allow use of ordinary functions as array parsers.
type ArrayParserFunc func(*Reader, interface{}) (interface{}, error)

// Len returns the number of values that the parser is expected to produce, which is
// -1 in the case of an ArrayParseFunc because the value is unknown.
func (f ArrayParserFunc) Len() int { return -1 }

// Parse calls f(r, hint).
func (f ArrayParserFunc) Parse(r *Reader, hint interface{}) (interface{}, error) { return f(r, hint) }

// ArrayParserLen adapts the parser p to make calls to its Len method returns n.
func ArrayParserLen(n int, p ArrayParser) ArrayParser { return &arrayParserLen{p: p, n: n} }

type arrayParserLen struct {
	p ArrayParser
	n int
	i int
}

func (a *arrayParserLen) Len() int { return a.n }

func (a *arrayParserLen) Parse(r *Reader, hint interface{}) (v interface{}, err error) {
	if a.i == a.n {
		err = io.EOF
	} else if v, err = a.p.Parse(r, hint); err == nil {
		a.i++
	}
	return
}

// NewArrayParser returns an array parser that reads values from the parser p.
//
// The returned parser array has an unknown length.
func NewArrayParser(p Parser) ArrayParser { return arrayParser{p} }

type arrayParser struct{ Parser }

func (a arrayParser) Len() int { return -1 }

// MapParser is an interface representing a parser of map values.
type MapParser interface {
	// Len returns the number of elements in the map.
	Len() int

	// ParseKey returns the next key of the array.
	//
	// The method is expected to return io.EOF if there was no more elements to
	// read, or any other error to indicate that something went wrong.
	ParseKey(r *Reader, hint interface{}) (interface{}, error)

	// ParseKey returns the next key of the array.
	//
	// The method must be called after ParseKey.
	ParseValue(r *Reader, hint interface{}) (interface{}, error)
}

// NewFixedMapParser returns an MapParser that uses p to parse values,
// expecting n entries in the map.
func NewFixedMapParser(p Parser, n int) MapParser { return &fixedMapParser{p: p, n: n} }

type fixedMapParser struct {
	p Parser
	n int
	i int
}

func (m *fixedMapParser) Len() int { return m.n }

func (m *fixedMapParser) ParseKey(r *Reader, hint interface{}) (v interface{}, err error) {
	if m.i == m.n {
		err = io.EOF
	} else {
		v, err = m.p.Parse(r, hint)
	}
	return
}

func (m *fixedMapParser) ParseValue(r *Reader, hint interface{}) (v interface{}, err error) {
	if v, err = m.p.Parse(r, hint); err == nil {
		m.i++
	}
	return
}

// The Parser interface must be implemented by types that provide decoding of a
// specific format (like json, resp, ...).
type Parser interface {
	// Parse returns the next value from r, using the hint to guess what type
	// the value will be decoded into.
	//
	// Note that hin will be nil if the value will be decoded into an
	// interface{}.
	//
	// The method must return values of types:
	// - bool
	// - int64
	// - uint64
	// - float64
	// - string
	// - []byte
	// - time.Time
	// - time.Duration
	// - error
	// - ArrayParser
	// - MapParser
	// - nil
	//
	// The method is expected to return io.EOF if there was no more elements to
	// read, or any other error to indicate that something went wrong.
	Parse(r *Reader, hint interface{}) (interface{}, error)
}

// RegisterParser adds a new parser factory under the given name.
func RegisterParser(format string, factory func() Parser) {
	parserMutex.Lock()
	parserStore[format] = factory
	parserMutex.Unlock()
}

// UnregisterParser removes the parser registered under the given name.
func UnregisterParser(format string) {
	parserMutex.Lock()
	delete(parserStore, format)
	parserMutex.Unlock()
}

// GetParser returns a new parser for the given format, or an error if no parser
// was registered for that format prior to the call.
func GetParser(format string) (p Parser, err error) {
	parserMutex.RLock()
	if f := parserStore[format]; f == nil {
		err = &UnsupportedFormatError{format}
	} else {
		p = f()
	}
	parserMutex.RUnlock()
	return
}

// NewParser returns a new parser for the given format, or panics if not parser
// was registered for that format prior to the call.
func NewParser(format string) Parser {
	if p, err := GetParser(format); err != nil {
		panic(err)
	} else {
		return p
	}
}

var (
	parserMutex sync.RWMutex
	parserStore = map[string](func() Parser){}
)
