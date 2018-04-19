package main

import (
	"flag"

	"github.com/segmentio/objconv/generate"
)

func main() {
	var pkg string
	var typ string
	var output string

	flag.StringVar(&pkg, "p", "", "The package containing the type to generate for")
	flag.StringVar(&typ, "t", "", "The name of the type to generate for")
	flag.StringVar(&output, "o", "", "The name of the file to generate into")

	flag.Parse()

	generate.GenerateDecode(pkg, typ)
}
