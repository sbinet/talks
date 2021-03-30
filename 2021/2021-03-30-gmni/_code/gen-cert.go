package main

import (
	"log"
	"time"

	"git.sr.ht/~adnano/go-gemini/certificate"
)

func main() {
	opts := certificate.CreateOptions{
		DNSNames: []string{"localhost"}, Duration: 10 * 24 * time.Hour,
	}
	cert, err := certificate.Create(opts)
	if err != nil {
		log.Fatalf("could not create certificate: %+v", err)
	}

	err = certificate.Write(cert, "./certs/cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("could not save certificate: %+v", err)
	}
}
