package objconv

import "reflect"

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

var (
	intType   = reflect.TypeOf(int(0))
	int8Type  = reflect.TypeOf(int8(0))
	int16Type = reflect.TypeOf(int16(0))
	int32Type = reflect.TypeOf(int32(0))
	int64Type = reflect.TypeOf(int64(0))

	uintType    = reflect.TypeOf(uint(0))
	uint8Type   = reflect.TypeOf(uint8(0))
	uint16Type  = reflect.TypeOf(uint16(0))
	uint32Type  = reflect.TypeOf(uint32(0))
	uint64Type  = reflect.TypeOf(uint64(0))
	uintptrType = reflect.TypeOf(uintptr(0))

	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))
)

func checkUint64Bounds(v uint64, max uint64, t reflect.Type) error {
	if v > max {
		return &OutOfBoundsError{
			Value: v,
			Type:  t,
		}
	}
	return nil
}

func convertUint64ToInt(v uint64) (int, error) {
	return int(v), checkUint64Bounds(v, uint64(IntMax), intType)
}

func convertUint64ToInt8(v uint64) (int8, error) {
	return int8(v), checkUint64Bounds(v, Int8Max, int8Type)
}

func convertUint64ToInt16(v uint64) (int16, error) {
	return int16(v), checkUint64Bounds(v, Int16Max, int16Type)
}

func convertUint64ToInt32(v uint64) (int32, error) {
	return int32(v), checkUint64Bounds(v, Int32Max, int32Type)
}

func convertUint64ToInt64(v uint64) (int64, error) {
	return int64(v), checkUint64Bounds(v, Uint64Max, int64Type)
}

func convertUint64ToUint(v uint64) (uint, error) {
	return uint(v), checkUint64Bounds(v, uint64(UintMax), uintType)
}

func convertUint64ToUint8(v uint64) (uint8, error) {
	return uint8(v), checkUint64Bounds(v, Uint8Max, uint8Type)
}

func convertUint64ToUint16(v uint64) (uint16, error) {
	return uint16(v), checkUint64Bounds(v, Uint16Max, uint16Type)
}

func convertUint64ToUint32(v uint64) (uint32, error) {
	return uint32(v), checkUint64Bounds(v, Uint32Max, uint32Type)
}

func convertUint64ToUintptr(v uint64) (uintptr, error) {
	return uintptr(v), checkUint64Bounds(v, uint64(UintptrMax), uintptrType)
}

func convertUint64ToFloat32(v uint64) (float32, error) {
	return float32(v), checkUint64Bounds(v, Float32IntMax, float32Type)
}

func convertUint64ToFloat64(v uint64) (float64, error) {
	return float64(v), checkUint64Bounds(v, Float64IntMax, float64Type)
}

func checkInt64Bounds(v int64, min int64, max uint64, t reflect.Type) error {
	if v < min || (v > 0 && uint64(v) > max) {
		return &OutOfBoundsError{
			Value: v,
			Type:  t,
		}
	}
	return nil
}

func convertInt64ToInt(v int64) (int, error) {
	return int(v), checkInt64Bounds(v, int64(IntMin), uint64(IntMax), intType)
}

func convertInt64ToInt8(v int64) (int8, error) {
	return int8(v), checkInt64Bounds(v, Int8Min, Int8Max, int8Type)
}

func convertInt64ToInt16(v int64) (int16, error) {
	return int16(v), checkInt64Bounds(v, Int16Min, Int16Max, int16Type)
}

func convertInt64ToInt32(v int64) (int32, error) {
	return int32(v), checkInt64Bounds(v, Int32Min, Int32Max, int32Type)
}

func convertInt64ToUint(v int64) (uint, error) {
	return uint(v), checkInt64Bounds(v, int64(UintMin), uint64(UintMax), uintType)
}

func convertInt64ToUint8(v int64) (uint8, error) {
	return uint8(v), checkInt64Bounds(v, Uint8Min, Uint8Max, uint8Type)
}

func convertInt64ToUint16(v int64) (uint16, error) {
	return uint16(v), checkInt64Bounds(v, Uint16Min, Uint16Max, uint16Type)
}

func convertInt64ToUint32(v int64) (uint32, error) {
	return uint32(v), checkInt64Bounds(v, Uint32Min, Uint32Max, uint32Type)
}

func convertInt64ToUint64(v int64) (uint64, error) {
	return uint64(v), checkInt64Bounds(v, Uint64Min, Uint64Max, uint64Type)
}

func convertInt64ToUintptr(v int64) (uintptr, error) {
	return uintptr(v), checkInt64Bounds(v, int64(UintptrMin), uint64(UintptrMax), uintptrType)
}

func convertInt64ToFloat32(v int64) (float32, error) {
	return float32(v), checkInt64Bounds(v, Float32IntMin, Float32IntMax, float32Type)
}

func convertInt64ToFloat64(v int64) (float64, error) {
	return float64(v), checkInt64Bounds(v, Float64IntMin, Float64IntMax, float64Type)
}
