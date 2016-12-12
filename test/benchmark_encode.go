package test

import (
	"math"
	"testing"

	"github.com/segmentio/objconv"
)

func BenchmarkEncodeNil(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, nil)
}

func BenchmarkEncodeBoolTrue(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, true)
}

func BenchmarkEncodeBoolFalse(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, false)
}

func BenchmarkEncodeIntZero(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, int64(0))
}

func BenchmarkEncodeIntShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, int64(100))
}

func BenchmarkEncodeIntLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, int64(math.MaxInt64))
}

func BenchmarkEncodeUintZero(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, uint64(0))
}

func BenchmarkEncodeUintShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, uint64(100))
}

func BenchmarkEncodeUintLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, uint64(math.MaxUint64))
}

func BenchmarkEncodeFloatZero(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, float64(0))
}

func BenchmarkEncodeFloatShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, float64(1.234))
}

func BenchmarkEncodeFloatLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, float64(math.MaxFloat64))
}

func BenchmarkEncodeStringEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, "")
}

func BenchmarkEncodeStringShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, "Hello World!")
}

func BenchmarkEncodeStringLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, longString)
}

func BenchmarkEncodeBytesEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, bytesEmpty)
}

func BenchmarkEncodeBytesShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, bytesShort)
}

func BenchmarkEncodeBytesLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, bytesLong)
}

func BenchmarkEncodeSliceInterfaceEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceInterfaceEmpty)
}

func BenchmarkEncodeSliceInterfaceShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceInterfaceShort)
}

func BenchmarkEncodeSliceInterfaceLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceInterfaceLong)
}

func BenchmarkEncodeSliceStringEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStringEmpty)
}

func BenchmarkEncodeSliceStringShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStringShort)
}

func BenchmarkEncodeSliceStringLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStringLong)
}

func BenchmarkEncodeSliceBytesEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceBytesEmpty)
}

func BenchmarkEncodeSliceBytesShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceBytesShort)
}

func BenchmarkEncodeSliceBytesLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceBytesLong)
}

func BenchmarkEncodeSliceStructEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStructEmpty)
}

func BenchmarkEncodeSliceStructShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStructShort)
}

func BenchmarkEncodeSliceStructLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, sliceStructLong)
}

func BenchmarkEncodeMapStringStringEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapStringStringEmpty)
}

func BenchmarkEncodeMapStringStringShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapStringStringShort)
}

func BenchmarkEncodeMapStringStringLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapStringStringLong)
}

func BenchmarkEncodeMapStringInterfaceEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapStringInterfaceEmpty)
}

func BenchmarkEncodeMapStringInterfaceShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapStringInterfaceShort)
}

func BenchmarkEncodeMapStringInterfaceLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapStringInterfaceLong)
}

func BenchmarkEncodeMapStringStructEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapStringStructEmpty)
}

func BenchmarkEncodeMapStringStructShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapStringStructShort)
}

func BenchmarkEncodeMapStringStructLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapStringStructLong)
}

func BenchmarkEncodeMapSliceEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapSliceEmpty)
}

func BenchmarkEncodeMapSliceShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapSliceShort)
}

func BenchmarkEncodeMapSliceLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, mapSliceLong)
}

func BenchmarkEncodeStructEmpty(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, structEmpty)
}

func BenchmarkEncodeStructShort(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, structShort)
}

func BenchmarkEncodeStructLong(b *testing.B, e objconv.Encoder) {
	BenchmarkEncode(b, e, structLong)
}

func BenchmarkEncode(b *testing.B, e objconv.Encoder, v interface{}) {
	b.ResetTimer()

	for i := 0; i != b.N; i++ {
		e.Encode(v)
	}
}

const longString = `Package json implements encoding and decoding of JSON objects as defined in RFC 4627. The mapping between JSON objects and Go values is described in the documentation for the Marshal and Unmarshal functions.`

