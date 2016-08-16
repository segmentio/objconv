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
// The Array has an unkonwn length and only one iterator which produces values
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
func (a ArraySlice) Iter() ArrayIter {
	return &arraySliceIter{
		v: reflect.Value(a),
		n: a.Len(),
	}
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
func (a ArrayValue) Iter() ArrayIter {
	return &arrayValueIter{
		v:  reflect.Value(a),
		ok: true,
	}
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

// NewArraySlice returns an ArraySlice that wraps v, assuming v is a slice or an
// array.
func NewArraySlice(v interface{}) ArraySlice {
	if a, ok := v.(ArraySlice); ok {
		return a
	}
	return ArraySlice(reflect.ValueOf(v))
}

// NewArrayValue returns an ArrayValue that wraps v.
func NewArrayValue(v interface{}) ArrayValue {
	if a, ok := v.(ArrayValue); ok {
		return a
	}
	return ArrayValue(reflect.ValueOf(v))
}
