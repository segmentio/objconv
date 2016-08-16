package objconv

import "reflect"

// Map is an interface representing iterable sequences of key/value pairs.
type Map interface {
	// Len rerturns the number of entries in the map.
	Len() int

	// Iter returns an iterator pointing to the first entry of the map.
	Iter() MapIter
}

// MapIter is an interface allowing the application to iterator over the entries
// of a map.
type MapIter interface {
	// Next returns the current map item, and advances the iterator.
	// The boolean will be false if there was no more values to return.
	Next() (MapItem, bool)
}

// MapItem represent a single mapping between a key and a value of arbitrary
// types.
type MapItem struct {
	Key   interface{}
	Value interface{}
}

// MapFunc is a type alias for a function that implements the Map interface.
//
// The Map has an unkonwn length and only one iterator which produces values
// returned by successive calls to the function.
type MapFunc func() (MapItem, bool)

// Len returns -1 because the map length is unknown.
func (f MapFunc) Len() int { return -1 }

// Iter returns the MapFunc itself which is its own iterator.
func (f MapFunc) Iter() MapIter { return f }

// Next calls the function, effectively advancing the iterator to the next
// element.
func (f MapFunc) Next() (MapItem, bool) { return f() }

// MapMap is type alias for a reflect.Value of type map that implements the Map
// interface.
type MapMap reflect.Value

// Len returns the number of elemnts in the map.
func (m MapMap) Len() int { return reflect.Value(m).Len() }

// Iter returns an iterator pointing to the first element of the map.
func (m MapMap) Iter() MapIter {
	v := reflect.Value(m)
	k := v.MapKeys()
	sortValues(v.Type().Key(), k)
	return &mapValueIter{
		k: k,
		m: v,
		n: v.Len(),
	}
}

type mapValueIter struct {
	k []reflect.Value
	m reflect.Value
	n int
	i int
}

func (it *mapValueIter) Next() (item MapItem, ok bool) {
	if ok = it.i < it.n; ok {
		item = MapItem{
			Key:   it.k[it.i].Interface(),
			Value: it.m.MapIndex(it.k[it.i]).Interface(),
		}
		it.i++
	}
	return
}

// MapSlice is a representation of a mapping as a slice of MapItem.
type MapSlice []MapItem

// Len returns the number of elemnts in the map.
func (m MapSlice) Len() int { return len(m) }

// Iter returns an iterator pointing to the first element of the map.
func (m MapSlice) Iter() MapIter {
	return &mapSliceIter{
		m: m,
		n: len(m),
	}
}

type mapSliceIter struct {
	m MapSlice
	n int
	i int
}

func (it *mapSliceIter) Next() (item MapItem, ok bool) {
	if ok = it.i < it.n; ok {
		item = it.m[it.i]
		it.i++
	}
	return
}

// MapStruct is a representation of a map obtained from a struct value where the
// keys are the names of the fields (after applying tag modifiers) and the
// values the values of the fields in the struct.
type MapStruct MapSlice

// NewMapStruct returns a new MapStruct value constructed from v which must be a
// struct, using t as the tag name to apply struct field modifiers from the
// struct field tags.
func NewMapStruct(t string, v interface{}) MapStruct {
	return newMapStruct(t, reflect.ValueOf(v))
}

func newMapStruct(t string, v reflect.Value) MapStruct {
	s := LookupStruct(v.Type())
	m := make(MapStruct, 0, len(s.Fields))
	it := s.IterValue(t, v, FilterUnexported|FilterAnonymous|FilterSkipped|FilterOmitempty)

	for {
		if k, v, ok := it.Next(); !ok {
			break
		} else {
			m = append(m, MapItem{Key: k, Value: v})
		}
	}

	return m
}

// Len returns the number of elemnts in the map.
func (m MapStruct) Len() int { return MapSlice(m).Len() }

// Iter returns an iterator pointing to the first element of the map.
func (m MapStruct) Iter() MapIter { return MapSlice(m).Iter() }