var (
	bytesEmpty = []byte{}
	bytesShort = []byte("Hello World!")
	bytesLong  = []byte(`Package json implements encoding and decoding of JSON objects as defined in RFC 4627. The mapping between JSON objects and Go values is described in the documentation for the Marshal and Unmarshal functions.`)

	sliceInterfaceEmpty = []interface{}{}
	sliceInterfaceShort = []interface{}{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil} // 10 items
	sliceInterfaceLong  = []interface{}{
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	} // 100 items

	sliceStringEmpty = []string{}
	sliceStringShort = []string{"", "", "", "", "", "", "", "", "", ""} // 10 items
	sliceStringLong  = []string{
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "",
	} // 100 items

	sliceBytesEmpty = [][]byte{}
	sliceBytesShort = [][]byte{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil} // 10 items
	sliceBytesLong  = [][]byte{
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	} // 100 items

	sliceStructEmpty = []struct{}{}
	sliceStructShort = []struct{}{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}} // 10 items
	sliceStructLong  = []struct{}{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
	} // 100 items

	structEmpty = struct{}{}
	structShort = struct {
		A int
		B int
		C int
	}{}
	structLong = struct {
		A int
		B int
		C int
		D int
		E int
		F int
		G int
		H int
		I int
		J int
		K int
		L int
		M int
		N int
		O int
		P int
		Q int
		R int
		S int
		T int
		U int
		V int
		W int
		X int
		Y int
		Z int
	}{}

	mapStringStringEmpty = map[string]string{}
	mapStringStringShort = map[string]string{
		"0": "", "1": "", "2": "", "3": "", "4": "", "5": "", "6": "", "7": "", "8": "", "9": "",
	} // 10 items
	mapStringStringLong = map[string]string{
		"0": "", "1": "", "2": "", "3": "", "4": "", "5": "", "6": "", "7": "", "8": "", "9": "",
		"10": "", "11": "", "12": "", "13": "", "14": "", "15": "", "16": "", "17": "", "18": "", "19": "",
		"20": "", "21": "", "22": "", "23": "", "24": "", "25": "", "26": "", "27": "", "28": "", "29": "",
		"30": "", "31": "", "32": "", "33": "", "34": "", "35": "", "36": "", "37": "", "38": "", "39": "",
		"40": "", "41": "", "42": "", "43": "", "44": "", "45": "", "46": "", "47": "", "48": "", "49": "",
		"50": "", "51": "", "52": "", "53": "", "54": "", "55": "", "56": "", "57": "", "58": "", "59": "",
		"60": "", "61": "", "62": "", "63": "", "64": "", "65": "", "66": "", "67": "", "68": "", "69": "",
		"70": "", "71": "", "72": "", "73": "", "74": "", "75": "", "76": "", "77": "", "78": "", "79": "",
		"80": "", "81": "", "82": "", "83": "", "84": "", "85": "", "86": "", "87": "", "88": "", "89": "",
		"90": "", "91": "", "92": "", "93": "", "94": "", "95": "", "96": "", "97": "", "98": "", "99": "",
	} // 100 items

	mapStringInterfaceEmpty = map[string]interface{}{}
	mapStringInterfaceShort = map[string]interface{}{
		"0": nil, "1": nil, "2": nil, "3": nil, "4": nil, "5": nil, "6": nil, "7": nil, "8": nil, "9": nil,
	} // 10 items
	mapStringInterfaceLong = map[string]interface{}{
		"0": nil, "1": nil, "2": nil, "3": nil, "4": nil, "5": nil, "6": nil, "7": nil, "8": nil, "9": nil,
		"10": nil, "11": nil, "12": nil, "13": nil, "14": nil, "15": nil, "16": nil, "17": nil, "18": nil, "19": nil,
		"20": nil, "21": nil, "22": nil, "23": nil, "24": nil, "25": nil, "26": nil, "27": nil, "28": nil, "29": nil,
		"30": nil, "31": nil, "32": nil, "33": nil, "34": nil, "35": nil, "36": nil, "37": nil, "38": nil, "39": nil,
		"40": nil, "41": nil, "42": nil, "43": nil, "44": nil, "45": nil, "46": nil, "47": nil, "48": nil, "49": nil,
		"50": nil, "51": nil, "52": nil, "53": nil, "54": nil, "55": nil, "56": nil, "57": nil, "58": nil, "59": nil,
		"60": nil, "61": nil, "62": nil, "63": nil, "64": nil, "65": nil, "66": nil, "67": nil, "68": nil, "69": nil,
		"70": nil, "71": nil, "72": nil, "73": nil, "74": nil, "75": nil, "76": nil, "77": nil, "78": nil, "79": nil,
		"80": nil, "81": nil, "82": nil, "83": nil, "84": nil, "85": nil, "86": nil, "87": nil, "88": nil, "89": nil,
		"90": nil, "91": nil, "92": nil, "93": nil, "94": nil, "95": nil, "96": nil, "97": nil, "98": nil, "99": nil,
	} // 100 items

	mapStringStructEmpty = map[string]struct{}{}
	mapStringStructShort = map[string]struct{}{
		"0": {}, "1": {}, "2": {}, "3": {}, "4": {}, "5": {}, "6": {}, "7": {}, "8": {}, "9": {},
	}
	mapStringStructLong = map[string]struct{}{
		"0": {}, "1": {}, "2": {}, "3": {}, "4": {}, "5": {}, "6": {}, "7": {}, "8": {}, "9": {},
		"10": {}, "11": {}, "12": {}, "13": {}, "14": {}, "15": {}, "16": {}, "17": {}, "18": {}, "19": {},
		"20": {}, "21": {}, "22": {}, "23": {}, "24": {}, "25": {}, "26": {}, "27": {}, "28": {}, "29": {},
		"30": {}, "31": {}, "32": {}, "33": {}, "34": {}, "35": {}, "36": {}, "37": {}, "38": {}, "39": {},
		"40": {}, "41": {}, "42": {}, "43": {}, "44": {}, "45": {}, "46": {}, "47": {}, "48": {}, "49": {},
		"50": {}, "51": {}, "52": {}, "53": {}, "54": {}, "55": {}, "56": {}, "57": {}, "58": {}, "59": {},
		"60": {}, "61": {}, "62": {}, "63": {}, "64": {}, "65": {}, "66": {}, "67": {}, "68": {}, "69": {},
		"70": {}, "71": {}, "72": {}, "73": {}, "74": {}, "75": {}, "76": {}, "77": {}, "78": {}, "79": {},
		"80": {}, "81": {}, "82": {}, "83": {}, "84": {}, "85": {}, "86": {}, "87": {}, "88": {}, "89": {},
		"90": {}, "91": {}, "92": {}, "93": {}, "94": {}, "95": {}, "96": {}, "97": {}, "98": {}, "99": {},
	}

	mapSliceEmpty = objconv.MapSlice{}
	mapSliceShort = objconv.MapSlice{
		{"0", nil}, {"1", nil}, {"2", nil}, {"3", nil}, {"4", nil}, {"5", nil}, {"6", nil}, {"7", nil}, {"8", nil}, {"9", nil},
	}
	mapSliceLong = objconv.MapSlice{
		{"0", nil}, {"1", nil}, {"2", nil}, {"3", nil}, {"4", nil}, {"5", nil}, {"6", nil}, {"7", nil}, {"8", nil}, {"9", nil},
		{"10", nil}, {"11", nil}, {"12", nil}, {"13", nil}, {"14", nil}, {"15", nil}, {"16", nil}, {"17", nil}, {"18", nil}, {"19", nil},
		{"20", nil}, {"21", nil}, {"22", nil}, {"23", nil}, {"24", nil}, {"25", nil}, {"26", nil}, {"27", nil}, {"28", nil}, {"29", nil},
		{"30", nil}, {"31", nil}, {"32", nil}, {"33", nil}, {"34", nil}, {"35", nil}, {"36", nil}, {"37", nil}, {"38", nil}, {"39", nil},
		{"40", nil}, {"41", nil}, {"42", nil}, {"43", nil}, {"44", nil}, {"45", nil}, {"46", nil}, {"47", nil}, {"48", nil}, {"49", nil},
		{"50", nil}, {"51", nil}, {"52", nil}, {"53", nil}, {"54", nil}, {"55", nil}, {"56", nil}, {"57", nil}, {"58", nil}, {"59", nil},
		{"60", nil}, {"61", nil}, {"62", nil}, {"63", nil}, {"64", nil}, {"65", nil}, {"66", nil}, {"67", nil}, {"68", nil}, {"69", nil},
		{"70", nil}, {"71", nil}, {"72", nil}, {"73", nil}, {"74", nil}, {"75", nil}, {"76", nil}, {"77", nil}, {"78", nil}, {"79", nil},
		{"80", nil}, {"81", nil}, {"82", nil}, {"83", nil}, {"84", nil}, {"85", nil}, {"86", nil}, {"87", nil}, {"88", nil}, {"89", nil},
		{"90", nil}, {"91", nil}, {"92", nil}, {"93", nil}, {"94", nil}, {"95", nil}, {"96", nil}, {"97", nil}, {"98", nil}, {"99", nil},
	}
)
