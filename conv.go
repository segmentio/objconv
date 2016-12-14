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

func checkInt64Bounds(v int64, min int64, max uint64, t reflect.Type) error {
	if v < min || (v > 0 && uint64(v) > max) {
		return &OutOfBoundsError{
			Value: v,
			Type:  t,
		}
	}
	return nil
}
