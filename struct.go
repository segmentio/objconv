package objconv

import (
	"reflect"
	"sync"
)

// Struct is used to represent a Go structure in internal data structures that
// cache meta information to make field lookups faster and avoid having to use
// reflection to lookup the same type information over and over again.
type Struct struct {
	// Type of the structure represented by the Struct value.
	Type reflect.Type

	// Fields is the list of field in the structure, represented as StructField
	// values.
	Fields []StructField
}

// LookupStruct behaves like MakeStruct but uses a global cache to avoid having
// to recreate the struct values when not needed.
//
// As much as possible you should be using this function instead of calling
// MakeStruct or maintaining your own cache so the program can efficiently make
// use of the cache and avoid storing duplicate information in different parts
// of the program.
func LookupStruct(t reflect.Type) *Struct {
	return structCache.Lookup(t)
}

// NewStruct takes a Go type as argument and extract information to make a new
// Struct value.
// The type has to be a struct type or a panic will be raised.
func NewStruct(t reflect.Type) *Struct {
	n := t.NumField()
	s := make([]StructField, 0, n)

	for i := 0; i != n; i++ {
		s = append(s, makeStructField(t.Field(i)))
	}

	return &Struct{
		Type:   t,
		Fields: s,
	}
}

// The StructSetter type is a mapping from field names as expected in their
// serialized form to the fields reflect.Value.
type StructSetter map[string]reflect.Value

// Setter returns a new struct setter that allows the program to lookup field
// values in the struct given as second argument based on their serialized name.
// The value must be a non-nil pointer to struct type or the function will
// panic.
func (s *Struct) Setter(tag string, val interface{}) StructSetter {
	return s.SetterValue(tag, reflect.ValueOf(val).Elem())
}

// SetterValue is like Setter but takes a reflect.Value instead of an interface.
func (s *Struct) SetterValue(tag string, val reflect.Value) StructSetter {
	st := make(StructSetter, len(s.Fields))
	it := s.IterValue(tag, val, FilterUnexported|FilterAnonymous|FilterSkipped)

	for {
		if name, val, ok := it.NextValue(); !ok {
			break
		} else {
			st[name] = val
		}
	}

	return st
}

// Iter returns an iterator that produces values from the fields in `val` and
// considers tags under the `tag` name to apply filtering to the results.
func (s *Struct) Iter(tag string, val interface{}, filters StructIterFilters) *StructIter {
	return s.IterValue(tag, reflect.ValueOf(val), filters)
}

// IterValue is like Iter but takes a reflect.Value instead of an interface.
func (s *Struct) IterValue(tag string, val reflect.Value, filters StructIterFilters) *StructIter {
	return &StructIter{
		value:   val,
		fields:  s.Fields,
		tag:     tag,
		index:   0,
		count:   len(s.Fields),
		filters: filters,
	}
}

// StructIterFilters are constants representing which filters a struct iterator
// should apply when
type StructIterFilters int

const (
	// FilterUnexported filters the unexported fields.
	FilterUnexported StructIterFilters = 1 << iota

	// FilterAnonymous filters the anonymous fields.
	FilterAnonymous

	// FilterSkipped filters fields with a tag name set to "-".
	FilterSkipped

	// FilterOmitempty filters empty fields that are marked with omitempty.
	FilterOmitempty
)

// StructIter is an iterator that allows the program to easily iterate over the
// fields of a struct value.
type StructIter struct {
	value   reflect.Value
	fields  []StructField
	tag     string
	index   int
	count   int
	filters StructIterFilters
}

// Next returns the name and value of the next field seen by the iterator,
// setting `ok` to true or false to indicated if a field was found or if the
// iterator reached the end.
func (it *StructIter) Next() (name string, val interface{}, ok bool) {
	var v reflect.Value
	if name, v, ok = it.NextValue(); v.IsValid() {
		val = v.Interface()
	}
	return
}

// NextValue is like Next but returns a reflect.Value instead of an empty
// interface.
func (it *StructIter) NextValue() (name string, val reflect.Value, ok bool) {
	for ; !ok && it.index != it.count; it.index++ {
		f := &it.fields[it.index]
		name = f.Name

		if (it.filters&FilterUnexported) != 0 && !f.Exported {
			continue
		}

		if (it.filters&FilterAnonymous) != 0 && f.Anonymous {
			continue
		}

		v := it.value.FieldByIndex(f.Index)

		if t, ok := f.Tags[it.tag]; ok {
			if (it.filters&FilterSkipped) != 0 && t.Skip {
				continue
			}

			if (it.filters&FilterOmitempty) != 0 && (t.Omitempty && isEmptyValue(v)) {
				continue
			}

			if len(t.Name) != 0 {
				name = t.Name
			}
		}

		val, ok = v, true
	}
	return
}

// StructCache is a simple cache for mapping Go types to Struct values.
type StructCache struct {
	mutex sync.RWMutex
	store map[reflect.Type]*Struct
}

// NewStructCache creates and returns a new StructCache value.
func NewStructCache() *StructCache {
	return &StructCache{store: make(map[reflect.Type]*Struct)}
}

// Lookup takes a Go type as argument and returns the matching Struct value,
// potentially creating it if it didn't already exist.
// This method is safe to call from multiple goroutines.
func (cache *StructCache) Lookup(t reflect.Type) (s *Struct) {
	cache.mutex.RLock()
	s = cache.store[t]
	cache.mutex.RUnlock()

	if s == nil {
		// There's a race confition here where this value may be generated
		// multiple times.
		// The impact in practice is really small as it's unlikely to happen
		// often, we take the appraoch of keeping the logic simple and avoid
		// a more complex synchronization logic required to solve this edge
		// case.
		s = NewStruct(t)
		cache.mutex.Lock()
		cache.store[t] = s
		cache.mutex.Unlock()
	}

	return
}

var (
	// This struct cache is used to avoid reusing reflection over and over when
	// the objconv functions are called. The performance improvements on iterating
	// over struct fields are huge, this is a really important optimization:
	//
	// benchmark                                   old ns/op     new ns/op     delta
	// BenchmarkLengthStructZero                   53.9          99.9          +85.34%
	// BenchmarkLengthStructNonZero                746           411           -44.91%
	// BenchmarkLengthStructOmitEmptyZero          779           174           -77.66%
	// BenchmarkLengthStructOmpitemptytNonZero     1119          425           -62.02%
	//
	// Note: Disregard the performance loss on the `StructZero` benchmark, this
	// is testing an empty struct with no field, which is just a baseline and not
	// actually useful in real-world use cases.
	//
	structCache = NewStructCache()
)
