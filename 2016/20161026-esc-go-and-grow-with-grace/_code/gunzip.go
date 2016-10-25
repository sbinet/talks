package main

import (
	"compress/gzip"
	"io"
	"os"
)

func main() {
	r, err := gzip.NewReader(os.Stdin)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		panic(err)
	}
}
