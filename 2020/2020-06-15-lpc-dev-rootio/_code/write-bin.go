package main

import (
	"encoding/binary"
	"io"
	"log"
	"math"
	"os"
)

// START OMIT
func main() {
	f, err := os.Create("data.raw")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := make([]byte, 8)
	bin := binary.LittleEndian

	bin.PutUint64(buf, math.Float64bits(10))
	write(f, buf)
	bin.PutUint64(buf, math.Float64bits(20))
	write(f, buf)
}

func write(w io.Writer, buf []byte) {
	_, err := w.Write(buf)
	if err != nil {
		log.Fatal(err)
	}
}

// END OMIT
