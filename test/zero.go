package test

import "reflect"

// NewZero creates a zero-value of the same type as v.
func NewZero(v interface{}) reflect.Value {
	if v == nil {
		var x interface{}
		return reflect.ValueOf(&x)
	}
	return reflect.New(reflect.TypeOf(v))
}
