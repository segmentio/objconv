package objconv

import (
	"fmt"
	"reflect"
)

const (
	// UintMax is the maximum value of a uint.
	UintMax = ^uint(0)

	// UintMin is the minimum value of a uint.
	UintMin = 0

	// Uint8Max is the maximum value of a uint8.
	Uint8Max = 255

	// Uint8Min is the minimum value of a uint8.
	Uint8Min = 0

	// Uint16Max is the maximum value of a uint16.
	Uint16Max = 65535

	// Uint16Min is the minimum value of a uint16.
	Uint16Min = 0

	// Uint32Max is the maximum value of a uint32.
	Uint32Max = 4294967295

	// Uint32Min is the minimum value of a uint32.
	Uint32Min = 0

	// Uint64Max is the maximum value of a uint64.
	Uint64Max = 18446744073709551615

	// Uint64Min is the minimum value of a uint64.
	Uint64Min = 0

	// UintptrMax is the maximum value of a uintptr.
	UintptrMax = ^uintptr(0)

	// UintptrMin is the minimum value of a uintptr.
	UintptrMin = 0

	// IntMax is the maximum value of a int.
	IntMax = int(UintMax >> 1)

	// IntMin is the minimum value of a int.
	IntMin = -IntMax - 1

	// Int8Max is the maximum value of a int8.
	Int8Max = 127

	// Int8Min is the minimum value of a int8.
	Int8Min = -128

	// Int16Max is the maximum value of a int16.
	Int16Max = 32767

	// Int16Min is the minimum value of a int16.
	Int16Min = -32768

	// Int32Max is the maximum value of a int32.
	Int32Max = 2147483647

	// Int32Min is the minimum value of a int32.
	Int32Min = -2147483648

	// Int64Max is the maximum value of a int64.
	Int64Max = 9223372036854775807

	// Int64Min is the minimum value of a int64.
	Int64Min = -9223372036854775808

	// Float32IntMax is the maximum consecutive integer value representable by a float32.
	Float32IntMax = 16777216

	// Float32IntMin is the minimum consecutive integer value representable by a float32.
	Float32IntMin = -16777216

	// Float64IntMax is the maximum consecutive integer value representable by a float64.
	Float64IntMax = 9007199254740992

	// Float64IntMin is the minimum consecutive integer value representable by a float64.
	Float64IntMin = -9007199254740992
)

func checkUint64Bounds(v uint64, max uint64, res interface{}) {
	if v > max {
		panic(&OutOfBoundsError{
			Value: v,
			Type:  reflect.TypeOf(res),
		})
	}
}

func convertUint64ToInt(v uint64) (res int) {
	checkUint64Bounds(v, uint64(IntMax), res)
	return int(v)
}

func convertUint64ToInt8(v uint64) (res int8) {
	checkUint64Bounds(v, Int8Max, res)
	return int8(v)
}

func convertUint64ToInt16(v uint64) (res int16) {
	checkUint64Bounds(v, Int16Max, res)
	return int16(v)
}

func convertUint64ToInt32(v uint64) (res int32) {
	checkUint64Bounds(v, Int32Max, res)
	return int32(v)
}

func convertUint64ToInt64(v uint64) (res int64) {
	checkUint64Bounds(v, Uint64Max, res)
	return int64(v)
}

func convertUint64ToUint(v uint64) (res uint) {
	checkUint64Bounds(v, uint64(UintMax), res)
	return uint(v)
}

func convertUint64ToUint8(v uint64) (res uint8) {
	checkUint64Bounds(v, Uint8Max, res)
	return uint8(v)
}

func convertUint64ToUint16(v uint64) (res uint16) {
	checkUint64Bounds(v, Uint16Max, res)
	return uint16(v)
}

func convertUint64ToUint32(v uint64) (res uint32) {
	checkUint64Bounds(v, Uint32Max, res)
	return uint32(v)
}

func convertUint64ToUintptr(v uint64) (res uintptr) {
	checkUint64Bounds(v, uint64(UintptrMax), res)
	return uintptr(v)
}

func convertUint64ToFloat32(v uint64) (res float32) {
	checkUint64Bounds(v, Float32IntMax, res)
	return float32(v)
}

func convertUint64ToFloat64(v uint64) (res float64) {
	checkUint64Bounds(v, Float64IntMax, res)
	return float64(v)
}

func checkInt64Bounds(v int64, min int64, max uint64, res interface{}) {
	if v < min || (v > 0 && uint64(v) > max) {
		panic(&OutOfBoundsError{
			Value: v,
			Type:  reflect.TypeOf(res),
		})
	}
}

func convertInt64ToInt(v int64) (res int) {
	checkInt64Bounds(v, int64(IntMin), uint64(IntMax), res)
	return int(v)
}

func convertInt64ToInt8(v int64) (res int8) {
	checkInt64Bounds(v, Int8Min, Int8Max, res)
	return int8(v)
}

func convertInt64ToInt16(v int64) (res int16) {
	checkInt64Bounds(v, Int16Min, Int16Max, res)
	return int16(v)
}

func convertInt64ToInt32(v int64) (res int32) {
	checkInt64Bounds(v, Int32Min, Int32Max, res)
	return int32(v)
}

func convertInt64ToUint(v int64) (res uint) {
	checkInt64Bounds(v, int64(UintMin), uint64(UintMax), res)
	return uint(v)
}

func convertInt64ToUint8(v int64) (res uint8) {
	checkInt64Bounds(v, Uint8Min, Uint8Max, res)
	return uint8(v)
}

func convertInt64ToUint16(v int64) (res uint16) {
	checkInt64Bounds(v, Uint16Min, Uint16Max, res)
	return uint16(v)
}

func convertInt64ToUint32(v int64) (res uint32) {
	checkInt64Bounds(v, Uint32Min, Uint32Max, res)
	return uint32(v)
}

func convertInt64ToUint64(v int64) (res uint64) {
	checkInt64Bounds(v, Uint64Min, Uint64Max, res)
	return uint64(v)
}

func convertInt64ToUintptr(v int64) (res uintptr) {
	checkInt64Bounds(v, int64(UintptrMin), uint64(UintptrMax), res)
	return uintptr(v)
}

func convertInt64ToFloat32(v int64) (res float32) {
	checkInt64Bounds(v, Float32IntMin, Float32IntMax, res)
	return float32(v)
}

func convertInt64ToFloat64(v int64) (res float64) {
	checkInt64Bounds(v, Float64IntMin, Float64IntMax, res)
	return float64(v)
}

func convertPanicToError(v interface{}) error {
	if v == nil {
		return nil
	}
	switch e := v.(type) {
	case error:
		return e
	default:
		return fmt.Errorf("objconv: %v", v)
	}
}
