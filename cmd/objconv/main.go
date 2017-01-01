package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/segmentio/objconv"
	_ "github.com/segmentio/objconv/json"
	"github.com/segmentio/objconv/mimetype"
	_ "github.com/segmentio/objconv/msgpack"
	_ "github.com/segmentio/objconv/resp"
)

func main() {
	var r = bufio.NewReader(os.Stdin)
	var w = bufio.NewWriter(os.Stdout)
	var input string
	var output string

	flag.StringVar(&input, "i", "json", "The format of the input stream")
	flag.StringVar(&output, "o", "json", "The format of the output stream")
	flag.Parse()

	if err := conv(w, output, r, input); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	w.Flush()
}

func conv(w io.Writer, output string, r io.Reader, input string) (err error) {
	var ic mimetype.Codec
	var oc mimetype.Codec
	var ok bool

	if ic, ok = mimetype.Lookup(input); !ok {
		err = fmt.Errorf("unknown input format: %s", input)
		return
	}

	if oc, ok = mimetype.Lookup(output); !ok {
		err = fmt.Errorf("unknown output format: %s", output)
		return
	}

	var d = objconv.NewStreamDecoder(ic.NewParser(r))
	var e *objconv.StreamEncoder
	var v interface{}

	if e, err = d.Encoder(oc.NewEmitter(w)); err != nil {
		if err == io.EOF { // empty input
			err = nil
		}
		return
	}

	for d.Decode(&v) == nil {
		if err = e.Encode(v); err != nil {
			return
		}
		v = nil
	}

	if err = e.Close(); err != nil {
		return
	}

	// Not ideal but does the job, if the output is JSON we add a newline
	// character at the end to make it easier to read in terminals.
	if strings.Contains(output, "json") {
		fmt.Fprintln(w)
	}

	err = d.Err()
	return
}
