package cbor

import "encoding/binary"

const (
	MajorType0 byte = iota
	MajorType1
	MajorType2
	MajorType3
	MajorType4
	MajorType5
	MajorType6
	MajorType7
)

const ( // usigned integer types
	Uint8 = 24 + iota
	Uint16
	Uint32
	Uint64
)

const ( // simple values
	False = 20 + iota
	True
	Null
	Undefined
	Extension
	Float16
	Float32
	Float64
	Break = 31
)

func majorByte(maj byte, val byte) byte {
	return (maj << 5) | val
}

func majorTypeOf(b byte) (maj byte, val byte) {
	const mask = byte(0xE0)
	return ((b & mask) >> 5), (b & ^mask)
}

func putUint16(b []byte, v uint16) {
	binary.BigEndian.PutUint16(b, v)
}

func putUint32(b []byte, v uint32) {
	binary.BigEndian.PutUint32(b, v)
}

func putUint64(b []byte, v uint64) {
	binary.BigEndian.PutUint64(b, v)
}

func getUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func getUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func getUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func align(n int, a int) int {
	if (n % a) == 0 {
		return n
	}
	return ((n / a) + 1) * a
}

// ported from http://stderr.org/doc/ogre-doc/api/OgreBitwise_8h-source.html
func f16tof32bits(yy uint16) (d uint32) {
	y := uint32(yy)
	s := (y >> 15) & 0x01
	e := (y >> 10) & 0x1f
	m := y & 0x03ff

	if e == 0 {
		if m == 0 { // plus or minus 0
			return s << 31
		} else { // Denormalized number -- renormalize it
			for (m & 0x00000400) == 0 {
				m <<= 1
				e -= 1
			}
			e += 1
			const zz uint32 = 0x0400
			m &= ^zz
		}
	} else if e == 31 {
		if m == 0 { // Inf
			return (s << 31) | 0x7f800000
		} else { // NaN
			return (s << 31) | 0x7f800000 | (m << 13)
		}
	}

	e = e + (127 - 15)
	m = m << 13
	return (s << 31) | (e << 23) | m
}
