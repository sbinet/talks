package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Create("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "px=%v, py=%v, pz=%v, e=%v\n", 10, 20, 30, 40)
	fmt.Fprintf(f, "px=%v, py=%v, pz=%v, e=%v\n", 12, -23, -42, 600)
}
