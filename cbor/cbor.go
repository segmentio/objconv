package cbor

import "encoding/binary"

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
