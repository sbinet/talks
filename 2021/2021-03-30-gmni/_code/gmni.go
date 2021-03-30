package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
)

func main() {
	certs := &certificate.Store{}
	certs.Register("localhost")
	if err := certs.Load("./certs"); err != nil {
		log.Fatal(err)
	}

	mux := &gemini.Mux{}
	mux.HandleFunc("/", func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
		fmt.Fprint(w, fmt.Sprintf(`# Heading 1
You have requested (at %s):
=> %s
`, time.Now().UTC().Format(time.RFC3339Nano), r.URL.Path))
	})

	server := &gemini.Server{
		Handler:        gemini.LoggingMiddleware(mux),
		GetCertificate: certs.Get,
	}

	err := server.ListenAndServe(context.Background())
	if err != nil {
		log.Fatalf("could not run gmni server: %+v", err)
	}
}
