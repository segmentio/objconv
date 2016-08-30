objconv [![CircleCI](https://circleci.com/gh/segmentio/objconv.svg?style=shield)](https://circleci.com/gh/segmentio/objconv) [![Go Report Card](https://goreportcard.com/badge/github.com/segmentio/objconv)](https://goreportcard.com/report/github.com/segmentio/objconv) [![GoDoc](https://godoc.org/github.com/segmentio/objconv?status.svg)](https://godoc.org/github.com/segmentio/objconv)
=======

A Go package exposing encoder and decoders that support data streaming to and
from multiple formats.

Installation
------------

```shell
go get github.com/segmentio/objconv
```

Encoder
-------

The package exposes a generic `Encoder` interface that let's the program
serialize native values into various formats.

Here's an example of how to serialize a structure to JSON:
```go
package main

import (
    "os"

    "github.com/segmentio/objconv"
)

func main() {
    objconv.Encode(os.Stdout, "json", struct{
        Hello `objconv:"hello"`
    }{"world"}) // prints {"hello":"world"}
}
```

Decoder
-------

Streaming
---------
