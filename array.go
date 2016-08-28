package objconv

import "reflect"

// Array is an interface representing iterable sequences of values.
type Array interface {
	// Len returns the number of elements in the array.
	Len() int

	// Iter returns an iterator pointing to the first element of the array.
	Iter() ArrayIter
}

// ArrayIter is an interface allowing the application to iterate over the
// elements of an array.
type ArrayIter interface {
	// Next returns the current value, and advances the iterator. The boolean
	// will be false if there was no more values to return.
	Next() (interface{}, bool)
}

// ArrayFunc is a type alias for a function that implements the Array interface.
//
// The Array has an unknown length and only one iterator which produces values
// returned by successive calls to the function.
type ArrayFunc func() (interface{}, bool)

// Len returns -1 because the array length is unknown.
func (f ArrayFunc) Len() int { return -1 }

// Iter returns the ArrayFunc itself which is its own iterator.
func (f ArrayFunc) Iter() ArrayIter { return f }

// Next calls the function, effectively advancing the iterator to the next
// element.
func (f ArrayFunc) Next() (interface{}, bool) { return f() }

// ArraySlice is a type alias for a reflect.Value of type slice or array that
// implements the Array interface.
type ArraySlice reflect.Value

// Len returns the number of elements in the array.
func (a ArraySlice) Len() int { return reflect.Value(a).Len() }

// Iter returns an iterator pointing to the first element of the array.
func (a ArraySlice) Iter() ArrayIter { return &arraySliceIter{v: reflect.Value(a), n: a.Len()} }

// NewArraySlice returns an ArraySlice that wraps v, assuming v is a slice or an
// array.
func NewArraySlice(v interface{}) ArraySlice {
	if a, ok := v.(ArraySlice); ok {
		return a
	}
	return ArraySlice(reflect.ValueOf(v))
}

type arraySliceIter struct {
	v reflect.Value
	n int
	i int
}

func (it *arraySliceIter) Next() (v interface{}, ok bool) {
	if ok = it.i < it.n; ok {
		v = it.v.Index(it.i).Interface()
		it.i++
	}
	return
}

// ArrayValue is a type alias for a reflect.Value of any type. The type
// implements the Array interface and represents an array returning a
// single value.
type ArrayValue reflect.Value

// Len returns the number of elements in the array, which is 1.
func (a ArrayValue) Len() int { return 1 }

// Iter returns an iterator pointing to the first element of the array.
func (a ArrayValue) Iter() ArrayIter { return &arrayValueIter{v: reflect.Value(a), ok: true} }

// NewArrayValue returns an ArrayValue that wraps v.
func NewArrayValue(v interface{}) ArrayValue {
	if a, ok := v.(ArrayValue); ok {
		return a
	}
	return ArrayValue(reflect.ValueOf(v))
}

type arrayValueIter struct {
	v  reflect.Value
	ok bool
}

func (it *arrayValueIter) Next() (v interface{}, ok bool) {
	if ok, it.ok = it.ok, false; ok {
		v = it.v.Interface()
	}
	return
}

// ArrayStream is an adapter for a StreamDecoder that implements the Array
// interface.
//
// The array has a single iterator and its length is the initial number of
// values that can be read from the stream.
type ArrayStream struct {
	d StreamDecoder
	n int
}

// Len returns the number of elements in the array.
func (a ArrayStream) Len() int { return a.n }

// Iter returns an iterator pointing to the first element of the array.
func (a ArrayStream) Iter() ArrayIter { return arrayStreamIter{a.d} }

// NewArrayStream returns an ArrayStream that adapts the StreamDecoder d.
func NewArrayStream(d StreamDecoder) ArrayStream { return ArrayStream{d, d.Len()} }

type arrayStreamIter struct{ StreamDecoder }

func (it arrayStreamIter) Next() (v interface{}, ok bool) {
	ok = it.Decode(&v) == nil
	return
}

// ArrayLen returns an Array value that wraps around the array a, but sets the
// length to n.
func ArrayLen(n int, a Array) Array { return arrayLen{a, n} }

type arrayLen struct {
	Array
	n int
}

func (a arrayLen) Len() int { return a.n }

// NewArray returns an Array that exposes the values of v.
func NewArray(v interface{}) Array {
	switch a := v.(type) {
	case Array:
		return a

	case func() (interface{}, bool):
		return ArrayFunc(a)
	}

	switch a := reflect.ValueOf(v); a.Kind() {
	case reflect.Slice, reflect.Array:
		return ArraySlice(a)

	default:
		return ArrayValue(a)
	}
}

// MultiArray combines multiple arrays into one.
func MultiArray(a ...Array) Array { return multiArray(a) }

type multiArray []Array

func (m multiArray) Len() (n int) {
	for _, a := range m {
		n += a.Len()
	}
	return
}

func (m multiArray) Iter() ArrayIter {
	it := make([]ArrayIter, len(m))
	for i, a := range m {
		it[i] = a.Iter()
	}
	return &multiArrayIter{it}
}

type multiArrayIter struct {
	it []ArrayIter
}

func (m *multiArrayIter) Next() (v interface{}, ok bool) {
	for len(m.it) != 0 {
		if v, ok = m.it[0].Next(); ok {
			break
		}
		m.it = m.it[1:]
	}
	return
}
