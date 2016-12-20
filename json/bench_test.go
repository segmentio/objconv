// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Large data benchmark.
// The JSON data is a summary of agl's changes in the
// go, webkit, and chromium open source projects.
// We benchmark converting between the JSON form
// and in-memory data structures.

package json

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

type codeResponse struct {
	Tree     *codeNode `objconv:"tree"`
	Username string    `objconv:"username"`
}

type codeNode struct {
	Name     string      `objconv:"name"`
	Kids     []*codeNode `objconv:"kids"`
	CLWeight float64     `objconv:"cl_weight"`
	Touches  int         `objconv:"touches"`
	MinT     int64       `objconv:"min_t"`
	MaxT     int64       `objconv:"max_t"`
	MeanT    int64       `objconv:"mean_t"`
}

var codeJSON []byte
var codeStruct codeResponse

func init() {
	f, err := os.Open(os.Getenv("GOROOT") + "/src/encoding/json/testdata/code.json.gz")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	codeJSON = data

	if err := json.Unmarshal(codeJSON, &codeStruct); err != nil {
		panic("unmarshal code.json: " + err.Error())
	}
}

func BenchmarkCodeEncoder(b *testing.B) {
	enc := NewEncoder(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		if err := enc.Encode(&codeStruct); err != nil {
			b.Fatal("Encode:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Marshal(&codeStruct); err != nil {
			b.Fatal("Marshal:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

/*
func BenchmarkCodeDecoder(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	var buf bytes.Buffer
	dec := NewDecoder(&buf)
	var r codeResponse
	for i := 0; i < b.N; i++ {
		buf.Write(codeJSON)
		// hide EOF
		buf.WriteByte('\n')
		buf.WriteByte('\n')
		buf.WriteByte('\n')
		if err := dec.Decode(&r); err != nil {
			b.Fatal("Decode:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkDecoderStream(b *testing.B) {
	b.StopTimer()
	var buf bytes.Buffer
	dec := NewDecoder(&buf)
	buf.WriteString(`"` + strings.Repeat("x", 1000000) + `"` + "\n\n\n")
	var x interface{}
	if err := dec.Decode(&x); err != nil {
		b.Fatal("Decode:", err)
	}
	ones := strings.Repeat(" 1\n", 300000) + "\n\n\n"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if i%300000 == 0 {
			buf.WriteString(ones)
		}
		x = nil
		if err := dec.Decode(&x); err != nil || x != 1.0 {
			b.Fatalf("Decode: %v after %d", err, i)
		}
	}
}

func BenchmarkCodeUnmarshal(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	for i := 0; i < b.N; i++ {
		var r codeResponse
		if err := Unmarshal(codeJSON, &r); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
	b.SetBytes(int64(len(codeJSON)))
}

func BenchmarkCodeUnmarshalReuse(b *testing.B) {
	if codeJSON == nil {
		b.StopTimer()
		codeInit()
		b.StartTimer()
	}
	var r codeResponse
	for i := 0; i < b.N; i++ {
		if err := Unmarshal(codeJSON, &r); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalString(b *testing.B) {
	data := []byte(`"hello, world"`)
	var s string

	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, &s); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalFloat64(b *testing.B) {
	var f float64
	data := []byte(`3.14`)

	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, &f); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkUnmarshalInt64(b *testing.B) {
	var x int64
	data := []byte(`3`)

	for i := 0; i < b.N; i++ {
		if err := Unmarshal(data, &x); err != nil {
			b.Fatal("Unmarshal:", err)
		}
	}
}

func BenchmarkIssue10335(b *testing.B) {
	b.ReportAllocs()
	var s struct{}
	j := []byte(`{"a":{ }}`)
	for n := 0; n < b.N; n++ {
		if err := Unmarshal(j, &s); err != nil {
			b.Fatal(err)
		}
	}
}
*/
