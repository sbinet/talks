package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("error while running web server: %+v", err)
	}
}
