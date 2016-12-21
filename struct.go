package objconv

import (
	"reflect"
	"sync"
)

// StructField represents a single field of a struct and carries information
// useful to the algorithms of the objconv package.
type StructField struct {
	// The index of the field in the structure.
	Index []int

	// The name of the field in the structure.
	Name string

	// Omitempty is set to true when the field should be omitted if it has an
	// empty value.
	Omitempty bool

	// Omitzero is set to true when the field should be omitted if it has a zero
	// value.
	Omitzero bool

	// cache for the encoder and decoder methods
	encode encodeFunc
	decode decodeFunc
}

func makeStructField(f reflect.StructField, c map[reflect.Type]*Struct) StructField {
	t := ParseTag(f.Tag.Get("objconv"))
	s := StructField{
		Index:     f.Index,
		Name:      f.Name,
		Omitempty: t.Omitempty,
		Omitzero:  t.Omitzero,

		encode: makeEncodeFunc(f.Type, encodeFuncOpts{
			recurse: true,
			structs: c,
		}),

		decode: makeDecodeFunc(f.Type, decodeFuncOpts{
			recurse: true,
			structs: c,
		}),
	}

	if len(t.Name) != 0 {
		s.Name = t.Name
	}

	return s
}

func (f *StructField) omit(v reflect.Value) bool {
	return (f.Omitempty && isEmptyValue(v)) || (f.Omitzero && isZeroValue(v))
}

// Struct is used to represent a Go structure in internal data structures that
// cache meta information to make field lookups faster and avoid having to use
// reflection to lookup the same type information over and over again.
type Struct struct {
	Fields       []StructField           // the serializable fields of the struct
	FieldsByName map[string]*StructField // cache of fields by name
}

// LookupStruct behaves like MakeStruct but uses a global cache to avoid having
// to recreate the struct values when not needed.
//
// As much as possible you should be using this function instead of calling
// MakeStruct or maintaining your own cache so the program can efficiently make
// use of the cache and avoid storing duplicate information in different parts
// of the program.
func LookupStruct(t reflect.Type) *Struct { return structCache.Lookup(t) }

// NewStruct takes a Go type as argument and extract information to make a new
// Struct value.
// The type has to be a struct type or a panic will be raised.
func NewStruct(t reflect.Type) *Struct {
	return newStruct(t, map[reflect.Type]*Struct{})
}

func newStruct(t reflect.Type, c map[reflect.Type]*Struct) *Struct {
	if s := c[t]; s != nil {
		return s
	}

	n := t.NumField()
	s := &Struct{
		Fields:       make([]StructField, 0, n),
		FieldsByName: make(map[string]*StructField),
	}
	c[t] = s

	for i := 0; i != n; i++ {
		ft := t.Field(i)

		if ft.Anonymous || len(ft.PkgPath) != 0 { // anonymous or non-exported
			continue
		}

		sf := makeStructField(ft, c)

		if sf.Name == "-" { // skip
			continue
		}

		s.Fields = append(s.Fields, sf)
		s.FieldsByName[sf.Name] = &s.Fields[len(s.Fields)-1]
	}

	return s
}

// StructCache is a simple cache for mapping Go types to Struct values.
type StructCache struct {
	mutex sync.RWMutex
	store map[reflect.Type]*Struct
}

// NewStructCache creates and returns a new StructCache value.
func NewStructCache() *StructCache {
	return &StructCache{store: make(map[reflect.Type]*Struct, 20)}
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
	structCache = NewStructCache()
)
