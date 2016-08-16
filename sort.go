package objconv

import (
	"reflect"
	"sort"
)

type sortIntValues []reflect.Value

func (s sortIntValues) Len() int               { return len(s) }
func (s sortIntValues) Swap(i int, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortIntValues) Less(i int, j int) bool { return s[i].Int() < s[j].Int() }

type sortUintValues []reflect.Value

func (s sortUintValues) Len() int               { return len(s) }
func (s sortUintValues) Swap(i int, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortUintValues) Less(i int, j int) bool { return s[i].Uint() < s[j].Uint() }

type sortFloatValues []reflect.Value

func (s sortFloatValues) Len() int               { return len(s) }
func (s sortFloatValues) Swap(i int, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortFloatValues) Less(i int, j int) bool { return s[i].Float() < s[j].Float() }

type sortStringValues []reflect.Value

func (s sortStringValues) Len() int               { return len(s) }
func (s sortStringValues) Swap(i int, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortStringValues) Less(i int, j int) bool { return s[i].String() < s[j].String() }

func sortValues(typ reflect.Type, v []reflect.Value) {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sort.Sort(sortIntValues(v))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		sort.Sort(sortUintValues(v))

	case reflect.Float32, reflect.Float64:
		sort.Sort(sortFloatValues(v))

	case reflect.String:
		sort.Sort(sortStringValues(v))

	default:
		// For all other types we give up on trying to sort the values,
		// anyway it's likely not gonna be a serializable type, or something
		// that doesn't make sense.
	}
}
