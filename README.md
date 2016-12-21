objconv [![CircleCI](https://circleci.com/gh/segmentio/objconv.svg?style=shield)](https://circleci.com/gh/segmentio/objconv) [![Go Report Card](https://goreportcard.com/badge/github.com/segmentio/objconv)](https://goreportcard.com/report/github.com/segmentio/objconv) [![GoDoc](https://godoc.org/github.com/segmentio/objconv?status.svg)](https://godoc.org/github.com/segmentio/objconv)
=======

This Go package provides the implementation of high performance encoder and
decoders for JSON-like object representations.

The top-level package exposes the generic types and algorithms for encoding and
decoding values, while each sub-package implements the parser and emitters for
specific types.

Compatibility with the standard library
---------------------------------------

The sub-packages providing implementation for specific formats also expose APIs
that mirror those of the standard library to make it easy to integrate with the
objconv package. However there are a couple of differences that need to be taken
in consideration:

- When struct tags are used to define the behavior of marshaling or unmarshaling
structs, the tag name needs to be `objconv` (instead of `json` for example).

- Encoder and Decoder types are not exposed in the objconv sub-packages, instead
the types from the top-level package are used. For example, variables declared
with the `json.Encoder` type would have to be replaced with `objconv.Encoder`.

- Interfaces like `json.Marshaler` or `json.Unmarshaler` are not supported.
However the `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces
are.

Encoder
-------

The package exposes a generic encoder API that let's the program serialize
native values into various formats.

Here's an example of how to serialize a structure to JSON:
```go
package main

import (
    "os"

    "github.com/segmentio/objconv/json"
)

func main() {
    enc := json.NewEncoder(os.Stdout)
    enc.Encode(struct{
        Hello string
    }{"World"})
}
```
```
$ go run ./example.go
{"Hello":"World"}
```

Note that this code is fully compatible with the standard `encoding/json`
package.

Decoder
-------

Here's an example of how to use a JSON decoder:
```go
package main

import (
    "fmt"
    "os"

    "github.com/segmentio/objconv/json"
)

func main() {
    v := struct{
        Message string
    }{}

    dec := json.NewDecoder(os.Stdin)
    dec.Decode(&v)

    fmt.Println(v.Message)
}
```
```
$ echo '{ "Message": "Hello World!" }' | go run ./example.go
Hello World!
```

Streaming
---------

One of the interesting features of the `objconv` package is the ability to read
and write streams of data. This has several advantages in terms of memory usage
and latency when passing data from service to service.  
The package exposes the `StreamEncoder` and `StreamDecoder` types for this
purpose.

For example the JSON stream encoder and decoder can produce a JSON array as a
stream where data are produced and consumed on the fly as they become available,
here's an example:
```go
package main

import (
    "io"

    "github.com/segmentio/objconv/json"
)

func main() {
     r, w := io.Pipe()

    go func() {
        defer w.Close()

        enc := json.NewStreamEncoder(w)

        // Produce values to the JSON stream.
        for i := 0; i != 1000; i++ {
            enc.Encode(i)
        }

        enc.Close()
    }()

    dec := json.NewStreamDecoder(r)

    // Consume values from the JSON stream.
    var v interface{}

    for dec.Decode(&v) == nil {
        // v => {0..999}
        // ...
        v = nil
    }
}
```

Stream decoders are capable of reading values from either arrays or single
values, this is very convenient when an program cannot predict the structure of
the stream. If the actual data representation is not an array the stream decoder
will simply behave like a normal decoder and produce a single value.

Encoding and decoding custom types
----------------------------------

To override the default encoder and decoder behaviors a type may implement the
`ValueEncoder` or `ValueDecoder` interface. The method on these interfaces are
called to customize the default behavior.

This can prove very useful to represent slice of pairs as maps for example:
```go
type KV struct {
    K string
    V interface{}
}

type M []KV

// Imlpement the ValueEncoder interface to provide a custom encoding.
func (m M) ValueEncode(e objconv.Encoder) error {
    i := 0
    return e.EncodeMap(func(k objconv.Encoder, v objconv.Encoder) (err error) {
        if i == len(m) {
            return objconv.End // done
        }

        if err = k.Encode(m[i].K); err != nil {
            return
        }

        if err = v.Encode(m[i].V); err != nil {
            return
        }

        i++
        return
    })
}
```

Mime Types
----------

The `mimetype` sub-package exposes APIs for creating encoders and decoders for
specific mime types. When an objconv package for a specific format is imported
it registers itself on the `mimetype` registry to be later referred by name.

```go
import (
    "github.com/segmentio/objconv/mimetype"
    _ "github.com/segmentio/objconv/json" // registers the JSON codec
)

func main() {
    // Creates an encoder for the "application/json" mime type.
    enc := mimetype.NewEncoder("application/json", os.Stdout)

    if enc == nil {
        // no encoder for the specified mime type exists
    }

    // ...
}
```
